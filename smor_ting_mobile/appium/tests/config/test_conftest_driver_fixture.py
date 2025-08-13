from pathlib import Path
import importlib.util
import pytest

import config.appium_config as appium_config


def test_driver_fixture_is_defined_and_callable():
    # Load conftest.py directly by file path and assert the fixture callable exists
    conftest_path = Path(__file__).resolve().parents[2] / "conftest.py"
    spec = importlib.util.spec_from_file_location("qa_conftest", conftest_path)
    assert spec and spec.loader
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)  # type: ignore
    assert hasattr(module, "driver"), "driver fixture missing"
    assert callable(getattr(module, "driver"))


def test_get_config_returns_appium_config_instance(monkeypatch):
    monkeypatch.setenv("PLATFORM", "android")
    cfg = appium_config.get_config()
    assert isinstance(cfg, appium_config.AppiumConfig)


def test_base_test_has_alias_wait_method():
    # Ensure the alias method exists for compatibility with tests
    from tests.base_test import BaseTest
    assert hasattr(BaseTest, "wait_for_element_to_be_clickable")


