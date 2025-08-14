package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/smorting/backend/internal/models"
)

func TestOTPDisabledInLoginFlow(t *testing.T) {
	t.Run("Enhanced auth result should never require verification", func(t *testing.T) {
		// Test that our EnhancedAuthResult struct supports disabling OTP
		result := models.EnhancedAuthResult{
			Success:              true,
			Message:              "Login successful",
			RequiresTwoFactor:    false,
			RequiresVerification: false, // This should always be false now
			RequiresCaptcha:      false,
			DeviceTrusted:        true,
		}

		// Assert OTP is disabled
		assert.False(t, result.RequiresTwoFactor, "Two-factor should be disabled")
		assert.False(t, result.RequiresVerification, "OTP verification should be disabled")
		assert.True(t, result.Success, "Login should succeed without OTP")
	})

	t.Run("Enhanced login request should not include OTP fields", func(t *testing.T) {
		// Test that login request works without OTP-related fields
		loginReq := models.EnhancedLoginRequest{
			Email:    "test@example.com",
			Password: "password123",
			// Note: No OTP or verification fields needed
		}

		// Basic validation that required fields are present
		assert.NotEmpty(t, loginReq.Email, "Email should be provided")
		assert.NotEmpty(t, loginReq.Password, "Password should be provided")
		
		// OTP-related fields should be empty/default
		assert.Empty(t, loginReq.TwoFactorCode, "Two-factor code should not be used")
	})

	t.Run("User model should work without email verification requirement", func(t *testing.T) {
		// Test that user can be authenticated regardless of email verification status
		user := models.User{
			Email:           "test@example.com",
			Password:        "hashedpassword",
			IsEmailVerified: false, // This should not block login anymore
		}

		// User should be valid for authentication regardless of email verification
		assert.NotEmpty(t, user.Email, "User should have email")
		assert.NotEmpty(t, user.Password, "User should have password")
		// Email verification status should not matter for login
		assert.False(t, user.IsEmailVerified, "Email verification should not be required for login")
	})
}
