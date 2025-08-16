#!/bin/bash

# Monitor deployment status

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

PACKAGE_NAME="com.smorting.app.smorTingMobile"
SERVICE_ACCOUNT_FILE="$SCRIPT_DIR/service-account.json"

# Colors
GREEN='\033[0;32m'
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

cd "$PROJECT_ROOT"

if [ ! -f "$SERVICE_ACCOUNT_FILE" ]; then
    echo "❌ Service account file not found"
    exit 1
fi

export GOOGLE_APPLICATION_CREDENTIALS="$SERVICE_ACCOUNT_FILE"

print_status "Monitoring deployment status..."

# Install required packages if needed
if ! python3 -c "import googleapiclient" 2>/dev/null; then
    pip3 install google-api-python-client google-auth-httplib2 google-auth-oauthlib
fi

python3 -c "
import json
from google.oauth2 import service_account
from googleapiclient.discovery import build

try:
    credentials = service_account.Credentials.from_service_account_file(
        '$SERVICE_ACCOUNT_FILE',
        scopes=['https://www.googleapis.com/auth/androidpublisher']
    )
    
    service = build('androidpublisher', 'v3', credentials=credentials)
    
    # Get app details
    app = service.applications().get(packageName='$PACKAGE_NAME').execute()
    print(f'📱 App: {app[\"title\"]}')
    print(f'📦 Package: {app[\"packageName\"]}')
    
    # Get tracks
    tracks = service.edits().tracks().list(packageName='$PACKAGE_NAME', editId='0').execute()
    
    for track in tracks.get('tracks', []):
        track_name = track['track']
        releases = track.get('releases', [])
        
        if releases:
            latest_release = releases[0]
            version_codes = latest_release.get('versionCodes', [])
            status = latest_release.get('status', 'unknown')
            
            print(f'🎯 Track: {track_name}')
            print(f'   Version: {version_codes[0] if version_codes else \"N/A\"}')
            print(f'   Status: {status}')
            print(f'   Name: {latest_release.get(\"name\", \"N/A\")}')
        else:
            print(f'🎯 Track: {track_name} (No releases)')
            
except Exception as e:
    print(f'❌ Error: {e}')
"

print_status "Monitoring completed"
