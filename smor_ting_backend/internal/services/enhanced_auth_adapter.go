package services

import (
	"context"
	"fmt"

	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// EnhancedAuthServiceAdapter adapts the internal EnhancedAuthService to the handler's expected interface
type EnhancedAuthServiceAdapter struct {
	service *EnhancedAuthService
	repo    database.Repository
	authSvc *auth.MongoDBService
	logger  *zap.Logger
}

func NewEnhancedAuthServiceAdapter(service *EnhancedAuthService, repo database.Repository, authSvc *auth.MongoDBService, logger *zap.Logger) *EnhancedAuthServiceAdapter {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &EnhancedAuthServiceAdapter{service: service, repo: repo, authSvc: authSvc, logger: logger}
}

// EnhancedLogin performs full enhanced login and returns models.EnhancedAuthResult
func (a *EnhancedAuthServiceAdapter) EnhancedLogin(req *models.EnhancedLoginRequest, clientIP string) (*models.EnhancedAuthResult, error) {
	ctx := context.Background()

	// Load user by email
	user, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return &models.EnhancedAuthResult{Success: false, Message: "Invalid email or password"}, nil
	}

	// Verify password using bcrypt against stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &models.EnhancedAuthResult{Success: false, Message: "Invalid email or password"}, nil
	}

	// Map device info (handle nil)
	var devInfo DeviceFingerprint
	if req.DeviceInfo != nil {
		devInfo = DeviceFingerprint{
			DeviceID:        req.DeviceInfo.DeviceID,
			Platform:        req.DeviceInfo.Platform,
			OSVersion:       req.DeviceInfo.OSVersion,
			AppVersion:      req.DeviceInfo.AppVersion,
			IsTrusted:       false,
			IsJailbroken:    false,
			TrustScore:      0,
			LastVerified:    user.CreatedAt, // placeholder, will be updated by service
			AttestationData: "",
		}
	}

	// Build auth request for real service
	authReq := &AuthRequest{
		Email:         req.Email,
		Password:      req.Password,
		RememberMe:    req.RememberMe,
		DeviceInfo:    devInfo,
		IPAddress:     clientIP,
		UserAgent:     req.UserAgent,
		CaptchaToken:  req.CaptchaToken,
		TwoFactorCode: req.TwoFactorCode,
	}

	result, err := a.service.Authenticate(ctx, authReq, user)
	if err != nil {
		return nil, err
	}

	// Map to models.EnhancedAuthResult
	return &models.EnhancedAuthResult{
		Success:              true,
		Message:              "Login successful",
		User:                 result.User,
		AccessToken:          result.AccessToken,
		RefreshToken:         result.RefreshToken,
		SessionID:            result.SessionID,
		TokenExpiresAt:       &result.TokenExpiresAt,
		RefreshExpiresAt:     &result.RefreshExpiresAt,
		RequiresTwoFactor:    false, // OTP/2FA DISABLED - always return false
		RequiresVerification: result.RequiresVerification,
		RequiresCaptcha:      false,
		DeviceTrusted:        result.DeviceTrusted,
		IsRestoredSession:    false,
	}, nil
}

func (a *EnhancedAuthServiceAdapter) BiometricLogin(sessionID, biometricData string) (*models.EnhancedAuthResult, error) {
	ctx := context.Background()
	// Load session and derive user
	session, err := a.service.sessionStore.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session: %w", err)
	}
	oid, err := primitive.ObjectIDFromHex(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in session: %w", err)
	}
	user, err := a.repo.GetUserByID(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	ar, err := a.service.GenerateTokensForExistingSession(ctx, session, user)
	if err != nil {
		return nil, err
	}
	return &models.EnhancedAuthResult{
		Success:           true,
		Message:           "Biometric login successful",
		User:              ar.User,
		AccessToken:       ar.AccessToken,
		RefreshToken:      ar.RefreshToken,
		SessionID:         ar.SessionID,
		TokenExpiresAt:    &ar.TokenExpiresAt,
		RefreshExpiresAt:  &ar.RefreshExpiresAt,
		DeviceTrusted:     ar.DeviceTrusted,
		IsRestoredSession: true,
	}, nil
}

func (a *EnhancedAuthServiceAdapter) GetUserSessions(userID string) ([]*models.SessionInfo, error) {
	ctx := context.Background()
	sessions, err := a.service.sessionStore.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, err
	}
	// Map to models.SessionInfo
	result := make([]*models.SessionInfo, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, &models.SessionInfo{
			SessionID:    s.SessionID,
			UserID:       s.UserID,
			DeviceName:   s.DeviceInfo.Platform,
			DeviceType:   s.DeviceInfo.Platform,
			IPAddress:    s.IPAddress,
			UserAgent:    s.UserAgent,
			LastActivity: s.LastActivity,
			CreatedAt:    s.CreatedAt,
			ExpiresAt:    s.ExpiresAt,
			IsCurrent:    !s.Revoked,
			IsRevoked:    s.Revoked,
		})
	}
	return result, nil
}

func (a *EnhancedAuthServiceAdapter) GetSessionByID(sessionID string) (*models.SessionInfo, error) {
	ctx := context.Background()
	s, err := a.service.sessionStore.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return &models.SessionInfo{
		SessionID:    s.SessionID,
		UserID:       s.UserID,
		DeviceName:   s.DeviceInfo.Platform,
		DeviceType:   s.DeviceInfo.Platform,
		IPAddress:    s.IPAddress,
		UserAgent:    s.UserAgent,
		LastActivity: s.LastActivity,
		CreatedAt:    s.CreatedAt,
		ExpiresAt:    s.ExpiresAt,
		IsCurrent:    !s.Revoked,
		IsRevoked:    s.Revoked,
	}, nil
}

func (a *EnhancedAuthServiceAdapter) RevokeSession(sessionID string) error {
	ctx := context.Background()
	return a.service.sessionStore.RevokeSession(ctx, sessionID)
}

func (a *EnhancedAuthServiceAdapter) RevokeAllSessions(userID string) error {
	ctx := context.Background()
	return a.service.sessionStore.RevokeAllUserSessions(ctx, userID)
}

func (a *EnhancedAuthServiceAdapter) SignOutAllDevices(userID string) error {
	// Same as revoking all sessions in this context
	return a.RevokeAllSessions(userID)
}

func (a *EnhancedAuthServiceAdapter) RefreshTokenWithSession(refreshToken, sessionID string) (*models.EnhancedAuthResult, error) {
	ctx := context.Background()
	ar, err := a.service.RefreshTokenWithSession(ctx, refreshToken, sessionID)
	if err != nil {
		return nil, err
	}
	return &models.EnhancedAuthResult{
		Success:           true,
		Message:           "Token refreshed",
		User:              ar.User,
		AccessToken:       ar.AccessToken,
		RefreshToken:      ar.RefreshToken,
		SessionID:         ar.SessionID,
		TokenExpiresAt:    &ar.TokenExpiresAt,
		RefreshExpiresAt:  &ar.RefreshExpiresAt,
		DeviceTrusted:     ar.DeviceTrusted,
		IsRestoredSession: true,
	}, nil
}

func (a *EnhancedAuthServiceAdapter) VerifyDeviceFingerprint(existing, provided *models.DeviceFingerprint) bool {
	if existing == nil || provided == nil {
		return false
	}
	// Basic comparison by device ID and core platform info
	return existing.DeviceID == provided.DeviceID || (existing.Platform == provided.Platform && existing.OSVersion == provided.OSVersion && existing.AppVersion == provided.AppVersion)
}

func (a *EnhancedAuthServiceAdapter) GenerateTokensForExistingSession(sessionID string) (*models.EnhancedAuthResult, error) {
	ctx := context.Background()
	session, err := a.service.sessionStore.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	// Load user
	oid, err := primitive.ObjectIDFromHex(session.UserID)
	if err != nil {
		return nil, err
	}
	user, err := a.repo.GetUserByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	ar, err := a.service.GenerateTokensForExistingSession(ctx, session, user)
	if err != nil {
		return nil, err
	}
	return &models.EnhancedAuthResult{
		Success:           true,
		Message:           "Tokens generated",
		User:              ar.User,
		AccessToken:       ar.AccessToken,
		RefreshToken:      ar.RefreshToken,
		SessionID:         ar.SessionID,
		TokenExpiresAt:    &ar.TokenExpiresAt,
		RefreshExpiresAt:  &ar.RefreshExpiresAt,
		DeviceTrusted:     ar.DeviceTrusted,
		IsRestoredSession: true,
	}, nil
}

func (a *EnhancedAuthServiceAdapter) UpdateSessionActivity(sessionID string) error {
	// Best-effort; no IP/UserAgent context here
	return a.service.UpdateSessionActivity(context.Background(), sessionID, "", "")
}
