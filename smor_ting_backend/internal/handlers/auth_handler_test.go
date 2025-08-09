package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func newAuthHandler(t *testing.T) (*handlers.AuthHandler, database.Repository, *services.JWTRefreshService) {
	t.Helper()
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	// JWT service
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	jwt := services.NewJWTRefreshService(access, refresh, lg.Logger)
	// encryption service
	encKey := make([]byte, 32)
	enc, _ := services.NewEncryptionService(encKey)
	// auth service (Mongo-like)
	// Minimal config is fine for this test path; we won't call Login here
	// Create a minimal valid auth config for MongoDB service
	authCfg := &configs.AuthConfig{JWTSecret: "test", JWTExpiration: 0, BCryptCost: 10}
	a, _ := auth.NewMongoDBService(repo, authCfg, lg)
	h := handlers.NewAuthHandler(jwt, enc, lg, a)
	return h, repo, jwt
}

func TestRevokeThenRefreshFails(t *testing.T) {
	h, _, jwt := newAuthHandler(t)
	app := fiber.New()
	app.Post("/refresh", h.RefreshToken)
	app.Post("/revoke", h.RevokeToken)

	// Create a token pair directly from the same JWT service instance
	user := &models.User{ID: primitive.NewObjectID(), Email: "u@example.com", Role: models.CustomerRole}
	pair, err := jwt.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("generate pair: %v", err)
	}

	// Revoke
	revokeReq := httptest.NewRequest(http.MethodPost, "/revoke", strings.NewReader(`{"refresh_token":"`+pair.RefreshToken+`"}`))
	revokeReq.Header.Set("Content-Type", "application/json")
	revokeResp, _ := app.Test(revokeReq)
	if revokeResp.StatusCode != http.StatusOK {
		t.Fatalf("revoke expected 200, got %d", revokeResp.StatusCode)
	}

	// Refresh should fail now
	refreshReq := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(`{"refresh_token":"`+pair.RefreshToken+`"}`))
	refreshReq.Header.Set("Content-Type", "application/json")
	refreshResp, _ := app.Test(refreshReq)
	if refreshResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("refresh expected 401, got %d", refreshResp.StatusCode)
	}
}
