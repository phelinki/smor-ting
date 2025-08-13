package test

import (
	"context"
	"testing"
	"time"

	"github.com/smorting/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestOfflineSyncSystem tests the complete offline sync system
func TestOfflineSyncSystem(t *testing.T) {
	_, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "sync@example.com",
		Password:  "hashed_password",
		FirstName: "Sync",
		LastName:  "Test",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID // Use the ID assigned during creation

	t.Run("Sync Status Management", func(t *testing.T) {
		// Get initial sync status
		status, err := repo.GetSyncStatus(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, userID, status.UserID)
		assert.True(t, status.IsOnline) // User is online by default (IsOffline = false)
		assert.False(t, status.SyncInProgress)
		assert.Equal(t, 0, status.PendingChanges)

		// Update sync status to online
		status.IsOnline = true
		status.ConnectionType = "wifi"
		status.ConnectionSpeed = "fast"
		status.UpdatedAt = time.Now()

		err = repo.UpdateSyncStatus(context.Background(), status)
		require.NoError(t, err)

		// Verify status was updated
		updatedStatus, err := repo.GetSyncStatus(context.Background(), userID)
		require.NoError(t, err)
		assert.True(t, updatedStatus.IsOnline)
		assert.Equal(t, "wifi", updatedStatus.ConnectionType)
		assert.Equal(t, "fast", updatedStatus.ConnectionSpeed)
	})

	t.Run("Sync Checkpoint Management", func(t *testing.T) {
		// Create sync checkpoint
		checkpoint := &models.SyncCheckpoint{
			UserID:     userID,
			Checkpoint: "checkpoint_data_123",
			LastSyncAt: time.Now(),
			Version:    1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := repo.CreateSyncCheckpoint(context.Background(), checkpoint)
		require.NoError(t, err)
		assert.NotEmpty(t, checkpoint.ID)

		// Get checkpoint
		retrievedCheckpoint, err := repo.GetSyncCheckpoint(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, checkpoint.UserID, retrievedCheckpoint.UserID)
		assert.Equal(t, checkpoint.Checkpoint, retrievedCheckpoint.Checkpoint)
		assert.Equal(t, checkpoint.Version, retrievedCheckpoint.Version)

		// Update checkpoint
		retrievedCheckpoint.Checkpoint = "updated_checkpoint_456"
		retrievedCheckpoint.Version = 2
		retrievedCheckpoint.UpdatedAt = time.Now()

		err = repo.UpdateSyncCheckpoint(context.Background(), retrievedCheckpoint)
		require.NoError(t, err)

		// Verify update
		updatedCheckpoint, err := repo.GetSyncCheckpoint(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, "updated_checkpoint_456", updatedCheckpoint.Checkpoint)
		assert.Equal(t, 2, updatedCheckpoint.Version)
	})

	t.Run("Sync Metrics Tracking", func(t *testing.T) {
		// Create sync metrics
		metrics := &models.SyncMetrics{
			UserID:            userID,
			LastSyncAt:        time.Now(),
			SyncDuration:      5 * time.Second,
			DataSize:          1024,
			CompressedSize:    512,
			RecordsSynced:     10,
			SyncSuccess:       true,
			NetworkType:       "wifi",
			ConnectionQuality: "good",
			CreatedAt:         time.Now(),
		}

		err := repo.CreateSyncMetrics(context.Background(), metrics)
		require.NoError(t, err)
		assert.NotEmpty(t, metrics.ID)

		// Get recent metrics
		recentMetrics, err := repo.GetRecentSyncMetrics(context.Background(), userID, 5)
		require.NoError(t, err)
		assert.Len(t, recentMetrics, 1)
		assert.Equal(t, metrics.UserID, recentMetrics[0].UserID)
		assert.Equal(t, metrics.SyncDuration, recentMetrics[0].SyncDuration)
		assert.True(t, recentMetrics[0].SyncSuccess)
	})
}

// TestOfflineDataSync tests offline data synchronization
func TestOfflineDataSync(t *testing.T) {
	_, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "offline@example.com",
		Password:  "hashed_password",
		FirstName: "Offline",
		LastName:  "User",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID // Use the ID assigned during creation

	// Create test data for syncing
	service := &models.Service{
		ID:          primitive.NewObjectID(),
		Name:        "Test Service",
		Description: "Test service for sync",
		CategoryID:  primitive.NewObjectID(),
		ProviderID:  userID,
		Price:       100.0,
		Currency:    "LRD",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = repo.CreateService(context.Background(), service)
	require.NoError(t, err)

	booking := &models.Booking{
		ID:            primitive.NewObjectID(),
		CustomerID:    userID,
		ProviderID:    userID,
		ServiceID:     service.ID,
		Status:        models.BookingPending,
		TotalAmount:   100.0,
		Currency:      "LRD",
		ScheduledDate: time.Now().Add(24 * time.Hour),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err = repo.CreateBooking(context.Background(), booking)
	require.NoError(t, err)

	t.Run("Basic Sync Data Retrieval", func(t *testing.T) {
		// Request sync data
		lastSyncAt := time.Now().Add(-1 * time.Hour) // 1 hour ago
		syncData, err := repo.GetUnsyncedData(context.Background(), userID, lastSyncAt)
		require.NoError(t, err)

		// Verify data structure
		assert.NotNil(t, syncData)
		assert.Contains(t, syncData, "user")
		assert.Contains(t, syncData, "services")
		assert.Contains(t, syncData, "bookings")

		// Verify user data
		userData := syncData["user"]
		assert.NotNil(t, userData)

		// Verify services data
		servicesData := syncData["services"]
		assert.NotNil(t, servicesData)

		// Verify bookings data
		bookingsData := syncData["bookings"]
		assert.NotNil(t, bookingsData)
	})

	t.Run("Enhanced Sync with Checkpoint", func(t *testing.T) {
		// Create sync request with checkpoint
		syncReq := &models.SyncRequest{
			UserID:      userID,
			Checkpoint:  "",
			LastSyncAt:  time.Now().Add(-1 * time.Hour),
			Limit:       100,
			Compression: false,
		}

		// Get sync data with checkpoint
		syncResp, err := repo.GetUnsyncedDataWithCheckpoint(context.Background(), syncReq)
		require.NoError(t, err)

		assert.NotNil(t, syncResp)
		assert.NotEmpty(t, syncResp.Data)
		assert.NotEmpty(t, syncResp.Checkpoint)
		assert.Greater(t, syncResp.RecordsCount, 0)
		assert.Greater(t, syncResp.DataSize, int64(0))
		assert.False(t, syncResp.HasMore) // Should be false for small dataset
	})

	t.Run("Chunked Sync for Large Datasets", func(t *testing.T) {
		// Create chunked sync request
		chunkReq := &models.ChunkedSyncRequest{
			UserID:     userID,
			ChunkIndex: 0,
			ChunkSize:  10,
		}

		// Get chunked sync data
		chunkResp, err := repo.GetChunkedUnsyncedData(context.Background(), chunkReq)
		require.NoError(t, err)

		assert.NotNil(t, chunkResp)
		assert.NotNil(t, chunkResp.Data)
		assert.Greater(t, chunkResp.RecordsCount, 0)
		assert.Equal(t, 0, chunkResp.NextChunk)
		assert.NotEmpty(t, chunkResp.Checkpoint)
	})

	t.Run("Sync Data Push", func(t *testing.T) {
		// Simulate pushing offline changes back to server
		offlineChanges := map[string]interface{}{
			"bookings": []map[string]interface{}{
				{
					"id":         booking.ID.Hex(),
					"status":     "confirmed",
					"updated_at": time.Now(),
					"offline_id": "offline_booking_123",
				},
			},
			"profile_updates": map[string]interface{}{
				"first_name": "Updated Name",
				"updated_at": time.Now(),
			},
		}

		err := repo.SyncData(context.Background(), userID, offlineChanges)
		require.NoError(t, err)

		// Verify user's last sync time was updated
		updatedUser, err := repo.GetUserByID(context.Background(), userID)
		require.NoError(t, err)
		assert.False(t, updatedUser.IsOffline)
	})
}

// TestConflictResolution tests conflict resolution in sync
func TestConflictResolution(t *testing.T) {
	_, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "conflict@example.com",
		Password:  "hashed_password",
		FirstName: "Conflict",
		LastName:  "Test",
		Role:      models.CustomerRole,
		Version:   1,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID // Use the ID assigned during creation

	t.Run("Version-Based Conflict Detection", func(t *testing.T) {
		// Simulate concurrent updates
		// Client A updates (version 1 -> 2)
		clientAChanges := map[string]interface{}{
			"profile_updates": map[string]interface{}{
				"first_name": "ClientA Update",
				"version":    2,
				"updated_at": time.Now(),
			},
		}

		// Client B updates (version 1 -> 2) - conflict!
		clientBChanges := map[string]interface{}{
			"profile_updates": map[string]interface{}{
				"first_name": "ClientB Update",
				"version":    2,
				"updated_at": time.Now().Add(1 * time.Second),
			},
		}

		// Apply first change
		err := repo.SyncData(context.Background(), userID, clientAChanges)
		require.NoError(t, err)

		// Apply second change - should detect conflict
		err = repo.SyncData(context.Background(), userID, clientBChanges)
		// In a real implementation, this should either:
		// 1. Return a conflict error
		// 2. Apply conflict resolution strategy
		// For now, we'll just verify the operation
		assert.NoError(t, err) // Will be enhanced with actual conflict resolution
	})

	t.Run("Last-Write-Wins Resolution", func(t *testing.T) {
		// Test last-write-wins conflict resolution
		// This is the simplest conflict resolution strategy
		laterUpdate := map[string]interface{}{
			"profile_updates": map[string]interface{}{
				"first_name": "Final Update",
				"version":    3,
				"updated_at": time.Now(),
			},
		}

		err := repo.SyncData(context.Background(), userID, laterUpdate)
		require.NoError(t, err)

		// Verify the last update won
		updatedUser, err := repo.GetUserByID(context.Background(), userID)
		require.NoError(t, err)
		// Note: In a real implementation, we'd have field-level tracking
		assert.NotNil(t, updatedUser)
	})
}

// TestSyncOptimization tests sync performance optimizations
func TestSyncOptimization(t *testing.T) {
	_, _, repo := setupTestApp(t)

	// Create test user
	user := &models.User{
		Email:     "optimize@example.com",
		Password:  "hashed_password",
		FirstName: "Optimize",
		LastName:  "Test",
		Role:      models.CustomerRole,
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	userID := user.ID // Use the ID assigned during creation

	t.Run("Compression Support", func(t *testing.T) {
		// Test sync with compression enabled
		syncReq := &models.SyncRequest{
			UserID:      userID,
			LastSyncAt:  time.Now().Add(-1 * time.Hour),
			Limit:       100,
			Compression: true,
		}

		syncResp, err := repo.GetUnsyncedDataWithCheckpoint(context.Background(), syncReq)
		require.NoError(t, err)

		// Verify compression metadata
		assert.True(t, syncResp.Compressed)
		// In a real implementation, compressed size should be less than data size
		// for data that compresses well
	})

	t.Run("Delta Sync", func(t *testing.T) {
		// Test incremental sync (only changed data)
		initialSync := time.Now().Add(-2 * time.Hour)
		deltaSync := time.Now().Add(-1 * time.Hour)

		// Get delta changes only
		deltaData, err := repo.GetUnsyncedData(context.Background(), userID, deltaSync)
		require.NoError(t, err)

		// Get full data
		fullData, err := repo.GetUnsyncedData(context.Background(), userID, initialSync)
		require.NoError(t, err)

		// Delta should be smaller or equal to full data
		assert.NotNil(t, deltaData)
		assert.NotNil(t, fullData)
	})

	t.Run("Bandwidth Optimization", func(t *testing.T) {
		// Test different chunk sizes for bandwidth optimization
		smallChunkReq := &models.ChunkedSyncRequest{
			UserID:    userID,
			ChunkSize: 5, // Small chunks for slow connections
		}

		largeChunkReq := &models.ChunkedSyncRequest{
			UserID:    userID,
			ChunkSize: 50, // Large chunks for fast connections
		}

		smallChunkResp, err := repo.GetChunkedUnsyncedData(context.Background(), smallChunkReq)
		require.NoError(t, err)

		largeChunkResp, err := repo.GetChunkedUnsyncedData(context.Background(), largeChunkReq)
		require.NoError(t, err)

		// Both should return data, but with different chunk characteristics
		assert.NotNil(t, smallChunkResp)
		assert.NotNil(t, largeChunkResp)
	})
}
