#!/bin/bash

# Simple Smor-Ting Build Script for Testing
# This script builds a basic APK for internal testing

set -e

echo "🔨 Building Smor-Ting APK for testing..."

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

print_status "Cleaning previous builds..."
flutter clean

print_status "Getting dependencies..."
flutter pub get

print_status "Building APK for testing..."
flutter build apk

# Check if build was successful
if [ -f "build/app/outputs/flutter-apk/app-release.apk" ]; then
    print_success "Release APK created successfully!"
    print_status "Location: build/app/outputs/flutter-apk/app-release.apk"
    print_status "Size: $(du -h build/app/outputs/flutter-apk/app-release.apk | cut -f1)"
    
    # Also check for debug APK
    if [ -f "build/app/outputs/flutter-apk/app-debug.apk" ]; then
        print_success "Debug APK also created!"
        print_status "Location: build/app/outputs/flutter-apk/app-debug.apk"
        print_status "Size: $(du -h build/app/outputs/flutter-apk/app-debug.apk | cut -f1)"
    fi
else
    print_error "Failed to create APK"
    exit 1
fi

# Create simple testing guide
cat > SIMPLE_TESTING_GUIDE.md << EOF
# Simple Smor-Ting Testing Guide

## APK Files Ready for Testing:
- **Release APK**: \`build/app/outputs/flutter-apk/app-release.apk\`
- **Debug APK**: \`build/app/outputs/flutter-apk/app-debug.apk\` (if available)

## Quick Testing Options:

### Option 1: Direct APK Installation
1. **Share the APK file** directly with testers
2. **Testers enable "Install from unknown sources"** in Android settings
3. **Install the APK directly** on their devices

### Option 2: Google Play Console Internal Testing
1. Go to [play.google.com/console](https://play.google.com/console)
2. Sign in with **pkaleewoun@gmail.com**
3. Create app: "Smor-Ting"
4. Go to **Testing** → **Internal testing**
5. Upload the APK file
6. Add testers and share the testing link

### Option 3: Firebase App Distribution (Alternative)
1. Set up Firebase project
2. Upload APK to Firebase App Distribution
3. Invite testers via email
4. Testers get email with download link

## Tester Instructions for Direct APK:
1. **Enable unknown sources**:
   - Settings → Security → Unknown sources
   - Or Settings → Apps → Special app access → Install unknown apps
2. **Download the APK** file
3. **Tap the APK file** to install
4. **Follow installation prompts**
5. **Launch and test** the app

## Testing Checklist:
- [ ] App installs successfully
- [ ] App launches without crashes
- [ ] All features work as expected
- [ ] Offline functionality works
- [ ] Payment integration works
- [ ] Location services work
- [ ] User registration/login works
- [ ] Service booking flow works

## Next Steps:
1. Test the APK thoroughly
2. Fix any issues found
3. Build updated APK
4. Upload to Google Play Console for internal testing
5. Share with testers

## Important Notes:
- Debug APK is larger but easier to debug
- Release APK is optimized but may have issues
- Always test on multiple devices
- Keep track of issues and feedback

EOF

print_success "Build completed successfully!"
print_status "Check SIMPLE_TESTING_GUIDE.md for testing instructions"
print_status "Your APK is ready for testing!"
