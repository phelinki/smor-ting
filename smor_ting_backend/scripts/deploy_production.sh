#!/bin/bash

# Production Deployment Script for Smor-Ting API
# Usage: ./scripts/deploy_production.sh [version]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="smor-ting-api"
SERVICE_NAME="smor-ting-api"
APP_DIR="/opt/smor-ting"
BACKUP_DIR="/opt/smor-ting/backups"
LOG_FILE="/var/log/smor-ting/deploy.log"
HEALTH_CHECK_URL="http://localhost:8080/health"

# Version (default to timestamp if not provided)
VERSION=${1:-$(date +%Y%m%d_%H%M%S)}

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}" | tee -a "$LOG_FILE"
    exit 1
}

warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}" | tee -a "$LOG_FILE"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}" | tee -a "$LOG_FILE"
}

# Check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        error "This script should not be run as root"
    fi
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if service exists
    if ! systemctl list-unit-files | grep -q "$SERVICE_NAME"; then
        error "Service $SERVICE_NAME not found. Please install the service first."
    fi
    
    # Check if app directory exists
    if [[ ! -d "$APP_DIR" ]]; then
        error "Application directory $APP_DIR not found"
    fi
    
    # Check if backup directory exists
    if [[ ! -d "$BACKUP_DIR" ]]; then
        log "Creating backup directory..."
        sudo mkdir -p "$BACKUP_DIR"
        sudo chown $USER:$USER "$BACKUP_DIR"
    fi
    
    # Check if log directory exists
    if [[ ! -d "/var/log/smor-ting" ]]; then
        log "Creating log directory..."
        sudo mkdir -p "/var/log/smor-ting"
        sudo chown $USER:$USER "/var/log/smor-ting"
    fi
}

# Health check function
health_check() {
    local max_attempts=30
    local attempt=1
    
    log "Performing health check..."
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -f -s "$HEALTH_CHECK_URL" > /dev/null; then
            log "Health check passed!"
            return 0
        else
            warning "Health check attempt $attempt/$max_attempts failed"
            sleep 2
            ((attempt++))
        fi
    done
    
    error "Health check failed after $max_attempts attempts"
}

# Backup current version
backup_current() {
    log "Creating backup of current version..."
    
    local backup_file="$BACKUP_DIR/${APP_NAME}_$(date +%Y%m%d_%H%M%S).backup"
    
    if [[ -f "$APP_DIR/$APP_NAME" ]]; then
        sudo cp "$APP_DIR/$APP_NAME" "$backup_file"
        log "Backup created: $backup_file"
    else
        warning "No existing binary found to backup"
    fi
}

# Build new version
build_new_version() {
    log "Building new version..."
    
    # Check if we're in the right directory
    if [[ ! -f "cmd/main.go" ]]; then
        error "Please run this script from the project root directory"
    fi
    
    # Build the application
    log "Building for Linux..."
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$APP_NAME" cmd/main.go
    
    if [[ ! -f "$APP_NAME" ]]; then
        error "Build failed - binary not created"
    fi
    
    log "Build completed successfully"
}

# Deploy new version
deploy_new_version() {
    log "Deploying new version..."
    
    # Stop the service
    log "Stopping service..."
    sudo systemctl stop "$SERVICE_NAME" || warning "Failed to stop service (might not be running)"
    
    # Wait a moment for graceful shutdown
    sleep 3
    
    # Copy new binary
    log "Installing new binary..."
    sudo cp "$APP_NAME" "$APP_DIR/$APP_NAME"
    sudo chown smor-ting:smor-ting "$APP_DIR/$APP_NAME"
    sudo chmod +x "$APP_DIR/$APP_NAME"
    
    # Start the service
    log "Starting service..."
    sudo systemctl start "$SERVICE_NAME"
    
    # Wait for service to start
    sleep 5
    
    # Check service status
    if ! sudo systemctl is-active --quiet "$SERVICE_NAME"; then
        error "Service failed to start"
    fi
    
    log "Service started successfully"
}

# Rollback function
rollback() {
    local backup_file="$1"
    
    error "Rolling back to previous version..."
    
    # Stop service
    sudo systemctl stop "$SERVICE_NAME"
    
    # Restore backup
    if [[ -f "$backup_file" ]]; then
        sudo cp "$backup_file" "$APP_DIR/$APP_NAME"
        sudo chown smor-ting:smor-ting "$APP_DIR/$APP_NAME"
        sudo chmod +x "$APP_DIR/$APP_NAME"
        log "Backup restored: $backup_file"
    else
        error "Backup file not found: $backup_file"
    fi
    
    # Start service
    sudo systemctl start "$SERVICE_NAME"
    
    # Health check
    health_check
    
    log "Rollback completed successfully"
}

# Cleanup old backups
cleanup_backups() {
    log "Cleaning up old backups (keeping last 5)..."
    
    # Keep only the last 5 backups
    cd "$BACKUP_DIR"
    ls -t ${APP_NAME}_*.backup 2>/dev/null | tail -n +6 | xargs -r rm -f
    
    log "Cleanup completed"
}

# Main deployment function
main() {
    log "Starting production deployment - Version: $VERSION"
    
    # Check if running as root
    check_root
    
    # Check prerequisites
    check_prerequisites
    
    # Create backup
    backup_current
    
    # Build new version
    build_new_version
    
    # Deploy new version
    deploy_new_version
    
    # Health check
    health_check
    
    # Cleanup
    cleanup_backups
    
    log "Deployment completed successfully!"
    log "Version: $VERSION"
    log "Health check URL: $HEALTH_CHECK_URL"
    
    # Show service status
    sudo systemctl status "$SERVICE_NAME" --no-pager
}

# Handle errors
trap 'error "Deployment failed at line $LINENO"' ERR

# Run main function
main "$@" 