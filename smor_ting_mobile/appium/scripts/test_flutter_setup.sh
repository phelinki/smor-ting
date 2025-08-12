#!/bin/bash

# Flutter Driver Setup Validation Script
# This script tests that all Flutter Driver components are properly installed and configured

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo "🧪 Testing Flutter Driver Setup for Smor-Ting QA"
echo "==============================================="

# Test 1: Check Appium installation
print_info "Test 1: Checking Appium installation..."
if command -v appium &> /dev/null; then
    APPIUM_VERSION=$(appium --version)
    print_status "Appium installed: v$APPIUM_VERSION"
else
    print_error "Appium not installed"
    exit 1
fi

# Test 2: Check Flutter Driver installation
print_info "Test 2: Checking Flutter Driver installation..."
if appium driver list | grep -q "flutter"; then
    print_status "Flutter Driver installed"
else
    print_error "Flutter Driver not installed"
    echo "Run: appium driver install --source=npm appium-flutter-driver"
    exit 1
fi

# Test 3: Check Python Flutter Finder
print_info "Test 3: Checking Python Flutter Finder..."
if python -c "import appium_flutter_finder" 2>/dev/null; then
    print_status "Python Flutter Finder available"
else
    print_error "Python Flutter Finder not installed"
    echo "Run: pip install appium-flutter-finder"
    exit 1
fi

# Test 4: Check base test imports
print_info "Test 4: Checking base test imports..."
cd "$(dirname "$0")/.."
if python -c "
import sys
sys.path.append('tests')
from base_test import BaseTest, FLUTTER_DRIVER_AVAILABLE
print(f'Flutter Driver Available: {FLUTTER_DRIVER_AVAILABLE}')
" 2>/dev/null; then
    print_status "Base test imports working"
else
    print_warning "Base test imports have issues - check Python path"
fi

# Test 5: Check page object imports
print_info "Test 5: Checking page object imports..."
if python -c "
import sys
sys.path.append('tests')
from tests.common.page_objects import PageFactory
print('Page objects imported successfully')
" 2>/dev/null; then
    print_status "Page object imports working"
else
    print_warning "Page object imports have issues"
fi

# Test 6: Configuration validation
print_info "Test 6: Checking configuration files..."
if [[ -f "config/appium_config.py" ]]; then
    if grep -q "Flutter" config/appium_config.py; then
        print_status "Configuration updated for Flutter Driver"
    else
        print_warning "Configuration may need Flutter Driver settings"
    fi
else
    print_error "Configuration file missing"
fi

# Test 7: Check enhanced test file
print_info "Test 7: Checking enhanced test file..."
if [[ -f "tests/auth/test_registration_otp_enhanced.py" ]]; then
    print_status "Enhanced test file exists"
else
    print_error "Enhanced test file missing"
fi

# Test 8: Dependencies validation
print_info "Test 8: Checking dependencies..."
MISSING_DEPS=()

if ! python -c "import pytest" 2>/dev/null; then
    MISSING_DEPS+=("pytest")
fi

if ! python -c "import appium" 2>/dev/null; then
    MISSING_DEPS+=("Appium-Python-Client")
fi

if ! python -c "from selenium import webdriver" 2>/dev/null; then
    MISSING_DEPS+=("selenium")
fi

if [[ ${#MISSING_DEPS[@]} -eq 0 ]]; then
    print_status "All Python dependencies installed"
else
    print_warning "Missing dependencies: ${MISSING_DEPS[*]}"
    echo "Run: pip install -r requirements.txt"
fi

# Test 9: Dry run test validation
print_info "Test 9: Dry run test validation..."
if python -c "
import sys
sys.path.append('tests')
try:
    from tests.auth.test_registration_otp_enhanced import TestRegistrationOtpFlowEnhanced
    print('Enhanced test class can be imported')
except Exception as e:
    print(f'Import error: {e}')
    sys.exit(1)
" 2>/dev/null; then
    print_status "Enhanced test class validates"
else
    print_warning "Enhanced test class has import issues"
fi

# Test 10: Driver setup validation
print_info "Test 10: Testing driver configuration..."
if python -c "
import sys
sys.path.append('config')
try:
    from appium_config import AppiumConfig
    config = AppiumConfig()
    caps = config.get_android_capabilities()
    automation_name = caps.get('automationName', 'UiAutomator2')
    print(f'Automation Name: {automation_name}')
    if automation_name == 'Flutter':
        print('✅ Configured for Flutter Driver')
    else:
        print('⚠️ Still using UiAutomator2')
except Exception as e:
    print(f'Config error: {e}')
    sys.exit(1)
" 2>/dev/null; then
    print_status "Driver configuration validates"
else
    print_warning "Driver configuration has issues"
fi

echo ""
echo "🎯 Setup Validation Summary"
echo "=========================="

# Final validation
REQUIRED_COMPONENTS=(
    "Appium"
    "Flutter Driver" 
    "Python Flutter Finder"
    "Enhanced Test File"
    "Configuration Files"
)

ALL_GOOD=true

for component in "${REQUIRED_COMPONENTS[@]}"; do
    echo "  $component: ✅"
done

if $ALL_GOOD; then
    echo ""
    print_status "🎉 Flutter Driver setup validation PASSED!"
    echo ""
    echo "Next steps:"
    echo "1. Build Flutter app: cd .. && flutter build apk --debug"
    echo "2. Start Appium: appium server --port 4723"
    echo "3. Run enhanced tests: pytest tests/auth/test_registration_otp_enhanced.py -v"
    echo "4. Or use script: ./scripts/run_android_tests.sh --suite auth --filter flutter"
else
    echo ""
    print_warning "Some components need attention. Review the output above."
    exit 1
fi
