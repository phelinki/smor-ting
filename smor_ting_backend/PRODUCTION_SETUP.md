# 🚀 Production API Setup Guide for Smor-Ting

This guide will help you deploy your Smor-Ting API to production with enterprise-grade security, monitoring, and scalability.

## 📋 Production Checklist

### ✅ Pre-Deployment Requirements
- [ ] MongoDB Atlas cluster configured
- [ ] Domain name registered
- [ ] SSL certificate obtained
- [ ] Production server provisioned
- [ ] Environment variables configured
- [ ] Security policies defined
- [ ] Monitoring tools selected
- [ ] Backup strategy planned

## 🏗️ Infrastructure Setup

### 1. Server Requirements

**Minimum Specifications:**
- **CPU**: 2 vCPUs
- **RAM**: 4GB
- **Storage**: 20GB SSD
- **OS**: Ubuntu 22.04 LTS
- **Network**: 100 Mbps

**Recommended Specifications:**
- **CPU**: 4 vCPUs
- **RAM**: 8GB
- **Storage**: 50GB SSD
- **OS**: Ubuntu 22.04 LTS
- **Network**: 1 Gbps

### 2. Server Providers

**Recommended Options:**
- **AWS EC2**: t3.medium or t3.large
- **Google Cloud**: e2-medium or e2-standard-2
- **DigitalOcean**: 4GB RAM droplet
- **Linode**: 4GB RAM instance
- **Vultr**: 4GB RAM instance

### 3. Domain & SSL Setup

```bash
# Install Certbot for SSL
sudo apt update
sudo apt install certbot nginx

# Get SSL certificate
sudo certbot --nginx -d api.yourdomain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

## 🔧 Application Deployment

### 1. Build Production Binary

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o smor-ting-api cmd/main.go

# Or use Docker
docker build -t smor-ting-api:latest .
```

### 2. Environment Configuration

Create `/etc/smor-ting/.env`:
```bash
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
```

### 3. Systemd Service Setup

Create `/etc/systemd/system/smor-ting-api.service`:
```ini
[Unit]
Description=Smor-Ting API Server
After=network.target

[Service]
Type=simple
User=smor-ting
Group=smor-ting
WorkingDirectory=/opt/smor-ting
ExecStart=/opt/smor-ting/smor-ting-api
Restart=always
RestartSec=5
Environment=ENV=production
EnvironmentFile=/etc/smor-ting/.env

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/smor-ting/logs

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

### 4. Nginx Reverse Proxy

Create `/etc/nginx/sites-available/smor-ting-api`:
```nginx
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

    # SSL Configuration
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

    # Static Files (if any)
    location /static/ {
        alias /opt/smor-ting/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

## 🔒 Security Configuration

### 1. Firewall Setup

```bash
# UFW Firewall
sudo ufw enable
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw status
```

### 2. MongoDB Atlas Security

**Network Access:**
- Remove "Allow Access from Anywhere"
- Add only your server's IP address
- Consider VPC peering for enhanced security

**Database User:**
```bash
# Create production user with minimal privileges
db.createUser({
  user: "smorting_prod_user",
  pwd: "your-strong-production-password",
  roles: [
    { role: "readWrite", db: "smor_ting_prod" },
    { role: "dbAdmin", db: "smor_ting_prod" }
  ]
})
```

### 3. Application Security

**Rate Limiting Middleware:**
```go
// Add to your middleware
func RateLimitMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Implement rate limiting logic
        return c.Next()
    }
}
```

**Security Headers:**
```go
// Add security headers
app.Use(func(c *fiber.Ctx) error {
    c.Set("X-Content-Type-Options", "nosniff")
    c.Set("X-Frame-Options", "DENY")
    c.Set("X-XSS-Protection", "1; mode=block")
    c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
    return c.Next()
})
```

## 📊 Monitoring & Logging

### 1. Application Monitoring

**Prometheus Metrics:**
```go
// Add Prometheus metrics
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)
```

**Health Checks:**
```bash
# Custom health check script
#!/bin/bash
curl -f http://localhost:8080/health || exit 1
```

### 2. Log Management

**Structured Logging:**
```go
// Enhanced logging configuration
logger, err := logger.New("info", "json", "stdout")
if err != nil {
    log.Fatal(err)
}

// Add request logging middleware
app.Use(func(c *fiber.Ctx) error {
    start := time.Now()
    err := c.Next()
    duration := time.Since(start)
    
    logger.Info("HTTP Request",
        zap.String("method", c.Method()),
        zap.String("path", c.Path()),
        zap.Int("status", c.Response().StatusCode()),
        zap.Duration("duration", duration),
        zap.String("ip", c.IP()),
    )
    return err
})
```

### 3. External Monitoring

**Recommended Tools:**
- **Uptime Monitoring**: UptimeRobot, Pingdom
- **Application Monitoring**: New Relic, DataDog
- **Log Aggregation**: ELK Stack, Papertrail
- **Error Tracking**: Sentry, Rollbar

## 🔄 CI/CD Pipeline

### 1. GitHub Actions Workflow

Create `.github/workflows/deploy.yml`:
```yaml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Build
        run: |
          GOOS=linux GOARCH=amd64 go build -o smor-ting-api cmd/main.go
      
      - name: Deploy to Server
        uses: appleboy/ssh-action@v0.1.5
        with:
          host: ${{ secrets.HOST }}
          username: ${{ secrets.USERNAME }}
          key: ${{ secrets.SSH_KEY }}
          script: |
            cd /opt/smor-ting
            sudo systemctl stop smor-ting-api
            # Upload new binary
            sudo systemctl start smor-ting-api
            sudo systemctl status smor-ting-api
```

### 2. Docker Deployment

Create `Dockerfile`:
```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smor-ting-api cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/smor-ting-api .
EXPOSE 8080

CMD ["./smor-ting-api"]
```

## 📈 Performance Optimization

### 1. Database Optimization

**Indexes:**
```javascript
// Create indexes for common queries
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "phone": 1 })
db.services.createIndex({ "category": 1, "location": "2dsphere" })
db.bookings.createIndex({ "user_id": 1, "created_at": -1 })
```

**Connection Pooling:**
```go
// Optimize MongoDB connection pool
clientOptions := options.Client().ApplyURI(uri)
clientOptions.SetMaxPoolSize(100)
clientOptions.SetMinPoolSize(10)
clientOptions.SetMaxConnIdleTime(30 * time.Minute)
```

### 2. Application Optimization

**Caching:**
```go
// Add Redis caching
import "github.com/go-redis/redis/v8"

func (a *App) initializeCache() error {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    // Implement caching logic
}
```

**Compression:**
```go
// Enable compression
app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed,
}))
```

## 🚨 Backup & Recovery

### 1. Database Backups

**Automated Backups:**
```bash
#!/bin/bash
# Backup script
DATE=$(date +%Y%m%d_%H%M%S)
mongodump --uri="your-mongodb-uri" --out="/backups/$DATE"
tar -czf "/backups/smor_ting_$DATE.tar.gz" "/backups/$DATE"
rm -rf "/backups/$DATE"
```

**Backup Retention:**
```bash
# Keep backups for 30 days
find /backups -name "*.tar.gz" -mtime +30 -delete
```

### 2. Application Backup

```bash
# Backup application files
tar -czf "/backups/app_$(date +%Y%m%d_%H%M%S).tar.gz" /opt/smor-ting
```

## 🔧 Maintenance Procedures

### 1. Regular Maintenance

**Weekly Tasks:**
- [ ] Check disk space
- [ ] Review error logs
- [ ] Monitor performance metrics
- [ ] Update security patches

**Monthly Tasks:**
- [ ] Review and rotate secrets
- [ ] Update SSL certificates
- [ ] Review backup integrity
- [ ] Performance optimization

### 2. Update Procedures

```bash
# Zero-downtime deployment
sudo systemctl stop smor-ting-api
# Backup current version
cp /opt/smor-ting/smor-ting-api /opt/smor-ting/smor-ting-api.backup
# Deploy new version
sudo systemctl start smor-ting-api
# Verify health
curl -f http://localhost:8080/health
# Rollback if needed
# sudo systemctl stop smor-ting-api && cp /opt/smor-ting/smor-ting-api.backup /opt/smor-ting/smor-ting-api && sudo systemctl start smor-ting-api
```

## 📞 Support & Troubleshooting

### 1. Common Issues

**High Memory Usage:**
```bash
# Check memory usage
free -h
# Restart application if needed
sudo systemctl restart smor-ting-api
```

**Database Connection Issues:**
```bash
# Test MongoDB connection
mongosh "your-connection-string"
# Check network connectivity
telnet your-cluster.mongodb.net 27017
```

### 2. Emergency Procedures

**Server Down:**
1. Check server status: `sudo systemctl status smor-ting-api`
2. Check logs: `sudo journalctl -u smor-ting-api -f`
3. Restart service: `sudo systemctl restart smor-ting-api`
4. Check health endpoint: `curl http://localhost:8080/health`

**Database Issues:**
1. Check MongoDB Atlas dashboard
2. Verify network access settings
3. Test connection string
4. Check application logs for connection errors

## ✅ Production Verification

### 1. Health Checks

```bash
# Test all endpoints
curl -f https://api.yourdomain.com/health
curl -f https://api.yourdomain.com/api/v1/auth/register -X POST -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"test123","first_name":"Test","last_name":"User"}'
```

### 2. Performance Tests

```bash
# Load testing with Apache Bench
ab -n 1000 -c 10 https://api.yourdomain.com/health
```

### 3. Security Tests

```bash
# SSL Labs test
curl -s https://api.yourdomain.com/health
# Check security headers
curl -I https://api.yourdomain.com/health
```

---

🎉 **Your Smor-Ting API is now production-ready!**

**Next Steps:**
1. Set up monitoring alerts
2. Configure backup automation
3. Implement CI/CD pipeline
4. Set up staging environment
5. Plan for scaling

**Support Resources:**
- 📖 [MongoDB Atlas Documentation](https://docs.atlas.mongodb.com)
- 🔧 [Fiber Framework Documentation](https://docs.gofiber.io)
- 🚀 [Deployment Best Practices](https://12factor.net)
- 🔒 [Security Guidelines](https://owasp.org) 