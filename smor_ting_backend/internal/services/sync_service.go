package services

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// SyncService handles enhanced sync operations with checkpoint, compression, and metrics
type SyncService struct {
	db     *mongo.Database
	logger *logger.Logger
}

// NewSyncService creates a new sync service
func NewSyncService(db *mongo.Database, logger *logger.Logger) *SyncService {
	return &SyncService{
		db:     db,
		logger: logger,
	}
}

// GetUnsyncedDataWithCheckpoint gets unsynced data with checkpoint support
func (s *SyncService) GetUnsyncedDataWithCheckpoint(ctx context.Context, req *models.SyncRequest) (*models.SyncResponse, error) {
	startTime := time.Now()

	s.logger.Info("Starting sync with checkpoint",
		zap.String("user_id", req.UserID.Hex()),
		zap.String("checkpoint", req.Checkpoint),
		zap.Bool("compression", req.Compression),
		zap.Int("limit", req.Limit),
	)

	// Get or create checkpoint
	checkpoint, err := s.getOrCreateCheckpoint(ctx, req.UserID, req.Checkpoint)
	if err != nil {
		s.logger.Error("Failed to get checkpoint", err,
			zap.String("user_id", req.UserID.Hex()),
		)
		return nil, fmt.Errorf("failed to get checkpoint: %w", err)
	}

	// Get unsynced data
	data, err := s.getUnsyncedData(ctx, req.UserID, checkpoint.LastSyncAt, req.Limit)
	if err != nil {
		s.logger.Error("Failed to get unsynced data", err,
			zap.String("user_id", req.UserID.Hex()),
		)
		return nil, fmt.Errorf("failed to get unsynced data: %w", err)
	}

	// Calculate data size
	dataSize := s.calculateDataSize(data)

	// Compress data if requested
	var compressed bool
	var compressedData []byte
	if req.Compression {
		compressedData, err = s.compressData(data)
		if err != nil {
			s.logger.Error("Failed to compress data", err,
				zap.String("user_id", req.UserID.Hex()),
			)
			return nil, fmt.Errorf("failed to compress data: %w", err)
		}
		compressed = true
	}

	// Create new checkpoint
	newCheckpoint, err := s.createCheckpoint(ctx, req.UserID, data)
	if err != nil {
		s.logger.Error("Failed to create checkpoint", err,
			zap.String("user_id", req.UserID.Hex()),
		)
		return nil, fmt.Errorf("failed to create checkpoint: %w", err)
	}

	// Calculate sync duration
	syncDuration := time.Since(startTime)

	// Log sync metrics
	s.logSyncMetrics(ctx, req.UserID, syncDuration, dataSize, int64(len(compressedData)), len(data), true, "")

	s.logger.Info("Sync completed successfully",
		zap.String("user_id", req.UserID.Hex()),
		zap.Duration("duration", syncDuration),
		zap.Int64("data_size", dataSize),
		zap.Int64("compressed_size", int64(len(compressedData))),
		zap.Int("records_count", len(data)),
		zap.Bool("compressed", compressed),
	)

	response := &models.SyncResponse{
		Data:         data,
		Checkpoint:   newCheckpoint.Checkpoint,
		LastSyncAt:   time.Now(),
		HasMore:      len(data) >= req.Limit,
		Compressed:   compressed,
		DataSize:     dataSize,
		RecordsCount: len(data),
		SyncDuration: syncDuration,
	}

	return response, nil
}

// GetChunkedUnsyncedData gets unsynced data in chunks
func (s *SyncService) GetChunkedUnsyncedData(ctx context.Context, req *models.ChunkedSyncRequest) (*models.ChunkedSyncResponse, error) {
	startTime := time.Now()

	s.logger.Info("Starting chunked sync",
		zap.String("user_id", req.UserID.Hex()),
		zap.Int("chunk_index", req.ChunkIndex),
		zap.Int("chunk_size", req.ChunkSize),
	)

	// Get checkpoint
	checkpoint, err := s.getOrCreateCheckpoint(ctx, req.UserID, req.Checkpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpoint: %w", err)
	}

	// Get chunked data
	data, hasMore, err := s.getChunkedData(ctx, req.UserID, checkpoint.LastSyncAt, req.ChunkIndex, req.ChunkSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunked data: %w", err)
	}

	// Convert data to map for compression and checkpoint
	dataMap := make(map[string]interface{})
	for i, item := range data {
		dataMap[fmt.Sprintf("item_%d", i)] = item
	}

	// Compress data
	compressedData, err := s.compressData(dataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %w", err)
	}

	// Create new checkpoint
	newCheckpoint, err := s.createCheckpoint(ctx, req.UserID, dataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkpoint: %w", err)
	}

	// Generate resume token
	resumeToken, err := s.generateResumeToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate resume token: %w", err)
	}

	syncDuration := time.Since(startTime)
	dataSize := s.calculateDataSize(dataMap)

	// Log sync metrics
	s.logSyncMetrics(ctx, req.UserID, syncDuration, dataSize, int64(len(compressedData)), len(data), true, "")

	response := &models.ChunkedSyncResponse{
		Data:         data,
		HasMore:      hasMore,
		NextChunk:    req.ChunkIndex + 1,
		ResumeToken:  resumeToken,
		TotalChunks:  -1, // Will be calculated if needed
		Checkpoint:   newCheckpoint.Checkpoint,
		Compressed:   true,
		DataSize:     dataSize,
		RecordsCount: len(data),
	}

	return response, nil
}

// getOrCreateCheckpoint gets or creates a checkpoint for a user
func (s *SyncService) getOrCreateCheckpoint(ctx context.Context, userID primitive.ObjectID, checkpointStr string) (*models.SyncCheckpoint, error) {
	collection := s.db.Collection("sync_checkpoints")

	if checkpointStr != "" {
		// Try to find existing checkpoint
		var checkpoint models.SyncCheckpoint
		err := collection.FindOne(ctx, bson.M{
			"user_id":    userID,
			"checkpoint": checkpointStr,
		}).Decode(&checkpoint)

		if err == nil {
			return &checkpoint, nil
		}
	}

	// Create new checkpoint
	checkpoint := &models.SyncCheckpoint{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		Checkpoint: s.generateCheckpoint(),
		LastSyncAt: time.Now(),
		Version:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := collection.InsertOne(ctx, checkpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkpoint: %w", err)
	}

	return checkpoint, nil
}

// createCheckpoint creates a new checkpoint based on current data state
func (s *SyncService) createCheckpoint(ctx context.Context, userID primitive.ObjectID, data map[string]interface{}) (*models.SyncCheckpoint, error) {
	// Create checkpoint data
	checkpointData := map[string]interface{}{
		"user_id":   userID.Hex(),
		"timestamp": time.Now().Unix(),
		"data_keys": s.getDataKeys(data),
		"version":   1,
	}

	// Encode checkpoint data
	checkpointJSON, err := json.Marshal(checkpointData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal checkpoint data: %w", err)
	}

	checkpoint := base64.StdEncoding.EncodeToString(checkpointJSON)

	// Save checkpoint
	collection := s.db.Collection("sync_checkpoints")
	checkpointDoc := &models.SyncCheckpoint{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		Checkpoint: checkpoint,
		LastSyncAt: time.Now(),
		Version:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = collection.InsertOne(ctx, checkpointDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to save checkpoint: %w", err)
	}

	return checkpointDoc, nil
}

// getUnsyncedData gets unsynced data for a user
func (s *SyncService) getUnsyncedData(ctx context.Context, userID primitive.ObjectID, lastSyncAt time.Time, limit int) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// Get unsynced bookings
	bookingsCollection := s.db.Collection("bookings")
	limit64 := int64(limit)
	bookingsCursor, err := bookingsCollection.Find(ctx, bson.M{
		"customer_id":  userID,
		"last_sync_at": bson.M{"$gt": lastSyncAt},
	}, &options.FindOptions{
		Limit: &limit64,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get unsynced bookings: %w", err)
	}
	defer bookingsCursor.Close(ctx)

	var bookings []models.Booking
	if err = bookingsCursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("failed to decode bookings: %w", err)
	}

	// Get unsynced user data
	userCollection := s.db.Collection("users")
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	data["bookings"] = bookings
	data["user"] = user

	return data, nil
}

// getChunkedData gets data in chunks
func (s *SyncService) getChunkedData(ctx context.Context, userID primitive.ObjectID, lastSyncAt time.Time, chunkIndex, chunkSize int) ([]interface{}, bool, error) {
	// For now, return all data as a single chunk
	// In a real implementation, you would implement proper pagination
	data, err := s.getUnsyncedData(ctx, userID, lastSyncAt, chunkSize)
	if err != nil {
		return nil, false, err
	}

	// Convert to slice
	var result []interface{}
	for _, v := range data {
		result = append(result, v)
	}

	return result, false, nil
}

// compressData compresses data using gzip
func (s *SyncService) compressData(data map[string]interface{}) ([]byte, error) {
	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Compress with gzip
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(jsonData); err != nil {
		return nil, fmt.Errorf("failed to write to gzip: %w", err)
	}

	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip: %w", err)
	}

	return buf.Bytes(), nil
}

// DecompressData decompresses gzip data
func (s *SyncService) DecompressData(compressedData []byte) (map[string]interface{}, error) {
	// Create gzip reader
	gz, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gz.Close()

	// Read decompressed data
	decompressedData, err := io.ReadAll(gz)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	// Unmarshal JSON
	var data map[string]interface{}
	if err := json.Unmarshal(decompressedData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return data, nil
}

// calculateDataSize calculates the size of data in bytes
func (s *SyncService) calculateDataSize(data map[string]interface{}) int64 {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0
	}
	return int64(len(jsonData))
}

// generateCheckpoint generates a unique checkpoint string
func (s *SyncService) generateCheckpoint() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

// generateResumeToken generates a unique resume token
func (s *SyncService) generateResumeToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// getDataKeys extracts keys from data for checkpoint
func (s *SyncService) getDataKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

// logSyncMetrics logs sync metrics for monitoring
func (s *SyncService) logSyncMetrics(ctx context.Context, userID primitive.ObjectID, duration time.Duration, dataSize, compressedSize int64, recordsCount int, success bool, errorMsg string) {
	metrics := &models.SyncMetrics{
		ID:                primitive.NewObjectID(),
		UserID:            userID,
		LastSyncAt:        time.Now(),
		SyncDuration:      duration,
		DataSize:          dataSize,
		CompressedSize:    compressedSize,
		RecordsSynced:     recordsCount,
		SyncSuccess:       success,
		ErrorMessage:      errorMsg,
		NetworkType:       "unknown", // Will be set by client
		ConnectionQuality: "unknown", // Will be set by client
		CreatedAt:         time.Now(),
	}

	// Save metrics to database
	collection := s.db.Collection("sync_metrics")
	_, err := collection.InsertOne(ctx, metrics)
	if err != nil {
		s.logger.Error("Failed to save sync metrics", err,
			zap.String("user_id", userID.Hex()),
		)
	}

	// Log metrics
	s.logger.Info("Sync metrics",
		zap.String("user_id", userID.Hex()),
		zap.Duration("duration", duration),
		zap.Int64("data_size", dataSize),
		zap.Int64("compressed_size", compressedSize),
		zap.Int("records_count", recordsCount),
		zap.Bool("success", success),
		zap.String("error", errorMsg),
	)
}

// GetSyncStatus gets the current sync status for a user
func (s *SyncService) GetSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.SyncStatus, error) {
	// Get user's last sync info
	userCollection := s.db.Collection("users")
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Count pending changes
	bookingsCollection := s.db.Collection("bookings")
	pendingCount, err := bookingsCollection.CountDocuments(ctx, bson.M{
		"customer_id":  userID,
		"last_sync_at": bson.M{"$lt": user.LastSyncAt},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to count pending changes: %w", err)
	}

	status := &models.SyncStatus{
		UserID:          userID,
		IsOnline:        !user.IsOffline,
		LastSyncAt:      user.LastSyncAt,
		PendingChanges:  int(pendingCount),
		SyncInProgress:  false,     // Will be set by client
		ConnectionType:  "unknown", // Will be set by client
		ConnectionSpeed: "unknown", // Will be set by client
		UpdatedAt:       time.Now(),
	}

	return status, nil
}
