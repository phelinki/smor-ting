package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MongoDeviceStore implements DeviceStore using MongoDB
type MongoDeviceStore struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// DeviceDocument represents the MongoDB document structure for devices
type DeviceDocument struct {
	DeviceFingerprint `bson:",inline"`
	UserID            string    `bson:"user_id"`
	CreatedAt         time.Time `bson:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at"`
}

// NewMongoDeviceStore creates a new MongoDB device store
func NewMongoDeviceStore(database *mongo.Database, logger *zap.Logger) *MongoDeviceStore {
	collection := database.Collection("trusted_devices")

	store := &MongoDeviceStore{
		collection: collection,
		logger:     logger,
	}

	// Create indexes
	store.createIndexes()

	return store
}

// createIndexes creates necessary indexes for optimal performance
func (s *MongoDeviceStore) createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{primitive.E{Key: "device_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{primitive.E{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{primitive.E{Key: "is_trusted", Value: 1}},
		},
		{
			Keys: bson.D{primitive.E{Key: "platform", Value: 1}},
		},
		{
			Keys: bson.D{primitive.E{Key: "is_jailbroken", Value: 1}},
		},
		{
			Keys: bson.D{primitive.E{Key: "trust_score", Value: 1}},
		},
		{
			Keys: bson.D{primitive.E{Key: "last_verified", Value: 1}},
		},
		{
			Keys: bson.D{primitive.E{Key: "created_at", Value: 1}},
		},
	}

	_, err := s.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		s.logger.Error("Failed to create device indexes", zap.Error(err))
	}
}

// RegisterDevice registers a new device
func (s *MongoDeviceStore) RegisterDevice(ctx context.Context, device *DeviceFingerprint) error {
	doc := &DeviceDocument{
		DeviceFingerprint: *device,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Use upsert to handle duplicate device IDs
	filter := bson.M{"device_id": device.DeviceID}
	update := bson.M{
		"$set": doc,
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		s.logger.Error("Failed to register device",
			zap.String("device_id", device.DeviceID),
			zap.String("platform", device.Platform),
			zap.Error(err),
		)
		return fmt.Errorf("failed to register device: %w", err)
	}

	s.logger.Info("Device registered successfully",
		zap.String("device_id", device.DeviceID),
		zap.String("platform", device.Platform),
		zap.Bool("trusted", device.IsTrusted),
		zap.Float64("trust_score", device.TrustScore),
	)

	return nil
}

// GetDevice retrieves a device by ID
func (s *MongoDeviceStore) GetDevice(ctx context.Context, deviceID string) (*DeviceFingerprint, error) {
	var doc DeviceDocument

	filter := bson.M{"device_id": deviceID}

	err := s.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("device not found")
		}
		s.logger.Error("Failed to get device",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return &doc.DeviceFingerprint, nil
}

// UpdateDeviceTrust updates the trust status of a device
func (s *MongoDeviceStore) UpdateDeviceTrust(ctx context.Context, deviceID string, trusted bool, score float64) error {
	filter := bson.M{"device_id": deviceID}

	update := bson.M{
		"$set": bson.M{
			"is_trusted":    trusted,
			"trust_score":   score,
			"last_verified": time.Now(),
			"updated_at":    time.Now(),
		},
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error("Failed to update device trust",
			zap.String("device_id", deviceID),
			zap.Bool("trusted", trusted),
			zap.Float64("score", score),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update device trust: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("device not found for trust update")
	}

	s.logger.Info("Device trust updated",
		zap.String("device_id", deviceID),
		zap.Bool("trusted", trusted),
		zap.Float64("score", score),
	)

	return nil
}

// GetUserDevices returns all devices for a user
func (s *MongoDeviceStore) GetUserDevices(ctx context.Context, userID string) ([]*DeviceFingerprint, error) {
	filter := bson.M{"user_id": userID}

	// Sort by last verified (most recent first)
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "last_verified", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		s.logger.Error("Failed to get user devices",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user devices: %w", err)
	}
	defer cursor.Close(ctx)

	var devices []*DeviceFingerprint
	for cursor.Next(ctx) {
		var doc DeviceDocument
		if err := cursor.Decode(&doc); err != nil {
			s.logger.Warn("Failed to decode device", zap.Error(err))
			continue
		}
		devices = append(devices, &doc.DeviceFingerprint)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return devices, nil
}

// RevokeDevice marks a device as untrusted and removes it
func (s *MongoDeviceStore) RevokeDevice(ctx context.Context, deviceID string) error {
	filter := bson.M{"device_id": deviceID}

	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to revoke device",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to revoke device: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("device not found for revocation")
	}

	s.logger.Info("Device revoked",
		zap.String("device_id", deviceID),
	)

	return nil
}

// GetDeviceStats returns statistics about devices
func (s *MongoDeviceStore) GetDeviceStats(ctx context.Context) (*DeviceStats, error) {
	// Get trust status breakdown
	trustPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$is_trusted",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := s.collection.Aggregate(ctx, trustPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get device trust stats: %w", err)
	}
	defer cursor.Close(ctx)

	stats := &DeviceStats{}
	for cursor.Next(ctx) {
		var result struct {
			ID    bool `bson:"_id"`
			Count int  `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}

		if result.ID {
			stats.TrustedDevices = result.Count
		} else {
			stats.UntrustedDevices = result.Count
		}
	}

	// Get platform breakdown
	platformPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$platform",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	platformCursor, err := s.collection.Aggregate(ctx, platformPipeline)
	if err == nil {
		defer platformCursor.Close(ctx)
		stats.PlatformBreakdown = make(map[string]int)

		for platformCursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int    `bson:"count"`
			}
			if err := platformCursor.Decode(&result); err == nil {
				stats.PlatformBreakdown[result.ID] = result.Count
			}
		}
	}

	// Get jailbroken devices count
	jailbrokenFilter := bson.M{"is_jailbroken": true}
	jailbrokenCount, err := s.collection.CountDocuments(ctx, jailbrokenFilter)
	if err == nil {
		stats.JailbrokenDevices = int(jailbrokenCount)
	}

	// Get average trust score
	trustScorePipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":             nil,
				"avg_trust_score": bson.M{"$avg": "$trust_score"},
				"min_trust_score": bson.M{"$min": "$trust_score"},
				"max_trust_score": bson.M{"$max": "$trust_score"},
			},
		},
	}

	scoreCursor, err := s.collection.Aggregate(ctx, trustScorePipeline)
	if err == nil {
		defer scoreCursor.Close(ctx)

		if scoreCursor.Next(ctx) {
			var result struct {
				AvgTrustScore float64 `bson:"avg_trust_score"`
				MinTrustScore float64 `bson:"min_trust_score"`
				MaxTrustScore float64 `bson:"max_trust_score"`
			}
			if err := scoreCursor.Decode(&result); err == nil {
				stats.AvgTrustScore = result.AvgTrustScore
				stats.MinTrustScore = result.MinTrustScore
				stats.MaxTrustScore = result.MaxTrustScore
			}
		}
	}

	return stats, nil
}

// DeviceStats contains device statistics
type DeviceStats struct {
	TrustedDevices    int            `json:"trusted_devices"`
	UntrustedDevices  int            `json:"untrusted_devices"`
	JailbrokenDevices int            `json:"jailbroken_devices"`
	PlatformBreakdown map[string]int `json:"platform_breakdown"`
	AvgTrustScore     float64        `json:"avg_trust_score"`
	MinTrustScore     float64        `json:"min_trust_score"`
	MaxTrustScore     float64        `json:"max_trust_score"`
	TotalDevices      int            `json:"total_devices"`
}

// GetTotalDevices calculates total devices
func (stats *DeviceStats) GetTotalDevices() int {
	return stats.TrustedDevices + stats.UntrustedDevices
}

// CleanupOldDevices removes devices that haven't been used in a long time
func (s *MongoDeviceStore) CleanupOldDevices(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)

	filter := bson.M{
		"last_verified": bson.M{"$lt": cutoff},
		"is_trusted":    false, // Only remove untrusted old devices
	}

	result, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to cleanup old devices", zap.Error(err))
		return fmt.Errorf("failed to cleanup old devices: %w", err)
	}

	if result.DeletedCount > 0 {
		s.logger.Info("Cleaned up old devices",
			zap.Int64("deleted_count", result.DeletedCount),
			zap.Duration("older_than", olderThan),
		)
	}

	return nil
}

// MarkDeviceCompromised marks a device as potentially compromised
func (s *MongoDeviceStore) MarkDeviceCompromised(ctx context.Context, deviceID string, reason string) error {
	filter := bson.M{"device_id": deviceID}

	update := bson.M{
		"$set": bson.M{
			"is_trusted":    false,
			"trust_score":   0.0,
			"is_jailbroken": true, // Mark as jailbroken as a safety measure
			"last_verified": time.Now(),
			"updated_at":    time.Now(),
		},
		"$push": bson.M{
			"security_events": bson.M{
				"event":     "compromised",
				"reason":    reason,
				"timestamp": time.Now(),
			},
		},
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error("Failed to mark device as compromised",
			zap.String("device_id", deviceID),
			zap.String("reason", reason),
			zap.Error(err),
		)
		return fmt.Errorf("failed to mark device as compromised: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("device not found for compromise marking")
	}

	s.logger.Warn("Device marked as compromised",
		zap.String("device_id", deviceID),
		zap.String("reason", reason),
	)

	return nil
}
