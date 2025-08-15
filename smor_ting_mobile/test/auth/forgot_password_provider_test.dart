import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mocktail/mocktail.dart';

import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';

class MockApiService extends Mock implements ApiService {}

void main() {
  late MockApiService mockApiService;

  setUp(() {
    mockApiService = MockApiService();
  });

  test('requestPasswordReset sets state to PasswordResetEmailSent', () async {
    // Arrange
    when(() => mockApiService.requestPasswordReset(any()))
        .thenAnswer((_) async {});

    final container = ProviderContainer(overrides: [
      apiServiceProvider.overrideWithValue(mockApiService),
    ]);

    final notifier = container.read(authNotifierProvider.notifier);

    // Act
    await notifier.requestPasswordReset('user@example.com');

    // Assert
    final state = container.read(authNotifierProvider);
    expect(state, isA<PasswordResetEmailSent>());
    verify(() => mockApiService.requestPasswordReset('user@example.com')).called(1);
  });

  test('resetPassword sets state to PasswordResetSuccess', () async {
    // Arrange
    when(() => mockApiService.resetPassword(any(), any()))
        .thenAnswer((_) async {});

    final container = ProviderContainer(overrides: [
      apiServiceProvider.overrideWithValue(mockApiService),
    ]);

    final notifier = container.read(authNotifierProvider.notifier);

    // Act
    await notifier.resetPassword('user@example.com', '123456', 'NewPass123!');

    // Assert
    final state = container.read(authNotifierProvider);
    expect(state, isA<PasswordResetSuccess>());
    verify(() => mockApiService.resetPassword('user@example.com', '123456', 'NewPass123!')).called(1);
  });

  test('handles error during password reset request', () async {
    // Arrange
    when(() => mockApiService.requestPasswordReset(any()))
        .thenThrow(Exception('Network error'));

    final container = ProviderContainer(overrides: [
      apiServiceProvider.overrideWithValue(mockApiService),
    ]);

    final notifier = container.read(authNotifierProvider.notifier);

    // Act
    await notifier.requestPasswordReset('user@example.com');

    // Assert
    final state = container.read(authNotifierProvider);
    expect(state, isA<Error>());
    expect((state as Error).message, contains('Network error'));
  });

  test('handles error during password reset', () async {
    // Arrange
    when(() => mockApiService.resetPassword(any(), any()))
        .thenThrow(Exception('Invalid OTP'));

    final container = ProviderContainer(overrides: [
      apiServiceProvider.overrideWithValue(mockApiService),
    ]);

    final notifier = container.read(authNotifierProvider.notifier);

    // Act
    await notifier.resetPassword('user@example.com', 'invalid', 'NewPass123!');

    // Assert
    final state = container.read(authNotifierProvider);
    expect(state, isA<Error>());
    expect((state as Error).message, contains('Invalid OTP'));
  });
}