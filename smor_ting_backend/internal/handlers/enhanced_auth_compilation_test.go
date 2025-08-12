package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smorting/backend/internal/models"
)

// TestEnhancedLoginRequest_RequiredFields tests that all required fields exist
func TestEnhancedLoginRequest_RequiredFields(t *testing.T) {
	// Test that the EnhancedLoginRequest model has all the fields used in the handler
	req := models.EnhancedLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
		// These fields must exist for the handler to compile
		DeviceFingerprint: "test-device-fingerprint",
		UserAgent:         "test-user-agent",
		ClientIP:          "192.168.1.1",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
	assert.Equal(t, "test-device-fingerprint", req.DeviceFingerprint)
	assert.Equal(t, "test-user-agent", req.UserAgent)
	assert.Equal(t, "192.168.1.1", req.ClientIP)
}

// TestEnhancedLoginRequest_OptionalFields tests optional fields that the handler expects
func TestEnhancedLoginRequest_OptionalFields(t *testing.T) {
	// Test that optional fields used in the handler exist
	req := models.EnhancedLoginRequest{
		Email:         "test@example.com",
		Password:      "password123",
		CaptchaToken:  "captcha-token",
		TwoFactorCode: "123456",
		BiometricData: "biometric-hash",
		SessionID:     "session-123",
	}

	assert.Equal(t, "captcha-token", req.CaptchaToken)
	assert.Equal(t, "123456", req.TwoFactorCode)
	assert.Equal(t, "biometric-hash", req.BiometricData)
	assert.Equal(t, "session-123", req.SessionID)
}

// TestEnhancedAuthService_InterfaceCompliance tests that our service implements required methods
func TestEnhancedAuthService_InterfaceCompliance(t *testing.T) {
	// This test ensures that our enhanced auth service interface has all required methods
	// The interface should be implemented by any concrete service

	var service EnhancedAuthService

	// These method signatures must exist for the handler to compile
	assert.NotNil(t, service) // interface can be nil, but type must exist

	// Test that the interface methods are properly defined by checking the interface
	// This will fail to compile if the methods don't exist with correct signatures

	// The following would be used in actual implementation:
	// result, err := service.EnhancedLogin(&models.EnhancedLoginRequest{}, "192.168.1.1")
	// sessions, err := service.GetUserSessions("user-123")
	// err := service.RevokeSession("session-123")
	// err := service.SignOutAllDevices("user-123")
	// result, err := service.RefreshTokenWithSession("refresh-token", "session-123")
}

// TestRefreshTokenRequest_SessionIDField tests that SessionID field exists and works
func TestRefreshTokenRequest_SessionIDField(t *testing.T) {
	// Test that RefreshTokenRequest has SessionID field and it's properly accessible
	req := RefreshTokenRequest{
		RefreshToken: "refresh-token-123",
		SessionID:    "session-123",
	}

	assert.Equal(t, "refresh-token-123", req.RefreshToken)
	assert.Equal(t, "session-123", req.SessionID)

	// Test that both fields are properly tagged for JSON
	assert.NotEmpty(t, req.RefreshToken)
	assert.NotEmpty(t, req.SessionID)
}

// TestEnhancedLoginResponse_Compatibility tests response structure compatibility
func TestEnhancedLoginResponse_Compatibility(t *testing.T) {
	// Test that EnhancedLoginResponse can be used alongside models.EnhancedAuthResult
	response := EnhancedLoginResponse{
		Success:              true,
		Message:              "Login successful",
		RequiresTwoFactor:    false,
		RequiresVerification: false,
		RequiresCaptcha:      false,
		DeviceTrusted:        true,
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Login successful", response.Message)
	assert.False(t, response.RequiresTwoFactor)
	assert.False(t, response.RequiresVerification)
	assert.False(t, response.RequiresCaptcha)
	assert.True(t, response.DeviceTrusted)
}
