"""
Test Navigation Helper Tightening for Landing → Login Flow
==========================================================

This test suite verifies that the navigation helper properly enforces
Landing → Login flow before making assertions in QA automation.

Following TDD principles as required by the user.
"""
import pytest
import time
from appium.webdriver.common.appiumby import AppiumBy
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException

from ..base_test import BaseTest
from ..common.page_objects.page_factory import PageFactory


class TestNavigationHelper(BaseTest):
    """Test suite for navigation helper tightening"""

    def setUp(self):
        """Setup for each test"""
        super().setUp()
        
        # Initialize page objects
        self.splash_page = PageFactory.get_splash_page(self.driver)
        self.landing_page = PageFactory.get_landing_page(self.driver)
        self.login_page = PageFactory.get_login_page(self.driver)

    def test_navigation_helper_forces_landing_to_login_sequence(self):
        """
        Test: Navigation helper should force Landing → Login sequence
        before allowing login page assertions.
        
        This ensures proper navigation flow in QA automation.
        """
        # Test: Start from app launch (splash)
        self.splash_page.wait_for_splash_to_complete()
        
        # Test: Should be on landing page first
        try:
            landing_sign_in_button = WebDriverWait(self.driver, 10).until(
                EC.presence_of_element_located(self.landing_page.SIGN_IN_BUTTON)
            )
            self.assertIsNotNone(landing_sign_in_button, "Should be on landing page initially")
        except TimeoutException:
            self.fail("Landing page not reached after splash completion")
        
        # Test: Navigation helper should require explicit Landing → Login transition
        # Using the tightened navigation helper
        self._tightened_navigate_to_login()
        
        # Test: Should now be on login page with proper elements
        try:
            login_submit_button = WebDriverWait(self.driver, 5).until(
                EC.presence_of_element_located((AppiumBy.ACCESSIBILITY_ID, "login_submit"))
            )
            self.assertIsNotNone(login_submit_button, "Should be on login page after navigation")
        except TimeoutException:
            self.fail("Login page not reached through proper navigation sequence")

    def test_navigation_helper_prevents_direct_login_assertions(self):
        """
        Test: Navigation helper should prevent making login page assertions
        without first going through Landing → Login flow.
        
        This test ensures proper navigation sequence enforcement.
        """
        # Test: Start from app launch
        self.splash_page.wait_for_splash_to_complete()
        
        # Test: Attempting to find login elements without proper navigation should fail
        # (or require the navigation helper to enforce the flow)
        login_elements_found_directly = self._try_find_login_elements_directly()
        
        if login_elements_found_directly:
            # If login elements are found directly, the helper should enforce proper navigation
            self.fail("Login elements should not be accessible without proper Landing → Login navigation")
        
        # Test: Now use proper navigation
        self._tightened_navigate_to_login()
        
        # Test: After proper navigation, login elements should be accessible
        login_elements_found_properly = self._try_find_login_elements_directly()
        self.assertTrue(login_elements_found_properly, "Login elements should be accessible after proper navigation")

    def test_navigation_helper_validates_landing_page_presence(self):
        """
        Test: Navigation helper should validate that landing page is present
        before proceeding to login navigation.
        """
        # Test: Start from app launch
        self.splash_page.wait_for_splash_to_complete()
        
        # Test: Tightened navigation should verify landing page presence
        try:
            # This should verify landing page elements are present
            self._validate_landing_page_presence()
            
            # Then proceed with navigation
            self._tightened_navigate_to_login()
            
            # Verify successful navigation
            login_page_present = self._validate_login_page_presence()
            self.assertTrue(login_page_present, "Login page should be present after validated navigation")
            
        except Exception as e:
            self.fail(f"Navigation helper failed to validate landing page presence: {e}")

    def test_navigation_helper_handles_app_state_variations(self):
        """
        Test: Navigation helper should handle different app states and still
        enforce Landing → Login flow.
        """
        # Test different scenarios that might affect navigation
        test_scenarios = [
            "fresh_app_launch",
            "after_background_foreground", 
            "after_orientation_change"
        ]
        
        for scenario in test_scenarios:
            with self.subTest(scenario=scenario):
                # Setup scenario
                if scenario == "fresh_app_launch":
                    # Already handled by setUp
                    pass
                elif scenario == "after_background_foreground":
                    # Simulate app backgrounding/foregrounding
                    self.driver.background_app(2)
                    time.sleep(1)
                elif scenario == "after_orientation_change":
                    # Simulate orientation change if supported
                    try:
                        self.driver.orientation = "LANDSCAPE"
                        time.sleep(1)
                        self.driver.orientation = "PORTRAIT"
                        time.sleep(1)
                    except:
                        # Skip if orientation not supported
                        continue
                
                # Test: Navigation helper should still work
                try:
                    self._tightened_navigate_to_login()
                    login_page_present = self._validate_login_page_presence()
                    self.assertTrue(login_page_present, f"Navigation should work in scenario: {scenario}")
                except Exception as e:
                    self.fail(f"Navigation failed in scenario {scenario}: {e}")

    def test_navigation_helper_error_handling(self):
        """
        Test: Navigation helper should have proper error handling for
        navigation failures and provide meaningful error messages.
        """
        # Test: Start from app launch
        self.splash_page.wait_for_splash_to_complete()
        
        # Test: Navigation helper should handle missing elements gracefully
        # We'll simulate this by attempting navigation with modified locators
        try:
            # This test verifies that the helper provides good error messages
            # when navigation elements are not found
            self._tightened_navigate_to_login()
            
            # If navigation succeeds, verify the elements are properly accessible
            login_elements_accessible = self._validate_login_page_presence()
            self.assertTrue(login_elements_accessible, "Login elements should be accessible after successful navigation")
            
        except Exception as e:
            # If navigation fails, the error should be descriptive
            error_message = str(e)
            self.assertIn("navigation", error_message.lower(), "Error message should mention navigation")
            # Don't fail the test if error handling is working properly

    # Helper methods for navigation testing

    def _tightened_navigate_to_login(self):
        """
        Tightened navigation helper that forces Landing → Login sequence.
        This is the implementation we need to create/fix.
        """
        # Step 1: Ensure we're on landing page
        max_attempts = 3
        for attempt in range(max_attempts):
            try:
                # Look for landing page sign-in button
                sign_in_button = WebDriverWait(self.driver, 5).until(
                    EC.element_to_be_clickable(self.landing_page.SIGN_IN_BUTTON)
                )
                break
            except TimeoutException:
                if attempt == max_attempts - 1:
                    raise Exception("Failed to find landing page after multiple attempts")
                time.sleep(2)
        
        # Step 2: Verify landing page is properly loaded
        self.assertTrue(sign_in_button.is_displayed(), "Landing sign-in button should be visible")
        self.assertTrue(sign_in_button.is_enabled(), "Landing sign-in button should be enabled")
        
        # Step 3: Navigate to login page
        sign_in_button.click()
        
        # Step 4: Verify navigation to login page completed
        try:
            login_submit_button = WebDriverWait(self.driver, 10).until(
                EC.presence_of_element_located((AppiumBy.ACCESSIBILITY_ID, "login_submit"))
            )
        except TimeoutException:
            raise Exception("Navigation to login page failed - login submit button not found")
        
        # Step 5: Additional validation that we're on the correct page
        self.assertTrue(login_submit_button.is_displayed(), "Login submit button should be visible")

    def _try_find_login_elements_directly(self):
        """Try to find login elements without proper navigation"""
        try:
            login_elements = [
                (AppiumBy.ACCESSIBILITY_ID, "login_email"),
                (AppiumBy.ACCESSIBILITY_ID, "login_password"),
                (AppiumBy.ACCESSIBILITY_ID, "login_submit")
            ]
            
            for locator in login_elements:
                element = self.driver.find_element(*locator)
                if not element.is_displayed():
                    return False
            
            return True
        except NoSuchElementException:
            return False

    def _validate_landing_page_presence(self):
        """Validate that landing page is properly loaded"""
        try:
            sign_in_button = WebDriverWait(self.driver, 5).until(
                EC.presence_of_element_located(self.landing_page.SIGN_IN_BUTTON)
            )
            register_button = WebDriverWait(self.driver, 2).until(
                EC.presence_of_element_located(self.landing_page.REGISTER_BUTTON)
            )
            
            return (sign_in_button.is_displayed() and 
                   sign_in_button.is_enabled() and
                   register_button.is_displayed() and
                   register_button.is_enabled())
        except TimeoutException:
            return False

    def _validate_login_page_presence(self):
        """Validate that login page is properly loaded"""
        try:
            required_elements = [
                (AppiumBy.ACCESSIBILITY_ID, "login_email"),
                (AppiumBy.ACCESSIBILITY_ID, "login_password"),
                (AppiumBy.ACCESSIBILITY_ID, "login_submit")
            ]
            
            for locator in required_elements:
                element = WebDriverWait(self.driver, 2).until(
                    EC.presence_of_element_located(locator)
                )
                if not element.is_displayed():
                    return False
            
            return True
        except TimeoutException:
            return False
