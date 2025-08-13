package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SyncQueueItemStatus represents the status of a sync queue item
type SyncQueueItemStatus string

const (
	SyncQueuePending    SyncQueueItemStatus = "pending"
	SyncQueueProcessing SyncQueueItemStatus = "processing"
	SyncQueueCompleted  SyncQueueItemStatus = "completed"
	SyncQueueFailed     SyncQueueItemStatus = "failed"
	SyncQueueRetrying   SyncQueueItemStatus = "retrying"
	SyncQueueCancelled  SyncQueueItemStatus = "cancelled"
)

// SyncQueueItemType represents the type of sync operation
type SyncQueueItemType string

const (
	SyncTypeUpload   SyncQueueItemType = "upload"
	SyncTypeDownload SyncQueueItemType = "download"
	SyncTypeConflict SyncQueueItemType = "conflict_resolution"
)

// SyncQueueItem represents an item in the background sync queue
type SyncQueueItem struct {
	ID             primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID         primitive.ObjectID     `json:"user_id" bson:"user_id"`
	Type           SyncQueueItemType      `json:"type" bson:"type"`
	Status         SyncQueueItemStatus    `json:"status" bson:"status"`
	Priority       int                    `json:"priority" bson:"priority"` // Higher number = higher priority
	Data           map[string]interface{} `json:"data" bson:"data"`
	ConflictData   *ConflictResolution    `json:"conflict_data,omitempty" bson:"conflict_data,omitempty"`
	RetryCount     int                    `json:"retry_count" bson:"retry_count"`
	MaxRetries     int                    `json:"max_retries" bson:"max_retries"`
	NextRetryAt    time.Time              `json:"next_retry_at" bson:"next_retry_at"`
	LastAttemptAt  *time.Time             `json:"last_attempt_at,omitempty" bson:"last_attempt_at,omitempty"`
	LastError      string                 `json:"last_error,omitempty" bson:"last_error,omitempty"`
	ProcessingNode string                 `json:"processing_node,omitempty" bson:"processing_node,omitempty"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
	CompletedAt    *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
}

// ConflictResolution represents conflict resolution data
type ConflictResolution struct {
	ConflictType       string                 `json:"conflict_type" bson:"conflict_type"`
	ClientVersion      int                    `json:"client_version" bson:"client_version"`
	ServerVersion      int                    `json:"server_version" bson:"server_version"`
	ClientData         map[string]interface{} `json:"client_data" bson:"client_data"`
	ServerData         map[string]interface{} `json:"server_data" bson:"server_data"`
	ResolutionStrategy string                 `json:"resolution_strategy" bson:"resolution_strategy"` // "client_wins", "server_wins", "merge", "manual"
	ResolvedData       map[string]interface{} `json:"resolved_data,omitempty" bson:"resolved_data,omitempty"`
	RequiresUserInput  bool                   `json:"requires_user_input" bson:"requires_user_input"`
	UserDecision       string                 `json:"user_decision,omitempty" bson:"user_decision,omitempty"`
}

// SyncQueueConfig represents configuration for the sync queue
type SyncQueueConfig struct {
	MaxRetries        int           `json:"max_retries" bson:"max_retries"`
	BaseRetryDelay    time.Duration `json:"base_retry_delay" bson:"base_retry_delay"`
	MaxRetryDelay     time.Duration `json:"max_retry_delay" bson:"max_retry_delay"`
	BackoffMultiplier float64       `json:"backoff_multiplier" bson:"backoff_multiplier"`
	BatchSize         int           `json:"batch_size" bson:"batch_size"`
	ProcessorCount    int           `json:"processor_count" bson:"processor_count"`
}

// BackgroundSyncStatus represents the overall background sync status
type BackgroundSyncStatus struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsEnabled        bool               `json:"is_enabled" bson:"is_enabled"`
	LastSyncAt       time.Time          `json:"last_sync_at" bson:"last_sync_at"`
	PendingItems     int                `json:"pending_items" bson:"pending_items"`
	FailedItems      int                `json:"failed_items" bson:"failed_items"`
	ConflictItems    int                `json:"conflict_items" bson:"conflict_items"`
	AutoRetryEnabled bool               `json:"auto_retry_enabled" bson:"auto_retry_enabled"`
	NextScheduledRun time.Time          `json:"next_scheduled_run" bson:"next_scheduled_run"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

// RetryPolicy represents the retry policy for sync operations
type RetryPolicy struct {
	MaxRetries   int           `json:"max_retries"`
	BaseDelay    time.Duration `json:"base_delay"`
	MaxDelay     time.Duration `json:"max_delay"`
	Multiplier   float64       `json:"multiplier"`
	RandomJitter bool          `json:"random_jitter"`
}

// GetDefaultRetryPolicy returns the default retry policy
func GetDefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries:   3,
		BaseDelay:    1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		RandomJitter: true,
	}
}

// CalculateNextRetry calculates the next retry time based on retry count
func (r *RetryPolicy) CalculateNextRetry(retryCount int) time.Time {
	if retryCount >= r.MaxRetries {
		return time.Time{} // No more retries
	}

	delay := time.Duration(float64(r.BaseDelay) * float64(retryCount+1) * r.Multiplier)
	if delay > r.MaxDelay {
		delay = r.MaxDelay
	}

	// Add random jitter if enabled
	if r.RandomJitter {
		jitter := time.Duration(float64(delay) * 0.1) // 10% jitter
		delay += time.Duration(float64(jitter) * (2.0*float64(time.Now().UnixNano()%1000)/1000.0 - 1.0))
	}

	return time.Now().Add(delay)
}

// ShouldRetry determines if an item should be retried
func (item *SyncQueueItem) ShouldRetry() bool {
	return item.Status == SyncQueueFailed &&
		item.RetryCount < item.MaxRetries &&
		time.Now().After(item.NextRetryAt)
}

// MarkForRetry marks an item for retry with exponential backoff
func (item *SyncQueueItem) MarkForRetry(error string, retryPolicy RetryPolicy) {
	item.Status = SyncQueueRetrying
	item.RetryCount++
	item.LastError = error
	item.NextRetryAt = retryPolicy.CalculateNextRetry(item.RetryCount - 1)
	item.UpdatedAt = time.Now()

	if item.RetryCount >= item.MaxRetries {
		item.Status = SyncQueueFailed
	}
}

// MarkCompleted marks an item as completed
func (item *SyncQueueItem) MarkCompleted() {
	item.Status = SyncQueueCompleted
	now := time.Now()
	item.CompletedAt = &now
	item.UpdatedAt = now
}

// MarkFailed marks an item as permanently failed
func (item *SyncQueueItem) MarkFailed(error string) {
	item.Status = SyncQueueFailed
	item.LastError = error
	item.UpdatedAt = time.Now()
}
