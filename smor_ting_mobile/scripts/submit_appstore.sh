#!/bin/bash

# Smor-Ting App Store Submission Script
# This script guides you through the complete App Store submission process

set -e

echo "🍎 Smor-Ting App Store Submission"
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

print_header() {
    echo -e "${BLUE}📋 $1${NC}"
    echo "=================================="
}

# Configuration
APP_NAME="Smor Ting"
BUNDLE_ID="com.smorting.app.smorTingMobile"
VERSION="1.0.0"
BUILD_NUMBER="1"

print_header "App Store Submission Checklist"

echo ""
print_info "📋 Pre-Submission Requirements:"
echo ""

echo "✅ Technical Requirements:"
echo "   - App passes all App Store Review Guidelines"
echo "   - No crashes on launch"
echo "   - All features work as expected"
echo "   - Performance is acceptable"
echo "   - Privacy policy is implemented"
echo "   - Terms of service are included"
echo ""

echo "✅ Content Requirements:"
echo "   - App icon is high quality (1024x1024)"
echo "   - Screenshots are up to date"
echo "   - App description is complete"
echo "   - Keywords are optimized"
echo "   - Support URL is working"
echo ""

echo "✅ Legal Requirements:"
echo "   - Privacy policy is accessible"
echo "   - Terms of service are included"
echo "   - Data collection is disclosed"
echo "   - Third-party libraries are listed"
echo ""

print_header "Step-by-Step Submission Process"

echo ""
print_info "Step 1: Prepare Your App"
echo ""

echo "1. Update Production API URL:"
echo "   Edit lib/core/constants/api_config.dart"
echo "   - Set _productionBaseUrl to your actual server"
echo "   - Ensure _currentEnvironment is set to production"
echo ""

echo "2. Update App Version:"
echo "   Edit pubspec.yaml"
echo "   - Update version to $VERSION+$BUILD_NUMBER"
echo ""

echo "3. Build Production App:"
echo "   ./scripts/build_testflight.sh"
echo ""

print_info "Step 2: App Store Connect Setup"
echo ""

echo "1. Go to App Store Connect:"
echo "   https://appstoreconnect.apple.com"
echo ""

echo "2. Create New App:"
echo "   - Click '+' → New App"
echo "   - Platform: iOS
echo "   - Name: Smor Ting"
echo "   - Bundle ID: $BUNDLE_ID"
echo "   - SKU: smor-ting-ios"
echo "   - User Access: Full Access"
echo ""

echo "3. Fill App Information:"
echo "   - App Name: Smor Ting"
echo "   - Subtitle: Handyman & Service Marketplace"
echo "   - Keywords: handyman,services,marketplace,liberia"
echo "   - Description: [Your app description]"
echo "   - Support URL: [Your support website]"
echo "   - Marketing URL: [Your marketing website]"
echo ""

print_info "Step 3: Upload Build"
echo ""

echo "1. Open Xcode:"
echo "   open ios/Runner.xcworkspace"
echo ""

echo "2. Archive App:"
echo "   - Product → Archive"
echo "   - Wait for archive to complete"
echo ""

echo "3. Upload to App Store Connect:"
echo "   Option A: Command-line (recommended for CI/non-interactive)"
echo "     export ASC_KEY_PATH=</path/to/AuthKey_XXXXXX.p8>"
echo "     export ASC_KEY_ID=<KeyID>"
echo "     export ASC_ISSUER_ID=<IssuerID>"
echo "     (cd ios && xcodebuild -exportArchive \\\""
echo "        -archivePath build/Runner.xcarchive \\\""
echo "        -exportPath build/ipa \\\""
echo "        -exportOptionsPlist exportOptions.plist \\\""
echo "        -allowProvisioningUpdates \\\""
echo "        -authenticationKeyPath \"$ASC_KEY_PATH\" \\\""
echo "        -authenticationKeyID \"$ASC_KEY_ID\" \\\""
echo "        -authenticationKeyIssuerID \"$ASC_ISSUER_ID\")"
echo "   Option B: Xcode Organizer (manual)"
echo "     - Window → Organizer → Archives → Distribute App → App Store Connect"
echo ""

print_info "Step 4: App Store Connect Configuration"
echo ""

echo "1. App Information:"
echo "   - Fill out all required fields"
echo "   - Add app description"
echo "   - Add keywords"
echo "   - Add support URL"
echo ""

echo "2. App Screenshots:"
echo "   - iPhone 6.7" Display (1290 x 2796)
echo "   - iPhone 6.5" Display (1242 x 2688)"
echo "   - iPhone 5.5" Display (1242 x 2208)"
echo "   - iPad Pro 12.9" Display (2048 x 2732)"
echo ""

echo "3. App Review Information:"
echo "   - Demo Account (if required)"
echo "   - Demo Password (if required)"
echo "   - Notes for Review"
echo ""

print_info "Step 5: Submit for Review"
echo ""

echo "1. Final Checks:"
echo "   - All required fields completed"
echo "   - Screenshots uploaded"
echo "   - App description complete"
echo "   - Privacy policy accessible"
echo ""

echo "2. Submit for Review:"
echo "   - Click 'Submit for Review'"
echo "   - Confirm submission"
echo "   - Wait for review (1-7 days)"
echo ""

print_header "Required App Store Assets"

echo ""
print_info "App Icon (Required):"
echo "   - Size: 1024 x 1024 pixels"
echo "   - Format: PNG or JPEG"
echo "   - No transparency"
echo "   - No rounded corners"
echo ""

print_info "Screenshots (Required):"
echo "   - Minimum 1 screenshot per device"
echo "   - Maximum 10 screenshots per device"
echo "   - No device frames"
echo "   - No text overlays"
echo ""

print_info "App Description (Required):"
echo "   - Maximum 4000 characters"
echo "   - Clear and compelling"
echo "   - Highlight key features"
echo "   - Include call-to-action"
echo ""

print_header "App Store Review Guidelines"

echo ""
print_warning "Common Rejection Reasons:"
echo ""

echo "❌ App crashes on launch"
echo "❌ Incomplete app functionality"
echo "❌ Missing privacy policy"
echo "❌ Poor app performance"
echo "❌ Inappropriate content"
echo "❌ Misleading app description"
echo "❌ Broken links or URLs"
echo ""

print_info "Best Practices:"
echo ""

echo "✅ Test thoroughly before submission"
echo "✅ Follow Apple's Human Interface Guidelines"
echo "✅ Provide clear app description"
echo "✅ Include privacy policy"
echo "✅ Respond quickly to review feedback"
echo "✅ Keep app updated regularly"
echo ""

print_header "Post-Submission Process"

echo ""
print_info "Review Timeline:"
echo "   - Standard review: 1-7 days"
echo "   - Expedited review: 24-48 hours"
echo "   - Re-review: 1-3 days"
echo ""

print_info "Review Status:"
echo "   - Waiting for Review"
echo "   - In Review"
echo "   - Ready for Sale"
echo "   - Rejected (with feedback)"
echo ""

print_info "If Rejected:"
echo "   1. Read rejection feedback carefully"
echo "   2. Fix the issues mentioned"
echo "   3. Upload new build"
echo "   4. Resubmit for review"
echo ""

print_header "Launch Preparation"

echo ""
print_info "Pre-Launch Checklist:"
echo ""

echo "✅ App approved by Apple"
echo "✅ App Store listing complete"
echo "✅ Marketing materials ready"
echo "✅ Support system in place"
echo "✅ Analytics tracking configured"
echo "✅ Crash reporting set up"
echo "✅ User feedback system ready"
echo ""

print_info "Launch Day:"
echo ""

echo "1. Monitor App Store Connect"
echo "2. Track download metrics"
echo "3. Monitor crash reports"
echo "4. Respond to user reviews"
echo "5. Prepare for updates"
echo ""

print_status "🎉 Your Smor-Ting app is ready for App Store submission!"
echo ""
print_info "Follow this guide step by step for a successful submission."
echo ""
print_warning "Remember: App Store review can take 1-7 days."
echo ""
print_info "Good luck with your App Store submission! 🚀" 