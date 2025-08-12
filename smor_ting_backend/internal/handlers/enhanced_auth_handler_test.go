package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smorting/backend/internal/models"
)

func TestRefreshTokenRequest_Structure(t *testing.T) {
	// Test that RefreshTokenRequest includes the SessionID field
	// This ensures we fixed the duplicate declaration and missing field issues
	req := RefreshTokenRequest{
		RefreshToken: "refresh_token_123",
		SessionID:    "session_123",
	}

	assert.Equal(t, "refresh_token_123", req.RefreshToken)
	assert.Equal(t, "session_123", req.SessionID)
}

func TestUserModel_PasswordField(t *testing.T) {
	// Test that User model has the Password field (not PasswordHash)
	// This ensures we're using the correct field name in the handler
	user := &models.User{
		Email:    "test@example.com",
		Password: "$2a$10$hashedpassword",
	}

	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "$2a$10$hashedpassword", user.Password)
	
	// Verify the field exists by attempting to access it
	assert.NotEmpty(t, user.Password)
}

func TestEnhancedAuthModels_Exist(t *testing.T) {
	// Test that the new enhanced auth models are properly defined
	
	// Test EnhancedLoginRequest
	loginReq := models.EnhancedLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	assert.Equal(t, "test@example.com", loginReq.Email)
	assert.Equal(t, "password123", loginReq.Password)
	
	// Test EnhancedAuthResult
	authResult := models.EnhancedAuthResult{
		Success:              true,
		Message:              "Login successful",
		RequiresTwoFactor:    false,
		RequiresVerification: false,
		RequiresCaptcha:      false,
		DeviceTrusted:        true,
	}
	assert.True(t, authResult.Success)
	assert.Equal(t, "Login successful", authResult.Message)
	assert.False(t, authResult.RequiresTwoFactor)
	assert.False(t, authResult.RequiresVerification)
	assert.False(t, authResult.RequiresCaptcha)
	assert.True(t, authResult.DeviceTrusted)
	
	// Test SessionInfo
	sessionInfo := models.SessionInfo{
		SessionID:  "session_123",
		UserID:     "user_123",
		DeviceName: "iPhone",
		DeviceType: "mobile",
		IPAddress:  "192.168.1.1",
		IsCurrent:  true,
		IsRevoked:  false,
	}
	assert.Equal(t, "session_123", sessionInfo.SessionID)
	assert.Equal(t, "user_123", sessionInfo.UserID)
	assert.Equal(t, "iPhone", sessionInfo.DeviceName)
	assert.Equal(t, "mobile", sessionInfo.DeviceType)
	assert.Equal(t, "192.168.1.1", sessionInfo.IPAddress)
	assert.True(t, sessionInfo.IsCurrent)
	assert.False(t, sessionInfo.IsRevoked)
	
	// Test LockoutInfo
	lockoutInfo := models.LockoutInfo{
		RemainingAttempts: 3,
		LockoutReason:     "Too many failed attempts",
	}
	assert.Equal(t, 3, lockoutInfo.RemainingAttempts)
	assert.Equal(t, "Too many failed attempts", lockoutInfo.LockoutReason)
}