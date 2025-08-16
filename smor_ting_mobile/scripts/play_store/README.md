# Google Play Store CLI Deployment System

This directory contains a complete CLI system for automating Google Play Store deployments, similar to TestFlight for iOS.

## 🚀 Quick Start

### 1. Configure Service Account
Follow the instructions in `SERVICE_ACCOUNT_SETUP.md` to create and configure your Google Play Store service account.

### 2. Test Authentication
```bash
./scripts/play_store/test_auth.sh
```

### 3. Deploy to Internal Testing
```bash
./scripts/play_store/quick_deploy.sh
```

## 📋 Available Scripts

### Core Deployment Scripts

| Script | Description | Usage |
|--------|-------------|-------|
| `deploy.sh` | Main deployment script with full options | `./deploy.sh -t internal` |
| `quick_deploy.sh` | Quick deployment to internal testing | `./quick_deploy.sh` |
| `deploy_production.sh` | Production deployment with confirmation | `./deploy_production.sh` |

### Utility Scripts

| Script | Description | Usage |
|--------|-------------|-------|
| `test_auth.sh` | Test API authentication | `./test_auth.sh` |
| `monitor_deployment.sh` | Monitor deployment status | `./monitor_deployment.sh` |

## 🎯 Deployment Tracks

### Internal Testing
- **Purpose**: Quick testing with your team
- **Deployment Time**: ~5 minutes
- **Command**: `./quick_deploy.sh` or `./deploy.sh -t internal`

### Closed Testing
- **Purpose**: Testing with external testers
- **Deployment Time**: ~10 minutes
- **Command**: `./deploy.sh -t closed`

### Production
- **Purpose**: Public release
- **Deployment Time**: ~30 minutes (includes review)
- **Command**: `./deploy_production.sh`

## 🔧 Configuration

### Package Name
Default: `com.smorting.app.smorTingMobile`

To change the package name, edit the `PACKAGE_NAME` variable in `deploy.sh`:
```bash
PACKAGE_NAME="your.package.name"
```

### Service Account
Place your service account JSON file at:
```
scripts/play_store/service-account.json
```

## 📱 Deployment Process

### What Happens During Deployment

1. **Environment Check**
   - Validates Java 17 installation
   - Checks Android SDK configuration
   - Verifies keystore setup

2. **Build Process**
   - Cleans previous builds
   - Installs dependencies
   - Builds Android App Bundle (AAB)

3. **Upload Process**
   - Authenticates with Google Play Store API
   - Creates an edit session
   - Uploads the AAB file
   - Updates the specified track
   - Commits the changes

4. **Post-Deployment**
   - Generates deployment summary
   - Provides monitoring commands

### Deployment Output

```
🚀 Starting deployment to internal track...
ℹ️  Building Android App Bundle...
✅ App bundle built successfully: build/app/outputs/bundle/release/app-release.aab
ℹ️  Bundle size: 45.2MB
ℹ️  Deploying to Google Play Store internal track...
📱 Deploying com.smorting.app.smorTingMobile to internal track...
✅ Created edit: 123456789
📤 Uploading AAB file...
✅ AAB uploaded successfully. Version code: 42
🔄 Updating internal track...
✅ Track updated successfully
🔒 Committing changes...
🎉 Deployment completed successfully!
📱 Package: com.smorting.app.smorTingMobile
🎯 Track: internal
📦 Version: 42
✅ Deployment completed successfully!
```

## 🔍 Monitoring Deployments

### Check Deployment Status
```bash
./scripts/play_store/monitor_deployment.sh
```

### Expected Output
```
ℹ️  Monitoring deployment status...
📱 App: Smor Ting
📦 Package: com.smorting.app.smorTingMobile
🎯 Track: internal
   Version: 42
   Status: completed
   Name: Release 42
🎯 Track: closed
   Version: 41
   Status: completed
   Name: Release 41
```

## 🛠️ Troubleshooting

### Common Issues

#### 1. Service Account Authentication Failed
**Error**: `Authentication failed: Insufficient permissions`

**Solution**:
1. Check that the service account has the correct permissions in Google Play Console
2. Ensure the service account JSON file is properly formatted
3. Verify the Google Play Developer API is enabled

#### 2. Build Failed
**Error**: `App bundle not found`

**Solution**:
1. Check that `android/key.properties` exists
2. Verify the keystore file path is correct
3. Ensure Java 17 is installed and configured

#### 3. Upload Failed
**Error**: `HTTP Error: 403`

**Solution**:
1. Check service account permissions
2. Verify the package name matches your Play Console app
3. Ensure the app is properly configured in Play Console

### Debug Mode

To enable debug output, set the environment variable:
```bash
export DEBUG=1
./scripts/play_store/deploy.sh
```

## 📊 Deployment Summary

After each successful deployment, a summary is generated at:
```
build/play_store/DEPLOYMENT_SUMMARY.md
```

This includes:
- Deployment details (package, track, size, time)
- Next steps for testers
- Testing links
- Bundle information

## 🔐 Security Considerations

### Service Account Security
- Keep your service account JSON file secure
- Don't commit it to version control
- Use environment variables in production

### Keystore Security
- Store your keystore file securely
- Use strong passwords
- Back up your keystore file

## 📚 Integration with CI/CD

### GitHub Actions Example
```yaml
name: Deploy to Play Store
on:
  push:
    tags:
      - 'v*'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: subosito/flutter-action@v2
        with:
          flutter-version: '3.16.0'
      - run: ./scripts/play_store/deploy.sh -t internal
        env:
          GOOGLE_APPLICATION_CREDENTIALS: ${{ secrets.GOOGLE_APPLICATION_CREDENTIALS }}
```

## 🆚 Comparison with TestFlight

| Feature | TestFlight | Play Store CLI |
|---------|------------|----------------|
| Setup Time | 5 minutes | 10 minutes |
| Deployment Time | 15-30 minutes | 5-10 minutes |
| Authentication | Apple ID | Service Account |
| Tracks | Internal/External | Internal/Closed/Production |
| Automation | Manual upload | Fully automated |
| Monitoring | App Store Connect | CLI commands |

## 🎉 Benefits

### For Developers
- **Automated Deployments**: No manual uploads required
- **Fast Iteration**: Deploy in minutes, not hours
- **Consistent Process**: Same process every time
- **Error Handling**: Comprehensive error checking

### For Testers
- **Quick Access**: Receive builds faster
- **Easy Installation**: Direct from Play Store
- **Feedback Integration**: Built-in feedback system
- **Version Management**: Clear version tracking

## 📞 Support

If you encounter issues:

1. **Check the logs**: Look for detailed error messages
2. **Verify setup**: Run `./test_auth.sh` to check authentication
3. **Review configuration**: Ensure all paths and credentials are correct
4. **Check permissions**: Verify service account has correct permissions

## 🔄 Updates

To update the deployment system:

1. **Pull latest changes**: `git pull origin main`
2. **Test authentication**: `./test_auth.sh`
3. **Deploy test build**: `./quick_deploy.sh`

---

**Happy Deploying! 🚀**
