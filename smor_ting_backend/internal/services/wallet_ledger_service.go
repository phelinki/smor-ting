package services

import (
	"context"
	"errors"
	"time"

	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletLedgerService struct {
	repo   database.Repository
	secure *WalletLedgerSecureStore
}

func NewWalletLedgerService(repo database.Repository) *WalletLedgerService {
	return &WalletLedgerService{repo: repo}
}

// AttachSecureStore wires an encrypted MongoDB Atlas-backed store
func AttachSecureStore(svc *WalletLedgerService, store *WalletLedgerSecureStore) *WalletLedgerService {
	if svc != nil {
		svc.secure = store
	}
	return svc
}

// RecordEntry persists a ledger entry. For memory repo we’ll extend with minimal support.
func (s *WalletLedgerService) RecordEntry(ctx context.Context, entry *models.WalletLedgerEntry) error {
	if entry == nil {
		return errors.New("entry is nil")
	}
	entry.ID = primitive.NewObjectID()
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	// Persist encrypted copy in system-of-record (Mongo) when available
	if s.secure != nil {
		_ = s.secure.SaveEncrypted(ctx, entry)
	}
	// For now, store as user wallet transaction shadow in memory; real impl should use dedicated collection
	user, err := s.repo.GetUserByID(ctx, entry.UserID)
	if err != nil {
		return err
	}
	user.Wallet.Transactions = append(user.Wallet.Transactions, models.Transaction{
		ID:          entry.ID,
		Type:        string(entry.Type),
		Amount:      entry.Amount,
		Description: entry.ProviderRef,
		Reference:   entry.Reference,
		Status:      string(entry.Status),
		CreatedAt:   entry.CreatedAt,
	})
	// Update wallet balance for completed entries
	switch entry.Type {
	case models.LedgerTopup:
		if entry.Status == models.LedgerCompleted && entry.Direction == models.LedgerCredit {
			user.Wallet.Balance += entry.Amount
		}
	case models.LedgerWithdraw:
		if entry.Status == models.LedgerCompleted && entry.Direction == models.LedgerDebit {
			user.Wallet.Balance -= entry.Amount
		}
	case models.LedgerEscrowHold:
		// pending held handled at compute-time; no immediate balance change
	case models.LedgerEscrowRelease:
		if entry.Status == models.LedgerCompleted {
			user.Wallet.Balance += entry.Amount
			// Mark corresponding escrow hold as completed to remove from pending
			for i := range user.Wallet.Transactions {
				tx := &user.Wallet.Transactions[i]
				if tx.Type == string(models.LedgerEscrowHold) && tx.Reference == entry.Reference && tx.Status == string(models.LedgerPending) {
					tx.Status = string(models.LedgerCompleted)
				}
			}
		}
	}
	return s.repo.UpdateUser(ctx, user)
}

func (s *WalletLedgerService) ComputeBalances(ctx context.Context, userID primitive.ObjectID) (*models.WalletBalances, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var pending float64
	// In a real impl we would query ledger entries; for demo use transactions' status and type
	for _, tx := range user.Wallet.Transactions {
		if tx.Status == string(models.LedgerPending) && (tx.Type == string(models.LedgerEscrowHold) || tx.Type == string(models.LedgerPayment)) {
			pending += tx.Amount
		}
	}
	total := user.Wallet.Balance + pending
	return &models.WalletBalances{Available: user.Wallet.Balance, PendingHeld: pending, Total: total, Currency: user.Wallet.Currency}, nil
}
