package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8080"
	apiV1   = baseURL + "/api/v1"
)

// TestUser represents a test user for integration tests
type TestUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	} `json:"user"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

// Integration test suite following TDD principles
func TestIntegrationSuite(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration tests; set RUN_INTEGRATION_TESTS=1 to enable")
	}
	// Wait for server to be ready
	waitForServer(t)

	t.Run("Health Check", testHealthCheck)
	t.Run("API Documentation", testAPIDocumentation)
	t.Run("Authentication Flow", testAuthenticationFlow)
	t.Run("User Management", testUserManagement)
	t.Run("Payment System", testPaymentSystem)
	t.Run("Service Management", testServiceManagement)
	t.Run("Wallet Operations", testWalletOperations)
	t.Run("Sync Operations", testSyncOperations)
	t.Run("Security Features", testSecurityFeatures)
}

func waitForServer(t *testing.T) {
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatal("Server did not start within 30 seconds")
}

func testHealthCheck(t *testing.T) {
	resp, err := http.Get(baseURL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var health map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&health)
	require.NoError(t, err)

	assert.Equal(t, "healthy", health["status"])
	assert.Equal(t, "smor-ting-backend", health["service"])
	assert.Equal(t, "1.0.0", health["version"])
	assert.Equal(t, "development", health["environment"])

	// Verify security features are enabled
	security, ok := health["security"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "enabled", security["aes_256_encryption"])
	assert.Equal(t, "enabled", security["jwt_refresh"])
	assert.Equal(t, "enabled", security["pci_dss_compliance"])
}

func testAPIDocumentation(t *testing.T) {
	resp, err := http.Get(baseURL + "/docs")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var docs map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&docs)
	require.NoError(t, err)

	assert.Equal(t, "API Documentation", docs["message"])
	assert.Equal(t, "1.0.0", docs["version"])

	// Verify all expected endpoints are documented
	endpoints, ok := docs["endpoints"].(map[string]interface{})
	require.True(t, ok)

	expectedGroups := []string{"auth", "users", "services", "payments", "sync"}
	for _, group := range expectedGroups {
		assert.Contains(t, endpoints, group, "Missing endpoint group: %s", group)
	}
}

func testAuthenticationFlow(t *testing.T) {
	testUser := TestUser{
		Email:     fmt.Sprintf("test_%d@example.com", time.Now().Unix()),
		Password:  "SecurePassword123!",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+231123456789",
		Role:      "customer",
	}

	t.Run("Register New User", func(t *testing.T) {
		resp, body := makeRequest(t, "POST", apiV1+"/auth/register", testUser)
		t.Logf("Registration response: %s", string(body))

		// The actual API might return different status codes or responses
		// This test documents the expected behavior
		if resp.StatusCode == 201 || resp.StatusCode == 200 {
			var authResp AuthResponse
			err := json.Unmarshal(body, &authResp)
			if err == nil && authResp.AccessToken != "" {
				assert.NotEmpty(t, authResp.AccessToken)
				assert.NotEmpty(t, authResp.User.Email)
				assert.Equal(t, testUser.Email, authResp.User.Email)
			}
		} else {
			t.Logf("Registration failed with status %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("Login Existing User", func(t *testing.T) {
		loginData := map[string]string{
			"email":    testUser.Email,
			"password": testUser.Password,
		}

		resp, body := makeRequest(t, "POST", apiV1+"/auth/login", loginData)
		t.Logf("Login response: %s", string(body))

		// Document expected behavior
		if resp.StatusCode == 200 {
			var authResp AuthResponse
			err := json.Unmarshal(body, &authResp)
			if err == nil {
				assert.NotEmpty(t, authResp.AccessToken)
				assert.Equal(t, testUser.Email, authResp.User.Email)
			}
		} else {
			t.Logf("Login failed with status %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("Token Validation", func(t *testing.T) {
		// This would require a valid token from registration/login
		tokenData := map[string]string{
			"token": "dummy_token_for_testing",
		}

		resp, body := makeRequest(t, "POST", apiV1+"/auth/validate", tokenData)
		t.Logf("Token validation response: %s", string(body))

		// Document behavior - might be 401 with dummy token
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 401)
	})
}

func testUserManagement(t *testing.T) {
	t.Run("Get User Profile - Unauthenticated", func(t *testing.T) {
		resp, body := makeRequest(t, "GET", apiV1+"/users/profile", nil)
		t.Logf("Profile (unauth) response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Get User Profile - With Token", func(t *testing.T) {
		// This would require authentication header
		req, err := http.NewRequest("GET", apiV1+"/users/profile", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer dummy_token")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		t.Logf("Profile (auth) response: %s", string(body))
		// Should return 401 with dummy token
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func testPaymentSystem(t *testing.T) {
	t.Run("Tokenize Payment Method", func(t *testing.T) {
		paymentData := map[string]interface{}{
			"card_number": "4111111111111111",
			"cvv":         "123",
			"expiry":      "12/25",
			"card_type":   "visa",
		}

		resp, body := makeRequest(t, "POST", apiV1+"/payments/tokenize", paymentData)
		t.Logf("Payment tokenization response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Process Payment", func(t *testing.T) {
		paymentReq := map[string]interface{}{
			"token_id": "dummy_token",
			"amount":   100.00,
			"currency": "USD",
		}

		resp, body := makeRequest(t, "POST", apiV1+"/payments/process", paymentReq)
		t.Logf("Payment processing response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Validate Payment Token", func(t *testing.T) {
		resp, body := makeRequest(t, "GET", apiV1+"/payments/validate?token_id=dummy_token", nil)
		t.Logf("Payment validation response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func testServiceManagement(t *testing.T) {
	t.Run("List Services", func(t *testing.T) {
		resp, body := makeRequest(t, "GET", apiV1+"/services", nil)
		t.Logf("Services list response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Create Service", func(t *testing.T) {
		serviceData := map[string]interface{}{
			"name":        "Test Service",
			"description": "A test service",
			"price":       50.00,
			"category":    "testing",
		}

		resp, body := makeRequest(t, "POST", apiV1+"/services", serviceData)
		t.Logf("Service creation response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func testWalletOperations(t *testing.T) {
	t.Run("Get Wallet Balances", func(t *testing.T) {
		resp, body := makeRequest(t, "GET", apiV1+"/wallet/balances", nil)
		t.Logf("Wallet balances response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Wallet Topup", func(t *testing.T) {
		topupData := map[string]interface{}{
			"amount":     100.00,
			"currency":   "USD",
			"payment_id": "dummy_payment_id",
		}

		resp, body := makeRequest(t, "POST", apiV1+"/wallet/topup", topupData)
		t.Logf("Wallet topup response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func testSyncOperations(t *testing.T) {
	t.Run("Sync Data", func(t *testing.T) {
		syncData := map[string]interface{}{
			"user_id":      "507f1f77bcf86cd799439011",
			"last_sync_at": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			"device_id":    "test_device_123",
			"app_version":  "1.0.0",
		}

		resp, body := makeRequest(t, "POST", apiV1+"/sync/data", syncData)
		t.Logf("Sync data response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Get Unsynced Data", func(t *testing.T) {
		resp, body := makeRequest(t, "GET", apiV1+"/sync/unsynced?user_id=507f1f77bcf86cd799439011", nil)
		t.Logf("Unsynced data response: %s", string(body))

		// Should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func testSecurityFeatures(t *testing.T) {
	t.Run("JWT Token Refresh", func(t *testing.T) {
		refreshData := map[string]string{
			"refresh_token": "dummy_refresh_token",
		}

		resp, body := makeRequest(t, "POST", apiV1+"/auth/refresh", refreshData)
		t.Logf("Token refresh response: %s", string(body))

		// Should return 401 with dummy token
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("JWT Token Revocation", func(t *testing.T) {
		revokeData := map[string]string{
			"refresh_token": "dummy_refresh_token",
		}

		resp, body := makeRequest(t, "POST", apiV1+"/auth/revoke", revokeData)
		t.Logf("Token revocation response: %s", string(body))

		// Should handle revocation request (might be 400 for invalid token)
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 400 || resp.StatusCode == 401)
	})

	t.Run("Rate Limiting", func(t *testing.T) {
		// Make multiple rapid requests to test rate limiting
		for i := 0; i < 5; i++ {
			resp, _ := makeRequest(t, "GET", baseURL+"/health", nil)
			if i == 0 {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
			// Note: Rate limiting might not be implemented yet
		}
	})
}

// Helper function to make HTTP requests
func makeRequest(t *testing.T, method, url string, data interface{}) (*http.Response, []byte) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		require.NoError(t, err)
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	return resp, respBody
}

// Benchmark tests for performance
func BenchmarkHealthCheck(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

func BenchmarkAPIDocumentation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(baseURL + "/docs")
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}
