package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// SyncHandler handles offline sync API endpoints
type SyncHandler struct {
	syncService  *services.SyncService
	auditService *services.AuditService
	logger       *zap.Logger
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler(
	syncService *services.SyncService,
	auditService *services.AuditService,
	logger *zap.Logger,
) *SyncHandler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SyncHandler{
		syncService:  syncService,
		auditService: auditService,
		logger:       logger,
	}
}

// GetSyncStatus returns the current sync status for the authenticated user
func (h *SyncHandler) GetSyncStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	status, err := h.syncService.GetSyncStatus(c.Context(), userObjectID)
	if err != nil {
		h.logger.Error("Failed to get sync status", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to retrieve sync status",
		})
	}

	return c.JSON(status)
}

// UpdateSyncStatus updates the sync status for the authenticated user
func (h *SyncHandler) UpdateSyncStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	var statusUpdate models.SyncStatus
	if err := c.BodyParser(&statusUpdate); err != nil {
		h.logger.Error("Failed to parse sync status update", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse sync status update",
		})
	}

	err = h.syncService.UpdateSyncStatus(c.Context(), userObjectID, &statusUpdate)
	if err != nil {
		h.logger.Error("Failed to update sync status", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to update sync status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Sync status updated successfully",
	})
}

// SyncUp handles uploading offline changes to the server
func (h *SyncHandler) SyncUp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	var changes map[string]interface{}
	if err := c.BodyParser(&changes); err != nil {
		h.logger.Error("Failed to parse sync up data", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse sync data",
		})
	}

	err = h.syncService.SyncUp(c.Context(), userObjectID, changes)
	if err != nil {
		h.logger.Error("Failed to sync up data", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to sync data to server",
		})
	}

	h.logger.Info("Sync up completed", zap.String("userID", userID))

	return c.JSON(fiber.Map{
		"message": "Data synced successfully",
		"status":  "success",
	})
}

// SyncDown handles downloading server changes to the client
func (h *SyncHandler) SyncDown(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	var syncReq models.SyncRequest
	if err := c.BodyParser(&syncReq); err != nil {
		h.logger.Error("Failed to parse sync down request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse sync request",
		})
	}

	// Ensure user ID matches authenticated user
	syncReq.UserID = userObjectID

	// Set default values if not provided
	if syncReq.Limit <= 0 {
		syncReq.Limit = 100
	}

	response, err := h.syncService.SyncDown(c.Context(), &syncReq)
	if err != nil {
		h.logger.Error("Failed to sync down data", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to retrieve sync data",
		})
	}

	h.logger.Info("Sync down completed",
		zap.String("userID", userID),
		zap.Int("recordCount", response.RecordsCount))

	return c.JSON(response)
}

// SyncDownChunked handles downloading server changes in chunks
func (h *SyncHandler) SyncDownChunked(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	var chunkReq models.ChunkedSyncRequest
	if err := c.BodyParser(&chunkReq); err != nil {
		h.logger.Error("Failed to parse chunked sync request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse chunked sync request",
		})
	}

	// Ensure user ID matches authenticated user
	chunkReq.UserID = userObjectID

	// Set default values if not provided
	if chunkReq.ChunkSize <= 0 {
		chunkReq.ChunkSize = 50
	}

	response, err := h.syncService.SyncDownChunked(c.Context(), &chunkReq)
	if err != nil {
		h.logger.Error("Failed to sync chunked data down", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to retrieve chunked sync data",
		})
	}

	h.logger.Info("Chunked sync down completed",
		zap.String("userID", userID),
		zap.Int("chunkIndex", chunkReq.ChunkIndex),
		zap.Int("recordCount", response.RecordsCount))

	return c.JSON(response)
}

// GetSyncMetrics returns sync performance metrics for the authenticated user
func (h *SyncHandler) GetSyncMetrics(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	// Parse limit from query parameter
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Cap at 100 for performance
	}

	metrics, err := h.syncService.GetSyncMetrics(c.Context(), userObjectID, limit)
	if err != nil {
		h.logger.Error("Failed to get sync metrics", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to retrieve sync metrics",
		})
	}

	return c.JSON(fiber.Map{
		"metrics": metrics,
		"count":   len(metrics),
	})
}

// CreateSyncCheckpoint creates a sync checkpoint for resumable sync
func (h *SyncHandler) CreateSyncCheckpoint(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	var req struct {
		Checkpoint string `json:"checkpoint" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse checkpoint request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse checkpoint request",
		})
	}

	if req.Checkpoint == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": "Checkpoint is required",
		})
	}

	err = h.syncService.CreateSyncCheckpoint(c.Context(), userObjectID, req.Checkpoint)
	if err != nil {
		h.logger.Error("Failed to create sync checkpoint", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to create sync checkpoint",
		})
	}

	h.logger.Info("Sync checkpoint created", zap.String("userID", userID))

	return c.JSON(fiber.Map{
		"message": "Sync checkpoint created successfully",
	})
}

// GetOfflineCapabilities returns the offline capabilities supported by the server
func (h *SyncHandler) GetOfflineCapabilities(c *fiber.Ctx) error {
	capabilities := map[string]interface{}{
		"sync": map[string]interface{}{
			"chunked_sync":        true,
			"compression":         true,
			"checkpoints":         true,
			"conflict_resolution": "last_write_wins",
			"max_chunk_size":      1000,
			"supported_data_types": []string{
				"bookings",
				"services",
				"profile",
				"reviews",
				"payments",
			},
		},
		"offline_first": map[string]interface{}{
			"enabled":          true,
			"max_offline_days": 30,
			"auto_sync":        true,
			"background_sync":  true,
			"retry_policy": map[string]interface{}{
				"max_retries":    3,
				"base_delay_ms":  1000,
				"backoff_factor": 2.0,
			},
		},
		"performance": map[string]interface{}{
			"delta_sync":       true,
			"binary_diff":      false,
			"gzip_compression": true,
			"batch_operations": true,
		},
	}

	return c.JSON(fiber.Map{
		"capabilities": capabilities,
		"version":      "1.0",
		"timestamp": fiber.Map{
			"server_time": fiber.Map{
				"utc": fiber.Map{},
			},
		},
	})
}
