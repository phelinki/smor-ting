package handlers

import (
	"fmt"
	"net/http"
	"os"

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
	SessionID    string `json:"session_id,omitempty"` // Optional for backward compatibility
}

// RefreshTokenResponse represents a token refresh response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Register handles user registration with enhanced security
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse register request body", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse request body",
		})
	}

	// Validate request
	if err := h.validateRegisterRequest(&req); err != nil {
		h.logger.Warn("Invalid register request", zap.String("error", err.Error()))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": err.Error(),
		})
	}

	// Register user through MongoDB service
	response, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to register user", err, zap.String("email", req.Email), zap.String("error_details", err.Error()))

		// Handle specific errors
		if err.Error() == "user with email "+req.Email+" already exists" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error":   "User already exists",
				"message": "A user with this email already exists",
			})
		}

		// For debugging: return more detailed error information
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Registration failed",
			"message":     "Failed to register user",
			"debug_error": err.Error(), // Add this temporarily for debugging
		})
	}

	// If registration successful, generate enhanced token pair
	if response.AccessToken != "" {
		// User is already authenticated, generate new token pair with JWT service
		tokenPair, err := h.jwtService.GenerateTokenPair(&response.User)
		if err != nil {
			h.logger.Error("Failed to generate token pair after registration", err)
			// Continue with the original token from auth service
		} else {
			// Use enhanced tokens
			response.AccessToken = tokenPair.AccessToken
			response.RefreshToken = tokenPair.RefreshToken
		}
	}

	h.logger.Info("User registered successfully", zap.String("email", req.Email))
	return c.Status(http.StatusCreated).JSON(response)
}

// validateRegisterRequest validates a registration request
func (h *AuthHandler) validateRegisterRequest(req *models.RegisterRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if req.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if req.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if req.Phone == "" {
		return fmt.Errorf("phone is required")
	}
	if req.Role == "" {
		return fmt.Errorf("role is required")
	}
	// Validate role is one of the allowed values
	if req.Role != models.CustomerRole && req.Role != models.ProviderRole && req.Role != models.AdminRole {
		return fmt.Errorf("role must be 'customer', 'provider', or 'admin'")
	}
	return nil
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

	// Validate request fields
	if req.Email == "" {
		h.logger.Warn("Login validation failed", zap.String("error", "email is required"))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": "email is required",
		})
	}

	if req.Password == "" {
		h.logger.Warn("Login validation failed", zap.String("error", "password is required"))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": "password is required",
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

// ResendOTP triggers a new OTP email for a given email address when appropriate
func (h *AuthHandler) ResendOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		// Optional purpose; defaults decided by server based on user status
		Purpose string `json:"purpose"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// Determine purpose: if user exists and not verified -> login/registration; else accept provided
	purpose := req.Purpose
	if purpose == "" {
		if u, err := h.authService.GetUserByEmail(c.Context(), req.Email); err == nil && u != nil {
			if !u.IsEmailVerified {
				purpose = "login"
			} else {
				purpose = "login"
			}
		} else {
			purpose = "registration"
		}
	}
	// Create OTP via repository through service by calling CreateOTP on service path
	// Generate using underlying service's generator flow: call CreateOTP directly
	if err := h.authService.CreateOTP(c.Context(), req.Email, purpose); err != nil {
		h.logger.Error("Failed to create OTP", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create OTP"})
	}
	// Attempt email via EmailService
	emailSvc := services.NewEmailService()
	// We cannot retrieve the plaintext OTP from service; for production we would send within service itself.
	// For now, respond success without disclosing OTP; email service in Mongo service should handle sending.
	_ = emailSvc // placeholder; sending is handled in service if wired
	return c.JSON(fiber.Map{"message": "OTP sent if the email exists"})
}

// VerifyOTP verifies the OTP and returns tokens
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req models.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil || req.Email == "" || req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.authService.VerifyOTP(c.Context(), req.Email, req.OTP); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	}
	// Load user and issue token pair via JWT service
	user, err := h.authService.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}
	pair, err := h.jwtService.GenerateTokenPair(user)
	if err != nil {
		h.logger.Error("Failed to generate token pair", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create tokens"})
	}
	return c.JSON(models.AuthResponse{User: *user, AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken, RequiresOTP: false})
}

// TestGetLatestOTP exposes latest OTP for a given email in non-production environments only
func (h *AuthHandler) TestGetLatestOTP(c *fiber.Ctx) error {
	// Only allow in development or when explicitly enabled
	if os.Getenv("ENV") == "production" && os.Getenv("ENABLE_TEST_ENDPOINTS") != "true" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Not allowed"})
	}
	email := c.Query("email")
	if email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}
	otp, err := h.authService.GetLatestOTPByEmail(c.Context(), email)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "OTP not found"})
	}
	return c.JSON(fiber.Map{"otp": otp.OTP, "expires_at": otp.ExpiresAt})
}

// RequestPasswordReset creates an OTP and emails it for password reset
func (h *AuthHandler) RequestPasswordReset(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.authService.RequestPasswordReset(c.Context(), req.Email); err != nil {
		// Still return 200 to avoid enumeration; log error
		h.logger.Warn("Password reset request error", zap.String("email", req.Email))
	}
	return c.JSON(fiber.Map{"message": "If the email exists, a reset code has been sent"})
}

// ResetPassword verifies OTP and updates password
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" || req.OTP == "" || len(req.NewPassword) < 6 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.authService.ResetPassword(c.Context(), req.Email, req.OTP, req.NewPassword); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	}
	return c.JSON(fiber.Map{"message": "Password reset successful"})
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

	// Fetch complete user data from database
	userObjectID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		h.logger.Error("Failed to parse user ID from token", err, zap.String("user_id", claims.UserID))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Token contains invalid user identifier",
		})
	}

	user, err := h.authService.GetUserByID(c.Context(), userObjectID)
	if err != nil {
		h.logger.Error("Failed to fetch user data for token validation", err, zap.String("user_id", claims.UserID))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "User not found",
			"message": "User associated with token not found",
		})
	}

	// Build complete user object matching mobile app expectations
	userResponse := fiber.Map{
		"id":                user.ID.Hex(),
		"email":             user.Email,
		"first_name":        user.FirstName,
		"last_name":         user.LastName,
		"phone":             user.Phone,
		"role":              user.Role,
		"is_email_verified": user.IsEmailVerified,
		"profile_image":     user.ProfileImage,
		"created_at":        user.CreatedAt,
		"updated_at":        user.UpdatedAt,
	}

	// Add address if it exists
	if user.Address != nil {
		userResponse["address"] = fiber.Map{
			"street":    user.Address.Street,
			"city":      user.Address.City,
			"county":    user.Address.County,
			"country":   user.Address.Country,
			"latitude":  user.Address.Latitude,
			"longitude": user.Address.Longitude,
		}
	} else {
		// Provide empty address to match Flutter model expectations
		userResponse["address"] = fiber.Map{
			"street":    "",
			"city":      "",
			"county":    "",
			"country":   "",
			"latitude":  0.0,
			"longitude": 0.0,
		}
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
		"message":    "Token is valid",
		"user":       userResponse, // Return complete user data
		"token_info": tokenInfo,
		"permissions": fiber.Map{
			"is_customer": claims.Role == string(models.CustomerRole),
			"is_provider": claims.Role == string(models.ProviderRole),
			"is_admin":    claims.Role == string(models.AdminRole),
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
