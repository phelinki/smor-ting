package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
)

func TestComputeBalances_AvailableAndPending(t *testing.T) {
	repo := database.NewMemoryDatabase()
	svc := services.NewWalletLedgerService(repo)
	// Seed user and get assigned ID
	_ = repo.CreateUser(context.TODO(), &models.User{Email: "u@example.com", Wallet: models.Wallet{Currency: "USD"}})
	u, err := repo.GetUserByEmail(context.TODO(), "u@example.com")
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	userID := u.ID

	// Completed topup 1000 LRD (credit -> available)
	_ = svc.RecordEntry(context.TODO(), &models.WalletLedgerEntry{
		UserID:    userID,
		Type:      models.LedgerTopup,
		Direction: models.LedgerCredit,
		Amount:    1000,
		Currency:  "LRD",
		Status:    models.LedgerCompleted,
		CreatedAt: time.Now(),
	})
	// Escrow hold 200 LRD (credit to held)
	_ = svc.RecordEntry(context.TODO(), &models.WalletLedgerEntry{
		UserID:    userID,
		Type:      models.LedgerEscrowHold,
		Direction: models.LedgerCredit,
		Amount:    200,
		Currency:  "LRD",
		Status:    models.LedgerPending,
		IsEscrow:  true,
		CreatedAt: time.Now(),
	})

	bal, err := svc.ComputeBalances(context.TODO(), userID)
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if bal.Available != 1000 || bal.PendingHeld != 200 || bal.Total != 1200 {
		t.Fatalf("unexpected balances: %+v", bal)
	}
}
