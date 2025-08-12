# Test Cases for Smor-Ting Authentication Flow

## Overview
This document outlines comprehensive test cases for the Smor-Ting mobile application authentication system, covering all possible error scenarios and user flows.

## Test Categories

### 1. Registration Tests

#### TC_REG_001: Successful Registration
**Objective**: Verify successful user registration with valid data
**Prerequisites**: App is installed and launched
**Test Data**: Valid user details
**Steps**:
1. Navigate to registration page
2. Enter valid email, password, first name, last name, phone, and role
3. Tap "Register" button
4. Wait for registration to complete

**Expected Result**: 
- User is successfully registered
- Navigation to appropriate dashboard based on role
- Success message displayed (if any)

**Error Scenarios Covered**: None

---

#### TC_REG_002: Email Already Exists
**Objective**: Verify proper handling when user tries to register with existing email
**Prerequisites**: User with test email already exists in system
**Test Data**: Existing user email
**Steps**:
1. Navigate to registration page
2. Enter existing email with other valid details
3. Tap "Register" button
4. Wait for response

**Expected Result**:
- Custom error widget displayed with message "This email is already being used in our system"
- Two action buttons visible: "Create Another User" and "Login"
- No navigation occurs

**API Error**: HTTP 409 - "User already exists"

---

#### TC_REG_003: Create Another User Flow
**Objective**: Verify "Create Another User" button functionality
**Prerequisites**: Email already exists error is displayed
**Steps**:
1. Follow TC_REG_002 to trigger error
2. Tap "Create Another User" button
3. Verify form state

**Expected Result**:
- All form fields are cleared
- User can enter new registration details
- Error widget disappears

---

#### TC_REG_004: Login from Error Widget
**Objective**: Verify "Login" button functionality from error widget
**Prerequisites**: Email already exists error is displayed
**Steps**:
1. Follow TC_REG_002 to trigger error
2. Tap "Login" button
3. Verify navigation

**Expected Result**:
- Navigation to login page
- Login form is displayed

---

#### TC_REG_005: Missing Email Validation
**Objective**: Verify validation when email field is empty
**Test Data**: Empty email field
**Steps**:
1. Navigate to registration page
2. Leave email field empty
3. Fill other required fields
4. Tap "Register" button

**Expected Result**:
- Validation error displayed: "Email is required"
- No API call made
- Form remains on registration page

---

#### TC_REG_006: Missing Password Validation
**Objective**: Verify validation when password field is empty
**Test Data**: Empty password field
**Steps**:
1. Navigate to registration page
2. Fill email and other fields
3. Leave password field empty
4. Tap "Register" button

**Expected Result**:
- Validation error: "Password is required"
- Form validation prevents submission

---

#### TC_REG_007: Password Too Short
**Objective**: Verify password length validation
**Test Data**: Password with less than 8 characters
**Steps**:
1. Navigate to registration page
2. Enter password shorter than 8 characters
3. Fill other fields
4. Tap "Register" button

**Expected Result**:
- Frontend validation: Password complexity requirements shown
- Backend validation (if reached): "Password must be at least 6 characters long"

---

#### TC_REG_008: Password Complexity Validation
**Objective**: Verify password meets complexity requirements
**Test Data**: Passwords missing required elements
**Test Cases**:
- No uppercase letter
- No lowercase letter  
- No number
- No special character

**Expected Result**:
- Frontend validation shows requirements
- User cannot proceed until requirements met

---

#### TC_REG_009: Missing First Name
**Objective**: Verify first name validation
**Steps**:
1. Leave first name field empty
2. Fill other required fields
3. Attempt registration

**Expected Result**:
- Validation error: "First name is required"

---

#### TC_REG_010: Missing Last Name
**Objective**: Verify last name validation
**Expected Result**:
- Validation error: "Last name is required"

---

#### TC_REG_011: Missing Phone Number
**Objective**: Verify phone number validation
**Expected Result**:
- Validation error: "Phone is required"

---

#### TC_REG_012: Invalid Phone Format
**Objective**: Verify Liberian phone number format validation
**Test Data**: 
- Non-Liberian numbers
- Invalid format numbers
- Too short/long numbers

**Expected Result**:
- Validation error with format hint
- Example format shown: "231777123456"

---

#### TC_REG_013: Missing Role Selection
**Objective**: Verify role selection validation
**Expected Result**:
- Validation error: "Role is required"

---

#### TC_REG_014: Invalid Role Value
**Objective**: Verify role value validation (backend)
**Test Data**: Invalid role value sent to API
**Expected Result**:
- Backend error: "Role must be 'customer', 'provider', or 'admin'"

---

#### TC_REG_015: Password Mismatch
**Objective**: Verify password confirmation validation
**Test Data**: Different values in password and confirm password
**Expected Result**:
- Frontend validation: "Passwords do not match"

---

### 2. Login Tests

#### TC_LOGIN_001: Successful Login
**Objective**: Verify successful login with valid credentials
**Prerequisites**: Valid user account exists
**Test Data**: Valid email and password
**Steps**:
1. Navigate to login page
2. Enter valid email and password
3. Tap "Login" button
4. Wait for authentication

**Expected Result**:
- User successfully authenticated
- Navigation to appropriate dashboard
- Access token stored securely

---

#### TC_LOGIN_002: Invalid Email
**Objective**: Verify login with non-existent email
**Test Data**: Email not registered in system
**Steps**:
1. Enter non-existent email
2. Enter any password
3. Attempt login

**Expected Result**:
- Error message: "Invalid email or password"
- No navigation occurs
- API returns HTTP 401

---

#### TC_LOGIN_003: Wrong Password
**Objective**: Verify login with incorrect password
**Test Data**: Valid email, incorrect password
**Expected Result**:
- Error message: "Invalid email or password"
- API returns HTTP 401

---

#### TC_LOGIN_004: Empty Email Field
**Objective**: Verify validation for empty email
**Expected Result**:
- Validation error: "Email is required"
- Form prevents submission

---

#### TC_LOGIN_005: Empty Password Field
**Objective**: Verify validation for empty password
**Expected Result**:
- Validation error: "Password is required"

---

#### TC_LOGIN_006: Invalid Email Format
**Objective**: Verify email format validation
**Test Data**: Malformed email addresses
**Examples**: "invalid-email", "test@", "@domain.com"
**Expected Result**:
- Frontend validation prevents submission
- Error message about email format

---

### 3. Network and Connectivity Tests

#### TC_NETWORK_001: No Internet Connection
**Objective**: Verify app behavior when offline
**Prerequisites**: Device has no internet connectivity
**Steps**:
1. Disable internet connection
2. Attempt registration or login
3. Wait for timeout

**Expected Result**:
- Error message: "No internet connection. Please check your network."
- User guidance provided

---

#### TC_NETWORK_002: Connection Timeout
**Objective**: Verify timeout handling
**Prerequisites**: Slow/unstable network
**Expected Result**:
- Error message: "Connection timeout. Please check your internet connection."
- Option to retry

---

#### TC_NETWORK_003: Server Error (500)
**Objective**: Verify handling of server errors
**Prerequisites**: Backend returns 500 error
**Expected Result**:
- Error message: "Server error. Please try again later."
- No app crash

---

### 4. UI/UX Tests

#### TC_UI_001: Loading States
**Objective**: Verify loading indicators during API calls
**Steps**:
1. Initiate registration or login
2. Observe UI during processing

**Expected Result**:
- Loading indicator displayed
- Form buttons disabled during processing
- User cannot submit multiple requests

---

#### TC_UI_002: Error Message Display
**Objective**: Verify error messages are user-friendly
**Expected Results**:
- Clear, actionable error messages
- Appropriate color coding (red for errors)
- Messages disappear when user corrects input

---

#### TC_UI_003: Form Field Validation
**Objective**: Verify real-time form validation
**Expected Results**:
- Validation occurs on field blur
- Clear indication of valid/invalid fields
- Password strength indicator

---

#### TC_UI_004: Accessibility
**Objective**: Verify app accessibility features
**Requirements**:
- Screen reader compatibility
- Proper semantic labels
- Keyboard navigation support
- Color contrast compliance

---

### 5. Security Tests

#### TC_SEC_001: Password Masking
**Objective**: Verify password fields are masked
**Expected Result**:
- Password characters hidden by default
- Option to toggle visibility

---

#### TC_SEC_002: Input Sanitization
**Objective**: Verify protection against injection attacks
**Test Data**: SQL injection, XSS attempts
**Expected Result**:
- Malicious input properly escaped
- No security vulnerabilities

---

#### TC_SEC_003: Token Security
**Objective**: Verify secure token storage
**Expected Result**:
- Tokens stored in secure storage
- Automatic token refresh (when implemented)
- Proper token expiration handling

---

### 6. Performance Tests

#### TC_PERF_001: Registration Performance
**Objective**: Verify registration completes within acceptable time
**Acceptance Criteria**: Registration < 5 seconds
**Steps**:
1. Measure time from button tap to completion
2. Test with various network conditions

---

#### TC_PERF_002: Login Performance
**Objective**: Verify login completes within acceptable time
**Acceptance Criteria**: Login < 3 seconds

---

#### TC_PERF_003: Form Responsiveness
**Objective**: Verify UI remains responsive during operations
**Expected Result**:
- No UI freezing
- Smooth animations
- Responsive user feedback

---

## Test Data Matrix

### Valid Test Data
```json
{
  "valid_users": [
    {
      "email": "customer@test.com",
      "password": "TestPass123!",
      "first_name": "John",
      "last_name": "Doe",
      "phone": "231777123456",
      "role": "customer"
    },
    {
      "email": "provider@test.com", 
      "password": "ProviderPass123!",
      "first_name": "Jane",
      "last_name": "Smith",
      "phone": "231888123456",
      "role": "provider"
    }
  ]
}
```

### Invalid Test Data
```json
{
  "invalid_emails": [
    "",
    "invalid-email",
    "test@",
    "@domain.com",
    "spaces in@email.com"
  ],
  "invalid_passwords": [
    "",
    "123",
    "short",
    "nouppercaseorspecial",
    "NOLOWERCASEORSPECIAL"
  ],
  "invalid_phones": [
    "",
    "123",
    "1234567890123456",
    "abcdefghijk",
    "555-1234"
  ]
}
```

## Automation Priorities

### High Priority (Critical Path)
1. Successful registration and login flows
2. Email already exists error handling
3. Basic form validation
4. Network error handling

### Medium Priority (Important Features)
1. Password complexity validation
2. Phone number format validation
3. UI loading states
4. Performance tests

### Low Priority (Edge Cases)
1. Malformed API responses
2. Concurrent user sessions
3. Memory pressure scenarios
4. Accessibility edge cases

## Test Environment Requirements

### Staging Environment
- Clean database state for each test run
- Predictable test users
- Stable API endpoints
- Logging enabled for debugging

### Production-like Testing
- Real network conditions
- Actual device testing
- Performance monitoring
- Security scanning

## Reporting Requirements

### Test Execution Reports
- Pass/fail status for each test case
- Execution time metrics
- Screenshot evidence for failures
- Error logs and stack traces

### Coverage Reports
- Functional coverage by feature
- Code coverage (if applicable)
- Risk-based testing coverage
- Regression test coverage

This comprehensive test suite ensures thorough validation of your authentication system following TDD principles [[memory:5654749]] and covering all error scenarios documented in the API specification.
