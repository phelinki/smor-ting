import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mocktail/mocktail.dart';
import 'package:local_auth/local_auth.dart';

import 'package:smor_ting_mobile/core/services/enhanced_auth_service.dart';
import 'package:smor_ting_mobile/features/settings/presentation/pages/settings_page.dart';

// Mock classes
class MockEnhancedAuthService extends Mock implements EnhancedAuthService {}

void main() {
  group('Biometric Settings Tests', () {
    late MockEnhancedAuthService mockAuthService;

    setUp(() {
      mockAuthService = MockEnhancedAuthService();
    });

    Widget createTestWidget({bool biometricAvailable = true, bool biometricEnabled = false}) {
      return ProviderScope(
        overrides: [
          enhancedAuthServiceProvider.overrideWithValue(mockAuthService),
        ],
        child: MaterialApp(
          home: const SettingsPage(),
        ),
      );
    }

    group('Biometric Toggle Visibility', () {
      testWidgets('shows biometric toggle when biometrics are available', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => false);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        // Assert
        expect(find.text('Biometric Authentication'), findsOneWidget);
        expect(find.text('Use fingerprint or face unlock to secure your account'), findsOneWidget);
      });

      testWidgets('hides biometric toggle when biometrics are not available', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => false);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => <BiometricType>[]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => false);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        // Assert
        expect(find.text('Biometric Authentication'), findsNothing);
      });

      testWidgets('hides biometric toggle when no biometric types are available', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => <BiometricType>[]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => false);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        // Assert
        expect(find.text('Biometric Authentication'), findsNothing);
      });
    });

    group('Biometric Toggle State', () {
      testWidgets('shows toggle as enabled when biometric is enabled for user', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => true);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        // Assert
        final switchFinder = find.byType(Switch);
        expect(switchFinder, findsOneWidget);
        
        final switchWidget = tester.widget<Switch>(switchFinder);
        expect(switchWidget.value, isTrue);
      });

      testWidgets('shows toggle as disabled when biometric is disabled for user', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => false);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        // Assert
        final switchFinder = find.byType(Switch);
        expect(switchFinder, findsOneWidget);
        
        final switchWidget = tester.widget<Switch>(switchFinder);
        expect(switchWidget.value, isFalse);
      });
    });

    group('Biometric Toggle Actions', () {
      testWidgets('enables biometric authentication when toggle is turned on', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => false);
        when(() => mockAuthService.setBiometricEnabled(any(), true))
            .thenAnswer((_) async => true);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        final switchFinder = find.byType(Switch);
        await tester.tap(switchFinder);
        await tester.pumpAndSettle();

        // Assert
        verify(() => mockAuthService.setBiometricEnabled(any(), true)).called(1);
        expect(find.text('Biometric authentication enabled successfully'), findsOneWidget);
      });

      testWidgets('disables biometric authentication when toggle is turned off', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => true);
        when(() => mockAuthService.setBiometricEnabled(any(), false))
            .thenAnswer((_) async => true);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        final switchFinder = find.byType(Switch);
        await tester.tap(switchFinder);
        await tester.pumpAndSettle();

        // Assert
        verify(() => mockAuthService.setBiometricEnabled(any(), false)).called(1);
        expect(find.text('Biometric authentication disabled successfully'), findsOneWidget);
      });

      testWidgets('shows error message when enabling biometric authentication fails', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => false);
        when(() => mockAuthService.setBiometricEnabled(any(), true))
            .thenAnswer((_) async => false);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        final switchFinder = find.byType(Switch);
        await tester.tap(switchFinder);
        await tester.pumpAndSettle();

        // Assert
        verify(() => mockAuthService.setBiometricEnabled(any(), true)).called(1);
        expect(find.text('Failed to enable biometric authentication. Please try again.'), findsOneWidget);
      });

      testWidgets('shows error message when disabling biometric authentication fails', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => true);
        when(() => mockAuthService.setBiometricEnabled(any(), false))
            .thenAnswer((_) async => false);

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        final switchFinder = find.byType(Switch);
        await tester.tap(switchFinder);
        await tester.pumpAndSettle();

        // Assert
        verify(() => mockAuthService.setBiometricEnabled(any(), false)).called(1);
        expect(find.text('Failed to disable biometric authentication. Please try again.'), findsOneWidget);
      });

      testWidgets('shows error message when biometric service throws exception', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenAnswer((_) async => false);
        when(() => mockAuthService.setBiometricEnabled(any(), true))
            .thenThrow(Exception('Biometric service error'));

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        final switchFinder = find.byType(Switch);
        await tester.tap(switchFinder);
        await tester.pumpAndSettle();

        // Assert
        verify(() => mockAuthService.setBiometricEnabled(any(), true)).called(1);
        expect(find.text('Error: Exception: Biometric service error'), findsOneWidget);
      });
    });

    group('Initialization', () {
      testWidgets('handles errors gracefully during biometric availability check', (tester) async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenThrow(Exception('Platform error'));
        when(() => mockAuthService.getAvailableBiometrics())
            .thenThrow(Exception('Platform error'));
        when(() => mockAuthService.isBiometricEnabled(any()))
            .thenThrow(Exception('Platform error'));

        // Act
        await tester.pumpWidget(createTestWidget());
        await tester.pumpAndSettle();

        // Assert - should not crash and should not show biometric toggle
        expect(find.text('Biometric Authentication'), findsNothing);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
