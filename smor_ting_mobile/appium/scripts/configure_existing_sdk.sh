#!/bin/bash

echo "🔧 Configuring Existing Android SDK for CI/CD"
echo "============================================="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

print_status() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

# Find Android Studio SDK
print_info "Looking for existing Android Studio SDK..."

POSSIBLE_SDK_PATHS=(
    "$HOME/Library/Android/sdk"
    "$HOME/Android/Sdk"
    "/opt/android-sdk"
    "/usr/local/android-sdk"
)

SDK_PATH=""
for path in "${POSSIBLE_SDK_PATHS[@]}"; do
    if [ -d "$path" ]; then
        print_info "Found potential SDK at: $path"
        if [ -f "$path/platform-tools/adb" ]; then
            SDK_PATH="$path"
            print_status "Valid Android SDK found at: $SDK_PATH"
            break
        fi
    fi
done

if [ -z "$SDK_PATH" ]; then
    print_error "No valid Android SDK found. Opening Android Studio to install SDK..."
    print_info "In Android Studio:"
    print_info "1. Go to Preferences → Appearance & Behavior → System Settings → Android SDK"
    print_info "2. Install SDK Platform 30, 31, 33"
    print_info "3. Install SDK Build Tools"
    print_info "4. Install Android Emulator"
    
    open -a "Android Studio" 2>/dev/null || print_warning "Could not open Android Studio automatically"
    exit 1
fi

# Update environment variables for CI/CD
print_info "Updating environment variables..."

# Clean old entries
sed -i '' '/export ANDROID_HOME/d' ~/.bash_profile 2>/dev/null
sed -i '' '/export ANDROID_SDK_ROOT/d' ~/.bash_profile 2>/dev/null

# Add new environment variables
cat >> ~/.bash_profile << EOF

# Android SDK for CI/CD (configured $(date))
export ANDROID_HOME="$SDK_PATH"
export ANDROID_SDK_ROOT="$SDK_PATH"
export PATH=\$PATH:\$ANDROID_HOME/platform-tools
export PATH=\$PATH:\$ANDROID_HOME/tools
export PATH=\$PATH:\$ANDROID_HOME/tools/bin
export PATH=\$PATH:\$ANDROID_HOME/emulator
export PATH=\$PATH:\$ANDROID_HOME/cmdline-tools/latest/bin

# Java 11 for Android (CI/CD compatible)
export JAVA_HOME=\$(/usr/libexec/java_home -v 11)
EOF

# Apply for current session
export ANDROID_HOME="$SDK_PATH"
export ANDROID_SDK_ROOT="$SDK_PATH"
export PATH=$PATH:$ANDROID_HOME/platform-tools
export PATH=$PATH:$ANDROID_HOME/tools
export PATH=$PATH:$ANDROID_HOME/tools/bin
export PATH=$PATH:$ANDROID_HOME/emulator
export JAVA_HOME=$(/usr/libexec/java_home -v 11)

# Verify tools
print_info "Verifying Android tools..."

if command -v adb &> /dev/null; then
    ADB_VERSION=$(adb --version | head -n1)
    print_status "ADB found: $ADB_VERSION"
else
    print_error "ADB not found in PATH"
fi

# Check for emulator
if [ -f "$ANDROID_HOME/emulator/emulator" ]; then
    print_status "Android Emulator found"
    
    # List available AVDs
    print_info "Available Android Virtual Devices:"
    if [ -f "$ANDROID_HOME/emulator/emulator" ]; then
        $ANDROID_HOME/emulator/emulator -list-avds 2>/dev/null || echo "No AVDs found"
    fi
else
    print_warning "Android Emulator not found"
fi

# Check for sdkmanager
if [ -f "$ANDROID_HOME/cmdline-tools/latest/bin/sdkmanager" ]; then
    print_status "SDK Manager found"
elif [ -f "$ANDROID_HOME/tools/bin/sdkmanager" ]; then
    print_status "SDK Manager found (legacy location)"
else
    print_warning "SDK Manager not found - may need command line tools"
fi

print_status "Android SDK configuration completed!"
print_info "SDK Path: $ANDROID_HOME"
print_info "Java Path: $JAVA_HOME"
print_info ""
print_info "Test with: flutter doctor"
print_info "Or run: ./scripts/run_android_tests.sh"
