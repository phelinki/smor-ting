"""
Appium Configuration Module

This module provides centralized configuration management for the Smor-Ting
mobile application QA automation framework.
"""

from .appium_config import AppiumConfig, get_config, validate_config

__all__ = ["AppiumConfig", "get_config", "validate_config"]
