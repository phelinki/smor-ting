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
	authService    *services.EnhancedAuthService
	userService    UserService
	otpService     OTPService
	captchaService CaptchaService
	logger         *zap.Logger
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
	authService *services.EnhancedAuthService,
	userService UserService,
	otpService OTPService,
	captchaService CaptchaService,
	logger *zap.Logger,
) *EnhancedAuthHandler {
	return &EnhancedAuthHandler{
		authService:    authService,
		userService:    userService,
		otpService:     otpService,
		captchaService: captchaService,
		logger:         logger,
	}
}

// Enhanced login request structure
type EnhancedLoginRequest struct {
	Email         string                     `json:"email" validate:"required,email"`
	Password      string                     `json:"password" validate:"required"`
	RememberMe    bool                       `json:"remember_me"`
	DeviceInfo    services.DeviceFingerprint `json:"device_info"`
	CaptchaToken  string                     `json:"captcha_token,omitempty"`
	TwoFactorCode string                     `json:"two_factor_code,omitempty"`
}

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
	var req EnhancedLoginRequest

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

	// Check if CAPTCHA is required
	bruteForceProtector := h.getBruteForceProtector()
	requiresCaptcha := bruteForceProtector.RequiresCaptcha(req.Email, clientIP)

	if requiresCaptcha && req.CaptchaToken == "" {
		remainingAttempts := bruteForceProtector.GetRemainingAttempts(req.Email, clientIP)
		return c.Status(http.StatusTooManyRequests).JSON(EnhancedLoginResponse{
			Success:           false,
			Message:           "CAPTCHA verification required",
			RequiresCaptcha:   true,
			RemainingAttempts: remainingAttempts,
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
		bruteForceProtector.RecordFailure(req.Email, clientIP)

		// Return generic error message to prevent email enumeration
		return c.Status(http.StatusUnauthorized).JSON(EnhancedLoginResponse{
			Success: false,
			Message: "Invalid email or password",
		})
	}

	// Verify password
	if err := h.userService.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		h.logger.Warn("Login attempt with invalid password",
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
		)

		// Record failure for brute force protection
		bruteForceProtector.RecordFailure(req.Email, clientIP)

		// Get lockout info for response
		lockoutInfo := bruteForceProtector.GetLockoutInfo(req.Email, clientIP)
		remainingAttempts := bruteForceProtector.GetRemainingAttempts(req.Email, clientIP)

		return c.Status(http.StatusUnauthorized).JSON(EnhancedLoginResponse{
			Success:           false,
			Message:           "Invalid email or password",
			LockoutInfo:       lockoutInfo,
			RemainingAttempts: remainingAttempts,
			RequiresCaptcha:   bruteForceProtector.RequiresCaptcha(req.Email, clientIP),
		})
	}

	// Prepare auth request
	authReq := &services.AuthRequest{
		Email:         req.Email,
		Password:      req.Password,
		RememberMe:    req.RememberMe,
		DeviceInfo:    req.DeviceInfo,
		IPAddress:     clientIP,
		UserAgent:     userAgent,
		CaptchaToken:  req.CaptchaToken,
		TwoFactorCode: req.TwoFactorCode,
	}

	// Perform authentication
	authResult, err := h.authService.Authenticate(c.Context(), authReq, user)
	if err != nil {
		h.logger.Error("Authentication failed",
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
			zap.Error(err),
		)

		// Check if this is a brute force lockout
		lockoutInfo := bruteForceProtector.GetLockoutInfo(req.Email, clientIP)
		if lockoutInfo.EmailLocked || lockoutInfo.IPLocked {
			return c.Status(http.StatusTooManyRequests).JSON(EnhancedLoginResponse{
				Success:     false,
				Message:     err.Error(),
				LockoutInfo: lockoutInfo,
			})
		}

		return c.Status(http.StatusUnauthorized).JSON(EnhancedLoginResponse{
			Success: false,
			Message: "Authentication failed",
		})
	}

	// Handle 2FA requirement
	if authResult.RequiresTwoFactor {
		// Generate and send OTP
		if req.TwoFactorCode == "" {
			_, err := h.otpService.GenerateOTP(c.Context(), user.ID.Hex(), "2fa")
			if err != nil {
				h.logger.Error("Failed to generate 2FA OTP", zap.Error(err))
			}
		}

		return c.Status(http.StatusOK).JSON(EnhancedLoginResponse{
			Success:           true,
			Message:           "Two-factor authentication required",
			User:              authResult.User,
			RequiresTwoFactor: true,
			DeviceTrusted:     authResult.DeviceTrusted,
		})
	}

	// Update user's last login
	if err := h.userService.UpdateLastLogin(c.Context(), user.ID.Hex()); err != nil {
		h.logger.Warn("Failed to update last login", zap.Error(err))
	}

	// Set session cookie
	h.setSessionCookie(c, authResult.SessionID, authResult.RefreshExpiresAt)

	h.logger.Info("User logged in successfully",
		zap.String("user_id", user.ID.Hex()),
		zap.String("email", user.Email),
		zap.String("session_id", authResult.SessionID),
		zap.String("ip", clientIP),
		zap.Bool("device_trusted", authResult.DeviceTrusted),
		zap.Bool("remembered", req.RememberMe),
	)

	return c.Status(http.StatusOK).JSON(EnhancedLoginResponse{
		Success:              true,
		Message:              "Login successful",
		User:                 authResult.User,
		AccessToken:          authResult.AccessToken,
		RefreshToken:         authResult.RefreshToken,
		SessionID:            authResult.SessionID,
		TokenExpiresAt:       authResult.TokenExpiresAt,
		RefreshExpiresAt:     authResult.RefreshExpiresAt,
		RequiresVerification: authResult.RequiresVerification,
		DeviceTrusted:        authResult.DeviceTrusted,
	})
}

// RefreshTokenRequest for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	SessionID    string `json:"session_id" validate:"required"`
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
	authResult, err := h.authService.RefreshTokenWithSession(c.Context(), req.RefreshToken, req.SessionID)
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

	sessions, err := h.authService.GetUserSessions(c.Context(), userID)
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

	err := h.authService.RevokeSession(c.Context(), sessionID)
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

	err := h.authService.RevokeAllSessions(c.Context(), userID)
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
	Email      string                     `json:"email" validate:"required,email"`
	SessionID  string                     `json:"session_id" validate:"required"`
	DeviceInfo services.DeviceFingerprint `json:"device_info"`
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
	session, err := h.authService.GetSessionByID(c.Context(), req.SessionID)
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
	authResult, err := h.authService.GenerateTokensForExistingSession(c.Context(), session, user)
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
	if err := h.authService.UpdateSessionActivity(c.Context(), req.SessionID, clientIP, userAgent); err != nil {
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
		Success:          true,
		Message:          "Biometric authentication successful",
		User:             authResult.User,
		AccessToken:      authResult.AccessToken,
		RefreshToken:     authResult.RefreshToken,
		SessionID:        authResult.SessionID,
		TokenExpiresAt:   authResult.TokenExpiresAt,
		RefreshExpiresAt: authResult.RefreshExpiresAt,
		DeviceTrusted:    true, // Biometric auth implies trusted device
	})
}

// Helper methods

func (h *EnhancedAuthHandler) validateLoginRequest(req *EnhancedLoginRequest) error {
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
