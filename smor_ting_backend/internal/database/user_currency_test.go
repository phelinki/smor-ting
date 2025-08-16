package database

import (
	"context"
	"testing"
	"time"

	"github.com/smorting/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserDefaultCurrency(t *testing.T) {
	// Test with both MongoDB and Memory databases
	// Test with Memory database
	db := NewMemoryDatabase()

	// Create a new user
	user := &models.User{
		Email:     "currency_test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+1234567890",
		Role:      models.CustomerRole,
	}

	// Create the user
	err := db.CreateUser(context.Background(), user)
	require.NoError(t, err)

	// Verify the user was created with USD currency
	assert.Equal(t, "USD", user.Wallet.Currency, "New users should default to USD currency")

	// Fetch the user from database to double-check
	fetchedUser, err := db.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	assert.Equal(t, "USD", fetchedUser.Wallet.Currency, "Fetched user should have USD currency")
}

func TestUserCreationWithCustomCurrency(t *testing.T) {
	db := NewMemoryDatabase()

	// Create a user with custom currency
	user := &models.User{
		Email:     "custom_currency@example.com",
		Password:  "password123",
		FirstName: "Custom",
		LastName:  "User",
		Phone:     "+1234567890",
		Role:      models.CustomerRole,
		Wallet: models.Wallet{
			Balance:     0,
			Currency:    "EUR", // Custom currency
			LastUpdated: time.Now(),
		},
	}

	// Create the user
	err := db.CreateUser(context.Background(), user)
	require.NoError(t, err)

	// Verify the custom currency is preserved
	assert.Equal(t, "EUR", user.Wallet.Currency, "Custom currency should be preserved")

	// Fetch the user from database to double-check
	fetchedUser, err := db.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	assert.Equal(t, "EUR", fetchedUser.Wallet.Currency, "Fetched user should preserve custom currency")
}
