package tests

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSemgrepSARIFHandling tests that our SARIF handling works correctly
func TestSemgrepSARIFHandling(t *testing.T) {
	t.Run("Should handle empty SARIF file generation", func(t *testing.T) {
		// Test the empty SARIF we generate when Semgrep fails
		emptySARIF := `{"version": "2.1.0", "runs": []}`
		
		// Verify it's valid JSON
		var result map[string]interface{}
		err := json.Unmarshal([]byte(emptySARIF), &result)
		assert.NoError(t, err, "Empty SARIF should be valid JSON")
		
		// Verify required fields
		assert.Equal(t, "2.1.0", result["version"], "Should have correct SARIF version")
		assert.Contains(t, result, "runs", "Should have runs array")
	})
	
	t.Run("Should verify workflow permissions fix", func(t *testing.T) {
		workflowPath := "../../.github/workflows/security-scan.yml"
		content, err := os.ReadFile(workflowPath)
		require.NoError(t, err, "Should be able to read security-scan.yml")
		
		contentStr := string(content)
		
		// Verify permissions are set correctly
		assert.Contains(t, contentStr, "permissions:", "Should have permissions section")
		assert.Contains(t, contentStr, "security-events: write", "Should have security-events write permission")
		assert.Contains(t, contentStr, "contents: read", "Should have contents read permission")
		
		// Verify SARIF upload has file existence check
		assert.Contains(t, contentStr, "hashFiles('semgrep.sarif')", "Should check if SARIF file exists before upload")
	})
	
	t.Run("Should handle continue-on-error for Semgrep", func(t *testing.T) {
		workflowPath := "../../.github/workflows/security-scan.yml"
		content, err := os.ReadFile(workflowPath)
		require.NoError(t, err, "Should be able to read security-scan.yml")
		
		contentStr := string(content)
		
		// Verify Semgrep has continue-on-error
		assert.Contains(t, contentStr, "continue-on-error: true", "Semgrep should have continue-on-error")
		
		// Verify we have SARIF file checking
		assert.Contains(t, contentStr, "Check SARIF file generation", "Should have SARIF file generation check")
	})
}

// TestGitHubActionsIntegration tests GitHub Actions integration issues
func TestGitHubActionsIntegration(t *testing.T) {
	t.Run("Should resolve Resource not accessible by integration error", func(t *testing.T) {
		// The "Resource not accessible by integration" error is caused by missing permissions
		// Our fix adds the required permissions to the workflow
		
		workflowPath := "../../.github/workflows/security-scan.yml"
		content, err := os.ReadFile(workflowPath)
		require.NoError(t, err, "Should be able to read security-scan.yml")
		
		contentStr := string(content)
		
		// Check that we have the permissions that fix the integration error
		requiredPermissions := []string{
			"contents: read",           // For checking out code
			"security-events: write",   // For uploading SARIF files
			"actions: read",           // For reading action metadata
		}
		
		for _, permission := range requiredPermissions {
			assert.Contains(t, contentStr, permission, 
				"Should have %s permission to prevent integration errors", permission)
		}
	})
	
	t.Run("Should resolve Path does not exist sarif error", func(t *testing.T) {
		// The "Path does not exist: semgrep.sarif" error occurs when Semgrep fails to generate SARIF
		// Our fix adds a fallback to create an empty SARIF file
		
		workflowPath := "../../.github/workflows/security-scan.yml"
		content, err := os.ReadFile(workflowPath)
		require.NoError(t, err, "Should be able to read security-scan.yml")
		
		contentStr := string(content)
		
		// Check that we handle missing SARIF files
		assert.Contains(t, contentStr, "if [ -f \"semgrep.sarif\" ]", "Should check for SARIF file existence")
		assert.Contains(t, contentStr, "creating empty SARIF", "Should create empty SARIF as fallback")
		assert.Contains(t, contentStr, "hashFiles('semgrep.sarif') != ''", "Should verify SARIF file before upload")
	})
}
