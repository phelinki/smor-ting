package test

import (
	"context"
	"testing"
	"time"

	"github.com/smorting/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestDeviceSessionModel tests the DeviceSession model functionality
func TestDeviceSessionModel(t *testing.T) {
	t.Run("Create Device Session", func(t *testing.T) {
		userID := primitive.NewObjectID()

		session := models.DeviceSession{
			ID:           primitive.NewObjectID(),
			UserID:       userID,
			DeviceID:     "device-123",
			DeviceName:   "iPhone 12",
			Platform:     "ios",
			IPAddress:    "192.168.1.1",
			UserAgent:    "SmorTing/1.0 iOS",
			RefreshToken: "refresh-token-123",
			IsActive:     true,
			LastActivity: time.Now(),
			CreatedAt:    time.Now(),
		}

		assert.NotEmpty(t, session.ID)
		assert.Equal(t, userID, session.UserID)
		assert.Equal(t, "device-123", session.DeviceID)
		assert.Equal(t, "iPhone 12", session.DeviceName)
		assert.Equal(t, "ios", session.Platform)
		assert.True(t, session.IsActive)
	})

	t.Run("Revoke Device Session", func(t *testing.T) {
		session := models.DeviceSession{
			ID:           primitive.NewObjectID(),
			UserID:       primitive.NewObjectID(),
			DeviceID:     "device-123",
			IsActive:     true,
			LastActivity: time.Now(),
		}

		// Revoke the session
		session.RevokeSession()

		assert.False(t, session.IsActive)
		assert.NotZero(t, session.RevokedAt)
	})

	t.Run("Check Session Expiry", func(t *testing.T) {
		session := models.DeviceSession{
			LastActivity: time.Now().Add(-8 * 24 * time.Hour), // 8 days ago
		}

		expired := session.IsExpired(7 * 24 * time.Hour) // 7-day expiry
		assert.True(t, expired)

		session.LastActivity = time.Now().Add(-6 * 24 * time.Hour) // 6 days ago
		expired = session.IsExpired(7 * 24 * time.Hour)
		assert.False(t, expired)
	})
}

// TestTokenRevocationService tests token revocation functionality
func TestTokenRevocationService(t *testing.T) {
	_, _, repo := setupTestApp(t)

	t.Run("Revoke All Tokens on Password Change", func(t *testing.T) {
		// Create a user first
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

		// Create multiple device sessions
		sessions := []models.DeviceSession{
			{
				ID:           primitive.NewObjectID(),
				UserID:       userID,
				DeviceID:     "device-1",
				DeviceName:   "iPhone",
				RefreshToken: "token-1",
				IsActive:     true,
			},
			{
				ID:           primitive.NewObjectID(),
				UserID:       userID,
				DeviceID:     "device-2",
				DeviceName:   "Android",
				RefreshToken: "token-2",
				IsActive:     true,
			},
		}

		for _, session := range sessions {
			err := repo.CreateDeviceSession(context.Background(), &session)
			require.NoError(t, err)
		}

		// Revoke all tokens for the user
		err = repo.RevokeAllUserTokens(context.Background(), userID.Hex())
		require.NoError(t, err)

		// Verify all sessions are revoked
		userSessions, err := repo.GetUserDeviceSessions(context.Background(), userID.Hex())
		require.NoError(t, err)

		for _, session := range userSessions {
			assert.False(t, session.IsActive)
			assert.NotZero(t, session.RevokedAt)
		}
	})

	t.Run("Revoke Specific Device Session", func(t *testing.T) {
		userID := primitive.NewObjectID()
		sessionID := primitive.NewObjectID()

		session := models.DeviceSession{
			ID:           sessionID,
			UserID:       userID,
			DeviceID:     "device-specific",
			RefreshToken: "token-specific",
			IsActive:     true,
		}

		err := repo.CreateDeviceSession(context.Background(), &session)
		require.NoError(t, err)

		// Revoke specific session
		err = repo.RevokeDeviceSession(context.Background(), sessionID.Hex())
		require.NoError(t, err)

		// Verify session is revoked
		revokedSession, err := repo.GetDeviceSession(context.Background(), sessionID.Hex())
		require.NoError(t, err)
		assert.False(t, revokedSession.IsActive)
	})

	t.Run("Token Rotation on Refresh", func(t *testing.T) {
		userID := primitive.NewObjectID()
		sessionID := primitive.NewObjectID()

		oldRefreshToken := "old-refresh-token"
		newRefreshToken := "new-refresh-token"

		session := models.DeviceSession{
			ID:           sessionID,
			UserID:       userID,
			DeviceID:     "device-rotation",
			RefreshToken: oldRefreshToken,
			IsActive:     true,
		}

		err := repo.CreateDeviceSession(context.Background(), &session)
		require.NoError(t, err)

		// Rotate the refresh token
		err = repo.RotateRefreshToken(context.Background(), sessionID.Hex(), newRefreshToken)
		require.NoError(t, err)

		// Verify token was rotated
		updatedSession, err := repo.GetDeviceSession(context.Background(), sessionID.Hex())
		require.NoError(t, err)
		assert.Equal(t, newRefreshToken, updatedSession.RefreshToken)
		assert.NotEqual(t, oldRefreshToken, updatedSession.RefreshToken)
	})
}

// TestSecurityEventLogging tests security event logging for audit trails
func TestSecurityEventLogging(t *testing.T) {
	_, _, repo := setupTestApp(t)

	t.Run("Log Password Change Event", func(t *testing.T) {
		userID := primitive.NewObjectID()

		event := models.SecurityEvent{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			EventType: models.PasswordChangeEvent,
			IPAddress: "192.168.1.1",
			UserAgent: "SmorTing/1.0",
			DeviceID:  "device-123",
			Metadata: map[string]interface{}{
				"reason": "user_requested",
				"method": "password_reset",
			},
			Timestamp: time.Now(),
		}

		err := repo.LogSecurityEvent(context.Background(), &event)
		require.NoError(t, err)

		// Verify event was logged
		events, err := repo.GetUserSecurityEvents(context.Background(), userID.Hex(), 10)
		require.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, models.PasswordChangeEvent, events[0].EventType)
	})

	t.Run("Log 2FA Change Event", func(t *testing.T) {
		userID := primitive.NewObjectID()

		event := models.SecurityEvent{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			EventType: models.TwoFactorChangeEvent,
			IPAddress: "192.168.1.1",
			Metadata: map[string]interface{}{
				"action": "enabled",
				"method": "totp",
			},
			Timestamp: time.Now(),
		}

		err := repo.LogSecurityEvent(context.Background(), &event)
		require.NoError(t, err)

		events, err := repo.GetUserSecurityEvents(context.Background(), userID.Hex(), 10)
		require.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, models.TwoFactorChangeEvent, events[0].EventType)
	})
}
