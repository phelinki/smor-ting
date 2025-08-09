# ⚡ Quick Start: TestFlight Deployment

Get your Smor-Ting app on TestFlight in 10 minutes!

## 🚀 Quick Setup

### 1. Run Setup Script
```bash
cd smor_ting_mobile
./scripts/setup_testflight.sh
```

### 2. Update Production Configuration
Edit `lib/core/constants/api_config.dart`:
```dart
// Change to production
static const Environment _currentEnvironment = Environment.production;
static const String _productionBaseUrl = 'https://your-production-server.com/api/v1';
```

### 3. Update Team ID
Edit `ios/exportOptions.plist`:
```xml
<key>teamID</key>
<string>YOUR_ACTUAL_TEAM_ID</string>
```

### 4. Build for TestFlight
```bash
./scripts/build_testflight.sh
```

### 5. Upload to TestFlight
```bash
# Open Xcode
open ios/Runner.xcworkspace

# Then:
# 1. Window → Organizer
# 2. Select your archive
# 3. Distribute App → App Store Connect
# 4. Follow the wizard
```

## ✅ Success Indicators

Look for these in the build output:
```
✅ Flutter version: Flutter 3.x.x
✅ Xcode version: Xcode 15.x
✅ iOS build test successful!
✅ App configuration looks good!
```

## 🔧 Essential Configuration

### Bundle ID
- **Current**: `com.smorting.app.smorTingMobile`
- **App Name**: Smor Ting
- **Version**: 1.0.0+1

### Required Permissions
- ✅ Camera access
- ✅ Photo library access
- ✅ Location services
- ✅ Microphone (for future features)
- ✅ Face ID/Touch ID

### API Configuration
- ✅ Development: `http://localhost:8080/api/v1`
- ✅ Production: `https://your-production-server.com/api/v1`

## 📱 Testing Checklist

### Core Features
- [ ] App launches without crashes
- [ ] User registration works
- [ ] User login works
- [ ] Service browsing works
- [ ] Booking creation works
- [ ] Location services work
- [ ] Camera/photo upload works
- [ ] Offline functionality works

### Device Testing
- [ ] iPhone 14 Pro (latest)
- [ ] iPhone 12 (common)
- [ ] iPhone SE (small screen)
- [ ] iPad (if supported)

## 🆘 Need Help?

- 📖 Full guide: `TESTFLIGHT_DEPLOYMENT.md`
- 🔧 Setup script: `./scripts/setup_testflight.sh`
- 🚀 Build script: `./scripts/build_testflight.sh`
- 📞 Apple Developer: https://developer.apple.com/testflight/

## 🎯 Pro Tips

1. **Test Thoroughly**: Test all features before uploading
2. **Clear Instructions**: Provide detailed testing instructions to testers
3. **Quick Response**: Respond to tester feedback promptly
4. **Regular Updates**: Upload new builds frequently
5. **Monitor Feedback**: Track crash reports and user feedback

---

🎉 **Your Smor-Ting app is ready for TestFlight!**

Follow these steps and your testers will be able to provide valuable feedback to improve your app before the official App Store launch. 