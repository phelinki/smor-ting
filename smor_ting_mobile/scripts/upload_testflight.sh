#!/bin/bash

set -e

echo "🚀 Uploading IPA to App Store Connect (TestFlight)"

IPA_PATH=${1:-"ios/build/ipa/Runner.ipa"}

if [ ! -f "$IPA_PATH" ]; then
  echo "❌ IPA not found at: $IPA_PATH"
  echo "   Build first, e.g.: ./scripts/build_testflight.sh"
  exit 1
fi

# Prefer App Store Connect API key auth
if [ -n "$ASC_KEY_ID" ] && [ -n "$ASC_ISSUER_ID" ]; then
  echo "ℹ️  Using App Store Connect API key (Key ID: $ASC_KEY_ID)"
  echo "   Ensure the private key file exists at ~/.appstoreconnect/private_keys/AuthKey_${ASC_KEY_ID}.p8"
  xcrun iTMSTransporter -m upload \
    -assetFile "$IPA_PATH" \
    -apiKey "$ASC_KEY_ID" \
    -apiIssuer "$ASC_ISSUER_ID" \
    -v informational
  echo "✅ Upload initiated via iTMSTransporter"
  exit 0
fi

if [ -n "$APPLE_ID" ] && [ -n "$APP_SPECIFIC_PASSWORD" ]; then
  echo "ℹ️  Using Apple ID with app-specific password"
  xcrun iTMSTransporter -m upload \
    -assetFile "$IPA_PATH" \
    -u "$APPLE_ID" \
    -p "$APP_SPECIFIC_PASSWORD" \
    -v informational
  echo "✅ Upload initiated via iTMSTransporter"
  exit 0
fi

echo "⚠️  No App Store Connect credentials provided."
echo "   Provide either ASC_KEY_ID/ASC_ISSUER_ID (+ key file) or APPLE_ID/APP_SPECIFIC_PASSWORD."
echo "   Alternatively, upload manually via Xcode Organizer."
exit 2


