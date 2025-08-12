#!/bin/bash

echo "☕ Installing Java 11 for Android Development"
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

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    print_warning "Homebrew not found. Installing Homebrew first..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    
    if [ $? -ne 0 ]; then
        print_error "Failed to install Homebrew"
        exit 1
    fi
    print_status "Homebrew installed successfully"
fi

# Install Java 11
print_info "Installing Java 11 via Homebrew..."
brew install openjdk@11

if [ $? -ne 0 ]; then
    print_error "Failed to install Java 11"
    exit 1
fi

# Create symlink for system Java
print_info "Creating system symlink for Java 11..."
sudo ln -sfn /opt/homebrew/opt/openjdk@11/libexec/openjdk.jdk /Library/Java/JavaVirtualMachines/openjdk-11.jdk 2>/dev/null || \
sudo ln -sfn /usr/local/opt/openjdk@11/libexec/openjdk.jdk /Library/Java/JavaVirtualMachines/openjdk-11.jdk

# Update JAVA_HOME in bash profile
print_info "Updating JAVA_HOME in ~/.bash_profile..."

# Remove old JAVA_HOME entries
sed -i '' '/export JAVA_HOME/d' ~/.bash_profile

# Add new JAVA_HOME for Java 11
echo "" >> ~/.bash_profile
echo "# Java 11 for Android Development" >> ~/.bash_profile
echo "export JAVA_HOME=\$(/usr/libexec/java_home -v 11)" >> ~/.bash_profile

# Source the profile
source ~/.bash_profile

# Set for current session
export JAVA_HOME=$(/usr/libexec/java_home -v 11)

print_status "Java 11 installation completed!"

# Verify installation
print_info "Verifying Java installation..."
java -version
echo ""
echo "JAVA_HOME: $JAVA_HOME"

print_status "Java 11 is ready for Android development! ☕"
print_info "Now you can run: ./scripts/install_android_sdk.sh"
