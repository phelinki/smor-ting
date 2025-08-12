#!/bin/bash

echo "🚀 Smor-Ting Appium QA Setup"
echo "============================"

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

echo ""
print_info "This script will set up Appium QA automation for Smor-Ting"
echo ""

# Step 1: Check system requirements
print_step "Step 1: Checking system requirements..."

# Check operating system
if [[ "$OSTYPE" == "darwin"* ]]; then
    print_status "Running on macOS - iOS and Android testing supported"
    PLATFORM_SUPPORT="both"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    print_status "Running on Linux - Android testing supported"
    PLATFORM_SUPPORT="android"
else
    print_warning "Running on $OSTYPE - limited support"
    PLATFORM_SUPPORT="android"
fi

# Check Node.js
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    print_status "Node.js found: $NODE_VERSION"
else
    print_error "Node.js not found"
    print_info "Install Node.js from: https://nodejs.org/"
    exit 1
fi

# Check Python
if command -v python3 &> /dev/null; then
    PYTHON_VERSION=$(python3 --version)
    print_status "Python found: $PYTHON_VERSION"
else
    print_error "Python 3 not found"
    print_info "Install Python 3 from: https://www.python.org/"
    exit 1
fi

# Check Java
if command -v java &> /dev/null; then
    JAVA_VERSION=$(java -version 2>&1 | head -n1)
    print_status "Java found: $JAVA_VERSION"
else
    print_warning "Java not found - required for Android testing"
    print_info "Install Java JDK 8 or 11"
fi

echo ""

# Step 2: Install Appium
print_step "Step 2: Installing Appium..."

if command -v appium &> /dev/null; then
    APPIUM_VERSION=$(appium --version)
    print_status "Appium already installed: $APPIUM_VERSION"
    
    read -p "Do you want to update Appium? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        npm install -g appium@next
        print_status "Appium updated"
    fi
else
    print_info "Installing Appium globally..."
    npm install -g appium@next
    
    if [ $? -eq 0 ]; then
        print_status "Appium installed successfully"
    else
        print_error "Failed to install Appium"
        exit 1
    fi
fi

# Step 3: Install Appium drivers
print_step "Step 3: Installing Appium drivers..."

print_info "Installing UiAutomator2 driver for Android..."
appium driver install uiautomator2

print_info "Installing Flutter driver for better Flutter app support..."
appium driver install --source=npm appium-flutter-driver

if [[ "$PLATFORM_SUPPORT" == "both" ]]; then
    print_info "Installing XCUITest driver for iOS..."
    appium driver install xcuitest
fi

print_status "Appium drivers installed"

# Step 4: Install Python dependencies
print_step "Step 4: Installing Python dependencies..."

if [ -f "requirements.txt" ]; then
    print_info "Installing Python packages from requirements.txt..."
    pip3 install -r requirements.txt
    
    if [ $? -eq 0 ]; then
        print_status "Python dependencies installed"
    else
        print_warning "Some Python packages may have failed to install"
    fi
else
    print_warning "requirements.txt not found"
    print_info "Creating requirements.txt with essential packages..."
    
    cat > requirements.txt << EOF
# Appium QA Automation Dependencies
Appium-Python-Client==3.1.0
selenium==4.15.0
pytest==7.4.3
pytest-html==4.1.1
allure-pytest==2.13.2
EOF
    
    pip3 install -r requirements.txt
    print_status "Basic Python dependencies installed"
fi

# Step 5: Setup Android environment (if applicable)
if [[ "$PLATFORM_SUPPORT" == "both" ]] || [[ "$PLATFORM_SUPPORT" == "android" ]]; then
    print_step "Step 5: Checking Android environment..."
    
    if [ -z "$ANDROID_HOME" ]; then
        print_warning "ANDROID_HOME not set"
        
        # Check common Android SDK locations
        COMMON_PATHS=(
            "$HOME/Library/Android/sdk"
            "$HOME/Android/Sdk"
            "/usr/local/android-sdk"
        )
        
        for path in "${COMMON_PATHS[@]}"; do
            if [ -d "$path" ]; then
                print_info "Found Android SDK at: $path"
                print_info "Add this to your shell profile:"
                echo "export ANDROID_HOME=$path"
                echo "export PATH=\$PATH:\$ANDROID_HOME/platform-tools:\$ANDROID_HOME/tools"
                break
            fi
        done
        
        if [ -z "$ANDROID_HOME" ]; then
            print_warning "Android SDK not found"
            print_info "Install Android Studio or standalone SDK"
        fi
    else
        print_status "ANDROID_HOME set to: $ANDROID_HOME"
        
        if command -v adb &> /dev/null; then
            print_status "ADB found and available"
        else
            print_warning "ADB not found in PATH"
        fi
    fi
fi

# Step 6: Setup iOS environment (macOS only)
if [[ "$OSTYPE" == "darwin"* ]]; then
    print_step "Step 6: Checking iOS environment..."
    
    if command -v xcodebuild &> /dev/null; then
        XCODE_VERSION=$(xcodebuild -version | head -n1)
        print_status "Xcode found: $XCODE_VERSION"
    else
        print_warning "Xcode not found"
        print_info "Install Xcode from the App Store for iOS testing"
    fi
    
    # Check for iOS Simulator
    if command -v xcrun &> /dev/null; then
        SIMULATORS=$(xcrun simctl list devices iPhone | grep iPhone | wc -l)
        if [ "$SIMULATORS" -gt 0 ]; then
            print_status "$SIMULATORS iPhone simulators available"
        else
            print_warning "No iPhone simulators found"
            print_info "Install iOS simulators through Xcode"
        fi
    fi
fi

# Step 7: Create directory structure
print_step "Step 7: Setting up project structure..."

mkdir -p config tests/android tests/ios tests/common/page_objects reports/screenshots

print_status "Directory structure created"

# Step 8: Create sample configuration files
print_step "Step 8: Creating configuration files..."

# Create test config if it doesn't exist
if [ ! -f "config/test_config.json" ]; then
    cat > config/test_config.json << EOF
{
  "appium_server": {
    "host": "127.0.0.1",
    "port": 4723,
    "timeout": 30
  },
  "test_data": {
    "valid_user": {
      "email": "test@smorting.com",
      "password": "TestPass123!",
      "first_name": "Test",
      "last_name": "User",
      "phone": "231777123456"
    },
    "existing_user": {
      "email": "libworker@smorting.com",
      "password": "Smorting8&",
      "first_name": "Job",
      "last_name": "Test",
      "phone": "231999999999"
    }
  },
  "timeouts": {
    "implicit_wait": 10,
    "explicit_wait": 20,
    "page_load": 30
  },
  "environment": {
    "api_base_url": "https://api.smor-ting.com/api/v1",
    "test_environment": "staging"
  }
}
EOF
    print_status "Test configuration created"
fi

# Step 9: Verify installation
print_step "Step 9: Verifying installation..."

print_info "Testing Appium installation..."
appium driver list --installed

if [ $? -eq 0 ]; then
    print_status "Appium verification successful"
else
    print_error "Appium verification failed"
fi

# Final summary
echo ""
echo "🎉 Setup Complete!"
echo "=================="
echo ""
print_info "Next steps:"
echo "1. Build your Flutter app:"
echo "   cd .. && flutter build apk --debug  # For Android"
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "   cd .. && flutter build ios --simulator --debug  # For iOS"
fi
echo ""
echo "2. Run tests:"
echo "   ./scripts/run_android_tests.sh  # For Android"
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "   ./scripts/run_ios_tests.sh      # For iOS"
fi
echo ""
print_info "Configuration files created in config/"
print_info "Test reports will be generated in reports/"
echo ""

if [[ "$PLATFORM_SUPPORT" == "both" ]]; then
    print_status "Ready for Android and iOS testing! 🤖🍎"
else
    print_status "Ready for Android testing! 🤖"
fi
