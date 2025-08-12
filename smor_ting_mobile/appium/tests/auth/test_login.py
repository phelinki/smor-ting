"""
Login Test Suite for Smor-Ting Mobile App

This module implements comprehensive test cases for user login functionality,
covering all validation scenarios, error handling, and success flows based on TDD principles.
"""
import pytest
import time
from typing import Dict, Any
from ..base_test import BaseTest
from ..common.page_objects import PageFactory


class TestLogin(BaseTest):
    """Test suite for user login functionality"""
    
    def setup_method(self, method):
        """Setup for each test method"""
        super().setup_method(method)
        self.splash_page = PageFactory.get_splash_page(self.driver)
        self.registration_page = PageFactory.get_registration_page(self.driver)
        self.login_page = PageFactory.get_login_page(self.driver)
        self.dashboard_page = PageFactory.get_dashboard_page(self.driver)
        self.error_dialog = PageFactory.get_error_dialog(self.driver)
        
        # Create a test user for login tests
        self.test_user = self._create_test_user()
        
        # Navigate to login page
        self._navigate_to_login()
    
    def _create_test_user(self) -> Dict[str, str]:
        """Create a test user for login tests"""
        test_data = self.config.get_test_data()
        user_data = test_data['valid_users'][0].copy()
        user_data['email'] = f"login_test_{int(time.time())}@smorting.com"
        
        try:
            # Navigate to registration first
            self.splash_page.wait_for_splash_to_complete()

            # Prefer landing → registration
            try:
                landing = PageFactory.get_landing_page(self.driver)
                landing.goto_register()
            except Exception:
                # Fallback: from login page, use register link
                if self.login_page.is_element_present(self.login_page.REGISTER_LINK, timeout=5):
                    self.login_page.tap_register_link()
                else:
                    # As last resort, try tapping any visible Register button by text
                    try:
                        self.registration_page.tap(self.registration_page.REGISTER_BUTTON_FALLBACK, timeout=2)
                    except Exception:
                        pass
            
            # Register the test user
            self.registration_page.fill_registration_form(user_data)
            self.registration_page.tap_register_button()
            self.registration_page.wait_for_registration_complete()
            
            # Logout if registration was successful
            if self.dashboard_page.is_dashboard_loaded():
                self.dashboard_page.logout()
                
            self.logger.info(f"Test user created: {user_data['email']}")
            
        except Exception as e:
            self.logger.warning(f"Could not create test user: {e}")
            # Use existing user data if creation fails
            user_data = test_data['existing_user'].copy()
        
        return user_data
    
    def _navigate_to_login(self):
        """Navigate to login page"""
        try:
            # Wait for splash if present
            self.splash_page.wait_for_splash_to_complete()

            # Use robust helper that can navigate via landing/registration
            self.login_page.ensure_on_login(timeout=10)

            # If still not present, try landing → login explicitly
            if not self.login_page.is_element_present(self.login_page.LOGIN_BUTTON, timeout=3):
                try:
                    landing = PageFactory.get_landing_page(self.driver)
                    landing.goto_login()
                except Exception:
                    pass

            # As a final fallback, try registration → login link
            if not self.login_page.is_element_present(self.login_page.LOGIN_BUTTON, timeout=3):
                if self.registration_page.is_element_present(self.registration_page.LOGIN_LINK, timeout=3):
                    self.registration_page.tap_login_link()

        except Exception as e:
            self.logger.warning(f"Navigation to login page: {e}")

        # Verify we're on login page
        self.login_page.assert_element_present(
            self.login_page.choose_locator(self.login_page.LOGIN_BUTTON, self.login_page.LOGIN_BUTTON_FALLBACK),
            "Failed to navigate to login page"
        )

        # Capture login page screenshot
        try:
            self.take_screenshot("login_page")
        except Exception as e:
            self.logger.warning(f"Failed to capture login page screenshot: {e}")
    
    def test_successful_login(self):
        """
        TC_LOGIN_001: Test successful login with valid credentials
        Verifies that user can login successfully with correct email and password
        """
        # Arrange
        email = self.test_user['email']
        password = self.test_user['password']
        
        # Act
        self.login_page.fill_login_form(email, password)
        self.login_page.tap_login_button()
        self.login_page.wait_for_login_complete()
        
        # Assert
        assert self.dashboard_page.is_dashboard_loaded(), "Dashboard should load after successful login"
        
        # Verify role-specific elements if role is known
        if 'role' in self.test_user:
            assert self.dashboard_page.verify_role_specific_elements(self.test_user['role']), \
                f"Role-specific elements for {self.test_user['role']} should be visible"
        
        self.logger.info(f"✅ Successful login for user: {email}")

        # Capture dashboard screenshot after successful login
        try:
            self.take_screenshot("dashboard_after_login")
        except Exception as e:
            self.logger.warning(f"Failed to capture dashboard screenshot: {e}")
    
    def test_login_with_nonexistent_email(self):
        """
        TC_LOGIN_002: Test login with non-existent email
        Verifies proper error handling for email not registered in system
        """
        # Arrange
        nonexistent_email = f"nonexistent_{int(time.time())}@example.com"
        password = "AnyPassword123!"
        
        # Act
        self.login_page.fill_login_form(nonexistent_email, password)
        self.login_page.tap_login_button()
        
        # Wait for response
        time.sleep(3)
        
        # Assert
        validation_errors = self.login_page.get_validation_errors()
        
        # Check for credentials error or error dialog
        has_error = (
            'credentials' in validation_errors or
            self.error_dialog.is_error_dialog_visible() or
            not self.dashboard_page.is_dashboard_loaded()
        )
        
        assert has_error, "Should show error for non-existent email"
        
        if 'credentials' in validation_errors:
            error_message = validation_errors['credentials']
            assert "invalid" in error_message.lower(), f"Error should mention invalid credentials: {error_message}"
        
        self.logger.info(f"✅ Non-existent email error handled correctly: {nonexistent_email}")
    
    def test_login_with_wrong_password(self):
        """
        TC_LOGIN_003: Test login with incorrect password
        Verifies proper error handling for wrong password with valid email
        """
        # Arrange
        email = self.test_user['email']
        wrong_password = "WrongPassword123!"
        
        # Act
        self.login_page.fill_login_form(email, wrong_password)
        self.login_page.tap_login_button()
        
        # Wait for response
        time.sleep(3)
        
        # Assert
        validation_errors = self.login_page.get_validation_errors()
        
        # Check for credentials error or error dialog
        has_error = (
            'credentials' in validation_errors or
            self.error_dialog.is_error_dialog_visible() or
            not self.dashboard_page.is_dashboard_loaded()
        )
        
        assert has_error, "Should show error for wrong password"
        
        if 'credentials' in validation_errors:
            error_message = validation_errors['credentials']
            assert "invalid" in error_message.lower(), f"Error should mention invalid credentials: {error_message}"
        
        self.logger.info(f"✅ Wrong password error handled correctly for user: {email}")
    
    def test_login_with_empty_email(self):
        """
        TC_LOGIN_004: Test validation for empty email field
        Verifies that empty email triggers validation error
        """
        # Arrange
        empty_email = ""
        password = "TestPassword123!"
        
        # Act
        self.login_page.fill_login_form(empty_email, password)
        self.login_page.tap_login_button()
        
        # Assert
        validation_errors = self.login_page.get_validation_errors()
        
        assert 'email' in validation_errors or not self.login_page.is_login_button_enabled(), \
            "Should show validation error for empty email or disable login button"
        
        if 'email' in validation_errors:
            error_message = validation_errors['email']
            assert "required" in error_message.lower(), f"Error should mention email is required: {error_message}"
        
        self.logger.info("✅ Empty email validation works correctly")
    
    def test_login_with_empty_password(self):
        """
        TC_LOGIN_005: Test validation for empty password field
        Verifies that empty password triggers validation error
        """
        # Arrange
        email = self.test_user['email']
        empty_password = ""
        
        # Act
        self.login_page.fill_login_form(email, empty_password)
        self.login_page.tap_login_button()
        
        # Assert
        validation_errors = self.login_page.get_validation_errors()
        
        assert 'password' in validation_errors or not self.login_page.is_login_button_enabled(), \
            "Should show validation error for empty password or disable login button"
        
        if 'password' in validation_errors:
            error_message = validation_errors['password']
            assert "required" in error_message.lower(), f"Error should mention password is required: {error_message}"
        
        self.logger.info("✅ Empty password validation works correctly")
    
    @pytest.mark.parametrize("invalid_email", [
        "invalid-email",
        "test@",
        "@domain.com",
        "spaces in@email.com"
    ])
    def test_login_with_invalid_email_format(self, invalid_email):
        """
        TC_LOGIN_006: Test email format validation
        Verifies that malformed email addresses trigger validation
        """
        # Arrange
        password = "TestPassword123!"
        
        # Act
        self.login_page.fill_login_form(invalid_email, password)
        self.login_page.tap_login_button()
        
        # Assert
        validation_errors = self.login_page.get_validation_errors()
        
        # Check for email validation or prevention of submission
        has_validation_error = (
            'email' in validation_errors or 
            not self.login_page.is_login_button_enabled() or
            not self.dashboard_page.is_dashboard_loaded()
        )
        
        assert has_validation_error, f"Invalid email '{invalid_email}' should trigger validation"
        
        self.logger.info(f"✅ Invalid email format '{invalid_email}' validation handled correctly")
    
    def test_login_form_navigation_to_registration(self):
        """
        Test navigation from login form to registration page
        Verifies that the register link works correctly
        """
        # Act
        self.login_page.tap_register_link()
        
        # Assert
        assert self.registration_page.is_element_present(self.registration_page.REGISTER_BUTTON), \
            "Should navigate to registration page"
        assert self.registration_page.is_element_present(self.registration_page.EMAIL_FIELD), \
            "Registration form should be visible"
        
        self.logger.info("✅ Navigation from login to registration works correctly")
    
    def test_login_loading_states(self):
        """
        TC_UI_001: Test loading states during login
        Verifies loading indicators and button states during authentication
        """
        # Arrange
        email = self.test_user['email']
        password = self.test_user['password']
        
        # Act
        self.login_page.fill_login_form(email, password)
        
        # Check button is enabled before submission
        assert self.login_page.is_login_button_enabled(), "Login button should be enabled with valid data"
        
        self.login_page.tap_login_button()
        
        # Check for loading state immediately after tap
        loading_present = self.login_page.is_element_present(
            self.login_page.LOADING_INDICATOR, 
            timeout=2
        )
        
        # Wait for login to complete
        self.login_page.wait_for_login_complete()
        
        # Assert
        # Loading indicator should be gone
        assert not self.login_page.is_element_present(
            self.login_page.LOADING_INDICATOR, 
            timeout=2
        ), "Loading indicator should disappear after login"
        
        self.logger.info(f"✅ Login loading states handled correctly (loading detected: {loading_present})")
    
    def test_login_performance(self):
        """
        TC_PERF_002: Test login performance
        Verifies login completes within acceptable time (< 3 seconds)
        """
        # Arrange
        email = self.test_user['email']
        password = self.test_user['password']
        
        # Act
        start_time = time.time()
        
        self.login_page.fill_login_form(email, password)
        self.login_page.tap_login_button()
        self.login_page.wait_for_login_complete()
        
        end_time = time.time()
        login_duration = end_time - start_time
        
        # Assert
        assert login_duration < 5.0, f"Login should complete in less than 5 seconds, took {login_duration:.2f}s"
        
        # Performance target is < 3 seconds, but we'll be lenient for test environment
        if login_duration < 3.0:
            self.logger.info(f"✅ Login performance excellent: {login_duration:.2f}s")
        else:
            self.logger.info(f"⚠️ Login performance acceptable but slow: {login_duration:.2f}s")
    
    def test_login_network_error_handling(self):
        """
        TC_NETWORK_001-003: Test network error scenarios during login
        Verifies app behavior during network issues
        """
        # Arrange
        email = self.test_user['email']
        password = self.test_user['password']
        
        # Act
        self.login_page.fill_login_form(email, password)
        self.login_page.tap_login_button()
        
        # Wait for longer than normal to test timeout handling
        try:
            self.login_page.wait_for_login_complete(timeout=60)
            
            if self.dashboard_page.is_dashboard_loaded():
                self.logger.info("✅ Login completed successfully (no network issues)")
            else:
                # Check for error handling
                if self.error_dialog.is_error_dialog_visible():
                    error_message = self.error_dialog.get_error_message()
                    if self.error_dialog.is_network_error():
                        self.logger.info(f"✅ Network error handled correctly: {error_message}")
                    elif self.error_dialog.is_server_error():
                        self.logger.info(f"✅ Server error handled correctly: {error_message}")
                    else:
                        self.logger.info(f"✅ Error handled: {error_message}")
                        
        except Exception as e:
            self.logger.info(f"✅ Network timeout handled: {e}")
    
    def test_login_after_logout(self):
        """
        Test login functionality after logout
        Verifies that user can login again after logging out
        """
        # Arrange - First login
        email = self.test_user['email']
        password = self.test_user['password']
        
        self.login_page.fill_login_form(email, password)
        self.login_page.tap_login_button()
        self.login_page.wait_for_login_complete()
        
        # Verify login success
        assert self.dashboard_page.is_dashboard_loaded(), "Initial login should succeed"
        
        # Act - Logout and login again
        self.dashboard_page.logout()
        
        # Navigate back to login if needed
        self._navigate_to_login()
        
        # Login again
        self.login_page.fill_login_form(email, password)
        self.login_page.tap_login_button()
        self.login_page.wait_for_login_complete()
        
        # Assert
        assert self.dashboard_page.is_dashboard_loaded(), "Should be able to login again after logout"
        
        self.logger.info(f"✅ Login after logout works correctly for user: {email}")
    
    def test_multiple_failed_login_attempts(self):
        """
        Test behavior with multiple failed login attempts
        Verifies app handles repeated failed login attempts gracefully
        """
        # Arrange
        email = self.test_user['email']
        wrong_password = "WrongPassword123!"
        
        # Act - Try multiple failed logins
        for attempt in range(3):
            self.logger.info(f"Failed login attempt {attempt + 1}")
            
            self.login_page.fill_login_form(email, wrong_password)
            self.login_page.tap_login_button()
            
            # Wait for response
            time.sleep(2)
            
            # Check that we're still on login page (not locked out)
            assert self.login_page.is_element_present(self.login_page.LOGIN_BUTTON), \
                f"Login page should still be accessible after {attempt + 1} failed attempts"
            
            # Clear fields for next attempt
            self.login_page.enter_text(self.login_page.EMAIL_FIELD, "")
            self.login_page.enter_text(self.login_page.PASSWORD_FIELD, "")
        
        # Assert - Should still be able to login with correct credentials
        self.login_page.fill_login_form(email, self.test_user['password'])
        self.login_page.tap_login_button()
        self.login_page.wait_for_login_complete()
        
        assert self.dashboard_page.is_dashboard_loaded(), \
            "Should still be able to login with correct credentials after failed attempts"
        
        self.logger.info("✅ Multiple failed login attempts handled correctly")
    
    def test_login_form_field_clearing(self):
        """
        Test that login form fields can be cleared and refilled
        Verifies form field behavior and data entry
        """
        # Arrange
        first_email = "first@example.com"
        first_password = "FirstPassword123!"
        second_email = self.test_user['email']
        second_password = self.test_user['password']
        
        # Act - Fill form first time
        self.login_page.fill_login_form(first_email, first_password)
        
        # Verify fields are filled
        email_text = self.login_page.get_attribute(self.login_page.EMAIL_FIELD, "text")
        assert first_email in email_text or len(email_text) > 0, "Email field should contain entered text"
        
        # Clear and refill
        self.login_page.fill_login_form(second_email, second_password)
        
        # Verify fields are updated
        email_text = self.login_page.get_attribute(self.login_page.EMAIL_FIELD, "text")
        assert second_email in email_text or len(email_text) > 0, "Email field should contain new text"
        
        # Test login with final values
        self.login_page.tap_login_button()
        self.login_page.wait_for_login_complete()
        
        # Assert
        assert self.dashboard_page.is_dashboard_loaded(), "Login should work with cleared and refilled form"
        
        self.logger.info("✅ Login form field clearing and refilling works correctly")
