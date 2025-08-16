# Google Play Console API Access Guide (2024 Interface)

## 🚨 **Important: Interface Has Changed**

Google Play Console has significantly updated its interface in 2024. The API access option is no longer under "Setup" and may be in different locations.

## 🔍 **Current API Access Locations (2024)**

### **Method 1: Test & Release Section** ⭐ (Most Likely)
1. Go to [play.google.com/console](https://play.google.com/console)
2. Sign in with your account
3. Click on your app (Smor-Ting)
4. In the left menu, click **"Test and release"**
5. Look for:
   - **"Advanced settings"**
   - **"API access"**
   - **"Service accounts"**
   - **"API credentials"**

### **Method 2: Settings Section**
1. In your app, look for **"Settings"** in the left menu
2. Check for:
   - **"API access"**
   - **"Service accounts"**
   - **"Developer settings"**

### **Method 3: Account Level**
1. Click on your **profile/account** in the top right
2. Look for:
   - **"API access"**
   - **"Service accounts"**
   - **"Developer settings"**

### **Method 4: Search Function**
1. Use the **search bar** in Play Console
2. Search for:
   - **"API"**
   - **"service account"**
   - **"API access"**

### **Method 5: App-Specific Settings**
1. Click on your **app name** (Smor-Ting)
2. Look for **"Settings"** or **"Configuration"**
3. Check for API-related options

## 🛠️ **Alternative Approach: Google Cloud Console**

If you can't find API access in Play Console, you can set it up directly in Google Cloud Console:

### **Step 1: Enable API in Google Cloud Console**
1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Select your project
3. Go to **"APIs & Services"** → **"Library"**
4. Search for **"Google Play Developer API"**
5. Click **"Enable"**

### **Step 2: Create Service Account**
1. Go to **"IAM & Admin"** → **"Service Accounts"**
2. Click **"Create Service Account"**
3. Name: `play-store-deployer`
4. Description: `Service account for Play Store deployments`
5. Click **"Create and Continue"**

### **Step 3: Grant Permissions**
1. Role: **"Editor"** (or custom role with Play Console permissions)
2. Click **"Continue"**
3. Click **"Done"**

### **Step 4: Create and Download Key**
1. Click on your service account
2. Go to **"Keys"** tab
3. Click **"Add Key"** → **"Create new key"**
4. Choose **"JSON"**
5. Click **"Create"**
6. Download the JSON file

### **Step 5: Configure Play Console Permissions**
1. Go back to Play Console
2. Try to find API access using the methods above
3. If found, link your Google Cloud project
4. Grant the service account **"Release apps to testing tracks"** permission

## 📱 **Manual Play Console Setup (If API Access Found)**

### **Step 1: Link Google Cloud Project**
1. In API access section, click **"Link Google Cloud project"**
2. Select your project
3. Click **"Link"**

### **Step 2: Create Service Account**
1. Click **"Create service account"**
2. Name: `play-store-deployer`
3. Click **"Create"**

### **Step 3: Grant Permissions**
1. Select **"Release apps to testing tracks"**
2. Click **"Grant access"**

### **Step 4: Download Service Account Key**
1. Click **"Download JSON"**
2. Save as: `scripts/play_store/service-account.json`

## 🧪 **Test Your Setup**

After completing the setup, test it:

```bash
cd smor_ting_mobile
./scripts/play_store/test_auth.sh
```

## 📞 **If Still Not Found**

If you still can't find API access:

1. **Contact Google Play Console Support**
2. **Check Play Console Help**: [support.google.com/googleplay/android-developer](https://support.google.com/googleplay/android-developer)
3. **Use the Google Cloud Console approach** (Method 2 above)
4. **Check if your account has the right permissions**

## 🔄 **Interface Updates**

Google frequently updates Play Console. If this guide becomes outdated:

1. Use the search function in Play Console
2. Check the latest Play Console documentation
3. Contact Google support for current interface guidance

## ✅ **Success Indicators**

You'll know you've found the right place when you see:
- Options to link Google Cloud projects
- Service account management
- API credentials section
- Permission management for Play Console access
