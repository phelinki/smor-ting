package test

import (
	"context"
	"testing"
	"time"

	"github.com/smorting/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBackgroundSyncQueue tests the background sync queue functionality
func TestBackgroundSyncQueue(t *testing.T) {
	_, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "backgroundsync@example.com",
		Password:  "hashed_password",
		FirstName: "Background",
		LastName:  "Sync",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID

	t.Run("Sync Queue Item Management", func(t *testing.T) {
		// Create a sync queue item
		item := &models.SyncQueueItem{
			UserID:   userID,
			Type:     models.SyncTypeUpload,
			Status:   models.SyncQueuePending,
			Priority: 5,
			Data: map[string]interface{}{
				"type": "booking_update",
				"payload": map[string]interface{}{
					"booking_id": "test_booking_123",
					"status":     "completed",
				},
			},
			MaxRetries:  3,
			NextRetryAt: time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Add item to queue
		err := repo.CreateSyncQueueItem(context.Background(), item)
		require.NoError(t, err)
		assert.NotEmpty(t, item.ID)

		// Get item from queue
		retrievedItem, err := repo.GetSyncQueueItem(context.Background(), item.ID)
		require.NoError(t, err)
		assert.Equal(t, item.UserID, retrievedItem.UserID)
		assert.Equal(t, item.Type, retrievedItem.Type)
		assert.Equal(t, item.Status, retrievedItem.Status)
		assert.Equal(t, item.Priority, retrievedItem.Priority)

		// Update item status
		retrievedItem.Status = models.SyncQueueProcessing
		err = repo.UpdateSyncQueueItem(context.Background(), retrievedItem)
		require.NoError(t, err)

		// Verify update
		updatedItem, err := repo.GetSyncQueueItem(context.Background(), item.ID)
		require.NoError(t, err)
		assert.Equal(t, models.SyncQueueProcessing, updatedItem.Status)
	})

	t.Run("Queue Processing and Priorities", func(t *testing.T) {
		// Create multiple queue items with different priorities
		items := []*models.SyncQueueItem{
			{
				UserID:      userID,
				Type:        models.SyncTypeUpload,
				Status:      models.SyncQueuePending,
				Priority:    1, // Low priority
				Data:        map[string]interface{}{"test": "data1"},
				MaxRetries:  3,
				NextRetryAt: time.Now(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				UserID:      userID,
				Type:        models.SyncTypeDownload,
				Status:      models.SyncQueuePending,
				Priority:    10, // High priority
				Data:        map[string]interface{}{"test": "data2"},
				MaxRetries:  3,
				NextRetryAt: time.Now(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				UserID:      userID,
				Type:        models.SyncTypeConflict,
				Status:      models.SyncQueuePending,
				Priority:    5, // Medium priority
				Data:        map[string]interface{}{"test": "data3"},
				MaxRetries:  3,
				NextRetryAt: time.Now(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		// Add all items to queue
		for _, item := range items {
			err := repo.CreateSyncQueueItem(context.Background(), item)
			require.NoError(t, err)
		}

		// Get pending items ordered by priority
		pendingItems, err := repo.GetPendingSyncQueueItems(context.Background(), userID, 10)
		require.NoError(t, err)
		assert.Len(t, pendingItems, 3)

		// Verify priority ordering (highest priority first)
		assert.Equal(t, 10, pendingItems[0].Priority)
		assert.Equal(t, 5, pendingItems[1].Priority)
		assert.Equal(t, 1, pendingItems[2].Priority)
	})

	t.Run("Retry Logic and Exponential Backoff", func(t *testing.T) {
		// Create a failed item that needs retry
		failedItem := &models.SyncQueueItem{
			UserID:      userID,
			Type:        models.SyncTypeUpload,
			Status:      models.SyncQueueFailed,
			Priority:    5,
			Data:        map[string]interface{}{"test": "retry_data"},
			RetryCount:  1,
			MaxRetries:  3,
			NextRetryAt: time.Now().Add(-1 * time.Minute), // Past retry time
			LastError:   "Network timeout",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.CreateSyncQueueItem(context.Background(), failedItem)
		require.NoError(t, err)

		// Test retry logic
		assert.True(t, failedItem.ShouldRetry())

		// Mark for retry with exponential backoff
		retryPolicy := models.GetDefaultRetryPolicy()
		failedItem.MarkForRetry("Another network error", retryPolicy)

		assert.Equal(t, models.SyncQueueRetrying, failedItem.Status)
		assert.Equal(t, 2, failedItem.RetryCount)
		assert.True(t, failedItem.NextRetryAt.After(time.Now()))

		// Update in database
		err = repo.UpdateSyncQueueItem(context.Background(), failedItem)
		require.NoError(t, err)
	})

	t.Run("Background Sync Status Management", func(t *testing.T) {
		// Get initial background sync status
		status, err := repo.GetBackgroundSyncStatus(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, userID, status.UserID)
		assert.True(t, status.IsEnabled) // Should be enabled by default

		// Update background sync status
		status.AutoRetryEnabled = false
		status.NextScheduledRun = time.Now().Add(1 * time.Hour)
		status.UpdatedAt = time.Now()

		err = repo.UpdateBackgroundSyncStatus(context.Background(), status)
		require.NoError(t, err)

		// Verify update
		updatedStatus, err := repo.GetBackgroundSyncStatus(context.Background(), userID)
		require.NoError(t, err)
		assert.False(t, updatedStatus.AutoRetryEnabled)
		assert.True(t, updatedStatus.NextScheduledRun.After(time.Now()))
	})

	t.Run("Conflict Resolution Queue", func(t *testing.T) {
		// Create a conflict resolution item
		conflictData := &models.ConflictResolution{
			ConflictType:       "data_mismatch",
			ClientVersion:      2,
			ServerVersion:      3,
			ClientData:         map[string]interface{}{"name": "Client Name"},
			ServerData:         map[string]interface{}{"name": "Server Name"},
			ResolutionStrategy: "manual",
			RequiresUserInput:  true,
		}

		conflictItem := &models.SyncQueueItem{
			UserID:       userID,
			Type:         models.SyncTypeConflict,
			Status:       models.SyncQueuePending,
			Priority:     15, // High priority for conflicts
			Data:         map[string]interface{}{"resource_id": "user_profile"},
			ConflictData: conflictData,
			MaxRetries:   1, // Conflicts typically don't retry automatically
			NextRetryAt:  time.Now(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := repo.CreateSyncQueueItem(context.Background(), conflictItem)
		require.NoError(t, err)

		// Get conflict items specifically
		conflictItems, err := repo.GetConflictQueueItems(context.Background(), userID, 10)
		require.NoError(t, err)
		assert.Greater(t, len(conflictItems), 0) // Should have at least one conflict item

		// Find our specific conflict item
		var foundConflictItem *models.SyncQueueItem
		for i := range conflictItems {
			if conflictItems[i].Type == models.SyncTypeConflict &&
				conflictItems[i].ConflictData != nil &&
				conflictItems[i].ConflictData.ConflictType == "data_mismatch" {
				foundConflictItem = &conflictItems[i]
				break
			}
		}

		require.NotNil(t, foundConflictItem, "Should find the conflict item we created")
		assert.Equal(t, models.SyncTypeConflict, foundConflictItem.Type)
		assert.NotNil(t, foundConflictItem.ConflictData)
		assert.True(t, foundConflictItem.ConflictData.RequiresUserInput)
		assert.Equal(t, "data_mismatch", foundConflictItem.ConflictData.ConflictType)
	})

	t.Run("Queue Cleanup and Maintenance", func(t *testing.T) {
		// Create old completed items
		oldCompletedItem := &models.SyncQueueItem{
			UserID:      userID,
			Type:        models.SyncTypeUpload,
			Status:      models.SyncQueueCompleted,
			Priority:    5,
			Data:        map[string]interface{}{"test": "old_data"},
			MaxRetries:  3,
			NextRetryAt: time.Now(),
			CreatedAt:   time.Now().Add(-48 * time.Hour), // 2 days old
			UpdatedAt:   time.Now().Add(-48 * time.Hour),
		}
		now := time.Now().Add(-48 * time.Hour)
		oldCompletedItem.CompletedAt = &now

		err := repo.CreateSyncQueueItem(context.Background(), oldCompletedItem)
		require.NoError(t, err)

		// Cleanup old completed items (older than 24 hours)
		deletedCount, err := repo.CleanupCompletedQueueItems(context.Background(), 24*time.Hour)
		require.NoError(t, err)
		assert.Greater(t, deletedCount, int64(0))

		// Verify item was deleted
		_, err = repo.GetSyncQueueItem(context.Background(), oldCompletedItem.ID)
		assert.Error(t, err) // Should not be found
	})
}

// TestRetryPolicyLogic tests the retry policy calculations
func TestRetryPolicyLogic(t *testing.T) {
	policy := models.GetDefaultRetryPolicy()

	t.Run("Default Retry Policy", func(t *testing.T) {
		assert.Equal(t, 3, policy.MaxRetries)
		assert.Equal(t, 1*time.Second, policy.BaseDelay)
		assert.Equal(t, 30*time.Second, policy.MaxDelay)
		assert.Equal(t, 2.0, policy.Multiplier)
		assert.True(t, policy.RandomJitter)
	})

	t.Run("Retry Time Calculation", func(t *testing.T) {
		// Test retry time calculation
		retry1 := policy.CalculateNextRetry(0)
		retry2 := policy.CalculateNextRetry(1)
		retry3 := policy.CalculateNextRetry(2)

		// All should be in the future
		assert.True(t, retry1.After(time.Now()))
		assert.True(t, retry2.After(time.Now()))
		assert.True(t, retry3.After(time.Now()))

		// Each retry should be later than the previous (with some tolerance for jitter)
		assert.True(t, retry2.After(retry1.Add(-2*time.Second))) // Allow 2s tolerance
		assert.True(t, retry3.After(retry2.Add(-4*time.Second))) // Allow 4s tolerance

		// Exceeding max retries should return zero time
		noMoreRetries := policy.CalculateNextRetry(3)
		assert.True(t, noMoreRetries.IsZero())
	})

	t.Run("Item Retry Logic", func(t *testing.T) {
		item := &models.SyncQueueItem{
			Status:      models.SyncQueueFailed,
			RetryCount:  1,
			MaxRetries:  3,
			NextRetryAt: time.Now().Add(-1 * time.Minute), // Past time
		}

		// Should retry
		assert.True(t, item.ShouldRetry())

		// Mark for retry
		item.MarkForRetry("Test error", policy)
		assert.Equal(t, models.SyncQueueRetrying, item.Status)
		assert.Equal(t, 2, item.RetryCount)
		assert.Equal(t, "Test error", item.LastError)

		// Exhaust retries
		item.RetryCount = 3
		item.MarkForRetry("Final error", policy)
		assert.Equal(t, models.SyncQueueFailed, item.Status)

		// Should not retry when exhausted
		assert.False(t, item.ShouldRetry())
	})

	t.Run("Item State Transitions", func(t *testing.T) {
		item := &models.SyncQueueItem{
			Status: models.SyncQueuePending,
		}

		// Mark completed
		item.MarkCompleted()
		assert.Equal(t, models.SyncQueueCompleted, item.Status)
		assert.NotNil(t, item.CompletedAt)

		// Reset and mark failed
		item.Status = models.SyncQueuePending
		item.MarkFailed("Permanent failure")
		assert.Equal(t, models.SyncQueueFailed, item.Status)
		assert.Equal(t, "Permanent failure", item.LastError)
	})
}
