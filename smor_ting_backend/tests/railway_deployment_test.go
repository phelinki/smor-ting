package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRailwayDeploymentConfiguration tests Railway deployment setup
func TestRailwayDeploymentConfiguration(t *testing.T) {
	t.Run("Should have proper Railway configuration", func(t *testing.T) {
		// Check for Railway config files
		configs := []string{
			"../railway.toml",
			"../railway.json",
			"../../nixpacks.toml",
		}

		foundConfig := false
		for _, config := range configs {
			if _, err := os.Stat(config); err == nil {
				foundConfig = true
				t.Logf("Found Railway config: %s", config)
				break
			}
		}

		assert.True(t, foundConfig, "Should have Railway configuration file")
	})

	t.Run("Should have working directory configured for monorepo", func(t *testing.T) {
		nixpacksPath := "../../nixpacks.toml"
		if _, err := os.Stat(nixpacksPath); err != nil {
			t.Skip("nixpacks.toml not found - skipping monorepo test")
		}

		content, err := os.ReadFile(nixpacksPath)
		require.NoError(t, err, "Should be able to read nixpacks.toml")

		contentStr := string(content)
		assert.Contains(t, contentStr, "smor_ting_backend", "Should specify backend directory")
		assert.Contains(t, contentStr, "workingDirectory", "Should have working directory set")
	})

	t.Run("Should have proper build configuration", func(t *testing.T) {
		// Check build command in nixpacks.toml
		nixpacksPath := "../../nixpacks.toml"
		if _, err := os.Stat(nixpacksPath); err == nil {
			content, err := os.ReadFile(nixpacksPath)
			require.NoError(t, err, "Should be able to read nixpacks.toml")

			contentStr := string(content)
			assert.Contains(t, contentStr, "go build", "Should have Go build command")
			assert.Contains(t, contentStr, "./cmd", "Should build from cmd directory")
		}
	})
}

// TestDockerfileConfiguration tests Docker deployment setup
func TestDockerfileConfiguration(t *testing.T) {
	t.Run("Should have valid Dockerfile", func(t *testing.T) {
		dockerfilePath := "../Dockerfile"
		content, err := os.ReadFile(dockerfilePath)
		require.NoError(t, err, "Should be able to read Dockerfile")

		contentStr := string(content)

		// Check for proper base image
		assert.Contains(t, contentStr, "FROM golang:", "Should use Go base image")

		// Check for proper working directory
		assert.Contains(t, contentStr, "WORKDIR", "Should set working directory")

		// Check for build commands
		assert.Contains(t, contentStr, "go build", "Should have Go build command")

		// Check for expose port
		assert.Contains(t, contentStr, "EXPOSE", "Should expose port")

		// Check for entrypoint/cmd
		hasEntrypoint := assert.Contains(t, contentStr, "ENTRYPOINT", "Should have entrypoint") ||
			assert.Contains(t, contentStr, "CMD", "Should have cmd")
		assert.True(t, hasEntrypoint, "Should have either ENTRYPOINT or CMD")
	})
}

// TestRailwayEnvironmentVariables tests environment variable setup
func TestRailwayEnvironmentVariables(t *testing.T) {
	t.Run("Should handle Railway PORT environment variable", func(t *testing.T) {
		// Railway automatically sets PORT environment variable
		// Our app should read from PORT env var, not hardcode 8080

		// Check if we have proper port handling in config
		configPath := "../configs/config.go"
		content, err := os.ReadFile(configPath)
		require.NoError(t, err, "Should be able to read config.go")

		contentStr := string(content)
		assert.Contains(t, contentStr, "getEnv(\"PORT\"", "Should read PORT from environment")
	})

	t.Run("Should have environment template for Railway", func(t *testing.T) {
		// Should have .env.example for Railway variable setup
		envExamplePath := "../.env.example"
		_, err := os.Stat(envExamplePath)
		assert.NoError(t, err, ".env.example should exist for Railway setup")
	})
}

// TestRailwayHealthCheck tests health check configuration
func TestRailwayHealthCheck(t *testing.T) {
	t.Run("Should have health check endpoint", func(t *testing.T) {
		// Check railway config for health check
		configs := []string{
			"../railway.toml",
			"../railway.json",
		}

		for _, configPath := range configs {
			if _, err := os.Stat(configPath); err == nil {
				content, err := os.ReadFile(configPath)
				require.NoError(t, err, "Should be able to read Railway config")

				contentStr := string(content)
				assert.Contains(t, contentStr, "/health", "Should have health check path configured")
				break
			}
		}
	})

	t.Run("Should have reasonable health check timeout", func(t *testing.T) {
		railwayJsonPath := "../railway.json"
		if _, err := os.Stat(railwayJsonPath); err == nil {
			content, err := os.ReadFile(railwayJsonPath)
			require.NoError(t, err, "Should be able to read railway.json")

			contentStr := string(content)
			// Should have healthcheck timeout (300 seconds is reasonable)
			assert.Contains(t, contentStr, "healthcheckTimeout", "Should have health check timeout")
		}
	})
}

// TestBuildIssues tests common build issues that cause deployment failures
func TestBuildIssues(t *testing.T) {
	t.Run("Should not have conflicting build configurations", func(t *testing.T) {
		// Check if we have both Dockerfile and nixpacks.toml which can conflict
		hasDockerfile := false
		hasNixpacks := false

		if _, err := os.Stat("../Dockerfile"); err == nil {
			hasDockerfile = true
		}

		if _, err := os.Stat("../../nixpacks.toml"); err == nil {
			hasNixpacks = true
		}

		if hasDockerfile && hasNixpacks {
			t.Log("Warning: Both Dockerfile and nixpacks.toml found - Railway will prefer Dockerfile")
			t.Log("Consider removing one to avoid confusion")
		}

		assert.True(t, hasDockerfile || hasNixpacks, "Should have either Dockerfile or nixpacks.toml")
	})

	t.Run("Should have proper Go module setup", func(t *testing.T) {
		// Check go.mod exists
		goModPath := "../go.mod"
		_, err := os.Stat(goModPath)
		assert.NoError(t, err, "go.mod should exist for proper Go builds")

		if err == nil {
			content, err := os.ReadFile(goModPath)
			require.NoError(t, err, "Should be able to read go.mod")

			contentStr := string(content)
			assert.Contains(t, contentStr, "module", "go.mod should have module declaration")
			assert.Contains(t, contentStr, "go ", "go.mod should specify Go version")
		}
	})
}
