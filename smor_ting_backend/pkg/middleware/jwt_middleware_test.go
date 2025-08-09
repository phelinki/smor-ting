package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/configs"
	authpkg "github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	mw "github.com/smorting/backend/pkg/middleware"
	"golang.org/x/crypto/bcrypt"
)

func newJWTDeps(t *testing.T) (*services.JWTRefreshService, database.Repository, *logger.Logger, *authpkg.MongoDBService) {
	t.Helper()
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	// use 32-byte secrets
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	jwtSvc := services.NewJWTRefreshService(access, refresh, lg.Logger)
	authSvc, _ := authpkg.NewMongoDBService(repo, &configs.AuthConfig{BCryptCost: 10, JWTSecret: "x"}, lg)
	return jwtSvc, repo, lg, authSvc
}

func createUserWithPassword(t *testing.T, repo database.Repository, email, password string) *models.User {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := &models.User{Email: email, Password: string(hash), Role: models.CustomerRole}
	if err := repo.CreateUser(context.TODO(), user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	return user
}

func TestJWTMiddleware_BlocksMissingToken(t *testing.T) {
	jwtSvc, repo, lg, _ := newJWTDeps(t)
	m, _ := mw.NewJWTAuthMiddleware(jwtSvc, repo, lg)

	app := fiber.New()
	app.Get("/protected", m.Authenticate(), func(c *fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestJWTMiddleware_AllowsValidTokenAndSetsUser(t *testing.T) {
	jwtSvc, repo, lg, _ := newJWTDeps(t)
	m, _ := mw.NewJWTAuthMiddleware(jwtSvc, repo, lg)

	// Create user and token
	user := createUserWithPassword(t, repo, "mw@example.com", "Password1!")
	pair, err := jwtSvc.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("failed to generate tokens: %v", err)
	}

	app := fiber.New()
	app.Get("/protected", m.Authenticate(), func(c *fiber.Ctx) error {
		if u, ok := mw.GetUserFromContextModels(c); !ok || u == nil || u.Email != user.Email {
			t.Fatalf("user not in context")
		}
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
