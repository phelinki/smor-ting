package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
)

func TestMomoWebhook_Topup_UpdatesBalances(t *testing.T) {
	// Setup repo and user
	repo := database.NewMemoryDatabase()
	user := &models.User{Email: "webhook@example.com", Wallet: models.Wallet{Currency: "LRD"}}
	if err := repo.CreateUser(context.TODO(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	svc := services.NewWalletLedgerService(repo)

	// Webhook handler with ledger
	lg, _ := logger.New("debug", "console", "stdout")
	wh := handlers.NewWalletWebhookHandlerWithLedger(lg, svc)

	app := fiber.New()
	app.Post("/webhooks/momo", wh.MomoCallback)

	payload := map[string]any{
		"type":        "topup",
		"status":      "SUCCESSFUL",
		"amount":      500.0,
		"currency":    "LRD",
		"user_id":     user.ID.Hex(),
		"referenceId": "ref-123",
	}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/momo", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("webhook expected 200, got %d", resp.StatusCode)
	}

	// Compute balances
	bal, err := svc.ComputeBalances(context.TODO(), user.ID)
	if err != nil {
		t.Fatalf("balances: %v", err)
	}
	if bal.Available != 500 || bal.PendingHeld != 0 || bal.Total != 500 {
		t.Fatalf("unexpected balances: %+v", bal)
	}
}

func TestMomoWebhook_EscrowHoldAndRelease(t *testing.T) {
	repo := database.NewMemoryDatabase()
	user := &models.User{Email: "escrow@example.com", Wallet: models.Wallet{Currency: "LRD"}}
	_ = repo.CreateUser(context.TODO(), user)
	svc := services.NewWalletLedgerService(repo)
	lg, _ := logger.New("debug", "console", "stdout")
	wh := handlers.NewWalletWebhookHandlerWithLedger(lg, svc)

	app := fiber.New()
	app.Post("/webhooks/momo", wh.MomoCallback)

	// Escrow hold 200
	hold := map[string]any{"type": "escrow_hold", "status": "SUCCESSFUL", "amount": 200.0, "currency": "LRD", "user_id": user.ID.Hex(), "referenceId": "task-1"}
	b1, _ := json.Marshal(hold)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/momo", bytes.NewReader(b1))
	req.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(req)

	bal, _ := svc.ComputeBalances(context.TODO(), user.ID)
	if bal.Available != 0 || bal.PendingHeld != 200 || bal.Total != 200 {
		t.Fatalf("after hold unexpected balances: %+v", bal)
	}

	// Escrow release 200
	rel := map[string]any{"type": "escrow_release", "status": "SUCCESSFUL", "amount": 200.0, "currency": "LRD", "user_id": user.ID.Hex(), "referenceId": "task-1"}
	b2, _ := json.Marshal(rel)
	req2 := httptest.NewRequest(http.MethodPost, "/webhooks/momo", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(req2)

	bal2, _ := svc.ComputeBalances(context.TODO(), user.ID)
	if bal2.Available != 200 || bal2.PendingHeld != 0 || bal2.Total != 200 {
		t.Fatalf("after release unexpected balances: %+v", bal2)
	}
}
