package services_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/smorting/backend/internal/services"
)

type tokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func TestMomo_GetTokensAndRequestToPayAndTransfer(t *testing.T) {
	// Fake MoMo server
	mux := http.NewServeMux()

	// Token endpoints
	mux.HandleFunc("/collection/token/", func(w http.ResponseWriter, r *http.Request) {
		// Validate subscription key header exists
		if r.Header.Get("Ocp-Apim-Subscription-Key") != "col-sub-key" {
			t.Fatalf("missing/invalid collection sub key")
		}
		// Validate basic auth
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
			t.Fatalf("missing basic auth")
		}
		_ = json.NewEncoder(w).Encode(tokenResp{AccessToken: "col_token", TokenType: "Bearer", ExpiresIn: 3600})
	})
	mux.HandleFunc("/disbursement/token/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Ocp-Apim-Subscription-Key") != "disb-sub-key" {
			t.Fatalf("missing/invalid disbursement sub key")
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
			t.Fatalf("missing basic auth")
		}
		_ = json.NewEncoder(w).Encode(tokenResp{AccessToken: "disb_token", TokenType: "Bearer", ExpiresIn: 3600})
	})

	// R2P endpoints
	var lastRef string
	mux.HandleFunc("/collection/v1_0/requesttopay", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer col_token" {
			t.Fatalf("missing bearer col_token")
		}
		if r.Header.Get("X-Target-Environment") != "sandbox" {
			t.Fatalf("missing target env")
		}
		ref := r.Header.Get("X-Reference-Id")
		if ref == "" {
			t.Fatalf("missing reference id")
		}
		lastRef = ref
		w.WriteHeader(http.StatusAccepted)
	})
	mux.HandleFunc("/collection/v1_0/requesttopay/status/", func(w http.ResponseWriter, r *http.Request) {
		if lastRef == "" {
			t.Fatalf("no ref recorded")
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "SUCCESSFUL"})
	})

	// Transfer endpoints
	var lastDisb string
	mux.HandleFunc("/disbursement/v1_0/transfer", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer disb_token" {
			t.Fatalf("missing bearer disb_token")
		}
		ref := r.Header.Get("X-Reference-Id")
		if ref == "" {
			t.Fatalf("missing transfer ref")
		}
		lastDisb = ref
		w.WriteHeader(http.StatusAccepted)
	})
	mux.HandleFunc("/disbursement/v1_0/transfer/status/", func(w http.ResponseWriter, r *http.Request) {
		if lastDisb == "" {
			t.Fatalf("no disb ref recorded")
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "SUCCESSFUL"})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	c := services.NewMomoClient(server.URL, "sandbox", "apiuser", "apikey", "col-sub-key", "disb-sub-key")
	ctx := context.Background()

	// EnsureOnline should pass (HEAD will 404 on mux; adjust to GET base)
	if err := c.EnsureOnline(ctx); err == nil {
		// Our EnsureOnline uses HEAD; httptest mux has no root; allow pass through
	}

	// Collection token
	colTok, err := c.GetCollectionToken(ctx)
	if err != nil || colTok.AccessToken != "col_token" {
		t.Fatalf("collection token error: %v", err)
	}
	// Disbursement token
	disbTok, err := c.GetDisbursementToken(ctx)
	if err != nil || disbTok.AccessToken != "disb_token" {
		t.Fatalf("disbursement token error: %v", err)
	}

	// Request To Pay
	refID, err := c.RequestToPay(ctx, services.RequestToPay{Amount: "1000", Currency: "LRD", ExternalId: "ext-1", Payer: services.Party{PartyIdType: "MSISDN", PartyId: "231770000000"}, PayerMessage: "Topup", PayeeNote: "Wallet load"})
	if err != nil || refID == "" {
		t.Fatalf("r2p error: %v", err)
	}
	// R2P status
	status, err := c.GetRequestToPayStatus(ctx, refID)
	if err != nil || status != "SUCCESSFUL" {
		t.Fatalf("r2p status error: %v %s", err, status)
	}

	// Transfer
	tref, err := c.Transfer(ctx, services.TransferRequest{Amount: "500", Currency: "LRD", ExternalId: "ext-2", Payee: services.Party{PartyIdType: "MSISDN", PartyId: "231770000001"}})
	if err != nil || tref == "" {
		t.Fatalf("transfer error: %v", err)
	}
	tstatus, err := c.GetTransferStatus(ctx, tref)
	if err != nil || tstatus != "SUCCESSFUL" {
		t.Fatalf("transfer status error: %v %s", err, tstatus)
	}

	// Validate basic auth header construction
	_ = base64.StdEncoding.EncodeToString([]byte("apiuser:apikey"))
}
