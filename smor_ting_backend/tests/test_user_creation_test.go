package tests

import (
	"context"
	"os"
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

// setupAuthService creates a real auth service for production database
func setupAuthService(t *testing.T) *auth.MongoDBService {
	// Use production MongoDB connection string
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		t.Skip("MONGODB_URI not set, skipping production database tests")
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	require.NoError(t, err, "Should connect to MongoDB")

	// Use production database (same name as shown in Atlas)
	db := client.Database("smor_ting")

	// Create logger
	logger, err := logger.New("info", "json", "stdout")
	require.NoError(t, err, "Should create logger")

	// Create repository
	repository := database.NewMongoDBRepository(db, logger)

	// Create auth config
	authConfig := &configs.AuthConfig{
		JWTSecret:        "test-secret-key-for-qa-user-creation-testing-32-chars-long",
		JWTAccessSecret:  "test-access-secret-key-for-qa-32b",
		JWTRefreshSecret: "test-refresh-secret-key-for-qa-32",
		BCryptCost:       10, // Lower cost for faster testing
	}

	// Create auth service
	authService, err := auth.NewMongoDBService(repository, authConfig, logger)
	require.NoError(t, err, "Should create auth service")

	return authService
}

// TestCreateQATestUsers tests the creation of QA test users in the production database
func TestCreateQATestUsers(t *testing.T) {
	// Setup test environment
	ctx := context.Background()
	authService := setupAuthService(t)

	testUsers := []models.RegisterRequest{
		{
			Email:     "qa_customer@smorting.com",
			Password:  "TestPass123!",
			FirstName: "QA",
			LastName:  "Customer",
			Phone:     "231777123456",
			Role:      models.CustomerRole,
		},
		{
			Email:     "qa_provider@smorting.com",
			Password:  "ProviderPass123!",
			FirstName: "QA",
			LastName:  "Provider",
			Phone:     "231888123456",
			Role:      models.ProviderRole,
		},
		{
			Email:     "qa_admin@smorting.com",
			Password:  "AdminPass123!",
			FirstName: "QA",
			LastName:  "Admin",
			Phone:     "231999123456",
			Role:      models.AdminRole,
		},
	}

	for _, testUser := range testUsers {
		t.Run("Create "+string(testUser.Role)+" test user", func(t *testing.T) {
			// Test: User should be created successfully
			response, err := authService.Register(ctx, &testUser)

			// Assertions following TDD principles
			require.NoError(t, err, "Should create test user without error")
			require.NotNil(t, response, "Should return auth response")
			require.NotEmpty(t, response.AccessToken, "Should return access token")

			assert.Equal(t, testUser.Email, response.User.Email, "Should have correct email")
			assert.Equal(t, testUser.FirstName, response.User.FirstName, "Should have correct first name")
			assert.Equal(t, testUser.LastName, response.User.LastName, "Should have correct last name")
			assert.Equal(t, testUser.Phone, response.User.Phone, "Should have correct phone")
			assert.Equal(t, testUser.Role, response.User.Role, "Should have correct role")
			assert.False(t, response.User.IsEmailVerified, "Email should not be verified initially")
			assert.NotEmpty(t, response.User.ID, "Should have user ID")
			assert.NotZero(t, response.User.CreatedAt, "Should have creation timestamp")

			// Test: User should be able to login with created credentials
			loginReq := &models.LoginRequest{
				Email:    testUser.Email,
				Password: testUser.Password,
			}
			loginResponse, err := authService.Login(ctx, loginReq)
			require.NoError(t, err, "Should be able to login with created credentials")
			require.NotNil(t, loginResponse, "Should return login response")
			assert.Equal(t, response.User.ID, loginResponse.User.ID, "Login should return same user")
		})
	}

	// Test: Verify all test users can login
	t.Run("Verify all test users can login", func(t *testing.T) {
		testCredentials := []struct {
			email    string
			password string
		}{
			{"qa_customer@smorting.com", "TestPass123!"},
			{"qa_provider@smorting.com", "ProviderPass123!"},
			{"qa_admin@smorting.com", "AdminPass123!"},
		}

		for _, cred := range testCredentials {
			loginReq := &models.LoginRequest{
				Email:    cred.email,
				Password: cred.password,
			}
			response, err := authService.Login(ctx, loginReq)
			require.NoError(t, err, "Test user %s should be able to login", cred.email)
			assert.Equal(t, cred.email, response.User.Email, "User email should match")
		}
	})
}

// TestQATestUsersDataIntegrity verifies the integrity of QA test user data by login
func TestQATestUsersDataIntegrity(t *testing.T) {
	ctx := context.Background()
	authService := setupAuthService(t)

	t.Run("QA Customer user data integrity", func(t *testing.T) {
		loginReq := &models.LoginRequest{
			Email:    "qa_customer@smorting.com",
			Password: "TestPass123!",
		}
		response, err := authService.Login(ctx, loginReq)
		require.NoError(t, err, "QA customer should exist and login")

		user := response.User
		assert.Equal(t, models.CustomerRole, user.Role, "Should have customer role")
		assert.Equal(t, "QA", user.FirstName, "Should have correct first name")
		assert.Equal(t, "Customer", user.LastName, "Should have correct last name")
		assert.Equal(t, "231777123456", user.Phone, "Should have correct phone")
		assert.NotEmpty(t, user.Wallet.Currency, "Should have wallet currency")
		assert.Equal(t, float64(0), user.Wallet.Balance, "Initial wallet balance should be 0")
	})

	t.Run("QA Provider user data integrity", func(t *testing.T) {
		loginReq := &models.LoginRequest{
			Email:    "qa_provider@smorting.com",
			Password: "ProviderPass123!",
		}
		response, err := authService.Login(ctx, loginReq)
		require.NoError(t, err, "QA provider should exist and login")

		user := response.User
		assert.Equal(t, models.ProviderRole, user.Role, "Should have provider role")
		assert.Equal(t, "QA", user.FirstName, "Should have correct first name")
		assert.Equal(t, "Provider", user.LastName, "Should have correct last name")
		assert.Equal(t, "231888123456", user.Phone, "Should have correct phone")
	})

	t.Run("QA Admin user data integrity", func(t *testing.T) {
		loginReq := &models.LoginRequest{
			Email:    "qa_admin@smorting.com",
			Password: "AdminPass123!",
		}
		response, err := authService.Login(ctx, loginReq)
		require.NoError(t, err, "QA admin should exist and login")

		user := response.User
		assert.Equal(t, models.AdminRole, user.Role, "Should have admin role")
		assert.Equal(t, "QA", user.FirstName, "Should have correct first name")
		assert.Equal(t, "Admin", user.LastName, "Should have correct last name")
		assert.Equal(t, "231999123456", user.Phone, "Should have correct phone")
	})
}

// TestQATestUsersAuthentication verifies authentication works for all test users
func TestQATestUsersAuthentication(t *testing.T) {
	ctx := context.Background()
	authService := setupAuthService(t)

	testCases := []struct {
		name     string
		email    string
		password string
		role     string
	}{
		{
			name:     "QA Customer Authentication",
			email:    "qa_customer@smorting.com",
			password: "TestPass123!",
			role:     "customer",
		},
		{
			name:     "QA Provider Authentication",
			email:    "qa_provider@smorting.com",
			password: "ProviderPass123!",
			role:     "provider",
		},
		{
			name:     "QA Admin Authentication",
			email:    "qa_admin@smorting.com",
			password: "AdminPass123!",
			role:     "admin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loginReq := &models.LoginRequest{
				Email:    tc.email,
				Password: tc.password,
			}

			response, err := authService.Login(ctx, loginReq)
			require.NoError(t, err, "Should authenticate successfully")
			require.NotNil(t, response, "Should return auth response")
			require.NotEmpty(t, response.AccessToken, "Should return access token")

			assert.Equal(t, tc.email, response.User.Email, "Should return correct user")
			assert.Equal(t, tc.role, response.User.Role, "Should have correct role")
			assert.NotEmpty(t, response.User.ID, "Should have user ID")
		})
	}
}
