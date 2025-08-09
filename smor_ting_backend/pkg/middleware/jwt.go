package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// JWTAuthMiddleware validates access tokens and loads the Mongo-backed user
type JWTAuthMiddleware struct {
	jwtService *services.JWTRefreshService
	repo       database.Repository
	logger     *logger.Logger
}

// NewJWTAuthMiddleware creates a new JWTAuthMiddleware
func NewJWTAuthMiddleware(jwtService *services.JWTRefreshService, repo database.Repository, logger *logger.Logger) (*JWTAuthMiddleware, error) {
	if jwtService == nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "jwt service is required")
	}
	if repo == nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "repository is required")
	}
	if logger == nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "logger is required")
	}
	return &JWTAuthMiddleware{jwtService: jwtService, repo: repo, logger: logger}, nil
}

// Authenticate enforces a valid access token and attaches models.User to context
func (m *JWTAuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing Authorization header", zap.String("path", c.Path()))
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Missing token",
				"message": "Authorization header is required",
			})
		}

		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			m.logger.Warn("Invalid Authorization header format", zap.String("path", c.Path()))
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token format",
				"message": "Token must be in 'Bearer <token>' format",
			})
		}

		// Validate token and load claims
		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			m.logger.Warn("Invalid access token", zap.Error(err))
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"message": "The access token is invalid or expired",
			})
		}

		// Load user from repository
		userObjectID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			m.logger.Warn("Invalid user id in token", zap.String("user_id", claims.UserID))
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"message": "Token contains invalid user identifier",
			})
		}

		user, err := m.repo.GetUserByID(c.Context(), userObjectID)
		if err != nil {
			m.logger.Warn("User not found for token", zap.String("user_id", userObjectID.Hex()))
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"message": "Associated user not found",
			})
		}

		// Attach user to context
		ctx := context.WithValue(c.Context(), "user", user)
		c.Locals("user", user)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// OptionalAuth tries to authenticate, but proceeds without error if missing/invalid
func (m *JWTAuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return c.Next()
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			return c.Next()
		}
		userObjectID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			return c.Next()
		}
		user, err := m.repo.GetUserByID(c.Context(), userObjectID)
		if err != nil {
			return c.Next()
		}

		ctx := context.WithValue(c.Context(), "user", user)
		c.Locals("user", user)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

// RequireRoles enforces that the authenticated user has one of the allowed roles
func (m *JWTAuthMiddleware) RequireRoles(allowedRoles ...models.UserRole) fiber.Handler {
	allowed := make(map[models.UserRole]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok || user == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "User not found in context",
			})
		}

		if _, ok := allowed[user.Role]; !ok {
			// Build a readable list of allowed roles for logs/response
			roleStrings := make([]string, 0, len(allowedRoles))
			for _, r := range allowedRoles {
				roleStrings = append(roleStrings, string(r))
			}
			m.logger.Warn("Forbidden: insufficient role",
				zap.String("user_id", user.ID.Hex()),
				zap.String("user_role", string(user.Role)),
				zap.String("allowed", strings.Join(roleStrings, ",")),
				zap.String("path", c.Path()),
			)
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden",
				"message": "You do not have permission to access this resource",
			})
		}

		return c.Next()
	}
}

// GetUserFromContextModels extracts a models.User from Fiber context
func GetUserFromContextModels(c *fiber.Ctx) (*models.User, bool) {
	user, ok := c.Locals("user").(*models.User)
	return user, ok
}

// GetUserFromContextStdModels extracts models.User from standard context
func GetUserFromContextStdModels(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value("user").(*models.User)
	return user, ok
}
