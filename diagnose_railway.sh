#!/bin/bash

echo "🔍 Railway Deployment Diagnostic Script"
echo "======================================"
echo

# Check if Railway CLI is installed
echo "1. Checking Railway CLI..."
if command -v railway &> /dev/null; then
    railway --version
    echo "✅ Railway CLI is installed"
else
    echo "❌ Railway CLI not found. Install with: npm install -g @railway/cli"
fi
echo

# Check build configurations
echo "2. Checking build configurations..."
if [ -f "smor_ting_backend/Dockerfile" ]; then
    echo "✅ Found Dockerfile in smor_ting_backend/"
    echo "   Railway will use this for building"
else
    echo "❌ No Dockerfile found in smor_ting_backend/"
fi

if [ -f "railway.toml" ]; then
    echo "✅ Found railway.toml (root level)"
    echo "   Content:"
    cat railway.toml | head -10
else
    echo "⚠️  No railway.toml in root (this is OK if using default settings)"
fi

if [ -f "smor_ting_backend/railway.toml" ]; then
    echo "✅ Found railway.toml in backend directory"
    echo "   Content:"
    cat smor_ting_backend/railway.toml | head -10
else
    echo "⚠️  No railway.toml in backend directory"
fi

if [ -f "nixpacks.toml" ]; then
    echo "⚠️  Found nixpacks.toml - this might conflict with Dockerfile"
    echo "   Railway prefers Dockerfile over Nixpacks when both exist"
fi
echo

# Check Go build
echo "3. Testing Go build locally..."
cd smor_ting_backend
if go mod tidy && go build -o test-server ./cmd; then
    echo "✅ Go build successful locally"
    rm -f test-server
else
    echo "❌ Go build failed locally - this will also fail on Railway"
    cd ..
    exit 1
fi
cd ..
echo

# Check environment variables template
echo "4. Checking environment configuration..."
if [ -f "smor_ting_backend/.env.example" ]; then
    echo "✅ Found .env.example template"
    echo "   Required variables for Railway:"
    grep -E "^[A-Z_]+=" smor_ting_backend/.env.example | head -10
else
    echo "❌ No .env.example found - you'll need to configure variables manually"
fi
echo

# Check for common Railway issues
echo "5. Checking for common Railway deployment issues..."

# Check for start command consistency
echo "   Checking start command consistency..."
if grep -r "smor-ting-api" railway.toml smor_ting_backend/railway.toml smor_ting_backend/Dockerfile 2>/dev/null; then
    echo "   ✅ Start commands reference correct binary name"
else
    echo "   ⚠️  Start command mismatch detected"
fi

# Check for port configuration
echo "   Checking port configuration..."
if grep -r "PORT" smor_ting_backend/configs/config.go; then
    echo "   ✅ App reads PORT from environment (required for Railway)"
else
    echo "   ❌ App doesn't read PORT from environment"
fi

# Check health endpoint
echo "   Checking health endpoint..."
if grep -r "/health" smor_ting_backend/cmd/main.go; then
    echo "   ✅ Health endpoint found in code"
else
    echo "   ❌ No health endpoint found"
fi
echo

# Railway status check
echo "6. Checking Railway project status..."
if command -v railway &> /dev/null; then
    echo "   Current Railway context:"
    railway status 2>/dev/null || echo "   ⚠️  Not connected to Railway project or not logged in"
    echo
    
    echo "   Recent deployments:"
    railway logs --limit 20 2>/dev/null || echo "   ⚠️  Cannot fetch deployment logs"
else
    echo "   ❌ Railway CLI not available for status check"
fi
echo

echo "7. Railway deployment troubleshooting suggestions:"
echo "   Based on 'Deployment created' but not running issue:"
echo
echo "   🔧 Common fixes:"
echo "   1. Check Railway environment variables are set correctly"
echo "   2. Ensure MONGODB_URI and JWT_SECRET are configured"
echo "   3. Verify Railway is using the correct build configuration"
echo "   4. Check Railway logs for build/runtime errors"
echo "   5. Ensure your Railway plan has sufficient resources"
echo
echo "   💡 Immediate actions to try:"
echo "   1. railway redeploy --force"
echo "   2. railway logs --follow"
echo "   3. Check Railway dashboard for resource limits"
echo "   4. Verify environment variables in Railway settings"
echo
echo "   🚀 If still stuck:"
echo "   1. Try deploying with Nixpacks instead of Dockerfile"
echo "   2. Simplify the Dockerfile to debug build issues"
echo "   3. Check Railway status page for platform issues"

echo
echo "Diagnostic complete! 🏁"
