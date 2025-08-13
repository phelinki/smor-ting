package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestAuthRefreshEndpointMismatch tests that mobile calls /auth/refresh-token but backend only has /auth/refresh
// This test exposes the API mismatch issue
func TestAuthRefreshEndpointMismatch(t *testing.T) {
	app := fiber.New()

	// Mock handler that returns success for any request
	mockHandler := func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":      true,
			"access_token": "new_token",
		})
	}

	// Setup routes as they currently are in main.go (only /auth/refresh)
	auth := app.Group("/auth")
	auth.Post("/refresh", mockHandler)

	// Test that /auth/refresh endpoint exists
	t.Run("current_auth_refresh_endpoint_exists", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"refresh_token": "valid_refresh_token",
			"session_id":    "session123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "/auth/refresh endpoint should exist")
	})

	// Test that /auth/refresh-token endpoint is missing (exposing the mismatch)
	t.Run("mobile_auth_refresh_token_endpoint_missing", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"refresh_token": "valid_refresh_token",
			"session_id":    "session123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/refresh-token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		// This should be 404, exposing the mismatch between mobile and backend
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "/auth/refresh-token endpoint should be missing, exposing the mismatch")
	})
}

// TestMobileAPICompatibilityFix tests that we can fix the mismatch by adding the missing route
func TestMobileAPICompatibilityFix(t *testing.T) {
	app := fiber.New()

	// Mock handler that returns success for any request
	mockHandler := func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":      true,
			"access_token": "new_token",
		})
	}

	// Setup routes - add both routes to fix the mismatch
	auth := app.Group("/auth")
	auth.Post("/refresh", mockHandler)       // Current backend route
	auth.Post("/refresh-token", mockHandler) // Mobile expected route

	t.Run("both_refresh_endpoints_work_after_fix", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"refresh_token": "valid_refresh_token",
			"session_id":    "session123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		// Test /auth/refresh (current backend)
		req1 := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req1.Header.Set("Content-Type", "application/json")
		resp1, err1 := app.Test(req1, -1)
		assert.NoError(t, err1)
		assert.Equal(t, http.StatusOK, resp1.StatusCode, "/auth/refresh should work")

		// Test /auth/refresh-token (mobile expected)
		req2 := httptest.NewRequest("POST", "/auth/refresh-token", bytes.NewBuffer(jsonBody))
		req2.Header.Set("Content-Type", "application/json")
		resp2, err2 := app.Test(req2, -1)
		assert.NoError(t, err2)
		assert.Equal(t, http.StatusOK, resp2.StatusCode, "/auth/refresh-token should work after fix")
	})
}
