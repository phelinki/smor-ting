# iOS QA Automation Test Plan

## Scope

- Functional E2E tests for authentication flows: registration, login, OTP, and error UX
- Non-functional tests: performance checks on login/registration, basic security input validation, and UX checks (loading states, accessibility ids presence)

## Test Environment

- Device: iOS Simulator (e.g., iPhone 16 Pro Max, iOS 18.5)
- App build: Flutter simulator build (`Runner.app`)
- Tools: Appium 2.x (XCUITest), Pytest, pytest-html, Allure

## Entry Criteria

- iOS simulator build succeeds: `flutter build ios --simulator --debug` → `Runner.app` present
- Appium server reachable on `127.0.0.1:4723`

## Exit Criteria

- All smoke tests pass (no regressions in critical auth flows)
- Performance thresholds met (see below)

## Test Suites and Cases

### Authentication – Registration

- TC_IOS_REG_001: Successful registration (customer)
- TC_IOS_REG_002: Email already exists → custom error widget visible; CTA to login
- TC_IOS_REG_003: Create another user → form reset
- TC_IOS_REG_004: Missing required fields validation (email, password, first/last name, phone, role)
- TC_IOS_REG_005: Invalid formats (email/phone/password complexity)
- TC_IOS_REG_006: OTP flow (enter, resend, invalid attempt)
- TC_IOS_REG_007: Loading indicators show/hide appropriately
- TC_IOS_REG_008: Network error surfaced with friendly message

### Authentication – Login

- TC_IOS_LOGIN_001: Successful login → navigate to dashboard
- TC_IOS_LOGIN_002: Wrong password → invalid credentials UX
- TC_IOS_LOGIN_003: Empty field validation
- TC_IOS_LOGIN_004: Invalid email format blocked
- TC_IOS_LOGIN_005: Navigate to registration
- TC_IOS_LOGIN_006: Loading state & retry UX on network failure

### Non-Functional

- TC_IOS_PERF_001: Login completes in ≤ 3s on simulator
- TC_IOS_PERF_002: Registration completes in ≤ 5s (excluding OTP wait)
- TC_IOS_SEC_001: Inputs sanitized (no obvious script injection in text fields)
- TC_IOS_UI_001: Accessibility IDs present for key controls

## Execution Instructions

```bash
cd smor_ting_mobile/appium
./scripts/run_ios_tests.sh --suite all --device "iPhone 16 Pro Max" --ios-version 18.5
```

Use markers:

```bash
./scripts/run_ios_tests.sh --markers "smoke"
```

## Reporting

- HTML: `reports/ios-report.html`
- JUnit: `reports/ios-junit.xml`
- Logs: `reports/appium.log`
- Screenshots on failure: `reports/screenshots/`

## Risks/Assumptions

- The app must expose stable accessibility IDs per the page objects.
- OTP delivery is stubbed/mocked or available via test API for automation.

## Current Status

- Environment ready. Build currently blocked by Flutter compile-time errors in app source. Resolve build errors, then execute plan to completion.


