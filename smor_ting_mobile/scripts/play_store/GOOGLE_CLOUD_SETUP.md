# Google Cloud Console Setup for Play Console API

## ✅ **What You've Done Right:**
- ✅ Created service account in Google Cloud Console
- ✅ Assigned Editor role
- ✅ Assigned API Key Admin role
- ✅ Enabled Google Play Developer API

## 🚨 **What's Missing: Play Console Permissions**

The roles you've assigned are for **Google Cloud Console access**, but we need **Google Play Console API permissions**.

## 🔧 **Next Steps in Google Cloud Console:**

### **Step 1: Check Service Account Permissions**
1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Navigate to **"IAM & Admin"** → **"Service Accounts"**
3. Click on your service account: `play-store-deployer`
4. Go to **"Permissions"** tab
5. Check if you see these roles:
   - ✅ **Editor** (you have this)
   - ✅ **API Key Admin** (you have this)
   - ❌ **Play Console API permissions** (we need this)

### **Step 2: Add Play Console API Role**
1. In the **"Permissions"** tab, click **"Grant Access"**
2. Add the service account email: `play-store-deployer@your-project.iam.gserviceaccount.com`
3. Look for roles related to:
   - **"Play Console API"**
   - **"Android Publisher"**
   - **"Google Play Developer API"**

### **Step 3: Alternative - Create Custom Role**
If you don't see Play Console specific roles, create a custom role:
1. Go to **"IAM & Admin"** → **"Roles"**
2. Click **"Create Role"**
3. Add these permissions:
   - `androidpublisher.applications.get`
   - `androidpublisher.edits.create`
   - `androidpublisher.edits.commit`
   - `androidpublisher.bundles.upload`
   - `androidpublisher.tracks.update`

## 🎯 **The Real Issue: Play Console Configuration**

The main issue is that **Google Play Console needs to grant permissions** to your service account. This is done in Play Console, not Google Cloud Console.

### **In Google Play Console, Look For:**
1. **"Settings"** → **"API access"**
2. **"Configuration"** → **"Service accounts"**
3. **"Developer settings"** → **"API access"**
4. **"Account"** → **"API access"**

### **What You Need to Find:**
- **"Link Google Cloud project"**
- **"Grant service account access"**
- **"Release apps to testing tracks"** permission

## 🧪 **Test Current Setup:**

Let's test if the current Google Cloud setup works:
```bash
python3 scripts/play_store/deploy_aab.py
```

If it still says "Permission denied", then we definitely need to configure Play Console permissions.

## 📞 **If Still Not Working:**

The issue is likely that **Play Console hasn't granted permissions** to your service account. We need to find the API access settings in Play Console to complete the setup.

---

**Next:** Try the deployment test, and if it fails, we need to find the API access settings in Play Console.
