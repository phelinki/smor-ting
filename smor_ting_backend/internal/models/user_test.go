package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestAuthResponse_JSON tests that AuthResponse properly marshals/unmarshals JSON
// This test ensures the requires_otp field is always included in the JSON response
func TestAuthResponse_JSON(t *testing.T) {
	t.Run("RequiresOTP true", func(t *testing.T) {
		user := User{
			ID:              primitive.NewObjectID(),
			Email:           "test@example.com",
			FirstName:       "Test",
			LastName:        "User",
			Phone:           "1234567890",
			Role:            CustomerRole,
			IsEmailVerified: false,
		}

		authResponse := AuthResponse{
			User:         user,
			AccessToken:  "test_access_token",
			RefreshToken: "test_refresh_token",
			RequiresOTP:  true,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(authResponse)
		require.NoError(t, err)

		// Verify requires_otp field is present
		assert.Contains(t, string(jsonData), `"requires_otp":true`)

		// Unmarshal back to struct
		var unmarshaledResponse AuthResponse
		err = json.Unmarshal(jsonData, &unmarshaledResponse)
		require.NoError(t, err)

		assert.Equal(t, authResponse.RequiresOTP, unmarshaledResponse.RequiresOTP)
		assert.True(t, unmarshaledResponse.RequiresOTP)
	})

	t.Run("RequiresOTP false", func(t *testing.T) {
		user := User{
			ID:              primitive.NewObjectID(),
			Email:           "test@example.com",
			FirstName:       "Test",
			LastName:        "User",
			Phone:           "1234567890",
			Role:            CustomerRole,
			IsEmailVerified: false,
		}

		authResponse := AuthResponse{
			User:         user,
			AccessToken:  "test_access_token",
			RefreshToken: "test_refresh_token",
			RequiresOTP:  false,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(authResponse)
		require.NoError(t, err)

		// Verify requires_otp field is present even when false
		assert.Contains(t, string(jsonData), `"requires_otp":false`)

		// Unmarshal back to struct
		var unmarshaledResponse AuthResponse
		err = json.Unmarshal(jsonData, &unmarshaledResponse)
		require.NoError(t, err)

		assert.Equal(t, authResponse.RequiresOTP, unmarshaledResponse.RequiresOTP)
		assert.False(t, unmarshaledResponse.RequiresOTP)
	})

	t.Run("JSON from Mobile App Format", func(t *testing.T) {
		// This tests the exact JSON format expected by the mobile app
		mobileAppJSON := `{
			"user": {
				"id": "6898ac7c3dc75dac76cd788c",
				"email": "test@example.com",
				"first_name": "Test",
				"last_name": "User", 
				"phone": "1234567890",
				"role": "customer",
				"is_email_verified": false,
				"profile_image": "",
				"address": {
					"street": "",
					"city": "",
					"county": "",
					"country": "",
					"latitude": 0,
					"longitude": 0
				},
				"wallet": {
					"balance": 0,
					"currency": "LRD",
					"last_updated": "2025-08-10T14:28:12.834256097Z"
				},
				"last_sync_at": "2025-08-10T14:28:12.83425601Z",
				"is_offline": false,
				"version": 1,
				"created_at": "2025-08-10T14:28:12.834255791Z",
				"updated_at": "2025-08-10T14:28:12.834255912Z"
			},
			"access_token": "test_access_token",
			"refresh_token": "test_refresh_token",
			"requires_otp": false
		}`

		// This should parse successfully
		var authResponse AuthResponse
		err := json.Unmarshal([]byte(mobileAppJSON), &authResponse)
		require.NoError(t, err)

		assert.Equal(t, "test@example.com", authResponse.User.Email)
		assert.Equal(t, "test_access_token", authResponse.AccessToken)
		assert.Equal(t, "test_refresh_token", authResponse.RefreshToken)
		assert.False(t, authResponse.RequiresOTP)
	})
}

// New test to ensure empty address is omitted (omitempty) and registration can work without address
func TestUser_AddressOmittedWhenEmpty(t *testing.T) {
	user := User{
		ID:        primitive.NewObjectID(),
		Email:     "noaddress@example.com",
		FirstName: "No",
		LastName:  "Address",
		Phone:     "",
		Role:      CustomerRole,
	}

	data, err := json.Marshal(user)
	require.NoError(t, err)
	// Address key should not be present when nil and omitempty is set
	assert.NotContains(t, string(data), "\"address\":")
}
