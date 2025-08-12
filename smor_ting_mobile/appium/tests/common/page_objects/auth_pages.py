"""
Authentication Page Objects for Smor-Ting Mobile App

This module contains page objects for all authentication-related screens
including registration, login, and error handling flows.
"""
from typing import Dict, Any
from appium.webdriver.common.appiumby import AppiumBy
from .base_page import BasePage


class SplashPage(BasePage):
    """Page object for the splash screen"""
    
    # Locators
    LOGO = (AppiumBy.ACCESSIBILITY_ID, "smor_ting_logo")
    LOADING_INDICATOR = (AppiumBy.CLASS_NAME, "CircularProgressIndicator")
    
    def wait_for_splash_to_complete(self, timeout: int = 10):
        """Wait for splash screen to complete"""
        try:
            # Wait for splash screen elements
            self.wait_for_element_visible(self.LOGO, timeout=5)
            # Wait for splash to disappear
            self.wait_for_element_to_disappear(self.LOGO, timeout)
        except Exception:
            # Splash might be very quick or already gone
            pass


class LandingPage(BasePage):
    """Landing screen shown after splash, with Sign In / Register actions"""

    SIGN_IN_BUTTON = (
        AppiumBy.ACCESSIBILITY_ID,
        "landing_sign_in",
    )
    SIGN_IN_FALLBACK = (
        AppiumBy.XPATH,
        "//android.widget.Button[contains(@text, 'Sign In') or contains(@text, 'Login') or contains(@text, 'Returning User')]",
    )

    REGISTER_BUTTON = (
        AppiumBy.ACCESSIBILITY_ID,
        "landing_register",
    )
    REGISTER_FALLBACK = (
        AppiumBy.XPATH,
        "//android.widget.Button[contains(@text, 'Register') or contains(@text, 'Sign Up') or contains(@text, 'New User')]",
    )

    def goto_login(self):
        locator = self.choose_locator(self.SIGN_IN_BUTTON, self.SIGN_IN_FALLBACK)
        self.tap(locator)
        return PageFactory.get_login_page(self.driver)

    def goto_register(self):
        locator = self.choose_locator(self.REGISTER_BUTTON, self.REGISTER_FALLBACK)
        self.tap(locator)
        return PageFactory.get_registration_page(self.driver)


class RegistrationPage(BasePage):
    """Page object for the registration screen with Flutter Driver support"""
    
    # Flutter keys and fallback locators for better element discovery
    EMAIL_FLUTTER_KEY = "register_email"
    EMAIL_FIELD = (AppiumBy.ACCESSIBILITY_ID, "register_email")
    EMAIL_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'email') or contains(@hint, 'Email')]")
    
    PASSWORD_FLUTTER_KEY = "register_password"
    PASSWORD_FIELD = (AppiumBy.ACCESSIBILITY_ID, "register_password")
    PASSWORD_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'password') or contains(@hint, 'Password')]")
    
    CONFIRM_PASSWORD_FLUTTER_KEY = "register_confirm_password"
    CONFIRM_PASSWORD_FIELD = (AppiumBy.ACCESSIBILITY_ID, "register_confirm_password")
    CONFIRM_PASSWORD_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'confirm') or contains(@hint, 'Confirm')]")
    
    FIRST_NAME_FLUTTER_KEY = "register_first_name"
    FIRST_NAME_FIELD = (AppiumBy.ACCESSIBILITY_ID, "register_first_name")
    FIRST_NAME_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'first') or contains(@hint, 'First')]")
    
    LAST_NAME_FLUTTER_KEY = "register_last_name"
    LAST_NAME_FIELD = (AppiumBy.ACCESSIBILITY_ID, "register_last_name")
    LAST_NAME_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'last') or contains(@hint, 'Last')]")
    
    PHONE_FLUTTER_KEY = "register_phone"
    PHONE_FIELD = (AppiumBy.ACCESSIBILITY_ID, "register_phone")
    PHONE_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'phone') or contains(@hint, 'Phone')]")
    
    # Role selection locators
    ROLE_DROPDOWN = (AppiumBy.XPATH, "//android.widget.Spinner[contains(@content-desc, 'role') or contains(@hint, 'Role')]")
    CUSTOMER_ROLE = (AppiumBy.XPATH, "//*[contains(@text, 'Customer')]")
    PROVIDER_ROLE = (AppiumBy.XPATH, "//*[contains(@text, 'Provider')]")
    ADMIN_ROLE = (AppiumBy.XPATH, "//*[contains(@text, 'Admin')]")
    
    # Button locators with Flutter keys
    REGISTER_FLUTTER_KEY = "register_submit"
    REGISTER_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "register_submit")
    REGISTER_BUTTON_FALLBACK = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Register') or contains(@content-desc, 'Register')]")
    REGISTER_AS_CUSTOMER_BUTTON = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Register as Customer')]")
    REGISTER_AS_AGENT_BUTTON = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Register as Agent')]")
    LOGIN_LINK = (AppiumBy.ACCESSIBILITY_ID, "register_to_login")
    LOGIN_LINK_FALLBACK = (
        AppiumBy.XPATH,
        "//*[contains(@text, 'Login') or contains(@text, 'Sign In')]",
    )
    
    # Validation error locators
    EMAIL_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'email') and (contains(@text, 'required') or contains(@text, 'invalid'))]")
    PASSWORD_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'password') and (contains(@text, 'required') or contains(@text, 'short'))]")
    FIRST_NAME_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'first name') and contains(@text, 'required')]")
    LAST_NAME_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'last name') and contains(@text, 'required')]")
    PHONE_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'phone') and contains(@text, 'required')]")
    ROLE_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'role') and contains(@text, 'required')]")
    PASSWORD_MISMATCH_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'password') and contains(@text, 'match')]")
    
    # Loading states
    LOADING_INDICATOR = (AppiumBy.CLASS_NAME, "CircularProgressIndicator")
    
    def fill_registration_form(self, user_data: Dict[str, Any]):
        """Fill the registration form with user data using Flutter-first approach"""
        if 'email' in user_data:
            self.enter_text_flutter_first(self.EMAIL_FLUTTER_KEY, self.EMAIL_FALLBACK, user_data['email'])
        
        if 'password' in user_data:
            self.enter_text_flutter_first(self.PASSWORD_FLUTTER_KEY, self.PASSWORD_FALLBACK, user_data['password'])
            
        if 'confirm_password' in user_data:
            self.enter_text_flutter_first(self.CONFIRM_PASSWORD_FLUTTER_KEY, self.CONFIRM_PASSWORD_FALLBACK, user_data['confirm_password'])
        elif 'password' in user_data:
            # Use same password for confirmation if not specified
            self.enter_text_flutter_first(self.CONFIRM_PASSWORD_FLUTTER_KEY, self.CONFIRM_PASSWORD_FALLBACK, user_data['password'])
        
        if 'first_name' in user_data:
            self.enter_text_flutter_first(self.FIRST_NAME_FLUTTER_KEY, self.FIRST_NAME_FALLBACK, user_data['first_name'])
        
        if 'last_name' in user_data:
            self.enter_text_flutter_first(self.LAST_NAME_FLUTTER_KEY, self.LAST_NAME_FALLBACK, user_data['last_name'])
        
        if 'phone' in user_data:
            self.enter_text_flutter_first(self.PHONE_FLUTTER_KEY, self.PHONE_FALLBACK, user_data['phone'])
        
        if 'role' in user_data:
            self.select_role(user_data['role'])
    
    def select_role(self, role: str):
        """Backward-compatible role selector; taps specific register button instead."""
        self.tap_register_as(role)
    
    def tap_register_button(self):
        """Tap the register button using Flutter-first approach"""
        if self.is_element_present(self.REGISTER_AS_CUSTOMER_BUTTON, timeout=2):
            self.tap(self.REGISTER_AS_CUSTOMER_BUTTON)
        else:
            self.tap_flutter_first(self.REGISTER_FLUTTER_KEY, self.REGISTER_BUTTON_FALLBACK)

    def tap_register_as(self, role: str):
        """Tap the specific register button based on role"""
        role = role.lower()
        if role in ('customer', 'user') and self.is_element_present(self.REGISTER_AS_CUSTOMER_BUTTON, timeout=2):
            self.tap(self.REGISTER_AS_CUSTOMER_BUTTON)
        elif role in ('provider', 'agent') and self.is_element_present(self.REGISTER_AS_AGENT_BUTTON, timeout=2):
            self.tap(self.REGISTER_AS_AGENT_BUTTON)
        else:
            # Fallback to generic button
            self.tap(self.REGISTER_BUTTON)
    
    def tap_login_link(self):
        """Tap the login link to navigate to login page"""
        self.tap(self.choose_locator(self.LOGIN_LINK, self.LOGIN_LINK_FALLBACK))
    
    def wait_for_registration_complete(self, timeout: int = 30):
        """Wait for registration to complete"""
        self.wait_for_loading_to_complete(timeout)
    
    def get_validation_errors(self) -> Dict[str, str]:
        """Get all visible validation errors"""
        errors = {}
        
        error_mappings = {
            'email': self.EMAIL_ERROR,
            'password': self.PASSWORD_ERROR,
            'first_name': self.FIRST_NAME_ERROR,
            'last_name': self.LAST_NAME_ERROR,
            'phone': self.PHONE_ERROR,
            'role': self.ROLE_ERROR,
            'password_mismatch': self.PASSWORD_MISMATCH_ERROR
        }
        
        for field, locator in error_mappings.items():
            if self.is_element_present(locator, timeout=2):
                errors[field] = self.get_text(locator)
        
        return errors
    
    def is_register_button_enabled(self) -> bool:
        """Check if register button is enabled"""
        return self.is_element_enabled(self.REGISTER_BUTTON)


class LoginPage(BasePage):
    """Page object for the login screen"""
    
    # Form field locators
    EMAIL_FIELD = (AppiumBy.ACCESSIBILITY_ID, "login_email")
    EMAIL_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'email') or contains(@hint, 'Email')]")
    PASSWORD_FIELD = (AppiumBy.ACCESSIBILITY_ID, "login_password")
    PASSWORD_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'password') or contains(@hint, 'Password')]")
    
    # Button locators
    LOGIN_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "login_submit")
    LOGIN_BUTTON_FALLBACK = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Login') or contains(@text, 'Sign In')]")
    REGISTER_LINK = (AppiumBy.ACCESSIBILITY_ID, "login_register_link")
    REGISTER_LINK_FALLBACK = (AppiumBy.XPATH, "//*[contains(@text, 'Register') or contains(@text, 'Sign Up')]")
    # Forgot password now has a Semantics id in Flutter
    FORGOT_PASSWORD_LINK = (AppiumBy.ACCESSIBILITY_ID, "login_forgot_password")
    FORGOT_PASSWORD_ANDROID_UIA = (AppiumBy.ANDROID_UIAUTOMATOR, 'new UiSelector().textContains("Forgot Password")')
    FORGOT_PASSWORD_FALLBACK = (AppiumBy.XPATH, "//*[contains(@text, 'Forgot Password')]")
    
    # Validation error locators
    EMAIL_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'email') and contains(@text, 'required')]")
    PASSWORD_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'password') and contains(@text, 'required')]")
    INVALID_CREDENTIALS_ERROR = (AppiumBy.XPATH, "//*[contains(@text, 'Invalid') and (contains(@text, 'email') or contains(@text, 'password') or contains(@text, 'credentials'))]")
    
    # Loading states
    LOADING_INDICATOR = (AppiumBy.CLASS_NAME, "CircularProgressIndicator")
    
    def _find_edit_text_within_semantics(self, semantics_id: str):
        try:
            container = self.find_element_by_accessibility_id(semantics_id)
            try:
                return container.find_element(AppiumBy.CLASS_NAME, "android.widget.EditText")
            except Exception:
                # Fallback to descendant search via XPath scoped to container
                return container.find_element(AppiumBy.XPATH, ".//android.widget.EditText")
        except Exception:
            return None

    def _find_input_by_label_text(self, label_text_candidates: list):
        for candidate in label_text_candidates:
            try:
                label = self.find_element_by_xpath(f"//*[contains(@text, '{candidate}')]")
                # Assume the input follows the label in the view hierarchy
                return label.find_element(AppiumBy.XPATH, "following::android.widget.EditText[1]")
            except Exception:
                continue
        return None

    def fill_login_form(self, email: str, password: str):
        """Fill login form with credentials using robust strategies."""
        # Email/username field
        email_input = (
            self._find_edit_text_within_semantics("login_email")
            or self._find_input_by_label_text(["Username or Email", "Email", "Username"])  # label-based
        )
        if email_input is not None:
            email_input.click()
            try:
                email_input.clear()
            except Exception:
                pass
            email_input.send_keys(email)
        else:
            # Final fallback to broad XPath
            self.enter_text((AppiumBy.XPATH, "//android.widget.EditText"), email)

        # Password field
        password_input = (
            self._find_edit_text_within_semantics("login_password")
            or self._find_input_by_label_text(["Password"])  # label-based
        )
        if password_input is not None:
            password_input.click()
            try:
                password_input.clear()
            except Exception:
                pass
            password_input.send_keys(password)
        else:
            # Final fallback: last EditText is typically password
            try:
                fields = self.driver.find_elements(AppiumBy.CLASS_NAME, "android.widget.EditText")
                if fields:
                    fields[-1].click()
                    try:
                        fields[-1].clear()
                    except Exception:
                        pass
                    fields[-1].send_keys(password)
                else:
                    # As a last resort, use generic locator
                    self.enter_text((AppiumBy.XPATH, "(//android.widget.EditText)[last()]"), password)
            except Exception:
                self.enter_text((AppiumBy.XPATH, "(//android.widget.EditText)[last()]"), password)
    
    def tap_login_button(self):
        """Tap the login button"""
        self.tap(self.choose_locator(self.LOGIN_BUTTON, self.LOGIN_BUTTON_FALLBACK))
    
    def tap_register_link(self):
        """Tap the register link to navigate to registration page"""
        self.tap(self.choose_locator(self.REGISTER_LINK, self.REGISTER_LINK_FALLBACK))
    
    def tap_forgot_password_link(self):
        """Tap the forgot password link"""
        # Hide keyboard if obstructing
        self.hide_keyboard()
        locator = self.choose_locator(
            self.FORGOT_PASSWORD_LINK,
            self.FORGOT_PASSWORD_ANDROID_UIA,
            self.FORGOT_PASSWORD_FALLBACK,
        )
        try:
            self.tap(locator)
        except Exception:
            # Attempt a small scroll and retry
            try:
                self.scroll_down(400)
                self.tap(locator)
            except Exception:
                # Final attempt: use text contains generic XPath
                alt = (AppiumBy.XPATH, "//*[contains(@text, 'Forgot')]")
                self.tap(alt)

    def ensure_on_login(self, timeout: int = 10):
        """Ensure we're on the login screen by navigating from landing or registration if necessary."""
        if self.is_element_present(self.LOGIN_BUTTON, timeout=3):
            return
        # Try landing → login
        try:
            landing = PageFactory.get_landing_page(self.driver)
            landing.goto_login()
        except Exception:
            pass
        # If still not, try registration → login link
        if not self.is_element_present(self.LOGIN_BUTTON, timeout=3):
            try:
                registration = PageFactory.get_registration_page(self.driver)
                registration.tap_login_link()
            except Exception:
                pass
        # Wait briefly for login button
        self.is_element_present(self.LOGIN_BUTTON, timeout=timeout)
    
    def wait_for_login_complete(self, timeout: int = 30):
        """Wait for login to complete"""
        self.wait_for_loading_to_complete(timeout)
    
    def get_validation_errors(self) -> Dict[str, str]:
        """Get all visible validation errors"""
        errors = {}
        
        if self.is_element_present(self.EMAIL_ERROR, timeout=2):
            errors['email'] = self.get_text(self.EMAIL_ERROR)
        
        if self.is_element_present(self.PASSWORD_ERROR, timeout=2):
            errors['password'] = self.get_text(self.PASSWORD_ERROR)
        
        if self.is_element_present(self.INVALID_CREDENTIALS_ERROR, timeout=2):
            errors['credentials'] = self.get_text(self.INVALID_CREDENTIALS_ERROR)
        
        return errors
    
    def is_login_button_enabled(self) -> bool:
        """Check if login button is enabled"""
        return self.is_element_enabled(self.LOGIN_BUTTON)


class EmailExistsErrorWidget(BasePage):
    """Page object for the email exists error widget"""
    
    # Error message locators
    ERROR_MESSAGE = (AppiumBy.XPATH, "//*[contains(@text, 'email') and contains(@text, 'already') and contains(@text, 'used')]")
    ERROR_TITLE = (AppiumBy.XPATH, "//*[contains(@text, 'Email Already Exists') or contains(@text, 'User already exists')]")
    
    # Action button locators
    CREATE_ANOTHER_USER_BUTTON = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Create Another User')]")
    LOGIN_BUTTON = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Login')]")
    CLOSE_BUTTON = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Close') or contains(@content-desc, 'close')]")
    
    def is_error_widget_visible(self) -> bool:
        """Check if error widget is visible"""
        return self.is_element_visible(self.ERROR_MESSAGE, timeout=10)
    
    def get_error_message(self) -> str:
        """Get the error message text"""
        return self.get_text(self.ERROR_MESSAGE)
    
    def tap_create_another_user(self):
        """Tap the 'Create Another User' button"""
        self.tap(self.CREATE_ANOTHER_USER_BUTTON)
    
    def tap_login_button(self):
        """Tap the 'Login' button"""
        self.tap(self.LOGIN_BUTTON)
    
    def tap_close_button(self):
        """Tap the close button"""
        self.tap(self.CLOSE_BUTTON)
    
    def wait_for_error_widget_to_appear(self, timeout: int = 10):
        """Wait for error widget to appear"""
        self.wait_for_element_visible(self.ERROR_MESSAGE, timeout)
    
    def wait_for_error_widget_to_disappear(self, timeout: int = 10):
        """Wait for error widget to disappear"""
        self.wait_for_element_to_disappear(self.ERROR_MESSAGE, timeout)


class DashboardPage(BasePage):
    """Page object for dashboard/home screen after successful auth"""
    
    # Common dashboard elements
    WELCOME_MESSAGE = (AppiumBy.XPATH, "//*[contains(@text, 'Welcome') or contains(@text, 'Dashboard')]")
    USER_AVATAR = (AppiumBy.XPATH, "//android.widget.ImageView[contains(@content-desc, 'avatar') or contains(@content-desc, 'profile')]")
    NAVIGATION_DRAWER = (AppiumBy.XPATH, "//android.widget.Button[contains(@content-desc, 'menu') or contains(@content-desc, 'drawer')]")
    LOGOUT_BUTTON = (AppiumBy.XPATH, "//*[contains(@text, 'Logout') or contains(@text, 'Sign Out')]")
    
    # Role-specific elements
    CUSTOMER_SERVICES_TAB = (AppiumBy.XPATH, "//*[contains(@text, 'Services') or contains(@text, 'Browse')]")
    PROVIDER_DASHBOARD_TAB = (AppiumBy.XPATH, "//*[contains(@text, 'My Services') or contains(@text, 'Provider')]")
    ADMIN_PANEL_TAB = (AppiumBy.XPATH, "//*[contains(@text, 'Admin') or contains(@text, 'Management')]")
    
    def is_dashboard_loaded(self) -> bool:
        """Check if dashboard is loaded"""
        return self.is_element_present(self.WELCOME_MESSAGE, timeout=15)
    
    def get_welcome_message(self) -> str:
        """Get welcome message text"""
        return self.get_text(self.WELCOME_MESSAGE)
    
    def tap_navigation_drawer(self):
        """Open navigation drawer"""
        self.tap(self.NAVIGATION_DRAWER)
    
    def logout(self):
        """Logout from the app"""
        self.tap_navigation_drawer()
        self.tap(self.LOGOUT_BUTTON)
    
    def verify_role_specific_elements(self, role: str) -> bool:
        """Verify role-specific elements are visible"""
        role_elements = {
            'customer': self.CUSTOMER_SERVICES_TAB,
            'provider': self.PROVIDER_DASHBOARD_TAB,
            'admin': self.ADMIN_PANEL_TAB
        }
        
        element = role_elements.get(role.lower())
        if element:
            return self.is_element_present(element, timeout=10)
        return False


class ErrorDialog(BasePage):
    """Page object for generic error dialogs"""
    
    # Generic error dialog locators
    ERROR_DIALOG = (AppiumBy.XPATH, "//android.widget.LinearLayout[contains(@resource-id, 'dialog') or @class='android.app.AlertDialog']")
    ERROR_TITLE = (AppiumBy.XPATH, "//*[contains(@resource-id, 'title') or contains(@text, 'Error')]")
    ERROR_MESSAGE = (AppiumBy.XPATH, "//*[contains(@resource-id, 'message') or contains(@resource-id, 'content')]")
    OK_BUTTON = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'OK') or contains(@text, 'Ok')]")
    CANCEL_BUTTON = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Cancel')]")
    RETRY_BUTTON = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Retry') or contains(@text, 'Try Again')]")
    
    # Network error specific
    NETWORK_ERROR_MESSAGE = (AppiumBy.XPATH, "//*[contains(@text, 'network') or contains(@text, 'connection') or contains(@text, 'internet')]")
    SERVER_ERROR_MESSAGE = (AppiumBy.XPATH, "//*[contains(@text, 'server') or contains(@text, 'Server')]")
    
    def is_error_dialog_visible(self) -> bool:
        """Check if error dialog is visible"""
        return self.is_element_visible(self.ERROR_DIALOG, timeout=10)
    
    def get_error_title(self) -> str:
        """Get error dialog title"""
        return self.get_text(self.ERROR_TITLE)
    
    def get_error_message(self) -> str:
        """Get error dialog message"""
        return self.get_text(self.ERROR_MESSAGE)
    
    def tap_ok_button(self):
        """Tap OK button"""
        self.tap(self.OK_BUTTON)
    
    def tap_cancel_button(self):
        """Tap Cancel button"""
        self.tap(self.CANCEL_BUTTON)
    
    def tap_retry_button(self):
        """Tap Retry button"""
        self.tap(self.RETRY_BUTTON)
    
    def is_network_error(self) -> bool:
        """Check if this is a network error"""
        return self.is_element_present(self.NETWORK_ERROR_MESSAGE, timeout=5)
    
    def is_server_error(self) -> bool:
        """Check if this is a server error"""
        return self.is_element_present(self.SERVER_ERROR_MESSAGE, timeout=5)
    
    def dismiss_error(self):
        """Dismiss error dialog by tapping OK or Cancel"""
        if self.is_element_present(self.OK_BUTTON, timeout=5):
            self.tap_ok_button()
        elif self.is_element_present(self.CANCEL_BUTTON, timeout=5):
            self.tap_cancel_button()


# Page factory for creating page objects
class PageFactory:
    """Factory class for creating page objects"""
    
    @staticmethod
    def get_splash_page(driver) -> SplashPage:
        return SplashPage(driver)
    
    @staticmethod
    def get_registration_page(driver) -> RegistrationPage:
        return RegistrationPage(driver)
    
    @staticmethod
    def get_login_page(driver) -> LoginPage:
        return LoginPage(driver)

    @staticmethod
    def get_landing_page(driver) -> 'LandingPage':
        return LandingPage(driver)
    
    @staticmethod
    def get_dashboard_page(driver) -> DashboardPage:
        return DashboardPage(driver)
    
    @staticmethod
    def get_email_exists_error_widget(driver) -> EmailExistsErrorWidget:
        return EmailExistsErrorWidget(driver)
    
    @staticmethod
    def get_error_dialog(driver) -> ErrorDialog:
        return ErrorDialog(driver)

    @staticmethod
    def get_otp_page(driver):
        from .auth_pages import OTPVerificationPage  # Local import to avoid early reference
        return OTPVerificationPage(driver)

    @staticmethod
    def get_forgot_password_page(driver):
        from .auth_pages import ForgotPasswordPage  # Local import to avoid early reference
        return ForgotPasswordPage(driver)


class OTPVerificationPage(BasePage):
    """Page object for the OTP verification screen with Flutter Driver support"""

    # Flutter keys for better element discovery
    OTP_FIELD_FLUTTER_KEY = "otp_field"
    RESEND_FLUTTER_KEY = "otp_resend_button"
    VERIFY_FLUTTER_KEY = "otp_verify_button"
    COUNTDOWN_FLUTTER_KEY = "otp_countdown_label"

    # OTP input fields (6 digits) - fallbacks
    OTP_EDIT_TEXTS = (AppiumBy.CLASS_NAME, "android.widget.EditText")
    RESEND_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "otp_resend_button")
    RESEND_FALLBACK = (AppiumBy.XPATH, "//*[contains(@text, 'Resend') or contains(@text, 'Send Again')]")
    VERIFY_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "otp_verify_button")
    VERIFY_FALLBACK = (AppiumBy.XPATH, "//*[contains(@text, 'Verify') or contains(@text, 'Confirm')]")
    COUNTDOWN_LABEL = (AppiumBy.ACCESSIBILITY_ID, "otp_countdown_label")
    COUNTDOWN_FALLBACK = (AppiumBy.XPATH, "//*[contains(@text, 'expires') or contains(@text, 'resend')]")

    def enter_otp(self, otp: str):
        """Enter the 6-digit OTP code using Flutter-first approach"""
        assert len(otp) == 6, "OTP must be 6 digits"
        
        # Try Flutter approach first
        try:
            if self.flutter_finder and FLUTTER_DRIVER_AVAILABLE:
                # Flutter apps often have a single OTP field or auto-advance
                self.enter_text_flutter_first(self.OTP_FIELD_FLUTTER_KEY, self.OTP_EDIT_TEXTS, otp)
                return
        except Exception:
            pass
        
        # Fallback to traditional multi-field approach
        fields = self.driver.find_elements(*self.OTP_EDIT_TEXTS)
        for index, digit in enumerate(otp):
            if index < len(fields):
                fields[index].click()
                try:
                    fields[index].clear()
                except Exception:
                    pass
                fields[index].send_keys(digit)
            else:
                elem = self.driver.find_elements(*self.OTP_EDIT_TEXTS)[-1]
                elem.click()
                elem.send_keys(digit)

    def is_resend_enabled(self) -> bool:
        """Check if resend button is enabled using Flutter-first approach"""
        try:
            element = self.find_element_flutter_first(self.RESEND_FLUTTER_KEY, self.RESEND_FALLBACK, timeout=2)
            return element.is_enabled()
        except Exception:
            return False

    def get_countdown_text(self) -> str:
        """Get countdown text using Flutter-first approach"""
        try:
            return self.get_text_flutter_first(self.COUNTDOWN_FLUTTER_KEY, self.COUNTDOWN_FALLBACK, timeout=2)
        except Exception:
            try:
                return self.get_text((AppiumBy.XPATH, "//*[contains(@text, 'Code expires in') or contains(@text, 'expires in') ]"))
            except Exception:
                return ""

    def tap_verify(self):
        """Tap verify button using Flutter-first approach"""
        self.tap_flutter_first(self.VERIFY_FLUTTER_KEY, self.VERIFY_FALLBACK)


class ForgotPasswordPage(BasePage):
    """Page object for the forgot password screen"""

    EMAIL_FIELD = (AppiumBy.ACCESSIBILITY_ID, "forgot_email")
    EMAIL_FALLBACK = (AppiumBy.XPATH, "//android.widget.EditText[contains(@content-desc, 'email') or contains(@hint, 'Email')]")
    SUBMIT_BUTTON = (AppiumBy.ACCESSIBILITY_ID, "forgot_submit")
    SUBMIT_FALLBACK = (AppiumBy.XPATH, "//android.widget.Button[contains(@text, 'Submit') or contains(@text, 'Continue')]")

    def wait_for_loaded(self, timeout: int = 10):
        try:
            self.wait_for_element(self.EMAIL_FIELD, timeout)
        except Exception:
            # Fallback: wait for any EditText to appear
            self.wait_for_element((AppiumBy.CLASS_NAME, "android.widget.EditText"), timeout)

    def fill_email(self, email: str):
        # Prefer the explicit email locator, otherwise use the first EditText on the page
        try:
            locator = self.choose_locator(self.EMAIL_FIELD, self.EMAIL_FALLBACK)
            self.enter_text(locator, email)
            return
        except Exception:
            pass
        fields = self.driver.find_elements(AppiumBy.CLASS_NAME, "android.widget.EditText")
        if fields:
            fields[0].click()
            try:
                fields[0].clear()
            except Exception:
                pass
            fields[0].send_keys(email)
        else:
            raise NoSuchElementException("No input field found on Forgot Password page")

    def submit(self):
        locator = self.choose_locator(self.SUBMIT_BUTTON, self.SUBMIT_FALLBACK)
        self.tap(locator)
