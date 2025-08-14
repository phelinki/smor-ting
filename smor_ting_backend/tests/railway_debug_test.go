package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRailwayConfigurationDebugging tests to debug why Railway isn't detecting changes
func TestRailwayConfigurationDebugging(t *testing.T) {
	t.Run("Should debug Railway configuration location", func(t *testing.T) {
		// Check which Railway config files exist
		configs := []string{
			"../../railway.toml",
			"../../railway.json", 
			"../../nixpacks.toml",
			"../railway.toml",
			"../railway.json",
		}
		
		var foundConfigs []string
		for _, config := range configs {
			if _, err := os.Stat(config); err == nil {
				foundConfigs = append(foundConfigs, config)
				t.Logf("Found config: %s", config)
				
				// Read and show content
				content, err := os.ReadFile(config)
				require.NoError(t, err)
				t.Logf("Content of %s:\n%s", config, string(content))
			}
		}
		
		assert.NotEmpty(t, foundConfigs, "Should have at least one Railway config")
		
		// Warn if multiple configs
		if len(foundConfigs) > 1 {
			t.Logf("WARNING: Multiple Railway configs found: %v", foundConfigs)
			t.Log("Railway may be using a different config than expected")
		}
	})

	t.Run("Should test Railway watch pattern syntax", func(t *testing.T) {
		configPath := "../../railway.toml"
		if _, err := os.Stat(configPath); err != nil {
			t.Skip("No root railway.toml found")
		}
		
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)
		
		contentStr := string(content)
		t.Logf("Railway.toml content:\n%s", contentStr)
		
		// Check if watchPatterns is in array format
		if strings.Contains(contentStr, "watchPatterns = [") {
			t.Log("✅ Using array format for watchPatterns")
		} else if strings.Contains(contentStr, "watchPatterns") {
			t.Log("❌ watchPatterns not in array format")
			assert.Fail(t, "watchPatterns must be in array format")
		} else {
			t.Log("❌ No watchPatterns found")
			assert.Fail(t, "No watchPatterns configured")
		}
		
		// Test each pattern individually
		patterns := []string{
			"smor_ting_backend/**/*",
			"smor_ting_backend/*",
			"smor_ting_backend/cmd/**/*",
			"smor_ting_backend/configs/**/*",
			"smor_ting_backend/Dockerfile",
			"smor_ting_backend/go.mod",
			"smor_ting_backend/go.sum",
		}
		
		for _, pattern := range patterns {
			if strings.Contains(contentStr, pattern) {
				t.Logf("✅ Found pattern: %s", pattern)
			} else {
				t.Logf("❌ Missing pattern: %s", pattern)
			}
		}
	})

	t.Run("Should check if Railway is using correct working directory", func(t *testing.T) {
		// Check current working directory
		wd, err := os.Getwd()
		require.NoError(t, err)
		t.Logf("Current working directory: %s", wd)
		
		// Check if Railway commands work from root
		rootDir := filepath.Join(wd, "../..")
		t.Logf("Root directory: %s", rootDir)
		
		// List files in backend directory to verify structure
		backendDir := filepath.Join(rootDir, "smor_ting_backend")
		entries, err := os.ReadDir(backendDir)
		require.NoError(t, err)
		
		t.Log("Files in smor_ting_backend/:")
		for _, entry := range entries {
			t.Logf("  %s", entry.Name())
		}
	})

	t.Run("Should test Railway project configuration", func(t *testing.T) {
		// Check Railway project info
		cmd := exec.Command("railway", "status")
		cmd.Dir = "../../" // Run from root directory
		output, err := cmd.Output()
		if err != nil {
			t.Logf("Could not get Railway status: %v", err)
			t.Skip("Railway not available or not logged in")
		}
		
		statusStr := string(output)
		t.Logf("Railway status:\n%s", statusStr)
		
		// Check if we're in the right project
		assert.Contains(t, statusStr, "smor-ting", "Should be in smor-ting project")
	})

	t.Run("Should verify file patterns would match actual files", func(t *testing.T) {
		// Get list of actual files in backend directory
		backendDir := "../"
		var actualFiles []string
		
		err := filepath.Walk(backendDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if !info.IsDir() {
				// Convert to relative path from root
				relPath := strings.TrimPrefix(path, "../")
				relPath = "smor_ting_backend/" + relPath
				actualFiles = append(actualFiles, relPath)
			}
			
			return nil
		})
		require.NoError(t, err)
		
		t.Logf("Found %d files in backend directory", len(actualFiles))
		
		// Test patterns against actual files
		patterns := []string{
			"smor_ting_backend/**/*",
			"smor_ting_backend/*",
			"smor_ting_backend/cmd/**/*",
			"smor_ting_backend/configs/**/*",
		}
		
		for _, pattern := range patterns {
			matchCount := 0
			for _, file := range actualFiles {
				if matchesPattern(pattern, file) {
					matchCount++
				}
			}
			t.Logf("Pattern '%s' would match %d files", pattern, matchCount)
			
			// Some patterns should match files
			if strings.Contains(pattern, "/**/*") {
				assert.Greater(t, matchCount, 0, "Pattern %s should match some files", pattern)
			}
		}
	})
}

// matchesPattern is a simple pattern matcher for testing
func matchesPattern(pattern, filename string) bool {
	// Simple wildcard matching for testing purposes
	if strings.HasSuffix(pattern, "/**/*") {
		prefix := strings.TrimSuffix(pattern, "/**/*")
		return strings.HasPrefix(filename, prefix+"/")
	}
	
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(filename, prefix+"/") && 
			   !strings.Contains(strings.TrimPrefix(filename, prefix+"/"), "/")
	}
	
	return filename == pattern
}

// TestRailwayAlternativeConfiguration tests alternative Railway configurations
func TestRailwayAlternativeConfiguration(t *testing.T) {
	t.Run("Should test if nixpacks.toml is interfering", func(t *testing.T) {
		nixpacksPath := "../../nixpacks.toml"
		if _, err := os.Stat(nixpacksPath); err == nil {
			t.Log("Found nixpacks.toml - this might override railway.toml")
			
			content, err := os.ReadFile(nixpacksPath)
			require.NoError(t, err)
			t.Logf("nixpacks.toml content:\n%s", string(content))
			
			// Check if it has watchPatterns
			contentStr := string(content)
			if strings.Contains(contentStr, "watchPatterns") {
				t.Log("❌ nixpacks.toml has watchPatterns - this might conflict")
			} else {
				t.Log("✅ nixpacks.toml doesn't have watchPatterns")
			}
		} else {
			t.Log("✅ No nixpacks.toml found")
		}
	})

	t.Run("Should check Railway service configuration", func(t *testing.T) {
		// Check if there are any .railwayignore files
		ignoreFiles := []string{
			"../../.railwayignore",
			"../.railwayignore",
		}
		
		for _, ignorePath := range ignoreFiles {
			if _, err := os.Stat(ignorePath); err == nil {
				content, err := os.ReadFile(ignorePath)
				require.NoError(t, err)
				t.Logf("Found .railwayignore at %s:\n%s", ignorePath, string(content))
			}
		}
	})
}
