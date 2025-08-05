package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/pkg/logger"
	"go.uber.org/zap"
)

// AuthMiddleware represents authentication middleware
type AuthMiddleware struct {
	authService *auth.Service
	logger      *logger.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService *auth.Service, logger *logger.Logger) (*AuthMiddleware, error) {
	if authService == nil {
		return nil, fmt.Errorf("auth service is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}, nil
}

// Authenticate validates JWT tokens and adds user to context
func (am *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip authentication for certain paths
		if am.shouldSkipAuth(c.Path()) {
			return c.Next()
		}

		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			am.logger.Warn("Missing Authorization header", zap.String("path", c.Path()))
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
			am.logger.Warn("Invalid Authorization header format", zap.String("path", c.Path()))
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token format",
				"message": "Token must be in 'Bearer <token>' format",
			})
		}

		// Validate token
		user, err := am.authService.ValidateToken(token)
		if err != nil {
			am.logger.Error("Failed to validate token", err, zap.String("path", c.Path()))
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"message": "Token is invalid or expired",
			})
		}

		// Add user to context
		ctx := context.WithValue(c.Context(), "user", user)
		c.Locals("user", user)
		c.SetUserContext(ctx)

		am.logger.Debug("User authenticated", zap.Int("user_id", user.ID), zap.String("path", c.Path()))
		return c.Next()
	}
}

// OptionalAuth validates JWT tokens if present but doesn't require them
func (am *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			return c.Next()
		}

		// Extract token from "Bearer <token>" format
		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			// Invalid format, continue without authentication
			return c.Next()
		}

		// Validate token
		user, err := am.authService.ValidateToken(token)
		if err != nil {
			// Invalid token, continue without authentication
			return c.Next()
		}

		// Add user to context
		ctx := context.WithValue(c.Context(), "user", user)
		c.Locals("user", user)
		c.SetUserContext(ctx)

		am.logger.Debug("User authenticated (optional)", zap.Int("user_id", user.ID), zap.String("path", c.Path()))
		return c.Next()
	}
}

// shouldSkipAuth checks if the path should skip authentication
func (am *AuthMiddleware) shouldSkipAuth(path string) bool {
	skipPaths := []string{
		"/health",
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/validate",
		"/docs",
		"/swagger",
	}

	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}

	return false
}

// GetUserFromContext extracts user from Fiber context
func GetUserFromContext(c *fiber.Ctx) (*auth.User, bool) {
	user, ok := c.Locals("user").(*auth.User)
	return user, ok
}

// GetUserFromContext extracts user from standard context
func GetUserFromContextStd(ctx context.Context) (*auth.User, bool) {
	user, ok := ctx.Value("user").(*auth.User)
	return user, ok
}
