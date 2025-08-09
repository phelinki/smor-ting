# 🍎 App Store Submission Guide

Complete guide for submitting Smor-Ting to the App Store.

## 🚀 Quick Start

```bash
# Run the submission guide
./scripts/submit_appstore.sh

# Build for App Store
./scripts/build_testflight.sh

# Open Xcode for submission
open ios/Runner.xcworkspace
```

## 📋 Pre-Submission Checklist

### ✅ Technical Requirements
- [ ] App passes all App Store Review Guidelines
- [ ] No crashes on launch
- [ ] All features work as expected
- [ ] Performance is acceptable
- [ ] Privacy policy is implemented
- [ ] Terms of service are included

### ✅ Content Requirements
- [ ] App icon is high quality (1024x1024)
- [ ] Screenshots are up to date
- [ ] App description is complete
- [ ] Keywords are optimized
- [ ] Support URL is working

### ✅ Legal Requirements
- [ ] Privacy policy is accessible
- [ ] Terms of service are included
- [ ] Data collection is disclosed
- [ ] Third-party libraries are listed

## 🔧 Step-by-Step Submission

### Step 1: Prepare Your App

1. **Update Production API URL**
   ```dart
   // Edit lib/core/constants/api_config.dart
   static const Environment _currentEnvironment = Environment.production;
   static const String _productionBaseUrl = 'https://your-production-server.com/api/v1';
   ```

2. **Update App Version**
   ```yaml
   # Edit pubspec.yaml
   version: 1.0.0+1
   ```

3. **Build Production App**
   ```bash
   ./scripts/build_testflight.sh
   ```

### Step 2: App Store Connect Setup

1. **Go to App Store Connect**
   - Visit: https://appstoreconnect.apple.com
   - Sign in with your Apple Developer account

2. **Create New App**
   - Click '+' → New App
   - Platform: iOS
   - Name: Smor Ting
   - Bundle ID: `com.smorting.app.smorTingMobile`
   - SKU: `smor-ting-ios`
   - User Access: Full Access

3. **Fill App Information**
   - App Name: Smor Ting
   - Subtitle: Handyman & Service Marketplace
   - Keywords: handyman,services,marketplace,liberia
   - Description: [Your app description]
   - Support URL: [Your support website]
   - Marketing URL: [Your marketing website]

### Step 3: Upload Build

1. **Open Xcode**
   ```bash
   open ios/Runner.xcworkspace
   ```

2. **Archive App**
   - Product → Archive
   - Wait for archive to complete

3. **Upload to App Store Connect**
   - Window → Organizer
   - Select your archive
   - Distribute App → App Store Connect
   - Follow the wizard

### Step 4: App Store Connect Configuration

1. **App Information**
   - Fill out all required fields
   - Add app description
   - Add keywords
   - Add support URL

2. **App Screenshots**
   - iPhone 6.7" Display (1290 x 2796)
   - iPhone 6.5" Display (1242 x 2688)
   - iPhone 5.5" Display (1242 x 2208)
   - iPad Pro 12.9" Display (2048 x 2732)

3. **App Review Information**
   - Demo Account (if required)
   - Demo Password (if required)
   - Notes for Review

### Step 5: Submit for Review

1. **Final Checks**
   - All required fields completed
   - Screenshots uploaded
   - App description complete
   - Privacy policy accessible

2. **Submit for Review**
   - Click 'Submit for Review'
   - Confirm submission
   - Wait for review (1-7 days)

## 📱 Required App Store Assets

### App Icon (Required)
- Size: 1024 x 1024 pixels
- Format: PNG or JPEG
- No transparency
- No rounded corners

### Screenshots (Required)
- Minimum 1 screenshot per device
- Maximum 10 screenshots per device
- No device frames
- No text overlays

### App Description (Required)
- Maximum 4000 characters
- Clear and compelling
- Highlight key features
- Include call-to-action

## 📋 App Store Review Guidelines

### Common Rejection Reasons
- ❌ App crashes on launch
- ❌ Incomplete app functionality
- ❌ Missing privacy policy
- ❌ Poor app performance
- ❌ Inappropriate content
- ❌ Misleading app description
- ❌ Broken links or URLs

### Best Practices
- ✅ Test thoroughly before submission
- ✅ Follow Apple's Human Interface Guidelines
- ✅ Provide clear app description
- ✅ Include privacy policy
- ✅ Respond quickly to review feedback
- ✅ Keep app updated regularly

## 📊 Post-Submission Process

### Review Timeline
- Standard review: 1-7 days
- Expedited review: 24-48 hours
- Re-review: 1-3 days

### Review Status
- Waiting for Review
- In Review
- Ready for Sale
- Rejected (with feedback)

### If Rejected
1. Read rejection feedback carefully
2. Fix the issues mentioned
3. Upload new build
4. Resubmit for review

## 🎯 Launch Preparation

### Pre-Launch Checklist
- [ ] App approved by Apple
- [ ] App Store listing complete
- [ ] Marketing materials ready
- [ ] Support system in place
- [ ] Analytics tracking configured
- [ ] Crash reporting set up
- [ ] User feedback system ready

### Launch Day
1. Monitor App Store Connect
2. Track download metrics
3. Monitor crash reports
4. Respond to user reviews
5. Prepare for updates

## 📝 App Store Listing Content

### App Description Template
```
Smor Ting - Handyman & Service Marketplace

Connect with trusted service providers in Liberia. Find reliable handymen, cleaners, electricians, and more for all your home and business needs.

Key Features:
• Browse local service providers
• Book appointments instantly
• Secure payment processing
• Real-time tracking
• Offline functionality
• User reviews and ratings

Perfect for:
• Homeowners needing repairs
• Businesses requiring services
• Service providers looking for clients
• Anyone seeking reliable local services

Download Smor Ting today and experience the easiest way to find and book trusted services in Liberia!
```

### Keywords
```
handyman,services,marketplace,liberia,home,repair,cleaning,electrician,plumber,carpenter,booking,appointment,local,trusted,reliable
```

### Support Information
- Support URL: [Your support website]
- Marketing URL: [Your marketing website]
- Privacy Policy: [Your privacy policy URL]
- Terms of Service: [Your terms URL]

## 🔧 Troubleshooting

### Common Issues

1. **Build Upload Fails**
   ```
   Error: Invalid bundle
   ```
   **Solution**: Check bundle ID and certificates

2. **App Review Rejected**
   ```
   Error: App crashes on launch
   ```
   **Solution**: Test thoroughly on multiple devices

3. **Metadata Rejected**
   ```
   Error: Inappropriate content
   ```
   **Solution**: Review app description and screenshots

### Debug Steps

1. **Check Certificates**
   ```bash
   security find-identity -v -p codesigning
   ```

2. **Verify Provisioning Profiles**
   ```bash
   ls ~/Library/MobileDevice/Provisioning\ Profiles/
   ```

3. **Test App Thoroughly**
   ```bash
   flutter test
   flutter build ios --release
   ```

## 📞 Support Resources

- **Apple Developer Documentation**: https://developer.apple.com/app-store/
- **App Store Connect Help**: https://help.apple.com/app-store-connect/
- **App Store Review Guidelines**: https://developer.apple.com/app-store/review/guidelines/
- **Flutter iOS Deployment**: https://docs.flutter.dev/deployment/ios

## 🎉 Success Checklist

- [ ] App built successfully
- [ ] App Store Connect account active
- [ ] App information completed
- [ ] Screenshots uploaded
- [ ] App description written
- [ ] Privacy policy included
- [ ] Build uploaded successfully
- [ ] App submitted for review
- [ ] Review feedback addressed (if any)
- [ ] App approved and live

---

🎉 **Your Smor-Ting app is ready for App Store submission!**

Follow this guide step by step for a successful submission. Remember that App Store review can take 1-7 days, so be patient and responsive to any feedback from Apple's review team. 