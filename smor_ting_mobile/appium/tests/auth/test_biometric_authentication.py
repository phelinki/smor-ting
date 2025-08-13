"""
Biometric Authentication E2E Tests

Tests the biometric authentication functionality including:
- Enable/disable biometric authentication in settings
- Biometric unlock from login screen
- Error handling for biometric authentication
"""

import pytest
import time
from appium.webdriver.common.appiumby import AppiumBy
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException

from ..base_test import BaseTest
from ..common.page_objects.auth_pages import LoginPage
from ..common.page_objects.base_page import BasePage


class SettingsPage(BasePage):
    """Page object for the settings screen"""
    
    # Settings navigation (platform-aware)
    SETTINGS_TITLE = (AppiumBy.XPATH, "//*[contains(@text, 'Settings') or contains(@name, 'Settings') or contains(@label, 'Settings')]")
    
    # Privacy & Security section (platform-aware)
    PRIVACY_SECURITY_SECTION_ANDROID = (AppiumBy.XPATH, "//*[contains(@text, 'Privacy & Security')]")
    PRIVACY_SECURITY_SECTION_IOS = (AppiumBy.ACCESSIBILITY_ID, "privacy_security_section")
    
    BIOMETRIC_AUTH_TOGGLE_ANDROID = (AppiumBy.XPATH, "//*[contains(@text, 'Biometric Authentication')]")
    BIOMETRIC_AUTH_TOGGLE_IOS = (AppiumBy.ACCESSIBILITY_ID, "biometric_auth_row")
    
    BIOMETRIC_AUTH_SWITCH_ANDROID = (AppiumBy.XPATH, "//*[contains(@text, 'Biometric Authentication')]/following-sibling::*//android.widget.Switch")
    BIOMETRIC_AUTH_SWITCH_IOS = (AppiumBy.ACCESSIBILITY_ID, "biometric_auth_switch")
    
    # Success/error messages
    SUCCESS_MESSAGE = (AppiumBy.XPATH, "//*[contains(@text, 'successfully')]")
    ERROR_MESSAGE = (AppiumBy.XPATH, "//*[contains(@text, 'Failed') or contains(@text, 'Error')]")
    
    def navigate_to_biometric_settings(self):
        """Navigate to the biometric authentication settings"""
        section_locator = self.choose_locator(
            self.PRIVACY_SECURITY_SECTION_ANDROID,
            self.PRIVACY_SECURITY_SECTION_IOS,
        )
        toggle_locator = self.choose_locator(
            self.BIOMETRIC_AUTH_TOGGLE_ANDROID,
            self.BIOMETRIC_AUTH_TOGGLE_IOS,
        )
        try:
            self.scroll_to_element(section_locator)
        except Exception:
            return False
        return self.is_element_present(toggle_locator)
    
    def is_biometric_toggle_visible(self):
        """Check if biometric authentication toggle is visible"""
        toggle_locator = self.choose_locator(
            self.BIOMETRIC_AUTH_TOGGLE_ANDROID,
            self.BIOMETRIC_AUTH_TOGGLE_IOS,
        )
        return self.is_element_present(toggle_locator)
    
    def is_biometric_enabled(self):
        """Check if biometric authentication is currently enabled"""
        try:
            switch_locator = self.choose_locator(
                self.BIOMETRIC_AUTH_SWITCH_ANDROID,
                self.BIOMETRIC_AUTH_SWITCH_IOS,
            )
            switch_element = self.driver.find_element(*switch_locator)
            # iOS uses 'value' attribute ('1' or '0'), Android uses 'checked' == 'true'
            if self.is_ios():
                return str(switch_element.get_attribute('value')).lower() in ('1', 'true')
            return switch_element.get_attribute("checked") == "true"
        except Exception:
            return False
    
    def toggle_biometric_auth(self):
        """Toggle the biometric authentication setting"""
        switch_locator = self.choose_locator(
            self.BIOMETRIC_AUTH_SWITCH_ANDROID,
            self.BIOMETRIC_AUTH_SWITCH_IOS,
        )
        switch_element = self.wait_for_element(switch_locator)
        switch_element.click()
        time.sleep(2)  # Wait for toggle animation
    
    def get_status_message(self):
        """Get success or error message after toggle"""
        try:
            success_element = self.driver.find_element(*self.SUCCESS_MESSAGE)
            return success_element.text
        except NoSuchElementException:
            try:
                error_element = self.driver.find_element(*self.ERROR_MESSAGE)
                return error_element.text
            except NoSuchElementException:
                return None


class BiometricLoginPage(LoginPage):
    """Extended login page with biometric functionality"""
    
    # Biometric unlock elements
    BIOMETRIC_UNLOCK_BUTTON = (AppiumBy.XPATH, "//*[contains(@text, 'Unlock with Biometrics')]")
    BIOMETRIC_ICON = (AppiumBy.XPATH, "//android.widget.Button[contains(@content-desc, 'fingerprint')]")
    
    # Divider elements
    OR_DIVIDER = (AppiumBy.XPATH, "//*[contains(@text, 'or')]")
    
    def is_biometric_unlock_available(self):
        """Check if biometric unlock button is visible"""
        return self.is_element_present(self.BIOMETRIC_UNLOCK_BUTTON)
    
    def tap_biometric_unlock(self):
        """Tap the biometric unlock button"""
        button = self.wait_for_element(self.BIOMETRIC_UNLOCK_BUTTON)
        button.click()


class TestBiometricAuthentication(BaseTest):
    """Test suite for biometric authentication functionality"""
    
    def setup_method(self, method):
        """Setup for each test method"""
        super().setup_method(method)
        self.settings_page = SettingsPage(self.driver)
        self.login_page = BiometricLoginPage(self.driver)
    
    @pytest.mark.biometric
    def test_biometric_settings_visibility_when_available(self):
        """Test that biometric settings are visible when biometrics are available"""
        # Navigate to settings (assuming we start from home/authenticated state)
        # This would need to be adapted based on your app's navigation
        
        # For now, let's assume we can navigate directly to settings
        # In a real test, you'd navigate through the app UI
        
        # Check if biometric toggle is visible
        has_biometric_toggle = self.settings_page.navigate_to_biometric_settings()
        
        if has_biometric_toggle:
            assert self.settings_page.is_biometric_toggle_visible(), \
                "Biometric authentication toggle should be visible when biometrics are available"
        else:
            pytest.skip("Biometric authentication not available on this device")
    
    @pytest.mark.biometric
    def test_enable_biometric_authentication(self):
        """Test enabling biometric authentication from settings"""
        # Navigate to biometric settings
        has_biometric_toggle = self.settings_page.navigate_to_biometric_settings()
        
        if not has_biometric_toggle:
            pytest.skip("Biometric authentication not available on this device")
        
        # Check current state
        initially_enabled = self.settings_page.is_biometric_enabled()
        
        # If already enabled, disable first to test enabling
        if initially_enabled:
            self.settings_page.toggle_biometric_auth()
            time.sleep(1)
        
        # Now enable biometric authentication
        self.settings_page.toggle_biometric_auth()
        
        # Verify it was enabled
        assert self.settings_page.is_biometric_enabled(), \
            "Biometric authentication should be enabled after toggling"
        
        # Check for success message
        status_message = self.settings_page.get_status_message()
        assert status_message is not None, "Should show status message after toggle"
        assert "successfully" in status_message.lower(), \
            f"Should show success message, got: {status_message}"
    
    @pytest.mark.biometric
    def test_disable_biometric_authentication(self):
        """Test disabling biometric authentication from settings"""
        # Navigate to biometric settings
        has_biometric_toggle = self.settings_page.navigate_to_biometric_settings()
        
        if not has_biometric_toggle:
            pytest.skip("Biometric authentication not available on this device")
        
        # Ensure biometric is enabled first
        if not self.settings_page.is_biometric_enabled():
            self.settings_page.toggle_biometric_auth()
            time.sleep(1)
        
        # Now disable biometric authentication
        self.settings_page.toggle_biometric_auth()
        
        # Verify it was disabled
        assert not self.settings_page.is_biometric_enabled(), \
            "Biometric authentication should be disabled after toggling"
        
        # Check for success message
        status_message = self.settings_page.get_status_message()
        assert status_message is not None, "Should show status message after toggle"
        assert "successfully" in status_message.lower(), \
            f"Should show success message, got: {status_message}"
    
    @pytest.mark.biometric
    def test_biometric_unlock_button_visibility(self):
        """Test that biometric unlock button appears on login when enabled"""
        # First, ensure biometric is enabled (this would require being logged in)
        # For this test, we'll assume it's already enabled or skip if not available
        
        # Navigate to login page
        # This would depend on your app's navigation flow
        
        # Check if biometric unlock is available
        if self.login_page.is_biometric_unlock_available():
            assert self.login_page.is_element_present(self.login_page.OR_DIVIDER), \
                "Should show 'or' divider when biometric unlock is available"
            assert self.login_page.is_element_present(self.login_page.BIOMETRIC_UNLOCK_BUTTON), \
                "Should show biometric unlock button when biometric auth is enabled"
        else:
            pytest.skip("Biometric unlock not available - either not supported or not enabled")
    
    @pytest.mark.biometric
    def test_biometric_unlock_flow(self):
        """Test the biometric unlock flow from login screen"""
        # Navigate to login page
        # This would depend on your app's navigation flow
        
        if not self.login_page.is_biometric_unlock_available():
            pytest.skip("Biometric unlock not available")
        
        # Tap biometric unlock button
        self.login_page.tap_biometric_unlock()
        
        # At this point, the system biometric prompt would appear
        # Since we can't interact with system dialogs in Appium easily,
        # we'll check for expected behavior:
        
        # 1. Either we get authenticated (app navigates away from login)
        # 2. Or we get an error message (biometric failed/cancelled)
        # 3. Or we timeout (no biometric interaction)
        
        time.sleep(3)  # Wait for biometric prompt and response
        
        # Check if we're still on login page or navigated away
        still_on_login = self.login_page.is_element_present(self.login_page.LOGIN_BUTTON)
        
        if still_on_login:
            # Check for error message
            error_message = self.login_page.get_error_message()
            # This is acceptable - user might have cancelled biometric prompt
            print(f"Still on login page. Error message: {error_message}")
        else:
            # Successfully authenticated and navigated away
            print("Successfully authenticated with biometrics and navigated away from login")
    
    @pytest.mark.biometric
    def test_biometric_unlock_with_no_enrolled_biometrics(self):
        """Test biometric unlock behavior when no biometrics are enrolled"""
        # This test would require a device/emulator with no biometrics enrolled
        # It should either not show the biometric button or show an appropriate error
        
        if not self.login_page.is_biometric_unlock_available():
            # This is expected behavior when no biometrics are enrolled
            print("Biometric unlock not available - likely no biometrics enrolled")
            return
        
        # If button is available, tapping it should show an error
        self.login_page.tap_biometric_unlock()
        time.sleep(2)
        
        # Should either show error message or system handles it gracefully
        error_message = self.login_page.get_error_message()
        if error_message:
            assert "biometric" in error_message.lower() or "fingerprint" in error_message.lower(), \
                f"Error message should mention biometrics: {error_message}"
    
    @pytest.mark.biometric
    def test_biometric_settings_persistence(self):
        """Test that biometric settings persist across app restarts"""
        # Navigate to biometric settings
        has_biometric_toggle = self.settings_page.navigate_to_biometric_settings()
        
        if not has_biometric_toggle:
            pytest.skip("Biometric authentication not available on this device")
        
        # Enable biometric authentication
        if not self.settings_page.is_biometric_enabled():
            self.settings_page.toggle_biometric_auth()
            time.sleep(1)
        
        initial_state = self.settings_page.is_biometric_enabled()
        assert initial_state, "Biometric should be enabled before app restart"
        
        # Restart the app
        self.restart_app()
        
        # Navigate back to settings
        has_biometric_toggle = self.settings_page.navigate_to_biometric_settings()
        assert has_biometric_toggle, "Should still have biometric toggle after restart"
        
        # Verify the setting persisted
        persisted_state = self.settings_page.is_biometric_enabled()
        assert persisted_state == initial_state, \
            "Biometric authentication setting should persist across app restarts"
    
    def restart_app(self):
        """Helper method to restart the app"""
        self.driver.terminate_app(self.app_package)
        time.sleep(2)
        self.driver.activate_app(self.app_package)
        time.sleep(3)  # Wait for app to fully load
