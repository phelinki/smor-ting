#!/bin/bash

# Quick deployment script for internal testing

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo "🚀 Quick deployment to internal testing..."

# Run the main deployment script
./scripts/play_store/deploy.sh -t internal

echo "✅ Quick deployment completed!"
