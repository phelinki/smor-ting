#!/bin/bash

# Production Setup Script for Smor-Ting API
# This script sets up a production server with all necessary components

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
USER_NAME="smor-ting"
GROUP_NAME="smor-ting"
DOMAIN="api.yourdomain.com"  # Change this to your domain

# Logging function
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

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root (use sudo)"
    fi
}

# Update system
update_system() {
    log "Updating system packages..."
    apt update
    apt upgrade -y
    log "System updated successfully"
}

# Install required packages
install_packages() {
    log "Installing required packages..."
    
    apt install -y \
        curl \
        wget \
        git \
        build-essential \
        nginx \
        certbot \
        python3-certbot-nginx \
        ufw \
        htop \
        unzip \
        software-properties-common \
        apt-transport-https \
        ca-certificates \
        gnupg \
        lsb-release
    
    log "Packages installed successfully"
}

# Install Go
install_go() {
    log "Installing Go..."
    
    # Check if Go is already installed
    if command -v go &> /dev/null; then
        log "Go is already installed"
        return
    fi
    
    # Download and install Go
    GO_VERSION="1.23.0"
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    
    # Add Go to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    source /etc/profile
    
    # Cleanup
    rm go${GO_VERSION}.linux-amd64.tar.gz
    
    log "Go installed successfully"
}

# Create application user
create_user() {
    log "Creating application user..."
    
    # Create user and group
    if ! id "$USER_NAME" &>/dev/null; then
        useradd -r -s /bin/bash -d "$APP_DIR" -m "$USER_NAME"
        log "User $USER_NAME created"
    else
        log "User $USER_NAME already exists"
    fi
    
    # Create application directory
    mkdir -p "$APP_DIR"
    chown "$USER_NAME:$GROUP_NAME" "$APP_DIR"
    
    # Create necessary directories
    mkdir -p "$APP_DIR/logs"
    mkdir -p "$APP_DIR/backups"
    mkdir -p "$APP_DIR/static"
    chown -R "$USER_NAME:$GROUP_NAME" "$APP_DIR"
    
    log "Application user and directories created"
}

# Setup firewall
setup_firewall() {
    log "Setting up firewall..."
    
    ufw --force enable
    ufw default deny incoming
    ufw default allow outgoing
    ufw allow ssh
    ufw allow 80/tcp
    ufw allow 443/tcp
    
    log "Firewall configured successfully"
}

# Setup Nginx
setup_nginx() {
    log "Setting up Nginx..."
    
    # Create Nginx configuration
    cat > /etc/nginx/sites-available/smor-ting-api << 'EOF'
upstream smor_ting_api {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL Configuration (will be updated by Certbot)
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security Headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin";

    # Rate Limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20 nodelay;

    # Gzip Compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml;

    # Proxy Configuration
    location / {
        proxy_pass http://smor_ting_api;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 86400;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
    }

    # Health Check
    location /health {
        access_log off;
        proxy_pass http://smor_ting_api;
    }

    # Static Files
    location /static/ {
        alias /opt/smor-ting/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
EOF

    # Enable site
    ln -sf /etc/nginx/sites-available/smor-ting-api /etc/nginx/sites-enabled/
    
    # Remove default site
    rm -f /etc/nginx/sites-enabled/default
    
    # Test Nginx configuration
    nginx -t
    
    # Start Nginx
    systemctl enable nginx
    systemctl start nginx
    
    log "Nginx configured successfully"
}

# Create systemd service
create_service() {
    log "Creating systemd service..."
    
    cat > /etc/systemd/system/smor-ting-api.service << EOF
[Unit]
Description=Smor-Ting API Server
After=network.target

[Service]
Type=simple
User=$USER_NAME
Group=$GROUP_NAME
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/$APP_NAME
Restart=always
RestartSec=5
Environment=ENV=production
EnvironmentFile=/etc/smor-ting/.env

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$APP_DIR/logs

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
EOF

    # Reload systemd
    systemctl daemon-reload
    
    log "Systemd service created successfully"
}

# Create environment file template
create_env_template() {
    log "Creating environment file template..."
    
    mkdir -p /etc/smor-ting
    
    cat > /etc/smor-ting/.env.template << 'EOF'
# Production Environment
ENV=production

# Server Configuration
PORT=8080
HOST=127.0.0.1
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
IDLE_TIMEOUT=120s

# MongoDB Atlas Configuration
DB_DRIVER=mongodb
DB_HOST=your-cluster.mongodb.net
DB_PORT=27017
DB_NAME=smor_ting_prod
DB_USERNAME=smorting_prod_user
DB_PASSWORD=your-strong-production-password
DB_SSL_MODE=require
DB_IN_MEMORY=false
MONGODB_ATLAS=true
MONGODB_URI=mongodb+srv://smorting_prod_user:password@cluster.mongodb.net/smor_ting_prod?retryWrites=true&w=majority

# JWT Configuration (CHANGE THIS!)
JWT_SECRET=your-super-secret-production-jwt-key-min-32-chars
JWT_EXPIRATION=24h
BCRYPT_COST=12

# CORS Configuration
CORS_ALLOW_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization,X-Requested-With
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Security
BCRYPT_COST=12
SESSION_SECRET=your-session-secret-key
EOF

    chown "$USER_NAME:$GROUP_NAME" /etc/smor-ting/.env.template
    
    log "Environment template created at /etc/smor-ting/.env.template"
    warning "Please update the environment file with your actual values"
}

# Setup SSL certificate
setup_ssl() {
    log "Setting up SSL certificate..."
    
    # Check if domain is provided
    if [[ "$DOMAIN" == "api.yourdomain.com" ]]; then
        warning "Please update the DOMAIN variable in this script before running SSL setup"
        warning "You can run SSL setup manually with: sudo certbot --nginx -d your-domain.com"
        return
    fi
    
    # Get SSL certificate
    certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos --email admin@yourdomain.com
    
    # Setup auto-renewal
    (crontab -l 2>/dev/null; echo "0 12 * * * /usr/bin/certbot renew --quiet") | crontab -
    
    log "SSL certificate configured successfully"
}

# Setup monitoring
setup_monitoring() {
    log "Setting up basic monitoring..."
    
    # Create log rotation
    cat > /etc/logrotate.d/smor-ting << EOF
$APP_DIR/logs/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 $USER_NAME $GROUP_NAME
    postrotate
        systemctl reload smor-ting-api
    endscript
}
EOF

    # Create basic monitoring script
    cat > /usr/local/bin/monitor-smor-ting << 'EOF'
#!/bin/bash

# Basic monitoring script for Smor-Ting API

LOG_FILE="/var/log/smor-ting/monitor.log"
HEALTH_URL="http://localhost:8080/health"

# Check if service is running
if ! systemctl is-active --quiet smor-ting-api; then
    echo "$(date): Service is down, attempting restart" >> "$LOG_FILE"
    systemctl restart smor-ting-api
fi

# Check health endpoint
if ! curl -f -s "$HEALTH_URL" > /dev/null; then
    echo "$(date): Health check failed" >> "$LOG_FILE"
fi

# Check disk space
DISK_USAGE=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 80 ]; then
    echo "$(date): Disk usage is high: ${DISK_USAGE}%" >> "$LOG_FILE"
fi

# Check memory usage
MEM_USAGE=$(free | awk 'NR==2{printf "%.0f", $3*100/$2}')
if [ "$MEM_USAGE" -gt 80 ]; then
    echo "$(date): Memory usage is high: ${MEM_USAGE}%" >> "$LOG_FILE"
fi
EOF

    chmod +x /usr/local/bin/monitor-smor-ting
    
    # Add to crontab (run every 5 minutes)
    (crontab -l 2>/dev/null; echo "*/5 * * * * /usr/local/bin/monitor-smor-ting") | crontab -
    
    log "Basic monitoring configured"
}

# Print next steps
print_next_steps() {
    echo
    echo -e "${GREEN}🎉 Production setup completed successfully!${NC}"
    echo
    echo -e "${BLUE}Next steps:${NC}"
    echo "1. Update domain in Nginx configuration: /etc/nginx/sites-available/smor-ting-api"
    echo "2. Configure environment variables: /etc/smor-ting/.env.template"
    echo "3. Set up SSL certificate: sudo certbot --nginx -d your-domain.com"
    echo "4. Deploy your application: ./scripts/deploy_production.sh"
    echo "5. Set up MongoDB Atlas cluster"
    echo "6. Configure monitoring and alerts"
    echo
    echo -e "${YELLOW}Important security notes:${NC}"
    echo "- Change default JWT secret in environment file"
    echo "- Configure MongoDB Atlas network access"
    echo "- Set up proper backup strategy"
    echo "- Monitor logs regularly"
    echo
    echo -e "${GREEN}Your API will be available at: https://your-domain.com${NC}"
}

# Main function
main() {
    log "Starting production setup..."
    
    # Check if running as root
    check_root
    
    # Update system
    update_system
    
    # Install packages
    install_packages
    
    # Install Go
    install_go
    
    # Create user
    create_user
    
    # Setup firewall
    setup_firewall
    
    # Setup Nginx
    setup_nginx
    
    # Create service
    create_service
    
    # Create environment template
    create_env_template
    
    # Setup monitoring
    setup_monitoring
    
    # Print next steps
    print_next_steps
}

# Run main function
main "$@" 