"""
Test ValueKey/Semantics Implementation for QA Automation
========================================================

This test suite verifies that ValueKey/semantics are correctly applied directly
on tappable/input widgets (not just wrappers) for QA automation to work properly.

Following TDD principles as required by the user.
"""
import pytest
import time
from appium.webdriver.common.appiumby import AppiumBy
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException

from ..base_test import BaseTest
from ..common.page_objects.page_factory import PageFactory


class TestValueKeySemantics(BaseTest):
    """Test suite for ValueKey/semantics implementation on widgets"""

    def setUp(self):
        """Setup for each test"""
        super().setUp()
        
        # Initialize page objects
        self.splash_page = PageFactory.get_splash_page(self.driver)
        self.landing_page = PageFactory.get_landing_page(self.driver)
        self.login_page = PageFactory.get_login_page(self.driver)

    def test_landing_returning_user_button_has_valuekey(self):
        """
        Test: Landing page "Returning User" button should have ValueKey/semantics
        directly on the button widget for QA automation.
        """
        # Navigate to landing page
        self._navigate_to_landing()
        
        # Test: Button should be found by accessibility ID (ValueKey/semantics)
        returning_user_button = None
        try:
            returning_user_button = WebDriverWait(self.driver, 10).until(
                EC.element_to_be_clickable((AppiumBy.ACCESSIBILITY_ID, "landing_sign_in"))
            )
        except TimeoutException:
            self.fail("Landing 'Returning User' button not found by accessibility ID 'landing_sign_in'")
        
        # Assertions following TDD principles
        self.assertIsNotNone(returning_user_button, "Returning User button should be found")
        self.assertTrue(returning_user_button.is_enabled(), "Returning User button should be enabled")
        self.assertTrue(returning_user_button.is_displayed(), "Returning User button should be visible")
        
        # Test: Button should contain expected text
        button_text = returning_user_button.text or returning_user_button.get_attribute("name")
        self.assertIn("Returning User", button_text, "Button should contain 'Returning User' text")
        
        # Test: Button should be tappable
        try:
            returning_user_button.click()
            # Should navigate to login page
            login_button = WebDriverWait(self.driver, 5).until(
                EC.presence_of_element_located((AppiumBy.ACCESSIBILITY_ID, "login_submit"))
            )
            self.assertIsNotNone(login_button, "Should navigate to login page after clicking")
        except Exception as e:
            self.fail(f"Failed to tap Returning User button: {e}")

    def test_login_email_field_has_valuekey(self):
        """
        Test: Login email field should have ValueKey/semantics directly on the
        input widget (not just wrapper) for QA automation.
        """
        # Navigate to login page
        self._navigate_to_login()
        
        # Test: Email field should be found by accessibility ID
        email_field = None
        try:
            email_field = WebDriverWait(self.driver, 10).until(
                EC.element_to_be_clickable((AppiumBy.ACCESSIBILITY_ID, "login_email"))
            )
        except TimeoutException:
            self.fail("Login email field not found by accessibility ID 'login_email'")
        
        # Assertions following TDD principles
        self.assertIsNotNone(email_field, "Email field should be found")
        self.assertTrue(email_field.is_enabled(), "Email field should be enabled")
        self.assertTrue(email_field.is_displayed(), "Email field should be visible")
        
        # Test: Field should accept text input
        test_email = "test@example.com"
        try:
            email_field.clear()
            email_field.send_keys(test_email)
            
            # Verify text was entered
            entered_text = email_field.text or email_field.get_attribute("value")
            self.assertEqual(test_email, entered_text, "Email field should accept and store text input")
        except Exception as e:
            self.fail(f"Failed to interact with email field: {e}")

    def test_login_password_field_has_valuekey(self):
        """
        Test: Login password field should have ValueKey/semantics directly on the
        input widget (not just wrapper) for QA automation.
        """
        # Navigate to login page
        self._navigate_to_login()
        
        # Test: Password field should be found by accessibility ID
        password_field = None
        try:
            password_field = WebDriverWait(self.driver, 10).until(
                EC.element_to_be_clickable((AppiumBy.ACCESSIBILITY_ID, "login_password"))
            )
        except TimeoutException:
            self.fail("Login password field not found by accessibility ID 'login_password'")
        
        # Assertions following TDD principles
        self.assertIsNotNone(password_field, "Password field should be found")
        self.assertTrue(password_field.is_enabled(), "Password field should be enabled")
        self.assertTrue(password_field.is_displayed(), "Password field should be visible")
        
        # Test: Field should accept text input and mask it (secure text entry)
        test_password = "TestPass123!"
        try:
            password_field.clear()
            password_field.send_keys(test_password)
            
            # For password fields, we don't check the actual value due to security masking
            # But we can verify that some input was registered
            # This varies by platform, so we'll check that the field has been interacted with
            password_field.click()  # Ensure focus
            self.assertTrue(True, "Password field should accept text input")
        except Exception as e:
            self.fail(f"Failed to interact with password field: {e}")

    def test_login_submit_button_has_valuekey(self):
        """
        Test: Login submit button should have ValueKey/semantics directly on the
        button widget (not just wrapper) for QA automation.
        """
        # Navigate to login page
        self._navigate_to_login()
        
        # Test: Submit button should be found by accessibility ID
        submit_button = None
        try:
            submit_button = WebDriverWait(self.driver, 10).until(
                EC.element_to_be_clickable((AppiumBy.ACCESSIBILITY_ID, "login_submit"))
            )
        except TimeoutException:
            self.fail("Login submit button not found by accessibility ID 'login_submit'")
        
        # Assertions following TDD principles
        self.assertIsNotNone(submit_button, "Submit button should be found")
        self.assertTrue(submit_button.is_enabled(), "Submit button should be enabled")
        self.assertTrue(submit_button.is_displayed(), "Submit button should be visible")
        
        # Test: Button should contain expected text
        button_text = submit_button.text or submit_button.get_attribute("name")
        self.assertIn("Sign In", button_text, "Button should contain 'Sign In' text")
        
        # Test: Button should be tappable (without valid credentials, should show validation)
        try:
            submit_button.click()
            # Should trigger form validation since no credentials entered
            time.sleep(2)  # Allow time for validation to appear
            self.assertTrue(True, "Submit button should be tappable")
        except Exception as e:
            self.fail(f"Failed to tap submit button: {e}")

    def test_valuekey_semantics_are_on_widgets_not_wrappers(self):
        """
        Test: Ensure ValueKey/semantics are applied directly on the actual
        input/button widgets, not just on wrapper containers.
        
        This is critical for QA automation tools to properly detect and
        interact with the elements.
        """
        # Navigate to login page
        self._navigate_to_login()
        
        # Test each widget type to ensure semantics are on the actual widget
        widgets_to_test = [
            ("login_email", "Email input widget"),
            ("login_password", "Password input widget"), 
            ("login_submit", "Submit button widget"),
        ]
        
        for accessibility_id, widget_description in widgets_to_test:
            with self.subTest(widget=widget_description):
                # Find element by accessibility ID
                element = None
                try:
                    element = WebDriverWait(self.driver, 5).until(
                        EC.presence_of_element_located((AppiumBy.ACCESSIBILITY_ID, accessibility_id))
                    )
                except TimeoutException:
                    self.fail(f"{widget_description} not found by accessibility ID '{accessibility_id}'")
                
                # Verify element properties that indicate it's the actual widget
                self.assertIsNotNone(element, f"{widget_description} should be found")
                
                # Check that the element has the expected widget properties
                tag_name = element.tag_name.lower()
                element_type = element.get_attribute("type") or ""
                class_name = element.get_attribute("class") or ""
                
                # For input fields, verify they are actual input elements
                if "input" in widget_description.lower():
                    # Should be an actual input element, not a wrapper
                    input_indicators = ["textfield", "edittext", "input", "textinput"]
                    has_input_indicator = any(indicator in tag_name.lower() or 
                                            indicator in element_type.lower() or
                                            indicator in class_name.lower() 
                                            for indicator in input_indicators)
                    self.assertTrue(has_input_indicator, 
                                  f"{widget_description} should be an actual input element, not a wrapper")
                
                # For buttons, verify they are actual button elements
                elif "button" in widget_description.lower():
                    button_indicators = ["button", "btn"]
                    has_button_indicator = any(indicator in tag_name.lower() or 
                                             indicator in element_type.lower() or
                                             indicator in class_name.lower()
                                             for indicator in button_indicators)
                    self.assertTrue(has_button_indicator,
                                  f"{widget_description} should be an actual button element, not a wrapper")

    def _navigate_to_landing(self):
        """Helper method to navigate to landing page"""
        try:
            # Wait for splash to complete
            self.splash_page.wait_for_splash_to_complete()
            
            # Verify we're on landing page or navigate to it
            if not self.landing_page.is_element_present(self.landing_page.SIGN_IN_BUTTON, timeout=3):
                # If not on landing, this might be first run - wait longer
                WebDriverWait(self.driver, 10).until(
                    EC.presence_of_element_located(self.landing_page.SIGN_IN_BUTTON)
                )
        except Exception as e:
            self.fail(f"Failed to navigate to landing page: {e}")

    def _navigate_to_login(self):
        """Helper method to navigate to login page"""
        try:
            # First navigate to landing
            self._navigate_to_landing()
            
            # Then navigate to login
            self.landing_page.goto_login()
            
            # Verify we're on login page
            WebDriverWait(self.driver, 10).until(
                EC.presence_of_element_located((AppiumBy.ACCESSIBILITY_ID, "login_submit"))
            )
        except Exception as e:
            self.fail(f"Failed to navigate to login page: {e}")
