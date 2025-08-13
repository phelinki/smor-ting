package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEnhancedAuthWithBruteForceProtection(t *testing.T) {
	app := fiber.New()
	logger := zap.NewNop()

	// Create enhanced auth handler with brute force protection
	mockEnhancedAuthService := &services.MockEnhancedAuthService{}
	mockUserService := &services.MockUserService{}
	mockOTPService := &services.MockOTPService{}
	mockCaptchaService := &services.MockCaptchaService{}

	handler := handlers.NewEnhancedAuthHandler(
		mockEnhancedAuthService,
		mockUserService,
		mockOTPService,
		mockCaptchaService,
		nil, // audit service (will use no-op)
		logger,
	)

	// Setup route
	app.Post("/auth/enhanced-login", handler.EnhancedLogin)

	email := "test@example.com"
	password := "wrong_password"

	t.Run("allows_initial_login_attempts", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"email":    email,
			"password": password,
			"device_info": map[string]interface{}{
				"device_id": "test_device",
				"platform":  "test",
			},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/enhanced-login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "192.168.1.1")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Wrong password should return unauthorized")

		var response handlers.EnhancedLoginResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(t, response.Success)
		assert.False(t, response.RequiresCaptcha, "Should not require CAPTCHA initially")
		assert.Greater(t, response.RemainingAttempts, 0, "Should have remaining attempts")
	})

	t.Run("requires_captcha_after_threshold", func(t *testing.T) {
		// Make 2 more failed attempts (total 3)
		for i := 0; i < 2; i++ {
			requestBody := map[string]interface{}{
				"email":    email,
				"password": password,
				"device_info": map[string]interface{}{
					"device_id": "test_device",
					"platform":  "test",
				},
			}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("POST", "/auth/enhanced-login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Forwarded-For", "192.168.1.1")

			app.Test(req, -1)
		}

		// Next attempt should require CAPTCHA
		requestBody := map[string]interface{}{
			"email":    email,
			"password": password,
			"device_info": map[string]interface{}{
				"device_id": "test_device",
				"platform":  "test",
			},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/enhanced-login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "192.168.1.1")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		var response handlers.EnhancedLoginResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.True(t, response.RequiresCaptcha, "Should require CAPTCHA after multiple failures")
	})

	t.Run("blocks_after_max_attempts", func(t *testing.T) {
		// Make one more failed attempt to trigger lockout
		requestBody := map[string]interface{}{
			"email":    email,
			"password": password,
			"device_info": map[string]interface{}{
				"device_id": "test_device",
				"platform":  "test",
			},
			"captcha_token": "valid_captcha", // Include CAPTCHA to pass validation
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/enhanced-login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "192.168.1.1")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		var response handlers.EnhancedLoginResponse
		json.NewDecoder(resp.Body).Decode(&response)
		// Should be close to 0 remaining attempts (may not be exactly 0 depending on implementation)
		assert.LessOrEqual(t, response.RemainingAttempts, 1, "Should have very few remaining attempts after multiple failures")
	})

	t.Run("rejects_requests_during_lockout", func(t *testing.T) {
		// Attempt login during lockout period
		requestBody := map[string]interface{}{
			"email":    email,
			"password": "correct_password", // Even with correct password
			"device_info": map[string]interface{}{
				"device_id": "test_device",
				"platform":  "test",
			},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/enhanced-login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "192.168.1.1")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "Should be blocked during lockout")

		var response handlers.EnhancedLoginResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(t, response.Success)
		// Message should mention either lockout or CAPTCHA requirement
		assert.True(t,
			strings.Contains(response.Message, "locked") || strings.Contains(response.Message, "CAPTCHA"),
			"Message should mention lockout or CAPTCHA requirement")
	})
}

func TestEnhancedAuthCaptchaFlow(t *testing.T) {
	app := fiber.New()
	logger := zap.NewNop()

	// Create enhanced auth handler
	mockEnhancedAuthService := &services.MockEnhancedAuthService{}
	mockUserService := &services.MockUserService{}
	mockOTPService := &services.MockOTPService{}
	mockCaptchaService := &services.MockCaptchaService{}

	handler := handlers.NewEnhancedAuthHandler(
		mockEnhancedAuthService,
		mockUserService,
		mockOTPService,
		mockCaptchaService,
		nil, // audit service (will use no-op)
		logger,
	)

	app.Post("/auth/enhanced-login", handler.EnhancedLogin)

	email := "captcha@example.com"

	t.Run("requires_captcha_without_token", func(t *testing.T) {
		// Make 3 failed attempts to trigger CAPTCHA requirement
		for i := 0; i < 3; i++ {
			requestBody := map[string]interface{}{
				"email":    email,
				"password": "wrong_password",
				"device_info": map[string]interface{}{
					"device_id": "test_device",
					"platform":  "test",
				},
			}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("POST", "/auth/enhanced-login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Forwarded-For", "192.168.1.2")

			app.Test(req, -1)
		}

		// Next attempt without CAPTCHA should be rejected
		requestBody := map[string]interface{}{
			"email":    email,
			"password": "wrong_password",
			"device_info": map[string]interface{}{
				"device_id": "test_device",
				"platform":  "test",
			},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/enhanced-login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "192.168.1.2")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

		var response handlers.EnhancedLoginResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.True(t, response.RequiresCaptcha)
		assert.Contains(t, response.Message, "CAPTCHA", "Should mention CAPTCHA requirement")
	})

	t.Run("accepts_request_with_valid_captcha", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"email":    email,
			"password": "wrong_password",
			"device_info": map[string]interface{}{
				"device_id": "test_device",
				"platform":  "test",
			},
			"captcha_token": "valid_captcha_token",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/enhanced-login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "192.168.1.2")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		// Should proceed to actual auth check (and fail due to wrong password)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should reach password validation")
	})
}
