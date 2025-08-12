package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthenticationIntegration tests all authentication endpoints comprehensively
func TestAuthenticationIntegration(t *testing.T) {
	// Setup test environment
	config := &configs.Config{
		Database: configs.DatabaseConfig{
			InMemory: true,
			Driver:   "mongodb",
		},
	}

	// TODO: Fix integration test setup
	// app, err := cmd.NewApp(config)
	app := fiber.New()
	err := error(nil)
	require.NoError(t, err)
	defer app.Close()

	// Test registration scenarios
	t.Run("Registration Tests", func(t *testing.T) {
		testRegistrationScenarios(t, app.GetFiberApp())
	})

	// Test login scenarios
	t.Run("Login Tests", func(t *testing.T) {
		testLoginScenarios(t, app.GetFiberApp())
	})

	// Test token validation scenarios
	t.Run("Token Validation Tests", func(t *testing.T) {
		testTokenValidationScenarios(t, app.GetFiberApp())
	})
}

func testRegistrationScenarios(t *testing.T, app *fiber.App) {
	validUser := models.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "231777123456",
		Role:      "customer",
	}

	testCases := []struct {
		name            string
		request         models.RegisterRequest
		expectedStatus  int
		expectedError   string
		expectedMessage string
	}{
		{
			name:           "Valid Registration",
			request:        validUser,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing Email",
			request: models.RegisterRequest{
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "231777123456",
				Role:      "customer",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "email is required",
		},
		{
			name: "Missing Password",
			request: models.RegisterRequest{
				Email:     "test2@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "231777123456",
				Role:      "customer",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "password is required",
		},
		{
			name: "Password Too Short",
			request: models.RegisterRequest{
				Email:     "test3@example.com",
				Password:  "12345",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "231777123456",
				Role:      "customer",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "password must be at least 6 characters long",
		},
		{
			name: "Missing First Name",
			request: models.RegisterRequest{
				Email:    "test4@example.com",
				Password: "password123",
				LastName: "Doe",
				Phone:    "231777123456",
				Role:     "customer",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "first name is required",
		},
		{
			name: "Missing Last Name",
			request: models.RegisterRequest{
				Email:     "test5@example.com",
				Password:  "password123",
				FirstName: "John",
				Phone:     "231777123456",
				Role:      "customer",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "last name is required",
		},
		{
			name: "Missing Phone",
			request: models.RegisterRequest{
				Email:     "test6@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "customer",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "phone is required",
		},
		{
			name: "Missing Role",
			request: models.RegisterRequest{
				Email:     "test7@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "231777123456",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "role is required",
		},
		{
			name: "Invalid Role",
			request: models.RegisterRequest{
				Email:     "test8@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "231777123456",
				Role:      "invalid_role",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "role must be 'customer', 'provider', or 'admin'",
		},
		{
			name:            "Email Already Exists",
			request:         validUser, // Same as first test
			expectedStatus:  http.StatusConflict,
			expectedError:   "User already exists",
			expectedMessage: "A user with this email already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus != http.StatusCreated {
				var errorResp struct {
					Error   string `json:"error"`
					Message string `json:"message"`
				}
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedError, errorResp.Error)
				assert.Equal(t, tc.expectedMessage, errorResp.Message)
			} else {
				var authResp models.AuthResponse
				err = json.NewDecoder(resp.Body).Decode(&authResp)
				require.NoError(t, err)

				assert.Equal(t, tc.request.Email, authResp.User.Email)
				assert.Equal(t, tc.request.FirstName, authResp.User.FirstName)
				assert.Equal(t, tc.request.LastName, authResp.User.LastName)
				assert.NotEmpty(t, authResp.AccessToken)
				assert.False(t, authResp.RequiresOTP)
			}
		})
	}
}

func testLoginScenarios(t *testing.T, app *fiber.App) {
	// First register a user for login tests
	registerUser := models.RegisterRequest{
		Email:     "login_test@example.com",
		Password:  "password123",
		FirstName: "Login",
		LastName:  "Test",
		Phone:     "231777123456",
		Role:      "customer",
	}

	body, _ := json.Marshal(registerUser)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	resp.Body.Close()

	testCases := []struct {
		name            string
		email           string
		password        string
		expectedStatus  int
		expectedError   string
		expectedMessage string
	}{
		{
			name:           "Valid Login",
			email:          "login_test@example.com",
			password:       "password123",
			expectedStatus: http.StatusOK,
		},
		{
			name:            "Missing Email",
			email:           "",
			password:        "password123",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "email is required",
		},
		{
			name:            "Missing Password",
			email:           "login_test@example.com",
			password:        "",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Validation failed",
			expectedMessage: "password is required",
		},
		{
			name:            "Non-existent Email",
			email:           "nonexistent@example.com",
			password:        "password123",
			expectedStatus:  http.StatusUnauthorized,
			expectedError:   "Invalid credentials",
			expectedMessage: "Invalid email or password",
		},
		{
			name:            "Wrong Password",
			email:           "login_test@example.com",
			password:        "wrongpassword",
			expectedStatus:  http.StatusUnauthorized,
			expectedError:   "Invalid credentials",
			expectedMessage: "Invalid email or password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loginReq := models.LoginRequest{
				Email:    tc.email,
				Password: tc.password,
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus != http.StatusOK {
				var errorResp struct {
					Error   string `json:"error"`
					Message string `json:"message"`
				}
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedError, errorResp.Error)
				assert.Equal(t, tc.expectedMessage, errorResp.Message)
			} else {
				var authResp models.AuthResponse
				err = json.NewDecoder(resp.Body).Decode(&authResp)
				require.NoError(t, err)

				assert.Equal(t, tc.email, authResp.User.Email)
				assert.NotEmpty(t, authResp.AccessToken)
				assert.False(t, authResp.RequiresOTP)
			}
		})
	}
}

func testTokenValidationScenarios(t *testing.T, app *fiber.App) {
	// First register and login to get a valid token
	registerUser := models.RegisterRequest{
		Email:     "token_test@example.com",
		Password:  "password123",
		FirstName: "Token",
		LastName:  "Test",
		Phone:     "231777123456",
		Role:      "customer",
	}

	body, _ := json.Marshal(registerUser)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var authResp models.AuthResponse
	json.NewDecoder(resp.Body).Decode(&authResp)
	resp.Body.Close()

	validToken := authResp.AccessToken

	testCases := []struct {
		name            string
		token           string
		expectedStatus  int
		expectedError   string
		expectedMessage string
		malformedJSON   bool
	}{
		{
			name:           "Valid Token",
			token:          validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:            "Invalid Token",
			token:           "invalid.token.here",
			expectedStatus:  http.StatusUnauthorized,
			expectedError:   "Invalid token",
			expectedMessage: "The provided token is invalid or expired",
		},
		{
			name:            "Empty Token",
			token:           "",
			expectedStatus:  http.StatusUnauthorized,
			expectedError:   "Invalid token",
			expectedMessage: "The provided token is invalid or expired",
		},
		{
			name:            "Malformed Request",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "Invalid request body",
			expectedMessage: "Failed to parse request body",
			malformedJSON:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request

			if tc.malformedJSON {
				req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/validate", bytes.NewReader([]byte("invalid json")))
			} else {
				tokenReq := struct {
					Token string `json:"token"`
				}{
					Token: tc.token,
				}
				body, _ := json.Marshal(tokenReq)
				req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/validate", bytes.NewReader(body))
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus != http.StatusOK {
				var errorResp struct {
					Error   string `json:"error"`
					Message string `json:"message"`
				}
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedError, errorResp.Error)
				assert.Equal(t, tc.expectedMessage, errorResp.Message)
			} else {
				var validationResp struct {
					Valid bool        `json:"valid"`
					User  models.User `json:"user"`
				}
				err = json.NewDecoder(resp.Body).Decode(&validationResp)
				require.NoError(t, err)

				assert.True(t, validationResp.Valid)
				assert.Equal(t, "token_test@example.com", validationResp.User.Email)
			}
		})
	}
}

// TestMalformedRequestBodies tests various malformed request scenarios
func TestMalformedRequestBodies(t *testing.T) {
	config := &configs.Config{
		Database: configs.DatabaseConfig{
			InMemory: true,
			Driver:   "mongodb",
		},
	}

	// TODO: Fix integration test setup
	// app, err := cmd.NewApp(config)
	app := fiber.New()
	err := error(nil)
	require.NoError(t, err)
	defer app.Close()

	endpoints := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/validate",
	}

	malformedBodies := []struct {
		name string
		body string
	}{
		{"Empty Body", ""},
		{"Invalid JSON", "{invalid json}"},
		{"Incomplete JSON", `{"email": "test@example.com"`},
		{"Wrong Content Type", `{"email": "test@example.com", "password": "test123"}`},
	}

	for _, endpoint := range endpoints {
		for _, test := range malformedBodies {
			t.Run(fmt.Sprintf("%s - %s", endpoint, test.name), func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader([]byte(test.body)))

				if test.name != "Wrong Content Type" {
					req.Header.Set("Content-Type", "application/json")
				}

				resp, err := app.Test(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var errorResp struct {
					Error   string `json:"error"`
					Message string `json:"message"`
				}
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				require.NoError(t, err)

				assert.Equal(t, "Invalid request body", errorResp.Error)
				assert.Equal(t, "Failed to parse request body", errorResp.Message)
			})
		}
	}
}

// BenchmarkAuthEndpoints benchmarks authentication endpoints performance
func BenchmarkAuthEndpoints(b *testing.B) {
	config := &configs.Config{
		Database: configs.DatabaseConfig{
			InMemory: true,
			Driver:   "mongodb",
		},
	}

	// TODO: Fix integration test setup
	// app, err := cmd.NewApp(config)
	app := fiber.New()
	err := error(nil)
	require.NoError(b, err)
	defer app.Close()

	b.Run("Registration", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := models.RegisterRequest{
				Email:     fmt.Sprintf("bench%d@example.com", i),
				Password:  "password123",
				FirstName: "Bench",
				LastName:  "Test",
				Phone:     "231777123456",
				Role:      "customer",
			}

			body, _ := json.Marshal(user)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)
			resp.Body.Close()
		}
	})

	// Register a user for login benchmark
	user := models.RegisterRequest{
		Email:     "benchlogin@example.com",
		Password:  "password123",
		FirstName: "Bench",
		LastName:  "Login",
		Phone:     "231777123456",
		Role:      "customer",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	resp.Body.Close()

	b.Run("Login", func(b *testing.B) {
		loginReq := models.LoginRequest{
			Email:    "benchlogin@example.com",
			Password: "password123",
		}

		for i := 0; i < b.N; i++ {
			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)
			resp.Body.Close()
		}
	})
}
