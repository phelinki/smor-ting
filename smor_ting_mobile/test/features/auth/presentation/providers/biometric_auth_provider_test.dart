import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mocktail/mocktail.dart';
import 'package:local_auth/local_auth.dart';

import 'package:smor_ting_mobile/core/services/enhanced_auth_service.dart';
import 'package:smor_ting_mobile/core/models/user.dart';
import 'package:smor_ting_mobile/core/models/enhanced_auth_models.dart' as models;
import 'package:smor_ting_mobile/features/auth/presentation/providers/biometric_auth_provider.dart';

// Mock classes
class MockEnhancedAuthService extends Mock implements EnhancedAuthService {}
class MockLocalAuthentication extends Mock implements LocalAuthentication {}

void main() {
  group('BiometricAuthProvider', () {
    late MockEnhancedAuthService mockAuthService;
    late ProviderContainer container;

    setUp(() {
      mockAuthService = MockEnhancedAuthService();
      
      container = ProviderContainer(
        overrides: [
          enhancedAuthServiceProvider.overrideWithValue(mockAuthService),
        ],
      );
    });

    tearDown(() {
      container.dispose();
    });

    group('canUseBiometrics', () {
      test('returns true when biometrics are available', () async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.canUseBiometrics();

        // Assert
        expect(result, isTrue);
        verify(() => mockAuthService.canUseBiometrics()).called(1);
      });

      test('returns false when biometrics are not available', () async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => false);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.canUseBiometrics();

        // Assert
        expect(result, isFalse);
        verify(() => mockAuthService.canUseBiometrics()).called(1);
      });
    });

    group('getAvailableBiometrics', () {
      test('returns available biometric types', () async {
        // Arrange
        const expectedBiometrics = [BiometricType.face, BiometricType.fingerprint];
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => expectedBiometrics);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.getAvailableBiometrics();

        // Assert
        expect(result, equals(expectedBiometrics));
        verify(() => mockAuthService.getAvailableBiometrics()).called(1);
      });

      test('returns empty list when no biometrics available', () async {
        // Arrange
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => <BiometricType>[]);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.getAvailableBiometrics();

        // Assert
        expect(result, isEmpty);
        verify(() => mockAuthService.getAvailableBiometrics()).called(1);
      });
    });

    group('isBiometricEnabled', () {
      test('returns true when biometric is enabled for user', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.isBiometricEnabled(email))
            .thenAnswer((_) async => true);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.isBiometricEnabled(email);

        // Assert
        expect(result, isTrue);
        verify(() => mockAuthService.isBiometricEnabled(email)).called(1);
      });

      test('returns false when biometric is disabled for user', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.isBiometricEnabled(email))
            .thenAnswer((_) async => false);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.isBiometricEnabled(email);

        // Assert
        expect(result, isFalse);
        verify(() => mockAuthService.isBiometricEnabled(email)).called(1);
      });
    });

    group('enableBiometric', () {
      test('successfully enables biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.setBiometricEnabled(email, true))
            .thenAnswer((_) async => true);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.enableBiometric(email);

        // Assert
        expect(result, isTrue);
        verify(() => mockAuthService.setBiometricEnabled(email, true)).called(1);
      });

      test('fails to enable biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.setBiometricEnabled(email, true))
            .thenAnswer((_) async => false);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.enableBiometric(email);

        // Assert
        expect(result, isFalse);
        verify(() => mockAuthService.setBiometricEnabled(email, true)).called(1);
      });
    });

    group('disableBiometric', () {
      test('successfully disables biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.setBiometricEnabled(email, false))
            .thenAnswer((_) async => true);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.disableBiometric(email);

        // Assert
        expect(result, isTrue);
        verify(() => mockAuthService.setBiometricEnabled(email, false)).called(1);
      });

      test('fails to disable biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.setBiometricEnabled(email, false))
            .thenAnswer((_) async => false);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        final result = await provider.disableBiometric(email);

        // Assert
        expect(result, isFalse);
        verify(() => mockAuthService.setBiometricEnabled(email, false)).called(1);
      });
    });

    group('authenticateWithBiometrics', () {
      test('successfully authenticates with biometrics', () async {
        // Arrange
        const email = 'test@example.com';
        final mockUser = User(
          id: '123',
          email: email,
          firstName: 'Test',
          lastName: 'User',
          phone: '+123456789',
          role: UserRole.customer,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final mockResult = models.EnhancedAuthResult(
          success: true,
          user: mockUser,
          accessToken: 'test_token',
          sessionId: 'session_123',
          deviceTrusted: true,
          isRestoredSession: true,
        );
        
        when(() => mockAuthService.authenticateWithBiometrics(email))
            .thenAnswer((_) async => mockResult);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        await provider.authenticateWithBiometrics(email);

        // Assert
        final state = container.read(biometricAuthProvider);
        expect(state, isA<BiometricAuthSuccess>());
        final successState = state as BiometricAuthSuccess;
        expect(successState.user.email, equals(email));
        expect(successState.accessToken, equals('test_token'));
        verify(() => mockAuthService.authenticateWithBiometrics(email)).called(1);
      });

      test('fails to authenticate with biometrics', () async {
        // Arrange
        const email = 'test@example.com';
        final mockResult = models.EnhancedAuthResult(
          success: false,
          message: 'Biometric authentication failed',
        );
        
        when(() => mockAuthService.authenticateWithBiometrics(email))
            .thenAnswer((_) async => mockResult);

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        await provider.authenticateWithBiometrics(email);

        // Assert
        final state = container.read(biometricAuthProvider);
        expect(state, isA<BiometricAuthError>());
        final errorState = state as BiometricAuthError;
        expect(errorState.message, equals('Biometric authentication failed'));
        verify(() => mockAuthService.authenticateWithBiometrics(email)).called(1);
      });

      test('handles exception during biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.authenticateWithBiometrics(email))
            .thenThrow(Exception('Network error'));

        // Act
        final provider = container.read(biometricAuthProvider.notifier);
        await provider.authenticateWithBiometrics(email);

        // Assert
        final state = container.read(biometricAuthProvider);
        expect(state, isA<BiometricAuthError>());
        final errorState = state as BiometricAuthError;
        expect(errorState.message, contains('Network error'));
        verify(() => mockAuthService.authenticateWithBiometrics(email)).called(1);
      });
    });
  });
}
