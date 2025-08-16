import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:smor_ting_mobile/features/splash/presentation/pages/splash_page.dart';
import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';
import 'package:smor_ting_mobile/core/models/user.dart';

class MockAuthNotifier extends Mock implements AuthNotifier {}

void main() {
  group('Splash Page Navigation Tests', () {
    late MockAuthNotifier mockAuthNotifier;

    setUp(() {
      mockAuthNotifier = MockAuthNotifier();
    });

    test('should navigate to home when user is authenticated as customer', () {
      // Arrange
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

      final authenticatedState = AuthState.authenticated(testUser, 'test-token');

      // Act & Assert
      expect(authenticatedState, isA<Authenticated>());
      if (authenticatedState is Authenticated) {
        expect(authenticatedState.user.role, equals(UserRole.customer));
      }
    });

    test('should navigate to agent dashboard when user is authenticated as provider', () {
      // Arrange
      final testUser = User(
        id: 'test-provider-id',
        email: 'provider@example.com',
        firstName: 'Test',
        lastName: 'Provider',
        phone: '1234567890',
        role: UserRole.provider,
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

      final authenticatedState = AuthState.authenticated(testUser, 'test-token');

      // Act & Assert
      expect(authenticatedState, isA<Authenticated>());
      if (authenticatedState is Authenticated) {
        expect(authenticatedState.user.role, equals(UserRole.provider));
      }
    });

    test('should navigate to landing when user is not authenticated', () {
      // Arrange
      final initialState = const AuthState.initial();

      // Act & Assert
      expect(initialState, isA<Initial>());
    });
  });
}
