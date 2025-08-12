# 🚀 Quick Start: MongoDB Atlas + Smor-Ting

Get your Smor-Ting app connected to MongoDB Atlas in 5 minutes!

## ⚡ Quick Setup

### 1. Run Setup Script
```bash
cd smor_ting_backend
./scripts/setup_atlas.sh
```

### 2. Create MongoDB Atlas Cluster
1. Go to [cloud.mongodb.com](https://cloud.mongodb.com)
2. Create new project: "Smor-Ting"
3. Build database → FREE tier (M0)
4. Create database user: `smorting_user` + strong password
5. Network access → "Allow Access from Anywhere" (for dev)
6. Get connection string from "Connect" button

### 3. Update Environment
Edit `.env` file with your connection string:
```bash
MONGODB_URI=mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@YOUR_CLUSTER.mongodb.net/YOUR_DATABASE?retryWrites=true&w=majority
MONGODB_ATLAS=true
JWT_SECRET=YOUR_JWT_SECRET_MIN_32_CHARS
```

### 4. Test Connection
```bash
./scripts/test_connection.sh
```

### 5. Start Your App
```bash
go run cmd/main.go
```

## ✅ Success Indicators

Look for these logs:
```
✅ Connected to MongoDB
✅ MongoDB indexes setup completed
✅ Migrations completed successfully
✅ Change stream service initialized successfully
```

## 🔗 Test Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","first_name":"Test","last_name":"User"}'
```

## 🆘 Need Help?

- 📖 Full guide: `ATLAS_SETUP.md`
- 🧪 Test script: `./scripts/test_connection.sh`
- 🔧 Setup script: `./scripts/setup_atlas.sh`

---

🎉 **You're all set!** Your Smor-Ting app now has:
- ✅ MongoDB Atlas cloud database
- ✅ Offline-first architecture
- ✅ Real-time synchronization
- ✅ Production-ready security
- ✅ Automatic backups
- ✅ Scalable infrastructure 