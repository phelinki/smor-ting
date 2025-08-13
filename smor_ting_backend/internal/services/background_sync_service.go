package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// BackgroundSyncService handles background sync operations and queue processing
type BackgroundSyncService struct {
	repo         database.Repository
	syncService  *SyncService
	auditService *AuditService
	logger       *zap.Logger
	retryPolicy  models.RetryPolicy
	isRunning    bool
	stopChan     chan struct{}
	wg           sync.WaitGroup
	mu           sync.RWMutex
}

// NewBackgroundSyncService creates a new background sync service
func NewBackgroundSyncService(
	repo database.Repository,
	syncService *SyncService,
	auditService *AuditService,
	logger *zap.Logger,
) *BackgroundSyncService {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &BackgroundSyncService{
		repo:         repo,
		syncService:  syncService,
		auditService: auditService,
		logger:       logger,
		retryPolicy:  models.GetDefaultRetryPolicy(),
		stopChan:     make(chan struct{}),
	}
}

// Start starts the background sync processing
func (bs *BackgroundSyncService) Start(ctx context.Context) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.isRunning {
		return fmt.Errorf("background sync service is already running")
	}

	bs.isRunning = true
	bs.wg.Add(1)

	go bs.processQueue(ctx)

	bs.logger.Info("Background sync service started")
	return nil
}

// Stop stops the background sync processing
func (bs *BackgroundSyncService) Stop() error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if !bs.isRunning {
		return fmt.Errorf("background sync service is not running")
	}

	close(bs.stopChan)
	bs.wg.Wait()
	bs.isRunning = false

	bs.logger.Info("Background sync service stopped")
	return nil
}

// AddToQueue adds an item to the background sync queue
func (bs *BackgroundSyncService) AddToQueue(ctx context.Context, item *models.SyncQueueItem) error {
	// Set default values
	if item.MaxRetries == 0 {
		item.MaxRetries = bs.retryPolicy.MaxRetries
	}
	if item.NextRetryAt.IsZero() {
		item.NextRetryAt = time.Now()
	}

	err := bs.repo.CreateSyncQueueItem(ctx, item)
	if err != nil {
		bs.logger.Error("Failed to add item to sync queue", zap.Error(err))
		return err
	}

	bs.logger.Info("Item added to sync queue",
		zap.String("userID", item.UserID.Hex()),
		zap.String("type", string(item.Type)),
		zap.Int("priority", item.Priority))

	return nil
}

// ProcessUserQueue processes all pending items for a specific user
func (bs *BackgroundSyncService) ProcessUserQueue(ctx context.Context, userID primitive.ObjectID) error {
	pendingItems, err := bs.repo.GetPendingSyncQueueItems(ctx, userID, 50)
	if err != nil {
		return fmt.Errorf("failed to get pending items: %w", err)
	}

	for _, item := range pendingItems {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := bs.processItem(ctx, &item)
			if err != nil {
				bs.logger.Error("Failed to process queue item",
					zap.Error(err),
					zap.String("itemID", item.ID.Hex()),
					zap.String("userID", item.UserID.Hex()))
			}
		}
	}

	return nil
}

// GetQueueStatus returns the current queue status for a user
func (bs *BackgroundSyncService) GetQueueStatus(ctx context.Context, userID primitive.ObjectID) (*models.BackgroundSyncStatus, error) {
	status, err := bs.repo.GetBackgroundSyncStatus(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update counters
	pendingItems, _ := bs.repo.GetPendingSyncQueueItems(ctx, userID, 1000)
	status.PendingItems = len(pendingItems)

	// Count failed items (simplified)
	status.FailedItems = 0
	status.ConflictItems = 0
	for _, item := range pendingItems {
		if item.Status == models.SyncQueueFailed {
			status.FailedItems++
		}
		if item.Type == models.SyncTypeConflict {
			status.ConflictItems++
		}
	}

	return status, nil
}

// ResolveConflict resolves a conflict queue item with user input
func (bs *BackgroundSyncService) ResolveConflict(ctx context.Context, itemID primitive.ObjectID, resolution string, resolvedData map[string]interface{}) error {
	item, err := bs.repo.GetSyncQueueItem(ctx, itemID)
	if err != nil {
		return fmt.Errorf("conflict item not found: %w", err)
	}

	if item.Type != models.SyncTypeConflict {
		return fmt.Errorf("item is not a conflict type")
	}

	if item.ConflictData == nil {
		return fmt.Errorf("conflict data is missing")
	}

	// Update conflict resolution
	item.ConflictData.UserDecision = resolution
	item.ConflictData.ResolvedData = resolvedData
	item.ConflictData.RequiresUserInput = false

	// Mark item for processing
	item.Status = models.SyncQueuePending
	item.Priority = 20 // High priority for resolved conflicts

	err = bs.repo.UpdateSyncQueueItem(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to update conflict item: %w", err)
	}

	bs.logger.Info("Conflict resolved",
		zap.String("itemID", itemID.Hex()),
		zap.String("resolution", resolution))

	return nil
}

// CleanupOldItems removes old completed items from the queue
func (bs *BackgroundSyncService) CleanupOldItems(ctx context.Context, olderThan time.Duration) (int64, error) {
	deletedCount, err := bs.repo.CleanupCompletedQueueItems(ctx, olderThan)
	if err != nil {
		bs.logger.Error("Failed to cleanup old queue items", zap.Error(err))
		return 0, err
	}

	if deletedCount > 0 {
		bs.logger.Info("Cleaned up old queue items", zap.Int64("count", deletedCount))
	}

	return deletedCount, nil
}

// Private methods

func (bs *BackgroundSyncService) processQueue(ctx context.Context) {
	defer bs.wg.Done()

	ticker := time.NewTicker(30 * time.Second) // Process queue every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-bs.stopChan:
			return
		case <-ticker.C:
			bs.processAllPendingItems(ctx)
		}
	}
}

func (bs *BackgroundSyncService) processAllPendingItems(ctx context.Context) {
	// Get a sample of users with pending items (simplified approach)
	// In a real implementation, you'd maintain a list of users with pending items

	// For now, we'll process items as they come in individual requests
	bs.logger.Debug("Background sync processor tick")
}

func (bs *BackgroundSyncService) processItem(ctx context.Context, item *models.SyncQueueItem) error {
	// Mark item as processing
	item.Status = models.SyncQueueProcessing
	now := time.Now()
	item.LastAttemptAt = &now

	err := bs.repo.UpdateSyncQueueItem(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to mark item as processing: %w", err)
	}

	// Process based on type
	var processErr error
	switch item.Type {
	case models.SyncTypeUpload:
		processErr = bs.processSyncUp(ctx, item)
	case models.SyncTypeDownload:
		processErr = bs.processSyncDown(ctx, item)
	case models.SyncTypeConflict:
		processErr = bs.processConflict(ctx, item)
	default:
		processErr = fmt.Errorf("unknown sync type: %s", item.Type)
	}

	// Update item status based on result
	if processErr != nil {
		bs.handleProcessingError(ctx, item, processErr)
	} else {
		item.MarkCompleted()
		bs.repo.UpdateSyncQueueItem(ctx, item)

		bs.logger.Info("Queue item processed successfully",
			zap.String("itemID", item.ID.Hex()),
			zap.String("type", string(item.Type)))
	}

	return processErr
}

func (bs *BackgroundSyncService) processSyncUp(ctx context.Context, item *models.SyncQueueItem) error {
	if bs.syncService == nil {
		return fmt.Errorf("sync service not available")
	}

	// Extract data for sync up
	data, ok := item.Data["payload"].(map[string]interface{})
	if !ok {
		data = item.Data
	}

	return bs.syncService.SyncUp(ctx, item.UserID, data)
}

func (bs *BackgroundSyncService) processSyncDown(ctx context.Context, item *models.SyncQueueItem) error {
	if bs.syncService == nil {
		return fmt.Errorf("sync service not available")
	}

	// Create sync request
	syncReq := &models.SyncRequest{
		UserID:      item.UserID,
		LastSyncAt:  time.Now().Add(-24 * time.Hour), // Default to 24 hours
		Limit:       100,
		Compression: false,
	}

	// Override with item-specific settings if available
	if lastSync, ok := item.Data["last_sync_at"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, lastSync); err == nil {
			syncReq.LastSyncAt = parsed
		}
	}
	if limit, ok := item.Data["limit"].(float64); ok {
		syncReq.Limit = int(limit)
	}

	_, err := bs.syncService.SyncDown(ctx, syncReq)
	return err
}

func (bs *BackgroundSyncService) processConflict(ctx context.Context, item *models.SyncQueueItem) error {
	if item.ConflictData == nil {
		return fmt.Errorf("conflict data is missing")
	}

	if item.ConflictData.RequiresUserInput {
		// Cannot process without user input
		return fmt.Errorf("conflict requires user input")
	}

	// Apply resolved data
	if item.ConflictData.ResolvedData != nil {
		return bs.syncService.SyncUp(ctx, item.UserID, item.ConflictData.ResolvedData)
	}

	// Apply resolution strategy
	switch item.ConflictData.ResolutionStrategy {
	case "client_wins":
		return bs.syncService.SyncUp(ctx, item.UserID, item.ConflictData.ClientData)
	case "server_wins":
		// Server data is already current, just mark as resolved
		return nil
	case "merge":
		// Simple merge strategy (in practice, this would be more sophisticated)
		mergedData := make(map[string]interface{})
		for k, v := range item.ConflictData.ServerData {
			mergedData[k] = v
		}
		for k, v := range item.ConflictData.ClientData {
			mergedData[k] = v // Client data overwrites in simple merge
		}
		return bs.syncService.SyncUp(ctx, item.UserID, mergedData)
	default:
		return fmt.Errorf("unknown resolution strategy: %s", item.ConflictData.ResolutionStrategy)
	}
}

func (bs *BackgroundSyncService) handleProcessingError(ctx context.Context, item *models.SyncQueueItem, err error) {
	if item.ShouldRetry() {
		item.MarkForRetry(err.Error(), bs.retryPolicy)
		bs.logger.Warn("Queue item failed, scheduling retry",
			zap.String("itemID", item.ID.Hex()),
			zap.Int("retryCount", item.RetryCount),
			zap.Time("nextRetry", item.NextRetryAt),
			zap.Error(err))
	} else {
		item.MarkFailed(err.Error())
		bs.logger.Error("Queue item failed permanently",
			zap.String("itemID", item.ID.Hex()),
			zap.Int("retryCount", item.RetryCount),
			zap.Error(err))
	}

	bs.repo.UpdateSyncQueueItem(ctx, item)
}
