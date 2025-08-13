package services

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// SyncService handles offline sync operations
type SyncService struct {
	repo         database.Repository
	auditService *AuditService
	logger       *zap.Logger
}

// NewSyncService creates a new sync service
func NewSyncService(repo database.Repository, auditService *AuditService, logger *zap.Logger) *SyncService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SyncService{
		repo:         repo,
		auditService: auditService,
		logger:       logger,
	}
}

// GetSyncStatus returns the current sync status for a user
func (s *SyncService) GetSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.SyncStatus, error) {
	status, err := s.repo.GetSyncStatus(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get sync status", zap.Error(err), zap.String("userID", userID.Hex()))
		return nil, err
	}

	// Update the status to reflect current connectivity
	status.UpdatedAt = time.Now()

	return status, nil
}

// UpdateSyncStatus updates the sync status for a user
func (s *SyncService) UpdateSyncStatus(ctx context.Context, userID primitive.ObjectID, update *models.SyncStatus) error {
	update.UserID = userID
	update.UpdatedAt = time.Now()

	err := s.repo.UpdateSyncStatus(ctx, update)
	if err != nil {
		s.logger.Error("Failed to update sync status", zap.Error(err), zap.String("userID", userID.Hex()))
		return err
	}

	s.logger.Info("Sync status updated",
		zap.String("userID", userID.Hex()),
		zap.Bool("isOnline", update.IsOnline),
		zap.String("connectionType", update.ConnectionType))

	return nil
}

// SyncUp processes offline changes and merges them with server data
func (s *SyncService) SyncUp(ctx context.Context, userID primitive.ObjectID, changes map[string]interface{}) error {
	start := time.Now()

	// Start sync
	err := s.markSyncInProgress(ctx, userID, true)
	if err != nil {
		return err
	}

	defer func() {
		// Mark sync as complete
		s.markSyncInProgress(ctx, userID, false)
	}()

	// Apply changes
	err = s.repo.SyncData(ctx, userID, changes)
	if err != nil {
		s.logger.Error("Failed to sync data up", zap.Error(err), zap.String("userID", userID.Hex()))

		// Record failed sync metrics
		s.recordSyncMetrics(ctx, userID, false, 0, time.Since(start), err.Error())
		return err
	}

	// Record successful sync metrics
	recordCount := s.countRecords(changes)
	s.recordSyncMetrics(ctx, userID, true, recordCount, time.Since(start), "")

	// Log security event for data sync
	securityEvent := &models.SecurityEvent{
		UserID:    userID,
		EventType: "data_sync",
		Metadata: map[string]interface{}{
			"direction":     "up",
			"record_count":  recordCount,
			"sync_duration": time.Since(start).Milliseconds(),
		},
		Timestamp: time.Now(),
	}
	if err := s.repo.LogSecurityEvent(ctx, securityEvent); err != nil {
		s.logger.Warn("Failed to log sync security event", zap.Error(err))
	}

	s.logger.Info("Sync up completed",
		zap.String("userID", userID.Hex()),
		zap.Int("recordCount", recordCount),
		zap.Duration("duration", time.Since(start)))

	return nil
}

// SyncDown retrieves server changes for the user
func (s *SyncService) SyncDown(ctx context.Context, req *models.SyncRequest) (*models.SyncResponse, error) {
	start := time.Now()

	// Start sync
	err := s.markSyncInProgress(ctx, req.UserID, true)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Mark sync as complete
		s.markSyncInProgress(ctx, req.UserID, false)
	}()

	// Get unsynced data with checkpoint
	response, err := s.repo.GetUnsyncedDataWithCheckpoint(ctx, req)
	if err != nil {
		s.logger.Error("Failed to sync data down", zap.Error(err), zap.String("userID", req.UserID.Hex()))

		// Record failed sync metrics
		s.recordSyncMetrics(ctx, req.UserID, false, 0, time.Since(start), err.Error())
		return nil, err
	}

	// Update response with actual duration
	response.SyncDuration = time.Since(start)

	// Record successful sync metrics
	s.recordSyncMetrics(ctx, req.UserID, true, response.RecordsCount, time.Since(start), "")

	// Log security event for data sync
	securityEvent := &models.SecurityEvent{
		UserID:    req.UserID,
		EventType: "data_sync",
		Metadata: map[string]interface{}{
			"direction":     "down",
			"record_count":  response.RecordsCount,
			"data_size":     response.DataSize,
			"compressed":    response.Compressed,
			"sync_duration": response.SyncDuration.Milliseconds(),
		},
		Timestamp: time.Now(),
	}
	if err := s.repo.LogSecurityEvent(ctx, securityEvent); err != nil {
		s.logger.Warn("Failed to log sync security event", zap.Error(err))
	}

	s.logger.Info("Sync down completed",
		zap.String("userID", req.UserID.Hex()),
		zap.Int("recordCount", response.RecordsCount),
		zap.Duration("duration", response.SyncDuration))

	return response, nil
}

// SyncDownChunked retrieves server changes in chunks for large datasets
func (s *SyncService) SyncDownChunked(ctx context.Context, req *models.ChunkedSyncRequest) (*models.ChunkedSyncResponse, error) {
	start := time.Now()

	response, err := s.repo.GetChunkedUnsyncedData(ctx, req)
	if err != nil {
		s.logger.Error("Failed to sync chunked data down", zap.Error(err), zap.String("userID", req.UserID.Hex()))
		return nil, err
	}

	// Log chunked sync
	s.logger.Info("Chunked sync completed",
		zap.String("userID", req.UserID.Hex()),
		zap.Int("chunkIndex", req.ChunkIndex),
		zap.Int("recordCount", response.RecordsCount),
		zap.Bool("hasMore", response.HasMore),
		zap.Duration("duration", time.Since(start)))

	return response, nil
}

// Backwards-compat wrappers used by callers expecting these method names
// GetUnsyncedDataWithCheckpoint delegates to SyncDown
func (s *SyncService) GetUnsyncedDataWithCheckpoint(ctx context.Context, req *models.SyncRequest) (*models.SyncResponse, error) {
	return s.SyncDown(ctx, req)
}

// GetChunkedUnsyncedData delegates to SyncDownChunked
func (s *SyncService) GetChunkedUnsyncedData(ctx context.Context, req *models.ChunkedSyncRequest) (*models.ChunkedSyncResponse, error) {
	return s.SyncDownChunked(ctx, req)
}

// DecompressData provides simple gzip decompression and JSON decoding helper for API handler
func (s *SyncService) DecompressData(compressed []byte) (interface{}, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Try to decode JSON payload into generic structure
	var anyJSON interface{}
	if err := json.Unmarshal(decompressed, &anyJSON); err == nil {
		return anyJSON, nil
	}

	// Fallback: return raw string
	return string(decompressed), nil
}

// GetSyncMetrics returns recent sync metrics for a user
func (s *SyncService) GetSyncMetrics(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncMetrics, error) {
	metrics, err := s.repo.GetRecentSyncMetrics(ctx, userID, limit)
	if err != nil {
		s.logger.Error("Failed to get sync metrics", zap.Error(err), zap.String("userID", userID.Hex()))
		return nil, err
	}

	return metrics, nil
}

// CreateSyncCheckpoint creates a sync checkpoint for resumable sync
func (s *SyncService) CreateSyncCheckpoint(ctx context.Context, userID primitive.ObjectID, checkpoint string) error {
	checkpointModel := &models.SyncCheckpoint{
		UserID:     userID,
		Checkpoint: checkpoint,
		LastSyncAt: time.Now(),
		Version:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := s.repo.CreateSyncCheckpoint(ctx, checkpointModel)
	if err != nil {
		s.logger.Error("Failed to create sync checkpoint", zap.Error(err), zap.String("userID", userID.Hex()))
		return err
	}

	return nil
}

// Helper methods

func (s *SyncService) markSyncInProgress(ctx context.Context, userID primitive.ObjectID, inProgress bool) error {
	status, err := s.repo.GetSyncStatus(ctx, userID)
	if err != nil {
		return err
	}

	status.SyncInProgress = inProgress
	status.UpdatedAt = time.Now()

	return s.repo.UpdateSyncStatus(ctx, status)
}

func (s *SyncService) recordSyncMetrics(ctx context.Context, userID primitive.ObjectID, success bool, recordCount int, duration time.Duration, errorMsg string) {
	metrics := &models.SyncMetrics{
		UserID:        userID,
		LastSyncAt:    time.Now(),
		SyncDuration:  duration,
		RecordsSynced: recordCount,
		SyncSuccess:   success,
		ErrorMessage:  errorMsg,
		CreatedAt:     time.Now(),
	}

	err := s.repo.CreateSyncMetrics(ctx, metrics)
	if err != nil {
		s.logger.Warn("Failed to record sync metrics", zap.Error(err))
	}
}

func (s *SyncService) countRecords(changes map[string]interface{}) int {
	count := 0
	for key, value := range changes {
		switch key {
		case "bookings":
			if bookings, ok := value.([]interface{}); ok {
				count += len(bookings)
			}
		case "services":
			if services, ok := value.([]interface{}); ok {
				count += len(services)
			}
		case "profile_updates":
			count += 1
		default:
			if slice, ok := value.([]interface{}); ok {
				count += len(slice)
			} else {
				count += 1
			}
		}
	}
	return count
}
