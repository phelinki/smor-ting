#!/bin/bash

# Railway Deployment Health Check Script
# This script should be run periodically to ensure Railway deployment is healthy

set -e

echo "🔍 Railway Deployment Health Check"
echo "=================================="
echo

# Build and run the monitoring tool
if [ ! -f "railway_monitor" ]; then
    echo "🏗️ Building monitoring tool..."
    go build -o railway_monitor ./scripts/railway_deployment_monitor.go
    chmod +x railway_monitor
fi

echo "1. 🩺 Health Check"
./railway_monitor health

echo
echo "2. 📊 Railway Status"
./railway_monitor status

echo
echo "3. 🔍 Deployment Verification"
./railway_monitor verify

echo
echo "4. 🧪 Running Deployment Tests"
cd smor_ting_backend
go test ./tests/deployment_health_test.go -v
go test ./tests/railway_fix_test.go -v
cd ..

echo
echo "✅ All checks passed! Railway deployment is healthy."
echo "🔗 Production URL: https://smor-ting-production.up.railway.app"
echo "🔗 Custom Domain: https://smor-ting.com"
