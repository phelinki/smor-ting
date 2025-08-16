# Google Play Store Service Account Setup

## Step 1: Create a Service Account
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing project
3. Enable the Google Play Developer API
4. Go to "IAM & Admin" → "Service Accounts"
5. Click "Create Service Account"
6. Name it "play-store-deployer"
7. Grant "Editor" role
8. Create and download the JSON key file

## Step 2: Set up Google Play Console
1. Go to [Google Play Console](https://play.google.com/console)
2. Go to "Setup" → "API access"
3. Link your Google Cloud project
4. Create a new service account or use existing
5. Grant "Release apps to testing tracks" permission
6. Download the service account JSON key

## Step 3: Configure the deployment
1. Place the service account JSON file in: scripts/play_store/service-account.json
2. Update the package name in scripts/play_store/deploy.sh if needed
3. Run the deployment script

## Step 4: Test the setup
Run: ./scripts/play_store/test_auth.sh
