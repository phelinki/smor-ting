# Railway Deployment Guide for Smor-Ting Backend

This guide will walk you through deploying your Smor-Ting Go backend to Railway.

## 🚀 Quick Deploy

1. **Install Railway CLI**
   ```bash
   npm install -g @railway/cli
   ```

2. **Login to Railway**
   ```bash
   railway login
   ```

3. **Deploy from your project directory**
   ```bash
   cd smor_ting_backend
   railway init
   railway up
   ```

## 📋 Prerequisites

- Railway account (free at [railway.app](https://railway.app))
- MongoDB Atlas cluster (already set up)
- GitHub repository with your code

## 🔧 Step-by-Step Setup

### Step 1: Prepare Your Repository

1. **Ensure your Dockerfile is ready**
   ```bash
   # Your existing Dockerfile is already optimized
   docker build -t smor-ting-api .
   ```

2. **Check your environment variables**
   ```bash
   # Make sure these are set in Railway
   ENV=production
   PORT=8080
   HOST=0.0.0.0
   MONGODB_URI=your_mongodb_atlas_connection_string
   JWT_SECRET=your_jwt_secret
   # ... other variables
   ```

### Step 2: Connect to Railway

1. **Go to Railway Dashboard**
   - Visit [railway.app](https://railway.app)
   - Sign in with GitHub

2. **Create New Project**
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your `smor-ting` repository

3. **Configure Service**
   - Railway will detect your Dockerfile
   - Set the service name: `smor-ting-api`

### Step 3: Configure Environment Variables

In Railway dashboard, go to your service → Variables tab:

```bash
# Server Configuration
ENV=production
PORT=8080
HOST=0.0.0.0
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
IDLE_TIMEOUT=120s

# MongoDB Atlas Configuration
DB_DRIVER=mongodb
DB_HOST=your_cluster.mongodb.net
DB_PORT=27017
DB_NAME=smor_ting
DB_USERNAME=your_username
DB_PASSWORD=your_password
DB_SSL_MODE=require
DB_IN_MEMORY=false
MONGODB_ATLAS=true
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/smor_ting?retryWrites=true&w=majority

# Authentication
JWT_ACCESS_SECRET=your_base64_encoded_access_secret
JWT_REFRESH_SECRET=your_base64_encoded_refresh_secret
JWT_EXPIRATION=24h
BCRYPT_COST=12

# Security
ENCRYPTION_KEY=your_base64_encoded_encryption_key
PAYMENT_ENCRYPTION_KEY=your_base64_encoded_payment_key

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout

# CORS (Production)
CORS_ALLOW_ORIGINS=https://your-frontend-domain.com
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_CREDENTIALS=true
```

### Step 4: Deploy

1. **Automatic Deployment**
   - Railway will automatically deploy when you push to main branch
   - Or manually deploy from dashboard

2. **Check Deployment**
   ```bash
   # View logs
   railway logs
   
   # Check status
   railway status
   ```

### Step 5: Configure Custom Domain

1. **Add Custom Domain**
   - Go to your service → Settings → Domains
   - Add your domain: `api.smor-ting.com`
   - Railway will provide SSL certificate

2. **Update CORS Settings**
   ```bash
   CORS_ALLOW_ORIGINS=https://smor-ting.com,https://www.smor-ting.com
   ```

## 🔒 Security Configuration

### 1. Environment Variables
- ✅ All secrets are environment variables
- ✅ No hardcoded credentials
- ✅ Railway encrypts sensitive data

### 2. Network Security
- ✅ HTTPS enforced
- ✅ CORS properly configured
- ✅ MongoDB Atlas network restrictions

### 3. Application Security
- ✅ JWT tokens with proper expiration
- ✅ bcrypt password hashing
- ✅ Input validation and sanitization

## 📊 Monitoring & Logs

### 1. Railway Dashboard
- Real-time logs
- Deployment history
- Performance metrics
- Error tracking

### 2. Health Checks
```bash
# Your app already has health endpoint
curl https://your-railway-app.railway.app/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "smor-ting-backend",
  "version": "1.0.0",
  "database": "healthy",
  "environment": "production"
}
```

## 🚀 Production Checklist

- [ ] Railway project created
- [ ] GitHub repository connected
- [ ] Environment variables configured
- [ ] MongoDB Atlas connection working
- [ ] Custom domain configured
- [ ] SSL certificate active
- [ ] Health endpoint responding
- [ ] CORS properly configured
- [ ] Logs being captured
- [ ] Monitoring alerts set up

## 🔧 Troubleshooting

### Common Issues

1. **Build Fails**
   ```bash
   # Check Dockerfile
   docker build -t test .
   
   # Check logs
   railway logs
   ```

2. **Database Connection Issues**
   ```bash
   # Verify MongoDB Atlas settings
   # Check network access
   # Verify connection string
   ```

3. **Environment Variables**
   ```bash
   # List all variables
   railway variables
   
   # Check specific variable
   railway variables get MONGODB_URI
   ```

### Debug Commands

```bash
# View real-time logs
railway logs --follow

# Check service status
railway status

# Restart service
railway service restart

# View deployment history
railway deployments
```

## 💰 Cost Optimization

### Railway Pricing
- **Free Tier**: $5/month credit
- **Pay-as-you-go**: $0.000463 per GB-second
- **Predictable**: No hidden fees

### Cost Saving Tips
1. **Use free tier** for development
2. **Optimize Docker image** size
3. **Monitor usage** in dashboard
4. **Scale down** during low traffic

## 🔄 CI/CD Setup

### Automatic Deployments
1. **Connect GitHub repository**
2. **Railway auto-deploys** on push to main
3. **Preview deployments** for pull requests

### Manual Deployments
```bash
# Deploy specific branch
railway up --service smor-ting-api

# Deploy with specific environment
railway up --environment production
```

## 📈 Scaling

### Auto-Scaling
- Railway automatically scales based on traffic
- No manual configuration needed
- Pay only for what you use

### Manual Scaling
- Set minimum/maximum instances
- Configure resource limits
- Monitor performance metrics

## 🆘 Support

- **Railway Documentation**: https://docs.railway.app
- **Community Discord**: https://discord.gg/railway
- **GitHub Issues**: https://github.com/railwayapp/cli/issues

---

🎉 **Your Smor-Ting backend is now deployed on Railway with enterprise-grade reliability, automatic scaling, and zero-downtime deployments!**
