package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smorting/backend/internal/services"
)

func TestSmileID_SubmitKYC(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/kyc", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Partner-ID") != "pid" || r.Header.Get("X-API-Key") != "key" {
			t.Fatalf("missing headers")
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "PENDING", "reference": "ref-1"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := services.NewSmileIDClient(srv.URL, "pid", "key")
	res, err := c.SubmitKYC(context.Background(), services.KYCRequest{Country: "LR", IDType: "NIN", IDNumber: "123", FirstName: "A", LastName: "B", Phone: "231770000000"})
	if err != nil || res.Reference == "" {
		t.Fatalf("submit kyc failed: %v", err)
	}
}
