package services

import (
	"context"
	"errors"
	"time"

	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoPaymentTokenStore persists encrypted payment tokens with TTL on expires_at
type MongoPaymentTokenStore struct {
	coll   *mongo.Collection
	logger *logger.Logger
}

func NewMongoPaymentTokenStore(db *mongo.Database, logger *logger.Logger) (*MongoPaymentTokenStore, error) {
	s := &MongoPaymentTokenStore{coll: db.Collection("payment_tokens"), logger: logger}
	// Ensure TTL index on expires_at and unique _id
	_, _ = s.coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "expires_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)},
		{Keys: bson.D{{Key: "_id", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	return s, nil
}

func (m *MongoPaymentTokenStore) Save(tokenID, userID, encryptedData string, expiresAt time.Time) error {
	if tokenID == "" || encryptedData == "" {
		return errors.New("tokenID and encryptedData required")
	}
	doc := bson.M{
		"_id":            tokenID,
		"user_id":        userID,
		"encrypted_data": encryptedData,
		"created_at":     time.Now(),
		"last_used":      time.Now(),
		"expires_at":     expiresAt,
		// Metadata fields can be updated after save
	}
	_, err := m.coll.UpdateByID(context.Background(), tokenID, bson.M{"$set": doc}, options.Update().SetUpsert(true))
	return err
}

func (m *MongoPaymentTokenStore) Get(tokenID string) (*PaymentTokenRecord, error) {
	var out struct {
		ID            string    `bson:"_id"`
		UserID        string    `bson:"user_id"`
		EncryptedData string    `bson:"encrypted_data"`
		TokenType     string    `bson:"token_type"`
		LastFour      string    `bson:"last_four"`
		Brand         string    `bson:"brand"`
		CreatedAt     time.Time `bson:"created_at"`
		LastUsed      time.Time `bson:"last_used"`
		ExpiresAt     time.Time `bson:"expires_at"`
	}
	err := m.coll.FindOne(context.Background(), bson.M{"_id": tokenID}).Decode(&out)
	if err != nil {
		return nil, err
	}
	return &PaymentTokenRecord{
		TokenID:       out.ID,
		UserID:        out.UserID,
		EncryptedData: out.EncryptedData,
		TokenType:     out.TokenType,
		LastFour:      out.LastFour,
		Brand:         out.Brand,
		CreatedAt:     out.CreatedAt,
		LastUsed:      out.LastUsed,
		ExpiresAt:     out.ExpiresAt,
	}, nil
}

func (m *MongoPaymentTokenStore) Delete(tokenID string) error {
	_, err := m.coll.DeleteOne(context.Background(), bson.M{"_id": tokenID})
	return err
}

func (m *MongoPaymentTokenStore) TouchLastUsed(tokenID string, t time.Time) error {
	_, err := m.coll.UpdateByID(context.Background(), tokenID, bson.M{"$set": bson.M{"last_used": t}})
	return err
}

func (m *MongoPaymentTokenStore) PurgeExpired(ctx context.Context) error {
	// TTL index handles purge; optional manual purge if index is not supported
	_, _ = m.coll.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": time.Now()}})
	return nil
}
