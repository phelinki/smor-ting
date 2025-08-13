import os
import pytest

from config.appium_config import get_config


@pytest.fixture(autouse=True)
def clear_env(monkeypatch):
    for k in [
        "PLATFORM",
        "ENVIRONMENT",
        "ANDROID_NO_RESET",
    ]:
        monkeypatch.delenv(k, raising=False)


def test_local_env_defaults_to_noreset(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    monkeypatch.setenv("ENVIRONMENT", "local")
    caps = get_config().get_capabilities()
    assert caps.get("noReset") is True
    assert caps.get("fullReset") is False


def test_ci_env_defaults_to_fullreset(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    monkeypatch.setenv("ENVIRONMENT", "ci")
    caps = get_config().get_capabilities()
    assert caps.get("noReset") is False
    assert caps.get("fullReset") is True


def test_timeouts_are_at_least_300s(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    caps = get_config().get_capabilities()
    assert int(caps.get("adbExecTimeout", 0)) >= 300000
    assert int(caps.get("uiautomator2ServerLaunchTimeout", 0)) >= 300000
    assert int(caps.get("uiautomator2ServerInstallTimeout", 0)) >= 300000
    assert int(caps.get("androidInstallTimeout", 0)) >= 300000


