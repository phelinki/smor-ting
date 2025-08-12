# Appium QA Automation Setup for Smor-Ting

## Overview
This guide provides comprehensive instructions for setting up Appium automation testing for the Smor-Ting mobile application, covering both iOS and Android platforms.

## Prerequisites

### System Requirements
- **macOS**: 10.15+ (for iOS testing)
- **Node.js**: 16.0+ 
- **Java**: JDK 8 or 11
- **Android Studio**: Latest version
- **Xcode**: 12.0+ (for iOS testing)
- **Python**: 3.8+ (for test scripts)

### Required Tools
1. **Appium Server**: 2.0+
2. **Appium Inspector**: For element identification
3. **Android SDK**: Platform tools and build tools
4. **iOS Simulator**: or physical iOS device
5. **Android Emulator**: or physical Android device

## Installation

### 1. Install Node.js and Appium
```bash
# Install Node.js (if not already installed)
brew install node

# Install Appium globally
npm install -g appium@next

# Install Appium drivers
appium driver install uiautomator2  # For Android
appium driver install xcuitest      # For iOS

# Verify installation
appium --version
```

### 2. Install Appium Inspector
```bash
# Download from GitHub releases
https://github.com/appium/appium-inspector/releases

# Or install via npm
npm install -g @appium/inspector
```

### 3. Setup Android Environment
```bash
# Add to ~/.zshrc or ~/.bash_profile
export ANDROID_HOME=/Users/$(whoami)/Library/Android/sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools
export PATH=$PATH:$ANDROID_HOME/tools
export PATH=$PATH:$ANDROID_HOME/tools/bin

# Reload profile
source ~/.zshrc
```

### 4. Setup iOS Environment (macOS only)
```bash
# Install iOS dependencies
npm install -g ios-deploy
npm install -g ios-sim

# For real device testing
brew install libimobiledevice
brew install ideviceinstaller
```

### 5. Install Python Testing Framework
```bash
# Install Python dependencies
pip3 install appium-python-client
pip3 install pytest
pip3 install pytest-html
pip3 install allure-pytest
pip3 install selenium
```

## Project Structure

```
smor_ting_mobile/
├── appium/
│   ├── tests/
│   │   ├── android/
│   │   │   ├── test_auth_flow.py
│   │   │   ├── test_registration.py
│   │   │   └── test_login.py
│   │   ├── ios/
│   │   │   ├── test_auth_flow.py
│   │   │   ├── test_registration.py
│   │   │   └── test_login.py
│   │   └── common/
│   │       ├── base_test.py
│   │       ├── page_objects/
│   │       └── utils/
│   ├── config/
│   │   ├── android_caps.json
│   │   ├── ios_caps.json
│   │   └── test_config.json
│   ├── scripts/
│   │   ├── start_appium.sh
│   │   ├── run_android_tests.sh
│   │   └── run_ios_tests.sh
│   ├── reports/
│   └── requirements.txt
```

## Configuration Files

### Android Capabilities (config/android_caps.json)
```json
{
  "platformName": "Android",
  "platformVersion": "11.0",
  "deviceName": "Android Emulator",
  "app": "/path/to/app-debug.apk",
  "automationName": "UiAutomator2",
  "appPackage": "com.smorting.app.smor_ting_mobile",
  "appActivity": "com.smorting.app.smor_ting_mobile.MainActivity",
  "noReset": false,
  "fullReset": false,
  "newCommandTimeout": 300,
  "androidInstallTimeout": 90000,
  "adbExecTimeout": 20000,
  "autoGrantPermissions": true,
  "settings": {
    "waitForIdleTimeout": 100,
    "waitForSelectorTimeout": 10000
  }
}
```

### iOS Capabilities (config/ios_caps.json)
```json
{
  "platformName": "iOS",
  "platformVersion": "15.0",
  "deviceName": "iPhone 13",
  "app": "/path/to/Runner.app",
  "automationName": "XCUITest",
  "bundleId": "com.smorting.app.smorTingMobile",
  "noReset": false,
  "fullReset": false,
  "newCommandTimeout": 300,
  "wdaLaunchTimeout": 60000,
  "wdaConnectionTimeout": 60000,
  "autoAcceptAlerts": false,
  "autoDismissAlerts": false
}
```

### Test Configuration (config/test_config.json)
```json
{
  "appium_server": {
    "host": "127.0.0.1",
    "port": 4723,
    "timeout": 30
  },
  "test_data": {
    "valid_user": {
      "email": "test@smorting.com",
      "password": "TestPass123!",
      "first_name": "Test",
      "last_name": "User",
      "phone": "231777123456"
    },
    "invalid_user": {
      "email": "invalid@email",
      "password": "123",
      "first_name": "",
      "last_name": ""
    }
  },
  "timeouts": {
    "implicit_wait": 10,
    "explicit_wait": 20,
    "page_load": 30
  },
  "environment": {
    "api_base_url": "https://api.smor-ting.com/api/v1",
    "test_environment": "staging"
  }
}
```

## Page Object Model

### Base Page Object (common/page_objects/base_page.py)
```python
from appium.webdriver.common.appiumby import AppiumBy
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException
import json

class BasePage:
    def __init__(self, driver):
        self.driver = driver
        self.wait = WebDriverWait(driver, 20)
        
        # Load test configuration
        with open('config/test_config.json', 'r') as f:
            self.config = json.load(f)
    
    def find_element(self, locator, timeout=None):
        """Find element with explicit wait"""
        wait_time = timeout or self.config['timeouts']['explicit_wait']
        try:
            return WebDriverWait(self.driver, wait_time).until(
                EC.presence_of_element_located(locator)
            )
        except TimeoutException:
            raise Exception(f"Element not found: {locator}")
    
    def find_elements(self, locator, timeout=None):
        """Find multiple elements"""
        wait_time = timeout or self.config['timeouts']['explicit_wait']
        try:
            WebDriverWait(self.driver, wait_time).until(
                EC.presence_of_element_located(locator)
            )
            return self.driver.find_elements(*locator)
        except TimeoutException:
            return []
    
    def click_element(self, locator, timeout=None):
        """Click element with wait"""
        element = self.find_element(locator, timeout)
        element.click()
        return element
    
    def enter_text(self, locator, text, timeout=None):
        """Enter text in input field"""
        element = self.find_element(locator, timeout)
        element.clear()
        element.send_keys(text)
        return element
    
    def get_text(self, locator, timeout=None):
        """Get element text"""
        element = self.find_element(locator, timeout)
        return element.text
    
    def is_element_present(self, locator, timeout=5):
        """Check if element is present"""
        try:
            self.find_element(locator, timeout)
            return True
        except:
            return False
    
    def wait_for_element_to_disappear(self, locator, timeout=None):
        """Wait for element to disappear"""
        wait_time = timeout or self.config['timeouts']['explicit_wait']
        WebDriverWait(self.driver, wait_time).until_not(
            EC.presence_of_element_located(locator)
        )
    
    def scroll_to_element(self, locator):
        """Scroll to element"""
        element = self.find_element(locator)
        self.driver.execute_script("arguments[0].scrollIntoView();", element)
        return element
    
    def take_screenshot(self, name):
        """Take screenshot"""
        self.driver.save_screenshot(f"reports/screenshots/{name}.png")
```

### Registration Page Object (common/page_objects/registration_page.py)
```python
from appium.webdriver.common.appiumby import AppiumBy
from .base_page import BasePage

class RegistrationPage(BasePage):
    # Locators
    EMAIL_INPUT = (AppiumBy.ACCESSIBILITY_ID, "email_input")
    PASSWORD_INPUT = (AppiumBy.ACCESSIBILITY_ID, "password_input")
    CONFIRM_PASSWORD_INPUT = (AppiumBy.ACCESSIBILITY_ID, "confirm_password_input")
    FIRST_NAME_INPUT = (AppiumBy.ACCESSIBILITY_ID, "first_name_input")
    LAST_NAME_INPUT = (AppiumBy.ACCESSIBILITY_ID, "last_name_input")
    PHONE_INPUT = (AppiumBy.ACCESSIBILITY_ID, "phone_input")
    ROLE_DROPDOWN = (AppiumBy.ACCESSIBILITY_ID, "role_dropdown")
    REGISTER_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "register_button")
    
    # Error messages
    EMAIL_ERROR = (AppiumBy.ACCESSIBILITY_ID, "email_error")
    PASSWORD_ERROR = (AppiumBy.ACCESSIBILITY_ID, "password_error")
    EMAIL_EXISTS_ERROR = (AppiumBy.ACCESSIBILITY_ID, "email_exists_error")
    CREATE_ANOTHER_USER_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "create_another_user_button")
    LOGIN_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "login_button")
    
    # Loading indicator
    LOADING_INDICATOR = (AppiumBy.ACCESSIBILITY_ID, "loading_indicator")
    
    def enter_registration_details(self, user_data):
        """Enter all registration details"""
        self.enter_text(self.EMAIL_INPUT, user_data['email'])
        self.enter_text(self.PASSWORD_INPUT, user_data['password'])
        self.enter_text(self.CONFIRM_PASSWORD_INPUT, user_data['password'])
        self.enter_text(self.FIRST_NAME_INPUT, user_data['first_name'])
        self.enter_text(self.LAST_NAME_INPUT, user_data['last_name'])
        self.enter_text(self.PHONE_INPUT, user_data['phone'])
        
    def select_role(self, role):
        """Select user role"""
        self.click_element(self.ROLE_DROPDOWN)
        role_option = (AppiumBy.ACCESSIBILITY_ID, f"role_{role}")
        self.click_element(role_option)
        
    def click_register(self):
        """Click register button"""
        self.click_element(self.REGISTER_BUTTON)
        
    def wait_for_registration_complete(self):
        """Wait for registration to complete"""
        self.wait_for_element_to_disappear(self.LOADING_INDICATOR)
        
    def get_email_error_message(self):
        """Get email error message"""
        return self.get_text(self.EMAIL_ERROR)
        
    def get_password_error_message(self):
        """Get password error message"""
        return self.get_text(self.PASSWORD_ERROR)
        
    def is_email_exists_error_displayed(self):
        """Check if email exists error is displayed"""
        return self.is_element_present(self.EMAIL_EXISTS_ERROR)
        
    def click_create_another_user(self):
        """Click create another user button"""
        self.click_element(self.CREATE_ANOTHER_USER_BUTTON)
        
    def click_login_from_error(self):
        """Click login button from error widget"""
        self.click_element(self.LOGIN_BUTTON)
        
    def validate_registration_form_reset(self):
        """Validate that registration form is reset"""
        email_value = self.find_element(self.EMAIL_INPUT).get_attribute("text")
        password_value = self.find_element(self.PASSWORD_INPUT).get_attribute("text")
        
        assert email_value == "", "Email field should be empty"
        assert password_value == "", "Password field should be empty"
```

### Login Page Object (common/page_objects/login_page.py)
```python
from appium.webdriver.common.appiumby import AppiumBy
from .base_page import BasePage

class LoginPage(BasePage):
    # Locators
    EMAIL_INPUT = (AppiumBy.ACCESSIBILITY_ID, "login_email_input")
    PASSWORD_INPUT = (AppiumBy.ACCESSIBILITY_ID, "login_password_input")
    LOGIN_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "login_button")
    FORGOT_PASSWORD_LINK = (AppiumBy.ACCESSIBILITY_ID, "forgot_password_link")
    
    # Error messages
    INVALID_CREDENTIALS_ERROR = (AppiumBy.ACCESSIBILITY_ID, "invalid_credentials_error")
    EMAIL_VALIDATION_ERROR = (AppiumBy.ACCESSIBILITY_ID, "email_validation_error")
    PASSWORD_VALIDATION_ERROR = (AppiumBy.ACCESSIBILITY_ID, "password_validation_error")
    
    # Loading
    LOGIN_LOADING = (AppiumBy.ACCESSIBILITY_ID, "login_loading")
    
    def enter_credentials(self, email, password):
        """Enter login credentials"""
        self.enter_text(self.EMAIL_INPUT, email)
        self.enter_text(self.PASSWORD_INPUT, password)
        
    def click_login(self):
        """Click login button"""
        self.click_element(self.LOGIN_BUTTON)
        
    def wait_for_login_complete(self):
        """Wait for login to complete"""
        self.wait_for_element_to_disappear(self.LOGIN_LOADING)
        
    def is_invalid_credentials_error_displayed(self):
        """Check if invalid credentials error is displayed"""
        return self.is_element_present(self.INVALID_CREDENTIALS_ERROR)
        
    def get_invalid_credentials_message(self):
        """Get invalid credentials error message"""
        return self.get_text(self.INVALID_CREDENTIALS_ERROR)
```

## Test Implementation

### Authentication Flow Test (tests/common/test_auth_flow.py)
```python
import pytest
import json
from appium import webdriver
from common.page_objects.registration_page import RegistrationPage
from common.page_objects.login_page import LoginPage
from common.base_test import BaseTest

class TestAuthenticationFlow(BaseTest):
    
    def test_successful_registration_flow(self):
        """Test successful user registration"""
        registration_page = RegistrationPage(self.driver)
        
        # Use valid test data
        user_data = self.config['test_data']['valid_user']
        
        # Enter registration details
        registration_page.enter_registration_details(user_data)
        registration_page.select_role('customer')
        registration_page.click_register()
        
        # Wait for registration to complete
        registration_page.wait_for_registration_complete()
        
        # Verify navigation to dashboard or next screen
        # This depends on your app's navigation flow
        
    def test_email_already_exists_error(self):
        """Test email already exists error handling"""
        registration_page = RegistrationPage(self.driver)
        
        # Use existing user data
        existing_user = self.config['test_data']['valid_user']
        
        # Attempt to register with existing email
        registration_page.enter_registration_details(existing_user)
        registration_page.select_role('customer')
        registration_page.click_register()
        
        # Wait for error to appear
        registration_page.wait_for_registration_complete()
        
        # Verify error message is displayed
        assert registration_page.is_email_exists_error_displayed(), \
            "Email exists error should be displayed"
            
    def test_create_another_user_flow(self):
        """Test create another user button functionality"""
        registration_page = RegistrationPage(self.driver)
        
        # Trigger email exists error first
        existing_user = self.config['test_data']['valid_user']
        registration_page.enter_registration_details(existing_user)
        registration_page.click_register()
        registration_page.wait_for_registration_complete()
        
        # Click create another user
        registration_page.click_create_another_user()
        
        # Verify form is reset
        registration_page.validate_registration_form_reset()
        
    def test_login_from_error_widget(self):
        """Test login button from error widget"""
        registration_page = RegistrationPage(self.driver)
        
        # Trigger email exists error
        existing_user = self.config['test_data']['valid_user']
        registration_page.enter_registration_details(existing_user)
        registration_page.click_register()
        registration_page.wait_for_registration_complete()
        
        # Click login button
        registration_page.click_login_from_error()
        
        # Verify navigation to login page
        login_page = LoginPage(self.driver)
        assert login_page.is_element_present(login_page.EMAIL_INPUT), \
            "Should navigate to login page"
            
    def test_invalid_credentials_login(self):
        """Test login with invalid credentials"""
        login_page = LoginPage(self.driver)
        
        # Navigate to login page
        # (depends on your app navigation)
        
        # Enter invalid credentials
        login_page.enter_credentials("invalid@email.com", "wrongpassword")
        login_page.click_login()
        login_page.wait_for_login_complete()
        
        # Verify error message
        assert login_page.is_invalid_credentials_error_displayed(), \
            "Invalid credentials error should be displayed"
            
    @pytest.mark.parametrize("email,password,expected_error", [
        ("", "password123", "email is required"),
        ("test@example.com", "", "password is required"),
        ("invalid-email", "password123", "Please enter a valid email"),
        ("test@example.com", "123", "Password must be at least 6 characters"),
    ])
    def test_registration_validation_errors(self, email, password, expected_error):
        """Test various registration validation errors"""
        registration_page = RegistrationPage(self.driver)
        
        # Enter invalid data
        user_data = {
            'email': email,
            'password': password,
            'first_name': 'Test',
            'last_name': 'User',
            'phone': '231777123456'
        }
        
        registration_page.enter_registration_details(user_data)
        registration_page.click_register()
        
        # Verify appropriate error message
        # Implementation depends on how your app displays validation errors
```

### Base Test Class (common/base_test.py)
```python
import pytest
import json
from appium import webdriver
from appium.options.android import UiAutomator2Options
from appium.options.ios import XCUITestOptions

class BaseTest:
    driver = None
    config = None
    
    @classmethod
    def setup_class(cls):
        """Setup test class"""
        # Load configuration
        with open('config/test_config.json', 'r') as f:
            cls.config = json.load(f)
            
    def setup_method(self):
        """Setup each test method"""
        self.driver = self.create_driver()
        
    def teardown_method(self):
        """Teardown each test method"""
        if self.driver:
            self.driver.quit()
            
    def create_driver(self):
        """Create Appium driver based on platform"""
        platform = self.get_platform()
        
        if platform == 'android':
            return self.create_android_driver()
        elif platform == 'ios':
            return self.create_ios_driver()
        else:
            raise ValueError(f"Unsupported platform: {platform}")
            
    def get_platform(self):
        """Get platform from environment or default to Android"""
        import os
        return os.environ.get('PLATFORM', 'android').lower()
        
    def create_android_driver(self):
        """Create Android driver"""
        with open('config/android_caps.json', 'r') as f:
            caps = json.load(f)
            
        options = UiAutomator2Options().load_capabilities(caps)
        
        return webdriver.Remote(
            f"http://{self.config['appium_server']['host']}:{self.config['appium_server']['port']}",
            options=options
        )
        
    def create_ios_driver(self):
        """Create iOS driver"""
        with open('config/ios_caps.json', 'r') as f:
            caps = json.load(f)
            
        options = XCUITestOptions().load_capabilities(caps)
        
        return webdriver.Remote(
            f"http://{self.config['appium_server']['host']}:{self.config['appium_server']['port']}",
            options=options
        )
```

## Running Tests

### Start Appium Server (scripts/start_appium.sh)
```bash
#!/bin/bash
echo "Starting Appium server..."
appium --port 4723 --log-level info --log ./reports/appium.log &
echo "Appium server started on port 4723"
```

### Run Android Tests (scripts/run_android_tests.sh)
```bash
#!/bin/bash
export PLATFORM=android

# Start Appium server
./scripts/start_appium.sh

# Wait for server to start
sleep 5

# Run tests
pytest tests/ -v \
    --html=reports/android_report.html \
    --self-contained-html \
    --junitxml=reports/android_junit.xml \
    --alluredir=reports/allure-results

# Generate Allure report
allure generate reports/allure-results -o reports/allure-report --clean

echo "Android tests completed. Reports generated in reports/ directory"
```

### Run iOS Tests (scripts/run_ios_tests.sh)
```bash
#!/bin/bash
export PLATFORM=ios

# Start Appium server
./scripts/start_appium.sh

# Wait for server to start
sleep 5

# Run tests
pytest tests/ -v \
    --html=reports/ios_report.html \
    --self-contained-html \
    --junitxml=reports/ios_junit.xml \
    --alluredir=reports/allure-results

# Generate Allure report
allure generate reports/allure-results -o reports/allure-report --clean

echo "iOS tests completed. Reports generated in reports/ directory"
```

## CI/CD Integration

### GitHub Actions Workflow (.github/workflows/appium-tests.yml)
```yaml
name: Appium Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  android-tests:
    runs-on: macos-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        
    - name: Setup Java
      uses: actions/setup-java@v3
      with:
        distribution: 'temurin'
        java-version: '11'
        
    - name: Setup Android SDK
      uses: android-actions/setup-android@v2
      
    - name: Install Appium
      run: |
        npm install -g appium@next
        appium driver install uiautomator2
        
    - name: Setup Python
      uses: actions/setup-python@v4
      with:
        python-version: '3.9'
        
    - name: Install Python dependencies
      run: |
        cd smor_ting_mobile/appium
        pip install -r requirements.txt
        
    - name: Build Flutter app
      run: |
        cd smor_ting_mobile
        flutter build apk --debug
        
    - name: Start Android Emulator
      uses: reactivecircus/android-emulator-runner@v2
      with:
        api-level: 29
        script: |
          cd smor_ting_mobile/appium
          chmod +x scripts/run_android_tests.sh
          ./scripts/run_android_tests.sh
          
    - name: Upload test reports
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: android-test-reports
        path: smor_ting_mobile/appium/reports/
```

## Best Practices

### 1. Element Identification
- Use accessibility IDs for consistent element identification
- Add accessibility labels to Flutter widgets:
```dart
TextField(
  key: Key('email_input'),
  semanticsLabel: 'Email input field',
  // ...
)
```

### 2. Test Data Management
- Use external JSON files for test data
- Implement data-driven testing with parametrize
- Use separate test data for different environments

### 3. Error Handling
- Implement proper exception handling in tests
- Use screenshots for failed tests
- Add retry mechanisms for flaky tests

### 4. Reporting
- Generate HTML reports with pytest-html
- Use Allure for detailed test reporting
- Include screenshots in test reports

### 5. Maintenance
- Regular updates of Appium and drivers
- Keep test locators updated with app changes
- Monitor test stability and execution times

## Troubleshooting

### Common Issues

1. **Element not found**: Check accessibility IDs and locators
2. **Timeout errors**: Increase wait times or check element loading
3. **Driver connection issues**: Verify Appium server is running
4. **App not launching**: Check app path and capabilities
5. **Slow test execution**: Optimize waits and element interactions

### Debug Commands
```bash
# Check connected devices
adb devices

# Check Appium server logs
tail -f reports/appium.log

# Inspect app elements
appium-inspector
```

This comprehensive setup provides full QA automation capabilities for your Smor-Ting application following security, performance, and usability priorities [[memory:5639049]]. For iOS-specific setup and test execution, see `IOS_SETUP.md` and `TEST_PLAN_IOS.md`.
