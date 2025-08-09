package services

import (
	"context"
	"errors"
	"sync"
	"time"
)

// TokenRevocationStore abstracts persistence for revoked token IDs (JTI)
type TokenRevocationStore interface {
	// Revoke stores the tokenID as revoked until expiresAt
	Revoke(tokenID string, expiresAt time.Time) error
	// IsRevoked returns true if tokenID is already revoked and not expired
	IsRevoked(tokenID string) (bool, error)
	// PurgeExpired removes expired revocation records (optional best-effort)
	PurgeExpired(ctx context.Context) error
}

// memoryRevocationStore provides in-memory implementation for tests and dev
type memoryRevocationStore struct {
	mu      sync.RWMutex
	revoked map[string]time.Time
}

func NewMemoryRevocationStore() TokenRevocationStore {
	return &memoryRevocationStore{revoked: make(map[string]time.Time)}
}

func (m *memoryRevocationStore) Revoke(tokenID string, expiresAt time.Time) error {
	if tokenID == "" {
		return errors.New("tokenID required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.revoked[tokenID] = expiresAt
	return nil
}

func (m *memoryRevocationStore) IsRevoked(tokenID string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	exp, ok := m.revoked[tokenID]
	if !ok {
		return false, nil
	}
	if time.Now().After(exp) {
		return false, nil
	}
	return true, nil
}

func (m *memoryRevocationStore) PurgeExpired(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for id, exp := range m.revoked {
		if now.After(exp) {
			delete(m.revoked, id)
		}
	}
	return nil
}
