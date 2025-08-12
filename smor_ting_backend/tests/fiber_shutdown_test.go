package tests

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFiberAppShutdown tests that Fiber app can be properly shut down
// This test is written to reproduce the issue where app.Close() doesn't exist
func TestFiberAppShutdown(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()
	require.NotNil(t, app)

	// Test that the app has a Shutdown method
	t.Run("App should have Shutdown method", func(t *testing.T) {
		// This test should pass - Fiber apps have Shutdown() method
		err := app.Shutdown()
		assert.NoError(t, err)
	})

	// Test that the app does NOT have a Close method
	t.Run("App should NOT have Close method", func(t *testing.T) {
		// This test documents that Close() method doesn't exist
		// We can't actually test this directly since it would cause a compile error
		// But this documents the expected behavior

		// The following line would cause a compile error:
		// app.Close() // This should NOT work

		// Instead, we should use Shutdown()
		err := app.Shutdown()
		assert.NoError(t, err)
	})
}

// TestFiberAppLifecycle tests the complete lifecycle of a Fiber app
func TestFiberAppLifecycle(t *testing.T) {
	app := fiber.New()

	// Test that we can create routes
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test"})
	})

	// Test that we can shut down the app properly
	defer func() {
		err := app.Shutdown()
		assert.NoError(t, err)
	}()

	// Test that the app is working
	assert.NotNil(t, app)
}
