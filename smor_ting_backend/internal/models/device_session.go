package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeviceSession represents an active user session on a specific device
type DeviceSession struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	DeviceID     string             `json:"device_id" bson:"device_id"`     // Unique device identifier
	DeviceName   string             `json:"device_name" bson:"device_name"` // Human-readable device name
	Platform     string             `json:"platform" bson:"platform"`       // ios, android, web
	IPAddress    string             `json:"ip_address" bson:"ip_address"`
	UserAgent    string             `json:"user_agent" bson:"user_agent"`
	RefreshToken string             `json:"-" bson:"refresh_token"` // Current refresh token
	IsActive     bool               `json:"is_active" bson:"is_active"`
	LastActivity time.Time          `json:"last_activity" bson:"last_activity"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	RevokedAt    *time.Time         `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
	// Security metadata
	Location *GeoLocation `json:"location,omitempty" bson:"location,omitempty"`
	// Biometric security (for mobile)
	BiometricEnabled bool   `json:"biometric_enabled" bson:"biometric_enabled"`
	BiometricType    string `json:"biometric_type,omitempty" bson:"biometric_type,omitempty"` // touch_id, face_id, fingerprint
}

// GeoLocation represents approximate location information for security
type GeoLocation struct {
	Country string  `json:"country" bson:"country"`
	Region  string  `json:"region" bson:"region"`
	City    string  `json:"city" bson:"city"`
	Lat     float64 `json:"lat,omitempty" bson:"lat,omitempty"`
	Lon     float64 `json:"lon,omitempty" bson:"lon,omitempty"`
}

// RevokeSession marks the session as inactive and records revocation time
func (ds *DeviceSession) RevokeSession() {
	ds.IsActive = false
	now := time.Now()
	ds.RevokedAt = &now
}

// IsExpired checks if the session has expired based on the given duration
func (ds *DeviceSession) IsExpired(maxAge time.Duration) bool {
	return time.Since(ds.LastActivity) > maxAge
}

// UpdateActivity updates the last activity timestamp
func (ds *DeviceSession) UpdateActivity() {
	ds.LastActivity = time.Now()
}

// SecurityEventType represents different types of security events
type SecurityEventType string

const (
	PasswordChangeEvent  SecurityEventType = "password_change"
	TwoFactorChangeEvent SecurityEventType = "two_factor_change"
	LoginEvent           SecurityEventType = "login"
	LogoutEvent          SecurityEventType = "logout"
	TokenRefreshEvent    SecurityEventType = "token_refresh"
	SuspiciousLoginEvent SecurityEventType = "suspicious_login"
	DeviceRegistered     SecurityEventType = "device_registered"
	DeviceRevoked        SecurityEventType = "device_revoked"
	BiometricEnabled     SecurityEventType = "biometric_enabled"
	BiometricDisabled    SecurityEventType = "biometric_disabled"
)

// SecurityEvent represents a security-related event for audit trails
type SecurityEvent struct {
	ID        primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID     `json:"user_id" bson:"user_id"`
	EventType SecurityEventType      `json:"event_type" bson:"event_type"`
	IPAddress string                 `json:"ip_address" bson:"ip_address"`
	UserAgent string                 `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	DeviceID  string                 `json:"device_id,omitempty" bson:"device_id,omitempty"`
	SessionID *primitive.ObjectID    `json:"session_id,omitempty" bson:"session_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp" bson:"timestamp"`
	// Risk assessment
	RiskScore int    `json:"risk_score,omitempty" bson:"risk_score,omitempty"` // 0-100
	RiskLevel string `json:"risk_level,omitempty" bson:"risk_level,omitempty"` // low, medium, high
}

// DeviceSessionListResponse represents the response for listing user device sessions
type DeviceSessionListResponse struct {
	Sessions      []DeviceSession `json:"sessions"`
	CurrentDevice *DeviceSession  `json:"current_device,omitempty"`
	TotalSessions int             `json:"total_sessions"`
}

// TokenRotationRequest represents a request to rotate refresh tokens
type TokenRotationRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
}

// BiometricUnlockRequest represents a request to unlock with biometrics
type BiometricUnlockRequest struct {
	DeviceID          string `json:"device_id" validate:"required"`
	BiometricType     string `json:"biometric_type" validate:"required"`
	BiometricData     string `json:"biometric_data" validate:"required"` // Encrypted biometric signature
	Challenge         string `json:"challenge" validate:"required"`      // Server-provided challenge
	ChallengeResponse string `json:"challenge_response" validate:"required"`
}

// SecureTokenStorage represents secure token storage configuration for mobile
type SecureTokenStorage struct {
	Platform         string `json:"platform"` // ios, android
	BiometricEnabled bool   `json:"biometric_enabled"`
	BiometricType    string `json:"biometric_type"`   // touch_id, face_id, fingerprint
	KeychainEnabled  bool   `json:"keychain_enabled"` // iOS Keychain
	KeystoreEnabled  bool   `json:"keystore_enabled"` // Android Keystore
	EncryptionLevel  string `json:"encryption_level"` // hardware, software
}
