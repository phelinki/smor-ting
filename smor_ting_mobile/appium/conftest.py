"""
Pytest Configuration for Smor-Ting Mobile QA Automation

This module provides pytest configuration, fixtures, and hooks for the
automated testing framework.
"""
import os
import sys
import pytest
import logging
from pathlib import Path

# Add project root to Python path
project_root = Path(__file__).parent
sys.path.insert(0, str(project_root))

from config.appium_config import get_config, validate_config
from appium import webdriver
from appium.options.android import UiAutomator2Options
from appium.options.ios import XCUITestOptions


def pytest_addoption(parser):
    """Add custom command line options"""
    parser.addoption(
        "--platform",
        action="store",
        default="android",
        help="Platform to test: android or ios"
    )
    parser.addoption(
        "--environment",
        action="store", 
        default="local",
        help="Test environment: local, ci, staging, production"
    )
    parser.addoption(
        "--device-name",
        action="store",
        default=None,
        help="Specific device name to use for testing"
    )
    parser.addoption(
        "--app-path",
        action="store",
        default=None,
        help="Path to the app binary"
    )
    parser.addoption(
        "--appium-url",
        action="store",
        default="http://127.0.0.1:4723",
        help="Appium server URL"
    )


def pytest_configure(config):
    """Configure pytest"""
    # Set environment variables from command line options
    platform = config.getoption("--platform")
    environment = config.getoption("--environment")
    device_name = config.getoption("--device-name")
    app_path = config.getoption("--app-path")
    appium_url = config.getoption("--appium-url")
    
    os.environ["PLATFORM"] = platform
    os.environ["ENVIRONMENT"] = environment
    
    if device_name:
        if platform == "android":
            os.environ["ANDROID_DEVICE_NAME"] = device_name
        else:
            os.environ["IOS_DEVICE_NAME"] = device_name
    
    if app_path:
        os.environ["APP_PATH"] = app_path
    
    if appium_url:
        os.environ["APPIUM_URL"] = appium_url
    
    # Configure logging
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    
    # Create reports directory
    reports_dir = Path(__file__).parent / "reports"
    reports_dir.mkdir(exist_ok=True)
    (reports_dir / "screenshots").mkdir(exist_ok=True)
    (reports_dir / "logs").mkdir(exist_ok=True)
    (reports_dir / "allure-results").mkdir(exist_ok=True)

    # Register custom markers with descriptions
    marker_descriptions = [
        "smoke: Smoke tests - critical functionality",
        "regression: Regression tests - full test suite",
        "auth: Authentication tests",
        "registration: Registration tests",
        "login: Login tests",
        "biometric: Biometric authentication tests",
        "flutter: Flutter-specific tests",
        "stress: Stress tests",
        "performance: Performance tests",
        "network: Network and connectivity tests",
        "ui: UI/UX tests",
        "security: Security tests",
        "e2e: End-to-end tests",
    ]
    for desc in marker_descriptions:
        config.addinivalue_line("markers", desc)


def pytest_sessionstart(session):
    """Session start hook"""
    print(f"\n🚀 Starting Smor-Ting Mobile QA Automation")
    print(f"Platform: {os.getenv('PLATFORM', 'android')}")
    print(f"Environment: {os.getenv('ENVIRONMENT', 'local')}")
    
    # Validate configuration
    config = get_config()
    if not validate_config(config):
        pytest.exit("Configuration validation failed")


def pytest_sessionfinish(session, exitstatus):
    """Session finish hook"""
    print(f"\n✅ Test session completed with exit status: {exitstatus}")


def pytest_runtest_setup(item):
    """Test setup hook"""
    # Log test start
    logging.getLogger("pytest").info(f"Starting test: {item.name}")


def pytest_runtest_teardown(item, nextitem):
    """Test teardown hook"""
    # Log test completion
    logging.getLogger("pytest").info(f"Completed test: {item.name}")


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_makereport(item, call):
    """Create test report hook"""
    outcome = yield
    rep = outcome.get_result()
    
    # Add extra information to report
    if rep.when == "call":
        # Add platform info to report
        rep.platform = os.getenv('PLATFORM', 'android')
        rep.environment = os.getenv('ENVIRONMENT', 'local')
        
        # Log test result
        status = "PASSED" if rep.passed else "FAILED" if rep.failed else "SKIPPED"
        logging.getLogger("pytest").info(f"Test {item.name}: {status}")


def pytest_html_report_title(report):
    """Customize HTML report title"""
    platform = os.getenv('PLATFORM', 'android').title()
    environment = os.getenv('ENVIRONMENT', 'local').title()
    report.title = f"Smor-Ting Mobile QA Report - {platform} ({environment})"


def pytest_html_results_summary(prefix, summary, postfix):
    """Customize HTML report summary"""
    platform = os.getenv('PLATFORM', 'android').title()
    environment = os.getenv('ENVIRONMENT', 'local').title()
    
    prefix.extend([
        f"<p><strong>Platform:</strong> {platform}</p>",
        f"<p><strong>Environment:</strong> {environment}</p>",
        f"<p><strong>Appium Server:</strong> {os.getenv('APPIUM_URL', 'http://127.0.0.1:4723')}</p>"
    ])


# Remove nonstandard dynamic marker creation; markers are registered in pytest_configure


# Fixtures
@pytest.fixture(scope="session")
def test_config():
    """Get test configuration"""
    return get_config()


@pytest.fixture(scope="session") 
def test_data(test_config):
    """Get test data"""
    return test_config.get_test_data()


@pytest.fixture(scope="function")
def driver(test_config):
    """Provide a raw Appium driver for tests that request it (e.g., simple function tests)."""
    if not validate_config(test_config):
        pytest.skip("Invalid Appium configuration")

    capabilities = test_config.get_capabilities()
    platform_name = capabilities.get("platformName", "Android").lower()
    options = None
    caps = None
    if platform_name == "android":
        options = UiAutomator2Options()
        options.load_capabilities(capabilities)
    elif platform_name == "ios":
        options = XCUITestOptions()
        options.load_capabilities(capabilities)
    else:
        caps = capabilities

    drv = None
    try:
        if options is not None:
            drv = webdriver.Remote(
                command_executor=test_config.appium_server_url,
                options=options,
            )
        else:
            drv = webdriver.Remote(
                command_executor=test_config.appium_server_url,
                desired_capabilities=caps,
            )
        yield drv
    finally:
        try:
            if drv:
                drv.quit()
        except Exception:
            pass

@pytest.fixture(scope="function")
def screenshot_on_failure(request):
    """Take screenshot on test failure"""
    yield
    
    if request.node.rep_call.failed:
        # Screenshot will be taken by base test class
        pass


def pytest_collection_modifyitems(config, items):
    """Modify test collection"""
    # Add default markers based on test file location
    for item in items:
        # Add auth marker to all auth tests
        if "auth" in str(item.fspath):
            item.add_marker(pytest.mark.auth)
        
        # Add specific markers based on test name
        if "registration" in item.name:
            item.add_marker(pytest.mark.registration)
        elif "login" in item.name:
            item.add_marker(pytest.mark.login)
        
        if "performance" in item.name:
            item.add_marker(pytest.mark.performance)
        elif "network" in item.name:
            item.add_marker(pytest.mark.network)
        elif "ui" in item.name or "loading" in item.name:
            item.add_marker(pytest.mark.ui)
        elif "security" in item.name:
            item.add_marker(pytest.mark.security)


def pytest_runtest_setup(item):
    """Conditional skips prior to each test based on platform/environment"""
    platform = os.getenv('PLATFORM', 'android').lower()
    environment = os.getenv('ENVIRONMENT', 'local').lower()

    # Skip iOS tests if not on macOS host
    if platform == "ios" and sys.platform != "darwin":
        pytest.skip("iOS tests can only run on macOS")

    # Skip slow tests in CI environment
    if environment == 'ci' and 'slow' in item.keywords:
        pytest.skip("Skipping slow test in CI environment")
