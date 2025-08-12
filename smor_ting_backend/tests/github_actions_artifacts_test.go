package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeprecatedUploadArtifact tests for deprecated upload-artifact@v3 usage
func TestDeprecatedUploadArtifact(t *testing.T) {
	workflowFiles := []string{
		"../../.github/workflows/qa-automation.yml",
		"../../.github/workflows/deployment-gate.yml",
	}
	
	for _, workflowFile := range workflowFiles {
		t.Run("Should not use deprecated upload-artifact@v3 in "+workflowFile, func(t *testing.T) {
			content, err := os.ReadFile(workflowFile)
			require.NoError(t, err, "Should be able to read workflow file")
			
			contentStr := string(content)
			
			// Check for deprecated v3 usage
			deprecatedPattern := "upload-artifact@v3"
			if strings.Contains(contentStr, deprecatedPattern) {
				t.Errorf("Found deprecated %s in %s", deprecatedPattern, workflowFile)
				
				// Show lines containing deprecated usage
				lines := strings.Split(contentStr, "\n")
				for i, line := range lines {
					if strings.Contains(line, deprecatedPattern) {
						t.Logf("Line %d: %s", i+1, strings.TrimSpace(line))
					}
				}
			}
			
			// Should use v4 instead
			assert.NotContains(t, contentStr, "upload-artifact@v3", 
				"Should not use deprecated upload-artifact@v3")
		})
	}
}

// TestDownloadArtifactVersion tests for consistent artifact action versions
func TestDownloadArtifactVersion(t *testing.T) {
	workflowFiles := []string{
		"../../.github/workflows/qa-automation.yml", 
		"../../.github/workflows/deployment-gate.yml",
	}
	
	for _, workflowFile := range workflowFiles {
		t.Run("Should use consistent artifact action versions in "+workflowFile, func(t *testing.T) {
			content, err := os.ReadFile(workflowFile)
			require.NoError(t, err, "Should be able to read workflow file")
			
			contentStr := string(content)
			
			// If using upload-artifact@v4, should also use download-artifact@v4
			hasUploadV4 := strings.Contains(contentStr, "upload-artifact@v4")
			hasDownloadV3 := strings.Contains(contentStr, "download-artifact@v3")
			
			if hasUploadV4 && hasDownloadV3 {
				t.Error("Inconsistent artifact action versions - using both v3 and v4")
			}
		})
	}
}

// TestBackendCIConfiguration tests Backend CI workflow configuration
func TestBackendCIConfiguration(t *testing.T) {
	t.Run("Should have valid Backend CI configuration", func(t *testing.T) {
		workflowPath := "../../.github/workflows/backend-ci.yml"
		content, err := os.ReadFile(workflowPath)
		require.NoError(t, err, "Should be able to read backend-ci.yml")
		
		contentStr := string(content)
		
		// Check for Go setup
		assert.Contains(t, contentStr, "setup-go@v5", "Should use recent Go setup action")
		assert.Contains(t, contentStr, "1.23", "Should use correct Go version")
		
		// Check for proper test environment
		assert.Contains(t, contentStr, "ENV: development", "Should set development environment for tests")
		
		// Check for vet and test steps
		assert.Contains(t, contentStr, "go vet", "Should run go vet")
		assert.Contains(t, contentStr, "go test", "Should run go test")
		assert.Contains(t, contentStr, "go build", "Should run go build")
	})
}

// TestSemgrepScanConfiguration tests Semgrep scan configuration 
func TestSemgrepScanConfiguration(t *testing.T) {
	t.Run("Should have proper Semgrep configuration", func(t *testing.T) {
		workflowPath := "../../.github/workflows/security-scan.yml"
		content, err := os.ReadFile(workflowPath)
		require.NoError(t, err, "Should be able to read security-scan.yml")
		
		contentStr := string(content)
		
		// Check for Semgrep action
		assert.Contains(t, contentStr, "semgrep/semgrep-action", "Should use Semgrep action")
		
		// Check for SARIF generation
		assert.Contains(t, contentStr, "generate_sarif: true", "Should generate SARIF output")
		
		// Check for continue-on-error
		assert.Contains(t, contentStr, "continue-on-error: true", "Should have continue-on-error for robustness")
		
		// Check for proper permissions
		assert.Contains(t, contentStr, "security-events: write", "Should have security-events write permission")
	})
}

// TestWorkflowVersionConsistency tests for consistent action versions across workflows
func TestWorkflowVersionConsistency(t *testing.T) {
	t.Run("Should use consistent action versions", func(t *testing.T) {
		workflowDir := "../../.github/workflows"
		files, err := os.ReadDir(workflowDir)
		require.NoError(t, err, "Should be able to read workflows directory")
		
		checkoutVersions := make(map[string][]string)
		setupGoVersions := make(map[string][]string)
		
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".yml") && !strings.HasSuffix(file.Name(), ".yaml") {
				continue
			}
			
			filePath := workflowDir + "/" + file.Name()
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			
			contentStr := string(content)
			lines := strings.Split(contentStr, "\n")
			
			for _, line := range lines {
				line = strings.TrimSpace(line)
				
				if strings.Contains(line, "actions/checkout@") {
					version := extractVersion(line, "actions/checkout@")
					if version != "" {
						checkoutVersions[version] = append(checkoutVersions[version], file.Name())
					}
				}
				
				if strings.Contains(line, "actions/setup-go@") {
					version := extractVersion(line, "actions/setup-go@")
					if version != "" {
						setupGoVersions[version] = append(setupGoVersions[version], file.Name())
					}
				}
			}
		}
		
		// Report inconsistencies
		if len(checkoutVersions) > 1 {
			t.Logf("Inconsistent checkout versions found: %v", checkoutVersions)
		}
		
		if len(setupGoVersions) > 1 {
			t.Logf("Inconsistent setup-go versions found: %v", setupGoVersions)
		}
		
		// Ensure we're using latest versions
		for version := range checkoutVersions {
			assert.Contains(t, []string{"v4"}, version, "Should use checkout@v4")
		}
		
		for version := range setupGoVersions {
			assert.Contains(t, []string{"v4", "v5"}, version, "Should use setup-go@v4 or v5")
		}
	})
}

// extractVersion extracts version from action usage line
func extractVersion(line, actionPrefix string) string {
	if !strings.Contains(line, actionPrefix) {
		return ""
	}
	
	start := strings.Index(line, actionPrefix) + len(actionPrefix)
	remaining := line[start:]
	
	// Find end of version (space, newline, or other delimiter)
	end := len(remaining)
	for i, char := range remaining {
		if char == ' ' || char == '\n' || char == '\r' || char == '\t' {
			end = i
			break
		}
	}
	
	return remaining[:end]
}
