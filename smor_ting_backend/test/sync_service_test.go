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

// TestSyncService tests the sync service functionality
func TestSyncService(t *testing.T) {
	_, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "syncservice@example.com",
		Password:  "hashed_password",
		FirstName: "Sync",
		LastName:  "Service",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID

	// Create sync service
	logger, _ := logger.New("debug", "console", "stdout")
	syncService := services.NewSyncService(repo, nil, logger.Logger)

	t.Run("Sync Status Operations", func(t *testing.T) {
		// Get initial sync status
		status, err := syncService.GetSyncStatus(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, userID, status.UserID)
		assert.True(t, status.IsOnline)

		// Update sync status
		updateStatus := &models.SyncStatus{
			IsOnline:        false,
			ConnectionType:  "mobile",
			ConnectionSpeed: "slow",
			PendingChanges:  5,
		}

		err = syncService.UpdateSyncStatus(context.Background(), userID, updateStatus)
		require.NoError(t, err)

		// Verify update
		updatedStatus, err := syncService.GetSyncStatus(context.Background(), userID)
		require.NoError(t, err)
		assert.False(t, updatedStatus.IsOnline)
		assert.Equal(t, "mobile", updatedStatus.ConnectionType)
		assert.Equal(t, "slow", updatedStatus.ConnectionSpeed)
	})

	t.Run("Sync Up Operations", func(t *testing.T) {
		// Create offline changes
		changes := map[string]interface{}{
			"bookings": []map[string]interface{}{
				{
					"id":         "offline_booking_1",
					"status":     "completed",
					"updated_at": time.Now(),
				},
			},
			"profile_updates": map[string]interface{}{
				"first_name": "Updated Name",
				"updated_at": time.Now(),
			},
		}

		// Sync up changes
		err := syncService.SyncUp(context.Background(), userID, changes)
		require.NoError(t, err)

		// Verify sync metrics were recorded
		metrics, err := syncService.GetSyncMetrics(context.Background(), userID, 5)
		require.NoError(t, err)
		assert.Len(t, metrics, 1)
		assert.True(t, metrics[0].SyncSuccess)
		assert.Greater(t, metrics[0].RecordsSynced, 0)
	})

	t.Run("Sync Down Operations", func(t *testing.T) {
		// Create sync request
		syncReq := &models.SyncRequest{
			UserID:      userID,
			LastSyncAt:  time.Now().Add(-1 * time.Hour),
			Limit:       100,
			Compression: false,
		}

		// Sync down changes
		response, err := syncService.SyncDown(context.Background(), syncReq)
		require.NoError(t, err)

		assert.NotNil(t, response)
		assert.NotNil(t, response.Data)
		assert.NotEmpty(t, response.Checkpoint)
		assert.Greater(t, response.SyncDuration, time.Duration(0))
	})

	t.Run("Chunked Sync Operations", func(t *testing.T) {
		// Create chunked sync request
		chunkReq := &models.ChunkedSyncRequest{
			UserID:    userID,
			ChunkSize: 10,
		}

		// Sync down chunked data
		response, err := syncService.SyncDownChunked(context.Background(), chunkReq)
		require.NoError(t, err)

		assert.NotNil(t, response)
		assert.NotNil(t, response.Data)
		assert.NotEmpty(t, response.Checkpoint)
	})

	t.Run("Sync Checkpoint Operations", func(t *testing.T) {
		checkpoint := "test_checkpoint_123"

		// Create sync checkpoint
		err := syncService.CreateSyncCheckpoint(context.Background(), userID, checkpoint)
		require.NoError(t, err)

		// Verify checkpoint was created (by checking it doesn't error)
		// In a real implementation, you'd have a GetSyncCheckpoint method
		assert.NoError(t, err)
	})
}

// TestSyncHandler tests the sync handler endpoints
func TestSyncHandler(t *testing.T) {
	app, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "synchandler@example.com",
		Password:  "hashed_password",
		FirstName: "Sync",
		LastName:  "Handler",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID

	// Create sync service and handler
	logger, _ := logger.New("debug", "console", "stdout")
	syncService := services.NewSyncService(repo, nil, logger.Logger)
	syncHandler := handlers.NewSyncHandler(syncService, nil, logger.Logger)

	// Setup routes with middleware simulation
	api := app.Group("/api/v1")
	sync := api.Group("/sync")

	// Middleware to simulate authentication
	sync.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.Hex())
		return c.Next()
	})

	sync.Get("/status", syncHandler.GetSyncStatus)
	sync.Put("/status", syncHandler.UpdateSyncStatus)
	sync.Post("/up", syncHandler.SyncUp)
	sync.Post("/down", syncHandler.SyncDown)
	sync.Post("/down/chunked", syncHandler.SyncDownChunked)
	sync.Get("/metrics", syncHandler.GetSyncMetrics)
	sync.Post("/checkpoint", syncHandler.CreateSyncCheckpoint)
	sync.Get("/capabilities", syncHandler.GetOfflineCapabilities)

	t.Run("Get Sync Status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sync/status", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var status models.SyncStatus
		err = json.NewDecoder(resp.Body).Decode(&status)
		require.NoError(t, err)

		assert.Equal(t, userID, status.UserID)
		assert.True(t, status.IsOnline)
	})

	t.Run("Update Sync Status", func(t *testing.T) {
		statusUpdate := models.SyncStatus{
			IsOnline:        false,
			ConnectionType:  "wifi",
			ConnectionSpeed: "fast",
			PendingChanges:  3,
		}

		body, _ := json.Marshal(statusUpdate)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/sync/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Sync status updated successfully", response["message"])
	})

	t.Run("Sync Up", func(t *testing.T) {
		changes := map[string]interface{}{
			"bookings": []map[string]interface{}{
				{
					"id":      "offline_booking_1",
					"status":  "confirmed",
					"offline": true,
				},
			},
			"profile_updates": map[string]interface{}{
				"phone": "+231777123456",
			},
		}

		body, _ := json.Marshal(changes)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sync/up", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Data synced successfully", response["message"])
		assert.Equal(t, "success", response["status"])
	})

	t.Run("Sync Down", func(t *testing.T) {
		syncReq := models.SyncRequest{
			LastSyncAt:  time.Now().Add(-2 * time.Hour),
			Limit:       50,
			Compression: false,
		}

		body, _ := json.Marshal(syncReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sync/down", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.SyncResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.Data)
		assert.NotEmpty(t, response.Checkpoint)
		assert.Greater(t, response.SyncDuration, time.Duration(0))
	})

	t.Run("Chunked Sync Down", func(t *testing.T) {
		chunkReq := models.ChunkedSyncRequest{
			ChunkIndex: 0,
			ChunkSize:  25,
		}

		body, _ := json.Marshal(chunkReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sync/down/chunked", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.ChunkedSyncResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.Data)
		assert.NotEmpty(t, response.Checkpoint)
	})

	t.Run("Get Sync Metrics", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sync/metrics?limit=5", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "metrics")
		assert.Contains(t, response, "count")
	})

	t.Run("Create Sync Checkpoint", func(t *testing.T) {
		checkpointReq := map[string]interface{}{
			"checkpoint": "checkpoint_test_123",
		}

		body, _ := json.Marshal(checkpointReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sync/checkpoint", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Sync checkpoint created successfully", response["message"])
	})

	t.Run("Get Offline Capabilities", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sync/capabilities", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "capabilities")
		assert.Contains(t, response, "version")

		capabilities := response["capabilities"].(map[string]interface{})
		assert.Contains(t, capabilities, "sync")
		assert.Contains(t, capabilities, "offline_first")
		assert.Contains(t, capabilities, "performance")
	})
}
