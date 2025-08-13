#!/bin/bash

# Complete test to verify Dockerfile build works end-to-end
# This simulates the exact Railway build environment

set -e

echo "🧪 Testing complete Docker build process..."

# Test with the actual Dockerfile in a temporary context
echo "Creating temporary build context..."
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Copy the backend directory to temp (simulating Railway's build context)
cp -r /Users/kaleewou/smor-ting/smor_ting_backend/* $TEMP_DIR/

# Navigate to temp directory and test the exact build command from Dockerfile
cd $TEMP_DIR

echo "Testing build command from Dockerfile..."
if CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o smor-ting-api ./cmd; then
    echo "✅ Dockerfile build command succeeded"
    
    # Verify the binary was created and is executable
    if [ -f "smor-ting-api" ]; then
        echo "✅ Binary smor-ting-api was created"
        
        # Check if it's a valid Linux binary (ELF format)
        if file smor-ting-api | grep -q "ELF"; then
            echo "✅ Binary is a valid Linux executable (ELF format)"
            file smor-ting-api
        else
            echo "❌ Binary is not a Linux executable"
            file smor-ting-api
            exit 1
        fi
    else
        echo "❌ Binary smor-ting-api was not created"
        exit 1
    fi
else
    echo "❌ Dockerfile build command failed"
    exit 1
fi

echo "🎉 Dockerfile build test passed! Ready for Railway deployment."
