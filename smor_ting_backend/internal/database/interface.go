package database

import (
	"context"
	"time"

	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository defines the interface for data access operations
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error

	// OTP operations
	CreateOTP(ctx context.Context, otp *models.OTPRecord) error
	GetOTP(ctx context.Context, email, otpCode string) (*models.OTPRecord, error)
	MarkOTPAsUsed(ctx context.Context, id primitive.ObjectID) error
	// Test-only helper: latest OTP by email (unconsumed, unexpired)
	GetLatestOTPByEmail(ctx context.Context, email string) (*models.OTPRecord, error)

	// Service operations
	CreateService(ctx context.Context, service *models.Service) error
	GetServices(ctx context.Context, categoryID *primitive.ObjectID, location *models.Address, radius float64) ([]models.Service, error)

	// Booking operations
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetUserBookings(ctx context.Context, userID primitive.ObjectID) ([]models.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID primitive.ObjectID, status models.BookingStatus) error

	// Wallet operations
	UpdateWallet(ctx context.Context, userID primitive.ObjectID, transaction *models.Transaction) error

	// Offline-first sync operations
	GetUnsyncedData(ctx context.Context, userID primitive.ObjectID, lastSyncAt time.Time) (map[string]interface{}, error)
	SyncData(ctx context.Context, userID primitive.ObjectID, data map[string]interface{}) error

	// Enhanced sync operations with checkpoint and compression
	GetUnsyncedDataWithCheckpoint(ctx context.Context, req *models.SyncRequest) (*models.SyncResponse, error)
	GetChunkedUnsyncedData(ctx context.Context, req *models.ChunkedSyncRequest) (*models.ChunkedSyncResponse, error)
	GetSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.SyncStatus, error)
	UpdateSyncStatus(ctx context.Context, status *models.SyncStatus) error

	// Sync checkpoint operations
	CreateSyncCheckpoint(ctx context.Context, checkpoint *models.SyncCheckpoint) error
	GetSyncCheckpoint(ctx context.Context, userID primitive.ObjectID) (*models.SyncCheckpoint, error)
	UpdateSyncCheckpoint(ctx context.Context, checkpoint *models.SyncCheckpoint) error

	// Sync metrics operations
	CreateSyncMetrics(ctx context.Context, metrics *models.SyncMetrics) error
	GetRecentSyncMetrics(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncMetrics, error)

	// Background sync queue operations
	CreateSyncQueueItem(ctx context.Context, item *models.SyncQueueItem) error
	GetSyncQueueItem(ctx context.Context, itemID primitive.ObjectID) (*models.SyncQueueItem, error)
	UpdateSyncQueueItem(ctx context.Context, item *models.SyncQueueItem) error
	GetPendingSyncQueueItems(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncQueueItem, error)
	GetConflictQueueItems(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncQueueItem, error)
	CleanupCompletedQueueItems(ctx context.Context, olderThan time.Duration) (int64, error)

	// Background sync status operations
	GetBackgroundSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.BackgroundSyncStatus, error)
	UpdateBackgroundSyncStatus(ctx context.Context, status *models.BackgroundSyncStatus) error

	// Device session operations
	CreateDeviceSession(ctx context.Context, session *models.DeviceSession) error
	GetDeviceSession(ctx context.Context, sessionID string) (*models.DeviceSession, error)
	GetDeviceSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.DeviceSession, error)
	GetDeviceSessionByDeviceID(ctx context.Context, deviceID string) (*models.DeviceSession, error)
	GetUserDeviceSessions(ctx context.Context, userID string) ([]models.DeviceSession, error)
	UpdateDeviceSessionActivity(ctx context.Context, sessionID string) error
	RevokeDeviceSession(ctx context.Context, sessionID string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	RotateRefreshToken(ctx context.Context, sessionID string, newRefreshToken string) error
	CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) error

	// Security event operations
	LogSecurityEvent(ctx context.Context, event *models.SecurityEvent) error
	GetUserSecurityEvents(ctx context.Context, userID string, limit int) ([]models.SecurityEvent, error)
	GetSecurityEventsByType(ctx context.Context, userID string, eventType models.SecurityEventType, limit int) ([]models.SecurityEvent, error)

	// Setup operations
	SetupIndexes(ctx context.Context) error
	Close() error
}
