package services

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// BruteForceProtector implements sophisticated brute force protection
type BruteForceProtector struct {
	emailAttempts map[string]*AttempTracker
	ipAttempts    map[string]*AttempTracker
	mu            sync.RWMutex
	logger        *zap.Logger
}

// AttempTracker tracks failed attempts and lockout state
type AttempTracker struct {
	FailedAttempts int
	FirstFailure   time.Time
	LastFailure    time.Time
	LockedUntil    time.Time
	TotalLockouts  int
}

// BruteForceConfig defines protection parameters
type BruteForceConfig struct {
	MaxAttempts        int           // Max attempts before lockout
	LockoutDuration    time.Duration // Base lockout duration
	WindowDuration     time.Duration // Time window for attempt counting
	MaxLockoutDuration time.Duration // Maximum lockout duration
	ExponentialBackoff bool          // Enable exponential backoff
}

// DefaultBruteForceConfig returns sensible defaults
func DefaultBruteForceConfig() *BruteForceConfig {
	return &BruteForceConfig{
		MaxAttempts:        5,
		LockoutDuration:    15 * time.Minute,
		WindowDuration:     1 * time.Hour,
		MaxLockoutDuration: 24 * time.Hour,
		ExponentialBackoff: true,
	}
}

// NewBruteForceProtector creates a new brute force protection service
func NewBruteForceProtector(logger *zap.Logger) *BruteForceProtector {
	protector := &BruteForceProtector{
		emailAttempts: make(map[string]*AttempTracker),
		ipAttempts:    make(map[string]*AttempTracker),
		logger:        logger,
	}

	// Start cleanup goroutine
	go protector.cleanupExpiredEntries()

	return protector
}

// CheckAllowed checks if authentication attempt is allowed
func (bp *BruteForceProtector) CheckAllowed(email, ipAddress string) error {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	now := time.Now()

	// Check email-based lockout
	if emailTracker, exists := bp.emailAttempts[email]; exists {
		if now.Before(emailTracker.LockedUntil) {
			remaining := emailTracker.LockedUntil.Sub(now)
			bp.logger.Warn("Email locked due to brute force",
				zap.String("email", email),
				zap.Duration("remaining", remaining),
				zap.Int("total_lockouts", emailTracker.TotalLockouts),
			)
			return fmt.Errorf("account temporarily locked. Try again in %v", remaining.Round(time.Minute))
		}
	}

	// Check IP-based lockout
	if ipTracker, exists := bp.ipAttempts[ipAddress]; exists {
		if now.Before(ipTracker.LockedUntil) {
			remaining := ipTracker.LockedUntil.Sub(now)
			bp.logger.Warn("IP locked due to brute force",
				zap.String("ip", ipAddress),
				zap.Duration("remaining", remaining),
				zap.Int("total_lockouts", ipTracker.TotalLockouts),
			)
			return fmt.Errorf("too many failed attempts from this location. Try again in %v", remaining.Round(time.Minute))
		}
	}

	return nil
}

// RecordFailure records a failed authentication attempt
func (bp *BruteForceProtector) RecordFailure(email, ipAddress string) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	config := DefaultBruteForceConfig()
	now := time.Now()

	// Record email failure
	bp.recordEmailFailure(email, now, config)

	// Record IP failure
	bp.recordIPFailure(ipAddress, now, config)

	bp.logger.Info("Recorded authentication failure",
		zap.String("email", email),
		zap.String("ip", ipAddress),
		zap.Time("timestamp", now),
	)
}

// RecordSuccess records a successful authentication (resets counters)
func (bp *BruteForceProtector) RecordSuccess(email, ipAddress string) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	// Reset email attempts
	if tracker, exists := bp.emailAttempts[email]; exists {
		tracker.FailedAttempts = 0
		tracker.LockedUntil = time.Time{}
	}

	// Reset IP attempts
	if tracker, exists := bp.ipAttempts[ipAddress]; exists {
		tracker.FailedAttempts = 0
		tracker.LockedUntil = time.Time{}
	}

	bp.logger.Info("Reset brute force counters after successful auth",
		zap.String("email", email),
		zap.String("ip", ipAddress),
	)
}

// GetLockoutInfo returns current lockout information
func (bp *BruteForceProtector) GetLockoutInfo(email, ipAddress string) *LockoutInfo {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	info := &LockoutInfo{}
	now := time.Now()

	if tracker, exists := bp.emailAttempts[email]; exists {
		info.EmailLocked = now.Before(tracker.LockedUntil)
		if info.EmailLocked {
			info.EmailLockoutRemaining = tracker.LockedUntil.Sub(now)
		}
		info.EmailAttempts = tracker.FailedAttempts
	}

	if tracker, exists := bp.ipAttempts[ipAddress]; exists {
		info.IPLocked = now.Before(tracker.LockedUntil)
		if info.IPLocked {
			info.IPLockoutRemaining = tracker.LockedUntil.Sub(now)
		}
		info.IPAttempts = tracker.FailedAttempts
	}

	return info
}

// LockoutInfo contains current lockout status
type LockoutInfo struct {
	EmailLocked           bool          `json:"email_locked"`
	EmailLockoutRemaining time.Duration `json:"email_lockout_remaining"`
	EmailAttempts         int           `json:"email_attempts"`
	IPLocked              bool          `json:"ip_locked"`
	IPLockoutRemaining    time.Duration `json:"ip_lockout_remaining"`
	IPAttempts            int           `json:"ip_attempts"`
}

// Private helper methods

func (bp *BruteForceProtector) recordEmailFailure(email string, now time.Time, config *BruteForceConfig) {
	tracker, exists := bp.emailAttempts[email]
	if !exists {
		tracker = &AttempTracker{
			FirstFailure: now,
		}
		bp.emailAttempts[email] = tracker
	}

	// Reset counter if outside window
	if now.Sub(tracker.FirstFailure) > config.WindowDuration {
		tracker.FailedAttempts = 0
		tracker.FirstFailure = now
	}

	tracker.FailedAttempts++
	tracker.LastFailure = now

	// Check if lockout needed
	if tracker.FailedAttempts >= config.MaxAttempts {
		lockoutDuration := bp.calculateLockoutDuration(tracker.TotalLockouts, config)
		tracker.LockedUntil = now.Add(lockoutDuration)
		tracker.TotalLockouts++
		tracker.FailedAttempts = 0 // Reset counter after lockout

		bp.logger.Warn("Email locked due to repeated failures",
			zap.String("email", email),
			zap.Duration("lockout_duration", lockoutDuration),
			zap.Int("total_lockouts", tracker.TotalLockouts),
		)
	}
}

func (bp *BruteForceProtector) recordIPFailure(ipAddress string, now time.Time, config *BruteForceConfig) {
	tracker, exists := bp.ipAttempts[ipAddress]
	if !exists {
		tracker = &AttempTracker{
			FirstFailure: now,
		}
		bp.ipAttempts[ipAddress] = tracker
	}

	// Reset counter if outside window
	if now.Sub(tracker.FirstFailure) > config.WindowDuration {
		tracker.FailedAttempts = 0
		tracker.FirstFailure = now
	}

	tracker.FailedAttempts++
	tracker.LastFailure = now

	// Check if lockout needed
	if tracker.FailedAttempts >= config.MaxAttempts {
		lockoutDuration := bp.calculateLockoutDuration(tracker.TotalLockouts, config)
		tracker.LockedUntil = now.Add(lockoutDuration)
		tracker.TotalLockouts++
		tracker.FailedAttempts = 0 // Reset counter after lockout

		bp.logger.Warn("IP locked due to repeated failures",
			zap.String("ip", ipAddress),
			zap.Duration("lockout_duration", lockoutDuration),
			zap.Int("total_lockouts", tracker.TotalLockouts),
		)
	}
}

func (bp *BruteForceProtector) calculateLockoutDuration(lockoutCount int, config *BruteForceConfig) time.Duration {
	if !config.ExponentialBackoff {
		return config.LockoutDuration
	}

	// Exponential backoff: base * 2^lockoutCount
	duration := config.LockoutDuration
	for i := 0; i < lockoutCount; i++ {
		duration *= 2
		if duration > config.MaxLockoutDuration {
			return config.MaxLockoutDuration
		}
	}

	return duration
}

func (bp *BruteForceProtector) cleanupExpiredEntries() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		bp.mu.Lock()
		now := time.Now()

		// Clean up email attempts
		for email, tracker := range bp.emailAttempts {
			if now.After(tracker.LockedUntil) && now.Sub(tracker.LastFailure) > 24*time.Hour {
				delete(bp.emailAttempts, email)
			}
		}

		// Clean up IP attempts
		for ip, tracker := range bp.ipAttempts {
			if now.After(tracker.LockedUntil) && now.Sub(tracker.LastFailure) > 24*time.Hour {
				delete(bp.ipAttempts, ip)
			}
		}

		bp.mu.Unlock()

		bp.logger.Debug("Cleaned up expired brute force entries")
	}
}

// RequiresCaptcha determines if CAPTCHA should be required
func (bp *BruteForceProtector) RequiresCaptcha(email, ipAddress string) bool {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	// Require CAPTCHA after 3 failed attempts
	if tracker, exists := bp.emailAttempts[email]; exists && tracker.FailedAttempts >= 3 {
		return true
	}

	if tracker, exists := bp.ipAttempts[ipAddress]; exists && tracker.FailedAttempts >= 3 {
		return true
	}

	return false
}

// GetRemainingAttempts returns how many attempts are left before lockout
func (bp *BruteForceProtector) GetRemainingAttempts(email, ipAddress string) int {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	config := DefaultBruteForceConfig()
	minRemaining := config.MaxAttempts

	if tracker, exists := bp.emailAttempts[email]; exists {
		emailRemaining := config.MaxAttempts - tracker.FailedAttempts
		if emailRemaining < minRemaining {
			minRemaining = emailRemaining
		}
	}

	if tracker, exists := bp.ipAttempts[ipAddress]; exists {
		ipRemaining := config.MaxAttempts - tracker.FailedAttempts
		if ipRemaining < minRemaining {
			minRemaining = ipRemaining
		}
	}

	if minRemaining < 0 {
		return 0
	}
	return minRemaining
}
