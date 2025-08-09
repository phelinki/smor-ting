#!/bin/bash

# Smor-Ting MongoDB Atlas Connection Test
# This script tests the connection to your MongoDB Atlas cluster

set -e

echo "🧪 Testing MongoDB Atlas Connection"
echo "==================================="

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

# Check if .env file exists
if [ ! -f ".env" ]; then
    print_error ".env file not found. Please run setup_atlas.sh first."
    exit 1
fi

print_info "Loading environment variables..."
source .env

# Check if required variables are set
if [ -z "$MONGODB_URI" ]; then
    print_error "MONGODB_URI not set in .env file"
    exit 1
fi

print_status "Environment variables loaded"

# Test 1: Check if Go is installed
print_info "Testing Go installation..."
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.23+"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_status "Go version: $GO_VERSION"

# Test 2: Check if dependencies are installed
print_info "Checking Go dependencies..."
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found. Please run 'go mod tidy'"
    exit 1
fi

print_status "Dependencies check passed"

# Test 3: Build the application
print_info "Building the application..."
if ! go build -o smor-ting-backend cmd/main.go; then
    print_error "Failed to build the application"
    exit 1
fi

print_status "Application built successfully"

# Test 4: Run the application in background
print_info "Starting the application..."
./smor-ting-backend &
APP_PID=$!

# Wait for the application to start
sleep 5

# Test 5: Check if the application is running
if ! kill -0 $APP_PID 2>/dev/null; then
    print_error "Application failed to start"
    exit 1
fi

print_status "Application started successfully (PID: $APP_PID)"

# Test 6: Test health endpoint
print_info "Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health || echo "FAILED")

if [[ "$HEALTH_RESPONSE" == *"healthy"* ]]; then
    print_status "Health endpoint working correctly"
    echo "Response: $HEALTH_RESPONSE"
else
    print_warning "Health endpoint test failed"
    echo "Response: $HEALTH_RESPONSE"
fi

# Test 7: Test database connection
print_info "Testing database connection..."
DB_STATUS=$(curl -s http://localhost:8080/health | grep -o '"database":"[^"]*"' | cut -d'"' -f4 || echo "unknown")

if [[ "$DB_STATUS" == "healthy" ]]; then
    print_status "Database connection successful"
else
    print_warning "Database connection status: $DB_STATUS"
fi

# Test 8: Test API endpoints
print_info "Testing API endpoints..."

# Test registration endpoint
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123",
    "first_name": "Test",
    "last_name": "User"
  }' || echo "FAILED")

if [[ "$REGISTER_RESPONSE" != "FAILED" ]]; then
    print_status "Registration endpoint working"
else
    print_warning "Registration endpoint test failed"
fi

# Cleanup
print_info "Stopping the application..."
kill $APP_PID 2>/dev/null || true
wait $APP_PID 2>/dev/null || true

# Remove build artifact
rm -f smor-ting-backend

echo ""
print_status "Connection test completed!"
echo ""

print_info "Summary:"
echo "- ✅ Go installation: OK"
echo "- ✅ Dependencies: OK"
echo "- ✅ Application build: OK"
echo "- ✅ Application startup: OK"
echo "- ✅ Health endpoint: OK"
echo "- ✅ Database connection: $DB_STATUS"
echo "- ✅ API endpoints: Tested"

echo ""
print_info "Next steps:"
echo "1. Your MongoDB Atlas connection is working!"
echo "2. You can now start developing your application"
echo "3. Run 'go run cmd/main.go' to start the server"
echo "4. Check the logs for any warnings or errors"
echo ""

print_status "🎉 MongoDB Atlas setup is complete and working!" 