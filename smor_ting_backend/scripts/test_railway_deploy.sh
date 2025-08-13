#!/bin/bash

# Test Railway Deployment Script
# This script tests the Railway deployment without actually deploying

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"
}

# Test Railway CLI availability
test_railway_cli() {
    log "Testing Railway CLI availability..."
    
    if ! command -v railway &> /dev/null; then
        error "Railway CLI not found. Please install it with: npm install -g @railway/cli"
    fi
    
    log "✅ Railway CLI is available"
    
    # Test Railway status (without requiring auth for this test)
    info "Testing Railway status command..."
    if railway status 2>/dev/null; then
        log "✅ Railway is authenticated and working"
    else
        warning "Railway CLI not authenticated or no project configured"
        info "To authenticate: railway login"
        info "To link project: railway link"
    fi
}

# Test configuration files
test_config_files() {
    log "Testing configuration files..."
    
    # Check railway.toml
    if [[ ! -f "railway.toml" ]]; then
        error "railway.toml not found in current directory"
    fi
    log "✅ railway.toml found"
    
    # Check Dockerfile
    if [[ ! -f "Dockerfile" ]]; then
        error "Dockerfile not found in current directory"
    fi
    log "✅ Dockerfile found"
    
    # Validate railway.toml content
    if ! grep -q "builder = \"DOCKERFILE\"" railway.toml; then
        error "railway.toml missing required DOCKERFILE builder configuration"
    fi
    
    if ! grep -q "healthcheckPath = \"/health\"" railway.toml; then
        error "railway.toml missing health check configuration"
    fi
    
    log "✅ Configuration files are valid"
}

# Test build process
test_build() {
    log "Testing build process..."
    
    # Check if we're in the right directory
    if [[ ! -f "cmd/main.go" ]]; then
        error "cmd/main.go not found. Please run from the backend directory"
    fi
    
    # Test Go build
    info "Testing Go build..."
    if go build -o test-build cmd/main.go; then
        log "✅ Go build successful"
        rm -f test-build
    else
        error "Go build failed"
    fi
}

# Test Docker build (optional)
test_docker_build() {
    log "Testing Docker build..."
    
    if command -v docker &> /dev/null; then
        info "Building Docker image (test)..."
        if docker build -t smor-ting-test:latest .; then
            log "✅ Docker build successful"
            # Cleanup test image
            docker rmi smor-ting-test:latest 2>/dev/null || true
        else
            error "Docker build failed"
        fi
    else
        warning "Docker not available, skipping Docker build test"
    fi
}

# Simulate deployment
simulate_deployment() {
    log "Simulating deployment process..."
    
    local target_env=${1:-staging}
    
    info "Target environment: $target_env"
    info "Deployment command would be: railway up --service smor-ting-$target_env"
    
    # Test railway login status
    if railway whoami 2>/dev/null; then
        info "✅ Railway authentication verified"
        
        # Test project link
        if railway status 2>/dev/null; then
            info "✅ Railway project linked"
            log "🚀 Deployment simulation successful"
        else
            warning "Railway project not linked, but CLI is authenticated"
            info "Deployment would work after project linking"
        fi
    else
        warning "Railway not authenticated"
        info "Deployment would require authentication first"
    fi
}

# Main test function
main() {
    log "🧪 Starting Railway deployment tests..."
    
    test_railway_cli
    test_config_files
    test_build
    test_docker_build
    simulate_deployment "${1:-staging}"
    
    log "🎉 All deployment tests completed successfully!"
    log "✅ Ready for Railway deployment"
}

# Run tests
main "$@"
