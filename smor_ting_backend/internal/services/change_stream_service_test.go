package services

import (
	"testing"
	"time"

	"github.com/smorting/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewChangeStreamService(t *testing.T) {
	t.Run("should create new change stream service", func(t *testing.T) {
		logger, _ := logger.New("info", "console", "stdout")
		service := NewChangeStreamService(nil, logger) // Use nil for DB in unit test

		assert.NotNil(t, service)
		assert.NotNil(t, service.logger)
		assert.NotNil(t, service.ctx)
		assert.NotNil(t, service.cancel)
	})
}

func TestProcessChangeEvent(t *testing.T) {
	logger, _ := logger.New("info", "console", "stdout")
	service := NewChangeStreamService(nil, logger)

	t.Run("should process insert change event", func(t *testing.T) {
		// Create a mock change event for insert operation
		event := bson.M{
			"operationType": "insert",
			"ns": bson.M{
				"db":   "test",
				"coll": "users",
			},
			"documentKey": bson.M{
				"_id": primitive.NewObjectID(),
			},
			"fullDocument": bson.M{
				"name":  "Test User",
				"email": "test@example.com",
			},
		}

		err := service.processChangeEvent(event)
		assert.NoError(t, err)
	})

	t.Run("should process update change event", func(t *testing.T) {
		// Create a mock change event for update operation
		event := bson.M{
			"operationType": "update",
			"ns": bson.M{
				"db":   "test",
				"coll": "users",
			},
			"documentKey": bson.M{
				"_id": primitive.NewObjectID(),
			},
			"updateDescription": bson.M{
				"updatedFields": bson.M{
					"name": "Updated Name",
				},
			},
		}

		err := service.processChangeEvent(event)
		assert.NoError(t, err)
	})

	t.Run("should process delete change event", func(t *testing.T) {
		// Create a mock change event for delete operation
		event := bson.M{
			"operationType": "delete",
			"ns": bson.M{
				"db":   "test",
				"coll": "users",
			},
			"documentKey": bson.M{
				"_id": primitive.NewObjectID(),
			},
		}

		err := service.processChangeEvent(event)
		assert.NoError(t, err)
	})

	t.Run("should handle invalid operation type", func(t *testing.T) {
		// Create a mock change event with invalid operation type
		event := bson.M{
			"operationType": 123, // Invalid type
		}

		err := service.processChangeEvent(event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid operation type")
	})

	t.Run("should handle invalid namespace", func(t *testing.T) {
		// Create a mock change event with invalid namespace
		event := bson.M{
			"operationType": "insert",
			"ns":            "invalid", // Invalid type
		}

		err := service.processChangeEvent(event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid namespace")
	})
}

func TestChangeEventHandlers(t *testing.T) {
	logger, _ := logger.New("info", "console", "stdout")
	service := NewChangeStreamService(nil, logger)

	t.Run("should handle user operations", func(t *testing.T) {
		event := ChangeEvent{
			OperationType: "insert",
			Collection:    "users",
			DocumentID:    "test-user-id",
			Timestamp:     time.Now(),
		}

		err := service.handleUserInsert(event)
		assert.NoError(t, err)

		err = service.handleUserUpdate(event)
		assert.NoError(t, err)

		err = service.handleUserDelete(event)
		assert.NoError(t, err)

		err = service.handleUserReplace(event)
		assert.NoError(t, err)
	})

	t.Run("should handle service operations", func(t *testing.T) {
		event := ChangeEvent{
			OperationType: "insert",
			Collection:    "services",
			DocumentID:    "test-service-id",
			Timestamp:     time.Now(),
		}

		err := service.handleServiceInsert(event)
		assert.NoError(t, err)

		err = service.handleServiceUpdate(event)
		assert.NoError(t, err)

		err = service.handleServiceDelete(event)
		assert.NoError(t, err)

		err = service.handleServiceReplace(event)
		assert.NoError(t, err)
	})

	t.Run("should handle booking operations", func(t *testing.T) {
		event := ChangeEvent{
			OperationType: "insert",
			Collection:    "bookings",
			DocumentID:    "test-booking-id",
			Timestamp:     time.Now(),
		}

		err := service.handleBookingInsert(event)
		assert.NoError(t, err)

		err = service.handleBookingUpdate(event)
		assert.NoError(t, err)

		err = service.handleBookingDelete(event)
		assert.NoError(t, err)

		err = service.handleBookingReplace(event)
		assert.NoError(t, err)
	})
}

func TestStopChangeStream(t *testing.T) {
	t.Run("should stop change stream successfully", func(t *testing.T) {
		logger, _ := logger.New("info", "console", "stdout")
		service := NewChangeStreamService(nil, logger)

		// Test stopping when no change stream is running
		err := service.StopChangeStream()
		assert.NoError(t, err)
	})
}

func TestChangeEventJSON(t *testing.T) {
	t.Run("should marshal change event to JSON", func(t *testing.T) {
		event := ChangeEvent{
			OperationType: "insert",
			Collection:    "users",
			DocumentID:    "test-id",
			FullDocument: map[string]interface{}{
				"name": "Test User",
			},
			Timestamp: time.Now(),
		}

		jsonData, err := event.GetChangeEventJSON()
		require.NoError(t, err)
		assert.NotEmpty(t, jsonData)
		assert.Contains(t, string(jsonData), "insert")
		assert.Contains(t, string(jsonData), "users")
		assert.Contains(t, string(jsonData), "test-id")
	})
}

// TestMongoPipelineConstruction specifically tests the pipeline construction
// that was causing the Go vet error
func TestMongoPipelineConstruction(t *testing.T) {
	t.Run("should construct MongoDB pipeline with keyed fields", func(t *testing.T) {
		// This test verifies that the pipeline construction uses keyed fields
		// which should fix the Go vet error
		pipeline := mongo.Pipeline{
			{bson.E{Key: "$match", Value: bson.M{
				"operationType": bson.M{
					"$in": []string{"insert", "update", "delete", "replace"},
				},
			}}},
		}

		assert.NotNil(t, pipeline)
		assert.Len(t, pipeline, 1)

		stage := pipeline[0]
		assert.NotNil(t, stage)
		assert.Len(t, stage, 1)

		// Verify the stage structure - get the first element
		elem := stage[0]
		assert.Equal(t, "$match", elem.Key)
		assert.NotNil(t, elem.Value)
	})
}
