# 🚀 Production Deployment Guide

This guide provides a complete overview of deploying your Smor-Ting API to production.

## 📋 Quick Start

### 1. Server Setup
```bash
# On your production server
sudo ./scripts/setup_production.sh
```

### 2. Configure Environment
```bash
# Edit environment file
sudo nano /etc/smor-ting/.env

# Copy from template and update values
sudo cp /etc/smor-ting/.env.template /etc/smor-ting/.env
```

### 3. Deploy Application
```bash
# From your development machine
./scripts/deploy_production.sh
```

## 🏗️ Infrastructure Overview

### Server Requirements
- **CPU**: 2-4 vCPUs
- **RAM**: 4-8GB
- **Storage**: 20-50GB SSD
- **OS**: Ubuntu 22.04 LTS
- **Network**: 100+ Mbps

### Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │   Nginx Proxy   │    │  Smor-Ting API  │
│   (Optional)    │───▶│   (SSL/TLS)     │───▶│   (Port 8080)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │  MongoDB Atlas  │
                                              │   (Database)    │
                                              └─────────────────┘
```

## 🔧 Deployment Options

### Option 1: Traditional Server Deployment
```bash
# 1. Setup server
sudo ./scripts/setup_production.sh

# 2. Configure environment
sudo nano /etc/smor-ting/.env

# 3. Deploy application
./scripts/deploy_production.sh

# 4. Setup SSL
sudo certbot --nginx -d api.yourdomain.com
```

### Option 2: Docker Deployment
```bash
# 1. Build and run with Docker Compose
docker-compose up -d

# 2. Or build and run manually
docker build -t smor-ting-api .
docker run -d --name smor-ting-api -p 8080:8080 smor-ting-api
```

### Option 3: Cloud Platform Deployment
```bash
# AWS ECS/Fargate
aws ecs create-service --cluster smor-ting --service-name api

# Google Cloud Run
gcloud run deploy smor-ting-api --source .

# DigitalOcean App Platform
doctl apps create --spec app.yaml
```

## 🔒 Security Configuration

### Environment Variables
```bash
# Required for production
ENV=production
JWT_SECRET=YOUR_PRODUCTION_JWT_SECRET_MIN_32_CHARS
MONGODB_URI=mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@YOUR_CLUSTER.mongodb.net/YOUR_DATABASE
MONGODB_ATLAS=true

# Security headers
CORS_ALLOW_ORIGINS=https://yourdomain.com
BCRYPT_COST=12
```

### SSL/TLS Setup
```bash
# Automatic SSL with Let's Encrypt
sudo certbot --nginx -d api.yourdomain.com

# Manual SSL certificate
sudo cp your-cert.pem /etc/ssl/certs/
sudo cp your-key.pem /etc/ssl/private/
```

### Firewall Configuration
```bash
# UFW Firewall
sudo ufw enable
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw status
```

## 📊 Monitoring & Alerting

### Health Checks
```bash
# Manual health check
curl -f https://api.yourdomain.com/health

# Automated monitoring
./scripts/monitor.sh

# Check specific components
./scripts/monitor.sh service
./scripts/monitor.sh health
./scripts/monitor.sh disk
./scripts/monitor.sh memory
```

### Log Management
```bash
# View application logs
sudo journalctl -u smor-ting-api -f

# View monitoring logs
tail -f /var/log/smor-ting/monitor.log

# View alerts
tail -f /var/log/smor-ting/alerts.log
```

### Performance Monitoring
```bash
# System resources
htop
df -h
free -h

# Application metrics
curl https://api.yourdomain.com/health | jq
```

## 🔄 CI/CD Pipeline

### GitHub Actions Setup
1. Add secrets to your repository:
   - `HOST`: Your server IP
   - `USERNAME`: SSH username
   - `SSH_PRIVATE_KEY`: SSH private key
   - `DOMAIN`: Your domain name

2. Push to main branch for automatic deployment:
```bash
git push origin main
```

### Manual Deployment
```bash
# Deploy with version tag
./scripts/deploy_production.sh v1.2.3

# Rollback to previous version
sudo systemctl stop smor-ting-api
sudo cp /opt/smor-ting/backups/smor-ting-api_*.backup /opt/smor-ting/smor-ting-api
sudo systemctl start smor-ting-api
```

## 🗄️ Database Management

### MongoDB Atlas Setup
1. Create cluster in MongoDB Atlas
2. Configure network access (your server IP)
3. Create database user
4. Get connection string
5. Update environment variables

### Database Backups
```bash
# Manual backup
mongodump --uri="your-connection-string" --out="/backups/$(date +%Y%m%d)"

# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
mongodump --uri="$MONGODB_URI" --out="/backups/$DATE"
tar -czf "/backups/smor_ting_$DATE.tar.gz" "/backups/$DATE"
```

## 🚨 Troubleshooting

### Common Issues

**Service won't start:**
```bash
# Check logs
sudo journalctl -u smor-ting-api -f

# Check configuration
sudo systemctl status smor-ting-api

# Restart service
sudo systemctl restart smor-ting-api
```

**Database connection issues:**
```bash
# Test connection
mongosh "your-connection-string"

# Check network access
telnet your-cluster.mongodb.net 27017

# Verify environment variables
sudo cat /etc/smor-ting/.env
```

**SSL certificate issues:**
```bash
# Check certificate
openssl s_client -connect api.yourdomain.com:443

# Renew certificate
sudo certbot renew

# Check Nginx configuration
sudo nginx -t
```

### Emergency Procedures

**Server down:**
1. SSH into server
2. Check service status: `sudo systemctl status smor-ting-api`
3. Restart service: `sudo systemctl restart smor-ting-api`
4. Check health: `curl http://localhost:8080/health`

**Database down:**
1. Check MongoDB Atlas dashboard
2. Verify network access settings
3. Test connection string
4. Check application logs for errors

**SSL issues:**
1. Check certificate expiration
2. Renew certificate: `sudo certbot renew`
3. Restart Nginx: `sudo systemctl restart nginx`

## 📈 Scaling

### Vertical Scaling
```bash
# Upgrade server resources
# - Increase CPU/RAM
# - Upgrade storage
# - Optimize application settings
```

### Horizontal Scaling
```bash
# Load balancer setup
# - Multiple application instances
# - Database read replicas
# - CDN for static assets
```

### Performance Optimization
```bash
# Database indexes
db.users.createIndex({ "email": 1 }, { unique: true })
db.services.createIndex({ "category": 1, "location": "2dsphere" })

# Application caching
# - Redis for session storage
# - CDN for static files
# - Database query optimization
```

## 🔧 Maintenance

### Regular Tasks
- [ ] Weekly: Check disk space and logs
- [ ] Monthly: Update SSL certificates
- [ ] Quarterly: Review security settings
- [ ] Annually: Update dependencies

### Backup Strategy
```bash
# Database backups (daily)
0 2 * * * /usr/local/bin/backup-database.sh

# Application backups (weekly)
0 3 * * 0 /usr/local/bin/backup-application.sh

# Log rotation (daily)
0 4 * * * /usr/local/bin/rotate-logs.sh
```

## 📞 Support

### Useful Commands
```bash
# Service management
sudo systemctl start/stop/restart smor-ting-api
sudo systemctl status smor-ting-api
sudo systemctl enable smor-ting-api

# Log viewing
sudo journalctl -u smor-ting-api -f
sudo journalctl -u smor-ting-api --since "1 hour ago"

# Configuration
sudo nano /etc/smor-ting/.env
sudo nano /etc/systemd/system/smor-ting-api.service
sudo nano /etc/nginx/sites-available/smor-ting-api

# Monitoring
./scripts/monitor.sh
./scripts/monitor.sh report
```

### Documentation Links
- [MongoDB Atlas Setup](ATLAS_SETUP.md)
- [Production Setup](PRODUCTION_SETUP.md)
- [Quick Start](QUICK_START.md)
- [API Documentation](docs/api.md)

---

🎉 **Your Smor-Ting API is now production-ready!**

**Next Steps:**
1. Set up monitoring alerts — use `scripts/monitor.sh` via cron/systemd timer (example below)
2. Configure backup automation — use `scripts/backup.sh` daily with retention
3. Implement CI/CD pipeline — GitHub Actions workflow at `.github/workflows/backend-ci.yml`
4. Set up staging environment — set `ENV=staging` and provide base64 secrets; staging fails-closed like prod
5. Plan for scaling — see Scaling section

### Monitoring (example cron)
```
*/5 * * * * /usr/local/bin/bash /opt/smor-ting/scripts/monitor.sh >> /var/log/smor-ting/monitor.log 2>&1
```

### Backups (example cron)
```
0 2 * * * BACKUP_DIR=/var/backups/smor-ting DB_NAME=smor_ting /usr/local/bin/bash /opt/smor-ting/scripts/backup.sh run
```

### Staging Environment
- Set `ENV=staging`
- Provide base64-encoded `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `ENCRYPTION_KEY`, `PAYMENT_ENCRYPTION_KEY`
- Provide MoMo and SmileID credentials

### Scaling
- Horizontal: multiple app replicas behind a load balancer
- Database: increase MongoDB Atlas cluster tier; enable connection pool tuning
- Caching: add edge CDN and server-side caching for read-heavy endpoints
- Background jobs: offload long-running tasks to a worker queue

**Support Resources:**
- 📖 [MongoDB Atlas Documentation](https://docs.atlas.mongodb.com)
- 🔧 [Fiber Framework Documentation](https://docs.gofiber.io)
- 🚀 [Deployment Best Practices](https://12factor.net)
- 🔒 [Security Guidelines](https://owasp.org)
