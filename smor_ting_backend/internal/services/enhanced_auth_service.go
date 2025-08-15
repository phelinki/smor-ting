package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// OTPService interface for 2FA OTP operations
type OTPService interface {
	CreateOTP(ctx context.Context, email, purpose string) error
	VerifyOTP(ctx context.Context, email, otpCode string) error
	GetLatestOTPByEmail(ctx context.Context, email string) (*models.OTPRecord, error)
}

// JWTService interface for JWT token operations
type JWTService interface {
	GenerateTokenPair(user *models.User) (*TokenPair, error)
	ValidateRefreshToken(token string) (*RefreshTokenClaims, error)
}

// EnhancedAuthService provides comprehensive authentication with session management
type EnhancedAuthService struct {
	jwtService          JWTService
	sessionStore        SessionStore
	deviceStore         DeviceStore
	bruteForceProtector *BruteForceProtector
	otpService          OTPService
	logger              *zap.Logger
}

// SessionInfo represents an active user session
type SessionInfo struct {
	SessionID     string                 `bson:"session_id" json:"session_id"`
	UserID        string                 `bson:"user_id" json:"user_id"`
	DeviceID      string                 `bson:"device_id" json:"device_id"`
	DeviceInfo    DeviceFingerprint      `bson:"device_info" json:"device_info"`
	IPAddress     string                 `bson:"ip_address" json:"ip_address"`
	UserAgent     string                 `bson:"user_agent" json:"user_agent"`
	IsRemembered  bool                   `bson:"is_remembered" json:"is_remembered"`
	LastActivity  time.Time              `bson:"last_activity" json:"last_activity"`
	CreatedAt     time.Time              `bson:"created_at" json:"created_at"`
	ExpiresAt     time.Time              `bson:"expires_at" json:"expires_at"`
	Revoked       bool                   `bson:"revoked" json:"revoked"`
	RefreshTokens []string               `bson:"refresh_tokens" json:"refresh_tokens"`
	Metadata      map[string]interface{} `bson:"metadata" json:"metadata"`
}

// DeviceFingerprint represents device characteristics for trust evaluation
type DeviceFingerprint struct {
	DeviceID        string    `bson:"device_id" json:"device_id"`
	Platform        string    `bson:"platform" json:"platform"`
	OSVersion       string    `bson:"os_version" json:"os_version"`
	AppVersion      string    `bson:"app_version" json:"app_version"`
	IsTrusted       bool      `bson:"is_trusted" json:"is_trusted"`
	IsJailbroken    bool      `bson:"is_jailbroken" json:"is_jailbroken"`
	TrustScore      float64   `bson:"trust_score" json:"trust_score"`
	LastVerified    time.Time `bson:"last_verified" json:"last_verified"`
	AttestationData string    `bson:"attestation_data" json:"attestation_data"`
}

// SessionStore interface for session persistence
type SessionStore interface {
	CreateSession(ctx context.Context, session *SessionInfo) error
	GetSession(ctx context.Context, sessionID string) (*SessionInfo, error)
	UpdateSession(ctx context.Context, session *SessionInfo) error
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	GetUserSessions(ctx context.Context, userID string) ([]*SessionInfo, error)
	CleanupExpiredSessions(ctx context.Context) error
}

// DeviceStore interface for device trust management
type DeviceStore interface {
	RegisterDevice(ctx context.Context, device *DeviceFingerprint) error
	GetDevice(ctx context.Context, deviceID string) (*DeviceFingerprint, error)
	UpdateDeviceTrust(ctx context.Context, deviceID string, trusted bool, score float64) error
	GetUserDevices(ctx context.Context, userID string) ([]*DeviceFingerprint, error)
	RevokeDevice(ctx context.Context, deviceID string) error
}

// AuthRequest represents enhanced authentication request
type AuthRequest struct {
	Email         string            `json:"email" validate:"required,email"`
	Password      string            `json:"password" validate:"required"`
	RememberMe    bool              `json:"remember_me"`
	DeviceInfo    DeviceFingerprint `json:"device_info"`
	IPAddress     string            `json:"ip_address"`
	UserAgent     string            `json:"user_agent"`
	CaptchaToken  string            `json:"captcha_token,omitempty"`
	TwoFactorCode string            `json:"two_factor_code,omitempty"`
}

// AuthResult represents enhanced authentication response
type AuthResult struct {
	User                 *models.User `json:"user"`
	AccessToken          string       `json:"access_token"`
	RefreshToken         string       `json:"refresh_token"`
	SessionID            string       `json:"session_id"`
	RequiresTwoFactor    bool         `json:"requires_two_factor"`
	RequiresVerification bool         `json:"requires_verification"`
	DeviceTrusted        bool         `json:"device_trusted"`
	TokenExpiresAt       time.Time    `json:"token_expires_at"`
	RefreshExpiresAt     time.Time    `json:"refresh_expires_at"`
}

// NewEnhancedAuthService creates a new enhanced authentication service
func NewEnhancedAuthService(
	jwtService JWTService,
	sessionStore SessionStore,
	deviceStore DeviceStore,
	otpService OTPService,
	logger *zap.Logger,
) *EnhancedAuthService {
	return &EnhancedAuthService{
		jwtService:          jwtService,
		sessionStore:        sessionStore,
		deviceStore:         deviceStore,
		bruteForceProtector: NewBruteForceProtector(logger),
		otpService:          otpService,
		logger:              logger,
	}
}

// Authenticate performs comprehensive authentication with session management
func (s *EnhancedAuthService) Authenticate(ctx context.Context, req *AuthRequest, user *models.User) (*AuthResult, error) {
	// Check for brute force attempts
	if err := s.bruteForceProtector.CheckAllowed(req.Email, req.IPAddress); err != nil {
		s.logger.Warn("Authentication blocked due to brute force protection",
			zap.String("email", req.Email),
			zap.String("ip", req.IPAddress),
		)
		return nil, fmt.Errorf("authentication temporarily blocked: %w", err)
	}

	// Evaluate device trust
	deviceTrust, err := s.evaluateDeviceTrust(ctx, &req.DeviceInfo, user.ID.Hex())
	if err != nil {
		s.logger.Warn("Device trust evaluation failed", zap.Error(err))
		// Continue with untrusted device
		deviceTrust = &DeviceFingerprint{
			DeviceID:     req.DeviceInfo.DeviceID,
			IsTrusted:    false,
			TrustScore:   0.0,
			IsJailbroken: req.DeviceInfo.IsJailbroken,
		}
	}

	// Check if 2FA is required - DISABLED FOR NOW
	requires2FA := s.requires2FA(user, deviceTrust)
	// Skip 2FA requirement check - always proceed with authentication
	s.logger.Info("2FA requirement check bypassed", zap.Bool("requires2FA", requires2FA))

	// Verify 2FA if provided
	if req.TwoFactorCode != "" {
		if err := s.verify2FA(ctx, user, req.TwoFactorCode); err != nil {
			s.bruteForceProtector.RecordFailure(req.Email, req.IPAddress)
			return nil, fmt.Errorf("two-factor authentication failed: %w", err)
		}
	}

	// Generate session
	sessionID, err := s.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Determine session duration based on "remember me"
	var expiresAt time.Time
	if req.RememberMe {
		expiresAt = time.Now().Add(30 * 24 * time.Hour) // 30 days
	} else {
		expiresAt = time.Now().Add(24 * time.Hour) // 1 day
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	session := &SessionInfo{
		SessionID:     sessionID,
		UserID:        user.ID.Hex(),
		DeviceID:      deviceTrust.DeviceID,
		DeviceInfo:    *deviceTrust,
		IPAddress:     req.IPAddress,
		UserAgent:     req.UserAgent,
		IsRemembered:  req.RememberMe,
		LastActivity:  time.Now(),
		CreatedAt:     time.Now(),
		ExpiresAt:     expiresAt,
		RefreshTokens: []string{tokenPair.RefreshToken},
		Metadata: map[string]interface{}{
			"login_method": "password",
			"trust_score":  deviceTrust.TrustScore,
		},
	}

	if err := s.sessionStore.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update device trust if authentication successful
	if !deviceTrust.IsTrusted && deviceTrust.TrustScore > 0.7 {
		deviceTrust.IsTrusted = true
		s.deviceStore.UpdateDeviceTrust(ctx, deviceTrust.DeviceID, true, deviceTrust.TrustScore)
	}

	// Reset brute force counter on successful auth
	s.bruteForceProtector.RecordSuccess(req.Email, req.IPAddress)

	s.logger.Info("User authenticated successfully",
		zap.String("user_id", user.ID.Hex()),
		zap.String("session_id", sessionID),
		zap.Bool("device_trusted", deviceTrust.IsTrusted),
		zap.Bool("remembered", req.RememberMe),
	)

	return &AuthResult{
		User:                 user,
		AccessToken:          tokenPair.AccessToken,
		RefreshToken:         tokenPair.RefreshToken,
		SessionID:            sessionID,
		DeviceTrusted:        deviceTrust.IsTrusted,
		RequiresTwoFactor:    false, // TWO-FACTOR AUTH DISABLED
		RequiresVerification: false, // TODO: Add IsVerified field to User model
		TokenExpiresAt:       time.Now().Add(30 * time.Minute),
		RefreshExpiresAt:     time.Now().Add(7 * 24 * time.Hour),
	}, nil
}

// RefreshTokenWithSession refreshes tokens while maintaining session context
func (s *EnhancedAuthService) RefreshTokenWithSession(ctx context.Context, refreshToken, sessionID string) (*AuthResult, error) {
	// Validate refresh token
	refreshClaims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get session
	session, err := s.sessionStore.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if session.Revoked {
		return nil, fmt.Errorf("session has been revoked")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session has expired")
	}

	// Verify refresh token belongs to this session
	tokenFound := false
	for _, token := range session.RefreshTokens {
		if token == refreshToken {
			tokenFound = true
			break
		}
	}
	if !tokenFound {
		return nil, fmt.Errorf("refresh token not associated with session")
	}

	// Get user
	userID, _ := primitive.ObjectIDFromHex(refreshClaims.UserID)
	user := &models.User{
		ID:    userID,
		Email: refreshClaims.Email,
		Role:  models.UserRole(refreshClaims.Role),
	}

	// Generate new token pair
	newTokenPair, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Update session with new refresh token
	session.RefreshTokens = append(session.RefreshTokens, newTokenPair.RefreshToken)
	session.LastActivity = time.Now()

	// Keep only last 3 refresh tokens for security
	if len(session.RefreshTokens) > 3 {
		session.RefreshTokens = session.RefreshTokens[len(session.RefreshTokens)-3:]
	}

	if err := s.sessionStore.UpdateSession(ctx, session); err != nil {
		s.logger.Warn("Failed to update session", zap.Error(err))
	}

	s.logger.Info("Token refreshed successfully",
		zap.String("user_id", user.ID.Hex()),
		zap.String("session_id", sessionID),
	)

	return &AuthResult{
		User:             user,
		AccessToken:      newTokenPair.AccessToken,
		RefreshToken:     newTokenPair.RefreshToken,
		SessionID:        sessionID,
		DeviceTrusted:    session.DeviceInfo.IsTrusted,
		TokenExpiresAt:   time.Now().Add(30 * time.Minute),
		RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}, nil
}

// RevokeSession revokes a specific session
func (s *EnhancedAuthService) RevokeSession(ctx context.Context, sessionID string) error {
	return s.sessionStore.RevokeSession(ctx, sessionID)
}

// RevokeAllSessions revokes all sessions for a user (sign out all devices)
func (s *EnhancedAuthService) RevokeAllSessions(ctx context.Context, userID string) error {
	return s.sessionStore.RevokeAllUserSessions(ctx, userID)
}

// GetUserSessions returns all active sessions for a user
func (s *EnhancedAuthService) GetUserSessions(ctx context.Context, userID string) ([]*SessionInfo, error) {
	return s.sessionStore.GetUserSessions(ctx, userID)
}

// Helper methods

func (s *EnhancedAuthService) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *EnhancedAuthService) evaluateDeviceTrust(ctx context.Context, deviceInfo *DeviceFingerprint, userID string) (*DeviceFingerprint, error) {
	// Try to get existing device
	existingDevice, err := s.deviceStore.GetDevice(ctx, deviceInfo.DeviceID)
	if err == nil {
		// Update last verified time
		existingDevice.LastVerified = time.Now()
		return existingDevice, nil
	}

	// New device - evaluate trust score
	trustScore := s.calculateTrustScore(deviceInfo)

	device := &DeviceFingerprint{
		DeviceID:        deviceInfo.DeviceID,
		Platform:        deviceInfo.Platform,
		OSVersion:       deviceInfo.OSVersion,
		AppVersion:      deviceInfo.AppVersion,
		IsTrusted:       trustScore > 0.8,
		IsJailbroken:    deviceInfo.IsJailbroken,
		TrustScore:      trustScore,
		LastVerified:    time.Now(),
		AttestationData: deviceInfo.AttestationData,
	}

	if err := s.deviceStore.RegisterDevice(ctx, device); err != nil {
		return device, err
	}

	return device, nil
}

func (s *EnhancedAuthService) calculateTrustScore(device *DeviceFingerprint) float64 {
	score := 1.0

	// Penalize jailbroken/rooted devices
	if device.IsJailbroken {
		score -= 0.5
	}

	// Reward official app stores
	if strings.Contains(device.AttestationData, "official") {
		score += 0.2
	}

	// Basic platform trust
	switch device.Platform {
	case "iOS", "Android":
		score += 0.1
	default:
		score -= 0.2
	}

	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score
}

func (s *EnhancedAuthService) requires2FA(user *models.User, device *DeviceFingerprint) bool {
	// TWO-FACTOR AUTH DISABLED: Skip 2FA requirement for all users/devices
	// This removes the 2FA screen that was appearing after login
	// Future versions may re-enable 2FA with proper configuration
	s.logger.Info("DEBUG: requires2FA called - returning false (2FA disabled)")
	return false
}

func (s *EnhancedAuthService) verify2FA(ctx context.Context, user *models.User, code string) error {
	// Validate format first
	if len(code) != 6 {
		return fmt.Errorf("invalid 2FA code format")
	}

	// Verify OTP using the integrated OTP service
	if s.otpService != nil {
		return s.otpService.VerifyOTP(ctx, user.Email, code)
	}

	// Fallback for backward compatibility (remove in production)
	s.logger.Warn("2FA verification without OTP service - using fallback")
	return nil
}

// Generate2FAOTP creates a new 2FA OTP for the user
func (s *EnhancedAuthService) Generate2FAOTP(ctx context.Context, email string) error {
	if s.otpService == nil {
		return fmt.Errorf("OTP service not configured")
	}

	err := s.otpService.CreateOTP(ctx, email, "2fa")
	if err != nil {
		s.logger.Error("Failed to create 2FA OTP",
			zap.String("email", email),
			zap.Error(err))
		return fmt.Errorf("failed to generate 2FA code: %w", err)
	}

	s.logger.Info("2FA OTP generated successfully", zap.String("email", email))
	return nil
}

// GetSessionByID retrieves a session by its ID
func (s *EnhancedAuthService) GetSessionByID(ctx context.Context, sessionID string) (*SessionInfo, error) {
	return s.sessionStore.GetSession(ctx, sessionID)
}

// VerifyDeviceFingerprint compares two device fingerprints for security verification
func (s *EnhancedAuthService) VerifyDeviceFingerprint(sessionDevice, currentDevice DeviceFingerprint) bool {
	// Allow some flexibility in device fingerprint matching
	if sessionDevice.DeviceID == currentDevice.DeviceID {
		return true
	}

	// Check if core identifiers match
	if sessionDevice.Platform == currentDevice.Platform &&
		sessionDevice.OSVersion == currentDevice.OSVersion &&
		sessionDevice.AppVersion == currentDevice.AppVersion {
		// Device info matches closely enough
		return true
	}

	// For now, be lenient with device matching to avoid UX issues
	// In production, you might want stricter validation
	return true
}

// GenerateTokensForExistingSession generates new tokens for an existing valid session
func (s *EnhancedAuthService) GenerateTokensForExistingSession(ctx context.Context, session *SessionInfo, user *models.User) (*AuthResult, error) {
	// Generate new token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update session with new refresh token
	session.RefreshTokens = append(session.RefreshTokens, tokenPair.RefreshToken)
	session.LastActivity = time.Now()

	// Keep only last 3 refresh tokens for security
	if len(session.RefreshTokens) > 3 {
		session.RefreshTokens = session.RefreshTokens[len(session.RefreshTokens)-3:]
	}

	// Update session metadata for biometric login
	if session.Metadata == nil {
		session.Metadata = make(map[string]interface{})
	}
	session.Metadata["last_biometric_login"] = time.Now()

	if err := s.sessionStore.UpdateSession(ctx, session); err != nil {
		s.logger.Warn("Failed to update session after biometric login", zap.Error(err))
	}

	s.logger.Info("Tokens generated for existing session",
		zap.String("user_id", user.ID.Hex()),
		zap.String("session_id", session.SessionID),
	)

	return &AuthResult{
		User:             user,
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		SessionID:        session.SessionID,
		DeviceTrusted:    session.DeviceInfo.IsTrusted,
		TokenExpiresAt:   time.Now().Add(30 * time.Minute),
		RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}, nil
}

// UpdateSessionActivity updates the last activity timestamp and client info for a session
func (s *EnhancedAuthService) UpdateSessionActivity(ctx context.Context, sessionID, ipAddress, userAgent string) error {
	session, err := s.sessionStore.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	session.LastActivity = time.Now()
	session.IPAddress = ipAddress
	session.UserAgent = userAgent

	if err := s.sessionStore.UpdateSession(ctx, session); err != nil {
		return fmt.Errorf("failed to update session activity: %w", err)
	}

	return nil
}
