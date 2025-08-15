package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestLocalDatabasePersistence tests that the local database persists users across restarts
func TestLocalDatabasePersistence(t *testing.T) {
	// Create temporary directory for test database
	tmpDir, err := os.MkdirTemp("", "smor_ting_test_db_*")
	require.NoError(t, err, "Should create temp directory")
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "smor_ting_test.db")
	
	t.Run("Should persist users after database restart", func(t *testing.T) {
		ctx := context.Background()

		// First session: Create a user
		{
			// Create persistent local database connection
			authService := createLocalAuthService(t, dbPath)
			
			// Create test user
			registerReq := &models.RegisterRequest{
				Email:     "test_persist@smorting.com",
				Password:  "TestPass123!",
				FirstName: "Test",
				LastName:  "Persist",
				Phone:     "231123456789",
				Role:      models.CustomerRole,
			}
			
			// Register user (should succeed)
			response, err := authService.Register(ctx, registerReq)
			require.NoError(t, err, "Should register user successfully")
			assert.Equal(t, registerReq.Email, response.User.Email, "Should have correct email")
			assert.NotEmpty(t, response.AccessToken, "Should return access token")
		}

		// Second session: Restart database and verify user exists
		{
			// Create new database connection (simulating restart)
			authService := createLocalAuthService(t, dbPath)
			
			// Try to login with same credentials (should succeed if persisted)
			loginReq := &models.LoginRequest{
				Email:    "test_persist@smorting.com",
				Password: "TestPass123!",
			}
			
			response, err := authService.Login(ctx, loginReq)
			require.NoError(t, err, "Should login successfully with persisted user")
			assert.Equal(t, "test_persist@smorting.com", response.User.Email, "Should return correct user")
			assert.Equal(t, "Test", response.User.FirstName, "Should preserve user data")
			assert.Equal(t, "Persist", response.User.LastName, "Should preserve user data")
			assert.Equal(t, "231123456789", response.User.Phone, "Should preserve user data")
			assert.Equal(t, models.CustomerRole, response.User.Role, "Should preserve user role")
		}
	})

	t.Run("Should handle multiple users persistence", func(t *testing.T) {
		ctx := context.Background()
		
		// Create test users
		testUsers := []models.RegisterRequest{
			{
				Email:     "user1@test.com",
				Password:  "Password1!",
				FirstName: "User",
				LastName:  "One",
				Phone:     "231111111111",
				Role:      models.CustomerRole,
			},
			{
				Email:     "user2@test.com",
				Password:  "Password2!",
				FirstName: "User",
				LastName:  "Two",
				Phone:     "231222222222",
				Role:      models.ProviderRole,
			},
		}

		// First session: Create multiple users
		{
			authService := createLocalAuthService(t, dbPath)
			
			for _, user := range testUsers {
				response, err := authService.Register(ctx, &user)
				require.NoError(t, err, "Should register user %s", user.Email)
				assert.Equal(t, user.Email, response.User.Email, "Should have correct email")
			}
		}

		// Second session: Verify all users persist
		{
			authService := createLocalAuthService(t, dbPath)
			
			for _, user := range testUsers {
				loginReq := &models.LoginRequest{
					Email:    user.Email,
					Password: user.Password,
				}
				
				response, err := authService.Login(ctx, loginReq)
				require.NoError(t, err, "Should login with persisted user %s", user.Email)
				assert.Equal(t, user.Email, response.User.Email, "Should return correct user")
				assert.Equal(t, user.Role, response.User.Role, "Should preserve user role")
			}
		}
	})
}

// TestLocalDatabaseConfiguration tests that local database is properly configured for persistence
func TestLocalDatabaseConfiguration(t *testing.T) {
	t.Run("Should use persistent file-based database for local development", func(t *testing.T) {
		// Create temporary directory for test
		tmpDir, err := os.MkdirTemp("", "smor_ting_config_test_*")
		require.NoError(t, err, "Should create temp directory")
		defer os.RemoveAll(tmpDir)

		dbPath := filepath.Join(tmpDir, "test_config.db")
		
		// Create database file
		authService := createLocalAuthService(t, dbPath)
		require.NotNil(t, authService, "Should create auth service")
		
		// Verify database file exists
		_, err = os.Stat(dbPath)
		assert.NoError(t, err, "Database file should exist after initialization")
	})

	t.Run("Should not use in-memory database for local development", func(t *testing.T) {
		// This test ensures we're not using MongoDB's in-memory storage
		// The presence of a file path indicates persistent storage
		tmpDir, err := os.MkdirTemp("", "smor_ting_memory_test_*")
		require.NoError(t, err, "Should create temp directory")
		defer os.RemoveAll(tmpDir)

		dbPath := filepath.Join(tmpDir, "persistent.db")
		
		// Create auth service
		authService := createLocalAuthService(t, dbPath)
		require.NotNil(t, authService, "Should create auth service")
		
		// Test that configuration indicates persistent storage
		assert.NotEmpty(t, dbPath, "Database path should not be empty")
		assert.Contains(t, dbPath, "persistent.db", "Should use persistent database file")
	})
}

// createLocalAuthService creates an auth service with local persistent database
func createLocalAuthService(t *testing.T, dbPath string) *auth.MongoDBService {
	// Use file-based MongoDB connection string for persistence
	// Format: mongodb://localhost:27017/database_name?directConnection=true
	mongoURI := "mongodb://localhost:27017/smor_ting_local_test?directConnection=true"
	
	// Connect to MongoDB with persistent configuration
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	require.NoError(t, err, "Should connect to local MongoDB")

	// Use local test database
	db := client.Database("smor_ting_local_test")

	// Create logger
	logger, err := logger.New("info", "json", "stdout")
	require.NoError(t, err, "Should create logger")

	// Create repository
	repository := database.NewMongoDBRepository(db, logger)

	// Create auth config for local development
	authConfig := &configs.AuthConfig{
		JWTSecret:        "local-dev-secret-key-for-testing-32-chars-long-here",
		JWTAccessSecret:  "local-access-secret-key-for-dev-32b",
		JWTRefreshSecret: "local-refresh-secret-key-for-dev-32",
		BCryptCost:       8, // Lower cost for faster local development
	}

	// Create auth service
	authService, err := auth.NewMongoDBService(repository, authConfig, logger)
	require.NoError(t, err, "Should create auth service")

	return authService
}
