#!/bin/bash

# Install Android Command-line Tools
# This script helps install the missing Android SDK command-line tools

set -e

echo "🔧 Installing Android Command-line Tools..."

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# Check if ANDROID_HOME is set
if [ -z "$ANDROID_HOME" ]; then
    print_error "ANDROID_HOME is not set"
    print_status "Setting ANDROID_HOME to default location..."
    export ANDROID_HOME="$HOME/Library/Android/sdk"
fi

print_status "Android SDK location: $ANDROID_HOME"

# Create cmdline-tools directory
CMD_TOOLS_DIR="$ANDROID_HOME/cmdline-tools"
LATEST_DIR="$CMD_TOOLS_DIR/latest"

print_status "Creating cmdline-tools directory..."
mkdir -p "$LATEST_DIR"

# Download command-line tools
print_status "Downloading Android Command-line Tools..."

# Get the latest version
LATEST_VERSION=$(curl -s https://developer.android.com/studio#command-line-tools-only | grep -o 'commandlinetools-mac-[0-9]*' | head -1)

if [ -z "$LATEST_VERSION" ]; then
    print_warning "Could not determine latest version, using a known version..."
    LATEST_VERSION="commandlinetools-mac-11076708_latest"
fi

DOWNLOAD_URL="https://dl.google.com/android/repository/$LATEST_VERSION.zip"
DOWNLOAD_FILE="/tmp/$LATEST_VERSION.zip"

print_status "Downloading from: $DOWNLOAD_URL"

# Download the file
curl -L -o "$DOWNLOAD_FILE" "$DOWNLOAD_URL"

if [ ! -f "$DOWNLOAD_FILE" ]; then
    print_error "Failed to download command-line tools"
    print_status "Please install manually via Android Studio:"
    print_status "1. Open Android Studio"
    print_status "2. Go to Tools → SDK Manager"
    print_status "3. Click on 'SDK Tools' tab"
    print_status "4. Check 'Android SDK Command-line Tools (latest)'"
    print_status "5. Click 'Apply' and install"
    exit 1
fi

print_success "Download completed"

# Extract the file
print_status "Extracting command-line tools..."
cd "$CMD_TOOLS_DIR"
unzip -q "$DOWNLOAD_FILE"

# Move contents to latest directory
if [ -d "cmdline-tools" ]; then
    mv cmdline-tools/* "$LATEST_DIR/"
    rmdir cmdline-tools
fi

# Clean up
rm "$DOWNLOAD_FILE"

print_success "Command-line tools installed successfully!"

# Add to PATH
print_status "Adding to PATH..."
echo 'export PATH="$ANDROID_HOME/cmdline-tools/latest/bin:$PATH"' >> ~/.bash_profile
export PATH="$ANDROID_HOME/cmdline-tools/latest/bin:$PATH"

# Accept licenses
print_status "Accepting Android licenses..."
yes | sdkmanager --licenses

print_success "Android Command-line Tools setup completed!"
print_status "Please restart your terminal or run: source ~/.bash_profile"
