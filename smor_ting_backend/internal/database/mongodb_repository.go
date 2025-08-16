package database

import (
	"context"
	"fmt"
	"time"

	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MongoDBRepository implements the data access layer with offline-first capabilities
type MongoDBRepository struct {
	db     *mongo.Database
	logger *logger.Logger
}

// NewMongoDBRepository creates a new MongoDB repository
func NewMongoDBRepository(db *mongo.Database, logger *logger.Logger) *MongoDBRepository {
	return &MongoDBRepository{
		db:     db,
		logger: logger,
	}
}

// User operations with offline-first support
func (r *MongoDBRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastSyncAt = time.Now()
	user.Version = 1
	user.IsOffline = false

	// Initialize wallet with default currency if not already set
	if user.Wallet.Currency == "" {
		user.Wallet = models.Wallet{
			Balance:     0,
			Currency:    "USD",
			LastUpdated: time.Now(),
		}
	} else {
		// Preserve existing wallet settings but ensure LastUpdated is set
		if user.Wallet.LastUpdated.IsZero() {
			user.Wallet.LastUpdated = time.Now()
		}
	}

	collection := r.db.Collection("users")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("User created successfully", zap.String("user_id", user.ID.Hex()))
	return nil
}

func (r *MongoDBRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := r.db.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *MongoDBRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	collection := r.db.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *MongoDBRepository) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	user.LastSyncAt = time.Now()
	user.Version++

	collection := r.db.Collection("users")
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// OTP operations with TTL index support
func (r *MongoDBRepository) CreateOTP(ctx context.Context, otp *models.OTPRecord) error {
	otp.ID = primitive.NewObjectID()
	otp.CreatedAt = time.Now()
	otp.IsUsed = false

	collection := r.db.Collection("otp_records")
	_, err := collection.InsertOne(ctx, otp)
	if err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	return nil
}

func (r *MongoDBRepository) GetOTP(ctx context.Context, email, otpCode string) (*models.OTPRecord, error) {
	collection := r.db.Collection("otp_records")

	var otp models.OTPRecord
	err := collection.FindOne(ctx, bson.M{
		"email":      email,
		"otp":        otpCode,
		"is_used":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&otp)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("invalid or expired OTP")
		}
		return nil, fmt.Errorf("failed to get OTP: %w", err)
	}

	return &otp, nil
}

func (r *MongoDBRepository) MarkOTPAsUsed(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection("otp_records")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_used": true}})
	if err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	return nil
}

// GetLatestOTPByEmail returns the most recent unused, unexpired OTP for an email
func (r *MongoDBRepository) GetLatestOTPByEmail(ctx context.Context, email string) (*models.OTPRecord, error) {
	collection := r.db.Collection("otp_records")
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var otp models.OTPRecord
	err := collection.FindOne(ctx, bson.M{
		"email":      email,
		"is_used":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}, opts).Decode(&otp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no otp found")
		}
		return nil, fmt.Errorf("failed to get latest otp: %w", err)
	}
	return &otp, nil
}

// Service operations with embedded documents
func (r *MongoDBRepository) CreateService(ctx context.Context, service *models.Service) error {
	service.ID = primitive.NewObjectID()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.LastSyncAt = time.Now()
	service.Version = 1

	collection := r.db.Collection("services")
	_, err := collection.InsertOne(ctx, service)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	return nil
}

func (r *MongoDBRepository) GetServices(ctx context.Context, categoryID *primitive.ObjectID, location *models.Address, radius float64) ([]models.Service, error) {
	collection := r.db.Collection("services")

	filter := bson.M{"is_active": true}

	if categoryID != nil {
		filter["category_id"] = *categoryID
	}

	if location != nil && radius > 0 {
		// Geospatial query for nearby services
		filter["location"] = bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{location.Longitude, location.Latitude},
				},
				"$maxDistance": radius * 1000, // Convert km to meters
			},
		}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}
	defer cursor.Close(ctx)

	var services []models.Service
	if err = cursor.All(ctx, &services); err != nil {
		return nil, fmt.Errorf("failed to decode services: %w", err)
	}

	return services, nil
}

// Booking operations with embedded documents
func (r *MongoDBRepository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	booking.ID = primitive.NewObjectID()
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	booking.LastSyncAt = time.Now()
	booking.Version = 1

	// Start a session for transaction
	session, err := r.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	// Use transaction for booking creation
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Create booking
		collection := r.db.Collection("bookings")
		_, err := collection.InsertOne(sessCtx, booking)
		if err != nil {
			return nil, err
		}

		// Update user's bookings array
		userCollection := r.db.Collection("users")
		_, err = userCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": booking.CustomerID},
			bson.M{"$push": bson.M{"bookings": booking}},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}

func (r *MongoDBRepository) GetUserBookings(ctx context.Context, userID primitive.ObjectID) ([]models.Booking, error) {
	collection := r.db.Collection("bookings")

	cursor, err := collection.Find(ctx, bson.M{"customer_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("failed to decode bookings: %w", err)
	}

	return bookings, nil
}

func (r *MongoDBRepository) UpdateBookingStatus(ctx context.Context, bookingID primitive.ObjectID, status models.BookingStatus) error {
	collection := r.db.Collection("bookings")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": bookingID},
		bson.M{
			"$set": bson.M{
				"status":       status,
				"updated_at":   time.Now(),
				"last_sync_at": time.Now(),
			},
			"$inc": bson.M{"version": 1},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	return nil
}

// Wallet operations
func (r *MongoDBRepository) UpdateWallet(ctx context.Context, userID primitive.ObjectID, transaction *models.Transaction) error {
	transaction.ID = primitive.NewObjectID()
	transaction.CreatedAt = time.Now()

	collection := r.db.Collection("users")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$push": bson.M{"wallet.transactions": transaction},
			"$inc":  bson.M{"wallet.balance": transaction.Amount},
			"$set":  bson.M{"wallet.last_updated": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	return nil
}

// Offline-first sync operations
func (r *MongoDBRepository) GetUnsyncedData(ctx context.Context, userID primitive.ObjectID, lastSyncAt time.Time) (map[string]interface{}, error) {
	// Get unsynced bookings
	bookingsCollection := r.db.Collection("bookings")
	bookingsCursor, err := bookingsCollection.Find(ctx, bson.M{
		"customer_id":  userID,
		"last_sync_at": bson.M{"$gt": lastSyncAt},
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
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return map[string]interface{}{
		"bookings": bookings,
		"user":     user,
	}, nil
}

func (r *MongoDBRepository) SyncData(ctx context.Context, userID primitive.ObjectID, data map[string]interface{}) error {
	// Update user's last sync time
	collection := r.db.Collection("users")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"last_sync_at": time.Now(),
				"is_offline":   false,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update sync time: %w", err)
	}

	return nil
}

// Enhanced sync operations with checkpoint and compression
func (r *MongoDBRepository) GetUnsyncedDataWithCheckpoint(ctx context.Context, req *models.SyncRequest) (*models.SyncResponse, error) {
	// This method delegates to the sync service
	// In a real implementation, you would implement this directly in the repository
	// For now, we'll return a placeholder response
	return &models.SyncResponse{
		Data:         make(map[string]interface{}),
		Checkpoint:   "",
		LastSyncAt:   time.Now(),
		HasMore:      false,
		Compressed:   false,
		DataSize:     0,
		RecordsCount: 0,
		SyncDuration: 0,
	}, nil
}

func (r *MongoDBRepository) GetChunkedUnsyncedData(ctx context.Context, req *models.ChunkedSyncRequest) (*models.ChunkedSyncResponse, error) {
	// This method delegates to the sync service
	// In a real implementation, you would implement this directly in the repository
	// For now, we'll return a placeholder response
	return &models.ChunkedSyncResponse{
		Data:         make([]interface{}, 0),
		HasMore:      false,
		NextChunk:    0,
		ResumeToken:  "",
		TotalChunks:  0,
		Checkpoint:   "",
		Compressed:   false,
		DataSize:     0,
		RecordsCount: 0,
	}, nil
}

func (r *MongoDBRepository) GetSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.SyncStatus, error) {
	// Get user's last sync info
	userCollection := r.db.Collection("users")
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Count pending changes
	bookingsCollection := r.db.Collection("bookings")
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

// Setup indexes for performance
func (r *MongoDBRepository) SetupIndexes(ctx context.Context) error {
	// Users collection indexes
	usersCollection := r.db.Collection("users")

	// Email index (unique)
	emailIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	// Phone index (unique)
	phoneIndex := mongo.IndexModel{
		Keys:    bson.M{"phone": 1},
		Options: options.Index().SetUnique(true),
	}

	// Geospatial index for user location
	locationIndex := mongo.IndexModel{
		Keys: bson.M{"address": "2dsphere"},
	}

	_, err := usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, phoneIndex, locationIndex})
	if err != nil {
		r.logger.Warn("Failed to create user indexes", zap.Error(err))
	}

	// OTP collection indexes with TTL
	otpCollection := r.db.Collection("otp_records")

	// TTL index for auto-expiring OTPs
	ttlIndex := mongo.IndexModel{
		Keys:    bson.M{"expires_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	// Email + purpose compound index
	emailPurposeIndex := mongo.IndexModel{
		Keys: bson.M{"email": 1, "purpose": 1},
	}

	_, err = otpCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{ttlIndex, emailPurposeIndex})
	if err != nil {
		r.logger.Warn("Failed to create OTP indexes", zap.Error(err))
	}

	// Services collection indexes
	servicesCollection := r.db.Collection("services")

	// Category index
	categoryIndex := mongo.IndexModel{
		Keys: bson.M{"category_id": 1},
	}

	// Provider index
	providerIndex := mongo.IndexModel{
		Keys: bson.M{"provider_id": 1},
	}

	// Geospatial index for service location
	serviceLocationIndex := mongo.IndexModel{
		Keys: bson.M{"location": "2dsphere"},
	}

	// Compound index for active services by category
	activeCategoryIndex := mongo.IndexModel{
		Keys: bson.M{"is_active": 1, "category_id": 1},
	}

	_, err = servicesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		categoryIndex, providerIndex, serviceLocationIndex, activeCategoryIndex,
	})
	if err != nil {
		r.logger.Warn("Failed to create service indexes", zap.Error(err))
	}

	// Bookings collection indexes
	bookingsCollection := r.db.Collection("bookings")

	// Customer index
	customerIndex := mongo.IndexModel{
		Keys: bson.M{"customer_id": 1},
	}

	// Provider index
	providerBookingIndex := mongo.IndexModel{
		Keys: bson.M{"provider_id": 1},
	}

	// Status index
	statusIndex := mongo.IndexModel{
		Keys: bson.M{"status": 1},
	}

	// Compound index for customer bookings by status
	customerStatusIndex := mongo.IndexModel{
		Keys: bson.M{"customer_id": 1, "status": 1},
	}

	_, err = bookingsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		customerIndex, providerBookingIndex, statusIndex, customerStatusIndex,
	})
	if err != nil {
		r.logger.Warn("Failed to create booking indexes", zap.Error(err))
	}

	r.logger.Info("MongoDB indexes setup completed")
	return nil
}

// Device session operations
func (r *MongoDBRepository) CreateDeviceSession(ctx context.Context, session *models.DeviceSession) error {
	collection := r.db.Collection("device_sessions")

	if session.ID.IsZero() {
		session.ID = primitive.NewObjectID()
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	session.LastActivity = time.Now()

	_, err := collection.InsertOne(ctx, session)
	if err != nil {
		return fmt.Errorf("failed to create device session: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) GetDeviceSession(ctx context.Context, sessionID string) (*models.DeviceSession, error) {
	collection := r.db.Collection("device_sessions")

	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	var session models.DeviceSession
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&session)
	if err != nil {
		return nil, fmt.Errorf("device session not found: %w", err)
	}
	return &session, nil
}

func (r *MongoDBRepository) GetDeviceSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.DeviceSession, error) {
	collection := r.db.Collection("device_sessions")

	var session models.DeviceSession
	err := collection.FindOne(ctx, bson.M{
		"refresh_token": refreshToken,
		"is_active":     true,
	}).Decode(&session)
	if err != nil {
		return nil, fmt.Errorf("device session not found for refresh token: %w", err)
	}
	return &session, nil
}

func (r *MongoDBRepository) GetDeviceSessionByDeviceID(ctx context.Context, deviceID string) (*models.DeviceSession, error) {
	collection := r.db.Collection("device_sessions")

	var session models.DeviceSession
	err := collection.FindOne(ctx, bson.M{
		"device_id": deviceID,
		"is_active": true,
	}).Decode(&session)
	if err != nil {
		return nil, fmt.Errorf("device session not found for device ID: %w", err)
	}
	return &session, nil
}

func (r *MongoDBRepository) GetUserDeviceSessions(ctx context.Context, userID string) ([]models.DeviceSession, error) {
	collection := r.db.Collection("device_sessions")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	cursor, err := collection.Find(ctx, bson.M{"user_id": userObjectID})
	if err != nil {
		return nil, fmt.Errorf("failed to find user device sessions: %w", err)
	}
	defer cursor.Close(ctx)

	var sessions []models.DeviceSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, fmt.Errorf("failed to decode device sessions: %w", err)
	}
	return sessions, nil
}

func (r *MongoDBRepository) UpdateDeviceSessionActivity(ctx context.Context, sessionID string) error {
	collection := r.db.Collection("device_sessions")

	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"last_activity": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("failed to update device session activity: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) RevokeDeviceSession(ctx context.Context, sessionID string) error {
	collection := r.db.Collection("device_sessions")

	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	now := time.Now()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"is_active":  false,
			"revoked_at": now,
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to revoke device session: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	collection := r.db.Collection("device_sessions")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	now := time.Now()
	_, err = collection.UpdateMany(
		ctx,
		bson.M{
			"user_id":   userObjectID,
			"is_active": true,
		},
		bson.M{"$set": bson.M{
			"is_active":  false,
			"revoked_at": now,
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) RotateRefreshToken(ctx context.Context, sessionID string, newRefreshToken string) error {
	collection := r.db.Collection("device_sessions")

	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"refresh_token": newRefreshToken,
			"last_activity": time.Now(),
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to rotate refresh token: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) error {
	collection := r.db.Collection("device_sessions")

	expiredTime := time.Now().Add(-maxAge)
	_, err := collection.UpdateMany(
		ctx,
		bson.M{
			"last_activity": bson.M{"$lt": expiredTime},
			"is_active":     true,
		},
		bson.M{"$set": bson.M{
			"is_active":  false,
			"revoked_at": time.Now(),
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	return nil
}

// Security event operations
func (r *MongoDBRepository) LogSecurityEvent(ctx context.Context, event *models.SecurityEvent) error {
	collection := r.db.Collection("security_events")

	if event.ID.IsZero() {
		event.ID = primitive.NewObjectID()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to log security event: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) GetUserSecurityEvents(ctx context.Context, userID string, limit int) ([]models.SecurityEvent, error) {
	collection := r.db.Collection("security_events")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	cursor, err := collection.Find(ctx,
		bson.M{"user_id": userObjectID},
		options.Find().SetLimit(int64(limit)).SetSort(bson.M{"timestamp": -1}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find user security events: %w", err)
	}
	defer cursor.Close(ctx)

	var events []models.SecurityEvent
	if err = cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode security events: %w", err)
	}
	return events, nil
}

func (r *MongoDBRepository) GetSecurityEventsByType(ctx context.Context, userID string, eventType models.SecurityEventType, limit int) ([]models.SecurityEvent, error) {
	collection := r.db.Collection("security_events")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	cursor, err := collection.Find(ctx,
		bson.M{
			"user_id":    userObjectID,
			"event_type": eventType,
		},
		options.Find().SetLimit(int64(limit)).SetSort(bson.M{"timestamp": -1}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find security events by type: %w", err)
	}
	defer cursor.Close(ctx)

	var events []models.SecurityEvent
	if err = cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode security events: %w", err)
	}
	return events, nil
}

// Sync status operations
func (r *MongoDBRepository) UpdateSyncStatus(ctx context.Context, status *models.SyncStatus) error {
	collection := r.db.Collection("sync_statuses")

	status.UpdatedAt = time.Now()
	_, err := collection.ReplaceOne(
		ctx,
		bson.M{"user_id": status.UserID},
		status,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to update sync status: %w", err)
	}
	return nil
}

// Sync checkpoint operations
func (r *MongoDBRepository) CreateSyncCheckpoint(ctx context.Context, checkpoint *models.SyncCheckpoint) error {
	collection := r.db.Collection("sync_checkpoints")

	if checkpoint.ID.IsZero() {
		checkpoint.ID = primitive.NewObjectID()
	}
	if checkpoint.CreatedAt.IsZero() {
		checkpoint.CreatedAt = time.Now()
	}
	checkpoint.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, checkpoint)
	if err != nil {
		return fmt.Errorf("failed to create sync checkpoint: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) GetSyncCheckpoint(ctx context.Context, userID primitive.ObjectID) (*models.SyncCheckpoint, error) {
	collection := r.db.Collection("sync_checkpoints")

	var checkpoint models.SyncCheckpoint
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&checkpoint)
	if err != nil {
		return nil, fmt.Errorf("sync checkpoint not found: %w", err)
	}
	return &checkpoint, nil
}

func (r *MongoDBRepository) UpdateSyncCheckpoint(ctx context.Context, checkpoint *models.SyncCheckpoint) error {
	collection := r.db.Collection("sync_checkpoints")

	checkpoint.UpdatedAt = time.Now()
	_, err := collection.ReplaceOne(
		ctx,
		bson.M{"user_id": checkpoint.UserID},
		checkpoint,
	)
	if err != nil {
		return fmt.Errorf("failed to update sync checkpoint: %w", err)
	}
	return nil
}

// Sync metrics operations
func (r *MongoDBRepository) CreateSyncMetrics(ctx context.Context, metrics *models.SyncMetrics) error {
	collection := r.db.Collection("sync_metrics")

	if metrics.ID.IsZero() {
		metrics.ID = primitive.NewObjectID()
	}
	if metrics.CreatedAt.IsZero() {
		metrics.CreatedAt = time.Now()
	}

	_, err := collection.InsertOne(ctx, metrics)
	if err != nil {
		return fmt.Errorf("failed to create sync metrics: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) GetRecentSyncMetrics(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncMetrics, error) {
	collection := r.db.Collection("sync_metrics")

	cursor, err := collection.Find(ctx,
		bson.M{"user_id": userID},
		options.Find().SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find sync metrics: %w", err)
	}
	defer cursor.Close(ctx)

	var metrics []models.SyncMetrics
	if err = cursor.All(ctx, &metrics); err != nil {
		return nil, fmt.Errorf("failed to decode sync metrics: %w", err)
	}
	return metrics, nil
}

// Background sync queue operations
func (r *MongoDBRepository) CreateSyncQueueItem(ctx context.Context, item *models.SyncQueueItem) error {
	collection := r.db.Collection("sync_queue")

	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to create sync queue item: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) GetSyncQueueItem(ctx context.Context, itemID primitive.ObjectID) (*models.SyncQueueItem, error) {
	collection := r.db.Collection("sync_queue")

	var item models.SyncQueueItem
	err := collection.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
	if err != nil {
		return nil, fmt.Errorf("sync queue item not found: %w", err)
	}
	return &item, nil
}

func (r *MongoDBRepository) UpdateSyncQueueItem(ctx context.Context, item *models.SyncQueueItem) error {
	collection := r.db.Collection("sync_queue")

	item.UpdatedAt = time.Now()
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": item.ID}, item)
	if err != nil {
		return fmt.Errorf("failed to update sync queue item: %w", err)
	}
	return nil
}

func (r *MongoDBRepository) GetPendingSyncQueueItems(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncQueueItem, error) {
	collection := r.db.Collection("sync_queue")

	cursor, err := collection.Find(ctx,
		bson.M{
			"user_id": userID,
			"status":  models.SyncQueuePending,
		},
		options.Find().
			SetSort(bson.M{"priority": -1, "created_at": 1}).
			SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending sync queue items: %w", err)
	}
	defer cursor.Close(ctx)

	var items []models.SyncQueueItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode pending sync queue items: %w", err)
	}
	return items, nil
}

func (r *MongoDBRepository) GetConflictQueueItems(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncQueueItem, error) {
	collection := r.db.Collection("sync_queue")

	cursor, err := collection.Find(ctx,
		bson.M{
			"user_id": userID,
			"type":    models.SyncTypeConflict,
		},
		options.Find().
			SetSort(bson.M{"priority": -1, "created_at": 1}).
			SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find conflict queue items: %w", err)
	}
	defer cursor.Close(ctx)

	var items []models.SyncQueueItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode conflict queue items: %w", err)
	}
	return items, nil
}

func (r *MongoDBRepository) CleanupCompletedQueueItems(ctx context.Context, olderThan time.Duration) (int64, error) {
	collection := r.db.Collection("sync_queue")

	cutoffTime := time.Now().Add(-olderThan)
	result, err := collection.DeleteMany(ctx, bson.M{
		"status":       models.SyncQueueCompleted,
		"completed_at": bson.M{"$lt": cutoffTime},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup completed queue items: %w", err)
	}
	return result.DeletedCount, nil
}

// Background sync status operations
func (r *MongoDBRepository) GetBackgroundSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.BackgroundSyncStatus, error) {
	collection := r.db.Collection("background_sync_status")

	var status models.BackgroundSyncStatus
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&status)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create default status
			status = models.BackgroundSyncStatus{
				ID:               primitive.NewObjectID(),
				UserID:           userID,
				IsEnabled:        true,
				LastSyncAt:       time.Now(),
				PendingItems:     0,
				FailedItems:      0,
				ConflictItems:    0,
				AutoRetryEnabled: true,
				NextScheduledRun: time.Now().Add(5 * time.Minute),
				UpdatedAt:        time.Now(),
			}

			_, insertErr := collection.InsertOne(ctx, status)
			if insertErr != nil {
				return nil, fmt.Errorf("failed to create default background sync status: %w", insertErr)
			}
			return &status, nil
		}
		return nil, fmt.Errorf("background sync status not found: %w", err)
	}
	return &status, nil
}

func (r *MongoDBRepository) UpdateBackgroundSyncStatus(ctx context.Context, status *models.BackgroundSyncStatus) error {
	collection := r.db.Collection("background_sync_status")

	status.UpdatedAt = time.Now()
	_, err := collection.ReplaceOne(
		ctx,
		bson.M{"user_id": status.UserID},
		status,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to update background sync status: %w", err)
	}
	return nil
}

// Close closes the MongoDB repository
func (r *MongoDBRepository) Close() error {
	// MongoDB connection is managed by the database package
	return nil
}
