package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// TestEnhancedAuthIntegration_LoginShouldNotRequire2FA tests that the enhanced auth login
// endpoint does not return requires_two_factor: true for any user type
func TestEnhancedAuthIntegration_LoginShouldNotRequire2FA(t *testing.T) {
	// Test cases for different user roles
	testCases := []struct {
		name string
		role models.UserRole
	}{
		{"Customer login should not require 2FA", models.CustomerRole},
		{"Provider login should not require 2FA", models.ProviderRole},
		{"Admin login should not require 2FA", models.AdminRole},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			app, userService := setupEnhancedAuthTestApp(t)

			// Create test user with specific role
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
			user := &models.User{
				ID:              primitive.NewObjectID(),
				Email:           "test@example.com",
				Password:        string(hashedPassword),
				FirstName:       "Test",
				LastName:        "User",
				Role:            tc.role,
				IsEmailVerified: false, // Important: unverified email should still not require 2FA
			}
			userService.CreateUser(context.Background(), user)

			// Create login request
			loginReq := models.EnhancedLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
				DeviceInfo: &models.DeviceFingerprint{
					DeviceID:   "test-device-123",
					Platform:   "iOS",
					OSVersion:  "17.0",
					AppVersion: "1.0.0",
				},
			}

			// Make login request
			reqBody, _ := json.Marshal(loginReq)
			req := httptest.NewRequest("POST", "/api/v1/auth/enhanced-login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "SmorTing/1.0 iOS")

			resp, err := app.Test(req)
			require.NoError(t, err)

			// Parse response
			var loginResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&loginResp)
			require.NoError(t, err)

			// Debug: Print the full response to understand what's happening
			t.Logf("Login response for %s: %+v", tc.role, loginResp)

			// Assertions
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Login should succeed")
			assert.True(t, loginResp["success"].(bool), "Login should be successful")

			// The key assertion: requires_two_factor should be false
			if requiresTwoFactor, exists := loginResp["requires_two_factor"]; exists {
				assert.False(t, requiresTwoFactor.(bool),
					"Login should NOT require 2FA for %s role, but got requires_two_factor: %v",
					tc.role, requiresTwoFactor)
			}

			// Additional checks
			assert.NotEmpty(t, loginResp["access_token"], "Should provide access token immediately")
			assert.NotEmpty(t, loginResp["refresh_token"], "Should provide refresh token immediately")
			assert.NotNil(t, loginResp["user"], "Should return user information")
		})
	}
}

// TestEnhancedAuthIntegration_UnverifiedEmailShouldNotRequire2FA specifically tests
// that unverified emails do not trigger 2FA requirement
func TestEnhancedAuthIntegration_UnverifiedEmailShouldNotRequire2FA(t *testing.T) {
	// Setup
	app, userService := setupEnhancedAuthTestApp(t)

	// Create test user with unverified email
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:              primitive.NewObjectID(),
		Email:           "unverified@example.com",
		Password:        string(hashedPassword),
		FirstName:       "Unverified",
		LastName:        "User",
		Role:            models.CustomerRole,
		IsEmailVerified: false, // Explicitly unverified
	}
	userService.CreateUser(context.Background(), user)

	// Create login request
	loginReq := models.EnhancedLoginRequest{
		Email:    "unverified@example.com",
		Password: "password123",
		DeviceInfo: &models.DeviceFingerprint{
			DeviceID:   "test-device-456",
			Platform:   "Android",
			OSVersion:  "14.0",
			AppVersion: "1.0.0",
		},
	}

	// Make login request
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/enhanced-login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SmorTing/1.0 Android")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Parse response
	var loginResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(t, err)

	// Debug: Print the full response
	t.Logf("Unverified email login response: %+v", loginResp)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Login should succeed even with unverified email")
	assert.True(t, loginResp["success"].(bool), "Login should be successful")

	// The critical assertion: even unverified emails should not require 2FA
	if requiresTwoFactor, exists := loginResp["requires_two_factor"]; exists {
		assert.False(t, requiresTwoFactor.(bool),
			"Unverified email login should NOT require 2FA, but got requires_two_factor: %v",
			requiresTwoFactor)
	}

	// Should get tokens immediately
	assert.NotEmpty(t, loginResp["access_token"], "Should provide access token for unverified email")
	assert.NotEmpty(t, loginResp["refresh_token"], "Should provide refresh token for unverified email")
}

// setupEnhancedAuthTestApp creates a test Fiber app with enhanced auth handlers
func setupEnhancedAuthTestApp(t *testing.T) (*fiber.App, *services.MockUserService) {
	app := fiber.New()

	// Create mock services
	userService := services.NewMockUserService()

	// Create enhanced auth service with stubs
	logger := zap.NewNop()
	enhancedAuthService := services.NewStubEnhancedAuthService(logger)

	// Create handler
	// Create mock services for all dependencies
	otpService := &services.MockOTPService{}
	captchaService := &services.MockCaptchaService{}
	auditService := services.NewAuditService(nil, logger)

	authHandler := handlers.NewEnhancedAuthHandler(
		enhancedAuthService,
		userService,
		otpService,
		captchaService,
		auditService,
		logger,
	)

	// Setup routes
	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	auth.Post("/enhanced-login", authHandler.EnhancedLogin)

	return app, userService
}
