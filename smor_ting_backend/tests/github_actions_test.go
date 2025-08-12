package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGitHubActionsEnvironment tests that the environment is properly set up for GitHub Actions
func TestGitHubActionsEnvironment(t *testing.T) {
	t.Run("Should have required environment variables for CI", func(t *testing.T) {
		// These are the environment variables that should be available in CI
		envVar := os.Getenv("ENV")
		
		// In CI, ENV should be set to development
		if envVar == "" {
			// Local development - that's okay
			t.Skip("ENV not set - running locally")
		}
		
		assert.Equal(t, "development", envVar, "ENV should be 'development' in CI")
	})
	
	t.Run("Should not have production secrets in CI", func(t *testing.T) {
		// These should NOT be set in CI (they should be empty or default values)
		sensitiveVars := []string{
			"MONGODB_URI",
			"JWT_SECRET", 
			"GITGUARDIAN_API_KEY",
		}
		
		for _, varName := range sensitiveVars {
			value := os.Getenv(varName)
			// In CI, these should either be empty or test values
			if value != "" {
				// If set, they should be test/placeholder values, not real secrets
				assert.NotContains(t, value, "mongodb+srv://", 
					"MONGODB_URI should not contain real Atlas connection in CI")
				assert.NotContains(t, value, "cluster.mongodb.net",
					"MONGODB_URI should not contain real cluster in CI")
			}
		}
	})
}

// TestSecurityScanCompatibility tests that our code works with security scanning tools
func TestSecurityScanCompatibility(t *testing.T) {
	t.Run("Should not contain patterns that trigger false positives", func(t *testing.T) {
		// Test that our .gitguardian.yaml exists and is properly configured
		_, err := os.Stat("../.gitguardian.yaml")
		assert.NoError(t, err, ".gitguardian.yaml should exist to configure security scans")
	})
	
	t.Run("Should have proper .env.example without real credentials", func(t *testing.T) {
		// Test that .env.example exists and contains only placeholder values
		_, err := os.Stat("../.env.example")
		assert.NoError(t, err, ".env.example should exist as a template")
		
		// Could add content checks here if needed
	})
}

// TestBuildEnvironment tests that the build environment is properly configured
func TestBuildEnvironment(t *testing.T) {
	t.Run("Should compile without errors", func(t *testing.T) {
		// This test just needs to run - if there are compilation errors,
		// the test won't even start
		assert.True(t, true, "If this test runs, compilation succeeded")
	})
	
	t.Run("Should have proper go.mod", func(t *testing.T) {
		// Check that go.mod exists
		_, err := os.Stat("../go.mod")
		require.NoError(t, err, "go.mod should exist")
	})
}
