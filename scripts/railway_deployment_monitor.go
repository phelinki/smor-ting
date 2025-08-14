package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// HealthResponse represents the expected health check response
type HealthResponse struct {
	Status      string                 `json:"status"`
	Service     string                 `json:"service"`
	Version     string                 `json:"version"`
	Database    string                 `json:"database"`
	Environment string                 `json:"environment"`
	Timestamp   string                 `json:"timestamp"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// DeploymentMonitor monitors Railway deployment status
type DeploymentMonitor struct {
	RailwayURL string
	Timeout    time.Duration
}

// NewDeploymentMonitor creates a new deployment monitor
func NewDeploymentMonitor() *DeploymentMonitor {
	railwayURL := os.Getenv("RAILWAY_URL")
	if railwayURL == "" {
		railwayURL = "https://smor-ting-production.up.railway.app"
	}

	return &DeploymentMonitor{
		RailwayURL: railwayURL,
		Timeout:    30 * time.Second,
	}
}

// CheckHealth checks if the Railway deployment is healthy
func (dm *DeploymentMonitor) CheckHealth() (*HealthResponse, error) {
	client := &http.Client{
		Timeout: dm.Timeout,
	}

	resp, err := client.Get(dm.RailwayURL + "/health")
	if err != nil {
		return nil, fmt.Errorf("failed to reach health endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health endpoint returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var health HealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	return &health, nil
}

// CheckRailwayStatus checks the current Railway deployment status
func (dm *DeploymentMonitor) CheckRailwayStatus() (string, error) {
	cmd := exec.Command("railway", "status")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get railway status: %w", err)
	}

	return string(output), nil
}

// ForceRedeploy forces a Railway redeployment
func (dm *DeploymentMonitor) ForceRedeploy() error {
	log.Println("🚀 Forcing Railway redeployment...")
	
	cmd := exec.Command("railway", "redeploy", "--force")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to force redeploy: %w\nOutput: %s", err, string(output))
	}

	log.Printf("✅ Redeploy initiated: %s", string(output))
	return nil
}

// MonitorDeployment continuously monitors deployment health
func (dm *DeploymentMonitor) MonitorDeployment(maxRetries int, retryInterval time.Duration) error {
	log.Printf("🔍 Starting deployment monitoring for %s", dm.RailwayURL)
	
	for i := 0; i < maxRetries; i++ {
		log.Printf("📊 Health check attempt %d/%d", i+1, maxRetries)
		
		health, err := dm.CheckHealth()
		if err != nil {
			log.Printf("❌ Health check failed: %v", err)
			
			if i == 0 {
				// On first failure, check Railway status
				status, statusErr := dm.CheckRailwayStatus()
				if statusErr != nil {
					log.Printf("⚠️  Could not get Railway status: %v", statusErr)
				} else {
					log.Printf("📋 Railway status: %s", status)
				}
			}
			
			if i == maxRetries-1 {
				return fmt.Errorf("deployment health check failed after %d attempts: %w", maxRetries, err)
			}
			
			log.Printf("⏳ Waiting %v before retry...", retryInterval)
			time.Sleep(retryInterval)
			continue
		}
		
		log.Printf("✅ Deployment is healthy!")
		log.Printf("   Status: %s", health.Status)
		log.Printf("   Service: %s", health.Service)
		log.Printf("   Environment: %s", health.Environment)
		log.Printf("   Database: %s", health.Database)
		
		return nil
	}
	
	return fmt.Errorf("deployment monitoring failed after %d attempts", maxRetries)
}

// VerifyDeploymentTriggered checks if Railway detected file changes
func (dm *DeploymentMonitor) VerifyDeploymentTriggered() error {
	log.Println("🔍 Checking if Railway detected file changes...")
	
	cmd := exec.Command("railway", "logs")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get railway logs: %w", err)
	}
	
	logStr := string(output)
	
	// Check for deployment indicators
	if strings.Contains(logStr, "No changes to watched files") {
		return fmt.Errorf("❌ Railway is not detecting file changes - check watch patterns")
	}
	
	if strings.Contains(logStr, "Starting Container") {
		log.Println("✅ Railway detected changes and started deployment")
		return nil
	}
	
	log.Println("⚠️  Deployment status unclear from logs")
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	monitor := NewDeploymentMonitor()
	
	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "health":
			health, err := monitor.CheckHealth()
			if err != nil {
				log.Fatalf("❌ Health check failed: %v", err)
			}
			fmt.Printf("✅ Deployment healthy: %+v\n", health)
			
		case "status":
			status, err := monitor.CheckRailwayStatus()
			if err != nil {
				log.Fatalf("❌ Status check failed: %v", err)
			}
			fmt.Printf("📋 Railway status:\n%s", status)
			
		case "redeploy":
			if err := monitor.ForceRedeploy(); err != nil {
				log.Fatalf("❌ Redeploy failed: %v", err)
			}
			
		case "verify":
			if err := monitor.VerifyDeploymentTriggered(); err != nil {
				log.Fatalf("❌ Verification failed: %v", err)
			}
			
		case "monitor":
			maxRetries := 10
			retryInterval := 30 * time.Second
			
			if err := monitor.MonitorDeployment(maxRetries, retryInterval); err != nil {
				log.Fatalf("❌ Monitoring failed: %v", err)
			}
			
		default:
			fmt.Printf("Usage: %s [health|status|redeploy|verify|monitor]\n", os.Args[0])
			os.Exit(1)
		}
	} else {
		// Default: quick health check
		health, err := monitor.CheckHealth()
		if err != nil {
			log.Printf("❌ Health check failed: %v", err)
			log.Println("🔧 Attempting automatic fix...")
			
			if err := monitor.ForceRedeploy(); err != nil {
				log.Fatalf("❌ Auto-fix failed: %v", err)
			}
			
			log.Println("⏳ Waiting for redeploy to complete...")
			time.Sleep(60 * time.Second)
			
			if err := monitor.MonitorDeployment(5, 30*time.Second); err != nil {
				log.Fatalf("❌ Auto-fix monitoring failed: %v", err)
			}
		} else {
			log.Printf("✅ Deployment healthy: %s", health.Status)
		}
	}
}
