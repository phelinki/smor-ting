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
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestDeviceSessionHandler tests the device session management endpoints
func TestDeviceSessionHandler(t *testing.T) {
	app, _, repo := setupTestApp(t)

	// Create a test user and device session
	userID := primitive.NewObjectID()
	user := &models.User{
		ID:        userID,
		Email:     "test@example.com",
		Password:  "hashed_password",
		FirstName: "Test",
		LastName:  "User",
		Role:      models.CustomerRole,
	}

	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	sessionID := primitive.NewObjectID()
	session := &models.DeviceSession{
		ID:           sessionID,
		UserID:       userID,
		DeviceID:     "test-device-123",
		DeviceName:   "iPhone 12",
		Platform:     "ios",
		IPAddress:    "192.168.1.1",
		UserAgent:    "SmorTing/1.0 iOS",
		RefreshToken: "test-refresh-token",
		IsActive:     true,
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
	}

	err = repo.CreateDeviceSession(context.Background(), session)
	require.NoError(t, err)

	// Setup device session handler
	logger, _ := logger.New("debug", "console", "stdout")
	deviceHandler := handlers.NewDeviceSessionHandler(repo, nil, nil, logger.Logger)

	// Setup routes with middleware simulation
	api := app.Group("/api/v1")
	sessions := api.Group("/sessions")

	// Middleware to simulate authentication
	sessions.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.Hex())
		c.Locals("sessionID", sessionID.Hex())
		return c.Next()
	})

	sessions.Get("/", deviceHandler.GetDeviceSessions)
	sessions.Delete("/:sessionID", deviceHandler.RevokeDeviceSession)
	sessions.Delete("/", deviceHandler.RevokeAllDeviceSessions)

	t.Run("Get Device Sessions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sessions/", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.DeviceSessionListResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Sessions, 1)
		assert.Equal(t, session.DeviceID, response.Sessions[0].DeviceID)
		assert.Equal(t, session.DeviceName, response.Sessions[0].DeviceName)
		assert.Equal(t, 1, response.TotalSessions)
	})

	t.Run("Revoke Device Session", func(t *testing.T) {
		// Create another session to revoke
		anotherSessionID := primitive.NewObjectID()
		anotherSession := &models.DeviceSession{
			ID:           anotherSessionID,
			UserID:       userID,
			DeviceID:     "test-device-456",
			DeviceName:   "Android Phone",
			Platform:     "android",
			IPAddress:    "192.168.1.2",
			RefreshToken: "another-refresh-token",
			IsActive:     true,
		}

		err := repo.CreateDeviceSession(context.Background(), anotherSession)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/"+anotherSessionID.Hex(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Device session revoked successfully", response["message"])

		// Verify session is revoked
		revokedSession, err := repo.GetDeviceSession(context.Background(), anotherSessionID.Hex())
		require.NoError(t, err)
		assert.False(t, revokedSession.IsActive)
		assert.NotNil(t, revokedSession.RevokedAt)
	})

	t.Run("Revoke All Device Sessions", func(t *testing.T) {
		// Create multiple additional sessions
		for i := 0; i < 3; i++ {
			sessionID := primitive.NewObjectID()
			testSession := &models.DeviceSession{
				ID:           sessionID,
				UserID:       userID,
				DeviceID:     "test-device-" + sessionID.Hex(),
				DeviceName:   "Test Device",
				Platform:     "test",
				RefreshToken: "refresh-token-" + sessionID.Hex(),
				IsActive:     true,
			}
			err := repo.CreateDeviceSession(context.Background(), testSession)
			require.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "All other device sessions revoked successfully", response["message"])
		assert.NotZero(t, response["revoked_count"])
	})

	t.Run("Revoke Session - Unauthorized", func(t *testing.T) {
		// Create session for different user
		otherUserID := primitive.NewObjectID()
		otherUser := &models.User{
			ID:        otherUserID,
			Email:     "other@example.com",
			Password:  "password",
			FirstName: "Other",
			LastName:  "User",
			Role:      models.CustomerRole,
		}
		err := repo.CreateUser(context.Background(), otherUser)
		require.NoError(t, err)

		otherSessionID := primitive.NewObjectID()
		otherSession := &models.DeviceSession{
			ID:           otherSessionID,
			UserID:       otherUserID, // Different user
			DeviceID:     "other-device",
			RefreshToken: "other-refresh-token",
			IsActive:     true,
		}
		err = repo.CreateDeviceSession(context.Background(), otherSession)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/"+otherSessionID.Hex(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Forbidden", response["error"])
	})
}

// TestTokenRotation tests the token rotation functionality
func TestTokenRotation(t *testing.T) {
	app, _, repo := setupTestApp(t)

	// Create user and session
	userID := primitive.NewObjectID()
	user := &models.User{
		ID:        userID,
		Email:     "rotation@example.com",
		Password:  "password",
		FirstName: "Rotation",
		LastName:  "Test",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	sessionID := primitive.NewObjectID()
	refreshToken := "initial-refresh-token"
	session := &models.DeviceSession{
		ID:           sessionID,
		UserID:       userID,
		DeviceID:     "rotation-device",
		DeviceName:   "Test Device",
		Platform:     "test",
		RefreshToken: refreshToken,
		IsActive:     true,
	}
	err = repo.CreateDeviceSession(context.Background(), session)
	require.NoError(t, err)

	// Create device session handler with JWT service
	accessSecret := make([]byte, 32)
	refreshSecret := make([]byte, 32)
	for i := range accessSecret {
		accessSecret[i] = byte(i + 1)
	}
	for i := range refreshSecret {
		refreshSecret[i] = byte(i + 32)
	}

	jwtService := createTestJWTService(accessSecret, refreshSecret)
	logger, _ := logger.New("debug", "console", "stdout")
	deviceHandler := handlers.NewDeviceSessionHandler(repo, jwtService, nil, logger.Logger)

	// Setup route
	app.Post("/api/v1/auth/rotate-token", deviceHandler.RefreshTokenWithRotation)

	t.Run("Successful Token Rotation", func(t *testing.T) {
		rotationReq := models.TokenRotationRequest{
			RefreshToken: refreshToken,
			DeviceID:     "rotation-device",
		}

		body, _ := json.Marshal(rotationReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/rotate-token", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response["access_token"])
		assert.NotEmpty(t, response["refresh_token"])
		assert.Equal(t, "Bearer", response["token_type"])
		assert.Equal(t, "Token rotated successfully", response["message"])

		// Verify the refresh token was updated in the session
		updatedSession, err := repo.GetDeviceSession(context.Background(), sessionID.Hex())
		require.NoError(t, err)
		assert.NotEqual(t, refreshToken, updatedSession.RefreshToken)
		assert.Equal(t, response["refresh_token"], updatedSession.RefreshToken)
	})

	t.Run("Invalid Refresh Token", func(t *testing.T) {
		rotationReq := models.TokenRotationRequest{
			RefreshToken: "invalid-refresh-token",
			DeviceID:     "rotation-device",
		}

		body, _ := json.Marshal(rotationReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/rotate-token", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Invalid refresh token", response["error"])
	})

	t.Run("Device ID Mismatch", func(t *testing.T) {
		// Get current refresh token from session
		currentSession, err := repo.GetDeviceSession(context.Background(), sessionID.Hex())
		require.NoError(t, err)

		rotationReq := models.TokenRotationRequest{
			RefreshToken: currentSession.RefreshToken,
			DeviceID:     "wrong-device-id", // Wrong device ID
		}

		body, _ := json.Marshal(rotationReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/rotate-token", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Device mismatch", response["error"])
	})
}
