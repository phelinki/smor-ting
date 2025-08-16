#!/bin/bash

# Google Play Console Setup Diagnostic Script (Updated for 2024 Interface)

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

echo "🔍 Google Play Console Setup Diagnostic (2024 Interface)"
echo "======================================================"

print_status "Checking your current setup..."

# Check if service account file exists
if [ -f "scripts/play_store/service-account.json" ]; then
    print_success "Service account file found"
else
    print_warning "Service account file not found"
    print_status "You'll need to create one after finding API access"
fi

echo ""
echo "📋 Current Google Play Console Interface (2024):"
echo "================================================"
echo ""
echo "1. Go to: https://play.google.com/console"
echo "2. Sign in with: pkaleewoun@gmail.com"
echo "3. Check if you see 'Smor-Ting' in the app list"
echo ""
echo "If NO app exists:"
echo "  - Click 'Create app'"
echo "  - Name: Smor-Ting"
echo "  - Language: English"
echo "  - Type: App"
echo "  - Price: Free"
echo ""
echo "🔍 CURRENT API ACCESS LOCATION (2024):"
echo "======================================"
echo ""
echo "After creating/finding your app, look for:"
echo ""
echo "Method 1: Test & Release Section"
echo "  - Click 'Test and release' in left menu"
echo "  - Look for 'Advanced settings' or 'API access'"
echo "  - Check for 'Service accounts' or 'API credentials'"
echo ""
echo "Method 2: Settings Section"
echo "  - Look for 'Settings' or 'Configuration' in left menu"
echo "  - Check for 'API access' or 'Service accounts'"
echo ""
echo "Method 3: Account Level"
echo "  - Click on your account/profile in top right"
echo "  - Look for 'API access' or 'Service accounts'"
echo ""
echo "Method 4: Search Function"
echo "  - Use the search bar in Play Console"
echo "  - Search for 'API' or 'service account'"
echo ""
echo "Method 5: App-Specific Settings"
echo "  - Click on your app name (Smor-Ting)"
echo "  - Look for 'Settings' or 'Configuration'"
echo "  - Check for 'API access' or 'Service accounts'"
echo ""

print_status "Common reasons API access is not visible (2024):"
echo "1. No app created yet"
echo "2. Not the account owner or admin"
echo "3. Account not fully verified"
echo "4. Need to complete initial app setup"
echo "5. Interface changes - API access moved to different location"
echo ""

print_status "Alternative approach if API access is not found:"
echo "1. Use Google Cloud Console directly"
echo "2. Enable Google Play Developer API in your project"
echo "3. Create service account in Google Cloud Console"
echo "4. Grant Play Console permissions manually"
echo ""

print_status "Next steps after finding API access:"
echo "1. Link your Google Cloud project"
echo "2. Create service account"
echo "3. Grant 'Release apps to testing tracks' permission"
echo "4. Download service account JSON"
echo "5. Place it in: scripts/play_store/service-account.json"
echo ""

print_success "Run this script again after completing the setup to verify everything is working"
