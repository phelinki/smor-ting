"""
Page Objects Module for Smor-Ting Mobile QA Automation

This module provides all page objects for the mobile application testing.
"""

from .base_page import BasePage
from .auth_pages import (
    SplashPage,
    LandingPage,
    RegistrationPage,
    LoginPage,
    DashboardPage,
    EmailExistsErrorWidget,
    ErrorDialog,
    PageFactory,
    OTPVerificationPage,
)

__all__ = [
    "BasePage",
    "SplashPage", 
    "LandingPage",
    "RegistrationPage",
    "LoginPage",
    "DashboardPage",
    "EmailExistsErrorWidget",
    "ErrorDialog",
    "PageFactory",
    "OTPVerificationPage",
]
