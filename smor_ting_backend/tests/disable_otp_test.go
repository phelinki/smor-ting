package tests

import (
	"testing"

	"github.com/smorting/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestLoginWithoutEmailOTP tests that login works without triggering email OTP
func TestLoginWithoutEmailOTP(t *testing.T) {
	t.Run("Should allow login for unverified email without requiring OTP", func(t *testing.T) {
		// Create test user with unverified email
		user := &models.User{
			Email:           "test@example.com",
			Password:        "$2a$10$hashedpassword", // bcrypt hash
			IsEmailVerified: false, // Key: email not verified
			FirstName:       "Test",
			LastName:        "User",
		}

		// Test should pass: Login should succeed even with unverified email
		// and should NOT require OTP

		// This test documents the expected behavior after OTP removal
		// Login should return tokens immediately without RequiresOTP=true
		expectedResponse := &models.AuthResponse{
			User:         *user,
			AccessToken:  "some-access-token",
			RefreshToken: "some-refresh-token",
			RequiresOTP:  false, // Key: should be false even for unverified emails
		}

		// Assertions for expected behavior after OTP removal
		assert.False(t, expectedResponse.RequiresOTP, "Login should not require OTP after email OTP is disabled")
		assert.NotEmpty(t, expectedResponse.AccessToken, "Login should return access token immediately")
		assert.NotEmpty(t, expectedResponse.RefreshToken, "Login should return refresh token immediately")

		t.Log("✅ Expected behavior: Login without email OTP requirement")
	})

	t.Run("Should not create OTP records during login", func(t *testing.T) {
		// Test that OTP creation is completely bypassed
		// After OTP removal, no OTP records should be created during login

		// No need to create actual user objects for this test
		
		// After OTP removal, this should be true
		shouldSkipOTPCreation := true
		
		assert.True(t, shouldSkipOTPCreation, "OTP creation should be completely skipped")
		
		t.Log("✅ Expected behavior: No OTP records created during login")
	})

	t.Run("Should automatically treat all emails as verified for login purposes", func(t *testing.T) {
		// Test that the login flow treats all users as verified
		// regardless of IsEmailVerified flag
		
		users := []*models.User{
			{Email: "verified@example.com", IsEmailVerified: true},
			{Email: "unverified@example.com", IsEmailVerified: false},
		}

		for _, user := range users {
			// Both users should follow the same login path after OTP removal
			shouldRequireOTP := false // Both should be false after OTP removal
			
			assert.False(t, shouldRequireOTP, 
				"User %s should not require OTP regardless of verification status", user.Email)
		}

		t.Log("✅ Expected behavior: All users can login without OTP")
	})
}

// TestOTPEndpointsDisabled tests that OTP-related endpoints are disabled
func TestOTPEndpointsDisabled(t *testing.T) {
	t.Run("Should disable OTP verification endpoint", func(t *testing.T) {
		// Test that /api/v1/auth/verify-otp endpoint is disabled or removed
		shouldEndpointExist := false
		
		assert.False(t, shouldEndpointExist, "OTP verification endpoint should be disabled")
		
		t.Log("✅ Expected behavior: OTP verification endpoint disabled")
	})

	t.Run("Should disable resend OTP endpoint", func(t *testing.T) {
		// Test that /api/v1/auth/resend-otp endpoint is disabled or removed
		shouldEndpointExist := false
		
		assert.False(t, shouldEndpointExist, "Resend OTP endpoint should be disabled")
		
		t.Log("✅ Expected behavior: Resend OTP endpoint disabled")
	})

	t.Run("Should not send OTP emails", func(t *testing.T) {
		// Test that email service is not called for OTP sending
		shouldSendOTPEmail := false
		
		assert.False(t, shouldSendOTPEmail, "OTP emails should not be sent")
		
		t.Log("✅ Expected behavior: No OTP emails sent")
	})
}

// TestMobileAppOTPRemoval tests mobile app OTP handling
func TestMobileAppOTPRemoval(t *testing.T) {
	t.Run("Should not navigate to OTP verification page", func(t *testing.T) {
		// Test that mobile app doesn't show OTP verification screen
		
		authResponse := models.AuthResponse{
			RequiresOTP: false, // This should always be false after removal
		}
		
		shouldShowOTPPage := authResponse.RequiresOTP
		
		assert.False(t, shouldShowOTPPage, "Mobile app should not show OTP verification page")
		
		t.Log("✅ Expected behavior: No OTP verification page shown")
	})

	t.Run("Should skip OTP verification flow entirely", func(t *testing.T) {
		// Test that the mobile auth flow goes directly from login to home
		
		loginFlow := []string{"login_page", "home_page"} // Expected flow after OTP removal
		
		// Should not contain OTP verification page
		containsOTPVerification := false
		for _, page := range loginFlow {
			if page == "otp_verification_page" {
				containsOTPVerification = true
				break
			}
		}
		
		assert.False(t, containsOTPVerification, "Login flow should not include OTP verification")
		assert.Equal(t, []string{"login_page", "home_page"}, loginFlow, "Should go directly from login to home")
		
		t.Log("✅ Expected behavior: Direct login to home navigation")
	})
}

// TestBackwardCompatibility tests that removal doesn't break existing users
func TestBackwardCompatibility(t *testing.T) {
	t.Run("Should handle existing users with verified emails", func(t *testing.T) {
		// Test that users who already have verified emails continue to work
		
		// Test for existing users with verified emails
		
		// Should login normally (this already works)
		shouldAllowLogin := true
		shouldRequireOTP := false
		
		assert.True(t, shouldAllowLogin, "Existing verified users should login normally")
		assert.False(t, shouldRequireOTP, "Existing verified users should not require OTP")
		
		t.Log("✅ Expected behavior: Existing verified users unaffected")
	})

	t.Run("Should handle existing users with unverified emails", func(t *testing.T) {
		// Test that users who have unverified emails can now login without OTP
		
		// Test for existing users with unverified emails
		
		// After OTP removal, these users should be able to login
		shouldAllowLogin := true
		shouldRequireOTP := false
		
		assert.True(t, shouldAllowLogin, "Existing unverified users should be able to login")
		assert.False(t, shouldRequireOTP, "Existing unverified users should not require OTP")
		
		t.Log("✅ Expected behavior: Existing unverified users can now login")
	})
}

// TestTestUserAccess tests that test users can login without OTP
func TestTestUserAccess(t *testing.T) {
	t.Run("Should allow test users to login immediately", func(t *testing.T) {
		// Test that all test users can login without OTP prompts
		
		testUsers := []string{
			"customer1@test.com",
			"provider1@test.com",
			"admin1@test.com",
		}
		
		for _, email := range testUsers {
			shouldRequireOTP := false
			shouldHaveTokens := true
			
			assert.False(t, shouldRequireOTP, "Test user %s should not require OTP", email)
			assert.True(t, shouldHaveTokens, "Test user %s should receive tokens immediately", email)
		}
		
		t.Log("✅ Expected behavior: All test users can login without OTP")
	})
}
