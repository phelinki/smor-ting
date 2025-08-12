# Quick Start Guide - Appium QA for Smor-Ting

## 🚀 Getting Started in 5 Minutes

### Prerequisites
- macOS (for iOS testing) or Linux/Windows (Android only)
- Node.js 16+ installed
- Python 3.8+ installed

### Step 1: Initial Setup
```bash
# Navigate to the appium directory
cd smor_ting_mobile/appium

# Run the setup script
./scripts/setup_appium.sh
```

This script will:
- ✅ Install Appium and required drivers
- ✅ Install Python dependencies
- ✅ Check your system configuration
- ✅ Create necessary directories and config files

### Step 2: Install Android SDK (for Android testing)

#### Option A: Install Android Studio (Recommended)
1. Download from: https://developer.android.com/studio
2. Install and open Android Studio
3. Go to SDK Manager and install:
   - Android SDK Platform-Tools
   - Android SDK Build-Tools
   - At least one Android platform (API 29+)

#### Option B: Command Line Tools Only
```bash
# Download Android command line tools
# Add to your shell profile (~/.zshrc or ~/.bash_profile):
export ANDROID_HOME=$HOME/Library/Android/sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools
export PATH=$PATH:$ANDROID_HOME/tools
```

### Step 3: Build Your Flutter App
```bash
# For Android
cd ..  # Go back to smor_ting_mobile directory
flutter build apk --debug

# For iOS (macOS only)
flutter build ios --simulator --debug
```

### Step 4: Run Tests

#### Android Tests
```bash
cd appium
./scripts/run_android_tests.sh
```

#### iOS Tests (macOS only)
```bash
./scripts/run_ios_tests.sh
```

## 🔧 What Each Script Does

### `setup_appium.sh`
- Installs and configures Appium
- Checks system requirements
- Creates project structure
- Installs Python dependencies

### `run_android_tests.sh`
- Checks Android environment
- Starts Android emulator if needed
- Builds Flutter app if needed
- Starts Appium server
- Runs test suite
- Generates HTML reports

### `run_ios_tests.sh`
- Checks iOS environment
- Starts iOS Simulator if needed
- Builds Flutter app if needed
- Runs iOS test suite

### `start_appium.sh`
- Starts Appium server on port 4723
- Configures logging
- Health check verification

## 📱 Device Setup

### Android Emulator
```bash
# List available AVDs
emulator -list-avds

# Start specific emulator
emulator -avd Pixel_4_API_30

# Or let the script auto-start one
./scripts/run_android_tests.sh
```

### iOS Simulator
```bash
# List available simulators
xcrun simctl list devices

# Boot a specific simulator
xcrun simctl boot "iPhone 13"

# Or let the script auto-start one
./scripts/run_ios_tests.sh
```

## 📊 Test Reports

After running tests, you'll find reports in:
```
appium/reports/
├── android_report.html      # HTML test report
├── ios_report.html         # iOS HTML report
├── android_junit.xml       # JUnit XML format
├── allure-report/          # Allure reports (if installed)
├── screenshots/            # Test screenshots
└── appium.log             # Appium server logs
```

## 🐛 Troubleshooting

### Common Issues

#### "Appium server failed to start"
```bash
# Kill any existing Appium processes
pkill -f appium

# Check port availability
lsof -i :4723

# Restart manually
appium server --port 4723
```

#### "No devices connected"
```bash
# For Android
adb devices

# Start emulator manually
emulator -avd YOUR_AVD_NAME

# For iOS
xcrun simctl list devices
```

#### "App not found"
```bash
# Rebuild Flutter app
cd ../smor_ting_mobile
flutter clean
flutter build apk --debug          # Android
flutter build ios --simulator      # iOS
```

#### "Python dependencies missing"
```bash
pip3 install -r requirements.txt

# Or install specific packages
pip3 install appium-python-client pytest pytest-html
```

### Environment Variables
Add these to your shell profile (`~/.zshrc` or `~/.bash_profile`):

```bash
# Android (required for Android testing)
export ANDROID_HOME=$HOME/Library/Android/sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools
export PATH=$PATH:$ANDROID_HOME/tools

# Java (required for Android)
export JAVA_HOME=/Library/Java/JavaVirtualMachines/openjdk-11.jdk/Contents/Home

# Node.js (if installed via nvm)
export PATH=$PATH:$HOME/.nvm/versions/node/v18.17.0/bin
```

## 🧪 Writing Your First Test

Create `tests/my_first_test.py`:

```python
import pytest
from appium import webdriver
from appium.options.android import UiAutomator2Options

def test_app_launches():
    """Test that the Smor-Ting app launches successfully"""
    
    # Configure Android capabilities
    options = UiAutomator2Options()
    options.platform_name = "Android"
    options.device_name = "Android Emulator"
    options.app = "/path/to/your/app.apk"
    options.automation_name = "UiAutomator2"
    
    # Create driver
    driver = webdriver.Remote(
        "http://127.0.0.1:4723",
        options=options
    )
    
    try:
        # Test that app launches
        assert driver.current_package == "com.smorting.app.smor_ting_mobile"
        print("✅ App launched successfully!")
        
    finally:
        driver.quit()
```

Run your test:
```bash
pytest tests/my_first_test.py -v
```

## 📚 Next Steps

1. **Study the Test Cases**: Review `TEST_CASES.md` for comprehensive test scenarios
2. **Read the API Documentation**: Check `../docs/API_ERRORS.md` for all error handling
3. **Implement Page Objects**: Use the page object pattern in `tests/common/page_objects/`
4. **Add CI/CD**: Set up automated testing in your build pipeline
5. **Monitor Results**: Set up test result monitoring and alerting

## 🆘 Getting Help

- **Appium Documentation**: https://appium.io/docs/en/latest/
- **Flutter Testing**: https://docs.flutter.dev/testing
- **Pytest Documentation**: https://docs.pytest.org/
- **GitHub Issues**: Report bugs or request features

---

**🎉 You're ready to start automated testing!**

The setup provides a solid foundation for comprehensive mobile app testing with proper error handling, reporting, and CI/CD integration.
