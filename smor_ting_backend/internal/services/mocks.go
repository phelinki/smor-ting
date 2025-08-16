package services

import (
	"context"
	"fmt"
	"time"

	"github.com/smorting/backend/internal/models"
)

// MockJWTRefreshService is a mock implementation for testing
type MockJWTRefreshService struct{}

func (m *MockJWTRefreshService) GenerateTokenPair(user *models.User) (*TokenPair, error) {
	return &TokenPair{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

func (m *MockJWTRefreshService) ValidateRefreshToken(token string) (*RefreshTokenClaims, error) {
	return &RefreshTokenClaims{
		UserID:  "user123",
		Email:   "test@example.com",
		Role:    "customer",
		TokenID: "token123",
	}, nil
}

func (m *MockJWTRefreshService) RefreshAccessToken(refreshToken string, user *models.User) (*TokenPair, error) {
	return m.GenerateTokenPair(user)
}

func (m *MockJWTRefreshService) ValidateAccessToken(token string) (*AccessTokenClaims, error) {
	return &AccessTokenClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Role:   "customer",
	}, nil
}

func (m *MockJWTRefreshService) RevokeRefreshToken(tokenID string) error {
	return nil
}

func (m *MockJWTRefreshService) IsTokenExpired(token string, isRefreshToken bool) (bool, error) {
	return false, nil
}

// MockEncryptionService is a mock implementation for testing
type MockEncryptionService struct{}

func (m *MockEncryptionService) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (m *MockEncryptionService) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (m *MockEncryptionService) Hash(data string) (string, error) {
	return "hashed_" + data, nil
}

func (m *MockEncryptionService) VerifyHash(data, hash string) bool {
	return hash == "hashed_"+data
}

// MockAuthService is a mock implementation for testing
type MockAuthService struct{}

func (m *MockAuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return &models.User{
		Email: email,
		Role:  models.CustomerRole,
	}, nil
}

func (m *MockAuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return &models.User{
		Email: "test@example.com",
		Role:  models.CustomerRole,
	}, nil
}

func (m *MockAuthService) VerifyPassword(password, hash string) error {
	return nil
}

// MockEnhancedAuthService is a mock implementation for testing
type MockEnhancedAuthService struct{}

func (m *MockEnhancedAuthService) EnhancedLogin(req *models.EnhancedLoginRequest, clientIP string) (*models.EnhancedAuthResult, error) {
	return &models.EnhancedAuthResult{
		Success:      true,
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		SessionID:    "mock_session_id",
	}, nil
}

func (m *MockEnhancedAuthService) BiometricLogin(sessionID, biometricData string) (*models.EnhancedAuthResult, error) {
	return &models.EnhancedAuthResult{
		Success:      true,
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		SessionID:    sessionID,
	}, nil
}

func (m *MockEnhancedAuthService) RefreshTokenWithSession(refreshToken, sessionID string) (*models.EnhancedAuthResult, error) {
	return &models.EnhancedAuthResult{
		Success:      true,
		AccessToken:  "new_mock_access_token",
		RefreshToken: "new_mock_refresh_token",
		SessionID:    sessionID,
	}, nil
}

func (m *MockEnhancedAuthService) GetUserSessions(userID string) ([]*models.SessionInfo, error) {
	return []*models.SessionInfo{
		{
			SessionID:    "session1",
			UserID:       userID,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			LastActivity: time.Now(),
		},
	}, nil
}

func (m *MockEnhancedAuthService) GetSessionByID(sessionID string) (*models.SessionInfo, error) {
	return &models.SessionInfo{
		SessionID:    sessionID,
		UserID:       "user123",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		LastActivity: time.Now(),
	}, nil
}

func (m *MockEnhancedAuthService) RevokeSession(sessionID string) error {
	return nil
}

func (m *MockEnhancedAuthService) RevokeAllSessions(userID string) error {
	return nil
}

func (m *MockEnhancedAuthService) SignOutAllDevices(userID string) error {
	return nil
}

func (m *MockEnhancedAuthService) VerifyDeviceFingerprint(existing, provided *models.DeviceFingerprint) bool {
	return true
}

func (m *MockEnhancedAuthService) GenerateTokensForExistingSession(sessionID string) (*models.EnhancedAuthResult, error) {
	return &models.EnhancedAuthResult{
		Success:      true,
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		SessionID:    sessionID,
	}, nil
}

func (m *MockEnhancedAuthService) UpdateSessionActivity(sessionID string) error {
	return nil
}

// MockUserService is a mock implementation for testing
type MockUserService struct {
	users map[string]*models.User
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return &models.User{
		Email: "test@example.com",
		Role:  models.CustomerRole,
	}, nil
}

func (m *MockUserService) VerifyPassword(password, hash string) error {
	// Simulate password verification - fail for "wrong_password"
	if password == "wrong_password" {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func (m *MockUserService) UpdateLastLogin(ctx context.Context, userID string) error {
	return nil
}

// MockOTPService is a mock implementation for testing
type MockOTPService struct{}

func (m *MockOTPService) GenerateOTP(ctx context.Context, userID, purpose string) (string, error) {
	return "123456", nil
}

func (m *MockOTPService) VerifyOTP(ctx context.Context, userID, otp, purpose string) error {
	return nil
}

func (m *MockOTPService) SendOTP(ctx context.Context, userID, purpose string) error {
	return nil
}

// MockCaptchaService is a mock implementation for testing
type MockCaptchaService struct{}

func (m *MockCaptchaService) GenerateCaptcha(ctx context.Context) (string, string, error) {
	return "captcha_id", "ABCD", nil
}

func (m *MockCaptchaService) VerifyCaptcha(token, clientIP string) error {
	return nil
}
