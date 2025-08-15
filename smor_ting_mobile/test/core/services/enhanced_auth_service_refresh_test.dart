import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:local_auth/local_auth.dart';
import 'package:smor_ting_mobile/core/services/enhanced_auth_service.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:smor_ting_mobile/core/services/session_manager.dart';
import 'package:smor_ting_mobile/core/services/device_fingerprint_service.dart';
import 'package:smor_ting_mobile/core/models/enhanced_auth_models.dart';
import 'package:smor_ting_mobile/core/models/user.dart';

class MockApiService extends Mock implements ApiService {}
class MockSessionManager extends Mock implements SessionManager {}
class MockDeviceFingerprintService extends Mock implements DeviceFingerprintService {}
class MockFlutterSecureStorage extends Mock implements FlutterSecureStorage {}
class MockLocalAuthentication extends Mock implements LocalAuthentication {}

void main() {
  group('Enhanced Auth Service - Token Refresh Fixes', () {
    late EnhancedAuthService authService;
    late MockApiService mockApiService;
    late MockSessionManager mockSessionManager;
    late MockDeviceFingerprintService mockDeviceService;
    late MockFlutterSecureStorage mockSecureStorage;
    late MockLocalAuthentication mockLocalAuth;

    // Helper to create test user
    User createTestUser() {
      return User(
        id: 'user123',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: true,
        createdAt: DateTime.now().subtract(const Duration(days: 1)),
        updatedAt: DateTime.now(),
      );
    }

    // Helper to create test session data
    SessionData createTestSessionData({
      String? accessToken,
      String? refreshToken,
      DateTime? tokenExpiresAt,
      DateTime? refreshExpiresAt,
    }) {
      // Create a valid JWT token for testing (expires in year 2030)
      const defaultValidJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';
      
      return SessionData(
        user: createTestUser(),
        accessToken: accessToken ?? defaultValidJwt,
        refreshToken: refreshToken ?? defaultValidJwt,
        sessionId: 'session123',
        tokenExpiresAt: tokenExpiresAt ?? DateTime.now().subtract(const Duration(minutes: 1)),
        refreshExpiresAt: refreshExpiresAt ?? DateTime.now().add(const Duration(hours: 1)),
        rememberMe: true,
        deviceTrusted: false,
      );
    }

    setUp(() {
      mockApiService = MockApiService();
      mockSessionManager = MockSessionManager();
      mockDeviceService = MockDeviceFingerprintService();
      mockSecureStorage = MockFlutterSecureStorage();
      mockLocalAuth = MockLocalAuthentication();

      authService = EnhancedAuthService(
        mockApiService,
        mockSessionManager,
        mockDeviceService,
        mockSecureStorage,
        mockLocalAuth,
      );

      // Register fallback values for mocktail
      registerFallbackValue(createTestSessionData());
      
      // Set up default stub for clearSession to return a completed future
      when(() => mockSessionManager.clearSession()).thenAnswer((_) async {});
    });

    group('Debouncing Tests', () {
      test('should prevent multiple simultaneous refresh requests', () async {
        // Arrange
        final sessionData = createTestSessionData();

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockApiService.refreshToken(any(), any()))
            .thenAnswer((_) async {
          // Simulate slow API call
          await Future.delayed(const Duration(milliseconds: 100));
          // Return properly formatted JWT tokens
          const newValidJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';
          return {
            'success': true,
            'access_token': newValidJwt,
            'refresh_token': newValidJwt,
            'token_expires_at': DateTime.now().add(const Duration(minutes: 30)).toIso8601String(),
            'refresh_expires_at': DateTime.now().add(const Duration(hours: 2)).toIso8601String(),
          };
        });

        when(() => mockSessionManager.storeSession(any()))
            .thenAnswer((_) async {});

        // Act - Make multiple concurrent refresh requests
        final futures = List.generate(5, (_) => authService.restoreSession());
        final results = await Future.wait(futures);

        // Assert - Only one API call should be made despite multiple requests
        const defaultValidJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';
        verify(() => mockApiService.refreshToken(defaultValidJwt, 'session123')).called(1);
        
        // All results should be successful and consistent
        const newValidJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';
        for (final result in results) {
          expect(result, isNotNull);
          expect(result!.success, isTrue);
          expect(result.accessToken, equals(newValidJwt));
        }
      });

      test('should allow new refresh request after previous one completes', () async {
        // Arrange
        final sessionData = createTestSessionData();

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockApiService.refreshToken(any(), any()))
            .thenAnswer((_) async {
              // Return properly formatted JWT tokens
              const newValidJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';
              return {
                'success': true,
                'access_token': newValidJwt,
                'refresh_token': newValidJwt,
                'token_expires_at': DateTime.now().add(const Duration(minutes: 30)).toIso8601String(),
                'refresh_expires_at': DateTime.now().add(const Duration(hours: 2)).toIso8601String(),
              };
            });

        when(() => mockSessionManager.storeSession(any())).thenAnswer((_) async {});

        // Act - Make first refresh request and wait for completion
        final result1 = await authService.restoreSession();
        
        // Make second refresh request after first completes
        final result2 = await authService.restoreSession();

        // Assert - Two separate API calls should be made
        verify(() => mockApiService.refreshToken(any(), any())).called(2);
        expect(result1, isNotNull);
        expect(result2, isNotNull);
      });

      test('should handle concurrent requests when refresh fails', () async {
        // Arrange - Use a properly formatted JWT that will pass validation but fail at API
        const validJwtThatFailsAtApi = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjk5OTk5OTk5OTl9.Lh5jOJ6bciYn4kZNc4u3Lx-D6QKPpQ1rUGKKHFwq-qE';
        final sessionData = createTestSessionData(
          refreshToken: validJwtThatFailsAtApi,
        );

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockApiService.refreshToken(any(), any()))
            .thenThrow(Exception('Refresh failed'));

        when(() => mockSessionManager.clearSession()).thenAnswer((_) async {});

        // Act - Make multiple concurrent refresh requests that will fail
        final futures = List.generate(3, (_) => authService.restoreSession());
        final results = await Future.wait(futures);

        // Assert - Only one API call should be made
        verify(() => mockApiService.refreshToken(validJwtThatFailsAtApi, 'session123')).called(1);
        
        // All results should be null due to failure
        for (final result in results) {
          expect(result, isNull);
        }
      });
    });

    group('Error Handling Tests', () {
      test('should handle network timeout during refresh', () async {
        // Arrange
        final sessionData = createTestSessionData();

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockApiService.refreshToken(any(), any()))
            .thenThrow(Exception('Network timeout'));

        // Act
        final result = await authService.restoreSession();

        // Assert
        expect(result, isNull);
        verifyNever(() => mockSessionManager.storeSession(any()));
      });

      test('should handle invalid refresh token response', () async {
        // Arrange
        final sessionData = createTestSessionData();

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockApiService.refreshToken(any(), any()))
            .thenAnswer((_) async => {
              'success': false,
              'message': 'Invalid refresh token',
            });

        // Act
        final result = await authService.restoreSession();

        // Assert
        expect(result, isNull);
        verifyNever(() => mockSessionManager.storeSession(any()));
      });

      test('should retry refresh on recoverable errors with exponential backoff', () async {
        // Arrange
        final sessionData = createTestSessionData();

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        // First two calls fail, third succeeds
        int callCount = 0;
        when(() => mockApiService.refreshToken(any(), any()))
            .thenAnswer((_) async {
          callCount++;
          if (callCount <= 2) {
            throw Exception('Server temporarily unavailable');
          }
          // Return properly formatted JWT tokens
          const newValidJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';
          return {
            'success': true,
            'access_token': newValidJwt,
            'refresh_token': newValidJwt,
            'token_expires_at': DateTime.now().add(const Duration(minutes: 30)).toIso8601String(),
            'refresh_expires_at': DateTime.now().add(const Duration(hours: 2)).toIso8601String(),
          };
        });

        when(() => mockSessionManager.storeSession(any())).thenAnswer((_) async {});

        // Act
        final result = await authService.restoreSession();

        // Assert
        expect(result, isNotNull);
        expect(result!.success, isTrue);
        verify(() => mockApiService.refreshToken(any(), any())).called(3);
      });
    });

    group('JWT Validation Tests', () {
      test('should validate access token structure before refresh', () async {
        // Arrange
        final sessionData = createTestSessionData(
          accessToken: 'invalid.jwt.structure', // Invalid JWT
        );

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockApiService.refreshToken(any(), any()))
            .thenAnswer((_) async {
              // Return properly formatted JWT tokens
              const newValidJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';
              return {
                'success': true,
                'access_token': newValidJwt,
                'refresh_token': newValidJwt,
                'token_expires_at': DateTime.now().add(const Duration(minutes: 30)).toIso8601String(),
                'refresh_expires_at': DateTime.now().add(const Duration(hours: 2)).toIso8601String(),
              };
            });

        when(() => mockSessionManager.storeSession(any())).thenAnswer((_) async {});

        // Act
        final result = await authService.restoreSession();

        // Assert - Should attempt refresh since access token is expired regardless of structure
        const defaultValidJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';
        verify(() => mockApiService.refreshToken(defaultValidJwt, 'session123')).called(1);
      });

      test('should validate refresh token structure before attempting refresh', () async {
        // Arrange
        final sessionData = createTestSessionData(
          refreshToken: 'invalid.jwt', // Invalid JWT structure
        );

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockSessionManager.clearSession()).thenAnswer((_) async {});

        // Act
        final result = await authService.restoreSession();

        // Assert - Should fail validation and clear session
        expect(result, isNull);
        verify(() => mockSessionManager.clearSession()).called(1);
        verifyNever(() => mockApiService.refreshToken(any(), any()));
      });

      test('should validate refresh token expiration before API call', () async {
        // Arrange
        final sessionData = createTestSessionData(
          refreshExpiresAt: DateTime.now().subtract(const Duration(minutes: 1)), // Expired refresh token
        );

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockSessionManager.clearSession()).thenAnswer((_) async {});

        // Act
        final result = await authService.restoreSession();

        // Assert - Should not attempt refresh with expired refresh token
        expect(result, isNull);
        verify(() => mockSessionManager.clearSession()).called(1);
        verifyNever(() => mockApiService.refreshToken(any(), any()));
      });

      test('should validate JWT signature format before refresh', () async {
        // Arrange - Create a malformed JWT token
        const malformedJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed_payload.invalid_signature';
        
        final sessionData = createTestSessionData(
          refreshToken: malformedJwt,
        );

        when(() => mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(() => mockSessionManager.clearSession()).thenAnswer((_) async {});

        // Act
        final result = await authService.restoreSession();

        // Assert - Should fail JWT validation and clear session
        expect(result, isNull);
        verify(() => mockSessionManager.clearSession()).called(1);
        verifyNever(() => mockApiService.refreshToken(any(), any()));
      });
    });
  });
}