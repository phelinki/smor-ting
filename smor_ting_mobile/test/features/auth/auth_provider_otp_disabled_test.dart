import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:smor_ting_mobile/core/models/user.dart';

// Mock classes
class MockApiService extends Mock implements ApiService {}

// Fake classes for mocktail
class FakeLoginRequest extends Fake implements LoginRequest {}
class FakeRegisterRequest extends Fake implements RegisterRequest {}

void main() {
  setUpAll(() {
    registerFallbackValue(FakeLoginRequest());
    registerFallbackValue(FakeRegisterRequest());
  });
  group('AuthProvider - OTP Disabled Tests', () {
    late MockApiService mockApiService;
    late ProviderContainer container;

    setUp(() {
      mockApiService = MockApiService();
      container = ProviderContainer(
        overrides: [
          apiServiceProvider.overrideWithValue(mockApiService),
        ],
      );
    });

    tearDown(() {
      container.dispose();
    });

    test('login should never transition to RequiresOTP state', () async {
      // Arrange
      final user = User(
        id: 'test123',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: false, // Even if not verified
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final authResponse = AuthResponse(
        user: user,
        accessToken: 'access_token_123',
        refreshToken: 'refresh_token_123',
        requiresOTP: false, // OTP is disabled
      );

      when(() => mockApiService.login(any())).thenAnswer((_) async => authResponse);
      when(() => mockApiService.setAuthToken(any())).thenReturn(null);

      final authNotifier = container.read(authNotifierProvider.notifier);

      // Act
      await authNotifier.login('test@example.com', 'password123');

      // Assert
      final state = container.read(authNotifierProvider);
      expect(state, isA<Authenticated>(), reason: 'Should go directly to authenticated state');
      expect(state, isNot(isA<RequiresOTP>()), reason: 'Should never require OTP');

      // Verify API calls
      verify(() => mockApiService.login(any())).called(1);
      verify(() => mockApiService.setAuthToken('access_token_123')).called(1);
    });

    test('register should never transition to RequiresOTP state', () async {
      // Arrange
      final user = User(
        id: 'test123',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: false, // Even if not verified
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final authResponse = AuthResponse(
        user: user,
        accessToken: 'access_token_123',
        refreshToken: 'refresh_token_123',
        requiresOTP: false, // OTP is disabled
      );

      when(() => mockApiService.register(any())).thenAnswer((_) async => authResponse);
      when(() => mockApiService.setAuthToken(any())).thenReturn(null);

      final authNotifier = container.read(authNotifierProvider.notifier);

      // Act
      await authNotifier.register(
        firstName: 'Test',
        lastName: 'User',
        email: 'test@example.com',
        phone: '+1234567890',
        password: 'password123',
      );

      // Assert
      final state = container.read(authNotifierProvider);
      expect(state, isA<Authenticated>(), reason: 'Should go directly to authenticated state');
      expect(state, isNot(isA<RequiresOTP>()), reason: 'Should never require OTP');

      // Verify API calls
      verify(() => mockApiService.register(any())).called(1);
      verify(() => mockApiService.setAuthToken('access_token_123')).called(1);
    });

    test('verifyOTP method should throw error when called', () async {
      // Arrange
      final authNotifier = container.read(authNotifierProvider.notifier);

      // Act & Assert
      expect(
        () async => await authNotifier.verifyOTP('test@example.com', '123456'),
        throwsA(isA<UnsupportedError>()),
        reason: 'verifyOTP should throw error since OTP is disabled',
      );
    });
  });
}
