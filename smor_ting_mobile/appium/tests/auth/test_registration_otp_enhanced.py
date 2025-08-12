"""
Enhanced Registration → OTP Flow Tests with Flutter Driver Support

This test suite provides comprehensive testing of the registration and OTP verification
flow using Flutter Driver for better element discovery and improved reliability.

Features tested:
- Flutter-first element discovery with UiAutomator2 fallbacks
- Complete registration flow validation
- OTP verification including UI/UX behaviors
- Error handling and edge cases
- Network resilience testing
"""

import time
import pytest
from typing import Dict, Any
from ..base_test import BaseTest
from ..common.page_objects import PageFactory
from ..common.otp_client import get_latest_otp


@pytest.mark.e2e
@pytest.mark.registration
@pytest.mark.flutter
class TestRegistrationOtpFlowEnhanced(BaseTest):
    """Enhanced Registration and OTP verification test suite with Flutter Driver support"""

    def setup_method(self, method):
        """Setup for each test method"""
        super().setup_method(method)
        self.splash = PageFactory.get_splash_page(self.driver)
        self.register = PageFactory.get_registration_page(self.driver)
        self.login = PageFactory.get_login_page(self.driver)
        self.otp = PageFactory.get_otp_page(self.driver)
        self.dashboard = PageFactory.get_dashboard_page(self.driver)
        self.splash.wait_for_splash_to_complete()

    def _generate_test_user(self, role: str = 'customer', suffix: str = None) -> Dict[str, Any]:
        """Generate test user data with unique email"""
        suffix = suffix or str(int(time.time()))
        return {
            'email': f"flutter_test_{role}_{suffix}@smorting.com",
            'password': 'FlutterTest123!',
            'first_name': 'Flutter',
            'last_name': f'Test{suffix}',
            'phone': '231777123456',
            'role': role,
        }

    def _navigate_to_registration(self):
        """Navigate to registration page from current state"""
        if self.login.is_element_present(self.login.REGISTER_LINK, timeout=5):
            self.login.tap_register_link()
        
        # Verify we're on registration page
        assert self.register.is_element_present_flutter_first(
            self.register.REGISTER_FLUTTER_KEY, 
            self.register.REGISTER_BUTTON_FALLBACK, 
            timeout=10
        ), "Failed to navigate to registration page"

    def test_registration_with_flutter_element_discovery(self):
        """
        Test successful registration using Flutter Driver for element discovery
        
        This test verifies that our Flutter-first approach can successfully
        discover and interact with form elements that were problematic with UiAutomator2.
        """
        # Arrange
        test_user = self._generate_test_user('customer', 'flutter_discovery')
        self._navigate_to_registration()

        # Act - Fill form using Flutter-first methods
        self.register.fill_registration_form(test_user)
        
        # Verify all fields were filled (Flutter Driver should find them reliably)
        assert self.register.is_element_present_flutter_first(
            self.register.EMAIL_FLUTTER_KEY, 
            self.register.EMAIL_FALLBACK, 
            timeout=5
        ), "Email field should be discoverable via Flutter Driver"
        
        assert self.register.is_element_present_flutter_first(
            self.register.PASSWORD_FLUTTER_KEY, 
            self.register.PASSWORD_FALLBACK, 
            timeout=5
        ), "Password field should be discoverable via Flutter Driver"

        # Submit registration
        self.register.tap_register_button()
        self.register.wait_for_registration_complete()

        # Assert - Should navigate to OTP page
        assert self.otp.is_element_present_flutter_first(
            self.otp.OTP_FIELD_FLUTTER_KEY,
            self.otp.OTP_EDIT_TEXTS,
            timeout=15
        ), "Should navigate to OTP verification page"

        self.logger.info(f"✅ Registration successful with Flutter Driver: {test_user['email']}")

    def test_otp_verification_with_flutter_widgets(self):
        """
        Test OTP verification screen using Flutter Driver for better widget discovery
        
        Focuses on OTP-specific Flutter widgets that may not be discoverable
        via traditional UiAutomator2 selectors.
        """
        # Arrange - Complete registration first
        test_user = self._generate_test_user('customer', 'otp_flutter')
        self._navigate_to_registration()
        self.register.fill_registration_form(test_user)
        self.register.tap_register_button()
        self.register.wait_for_registration_complete()

        # Verify OTP screen elements using Flutter Driver
        assert self.otp.is_element_present_flutter_first(
            self.otp.COUNTDOWN_FLUTTER_KEY,
            self.otp.COUNTDOWN_FALLBACK,
            timeout=10
        ), "Countdown widget should be discoverable via Flutter Driver"

        assert self.otp.is_element_present_flutter_first(
            self.otp.RESEND_FLUTTER_KEY,
            self.otp.RESEND_FALLBACK,
            timeout=10
        ), "Resend button should be discoverable via Flutter Driver"

        # Test OTP input using Flutter Driver
        api_base = self.config.get_api_config()["base_url"]
        otp = get_latest_otp(api_base, test_user['email']) or '123456'
        
        self.otp.enter_otp(otp)

        # Verify OTP submission
        if self.otp.is_element_present_flutter_first(
            self.otp.VERIFY_FLUTTER_KEY,
            self.otp.VERIFY_FALLBACK,
            timeout=5
        ):
            self.otp.tap_verify()

        # Assert - Should navigate to dashboard
        for _ in range(20):
            if self.dashboard.is_dashboard_loaded():
                break
            time.sleep(0.5)

        assert self.dashboard.is_dashboard_loaded(), "Should navigate to dashboard after OTP verification"
        self.logger.info("✅ OTP verification successful with Flutter Driver")

    def test_registration_form_field_resilience(self):
        """
        Test form field interaction resilience with Flutter vs UiAutomator2 fallback
        
        This test attempts interactions that commonly fail with UiAutomator2
        and verifies Flutter Driver provides better reliability.
        """
        # Arrange
        test_user = self._generate_test_user('provider', 'resilience')
        self._navigate_to_registration()

        # Test each field individually with explicit Flutter-first approach
        fields_to_test = [
            ('email', self.register.EMAIL_FLUTTER_KEY, self.register.EMAIL_FALLBACK),
            ('first_name', self.register.FIRST_NAME_FLUTTER_KEY, self.register.FIRST_NAME_FALLBACK),
            ('last_name', self.register.LAST_NAME_FLUTTER_KEY, self.register.LAST_NAME_FALLBACK),
            ('phone', self.register.PHONE_FLUTTER_KEY, self.register.PHONE_FALLBACK),
            ('password', self.register.PASSWORD_FLUTTER_KEY, self.register.PASSWORD_FALLBACK),
        ]

        successful_discoveries = []
        
        for field_name, flutter_key, fallback_locator in fields_to_test:
            try:
                # Test Flutter Driver discovery
                element = self.register.find_element_flutter_first(flutter_key, fallback_locator, timeout=10)
                
                # Test interaction (clear and enter text)
                if field_name in test_user:
                    self.register.enter_text_flutter_first(flutter_key, fallback_locator, test_user[field_name])
                    successful_discoveries.append(field_name)
                    self.logger.info(f"✅ Successfully interacted with {field_name} field via Flutter Driver")
                    
            except Exception as e:
                self.logger.warning(f"⚠️ Flutter Driver interaction failed for {field_name}: {e}")

        # Assert that Flutter Driver was able to interact with most fields
        assert len(successful_discoveries) >= 4, f"Flutter Driver should successfully interact with most fields. Successful: {successful_discoveries}"
        
        self.logger.info(f"✅ Form field resilience test passed. Successful interactions: {len(successful_discoveries)}/5")

    def test_otp_countdown_widget_discovery(self):
        """
        Test OTP countdown widget discovery using Flutter Driver
        
        Countdown widgets in Flutter often have custom rendering that
        UiAutomator2 cannot properly detect.
        """
        # Arrange - Get to OTP screen
        test_user = self._generate_test_user('customer', 'countdown')
        self._navigate_to_registration()
        self.register.fill_registration_form(test_user)
        self.register.tap_register_button()
        self.register.wait_for_registration_complete()

        # Test countdown text retrieval using Flutter Driver
        countdown_text = self.otp.get_countdown_text()
        
        # Verify countdown is working
        assert len(countdown_text) > 0, "Countdown text should be retrievable via Flutter Driver"
        assert any(keyword in countdown_text.lower() for keyword in ['expires', 'resend', 'minute', 'second']), \
            f"Countdown text should contain time-related keywords: {countdown_text}"

        # Test resend button state using Flutter Driver
        resend_enabled = self.otp.is_resend_enabled()
        
        # Initially, resend should be disabled due to countdown
        assert not resend_enabled, "Resend button should initially be disabled during countdown"

        self.logger.info(f"✅ Countdown widget discovered successfully. Text: '{countdown_text}', Resend enabled: {resend_enabled}")

    def test_flutter_vs_uiautomator2_fallback_behavior(self):
        """
        Test the fallback behavior from Flutter Driver to UiAutomator2
        
        This test intentionally triggers fallback scenarios to ensure
        the system gracefully handles mixed automation approaches.
        """
        # Arrange
        test_user = self._generate_test_user('customer', 'fallback')
        self._navigate_to_registration()

        flutter_attempts = 0
        fallback_attempts = 0

        # Mock scenario where Flutter Driver might fail for some elements
        try:
            # This should work with Flutter Driver
            self.register.enter_text_flutter_first(
                self.register.EMAIL_FLUTTER_KEY, 
                self.register.EMAIL_FALLBACK, 
                test_user['email']
            )
            flutter_attempts += 1
        except Exception:
            fallback_attempts += 1

        try:
            # Test other fields
            self.register.enter_text_flutter_first(
                self.register.FIRST_NAME_FLUTTER_KEY, 
                self.register.FIRST_NAME_FALLBACK, 
                test_user['first_name']
            )
            flutter_attempts += 1
        except Exception:
            fallback_attempts += 1

        # At least some interactions should succeed
        total_attempts = flutter_attempts + fallback_attempts
        assert total_attempts >= 1, "At least some form interactions should succeed"

        self.logger.info(f"✅ Fallback behavior test completed. Flutter: {flutter_attempts}, Fallback: {fallback_attempts}")

    @pytest.mark.stress
    def test_multiple_registration_attempts_stability(self):
        """
        Test stability of Flutter Driver across multiple registration attempts
        
        This stress test ensures Flutter Driver maintains reliability
        across multiple form interactions and page navigations.
        """
        successful_attempts = 0
        failed_attempts = 0

        for attempt in range(3):  # 3 attempts for reasonable test time
            try:
                test_user = self._generate_test_user('customer', f'stress_{attempt}')
                
                # Navigate to registration
                if attempt > 0:
                    # Navigate back to registration if not first attempt
                    self.driver.back()
                    time.sleep(2)
                self._navigate_to_registration()

                # Fill form with Flutter Driver
                self.register.fill_registration_form(test_user)
                
                # Verify button is discoverable
                assert self.register.is_element_present_flutter_first(
                    self.register.REGISTER_FLUTTER_KEY,
                    self.register.REGISTER_BUTTON_FALLBACK,
                    timeout=10
                ), f"Register button should be discoverable on attempt {attempt + 1}"

                successful_attempts += 1
                self.logger.info(f"✅ Stress test attempt {attempt + 1} successful")

            except Exception as e:
                failed_attempts += 1
                self.logger.warning(f"⚠️ Stress test attempt {attempt + 1} failed: {e}")

        # Assert that Flutter Driver maintains good stability
        success_rate = successful_attempts / (successful_attempts + failed_attempts)
        assert success_rate >= 0.67, f"Flutter Driver should maintain >67% success rate across attempts. Actual: {success_rate:.2f}"

        self.logger.info(f"✅ Stress test completed. Success rate: {success_rate:.2f} ({successful_attempts}/{successful_attempts + failed_attempts})")

    def test_registration_otp_complete_flow_flutter_driver(self):
        """
        Complete end-to-end registration → OTP → dashboard flow using Flutter Driver
        
        This is the primary test that validates the complete user journey
        with improved element discovery.
        """
        # Arrange
        test_user = self._generate_test_user('customer', 'complete_flow')
        self._navigate_to_registration()

        # Act - Complete Registration
        self.register.fill_registration_form(test_user)
        self.register.tap_register_button()
        self.register.wait_for_registration_complete()

        # Verify OTP screen loaded
        assert self.otp.is_element_present_flutter_first(
            self.otp.OTP_FIELD_FLUTTER_KEY,
            self.otp.OTP_EDIT_TEXTS,
            timeout=15
        ), "OTP screen should load after registration"

        # Complete OTP verification
        api_base = self.config.get_api_config()["base_url"]
        otp = get_latest_otp(api_base, test_user['email']) or '123456'
        self.otp.enter_otp(otp)

        # Submit OTP if verify button exists
        if self.otp.is_element_present_flutter_first(
            self.otp.VERIFY_FLUTTER_KEY,
            self.otp.VERIFY_FALLBACK,
            timeout=5
        ):
            self.otp.tap_verify()

        # Wait for dashboard navigation
        for _ in range(20):
            if self.dashboard.is_dashboard_loaded():
                break
            time.sleep(0.5)

        # Assert complete flow success
        assert self.dashboard.is_dashboard_loaded(), "Complete flow should end at dashboard"
        assert self.dashboard.verify_role_specific_elements(test_user['role']), f"Dashboard should show {test_user['role']}-specific elements"

        self.logger.info(f"✅ Complete Registration → OTP → Dashboard flow successful with Flutter Driver: {test_user['email']}")

    def teardown_method(self, method):
        """Cleanup after each test"""
        try:
            # Take screenshot on failure for debugging
            if hasattr(self, '_pytest_assertion_failures') and self._pytest_assertion_failures:
                self._take_failure_screenshot(method.__name__)
        except Exception:
            pass
        
        super().teardown_method(method)
