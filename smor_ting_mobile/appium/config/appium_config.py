"""
Appium Configuration for Smor-Ting Mobile QA Automation

This module provides comprehensive configuration management for Appium testing
across Android and iOS platforms, supporting both local and CI/CD environments.
"""
import os
import platform
from pathlib import Path
from typing import Dict, Any, Optional


class AppiumConfig:
    """Centralized configuration management for Appium testing"""
    
    def __init__(self, platform_name: str = "android", environment: str = "local"):
        self.platform_name = platform_name.lower()
        self.environment = environment.lower()
        self.project_root = Path(__file__).parent.parent.parent
        self.app_build_path = self._get_app_build_path()
        
    def _get_app_build_path(self) -> str:
        """Get the path to the built Flutter app"""
        # Allow override via APP_PATH for flexibility in CI/local
        env_app_path = os.getenv("APP_PATH")
        if env_app_path:
            return env_app_path

        if self.platform_name == "android":
            return str(self.project_root / "build" / "app" / "outputs" / "flutter-apk" / "app-debug.apk")
        else:  # iOS
            return str(self.project_root / "build" / "ios" / "iphonesimulator" / "Runner.app")
    
    def get_android_capabilities(self) -> Dict[str, Any]:
        """Get Android-specific capabilities"""
        # Use UiAutomator2 for native Android UI testing, Flutter for Flutter-specific tests
        automation_name = os.getenv("ANDROID_AUTOMATION_NAME", "UiAutomator2")
        capabilities = {
            "platformName": "Android", 
            "platformVersion": os.getenv("ANDROID_API_LEVEL", "30"),
            "deviceName": os.getenv("ANDROID_DEVICE_NAME", "Android Emulator"),
            "automationName": automation_name,
            "app": self.app_build_path,
            "appPackage": "com.smorting.app.smor_ting_mobile",
            "appActivity": "com.smorting.app.smor_ting_mobile.MainActivity",
            "autoGrantPermissions": True,
            "noReset": False,
            "fullReset": True,
            "newCommandTimeout": 300,
            "androidInstallTimeout": 120000,
            "uiautomator2ServerInstallTimeout": 120000,
            "uiautomator2ServerLaunchTimeout": 120000,
            "adbExecTimeout": 120000,
            "avd": os.getenv("ANDROID_AVD_NAME", os.getenv("ANDROID_DEVICE_NAME", "Medium_Phone_API_36.0")),
            "avdLaunchTimeout": 180000,
            "avdReadyTimeout": 180000,
            "systemPort": self._get_system_port(),
            # Avoid preinstalling Appium Settings/Unlock if emulator policy blocks it
            "skipDeviceInitialization": True,
            "ignoreHiddenApiPolicyError": True,
            "disableWindowAnimation": True,
        }
        
        # Add CI-specific capabilities
        if self.environment == "ci":
            capabilities.update({
                "isHeadless": True,
                "avdArgs": "-no-audio -no-window -gpu swiftshader_indirect",
                "deviceReadyTimeout": 120,
                "androidDeviceReadyTimeout": 120,
            })
            
        return capabilities
    
    def get_ios_capabilities(self) -> Dict[str, Any]:
        """Get iOS-specific capabilities"""
        capabilities = {
            "platformName": "iOS",
            "platformVersion": os.getenv("IOS_VERSION", "16.4"),
            "deviceName": os.getenv("IOS_DEVICE_NAME", "iPhone 13"),
            "automationName": "XCUITest",
            "app": self.app_build_path,
            "bundleId": "com.smorting.app.smor-ting-mobile",
            "noReset": False,
            "fullReset": True,
            "newCommandTimeout": 300,
            "wdaLaunchTimeout": 60000,
            "wdaConnectionTimeout": 60000,
            "iosInstallPause": 8000,
            "xcodeOrgId": os.getenv("XCODE_ORG_ID"),
            "xcodeSigningId": os.getenv("XCODE_SIGNING_ID", "iPhone Developer"),
            "updatedWDABundleId": os.getenv("WDA_BUNDLE_ID"),
            "usePrebuiltWDA": True,
            "shouldUseSingletonTestManager": False,
        }
        
        # Add simulator-specific settings
        if "Simulator" in capabilities["deviceName"] or self.environment == "ci":
            capabilities.update({
                "isSimulator": True,
                "simulatorStartupTimeout": 120000,
                "useSimulatorPasteboard": True,
            })
            
        return capabilities
    
    def get_capabilities(self) -> Dict[str, Any]:
        """Get platform-specific capabilities"""
        base_capabilities = {
            "orientation": "PORTRAIT",
            "unicodeKeyboard": True,
            "resetKeyboard": True,
            "clearSystemFiles": True,
        }
        
        if self.platform_name == "android":
            platform_caps = self.get_android_capabilities()
        else:
            platform_caps = self.get_ios_capabilities()
            
        # Merge capabilities
        base_capabilities.update(platform_caps)
        return base_capabilities
    
    def _get_system_port(self) -> int:
        """Get available system port for Android"""
        base_port = 8200
        # Prefer explicit offset when provided (e.g., CI matrix)
        explicit_offset = os.getenv("SYSTEM_PORT_OFFSET")
        if explicit_offset and explicit_offset.isdigit():
            return base_port + int(explicit_offset)

        # Derive from pytest-xdist worker id when running in parallel
        worker_id = os.getenv("PYTEST_XDIST_WORKER_ID", "master")
        worker_num = int(worker_id.replace("gw", "")) if "gw" in worker_id else 0

        # Add a small per-process entropy to avoid collisions when multiple sessions run serially
        try:
            pid_entropy = os.getpid() % 97  # prime-ish spread within a small band
        except Exception:
            pid_entropy = 0

        return base_port + worker_num + pid_entropy
    
    # Removed ChromeDriver port management since we're not driving webviews directly
    
    @property
    def appium_server_url(self) -> str:
        """Get Appium server URL"""
        host = os.getenv("APPIUM_HOST", "127.0.0.1")
        port = os.getenv("APPIUM_PORT", "4723")
        return f"http://{host}:{port}"
    
    @property
    def test_timeout(self) -> int:
        """Get test timeout in seconds"""
        return int(os.getenv("TEST_TIMEOUT", "300"))
    
    @property
    def implicit_wait(self) -> int:
        """Get implicit wait time in seconds"""
        return int(os.getenv("IMPLICIT_WAIT", "10"))
    
    @property
    def explicit_wait(self) -> int:
        """Get explicit wait time in seconds"""
        return int(os.getenv("EXPLICIT_WAIT", "30"))
    
    def get_test_data(self) -> Dict[str, Any]:
        """Get test data configuration"""
        return {
            "valid_users": [
                {
                    "email": "qa_customer@smorting.com",
                    "password": "TestPass123!",
                    "first_name": "QA",
                    "last_name": "Customer",
                    "phone": "231777123456",
                    "role": "customer"
                },
                {
                    "email": "qa_provider@smorting.com",
                    "password": "ProviderPass123!",
                    "first_name": "QA",
                    "last_name": "Provider", 
                    "phone": "231888123456",
                    "role": "provider"
                }
            ],
            "existing_user": {
                "email": "existing@smorting.com",
                "password": "ExistingPass123!",
                "first_name": "Existing",
                "last_name": "User",
                "phone": "231999123456",
                "role": "customer"
            },
            "invalid_data": {
                "emails": [
                    "",
                    "invalid-email",
                    "test@",
                    "@domain.com",
                    "spaces in@email.com"
                ],
                "passwords": [
                    "",
                    "123",
                    "short",
                    "nouppercaseorspecial",
                    "NOLOWERCASEORSPECIAL"
                ],
                "phones": [
                    "",
                    "123",
                    "1234567890123456",
                    "abcdefghijk",
                    "555-1234"
                ]
            }
        }
    
    def get_api_config(self) -> Dict[str, str]:
        """Get API configuration"""
        return {
            "base_url": os.getenv("API_BASE_URL", "http://localhost:8080/api/v1"),
            "timeout": os.getenv("API_TIMEOUT", "30"),
            "retry_count": os.getenv("API_RETRY_COUNT", "3"),
        }
    
    def get_environment_info(self) -> Dict[str, str]:
        """Get environment information for reporting"""
        return {
            "platform": self.platform_name,
            "environment": self.environment,
            "os": platform.system(),
            "python_version": platform.python_version(),
            "project_root": str(self.project_root),
            "app_path": self.app_build_path,
            "appium_url": self.appium_server_url,
        }


# Global configuration instances
def get_config(platform_name: str = None, environment: str = None) -> AppiumConfig:
    """Get configuration instance"""
    platform_name = platform_name or os.getenv("PLATFORM", "android")
    environment = environment or os.getenv("ENVIRONMENT", "local")
    return AppiumConfig(platform_name, environment)


# Configuration validation
def validate_config(config: AppiumConfig) -> bool:
    """Validate configuration setup"""
    errors = []
    
    # Check if app exists
    if not Path(config.app_build_path).exists():
        errors.append(f"App not found at: {config.app_build_path}")
    
    # Check platform-specific requirements
    if config.platform_name == "android":
        android_home = os.getenv("ANDROID_HOME")
        if not android_home:
            errors.append("ANDROID_HOME environment variable not set")
        elif not Path(android_home).exists():
            errors.append(f"ANDROID_HOME path does not exist: {android_home}")
    
    elif config.platform_name == "ios":
        if platform.system() != "Darwin":
            errors.append("iOS testing requires macOS")
    
    # Report errors
    if errors:
        allow_missing = os.getenv("ALLOW_MISSING_APP", "0") == "1"
        print("Configuration validation warnings:")
        for error in errors:
            print(f"  ⚠️ {error}")
        if allow_missing:
            print("Proceeding despite warnings because ALLOW_MISSING_APP=1")
            return True
        print("Set ALLOW_MISSING_APP=1 to bypass this check for dry-runs.")
        return False
    
    print("✅ Configuration validation passed")
    return True


if __name__ == "__main__":
    # Test configuration
    config = get_config()
    print("Current configuration:")
    for key, value in config.get_environment_info().items():
        print(f"  {key}: {value}")
    
    print("\nValidating configuration...")
    validate_config(config)
