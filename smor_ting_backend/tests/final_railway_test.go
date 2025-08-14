package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFinalRailwayFix verifies that all conflicting configurations are removed
func TestFinalRailwayFix(t *testing.T) {
	t.Run("Should have only ONE Railway configuration", func(t *testing.T) {
		// Check that conflicting configs are removed
		conflictingConfigs := []string{
			"../../nixpacks.toml",
			"../railway.json",
		}
		
		for _, config := range conflictingConfigs {
			_, err := os.Stat(config)
			assert.Error(t, err, "Conflicting config should be removed: %s", config)
		}
		
		// Check that the correct config exists
		correctConfig := "../../railway.toml"
		_, err := os.Stat(correctConfig)
		assert.NoError(t, err, "Correct Railway config should exist")
	})

	t.Run("Should have proper railwayignore that doesn't block critical files", func(t *testing.T) {
		railwayignorePath := "../.railwayignore"
		content, err := os.ReadFile(railwayignorePath)
		require.NoError(t, err)
		
		contentStr := string(content)
		
		// Should explicitly allow critical files
		criticalFiles := []string{
			"!go.mod",
			"!go.sum", 
			"!Dockerfile",
			"!cmd/",
			"!configs/",
			"!api/",
			"!internal/",
		}
		
		for _, file := range criticalFiles {
			assert.Contains(t, contentStr, file, "Should explicitly allow: %s", file)
		}
		
		// Should not ignore all scripts anymore
		assert.NotContains(t, contentStr, "scripts/\n!scripts/deploy_production.sh", 
			"Should not broadly ignore scripts directory")
	})

	t.Run("Should have comprehensive watch patterns in the only config", func(t *testing.T) {
		configPath := "../../railway.toml"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)
		
		contentStr := string(content)
		
		// Verify it's the authoritative config
		assert.Contains(t, contentStr, "ONLY Railway configuration file", 
			"Should mark as the only config")
		
		// Verify comprehensive patterns
		patterns := []string{
			"smor_ting_backend/**/*",
			"smor_ting_backend/cmd/**/*", 
			"smor_ting_backend/configs/**/*",
			"smor_ting_backend/Dockerfile",
			"smor_ting_backend/go.mod",
			"smor_ting_backend/go.sum",
		}
		
		for _, pattern := range patterns {
			assert.Contains(t, contentStr, pattern, "Should have pattern: %s", pattern)
		}
	})
}

// TestRailwayConfigurationIsAuthoritative verifies Railway will use our config
func TestRailwayConfigurationIsAuthoritative(t *testing.T) {
	t.Run("Should have no competing build configurations", func(t *testing.T) {
		// Count how many files could control Railway builds
		buildConfigs := []string{
			"../../railway.toml",      // This should exist
			"../../railway.json",      // This should NOT exist
			"../../nixpacks.toml",     // This should NOT exist  
			"../railway.toml",         // This should NOT exist
			"../railway.json",         // This should NOT exist
		}
		
		existingConfigs := 0
		var existingConfigsList []string
		
		for _, config := range buildConfigs {
			if _, err := os.Stat(config); err == nil {
				existingConfigs++
				existingConfigsList = append(existingConfigsList, config)
			}
		}
		
		assert.Equal(t, 1, existingConfigs, 
			"Should have exactly 1 Railway config, found: %v", existingConfigsList)
		
		// The one config should be the root railway.toml
		_, err := os.Stat("../../railway.toml")
		assert.NoError(t, err, "Root railway.toml should be the only config")
	})
}
