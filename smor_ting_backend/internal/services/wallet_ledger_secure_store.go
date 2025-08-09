package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// WalletLedgerSecureStore persists wallet ledger entries with encryption-at-rest.
// MongoDB Atlas is the system of record.
type WalletLedgerSecureStore struct {
	coll   *mongo.Collection
	enc    *EncryptionService
	logger *logger.Logger
}

func NewWalletLedgerSecureStore(db *mongo.Database, enc *EncryptionService, logger *logger.Logger) *WalletLedgerSecureStore {
	return &WalletLedgerSecureStore{coll: db.Collection("wallet_ledger"), enc: enc, logger: logger}
}

// EncryptLedgerEntryForStorage serializes and encrypts a ledger entry for storage
func EncryptLedgerEntryForStorage(enc *EncryptionService, entry *models.WalletLedgerEntry) (string, error) {
	if entry == nil {
		return "", fmt.Errorf("entry is nil")
	}
	raw, err := json.Marshal(entry)
	if err != nil {
		return "", fmt.Errorf("marshal entry: %w", err)
	}
	return enc.Encrypt(raw)
}

// SaveEncrypted stores an entry encrypted, alongside minimal metadata fields for querying
func (s *WalletLedgerSecureStore) SaveEncrypted(ctx context.Context, entry *models.WalletLedgerEntry) error {
	cipher, err := EncryptLedgerEntryForStorage(s.enc, entry)
	if err != nil {
		return err
	}
	doc := bson.M{
		"user_id":      entry.UserID,
		"type":         entry.Type,
		"direction":    entry.Direction,
		"amount":       entry.Amount,
		"currency":     entry.Currency,
		"status":       entry.Status,
		"is_escrow":    entry.IsEscrow,
		"reference":    entry.Reference,
		"provider_ref": entry.ProviderRef,
		"created_at":   entry.CreatedAt,
		"updated_at":   entry.UpdatedAt,
		"encrypted":    cipher,
	}
	_, err = s.coll.InsertOne(ctx, doc)
	return err
}
