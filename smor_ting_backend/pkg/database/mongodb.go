package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MongoDB represents a MongoDB database connection
type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
	Config *configs.DatabaseConfig
	Logger *logger.Logger
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(config *configs.DatabaseConfig, log *logger.Logger) (*MongoDB, error) {
	mongoDB := &MongoDB{
		Config: config,
		Logger: log,
	}

	if err := mongoDB.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := mongoDB.ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	if err := mongoDB.setupCollections(); err != nil {
		return nil, fmt.Errorf("failed to setup collections: %w", err)
	}

	return mongoDB, nil
}

// connect establishes a connection to MongoDB
func (m *MongoDB) connect() error {
	var uri string

	// Check if we have a MongoDB Atlas connection string
	if m.Config.ConnectionString != "" {
		uri = m.Config.ConnectionString
		m.Logger.Info("Using MongoDB Atlas connection string")
	} else if m.Config.Driver == "mongodb" {
		// MongoDB connection string format
		if m.Config.Username != "" && m.Config.Password != "" {
			uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
				m.Config.Username, m.Config.Password, m.Config.Host, m.Config.Port, m.Config.Database)
		} else {
			uri = fmt.Sprintf("mongodb://%s:%s/%s",
				m.Config.Host, m.Config.Port, m.Config.Database)
		}

		// Add SSL mode if specified
		if m.Config.SSLMode == "require" {
			uri += "?ssl=true"
		}
	} else {
		return fmt.Errorf("unsupported database driver: %s", m.Config.Driver)
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// For MongoDB Atlas, we need to set additional options
	if m.Config.AtlasCluster || strings.Contains(uri, "mongodb+srv://") {
		clientOptions.SetRetryWrites(true)
		clientOptions.SetRetryReads(true)
		clientOptions.SetMaxPoolSize(10)
		clientOptions.SetMinPoolSize(1)
		clientOptions.SetMaxConnIdleTime(30 * time.Second)
	}

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	m.Client = client
	m.DB = client.Database(m.Config.Database)

	m.Logger.Info("Connected to MongoDB",
		zap.String("host", m.Config.Host),
		zap.String("database", m.Config.Database),
		zap.Bool("atlas_cluster", m.Config.AtlasCluster))

	return nil
}

// ping verifies the MongoDB connection
func (m *MongoDB) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.Client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping failed: %w", err)
	}

	m.Logger.Info("MongoDB connection established successfully")
	return nil
}

// setupCollections creates the necessary collections and indexes
func (m *MongoDB) setupCollections() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create collections if they don't exist
	collections := []string{"users", "services", "bookings", "payments", "agents"}

	for _, collectionName := range collections {
		// Create collection (MongoDB will create it if it doesn't exist)
		err := m.DB.CreateCollection(ctx, collectionName)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			m.Logger.Warn("Collection might already exist", zap.String("collection", collectionName))
		}
	}

	// Create indexes for better performance
	if err := m.createIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	m.Logger.Info("MongoDB collections and indexes setup completed")
	return nil
}

// createIndexes creates necessary indexes for performance
func (m *MongoDB) createIndexes(ctx context.Context) error {
	// Users collection indexes
	usersCollection := m.DB.Collection("users")

	// Email index (unique)
	emailIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	// Phone index (unique)
	phoneIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"phone": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, phoneIndex})
	if err != nil {
		m.Logger.Warn("Failed to create user indexes", zap.Error(err))
	}

	// Services collection indexes
	servicesCollection := m.DB.Collection("services")

	// Category index
	categoryIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"category": 1,
		},
	}

	// Location index for geospatial queries
	locationIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"location": "2dsphere",
		},
	}

	_, err = servicesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{categoryIndex, locationIndex})
	if err != nil {
		m.Logger.Warn("Failed to create service indexes", zap.Error(err))
	}

	// Bookings collection indexes
	bookingsCollection := m.DB.Collection("bookings")

	// User ID index
	userIDIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"user_id": 1,
		},
	}

	// Service provider ID index
	providerIDIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"provider_id": 1,
		},
	}

	// Status index
	statusIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"status": 1,
		},
	}

	_, err = bookingsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{userIDIndex, providerIDIndex, statusIndex})
	if err != nil {
		m.Logger.Warn("Failed to create booking indexes", zap.Error(err))
	}

	return nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	if m.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := m.Client.Disconnect(ctx); err != nil {
			m.Logger.Error("Failed to close MongoDB connection", err)
			return err
		}
		m.Logger.Info("MongoDB connection closed")
	}
	return nil
}

// GetDB returns the MongoDB database instance
func (m *MongoDB) GetDB() *mongo.Database {
	return m.DB
}

// GetClient returns the MongoDB client instance
func (m *MongoDB) GetClient() *mongo.Client {
	return m.Client
}

// HealthCheck performs a health check on the MongoDB connection
func (m *MongoDB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Ping(ctx, nil)
}

// IsInMemory returns false for MongoDB (always persistent)
func (m *MongoDB) IsInMemory() bool {
	return false
}

// Device session operations (placeholder implementations for now)
func (m *MongoDB) CreateDeviceSession(ctx context.Context, session *models.DeviceSession) error {
	// TODO: Implement MongoDB device session creation
	return nil
}

func (m *MongoDB) GetDeviceSession(ctx context.Context, sessionID string) (*models.DeviceSession, error) {
	// TODO: Implement MongoDB device session retrieval
	return nil, nil
}

func (m *MongoDB) GetDeviceSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.DeviceSession, error) {
	// TODO: Implement MongoDB device session retrieval by refresh token
	return nil, nil
}

func (m *MongoDB) GetDeviceSessionByDeviceID(ctx context.Context, deviceID string) (*models.DeviceSession, error) {
	// TODO: Implement MongoDB device session retrieval by device ID
	return nil, nil
}

func (m *MongoDB) GetUserDeviceSessions(ctx context.Context, userID string) ([]models.DeviceSession, error) {
	// TODO: Implement MongoDB user device sessions retrieval
	return nil, nil
}

func (m *MongoDB) UpdateDeviceSessionActivity(ctx context.Context, sessionID string) error {
	// TODO: Implement MongoDB device session activity update
	return nil
}

func (m *MongoDB) RevokeDeviceSession(ctx context.Context, sessionID string) error {
	// TODO: Implement MongoDB device session revocation
	return nil
}

func (m *MongoDB) RevokeAllUserTokens(ctx context.Context, userID string) error {
	// TODO: Implement MongoDB user token revocation
	return nil
}

func (m *MongoDB) RotateRefreshToken(ctx context.Context, sessionID string, newRefreshToken string) error {
	// TODO: Implement MongoDB refresh token rotation
	return nil
}

func (m *MongoDB) CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) error {
	// TODO: Implement MongoDB expired session cleanup
	return nil
}

// Security event operations (placeholder implementations for now)
func (m *MongoDB) LogSecurityEvent(ctx context.Context, event *models.SecurityEvent) error {
	// TODO: Implement MongoDB security event logging
	return nil
}

func (m *MongoDB) GetUserSecurityEvents(ctx context.Context, userID string, limit int) ([]models.SecurityEvent, error) {
	// TODO: Implement MongoDB user security events retrieval
	return nil, nil
}

func (m *MongoDB) GetSecurityEventsByType(ctx context.Context, userID string, eventType models.SecurityEventType, limit int) ([]models.SecurityEvent, error) {
	// TODO: Implement MongoDB security events retrieval by type
	return nil, nil
}

// Sync operations (placeholder implementations for now)
func (m *MongoDB) UpdateSyncStatus(ctx context.Context, status *models.SyncStatus) error {
	// TODO: Implement MongoDB sync status update
	return nil
}

func (m *MongoDB) CreateSyncCheckpoint(ctx context.Context, checkpoint *models.SyncCheckpoint) error {
	// TODO: Implement MongoDB sync checkpoint creation
	return nil
}

func (m *MongoDB) GetSyncCheckpoint(ctx context.Context, userID primitive.ObjectID) (*models.SyncCheckpoint, error) {
	// TODO: Implement MongoDB sync checkpoint retrieval
	return nil, nil
}

func (m *MongoDB) UpdateSyncCheckpoint(ctx context.Context, checkpoint *models.SyncCheckpoint) error {
	// TODO: Implement MongoDB sync checkpoint update
	return nil
}

func (m *MongoDB) CreateSyncMetrics(ctx context.Context, metrics *models.SyncMetrics) error {
	// TODO: Implement MongoDB sync metrics creation
	return nil
}

func (m *MongoDB) GetRecentSyncMetrics(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncMetrics, error) {
	// TODO: Implement MongoDB sync metrics retrieval
	return nil, nil
}

// Background sync queue operations (placeholder implementations for now)
func (m *MongoDB) CreateSyncQueueItem(ctx context.Context, item *models.SyncQueueItem) error {
	// TODO: Implement MongoDB sync queue item creation
	return nil
}

func (m *MongoDB) GetSyncQueueItem(ctx context.Context, itemID primitive.ObjectID) (*models.SyncQueueItem, error) {
	// TODO: Implement MongoDB sync queue item retrieval
	return nil, nil
}

func (m *MongoDB) UpdateSyncQueueItem(ctx context.Context, item *models.SyncQueueItem) error {
	// TODO: Implement MongoDB sync queue item update
	return nil
}

func (m *MongoDB) GetPendingSyncQueueItems(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncQueueItem, error) {
	// TODO: Implement MongoDB pending sync queue items retrieval
	return nil, nil
}

func (m *MongoDB) GetConflictQueueItems(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncQueueItem, error) {
	// TODO: Implement MongoDB conflict queue items retrieval
	return nil, nil
}

func (m *MongoDB) CleanupCompletedQueueItems(ctx context.Context, olderThan time.Duration) (int64, error) {
	// TODO: Implement MongoDB completed queue items cleanup
	return 0, nil
}

func (m *MongoDB) GetBackgroundSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.BackgroundSyncStatus, error) {
	// TODO: Implement MongoDB background sync status retrieval
	return nil, nil
}

func (m *MongoDB) UpdateBackgroundSyncStatus(ctx context.Context, status *models.BackgroundSyncStatus) error {
	// TODO: Implement MongoDB background sync status update
	return nil
}
