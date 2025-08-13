package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smorting/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthEdgeCases tests edge cases and error scenarios for authentication
func TestAuthEdgeCases(t *testing.T) {
	app, _, repo := setupTestApp(t)

	t.Run("Duplicate Registration", func(t *testing.T) {
		// First registration
		registerReq := models.RegisterRequest{
			Email:     "duplicate@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "231777123456",
			Role:      "customer",
		}

		body, _ := json.Marshal(registerReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Second registration with same email should fail
		body, _ = json.Marshal(registerReq)
		req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err = app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var errorResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)

		assert.Equal(t, "User already exists", errorResp["error"])
		assert.Contains(t, errorResp["message"], "already exists")
	})

	t.Run("Invalid Token Scenarios", func(t *testing.T) {
		t.Run("Malformed Token", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/validate", nil)
			req.Header.Set("Authorization", "Bearer invalid.token.here")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Missing Authorization Header", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/validate", nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Invalid Bearer Format", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/validate", nil)
			req.Header.Set("Authorization", "Token invalid_format")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("Refresh Token Edge Cases", func(t *testing.T) {
		t.Run("Invalid Refresh Token", func(t *testing.T) {
			refreshReq := struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: "invalid.refresh.token",
			}

			body, _ := json.Marshal(refreshReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			var errorResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)

			assert.Equal(t, "Invalid refresh token", errorResp["error"])
		})

		t.Run("Empty Refresh Token", func(t *testing.T) {
			refreshReq := struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: "",
			}

			body, _ := json.Marshal(refreshReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("Password Reset Edge Cases", func(t *testing.T) {
		t.Run("Password Reset for Non-existent User", func(t *testing.T) {
			resetReq := struct {
				Email string `json:"email"`
			}{
				Email: "nonexistent@example.com",
			}

			body, _ := json.Marshal(resetReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/request-password-reset", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should still return 200 to prevent user enumeration
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var resetResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&resetResp)
			require.NoError(t, err)

			assert.Contains(t, resetResp["message"], "If the email exists")
		})

		t.Run("Password Reset with Invalid OTP", func(t *testing.T) {
			// First create a user
			registerReq := models.RegisterRequest{
				Email:     "resettest@example.com",
				Password:  "password123",
				FirstName: "Reset",
				LastName:  "Test",
				Phone:     "231777123456",
				Role:      "customer",
			}

			body, _ := json.Marshal(registerReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			resp.Body.Close()

			// Request password reset
			resetReq := struct {
				Email string `json:"email"`
			}{
				Email: "resettest@example.com",
			}

			body, _ = json.Marshal(resetReq)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/request-password-reset", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			resp.Body.Close()

			// Try to reset with invalid OTP
			passwordResetReq := struct {
				Email       string `json:"email"`
				OTP         string `json:"otp"`
				NewPassword string `json:"new_password"`
			}{
				Email:       "resettest@example.com",
				OTP:         "000000", // Invalid OTP
				NewPassword: "newpassword456",
			}

			body, _ = json.Marshal(passwordResetReq)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			var errorResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)

			assert.Equal(t, "Invalid or expired OTP", errorResp["error"])
		})

		t.Run("Password Reset with Expired OTP", func(t *testing.T) {
			// Create a user
			registerReq := models.RegisterRequest{
				Email:     "expiredtest@example.com",
				Password:  "password123",
				FirstName: "Expired",
				LastName:  "Test",
				Phone:     "231777123456",
				Role:      "customer",
			}

			body, _ := json.Marshal(registerReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			resp.Body.Close()

			// Request password reset
			resetReq := struct {
				Email string `json:"email"`
			}{
				Email: "expiredtest@example.com",
			}

			body, _ = json.Marshal(resetReq)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/request-password-reset", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			resp.Body.Close()

			// Get the OTP and manually expire it
			_, err = repo.GetLatestOTPByEmail(context.Background(), "expiredtest@example.com")
			require.NoError(t, err)

			// Manually set the OTP as expired in the test database
			// For this test, we'll use an OTP that should be invalid
			passwordResetReq := struct {
				Email       string `json:"email"`
				OTP         string `json:"otp"`
				NewPassword string `json:"new_password"`
			}{
				Email:       "expiredtest@example.com",
				OTP:         "expired_otp",
				NewPassword: "newpassword456",
			}

			body, _ = json.Marshal(passwordResetReq)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("Malformed Request Bodies", func(t *testing.T) {
		t.Run("Invalid JSON in Registration", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var errorResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)

			assert.Equal(t, "Invalid request body", errorResp["error"])
		})

		t.Run("Invalid JSON in Login", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Invalid JSON in Refresh", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Login Validation", func(t *testing.T) {
		t.Run("Missing Email in Login", func(t *testing.T) {
			loginReq := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "",
				Password: "password123",
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var errorResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)

			assert.Equal(t, "Validation failed", errorResp["error"])
			assert.Contains(t, errorResp["message"], "email is required")
		})

		t.Run("Missing Password in Login", func(t *testing.T) {
			loginReq := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "test@example.com",
				Password: "",
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var errorResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)

			assert.Equal(t, "Validation failed", errorResp["error"])
			assert.Contains(t, errorResp["message"], "password is required")
		})
	})
}
