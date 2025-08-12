#!/bin/bash

# Smor-Ting TestFlight Deployment Script
# This script builds and uploads your app to TestFlight via CLI

set -e

echo "🚀 Smor-Ting TestFlight Deployment"
echo "=================================="

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

# Check if we're in the right directory
if [ ! -f "pubspec.yaml" ]; then
    print_error "Please run this script from the smor_ting_mobile directory"
    exit 1
fi

# Get current version from pubspec.yaml
CURRENT_VERSION=$(grep '^version:' pubspec.yaml | sed 's/version: //')
IFS='+' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
VERSION_NAME=${VERSION_PARTS[0]}
BUILD_NUMBER=${VERSION_PARTS[1]}
NEW_BUILD_NUMBER=$((BUILD_NUMBER + 1))
NEW_VERSION="$VERSION_NAME+$NEW_BUILD_NUMBER"

print_info "Current version: $CURRENT_VERSION"
print_info "New version: $NEW_VERSION"

# Step 1: Update build number
print_info "Updating build number to $NEW_BUILD_NUMBER..."
sed -i '' "s/version: .*/version: $NEW_VERSION/" pubspec.yaml

# Step 2: Clean and get dependencies
print_info "Cleaning previous builds..."
flutter clean
flutter pub get

# Step 3: Build IPA
print_info "Building iOS app for release..."
flutter build ipa --release

# Step 4: Upload to TestFlight
print_info "Uploading to TestFlight..."
cd ios
xcodebuild -exportArchive \
    -archivePath ../build/ios/archive/Runner.xcarchive \
    -exportPath ../build/ios/export \
    -exportOptionsPlist exportOptions.plist \
    -allowProvisioningUpdates

cd ..

print_status "🎉 Successfully uploaded to TestFlight!"
echo ""
print_info "Build details:"
echo "📱 Version: $NEW_VERSION"
echo "📁 Archive: build/ios/archive/Runner.xcarchive"
echo "📦 IPA: build/ios/ipa/smor_ting_mobile.ipa"
echo ""
print_info "The build will appear in TestFlight within 5-10 minutes."
print_info "You can now distribute it to internal or external testers."
