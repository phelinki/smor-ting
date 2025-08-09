#!/bin/bash

# Production Monitoring Script for Smor-Ting API
# This script monitors the health of your production API and sends alerts

set -e

# Configuration
APP_NAME="smor-ting-api"
SERVICE_NAME="smor-ting-api"
HEALTH_URL="http://localhost:8080/health"
API_URL="https://api.yourdomain.com"
LOG_FILE="/var/log/smor-ting/monitor.log"
ALERT_LOG="/var/log/smor-ting/alerts.log"
DISK_THRESHOLD=80
MEMORY_THRESHOLD=80
CPU_THRESHOLD=80

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}" | tee -a "$LOG_FILE"
    echo "$(date +'%Y-%m-%d %H:%M:%S') - ERROR: $1" >> "$ALERT_LOG"
}

warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}" | tee -a "$LOG_FILE"
    echo "$(date +'%Y-%m-%d %H:%M:%S') - WARNING: $1" >> "$ALERT_LOG"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}" | tee -a "$LOG_FILE"
}

# Check if service is running
check_service() {
    if ! systemctl is-active --quiet "$SERVICE_NAME"; then
        error "Service $SERVICE_NAME is down"
        systemctl restart "$SERVICE_NAME"
        sleep 5
        if ! systemctl is-active --quiet "$SERVICE_NAME"; then
            error "Failed to restart service $SERVICE_NAME"
            return 1
        else
            info "Service $SERVICE_NAME restarted successfully"
        fi
    else
        info "Service $SERVICE_NAME is running"
    fi
}

# Check health endpoint
check_health() {
    local max_attempts=3
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -f -s "$HEALTH_URL" > /dev/null; then
            info "Health check passed"
            return 0
        else
            warning "Health check attempt $attempt/$max_attempts failed"
            sleep 2
            ((attempt++))
        fi
    done
    
    error "Health check failed after $max_attempts attempts"
    return 1
}

# Check external API endpoint
check_external_api() {
    if curl -f -s "$API_URL/health" > /dev/null; then
        info "External API health check passed"
    else
        error "External API health check failed"
        return 1
    fi
}

# Check disk usage
check_disk() {
    local disk_usage=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
    
    if [ "$disk_usage" -gt "$DISK_THRESHOLD" ]; then
        warning "Disk usage is high: ${disk_usage}%"
        return 1
    else
        info "Disk usage: ${disk_usage}%"
    fi
}

# Check memory usage
check_memory() {
    local mem_usage=$(free | awk 'NR==2{printf "%.0f", $3*100/$2}')
    
    if [ "$mem_usage" -gt "$MEMORY_THRESHOLD" ]; then
        warning "Memory usage is high: ${mem_usage}%"
        return 1
    else
        info "Memory usage: ${mem_usage}%"
    fi
}

# Check CPU usage
check_cpu() {
    local cpu_usage=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | awk -F'%' '{print $1}')
    
    if [ "$cpu_usage" -gt "$CPU_THRESHOLD" ]; then
        warning "CPU usage is high: ${cpu_usage}%"
        return 1
    else
        info "CPU usage: ${cpu_usage}%"
    fi
}

# Check MongoDB connection
check_database() {
    # This would require MongoDB tools to be installed
    # For now, we'll check if the application can connect to the database
    # by checking the health endpoint which includes database status
    if curl -f -s "$HEALTH_URL" | grep -q '"database":"healthy"'; then
        info "Database connection is healthy"
    else
        error "Database connection check failed"
        return 1
    fi
}

# Check log files for errors
check_logs() {
    local error_count=$(journalctl -u "$SERVICE_NAME" --since "5 minutes ago" | grep -i "error\|fatal\|panic" | wc -l)
    
    if [ "$error_count" -gt 0 ]; then
        warning "Found $error_count errors in recent logs"
        journalctl -u "$SERVICE_NAME" --since "5 minutes ago" | grep -i "error\|fatal\|panic" | tail -5 >> "$ALERT_LOG"
    else
        info "No recent errors found in logs"
    fi
}

# Check SSL certificate expiration
check_ssl() {
    if command -v openssl &> /dev/null; then
        local cert_expiry=$(echo | openssl s_client -servername "$(echo $API_URL | sed 's/https:\/\///')" -connect "$(echo $API_URL | sed 's/https:\/\///')":443 2>/dev/null | openssl x509 -noout -dates | grep notAfter | cut -d= -f2)
        
        if [ -n "$cert_expiry" ]; then
            local expiry_date=$(date -d "$cert_expiry" +%s)
            local current_date=$(date +%s)
            local days_until_expiry=$(( (expiry_date - current_date) / 86400 ))
            
            if [ "$days_until_expiry" -lt 30 ]; then
                warning "SSL certificate expires in $days_until_expiry days"
            else
                info "SSL certificate expires in $days_until_expiry days"
            fi
        fi
    fi
}

# Send alert (placeholder for integration with alerting services)
send_alert() {
    local message="$1"
    local level="$2"
    
    # Log the alert
    echo "$(date +'%Y-%m-%d %H:%M:%S') - $level: $message" >> "$ALERT_LOG"
    
    # Here you can integrate with various alerting services:
    # - Slack webhook
    # - Email
    # - PagerDuty
    # - Discord webhook
    # - SMS via Twilio
    
    # Example Slack webhook (uncomment and configure):
    # if [ -n "$SLACK_WEBHOOK_URL" ]; then
    #     curl -X POST -H 'Content-type: application/json' \
    #         --data "{\"text\":\"[$level] $message\"}" \
    #         "$SLACK_WEBHOOK_URL"
    # fi
    
    # Example email alert (uncomment and configure):
    # if [ -n "$ALERT_EMAIL" ]; then
    #     echo "$message" | mail -s "[$level] Smor-Ting API Alert" "$ALERT_EMAIL"
    # fi
}

# Generate system report
generate_report() {
    local report_file="/tmp/smor-ting-system-report-$(date +%Y%m%d-%H%M%S).txt"
    
    {
        echo "=== Smor-Ting API System Report ==="
        echo "Generated: $(date)"
        echo
        echo "=== Service Status ==="
        systemctl status "$SERVICE_NAME" --no-pager
        echo
        echo "=== System Resources ==="
        echo "Disk Usage:"
        df -h
        echo
        echo "Memory Usage:"
        free -h
        echo
        echo "CPU Usage:"
        top -bn1 | head -5
        echo
        echo "=== Recent Logs ==="
        journalctl -u "$SERVICE_NAME" --since "1 hour ago" | tail -20
        echo
        echo "=== Network Connections ==="
        netstat -tuln | grep :8080
        echo
        echo "=== Process Information ==="
        ps aux | grep "$SERVICE_NAME" | grep -v grep
    } > "$report_file"
    
    info "System report generated: $report_file"
}

# Cleanup old logs
cleanup_logs() {
    # Keep only last 30 days of logs
    find /var/log/smor-ting -name "*.log" -mtime +30 -delete
    
    # Keep only last 7 days of alert logs
    find /var/log/smor-ting -name "alerts.log" -mtime +7 -exec truncate -s 0 {} \;
    
    info "Log cleanup completed"
}

# Main monitoring function
main() {
    local exit_code=0
    
    log "Starting monitoring check..."
    
    # Check service status
    if ! check_service; then
        exit_code=1
    fi
    
    # Check health endpoint
    if ! check_health; then
        exit_code=1
    fi
    
    # Check external API (if configured)
    if [ "$API_URL" != "https://api.yourdomain.com" ]; then
        if ! check_external_api; then
            exit_code=1
        fi
    fi
    
    # Check system resources
    if ! check_disk; then
        exit_code=1
    fi
    
    if ! check_memory; then
        exit_code=1
    fi
    
    if ! check_cpu; then
        exit_code=1
    fi
    
    # Check database
    if ! check_database; then
        exit_code=1
    fi
    
    # Check logs
    check_logs
    
    # Check SSL certificate
    check_ssl
    
    # Cleanup old logs
    cleanup_logs
    
    # Generate report if there are issues
    if [ $exit_code -ne 0 ]; then
        generate_report
        send_alert "Multiple issues detected. Check system report." "ERROR"
    fi
    
    if [ $exit_code -eq 0 ]; then
        log "All monitoring checks passed"
    else
        error "Some monitoring checks failed"
    fi
    
    return $exit_code
}

# Handle command line arguments
case "${1:-}" in
    "service")
        check_service
        ;;
    "health")
        check_health
        ;;
    "disk")
        check_disk
        ;;
    "memory")
        check_memory
        ;;
    "cpu")
        check_cpu
        ;;
    "database")
        check_database
        ;;
    "logs")
        check_logs
        ;;
    "ssl")
        check_ssl
        ;;
    "report")
        generate_report
        ;;
    "cleanup")
        cleanup_logs
        ;;
    *)
        main
        ;;
esac
