#!/bin/bash

# Create QA Test Users in Production Database
# This script creates the required test users for QA automation testing

set -e  # Exit on any error

echo "🎯 Creating QA Test Users in Production Database"
echo "=================================================="

# Check if MONGODB_URI is set
if [ -z "$MONGODB_URI" ]; then
    echo "❌ ERROR: MONGODB_URI environment variable is not set"
    echo "Please set the MongoDB production connection string:"
    echo "  export MONGODB_URI='mongodb+srv://...'"
    exit 1
fi

echo "✅ MongoDB URI detected: ${MONGODB_URI:0:20}..."

# Check if Railway production environment
if [ -z "$RAILWAY_ENVIRONMENT" ] || [ "$RAILWAY_ENVIRONMENT" != "production" ]; then
    echo "⚠️  WARNING: Not in Railway production environment"
    echo "Current environment: ${RAILWAY_ENVIRONMENT:-local}"
    
    # Ask for confirmation
    read -p "Are you sure you want to create test users? (y/N): " confirm
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        echo "❌ Operation cancelled"
        exit 1
    fi
fi

echo "🔧 Running QA test user creation tests..."

# Run the test to create users
export MONGODB_URI="$MONGODB_URI"
go test -v ./tests -run TestCreateQATestUsers

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ QA Test Users Created Successfully!"
    echo "====================================="
    echo ""
    echo "📋 Created Users:"
    echo "  🧑‍💼 Customer: qa_customer@smorting.com (TestPass123!)"
    echo "  🛠️  Provider: qa_provider@smorting.com (ProviderPass123!)"
    echo "  👤 Admin:    qa_admin@smorting.com (AdminPass123!)"
    echo ""
    echo "🔍 These users can now be used for QA automation testing"
    echo "📱 Mobile app tests should use these credentials"
    echo ""
else
    echo ""
    echo "❌ Failed to create QA test users"
    echo "Check the error messages above for details"
    exit 1
fi
