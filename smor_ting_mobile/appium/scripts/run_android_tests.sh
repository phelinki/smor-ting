#!/bin/bash

echo "🤖 Running Android Tests for Smor-Ting"
echo "======================================"

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

# Parse command line arguments
TEST_SUITE="all"
ENVIRONMENT="local"
PARALLEL=""
MARKERS=""
DEVICE_NAME=""

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
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --suite <all|auth|registration|login>  Test suite to run (default: all)"
            echo "  --environment <local|ci|staging>       Test environment (default: local)"
            echo "  --parallel                              Run tests in parallel"
            echo "  --markers <marker>                      Run tests with specific markers"
            echo "  --device <device_name>                  Specific device to use"
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

# Step 1: Check prerequisites
print_step "Step 1: Checking prerequisites..."

# Check if we're in the right directory
if [ ! -f "requirements.txt" ]; then
    print_error "Not in appium directory. Please run from smor_ting_mobile/appium/"
    exit 1
fi

# Check Android SDK
if [ -z "$ANDROID_HOME" ]; then
    print_error "ANDROID_HOME not set"
    print_info "Please set ANDROID_HOME environment variable"
    exit 1
fi

if ! command -v adb &> /dev/null; then
    print_error "ADB not found in PATH"
    print_info "Please add Android SDK platform-tools to PATH"
    exit 1
fi

print_status "Prerequisites check passed"

# Step 2: Check if Flutter app is built
print_step "Step 2: Checking Flutter app..."

APP_PATH="../build/app/outputs/flutter-apk/app-debug.apk"
if [ ! -f "$APP_PATH" ]; then
    print_warning "App not found at $APP_PATH"
    print_info "Building Flutter app..."
    cd .. && flutter build apk --debug
    cd appium
    
    if [ ! -f "$APP_PATH" ]; then
        print_error "Failed to build Flutter app"
        exit 1
    fi
fi

print_status "Flutter app found: $APP_PATH"

# Step 3: Start Android emulator if needed
print_step "Step 3: Checking Android emulator..."

# Check if any device is connected
DEVICES=$(adb devices | grep -v "List of devices" | grep -v "^$" | wc -l)

if [ "$DEVICES" -eq 0 ]; then
    print_warning "No Android devices connected"
    print_info "Starting Android emulator..."
    
    # List available AVDs
    AVDS=$(emulator -list-avds)
    
    if [ -z "$AVDS" ]; then
        print_error "No Android Virtual Devices found"
        print_info "Create an AVD in Android Studio or use avdmanager"
        exit 1
    fi
    
    # Use specified device or first available AVD
    if [ ! -z "$DEVICE_NAME" ]; then
        SELECTED_AVD="$DEVICE_NAME"
    else
        SELECTED_AVD=$(echo "$AVDS" | head -n1)
    fi
    
    print_info "Starting emulator: $SELECTED_AVD"
    
    # Start emulator with appropriate settings for CI/local
    if [ "$ENVIRONMENT" = "ci" ]; then
        emulator -avd "$SELECTED_AVD" -no-snapshot-save -no-audio -no-window -gpu swiftshader_indirect &
    else
        emulator -avd "$SELECTED_AVD" -no-snapshot-save -no-audio &
    fi
    
    EMULATOR_PID=$!
    # Ensure tests target this AVD
    export ANDROID_AVD_NAME="$SELECTED_AVD"
    
    print_info "Waiting for emulator to boot..."
    adb wait-for-device
    
    # Wait for emulator to be fully booted
    BOOT_TIMEOUT=300  # 5 minutes
    BOOT_COUNTER=0
    while [ "$(adb shell getprop sys.boot_completed 2>/dev/null)" != "1" ]; do
        sleep 5
        BOOT_COUNTER=$((BOOT_COUNTER + 5))
        if [ $BOOT_COUNTER -ge $BOOT_TIMEOUT ]; then
            print_error "Emulator boot timeout"
            exit 1
        fi
        print_info "Waiting for emulator to complete boot... (${BOOT_COUNTER}s)"
    done
    
    # Additional wait for system to stabilize
    sleep 10
    print_status "Emulator is ready"
else
    print_status "$DEVICES Android device(s) connected"
fi

# Step 4: Start Appium server
print_step "Step 4: Starting Appium server..."

# Kill any existing Appium processes
pkill -f appium || true
sleep 2

# Ensure UiAutomator2 driver installed
if ! appium driver list --installed | grep -q uiautomator2; then
  print_info "Installing Appium UiAutomator2 driver..."
  appium driver install uiautomator2 >/dev/null 2>&1 || true
fi

# Start Appium server in background
print_info "Starting Appium server on port 4723..."
appium server --port 4723 --log reports/appium.log --log-level info &
APPIUM_PID=$!

# Wait for Appium to start
sleep 10

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
print_step "Step 6: Running Android tests..."

# Set environment variables
export PLATFORM=android
export ENVIRONMENT="$ENVIRONMENT"
export APP_PATH="$APP_PATH"
if [ ! -z "$DEVICE_NAME" ]; then
    export ANDROID_DEVICE_NAME="$DEVICE_NAME"
    export ANDROID_AVD_NAME="$DEVICE_NAME"
fi

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

# Run pytest with comprehensive reporting
print_info "Executing test suite: $TEST_SUITE"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_NAME="android-report-${TIMESTAMP}.html"
JUNIT_NAME="android-junit-${TIMESTAMP}.xml"

pytest $PYTEST_ARGS \
    -v \
    --platform=android \
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
ln -sf "$REPORT_NAME" reports/android-report.html
ln -sf "$JUNIT_NAME" reports/android-junit.xml

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

# Stop emulator if we started it
if [ ! -z "$EMULATOR_PID" ]; then
    kill $EMULATOR_PID 2>/dev/null || true
    sleep 5
    print_status "Emulator stopped"
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
    tar -czf "reports/android-test-results-${TIMESTAMP}.tar.gz" reports/
fi

echo ""
exit $TEST_EXIT_CODE