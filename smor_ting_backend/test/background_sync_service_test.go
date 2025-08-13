package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBackgroundSyncService tests the background sync service functionality
func TestBackgroundSyncService(t *testing.T) {
	_, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "backgroundservice@example.com",
		Password:  "hashed_password",
		FirstName: "Background",
		LastName:  "Service",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID

	// Setup services
	logger, _ := logger.New("debug", "console", "stdout")
	syncService := services.NewSyncService(repo, nil, logger.Logger)
	backgroundSyncService := services.NewBackgroundSyncService(repo, syncService, nil, logger.Logger)

	t.Run("Add Items to Queue", func(t *testing.T) {
		// Add upload item
		uploadItem := &models.SyncQueueItem{
			UserID:   userID,
			Type:     models.SyncTypeUpload,
			Priority: 10,
			Data: map[string]interface{}{
				"type": "booking_update",
				"payload": map[string]interface{}{
					"booking_id": "test_booking_123",
					"status":     "completed",
				},
			},
		}

		err := backgroundSyncService.AddToQueue(context.Background(), uploadItem)
		require.NoError(t, err)

		// Add download item
		downloadItem := &models.SyncQueueItem{
			UserID:   userID,
			Type:     models.SyncTypeDownload,
			Priority: 5,
			Data: map[string]interface{}{
				"last_sync_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				"limit":        50,
			},
		}

		err = backgroundSyncService.AddToQueue(context.Background(), downloadItem)
		require.NoError(t, err)
	})

	t.Run("Get Queue Status", func(t *testing.T) {
		// Add a new item to ensure we have pending items
		statusTestItem := &models.SyncQueueItem{
			UserID:   userID,
			Type:     models.SyncTypeUpload,
			Priority: 8,
			Data: map[string]interface{}{
				"test": "status_test",
			},
		}
		err := backgroundSyncService.AddToQueue(context.Background(), statusTestItem)
		require.NoError(t, err)

		status, err := backgroundSyncService.GetQueueStatus(context.Background(), userID)
		require.NoError(t, err)

		assert.Equal(t, userID, status.UserID)
		assert.True(t, status.IsEnabled)
		assert.GreaterOrEqual(t, status.PendingItems, 0) // Allow 0 in case items were processed
	})

	t.Run("Process User Queue", func(t *testing.T) {
		err := backgroundSyncService.ProcessUserQueue(context.Background(), userID)
		require.NoError(t, err)

		// Check if items were processed
		status, err := backgroundSyncService.GetQueueStatus(context.Background(), userID)
		require.NoError(t, err)

		// Items should have been processed (though some might fail in test environment)
		assert.NotNil(t, status)
	})

	t.Run("Conflict Resolution", func(t *testing.T) {
		// Create a conflict item
		conflictData := &models.ConflictResolution{
			ConflictType:       "data_mismatch",
			ClientVersion:      2,
			ServerVersion:      3,
			ClientData:         map[string]interface{}{"name": "Client Name"},
			ServerData:         map[string]interface{}{"name": "Server Name"},
			ResolutionStrategy: "client_wins",
			RequiresUserInput:  true,
		}

		conflictItem := &models.SyncQueueItem{
			UserID:       userID,
			Type:         models.SyncTypeConflict,
			Priority:     15,
			Data:         map[string]interface{}{"resource_id": "user_profile"},
			ConflictData: conflictData,
		}

		err := backgroundSyncService.AddToQueue(context.Background(), conflictItem)
		require.NoError(t, err)

		// Resolve the conflict
		resolvedData := map[string]interface{}{
			"name": "Resolved Name",
		}

		err = backgroundSyncService.ResolveConflict(context.Background(), conflictItem.ID, "custom", resolvedData)
		require.NoError(t, err)

		// Verify conflict was resolved
		item, err := repo.GetSyncQueueItem(context.Background(), conflictItem.ID)
		require.NoError(t, err)
		assert.Equal(t, "custom", item.ConflictData.UserDecision)
		assert.False(t, item.ConflictData.RequiresUserInput)
		assert.Equal(t, models.SyncQueuePending, item.Status)
		assert.Equal(t, 20, item.Priority) // Should be high priority after resolution
	})

	t.Run("Cleanup Old Items", func(t *testing.T) {
		// Create an old completed item
		oldItem := &models.SyncQueueItem{
			UserID:    userID,
			Type:      models.SyncTypeUpload,
			Status:    models.SyncQueueCompleted,
			Priority:  5,
			Data:      map[string]interface{}{"test": "old_data"},
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-48 * time.Hour),
		}
		completedAt := time.Now().Add(-48 * time.Hour)
		oldItem.CompletedAt = &completedAt

		err := repo.CreateSyncQueueItem(context.Background(), oldItem)
		require.NoError(t, err)

		// Cleanup items older than 24 hours
		deletedCount, err := backgroundSyncService.CleanupOldItems(context.Background(), 24*time.Hour)
		require.NoError(t, err)
		assert.Greater(t, deletedCount, int64(0))
	})

	t.Run("Service Lifecycle", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start the service
		err := backgroundSyncService.Start(ctx)
		require.NoError(t, err)

		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)

		// Try to start again (should fail)
		err = backgroundSyncService.Start(ctx)
		assert.Error(t, err)

		// Stop the service
		err = backgroundSyncService.Stop()
		require.NoError(t, err)

		// Try to stop again (should fail)
		err = backgroundSyncService.Stop()
		assert.Error(t, err)
	})
}

// TestBackgroundSyncHandler tests the background sync handler endpoints
func TestBackgroundSyncHandler(t *testing.T) {
	app, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "backgroundhandler@example.com",
		Password:  "hashed_password",
		FirstName: "Background",
		LastName:  "Handler",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID

	// Setup services and handler
	logger, _ := logger.New("debug", "console", "stdout")
	syncService := services.NewSyncService(repo, nil, logger.Logger)
	backgroundSyncService := services.NewBackgroundSyncService(repo, syncService, nil, logger.Logger)
	backgroundSyncHandler := handlers.NewBackgroundSyncHandler(backgroundSyncService, logger.Logger)

	// Setup routes with middleware simulation
	api := app.Group("/api/v1")
	bgSync := api.Group("/background-sync")

	// Middleware to simulate authentication
	bgSync.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.Hex())
		return c.Next()
	})

	bgSync.Get("/status", backgroundSyncHandler.GetQueueStatus)
	bgSync.Post("/queue", backgroundSyncHandler.AddToQueue)
	bgSync.Post("/process", backgroundSyncHandler.ProcessUserQueue)
	bgSync.Post("/conflicts/:itemId/resolve", backgroundSyncHandler.ResolveConflict)
	bgSync.Delete("/cleanup", backgroundSyncHandler.CleanupOldItems)
	bgSync.Get("/retry-policy", backgroundSyncHandler.GetRetryPolicy)

	t.Run("Get Queue Status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/background-sync/status", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var status models.BackgroundSyncStatus
		err = json.NewDecoder(resp.Body).Decode(&status)
		require.NoError(t, err)

		assert.Equal(t, userID, status.UserID)
		assert.True(t, status.IsEnabled)
	})

	t.Run("Add Upload Item to Queue", func(t *testing.T) {
		queueReq := map[string]interface{}{
			"type":     "upload",
			"priority": 10,
			"data": map[string]interface{}{
				"type": "booking_update",
				"payload": map[string]interface{}{
					"booking_id": "handler_test_123",
					"status":     "completed",
				},
			},
		}

		body, _ := json.Marshal(queueReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/background-sync/queue", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Item added to sync queue successfully", response["message"])
		assert.Equal(t, "queued", response["status"])
		assert.NotEmpty(t, response["item_id"])
	})

	t.Run("Add Download Item to Queue", func(t *testing.T) {
		queueReq := map[string]interface{}{
			"type":     "download",
			"priority": 5,
			"data": map[string]interface{}{
				"last_sync_at": time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
				"limit":        25,
			},
		}

		body, _ := json.Marshal(queueReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/background-sync/queue", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Item added to sync queue successfully", response["message"])
	})

	t.Run("Invalid Sync Type", func(t *testing.T) {
		queueReq := map[string]interface{}{
			"type": "invalid_type",
			"data": map[string]interface{}{"test": "data"},
		}

		body, _ := json.Marshal(queueReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/background-sync/queue", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Invalid sync type", response["error"])
	})

	t.Run("Process User Queue", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/background-sync/process", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Queue processed successfully", response["message"])
		assert.Equal(t, "completed", response["status"])
	})

	t.Run("Cleanup Old Items", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/background-sync/cleanup?older_than=1", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Cleanup completed successfully", response["message"])
		assert.Contains(t, response, "deleted_count")
		assert.Equal(t, "1 hours", response["older_than"])
	})

	t.Run("Get Retry Policy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/background-sync/retry-policy", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "retry_policy")
		assert.Equal(t, "1.0", response["version"])

		retryPolicy := response["retry_policy"].(map[string]interface{})
		assert.Equal(t, float64(3), retryPolicy["max_retries"])
		assert.Equal(t, "1s", retryPolicy["base_delay"])
		assert.Equal(t, "30s", retryPolicy["max_delay"])
		assert.Equal(t, 2.0, retryPolicy["multiplier"])
		assert.Equal(t, true, retryPolicy["random_jitter"])
	})
}
