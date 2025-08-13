package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCoreAuthFlowIntegration tests the complete authentication flow end-to-end
func TestCoreAuthFlowIntegration(t *testing.T) {
	// Setup test environment with memory database
	app, _, repo := setupTestApp(t)

	t.Run("Complete Auth Flow", func(t *testing.T) {
		// Test 1: Register a new user
		t.Run("User Registration", func(t *testing.T) {
			registerReq := models.RegisterRequest{
				Email:     "flowtest@example.com",
				Password:  "password123",
				FirstName: "Flow",
				LastName:  "Test",
				Phone:     "231777123456",
				Role:      models.CustomerRole,
			}

			body, _ := json.Marshal(registerReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var authResp models.AuthResponse
			err = json.NewDecoder(resp.Body).Decode(&authResp)
			require.NoError(t, err)

			assert.Equal(t, registerReq.Email, authResp.User.Email)
			assert.Equal(t, registerReq.FirstName, authResp.User.FirstName)
			assert.Equal(t, registerReq.LastName, authResp.User.LastName)
			assert.NotEmpty(t, authResp.AccessToken)
			assert.NotEmpty(t, authResp.RefreshToken)
			assert.False(t, authResp.RequiresOTP)
		})

		// Test 2: Login with correct credentials
		t.Run("Successful Login", func(t *testing.T) {
			loginReq := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "flowtest@example.com",
				Password: "password123",
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var authResp models.AuthResponse
			err = json.NewDecoder(resp.Body).Decode(&authResp)
			require.NoError(t, err)

			assert.Equal(t, loginReq.Email, authResp.User.Email)
			assert.NotEmpty(t, authResp.AccessToken)
			assert.NotEmpty(t, authResp.RefreshToken)
		})

		// Test 3: Login with wrong credentials
		t.Run("Failed Login - Wrong Password", func(t *testing.T) {
			loginReq := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "flowtest@example.com",
				Password: "wrongpassword",
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			var errorResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)

			assert.Equal(t, "Invalid credentials", errorResp["error"])
		})

		// Test 4: Token validation
		t.Run("Token Validation", func(t *testing.T) {
			// First login to get a token
			loginReq := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "flowtest@example.com",
				Password: "password123",
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)

			var authResp models.AuthResponse
			err = json.NewDecoder(resp.Body).Decode(&authResp)
			require.NoError(t, err)
			resp.Body.Close()

			// Now validate the token
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/validate", nil)
			req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)

			resp, err = app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var validationResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&validationResp)
			require.NoError(t, err)

			assert.Equal(t, "Token is valid", validationResp["message"])

			data := validationResp["data"].(map[string]interface{})
			assert.Equal(t, loginReq.Email, data["email"])
			assert.Equal(t, "customer", data["role"])
		})

		// Test 5: Token refresh
		t.Run("Token Refresh", func(t *testing.T) {
			// First login to get tokens
			loginReq := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "flowtest@example.com",
				Password: "password123",
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)

			var authResp models.AuthResponse
			err = json.NewDecoder(resp.Body).Decode(&authResp)
			require.NoError(t, err)
			resp.Body.Close()

			// Now refresh the token
			refreshReq := struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: authResp.RefreshToken,
			}

			body, _ = json.Marshal(refreshReq)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var refreshResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&refreshResp)
			require.NoError(t, err)

			assert.Equal(t, "Token refreshed successfully", refreshResp["message"])

			data := refreshResp["data"].(map[string]interface{})
			assert.NotEmpty(t, data["access_token"])
			assert.NotEmpty(t, data["refresh_token"])
			assert.NotEqual(t, authResp.AccessToken, data["access_token"]) // Should be a new token
		})

		// Test 6: Password reset flow
		t.Run("Password Reset Flow", func(t *testing.T) {
			// Request password reset
			resetReq := struct {
				Email string `json:"email"`
			}{
				Email: "flowtest@example.com",
			}

			body, _ := json.Marshal(resetReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/request-password-reset", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Get the OTP from the repository (for testing)
			otp, err := repo.GetLatestOTPByEmail(context.Background(), "flowtest@example.com")
			require.NoError(t, err)
			require.NotNil(t, otp)

			// Reset password with OTP
			passwordResetReq := struct {
				Email       string `json:"email"`
				OTP         string `json:"otp"`
				NewPassword string `json:"new_password"`
			}{
				Email:       "flowtest@example.com",
				OTP:         otp.OTP,
				NewPassword: "newpassword456",
			}

			body, _ = json.Marshal(passwordResetReq)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var resetResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&resetResp)
			require.NoError(t, err)

			assert.Equal(t, "Password reset successful", resetResp["message"])

			// Verify login with new password works
			loginReq := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "flowtest@example.com",
				Password: "newpassword456",
			}

			body, _ = json.Marshal(loginReq)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify old password no longer works
			loginReq.Password = "password123"
			body, _ = json.Marshal(loginReq)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}

// TestAuthValidationRules tests validation rules for authentication endpoints
func TestAuthValidationRules(t *testing.T) {
	app, _, _ := setupTestApp(t)

	t.Run("Registration Validation", func(t *testing.T) {
		testCases := []struct {
			name           string
			request        models.RegisterRequest
			expectedStatus int
			expectedError  string
		}{
			{
				name: "Missing Email",
				request: models.RegisterRequest{
					Password:  "password123",
					FirstName: "John",
					LastName:  "Doe",
					Phone:     "231777123456",
					Role:      models.CustomerRole,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "email is required",
			},
			{
				name: "Short Password",
				request: models.RegisterRequest{
					Email:     "test@example.com",
					Password:  "123",
					FirstName: "John",
					LastName:  "Doe",
					Phone:     "231777123456",
					Role:      models.CustomerRole,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "password must be at least 6 characters long",
			},
			{
				name: "Invalid Role",
				request: models.RegisterRequest{
					Email:     "test@example.com",
					Password:  "password123",
					FirstName: "John",
					LastName:  "Doe",
					Phone:     "231777123456",
					Role:      "invalid_role",
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "role must be 'customer', 'provider', or 'admin'",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.request)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, tc.expectedStatus, resp.StatusCode)

				var errorResp map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				require.NoError(t, err)

				assert.Contains(t, errorResp["message"], tc.expectedError)
			})
		}
	})
}

// setupTestApp creates a test Fiber app with in-memory database for testing
func setupTestApp(t *testing.T) (*fiber.App, *handlers.AuthHandler, *database.MemoryDatabase) {
	t.Helper()

	// Setup logger
	lg, err := logger.New("debug", "console", "stdout")
	require.NoError(t, err)

	// Setup memory database
	repo := database.NewMemoryDatabase()

	// Setup JWT service
	accessSecret := make([]byte, 32)
	refreshSecret := make([]byte, 32)
	for i := range accessSecret {
		accessSecret[i] = byte(i + 1)
	}
	for i := range refreshSecret {
		refreshSecret[i] = byte(i + 32)
	}
	jwtService := services.NewJWTRefreshService(accessSecret, refreshSecret, lg.Logger)

	// Setup encryption service
	encKey := make([]byte, 32)
	for i := range encKey {
		encKey[i] = byte(i + 64)
	}
	encryptionService, err := services.NewEncryptionService(encKey)
	require.NoError(t, err)

	// Setup auth config
	authConfig := &configs.AuthConfig{
		JWTSecret:     "test-secret-for-legacy-jwt-at-least-32-chars-long",
		JWTExpiration: 30 * time.Minute,
		BCryptCost:    10,
	}

	// Setup MongoDB auth service
	authSvc, err := auth.NewMongoDBService(repo, authConfig, lg)
	require.NoError(t, err)

	// Setup auth handler
	authHandler := handlers.NewAuthHandler(jwtService, encryptionService, lg, authSvc)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   "Request failed",
				"message": err.Error(),
			})
		},
	})

	// Setup routes
	api := app.Group("/api/v1")
	auth := api.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/validate", authHandler.ValidateToken)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/verify-otp", authHandler.VerifyOTP)
	auth.Post("/resend-otp", authHandler.ResendOTP)
	auth.Post("/request-password-reset", authHandler.RequestPasswordReset)
	auth.Post("/reset-password", authHandler.ResetPassword)
	auth.Get("/test/get-latest-otp", authHandler.TestGetLatestOTP)

	return app, authHandler, repo
}

// createTestJWTService creates a JWT service for testing
func createTestJWTService(accessSecret, refreshSecret []byte) *services.JWTRefreshService {
	logger, _ := logger.New("debug", "console", "stdout")
	return services.NewJWTRefreshService(accessSecret, refreshSecret, logger.Logger)
}
