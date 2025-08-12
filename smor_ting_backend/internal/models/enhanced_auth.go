package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EnhancedLoginRequest represents an enhanced login request with additional security features
type EnhancedLoginRequest struct {
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=6"`
	DeviceFingerprint string `json:"device_fingerprint,omitempty"`
	ClientIP         string `json:"client_ip,omitempty"`
	UserAgent        string `json:"user_agent,omitempty"`
	CaptchaToken     string `json:"captcha_token,omitempty"`
	TwoFactorCode    string `json:"two_factor_code,omitempty"`
	BiometricData    string `json:"biometric_data,omitempty"`
	SessionID        string `json:"session_id,omitempty"`
}

// EnhancedAuthResult represents the result of an enhanced authentication attempt
type EnhancedAuthResult struct {
	Success              bool       `json:"success"`
	Message              string     `json:"message,omitempty"`
	User                 *User      `json:"user,omitempty"`
	AccessToken          string     `json:"access_token,omitempty"`
	RefreshToken         string     `json:"refresh_token,omitempty"`
	SessionID            string     `json:"session_id,omitempty"`
	TokenExpiresAt       *time.Time `json:"token_expires_at,omitempty"`
	RefreshExpiresAt     *time.Time `json:"refresh_expires_at,omitempty"`
	RequiresTwoFactor    bool       `json:"requires_two_factor"`
	RequiresVerification bool       `json:"requires_verification"`
	RequiresCaptcha      bool       `json:"requires_captcha"`
	DeviceTrusted        bool       `json:"device_trusted"`
	IsRestoredSession    bool       `json:"is_restored_session"`
	RemainingAttempts    *int       `json:"remaining_attempts,omitempty"`
	LockoutInfo          *LockoutInfo `json:"lockout_info,omitempty"`
}

// SessionInfo represents information about a user session
type SessionInfo struct {
	SessionID    string    `json:"session_id" bson:"session_id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	DeviceName   string    `json:"device_name" bson:"device_name"`
	DeviceType   string    `json:"device_type" bson:"device_type"`
	IPAddress    string    `json:"ip_address" bson:"ip_address"`
	UserAgent    string    `json:"user_agent" bson:"user_agent"`
	LastActivity time.Time `json:"last_activity" bson:"last_activity"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt    time.Time `json:"expires_at" bson:"expires_at"`
	IsCurrent    bool      `json:"is_current" bson:"is_current"`
	IsRevoked    bool      `json:"is_revoked" bson:"is_revoked"`
}

// LockoutInfo represents information about account lockout
type LockoutInfo struct {
	LockedUntil       time.Time `json:"locked_until" bson:"locked_until"`
	RemainingAttempts int       `json:"remaining_attempts" bson:"remaining_attempts"`
	TimeUntilUnlock   *int      `json:"time_until_unlock,omitempty" bson:"time_until_unlock,omitempty"`
	LockoutReason     string    `json:"lockout_reason,omitempty" bson:"lockout_reason,omitempty"`
}

// BiometricLoginRequest represents a biometric login request
type BiometricLoginRequest struct {
	SessionID     string `json:"session_id" validate:"required"`
	BiometricData string `json:"biometric_data" validate:"required"`
	DeviceID      string `json:"device_id,omitempty"`
}

// Session represents a user session in the database
type Session struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SessionID    string             `json:"session_id" bson:"session_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	DeviceName   string             `json:"device_name" bson:"device_name"`
	DeviceType   string             `json:"device_type" bson:"device_type"`
	IPAddress    string             `json:"ip_address" bson:"ip_address"`
	UserAgent    string             `json:"user_agent" bson:"user_agent"`
	RefreshToken string             `json:"-" bson:"refresh_token"` // Hidden from JSON
	LastActivity time.Time          `json:"last_activity" bson:"last_activity"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt    time.Time          `json:"expires_at" bson:"expires_at"`
	IsRevoked    bool               `json:"is_revoked" bson:"is_revoked"`
}

// GetSessionsResponse represents the response for getting user sessions
type GetSessionsResponse struct {
	Sessions []*SessionInfo `json:"sessions"`
	Current  *SessionInfo   `json:"current"`
}

// RevokeSessionRequest represents a request to revoke a session
type RevokeSessionRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}

// SignOutAllDevicesRequest represents a request to sign out from all devices
type SignOutAllDevicesRequest struct {
	UserID string `json:"user_id" validate:"required"`
}
