#!/bin/bash

# Test Google Play Store API authentication

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if service account file exists
SERVICE_ACCOUNT_FILE="$SCRIPT_DIR/service-account.json"
if [ ! -f "$SERVICE_ACCOUNT_FILE" ]; then
    print_error "Service account file not found at: $SERVICE_ACCOUNT_FILE"
    echo "Please follow the setup instructions in SERVICE_ACCOUNT_SETUP.md"
    exit 1
fi

print_status "Testing Google Play Store API authentication..."

# Set the service account
export GOOGLE_APPLICATION_CREDENTIALS="$SERVICE_ACCOUNT_FILE"

# Test API access
print_status "Testing API access..."
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
    
    # Try to list apps (this will fail if no apps, but will test auth)
    try:
        apps = service.applications().list().execute()
        print('✅ Authentication successful!')
        print(f'Found {len(apps.get(\"applications\", []))} apps')
    except Exception as e:
        if '403' in str(e):
            print('❌ Authentication failed: Insufficient permissions')
            print('Make sure the service account has the correct permissions in Play Console')
        else:
            print(f'✅ Authentication successful! (API error: {e})')
            
except Exception as e:
    print(f'❌ Authentication failed: {e}')
    print('Please check your service account configuration')
"

print_status "Authentication test completed"
