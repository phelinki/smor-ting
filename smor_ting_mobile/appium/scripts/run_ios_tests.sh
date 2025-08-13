#!/bin/bash

echo "🍎 Running iOS Tests for Smor-Ting"
echo "================================="

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
    echo -e "${BLUE}🔧 $1${NC}"
}

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    print_error "iOS testing requires macOS"
    exit 1
fi

# Parse command line arguments
TEST_SUITE="all"
ENVIRONMENT="local"
PARALLEL=""
MARKERS=""
DEVICE_NAME="iPhone 13"
IOS_VERSION="16.4"

while [[ $# -gt 0 ]]; do
    case $1 in
        --suite)
            TEST_SUITE="$2"
            shift 2
            ;;
        --environment)
            ENVIRONMENT="$2" 
            shift 2
            ;;
        --parallel)
            PARALLEL="--numprocesses auto"
            shift
            ;;
        --markers)
            MARKERS="-m $2"
            shift 2
            ;;
        --device)
            DEVICE_NAME="$2"
            shift 2
            ;;
        --ios-version)
            IOS_VERSION="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --suite <all|auth|registration|login>  Test suite to run (default: all)"
            echo "  --environment <local|ci|staging>       Test environment (default: local)"
            echo "  --parallel                              Run tests in parallel"
            echo "  --markers <marker>                      Run tests with specific markers"
            echo "  --device <device_name>                  iOS device to use (default: iPhone 13)"
            echo "  --ios-version <version>                 iOS version (default: 16.4)"
            echo "  --help                                  Show this help message"
            exit 0
            ;;
        *)
            print_warning "Unknown option: $1"
            shift
            ;;
    esac
done

print_info "Test suite: $TEST_SUITE"
print_info "Environment: $ENVIRONMENT"
print_info "Device: $DEVICE_NAME"
print_info "iOS Version: $IOS_VERSION"

# Step 1: Check prerequisites
print_step "Step 1: Checking prerequisites..."

# Check if we're in the right directory
if [ ! -f "requirements.txt" ]; then
    print_error "Not in appium directory. Please run from smor_ting_mobile/appium/"
    exit 1
fi

# Check Xcode
if ! command -v xcodebuild &> /dev/null; then
    print_error "Xcode not found"
    print_info "Please install Xcode from the App Store"
    exit 1
fi

# Check xcrun
if ! command -v xcrun &> /dev/null; then
    print_error "Xcode command line tools not found"
    print_info "Install with: xcode-select --install"
    exit 1
fi

print_status "Prerequisites check passed"

# Step 2: Check if Flutter app is built
print_step "Step 2: Checking Flutter app..."

APP_PATH="../build/ios/iphonesimulator/Runner.app"
if [ ! -d "$APP_PATH" ]; then
    print_warning "iOS app not found at $APP_PATH"
    print_info "Building Flutter app for iOS..."
    cd .. && flutter build ios --simulator --debug
    cd appium
    
    if [ ! -d "$APP_PATH" ]; then
        print_error "Failed to build Flutter iOS app"
        exit 1
    fi
fi

print_status "Flutter iOS app found: $APP_PATH"

# Step 3: Setup iOS Simulator
print_step "Step 3: Setting up iOS Simulator..."

# List available simulators
print_info "Available iOS simulators:"
xcrun simctl list devices available | grep iPhone

# Find or create simulator
SIMULATOR_ID=""
SIMULATOR_LIST=$(xcrun simctl list devices available)

# Try to find existing simulator
EXISTING_SIM=$(echo "$SIMULATOR_LIST" | grep "$DEVICE_NAME" | grep "$IOS_VERSION" | head -1)
if [ ! -z "$EXISTING_SIM" ]; then
    SIMULATOR_ID=$(echo "$EXISTING_SIM" | grep -o '[A-F0-9-]\{36\}')
    print_info "Found existing simulator: $SIMULATOR_ID"
else
    print_info "Creating new simulator: $DEVICE_NAME with iOS $IOS_VERSION"
    
    # Get device type and runtime identifiers
    DEVICE_TYPE=$(xcrun simctl list devicetypes | grep "$DEVICE_NAME" | head -1 | sed 's/.*(\(.*\)).*/\1/')
    IOS_RUNTIME=$(xcrun simctl list runtimes | grep "iOS $IOS_VERSION" | head -1 | sed 's/.*(\(.*\)).*/\1/')
    
    if [ -z "$DEVICE_TYPE" ] || [ -z "$IOS_RUNTIME" ]; then
        print_error "Could not find device type or iOS runtime"
        print_info "Available device types:"
        xcrun simctl list devicetypes | grep iPhone
        print_info "Available runtimes:"
        xcrun simctl list runtimes | grep iOS
        exit 1
    fi
    
    SIMULATOR_ID=$(xcrun simctl create "QA-$DEVICE_NAME-$IOS_VERSION" "$DEVICE_TYPE" "$IOS_RUNTIME")
    print_status "Created simulator: $SIMULATOR_ID"
fi

# Boot simulator
print_info "Booting iOS simulator..."
xcrun simctl boot "$SIMULATOR_ID" 2>/dev/null || true

# Wait for simulator to boot
print_info "Waiting for simulator to boot..."
xcrun simctl bootstatus "$SIMULATOR_ID" -b

print_status "iOS simulator is ready"

# Step 4: Start Appium server
print_step "Step 4: Starting Appium server..."

# Kill any existing Appium processes
pkill -f appium || true
sleep 2

# Start Appium server in background
print_info "Starting Appium server on port 4723..."
if command -v appium &> /dev/null; then
    appium server --port 4723 --log reports/appium.log --log-level info &
else
    npx --yes appium server --port 4723 --log reports/appium.log --log-level info &
fi
APPIUM_PID=$!

# Wait for Appium to start
sleep 15

# Check if Appium is running
APPIUM_READY=false
for i in {1..30}; do
    if curl -s http://127.0.0.1:4723/status > /dev/null; then
        APPIUM_READY=true
        break
    fi
    sleep 2
    print_info "Waiting for Appium server... (${i}/30)"
done

if [ "$APPIUM_READY" = false ]; then
    print_error "Failed to start Appium server"
    cat reports/appium.log 2>/dev/null || echo "No log file found"
    exit 1
fi

print_status "Appium server running (PID: $APPIUM_PID)"

# Step 5: Install Python dependencies
print_step "Step 5: Installing Python dependencies..."

pip3 install -r requirements.txt > /dev/null 2>&1
print_status "Python dependencies installed"

# Step 6: Run tests
print_step "Step 6: Running iOS tests..."

# Set environment variables
export PLATFORM=ios
export ENVIRONMENT="$ENVIRONMENT"
export APP_PATH="$APP_PATH"
export IOS_DEVICE_NAME="$DEVICE_NAME"
export IOS_VERSION="$IOS_VERSION"
export SIMULATOR_UDID="$SIMULATOR_ID"

# Create reports directory
mkdir -p reports/screenshots
mkdir -p reports/logs

# Build pytest command based on test suite
PYTEST_ARGS="tests/"
case $TEST_SUITE in
    "auth")
        PYTEST_ARGS="tests/auth/"
        ;;
    "registration") 
        PYTEST_ARGS="tests/auth/test_registration.py"
        ;;
    "login")
        PYTEST_ARGS="tests/auth/test_login.py"
        ;;
    "all")
        PYTEST_ARGS="tests/"
        ;;
    *)
        print_warning "Unknown test suite: $TEST_SUITE, running all tests"
        PYTEST_ARGS="tests/"
        ;;
esac

# Run pytest with comprehensive reporting (use python -m to avoid PATH issues)
print_info "Executing test suite: $TEST_SUITE"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_NAME="ios-report-${TIMESTAMP}.html"
JUNIT_NAME="ios-junit-${TIMESTAMP}.xml"

python3 -m pytest $PYTEST_ARGS \
    -v \
    --platform=ios \
    --environment="$ENVIRONMENT" \
    --html="reports/$REPORT_NAME" \
    --self-contained-html \
    --junitxml="reports/$JUNIT_NAME" \
    --tb=short \
    --maxfail=5 \
    --capture=tee-sys \
    $PARALLEL \
    $MARKERS

TEST_EXIT_CODE=$?

# Generate test summary
print_step "Generating test summary..."

# Create latest symlinks
ln -sf "$REPORT_NAME" reports/ios-report.html
ln -sf "$JUNIT_NAME" reports/ios-junit.xml

# Count test results from JUnit XML
if [ -f "reports/$JUNIT_NAME" ]; then
    TOTAL_TESTS=$(grep -o 'tests="[0-9]*"' "reports/$JUNIT_NAME" | grep -o '[0-9]*' | head -1)
    FAILED_TESTS=$(grep -o 'failures="[0-9]*"' "reports/$JUNIT_NAME" | grep -o '[0-9]*' | head -1)
    ERROR_TESTS=$(grep -o 'errors="[0-9]*"' "reports/$JUNIT_NAME" | grep -o '[0-9]*' | head -1)
    SKIPPED_TESTS=$(grep -o 'skipped="[0-9]*"' "reports/$JUNIT_NAME" | grep -o '[0-9]*' | head -1)
    
    PASSED_TESTS=$((TOTAL_TESTS - FAILED_TESTS - ERROR_TESTS - SKIPPED_TESTS))
    
    echo "📊 Test Summary:"
    echo "  Total: ${TOTAL_TESTS:-0}"
    echo "  Passed: ${PASSED_TESTS:-0}"
    echo "  Failed: ${FAILED_TESTS:-0}"
    echo "  Errors: ${ERROR_TESTS:-0}"
    echo "  Skipped: ${SKIPPED_TESTS:-0}"
fi

# Step 7: Cleanup
print_step "Step 7: Cleaning up..."

# Stop Appium server
if [ ! -z "$APPIUM_PID" ]; then
    kill $APPIUM_PID 2>/dev/null || true
    sleep 2
    print_status "Appium server stopped"
fi

# Shutdown simulator
if [ ! -z "$SIMULATOR_ID" ]; then
    xcrun simctl shutdown "$SIMULATOR_ID" 2>/dev/null || true
    print_status "iOS simulator stopped"
fi

# Final results
echo ""
echo "🏁 Test Execution Complete!"
echo "============================"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    print_status "All tests passed! ✨"
else
    print_error "Some tests failed (exit code: $TEST_EXIT_CODE)"
fi

print_info "Test report: reports/$REPORT_NAME"
print_info "JUnit XML: reports/$JUNIT_NAME"
print_info "Appium logs: reports/appium.log"
print_info "Screenshots: reports/screenshots/"

# Archive results if in CI
if [ "$ENVIRONMENT" = "ci" ]; then
    print_info "Archiving test results..."
    tar -czf "reports/ios-test-results-${TIMESTAMP}.tar.gz" reports/
fi

# Open report in browser (macOS)
if command -v open &> /dev/null && [ -f "reports/$REPORT_NAME" ]; then
    print_info "Opening test report in browser..."
    open "reports/$REPORT_NAME"
fi

echo ""
exit $TEST_EXIT_CODE