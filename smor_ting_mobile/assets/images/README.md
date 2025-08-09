# Images Directory

## Adding the Smor-Ting Logo

1. Save your custom logo image as `smor_ting_logo.png`
2. Place it in this directory (`assets/images/`)
3. The logo should be a square image, ideally 512x512 pixels or larger
4. The app will automatically scale it to fit the 120x120 container

## Current Logo
The app currently uses a fallback handyman icon if the custom logo is not found.

## File Structure
```
assets/
  images/
    smor_ting_logo.png  <- Add your logo here
    README.md
```

## Usage in Code
The logo is referenced in `lib/features/auth/presentation/pages/landing_page.dart`:
```dart
Image.asset(
  'assets/images/smor_ting_logo.png',
  width: 80,
  height: 80,
  fit: BoxFit.contain,
)
``` 