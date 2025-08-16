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

// Mock implementations for testing

type mockSessionStore struct {
	sessions map[string]*SessionInfo
}

func newMockSessionStore() *mockSessionStore {
	return &mockSessionStore{
		sessions: make(map[string]*SessionInfo),
	}
}

func (m *mockSessionStore) CreateSession(ctx context.Context, session *SessionInfo) error {
	m.sessions[session.SessionID] = session
	return nil
}

func (m *mockSessionStore) GetSession(ctx context.Context, sessionID string) (*SessionInfo, error) {
	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (m *mockSessionStore) UpdateSession(ctx context.Context, session *SessionInfo) error {
	m.sessions[session.SessionID] = session
	return nil
}

func (m *mockSessionStore) RevokeSession(ctx context.Context, sessionID string) error {
	if session, exists := m.sessions[sessionID]; exists {
		session.Revoked = true
		return nil
	}
	return fmt.Errorf("session not found")
}

func (m *mockSessionStore) RevokeAllUserSessions(ctx context.Context, userID string) error {
	for _, session := range m.sessions {
		if session.UserID == userID {
			session.Revoked = true
		}
	}
	return nil
}

func (m *mockSessionStore) GetUserSessions(ctx context.Context, userID string) ([]*SessionInfo, error) {
	var sessions []*SessionInfo
	for _, session := range m.sessions {
		if session.UserID == userID && !session.Revoked {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (m *mockSessionStore) CleanupExpiredSessions(ctx context.Context) error {
	now := time.Now()
	for sessionID, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, sessionID)
		}
	}
	return nil
}

type mockDeviceStore struct {
	devices map[string]*DeviceFingerprint
}

func newMockDeviceStore() *mockDeviceStore {
	return &mockDeviceStore{
		devices: make(map[string]*DeviceFingerprint),
	}
}

func (m *mockDeviceStore) RegisterDevice(ctx context.Context, device *DeviceFingerprint) error {
	m.devices[device.DeviceID] = device
	return nil
}

func (m *mockDeviceStore) GetDevice(ctx context.Context, deviceID string) (*DeviceFingerprint, error) {
	device, exists := m.devices[deviceID]
	if !exists {
		return nil, fmt.Errorf("device not found")
	}
	return device, nil
}

func (m *mockDeviceStore) UpdateDeviceTrust(ctx context.Context, deviceID string, trusted bool, score float64) error {
	if device, exists := m.devices[deviceID]; exists {
		device.IsTrusted = trusted
		device.TrustScore = score
		device.LastVerified = time.Now()
		return nil
	}
	return fmt.Errorf("device not found")
}

func (m *mockDeviceStore) GetUserDevices(ctx context.Context, userID string) ([]*DeviceFingerprint, error) {
	var devices []*DeviceFingerprint
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices, nil
}

func (m *mockDeviceStore) RevokeDevice(ctx context.Context, deviceID string) error {
	delete(m.devices, deviceID)
	return nil
}

// mockJWTService is now defined in enhanced_auth_service_2fa_test.go and shared

func setupTestService(t *testing.T) (*EnhancedAuthService, *mockSessionStore, *mockDeviceStore) {
	logger := zaptest.NewLogger(t)
	sessionStore := newMockSessionStore()
	deviceStore := newMockDeviceStore()

	service := &EnhancedAuthService{
		jwtService:          &mockJWTService{},
		sessionStore:        sessionStore,
		deviceStore:         deviceStore,
		bruteForceProtector: NewBruteForceProtector(logger),
		logger:              logger,
		otpService:          nil, // No OTP service in regular tests
	}

	return service, sessionStore, deviceStore
}

func createTestUser() *models.User {
	return &models.User{
		ID:              primitive.NewObjectID(),
		Email:           "test@example.com",
		Password:        "hashed_password",
		Role:            models.CustomerRole,
		IsEmailVerified: true,
	}
}

func createTestDeviceInfo() DeviceFingerprint {
	return DeviceFingerprint{
		DeviceID:        "device123",
		Platform:        "Android",
		OSVersion:       "11.0",
		AppVersion:      "1.0.0",
		IsTrusted:       false,
		IsJailbroken:    false,
		TrustScore:      0.8,
		LastVerified:    time.Now(),
		AttestationData: "official",
	}
}

func TestEnhancedAuthService_Authenticate_Success(t *testing.T) {
	// Arrange
	service, sessionStore, deviceStore := setupTestService(t)
	user := createTestUser()
	deviceInfo := createTestDeviceInfo()

	authReq := &AuthRequest{
		Email:      user.Email,
		Password:   "password",
		RememberMe: false,
		DeviceInfo: deviceInfo,
		IPAddress:  "192.168.1.1",
		UserAgent:  "TestAgent/1.0",
	}

	// Act
	result, err := service.Authenticate(context.Background(), authReq, user)

	// Assert
	if err != nil {
		t.Fatalf("Expected successful authentication, got error: %v", err)
	}

	if result.User.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID.Hex(), result.User.ID.Hex())
	}

	if result.AccessToken == "" {
		t.Error("Expected access token to be set")
	}

	if result.RefreshToken == "" {
		t.Error("Expected refresh token to be set")
	}

	if result.SessionID == "" {
		t.Error("Expected session ID to be set")
	}

	// Verify session was created
	sessions := sessionStore.sessions
	if len(sessions) != 1 {
		t.Errorf("Expected 1 session, got %d", len(sessions))
	}

	// Verify device was registered
	devices := deviceStore.devices
	if len(devices) != 1 {
		t.Errorf("Expected 1 device, got %d", len(devices))
	}
}

func TestEnhancedAuthService_Authenticate_WithRememberMe(t *testing.T) {
	// Arrange
	service, sessionStore, _ := setupTestService(t)
	user := createTestUser()
	deviceInfo := createTestDeviceInfo()

	authReq := &AuthRequest{
		Email:      user.Email,
		Password:   "password",
		RememberMe: true,
		DeviceInfo: deviceInfo,
		IPAddress:  "192.168.1.1",
		UserAgent:  "TestAgent/1.0",
	}

	// Act
	result, err := service.Authenticate(context.Background(), authReq, user)

	// Assert
	if err != nil {
		t.Fatalf("Expected successful authentication, got error: %v", err)
	}

	// Verify session has extended expiry
	session := sessionStore.sessions[result.SessionID]
	expectedMinExpiry := time.Now().Add(29 * 24 * time.Hour) // Allow some tolerance
	if session.ExpiresAt.Before(expectedMinExpiry) {
		t.Errorf("Expected session to expire after %v, got %v", expectedMinExpiry, session.ExpiresAt)
	}

	if !session.IsRemembered {
		t.Error("Expected session to be marked as remembered")
	}
}

func TestEnhancedAuthService_Authenticate_AdminUserNoLongerRequires2FA(t *testing.T) {
	// Arrange
	service, _, _ := setupTestService(t)
	user := createTestUser()
	user.Role = models.AdminRole // Admin users no longer require 2FA

	deviceInfo := createTestDeviceInfo()
	authReq := &AuthRequest{
		Email:      user.Email,
		Password:   "password",
		DeviceInfo: deviceInfo,
		IPAddress:  "192.168.1.1",
		UserAgent:  "TestAgent/1.0",
	}

	// Act
	result, err := service.Authenticate(context.Background(), authReq, user)

	// Assert
	if err != nil {
		t.Fatalf("Expected successful authentication, got error: %v", err)
	}

	// 2FA is now disabled for all users, including admin users
	if result.RequiresTwoFactor {
		t.Error("Expected 2FA to be disabled for admin user (2FA is now disabled globally)")
	}

	// Should get access token immediately since 2FA is disabled
	if result.AccessToken == "" {
		t.Error("Expected access token to be provided immediately when 2FA is disabled")
	}

	// Should get refresh token too
	if result.RefreshToken == "" {
		t.Error("Expected refresh token to be provided immediately when 2FA is disabled")
	}

	// Should have a session ID
	if result.SessionID == "" {
		t.Error("Expected session ID to be provided")
	}
}

func TestEnhancedAuthService_Authenticate_JailbrokenDevice_2FA_Disabled(t *testing.T) {
	// Arrange
	service, _, _ := setupTestService(t)
	user := createTestUser()
	deviceInfo := createTestDeviceInfo()
	deviceInfo.IsJailbroken = true

	authReq := &AuthRequest{
		Email:         user.Email,
		Password:      "password",
		DeviceInfo:    deviceInfo,
		IPAddress:     "192.168.1.1",
		UserAgent:     "TestAgent/1.0",
		TwoFactorCode: "",
	}

	// Act
	result, err := service.Authenticate(context.Background(), authReq, user)

	// Assert
	if err != nil {
		t.Fatalf("Expected successful authentication, got error: %v", err)
	}

	// 2FA is globally disabled, so even jailbroken devices don't require it
	if result.RequiresTwoFactor {
		t.Error("Expected 2FA to be disabled for jailbroken device (2FA globally disabled)")
	}

	// Should still get tokens immediately even for jailbroken devices
	if result.AccessToken == "" {
		t.Error("Expected access token to be provided immediately for jailbroken device")
	}

	if result.RefreshToken == "" {
		t.Error("Expected refresh token to be provided immediately for jailbroken device")
	}

	// Device may still be untrusted due to jailbreaking, but this doesn't affect login anymore
	// (The trust status is informational only when 2FA is disabled)
}

func TestEnhancedAuthService_RefreshTokenWithSession_Success(t *testing.T) {
	// Arrange
	service, sessionStore, _ := setupTestService(t)
	user := createTestUser()

	// Create a session with refresh token
	sessionID := "session123"
	refreshToken := "refresh_token_123"
	session := &SessionInfo{
		SessionID:     sessionID,
		UserID:        user.ID.Hex(),
		DeviceID:      "device123",
		RefreshTokens: []string{refreshToken},
		ExpiresAt:     time.Now().Add(24 * time.Hour),
		Revoked:       false,
		LastActivity:  time.Now(),
	}
	sessionStore.sessions[sessionID] = session

	// Act
	result, err := service.RefreshTokenWithSession(context.Background(), refreshToken, sessionID)

	// Assert
	if err != nil {
		t.Fatalf("Expected successful token refresh, got error: %v", err)
	}

	if result.AccessToken == "" {
		t.Error("Expected new access token to be set")
	}

	if result.RefreshToken == "" {
		t.Error("Expected new refresh token to be set")
	}

	// Verify session was updated with new refresh token
	updatedSession := sessionStore.sessions[sessionID]
	if len(updatedSession.RefreshTokens) < 2 {
		t.Error("Expected session to have multiple refresh tokens")
	}
}

func TestEnhancedAuthService_RefreshTokenWithSession_RevokedSession(t *testing.T) {
	// Arrange
	service, sessionStore, _ := setupTestService(t)
	user := createTestUser()

	// Create a revoked session
	sessionID := "session123"
	refreshToken := "refresh_token_123"
	session := &SessionInfo{
		SessionID:     sessionID,
		UserID:        user.ID.Hex(),
		RefreshTokens: []string{refreshToken},
		Revoked:       true,
	}
	sessionStore.sessions[sessionID] = session

	// Act
	_, err := service.RefreshTokenWithSession(context.Background(), refreshToken, sessionID)

	// Assert
	if err == nil {
		t.Error("Expected error for revoked session")
	}

	expectedError := "session has been revoked"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEnhancedAuthService_RefreshTokenWithSession_ExpiredSession(t *testing.T) {
	// Arrange
	service, sessionStore, _ := setupTestService(t)
	user := createTestUser()

	// Create an expired session
	sessionID := "session123"
	refreshToken := "refresh_token_123"
	session := &SessionInfo{
		SessionID:     sessionID,
		UserID:        user.ID.Hex(),
		RefreshTokens: []string{refreshToken},
		ExpiresAt:     time.Now().Add(-1 * time.Hour), // Expired
		Revoked:       false,
	}
	sessionStore.sessions[sessionID] = session

	// Act
	_, err := service.RefreshTokenWithSession(context.Background(), refreshToken, sessionID)

	// Assert
	if err == nil {
		t.Error("Expected error for expired session")
	}

	expectedError := "session has expired"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEnhancedAuthService_RevokeSession_Success(t *testing.T) {
	// Arrange
	service, sessionStore, _ := setupTestService(t)
	sessionID := "session123"

	session := &SessionInfo{
		SessionID: sessionID,
		UserID:    "user123",
		Revoked:   false,
	}
	sessionStore.sessions[sessionID] = session

	// Act
	err := service.RevokeSession(context.Background(), sessionID)

	// Assert
	if err != nil {
		t.Fatalf("Expected successful session revocation, got error: %v", err)
	}

	// Verify session was revoked
	revokedSession := sessionStore.sessions[sessionID]
	if !revokedSession.Revoked {
		t.Error("Expected session to be marked as revoked")
	}
}

func TestEnhancedAuthService_RevokeAllSessions_Success(t *testing.T) {
	// Arrange
	service, sessionStore, _ := setupTestService(t)
	userID := "user123"

	// Create multiple sessions for the user
	session1 := &SessionInfo{SessionID: "session1", UserID: userID, Revoked: false}
	session2 := &SessionInfo{SessionID: "session2", UserID: userID, Revoked: false}
	session3 := &SessionInfo{SessionID: "session3", UserID: "otheruser", Revoked: false}

	sessionStore.sessions["session1"] = session1
	sessionStore.sessions["session2"] = session2
	sessionStore.sessions["session3"] = session3

	// Act
	err := service.RevokeAllSessions(context.Background(), userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected successful session revocation, got error: %v", err)
	}

	// Verify user's sessions were revoked
	if !sessionStore.sessions["session1"].Revoked {
		t.Error("Expected session1 to be revoked")
	}

	if !sessionStore.sessions["session2"].Revoked {
		t.Error("Expected session2 to be revoked")
	}

	// Verify other user's session was not affected
	if sessionStore.sessions["session3"].Revoked {
		t.Error("Expected session3 to remain active")
	}
}

func TestEnhancedAuthService_GetUserSessions_Success(t *testing.T) {
	// Arrange
	service, sessionStore, _ := setupTestService(t)
	userID := "user123"

	// Create sessions for the user
	activeSession := &SessionInfo{
		SessionID: "active_session",
		UserID:    userID,
		Revoked:   false,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	revokedSession := &SessionInfo{
		SessionID: "revoked_session",
		UserID:    userID,
		Revoked:   true,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	sessionStore.sessions["active_session"] = activeSession
	sessionStore.sessions["revoked_session"] = revokedSession

	// Act
	sessions, err := service.GetUserSessions(context.Background(), userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected successful session retrieval, got error: %v", err)
	}

	if len(sessions) != 1 {
		t.Errorf("Expected 1 active session, got %d", len(sessions))
	}

	if sessions[0].SessionID != "active_session" {
		t.Errorf("Expected active session, got %s", sessions[0].SessionID)
	}
}

func TestDeviceTrustCalculation(t *testing.T) {
	tests := []struct {
		name          string
		device        DeviceFingerprint
		expectedScore float64
		expectedTrust bool
	}{
		{
			name: "Trusted iOS device",
			device: DeviceFingerprint{
				Platform:        "iOS",
				IsJailbroken:    false,
				AttestationData: "official",
			},
			expectedScore: 1.0,
			expectedTrust: true,
		},
		{
			name: "Jailbroken device",
			device: DeviceFingerprint{
				Platform:        "iOS",
				IsJailbroken:    true,
				AttestationData: "official",
			},
			expectedScore: 0.8,   // 1.0 - 0.5 + 0.2 + 0.1 = 0.8
			expectedTrust: false, // 0.8 is not > 0.8
		},
		{
			name: "Unknown platform",
			device: DeviceFingerprint{
				Platform:        "Unknown",
				IsJailbroken:    false,
				AttestationData: "",
			},
			expectedScore: 0.8,   // 1.0 - 0.2 = 0.8
			expectedTrust: false, // 0.8 is not > 0.8
		},
	}

	service, _, _ := setupTestService(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.calculateTrustScore(&tt.device)

			if score-tt.expectedScore > 0.001 || tt.expectedScore-score > 0.001 {
				t.Errorf("Expected trust score %f, got %f", tt.expectedScore, score)
			}

			trusted := score > 0.8
			if trusted != tt.expectedTrust {
				t.Errorf("Expected trusted %v, got %v", tt.expectedTrust, trusted)
			}
		})
	}
}
