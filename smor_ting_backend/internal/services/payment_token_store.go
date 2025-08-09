package services

import (
	"context"
	"errors"
	"sync"
	"time"
)

// PaymentTokenRecord is the stored representation of a tokenized method
// Note: EncryptedData must be ciphertext (AES-256-GCM base64)
// No plaintext PAN/CVV should be here
// Metadata fields are non-sensitive summary derived from PAN
// (e.g., last4, brand)
type PaymentTokenRecord struct {
	TokenID       string
	UserID        string
	EncryptedData string
	TokenType     string
	LastFour      string
	Brand         string
	CreatedAt     time.Time
	LastUsed      time.Time
	ExpiresAt     time.Time
}

// PaymentTokenStore abstracts secure storage with TTL/expiry
// Implementations should enforce encryption-at-rest and TTL via DB config (e.g., Mongo TTL index)
type PaymentTokenStore interface {
	Save(tokenID, userID, encryptedData string, expiresAt time.Time) error
	Get(tokenID string) (*PaymentTokenRecord, error)
	Delete(tokenID string) error
	TouchLastUsed(tokenID string, t time.Time) error
	PurgeExpired(ctx context.Context) error
}

// In-memory implementation for tests/dev
type memoryPaymentTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]*PaymentTokenRecord
}

func NewMemoryPaymentTokenStore() PaymentTokenStore {
	return &memoryPaymentTokenStore{tokens: make(map[string]*PaymentTokenRecord)}
}

func (m *memoryPaymentTokenStore) Save(tokenID, userID, encryptedData string, expiresAt time.Time) error {
	if tokenID == "" || encryptedData == "" {
		return errors.New("tokenID and encryptedData required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	rec := &PaymentTokenRecord{
		TokenID:       tokenID,
		UserID:        userID,
		EncryptedData: encryptedData,
		CreatedAt:     time.Now(),
		LastUsed:      time.Now(),
		ExpiresAt:     expiresAt,
	}
	m.tokens[tokenID] = rec
	return nil
}

func (m *memoryPaymentTokenStore) Get(tokenID string) (*PaymentTokenRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rec, ok := m.tokens[tokenID]
	if !ok {
		return nil, errors.New("not found")
	}
	return rec, nil
}

func (m *memoryPaymentTokenStore) Delete(tokenID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tokens, tokenID)
	return nil
}

func (m *memoryPaymentTokenStore) TouchLastUsed(tokenID string, t time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.tokens[tokenID]
	if !ok {
		return errors.New("not found")
	}
	rec.LastUsed = t
	return nil
}

func (m *memoryPaymentTokenStore) PurgeExpired(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for id, rec := range m.tokens {
		if now.After(rec.ExpiresAt) {
			delete(m.tokens, id)
		}
	}
	return nil
}
