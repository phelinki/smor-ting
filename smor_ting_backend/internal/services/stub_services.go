package services

import (
	"context"
	"fmt"

	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/models"
	"go.uber.org/zap"
)

// UserServiceAdapter adapts the existing auth service to implement UserService interface
type UserServiceAdapter struct {
	authService *auth.MongoDBService
}

func NewUserServiceAdapter(authService *auth.MongoDBService) *UserServiceAdapter {
	return &UserServiceAdapter{authService: authService}
}

func (u *UserServiceAdapter) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// For now, return a stub implementation
	// TODO: Implement proper user retrieval from auth service
	return nil, fmt.Errorf("GetUserByEmail not implemented yet")
}

func (u *UserServiceAdapter) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	// For now, return a stub implementation
	// TODO: Implement proper user retrieval from auth service
	return nil, fmt.Errorf("GetUserByID not implemented yet")
}

func (u *UserServiceAdapter) VerifyPassword(password, hash string) error {
	// For now, return a stub implementation
	// TODO: Implement proper password verification
	return fmt.Errorf("VerifyPassword not implemented yet")
}

func (u *UserServiceAdapter) UpdateLastLogin(ctx context.Context, userID string) error {
	// For now, return a stub implementation
	// TODO: Implement proper last login update
	return fmt.Errorf("UpdateLastLogin not implemented yet")
}

// StubOTPService provides a stub implementation of OTPService
type StubOTPService struct {
	logger *zap.Logger
}

func NewStubOTPService(logger *zap.Logger) *StubOTPService {
	return &StubOTPService{logger: logger}
}

func (s *StubOTPService) GenerateOTP(ctx context.Context, userID, purpose string) (string, error) {
	s.logger.Info("Stub OTP generation", zap.String("userID", userID), zap.String("purpose", purpose))
	return "123456", nil // Stub OTP
}

func (s *StubOTPService) VerifyOTP(ctx context.Context, userID, otp, purpose string) error {
	s.logger.Info("Stub OTP verification", zap.String("userID", userID), zap.String("otp", otp), zap.String("purpose", purpose))
	if otp == "123456" {
		return nil // Accept stub OTP
	}
	return fmt.Errorf("invalid OTP")
}

// StubCaptchaService provides a stub implementation of CaptchaService
type StubCaptchaService struct {
	logger *zap.Logger
}

func NewStubCaptchaService(logger *zap.Logger) *StubCaptchaService {
	return &StubCaptchaService{logger: logger}
}

func (s *StubCaptchaService) VerifyCaptcha(token, clientIP string) error {
	s.logger.Info("Stub CAPTCHA verification", zap.String("token", token), zap.String("clientIP", clientIP))
	// For now, always pass CAPTCHA verification in development
	return nil
}

// StubEnhancedAuthService provides a stub implementation of EnhancedAuthService
type StubEnhancedAuthService struct {
	logger *zap.Logger
}

func NewStubEnhancedAuthService(logger *zap.Logger) *StubEnhancedAuthService {
	return &StubEnhancedAuthService{logger: logger}
}

func (s *StubEnhancedAuthService) EnhancedLogin(req *models.EnhancedLoginRequest, clientIP string) (*models.EnhancedAuthResult, error) {
	s.logger.Info("Stub enhanced login", zap.String("email", req.Email), zap.String("clientIP", clientIP))

	// Return a stub successful result
	return &models.EnhancedAuthResult{
		Success:              true,
		Message:              "Login successful (stub)",
		RequiresTwoFactor:    false,
		RequiresVerification: false,
		RequiresCaptcha:      false,
		DeviceTrusted:        true,
		IsRestoredSession:    false,
	}, nil
}

func (s *StubEnhancedAuthService) BiometricLogin(sessionID, biometricData string) (*models.EnhancedAuthResult, error) {
	s.logger.Info("Stub biometric login", zap.String("sessionID", sessionID))

	return &models.EnhancedAuthResult{
		Success:              true,
		Message:              "Biometric login successful (stub)",
		RequiresTwoFactor:    false,
		RequiresVerification: false,
		RequiresCaptcha:      false,
		DeviceTrusted:        true,
		IsRestoredSession:    true,
	}, nil
}

func (s *StubEnhancedAuthService) GetUserSessions(userID string) ([]*models.SessionInfo, error) {
	s.logger.Info("Stub get user sessions", zap.String("userID", userID))

	// Return stub session data
	return []*models.SessionInfo{
		{
			SessionID:  "session_1",
			UserID:     userID,
			DeviceName: "iPhone 13",
			DeviceType: "mobile",
			IPAddress:  "192.168.1.100",
			IsCurrent:  true,
			IsRevoked:  false,
		},
	}, nil
}

func (s *StubEnhancedAuthService) RevokeSession(sessionID string) error {
	s.logger.Info("Stub revoke session", zap.String("sessionID", sessionID))
	return nil
}

func (s *StubEnhancedAuthService) SignOutAllDevices(userID string) error {
	s.logger.Info("Stub sign out all devices", zap.String("userID", userID))
	return nil
}

func (s *StubEnhancedAuthService) RefreshTokenWithSession(refreshToken, sessionID string) (*models.EnhancedAuthResult, error) {
	s.logger.Info("Stub refresh token with session", zap.String("sessionID", sessionID))

	return &models.EnhancedAuthResult{
		Success:              true,
		Message:              "Token refreshed (stub)",
		AccessToken:          "new_access_token_stub",
		RefreshToken:         "new_refresh_token_stub",
		RequiresTwoFactor:    false,
		RequiresVerification: false,
		RequiresCaptcha:      false,
		DeviceTrusted:        true,
		IsRestoredSession:    true,
	}, nil
}

func (s *StubEnhancedAuthService) GetSessionByID(sessionID string) (*models.SessionInfo, error) {
	s.logger.Info("Stub get session by ID", zap.String("sessionID", sessionID))

	return &models.SessionInfo{
		SessionID:  sessionID,
		UserID:     "user-123",
		DeviceName: "Test Device",
		DeviceType: "mobile",
		IPAddress:  "192.168.1.100",
		IsCurrent:  true,
		IsRevoked:  false,
	}, nil
}

func (s *StubEnhancedAuthService) RevokeAllSessions(userID string) error {
	s.logger.Info("Stub revoke all sessions", zap.String("userID", userID))
	return nil
}

func (s *StubEnhancedAuthService) VerifyDeviceFingerprint(existing, provided *models.DeviceFingerprint) bool {
	s.logger.Info("Stub verify device fingerprint")
	// For stub, always return true
	return true
}

func (s *StubEnhancedAuthService) GenerateTokensForExistingSession(sessionID string) (*models.EnhancedAuthResult, error) {
	s.logger.Info("Stub generate tokens for existing session", zap.String("sessionID", sessionID))

	return &models.EnhancedAuthResult{
		Success:              true,
		Message:              "Tokens generated (stub)",
		AccessToken:          "new_access_token_stub",
		RefreshToken:         "new_refresh_token_stub",
		RequiresTwoFactor:    false,
		RequiresVerification: false,
		RequiresCaptcha:      false,
		DeviceTrusted:        true,
		IsRestoredSession:    true,
	}, nil
}

func (s *StubEnhancedAuthService) UpdateSessionActivity(sessionID string) error {
	s.logger.Info("Stub update session activity", zap.String("sessionID", sessionID))
	return nil
}
