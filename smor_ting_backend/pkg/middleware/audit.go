package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"go.uber.org/zap"
)

// AuditMiddleware provides audit logging for sensitive operations
type AuditMiddleware struct {
	auditService *services.AuditService
	logger       *zap.Logger
}

// NewAuditMiddleware creates a new audit middleware
func NewAuditMiddleware(auditService *services.AuditService, logger *zap.Logger) *AuditMiddleware {
	return &AuditMiddleware{
		auditService: auditService,
		logger:       logger,
	}
}

// AuditConfig defines configuration for audit logging
type AuditConfig struct {
	Action        services.AuditAction
	Resource      string
	SensitiveOp   bool                    // Whether this is a sensitive operation requiring audit
	GetResourceID func(*fiber.Ctx) string // Function to extract resource ID from context
}

// Audit creates middleware that logs the specified action
func (m *AuditMiddleware) Audit(config AuditConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user information if available
		user, _ := GetUserFromContextModels(c)

		// Get request details
		ipAddress := c.IP()
		userAgent := c.Get("User-Agent")
		method := c.Method()
		path := c.Path()

		// Get resource ID if function provided
		resourceID := ""
		if config.GetResourceID != nil {
			resourceID = config.GetResourceID(c)
		}

		// Record start time
		startTime := time.Now()

		// Execute the request
		err := c.Next()

		// Record end time and success status
		endTime := time.Now()
		success := err == nil && c.Response().StatusCode() < 400

		// Prepare audit details
		details := map[string]interface{}{
			"method":        method,
			"path":          path,
			"status_code":   c.Response().StatusCode(),
			"duration_ms":   endTime.Sub(startTime).Milliseconds(),
			"request_size":  len(c.Body()),
			"response_size": len(c.Response().Body()),
		}

		if resourceID != "" {
			details["resource_id"] = resourceID
		}

		// Add error details if failed
		if err != nil {
			details["error"] = err.Error()
		}

		// Log the audit entry
		var auditErr error
		if user != nil {
			auditErr = m.auditService.LogUserAction(
				c.Context(),
				user,
				config.Action,
				config.Resource,
				ipAddress,
				userAgent,
				success,
				details,
			)
		} else {
			// For actions without authenticated user context
			auditEntry := &services.AuditEntry{
				UserID:    "anonymous",
				Action:    config.Action,
				Resource:  config.Resource,
				IPAddress: ipAddress,
				UserAgent: userAgent,
				Success:   success,
				Details:   details,
			}
			if resourceID != "" {
				auditEntry.ResourceID = resourceID
			}
			auditErr = m.auditService.LogAction(c.Context(), auditEntry)
		}

		if auditErr != nil {
			m.logger.Error("Failed to create audit log",
				zap.Error(auditErr),
				zap.String("action", string(config.Action)),
				zap.String("resource", config.Resource),
			)
		}

		return err
	}
}

// AuditSensitiveOperation creates middleware for highly sensitive operations
func (m *AuditMiddleware) AuditSensitiveOperation(action services.AuditAction, resource string) fiber.Handler {
	return m.Audit(AuditConfig{
		Action:      action,
		Resource:    resource,
		SensitiveOp: true,
	})
}

// AuditWithResourceID creates middleware that extracts resource ID from URL params
func (m *AuditMiddleware) AuditWithResourceID(action services.AuditAction, resource string, paramName string) fiber.Handler {
	return m.Audit(AuditConfig{
		Action:   action,
		Resource: resource,
		GetResourceID: func(c *fiber.Ctx) string {
			return c.Params(paramName)
		},
	})
}

// AdminActionAudit creates middleware specifically for admin actions
func (m *AuditMiddleware) AdminActionAudit(action services.AuditAction, resource string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First check if user is admin
		user, ok := GetUserFromContextModels(c)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "User not found in context",
			})
		}

		if user.Role != models.AdminRole {
			// Log unauthorized admin access attempt
			m.auditService.LogUserAction(
				c.Context(),
				user,
				action,
				resource,
				c.IP(),
				c.Get("User-Agent"),
				false,
				map[string]interface{}{
					"reason":        "insufficient_privileges",
					"required_role": string(models.AdminRole),
					"user_role":     string(user.Role),
				},
			)

			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden",
				"message": "Administrative privileges required",
			})
		}

		// Apply audit logging
		return m.Audit(AuditConfig{
			Action:      action,
			Resource:    resource,
			SensitiveOp: true,
		})(c)
	}
}

// LogSecurityEvent logs security-related events (brute force, suspicious activity, etc.)
func (m *AuditMiddleware) LogSecurityEvent(action services.AuditAction, email string, ipAddress string, userAgent string, details map[string]interface{}) {
	err := m.auditService.LogSecurityEvent(
		context.Background(),
		action,
		email,
		ipAddress,
		userAgent,
		details,
	)
	if err != nil {
		m.logger.Error("Failed to log security event",
			zap.Error(err),
			zap.String("action", string(action)),
			zap.String("email", email),
		)
	}
}

// LogBruteForceAttempt logs brute force protection events
func (m *AuditMiddleware) LogBruteForceAttempt(email string, ipAddress string, userAgent string, blocked bool, attemptCount int) {
	details := map[string]interface{}{
		"blocked":       blocked,
		"attempt_count": attemptCount,
	}

	m.LogSecurityEvent(services.ActionBruteForceBlock, email, ipAddress, userAgent, details)
}
