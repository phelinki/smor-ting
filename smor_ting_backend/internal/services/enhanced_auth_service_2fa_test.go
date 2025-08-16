package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap/zaptest"
)

// Mock OTP service for testing 2FA integration
type mockOTPService struct {
	otps map[string]*models.OTPRecord
}

func newMockOTPService() *mockOTPService {
	return &mockOTPService{
		otps: make(map[string]*models.OTPRecord),
	}
}

func (m *mockOTPService) CreateOTP(ctx context.Context, email, purpose string) error {
	otp := &models.OTPRecord{
		ID:        primitive.NewObjectID(),
		Email:     email,
		OTP:       "123456", // Fixed OTP for testing
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
		IsUsed:    false,
	}

	// Store by email+purpose key
	m.otps[email+":"+purpose] = otp
	return nil
}

func (m *mockOTPService) VerifyOTP(ctx context.Context, email, otpCode string) error {
	// Look for any unused OTP for this email
	for _, otp := range m.otps {
		if otp.Email == email && otp.OTP == otpCode && !otp.IsUsed && otp.ExpiresAt.After(time.Now()) {
			otp.IsUsed = true
			return nil
		}
	}
	return fmt.Errorf("invalid or expired OTP")
}

func (m *mockOTPService) GetLatestOTPByEmail(ctx context.Context, email string) (*models.OTPRecord, error) {
	for _, otp := range m.otps {
		if otp.Email == email && !otp.IsUsed {
			return otp, nil
		}
	}
	return nil, fmt.Errorf("no OTP found")
}

// Mock JWT service for testing
type mockJWTService struct{}

func (m *mockJWTService) GenerateTokenPair(user *models.User) (*TokenPair, error) {
	return &TokenPair{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    1800,
	}, nil
}

func (m *mockJWTService) ValidateRefreshToken(token string) (*RefreshTokenClaims, error) {
	return &RefreshTokenClaims{
		UserID:  "user123",
		Email:   "test@example.com",
		Role:    "customer",
		TokenID: "token123",
	}, nil
}

func setupTestServiceWith2FA(t *testing.T) (*EnhancedAuthService, *mockOTPService) {
	logger := zaptest.NewLogger(t)
	sessionStore := newMockSessionStore()
	deviceStore := newMockDeviceStore()
	otpService := newMockOTPService()

	service := &EnhancedAuthService{
		jwtService:          &mockJWTService{},
		sessionStore:        sessionStore,
		deviceStore:         deviceStore,
		bruteForceProtector: NewBruteForceProtector(logger),
		logger:              logger,
		otpService:          otpService,
	}

	return service, otpService
}

func TestEnhancedAuthService_Verify2FA_Success(t *testing.T) {
	// Arrange
	service, otpService := setupTestServiceWith2FA(t)
	user := createTestUser()

	// Create a 2FA OTP
	err := otpService.CreateOTP(context.Background(), user.Email, "2fa")
	if err != nil {
		t.Fatalf("Failed to create OTP: %v", err)
	}

	// Act
	err = service.verify2FA(context.Background(), user, "123456")

	// Assert
	if err != nil {
		t.Errorf("Expected successful 2FA verification, got error: %v", err)
	}
}

func TestEnhancedAuthService_Verify2FA_InvalidCode(t *testing.T) {
	// Arrange
	service, otpService := setupTestServiceWith2FA(t)
	user := createTestUser()

	// Create a 2FA OTP
	err := otpService.CreateOTP(context.Background(), user.Email, "2fa")
	if err != nil {
		t.Fatalf("Failed to create OTP: %v", err)
	}

	// Act
	err = service.verify2FA(context.Background(), user, "654321")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid 2FA code")
	}

	expectedError := "invalid or expired OTP"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEnhancedAuthService_Verify2FA_ExpiredOTP(t *testing.T) {
	// Arrange
	service, otpService := setupTestServiceWith2FA(t)
	user := createTestUser()

	// Create an expired OTP manually
	expiredOTP := &models.OTPRecord{
		ID:        primitive.NewObjectID(),
		Email:     user.Email,
		OTP:       "123456",
		Purpose:   "2fa",
		ExpiresAt: time.Now().Add(-5 * time.Minute), // Expired
		CreatedAt: time.Now().Add(-15 * time.Minute),
		IsUsed:    false,
	}
	otpService.otps[user.Email+":2fa"] = expiredOTP

	// Act
	err := service.verify2FA(context.Background(), user, "123456")

	// Assert
	if err == nil {
		t.Error("Expected error for expired OTP")
	}
}

func TestEnhancedAuthService_Verify2FA_InvalidFormat(t *testing.T) {
	// Arrange
	service, _ := setupTestServiceWith2FA(t)
	user := createTestUser()

	// Act
	err := service.verify2FA(context.Background(), user, "12345") // Too short

	// Assert
	if err == nil {
		t.Error("Expected error for invalid format")
	}

	expectedError := "invalid 2FA code format"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEnhancedAuthService_Verify2FA_UsedOTP(t *testing.T) {
	// Arrange
	service, otpService := setupTestServiceWith2FA(t)
	user := createTestUser()

	// Create and use an OTP
	err := otpService.CreateOTP(context.Background(), user.Email, "2fa")
	if err != nil {
		t.Fatalf("Failed to create OTP: %v", err)
	}

	// Use the OTP once
	err = service.verify2FA(context.Background(), user, "123456")
	if err != nil {
		t.Fatalf("Failed to verify OTP first time: %v", err)
	}

	// Act - try to use same OTP again
	err = service.verify2FA(context.Background(), user, "123456")

	// Assert
	if err == nil {
		t.Error("Expected error for already used OTP")
	}
}

func TestEnhancedAuthService_Generate2FAOTP_Success(t *testing.T) {
	// Arrange
	service, otpService := setupTestServiceWith2FA(t)
	user := createTestUser()

	// Act
	err := service.Generate2FAOTP(context.Background(), user.Email)

	// Assert
	if err != nil {
		t.Errorf("Expected successful OTP generation, got error: %v", err)
	}

	// Verify OTP was created
	otp, err := otpService.GetLatestOTPByEmail(context.Background(), user.Email)
	if err != nil {
		t.Errorf("Expected OTP to be created, got error: %v", err)
	}

	if otp.Purpose != "2fa" {
		t.Errorf("Expected OTP purpose to be '2fa', got '%s'", otp.Purpose)
	}

	if len(otp.OTP) != 6 {
		t.Errorf("Expected OTP to be 6 digits, got %d", len(otp.OTP))
	}
}

func TestEnhancedAuthService_FullWorkflow_2FA_Disabled(t *testing.T) {
	// Arrange
	service, _ := setupTestServiceWith2FA(t)
	user := createTestUser()
	user.Role = models.AdminRole // Admin users no longer require 2FA

	deviceInfo := createTestDeviceInfo()
	deviceInfo.IsTrusted = false // Even untrusted devices no longer require 2FA

	authReq := &AuthRequest{
		Email:      user.Email,
		Password:   "password",
		DeviceInfo: deviceInfo,
		IPAddress:  "192.168.1.1",
		UserAgent:  "TestAgent/1.0",
	}

	// Act - Authentication should complete without requiring 2FA
	result, err := service.Authenticate(context.Background(), authReq, user)
	if err != nil {
		t.Fatalf("Expected successful authentication, got error: %v", err)
	}

	// 2FA is now globally disabled
	if result.RequiresTwoFactor {
		t.Error("Expected 2FA to be disabled globally - should not require 2FA for any user")
	}

	// Should get tokens immediately without 2FA
	if result.AccessToken == "" {
		t.Error("Expected access token to be provided immediately (2FA disabled)")
	}

	if result.RefreshToken == "" {
		t.Error("Expected refresh token to be provided immediately (2FA disabled)")
	}

	if result.SessionID == "" {
		t.Error("Expected session ID to be provided immediately (2FA disabled)")
	}

	// Test that even providing a 2FA code doesn't change the behavior
	authReq.TwoFactorCode = "123456"
	result2, err := service.Authenticate(context.Background(), authReq, user)

	// Assert - behavior should be the same whether 2FA code is provided or not
	if err != nil {
		t.Fatalf("Expected successful authentication even with 2FA code, got error: %v", err)
	}

	if result2.RequiresTwoFactor {
		t.Error("Expected 2FA to remain disabled even when 2FA code is provided")
	}
}
