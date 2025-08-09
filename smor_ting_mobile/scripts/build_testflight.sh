#!/bin/bash

# Smor-Ting TestFlight Build Script
# This script builds and prepares your app for TestFlight deployment

set -e

echo "🚀 Smor-Ting TestFlight Build"
echo "=============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Configuration
APP_NAME="Smor Ting"
BUNDLE_ID="com.smorting.app.smorTingMobile"
VERSION="1.0.0"
BUILD_NUMBER="1"

# Check if we're in the right directory
if [ ! -f "pubspec.yaml" ]; then
    print_error "Please run this script from the smor_ting_mobile directory"
    exit 1
fi

print_info "Starting TestFlight build process..."

# Step 1: Clean previous builds
print_info "Cleaning previous builds..."
flutter clean
flutter pub get

# Step 2: Update version and build number
print_info "Updating version to $VERSION+$BUILD_NUMBER..."
sed -i '' "s/version: .*/version: $VERSION+$BUILD_NUMBER/" pubspec.yaml

# Step 3: Generate code
print_info "Generating code..."
flutter packages pub run build_runner build --delete-conflicting-outputs

# Step 4: Build iOS app
print_info "Building iOS app for release..."
flutter build ios --release --no-codesign

# Step 5: Check if Xcode is available
if ! command -v xcodebuild &> /dev/null; then
    print_error "Xcode is not installed or not in PATH"
    print_info "Please install Xcode from the App Store"
    exit 1
fi

# Step 6: Archive the app
print_info "Creating archive..."
cd ios

# Optional App Store Connect API Key auth
ASC_FLAGS=()
if [ -n "$ASC_KEY_PATH" ] && [ -n "$ASC_KEY_ID" ] && [ -n "$ASC_ISSUER_ID" ]; then
    print_info "Using App Store Connect API key authentication"
    ASC_FLAGS=(
        -authenticationKeyPath "$ASC_KEY_PATH"
        -authenticationKeyID "$ASC_KEY_ID"
        -authenticationKeyIssuerID "$ASC_ISSUER_ID"
    )
else
    print_warning "No ASC API key provided (ASC_KEY_PATH, ASC_KEY_ID, ASC_ISSUER_ID). Will use Xcode account if available."
fi
xcodebuild -workspace Runner.xcworkspace \
    -scheme Runner \
    -configuration Release \
    -destination generic/platform=iOS \
    -archivePath build/Runner.xcarchive \
    clean archive

# Step 7: Export IPA
print_info "Exporting IPA..."
xcodebuild -exportArchive \
    -archivePath build/Runner.xcarchive \
    -exportPath build/ipa \
    -exportOptionsPlist exportOptions.plist \
    -allowProvisioningUpdates \
    ${ASC_FLAGS[@]}

print_status "Build completed successfully!"
echo ""

print_info "Next steps for TestFlight deployment:"
echo "1. Open Xcode"
echo "2. Go to Window → Organizer"
echo "3. Select your archive"
echo "4. Click 'Distribute App'"
echo "5. Choose 'App Store Connect'"
echo "6. Follow the distribution wizard"
echo ""

print_info "Archive location:"
echo "📁 ios/build/Runner.xcarchive"
echo ""

print_info "IPA location:"
echo "📁 ios/build/ipa/Runner.ipa"
echo ""

print_warning "Important notes:"
echo "- Make sure your Apple Developer account is active"
echo "- Ensure your app is registered in App Store Connect"
echo "- Verify your provisioning profiles are up to date"
echo "- Test the app thoroughly before uploading to TestFlight"
echo ""

print_status "🎉 Your app is ready for TestFlight deployment!" 