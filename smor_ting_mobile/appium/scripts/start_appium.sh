#!/bin/bash

echo "🚀 Starting Appium Server for Smor-Ting QA..."

# Check if Appium is installed
if ! command -v appium &> /dev/null; then
    echo "❌ Appium is not installed. Please install with: npm install -g appium@next"
    exit 1
fi

# Create reports directory if it doesn't exist
mkdir -p ../reports

# Kill any existing Appium processes
echo "🔄 Stopping any existing Appium processes..."
pkill -f "appium" || true

# Start Appium server
echo "▶️  Starting Appium server on port 4723..."
appium server \
    --port 4723 \
    --log-level info \
    --log ../reports/appium.log \
    --log-timestamp \
    --local-timezone &

# Wait for server to start
echo "⏳ Waiting for Appium server to initialize..."
sleep 8

# Check if server is running
if curl -s http://127.0.0.1:4723/status > /dev/null; then
    echo "✅ Appium server is running successfully on http://127.0.0.1:4723"
    echo "📝 Logs are being written to reports/appium.log"
else
    echo "❌ Failed to start Appium server"
    exit 1
fi

echo "🎯 Ready for test execution!"
