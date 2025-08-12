"""
Base Page Object for Smor-Ting Mobile App

This module provides the base page object class that all page objects inherit from,
implementing common page interactions and utilities with Flutter Driver support.
"""
import time
from typing import Any, Optional, Tuple
from appium.webdriver.common.appiumby import AppiumBy
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException

# Flutter Driver imports
try:
    from appium_flutter_finder.flutter_finder import FlutterFinder
    FLUTTER_DRIVER_AVAILABLE = True
except ImportError:
    FLUTTER_DRIVER_AVAILABLE = False
    FlutterFinder = None


class BasePage:
    """Base page object providing common functionality for all pages"""
    
    def __init__(self, driver):
        self.driver = driver
        self.wait = WebDriverWait(driver, 30)  # Default wait time
        self.flutter_finder = FlutterFinder() if FLUTTER_DRIVER_AVAILABLE else None
        
    # Element location strategies
    def find_element_by_id(self, element_id: str) -> Any:
        """Find element by resource ID"""
        return self.driver.find_element(AppiumBy.ID, element_id)
    
    def find_element_by_xpath(self, xpath: str) -> Any:
        """Find element by XPath"""
        return self.driver.find_element(AppiumBy.XPATH, xpath)
    
    def find_element_by_text(self, text: str) -> Any:
        """Find element by exact text"""
        return self.driver.find_element(AppiumBy.XPATH, f"//*[@text='{text}']")
    
    def find_element_by_partial_text(self, partial_text: str) -> Any:
        """Find element by partial text"""
        return self.driver.find_element(AppiumBy.XPATH, f"//*[contains(@text, '{partial_text}')]")
    
    def find_element_by_class(self, class_name: str) -> Any:
        """Find element by class name"""
        return self.driver.find_element(AppiumBy.CLASS_NAME, class_name)
    
    def find_elements_by_class(self, class_name: str) -> list:
        """Find multiple elements by class name"""
        return self.driver.find_elements(AppiumBy.CLASS_NAME, class_name)
    
    def find_element_by_accessibility_id(self, accessibility_id: str) -> Any:
        """Find element by accessibility ID"""
        return self.driver.find_element(AppiumBy.ACCESSIBILITY_ID, accessibility_id)

    def choose_locator(self, *locators: Tuple[str, str], timeout_each: int = 1) -> Tuple[str, str]:
        """Return the first locator that resolves to a present element, otherwise the first locator.

        This enables resilient locator strategies with accessibility-id first and XPath/text fallbacks.
        """
        for locator in locators:
            try:
                if self.is_element_present(locator, timeout=timeout_each):
                    return locator
            except Exception:
                continue
        # Fallback to first provided locator when none found immediately
        return locators[0]
    
    # Wait methods
    def wait_for_element(self, locator: Tuple[str, str], timeout: int = 30) -> Any:
        """Wait for element to be present"""
        wait = WebDriverWait(self.driver, timeout)
        return wait.until(EC.presence_of_element_located(locator))
    
    def wait_for_element_clickable(self, locator: Tuple[str, str], timeout: int = 30) -> Any:
        """Wait for element to be clickable"""
        wait = WebDriverWait(self.driver, timeout)
        return wait.until(EC.element_to_be_clickable(locator))
    
    def wait_for_element_visible(self, locator: Tuple[str, str], timeout: int = 30) -> Any:
        """Wait for element to be visible"""
        wait = WebDriverWait(self.driver, timeout)
        return wait.until(EC.visibility_of_element_located(locator))
    
    def wait_for_text_in_element(self, locator: Tuple[str, str], text: str, timeout: int = 30) -> bool:
        """Wait for text to appear in element"""
        wait = WebDriverWait(self.driver, timeout)
        return wait.until(EC.text_to_be_present_in_element(locator, text))
    
    def wait_for_element_to_disappear(self, locator: Tuple[str, str], timeout: int = 30):
        """Wait for element to disappear"""
        wait = WebDriverWait(self.driver, timeout)
        wait.until_not(EC.presence_of_element_located(locator))
    
    # Interaction methods
    def tap(self, locator: Tuple[str, str], timeout: int = 30):
        """Tap on element"""
        element = self.wait_for_element_clickable(locator, timeout)
        element.click()
    
    def enter_text(self, locator: Tuple[str, str], text: str, clear_first: bool = True, timeout: int = 30):
        """Enter text into element"""
        element = self.wait_for_element(locator, timeout)
        if clear_first:
            element.clear()
        element.send_keys(text)
    
    def get_text(self, locator: Tuple[str, str], timeout: int = 30) -> str:
        """Get text from element"""
        element = self.wait_for_element(locator, timeout)
        return element.text
    
    def get_attribute(self, locator: Tuple[str, str], attribute: str, timeout: int = 30) -> str:
        """Get attribute value from element"""
        element = self.wait_for_element(locator, timeout)
        return element.get_attribute(attribute)
    
    def is_element_present(self, locator: Tuple[str, str], timeout: int = 5) -> bool:
        """Check if element is present"""
        try:
            self.wait_for_element(locator, timeout)
            return True
        except TimeoutException:
            return False
    
    def is_element_visible(self, locator: Tuple[str, str], timeout: int = 5) -> bool:
        """Check if element is visible"""
        try:
            self.wait_for_element_visible(locator, timeout)
            return True
        except TimeoutException:
            return False
    
    def is_element_enabled(self, locator: Tuple[str, str], timeout: int = 30) -> bool:
        """Check if element is enabled"""
        try:
            element = self.wait_for_element(locator, timeout)
            return element.is_enabled()
        except TimeoutException:
            return False
    
    # Scrolling methods
    def scroll_down(self, duration: int = 1000):
        """Scroll down on the screen"""
        size = self.driver.get_window_size()
        start_x = size['width'] // 2
        start_y = size['height'] * 0.8
        end_y = size['height'] * 0.2
        
        self.driver.swipe(start_x, start_y, start_x, end_y, duration)
    
    def scroll_up(self, duration: int = 1000):
        """Scroll up on the screen"""
        size = self.driver.get_window_size()
        start_x = size['width'] // 2
        start_y = size['height'] * 0.2
        end_y = size['height'] * 0.8
        
        self.driver.swipe(start_x, start_y, start_x, end_y, duration)
    
    def scroll_to_element(self, locator: Tuple[str, str], max_scrolls: int = 10) -> Any:
        """Scroll to find and return element"""
        for _ in range(max_scrolls):
            if self.is_element_present(locator, timeout=2):
                return self.wait_for_element(locator)
            self.scroll_down()
        
        raise NoSuchElementException(f"Element not found after {max_scrolls} scrolls: {locator}")
    
    # Navigation methods
    def go_back(self):
        """Press back button"""
        self.driver.back()
    
    def hide_keyboard(self):
        """Hide keyboard if visible"""
        try:
            if self.driver.is_keyboard_shown():
                self.driver.hide_keyboard()
        except Exception:
            # Some platforms don't support keyboard detection
            pass
    
    # Assertion helpers
    def assert_element_present(self, locator: Tuple[str, str], message: str = None):
        """Assert element is present"""
        assert self.is_element_present(locator), message or f"Element not present: {locator}"
    
    def assert_element_visible(self, locator: Tuple[str, str], message: str = None):
        """Assert element is visible"""
        assert self.is_element_visible(locator), message or f"Element not visible: {locator}"
    
    def assert_text_equals(self, locator: Tuple[str, str], expected_text: str, message: str = None):
        """Assert element text equals expected text"""
        actual_text = self.get_text(locator)
        assert actual_text == expected_text, message or f"Expected '{expected_text}', got '{actual_text}'"
    
    def assert_text_contains(self, locator: Tuple[str, str], expected_text: str, message: str = None):
        """Assert element text contains expected text"""
        actual_text = self.get_text(locator)
        assert expected_text in actual_text, message or f"'{expected_text}' not found in '{actual_text}'"
    
    def assert_element_enabled(self, locator: Tuple[str, str], message: str = None):
        """Assert element is enabled"""
        assert self.is_element_enabled(locator), message or f"Element not enabled: {locator}"
    
    def assert_element_disabled(self, locator: Tuple[str, str], message: str = None):
        """Assert element is disabled"""
        assert not self.is_element_enabled(locator), message or f"Element should be disabled: {locator}"
    
    # Utility methods
    def take_screenshot(self, name: str = None) -> str:
        """Take screenshot and return filename"""
        if not name:
            timestamp = int(time.time())
            name = f"screenshot_{timestamp}"
        
        filename = f"{name}.png"
        self.driver.save_screenshot(filename)
        return filename
    
    def get_page_source(self) -> str:
        """Get current page source"""
        return self.driver.page_source
    
    def wait_for_loading_to_complete(self, timeout: int = 30):
        """Wait for loading indicators to disappear"""
        # Common loading indicators
        loading_indicators = [
            (AppiumBy.CLASS_NAME, "CircularProgressIndicator"),
            (AppiumBy.XPATH, "//*[contains(@text, 'Loading')]"),
            (AppiumBy.XPATH, "//*[contains(@text, 'Please wait')]"),
            (AppiumBy.ACCESSIBILITY_ID, "loading"),
        ]
        
        for locator in loading_indicators:
            try:
                # Wait for loading indicator to appear first (optional)
                self.wait_for_element(locator, timeout=5)
                # Then wait for it to disappear
                self.wait_for_element_to_disappear(locator, timeout)
            except TimeoutException:
                # Loading indicator might not appear or already gone
                continue
    
    def wait_for_network_idle(self, timeout: int = 10):
        """Wait for network requests to complete"""
        # Simple wait - can be enhanced with network monitoring
        time.sleep(2)
        self.wait_for_loading_to_complete(timeout)
    
    # Platform-specific helpers
    def is_android(self) -> bool:
        """Check if running on Android"""
        return self.driver.capabilities.get('platformName', '').lower() == 'android'
    
    def is_ios(self) -> bool:
        """Check if running on iOS"""
        return self.driver.capabilities.get('platformName', '').lower() == 'ios'
    
    # Flutter-specific methods with fallbacks
    def find_element_flutter_first(self, flutter_key: str, fallback_locator: Tuple[str, str], timeout: int = 30) -> Any:
        """Find element using Flutter Driver first, fallback to UiAutomator2"""
        if self.flutter_finder and FLUTTER_DRIVER_AVAILABLE:
            try:
                element_locator = self.flutter_finder.by_value_key(flutter_key)
                return WebDriverWait(self.driver, timeout).until(
                    lambda driver: driver.find_element(AppiumBy.FLUTTER, element_locator)
                )
            except Exception:
                # Fall back to traditional locator
                pass
        
        # Use traditional locator strategy
        return self.wait_for_element(fallback_locator, timeout)
    
    def tap_flutter_first(self, flutter_key: str, fallback_locator: Tuple[str, str], timeout: int = 30):
        """Tap element using Flutter Driver first, fallback to UiAutomator2"""
        element = self.find_element_flutter_first(flutter_key, fallback_locator, timeout)
        element.click()
    
    def enter_text_flutter_first(self, flutter_key: str, fallback_locator: Tuple[str, str], text: str, clear_first: bool = True, timeout: int = 30):
        """Enter text using Flutter Driver first, fallback to UiAutomator2"""
        element = self.find_element_flutter_first(flutter_key, fallback_locator, timeout)
        if clear_first:
            element.clear()
        element.send_keys(text)
    
    def get_text_flutter_first(self, flutter_key: str, fallback_locator: Tuple[str, str], timeout: int = 30) -> str:
        """Get text using Flutter Driver first, fallback to UiAutomator2"""
        element = self.find_element_flutter_first(flutter_key, fallback_locator, timeout)
        return element.text or element.get_attribute('text') or ""
    
    def is_element_present_flutter_first(self, flutter_key: str, fallback_locator: Tuple[str, str], timeout: int = 5) -> bool:
        """Check if element is present using Flutter Driver first, fallback to UiAutomator2"""
        try:
            self.find_element_flutter_first(flutter_key, fallback_locator, timeout)
            return True
        except Exception:
            return False
    
    def get_platform_specific_locator(self, android_locator: Tuple[str, str], ios_locator: Tuple[str, str]) -> Tuple[str, str]:
        """Get platform-specific locator"""
        return android_locator if self.is_android() else ios_locator
