"""
Base Test Class for Appium QA Automation

This module provides the base test class that all test classes inherit from,
implementing common setup, teardown, and utility methods following TDD principles.
"""
import os
import sys
import time
import pytest
import logging
from typing import Optional, Dict, Any
from pathlib import Path
from appium import webdriver
from appium.options.android import UiAutomator2Options
from appium.options.ios import XCUITestOptions
from appium.webdriver.common.appiumby import AppiumBy
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, WebDriverException

# Flutter Driver imports
try:
    from appium_flutter_finder.flutter_finder import FlutterFinder
    from appium_flutter_finder.flutter_element import FlutterElement
    FLUTTER_DRIVER_AVAILABLE = True
except ImportError:
    FLUTTER_DRIVER_AVAILABLE = False
    FlutterFinder = None
    FlutterElement = None

# Add config to path
sys.path.append(str(Path(__file__).parent.parent))
from config.appium_config import get_config, validate_config


class BaseTest:
    """Base test class providing common functionality for all tests"""
    
    driver: Optional[webdriver.Remote] = None
    config = None
    logger = None
    
    @classmethod
    def setup_class(cls):
        """Class-level setup - runs once per test class"""
        cls.config = get_config()
        cls.logger = cls._setup_logging()
        cls.logger.info(f"Starting test class: {cls.__name__}")
        
        # Validate configuration
        if not validate_config(cls.config):
            pytest.skip("Configuration validation failed")
    
    @classmethod
    def teardown_class(cls):
        """Class-level teardown - runs once per test class"""
        if cls.logger:
            cls.logger.info(f"Completed test class: {cls.__name__}")
    
    def setup_method(self, method):
        """Method-level setup - runs before each test method"""
        self.logger.info(f"Starting test: {method.__name__}")
        self.start_time = time.time()
        
        # Start driver
        self.driver = self._create_driver()
        if not self.driver:
            pytest.fail("Failed to create Appium driver")
        
        # Set implicit wait
        self.driver.implicitly_wait(self.config.implicit_wait)
        
        # Wait for app to launch
        self._wait_for_app_launch()
    
    def teardown_method(self, method):
        """Method-level teardown - runs after each test method"""
        test_duration = time.time() - self.start_time
        self.logger.info(f"Test {method.__name__} completed in {test_duration:.2f}s")
        
        # Take screenshot on failure
        if hasattr(self, '_outcome') and self._outcome.errors:
            self._take_failure_screenshot(method.__name__)
        
        # Quit driver
        if self.driver:
            try:
                self.driver.quit()
            except Exception as e:
                self.logger.error(f"Error quitting driver: {e}")
            finally:
                self.driver = None
    
    def _create_driver(self) -> Optional[webdriver.Remote]:
        """Create and return Appium WebDriver instance"""
        try:
            capabilities = self.config.get_capabilities()
            self.logger.info(f"Creating driver with capabilities: {capabilities}")

            # Choose platform-specific driver options
            if self.config.platform_name.lower() == "ios":
                options = XCUITestOptions()
            else:
                options = UiAutomator2Options()
            options.load_capabilities(capabilities)

            driver = webdriver.Remote(
                command_executor=self.config.appium_server_url,
                options=options,
            )

            self.logger.info("Driver created successfully")
            return driver

        except Exception as e:
            self.logger.error(f"Failed to create driver: {e}")
            return None
    
    def _wait_for_app_launch(self, timeout: int = 30):
        """Wait for app to launch completely"""
        try:
            # Wait for any element to appear (app has launched)
            WebDriverWait(self.driver, timeout).until(
                lambda driver: len(driver.find_elements(AppiumBy.XPATH, "//*")) > 0
            )
            self.logger.info("App launched successfully")
            
            # Additional wait for app to stabilize
            time.sleep(2)
            
            # Make startup deterministic: skip onboarding and ensure landing is reachable
            try:
                from .common.page_objects.auth_pages import PageFactory
                splash = PageFactory.get_splash_page(self.driver)
                splash.wait_for_splash_to_complete()
                onboarding = PageFactory.get_onboarding_page(self.driver)
                onboarding.skip_if_present()
                landing = PageFactory.get_landing_page(self.driver)
                if not landing.ensure_loaded(timeout=5):
                    # Try once more after brief delay
                    time.sleep(2)
                    onboarding.skip_if_present()
                    landing.ensure_loaded(timeout=5)
            except Exception as e:
                self.logger.warning(f"Deterministic startup sequence encountered an issue: {e}")

        except TimeoutException:
            self.logger.error("Timeout waiting for app to launch")
            raise
    
    def _take_failure_screenshot(self, test_name: str):
        """Take screenshot on test failure"""
        try:
            if self.driver:
                screenshot_dir = Path(__file__).parent.parent / "reports" / "screenshots"
                screenshot_dir.mkdir(parents=True, exist_ok=True)
                
                timestamp = int(time.time())
                filename = f"{test_name}_failure_{timestamp}.png"
                filepath = screenshot_dir / filename
                
                self.driver.save_screenshot(str(filepath))
                self.logger.info(f"Failure screenshot saved: {filepath}")
                
                # Also save page source for debugging
                source_file = filepath.with_suffix('.xml')
                with open(source_file, 'w', encoding='utf-8') as f:
                    f.write(self.driver.page_source)
                
        except Exception as e:
            self.logger.error(f"Failed to capture failure screenshot: {e}")
    
    @classmethod
    def _setup_logging(cls) -> logging.Logger:
        """Setup logging for the test class"""
        logger = logging.getLogger(cls.__name__)
        logger.setLevel(logging.INFO)
        
        # Create handler if not exists
        if not logger.handlers:
            # File handler
            log_dir = Path(__file__).parent.parent / "reports" / "logs"
            log_dir.mkdir(parents=True, exist_ok=True)
            
            log_file = log_dir / f"{cls.__name__}.log"
            file_handler = logging.FileHandler(log_file)
            file_handler.setLevel(logging.INFO)
            
            # Console handler
            console_handler = logging.StreamHandler()
            console_handler.setLevel(logging.INFO)
            
            # Formatter
            formatter = logging.Formatter(
                '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
            )
            file_handler.setFormatter(formatter)
            console_handler.setFormatter(formatter)
            
            logger.addHandler(file_handler)
            logger.addHandler(console_handler)
        
        return logger
    
    # Utility methods for common operations
    
    def wait_for_element(self, locator: tuple, timeout: int = None) -> Any:
        """Wait for element to be present and return it"""
        timeout = timeout or self.config.explicit_wait
        try:
            element = WebDriverWait(self.driver, timeout).until(
                EC.presence_of_element_located(locator)
            )
            return element
        except TimeoutException:
            self.logger.error(f"Element not found: {locator}")
            raise
    
    def wait_for_element_clickable(self, locator: tuple, timeout: int = None) -> Any:
        """Wait for element to be clickable and return it"""
        timeout = timeout or self.config.explicit_wait
        try:
            element = WebDriverWait(self.driver, timeout).until(
                EC.element_to_be_clickable(locator)
            )
            return element
        except TimeoutException:
            self.logger.error(f"Element not clickable: {locator}")
            raise

    # Backwards-compatible alias used in some tests
    def wait_for_element_to_be_clickable(self, locator: tuple, timeout: int = None) -> Any:
        return self.wait_for_element_clickable(locator, timeout)
    
    def wait_for_text_in_element(self, locator: tuple, text: str, timeout: int = None) -> bool:
        """Wait for specific text to appear in element"""
        timeout = timeout or self.config.explicit_wait
        try:
            return WebDriverWait(self.driver, timeout).until(
                EC.text_to_be_present_in_element(locator, text)
            )
        except TimeoutException:
            self.logger.error(f"Text '{text}' not found in element: {locator}")
            return False
    
    def scroll_to_element(self, locator: tuple) -> Any:
        """Scroll to element and return it"""
        try:
            element = self.driver.find_element(*locator)
            self.driver.execute_script("arguments[0].scrollIntoView();", element)
            return element
        except Exception as e:
            self.logger.error(f"Failed to scroll to element: {e}")
            raise
    
    def tap_element(self, locator: tuple, timeout: int = None):
        """Tap on element with wait"""
        element = self.wait_for_element_clickable(locator, timeout)
        element.click()
        self.logger.info(f"Tapped element: {locator}")
    
    def enter_text(self, locator: tuple, text: str, clear_first: bool = True):
        """Enter text into element"""
        element = self.wait_for_element(locator)
        if clear_first:
            element.clear()
        element.send_keys(text)
        self.logger.info(f"Entered text '{text}' into element: {locator}")
    
    def get_element_text(self, locator: tuple) -> str:
        """Get text from element"""
        element = self.wait_for_element(locator)
        text = element.text
        self.logger.info(f"Got text '{text}' from element: {locator}")
        return text
    
    def is_element_present(self, locator: tuple, timeout: int = 5) -> bool:
        """Check if element is present without throwing exception"""
        try:
            WebDriverWait(self.driver, timeout).until(
                EC.presence_of_element_located(locator)
            )
            return True
        except TimeoutException:
            return False
    
    def is_element_visible(self, locator: tuple, timeout: int = 5) -> bool:
        """Check if element is visible"""
        try:
            WebDriverWait(self.driver, timeout).until(
                EC.visibility_of_element_located(locator)
            )
            return True
        except TimeoutException:
            return False
    
    def wait_for_element_to_disappear(self, locator: tuple, timeout: int = None):
        """Wait for element to disappear"""
        timeout = timeout or self.config.explicit_wait
        try:
            WebDriverWait(self.driver, timeout).until_not(
                EC.presence_of_element_located(locator)
            )
        except TimeoutException:
            self.logger.error(f"Element did not disappear: {locator}")
            raise
    
    def get_current_activity(self) -> str:
        """Get current activity (Android) or app state (iOS)"""
        if self.config.platform_name == "android":
            return self.driver.current_activity
        else:
            return self.driver.query_app_state(self.config.get_capabilities()["bundleId"])
    
    def restart_app(self):
        """Restart the app"""
        try:
            self.driver.terminate_app(self.config.get_capabilities().get("appPackage", ""))
            time.sleep(2)
            self.driver.activate_app(self.config.get_capabilities().get("appPackage", ""))
            self._wait_for_app_launch()
        except Exception as e:
            self.logger.error(f"Failed to restart app: {e}")
    
    # Flutter-specific methods
    def get_flutter_finder(self) -> Optional[FlutterFinder]:
        """Get Flutter finder instance for Flutter-specific element finding"""
        if not FLUTTER_DRIVER_AVAILABLE:
            self.logger.warning("Flutter driver not available, falling back to UiAutomator2")
            return None
        try:
            return FlutterFinder()
        except Exception as e:
            self.logger.error(f"Failed to create Flutter finder: {e}")
            return None
    
    def find_flutter_element_by_key(self, key: str, timeout: int = 10):
        """Find Flutter element by key"""
        flutter_finder = self.get_flutter_finder()
        if not flutter_finder:
            # Fallback to accessibility ID for UiAutomator2
            return self.wait_for_element((AppiumBy.ACCESSIBILITY_ID, key), timeout)
        
        try:
            element_locator = flutter_finder.by_value_key(key)
            return WebDriverWait(self.driver, timeout).until(
                lambda driver: driver.find_element(AppiumBy.FLUTTER, element_locator)
            )
        except Exception as e:
            self.logger.warning(f"Flutter element not found by key '{key}', trying accessibility ID: {e}")
            return self.wait_for_element((AppiumBy.ACCESSIBILITY_ID, key), timeout)
    
    def find_flutter_element_by_text(self, text: str, timeout: int = 10):
        """Find Flutter element by text"""
        flutter_finder = self.get_flutter_finder()
        if not flutter_finder:
            # Fallback to XPath for UiAutomator2
            return self.wait_for_element((AppiumBy.XPATH, f"//*[@text='{text}']"), timeout)
        
        try:
            element_locator = flutter_finder.by_text(text)
            return WebDriverWait(self.driver, timeout).until(
                lambda driver: driver.find_element(AppiumBy.FLUTTER, element_locator)
            )
        except Exception as e:
            self.logger.warning(f"Flutter element not found by text '{text}', trying XPath: {e}")
            return self.wait_for_element((AppiumBy.XPATH, f"//*[@text='{text}']"), timeout)
    
    def find_flutter_element_by_type(self, widget_type: str, timeout: int = 10):
        """Find Flutter element by widget type"""
        flutter_finder = self.get_flutter_finder()
        if not flutter_finder:
            # Fallback to class name for UiAutomator2
            return self.wait_for_element((AppiumBy.CLASS_NAME, widget_type), timeout)
        
        try:
            element_locator = flutter_finder.by_type(widget_type)
            return WebDriverWait(self.driver, timeout).until(
                lambda driver: driver.find_element(AppiumBy.FLUTTER, element_locator)
            )
        except Exception as e:
            self.logger.warning(f"Flutter element not found by type '{widget_type}': {e}")
            raise
    
    def tap_flutter_element_by_key(self, key: str, timeout: int = 10):
        """Tap Flutter element by key"""
        element = self.find_flutter_element_by_key(key, timeout)
        element.click()
        self.logger.info(f"Tapped Flutter element with key: {key}")
    
    def enter_text_flutter_by_key(self, key: str, text: str, clear_first: bool = True, timeout: int = 10):
        """Enter text into Flutter element by key"""
        element = self.find_flutter_element_by_key(key, timeout)
        if clear_first:
            element.clear()
        element.send_keys(text)
        self.logger.info(f"Entered text '{text}' into Flutter element with key: {key}")
    
    def is_flutter_element_present(self, key: str, timeout: int = 5) -> bool:
        """Check if Flutter element is present by key"""
        try:
            self.find_flutter_element_by_key(key, timeout)
            return True
        except Exception:
            return False
            raise
    
    def reset_app(self):
        """Reset app to initial state"""
        try:
            self.driver.reset()
            self._wait_for_app_launch()
        except Exception as e:
            self.logger.error(f"Failed to reset app: {e}")
            raise
    
    def take_screenshot(self, name: str = None) -> str:
        """Take screenshot and return path"""
        try:
            screenshot_dir = Path(__file__).parent.parent / "reports" / "screenshots"
            screenshot_dir.mkdir(parents=True, exist_ok=True)
            
            timestamp = int(time.time())
            filename = f"{name or 'screenshot'}_{timestamp}.png"
            filepath = screenshot_dir / filename
            
            self.driver.save_screenshot(str(filepath))
            self.logger.info(f"Screenshot saved: {filepath}")
            return str(filepath)
            
        except Exception as e:
            self.logger.error(f"Failed to take screenshot: {e}")
            return ""
    
    def assert_element_present(self, locator: tuple, message: str = None):
        """Assert that element is present"""
        assert self.is_element_present(locator), message or f"Element not present: {locator}"
    
    def assert_element_visible(self, locator: tuple, message: str = None):
        """Assert that element is visible"""
        assert self.is_element_visible(locator), message or f"Element not visible: {locator}"
    
    def assert_text_in_element(self, locator: tuple, expected_text: str, message: str = None):
        """Assert that element contains expected text"""
        actual_text = self.get_element_text(locator)
        assert expected_text in actual_text, message or f"Expected '{expected_text}' in '{actual_text}'"
    
    def assert_element_not_present(self, locator: tuple, message: str = None):
        """Assert that element is not present"""
        assert not self.is_element_present(locator), message or f"Element should not be present: {locator}"


# Pytest fixtures
@pytest.fixture(scope="session", autouse=True)
def session_setup():
    """Session-level setup"""
    # Create reports directory
    reports_dir = Path(__file__).parent.parent / "reports"
    reports_dir.mkdir(exist_ok=True)
    (reports_dir / "screenshots").mkdir(exist_ok=True)
    (reports_dir / "logs").mkdir(exist_ok=True)


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_makereport(item, call):
    """Hook to capture test outcome for screenshot on failure"""
    outcome = yield
    rep = outcome.get_result()
    setattr(item, "rep_" + rep.when, rep)
    
    # Store outcome in test instance for access in teardown
    if hasattr(item.instance, 'teardown_method'):
        item.instance._outcome = rep
