package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateToken_ReturnsCompleteUserData(t *testing.T) {
	// Setup test dependencies
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	
	// JWT service
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	jwtService := services.NewJWTRefreshService(access, refresh, lg.Logger)
	
	// Encryption service
	encKey := make([]byte, 32)
	encService, err := services.NewEncryptionService(encKey)
	require.NoError(t, err)
	
	// Auth service
	authCfg := &configs.AuthConfig{
		JWTSecret:     "test-secret-key",
		JWTExpiration: time.Hour,
		BCryptCost:    10,
	}
	authService, err := auth.NewMongoDBService(repo, authCfg, lg)
	require.NoError(t, err)
	
	// Create auth handler
	authHandler := handlers.NewAuthHandler(jwtService, encService, lg, authService)
	
	// Create test user with complete data
	testUser := &models.User{
		ID:              primitive.NewObjectID(),
		Email:           "test@example.com",
		FirstName:       "John",
		LastName:        "Doe",
		Phone:           "1234567890",
		Role:            models.CustomerRole,
		IsEmailVerified: true,
		ProfileImage:    "https://example.com/profile.jpg",
		Address: &models.Address{
			Street:    "123 Test St",
			City:      "Test City",
			County:    "Test County",
			Country:   "Test Country",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	
	// Store user in repository
	err = repo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	
	// Generate token for the user
	tokenPair, err := jwtService.GenerateTokenPair(testUser)
	require.NoError(t, err)
	
	// Setup Fiber app
	app := fiber.New()
	app.Post("/validate", authHandler.ValidateToken)
	
	// Test cases
	tests := []struct {
		name         string
		setupUser    func() *models.User
		authHeader   string
		expectedCode int
		checkResp    func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "Valid token returns complete user data",
			setupUser: func() *models.User {
				return testUser
			},
			authHeader:   "Bearer " + tokenPair.AccessToken,
			expectedCode: http.StatusOK,
			checkResp: func(t *testing.T, body map[string]interface{}) {
				// Check response structure
				assert.Equal(t, "Token is valid", body["message"])
				assert.Contains(t, body, "user")
				assert.Contains(t, body, "token_info")
				assert.Contains(t, body, "permissions")
				
				// Check complete user data
				user := body["user"].(map[string]interface{})
				assert.Equal(t, testUser.ID.Hex(), user["id"])
				assert.Equal(t, testUser.Email, user["email"])
				assert.Equal(t, testUser.FirstName, user["first_name"])
				assert.Equal(t, testUser.LastName, user["last_name"])
				assert.Equal(t, testUser.Phone, user["phone"])
				assert.Equal(t, string(testUser.Role), user["role"])
				assert.Equal(t, testUser.IsEmailVerified, user["is_email_verified"])
				assert.Equal(t, testUser.ProfileImage, user["profile_image"])
				assert.Contains(t, user, "created_at")
				assert.Contains(t, user, "updated_at")
				
				// Check address data
				assert.Contains(t, user, "address")
				address := user["address"].(map[string]interface{})
				assert.Equal(t, testUser.Address.Street, address["street"])
				assert.Equal(t, testUser.Address.City, address["city"])
				assert.Equal(t, testUser.Address.County, address["county"])
				assert.Equal(t, testUser.Address.Country, address["country"])
				assert.Equal(t, testUser.Address.Latitude, address["latitude"])
				assert.Equal(t, testUser.Address.Longitude, address["longitude"])
				
				// Check permissions
				permissions := body["permissions"].(map[string]interface{})
				assert.Equal(t, true, permissions["is_customer"])
				assert.Equal(t, false, permissions["is_provider"])
				assert.Equal(t, false, permissions["is_admin"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("POST", "/validate", nil)
			req.Header.Set("Authorization", tt.authHeader)
			req.Header.Set("Content-Type", "application/json")
			
			// Send request
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			
			// Check status code
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			
			// Parse and check response
			if tt.checkResp != nil && resp.StatusCode == http.StatusOK {
				var body map[string]interface{}
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				err = json.Unmarshal(bodyBytes, &body)
				require.NoError(t, err)
				tt.checkResp(t, body)
			}
		})
	}
}

func TestValidateToken_UserWithoutAddress(t *testing.T) {
	// Setup test dependencies
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	
	// JWT service
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	jwtService := services.NewJWTRefreshService(access, refresh, lg.Logger)
	
	// Encryption service
	encKey := make([]byte, 32)
	encService, err := services.NewEncryptionService(encKey)
	require.NoError(t, err)
	
	// Auth service
	authCfg := &configs.AuthConfig{
		JWTSecret:     "test-secret-key",
		JWTExpiration: time.Hour,
		BCryptCost:    10,
	}
	authService, err := auth.NewMongoDBService(repo, authCfg, lg)
	require.NoError(t, err)
	
	// Create auth handler
	authHandler := handlers.NewAuthHandler(jwtService, encService, lg, authService)
	
	// Create test user without address
	testUser := &models.User{
		ID:              primitive.NewObjectID(),
		Email:           "noaddress@example.com",
		FirstName:       "Jane",
		LastName:        "Smith",
		Phone:           "9876543210",
		Role:            models.ProviderRole,
		IsEmailVerified: false,
		ProfileImage:    "",
		Address:         nil, // No address
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	
	// Store user in repository
	err = repo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	
	// Generate token for the user
	tokenPair, err := jwtService.GenerateTokenPair(testUser)
	require.NoError(t, err)
	
	// Setup Fiber app
	app := fiber.New()
	app.Post("/validate", authHandler.ValidateToken)
	
	// Create request
	req := httptest.NewRequest("POST", "/validate", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Parse response
	var body map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(bodyBytes, &body)
	require.NoError(t, err)
	
	// Check that empty address is provided
	user := body["user"].(map[string]interface{})
	assert.Contains(t, user, "address")
	address := user["address"].(map[string]interface{})
	assert.Equal(t, "", address["street"])
	assert.Equal(t, "", address["city"])
	assert.Equal(t, "", address["county"])
	assert.Equal(t, "", address["country"])
	assert.Equal(t, 0.0, address["latitude"])
	assert.Equal(t, 0.0, address["longitude"])
	
	// Check permissions for provider role
	permissions := body["permissions"].(map[string]interface{})
	assert.Equal(t, false, permissions["is_customer"])
	assert.Equal(t, true, permissions["is_provider"])
	assert.Equal(t, false, permissions["is_admin"])
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Setup test dependencies
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	
	// JWT service
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	jwtService := services.NewJWTRefreshService(access, refresh, lg.Logger)
	
	// Encryption service
	encKey := make([]byte, 32)
	encService, err := services.NewEncryptionService(encKey)
	require.NoError(t, err)
	
	// Auth service
	authCfg := &configs.AuthConfig{
		JWTSecret:     "test-secret-key",
		JWTExpiration: time.Hour,
		BCryptCost:    10,
	}
	authService, err := auth.NewMongoDBService(repo, authCfg, lg)
	require.NoError(t, err)
	
	// Create auth handler
	authHandler := handlers.NewAuthHandler(jwtService, encService, lg, authService)
	
	// Setup Fiber app
	app := fiber.New()
	app.Post("/validate", authHandler.ValidateToken)
	
	// Test cases for invalid tokens
	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
	}{
		{
			name:         "Missing authorization header",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid token format",
			authHeader:   "invalid-token",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid token",
			authHeader:   "Bearer invalid-token",
			expectedCode: http.StatusUnauthorized,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/validate", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}