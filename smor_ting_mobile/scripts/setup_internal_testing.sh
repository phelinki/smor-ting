#!/bin/bash

# Smor-Ting Internal Testing Setup Script
# This script sets up internal testing distribution (Android equivalent to TestFlight)

set -e

echo "🧪 Setting up Smor-Ting Internal Testing (Android TestFlight equivalent)..."

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

# Check if we're in the right directory
if [ ! -f "pubspec.yaml" ]; then
    print_error "Please run this script from the smor_ting_mobile directory"
    exit 1
fi

print_status "Building app for internal testing..."

# Build APK for testing (more reliable than app bundle for testing)
print_status "Building APK for direct installation..."
flutter build apk --release

# Check if build was successful
if [ -f "build/app/outputs/flutter-apk/app-release.apk" ]; then
    print_success "APK created successfully!"
    print_status "Location: build/app/outputs/flutter-apk/app-release.apk"
    print_status "Size: $(du -h build/app/outputs/flutter-apk/app-release.apk | cut -f1)"
else
    print_error "Failed to create APK"
    exit 1
fi

# Create internal testing guide
cat > INTERNAL_TESTING_GUIDE.md << EOF
# Smor-Ting Internal Testing Setup Guide

## Android Equivalent to TestFlight

Google Play Console offers **Internal Testing** and **Closed Testing** tracks that work similarly to iOS TestFlight.

## Files Ready for Upload:
- **APK**: \`build/app/outputs/flutter-apk/app-release.apk\`

## Step-by-Step Setup:

### 1. Access Google Play Console
- Go to [play.google.com/console](https://play.google.com/console)
- Sign in with: **pkaleewoun@gmail.com**

### 2. Create App (if not exists)
- Click **"Create app"**
- App name: **Smor-Ting**
- Default language: **English**
- App or game: **App**
- Free or paid: **Free**
- Click **"Create"**

### 3. Set Up Internal Testing Track
1. In the left menu, click **"Testing"**
2. Click **"Internal testing"**
3. Click **"Create new release"**
4. Upload your APK: \`app-release.apk\`
5. Add release notes:
   \`\`\`
   Internal testing build for Smor-Ting
   
   Features:
   - User authentication
   - Service booking
   - Location tracking
   - Payment integration
   - Offline functionality
   \`\`\`
6. Click **"Save"**

### 4. Add Testers
1. In **"Internal testing"**, click **"Testers"**
2. Click **"Create email list"**
3. Add tester emails (one per line):
   \`\`\`
   tester1@example.com
   tester2@example.com
   pkaleewoun@gmail.com
   \`\`\`
4. Click **"Save"**

### 5. Get Testing Link
1. Click **"Get link"** in the Internal testing section
2. Copy the testing URL
3. Share this URL with your testers

## Alternative: Direct APK Distribution
For immediate testing, you can also:
1. **Share the APK file directly** with testers
2. **Testers enable "Install from unknown sources"** in Android settings
3. **Install the APK directly** on their devices

## Tester Instructions:
1. **Receive testing link** from you
2. **Click the link** on their Android device
3. **Accept the invitation** to become a tester
4. **Wait 10-15 minutes** for Google to process
5. **Download the app** from the testing link
6. **Install and test** the app

## Benefits of Internal Testing:
- ✅ **No review process** (instant availability)
- ✅ **Up to 100 testers** per track
- ✅ **Easy updates** (just upload new APK)
- ✅ **Crash reporting** and analytics
- ✅ **Feedback collection** from testers
- ✅ **No public visibility** (only testers can see)

## Testing Checklist:
- [ ] App installs successfully
- [ ] App launches without crashes
- [ ] All features work as expected
- [ ] Offline functionality works
- [ ] Payment integration works
- [ ] Location services work
- [ ] User registration/login works
- [ ] Service booking flow works

## Troubleshooting:
- **Tester can't install**: Wait 10-15 minutes after invitation
- **App crashes**: Check crash reports in Play Console
- **Features not working**: Test on different devices
- **Upload fails**: Ensure proper signing configuration

## Next Steps:
1. Set up internal testing track
2. Add testers
3. Share testing link
4. Collect feedback
5. Fix issues
6. Upload updates as needed
7. When ready, move to production

## Important Notes:
- Internal testing builds are **not reviewed** by Google
- Testers need to **accept invitation** before downloading
- Updates are **instant** (no waiting period)
- You can have **multiple testing tracks** (Internal, Alpha, Beta)
- **No cost** for testing (only $25 one-time developer fee)

EOF

print_success "Internal testing setup completed!"
print_status "Check INTERNAL_TESTING_GUIDE.md for detailed instructions"
print_status "Your APK is ready for internal testing distribution"
print_status "Upload the .apk file to Google Play Console Internal Testing track"
