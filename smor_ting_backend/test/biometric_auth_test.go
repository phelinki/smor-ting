package test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestBiometricAuthentication tests the complete biometric authentication flow
func TestBiometricAuthentication(t *testing.T) {
	app, _, repo := setupTestApp(t)

	// Create test user and device session
	userID := primitive.NewObjectID()
	user := &models.User{
		ID:        userID,
		Email:     "biometric@example.com",
		Password:  "hashed_password",
		FirstName: "Bio",
		LastName:  "Test",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	sessionID := primitive.NewObjectID()
	session := &models.DeviceSession{
		ID:           sessionID,
		UserID:       userID,
		DeviceID:     "bio-device-123",
		DeviceName:   "iPhone 14",
		Platform:     "ios",
		IPAddress:    "192.168.1.1",
		RefreshToken: "bio-refresh-token",
		IsActive:     true,
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
	}
	err = repo.CreateDeviceSession(context.Background(), session)
	require.NoError(t, err)

	// Setup biometric auth handler
	logger, _ := logger.New("debug", "console", "stdout")

	// Create encryption service
	encKey := make([]byte, 32)
	for i := range encKey {
		encKey[i] = byte(i + 64)
	}
	encryptionService, err := services.NewEncryptionService(encKey)
	require.NoError(t, err)

	// Create JWT service
	accessSecret := make([]byte, 32)
	refreshSecret := make([]byte, 32)
	for i := range accessSecret {
		accessSecret[i] = byte(i + 1)
	}
	for i := range refreshSecret {
		refreshSecret[i] = byte(i + 32)
	}
	jwtService := createTestJWTService(accessSecret, refreshSecret)

	// Create biometric secret
	biometricSecret := make([]byte, 32)
	rand.Read(biometricSecret)

	biometricHandler := handlers.NewBiometricAuthHandler(
		repo, jwtService, encryptionService, nil, logger.Logger, biometricSecret)

	// Setup routes with middleware simulation
	api := app.Group("/api/v1")
	bio := api.Group("/biometric")

	// Middleware to simulate authentication
	bio.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.Hex())
		c.Locals("sessionID", sessionID.Hex())
		return c.Next()
	})

	bio.Post("/enable", biometricHandler.EnableBiometricAuth)
	bio.Post("/disable", biometricHandler.DisableBiometricAuth)
	bio.Get("/challenge", biometricHandler.GenerateBiometricChallenge)
	bio.Post("/unlock", biometricHandler.UnlockWithBiometric)
	bio.Get("/config", biometricHandler.GetSecureStorageConfig)

	t.Run("Generate Biometric Challenge", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/biometric/challenge?device_id=bio-device-123", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response["challenge"])
		assert.NotEmpty(t, response["signature"])
		assert.NotEmpty(t, response["expires_at"])
		assert.Equal(t, "bio-device-123", response["device_id"])
	})

	t.Run("Enable Biometric Authentication", func(t *testing.T) {
		enableReq := map[string]interface{}{
			"biometric_type": "face_id",
			"device_id":      "bio-device-123",
			"challenge":      "test-challenge",
			"signature":      "test-signature",
		}

		body, _ := json.Marshal(enableReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/biometric/enable", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Biometric authentication enabled successfully", response["message"])
		assert.Equal(t, "face_id", response["biometric_type"])
		assert.Equal(t, "bio-device-123", response["device_id"])

		// Verify session was updated
		updatedSession, err := repo.GetDeviceSession(context.Background(), sessionID.Hex())
		require.NoError(t, err)
		assert.True(t, updatedSession.BiometricEnabled)
		assert.Equal(t, "face_id", updatedSession.BiometricType)
	})

	t.Run("Get Secure Storage Config", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/biometric/config?device_id=bio-device-123&platform=ios", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var config models.SecureTokenStorage
		err = json.NewDecoder(resp.Body).Decode(&config)
		require.NoError(t, err)

		assert.Equal(t, "ios", config.Platform)
		assert.True(t, config.BiometricEnabled)
		assert.Equal(t, "face_id", config.BiometricType)
		assert.True(t, config.KeychainEnabled)
		assert.False(t, config.KeystoreEnabled)
		assert.Equal(t, "hardware", config.EncryptionLevel)
	})

	t.Run("Biometric Unlock Flow", func(t *testing.T) {
		// First, generate a valid challenge
		challengeReq := httptest.NewRequest(http.MethodGet, "/api/v1/biometric/challenge?device_id=bio-device-123", nil)
		challengeResp, err := app.Test(challengeReq)
		require.NoError(t, err)

		var challengeData map[string]interface{}
		err = json.NewDecoder(challengeResp.Body).Decode(&challengeData)
		require.NoError(t, err)
		challengeResp.Body.Close()

		challenge := challengeData["challenge"].(string)
		signature := challengeData["signature"].(string)

		// Create unlock request
		unlockReq := models.BiometricUnlockRequest{
			DeviceID:          "bio-device-123",
			BiometricType:     "face_id",
			BiometricData:     "encrypted-biometric-data",
			Challenge:         challenge,
			ChallengeResponse: signature,
		}

		body, _ := json.Marshal(unlockReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/biometric/unlock", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Biometric unlock successful", response["message"])
		assert.NotEmpty(t, response["access_token"])
		assert.NotEmpty(t, response["refresh_token"])
		assert.Equal(t, "Bearer", response["token_type"])
		assert.NotEmpty(t, response["user"])
		assert.NotEmpty(t, response["secure_storage"])

		// Verify secure storage config in response
		secureStorage := response["secure_storage"].(map[string]interface{})
		assert.Equal(t, "ios", secureStorage["platform"])
		assert.Equal(t, true, secureStorage["biometric_enabled"])
		assert.Equal(t, "face_id", secureStorage["biometric_type"])
	})

	t.Run("Disable Biometric Authentication", func(t *testing.T) {
		disableReq := map[string]interface{}{
			"device_id": "bio-device-123",
			"password":  "password123", // In production, this would be validated
		}

		body, _ := json.Marshal(disableReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/biometric/disable", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Biometric authentication disabled successfully", response["message"])

		// Verify session was updated
		updatedSession, err := repo.GetDeviceSession(context.Background(), sessionID.Hex())
		require.NoError(t, err)
		assert.False(t, updatedSession.BiometricEnabled)
		assert.Equal(t, "", updatedSession.BiometricType)
	})
}

// TestBiometricValidation tests validation scenarios
func TestBiometricValidation(t *testing.T) {
	app, _, repo := setupTestApp(t)

	logger, _ := logger.New("debug", "console", "stdout")
	biometricSecret := make([]byte, 32)
	rand.Read(biometricSecret)

	biometricHandler := handlers.NewBiometricAuthHandler(
		repo, nil, nil, nil, logger.Logger, biometricSecret)

	api := app.Group("/api/v1")
	bio := api.Group("/biometric")

	bio.Get("/challenge", biometricHandler.GenerateBiometricChallenge)
	bio.Post("/unlock", biometricHandler.UnlockWithBiometric)

	t.Run("Challenge Generation - Missing Device ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/biometric/challenge", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Validation failed", response["error"])
		assert.Equal(t, "Device ID is required", response["message"])
	})

	t.Run("Biometric Unlock - Invalid Challenge", func(t *testing.T) {
		unlockReq := models.BiometricUnlockRequest{
			DeviceID:          "test-device",
			BiometricType:     "face_id",
			BiometricData:     "test-data",
			Challenge:         "invalid-challenge",
			ChallengeResponse: "invalid-response",
		}

		body, _ := json.Marshal(unlockReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/biometric/unlock", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Invalid challenge", response["error"])
		assert.Equal(t, "Challenge verification failed", response["message"])
	})

	t.Run("Biometric Unlock - Valid Challenge But No Session", func(t *testing.T) {
		// Generate valid challenge
		deviceID := "non-existent-device"
		challenge := "test-challenge"

		// Create valid signature
		hmacHash := hmac.New(sha256.New, biometricSecret)
		hmacHash.Write([]byte(deviceID + ":" + challenge))
		signature := base64.URLEncoding.EncodeToString(hmacHash.Sum(nil))

		unlockReq := models.BiometricUnlockRequest{
			DeviceID:          deviceID,
			BiometricType:     "face_id",
			BiometricData:     "test-data",
			Challenge:         challenge,
			ChallengeResponse: signature,
		}

		body, _ := json.Marshal(unlockReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/biometric/unlock", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Session not found", response["error"])
		assert.Contains(t, response["message"], "No biometric-enabled session found")
	})
}
