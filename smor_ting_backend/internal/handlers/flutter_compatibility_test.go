package handlers_test

import (
	"context"
	"encoding/json"
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

func TestValidateToken_FlutterAppCompatibility(t *testing.T) {
	// This test ensures the ValidateToken response matches the exact format expected by the Flutter app
	
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
	
	// Create test user with all fields as expected by Flutter
	testUser := &models.User{
		ID:              primitive.NewObjectID(),
		Email:           "flutter@example.com",
		FirstName:       "Flutter",
		LastName:        "User",
		Phone:           "+1234567890",
		Role:            models.CustomerRole,
		IsEmailVerified: true,
		ProfileImage:    "https://example.com/profile.jpg",
		Address: &models.Address{
			Street:    "123 Flutter St",
			City:      "Mobile City",
			County:    "App County", 
			Country:   "DevLand",
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
		CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
	}
	
	// Store user in repository
	err = repo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	
	// Generate token for the user
	tokenPair, err := jwtService.GenerateTokenPair(testUser)
	require.NoError(t, err)
	
	// Setup Fiber app
	app := fiber.New()
	app.Post("/auth/validate", authHandler.ValidateToken)
	
	// Create request as Flutter app would send
	req := httptest.NewRequest("POST", "/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// Should return 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Parse the actual response 
	var actualResponse map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&actualResponse)
	require.NoError(t, err)
	
	// Test the exact response structure expected by Flutter
	expectedStructure := map[string]interface{}{
		"message": "Token is valid",
		"user": map[string]interface{}{
			"id":                testUser.ID.Hex(),
			"email":             "flutter@example.com",
			"first_name":        "Flutter",
			"last_name":         "User",
			"phone":             "+1234567890",
			"role":              "customer",
			"is_email_verified": true,
			"profile_image":     "https://example.com/profile.jpg",
			"address": map[string]interface{}{
				"street":    "123 Flutter St",
				"city":      "Mobile City",
				"county":    "App County",
				"country":   "DevLand",
				"latitude":  37.7749,
				"longitude": -122.4194,
			},
		},
		"token_info":  map[string]interface{}{}, // Structure varies, just check it exists
		"permissions": map[string]interface{}{
			"is_customer": true,
			"is_provider": false,
			"is_admin":    false,
		},
	}
	
	// Verify top-level structure
	assert.Equal(t, expectedStructure["message"], actualResponse["message"])
	assert.Contains(t, actualResponse, "user")
	assert.Contains(t, actualResponse, "token_info")
	assert.Contains(t, actualResponse, "permissions")
	
	// Verify user object structure
	actualUser := actualResponse["user"].(map[string]interface{})
	expectedUser := expectedStructure["user"].(map[string]interface{})
	
	assert.Equal(t, expectedUser["id"], actualUser["id"])
	assert.Equal(t, expectedUser["email"], actualUser["email"])
	assert.Equal(t, expectedUser["first_name"], actualUser["first_name"])
	assert.Equal(t, expectedUser["last_name"], actualUser["last_name"])
	assert.Equal(t, expectedUser["phone"], actualUser["phone"])
	assert.Equal(t, expectedUser["role"], actualUser["role"])
	assert.Equal(t, expectedUser["is_email_verified"], actualUser["is_email_verified"])
	assert.Equal(t, expectedUser["profile_image"], actualUser["profile_image"])
	
	// Verify timestamps are present
	assert.Contains(t, actualUser, "created_at")
	assert.Contains(t, actualUser, "updated_at")
	assert.IsType(t, "", actualUser["created_at"]) // Should be string (ISO format)
	assert.IsType(t, "", actualUser["updated_at"]) // Should be string (ISO format)
	
	// Verify address structure
	assert.Contains(t, actualUser, "address")
	actualAddress := actualUser["address"].(map[string]interface{})
	expectedAddress := expectedUser["address"].(map[string]interface{})
	
	assert.Equal(t, expectedAddress["street"], actualAddress["street"])
	assert.Equal(t, expectedAddress["city"], actualAddress["city"])
	assert.Equal(t, expectedAddress["county"], actualAddress["county"])
	assert.Equal(t, expectedAddress["country"], actualAddress["country"])
	assert.Equal(t, expectedAddress["latitude"], actualAddress["latitude"])
	assert.Equal(t, expectedAddress["longitude"], actualAddress["longitude"])
	
	// Verify permissions structure
	actualPermissions := actualResponse["permissions"].(map[string]interface{})
	expectedPermissions := expectedStructure["permissions"].(map[string]interface{})
	
	assert.Equal(t, expectedPermissions["is_customer"], actualPermissions["is_customer"])
	assert.Equal(t, expectedPermissions["is_provider"], actualPermissions["is_provider"])
	assert.Equal(t, expectedPermissions["is_admin"], actualPermissions["is_admin"])
	
	// Print the actual response for manual verification
	t.Logf("Actual response structure: %+v", actualResponse)
}

func TestValidateToken_EmptyAddressHandling(t *testing.T) {
	// Test that users without address get the empty address structure expected by Flutter
	
	// Setup test dependencies (same as above)
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	jwtService := services.NewJWTRefreshService(access, refresh, lg.Logger)
	
	encKey := make([]byte, 32)
	encService, err := services.NewEncryptionService(encKey)
	require.NoError(t, err)
	
	authCfg := &configs.AuthConfig{
		JWTSecret:     "test-secret-key",
		JWTExpiration: time.Hour,
		BCryptCost:    10,
	}
	authService, err := auth.NewMongoDBService(repo, authCfg, lg)
	require.NoError(t, err)
	
	authHandler := handlers.NewAuthHandler(jwtService, encService, lg, authService)
	
	// Create test user WITHOUT address
	testUser := &models.User{
		ID:              primitive.NewObjectID(),
		Email:           "no-address@example.com",
		FirstName:       "No",
		LastName:        "Address",
		Phone:           "",
		Role:            models.CustomerRole,
		IsEmailVerified: false,
		ProfileImage:    "",
		Address:         nil, // No address
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	
	err = repo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	
	tokenPair, err := jwtService.GenerateTokenPair(testUser)
	require.NoError(t, err)
	
	app := fiber.New()
	app.Post("/auth/validate", authHandler.ValidateToken)
	
	req := httptest.NewRequest("POST", "/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var actualResponse map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&actualResponse)
	require.NoError(t, err)
	
	// Verify empty address structure is provided for Flutter compatibility
	user := actualResponse["user"].(map[string]interface{})
	assert.Contains(t, user, "address")
	
	address := user["address"].(map[string]interface{})
	expectedEmptyAddress := map[string]interface{}{
		"street":    "",
		"city":      "",
		"county":    "",
		"country":   "",
		"latitude":  0.0,
		"longitude": 0.0,
	}
	
	assert.Equal(t, expectedEmptyAddress["street"], address["street"])
	assert.Equal(t, expectedEmptyAddress["city"], address["city"])
	assert.Equal(t, expectedEmptyAddress["county"], address["county"])
	assert.Equal(t, expectedEmptyAddress["country"], address["country"])
	assert.Equal(t, expectedEmptyAddress["latitude"], address["latitude"])
	assert.Equal(t, expectedEmptyAddress["longitude"], address["longitude"])
}