# Railway Configuration Guide

## Simple Railway Setup (No TOML file)

Instead of using a complex `railway.toml` file, configure everything directly in the Railway dashboard for simplicity and reliability.

### Railway Dashboard Settings

#### 1. **Deploy Settings**
- **Root Directory**: `smor_ting_backend`
- **Build Command**: (leave empty - Railway will auto-detect Dockerfile)
- **Start Command**: `./smor-ting-api`

#### 2. **Build Settings**
- **Builder**: Dockerfile
- **Dockerfile Path**: `Dockerfile` (relative to root directory)

#### 3. **Health Check**
- **Health Check Path**: `/health`
- **Health Check Timeout**: 300 seconds

#### 4. **Environment Variables**
Set these in Railway dashboard Variables section:
```
PORT=8080
ENV=production
HOST=0.0.0.0
MONGODB_URI=(your MongoDB connection string)
JWT_SECRET=(your JWT secret)
```

#### 5. **Watch Patterns** 
Railway will automatically watch the `smor_ting_backend` directory when it's set as the root directory.

### Benefits of This Approach
✅ **Simplified Configuration**: No TOML file conflicts
✅ **Visual Management**: Easy to see and modify settings in dashboard  
✅ **Reliable Deployments**: Railway handles file watching automatically
✅ **Less Configuration Drift**: Settings are centralized in one place
✅ **Easier Debugging**: Clear separation between code and deployment config

### Deployment Process
1. Push code changes to `main` branch
2. Railway automatically detects changes in `smor_ting_backend/` 
3. Builds using `smor_ting_backend/Dockerfile`
4. Deploys with settings configured in dashboard

### Files to Remove
- ✅ `railway.toml` (removed)
- ✅ `smor_ting_backend/railway.json` (already removed)
- ✅ `nixpacks.toml` (already removed)

This setup eliminates configuration complexity while maintaining all the functionality we need.
