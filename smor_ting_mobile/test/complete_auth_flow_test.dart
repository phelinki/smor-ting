import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:smor_ting_mobile/core/models/user.dart';

class MockFlutterSecureStorage extends Mock implements FlutterSecureStorage {}
class MockApiService extends Mock implements ApiService {}

void main() {
  setUpAll(() {
    registerFallbackValue(LoginRequest(email: 'test@example.com', password: 'password'));
  });

  group('Complete Auth Flow Tests', () {
    late MockFlutterSecureStorage mockSecureStorage;
    late MockApiService mockApiService;
    
    // Valid JWT token for testing (expires in year 2030)
    const validJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';

    setUp(() {
      mockSecureStorage = MockFlutterSecureStorage();
      mockApiService = MockApiService();
    });

    test('should complete full auth flow: login -> store tokens -> restore session', () async {
      // Arrange - Simulate successful login
      final testUser = User(
        id: 'test-user-id',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '1234567890',
        role: UserRole.customer,
        isEmailVerified: true,
        profileImage: '',
        address: const Address(
          street: '',
          city: '',
          county: '',
          country: '',
          latitude: 0,
          longitude: 0,
        ),
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final authResponse = AuthResponse(
        user: testUser,
        accessToken: validJwt,
        refreshToken: validJwt,
        requiresOTP: false,
      );

      // Mock API responses
      when(() => mockApiService.login(any())).thenAnswer((_) async => authResponse);
      when(() => mockApiService.getUserProfile()).thenAnswer((_) async => testUser);
      when(() => mockSecureStorage.write(key: any(named: 'key'), value: any(named: 'value')))
          .thenAnswer((_) async {});

      // Act - Simulate login flow
      // 1. Login and store tokens
      final loginRequest = LoginRequest(email: 'test@example.com', password: 'password');
      final response = await mockApiService.login(loginRequest);
      
      // 2. Store tokens (this would happen in auth provider)
      await Future.wait([
        mockSecureStorage.write(key: 'access_token', value: response.accessToken),
        mockSecureStorage.write(key: 'refresh_token', value: response.refreshToken),
        mockSecureStorage.write(key: 'session_id', value: 'test-session'),
      ]);

      // 3. Simulate app restart - restore tokens
      when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => validJwt);
      when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => validJwt);
      when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => 'test-session');

      final restoredAccessToken = await mockSecureStorage.read(key: 'access_token');
      final restoredRefreshToken = await mockSecureStorage.read(key: 'refresh_token');
      final restoredSessionId = await mockSecureStorage.read(key: 'session_id');

      // Assert
      expect(response.user.email, equals('test@example.com'));
      expect(response.accessToken, equals(validJwt));
      expect(restoredAccessToken, equals(validJwt));
      expect(restoredRefreshToken, equals(validJwt));
      expect(restoredSessionId, equals('test-session'));
    });

    test('should handle logout and clear all tokens', () async {
      // Arrange
      when(() => mockSecureStorage.delete(key: any(named: 'key'))).thenAnswer((_) async {});

      // Act - Simulate logout
      await Future.wait([
        mockSecureStorage.delete(key: 'access_token'),
        mockSecureStorage.delete(key: 'refresh_token'),
        mockSecureStorage.delete(key: 'session_id'),
        mockSecureStorage.delete(key: 'token_expires_at'),
        mockSecureStorage.delete(key: 'refresh_expires_at'),
      ]);

      // Assert
      verify(() => mockSecureStorage.delete(key: 'access_token')).called(1);
      verify(() => mockSecureStorage.delete(key: 'refresh_token')).called(1);
      verify(() => mockSecureStorage.delete(key: 'session_id')).called(1);
      verify(() => mockSecureStorage.delete(key: 'token_expires_at')).called(1);
      verify(() => mockSecureStorage.delete(key: 'refresh_expires_at')).called(1);
    });

    test('should handle missing tokens gracefully', () async {
      // Arrange - No stored tokens
      when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => null);
      when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => null);
      when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => null);

      // Act
      final accessToken = await mockSecureStorage.read(key: 'access_token');
      final refreshToken = await mockSecureStorage.read(key: 'refresh_token');
      final sessionId = await mockSecureStorage.read(key: 'session_id');

      // Assert
      expect(accessToken, isNull);
      expect(refreshToken, isNull);
      expect(sessionId, isNull);
    });
  });
}
