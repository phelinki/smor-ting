package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
)

type fakeMomo struct{ online bool }

func (f *fakeMomo) EnsureOnline(ctx context.Context) error {
	if !f.online {
		return errors.New("offline")
	}
	return nil
}
func (f *fakeMomo) RequestToPay(ctx context.Context, body services.RequestToPay) (string, error) {
	return "ref-topup", nil
}
func (f *fakeMomo) GetRequestToPayStatus(ctx context.Context, id string) (string, error) {
	return "SUCCESSFUL", nil
}
func (f *fakeMomo) Transfer(ctx context.Context, body services.TransferRequest) (string, error) {
	return "ref-withdraw", nil
}
func (f *fakeMomo) GetTransferStatus(ctx context.Context, id string) (string, error) {
	return "SUCCESSFUL", nil
}

func TestWallet_Topup_OnlineOnly(t *testing.T) {
	lg, _ := logger.New("debug", "console", "stdout")

	// offline
	app := fiber.New()
	h := handlers.NewWalletHandler(&fakeMomo{online: false}, lg)
	app.Post("/topup", h.Topup)

	req := httptest.NewRequest(http.MethodPost, "/topup", strings.NewReader(`{"amount":"100","currency":"LRD","msisdn":"231770000000"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 offline, got %d", resp.StatusCode)
	}

	// online
	app = fiber.New()
	h = handlers.NewWalletHandler(&fakeMomo{online: true}, lg)
	app.Post("/topup", h.Topup)
	req = httptest.NewRequest(http.MethodPost, "/topup", strings.NewReader(`{"amount":"100","currency":"LRD","msisdn":"231770000000"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 online, got %d", resp.StatusCode)
	}
}
