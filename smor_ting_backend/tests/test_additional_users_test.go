package tests

import (
	"context"
	"testing"

	"github.com/smorting/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateAdditionalTestUsers creates the additional test users requested by the user
func TestCreateAdditionalTestUsers(t *testing.T) {
	// Setup test environment
	ctx := context.Background()
	authService := setupAuthService(t)

	additionalUsers := []models.RegisterRequest{
		{
			Email:     "libworker@smorting.com",
			Password:  "Smorting8&",
			FirstName: "Agent",
			LastName:  "User",
			Phone:     "231555123456",
			Role:      models.ProviderRole, // Agent = Provider role
		},
		{
			Email:     "libhubby@smorting.com",
			Password:  "Smorting8&",
			FirstName: "Customer",
			LastName:  "User",
			Phone:     "231666123456",
			Role:      models.CustomerRole,
		},
		{
			Email:     "boss@smorting.com",
			Password:  "23$%ting&kukujumuku",
			FirstName: "Admin",
			LastName:  "User",
			Phone:     "231999888777",
			Role:      models.AdminRole,
		},
	}

	for _, testUser := range additionalUsers {
		t.Run("Create "+string(testUser.Role)+" user: "+testUser.Email, func(t *testing.T) {
			// Test: User should be created successfully
			response, err := authService.Register(ctx, &testUser)

			// Assertions following TDD principles
			require.NoError(t, err, "Should create test user without error")
			require.NotNil(t, response, "Should return auth response")
			require.NotEmpty(t, response.AccessToken, "Should return access token")

			assert.Equal(t, testUser.Email, response.User.Email, "Should have correct email")
			assert.Equal(t, testUser.FirstName, response.User.FirstName, "Should have correct first name")
			assert.Equal(t, testUser.LastName, response.User.LastName, "Should have correct last name")
			assert.Equal(t, testUser.Phone, response.User.Phone, "Should have correct phone")
			assert.Equal(t, testUser.Role, response.User.Role, "Should have correct role")
			assert.False(t, response.User.IsEmailVerified, "Email should not be verified initially")
			assert.NotEmpty(t, response.User.ID, "Should have user ID")
			assert.NotZero(t, response.User.CreatedAt, "Should have creation timestamp")

			// Test: User should be able to login with created credentials
			loginReq := &models.LoginRequest{
				Email:    testUser.Email,
				Password: testUser.Password,
			}
			loginResponse, err := authService.Login(ctx, loginReq)
			require.NoError(t, err, "Should be able to login with created credentials")
			require.NotNil(t, loginResponse, "Should return login response")
			assert.Equal(t, response.User.ID, loginResponse.User.ID, "Login should return same user")
		})
	}

	// Test: Verify all additional test users can login
	t.Run("Verify all additional test users can login", func(t *testing.T) {
		testCredentials := []struct {
			email    string
			password string
			role     string
		}{
			{"libworker@smorting.com", "Smorting8&", "provider"}, // Agent = Provider
			{"libhubby@smorting.com", "Smorting8&", "customer"},
			{"boss@smorting.com", "23$%ting&kukujumuku", "admin"},
		}

		for _, cred := range testCredentials {
			loginReq := &models.LoginRequest{
				Email:    cred.email,
				Password: cred.password,
			}
			response, err := authService.Login(ctx, loginReq)
			require.NoError(t, err, "Test user %s should be able to login", cred.email)
			assert.Equal(t, cred.email, response.User.Email, "User email should match")
			assert.Equal(t, cred.role, string(response.User.Role), "User role should match")
		}
	})
}

// setupAuthService is defined in test_user_creation_test.go
