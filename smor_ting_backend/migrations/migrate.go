package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Migration represents a database migration
type Migration struct {
	Version     int       `bson:"version"`
	Description string    `bson:"description"`
	AppliedAt   time.Time `bson:"applied_at"`
	Script      string    `bson:"script"`
}

// Migrator handles database migrations
type Migrator struct {
	db     *mongo.Database
	logger *logger.Logger
}

// NewMigrator creates a new migrator
func NewMigrator(db *mongo.Database, logger *logger.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// RunMigrations runs all pending migrations
func (m *Migrator) RunMigrations(ctx context.Context) error {
	m.logger.Info("Starting database migrations")

	// Create migrations collection if it doesn't exist
	collection := m.db.Collection("migrations")

	// Get current version
	currentVersion, err := m.getCurrentVersion(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	m.logger.Info("Current database version", zap.Int("version", currentVersion))

	// Define migrations
	migrations := []Migration{
		{
			Version:     1,
			Description: "Create initial collections and indexes",
			Script:      "create_initial_collections",
		},
		{
			Version:     2,
			Description: "Add offline-first fields to all collections",
			Script:      "add_offline_fields",
		},
		{
			Version:     3,
			Description: "Add wallet and transaction support",
			Script:      "add_wallet_support",
		},
		{
			Version:     4,
			Description: "Add geospatial indexes for location-based queries",
			Script:      "add_geospatial_indexes",
		},
		{
			Version:     5,
			Description: "Add sync collections and indexes for checkpoint and metrics",
			Script:      "add_sync_collections",
		},
		{
			Version:     6,
			Description: "Fix geospatial indexes to avoid indexing non-GeoJSON address",
			Script:      "fix_geospatial_indexes",
		},
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			if err := m.applyMigration(ctx, collection, migration); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
			}
			currentVersion = migration.Version
		}
	}

	m.logger.Info("Database migrations completed", zap.Int("final_version", currentVersion))
	return nil
}

// getCurrentVersion gets the current database version
func (m *Migrator) getCurrentVersion(ctx context.Context, collection *mongo.Collection) (int, error) {
	// Find the highest version
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"version": -1}).SetLimit(1))
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var migration Migration
		if err := cursor.Decode(&migration); err != nil {
			return 0, err
		}
		return migration.Version, nil
	}

	return 0, nil
}

// applyMigration applies a single migration
func (m *Migrator) applyMigration(ctx context.Context, collection *mongo.Collection, migration Migration) error {
	m.logger.Info("Applying migration",
		zap.Int("version", migration.Version),
		zap.String("description", migration.Description))

	// Execute migration without transactions to support standalone MongoDB
	// (Transactions require a replica set; dev often runs standalone.)
	switch migration.Script {
	case "create_initial_collections":
		if err := m.createInitialCollections(ctx); err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	case "add_offline_fields":
		if err := m.addOfflineFields(ctx); err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	case "add_wallet_support":
		if err := m.addWalletSupport(ctx); err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	case "add_geospatial_indexes":
		if err := m.addGeospatialIndexes(ctx); err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	case "add_sync_collections":
		if err := m.addSyncCollections(ctx); err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	case "fix_geospatial_indexes":
		if err := m.fixGeospatialIndexes(ctx); err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	default:
		return fmt.Errorf("unknown migration script: %s", migration.Script)
	}

	// Record migration
	migration.AppliedAt = time.Now()
	if _, insertErr := collection.InsertOne(ctx, migration); insertErr != nil {
		return fmt.Errorf("failed to record migration: %w", insertErr)
	}

	m.logger.Info("Migration applied successfully", zap.Int("version", migration.Version))
	return nil
}

// createInitialCollections creates the initial collections
func (m *Migrator) createInitialCollections(ctx context.Context) error {
	collections := []string{"users", "services", "bookings", "payments", "agents", "otp_records"}

	for _, collectionName := range collections {
		err := m.db.CreateCollection(ctx, collectionName)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("failed to create collection %s: %w", collectionName, err)
		}
	}

	return nil
}

// addOfflineFields adds offline-first fields to existing documents
func (m *Migrator) addOfflineFields(ctx context.Context) error {
	// Update users collection
	usersCollection := m.db.Collection("users")
	_, err := usersCollection.UpdateMany(
		ctx,
		bson.M{"last_sync_at": bson.M{"$exists": false}},
		bson.M{
			"$set": bson.M{
				"last_sync_at": time.Now(),
				"is_offline":   false,
				"version":      1,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update users: %w", err)
	}

	// Update services collection
	servicesCollection := m.db.Collection("services")
	_, err = servicesCollection.UpdateMany(
		ctx,
		bson.M{"last_sync_at": bson.M{"$exists": false}},
		bson.M{
			"$set": bson.M{
				"last_sync_at": time.Now(),
				"version":      1,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update services: %w", err)
	}

	// Update bookings collection
	bookingsCollection := m.db.Collection("bookings")
	_, err = bookingsCollection.UpdateMany(
		ctx,
		bson.M{"last_sync_at": bson.M{"$exists": false}},
		bson.M{
			"$set": bson.M{
				"last_sync_at": time.Now(),
				"version":      1,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update bookings: %w", err)
	}

	return nil
}

// addWalletSupport adds wallet support to users
func (m *Migrator) addWalletSupport(ctx context.Context) error {
	usersCollection := m.db.Collection("users")
	_, err := usersCollection.UpdateMany(
		ctx,
		bson.M{"wallet": bson.M{"$exists": false}},
		bson.M{
			"$set": bson.M{
				"wallet": bson.M{
					"balance":      0,
					"currency":     "LRD",
					"transactions": []interface{}{},
					"last_updated": time.Now(),
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add wallet support: %w", err)
	}

	return nil
}

// addGeospatialIndexes adds geospatial indexes for location-based queries
func (m *Migrator) addGeospatialIndexes(ctx context.Context) error {
	// Users collection geospatial index
	usersCollection := m.db.Collection("users")
    locationIndex := mongo.IndexModel{
        Keys: bson.M{"address": "2dsphere"},
        Options: options.Index().SetPartialFilterExpression(bson.M{"address.type": "Point"}),
    }
	_, err := usersCollection.Indexes().CreateOne(ctx, locationIndex)
	if err != nil {
		m.logger.Warn("Failed to create user location index", zap.Error(err))
	}

	// Services collection geospatial index
	servicesCollection := m.db.Collection("services")
	serviceLocationIndex := mongo.IndexModel{
		Keys: bson.M{"location": "2dsphere"},
	}
	_, err = servicesCollection.Indexes().CreateOne(ctx, serviceLocationIndex)
	if err != nil {
		m.logger.Warn("Failed to create service location index", zap.Error(err))
	}

	return nil
}

// fixGeospatialIndexes drops old address 2dsphere index and recreates with partial filter
func (m *Migrator) fixGeospatialIndexes(ctx context.Context) error {
    usersCollection := m.db.Collection("users")

    // List indexes and drop any 2dsphere on address without partial filter
    cur, err := usersCollection.Indexes().List(ctx)
    if err == nil {
        defer cur.Close(ctx)
        for cur.Next(ctx) {
            var idx bson.M
            if err := cur.Decode(&idx); err == nil {
                // keys: { address: "2dsphere" }
                if keys, ok := idx["key"].(bson.M); ok {
                    if v, ok2 := keys["address"]; ok2 && v == "2dsphere" {
                        if name, ok3 := idx["name"].(string); ok3 {
                            // Drop index
                            if _, dropErr := usersCollection.Indexes().DropOne(ctx, name); dropErr != nil {
                                m.logger.Warn("Failed to drop existing address index", zap.Error(dropErr))
                            }
                        }
                    }
                }
            }
        }
    }

    // Recreate with partial filter
    locationIndex := mongo.IndexModel{
        Keys: bson.M{"address": "2dsphere"},
        Options: options.Index().SetPartialFilterExpression(bson.M{"address.type": "Point"}),
    }
    if _, err := usersCollection.Indexes().CreateOne(ctx, locationIndex); err != nil {
        m.logger.Warn("Failed to recreate user location index with partial filter", zap.Error(err))
    }

    return nil
}

// addSyncCollections adds sync-related collections and indexes
func (m *Migrator) addSyncCollections(ctx context.Context) error {
	// Create sync_checkpoints collection with indexes
	syncCheckpointsCollection := m.db.Collection("sync_checkpoints")

	// Index on user_id and checkpoint (ordered)
	userCheckpointIndex := mongo.IndexModel{
		Keys: bson.D{bson.E{Key: "user_id", Value: 1}, bson.E{Key: "checkpoint", Value: 1}},
	}
	if _, err := syncCheckpointsCollection.Indexes().CreateOne(ctx, userCheckpointIndex); err != nil {
		return fmt.Errorf("failed to create sync_checkpoints index: %w", err)
	}

	// Index on user_id and created_at for cleanup (ordered)
	userCreatedIndex := mongo.IndexModel{
		Keys: bson.D{bson.E{Key: "user_id", Value: 1}, bson.E{Key: "created_at", Value: -1}},
	}
	if _, err := syncCheckpointsCollection.Indexes().CreateOne(ctx, userCreatedIndex); err != nil {
		return fmt.Errorf("failed to create sync_checkpoints created_at index: %w", err)
	}

	// Create sync_metrics collection with indexes
	syncMetricsCollection := m.db.Collection("sync_metrics")

	// Index on user_id and created_at (ordered)
	metricsUserIndex := mongo.IndexModel{
		Keys: bson.D{bson.E{Key: "user_id", Value: 1}, bson.E{Key: "created_at", Value: -1}},
	}
	if _, err := syncMetricsCollection.Indexes().CreateOne(ctx, metricsUserIndex); err != nil {
		return fmt.Errorf("failed to create sync_metrics index: %w", err)
	}

	// TTL index to automatically clean up old metrics (keep for 30 days)
	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{bson.E{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60), // 30 days
	}
	if _, err := syncMetricsCollection.Indexes().CreateOne(ctx, ttlIndex); err != nil {
		return fmt.Errorf("failed to create sync_metrics TTL index: %w", err)
	}

	return nil
}
