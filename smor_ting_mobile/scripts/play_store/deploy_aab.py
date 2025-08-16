#!/usr/bin/env python3

import os
import sys
import json
from google.oauth2 import service_account
from googleapiclient.discovery import build
from googleapiclient.http import MediaFileUpload
from googleapiclient.errors import HttpError

def deploy_to_play_store():
    """Deploy existing AAB to Google Play Store internal testing"""
    
    # Configuration
    package_name = "com.smorting.app.smor_ting_mobile"
    track = "internal"
    aab_file = "build/app/outputs/bundle/release/app-release.aab"
    service_account_file = "scripts/play_store/service-account.json"
    
    print(f"📱 Deploying {package_name} to {track} track...")
    
    # Check if files exist
    if not os.path.exists(aab_file):
        print(f"❌ AAB file not found: {aab_file}")
        return False
        
    if not os.path.exists(service_account_file):
        print(f"❌ Service account file not found: {service_account_file}")
        return False
    
    try:
        # Set up credentials
        credentials = service_account.Credentials.from_service_account_file(
            service_account_file,
            scopes=['https://www.googleapis.com/auth/androidpublisher']
        )
        
        # Build the service
        service = build('androidpublisher', 'v3', credentials=credentials)
        
        print("✅ Authentication successful")
        
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
        print(f"📱 App deployed to {track} track")
        print(f"📦 Version code: {version_code}")
        print(f"📁 AAB file: {aab_file}")
        
        return True
        
    except HttpError as e:
        if e.resp.status == 403:
            print("❌ Permission denied. Check service account permissions.")
            print("Make sure the service account has 'Release apps to testing tracks' permission.")
        elif e.resp.status == 404:
            print("❌ App not found. Make sure the app exists in Play Console.")
            print(f"Package name: {package_name}")
        else:
            print(f"❌ HTTP Error: {e}")
        return False
        
    except Exception as e:
        print(f"❌ Error: {e}")
        return False

if __name__ == "__main__":
    success = deploy_to_play_store()
    sys.exit(0 if success else 1)
