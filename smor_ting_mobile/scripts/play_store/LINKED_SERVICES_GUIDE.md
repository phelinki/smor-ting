# Google Play Console Linked Services Configuration

## 🎯 **Found It! Linked Services is the Right Place**

The "Linked services" section in Play Console is where you configure service accounts for API access.

## 📋 **Step-by-Step Configuration:**

### **Step 1: Access Linked Services**
1. You're already in Play Console Settings
2. Click on **"Linked services"**
3. Look for options related to:
   - Google Cloud projects
   - Service accounts
   - API access

### **Step 2: Link Your Google Cloud Project**
1. In Linked Services, look for **"Google Cloud"** or **"Cloud projects"**
2. Click **"Link project"** or **"Add project"**
3. Select your Google Cloud project (where you created the service account)
4. Click **"Link"** or **"Connect"**

### **Step 3: Configure Service Account**
1. After linking the project, look for **"Service accounts"** or **"API credentials"**
2. You should see your service account: `play-store-deployer`
3. Click on it or look for **"Grant access"** or **"Configure permissions"**

### **Step 4: Grant Permissions**
1. Look for permission options like:
   - **"Release apps to testing tracks"**
   - **"Manage app releases"**
   - **"Upload app bundles"**
2. Select the appropriate permissions
3. Click **"Save"** or **"Grant access"**

## 🔍 **What to Look For:**

### **Common Options in Linked Services:**
- **Google Cloud projects**
- **Service accounts**
- **API credentials**
- **Access management**
- **Permissions**

### **Permission Names to Look For:**
- **"Release apps to testing tracks"** ⭐ (Most important)
- **"Manage app releases"**
- **"Upload app bundles"**
- **"Access to Play Console API"**

## 🚨 **If You Don't See Your Service Account:**

### **Option 1: Create New Service Account**
1. Look for **"Create service account"** or **"Add service account"**
2. Name it: `play-store-deployer`
3. Grant the permissions mentioned above

### **Option 2: Import Existing Service Account**
1. Look for **"Import"** or **"Add existing"**
2. Select your Google Cloud project
3. Choose the `play-store-deployer` service account
4. Grant permissions

## ✅ **Success Indicators:**

You'll know it's working when you see:
- ✅ Your Google Cloud project linked
- ✅ Service account `play-store-deployer` listed
- ✅ Permissions granted (especially "Release apps to testing tracks")
- ✅ Status showing "Active" or "Connected"

## 🧪 **Test After Configuration:**

Once you've configured the permissions, test it:
```bash
python3 scripts/play_store/deploy_aab.py
```

## 📞 **If Still Having Issues:**

1. **Take a screenshot** of what you see in Linked Services
2. **Check the exact options** available
3. **Look for any error messages**
4. **Try refreshing the page** if options don't appear

---

**Next:** Configure the service account permissions in Linked Services, then we can deploy your app!
