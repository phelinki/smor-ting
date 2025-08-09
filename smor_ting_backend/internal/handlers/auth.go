package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// AuthHandler handles authentication requests with enhanced security
type AuthHandler struct {
	jwtService        *services.JWTRefreshService
	encryptionService *services.EncryptionService
	logger            *logger.Logger
	authService       *auth.MongoDBService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(jwtService *services.JWTRefreshService, encryptionService *services.EncryptionService, logger *logger.Logger, authService *auth.MongoDBService) *AuthHandler {
	return &AuthHandler{
		jwtService:        jwtService,
		encryptionService: encryptionService,
		logger:            logger,
		authService:       authService,
	}
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents a token refresh response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Login handles user login with enhanced security
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse login request", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse request body",
		})
	}

	// Authenticate against MongoDB-backed service
	user, err := h.authService.Authenticate(c.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Invalid credentials", zap.String("email", req.Email))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid credentials",
			"message": "Invalid email or password",
		})
	}

	// Generate token pair with 30-minute access token
	tokenPair, err := h.jwtService.GenerateTokenPair(user)
	if err != nil {
		h.logger.Error("Failed to generate token pair", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Authentication failed",
			"message": "Failed to generate authentication tokens",
		})
	}

	h.logger.Info("User logged in successfully",
		zap.String("email", req.Email),
		zap.String("user_id", user.ID.Hex()),
	)

	// Return response matching mobile AuthResponse
	return c.JSON(models.AuthResponse{
		User:         *user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		RequiresOTP:  false,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse refresh token request", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse request body",
		})
	}

	// Validate refresh token
	refreshClaims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Warn("Invalid refresh token", zap.Error(err))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid refresh token",
			"message": "The refresh token is invalid or expired",
		})
	}

	// Create user with proper ObjectID for token generation
	userID, _ := primitive.ObjectIDFromHex(refreshClaims.UserID)
	user := &models.User{
		ID:    userID,
		Email: refreshClaims.Email,
		Role:  models.UserRole(refreshClaims.Role),
	}

	// Generate new token pair
	tokenPair, err := h.jwtService.RefreshAccessToken(req.RefreshToken, user)
	if err != nil {
		h.logger.Error("Failed to refresh access token", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Token refresh failed",
			"message": "Failed to generate new access token",
		})
	}

	h.logger.Info("Token refreshed successfully",
		zap.String("user_id", user.ID.Hex()),
		zap.String("token_id", refreshClaims.TokenID),
	)

	return c.JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"data":    tokenPair,
	})
}

// ValidateToken validates an access token
func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Missing token",
			"message": "Authorization header is required",
		})
	}

	// Extract token from "Bearer <token>" format
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid token format",
			"message": "Token must be in 'Bearer <token>' format",
		})
	}

	// Validate access token
	claims, err := h.jwtService.ValidateAccessToken(token)
	if err != nil {
		h.logger.Warn("Invalid access token", zap.Error(err))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid token",
			"message": "The access token is invalid or expired",
		})
	}

	// Get token information
	tokenInfo, err := h.jwtService.GetTokenInfo(token, false)
	if err != nil {
		h.logger.Error("Failed to get token info", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Token validation failed",
			"message": "Failed to retrieve token information",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Token is valid",
		"data": fiber.Map{
			"user_id":    claims.UserID,
			"email":      claims.Email,
			"role":       claims.Role,
			"token_info": tokenInfo,
			"permissions": fiber.Map{
				"is_customer": claims.Role == string(models.CustomerRole),
				"is_provider": claims.Role == string(models.ProviderRole),
				"is_admin":    claims.Role == string(models.AdminRole),
			},
		},
	})
}

// RevokeToken revokes a refresh token
func (h *AuthHandler) RevokeToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse revoke token request", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse request body",
		})
	}

	// Validate refresh token to get token ID
	refreshClaims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Warn("Invalid refresh token for revocation", zap.Error(err))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid refresh token",
			"message": "The refresh token is invalid or expired",
		})
	}

	// Revoke the refresh token
	err = h.jwtService.RevokeRefreshToken(refreshClaims.TokenID)
	if err != nil {
		h.logger.Error("Failed to revoke refresh token", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Token revocation failed",
			"message": "Failed to revoke the refresh token",
		})
	}

	h.logger.Info("Token revoked successfully",
		zap.String("token_id", refreshClaims.TokenID),
		zap.String("user_id", refreshClaims.UserID),
	)

	return c.JSON(fiber.Map{
		"message": "Token revoked successfully",
	})
}

// GetTokenInfo returns information about a token
func (h *AuthHandler) GetTokenInfo(c *fiber.Ctx) error {
	tokenType := c.Query("type", "access") // "access" or "refresh"
	token := c.Query("token")

	if token == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing token",
			"message": "Token parameter is required",
		})
	}

	isRefreshToken := tokenType == "refresh"
	tokenInfo, err := h.jwtService.GetTokenInfo(token, isRefreshToken)
	if err != nil {
		h.logger.Warn("Failed to get token info", zap.Error(err))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid token",
			"message": "Failed to retrieve token information",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Token information retrieved",
		"data":    tokenInfo,
	})
}
