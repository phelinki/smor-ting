#!/bin/bash

# Production deployment script with confirmation

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo "⚠️  WARNING: This will deploy to PRODUCTION!"
echo "This action cannot be undone."
echo ""
read -p "Are you sure you want to deploy to production? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ Deployment cancelled"
    exit 1
fi

echo "🚀 Deploying to production..."

# Run the main deployment script
./scripts/play_store/deploy.sh -t production

echo "✅ Production deployment completed!"
