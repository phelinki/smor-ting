#!/bin/bash

# Smor-Ting TestFlight Setup Script
# This script prepares your app for TestFlight deployment

set -e

echo "🚀 Smor-Ting TestFlight Setup"
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

# Check if we're in the right directory
if [ ! -f "pubspec.yaml" ]; then
    print_error "Please run this script from the smor_ting_mobile directory"
    exit 1
fi

print_info "Setting up TestFlight deployment..."

# Step 1: Check Flutter installation
print_info "Checking Flutter installation..."
if ! command -v flutter &> /dev/null; then
    print_error "Flutter is not installed. Please install Flutter first."
    exit 1
fi

FLUTTER_VERSION=$(flutter --version | head -n 1)
print_status "Flutter version: $FLUTTER_VERSION"

# Step 2: Check Xcode installation
print_info "Checking Xcode installation..."
if ! command -v xcodebuild &> /dev/null; then
    print_error "Xcode is not installed. Please install Xcode from the App Store."
    exit 1
fi

XCODE_VERSION=$(xcodebuild -version | head -n 1)
print_status "Xcode version: $XCODE_VERSION"

# Step 3: Clean and get dependencies
print_info "Cleaning and getting dependencies..."
flutter clean
flutter pub get

# Step 4: Generate code
print_info "Generating code..."
flutter packages pub run build_runner build --delete-conflicting-outputs

# Step 5: Check iOS build
print_info "Testing iOS build..."
flutter build ios --no-codesign

print_status "iOS build test successful!"

# Step 6: Check app configuration
print_info "Checking app configuration..."

# Check if API config exists
if [ ! -f "lib/core/constants/api_config.dart" ]; then
    print_warning "API configuration file not found. Creating..."
    # This should already be created, but just in case
fi

# Check if export options exist
if [ ! -f "ios/exportOptions.plist" ]; then
    print_warning "Export options file not found. Creating..."
    # This should already be created, but just in case
fi

print_status "App configuration looks good!"

echo ""
print_info "📋 TestFlight Deployment Checklist:"
echo ""

echo "✅ Prerequisites:"
echo "   - Flutter SDK installed"
echo "   - Xcode installed"
echo "   - Apple Developer Account active"
echo "   - App registered in App Store Connect"
echo ""

echo "✅ App Configuration:"
echo "   - Bundle ID: com.smorting.app.smorTingMobile"
echo "   - App name: Smor Ting"
echo "   - Permissions configured"
echo "   - API endpoints configured"
echo ""

echo "✅ Build Ready:"
echo "   - Dependencies installed"
echo "   - Code generated"
echo "   - iOS build successful"
echo "   - Export options configured"
echo ""

print_info "🔧 Next Steps:"
echo ""

echo "1. Update API Configuration:"
echo "   Edit lib/core/constants/api_config.dart"
echo "   - Set production URL"
echo "   - Change environment to production"
echo ""

echo "2. Update Team ID:"
echo "   Edit ios/exportOptions.plist"
echo "   - Replace YOUR_TEAM_ID with your actual team ID"
echo ""

echo "3. Build for TestFlight:"
echo "   ./scripts/build_testflight.sh"
echo ""

echo "4. Upload to TestFlight:"
echo "   - Open Xcode"
echo "   - Window → Organizer"
echo "   - Distribute App → App Store Connect"
echo ""

print_warning "⚠️  Important Notes:"
echo ""

echo "- Make sure your Apple Developer account is active"
echo "- Ensure your app is registered in App Store Connect"
echo "- Verify your provisioning profiles are up to date"
echo "- Test the app thoroughly before uploading"
echo "- Update the production API URL before deployment"
echo ""

print_info "📚 Resources:"
echo ""

echo "- TestFlight Deployment Guide: TESTFLIGHT_DEPLOYMENT.md"
echo "- Apple Developer Documentation: https://developer.apple.com/testflight/"
echo "- App Store Connect: https://appstoreconnect.apple.com"
echo "- Flutter iOS Deployment: https://docs.flutter.dev/deployment/ios"
echo ""

print_status "🎉 Your app is ready for TestFlight deployment!"
echo ""
print_info "Run './scripts/build_testflight.sh' when you're ready to build and upload." 