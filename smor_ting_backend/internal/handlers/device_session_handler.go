package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// DeviceSessionHandler handles device session management endpoints
type DeviceSessionHandler struct {
	repo         database.Repository
	jwtService   *services.JWTRefreshService
	auditService *services.AuditService
	logger       *zap.Logger
}

// NewDeviceSessionHandler creates a new device session handler
func NewDeviceSessionHandler(
	repo database.Repository,
	jwtService *services.JWTRefreshService,
	auditService *services.AuditService,
	logger *zap.Logger,
) *DeviceSessionHandler {
	return &DeviceSessionHandler{
		repo:         repo,
		jwtService:   jwtService,
		auditService: auditService,
		logger:       logger,
	}
}

// GetDeviceSessions returns all device sessions for the authenticated user
func (h *DeviceSessionHandler) GetDeviceSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	sessions, err := h.repo.GetUserDeviceSessions(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user device sessions", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to retrieve device sessions",
		})
	}

	// Find current device session based on the request
	var currentDevice *models.DeviceSession
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		// Extract token and find matching session
		// This is simplified - in production you'd validate the token first
		for i := range sessions {
			if sessions[i].IsActive {
				// You could match by IP, User-Agent, or other factors
				if sessions[i].IPAddress == c.IP() {
					currentDevice = &sessions[i]
					break
				}
			}
		}
	}

	response := models.DeviceSessionListResponse{
		Sessions:      sessions,
		CurrentDevice: currentDevice,
		TotalSessions: len(sessions),
	}

	return c.JSON(response)
}

// RevokeDeviceSession revokes a specific device session
func (h *DeviceSessionHandler) RevokeDeviceSession(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	sessionID := c.Params("sessionID")

	if sessionID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": "Session ID is required",
		})
	}

	// Verify the session belongs to the user
	session, err := h.repo.GetDeviceSession(c.Context(), sessionID)
	if err != nil {
		h.logger.Warn("Device session not found", zap.String("sessionID", sessionID), zap.String("userID", userID))
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   "Session not found",
			"message": "The specified device session does not exist",
		})
	}

	if session.UserID.Hex() != userID {
		h.logger.Warn("Unauthorized session revocation attempt",
			zap.String("sessionID", sessionID),
			zap.String("userID", userID),
			zap.String("sessionUserID", session.UserID.Hex()))
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error":   "Forbidden",
			"message": "You don't have permission to revoke this session",
		})
	}

	// Revoke the session
	err = h.repo.RevokeDeviceSession(c.Context(), sessionID)
	if err != nil {
		h.logger.Error("Failed to revoke device session", zap.Error(err), zap.String("sessionID", sessionID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to revoke device session",
		})
	}

	// Log security event
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	sessionObjectID, _ := primitive.ObjectIDFromHex(sessionID)
	securityEvent := &models.SecurityEvent{
		UserID:    userObjectID,
		EventType: models.DeviceRevoked,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
		DeviceID:  session.DeviceID,
		SessionID: &sessionObjectID,
		Metadata: map[string]interface{}{
			"device_name": session.DeviceName,
			"platform":    session.Platform,
			"revoked_by":  "user",
		},
		Timestamp: time.Now(),
	}

	if err := h.repo.LogSecurityEvent(c.Context(), securityEvent); err != nil {
		h.logger.Warn("Failed to log security event", zap.Error(err))
	}

	h.logger.Info("Device session revoked",
		zap.String("sessionID", sessionID),
		zap.String("userID", userID),
		zap.String("deviceName", session.DeviceName))

	return c.JSON(fiber.Map{
		"message": "Device session revoked successfully",
	})
}

// RevokeAllDeviceSessions revokes all device sessions for the user except current
func (h *DeviceSessionHandler) RevokeAllDeviceSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	currentSessionID := c.Locals("sessionID") // This would be set by auth middleware

	// Get all user sessions
	sessions, err := h.repo.GetUserDeviceSessions(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user device sessions", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to retrieve device sessions",
		})
	}

	revokedCount := 0
	for _, session := range sessions {
		// Skip current session and already revoked sessions
		if session.ID.Hex() == currentSessionID || !session.IsActive {
			continue
		}

		err := h.repo.RevokeDeviceSession(c.Context(), session.ID.Hex())
		if err != nil {
			h.logger.Warn("Failed to revoke device session",
				zap.Error(err),
				zap.String("sessionID", session.ID.Hex()))
			continue
		}
		revokedCount++
	}

	// Log security event
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	securityEvent := &models.SecurityEvent{
		UserID:    userObjectID,
		EventType: models.DeviceRevoked,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
		Metadata: map[string]interface{}{
			"action":        "revoke_all_sessions",
			"revoked_count": revokedCount,
			"reason":        "user_requested",
		},
		Timestamp: time.Now(),
	}

	if err := h.repo.LogSecurityEvent(c.Context(), securityEvent); err != nil {
		h.logger.Warn("Failed to log security event", zap.Error(err))
	}

	h.logger.Info("All device sessions revoked",
		zap.String("userID", userID),
		zap.Int("revokedCount", revokedCount))

	return c.JSON(fiber.Map{
		"message":       "All other device sessions revoked successfully",
		"revoked_count": revokedCount,
	})
}

// RefreshTokenWithRotation performs refresh token rotation for enhanced security
func (h *DeviceSessionHandler) RefreshTokenWithRotation(c *fiber.Ctx) error {
	var req models.TokenRotationRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse token rotation request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse request",
		})
	}

	// Validate required fields
	if req.RefreshToken == "" || req.DeviceID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": "Refresh token and device ID are required",
		})
	}

	// Find device session by refresh token
	session, err := h.repo.GetDeviceSessionByRefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Warn("Invalid refresh token for rotation", zap.String("deviceID", req.DeviceID))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid refresh token",
			"message": "The provided refresh token is invalid or expired",
		})
	}

	// Verify device ID matches
	if session.DeviceID != req.DeviceID {
		h.logger.Warn("Device ID mismatch during token rotation",
			zap.String("providedDeviceID", req.DeviceID),
			zap.String("sessionDeviceID", session.DeviceID))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Device mismatch",
			"message": "Device ID does not match the session",
		})
	}

	// Get user for token generation
	user, err := h.repo.GetUserByID(c.Context(), session.UserID)
	if err != nil {
		h.logger.Error("Failed to get user for token rotation", zap.Error(err), zap.String("userID", session.UserID.Hex()))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to process token rotation",
		})
	}

	// Generate new token pair
	tokenPair, err := h.jwtService.GenerateTokenPair(user)
	if err != nil {
		h.logger.Error("Failed to generate new token pair", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to generate new tokens",
		})
	}

	// Rotate the refresh token in the session
	err = h.repo.RotateRefreshToken(c.Context(), session.ID.Hex(), tokenPair.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to rotate refresh token", zap.Error(err), zap.String("sessionID", session.ID.Hex()))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to rotate refresh token",
		})
	}

	// Update session activity
	err = h.repo.UpdateDeviceSessionActivity(c.Context(), session.ID.Hex())
	if err != nil {
		h.logger.Warn("Failed to update session activity", zap.Error(err))
	}

	// Log security event
	securityEvent := &models.SecurityEvent{
		UserID:    session.UserID,
		EventType: models.TokenRefreshEvent,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
		DeviceID:  session.DeviceID,
		SessionID: &session.ID,
		Metadata: map[string]interface{}{
			"device_name": session.DeviceName,
			"platform":    session.Platform,
		},
		Timestamp: time.Now(),
	}

	if err := h.repo.LogSecurityEvent(c.Context(), securityEvent); err != nil {
		h.logger.Warn("Failed to log token refresh event", zap.Error(err))
	}

	h.logger.Info("Token rotated successfully",
		zap.String("userID", user.ID.Hex()),
		zap.String("deviceID", session.DeviceID))

	return c.JSON(fiber.Map{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    int(30 * time.Minute.Seconds()), // 30 minutes
		"message":       "Token rotated successfully",
	})
}
