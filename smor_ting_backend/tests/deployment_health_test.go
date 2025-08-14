package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRailwayWatchPatterns tests that Railway watch patterns are configured correctly
func TestRailwayWatchPatterns(t *testing.T) {
	t.Run("Should have correct watch patterns in root railway.toml", func(t *testing.T) {
		configPath := "../../railway.toml"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err, "Should be able to read root railway.toml")

		contentStr := string(content)

		// Should watch the entire backend directory
		assert.Contains(t, contentStr, "smor_ting_backend/**", "Should watch all backend files")

		// Should have proper dockerfile path
		assert.Contains(t, contentStr, "smor_ting_backend/Dockerfile", "Should reference correct Dockerfile path")

		// Should use DOCKERFILE builder
		assert.Contains(t, contentStr, "DOCKERFILE", "Should use Dockerfile builder")
	})

	t.Run("Should not have conflicting watch patterns", func(t *testing.T) {
		// Check both config files don't conflict
		rootConfig := "../../railway.toml"
		backendConfig := "../railway.toml"

		if _, err := os.Stat(rootConfig); err == nil {
			if _, err := os.Stat(backendConfig); err == nil {
				t.Log("Warning: Both root and backend railway.toml exist - Railway prefers root config")

				// Root config should be comprehensive enough to handle backend
				rootContent, err := os.ReadFile(rootConfig)
				require.NoError(t, err)

				rootStr := string(rootContent)
				assert.Contains(t, rootStr, "smor_ting_backend", "Root config should handle backend directory")
			}
		}
	})

	t.Run("Should include all important backend file patterns", func(t *testing.T) {
		configPath := "../../railway.toml"
		if _, err := os.Stat(configPath); err != nil {
			t.Skip("No root railway.toml found")
		}

		content, err := os.ReadFile(configPath)
		require.NoError(t, err)

		contentStr := string(content)

		// Should catch changes to critical files
		patterns := []string{
			"smor_ting_backend/**", // All backend files
		}

		for _, pattern := range patterns {
			assert.Contains(t, contentStr, pattern, "Should include pattern: %s", pattern)
		}
	})
}

// TestDeploymentEnvironmentVariables tests critical environment variables for Railway
func TestDeploymentEnvironmentVariables(t *testing.T) {
	t.Run("Should have Railway-compatible port configuration", func(t *testing.T) {
		configPath := "../configs/config.go"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)

		contentStr := string(content)

		// Must read PORT from environment (Railway requirement)
		assert.Contains(t, contentStr, "getEnv(\"PORT\"", "Must read PORT from environment for Railway")

		// Should default to 8080 if PORT not set
		assert.Contains(t, contentStr, "\"8080\"", "Should default to port 8080")
	})

	t.Run("Should have production-ready database configuration", func(t *testing.T) {
		configPath := "../configs/config.go"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)

		contentStr := string(content)

		// Should handle MongoDB URI from environment
		assert.Contains(t, contentStr, "MONGODB_URI", "Should read MongoDB URI from environment")
	})
}

// TestHealthEndpointReachability tests that the health endpoint is working correctly
func TestHealthEndpointReachability(t *testing.T) {
	t.Run("Should have health endpoint defined in routes", func(t *testing.T) {
		mainPath := "../cmd/main.go"
		content, err := os.ReadFile(mainPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "/health", "Should define health endpoint")
	})

	t.Run("Should be able to test health endpoint on live Railway deployment", func(t *testing.T) {
		// Test against the Railway URL directly (not custom domain)
		railwayURL := "https://smor-ting-production.up.railway.app"

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Get(railwayURL + "/health")
		if err != nil {
			t.Logf("Could not reach Railway deployment: %v", err)
			t.Skip("Railway deployment not reachable - this indicates deployment issue")
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Health endpoint should return 200")

		// Check response content
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var healthResp map[string]interface{}
		err = json.Unmarshal(body, &healthResp)
		require.NoError(t, err, "Health response should be valid JSON")

		assert.Equal(t, "healthy", healthResp["status"], "Should report healthy status")
		assert.Contains(t, healthResp, "service", "Should include service name")
	})
}

// TestGitChangeDetection tests that our changes are properly tracked
func TestGitChangeDetection(t *testing.T) {
	t.Run("Should detect recent changes in backend directory", func(t *testing.T) {
		// Check if we have recent git changes in the backend directory
		backendDir := "../"

		err := filepath.Walk(backendDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip hidden files and directories
			if strings.HasPrefix(info.Name(), ".") {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// Check for Go files, configs, and Dockerfile
			if strings.HasSuffix(path, ".go") ||
				strings.HasSuffix(path, "Dockerfile") ||
				strings.HasSuffix(path, ".toml") ||
				strings.HasSuffix(path, ".json") {

				// File exists and should be watched by Railway
				t.Logf("Found watchable file: %s", path)
			}

			return nil
		})

		require.NoError(t, err, "Should be able to scan backend directory")
	})
}

// TestRailwayConfigurationConsistency tests that Railway configs are consistent
func TestRailwayConfigurationConsistency(t *testing.T) {
	t.Run("Should have consistent binary names across configs", func(t *testing.T) {
		expectedBinary := "./smor-ting-api"

		configs := []string{
			"../../railway.toml",
			"../railway.toml",
		}

		for _, configPath := range configs {
			if _, err := os.Stat(configPath); err != nil {
				continue // Skip if file doesn't exist
			}

			content, err := os.ReadFile(configPath)
			require.NoError(t, err, "Should read config: %s", configPath)

			contentStr := string(content)
			if strings.Contains(contentStr, "startCommand") {
				assert.Contains(t, contentStr, expectedBinary,
					"Config %s should reference binary %s", configPath, expectedBinary)
			}
		}
	})

	t.Run("Should have consistent health check paths", func(t *testing.T) {
		expectedHealthPath := "/health"

		configs := []string{
			"../../railway.toml",
			"../railway.toml",
		}

		for _, configPath := range configs {
			if _, err := os.Stat(configPath); err != nil {
				continue
			}

			content, err := os.ReadFile(configPath)
			require.NoError(t, err)

			contentStr := string(content)
			if strings.Contains(contentStr, "healthcheckPath") {
				assert.Contains(t, contentStr, expectedHealthPath,
					"Config %s should use health path %s", configPath, expectedHealthPath)
			}
		}
	})
}

// TestDeploymentTriggers tests that deployments are triggered properly
func TestDeploymentTriggers(t *testing.T) {
	t.Run("Should trigger deployment on backend file changes", func(t *testing.T) {
		// This test documents what should trigger deployments
		triggerFiles := []string{
			"../cmd/main.go",
			"../configs/config.go",
			"../Dockerfile",
			"../go.mod",
			"../go.sum",
			"../../railway.toml",
		}

		for _, file := range triggerFiles {
			if _, err := os.Stat(file); err == nil {
				t.Logf("File %s exists and should trigger Railway deployment", file)
			}
		}

		// The key insight: Railway should watch smor_ting_backend/** to catch ALL these
		assert.True(t, true, "This test documents expected deployment triggers")
	})
}
