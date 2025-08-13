#!/bin/bash

# Test script to verify Dockerfile build works correctly
# This tests that the cmd package (with multiple files) builds properly

set -e

echo "🧪 Testing Docker build process..."

# Test 1: Verify local build works with cmd package
echo "Test 1: Local build with cmd package"
cd /Users/kaleewou/smor-ting/smor_ting_backend
if go build -o test-binary ./cmd; then
    echo "✅ Local build with ./cmd succeeded"
    rm -f test-binary
else
    echo "❌ Local build with ./cmd failed"
    exit 1
fi

# Test 2: Verify local build fails with cmd/main.go (current Dockerfile approach)
echo "Test 2: Local build with cmd/main.go (should fail)"
if go build -o test-binary cmd/main.go 2>/dev/null; then
    echo "❌ Build with cmd/main.go unexpectedly succeeded"
    rm -f test-binary
    exit 1
else
    echo "✅ Build with cmd/main.go correctly failed (as expected)"
fi

# Test 3: Simulate the exact Docker build command locally
echo "Test 3: Simulate Docker build command"
if CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o smor-ting-api ./cmd; then
    echo "✅ Docker-style build with ./cmd succeeded"
    rm -f smor-ting-api
else
    echo "❌ Docker-style build with ./cmd failed"
    exit 1
fi

echo "🎉 All tests passed! Dockerfile should use './cmd' not 'cmd/main.go'"
