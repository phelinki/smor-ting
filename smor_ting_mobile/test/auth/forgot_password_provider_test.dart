import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:smor_ting_mobile/lib/features/auth/presentation/providers/auth_provider.dart';
import 'package:smor_ting_mobile/lib/core/services/api_service.dart';

class _FakeApiService extends ApiService {
  bool requested = false;
  bool reset = false;

  _FakeApiService(): super(baseUrl: 'http://test');

  @override
  Future<void> requestPasswordReset(String email) async {
    requested = true;
  }

  @override
  Future<void> resetPassword(String email, String otp, String newPassword) async {
    reset = true;
  }
}

void main() {
  test('requestPasswordReset sets state to PasswordResetEmailSent', () async {
    final fake = _FakeApiService();
    final container = ProviderContainer(overrides: [
      apiServiceProvider.overrideWithValue(fake),
    ]);

    final notifier = container.read(authNotifierProvider.notifier);
    await notifier.requestPasswordReset('user@example.com');

    final state = container.read(authNotifierProvider);
    expect(fake.requested, true);
    expect(state, isA<PasswordResetEmailSent>());
  });

  test('resetPassword sets state to PasswordResetSuccess', () async {
    final fake = _FakeApiService();
    final container = ProviderContainer(overrides: [
      apiServiceProvider.overrideWithValue(fake),
    ]);

    final notifier = container.read(authNotifierProvider.notifier);
    await notifier.resetPassword('user@example.com', '123456', 'NewPass123!');

    final state = container.read(authNotifierProvider);
    expect(fake.reset, true);
    expect(state, isA<PasswordResetSuccess>());
  });
}


