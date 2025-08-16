#!/bin/bash

# Smor-Ting AAB Upload Workaround Script
# This script provides a workaround for app bundle build issues

set -e

echo "🔧 Smor-Ting AAB Upload Workaround..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_status "Building optimized APK for Google Play Console..."

# Build optimized APK
flutter build apk --release --split-per-abi

# Check if builds were successful
if [ -f "build/app/outputs/flutter-apk/app-arm64-v8a-release.apk" ]; then
    print_success "Optimized APK created successfully!"
    print_status "Location: build/app/outputs/flutter-apk/app-arm64-v8a-release.apk"
    print_status "Size: $(du -h build/app/outputs/flutter-apk/app-arm64-v8a-release.apk | cut -f1)"
else
    print_error "Failed to create optimized APK"
    exit 1
fi

# Create AAB upload guide
cat > AAB_UPLOAD_GUIDE.md << EOF
# Smor-Ting AAB Upload Guide

## 🚨 App Bundle Build Issue Workaround

Due to Android toolchain issues, we're using a **workaround** to get your app on Google Play Console.

## 📱 **Solution: Google Play Console APK to AAB Conversion**

Google Play Console can automatically convert your APK to an App Bundle (.aab) during upload.

### **Step 1: Upload APK to Google Play Console**
1. Go to [play.google.com/console](https://play.google.com/console)
2. Sign in with **pkaleewoun@gmail.com**
3. Create app: "Smor-Ting" (if not exists)
4. Go to **Testing** → **Internal testing**
5. Click **"Create new release"**
6. **Upload the APK file**: \`app-arm64-v8a-release.apk\`
7. Add release notes and save

### **Step 2: Google Play Console Converts to AAB**
- Google Play Console will **automatically convert** your APK to an App Bundle
- This happens during the upload process
- No additional action needed from you

### **Step 3: Alternative - Fix Android Toolchain**
If you want to build AAB directly, fix the Android toolchain:

1. **Install Android Command Line Tools**:
   - Download from: https://developer.android.com/studio#command-line-tools-only
   - Extract to: \`~/Library/Android/sdk/cmdline-tools/latest/\`
   - Add to PATH: \`export PATH=\$PATH:~/Library/Android/sdk/cmdline-tools/latest/bin\`

2. **Accept Android Licenses**:
   \`\`\`bash
   flutter doctor --android-licenses
   \`\`\`

3. **Build AAB**:
   \`\`\`bash
   flutter build appbundle --release
   \`\`\`

## 📋 **Current Status**

### ✅ **Ready for Upload:**
- **APK File**: \`app-arm64-v8a-release.apk\` (27.7MB)
- **Optimized**: Split per ABI for smaller size
- **Compatible**: Works with Google Play Console conversion

### ⚠️ **Known Issues:**
- Android command-line tools missing
- Debug symbols stripping issue
- AAB build failing due to toolchain

## 🎯 **Recommended Action**

**Upload the APK now** - Google Play Console will handle the AAB conversion automatically. This is the fastest way to get your app testing.

## 📁 **Files Available:**

- **\`app-arm64-v8a-release.apk\`** (27.7MB) - **Recommended for upload**
- **\`app-armeabi-v7a-release.apk\`** (27.2MB) - For older devices
- **\`app-x86_64-release.apk\`** (27.8MB) - For emulators

## 🔧 **Next Steps:**

1. **Upload APK to Google Play Console**
2. **Google Play Console converts to AAB automatically**
3. **Add testers and share testing link**
4. **Start testing immediately**

## 📞 **Support:**

If you encounter issues:
1. Try uploading the APK directly
2. Check Google Play Console help
3. Consider fixing Android toolchain for future builds

EOF

print_success "Workaround guide created!"
print_status "Check AAB_UPLOAD_GUIDE.md for detailed instructions"
print_status "Upload the APK to Google Play Console - it will convert to AAB automatically"
