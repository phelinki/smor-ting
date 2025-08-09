# 🚀 TestFlight Deployment Guide

This guide will walk you through deploying your Smor-Ting app to TestFlight for beta testing.

## 📋 Prerequisites

### 1. Apple Developer Account
- ✅ Active Apple Developer Program membership ($99/year)
- ✅ Access to App Store Connect
- ✅ Valid provisioning profiles

### 2. Development Environment
- ✅ macOS with latest Xcode installed
- ✅ Flutter SDK (latest stable version)
- ✅ iOS Simulator or physical device for testing

### 3. App Store Connect Setup
- ✅ App registered in App Store Connect
- ✅ Bundle ID: `com.smorting.app.smorTingMobile`
- ✅ App information filled out

## 🔧 Step-by-Step Deployment

### Step 1: Prepare Your Environment

```bash
# Navigate to your mobile app directory
cd smor_ting_mobile

# Make the build script executable
chmod +x scripts/build_testflight.sh

# Clean and get dependencies
flutter clean
flutter pub get
```

### Step 2: Update Configuration

1. **Update API Configuration**
   Edit `lib/core/constants/api_config.dart`:
   ```dart
   // Change this for production
   static const Environment _currentEnvironment = Environment.production;
   static const String _productionBaseUrl = 'https://your-production-server.com/api/v1';
   ```

2. **Update Bundle Identifier** (if needed)
   - Open `ios/Runner.xcodeproj` in Xcode
   - Select Runner target
   - Update Bundle Identifier if needed

3. **Update Team ID**
   Edit `ios/exportOptions.plist`:
   ```xml
   <key>teamID</key>
   <string>YOUR_ACTUAL_TEAM_ID</string>
   ```

### Step 3: Build for TestFlight

```bash
# Run the automated build script
./scripts/build_testflight.sh
```

This script will:
- ✅ Clean previous builds
- ✅ Update version numbers
- ✅ Generate code
- ✅ Build iOS app
- ✅ Create archive
- ✅ Export IPA

### Step 4: Upload to TestFlight

#### Option A: Using Xcode (Recommended)

1. **Open Xcode**
   ```bash
   open ios/Runner.xcworkspace
   ```

2. **Go to Organizer**
   - Window → Organizer
   - Select your archive

3. **Distribute App**
   - Click "Distribute App"
   - Choose "App Store Connect"
   - Select "Upload"
   - Follow the wizard

#### Option B: Using Command Line

```bash
# Upload to App Store Connect
xcrun altool --upload-app \
  --type ios \
  --file ios/build/ipa/Runner.ipa \
  --username "your-apple-id@email.com" \
  --password "app-specific-password"
```

### Step 5: Configure TestFlight

1. **App Store Connect Setup**
   - Go to [App Store Connect](https://appstoreconnect.apple.com)
   - Select your app
   - Go to TestFlight tab

2. **Add Testers**
   - Internal Testers (up to 100)
   - External Testers (up to 10,000)

3. **Build Information**
   - Add build description
   - Include testing instructions
   - Add feedback email

## 📱 Testing Instructions for Testers

### Installation
1. Download TestFlight from App Store
2. Accept invitation via email
3. Install Smor Ting app

### Testing Checklist
- [ ] App launches without crashes
- [ ] User registration works
- [ ] User login works
- [ ] Service browsing works
- [ ] Booking creation works
- [ ] Location services work
- [ ] Camera/photo upload works
- [ ] Offline functionality works
- [ ] Push notifications work (if implemented)

### Feedback Collection
- Use TestFlight's built-in feedback
- Report bugs with screenshots
- Test on different iOS versions
- Test on different device sizes

## 🔧 Troubleshooting

### Common Issues

1. **Code Signing Errors**
   ```
   Error: No provisioning profile found
   ```
   **Solution**: Update provisioning profiles in Xcode

2. **Archive Fails**
   ```
   Error: Build failed
   ```
   **Solution**: Check for compilation errors, clean build

3. **Upload Rejected**
   ```
   Error: Invalid bundle
   ```
   **Solution**: Verify bundle ID and certificates

4. **TestFlight Build Issues**
   ```
   Error: Build processing failed
   ```
   **Solution**: Check app permissions and capabilities

### Debug Steps

1. **Check Certificates**
   ```bash
   security find-identity -v -p codesigning
   ```

2. **Verify Provisioning Profiles**
   ```bash
   ls ~/Library/MobileDevice/Provisioning\ Profiles/
   ```

3. **Clean Build**
   ```bash
   flutter clean
   flutter pub get
   cd ios && rm -rf build && cd ..
   ```

## 📊 Monitoring & Analytics

### TestFlight Analytics
- Track crash reports
- Monitor app performance
- Review tester feedback
- Track installation rates

### App Store Connect Metrics
- Build processing time
- TestFlight usage
- Crash reports
- Performance metrics

## 🔄 Continuous Deployment

### Automated Build Script
```bash
# Update version
VERSION="1.0.1"
BUILD_NUMBER="2"

# Run build
./scripts/build_testflight.sh

# Upload automatically (optional)
xcrun altool --upload-app \
  --type ios \
  --file ios/build/ipa/Runner.ipa \
  --username "$APPLE_ID" \
  --password "$APP_SPECIFIC_PASSWORD"
```

### Version Management
- Increment version for each TestFlight build
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Document changes in build notes

## 📋 Pre-Launch Checklist

### Technical Requirements
- [ ] App passes all App Store Review Guidelines
- [ ] No crashes on launch
- [ ] All features work as expected
- [ ] Performance is acceptable
- [ ] Privacy policy is implemented
- [ ] Terms of service are included

### Content Requirements
- [ ] App icon is high quality
- [ ] Screenshots are up to date
- [ ] App description is complete
- [ ] Keywords are optimized
- [ ] Support URL is working

### Legal Requirements
- [ ] Privacy policy is accessible
- [ ] Terms of service are included
- [ ] Data collection is disclosed
- [ ] Third-party libraries are listed

## 🎯 Best Practices

### For TestFlight
1. **Regular Updates**: Upload new builds frequently
2. **Clear Instructions**: Provide detailed testing instructions
3. **Quick Response**: Respond to tester feedback promptly
4. **Diverse Testing**: Test on different devices and iOS versions

### For Production
1. **Thorough Testing**: Test all features before App Store submission
2. **Performance**: Ensure app performs well on older devices
3. **Accessibility**: Implement accessibility features
4. **Localization**: Consider localizing for Liberia

## 📞 Support Resources

- **Apple Developer Documentation**: https://developer.apple.com/testflight/
- **App Store Connect Help**: https://help.apple.com/app-store-connect/
- **Flutter iOS Deployment**: https://docs.flutter.dev/deployment/ios
- **TestFlight Guidelines**: https://developer.apple.com/app-store/review/guidelines/

---

🎉 **Your Smor-Ting app is ready for TestFlight deployment!**

Follow this guide step by step, and your testers will be able to provide valuable feedback to improve your app before the official App Store launch. 