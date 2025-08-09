#!/bin/bash

# Generate Secure Encryption Keys for Smor-Ting Backend
# This script generates base64-encoded 32-byte keys for production use

set -e

echo "🔐 Generating Secure Encryption Keys for Smor-Ting Backend"
echo "=========================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Check if openssl is available
if ! command -v openssl &> /dev/null; then
    print_error "OpenSSL is not installed. Please install OpenSSL first."
    exit 1
fi

# Create keys directory
KEYS_DIR="./keys"
mkdir -p "$KEYS_DIR"

print_info "Generating 32-byte encryption keys..."

# Generate JWT Access Secret
JWT_ACCESS_SECRET=$(openssl rand -base64 32)
echo "JWT_ACCESS_SECRET=$JWT_ACCESS_SECRET" > "$KEYS_DIR/jwt_access_secret.env"

# Generate JWT Refresh Secret
JWT_REFRESH_SECRET=$(openssl rand -base64 32)
echo "JWT_REFRESH_SECRET=$JWT_REFRESH_SECRET" > "$KEYS_DIR/jwt_refresh_secret.env"

# Generate Encryption Key
ENCRYPTION_KEY=$(openssl rand -base64 32)
echo "ENCRYPTION_KEY=$ENCRYPTION_KEY" > "$KEYS_DIR/encryption_key.env"

# Generate Payment Encryption Key
PAYMENT_ENCRYPTION_KEY=$(openssl rand -base64 32)
echo "PAYMENT_ENCRYPTION_KEY=$PAYMENT_ENCRYPTION_KEY" > "$KEYS_DIR/payment_encryption_key.env"

# Set proper permissions (local only). Do not commit outputs.
chmod 600 "$KEYS_DIR"/*.env || true

print_status "Encryption keys generated successfully!"
echo ""
print_info "Generated keys (DO NOT COMMIT):"
echo "  📁 JWT Access Secret: $KEYS_DIR/jwt_access_secret.env"
echo "  📁 JWT Refresh Secret: $KEYS_DIR/jwt_refresh_secret.env"
echo "  📁 Encryption Key: $KEYS_DIR/encryption_key.env"
echo "  📁 Payment Encryption Key: $KEYS_DIR/payment_encryption_key.env"
echo ""

print_warning "IMPORTANT SECURITY NOTES:"
echo "  🔒 Keys are stored with 600 permissions (owner read/write only)"
echo "  🔒 Never commit these files to version control"
echo "  🔒 Use different keys for each environment (dev/staging/prod)"
echo "  🔒 Rotate keys regularly in production"
echo ""

print_info "To use these keys in your application:"
echo "  1. Copy the values to your .env file:"
echo "     source $KEYS_DIR/all_keys.env"
echo ""
echo "  2. Or manually add to your .env file:"
echo "     JWT_ACCESS_SECRET=$JWT_ACCESS_SECRET"
echo "     JWT_REFRESH_SECRET=$JWT_REFRESH_SECRET"
echo "     ENCRYPTION_KEY=$ENCRYPTION_KEY"
echo "     PAYMENT_ENCRYPTION_KEY=$PAYMENT_ENCRYPTION_KEY"
echo ""

print_warning "For production deployment:"
echo "  🚀 Store keys in environment variables or secure key management system"
echo "  🚀 Use different keys for each environment"
echo "  🚀 Implement key rotation strategy"
echo "  🚀 Monitor for key compromise"
echo ""

print_status "Key generation completed successfully!"
print_info "Your application is now ready with enterprise-grade security!"
