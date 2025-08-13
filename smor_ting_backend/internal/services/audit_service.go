package services

import (
	"context"
	"time"

	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// AuditService handles audit logging for sensitive operations
type AuditService struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// AuditAction represents different types of auditable actions
type AuditAction string

const (
	ActionLogin               AuditAction = "LOGIN"
	ActionLogout              AuditAction = "LOGOUT"
	ActionPasswordChange      AuditAction = "PASSWORD_CHANGE"
	ActionUserCreate          AuditAction = "USER_CREATE"
	ActionUserUpdate          AuditAction = "USER_UPDATE"
	ActionUserDelete          AuditAction = "USER_DELETE"
	ActionRoleChange          AuditAction = "ROLE_CHANGE"
	ActionServiceCreate       AuditAction = "SERVICE_CREATE"
	ActionServiceUpdate       AuditAction = "SERVICE_UPDATE"
	ActionServiceDelete       AuditAction = "SERVICE_DELETE"
	ActionPaymentProcess      AuditAction = "PAYMENT_PROCESS"
	ActionPaymentRefund       AuditAction = "PAYMENT_REFUND"
	ActionWalletCreate        AuditAction = "WALLET_CREATE"
	ActionWalletLink          AuditAction = "WALLET_LINK"
	ActionWalletUnlink        AuditAction = "WALLET_UNLINK"
	ActionWalletDelete        AuditAction = "WALLET_DELETE"
	ActionDataExport          AuditAction = "DATA_EXPORT"
	ActionSystemConfiguration AuditAction = "SYSTEM_CONFIG"
	ActionKYCUpdate           AuditAction = "KYC_UPDATE"
	ActionKYCApproval         AuditAction = "KYC_APPROVAL"
	ActionKYCRejection        AuditAction = "KYC_REJECTION"
	ActionSessionRevoke       AuditAction = "SESSION_REVOKE"
	ActionBruteForceBlock     AuditAction = "BRUTE_FORCE_BLOCK"
)

// AuditEntry represents a single audit log entry
type AuditEntry struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Timestamp    time.Time              `bson:"timestamp" json:"timestamp"`
	UserID       string                 `bson:"user_id" json:"user_id"`
	UserEmail    string                 `bson:"user_email" json:"user_email"`
	UserRole     string                 `bson:"user_role" json:"user_role"`
	Action       AuditAction            `bson:"action" json:"action"`
	Resource     string                 `bson:"resource" json:"resource"`
	ResourceID   string                 `bson:"resource_id,omitempty" json:"resource_id,omitempty"`
	IPAddress    string                 `bson:"ip_address" json:"ip_address"`
	UserAgent    string                 `bson:"user_agent" json:"user_agent"`
	Details      map[string]interface{} `bson:"details,omitempty" json:"details,omitempty"`
	Success      bool                   `bson:"success" json:"success"`
	ErrorMessage string                 `bson:"error_message,omitempty" json:"error_message,omitempty"`
	SessionID    string                 `bson:"session_id,omitempty" json:"session_id,omitempty"`
}

// NewAuditService creates a new audit service
func NewAuditService(db *mongo.Database, logger *zap.Logger) *AuditService {
	// Handle nil database gracefully for tests
	var collection *mongo.Collection
	if db != nil {
		collection = db.Collection("audit_logs")

		// Create indexes for efficient querying
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Create indexes for efficient querying
			indexModels := []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "timestamp", Value: -1}},
				},
				{
					Keys: bson.D{{Key: "user_id", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "action", Value: 1}},
				},
				{
					Keys: bson.D{
						{Key: "user_id", Value: 1},
						{Key: "timestamp", Value: -1},
					},
				},
			}

			_, err := collection.Indexes().CreateMany(ctx, indexModels)
			if err != nil {
				// Log error but don't fail - indexes are not critical for basic functionality
			}
		}()
	}

	// Handle nil logger gracefully
	if logger == nil {
		logger = zap.NewNop()
	}

	return &AuditService{
		collection: collection,
		logger:     logger,
	}
}

// LogAction logs an audit entry for a user action
func (a *AuditService) LogAction(ctx context.Context, entry *AuditEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	if entry.ID.IsZero() {
		entry.ID = primitive.NewObjectID()
	}

	// Skip database operation if collection is nil (test mode)
	if a.collection != nil {
		_, err := a.collection.InsertOne(ctx, entry)
		if err != nil {
			a.logger.Error("Failed to insert audit log entry",
				zap.Error(err),
				zap.String("action", string(entry.Action)),
				zap.String("user_id", entry.UserID),
			)
			return err
		}
	}

	a.logger.Info("Audit log entry created",
		zap.String("action", string(entry.Action)),
		zap.String("user_id", entry.UserID),
		zap.String("resource", entry.Resource),
		zap.Bool("success", entry.Success),
	)

	return nil
}

// LogUserAction is a convenience method for logging user-initiated actions
func (a *AuditService) LogUserAction(ctx context.Context, user *models.User, action AuditAction, resource string, ipAddress string, userAgent string, success bool, details map[string]interface{}) error {
	entry := &AuditEntry{
		UserID:    user.ID.Hex(),
		UserEmail: user.Email,
		UserRole:  string(user.Role),
		Action:    action,
		Resource:  resource,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
		Details:   details,
	}

	return a.LogAction(ctx, entry)
}

// LogSystemAction logs actions performed by the system (no user context)
func (a *AuditService) LogSystemAction(ctx context.Context, action AuditAction, resource string, details map[string]interface{}) error {
	entry := &AuditEntry{
		UserID:   "system",
		Action:   action,
		Resource: resource,
		Success:  true,
		Details:  details,
	}

	return a.LogAction(ctx, entry)
}

// LogSecurityEvent logs security-related events like brute force attempts
func (a *AuditService) LogSecurityEvent(ctx context.Context, action AuditAction, email string, ipAddress string, userAgent string, details map[string]interface{}) error {
	entry := &AuditEntry{
		UserEmail: email,
		Action:    action,
		Resource:  "security",
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   false, // Security events are typically failed attempts
		Details:   details,
	}

	return a.LogAction(ctx, entry)
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (a *AuditService) GetUserAuditLogs(ctx context.Context, userID string, limit int, offset int) ([]*AuditEntry, error) {
	filter := bson.M{"user_id": userID}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := a.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*AuditEntry
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetActionAuditLogs retrieves audit logs for a specific action type
func (a *AuditService) GetActionAuditLogs(ctx context.Context, action AuditAction, limit int, offset int) ([]*AuditEntry, error) {
	filter := bson.M{"action": action}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := a.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*AuditEntry
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetRecentAuditLogs retrieves recent audit logs across all users and actions
func (a *AuditService) GetRecentAuditLogs(ctx context.Context, limit int) ([]*AuditEntry, error) {
	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := a.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*AuditEntry
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetAuditLogsInTimeRange retrieves audit logs within a specific time range
func (a *AuditService) GetAuditLogsInTimeRange(ctx context.Context, startTime, endTime time.Time, limit int) ([]*AuditEntry, error) {
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := a.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*AuditEntry
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// CountAuditLogs returns the total count of audit logs matching the filter
func (a *AuditService) CountAuditLogs(ctx context.Context, filter bson.M) (int64, error) {
	return a.collection.CountDocuments(ctx, filter)
}

// DeleteOldAuditLogs removes audit logs older than the specified duration (for compliance/cleanup)
func (a *AuditService) DeleteOldAuditLogs(ctx context.Context, olderThan time.Time) (int64, error) {
	filter := bson.M{
		"timestamp": bson.M{
			"$lt": olderThan,
		},
	}

	result, err := a.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	a.logger.Info("Deleted old audit logs",
		zap.Int64("count", result.DeletedCount),
		zap.Time("older_than", olderThan),
	)

	return result.DeletedCount, nil
}
