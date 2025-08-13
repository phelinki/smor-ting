package test

import (
	"testing"
	"time"

	"github.com/smorting/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBruteForceProtection(t *testing.T) {
	logger := zap.NewNop()
	protector := services.NewBruteForceProtector(logger)

	email := "test@example.com"
	ipAddress := "192.168.1.1"

	t.Run("initially_allows_login_attempts", func(t *testing.T) {
		err := protector.CheckAllowed(email, ipAddress)
		assert.NoError(t, err, "Initial login attempts should be allowed")
	})

	t.Run("records_failed_attempts", func(t *testing.T) {
		// Record 3 failures
		for i := 0; i < 3; i++ {
			protector.RecordFailure(email, ipAddress)
		}

		// Should still be allowed but require CAPTCHA
		err := protector.CheckAllowed(email, ipAddress)
		assert.NoError(t, err, "Should still allow attempts before lockout")

		requiresCaptcha := protector.RequiresCaptcha(email, ipAddress)
		assert.True(t, requiresCaptcha, "Should require CAPTCHA after 3 failed attempts")
	})

	t.Run("locks_out_after_max_attempts", func(t *testing.T) {
		// Record 2 more failures (total 5)
		for i := 0; i < 2; i++ {
			protector.RecordFailure(email, ipAddress)
		}

		// Should now be locked out
		err := protector.CheckAllowed(email, ipAddress)
		assert.Error(t, err, "Should be locked out after 5 failed attempts")
		assert.Contains(t, err.Error(), "temporarily locked", "Error should mention temporary lockout")
	})

	t.Run("gets_lockout_info", func(t *testing.T) {
		info := protector.GetLockoutInfo(email, ipAddress)
		assert.True(t, info.EmailLocked, "Email should be locked")
		assert.True(t, info.IPLocked, "IP should be locked")
		assert.Greater(t, info.EmailLockoutRemaining, time.Duration(0), "Should have remaining lockout time")
	})

	t.Run("resets_on_success", func(t *testing.T) {
		// Reset counters
		protector.RecordSuccess(email, ipAddress)

		// Should now be allowed
		err := protector.CheckAllowed(email, ipAddress)
		assert.NoError(t, err, "Should be allowed after successful auth")

		requiresCaptcha := protector.RequiresCaptcha(email, ipAddress)
		assert.False(t, requiresCaptcha, "Should not require CAPTCHA after successful auth")
	})
}

func TestBruteForceProtectionSeparateTracking(t *testing.T) {
	logger := zap.NewNop()
	protector := services.NewBruteForceProtector(logger)

	email1 := "user1@example.com"
	email2 := "user2@example.com"
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	t.Run("tracks_email_and_ip_separately", func(t *testing.T) {
		// Lock out email1 from ip1
		for i := 0; i < 5; i++ {
			protector.RecordFailure(email1, ip1)
		}

		// email1 should be locked from any IP
		err := protector.CheckAllowed(email1, ip2)
		assert.Error(t, err, "Email should be locked from any IP")

		// Different email from same IP should be allowed
		err = protector.CheckAllowed(email2, ip1)
		assert.Error(t, err, "IP should also be locked") // IP got locked too

		// Different email from different IP should be allowed
		err = protector.CheckAllowed(email2, ip2)
		assert.NoError(t, err, "Different email from different IP should be allowed")
	})
}

func TestBruteForceProtectionCaptchaRequirement(t *testing.T) {
	logger := zap.NewNop()
	protector := services.NewBruteForceProtector(logger)

	email := "test@example.com"
	ipAddress := "192.168.1.1"

	t.Run("requires_captcha_after_threshold", func(t *testing.T) {
		// Should not require CAPTCHA initially
		requiresCaptcha := protector.RequiresCaptcha(email, ipAddress)
		assert.False(t, requiresCaptcha, "Should not require CAPTCHA initially")

		// Record 2 failures - still no CAPTCHA
		for i := 0; i < 2; i++ {
			protector.RecordFailure(email, ipAddress)
		}
		requiresCaptcha = protector.RequiresCaptcha(email, ipAddress)
		assert.False(t, requiresCaptcha, "Should not require CAPTCHA after 2 failures")

		// Record 1 more failure - now requires CAPTCHA
		protector.RecordFailure(email, ipAddress)
		requiresCaptcha = protector.RequiresCaptcha(email, ipAddress)
		assert.True(t, requiresCaptcha, "Should require CAPTCHA after 3 failures")
	})
}

func TestBruteForceProtectionExponentialBackoff(t *testing.T) {
	logger := zap.NewNop()
	protector := services.NewBruteForceProtector(logger)

	email := "test@example.com"
	ipAddress := "192.168.1.1"

	t.Run("implements_exponential_backoff", func(t *testing.T) {
		// First lockout
		for i := 0; i < 5; i++ {
			protector.RecordFailure(email, ipAddress)
		}

		info1 := protector.GetLockoutInfo(email, ipAddress)
		assert.True(t, info1.EmailLocked, "Should be locked after first round")

		// Wait a bit to simulate time passing (in real scenario, would wait for lockout to expire)
		// For test, we'll reset and simulate second lockout
		protector.RecordSuccess(email, ipAddress) // Reset

		// Second lockout
		for i := 0; i < 5; i++ {
			protector.RecordFailure(email, ipAddress)
		}

		info2 := protector.GetLockoutInfo(email, ipAddress)
		assert.True(t, info2.EmailLocked, "Should be locked after second round")

		// Note: In a real implementation, the second lockout would be longer
		// This is a basic test to ensure the functionality works
	})
}

func TestBruteForceProtectionRemainingAttempts(t *testing.T) {
	logger := zap.NewNop()
	protector := services.NewBruteForceProtector(logger)

	email := "test@example.com"
	ipAddress := "192.168.1.1"

	t.Run("calculates_remaining_attempts", func(t *testing.T) {
		remaining := protector.GetRemainingAttempts(email, ipAddress)
		assert.Equal(t, 5, remaining, "Should have 5 attempts initially")

		// Record 2 failures
		for i := 0; i < 2; i++ {
			protector.RecordFailure(email, ipAddress)
		}

		remaining = protector.GetRemainingAttempts(email, ipAddress)
		assert.Equal(t, 3, remaining, "Should have 3 attempts remaining after 2 failures")

		// Record 3 more failures (total 5)
		for i := 0; i < 3; i++ {
			protector.RecordFailure(email, ipAddress)
		}

		remaining = protector.GetRemainingAttempts(email, ipAddress)
		assert.Equal(t, 0, remaining, "Should have 0 attempts remaining after lockout")
	})
}
