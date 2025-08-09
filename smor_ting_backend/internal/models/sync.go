package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SyncCheckpoint represents a sync checkpoint for efficient resuming
type SyncCheckpoint struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Checkpoint string             `json:"checkpoint" bson:"checkpoint"` // Base64 encoded state
	LastSyncAt time.Time          `json:"last_sync_at" bson:"last_sync_at"`
	Version    int                `json:"version" bson:"version"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// SyncMetrics tracks sync performance and statistics
type SyncMetrics struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id"`
	LastSyncAt        time.Time          `json:"last_sync_at" bson:"last_sync_at"`
	SyncDuration      time.Duration      `json:"sync_duration" bson:"sync_duration"`
	DataSize          int64              `json:"data_size" bson:"data_size"`             // in bytes
	CompressedSize    int64              `json:"compressed_size" bson:"compressed_size"` // in bytes
	RecordsSynced     int                `json:"records_synced" bson:"records_synced"`
	SyncSuccess       bool               `json:"sync_success" bson:"sync_success"`
	ErrorMessage      string             `json:"error_message,omitempty" bson:"error_message,omitempty"`
	NetworkType       string             `json:"network_type" bson:"network_type"`             // wifi, mobile, etc.
	ConnectionQuality string             `json:"connection_quality" bson:"connection_quality"` // good, poor, etc.
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
}

// SyncRequest represents a sync request with checkpoint support
type SyncRequest struct {
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Checkpoint  string             `json:"checkpoint,omitempty" bson:"checkpoint,omitempty"`
	LastSyncAt  time.Time          `json:"last_sync_at" bson:"last_sync_at"`
	Limit       int                `json:"limit" bson:"limit"`
	Compression bool               `json:"compression" bson:"compression"`
}

// SyncResponse represents a sync response with checkpoint and metrics
type SyncResponse struct {
	Data         map[string]interface{} `json:"data" bson:"data"`
	Checkpoint   string                 `json:"checkpoint" bson:"checkpoint"`
	LastSyncAt   time.Time              `json:"last_sync_at" bson:"last_sync_at"`
	HasMore      bool                   `json:"has_more" bson:"has_more"`
	Compressed   bool                   `json:"compressed" bson:"compressed"`
	DataSize     int64                  `json:"data_size" bson:"data_size"`
	RecordsCount int                    `json:"records_count" bson:"records_count"`
	SyncDuration time.Duration          `json:"sync_duration" bson:"sync_duration"`
}

// ChunkedSyncRequest represents a chunked sync request
type ChunkedSyncRequest struct {
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	ChunkIndex  int                `json:"chunk_index" bson:"chunk_index"`
	ChunkSize   int                `json:"chunk_size" bson:"chunk_size"`
	ResumeToken string             `json:"resume_token,omitempty" bson:"resume_token,omitempty"`
	Checkpoint  string             `json:"checkpoint,omitempty" bson:"checkpoint,omitempty"`
}

// ChunkedSyncResponse represents a chunked sync response
type ChunkedSyncResponse struct {
	Data         []interface{} `json:"data" bson:"data"`
	HasMore      bool          `json:"has_more" bson:"has_more"`
	NextChunk    int           `json:"next_chunk" bson:"next_chunk"`
	ResumeToken  string        `json:"resume_token" bson:"resume_token"`
	TotalChunks  int           `json:"total_chunks" bson:"total_chunks"`
	Checkpoint   string        `json:"checkpoint" bson:"checkpoint"`
	Compressed   bool          `json:"compressed" bson:"compressed"`
	DataSize     int64         `json:"data_size" bson:"data_size"`
	RecordsCount int           `json:"records_count" bson:"records_count"`
}

// SyncStatus represents the current sync status
type SyncStatus struct {
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsOnline        bool               `json:"is_online" bson:"is_online"`
	LastSyncAt      time.Time          `json:"last_sync_at" bson:"last_sync_at"`
	PendingChanges  int                `json:"pending_changes" bson:"pending_changes"`
	SyncInProgress  bool               `json:"sync_in_progress" bson:"sync_in_progress"`
	ConnectionType  string             `json:"connection_type" bson:"connection_type"`
	ConnectionSpeed string             `json:"connection_speed" bson:"connection_speed"`
	LastError       string             `json:"last_error,omitempty" bson:"last_error,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}
