#!/bin/bash

echo "📱 Setting up Physical Device Testing"
echo "=================================="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo -e "${GREEN}✅ $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

# Install only ADB (no full SDK needed)
print_info "Installing minimal Android tools via Homebrew..."
brew install android-platform-tools

if [ $? -eq 0 ]; then
    print_status "ADB installed successfully!"
else
    echo "Failed to install ADB. You can also download it manually."
fi

# Update environment variables
echo "" >> ~/.bash_profile
echo "# Minimal Android tools for Appium" >> ~/.bash_profile
echo "export PATH=\$PATH:/opt/homebrew/bin:/usr/local/bin" >> ~/.bash_profile

source ~/.bash_profile

echo ""
print_info "Next steps:"
echo "1. Enable 'Developer Options' on your Android phone:"
echo "   Settings → About Phone → tap 'Build Number' 7 times"
echo ""
echo "2. Enable 'USB Debugging':"
echo "   Settings → Developer Options → USB Debugging → ON"
echo ""
echo "3. Connect your phone via USB"
echo ""
echo "4. Test connection:"
echo "   adb devices"
echo ""
echo "5. Run tests:"
echo "   ./scripts/run_android_tests.sh"

print_status "Physical device setup complete! 📱"
