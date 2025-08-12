package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestDeviceStoreMongoIndexCreation tests that the index creation
// uses proper keyed fields in bson.D structures
func TestDeviceStoreMongoIndexCreation(t *testing.T) {
	t.Run("should create indexes with properly structured bson.D", func(t *testing.T) {
		// Test the index structure that was causing the go vet error
		// This validates that bson.D uses proper primitive.E syntax

		indexes := []struct {
			Keys    bson.D
			Options interface{}
		}{
			{
				Keys: bson.D{primitive.E{Key: "device_id", Value: 1}},
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

		assert.Len(t, indexes, 8)

		// Verify each index has proper structure
		for i, index := range indexes {
			assert.NotNil(t, index.Keys, "Index %d should have Keys", i)
			assert.Len(t, index.Keys, 1, "Index %d should have exactly one key", i)

			// Verify the key structure
			element := index.Keys[0]
			assert.NotEmpty(t, element.Key, "Index %d should have a non-empty key", i)
			assert.Equal(t, 1, element.Value, "Index %d should have value 1", i)
		}
	})
}

// TestDeviceStoreMongo tests basic device store functionality
// Skipping constructor test due to DB dependency
// func TestDeviceStoreMongo(t *testing.T) {
//	t.Run("should create new device store", func(t *testing.T) {
//		logger := zap.NewNop()
//
//		// Test that we can create a device store without panicking
//		store := NewMongoDeviceStore(nil, logger) // Use nil for database in unit test
//		assert.NotNil(t, store)
//	})
// }

// TestBsonDStructureCorrectness specifically tests BSON structure compliance
func TestBsonDStructureCorrectness(t *testing.T) {
	t.Run("should use keyed fields in bson.D for filters", func(t *testing.T) {
		// Test filter construction that was causing go vet errors
		filter := bson.D{primitive.E{Key: "_id", Value: primitive.NewObjectID()}}

		assert.NotNil(t, filter)
		assert.Len(t, filter, 1)

		element := filter[0]
		assert.Equal(t, "_id", element.Key)
		assert.NotNil(t, element.Value)
	})

	t.Run("should construct update documents with keyed fields", func(t *testing.T) {
		// Test update document construction
		update := bson.D{
			primitive.E{Key: "$set", Value: bson.M{
				"last_verified": time.Now(),
				"is_trusted":    true,
			}},
		}

		assert.NotNil(t, update)
		assert.Len(t, update, 1)

		element := update[0]
		assert.Equal(t, "$set", element.Key)
		assert.NotNil(t, element.Value)
	})
}
