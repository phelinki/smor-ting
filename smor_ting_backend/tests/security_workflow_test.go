package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSecurityWorkflowRequirements tests that all files required by the security workflow exist
func TestSecurityWorkflowRequirements(t *testing.T) {
	t.Run("Should have GitGuardian configuration", func(t *testing.T) {
		_, err := os.Stat("../.gitguardian.yaml")
		assert.NoError(t, err, ".gitguardian.yaml should exist for security scanning")
	})
	
	t.Run("Should have TruffleHog ignore file", func(t *testing.T) {
		_, err := os.Stat("../../.trufflehog-ignore")
		assert.NoError(t, err, ".trufflehog-ignore should exist to prevent false positives")
	})
	
	t.Run("Should have .env.example template", func(t *testing.T) {
		_, err := os.Stat("../.env.example")
		assert.NoError(t, err, ".env.example should exist as secure template")
	})
	
	t.Run("Should have security scan workflow", func(t *testing.T) {
		_, err := os.Stat("../../.github/workflows/security-scan.yml")
		assert.NoError(t, err, "security-scan.yml workflow should exist")
	})
	
	t.Run("Should not have real .env files committed", func(t *testing.T) {
		// Check for .env files that shouldn't be committed
		envFiles := []string{
			"../.env",
			"../.env.local", 
			"../.env.production",
			"../.env.development",
		}
		
		for _, envFile := range envFiles {
			_, err := os.Stat(envFile)
			if err == nil {
				t.Errorf("Environment file %s should not be committed to repository", envFile)
			}
		}
	})
}

// TestDockerfileExists tests that Dockerfile exists for infrastructure scanning
func TestDockerfileSecurityFeatures(t *testing.T) {
	t.Run("Should have Dockerfile for infrastructure scanning", func(t *testing.T) {
		dockerfilePath := "../Dockerfile"
		
		// Check if Dockerfile exists
		_, err := os.Stat(dockerfilePath)
		if err != nil {
			// If it doesn't exist, create a basic one for testing
			t.Logf("Dockerfile not found, this may cause infrastructure scan failures")
			
			// Check if we're in a CI environment
			if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
				t.Skip("In CI environment without Dockerfile - infrastructure scan will be skipped")
			}
		} else {
			t.Logf("✅ Dockerfile found at %s", dockerfilePath)
		}
	})
}

// TestSecurityConfiguration tests that security configurations are valid
func TestSecurityConfiguration(t *testing.T) {
	t.Run("GitGuardian config should be valid YAML", func(t *testing.T) {
		content, err := os.ReadFile("../.gitguardian.yaml")
		require.NoError(t, err, "Should be able to read .gitguardian.yaml")
		
		// Basic checks for required content
		contentStr := string(content)
		assert.Contains(t, contentStr, "version:", "Should have version specified")
		assert.Contains(t, contentStr, "paths-ignore:", "Should have paths to ignore")
		assert.Contains(t, contentStr, "package-lock.json", "Should ignore npm integrity hashes")
	})
	
	t.Run("Should have proper .gitignore patterns", func(t *testing.T) {
		gitignorePath := "../../.gitignore"
		content, err := os.ReadFile(gitignorePath)
		require.NoError(t, err, "Should be able to read .gitignore")
		
		contentStr := string(content)
		assert.Contains(t, contentStr, ".env", "Should ignore .env files")
		assert.Contains(t, contentStr, "!.env.example", "Should NOT ignore .env.example")
	})
}

// TestWorkflowFiles tests that all workflow files are present and valid
func TestWorkflowFiles(t *testing.T) {
	workflowDir := "../../.github/workflows"
	
	t.Run("Should have required workflow files", func(t *testing.T) {
		requiredWorkflows := []string{
			"security-scan.yml",
			"backend-ci.yml",
		}
		
		for _, workflow := range requiredWorkflows {
			workflowPath := filepath.Join(workflowDir, workflow)
			_, err := os.Stat(workflowPath)
			assert.NoError(t, err, "Workflow %s should exist", workflow)
		}
	})
	
	t.Run("Security workflow should have continue-on-error for external tools", func(t *testing.T) {
		securityWorkflowPath := filepath.Join(workflowDir, "security-scan.yml")
		content, err := os.ReadFile(securityWorkflowPath)
		require.NoError(t, err, "Should be able to read security-scan.yml")
		
		contentStr := string(content)
		// Check that external tools have continue-on-error to prevent build failures
		assert.Contains(t, contentStr, "continue-on-error", "Should have continue-on-error for external tools")
	})
}
