#!/bin/bash

# Smor-Ting Google Play Store Build Script
# This script builds and prepares the app for Google Play Store submission

set -e  # Exit on any error

echo "🚀 Starting Smor-Ting Google Play Store build process..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Check if Flutter is installed
if ! command -v flutter &> /dev/null; then
    print_error "Flutter is not installed or not in PATH"
    exit 1
fi

# Check if key.properties exists
if [ ! -f "android/key.properties" ]; then
    print_warning "key.properties not found. Creating template..."
    cat > android/key.properties << EOF
# Add your keystore details here
storePassword=your_keystore_password
keyPassword=your_key_password
keyAlias=upload
storeFile=keystore/upload-keystore.jks
EOF
    print_error "Please update android/key.properties with your actual keystore details"
    print_error "Then run this script again"
    exit 1
fi

# Check if keystore exists
if [ ! -f "android/app/keystore/upload-keystore.jks" ]; then
    print_warning "Keystore not found. Creating keystore..."
    mkdir -p android/app/keystore
    
    echo "Please provide the following information for your keystore:"
    read -p "Keystore password: " store_password
    read -p "Key password: " key_password
    read -p "Your name: " name
    read -p "Your organization: " org
    read -p "Your city: " city
    read -p "Your state: " state
    read -p "Your country code (e.g., LR): " country
    
    # Update key.properties
    sed -i.bak "s/your_keystore_password/$store_password/g" android/key.properties
    sed -i.bak "s/your_key_password/$key_password/g" android/key.properties
    
    # Generate keystore
    keytool -genkey -v \
        -keystore android/app/keystore/upload-keystore.jks \
        -keyalg RSA \
        -keysize 2048 \
        -validity 10000 \
        -alias upload \
        -storepass "$store_password" \
        -keypass "$key_password" \
        -dname "CN=$name, OU=Development, O=$org, L=$city, ST=$state, C=$country"
    
    print_success "Keystore created successfully"
fi

# Clean previous builds
print_status "Cleaning previous builds..."
flutter clean

# Get dependencies
print_status "Getting dependencies..."
flutter pub get

# Run tests
print_status "Running tests..."
flutter test

# Build app bundle for Play Store (recommended)
print_status "Building app bundle for Google Play Store..."
flutter build appbundle --release

# Build APK as backup
print_status "Building APK as backup..."
flutter build apk --release

# Check if builds were successful
if [ -f "build/app/outputs/bundle/release/app-release.aab" ]; then
    print_success "App bundle created successfully!"
    print_status "Location: build/app/outputs/bundle/release/app-release.aab"
    print_status "Size: $(du -h build/app/outputs/bundle/release/app-release.aab | cut -f1)"
else
    print_error "Failed to create app bundle"
    exit 1
fi

if [ -f "build/app/outputs/flutter-apk/app-release.apk" ]; then
    print_success "APK created successfully!"
    print_status "Location: build/app/outputs/flutter-apk/app-release.apk"
    print_status "Size: $(du -h build/app/outputs/flutter-apk/app-release.apk | cut -f1)"
fi

# Create upload instructions
cat > UPLOAD_INSTRUCTIONS.md << EOF
# Google Play Store Upload Instructions

## Files Ready for Upload:
- **App Bundle (Recommended)**: \`build/app/outputs/bundle/release/app-release.aab\`
- **APK (Backup)**: \`build/app/outputs/flutter-apk/app-release.apk\`

## Next Steps:

1. **Go to Google Play Console**: https://play.google.com/console
2. **Sign in with**: pkaleewoun@gmail.com
3. **Create New App** (if not already created):
   - App name: Smor-Ting
   - Default language: English
   - App or game: App
   - Free or paid: Free
   - Category: Business

4. **Upload App Bundle**:
   - Go to Production track
   - Click "Create new release"
   - Upload the \`app-release.aab\` file
   - Add release notes
   - Save and review

5. **Complete Store Listing**:
   - App details
   - Graphics (icons, screenshots)
   - Content rating
   - Pricing & distribution

6. **Submit for Review**:
   - Review all information
   - Click "Start rollout to Production"

## Important Notes:
- Keep your keystore file safe (\`android/app/keystore/upload-keystore.jks\`)
- Store your keystore passwords securely
- The review process typically takes 1-3 days
- Monitor your email for any issues from Google

## Troubleshooting:
- If upload fails, check that your app meets Google Play policies
- Ensure all required permissions are properly declared
- Verify your app doesn't crash on startup
- Check that your app works offline (if claimed)

EOF

print_success "Build completed successfully!"
print_status "Check UPLOAD_INSTRUCTIONS.md for next steps"
print_status "Your app bundle is ready for upload to Google Play Store"
