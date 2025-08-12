#!/bin/bash

# Flutter Driver Setup Script for Smor-Ting QA Automation
# This script installs and configures Appium Flutter Driver for improved element discovery

set -e

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

print_step() {
    echo -e "\n${BLUE}==== $1 ====${NC}"
}

echo "🚀 Flutter Driver Setup for Smor-Ting QA Automation"
echo "=================================================="

# Check if Appium is installed
print_step "Step 1: Checking Appium installation"
if ! command -v appium &> /dev/null; then
    print_error "Appium is not installed. Please run setup_appium.sh first."
    exit 1
fi
print_status "Appium is installed"

# Check Appium version
APPIUM_VERSION=$(appium --version)
print_info "Appium version: $APPIUM_VERSION"

# Install Flutter Driver
print_step "Step 2: Installing Appium Flutter Driver"
print_info "Installing Flutter Driver from npm..."
appium driver install --source=npm appium-flutter-driver

# Verify Flutter Driver installation
print_step "Step 3: Verifying Flutter Driver installation"
if appium driver list | grep -q "flutter"; then
    print_status "Flutter Driver installed successfully"
else
    print_error "Flutter Driver installation failed"
    exit 1
fi

# Install Python Flutter Finder library
print_step "Step 4: Installing Python Flutter Finder library"
print_info "Installing appium-flutter-finder..."
pip install appium-flutter-finder==0.2.0

# Verify Python library installation
print_step "Step 5: Verifying Python Flutter Finder installation"
if python -c "import appium_flutter_finder; print('Flutter Finder available')" 2>/dev/null; then
    print_status "Python Flutter Finder installed successfully"
else
    print_error "Python Flutter Finder installation failed"
    exit 1
fi

# Run Flutter Driver doctor
print_step "Step 6: Running Flutter Driver diagnostics"
print_info "Running Flutter Driver doctor..."
appium driver doctor flutter || print_warning "Flutter Driver doctor completed with warnings (this is normal)"

# Test basic Flutter Driver functionality
print_step "Step 7: Testing Flutter Driver setup"
print_info "Creating test configuration..."

cat > /tmp/flutter_driver_test.py << 'EOF'
#!/usr/bin/env python3
"""
Quick test to verify Flutter Driver setup
"""
import sys

try:
    from appium_flutter_finder.flutter_finder import FlutterFinder
    print("✅ FlutterFinder import successful")
    
    finder = FlutterFinder()
    print("✅ FlutterFinder instance created")
    
    # Test basic locator creation
    locator = finder.by_value_key("test_key")
    print("✅ Flutter locator creation successful")
    
    print("🎉 Flutter Driver setup verification PASSED")
    sys.exit(0)
    
except ImportError as e:
    print(f"❌ Import error: {e}")
    sys.exit(1)
except Exception as e:
    print(f"❌ Setup verification failed: {e}")
    sys.exit(1)
EOF

python /tmp/flutter_driver_test.py
rm /tmp/flutter_driver_test.py

# Display configuration summary
print_step "Setup Summary"
echo "Installed components:"
print_status "✅ Appium Flutter Driver (npm package)"
print_status "✅ Python appium-flutter-finder library"
print_status "✅ Flutter locator capabilities"

echo ""
print_info "Flutter Driver is now ready for use in QA automation!"
print_info "Tests can now use Flutter-first element discovery with UiAutomator2 fallbacks."

echo ""
echo "Next steps:"
echo "1. Run the enhanced registration tests: ./scripts/run_android_tests.sh --suite auth --filter flutter"
echo "2. Check test results in: reports/android-report.html"
echo "3. Review Flutter Driver usage in: tests/auth/test_registration_otp_enhanced.py"

print_status "Flutter Driver setup completed successfully! 🎉"
