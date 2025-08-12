# MongoDB Atlas Setup Guide for Smor-Ting

This guide will walk you through setting up MongoDB Atlas and connecting it to your Smor-Ting application.

## 🚀 Quick Start

Run the setup script to get started:
```bash
cd smor_ting_backend
./scripts/setup_atlas.sh
```

## 📋 Prerequisites

- MongoDB Atlas account (free tier available)
- Go 1.23+ installed
- Your Smor-Ting backend code

## 🔧 Step-by-Step Setup

### Step 1: Create MongoDB Atlas Account

1. **Go to MongoDB Atlas**
   - Visit [cloud.mongodb.com](https://cloud.mongodb.com)
   - Sign up for a free account if you don't have one

2. **Create a New Project**
   - Click "New Project" in the top right
   - Name: `Smor-Ting`
   - Click "Create Project"

### Step 2: Create Your Database Cluster

1. **Build a Database**
   - Click "Build a Database"
   - Choose "FREE" tier (M0)
   - Select "AWS" as cloud provider
   - Choose region closest to your users:
     - **US East (N. Virginia)** for US users
     - **Europe (Ireland)** for European users
     - **Asia Pacific (Mumbai)** for Indian users
   - Click "Create"

2. **Wait for Cluster Creation**
   - This takes 2-3 minutes
   - You'll see a green checkmark when ready

### Step 3: Configure Database Access

1. **Create Database User**
   - Click "Database Access" in the left sidebar
   - Click "Add New Database User"
   - Authentication Method: "Password"
   - Username: `smorting_user`
   - Password: Generate a strong password (save this!)
   - Database User Privileges: "Read and write to any database"
   - Click "Add User"

### Step 4: Configure Network Access

1. **Add IP Address**
   - Click "Network Access" in the left sidebar
   - Click "Add IP Address"

2. **Choose Access Level**
   - **For Development**: Click "Allow Access from Anywhere" (0.0.0.0/0)
   - **For Production**: Add your server's specific IP address
   - Click "Confirm"

### Step 5: Get Connection String

1. **Connect to Your Cluster**
   - Click "Database" in the left sidebar
   - Click "Connect" on your cluster

2. **Choose Connection Method**
   - Select "Connect your application"
   - Choose "Go" as your driver
   - Copy the connection string

3. **Connection String Format**
   ```
   mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@YOUR_CLUSTER.mongodb.net/YOUR_DATABASE?retryWrites=true&w=majority
   ```

### Step 6: Configure Your Application

1. **Create Environment File**
   ```bash
   cd smor_ting_backend
   cp .env.example .env  # if you have an example file
   ```

2. **Update .env File**
   ```bash
   # MongoDB Atlas Configuration
   DB_DRIVER=mongodb
   DB_HOST=cluster0.xxxxx.mongodb.net
   DB_PORT=27017
   DB_NAME=smor_ting
   DB_USERNAME=YOUR_MONGODB_USERNAME
   DB_PASSWORD=YOUR_MONGODB_PASSWORD
   DB_SSL_MODE=require
   DB_IN_MEMORY=false
   MONGODB_ATLAS=true
   MONGODB_URI=mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@YOUR_CLUSTER.mongodb.net/YOUR_DATABASE?retryWrites=true&w=majority

   # JWT Configuration
   JWT_SECRET=YOUR_VERY_SECURE_JWT_SECRET_MIN_32_CHARS
   JWT_EXPIRATION=24h
   BCRYPT_COST=12

   # Server Configuration
   PORT=8080
   HOST=0.0.0.0
   ENV=production

   # Logging
   LOG_LEVEL=info
   LOG_FORMAT=json
   LOG_OUTPUT=stdout
   ```

### Step 7: Test Your Connection

1. **Run the Application**
   ```bash
   go run cmd/main.go
   ```

2. **Check Logs**
   Look for these success messages:
   ```
   ✅ Connected to MongoDB
   ✅ MongoDB indexes setup completed
   ✅ Migrations completed successfully
   ✅ Change stream service initialized successfully
   ```

3. **Test Health Endpoint**
   ```bash
   curl http://localhost:8080/health
   ```

Expected response:
```json
{
  "status": "healthy",
  "service": "smor-ting-backend",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z",
  "database": "healthy",
  "environment": "production"
}
```

## 🔒 Security Best Practices

### 1. Environment Variables
- ✅ Never commit `.env` files to version control
- ✅ Use strong, unique passwords
- ✅ Rotate passwords regularly

### 2. Network Security
- ✅ For production, restrict IP access to your server only
- ✅ Use VPC peering for enhanced security
- ✅ Enable MongoDB Atlas security features

### 3. Database Security
- ✅ Use database users with minimal required privileges
- ✅ Enable MongoDB Atlas security features
- ✅ Set up alerts for suspicious activity

## 🚀 Production Deployment

### 1. Update Network Access
- Remove "Allow Access from Anywhere"
- Add only your server's IP address
- Consider using VPC peering for enhanced security

### 2. Environment Variables
```bash
# Production settings
ENV=production
LOG_LEVEL=info
JWT_SECRET=YOUR_PRODUCTION_JWT_SECRET_MIN_32_CHARS
MONGODB_ATLAS=true
```

### 3. Monitoring Setup
- Enable MongoDB Atlas monitoring
- Set up alerts for:
  - Connection failures
  - High memory usage
  - Slow queries

## 🔧 Troubleshooting

### Common Issues

1. **Connection Timeout**
   ```
   Error: failed to connect to MongoDB: context deadline exceeded
   ```
   **Solution**: Check network access settings in Atlas

2. **Authentication Failed**
   ```
   Error: authentication failed
   ```
   **Solution**: Verify username/password in connection string

3. **SSL Certificate Issues**
   ```
   Error: x509: certificate signed by unknown authority
   ```
   **Solution**: Ensure `DB_SSL_MODE=require` is set

4. **Index Creation Failed**
   ```
   Warning: Failed to create indexes
   ```
   **Solution**: Check database user privileges

### Debug Steps

1. **Check Connection String**
   ```bash
   echo $MONGODB_URI
   ```

2. **Test Connection Manually**
   ```bash
   # Install MongoDB shell
   brew install mongodb/brew/mongodb-database-tools
   
   # Test connection
   mongosh "your-connection-string"
   ```

3. **Check Logs**
   ```bash
   go run cmd/main.go 2>&1 | grep -i mongo
   ```

## 📊 Monitoring Your Cluster

### 1. Atlas Dashboard
- Monitor cluster performance
- View connection metrics
- Check storage usage

### 2. Set Up Alerts
- Go to "Alerts" in Atlas
- Create alerts for:
  - Connection count
  - Memory usage
  - CPU usage
  - Storage usage

### 3. Performance Monitoring
- Use MongoDB Compass for query analysis
- Monitor slow queries
- Optimize indexes based on usage

## 🔄 Backup Strategy

### 1. Automated Backups
- MongoDB Atlas provides automatic backups
- Configure backup retention policy
- Test restore procedures

### 2. Point-in-Time Recovery
- Available on M10+ clusters
- Restore to any point in time
- Useful for data recovery

## 📈 Scaling Your Cluster

### 1. Upgrade Tier
- Start with M0 (free)
- Upgrade to M2/M5 for production
- Consider M10+ for advanced features

### 2. Sharding
- For high-traffic applications
- Distribute data across multiple shards
- Requires M10+ cluster

### 3. Read Replicas
- Improve read performance
- Distribute read load
- Available on M10+ clusters

## 🆘 Support Resources

- **MongoDB Atlas Documentation**: https://docs.atlas.mongodb.com
- **Connection String Format**: https://docs.mongodb.com/manual/reference/connection-string
- **Security Best Practices**: https://docs.atlas.mongodb.com/security
- **Performance Optimization**: https://docs.atlas.mongodb.com/performance

## ✅ Verification Checklist

- [ ] MongoDB Atlas account created
- [ ] Cluster created and running
- [ ] Database user created with proper privileges
- [ ] Network access configured
- [ ] Connection string obtained
- [ ] Environment variables configured
- [ ] Application connects successfully
- [ ] Health endpoint returns "healthy"
- [ ] Migrations run successfully
- [ ] Indexes created successfully
- [ ] Change streams working
- [ ] Security settings configured
- [ ] Monitoring alerts set up

---

🎉 **Congratulations!** Your Smor-Ting app is now connected to MongoDB Atlas with full offline-first capabilities, real-time synchronization, and production-ready security features. 