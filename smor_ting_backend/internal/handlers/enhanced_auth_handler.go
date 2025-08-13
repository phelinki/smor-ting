package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"go.uber.org/zap"
)

// EnhancedAuthHandler provides comprehensive authentication endpoints
type EnhancedAuthHandler struct {
	authService         EnhancedAuthService
	userService         UserService
	otpService          OTPService
	captchaService      CaptchaService
	bruteForceProtector *services.BruteForceProtector
	auditService        *services.AuditService
	logger              *zap.Logger
}

// EnhancedAuthService interface for enhanced authentication operations
type EnhancedAuthService interface {
	EnhancedLogin(req *models.EnhancedLoginRequest, clientIP string) (*models.EnhancedAuthResult, error)
	BiometricLogin(sessionID, biometricData string) (*models.EnhancedAuthResult, error)
	GetUserSessions(userID string) ([]*models.SessionInfo, error)
	GetSessionByID(sessionID string) (*models.SessionInfo, error)
	RevokeSession(sessionID string) error
	RevokeAllSessions(userID string) error
	SignOutAllDevices(userID string) error
	RefreshTokenWithSession(refreshToken, sessionID string) (*models.EnhancedAuthResult, error)
	VerifyDeviceFingerprint(existing, provided *models.DeviceFingerprint) bool
	GenerateTokensForExistingSession(sessionID string) (*models.EnhancedAuthResult, error)
	UpdateSessionActivity(sessionID string) error
}

// UserService interface for user operations
type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	VerifyPassword(password, hash string) error
	UpdateLastLogin(ctx context.Context, userID string) error
}

// OTPService interface for OTP operations
type OTPService interface {
	GenerateOTP(ctx context.Context, userID, purpose string) (string, error)
	VerifyOTP(ctx context.Context, userID, otp, purpose string) error
}

// CaptchaService interface for CAPTCHA validation
type CaptchaService interface {
	VerifyCaptcha(token, clientIP string) error
}

// NewEnhancedAuthHandler creates a new enhanced authentication handler
func NewEnhancedAuthHandler(
	authService EnhancedAuthService,
	userService UserService,
	otpService OTPService,
	captchaService CaptchaService,
	auditService *services.AuditService,
	logger *zap.Logger,
) *EnhancedAuthHandler {
	// Fallbacks for nil dependencies to keep tests simple
	if logger == nil {
		logger = zap.NewNop()
	}
	if auditService == nil {
		auditService = services.NewAuditService(nil, logger)
	}

	return &EnhancedAuthHandler{
		authService:         authService,
		userService:         userService,
		otpService:          otpService,
		captchaService:      captchaService,
		bruteForceProtector: services.NewBruteForceProtector(logger),
		auditService:        auditService,
		logger:              logger,
	}
}

// Enhanced login request structure
// Note: EnhancedLoginRequest is now defined in models/enhanced_auth.go

// Enhanced login response structure
type EnhancedLoginResponse struct {
	Success              bool                  `json:"success"`
	Message              string                `json:"message"`
	User                 *models.User          `json:"user,omitempty"`
	AccessToken          string                `json:"access_token,omitempty"`
	RefreshToken         string                `json:"refresh_token,omitempty"`
	SessionID            string                `json:"session_id,omitempty"`
	TokenExpiresAt       time.Time             `json:"token_expires_at,omitempty"`
	RefreshExpiresAt     time.Time             `json:"refresh_expires_at,omitempty"`
	RequiresTwoFactor    bool                  `json:"requires_two_factor"`
	RequiresVerification bool                  `json:"requires_verification"`
	RequiresCaptcha      bool                  `json:"requires_captcha"`
	DeviceTrusted        bool                  `json:"device_trusted"`
	LockoutInfo          *services.LockoutInfo `json:"lockout_info,omitempty"`
	RemainingAttempts    int                   `json:"remaining_attempts,omitempty"`
}

// EnhancedLogin handles comprehensive login with all security features
func (h *EnhancedAuthHandler) EnhancedLogin(c *fiber.Ctx) error {
	var req models.EnhancedLoginRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse login request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(EnhancedLoginResponse{
			Success: false,
			Message: "Invalid request format",
		})
	}

	// Validate request
	if err := h.validateLoginRequest(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(EnhancedLoginResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	// Extract client information
	clientIP := c.IP()
	userAgent := c.Get("User-Agent")

	// Set client info in device fingerprint
	req.DeviceInfo.Platform = h.extractPlatform(userAgent)

	// Check brute force protection
	if err := h.bruteForceProtector.CheckAllowed(req.Email, clientIP); err != nil {
		h.logger.Warn("Login blocked by brute force protection",
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
			zap.Error(err),
		)
		return c.Status(http.StatusTooManyRequests).JSON(EnhancedLoginResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	// Check if CAPTCHA is required
	requiresCaptcha := h.bruteForceProtector.RequiresCaptcha(req.Email, clientIP)

	if requiresCaptcha && req.CaptchaToken == "" {
		return c.Status(http.StatusTooManyRequests).JSON(EnhancedLoginResponse{
			Success:         false,
			Message:         "CAPTCHA verification required",
			RequiresCaptcha: true,
		})
	}

	// Verify CAPTCHA if provided
	if req.CaptchaToken != "" {
		if err := h.captchaService.VerifyCaptcha(req.CaptchaToken, clientIP); err != nil {
			h.logger.Warn("CAPTCHA verification failed",
				zap.String("email", req.Email),
				zap.String("ip", clientIP),
				zap.Error(err),
			)
			return c.Status(http.StatusBadRequest).JSON(EnhancedLoginResponse{
				Success:         false,
				Message:         "CAPTCHA verification failed",
				RequiresCaptcha: true,
			})
		}
	}

	// Get user by email
	user, err := h.userService.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		h.logger.Warn("Login attempt with non-existent email",
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
		)

		// Record failure for brute force protection
		h.bruteForceProtector.RecordFailure(req.Email, clientIP)

		// Log failed login attempt for audit
		h.auditService.LogSecurityEvent(
			c.Context(),
			services.ActionLogin,
			req.Email,
			clientIP,
			userAgent,
			map[string]interface{}{
				"reason": "user_not_found",
				"email":  req.Email,
			},
		)

		// Get lockout info and remaining attempts
		lockoutInfo := h.bruteForceProtector.GetLockoutInfo(req.Email, clientIP)
		remainingAttempts := h.bruteForceProtector.GetRemainingAttempts(req.Email, clientIP)

		// Return generic error message to prevent email enumeration
		return c.Status(http.StatusUnauthorized).JSON(EnhancedLoginResponse{
			Success:           false,
			Message:           "Invalid email or password",
			LockoutInfo:       lockoutInfo,
			RemainingAttempts: remainingAttempts,
			RequiresCaptcha:   h.bruteForceProtector.RequiresCaptcha(req.Email, clientIP),
		})
	}

	// Verify password
	if err := h.userService.VerifyPassword(req.Password, user.Password); err != nil {
		h.logger.Warn("Login attempt with invalid password",
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
		)

		// Record failure for brute force protection
		h.bruteForceProtector.RecordFailure(req.Email, clientIP)

		// Log failed login attempt for audit
		h.auditService.LogUserAction(
			c.Context(),
			user,
			services.ActionLogin,
			"authentication",
			clientIP,
			userAgent,
			false,
			map[string]interface{}{
				"reason": "invalid_password",
			},
		)

		// Get lockout info and remaining attempts
		lockoutInfo := h.bruteForceProtector.GetLockoutInfo(req.Email, clientIP)
		remainingAttempts := h.bruteForceProtector.GetRemainingAttempts(req.Email, clientIP)

		return c.Status(http.StatusUnauthorized).JSON(EnhancedLoginResponse{
			Success:           false,
			Message:           "Invalid email or password",
			LockoutInfo:       lockoutInfo,
			RemainingAttempts: remainingAttempts,
			RequiresCaptcha:   h.bruteForceProtector.RequiresCaptcha(req.Email, clientIP),
		})
	}

	// Call enhanced auth service directly
	authResult, err := h.authService.EnhancedLogin(&req, clientIP)
	if err != nil {
		h.logger.Error("Enhanced login failed", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(EnhancedLoginResponse{
			Success: false,
			Message: "Internal server error",
		})
	}

	if !authResult.Success {
		// Record failure for enhanced auth service failure
		h.bruteForceProtector.RecordFailure(req.Email, clientIP)

		// Get lockout info and remaining attempts
		lockoutInfo := h.bruteForceProtector.GetLockoutInfo(req.Email, clientIP)
		remainingAttempts := h.bruteForceProtector.GetRemainingAttempts(req.Email, clientIP)

		return c.Status(http.StatusUnauthorized).JSON(EnhancedLoginResponse{
			Success:              authResult.Success,
			Message:              authResult.Message,
			RequiresTwoFactor:    authResult.RequiresTwoFactor,
			RequiresVerification: authResult.RequiresVerification,
			RequiresCaptcha:      h.bruteForceProtector.RequiresCaptcha(req.Email, clientIP),
			DeviceTrusted:        authResult.DeviceTrusted,
			LockoutInfo:          lockoutInfo,
			RemainingAttempts:    remainingAttempts,
		})
	}

	// Success - record successful authentication and reset brute force counters
	h.bruteForceProtector.RecordSuccess(req.Email, clientIP)

	// Log successful login for audit
	h.auditService.LogUserAction(
		c.Context(),
		authResult.User,
		services.ActionLogin,
		"authentication",
		clientIP,
		userAgent,
		true,
		map[string]interface{}{
			"session_id":     authResult.SessionID,
			"device_trusted": authResult.DeviceTrusted,
		},
	)

	// Success - return successful login response
	return c.Status(http.StatusOK).JSON(EnhancedLoginResponse{
		Success:      authResult.Success,
		Message:      authResult.Message,
		User:         authResult.User,
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
		SessionID:    authResult.SessionID,
		TokenExpiresAt: func() time.Time {
			if authResult.TokenExpiresAt != nil {
				return *authResult.TokenExpiresAt
			}
			return time.Time{}
		}(),
		RefreshExpiresAt: func() time.Time {
			if authResult.RefreshExpiresAt != nil {
				return *authResult.RefreshExpiresAt
			}
			return time.Time{}
		}(),
		RequiresTwoFactor:    authResult.RequiresTwoFactor,
		RequiresVerification: authResult.RequiresVerification,
		RequiresCaptcha:      authResult.RequiresCaptcha,
		DeviceTrusted:        authResult.DeviceTrusted,
	})
}

// RefreshToken handles token refresh with session validation
func (h *EnhancedAuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse refresh token request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request format",
		})
	}

	// Perform token refresh with session validation
	authResult, err := h.authService.RefreshTokenWithSession(req.RefreshToken, req.SessionID)
	if err != nil {
		h.logger.Warn("Token refresh failed",
			zap.String("session_id", req.SessionID),
			zap.Error(err),
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token refresh failed",
		})
	}

	h.logger.Info("Token refreshed successfully",
		zap.String("user_id", authResult.User.ID.Hex()),
		zap.String("session_id", authResult.SessionID),
	)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success":            true,
		"message":            "Token refreshed successfully",
		"access_token":       authResult.AccessToken,
		"refresh_token":      authResult.RefreshToken,
		"token_expires_at":   authResult.TokenExpiresAt,
		"refresh_expires_at": authResult.RefreshExpiresAt,
	})
}

// GetSessions returns all active sessions for the authenticated user
func (h *EnhancedAuthHandler) GetSessions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	sessions, err := h.authService.GetUserSessions(userID)
	if err != nil {
		h.logger.Error("Failed to get user sessions",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve sessions",
		})
	}

	// Sanitize session data for client
	sanitizedSessions := make([]map[string]interface{}, len(sessions))
	for i, session := range sessions {
		sanitizedSessions[i] = map[string]interface{}{
			"session_id":    session.SessionID,
			"device_info":   session.DeviceInfo,
			"ip_address":    session.IPAddress,
			"user_agent":    session.UserAgent,
			"is_remembered": session.IsRemembered,
			"last_activity": session.LastActivity,
			"created_at":    session.CreatedAt,
			"expires_at":    session.ExpiresAt,
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success":  true,
		"sessions": sanitizedSessions,
	})
}

// RevokeSession revokes a specific session
func (h *EnhancedAuthHandler) RevokeSession(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	sessionID := c.Params("sessionId")

	if sessionID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Session ID is required",
		})
	}

	err := h.authService.RevokeSession(sessionID)
	if err != nil {
		h.logger.Error("Failed to revoke session",
			zap.String("user_id", userID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to revoke session",
		})
	}

	h.logger.Info("Session revoked",
		zap.String("user_id", userID),
		zap.String("session_id", sessionID),
	)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Session revoked successfully",
	})
}

// RevokeAllSessions revokes all sessions for the user (sign out all devices)
func (h *EnhancedAuthHandler) RevokeAllSessions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	err := h.authService.RevokeAllSessions(userID)
	if err != nil {
		h.logger.Error("Failed to revoke all sessions",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to revoke all sessions",
		})
	}

	h.logger.Info("All sessions revoked",
		zap.String("user_id", userID),
	)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "All sessions revoked successfully",
	})
}

// BiometricLoginRequest for biometric authentication
type BiometricLoginRequest struct {
	Email      string                    `json:"email" validate:"required,email"`
	SessionID  string                    `json:"session_id" validate:"required"`
	DeviceInfo *models.DeviceFingerprint `json:"device_info"`
}

// BiometricLogin handles biometric authentication for quick unlock
func (h *EnhancedAuthHandler) BiometricLogin(c *fiber.Ctx) error {
	var req BiometricLoginRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse biometric login request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request format",
		})
	}

	// Validate request
	if req.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email is required",
		})
	}
	if req.SessionID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Session ID is required",
		})
	}

	// Extract client information
	clientIP := c.IP()
	userAgent := c.Get("User-Agent")

	// Set client info in device fingerprint
	req.DeviceInfo.Platform = h.extractPlatform(userAgent)

	// Get user by email
	user, err := h.userService.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		h.logger.Warn("Biometric login attempt with non-existent email",
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or session",
		})
	}

	// Validate session exists and belongs to the user
	session, err := h.authService.GetSessionByID(req.SessionID)
	if err != nil || session.UserID != user.ID.Hex() {
		h.logger.Warn("Biometric login attempt with invalid session",
			zap.String("email", req.Email),
			zap.String("session_id", req.SessionID),
			zap.String("ip", clientIP),
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid session",
		})
	}

	// Check if session is still valid (not expired)
	if session.ExpiresAt.Before(time.Now()) {
		h.logger.Warn("Biometric login attempt with expired session",
			zap.String("email", req.Email),
			zap.String("session_id", req.SessionID),
			zap.String("ip", clientIP),
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Session expired",
		})
	}

	// Verify device fingerprint matches the session (optional security check)
	if !h.authService.VerifyDeviceFingerprint(session.DeviceInfo, req.DeviceInfo) {
		h.logger.Warn("Biometric login attempt with device fingerprint mismatch",
			zap.String("email", req.Email),
			zap.String("session_id", req.SessionID),
			zap.String("ip", clientIP),
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Device verification failed",
		})
	}

	// Generate new tokens for the biometric login
	authResult, err := h.authService.GenerateTokensForExistingSession(req.SessionID)
	if err != nil {
		h.logger.Error("Failed to generate tokens for biometric login",
			zap.String("email", req.Email),
			zap.String("session_id", req.SessionID),
			zap.Error(err),
		)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to authenticate",
		})
	}

	// Update session last activity
	if err := h.authService.UpdateSessionActivity(req.SessionID); err != nil {
		h.logger.Warn("Failed to update session activity", zap.Error(err))
	}

	// Update user's last login
	if err := h.userService.UpdateLastLogin(c.Context(), user.ID.Hex()); err != nil {
		h.logger.Warn("Failed to update last login", zap.Error(err))
	}

	h.logger.Info("User authenticated via biometric login",
		zap.String("user_id", user.ID.Hex()),
		zap.String("email", user.Email),
		zap.String("session_id", req.SessionID),
		zap.String("ip", clientIP),
	)

	return c.Status(http.StatusOK).JSON(EnhancedLoginResponse{
		Success:      true,
		Message:      "Biometric authentication successful",
		User:         authResult.User,
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
		SessionID:    authResult.SessionID,
		TokenExpiresAt: func() time.Time {
			if authResult.TokenExpiresAt != nil {
				return *authResult.TokenExpiresAt
			}
			return time.Time{}
		}(),
		RefreshExpiresAt: func() time.Time {
			if authResult.RefreshExpiresAt != nil {
				return *authResult.RefreshExpiresAt
			}
			return time.Time{}
		}(),
		DeviceTrusted: true, // Biometric auth implies trusted device
	})
}

// Helper methods

func (h *EnhancedAuthHandler) validateLoginRequest(req *models.EnhancedLoginRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (h *EnhancedAuthHandler) extractPlatform(userAgent string) string {
	// Simple platform detection based on user agent
	if userAgent == "" {
		return "unknown"
	}

	// Add more sophisticated detection as needed
	if len(userAgent) > 100 {
		return userAgent[:100] // Truncate for storage
	}
	return userAgent
}

func (h *EnhancedAuthHandler) getBruteForceProtector() *services.BruteForceProtector {
	// This would be injected in a real implementation
	return services.NewBruteForceProtector(h.logger)
}

func (h *EnhancedAuthHandler) setSessionCookie(c *fiber.Ctx, sessionID string, expiresAt time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Domain:   "", // Set appropriate domain
		Path:     "/",
	})
}
