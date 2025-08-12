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

// MongoSessionStore implements SessionStore using MongoDB
type MongoSessionStore struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewMongoSessionStore creates a new MongoDB session store
func NewMongoSessionStore(database *mongo.Database, logger *zap.Logger) *MongoSessionStore {
	collection := database.Collection("user_sessions")

	store := &MongoSessionStore{
		collection: collection,
		logger:     logger,
	}

	// Create indexes
	store.createIndexes()

	return store
}

// createIndexes creates necessary indexes for optimal performance
func (s *MongoSessionStore) createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{primitive.E{Key: "session_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{primitive.E{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{primitive.E{Key: "device_id", Value: 1}},
		},
		{
			Keys:    bson.D{primitive.E{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0), // TTL index
		},
		{
			Keys: bson.D{primitive.E{Key: "created_at", Value: 1}},
		},
		{
			Keys: bson.D{primitive.E{Key: "revoked", Value: 1}},
		},
	}

	_, err := s.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		s.logger.Error("Failed to create session indexes", zap.Error(err))
	}
}

// CreateSession creates a new session in the database
func (s *MongoSessionStore) CreateSession(ctx context.Context, session *SessionInfo) error {
	_, err := s.collection.InsertOne(ctx, session)
	if err != nil {
		s.logger.Error("Failed to create session",
			zap.String("session_id", session.SessionID),
			zap.String("user_id", session.UserID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info("Session created successfully",
		zap.String("session_id", session.SessionID),
		zap.String("user_id", session.UserID),
		zap.String("device_id", session.DeviceID),
	)

	return nil
}

// GetSession retrieves a session by ID
func (s *MongoSessionStore) GetSession(ctx context.Context, sessionID string) (*SessionInfo, error) {
	var session SessionInfo

	filter := bson.M{
		"session_id": sessionID,
		"revoked":    false,
	}

	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("session not found")
		}
		s.logger.Error("Failed to get session",
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

// UpdateSession updates an existing session
func (s *MongoSessionStore) UpdateSession(ctx context.Context, session *SessionInfo) error {
	filter := bson.M{"session_id": session.SessionID}

	update := bson.M{
		"$set": bson.M{
			"last_activity":  session.LastActivity,
			"refresh_tokens": session.RefreshTokens,
			"metadata":       session.Metadata,
		},
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error("Failed to update session",
			zap.String("session_id", session.SessionID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update session: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("session not found for update")
	}

	return nil
}

// RevokeSession marks a session as revoked
func (s *MongoSessionStore) RevokeSession(ctx context.Context, sessionID string) error {
	filter := bson.M{"session_id": sessionID}

	update := bson.M{
		"$set": bson.M{
			"revoked":       true,
			"last_activity": time.Now(),
		},
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error("Failed to revoke session",
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("session not found for revocation")
	}

	s.logger.Info("Session revoked",
		zap.String("session_id", sessionID),
	)

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (s *MongoSessionStore) RevokeAllUserSessions(ctx context.Context, userID string) error {
	filter := bson.M{
		"user_id": userID,
		"revoked": false,
	}

	update := bson.M{
		"$set": bson.M{
			"revoked":       true,
			"last_activity": time.Now(),
		},
	}

	result, err := s.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		s.logger.Error("Failed to revoke all user sessions",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}

	s.logger.Info("All user sessions revoked",
		zap.String("user_id", userID),
		zap.Int64("sessions_revoked", result.ModifiedCount),
	)

	return nil
}

// GetUserSessions returns all active sessions for a user
func (s *MongoSessionStore) GetUserSessions(ctx context.Context, userID string) ([]*SessionInfo, error) {
	filter := bson.M{
		"user_id":    userID,
		"revoked":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	// Sort by last activity (most recent first)
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "last_activity", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		s.logger.Error("Failed to get user sessions",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	defer cursor.Close(ctx)

	var sessions []*SessionInfo
	for cursor.Next(ctx) {
		var session SessionInfo
		if err := cursor.Decode(&session); err != nil {
			s.logger.Warn("Failed to decode session", zap.Error(err))
			continue
		}
		sessions = append(sessions, &session)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return sessions, nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *MongoSessionStore) CleanupExpiredSessions(ctx context.Context) error {
	filter := bson.M{
		"$or": []bson.M{
			{"expires_at": bson.M{"$lt": time.Now()}},
			{"revoked": true, "last_activity": bson.M{"$lt": time.Now().Add(-7 * 24 * time.Hour)}},
		},
	}

	result, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to cleanup expired sessions", zap.Error(err))
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	if result.DeletedCount > 0 {
		s.logger.Info("Cleaned up expired sessions",
			zap.Int64("deleted_count", result.DeletedCount),
		)
	}

	return nil
}

// GetSessionStats returns statistics about sessions
func (s *MongoSessionStore) GetSessionStats(ctx context.Context) (*SessionStats, error) {
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$revoked",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get session stats: %w", err)
	}
	defer cursor.Close(ctx)

	stats := &SessionStats{}
	for cursor.Next(ctx) {
		var result struct {
			ID    bool `bson:"_id"`
			Count int  `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}

		if result.ID {
			stats.RevokedSessions = result.Count
		} else {
			stats.ActiveSessions = result.Count
		}
	}

	// Get expired sessions count
	expiredFilter := bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
		"revoked":    false,
	}
	expiredCount, err := s.collection.CountDocuments(ctx, expiredFilter)
	if err == nil {
		stats.ExpiredSessions = int(expiredCount)
	}

	// Get sessions by device type
	devicePipeline := []bson.M{
		{
			"$match": bson.M{
				"revoked":    false,
				"expires_at": bson.M{"$gt": time.Now()},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$device_info.platform",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	deviceCursor, err := s.collection.Aggregate(ctx, devicePipeline)
	if err == nil {
		defer deviceCursor.Close(ctx)
		stats.DeviceBreakdown = make(map[string]int)

		for deviceCursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int    `bson:"count"`
			}
			if err := deviceCursor.Decode(&result); err == nil {
				stats.DeviceBreakdown[result.ID] = result.Count
			}
		}
	}

	return stats, nil
}

// SessionStats contains session statistics
type SessionStats struct {
	ActiveSessions  int            `json:"active_sessions"`
	RevokedSessions int            `json:"revoked_sessions"`
	ExpiredSessions int            `json:"expired_sessions"`
	DeviceBreakdown map[string]int `json:"device_breakdown"`
	TotalSessions   int            `json:"total_sessions"`
}

// GetTotalSessions calculates total sessions
func (stats *SessionStats) GetTotalSessions() int {
	return stats.ActiveSessions + stats.RevokedSessions + stats.ExpiredSessions
}
