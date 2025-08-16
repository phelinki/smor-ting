#!/bin/bash

# Google Play Store Deployment Script
# This script builds and deploys the app to Google Play Store testing tracks

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Configuration
PACKAGE_NAME="com.smorting.app.smor_ting_mobile"
TRACK="internal"  # internal, closed, or production
SERVICE_ACCOUNT_FILE="$SCRIPT_DIR/service-account.json"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -t, --track TRACK     Testing track (internal, closed, production) [default: internal]"
    echo "  -p, --package NAME    Package name [default: $PACKAGE_NAME]"
    echo "  -h, --help           Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Deploy to internal testing"
    echo "  $0 -t closed         # Deploy to closed testing"
    echo "  $0 -t production     # Deploy to production"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--track)
            TRACK="$2"
            shift 2
            ;;
        -p|--package)
            PACKAGE_NAME="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validate track
if [[ ! "$TRACK" =~ ^(internal|closed|production)$ ]]; then
    print_error "Invalid track: $TRACK. Must be internal, closed, or production"
    exit 1
fi

print_status "Starting deployment to $TRACK track..."

# Check if we're in the right directory
if [ ! -f "pubspec.yaml" ]; then
    print_error "This script must be run from the Flutter project root directory"
    exit 1
fi

# Check if service account file exists
if [ ! -f "$SERVICE_ACCOUNT_FILE" ]; then
    print_error "Service account file not found at: $SERVICE_ACCOUNT_FILE"
    echo "Please follow the setup instructions in scripts/play_store/SERVICE_ACCOUNT_SETUP.md"
    exit 1
fi

# Check if key.properties exists
if [ ! -f "android/key.properties" ]; then
    print_error "android/key.properties not found"
    exit 1
fi

# Set environment variables
export GOOGLE_APPLICATION_CREDENTIALS="$SERVICE_ACCOUNT_FILE"

# Build the app bundle
print_status "Building Android App Bundle..."
flutter clean
flutter pub get

# Build the app bundle
flutter build appbundle --release

# Check if build was successful
AAB_FILE="build/app/outputs/bundle/release/app-release.aab"
if [ ! -f "$AAB_FILE" ]; then
    print_error "App bundle not found at: $AAB_FILE"
    exit 1
fi

print_success "App bundle built successfully: $AAB_FILE"

# Get bundle size
BUNDLE_SIZE=$(du -h "$AAB_FILE" | cut -f1)
print_status "Bundle size: $BUNDLE_SIZE"

# Deploy to Google Play Store
print_status "Deploying to Google Play Store $TRACK track..."

# Install required Python packages if not available
if ! python3 -c "import googleapiclient" 2>/dev/null; then
    print_status "Installing required Python packages..."
    pip3 install google-api-python-client google-auth-httplib2 google-auth-oauthlib
fi

# Create deployment script
cat > scripts/play_store/deploy_to_play.py << 'PYTHON_EOF'
#!/usr/bin/env python3

import os
import sys
import json
from google.oauth2 import service_account
from googleapiclient.discovery import build
from googleapiclient.http import MediaFileUpload
from googleapiclient.errors import HttpError

def deploy_to_play_store(package_name, track, aab_file, service_account_file):
    """Deploy AAB to Google Play Store"""
    
    try:
        # Set up credentials
        credentials = service_account.Credentials.from_service_account_file(
            service_account_file,
            scopes=['https://www.googleapis.com/auth/androidpublisher']
        )
        
        # Build the service
        service = build('androidpublisher', 'v3', credentials=credentials)
        
        print(f"📱 Deploying {package_name} to {track} track...")
        
        # Create an edit
        edit_request = service.edits().insert(packageName=package_name, body={})
        edit_response = edit_request.execute()
        edit_id = edit_response['id']
        
        print(f"✅ Created edit: {edit_id}")
        
        # Upload the AAB
        print("📤 Uploading AAB file...")
        media = MediaFileUpload(aab_file, mimetype='application/octet-stream', resumable=True)
        
        upload_request = service.edits().bundles().upload(
            packageName=package_name,
            editId=edit_id,
            media_body=media
        )
        
        upload_response = upload_request.execute()
        version_code = upload_response['versionCode']
        
        print(f"✅ AAB uploaded successfully. Version code: {version_code}")
        
        # Update the track
        print(f"🔄 Updating {track} track...")
        
        track_body = {
            'releases': [{
                'name': f'Release {version_code}',
                'versionCodes': [str(version_code)],
                'status': 'completed'
            }]
        }
        
        track_request = service.edits().tracks().update(
            packageName=package_name,
            editId=edit_id,
            track=track,
            body=track_body
        )
        
        track_response = track_request.execute()
        
        print(f"✅ Track updated successfully")
        
        # Commit the edit
        print("🔒 Committing changes...")
        commit_request = service.edits().commit(
            packageName=package_name,
            editId=edit_id
        )
        
        commit_response = commit_request.execute()
        
        print("🎉 Deployment completed successfully!")
        print(f"📱 Package: {package_name}")
        print(f"🎯 Track: {track}")
        print(f"📦 Version: {version_code}")
        
        return True
        
    except HttpError as e:
        print(f"❌ HTTP Error: {e}")
        if e.resp.status == 403:
            print("💡 Make sure your service account has the correct permissions")
        return False
    except Exception as e:
        print(f"❌ Error: {e}")
        return False

if __name__ == "__main__":
    if len(sys.argv) != 5:
        print("Usage: python3 deploy_to_play.py <package_name> <track> <aab_file> <service_account_file>")
        sys.exit(1)
    
    package_name = sys.argv[1]
    track = sys.argv[2]
    aab_file = sys.argv[3]
    service_account_file = sys.argv[4]
    
    success = deploy_to_play_store(package_name, track, aab_file, service_account_file)
    sys.exit(0 if success else 1)
PYTHON_EOF

# Run the deployment
python3 scripts/play_store/deploy_to_play.py "$PACKAGE_NAME" "$TRACK" "$AAB_FILE" "$SERVICE_ACCOUNT_FILE"

if [ $? -eq 0 ]; then
    print_success "Deployment completed successfully!"
    
    # Create deployment summary
    cat > build/play_store/DEPLOYMENT_SUMMARY.md << EOF
# Play Store Deployment Summary

## Deployment Details
- **Package Name**: $PACKAGE_NAME
- **Track**: $TRACK
- **Bundle Size**: $BUNDLE_SIZE
- **Deployment Time**: $(date)

## Next Steps
1. Wait for Google Play Store to process the upload (5-15 minutes)
2. Check the Google Play Console for processing status
3. Add testers to the $TRACK track if needed
4. Testers will receive an email invitation

## Testing Links
- **Internal Testing**: https://play.google.com/console/u/0/developers/$(grep -o '[0-9]*' <<< "$PACKAGE_NAME" | head -1)/app/$(grep -o '[0-9]*' <<< "$PACKAGE_NAME" | tail -1)/testing/internal
- **Closed Testing**: https://play.google.com/console/u/0/developers/$(grep -o '[0-9]*' <<< "$PACKAGE_NAME" | head -1)/app/$(grep -o '[0-9]*' <<< "$PACKAGE_NAME" | tail -1)/testing/closed

## Bundle Information
- **File**: $AAB_FILE
- **Size**: $BUNDLE_SIZE
- **Version**: $(grep "versionCode" android/app/build.gradle.kts | grep -o '[0-9]*')
EOF
    
    print_success "Deployment summary saved to: build/play_store/DEPLOYMENT_SUMMARY.md"
else
    print_error "Deployment failed!"
    exit 1
fi
