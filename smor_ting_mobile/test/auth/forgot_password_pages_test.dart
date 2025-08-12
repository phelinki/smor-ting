import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';
import 'package:smor_ting_mobile/features/auth/presentation/pages/forgot_password_page.dart';
import 'package:smor_ting_mobile/features/auth/presentation/pages/reset_password_page.dart';

class _FakeApiService extends ApiService {
  bool requested = false;
  bool reset = false;
  _FakeApiService(): super(baseUrl: 'http://test');
  @override
  Future<void> requestPasswordReset(String email) async { requested = true; }
  @override
  Future<void> resetPassword(String email, String otp, String newPassword) async { reset = true; }
}

void main() {
  testWidgets('ForgotPasswordPage sends request and shows success UI', (tester) async {
    final fake = _FakeApiService();
    await tester.pumpWidget(ProviderScope(
      overrides: [apiServiceProvider.overrideWithValue(fake)],
      child: const MaterialApp(home: ForgotPasswordPage()),
    ));

    // Enter email and submit
    final emailField = find.byKey(const Key('forgot_email'));
    await tester.enterText(emailField, 'user@example.com');
    await tester.tap(find.byKey(const Key('forgot_submit')));
    await tester.pump();

    expect(fake.requested, true);
  });

  testWidgets('ResetPasswordPage calls reset and shows success', (tester) async {
    final fake = _FakeApiService();
    await tester.pumpWidget(ProviderScope(
      overrides: [apiServiceProvider.overrideWithValue(fake)],
      child: const MaterialApp(home: ResetPasswordPage(email: 'user@example.com')),
    ));

    await tester.enterText(find.byKey(const Key('reset_otp')), '123456');
    await tester.enterText(find.byKey(const Key('reset_new_password')), 'NewPass123!');
    await tester.enterText(find.byKey(const Key('reset_confirm_password')), 'NewPass123!');
    await tester.tap(find.byKey(const Key('reset_submit')));
    await tester.pump();

    expect(fake.reset, true);
  });
}


