package services_test

import (
	"testing"

	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"go.uber.org/zap"
)

func newJWT(t *testing.T) *services.JWTRefreshService {
	t.Helper()
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	lg, _ := zap.NewDevelopment()
	return services.NewJWTRefreshService(access, refresh, lg)
}

func TestGenerateAndValidateTokenPair_Succeeds(t *testing.T) {
	jwt := newJWT(t)
	user := &models.User{Email: "a@example.com"}
	pair, err := jwt.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("generate pair: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatalf("expected tokens")
	}
	if _, err := jwt.ValidateAccessToken(pair.AccessToken); err != nil {
		t.Fatalf("validate access: %v", err)
	}
	if _, err := jwt.ValidateRefreshToken(pair.RefreshToken); err != nil {
		t.Fatalf("validate refresh: %v", err)
	}
}

func TestValidateRefreshToken_FailsForRevoked(t *testing.T) {
	jwt := newJWT(t)
	user := &models.User{Email: "a@example.com"}
	pair, _ := jwt.GenerateTokenPair(user)

	// Extract token_id by decoding token info
	info, err := jwt.GetTokenInfo(pair.RefreshToken, true)
	if err != nil {
		t.Fatalf("token info: %v", err)
	}
	tokenID, _ := info["token_id"].(string)
	if tokenID == "" {
		t.Fatalf("missing token_id")
	}

	// Revoke and then validate should fail
	if err := jwt.RevokeRefreshToken(tokenID); err != nil {
		t.Fatalf("revoke: %v", err)
	}

	if _, err := jwt.ValidateRefreshToken(pair.RefreshToken); err == nil {
		t.Fatalf("expected revoked token to fail")
	}
}
