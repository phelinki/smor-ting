package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// BackgroundSyncHandler handles background sync API endpoints
type BackgroundSyncHandler struct {
	backgroundSyncService *services.BackgroundSyncService
	logger                *zap.Logger
}

// NewBackgroundSyncHandler creates a new background sync handler
func NewBackgroundSyncHandler(
	backgroundSyncService *services.BackgroundSyncService,
	logger *zap.Logger,
) *BackgroundSyncHandler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &BackgroundSyncHandler{
		backgroundSyncService: backgroundSyncService,
		logger:                logger,
	}
}

// GetQueueStatus returns the current queue status for the authenticated user
func (h *BackgroundSyncHandler) GetQueueStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	status, err := h.backgroundSyncService.GetQueueStatus(c.Context(), userObjectID)
	if err != nil {
		h.logger.Error("Failed to get queue status", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to retrieve queue status",
		})
	}

	return c.JSON(status)
}

// AddToQueue adds an item to the background sync queue
func (h *BackgroundSyncHandler) AddToQueue(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	var req struct {
		Type     models.SyncQueueItemType `json:"type" validate:"required"`
		Priority int                      `json:"priority"`
		Data     map[string]interface{}   `json:"data" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse add to queue request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse request body",
		})
	}

	// Validate sync type
	if req.Type != models.SyncTypeUpload && req.Type != models.SyncTypeDownload && req.Type != models.SyncTypeConflict {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid sync type",
			"message": "Sync type must be upload, download, or conflict_resolution",
		})
	}

	// Set default priority if not provided
	if req.Priority == 0 {
		switch req.Type {
		case models.SyncTypeConflict:
			req.Priority = 15 // High priority for conflicts
		case models.SyncTypeUpload:
			req.Priority = 10 // Medium priority for uploads
		case models.SyncTypeDownload:
			req.Priority = 5 // Lower priority for downloads
		}
	}

	item := &models.SyncQueueItem{
		UserID:      userObjectID,
		Type:        req.Type,
		Status:      models.SyncQueuePending,
		Priority:    req.Priority,
		Data:        req.Data,
		MaxRetries:  3,
		NextRetryAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.backgroundSyncService.AddToQueue(c.Context(), item)
	if err != nil {
		h.logger.Error("Failed to add item to queue", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to add item to sync queue",
		})
	}

	h.logger.Info("Item added to sync queue",
		zap.String("userID", userID),
		zap.String("type", string(req.Type)),
		zap.String("itemID", item.ID.Hex()))

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Item added to sync queue successfully",
		"item_id": item.ID.Hex(),
		"status":  "queued",
	})
}

// ProcessUserQueue manually triggers processing of the user's queue
func (h *BackgroundSyncHandler) ProcessUserQueue(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	err = h.backgroundSyncService.ProcessUserQueue(c.Context(), userObjectID)
	if err != nil {
		h.logger.Error("Failed to process user queue", zap.Error(err), zap.String("userID", userID))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to process sync queue",
		})
	}

	h.logger.Info("User queue processed", zap.String("userID", userID))

	return c.JSON(fiber.Map{
		"message": "Queue processed successfully",
		"status":  "completed",
	})
}

// ResolveConflict resolves a conflict with user input
func (h *BackgroundSyncHandler) ResolveConflict(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	itemIDStr := c.Params("itemId")

	itemID, err := primitive.ObjectIDFromHex(itemIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid item ID",
			"message": "Failed to parse item ID",
		})
	}

	var req struct {
		Resolution   string                 `json:"resolution" validate:"required"` // "client_wins", "server_wins", "merge", "custom"
		ResolvedData map[string]interface{} `json:"resolved_data,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse conflict resolution request", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse request body",
		})
	}

	// Validate resolution strategy
	validResolutions := []string{"client_wins", "server_wins", "merge", "custom"}
	isValid := false
	for _, valid := range validResolutions {
		if req.Resolution == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid resolution",
			"message": "Resolution must be one of: client_wins, server_wins, merge, custom",
		})
	}

	// For custom resolution, resolved_data is required
	if req.Resolution == "custom" && req.ResolvedData == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing resolved data",
			"message": "Resolved data is required for custom resolution",
		})
	}

	err = h.backgroundSyncService.ResolveConflict(c.Context(), itemID, req.Resolution, req.ResolvedData)
	if err != nil {
		h.logger.Error("Failed to resolve conflict", zap.Error(err),
			zap.String("userID", userID),
			zap.String("itemID", itemIDStr))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to resolve conflict",
		})
	}

	h.logger.Info("Conflict resolved",
		zap.String("userID", userID),
		zap.String("itemID", itemIDStr),
		zap.String("resolution", req.Resolution))

	return c.JSON(fiber.Map{
		"message":    "Conflict resolved successfully",
		"item_id":    itemIDStr,
		"resolution": req.Resolution,
		"status":     "resolved",
	})
}

// CleanupOldItems removes old completed items from the queue
func (h *BackgroundSyncHandler) CleanupOldItems(c *fiber.Ctx) error {
	// Parse optional olderThan parameter (in hours)
	olderThanStr := c.Query("older_than", "24")
	olderThanHours, err := strconv.Atoi(olderThanStr)
	if err != nil || olderThanHours < 1 {
		olderThanHours = 24
	}

	olderThan := time.Duration(olderThanHours) * time.Hour

	deletedCount, err := h.backgroundSyncService.CleanupOldItems(c.Context(), olderThan)
	if err != nil {
		h.logger.Error("Failed to cleanup old items", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to cleanup old items",
		})
	}

	h.logger.Info("Old items cleaned up", zap.Int64("deletedCount", deletedCount))

	return c.JSON(fiber.Map{
		"message":       "Cleanup completed successfully",
		"deleted_count": deletedCount,
		"older_than":    olderThanStr + " hours",
	})
}

// GetConflictItems returns pending conflict items for the user
func (h *BackgroundSyncHandler) GetConflictItems(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"message": "Failed to parse user ID",
		})
	}

	// Parse limit parameter
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50 // Cap at 50 for performance
	}

	// Get conflict items from repository directly
	// In a real implementation, this would be part of the background sync service
	return c.Status(http.StatusNotImplemented).JSON(fiber.Map{
		"error":   "Not implemented",
		"message": "Conflict items endpoint not yet implemented",
	})
}

// GetRetryPolicy returns the current retry policy configuration
func (h *BackgroundSyncHandler) GetRetryPolicy(c *fiber.Ctx) error {
	policy := models.GetDefaultRetryPolicy()

	return c.JSON(fiber.Map{
		"retry_policy": map[string]interface{}{
			"max_retries":   policy.MaxRetries,
			"base_delay":    policy.BaseDelay.String(),
			"max_delay":     policy.MaxDelay.String(),
			"multiplier":    policy.Multiplier,
			"random_jitter": policy.RandomJitter,
		},
		"version": "1.0",
	})
}
