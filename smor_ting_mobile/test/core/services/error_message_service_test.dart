import 'package:flutter_test/flutter_test.dart';

import '../../../lib/core/services/error_message_service.dart';

void main() {
  group('ErrorMessageService', () {
    group('getErrorMessage', () {
      test('should return mapped message for known error code', () {
        // Act
        final message = ErrorMessageService.getErrorMessage('AUTH_001');

        // Assert
        expect(message, 'Invalid email or password. Please check your credentials and try again.');
      });

      test('should return fallback message for unknown error code', () {
        // Act
        final message = ErrorMessageService.getErrorMessage('UNKNOWN_001', fallbackMessage: 'Custom fallback');

        // Assert
        expect(message, 'Custom fallback');
      });

      test('should return generic message for unknown code with known prefix', () {
        // Act
        final message = ErrorMessageService.getErrorMessage('AUTH_999');

        // Assert
        expect(message, 'Authentication error. Please try logging in again.');
      });

      test('should return server error message for null/empty code', () {
        // Act
        final message1 = ErrorMessageService.getErrorMessage(null);
        final message2 = ErrorMessageService.getErrorMessage('');

        // Assert
        expect(message1, 'Something went wrong. Please try again later.');
        expect(message2, 'Something went wrong. Please try again later.');
      });
    });

    group('getActionMessage', () {
      test('should return action message for codes with actions', () {
        // Act
        final message = ErrorMessageService.getActionMessage('AUTH_001');

        // Assert
        expect(message, 'Try resetting your password if you\'ve forgotten it.');
      });

      test('should return null for codes without actions', () {
        // Act
        final message = ErrorMessageService.getActionMessage('AUTH_014');

        // Assert
        expect(message, null);
      });
    });

    group('isRetryable', () {
      test('should return false for non-retryable error codes', () {
        // Act & Assert
        expect(ErrorMessageService.isRetryable('AUTH_003'), false); // Email exists
        expect(ErrorMessageService.isRetryable('USER_007'), false); // Account suspended
        expect(ErrorMessageService.isRetryable('PERM_001'), false); // No permission
      });

      test('should return true for retryable error codes', () {
        // Act & Assert
        expect(ErrorMessageService.isRetryable('AUTH_001'), true); // Invalid credentials
        expect(ErrorMessageService.isRetryable('NET_001'), true); // Network error
        expect(ErrorMessageService.isRetryable('PAY_001'), true); // Payment failed
      });

      test('should return true for null error code', () {
        // Act
        final result = ErrorMessageService.isRetryable(null);

        // Assert
        expect(result, true);
      });
    });

    group('getSeverity', () {
      test('should return appropriate severity for different error types', () {
        // Act & Assert
        expect(ErrorMessageService.getSeverity('AUTH_004'), ErrorSeverity.warning); // Not verified
        expect(ErrorMessageService.getSeverity('AUTH_005'), ErrorSeverity.info); // Session expired
        expect(ErrorMessageService.getSeverity('AUTH_001'), ErrorSeverity.error); // Invalid credentials
        expect(ErrorMessageService.getSeverity('PERM_001'), ErrorSeverity.error); // No permission
      });

      test('should return error severity for null code', () {
        // Act
        final severity = ErrorMessageService.getSeverity(null);

        // Assert
        expect(severity, ErrorSeverity.error);
      });
    });

    group('error code categories', () {
      test('should handle authentication errors correctly', () {
        // Act
        final message = ErrorMessageService.getErrorMessage('AUTH_001');
        final severity = ErrorMessageService.getSeverity('AUTH_001');
        final isRetryable = ErrorMessageService.isRetryable('AUTH_001');

        // Assert
        expect(message, contains('Invalid email or password'));
        expect(severity, ErrorSeverity.error);
        expect(isRetryable, true);
      });

      test('should handle permission errors correctly', () {
        // Act
        final message = ErrorMessageService.getErrorMessage('PERM_004');
        final severity = ErrorMessageService.getSeverity('PERM_004');
        final isRetryable = ErrorMessageService.isRetryable('PERM_004');

        // Assert
        expect(message, contains('Email verification required'));
        expect(severity, ErrorSeverity.warning);
        expect(isRetryable, true);
      });

      test('should handle network errors correctly', () {
        // Act
        final message = ErrorMessageService.getErrorMessage('NET_001');
        final severity = ErrorMessageService.getSeverity('NET_001');
        final isRetryable = ErrorMessageService.isRetryable('NET_001');

        // Assert
        expect(message, contains('Network connection failed'));
        expect(severity, ErrorSeverity.error);
        expect(isRetryable, true);
      });
    });
  });
}
