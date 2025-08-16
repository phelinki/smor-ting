# Google Play Store Upload Guide for Smor-Ting

## Prerequisites ✅
- Google Play Console account: **pkaleewoun@gmail.com**
- $25 developer registration fee paid
- App bundle (.aab) file ready

## Step-by-Step Upload Process

### Step 1: Access Google Play Console
1. Go to [play.google.com/console](https://play.google.com/console)
2. Sign in with **pkaleewoun@gmail.com**
3. If this is your first time, complete the developer registration

### Step 2: Create New App
1. Click **"Create app"** button
2. Fill in the app details:
   - **App name**: Smor-Ting
   - **Default language**: English
   - **App or game**: App
   - **Free or paid**: Free
   - Click **"Create"**

### Step 3: Complete App Setup
1. **App access**: Choose "All users" or "Testing" (recommend "Testing" first)
2. **App category**: Select "Business" or "Lifestyle"
3. **Tags**: Add relevant keywords like "handyman", "services", "Liberia"

### Step 4: Upload App Bundle
1. In the left menu, click **"Production"**
2. Click **"Create new release"**
3. Click **"Upload"** and select your `.aab` file:
   - File location: `build/app/outputs/bundle/release/app-release.aab`
4. Add **Release notes**:
   ```
   Initial release of Smor-Ting - Handyman and Service Marketplace for Liberia
   
   Features:
   - User authentication and registration
   - Service booking and management
   - Real-time location tracking
   - Payment integration
   - Offline-first functionality
   ```
5. Click **"Save"**

### Step 5: Complete Store Listing
1. Go to **"Store listing"** in the left menu
2. Fill in the required information:

#### App Details
- **App title**: Smor-Ting
- **Short description**: "Handyman and Service Marketplace for Liberia"
- **Full description**:
```
Smor-Ting is Liberia's premier handyman and service marketplace, connecting customers with skilled professionals for all their home and business needs.

Key Features:
• Find trusted handymen and service providers in your area
• Book services with real-time tracking and updates
• Secure payment processing with mobile money integration
• Offline functionality for areas with limited connectivity
• User reviews and ratings for quality assurance
• Emergency service requests with priority handling

Whether you need plumbing, electrical work, carpentry, cleaning, or any other service, Smor-Ting makes it easy to find reliable professionals in Liberia.

Download now and experience the convenience of on-demand services at your fingertips!
```

#### Graphics
- **App icon**: 512x512 px PNG (use your existing logo)
- **Feature graphic**: 1024x500 px PNG
- **Screenshots**: Upload 2-8 screenshots of your app
  - Phone screenshots: 1080x1920 px minimum
  - Show key features: login, home screen, booking, payment

### Step 6: Content Rating
1. Go to **"Content rating"**
2. Complete the content rating questionnaire
3. Answer questions about:
   - Violence
   - Sexuality
   - Language
   - Controlled substances
4. Submit for review

### Step 7: Pricing & Distribution
1. Go to **"Pricing & distribution"**
2. **Pricing**: Select "Free"
3. **Countries**: Select "All countries" or specific regions
4. **Content guidelines**: Accept the terms
5. **US export laws**: Accept compliance

### Step 8: App Signing
1. Go to **"Setup"** → **"App signing"**
2. Choose **"Upload key"** (you're uploading your own signed app)
3. Upload your keystore file if prompted

### Step 9: Review and Submit
1. Go back to **"Production"** → **"Releases"**
2. Click **"Review release"**
3. Review all information:
   - App bundle is uploaded
   - Release notes are added
   - Store listing is complete
   - Content rating is approved
4. Click **"Start rollout to Production"**

## Important Notes

### Security
- **Keep your keystore safe**: Store `android/app/keystore/upload-keystore.jks` securely
- **Backup passwords**: Store keystore passwords in a secure location
- **Never share credentials**: Don't commit keystore files to version control

### Review Process
- **Timeline**: 1-3 business days for initial review
- **Email notifications**: Monitor pkaleewoun@gmail.com for updates
- **Common issues**: 
  - App crashes on startup
  - Missing privacy policy
  - Inappropriate content
  - Permission violations

### After Publication
- **Monitor reviews**: Respond to user feedback
- **Track performance**: Use Google Play Console analytics
- **Update regularly**: Keep app updated with new features
- **Maintain quality**: Address bugs and issues promptly

## Troubleshooting

### Upload Issues
- **File too large**: Optimize images and assets
- **Invalid bundle**: Ensure proper signing configuration
- **Permission errors**: Check AndroidManifest.xml permissions

### Review Rejections
- **Policy violations**: Review Google Play policies
- **Technical issues**: Test app thoroughly before submission
- **Content issues**: Ensure appropriate content for all ages

### Contact Support
- **Google Play Console Help**: Use the help section in the console
- **Developer Support**: Contact Google developer support if needed

## Next Steps After Upload

1. **Monitor the review process** (1-3 days)
2. **Prepare marketing materials** for app launch
3. **Set up analytics** to track app performance
4. **Plan future updates** and feature releases
5. **Engage with users** through reviews and feedback

## Success Checklist

- [ ] Google Play Console account created
- [ ] App bundle (.aab) file generated
- [ ] App uploaded to Play Console
- [ ] Store listing completed
- [ ] Content rating approved
- [ ] App submitted for review
- [ ] Review process completed
- [ ] App published to Play Store

Your app will be available on the Google Play Store once the review process is complete!
