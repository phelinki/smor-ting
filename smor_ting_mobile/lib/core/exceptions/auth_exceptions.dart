/// Custom exceptions for authentication-related errors
class AuthException implements Exception {
  final String message;
  final String? code;
  final Map<String, dynamic>? details;

  const AuthException(this.message, {this.code, this.details});

  @override
  String toString() => message;
}

/// Exception thrown when an email is already registered
class EmailAlreadyExistsException extends AuthException {
  final String email;
  
  EmailAlreadyExistsException(this.email)
      : super(
          'This email is already being used in our system',
          code: 'EMAIL_ALREADY_EXISTS',
          details: {'email': email},
        );
}

/// Exception thrown when login credentials are invalid
class InvalidCredentialsException extends AuthException {
  const InvalidCredentialsException()
      : super(
          'Invalid email or password. Please check your credentials and try again.',
          code: 'INVALID_CREDENTIALS',
        );
}

/// Exception thrown when account is not verified
class AccountNotVerifiedException extends AuthException {
  final String email;
  
  AccountNotVerifiedException(this.email)
      : super(
          'Your account is not verified. Please check your email for verification instructions.',
          code: 'ACCOUNT_NOT_VERIFIED',
          details: {'email': email},
        );
}

/// Exception thrown when token authentication/refresh fails
class AuthenticationException extends AuthException {
  AuthenticationException(String message)
      : super(
          message,
          code: 'AUTHENTICATION_FAILED',
        );
}
