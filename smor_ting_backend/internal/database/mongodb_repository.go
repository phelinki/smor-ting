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

	// Initialize wallet
	user.Wallet = models.Wallet{
		Balance:     0,
		Currency:    "KES",
		LastUpdated: time.Now(),
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

// Close closes the MongoDB repository
func (r *MongoDBRepository) Close() error {
	// MongoDB connection is managed by the database package
	return nil
}
