package test

import (
	"testing"

	"github.com/smorting/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBruteForceIntegrationScenarios(t *testing.T) {
	logger := zap.NewNop()
	protector := services.NewBruteForceProtector(logger)

	email := "user@example.com"
	ip := "192.168.1.100"

	t.Run("complete_brute_force_protection_workflow", func(t *testing.T) {
		// Step 1: Initial attempts should be allowed
		err := protector.CheckAllowed(email, ip)
		assert.NoError(t, err, "Initial attempt should be allowed")

		requiresCaptcha := protector.RequiresCaptcha(email, ip)
		assert.False(t, requiresCaptcha, "Should not require CAPTCHA initially")

		remaining := protector.GetRemainingAttempts(email, ip)
		assert.Equal(t, 5, remaining, "Should have 5 attempts initially")

		// Step 2: Record 2 failures - still allowed but getting closer
		for i := 0; i < 2; i++ {
			protector.RecordFailure(email, ip)
		}

		err = protector.CheckAllowed(email, ip)
		assert.NoError(t, err, "Should still be allowed after 2 failures")

		remaining = protector.GetRemainingAttempts(email, ip)
		assert.Equal(t, 3, remaining, "Should have 3 attempts remaining")

		requiresCaptcha = protector.RequiresCaptcha(email, ip)
		assert.False(t, requiresCaptcha, "Should not require CAPTCHA yet")

		// Step 3: One more failure triggers CAPTCHA requirement
		protector.RecordFailure(email, ip)

		err = protector.CheckAllowed(email, ip)
		assert.NoError(t, err, "Should still be allowed after 3 failures")

		requiresCaptcha = protector.RequiresCaptcha(email, ip)
		assert.True(t, requiresCaptcha, "Should require CAPTCHA after 3 failures")

		remaining = protector.GetRemainingAttempts(email, ip)
		assert.Equal(t, 2, remaining, "Should have 2 attempts remaining")

		// Step 4: Two more failures trigger lockout
		for i := 0; i < 2; i++ {
			protector.RecordFailure(email, ip)
		}

		err = protector.CheckAllowed(email, ip)
		assert.Error(t, err, "Should be locked out after 5 failures")
		assert.Contains(t, err.Error(), "locked", "Error message should mention lockout")

		remaining = protector.GetRemainingAttempts(email, ip)
		assert.Equal(t, 0, remaining, "Should have 0 attempts remaining during lockout")

		// Step 5: Get lockout info
		lockoutInfo := protector.GetLockoutInfo(email, ip)
		assert.True(t, lockoutInfo.EmailLocked, "Email should be locked")
		assert.True(t, lockoutInfo.IPLocked, "IP should be locked")
		assert.Greater(t, lockoutInfo.EmailLockoutRemaining.Seconds(), float64(0), "Should have lockout time remaining")

		// Step 6: Success resets everything
		protector.RecordSuccess(email, ip)

		err = protector.CheckAllowed(email, ip)
		assert.NoError(t, err, "Should be allowed after successful auth")

		requiresCaptcha = protector.RequiresCaptcha(email, ip)
		assert.False(t, requiresCaptcha, "Should not require CAPTCHA after success")

		remaining = protector.GetRemainingAttempts(email, ip)
		assert.Equal(t, 5, remaining, "Should have full attempts after success")
	})

	t.Run("separate_email_and_ip_tracking", func(t *testing.T) {
		email1 := "user1@example.com"
		email2 := "user2@example.com"
		ip1 := "192.168.1.101"
		ip2 := "192.168.1.102"

		// Lock user1 from ip1
		for i := 0; i < 5; i++ {
			protector.RecordFailure(email1, ip1)
		}

		// user1 should be locked from any IP
		err := protector.CheckAllowed(email1, ip2)
		assert.Error(t, err, "user1 should be locked from any IP")

		// ip1 should be locked for any user
		err = protector.CheckAllowed(email2, ip1)
		assert.Error(t, err, "ip1 should be locked for any user")

		// Different user from different IP should work
		err = protector.CheckAllowed(email2, ip2)
		assert.NoError(t, err, "Different user from different IP should work")
	})

	t.Run("captcha_thresholds", func(t *testing.T) {
		testEmail := "captcha-test@example.com"
		testIP := "192.168.1.103"

		// Test CAPTCHA threshold - only test up to 4 failures (before lockout)
		for i := 1; i <= 4; i++ {
			protector.RecordFailure(testEmail, testIP)

			requiresCaptcha := protector.RequiresCaptcha(testEmail, testIP)
			if i < 3 {
				assert.False(t, requiresCaptcha, "Should not require CAPTCHA before 3 failures")
			} else {
				assert.True(t, requiresCaptcha, "Should require CAPTCHA after 3+ failures")
			}
		}
	})
}

func TestBruteForceProtectionConfiguration(t *testing.T) {
	t.Run("default_configuration_values", func(t *testing.T) {
		config := services.DefaultBruteForceConfig()

		assert.Equal(t, 5, config.MaxAttempts, "Should have 5 max attempts")
		assert.True(t, config.ExponentialBackoff, "Should enable exponential backoff")
		assert.Greater(t, config.LockoutDuration.Minutes(), float64(10), "Should have reasonable lockout duration")
		assert.Greater(t, config.MaxLockoutDuration.Hours(), float64(12), "Should have reasonable max lockout")
	})
}

func TestBruteForceProtectionConcurrentAccess(t *testing.T) {
	logger := zap.NewNop()
	protector := services.NewBruteForceProtector(logger)

	email := "concurrent@example.com"
	ip := "192.168.1.104"

	t.Run("concurrent_failure_recording", func(t *testing.T) {
		// This test ensures the brute force protector is thread-safe
		// In a real scenario, multiple goroutines might record failures simultaneously

		done := make(chan bool, 3)

		// Start 3 goroutines that record failures
		for i := 0; i < 3; i++ {
			go func() {
				protector.RecordFailure(email, ip)
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		// Check that failures were recorded
		remaining := protector.GetRemainingAttempts(email, ip)
		assert.Equal(t, 2, remaining, "Should have 2 attempts remaining after 3 concurrent failures")

		requiresCaptcha := protector.RequiresCaptcha(email, ip)
		assert.True(t, requiresCaptcha, "Should require CAPTCHA after 3 failures")
	})
}
