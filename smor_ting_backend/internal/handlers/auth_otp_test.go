package handlers_test

import (
	"context"
	"encoding/json"
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
)

// helper to build handler with memory repo we can inspect
func newAuthHandlerWithMemoryRepo(t *testing.T) (*handlers.AuthHandler, *database.MemoryDatabase, *services.JWTRefreshService) {
	t.Helper()
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	// jwt service
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 3
	}
	for i := range refresh {
		refresh[i] = 4
	}
	jwt := services.NewJWTRefreshService(access, refresh, lg.Logger)
	// enc service
	encKey := make([]byte, 32)
	enc, _ := services.NewEncryptionService(encKey)
	// mongo-like auth service
	authCfg := &configs.AuthConfig{JWTSecret: "test", JWTExpiration: 0, BCryptCost: 10}
	a, _ := auth.NewMongoDBService(repo, authCfg, lg)
	h := handlers.NewAuthHandler(jwt, enc, lg, a)
	return h, repo, jwt
}

func TestResendAndVerifyOTP_LoginFlow(t *testing.T) {
	h, repo, jwt := newAuthHandlerWithMemoryRepo(t)
	_ = jwt
	app := fiber.New()
	app.Post("/auth/resend-otp", h.ResendOTP)
	app.Post("/auth/verify-otp", h.VerifyOTP)

	// Seed a user not yet verified
	u := &models.User{Email: "otpuser@example.com", Password: "$2a$10$abcdefghijklmnopqrstuv", Role: models.CustomerRole, IsEmailVerified: false}
	if err := repo.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	// Resend OTP
	resendReq := httptest.NewRequest(http.MethodPost, "/auth/resend-otp", strings.NewReader(`{"email":"otpuser@example.com"}`))
	resendReq.Header.Set("Content-Type", "application/json")
	resendResp, _ := app.Test(resendReq)
	if resendResp.StatusCode != http.StatusOK {
		t.Fatalf("resend otp expected 200, got %d", resendResp.StatusCode)
	}

	// Fetch OTP from memory repo (there should be one)
	// iterate internal store through exported method
	otp, err := repo.GetOTP(context.Background(), "otpuser@example.com", "000000")
	if err == nil && otp != nil {
		// unlikely; we didn't know the code; ignore
	}
	// Since we can't know the OTP value, iterate by trying to find any valid OTP by temporarily exposing via helper
	// Instead, we create one ourselves directly and verify with that to keep the test deterministic
	// Create deterministic OTP
	customOTP := &models.OTPRecord{Email: "otpuser@example.com", OTP: "123456", Purpose: "login", IsUsed: false}
	if err := repo.CreateOTP(context.Background(), customOTP); err != nil {
		t.Fatalf("create otp: %v", err)
	}

	// Verify OTP
	verifyReq := httptest.NewRequest(http.MethodPost, "/auth/verify-otp", strings.NewReader(`{"email":"otpuser@example.com","otp":"123456"}`))
	verifyReq.Header.Set("Content-Type", "application/json")
	verifyResp, _ := app.Test(verifyReq)
	if verifyResp.StatusCode != http.StatusOK {
		t.Fatalf("verify otp expected 200, got %d", verifyResp.StatusCode)
	}
	var payload map[string]interface{}
	_ = json.NewDecoder(verifyResp.Body).Decode(&payload)
	if payload["access_token"] == "" {
		t.Fatalf("expected access_token in response")
	}
}

func TestPasswordResetFlow(t *testing.T) {
	h, repo, _ := newAuthHandlerWithMemoryRepo(t)
	app := fiber.New()
	app.Post("/auth/request-password-reset", h.RequestPasswordReset)
	app.Post("/auth/reset-password", h.ResetPassword)

	// Seed a user
	u := &models.User{Email: "reset@example.com", Password: "$2a$10$abcdefghijklmnopqrstuv", Role: models.CustomerRole, IsEmailVerified: true}
	if err := repo.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	// Request reset
	req := httptest.NewRequest(http.MethodPost, "/auth/request-password-reset", strings.NewReader(`{"email":"reset@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("request reset expected 200, got %d", resp.StatusCode)
	}

	// Create deterministic OTP and attempt reset
	customOTP := &models.OTPRecord{Email: "reset@example.com", OTP: "654321", Purpose: "password_reset", IsUsed: false}
	if err := repo.CreateOTP(context.Background(), customOTP); err != nil {
		t.Fatalf("create otp: %v", err)
	}

	resetReq := httptest.NewRequest(http.MethodPost, "/auth/reset-password", strings.NewReader(`{"email":"reset@example.com","otp":"654321","new_password":"NewPass123!"}`))
	resetReq.Header.Set("Content-Type", "application/json")
	resetResp, _ := app.Test(resetReq)
	if resetResp.StatusCode != http.StatusOK {
		t.Fatalf("reset expected 200, got %d", resetResp.StatusCode)
	}
}
