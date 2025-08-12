# API Error Documentation - Smor-Ting Backend

## Overview
This document provides comprehensive documentation of all possible error responses from the Smor-Ting API, organized by endpoint and error type.

## Error Response Format

All API errors follow a consistent JSON structure:

```json
{
  "error": "Error Type",
  "message": "Human-readable error description"
}
```

## HTTP Status Codes

- `200` - Success
- `201` - Created (successful registration)
- `400` - Bad Request (validation errors, malformed requests)
- `401` - Unauthorized (authentication failures)
- `409` - Conflict (resource already exists)
- `500` - Internal Server Error (unexpected server errors)

---

## Authentication Endpoints

### POST /api/v1/auth/register

#### Success Response (201)
```json
{
  "user": {
    "id": "64f123abc...",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "requires_otp": false
}
```

#### Error Responses

**400 - Invalid Request Body**
```json
{
  "error": "Invalid request body",
  "message": "Failed to parse request body"
}
```

**400 - Validation Errors**
```json
{
  "error": "Validation failed",
  "message": "email is required"
}
```

Possible validation messages:
- `"email is required"`
- `"password is required"`
- `"password must be at least 6 characters long"`
- `"first name is required"`
- `"last name is required"`
- `"phone is required"`
- `"role is required"`
- `"role must be 'customer', 'provider', or 'admin'"`

**409 - User Already Exists**
```json
{
  "error": "User already exists",
  "message": "A user with this email already exists"
}
```

**500 - Registration Failed**
```json
{
  "error": "Registration failed",
  "message": "Failed to register user"
}
```

---

### POST /api/v1/auth/login

#### Success Response (200)
```json
{
  "user": {
    "id": "64f123abc...",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "requires_otp": false
}
```

#### Error Responses

**400 - Invalid Request Body**
```json
{
  "error": "Invalid request body",
  "message": "Failed to parse request body"
}
```

**400 - Validation Errors**
```json
{
  "error": "Validation failed",
  "message": "email is required"
}
```

Possible validation messages:
- `"email is required"`
- `"password is required"`

**401 - Invalid Credentials**
```json
{
  "error": "Invalid credentials",
  "message": "Invalid email or password"
}
```

**500 - Authentication Failed**
```json
{
  "error": "Authentication failed",
  "message": "Failed to generate authentication tokens"
}
```

**500 - Login Failed**
```json
{
  "error": "Login failed",
  "message": "Failed to authenticate user"
}
```

---

### POST /api/v1/auth/validate

#### Success Response (200)
```json
{
  "valid": true,
  "user": {
    "id": "64f123abc...",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer"
  }
}
```

#### Error Responses

**400 - Invalid Request Body**
```json
{
  "error": "Invalid request body",
  "message": "Failed to parse request body"
}
```

**401 - Invalid Token**
```json
{
  "error": "Invalid token",
  "message": "The provided token is invalid or expired"
}
```

---

## Flutter Client Error Mapping

The Flutter client maps server errors to custom exceptions:

### EmailAlreadyExistsException
- **Triggered by**: 409 status with message containing "already exists" or "user already exists"
- **Client Message**: "This email is already being used in our system"
- **UI Behavior**: Shows custom error widget with "Create Another User" and "Login" buttons

### InvalidCredentialsException
- **Triggered by**: 401 status with message containing "invalid credentials"
- **Client Message**: "Invalid email or password. Please check your credentials and try again."
- **UI Behavior**: Shows error message in SnackBar

### AccountNotVerifiedException
- **Triggered by**: Custom implementation (future feature)
- **Client Message**: "Your account is not verified. Please check your email for verification instructions."
- **UI Behavior**: Navigate to verification page

### AuthException (Generic)
- **Triggered by**: All other authentication errors
- **Client Message**: Varies based on server response
- **UI Behavior**: Shows error message in SnackBar

---

## Network-Level Errors

### Connection Timeout
```
Client Message: "Connection timeout. Please check your internet connection."
```

### No Internet Connection
```
Client Message: "No internet connection. Please check your network."
```

### Request Cancelled
```
Client Message: "Request cancelled."
```

### Unexpected Error
```
Client Message: "An unexpected error occurred."
```

---

## Testing Scenarios

### Registration Endpoint Tests
1. ✅ Valid registration
2. ✅ Missing email
3. ✅ Missing password
4. ✅ Password too short (< 6 characters)
5. ✅ Missing first name
6. ✅ Missing last name
7. ✅ Missing phone
8. ✅ Missing role
9. ✅ Invalid role
10. ✅ Email already exists
11. ✅ Malformed JSON
12. ✅ Server error

### Login Endpoint Tests
1. ✅ Valid login
2. ✅ Missing email
3. ✅ Missing password
4. ✅ Invalid email (user doesn't exist)
5. ✅ Wrong password
6. ✅ Malformed JSON
7. ✅ Server error

### Token Validation Tests
1. ✅ Valid token
2. ✅ Invalid token
3. ✅ Expired token
4. ✅ Malformed token
5. ✅ Missing token
6. ✅ Malformed JSON

---

## Error Logging

All errors are logged with appropriate levels:

- **WARN**: Validation errors, invalid credentials
- **ERROR**: Server errors, database failures, token generation failures
- **INFO**: Successful operations

Example log entry:
```json
{
  "level": "warn",
  "timestamp": "2024-01-01T12:00:00.000Z",
  "message": "Invalid login request",
  "error": "email is required",
  "email": "",
  "endpoint": "/auth/login"
}
```

---

## Security Considerations

1. **Password Validation**: Minimum 6 characters (consider stronger requirements)
2. **Rate Limiting**: Implement to prevent brute force attacks
3. **Email Validation**: Server-side email format validation
4. **Input Sanitization**: Prevent injection attacks
5. **Token Expiration**: Implement proper token refresh mechanism

---

## Monitoring and Alerts

### Key Metrics to Monitor
- Authentication failure rate
- 409 conflicts (duplicate registrations)
- 500 errors (server failures)
- Response times
- Token validation failures

### Alert Thresholds
- Error rate > 5% over 5 minutes
- Response time > 2 seconds average
- 500 errors > 10 in 1 minute
