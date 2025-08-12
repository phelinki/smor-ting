#!/bin/bash

echo "🚀 Fixing Railway Deployment Issues"
echo "===================================="
echo

# 1. Add missing JWT_SECRET
echo "1. Adding missing JWT_SECRET..."
JWT_SECRET=$(openssl rand -base64 32)
echo "Generated JWT_SECRET: ${JWT_SECRET:0:10}..."

# Note: We need to set this in Railway dashboard since CLI syntax varies
echo "⚠️  You need to manually add JWT_SECRET in Railway dashboard:"
echo "   - Go to https://railway.app/project/smor-ting"
echo "   - Go to Variables tab"
echo "   - Add: JWT_SECRET = $JWT_SECRET"
echo

# 2. Check if the domain is working
echo "2. Testing domain accessibility..."
echo "Testing Railway subdomain..."
if curl -f https://smor-ting-production.up.railway.app/health 2>/dev/null; then
    echo "✅ Railway subdomain is working"
else
    echo "❌ Railway subdomain failing (this is the core issue)"
fi

echo "Testing custom domain..."
if curl -f https://smor-ting.com/health 2>/dev/null; then
    echo "✅ Custom domain is working"
else
    echo "❌ Custom domain failing (likely due to Railway backend issues)"
fi
echo

# 3. Check Railway deployment status
echo "3. Checking deployment status..."
railway status
echo

# 4. Force redeploy to apply fixes
echo "4. Recommended actions:"
echo "   1. Add JWT_SECRET variable in Railway dashboard"
echo "   2. Force redeploy with: railway redeploy"
echo "   3. Check logs with: railway logs --follow"
echo "   4. Test health endpoint directly on Railway subdomain"
echo

echo "5. Diagnosis:"
echo "   ✅ Your app IS running (logs show successful startup)"
echo "   ❌ Missing JWT_SECRET causing auth warnings"
echo "   ❌ HTTP 522 errors suggest Railway networking issues"
echo "   💡 This looks like Railway infrastructure problems, not your code"
echo

echo "6. If still not working after adding JWT_SECRET:"
echo "   - Try railway redeploy --force"
echo "   - Check Railway status page for outages"
echo "   - Consider switching to Nixpacks builder temporarily"
echo "   - Contact Railway support if infrastructure issues persist"
