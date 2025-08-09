package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smorting/backend/internal/models"
	"go.uber.org/zap"
)

// JWTRefreshService handles JWT token refresh with 30-minute access tokens
type JWTRefreshService struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	logger             *zap.Logger
	revocationStore    TokenRevocationStore
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// AccessTokenClaims represents claims for access tokens (30 minutes)
type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents claims for refresh tokens (7 days)
type RefreshTokenClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	TokenID string `json:"token_id"` // Unique identifier for refresh token
	jwt.RegisteredClaims
}

// NewJWTRefreshService creates a new JWT refresh service
func NewJWTRefreshService(accessSecret, refreshSecret []byte, logger *zap.Logger) *JWTRefreshService {
	svc := &JWTRefreshService{
		accessTokenSecret:  accessSecret,
		refreshTokenSecret: refreshSecret,
		logger:             logger,
	}
	// Default to in-memory store; production should override with persistent store
	svc.revocationStore = NewMemoryRevocationStore()
	return svc
}

// SetRevocationStore injects a persistent revocation store (e.g., Mongo-backed)
func (j *JWTRefreshService) SetRevocationStore(store TokenRevocationStore) {
	if store != nil {
		j.revocationStore = store
	}
}

// GenerateTokenPair generates access and refresh tokens
func (j *JWTRefreshService) GenerateTokenPair(user *models.User) (*TokenPair, error) {
	now := time.Now()

	// Generate unique token ID for refresh token
	tokenID, err := j.generateTokenID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token ID: %w", err)
	}

	// Access token - 30 minutes
	accessClaims := &AccessTokenClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "smor-ting-backend",
			Subject:   user.ID.Hex(),
			ID:        tokenID,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.accessTokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token - 7 days
	refreshClaims := &RefreshTokenClaims{
		UserID:  user.ID.Hex(),
		Email:   user.Email,
		Role:    string(user.Role),
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "smor-ting-backend",
			Subject:   user.ID.Hex(),
			ID:        tokenID,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.refreshTokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	j.logger.Info("Generated token pair",
		zap.String("user_id", user.ID.Hex()),
		zap.String("token_id", tokenID),
	)

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    30 * 60, // 30 minutes in seconds
	}, nil
}

// ValidateAccessToken validates an access token
func (j *JWTRefreshService) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.accessTokenSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		// Check revocation by token ID (JTI)
		if j.revocationStore != nil && claims.ID != "" {
			revoked, rerr := j.revocationStore.IsRevoked(claims.ID)
			if rerr != nil {
				return nil, fmt.Errorf("revocation check failed: %w", rerr)
			}
			if revoked {
				return nil, fmt.Errorf("token has been revoked")
			}
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid access token claims")
}

// ValidateRefreshToken validates a refresh token
func (j *JWTRefreshService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.refreshTokenSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		// Check revocation by TokenID
		if j.revocationStore != nil && claims.TokenID != "" {
			revoked, rerr := j.revocationStore.IsRevoked(claims.TokenID)
			if rerr != nil {
				return nil, fmt.Errorf("revocation check failed: %w", rerr)
			}
			if revoked {
				return nil, fmt.Errorf("token has been revoked")
			}
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid refresh token claims")
}

// RefreshAccessToken generates a new access token using a valid refresh token
func (j *JWTRefreshService) RefreshAccessToken(refreshTokenString string, user *models.User) (*TokenPair, error) {
	// Validate refresh token
	refreshClaims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Verify user ID matches
	if refreshClaims.UserID != user.ID.Hex() {
		return nil, fmt.Errorf("refresh token user ID mismatch")
	}

	// Generate new token pair
	return j.GenerateTokenPair(user)
}

// RevokeRefreshToken marks a refresh token as revoked (implement with database)
func (j *JWTRefreshService) RevokeRefreshToken(tokenID string) error {
	if tokenID == "" {
		return fmt.Errorf("tokenID is required")
	}
	if j.revocationStore == nil {
		return fmt.Errorf("revocation store not configured")
	}
	// Default expiration horizon for revocation records if unknown: 400 days
	expiresAt := time.Now().Add(400 * 24 * time.Hour)
	if err := j.revocationStore.Revoke(tokenID, expiresAt); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	j.logger.Info("Revoked token", zap.String("token_id", tokenID))
	return nil
}

// IsTokenExpired checks if a token is expired
func (j *JWTRefreshService) IsTokenExpired(tokenString string, isRefreshToken bool) (bool, error) {
	var claims jwt.Claims

	if isRefreshToken {
		claims = &RefreshTokenClaims{}
	} else {
		claims = &AccessTokenClaims{}
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if isRefreshToken {
			return j.refreshTokenSecret, nil
		}
		return j.accessTokenSecret, nil
	})

	if err != nil {
		return true, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return true, nil
	}

	return false, nil
}

// GetTokenExpiration returns the expiration time of a token
func (j *JWTRefreshService) GetTokenExpiration(tokenString string, isRefreshToken bool) (*time.Time, error) {
	var claims jwt.Claims

	if isRefreshToken {
		claims = &RefreshTokenClaims{}
	} else {
		claims = &AccessTokenClaims{}
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if isRefreshToken {
			return j.refreshTokenSecret, nil
		}
		return j.accessTokenSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract expiration time
	if refreshClaims, ok := claims.(*RefreshTokenClaims); ok {
		return &refreshClaims.ExpiresAt.Time, nil
	}

	if accessClaims, ok := claims.(*AccessTokenClaims); ok {
		return &accessClaims.ExpiresAt.Time, nil
	}

	return nil, fmt.Errorf("failed to extract expiration time")
}

// generateTokenID generates a unique token ID
func (j *JWTRefreshService) generateTokenID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetTokenInfo returns information about a token
func (j *JWTRefreshService) GetTokenInfo(tokenString string, isRefreshToken bool) (map[string]interface{}, error) {
	var claims jwt.Claims

	if isRefreshToken {
		claims = &RefreshTokenClaims{}
	} else {
		claims = &AccessTokenClaims{}
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if isRefreshToken {
			return j.refreshTokenSecret, nil
		}
		return j.accessTokenSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	info := map[string]interface{}{
		"valid": token.Valid,
	}

	// Safely extract claims information
	switch c := token.Claims.(type) {
	case *RefreshTokenClaims:
		info["expires_at"] = c.ExpiresAt
		info["issued_at"] = c.IssuedAt
		info["issuer"] = c.Issuer
		info["subject"] = c.Subject
	case *AccessTokenClaims:
		info["expires_at"] = c.ExpiresAt
		info["issued_at"] = c.IssuedAt
		info["issuer"] = c.Issuer
		info["subject"] = c.Subject
	}

	if refreshClaims, ok := claims.(*RefreshTokenClaims); ok {
		info["user_id"] = refreshClaims.UserID
		info["email"] = refreshClaims.Email
		info["role"] = refreshClaims.Role
		info["token_id"] = refreshClaims.TokenID
	} else if accessClaims, ok := claims.(*AccessTokenClaims); ok {
		info["user_id"] = accessClaims.UserID
		info["email"] = accessClaims.Email
		info["role"] = accessClaims.Role
	}

	return info, nil
}
