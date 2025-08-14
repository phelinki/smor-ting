package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRailwayFixImplementation tests that our Railway fix is properly implemented
func TestRailwayFixImplementation(t *testing.T) {
	t.Run("Should have only one railway.toml configuration", func(t *testing.T) {
		// Check that backend railway.toml is removed
		backendConfig := "../railway.toml"
		_, err := os.Stat(backendConfig)
		assert.Error(t, err, "Backend railway.toml should be removed to avoid conflicts")
		
		// Check that root railway.toml exists
		rootConfig := "../../railway.toml"
		_, err = os.Stat(rootConfig)
		assert.NoError(t, err, "Root railway.toml should exist")
	})

	t.Run("Should have comprehensive watch patterns", func(t *testing.T) {
		configPath := "../../railway.toml"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err, "Should be able to read root railway.toml")

		contentStr := string(content)
		
		// Should have multiple specific watch patterns
		requiredPatterns := []string{
			"smor_ting_backend/**/*",  // All files recursively
			"smor_ting_backend/*",     // Direct files in backend dir
			"smor_ting_backend/cmd/**/*", // All command files
			"smor_ting_backend/configs/**/*", // All config files
			"smor_ting_backend/Dockerfile", // Dockerfile specifically
			"smor_ting_backend/go.mod",     // Go module file
			"smor_ting_backend/go.sum",     // Go dependencies
		}

		for _, pattern := range requiredPatterns {
			assert.Contains(t, contentStr, pattern, "Should include watch pattern: %s", pattern)
		}
	})

	t.Run("Should have proper configuration structure", func(t *testing.T) {
		configPath := "../../railway.toml"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)

		contentStr := string(content)
		
		// Should use array format for watchPatterns
		assert.Contains(t, contentStr, "watchPatterns = [", "Should use array format for watch patterns")
		
		// Should have descriptive comments
		assert.Contains(t, contentStr, "Fixed watch patterns permanently", "Should document the fix")
		
		// Should still have proper build config
		assert.Contains(t, contentStr, "builder = \"DOCKERFILE\"", "Should use Dockerfile builder")
		assert.Contains(t, contentStr, "dockerfilePath = \"smor_ting_backend/Dockerfile\"", "Should reference correct Dockerfile")
	})
}

// TestDeploymentMonitoringTools tests that monitoring tools are available
func TestDeploymentMonitoringTools(t *testing.T) {
	t.Run("Should have deployment monitoring script", func(t *testing.T) {
		monitorPath := "../../scripts/railway_deployment_monitor.go"
		_, err := os.Stat(monitorPath)
		assert.NoError(t, err, "Should have deployment monitoring script")
		
		if err == nil {
			content, err := os.ReadFile(monitorPath)
			require.NoError(t, err)
			
			contentStr := string(content)
			
			// Should have essential monitoring functions
			assert.Contains(t, contentStr, "CheckHealth", "Should have health checking")
			assert.Contains(t, contentStr, "ForceRedeploy", "Should have force redeploy capability")
			assert.Contains(t, contentStr, "VerifyDeploymentTriggered", "Should verify deployment triggering")
		}
	})
}

// TestRailwayWatchPatternFunctionality tests that watch patterns will actually work
func TestRailwayWatchPatternFunctionality(t *testing.T) {
	t.Run("Should catch changes to critical backend files", func(t *testing.T) {
		configPath := "../../railway.toml"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)

		contentStr := string(content)
		
		// Test files that should trigger deployment
		testFiles := map[string]string{
			"../cmd/main.go":         "smor_ting_backend/cmd/**/*",
			"../configs/config.go":   "smor_ting_backend/configs/**/*",
			"../Dockerfile":          "smor_ting_backend/Dockerfile",
			"../go.mod":              "smor_ting_backend/go.mod",
			"../go.sum":              "smor_ting_backend/go.sum",
		}
		
		for file, expectedPattern := range testFiles {
			if _, err := os.Stat(file); err == nil {
				assert.Contains(t, contentStr, expectedPattern, 
					"Watch pattern for %s should be present: %s", file, expectedPattern)
				t.Logf("✅ File %s will be caught by pattern %s", file, expectedPattern)
			}
		}
	})

	t.Run("Should not have pattern conflicts", func(t *testing.T) {
		configPath := "../../railway.toml"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)

		contentStr := string(content)
		
		// Count watchPatterns sections - should only have one
		watchPatternSections := strings.Count(contentStr, "watchPatterns")
		assert.Equal(t, 1, watchPatternSections, "Should have exactly one watchPatterns section")
		
		// Should not have duplicate patterns
		lines := strings.Split(contentStr, "\n")
		var patterns []string
		inWatchSection := false
		
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "watchPatterns = [") {
				inWatchSection = true
				continue
			}
			if inWatchSection && line == "]" {
				break
			}
			if inWatchSection && strings.Contains(line, "smor_ting_backend") {
				patterns = append(patterns, line)
			}
		}
		
		// Check for duplicates
		seen := make(map[string]bool)
		for _, pattern := range patterns {
			assert.False(t, seen[pattern], "Pattern should not be duplicated: %s", pattern)
			seen[pattern] = true
		}
	})
}
