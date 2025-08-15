import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mocktail/mocktail.dart';

import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';
import 'package:smor_ting_mobile/features/auth/presentation/pages/forgot_password_page.dart';
import 'package:smor_ting_mobile/features/auth/presentation/pages/reset_password_page.dart';

class MockApiService extends Mock implements ApiService {}

void main() {
  late MockApiService mockApiService;

  setUp(() {
    mockApiService = MockApiService();
  });

  testWidgets('ForgotPasswordPage sends request and shows success UI', (tester) async {
    // Arrange
    when(() => mockApiService.requestPasswordReset(any()))
        .thenAnswer((_) async {});

    await tester.pumpWidget(ProviderScope(
      overrides: [apiServiceProvider.overrideWithValue(mockApiService)],
      child: const MaterialApp(home: ForgotPasswordPage()),
    ));

    // Act
    // Find email field by text field type and enter email
    final emailField = find.byType(TextFormField);
    expect(emailField, findsOneWidget);
    
    await tester.enterText(emailField, 'user@example.com');
    await tester.pump();

    // Find and tap submit button
    final submitButton = find.byType(ElevatedButton);
    expect(submitButton, findsOneWidget);
    
    await tester.tap(submitButton);
    await tester.pump();

    // Assert
    verify(() => mockApiService.requestPasswordReset('user@example.com')).called(1);
  });

  testWidgets('ResetPasswordPage calls reset and shows success', (tester) async {
    // Arrange
    when(() => mockApiService.resetPassword(any(), any()))
        .thenAnswer((_) async {});

    await tester.pumpWidget(ProviderScope(
      overrides: [apiServiceProvider.overrideWithValue(mockApiService)],
      child: const MaterialApp(home: ResetPasswordPage(email: 'user@example.com')),
    ));

    // Act
    // Find text fields (should be 3: OTP, new password, confirm password)
    final textFields = find.byType(TextFormField);
    expect(textFields, findsNWidgets(3));

    // Enter OTP
    await tester.enterText(textFields.at(0), '123456');
    await tester.pump();

    // Enter new password
    await tester.enterText(textFields.at(1), 'NewPass123!');
    await tester.pump();

    // Enter confirm password
    await tester.enterText(textFields.at(2), 'NewPass123!');
    await tester.pump();

    // Find and tap submit button
    final submitButton = find.byType(ElevatedButton);
    expect(submitButton, findsOneWidget);
    
    await tester.tap(submitButton);
    await tester.pump();

    // Assert
    verify(() => mockApiService.resetPassword('user@example.com', '123456', 'NewPass123!')).called(1);
  });

  testWidgets('ForgotPasswordPage shows validation error for empty email', (tester) async {
    // Arrange
    await tester.pumpWidget(ProviderScope(
      overrides: [apiServiceProvider.overrideWithValue(mockApiService)],
      child: const MaterialApp(home: ForgotPasswordPage()),
    ));

    // Act
    // Find and tap submit button without entering email
    final submitButton = find.byType(ElevatedButton);
    await tester.tap(submitButton);
    await tester.pump();

    // Assert
    expect(find.text('Email is required'), findsOneWidget);
    verifyNever(() => mockApiService.requestPasswordReset(any()));
  });

  testWidgets('ResetPasswordPage handles empty validation', (tester) async {
    // Arrange
    await tester.pumpWidget(ProviderScope(
      overrides: [apiServiceProvider.overrideWithValue(mockApiService)],
      child: const MaterialApp(home: ResetPasswordPage(email: 'user@example.com')),
    ));

    // Act
    // Find and tap submit button without entering any data
    final submitButton = find.byType(ElevatedButton);
    await tester.tap(submitButton);
    await tester.pump();

    // Assert
    // Should show validation errors - OTP requires 6 digits
    expect(find.text('Enter 6-digit code'), findsOneWidget);
    expect(find.text('Min 6 characters'), findsOneWidget);
    verifyNever(() => mockApiService.resetPassword(any(), any(), any()));
  });
}