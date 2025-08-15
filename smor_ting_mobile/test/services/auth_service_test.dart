import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:smor_ting_mobile/services/auth_service.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class MockApiService extends Mock implements ApiService {}
class MockFlutterSecureStorage extends Mock implements FlutterSecureStorage {}

void main() {
  group('AuthService - Cleaner Token Refresh Implementation', () {
    late AuthService authService;
    late MockApiService mockApiService;
    late MockFlutterSecureStorage mockSecureStorage;

    // Valid JWT token for testing (expires in year 2030)
    const validJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';

    // Expired JWT token for testing
    const expiredJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyMzkwMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';

    setUp(() {
      mockApiService = MockApiService();
      mockSecureStorage = MockFlutterSecureStorage();
      authService = AuthService(
        apiService: mockApiService,
        secureStorage: mockSecureStorage,
      );
      
      // Register fallback values
      registerFallbackValue(<String, dynamic>{});
    });

    group('Token Validation Tests', () {
      test('should return valid token when it has sufficient time remaining', () async {
        // Arrange
        when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => validJwt);

        // Act
        final result = await authService.getValidToken();

        // Assert
        expect(result, equals(validJwt));
        verifyNever(() => mockApiService.refreshToken(any(), any()));
      });

      test('should refresh token when it is expired or expiring soon', () async {
        // Arrange
        when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => expiredJwt);
        when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => validJwt);
        when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => 'test-session');
        
        when(() => mockApiService.refreshToken(any(), any())).thenAnswer((_) async => {
          'access_token': validJwt,
          'refresh_token': validJwt,
          'refresh_expires_at': DateTime.now().add(const Duration(days: 30)).toIso8601String(),
        });

        when(() => mockSecureStorage.write(key: any(named: 'key'), value: any(named: 'value')))
            .thenAnswer((_) async {});

        // Act
        final result = await authService.getValidToken();

        // Assert
        expect(result, equals(validJwt));
        verify(() => mockApiService.refreshToken(any(), any())).called(1);
      });

      test('should throw when no access token found', () async {
        // Arrange
        when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => null);

        // Act & Assert
        expect(
          () async => await authService.getValidToken(),
          throwsA(isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('No access token found'),
          )),
        );
      });
    });

    group('Concurrent Refresh Prevention Tests', () {
      test('should prevent multiple concurrent refresh requests', () async {
        // Arrange
        when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => expiredJwt);
        when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => validJwt);
        when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => 'test-session');
        
        // Mock API service to have slow response
        when(() => mockApiService.refreshToken(any(), any())).thenAnswer((_) async {
          await Future.delayed(const Duration(milliseconds: 100));
          return {
            'access_token': validJwt,
            'refresh_token': validJwt,
            'refresh_expires_at': DateTime.now().add(const Duration(days: 30)).toIso8601String(),
          };
        });

        when(() => mockSecureStorage.write(key: any(named: 'key'), value: any(named: 'value')))
            .thenAnswer((_) async {});

        // Act - Make 5 concurrent refresh requests
        final futures = List.generate(5, (_) => authService.refreshToken());
        final results = await Future.wait(futures);

        // Assert - Only one API call should be made
        verify(() => mockApiService.refreshToken(any(), any())).called(1);
        
        // All results should be successful and identical
        for (final result in results) {
          expect(result, equals(validJwt));
        }
      });

      test('should allow new refresh after previous completes', () async {
        // Arrange
        when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => validJwt);
        when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => 'test-session');

        when(() => mockApiService.refreshToken(any(), any())).thenAnswer((_) async => {
          'access_token': validJwt,
          'refresh_token': validJwt,
          'refresh_expires_at': DateTime.now().add(const Duration(days: 30)).toIso8601String(),
        });

        when(() => mockSecureStorage.write(key: any(named: 'key'), value: any(named: 'value')))
            .thenAnswer((_) async {});

        // Act - Make sequential refresh requests
        final result1 = await authService.refreshToken();
        final result2 = await authService.refreshToken();

        // Assert - Two separate API calls should be made
        verify(() => mockApiService.refreshToken(any(), any())).called(2);
        expect(result1, equals(validJwt));
        expect(result2, equals(validJwt));
      });
    });

    group('Error Handling Tests', () {
      test('should handle refresh token API failures', () async {
        // Arrange
        when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => validJwt);
        when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => 'test-session');

        when(() => mockApiService.refreshToken(any(), any())).thenThrow(Exception('Refresh failed'));

        // Act & Assert
        expect(
          () async => await authService.refreshToken(),
          throwsA(isA<Exception>()),
        );

        // The verify needs to be outside the expect since the exception prevents normal call tracking
      });

      test('should handle missing refresh token', () async {
        // Arrange
        when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => null);
        when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => null);

        // Act & Assert
        expect(
          () async => await authService.refreshToken(),
          throwsA(isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('No refresh token found'),
          )),
        );
      });
    });
  });
}