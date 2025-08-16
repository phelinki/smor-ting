#!/bin/bash

# Smor-Ting Google Play Store Setup Script
# This script helps set up the environment for Google Play Store submission

set -e

echo "🔧 Setting up Smor-Ting for Google Play Store..."

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

print_status "Checking prerequisites..."

# Check Flutter installation
if command -v flutter &> /dev/null; then
    FLUTTER_VERSION=$(flutter --version | head -n 1)
    print_success "Flutter found: $FLUTTER_VERSION"
else
    print_error "Flutter not found. Please install Flutter first."
    exit 1
fi

# Check Java installation
if command -v java &> /dev/null; then
    JAVA_VERSION=$(java -version 2>&1 | head -n 1)
    print_success "Java found: $JAVA_VERSION"
else
    print_error "Java not found. Please install Java 11 or higher."
    exit 1
fi

# Check Android SDK
if [ -n "$ANDROID_HOME" ]; then
    print_success "Android SDK found at: $ANDROID_HOME"
else
    print_warning "ANDROID_HOME not set. Please set it to your Android SDK path."
fi

# Check if key.properties exists
if [ -f "android/key.properties" ]; then
    print_success "key.properties found"
else
    print_warning "key.properties not found. Creating template..."
    cat > android/key.properties << EOF
# Add your keystore details here
storePassword=your_keystore_password
keyPassword=your_key_password
keyAlias=upload
storeFile=keystore/upload-keystore.jks
EOF
    print_status "Please update android/key.properties with your actual keystore details"
fi

# Check if keystore exists
if [ -f "android/app/keystore/upload-keystore.jks" ]; then
    print_success "Keystore found"
else
    print_warning "Keystore not found. Will be created during build process."
fi

# Create necessary directories
print_status "Creating necessary directories..."
mkdir -p android/app/keystore
mkdir -p build/app/outputs/bundle/release
mkdir -p build/app/outputs/flutter-apk

# Check app configuration
print_status "Checking app configuration..."

# Check if app name is set correctly
if grep -q "Smor-Ting" android/app/src/main/AndroidManifest.xml; then
    print_success "App name configured correctly"
else
    print_warning "App name may need updating in AndroidManifest.xml"
fi

# Check if permissions are properly declared
if grep -q "android.permission.INTERNET" android/app/src/main/AndroidManifest.xml; then
    print_success "Internet permission declared"
else
    print_warning "Internet permission not found in AndroidManifest.xml"
fi

# Check pubspec.yaml version
VERSION=$(grep "^version:" pubspec.yaml | cut -d' ' -f2)
print_status "Current app version: $VERSION"

# Create Google Play Console checklist
cat > GOOGLE_PLAY_CHECKLIST.md << EOF
# Google Play Console Setup Checklist

## Account Setup ✅
- [ ] Google Play Console account created (pkaleewoun@gmail.com)
- [ ] $25 developer registration fee paid
- [ ] Developer agreement accepted

## App Configuration ✅
- [ ] App name: Smor-Ting
- [ ] Package name: com.smorting.app.smor_ting_mobile
- [ ] Version: $VERSION
- [ ] Keystore configured
- [ ] App signing enabled

## Store Listing Requirements
- [ ] App title (80 characters max)
- [ ] Short description (80 characters max)
- [ ] Full description (4000 characters max)
- [ ] App icon (512x512 px)
- [ ] Feature graphic (1024x500 px)
- [ ] Screenshots (minimum 2, max 8)
- [ ] Content rating questionnaire completed
- [ ] Privacy policy URL (if required)

## Technical Requirements
- [ ] App bundle (.aab) file ready
- [ ] App targets API level 21 or higher
- [ ] App doesn't crash on startup
- [ ] App works offline (if claimed)
- [ ] All permissions properly declared
- [ ] No malware or harmful content

## Legal Requirements
- [ ] App complies with Google Play policies
- [ ] App doesn't violate intellectual property rights
- [ ] App doesn't contain inappropriate content
- [ ] App respects user privacy

## Upload Process
1. Go to https://play.google.com/console
2. Sign in with pkaleewoun@gmail.com
3. Create new app (if not exists)
4. Upload app bundle (.aab file)
5. Complete store listing
6. Submit for review

## Important Notes
- Review process takes 1-3 days
- Monitor email for Google notifications
- Keep keystore file secure
- Backup keystore and passwords

EOF

print_success "Setup completed!"
print_status "Check GOOGLE_PLAY_CHECKLIST.md for next steps"
print_status "Run ./scripts/build_play_store.sh to build your app"
