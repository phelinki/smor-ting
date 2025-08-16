# Google Play Console App Integrity Configuration

## 🎯 **Found It! App Integrity is the Right Place**

The "App integrity" section in Play Console is where you configure service accounts for API access. This is a newer location in the current interface.

## 📋 **What You Need to Configure:**

### **✅ Service Account Permissions** (What we need)
- Look for **"API access"** or **"Service accounts"**
- Configure permissions for deploying to testing tracks
- Grant "Release apps to testing tracks" permission

### **❌ Inline Installs** (Not what we need)
- This is for marketing and app discovery
- Allows other apps to install your app directly
- Not related to deploying your app to testing

## 🔍 **In App Integrity, Look For:**

### **Service Account Configuration:**
1. **"API access"** or **"Service accounts"**
2. **"Google Cloud project"** linking
3. **"Permissions"** or **"Access management"**
4. **"Release apps to testing tracks"** permission

### **What to Configure:**
1. **Link your Google Cloud project** (if not already done)
2. **Find your service account**: `play-store-deployer`
3. **Grant permissions**:
   - ✅ **"Release apps to testing tracks"** (Most important)
   - ✅ **"Manage app releases"**
   - ✅ **"Upload app bundles"**

## 🚨 **Skip Inline Installs:**

The inline installs feature you found is for:
- Marketing campaigns
- App discovery
- Cross-app installations
- **Not for deploying to testing tracks**

## ✅ **Success Indicators:**

You'll know you've found the right configuration when you see:
- ✅ Your Google Cloud project linked
- ✅ Service account `play-store-deployer` listed
- ✅ Permission options for "Release apps to testing tracks"
- ✅ Status showing "Active" or "Connected"

## 🧪 **Test After Configuration:**

Once you've configured the service account permissions (not inline installs), test it:
```bash
python3 scripts/play_store/deploy_aab.py
```

## 📞 **If You Can't Find Service Account Settings:**

Look for these terms in App Integrity:
- **"API access"**
- **"Service accounts"**
- **"Google Cloud integration"**
- **"Developer API"**
- **"Access management"**

---

**Next:** Find the service account configuration in App Integrity (not inline installs), then we can deploy your app!
