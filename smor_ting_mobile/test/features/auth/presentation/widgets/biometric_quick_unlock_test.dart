import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:go_router/go_router.dart';

import '../../../../../lib/features/auth/presentation/widgets/biometric_quick_unlock.dart';
import '../../../../../lib/features/auth/presentation/providers/biometric_auth_provider.dart';
import '../../../../../lib/core/services/enhanced_auth_service.dart';
import '../../../../../lib/core/models/user.dart';

import 'biometric_quick_unlock_test.mocks.dart';

@GenerateMocks([EnhancedAuthService, GoRouter])
void main() {
  group('BiometricQuickUnlock', () {
    late MockEnhancedAuthService mockAuthService;
    late MockGoRouter mockRouter;
    late ProviderContainer container;

    setUp(() {
      mockAuthService = MockEnhancedAuthService();
      mockRouter = MockGoRouter();
      container = ProviderContainer(
        overrides: [
          enhancedAuthServiceProvider.overrideWithValue(mockAuthService),
        ],
      );
    });

    tearDown(() {
      container.dispose();
    });

    Widget createWidget({String? userEmail}) {
      return ProviderScope(
        parent: container,
        child: MaterialApp(
          home: Scaffold(
            body: BiometricQuickUnlock(
              userEmail: userEmail ?? 'test@example.com',
              onSuccess: () => mockRouter.go('/home'),
              onCancel: () => mockRouter.go('/login'),
            ),
          ),
        ),
      );
    }

    testWidgets('should show biometric unlock button when available and enabled', (tester) async {
      // Arrange
      when(mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => true);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      // Assert
      expect(find.text('Quick Unlock'), findsOneWidget);
      expect(find.byIcon(Icons.fingerprint), findsOneWidget);
      expect(find.text('Use biometric authentication to unlock'), findsOneWidget);
    });

    testWidgets('should not show biometric unlock when not available', (tester) async {
      // Arrange
      when(mockAuthService.canUseBiometrics()).thenAnswer((_) async => false);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      // Assert
      expect(find.text('Quick Unlock'), findsNothing);
      expect(find.byIcon(Icons.fingerprint), findsNothing);
    });

    testWidgets('should not show biometric unlock when not enabled for user', (tester) async {
      // Arrange
      when(mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => false);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      // Assert
      expect(find.text('Quick Unlock'), findsNothing);
    });

    testWidgets('should trigger biometric authentication when unlock button is tapped', (tester) async {
      // Arrange
      when(mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => true);
      
      final mockUser = User(
        id: '1',
        email: 'test@example.com',
        firstName: 'John',
        lastName: 'Doe',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: true,
        isPhoneVerified: true,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );
      
      final mockResult = EnhancedAuthResult(
        success: true,
        user: mockUser,
        accessToken: 'token123',
        sessionId: 'session123',
        deviceTrusted: true,
      );
      
      when(mockAuthService.authenticateWithBiometrics('test@example.com'))
          .thenAnswer((_) async => mockResult);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.fingerprint));
      await tester.pumpAndSettle();

      // Assert
      verify(mockAuthService.authenticateWithBiometrics('test@example.com')).called(1);
      verify(mockRouter.go('/home')).called(1);
    });

    testWidgets('should show error message when biometric authentication fails', (tester) async {
      // Arrange
      when(mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => true);
      
      final mockResult = EnhancedAuthResult(
        success: false,
        message: 'Biometric authentication failed',
      );
      
      when(mockAuthService.authenticateWithBiometrics('test@example.com'))
          .thenAnswer((_) async => mockResult);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.fingerprint));
      await tester.pumpAndSettle();

      // Assert
      expect(find.byType(SnackBar), findsOneWidget);
      expect(find.text('Biometric authentication failed'), findsOneWidget);
    });

    testWidgets('should show loading state during authentication', (tester) async {
      // Arrange
      when(mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => true);
      
      when(mockAuthService.authenticateWithBiometrics('test@example.com'))
          .thenAnswer((_) => Future.delayed(const Duration(seconds: 1), () => EnhancedAuthResult(success: false)));

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.fingerprint));
      await tester.pump(); // Don't wait for completion

      // Assert
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('should show alternative login option', (tester) async {
      // Arrange
      when(mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => true);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      // Assert
      expect(find.text('Use Password Instead'), findsOneWidget);
    });

    testWidgets('should navigate to login when "Use Password Instead" is tapped', (tester) async {
      // Arrange
      when(mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => true);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      await tester.tap(find.text('Use Password Instead'));
      await tester.pumpAndSettle();

      // Assert
      verify(mockRouter.go('/login')).called(1);
    });
  });
}
