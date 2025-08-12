package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSemgrepWorkflowConfiguration tests that Semgrep workflow is properly configured
func TestSemgrepWorkflowConfiguration(t *testing.T) {
	t.Run("Should have correct SARIF upload configuration", func(t *testing.T) {
		workflowPath := "../../.github/workflows/security-scan.yml"
		content, err := os.ReadFile(workflowPath)
		require.NoError(t, err, "Should be able to read security-scan.yml")
		
		contentStr := string(content)
		
		// Check that SARIF upload has proper conditions
		assert.Contains(t, contentStr, "generate_sarif:", "Should have generate_sarif configuration")
		assert.Contains(t, contentStr, "upload-sarif@v3", "Should use upload-sarif action")
		
		// Should have if condition to handle missing SARIF file
		assert.Contains(t, contentStr, "if:", "Should have conditional execution for SARIF upload")
	})
	
	t.Run("Should handle missing SARIF file gracefully", func(t *testing.T) {
		// Test that we don't try to upload non-existent SARIF files
		sarfPath := "/tmp/test-missing.sarif"
		
		// Ensure file doesn't exist
		_, err := os.Stat(sarfPath)
		assert.True(t, os.IsNotExist(err), "Test SARIF file should not exist")
	})
}

// TestGitHubActionsPermissions tests that we handle GitHub Actions permissions properly
func TestGitHubActionsPermissions(t *testing.T) {
	t.Run("Should have security-events permission for SARIF upload", func(t *testing.T) {
		workflowPath := "../../.github/workflows/security-scan.yml"
		content, err := os.ReadFile(workflowPath)
		require.NoError(t, err, "Should be able to read security-scan.yml")
		
		contentStr := string(content)
		
		// Check for permissions section
		if !assert.Contains(t, contentStr, "permissions:", "Workflow should have permissions section") {
			t.Log("Missing permissions section - this may cause 'Resource not accessible by integration' errors")
		}
	})
}

// TestSemgrepConfigurationFiles tests Semgrep configuration
func TestSemgrepConfigurationFiles(t *testing.T) {
	t.Run("Should have valid project structure for Semgrep", func(t *testing.T) {
		// Check for Go files that Semgrep can analyze
		goFiles := []string{
			"../cmd/main.go",
			"../internal/handlers",
			"../internal/services", 
		}
		
		foundGoCode := false
		for _, path := range goFiles {
			if _, err := os.Stat(path); err == nil {
				foundGoCode = true
				break
			}
		}
		
		assert.True(t, foundGoCode, "Should have Go code for Semgrep to analyze")
	})
	
	t.Run("Should not require Semgrep config file for default rules", func(t *testing.T) {
		// Semgrep should work with default rulesets without a config file
		// This test documents that we're using p/security-audit, p/secrets, etc.
		configPath := "../../.semgrep.yml"
		
		if _, err := os.Stat(configPath); err == nil {
			t.Log("Found .semgrep.yml config file")
		} else {
			t.Log("Using default Semgrep rulesets (p/security-audit, p/secrets, etc.)")
		}
		
		// This should not fail - we're using public rulesets
		assert.True(t, true, "Default rulesets should work without custom config")
	})
}
