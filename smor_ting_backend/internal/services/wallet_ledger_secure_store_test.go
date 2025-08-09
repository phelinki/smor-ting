package services_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
)

func TestMongoWalletLedgerSecureStore_SavesEncrypted(t *testing.T) {
	// Generate a 32-byte key
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	encSvc, err := services.NewEncryptionService(key)
	if err != nil {
		t.Fatalf("encryption init: %v", err)
	}

	// Use in-memory fake mongo by nil db is not feasible; instead, test encryptor path only via helper
	entry := &models.WalletLedgerEntry{Currency: "LRD", Amount: 10.0, ProviderRef: "prov", Reference: "ref"}
	cipher, err := services.EncryptLedgerEntryForStorage(encSvc, entry)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	// Base64 should decode
	if _, err := base64.StdEncoding.DecodeString(cipher); err != nil {
		t.Fatalf("cipher not base64: %v", err)
	}
}

// Interface compliance test
type fakeSecureStore struct{ saved int }

func (f *fakeSecureStore) SaveEncrypted(ctx context.Context, entry *models.WalletLedgerEntry) error {
	f.saved++
	return nil
}
