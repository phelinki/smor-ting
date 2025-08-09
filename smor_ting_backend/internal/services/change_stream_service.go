package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// ChangeEvent represents a database change event
type ChangeEvent struct {
	OperationType string                 `json:"operation_type"`
	Collection    string                 `json:"collection"`
	DocumentID    string                 `json:"document_id"`
	FullDocument  map[string]interface{} `json:"full_document,omitempty"`
	UpdatedFields map[string]interface{} `json:"updated_fields,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
}

// ChangeStreamService handles MongoDB change streams for real-time sync
type ChangeStreamService struct {
	db           *mongo.Database
	logger       *logger.Logger
	changeStream *mongo.ChangeStream
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewChangeStreamService creates a new change stream service
func NewChangeStreamService(db *mongo.Database, logger *logger.Logger) *ChangeStreamService {
	ctx, cancel := context.WithCancel(context.Background())
	return &ChangeStreamService{
		db:     db,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// StartChangeStream starts monitoring database changes
func (cs *ChangeStreamService) StartChangeStream() error {
	// Watch all collections for changes
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"operationType": bson.M{
				"$in": []string{"insert", "update", "delete", "replace"},
			},
		}}},
	}

	// Create change stream
	changeStream, err := cs.db.Watch(cs.ctx, pipeline)
	if err != nil {
		return fmt.Errorf("failed to create change stream: %w", err)
	}

	cs.changeStream = changeStream
	cs.logger.Info("Change stream started successfully")

	// Start listening for changes
	go cs.listenForChanges()

	return nil
}

// listenForChanges listens for database changes
func (cs *ChangeStreamService) listenForChanges() {
	for cs.changeStream.Next(cs.ctx) {
		var changeEvent bson.M
		if err := cs.changeStream.Decode(&changeEvent); err != nil {
			cs.logger.Error("Failed to decode change event", err)
			continue
		}

		// Process the change event
		if err := cs.processChangeEvent(changeEvent); err != nil {
			cs.logger.Error("Failed to process change event", err)
		}
	}

	if err := cs.changeStream.Err(); err != nil {
		cs.logger.Error("Change stream error", err)
	}
}

// processChangeEvent processes a single change event
func (cs *ChangeStreamService) processChangeEvent(event bson.M) error {
	// Extract event details
	operationType, ok := event["operationType"].(string)
	if !ok {
		return fmt.Errorf("invalid operation type")
	}

	collection, ok := event["ns"].(bson.M)
	if !ok {
		return fmt.Errorf("invalid namespace")
	}

	collectionName, ok := collection["coll"].(string)
	if !ok {
		return fmt.Errorf("invalid collection name")
	}

	// Create change event
	changeEvent := ChangeEvent{
		OperationType: operationType,
		Collection:    collectionName,
		Timestamp:     time.Now(),
	}

	// Extract document ID
	if documentKey, ok := event["documentKey"].(bson.M); ok {
		if id, ok := documentKey["_id"]; ok {
			changeEvent.DocumentID = fmt.Sprintf("%v", id)
		}
	}

	// Extract full document for insert/replace operations
	if operationType == "insert" || operationType == "replace" {
		if fullDocument, ok := event["fullDocument"].(bson.M); ok {
			changeEvent.FullDocument = fullDocument
		}
	}

	// Extract updated fields for update operations
	if operationType == "update" {
		if updateDescription, ok := event["updateDescription"].(bson.M); ok {
			if updatedFields, ok := updateDescription["updatedFields"].(bson.M); ok {
				changeEvent.UpdatedFields = updatedFields
			}
		}
	}

	// Log the change event
	cs.logger.Info("Database change detected",
		zap.String("operation", operationType),
		zap.String("collection", collectionName),
		zap.String("document_id", changeEvent.DocumentID),
	)

	// Handle different types of changes
	switch operationType {
	case "insert":
		return cs.handleInsert(changeEvent)
	case "update":
		return cs.handleUpdate(changeEvent)
	case "delete":
		return cs.handleDelete(changeEvent)
	case "replace":
		return cs.handleReplace(changeEvent)
	}

	return nil
}

// handleInsert handles insert operations
func (cs *ChangeStreamService) handleInsert(event ChangeEvent) error {
	switch event.Collection {
	case "users":
		return cs.handleUserInsert(event)
	case "services":
		return cs.handleServiceInsert(event)
	case "bookings":
		return cs.handleBookingInsert(event)
	}
	return nil
}

// handleUpdate handles update operations
func (cs *ChangeStreamService) handleUpdate(event ChangeEvent) error {
	switch event.Collection {
	case "users":
		return cs.handleUserUpdate(event)
	case "services":
		return cs.handleServiceUpdate(event)
	case "bookings":
		return cs.handleBookingUpdate(event)
	}
	return nil
}

// handleDelete handles delete operations
func (cs *ChangeStreamService) handleDelete(event ChangeEvent) error {
	switch event.Collection {
	case "users":
		return cs.handleUserDelete(event)
	case "services":
		return cs.handleServiceDelete(event)
	case "bookings":
		return cs.handleBookingDelete(event)
	}
	return nil
}

// handleReplace handles replace operations
func (cs *ChangeStreamService) handleReplace(event ChangeEvent) error {
	switch event.Collection {
	case "users":
		return cs.handleUserReplace(event)
	case "services":
		return cs.handleServiceReplace(event)
	case "bookings":
		return cs.handleBookingReplace(event)
	}
	return nil
}

// User change handlers
func (cs *ChangeStreamService) handleUserInsert(event ChangeEvent) error {
	cs.logger.Info("New user created", zap.String("user_id", event.DocumentID))
	// Here you could notify other services, send push notifications, etc.
	return nil
}

func (cs *ChangeStreamService) handleUserUpdate(event ChangeEvent) error {
	cs.logger.Info("User updated", zap.String("user_id", event.DocumentID))

	// Check if wallet was updated
	if _, hasWalletUpdate := event.UpdatedFields["wallet"]; hasWalletUpdate {
		cs.logger.Info("User wallet updated", zap.String("user_id", event.DocumentID))
		// Handle wallet updates
	}

	return nil
}

func (cs *ChangeStreamService) handleUserDelete(event ChangeEvent) error {
	cs.logger.Info("User deleted", zap.String("user_id", event.DocumentID))
	return nil
}

func (cs *ChangeStreamService) handleUserReplace(event ChangeEvent) error {
	cs.logger.Info("User replaced", zap.String("user_id", event.DocumentID))
	return nil
}

// Service change handlers
func (cs *ChangeStreamService) handleServiceInsert(event ChangeEvent) error {
	cs.logger.Info("New service created", zap.String("service_id", event.DocumentID))
	return nil
}

func (cs *ChangeStreamService) handleServiceUpdate(event ChangeEvent) error {
	cs.logger.Info("Service updated", zap.String("service_id", event.DocumentID))
	return nil
}

func (cs *ChangeStreamService) handleServiceDelete(event ChangeEvent) error {
	cs.logger.Info("Service deleted", zap.String("service_id", event.DocumentID))
	return nil
}

func (cs *ChangeStreamService) handleServiceReplace(event ChangeEvent) error {
	cs.logger.Info("Service replaced", zap.String("service_id", event.DocumentID))
	return nil
}

// Booking change handlers
func (cs *ChangeStreamService) handleBookingInsert(event ChangeEvent) error {
	cs.logger.Info("New booking created", zap.String("booking_id", event.DocumentID))
	return nil
}

func (cs *ChangeStreamService) handleBookingUpdate(event ChangeEvent) error {
	cs.logger.Info("Booking updated", zap.String("booking_id", event.DocumentID))

	// Check if status was updated
	if _, hasStatusUpdate := event.UpdatedFields["status"]; hasStatusUpdate {
		cs.logger.Info("Booking status updated", zap.String("booking_id", event.DocumentID))
		// Handle status updates - could trigger notifications
	}

	return nil
}

func (cs *ChangeStreamService) handleBookingDelete(event ChangeEvent) error {
	cs.logger.Info("Booking deleted", zap.String("booking_id", event.DocumentID))
	return nil
}

func (cs *ChangeStreamService) handleBookingReplace(event ChangeEvent) error {
	cs.logger.Info("Booking replaced", zap.String("booking_id", event.DocumentID))
	return nil
}

// StopChangeStream stops the change stream
func (cs *ChangeStreamService) StopChangeStream() error {
	if cs.changeStream != nil {
		cs.cancel()
		return cs.changeStream.Close(cs.ctx)
	}
	return nil
}

// GetChangeEventJSON returns the change event as JSON
func (event ChangeEvent) GetChangeEventJSON() ([]byte, error) {
	return json.Marshal(event)
}
