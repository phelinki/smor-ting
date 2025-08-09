package services

import (
	"context"
	"time"

	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRevocationStore stores revoked token IDs in MongoDB with TTL
type MongoRevocationStore struct {
	coll   *mongo.Collection
	logger *logger.Logger
}

type revokedTokenDoc struct {
	TokenID   string    `bson:"_id"`
	ExpiresAt time.Time `bson:"expires_at"`
	RevokedAt time.Time `bson:"revoked_at"`
}

func NewMongoRevocationStore(db *mongo.Database, logger *logger.Logger) (*MongoRevocationStore, error) {
	store := &MongoRevocationStore{coll: db.Collection("revoked_tokens"), logger: logger}
	// Ensure TTL index on expires_at
	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	if _, err := store.coll.Indexes().CreateOne(context.Background(), idx); err != nil {
		// log only; not fatal to create store
		logger.Error("Failed to create TTL index for revoked_tokens", err)
	}
	return store, nil
}

func (m *MongoRevocationStore) Revoke(tokenID string, expiresAt time.Time) error {
	doc := revokedTokenDoc{
		TokenID:   tokenID,
		ExpiresAt: expiresAt,
		RevokedAt: time.Now(),
	}
	_, err := m.coll.UpdateByID(context.Background(), tokenID, bson.M{"$set": doc}, options.Update().SetUpsert(true))
	return err
}

func (m *MongoRevocationStore) IsRevoked(tokenID string) (bool, error) {
	err := m.coll.FindOne(context.Background(), bson.M{"_id": tokenID}).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return err == nil, err
}

func (m *MongoRevocationStore) PurgeExpired(ctx context.Context) error {
	// TTL handles purge; do nothing
	return nil
}
