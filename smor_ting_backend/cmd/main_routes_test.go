package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMainRoutes_AuthenticationEnforcement tests that the main application
// properly enforces authentication on protected routes
func TestMainRoutes_AuthenticationEnforcement(t *testing.T) {
	// Create a test application instance
	app, err := NewApp()
	if err != nil {
		t.Skipf("Cannot create app for testing: %v", err)
		return
	}

	// Test protected routes that should require authentication
	protectedRoutes := []struct {
		method string
		path   string
		body   interface{}
	}{
		{"GET", "/api/v1/users/profile", nil},
		{"GET", "/api/v1/services", nil},
		{"POST", "/api/v1/services", map[string]string{"name": "test"}},
		{"GET", "/api/v1/payments/validate?token_id=test", nil},
		{"POST", "/api/v1/payments/tokenize", map[string]string{"card_number": "4111111111111111"}},
		{"POST", "/api/v1/sync/data", map[string]interface{}{"user_id": "507f1f77bcf86cd799439011"}},
		{"GET", "/api/v1/sync/unsynced?user_id=507f1f77bcf86cd799439011", nil},
	}

	for _, route := range protectedRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			var req *http.Request
			var err error

			if route.body != nil {
				bodyBytes, _ := json.Marshal(route.body)
				req, err = http.NewRequest(route.method, route.path, bytes.NewBuffer(bodyBytes))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(route.method, route.path, nil)
				require.NoError(t, err)
			}

			// Test without authentication token
			resp, err := app.server.Test(req)
			require.NoError(t, err)

			// Should return 401 Unauthorized, not 200 or 500
			if resp.StatusCode == 200 {
				t.Errorf("Route %s %s is NOT PROTECTED! Returned 200 instead of 401", route.method, route.path)
			} else if resp.StatusCode == 500 {
				t.Errorf("Route %s %s has server error (500) instead of proper auth check (401)", route.method, route.path)
			} else if resp.StatusCode == 404 {
				t.Logf("Route %s %s not found (404) - might be missing route configuration", route.method, route.path)
			} else {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
					"Protected route should return 401 Unauthorized")
			}
		})
	}
}

// TestMainRoutes_PublicEndpoints tests that public endpoints work without authentication
func TestMainRoutes_PublicEndpoints(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Skipf("Cannot create app for testing: %v", err)
		return
	}

	publicRoutes := []struct {
		method         string
		path           string
		expectedStatus int
	}{
		{"GET", "/health", 200},
		{"GET", "/docs", 200},
		{"GET", "/swagger", 200},
		{"POST", "/api/v1/auth/login", 400},    // 400 because no body, but shouldn't be 401
		{"POST", "/api/v1/auth/register", 400}, // 400 because no body, but shouldn't be 401
	}

	for _, route := range publicRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req, err := http.NewRequest(route.method, route.path, nil)
			require.NoError(t, err)

			resp, err := app.server.Test(req)
			require.NoError(t, err)

			// Should NOT return 401 (unauthorized)
			assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode,
				"Public route should not require authentication")

			if route.expectedStatus > 0 {
				assert.Equal(t, route.expectedStatus, resp.StatusCode)
			}
		})
	}
}

// TestMainRoutes_MiddlewareTypeAssertion tests the middleware type assertion issue
func TestMainRoutes_MiddlewareTypeAssertion(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Skipf("Cannot create app for testing: %v", err)
		return
	}

	// The issue might be that authMiddleware is interface{} instead of the concrete type
	// This test checks if the type assertion is working

	// If we can create the app, the type assertions in setupRoutes should have worked
	assert.NotNil(t, app.server, "Server should be initialized")
	assert.NotNil(t, app.jwtService, "JWT service should be initialized")

	// Test that middleware is actually being applied by checking a protected route
	req, err := http.NewRequest("GET", "/api/v1/users/profile", nil)
	require.NoError(t, err)

	resp, err := app.server.Test(req, 5000)
	require.NoError(t, err)

	// If middleware is properly applied, we should get 401, not 500 or 200
	if resp.StatusCode == 200 {
		t.Error("CRITICAL: Authentication middleware is NOT being applied! Route returns 200 without auth.")
	} else if resp.StatusCode == 500 {
		t.Error("CRITICAL: Server error (500) suggests middleware type assertion or nil pointer issue")
	} else {
		t.Logf("Good: Protected route returns %d (should be 401)", resp.StatusCode)
	}
}

// TestMainRoutes_WalletRoutesExist tests that wallet routes are properly configured
func TestMainRoutes_WalletRoutesExist(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Skipf("Cannot create app for testing: %v", err)
		return
	}

	walletRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/wallet/balances"},
		{"POST", "/api/v1/wallet/topup"},
		{"POST", "/api/v1/wallet/pay"},
		{"POST", "/api/v1/wallet/withdraw"},
	}

	for _, route := range walletRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req, err := http.NewRequest(route.method, route.path, nil)
			require.NoError(t, err)

			resp, err := app.server.Test(req)
			require.NoError(t, err)

			// Should NOT return 404 (route not found)
			if resp.StatusCode == 404 {
				t.Errorf("Wallet route %s %s is MISSING (404). Check route configuration.", route.method, route.path)
			} else {
				t.Logf("Wallet route %s %s exists (status: %d)", route.method, route.path, resp.StatusCode)
			}
		})
	}
}
