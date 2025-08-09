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

	// Setup operations
	SetupIndexes(ctx context.Context) error
	Close() error
}
