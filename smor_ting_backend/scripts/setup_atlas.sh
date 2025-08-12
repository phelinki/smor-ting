#!/bin/bash

# Smor-Ting MongoDB Atlas Setup Script
# This script helps you set up MongoDB Atlas for your Smor-Ting application

set -e

echo "🚀 Smor-Ting MongoDB Atlas Setup"
echo "=================================="

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

# Check if .env file exists
if [ ! -f ".env" ]; then
    print_warning ".env file not found. Creating one..."
    cat > .env << EOF
# MongoDB Atlas Configuration
DB_DRIVER=mongodb
DB_HOST=your-cluster-host.mongodb.net
DB_PORT=27017
DB_NAME=smor_ting
DB_USERNAME=smorting_user
DB_PASSWORD=YOUR_MONGODB_PASSWORD
DB_SSL_MODE=require
DB_IN_MEMORY=false
MONGODB_ATLAS=true

# JWT Configuration
JWT_SECRET=YOUR_JWT_SECRET_MIN_32_CHARS
JWT_EXPIRATION=24h
BCRYPT_COST=12

# Server Configuration
PORT=8080
HOST=0.0.0.0
ENV=production

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
EOF
    print_status ".env file created"
else
    print_status ".env file already exists"
fi

echo ""
print_info "📋 MongoDB Atlas Setup Instructions:"
echo ""

echo "1. 🌐 Go to MongoDB Atlas:"
echo "   https://cloud.mongodb.com"
echo ""

echo "2. 📁 Create a New Project:"
echo "   - Click 'New Project'"
echo "   - Name: 'Smor-Ting'"
echo "   - Click 'Create Project'"
echo ""

echo "3. 🗄️  Build a Database:"
echo "   - Click 'Build a Database'"
echo "   - Choose 'FREE' tier (M0)"
echo "   - Select 'AWS' as cloud provider"
echo "   - Choose region:"
echo "     * For development (US): US East (N. Virginia)"
echo "     * For Liberia production: South Africa (Johannesburg)"
echo "     * Alternative: Europe (Ireland) - good balance"
echo "   - Click 'Create'"
echo ""

echo "4. 👤 Configure Database Access:"
echo "   - Create database user:"
echo "     Username: smorting_user"
echo "     Password: [generate strong password]"
echo "     Privileges: 'Read and write to any database'"
echo "   - Click 'Create User'"
echo ""

echo "5. 🌍 Configure Network Access:"
echo "   - Click 'Network Access'"
echo "   - Click 'Add IP Address'"
echo "   - For development: 'Allow Access from Anywhere' (0.0.0.0/0)"
echo "   - For production: Add your server's IP address"
echo "   - Click 'Confirm'"
echo ""

echo "6. 🔗 Get Connection String:"
echo "   - Click 'Database'"
echo "   - Click 'Connect' on your cluster"
echo "   - Choose 'Connect your application'"
echo "   - Copy the connection string"
echo ""

echo "7. ⚙️  Update Your .env File:"
echo "   - Replace the placeholder values in .env"
echo "   - Update MONGODB_URI with your connection string"
echo "   - Update DB_PASSWORD with your database user password"
echo ""

print_warning "⚠️  Important Security Notes:"
echo "   - Never commit your .env file to version control"
echo "   - Use strong passwords for database users"
echo "   - For production, restrict IP access to your server only"
echo "   - Change JWT_SECRET to a unique, strong value"
echo ""
print_info "🌍 Region Strategy for Liberia:"
echo "   - Start with US region for faster development"
echo "   - Migrate to South Africa when ready for production"
echo "   - Use ./scripts/migrate_region.sh for migration help"
echo ""

print_info "🔧 Next Steps:"
echo "1. Follow the instructions above to set up your Atlas cluster"
echo "2. Update the .env file with your actual values"
echo "3. Test the connection: go run cmd/main.go"
echo "4. Check the logs for successful connection"
echo ""

print_info "📚 Additional Resources:"
echo "- MongoDB Atlas Documentation: https://docs.atlas.mongodb.com"
echo "- Connection String Format: https://docs.mongodb.com/manual/reference/connection-string"
echo "- Security Best Practices: https://docs.atlas.mongodb.com/security"
echo ""

print_status "Setup script completed! Follow the instructions above to configure your MongoDB Atlas cluster." 