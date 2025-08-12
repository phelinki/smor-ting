#!/bin/bash

echo "🤖 Installing Android SDK Command Line Tools"
echo "==========================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

# Check if Android SDK already exists
if [ -d "$HOME/Library/Android/sdk" ]; then
    print_warning "Android SDK directory already exists"
    read -p "Do you want to proceed anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Installation cancelled"
        exit 0
    fi
fi

# Create Android SDK directory
print_info "Creating Android SDK directory..."
mkdir -p $HOME/Library/Android/sdk
cd $HOME/Library/Android/sdk

# Download Android command line tools
print_info "Downloading Android command line tools..."
CMDLINE_TOOLS_URL="https://dl.google.com/android/repository/commandlinetools-mac-9477386_latest.zip"

if command -v curl &> /dev/null; then
    curl -o cmdline-tools.zip $CMDLINE_TOOLS_URL
elif command -v wget &> /dev/null; then
    wget -O cmdline-tools.zip $CMDLINE_TOOLS_URL
else
    print_error "Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Extract command line tools
print_info "Extracting command line tools..."
unzip -q cmdline-tools.zip
rm cmdline-tools.zip

# Move to correct directory structure
mkdir -p cmdline-tools/latest
mv cmdline-tools/* cmdline-tools/latest/ 2>/dev/null || true

# Set up environment variables for this session
export ANDROID_HOME=$HOME/Library/Android/sdk
export PATH=$PATH:$ANDROID_HOME/cmdline-tools/latest/bin
export PATH=$PATH:$ANDROID_HOME/platform-tools
export PATH=$PATH:$ANDROID_HOME/tools
export PATH=$PATH:$ANDROID_HOME/emulator

print_info "Installing SDK components..."

# Accept all licenses
yes | sdkmanager --licenses > /dev/null 2>&1

# Install essential SDK components
print_info "Installing Android platforms and tools..."
sdkmanager "platform-tools" "platforms;android-30" "platforms;android-31" "platforms;android-33"
sdkmanager "build-tools;30.0.3" "build-tools;31.0.0" "build-tools;33.0.0"
sdkmanager "emulator" "tools"

# Install system images for emulator
print_info "Installing system images for emulator..."
sdkmanager "system-images;android-30;google_apis;x86_64"

# Create AVD
print_info "Creating Android Virtual Device..."
echo "no" | avdmanager create avd -n "Pixel_4_API_30" -k "system-images;android-30;google_apis;x86_64" --device "pixel_4"

print_status "Android SDK installation completed!"

# Verify installation
print_info "Verifying installation..."
if [ -f "$ANDROID_HOME/platform-tools/adb" ]; then
    print_status "ADB found: $($ANDROID_HOME/platform-tools/adb --version | head -n1)"
else
    print_error "ADB not found after installation"
fi

if [ -f "$ANDROID_HOME/emulator/emulator" ]; then
    print_status "Emulator found"
else
    print_error "Emulator not found after installation"
fi

# List available AVDs
print_info "Available Android Virtual Devices:"
$ANDROID_HOME/emulator/emulator -list-avds

echo ""
print_status "Setup complete! 🎉"
print_info "Environment variables are already added to ~/.bash_profile"
print_info "Restart your terminal or run: source ~/.bash_profile"
print_info "Then test with: ./scripts/run_android_tests.sh"
