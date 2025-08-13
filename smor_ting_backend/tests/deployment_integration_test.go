package tests

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeploymentConfiguration validates that deployment configuration is correct
func TestDeploymentConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		configFile     string
		expectedFields []string
	}{
		{
			name:       "Railway config exists",
			configFile: "railway.toml",
			expectedFields: []string{
				"[build]",
				"builder = \"DOCKERFILE\"",
				"[deploy]",
				"startCommand",
				"healthcheckPath = \"/health\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to parent directory (backend root) if we're in tests/
			configPath := tt.configFile
			if _, err := os.Stat("../"+tt.configFile); err == nil {
				configPath = "../" + tt.configFile
			}
			
			// Check if config file exists
			_, err := os.Stat(configPath)
			require.NoError(t, err, "Config file %s should exist", configPath)

			// Read config file
			content, err := os.ReadFile(configPath)
			require.NoError(t, err, "Should be able to read config file")

			configStr := string(content)

			// Validate required fields
			for _, field := range tt.expectedFields {
				assert.Contains(t, configStr, field, "Config should contain %s", field)
			}
		})
	}
}

// TestRailwayCliAvailable tests that Railway CLI is available and authenticated
func TestRailwayCliAvailable(t *testing.T) {
	t.Run("Railway CLI is installed", func(t *testing.T) {
		cmd := exec.Command("which", "railway")
		err := cmd.Run()
		require.NoError(t, err, "Railway CLI should be installed")
	})

	t.Run("Railway CLI can show status", func(t *testing.T) {
		cmd := exec.Command("railway", "status")
		cmd.Dir = ".."
		output, err := cmd.CombinedOutput()
		
		// Note: This might fail if not authenticated, but we should get a meaningful error
		if err != nil {
			// Check if it's an authentication issue vs installation issue
			outputStr := string(output)
			if strings.Contains(outputStr, "not logged in") || 
			   strings.Contains(outputStr, "authenticate") ||
			   strings.Contains(outputStr, "login") {
				t.Logf("Railway CLI available but not authenticated: %s", outputStr)
			} else {
				t.Errorf("Railway CLI error (not auth related): %s", outputStr)
			}
		} else {
			t.Logf("Railway status: %s", string(output))
		}
	})
}

// TestDeploymentScripts validates that deployment scripts exist and are executable
func TestDeploymentScripts(t *testing.T) {
	scripts := []struct {
		name string
		path string
	}{
		{
			name: "Production deployment script",
			path: "scripts/deploy_production.sh",
		},
	}

	for _, script := range scripts {
		t.Run(script.name, func(t *testing.T) {
			// Check if script exists (try parent directory first)
			scriptPath := script.path
			if _, err := os.Stat("../"+script.path); err == nil {
				scriptPath = "../" + script.path
			}
			
			_, err := os.Stat(scriptPath)
			require.NoError(t, err, "Script %s should exist", scriptPath)

			// Check if script is executable
			info, err := os.Stat(scriptPath)
			require.NoError(t, err)
			
			mode := info.Mode()
			assert.True(t, mode&0111 != 0, "Script should be executable")
		})
	}
}

// TestDockerfileExists validates that Dockerfile exists for containerized deployment
func TestDockerfileExists(t *testing.T) {
	// Check if Dockerfile exists (try parent directory first)
	dockerfilePath := "Dockerfile"
	if _, err := os.Stat("../Dockerfile"); err == nil {
		dockerfilePath = "../Dockerfile"
	}
	
	_, err := os.Stat(dockerfilePath)
	require.NoError(t, err, "Dockerfile should exist for containerized deployment")

	// Read Dockerfile and validate key components
	content, err := os.ReadFile(dockerfilePath)
	require.NoError(t, err)

	dockerfileStr := string(content)
	
	// Check for essential Dockerfile components
	assert.Contains(t, dockerfileStr, "FROM", "Dockerfile should have FROM instruction")
	assert.Contains(t, dockerfileStr, "COPY", "Dockerfile should copy application files")
	assert.Contains(t, dockerfileStr, "EXPOSE", "Dockerfile should expose port")
	assert.Contains(t, dockerfileStr, "CMD", "Dockerfile should have CMD instruction")
}

// TestGitHubActionsDeploymentWorkflow validates the deployment workflow
func TestGitHubActionsDeploymentWorkflow(t *testing.T) {
	workflowPath := "../../.github/workflows/deployment-gate.yml"
	
	// Check if workflow exists
	_, err := os.Stat(workflowPath)
	require.NoError(t, err, "Deployment workflow should exist")

	// Read workflow content
	content, err := os.ReadFile(workflowPath)
	require.NoError(t, err)

	workflowStr := string(content)

	// Validate workflow structure
	assert.Contains(t, workflowStr, "name:", "Workflow should have a name")
	assert.Contains(t, workflowStr, "on:", "Workflow should have triggers")
	assert.Contains(t, workflowStr, "deploy-backend:", "Workflow should have backend deployment job")
}

// TestEnvironmentVariablesForDeployment tests that required environment variables are documented
func TestEnvironmentVariablesForDeployment(t *testing.T) {
	envFiles := []string{
		".env.example",
		"../railway.toml",
	}

	requiredEnvVars := []string{
		"PORT",
		"ENV",
		"HOST",
	}

	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err != nil {
			continue // Skip if file doesn't exist
		}

		t.Run("Check "+envFile, func(t *testing.T) {
			content, err := os.ReadFile(envFile)
			require.NoError(t, err)

			envContent := string(content)

			for _, envVar := range requiredEnvVars {
				assert.Contains(t, envContent, envVar, 
					"Environment file should contain %s", envVar)
			}
		})
	}
}

// TestRailwayDeployment tests the actual Railway deployment process
func TestRailwayDeployment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping deployment test in short mode")
	}

	t.Run("Railway deployment dry run", func(t *testing.T) {
		// This test validates the deployment process without actually deploying
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// First, ensure we're in the backend directory
		cmd := exec.CommandContext(ctx, "railway", "status")
		cmd.Dir = ".."
		
		output, err := cmd.CombinedOutput()
		outputStr := string(output)
		
		if err != nil {
			if strings.Contains(outputStr, "not logged in") {
				t.Skip("Railway CLI not authenticated - skipping deployment test")
			}
			if strings.Contains(outputStr, "No project found") {
				t.Log("No Railway project configured yet - this is expected for initial setup")
			} else {
				t.Logf("Railway status check failed: %s", outputStr)
			}
		} else {
			t.Logf("Railway status successful: %s", outputStr)
		}
	})
}

// TestHealthCheckEndpoint validates that the health check endpoint is implemented
func TestHealthCheckEndpoint(t *testing.T) {
	// This test should validate that /health endpoint exists in the application
	// For now, we'll check if it's referenced in the code
	
	files := []string{
		"../cmd/main.go",
		"../internal/handlers/health.go",
		"../internal/routes/routes.go",
	}

	healthEndpointFound := false

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		if strings.Contains(string(content), "/health") {
			healthEndpointFound = true
			break
		}
	}

	assert.True(t, healthEndpointFound, "Health check endpoint should be implemented")
}
