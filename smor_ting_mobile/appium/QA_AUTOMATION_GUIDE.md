# 📱 Smor-Ting QA Automation Guide

## 🎯 Overview

This comprehensive guide covers the complete QA automation setup for the Smor-Ting mobile application, including local development testing and CI/CD integration with deployment gates.

## 🏗️ Architecture

```
smor-ting/
├── .github/workflows/           # CI/CD workflows
│   ├── qa-automation.yml       # Main QA pipeline
│   └── deployment-gate.yml     # Deployment control
├── smor_ting_mobile/appium/    # QA automation framework
│   ├── config/                 # Configuration management
│   ├── tests/                  # Test suites
│   ├── scripts/                # Automation scripts
│   └── reports/                # Test reports & artifacts
└── smor_ting_backend/tests/    # Backend API tests
```

## 🚀 Quick Start

### Prerequisites

- **Node.js** 16+ (for Appium)
- **Python** 3.8+ (for test framework)
- **Java JDK** 11+ (for Android testing)
- **Android Studio** (for Android SDK)
- **Xcode** (for iOS testing - macOS only)
- **Flutter** 3.16+

### 1. Setup Local Environment

```bash
# Navigate to appium directory
cd smor_ting_mobile/appium

# Run setup script
./scripts/setup_appium.sh

# Verify installation
appium doctor
```

### 2. Build Flutter App

```bash
# For Android
cd ../
flutter build apk --debug

# For iOS (macOS only)
flutter build ios --simulator --debug
```

### 3. Run Tests

```bash
# Run all Android tests
./scripts/run_android_tests.sh

# Run specific test suite
./scripts/run_android_tests.sh --suite auth --environment local

# Run tests in parallel
./scripts/run_android_tests.sh --parallel

# iOS tests (macOS only)
./scripts/run_ios_tests.sh
```

## 📋 Test Cases Documentation

### 🔐 Authentication Test Suite

Our QA automation covers **35+ comprehensive test scenarios** across registration and login flows, implementing all the TDD test cases from the backend with mobile-specific validations.

#### Registration Tests (18 scenarios)

| Test ID | Scenario | Expected Result | Error Handling |
|---------|----------|----------------|----------------|
| TC_REG_001 | Successful customer registration | Navigate to customer dashboard | ✅ Success flow |
| TC_REG_002 | Email already exists | Show custom error widget | 🚫 HTTP 409 + UI response |
| TC_REG_003 | Create another user flow | Clear form, hide error | 🔄 Form reset |
| TC_REG_004 | Login from error widget | Navigate to login page | 🔄 Navigation |
| TC_REG_005 | Missing email validation | Show "Email is required" | ❌ Client validation |
| TC_REG_006 | Missing password validation | Show "Password is required" | ❌ Client validation |
| TC_REG_007 | Password too short | Show length requirement | ❌ Client/Server validation |
| TC_REG_008 | Password complexity | Show complexity rules | ❌ Client validation |
| TC_REG_009 | Missing first name | Show "First name is required" | ❌ Client validation |
| TC_REG_010 | Missing last name | Show "Last name is required" | ❌ Client validation |
| TC_REG_011 | Missing phone number | Show "Phone is required" | ❌ Client validation |
| TC_REG_012 | Invalid phone format | Show format hint | ❌ Format validation |
| TC_REG_013 | Missing role selection | Show "Role is required" | ❌ Client validation |
| TC_REG_014 | Invalid role value | Backend error handling | ❌ Server validation |
| TC_REG_015 | Password mismatch | Show "Passwords don't match" | ❌ Client validation |
| TC_REG_016 | Invalid email formats | Prevent submission | ❌ Format validation |
| TC_REG_017 | Loading states | Show/hide loading indicators | 🔄 UI feedback |
| TC_REG_018 | Network error handling | Show network error message | 🌐 Connection issues |

#### Login Tests (12 scenarios)

| Test ID | Scenario | Expected Result | Error Handling |
|---------|----------|----------------|----------------|
| TC_LOGIN_001 | Successful login | Navigate to dashboard | ✅ Success flow |
| TC_LOGIN_002 | Non-existent email | Show "Invalid credentials" | 🚫 HTTP 401 |
| TC_LOGIN_003 | Wrong password | Show "Invalid credentials" | 🚫 HTTP 401 |
| TC_LOGIN_004 | Empty email field | Show validation error | ❌ Client validation |
| TC_LOGIN_005 | Empty password field | Show validation error | ❌ Client validation |
| TC_LOGIN_006 | Invalid email format | Prevent submission | ❌ Format validation |
| TC_LOGIN_007 | Navigation to registration | Show registration form | 🔄 Navigation |
| TC_LOGIN_008 | Loading states | Show/hide loading indicators | 🔄 UI feedback |
| TC_LOGIN_009 | Performance testing | Complete within 3 seconds | ⚡ Performance |
| TC_LOGIN_010 | Network error handling | Show appropriate error | 🌐 Connection issues |
| TC_LOGIN_011 | Multiple failed attempts | Handle gracefully | 🔄 Retry logic |
| TC_LOGIN_012 | Form field behavior | Clear and refill correctly | 🔄 Form handling |

#### UI/UX Tests (8 scenarios)

| Test Category | Scenarios | Coverage |
|---------------|-----------|----------|
| **Loading States** | Button disabling, progress indicators | ⏳ UI feedback |
| **Error Display** | Message formatting, color coding | 🎨 Error UX |
| **Form Validation** | Real-time validation, field highlighting | ✅ Input validation |
| **Accessibility** | Screen reader support, keyboard navigation | ♿ Accessibility |
| **Performance** | Response times, memory usage | ⚡ Performance |
| **Network Handling** | Offline behavior, timeout handling | 🌐 Connectivity |
| **Security** | Input sanitization, token storage | 🔒 Security |
| **Responsive Design** | Different screen sizes, orientations | 📱 Responsive |

### 🎯 Test Data Management

#### Valid Test Users
```json
{
  "customer": {
    "email": "qa_customer@smorting.com",
    "password": "TestPass123!",
    "first_name": "QA",
    "last_name": "Customer",
    "phone": "231777123456",
    "role": "customer"
  },
  "provider": {
    "email": "qa_provider@smorting.com", 
    "password": "ProviderPass123!",
    "first_name": "QA",
    "last_name": "Provider",
    "phone": "231888123456",
    "role": "provider"
  }
}
```

#### Invalid Test Data
```json
{
  "invalid_emails": [
    "", "invalid-email", "test@", "@domain.com"
  ],
  "invalid_passwords": [
    "", "123", "short", "nouppercaseorspecial"
  ],
  "invalid_phones": [
    "", "123", "1234567890123456", "abcdefghijk"
  ]
}
```

## 🔧 Configuration

### Environment Configuration

| Environment | API Base URL | Test Scope | Deployment |
|-------------|--------------|------------|------------|
| **Local** | `http://localhost:8080` | Full suite | Manual |
| **CI** | `https://api.smor-ting.com` | Full suite + Performance | Automatic |
| **Staging** | `https://staging-api.smor-ting.com` | Smoke tests | Automatic |
| **Production** | `https://api.smor-ting.com` | Smoke tests | Gated |

### Platform Configuration

#### Android Configuration
```json
{
  "platformName": "Android",
  "automationName": "UiAutomator2",
  "appPackage": "com.smorting.app.smor_ting_mobile",
  "appActivity": "com.smorting.app.smor_ting_mobile.MainActivity",
  "deviceName": "Android Emulator",
  "platformVersion": "11.0",
  "autoGrantPermissions": true,
  "noReset": false,
  "fullReset": true
}
```

#### iOS Configuration
```json
{
  "platformName": "iOS",
  "automationName": "XCUITest", 
  "bundleId": "com.smorting.app.smor-ting-mobile",
  "deviceName": "iPhone 13",
  "platformVersion": "16.4",
  "isSimulator": true,
  "noReset": false,
  "fullReset": true
}
```

## 🚦 CI/CD Integration

### GitHub Actions Workflow

Our CI/CD pipeline implements **enterprise-grade QA automation** with deployment gates:

#### 1. QA Automation Pipeline (`qa-automation.yml`)

```yaml
# Triggers
- Push to: main, develop, staging
- Pull requests to: main, develop
- Manual dispatch with options

# Matrix Testing
- Android: API levels 30, 33
- iOS: iPhone 13, 14 with iOS 15.5, 16.4
- Test suites: all, auth, registration, login
```

#### 2. Deployment Gate Pipeline (`deployment-gate.yml`)

```yaml
# Quality Gates
- QA tests must pass for production
- Automatic rollback on failure
- Performance thresholds enforced
- Security scans required
```

### Deployment Rules

| Environment | QA Gate | Auto Deploy | Rollback |
|-------------|---------|-------------|----------|
| **Development** | Optional | ✅ Always | Manual |
| **Staging** | Warning only | ✅ Always | Manual |
| **Production** | **Required** | 🚫 Gated | ✅ Automatic |

### Quality Metrics

- **Minimum Pass Rate**: 80%
- **Maximum Failures**: 5 tests
- **Performance Threshold**: < 5 seconds
- **Security Scan**: Required for production

## 📊 Test Reporting

### Report Types

1. **HTML Reports** - Visual test results with screenshots
2. **JUnit XML** - Machine-readable results for CI integration
3. **Allure Reports** - Advanced reporting with trends
4. **Screenshot Gallery** - Visual evidence of failures

### Report Locations

```
reports/
├── android-report.html          # Latest Android test report
├── ios-report.html             # Latest iOS test report  
├── android-junit.xml           # JUnit XML for Android
├── ios-junit.xml               # JUnit XML for iOS
├── screenshots/                # Failure screenshots
│   ├── test_registration_failure_1234567890.png
│   └── test_login_failure_1234567891.png
├── logs/                       # Test execution logs
│   ├── appium.log             # Appium server logs
│   └── test_execution.log     # Python test logs
└── allure-results/            # Allure test data
```

### GitHub Integration

- **PR Comments** - Automatic test result summaries
- **Check Status** - Pass/fail status in PR checks
- **Issue Creation** - Automatic issues for failures
- **Deployment Status** - Approval/blocking notifications

## 🛠️ Advanced Usage

### Running Specific Test Suites

```bash
# Authentication tests only
./scripts/run_android_tests.sh --suite auth

# Registration tests only  
./scripts/run_android_tests.sh --suite registration

# Login tests only
./scripts/run_android_tests.sh --suite login

# With specific markers
./scripts/run_android_tests.sh --markers "smoke"
./scripts/run_android_tests.sh --markers "regression"
```

### Parallel Test Execution

```bash
# Run tests in parallel
./scripts/run_android_tests.sh --parallel

# Specify number of workers
pytest tests/ --numprocesses 4
```

### Custom Device Testing

```bash
# Specific Android device
./scripts/run_android_tests.sh --device "Pixel_4_API_30"

# Specific iOS device  
./scripts/run_ios_tests.sh --device "iPhone 14"
```

### Environment-Specific Testing

```bash
# Local development
./scripts/run_android_tests.sh --environment local

# CI environment
./scripts/run_android_tests.sh --environment ci

# Staging environment
./scripts/run_android_tests.sh --environment staging
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Appium Server Issues
```bash
# Check if Appium is running
curl -s http://127.0.0.1:4723/status

# Kill existing processes
pkill -f appium

# Restart server
appium server --port 4723
```

#### 2. Android Emulator Issues
```bash
# List available AVDs
emulator -list-avds

# Check connected devices
adb devices

# Restart ADB
adb kill-server && adb start-server
```

#### 3. iOS Simulator Issues
```bash
# List available simulators
xcrun simctl list devices

# Boot specific simulator
xcrun simctl boot "iPhone 13"

# Reset simulator
xcrun simctl erase "iPhone 13"
```

#### 4. Flutter Build Issues
```bash
# Clean and rebuild
flutter clean
flutter pub get
flutter build apk --debug
```

### Debug Mode

```bash
# Enable verbose logging
export DEBUG=true
pytest tests/ -v -s --log-cli-level=DEBUG

# Capture screenshots on all steps
export CAPTURE_SCREENSHOTS=true
pytest tests/
```

### Performance Debugging

```bash
# Enable performance monitoring
export PERFORMANCE_MONITORING=true

# Set custom timeouts
export IMPLICIT_WAIT=15
export EXPLICIT_WAIT=30
```

## 📈 Metrics & Monitoring

### Test Execution Metrics

- **Total Test Coverage**: 35+ scenarios
- **Platform Coverage**: Android + iOS
- **API Error Scenarios**: 30+ covered
- **Performance Tests**: Response time < 5s
- **Security Tests**: Input validation + token security

### Quality Metrics

- **Pass Rate Target**: ≥ 80%
- **Execution Time**: < 30 minutes full suite
- **Flaky Test Rate**: < 5%
- **Code Coverage**: Backend API calls

### Monitoring Dashboard

Available in GitHub Actions:
- Test execution trends
- Pass/fail rates by environment
- Performance metrics over time
- Error categorization
- Device/platform breakdown

## 🔒 Security Considerations

### Data Privacy
- Test data is anonymized
- No production data in tests
- Secure token handling
- GDPR compliance verified

### Security Testing
- Input sanitization validation
- XSS prevention testing
- SQL injection protection
- Authentication bypass prevention

## 📚 Additional Resources

### Documentation
- [Appium Documentation](https://appium.io/docs/en/latest/)
- [Pytest Documentation](https://docs.pytest.org/)
- [Flutter Testing Guide](https://docs.flutter.dev/testing)

### Training Materials
- Test case design principles
- Page Object Model patterns
- CI/CD best practices
- Mobile testing strategies

### Support
- **Issues**: GitHub Issues for bug reports
- **Discussions**: Team Slack #qa-automation
- **Updates**: Automatic notifications via GitHub

---

## 🎉 Conclusion

This QA automation framework provides:

✅ **Comprehensive Coverage** - 35+ test scenarios covering all auth flows  
✅ **Enterprise Quality** - CI/CD integration with deployment gates  
✅ **Multi-Platform** - Android and iOS testing support  
✅ **Performance Focused** - Response time and load testing  
✅ **Security Aware** - Input validation and security testing  
✅ **Developer Friendly** - Easy setup and extensive documentation  

**Your builds won't run until all automated QA testing passes** - ensuring consistent quality and preventing broken deployments! 🚀
