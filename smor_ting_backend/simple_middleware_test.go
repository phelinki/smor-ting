package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"github.com/smorting/backend/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSimpleMiddleware_DirectApplication tests middleware application directly
func TestSimpleMiddleware_DirectApplication(t *testing.T) {
	// Create a simple test app and JWT middleware
	testApp := fiber.New()
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	jwtSvc := services.NewJWTRefreshService(access, refresh, lg.Logger)
	authMiddleware, err := middleware.NewJWTAuthMiddleware(jwtSvc, repo, lg)
	require.NoError(t, err)

	// Test 1: Apply middleware directly to a single route
	testApp.Get("/direct-protected", authMiddleware.Authenticate(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "direct protected route"})
	})

	// Test 2: Apply middleware to a group exactly like main.go
	api := testApp.Group("/api/v1")
	services := api.Group("/services")
	services.Use(authMiddleware.Authenticate())
	services.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "group protected route"})
	})

	// Test the direct route
	req1 := httptest.NewRequest("GET", "/direct-protected", nil)
	resp1, err := testApp.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, 401, resp1.StatusCode, "Direct middleware should block unauthenticated requests")

	// Test the group route
	req2 := httptest.NewRequest("GET", "/api/v1/services/", nil)
	resp2, err := testApp.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, 401, resp2.StatusCode, "Group middleware should block unauthenticated requests")

	// Also test equivalent route path without trailing slash
	req3 := httptest.NewRequest("GET", "/api/v1/services", nil)
	resp3, err := testApp.Test(req3)
	require.NoError(t, err)
	assert.Equal(t, 401, resp3.StatusCode)
}

// TestSimpleMiddleware_RouteInspection inspects the actual routes in the main app
func TestSimpleMiddleware_RouteInspection(t *testing.T) {
	// Build a small app with middleware for inspection
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	jwtSvc := services.NewJWTRefreshService(access, refresh, lg.Logger)
	authMiddleware, err := middleware.NewJWTAuthMiddleware(jwtSvc, repo, lg)
	require.NoError(t, err)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Get("/services", authMiddleware.Authenticate(), func(c *fiber.Ctx) error { return c.SendStatus(200) })
	api.Get("/users/profile", authMiddleware.Authenticate(), func(c *fiber.Ctx) error { return c.SendStatus(200) })
	api.Post("/payments/tokenize", authMiddleware.Authenticate(), func(c *fiber.Ctx) error { return c.SendStatus(200) })
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	app.Post("/api/v1/auth/login", func(c *fiber.Ctx) error { return c.SendStatus(400) })

	// Test various route paths to see which ones are protected
	routes := []string{
		"/api/v1/services",
		"/api/v1/services/",
		"/api/v1/users/profile",
		"/api/v1/payments/tokenize",
		"/health",
		"/api/v1/auth/login",
	}

	for _, route := range routes {
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		t.Logf("Route %-25s -> Status: %d", route, resp.StatusCode)

		// All /api/v1/ routes except /auth should be protected (401)
		if strings.HasPrefix(route, "/api/v1/") && !strings.HasPrefix(route, "/api/v1/auth/") {
			if resp.StatusCode == 200 {
				t.Errorf("🚨 UNPROTECTED: %s returns 200 (should be 401)", route)
			}
		}
	}
}
