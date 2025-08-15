#!/bin/bash

# Setup Local Persistent Database for Smor-Ting Development
# This script configures the backend to use a persistent local MongoDB database

set -e  # Exit on any error

echo "🗄️ Setting up Local Persistent Database for Smor-Ting"
echo "===================================================="

# Check if MongoDB is installed and running locally
echo "📋 Checking MongoDB installation..."

if ! command -v mongod &> /dev/null; then
    echo "❌ MongoDB is not installed. Installing MongoDB..."
    
    # Detect OS and install MongoDB
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            echo "🍺 Installing MongoDB via Homebrew..."
            brew tap mongodb/brew
            brew install mongodb-community
            
            echo "🚀 Starting MongoDB service..."
            brew services start mongodb/brew/mongodb-community
        else
            echo "❌ Homebrew not found. Please install MongoDB manually:"
            echo "   https://docs.mongodb.com/manual/installation/"
            exit 1
        fi
    else
        echo "❌ Please install MongoDB manually for your operating system:"
        echo "   https://docs.mongodb.com/manual/installation/"
        exit 1
    fi
else
    echo "✅ MongoDB is installed"
fi

# Check if MongoDB is running
echo "📋 Checking if MongoDB is running..."
if ! nc -z localhost 27017 &> /dev/null; then
    echo "🚀 Starting MongoDB..."
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew services start mongodb/brew/mongodb-community
        else
            mongod --dbpath /usr/local/var/mongodb --logpath /usr/local/var/log/mongodb/mongo.log --fork
        fi
    else
        # Linux
        sudo systemctl start mongod || mongod --dbpath /var/lib/mongodb --logpath /var/log/mongodb/mongod.log --fork
    fi
    
    # Wait for MongoDB to start
    echo "⏳ Waiting for MongoDB to start..."
    for i in {1..30}; do
        if nc -z localhost 27017 &> /dev/null; then
            echo "✅ MongoDB is running"
            break
        fi
        sleep 1
    done
    
    if ! nc -z localhost 27017 &> /dev/null; then
        echo "❌ MongoDB failed to start. Please check the installation and try again."
        exit 1
    fi
else
    echo "✅ MongoDB is already running"
fi

# Create .env file for local development
echo "📝 Creating local development configuration..."

cat > .env << EOF
# Local Development Configuration with Persistent Database
# This configuration ensures users persist across application restarts

# Server Configuration
PORT=8088
HOST=0.0.0.0
ENV=development

# Database Configuration - PERSISTENT LOCAL MongoDB
DB_DRIVER=mongodb
DB_HOST=localhost
DB_PORT=27017
DB_NAME=smor_ting_local
DB_USERNAME=
DB_PASSWORD=
DB_SSL_MODE=disable
DB_IN_MEMORY=false
MONGODB_ATLAS=false
MONGODB_URI=mongodb://localhost:27017/smor_ting_local

# Authentication Configuration
JWT_SECRET=local-development-jwt-secret-key-min-32-chars-for-security
JWT_ACCESS_SECRET=local-dev-access-secret-key-32chars
JWT_REFRESH_SECRET=local-dev-refresh-secret-key-32chars
JWT_EXPIRATION=24h
BCRYPT_COST=8

# CORS Configuration (Development - allows all origins)
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization,CF-Connecting-IP
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_CREDENTIALS=true

# Logging Configuration
LOG_LEVEL=debug
LOG_FORMAT=console
LOG_OUTPUT=stdout

# Rate Limiting
RATE_LIMIT_MAX_REQUESTS=100
RATE_LIMIT_WINDOW_SIZE=60

# Other Configuration
ENCRYPTION_KEY=local-dev-encryption-key-32-chars-long
PCI_COMPLIANCE_MODE=false
AUDIT_LOG_RETENTION_DAYS=30
EOF

echo "✅ Created .env file with persistent database configuration"

# Test database connection
echo "🧪 Testing database connection..."

# Create a simple test Go script
cat > test_db_connection.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/smor_ting_local"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("❌ Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	// Test ping
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("❌ Failed to ping MongoDB: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Successfully connected to MongoDB")
	fmt.Printf("📍 Database: %s\n", mongoURI)
	fmt.Println("💾 Data will persist across application restarts")
}
EOF

# Run the test (if go modules are initialized)
if [ -f "go.mod" ]; then
    echo "🏃 Running connection test..."
    source .env
    go run test_db_connection.go
    rm test_db_connection.go
else
    echo "⚠️  Skipping connection test (go.mod not found)"
    rm test_db_connection.go
fi

echo ""
echo "🎉 Local Persistent Database Setup Complete!"
echo "============================================"
echo ""
echo "✅ MongoDB is running locally on port 27017"
echo "✅ Database: smor_ting_local (persistent)"
echo "✅ Configuration: .env file created"
echo "✅ Users will persist across application restarts"
echo ""
echo "🚀 To start your application:"
echo "   cd smor_ting_backend"
echo "   go run cmd/main.go"
echo ""
echo "📱 To connect mobile app to local backend:"
echo "   Update smor_ting_mobile/lib/core/constants/api_config.dart"
echo "   Change development URL back to: http://127.0.0.1:8088"
echo ""
echo "🔍 To verify data persistence:"
echo "   1. Register a user"
echo "   2. Stop the backend (Ctrl+C)"
echo "   3. Restart the backend"
echo "   4. Login with the same user (should work!)"
echo ""
