package services_test

import (
	"testing"
	"time"

	"github.com/smorting/backend/internal/services"
	"go.uber.org/zap"
)

func newPCI(t *testing.T) *services.PCIDSSService {
	t.Helper()
	key := make([]byte, 32)
	lg, _ := zap.NewDevelopment()
	svc, err := services.NewPCIDSSService(key, lg)
	if err != nil {
		t.Fatalf("pci svc: %v", err)
	}
	svc.SetTokenStore(services.NewMemoryPaymentTokenStore())
	svc.SetTokenTTL(1 * time.Second)
	return svc
}

func TestTokenizeStoresEncryptedAndValidates(t *testing.T) {
	svc := newPCI(t)
	data := &services.SensitivePaymentData{CardNumber: "4111111111111111", CVV: "123", ExpiryMonth: "12", ExpiryYear: "2028"}
	token, err := svc.TokenizePaymentMethod(data, "user1")
	if err != nil {
		t.Fatalf("tokenize: %v", err)
	}
	if token.TokenID == "" || token.LastFour != "1111" {
		t.Fatalf("unexpected token metadata: %+v", token)
	}
	// Validate should succeed before expiry
	if _, err := svc.ValidatePaymentToken(token.TokenID); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestExpiredTokenInvalid(t *testing.T) {
	svc := newPCI(t)
	data := &services.SensitivePaymentData{CardNumber: "4111111111111111", CVV: "123", ExpiryMonth: "12", ExpiryYear: "2028"}
	token, err := svc.TokenizePaymentMethod(data, "user1")
	if err != nil {
		t.Fatalf("tokenize: %v", err)
	}
	time.Sleep(1500 * time.Millisecond)
	if _, err := svc.ValidatePaymentToken(token.TokenID); err == nil {
		t.Fatalf("expected expired token to be invalid")
	}
}

func TestDeleteTokenRemovesData(t *testing.T) {
	svc := newPCI(t)
	data := &services.SensitivePaymentData{CardNumber: "4111111111111111", CVV: "123", ExpiryMonth: "12", ExpiryYear: "2028"}
	token, err := svc.TokenizePaymentMethod(data, "user1")
	if err != nil {
		t.Fatalf("tokenize: %v", err)
	}
	if err := svc.DeletePaymentToken(token.TokenID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.ValidatePaymentToken(token.TokenID); err == nil {
		t.Fatalf("expected deleted token to be invalid")
	}
}
