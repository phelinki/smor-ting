package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
)

type fakeSmile struct{}

func (f *fakeSmile) SubmitKYC(ctx interface{}, req services.KYCRequest) (*services.KYCResult, error) {
	return &services.KYCResult{Status: "PENDING", Reference: "ref-1"}, nil
}

func TestKYC_Submit(t *testing.T) {
	lg, _ := logger.New("debug", "console", "stdout")
	// point to a test server to avoid nil pointer from real HTTP
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "PENDING", "reference": "ref-1"})
	}))
	defer srv.Close()
	client := services.NewSmileIDClient(srv.URL, "pid", "key")
	h := handlers.NewKYCHandler(client, lg)
	app := fiber.New()
	app.Post("/kyc/submit", h.Submit)

	body := map[string]string{"country": "LR", "id_type": "NIN", "id_number": "123", "first_name": "A", "last_name": "B", "phone": "231770000000"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/kyc/submit", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadGateway && resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}
