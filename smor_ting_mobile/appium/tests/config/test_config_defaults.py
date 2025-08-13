import os
from pathlib import Path

import pytest

from config.appium_config import get_config


@pytest.fixture(autouse=True)
def clear_android_env(monkeypatch):
    keys = [
        "ANDROID_AVD_NAME",
        "ANDROID_DEVICE_NAME",
        "ANDROID_API_LEVEL",
        "ANDROID_AUTOMATION_NAME",
        "APP_PATH",
        "PLATFORM",
        "ENVIRONMENT",
    ]
    for key in keys:
        monkeypatch.delenv(key, raising=False)


def test_default_avd_name_is_medium_phone_api_36(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    cfg = get_config()
    caps = cfg.get_capabilities()
    assert caps.get("avd") == "Medium_Phone_API_36.0"


def test_env_overrides_avd_name(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    monkeypatch.setenv("ANDROID_AVD_NAME", "My_Custom_AVD")
    cfg = get_config()
    caps = cfg.get_capabilities()
    assert caps.get("avd") == "My_Custom_AVD"


def test_default_automation_is_uiautomator2(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    cfg = get_config()
    caps = cfg.get_capabilities()
    assert caps.get("automationName") == "UiAutomator2"


def test_app_path_can_be_overridden_by_env(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    absolute_app_path = str(
        Path(__file__).parent.parent.parent.parent
        / "build"
        / "app"
        / "outputs"
        / "flutter-apk"
        / "app-debug.apk"
    )
    monkeypatch.setenv("APP_PATH", absolute_app_path)
    cfg = get_config()
    caps = cfg.get_capabilities()
    assert caps.get("app") == absolute_app_path


def test_default_android_platform_version_is_34(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    cfg = get_config()
    caps = cfg.get_capabilities()
    assert caps.get("platformVersion") == "34"


def test_system_port_is_dynamic_and_integer(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    cfg = get_config()
    caps = cfg.get_capabilities()
    port = caps.get("systemPort")
    assert isinstance(port, int)
    assert 8000 <= port <= 9000


