import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:go_router/go_router.dart';

import '../../../../../lib/features/splash/presentation/pages/splash_page.dart';
import '../../../../../lib/features/auth/presentation/providers/enhanced_auth_provider.dart';
import '../../../../../lib/core/services/enhanced_auth_service.dart';
import '../../../../../lib/core/models/user.dart';
import '../../../../../lib/core/models/enhanced_auth_models.dart';

class MockEnhancedAuthService extends Mock implements EnhancedAuthService {}
class MockGoRouter extends Mock implements GoRouter {}
void main() {
  group('SplashPage', () {
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

    Widget createWidget() {
      return ProviderScope(
        parent: container,
        child: MaterialApp.router(
          routerConfig: GoRouter(routes: [
            GoRoute(path: '/', builder: (context, state) => const SplashPage()),
          ]),
        ),
      );
    }

    testWidgets('should show app logo and loading indicator', (tester) async {
      // Arrange
      when(() => mockAuthService.restoreSession()).thenAnswer((_) async => null);

      // Act
      await tester.pumpWidget(createWidget());

      // Assert
      expect(find.text('Smor-Ting'), findsOneWidget);
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('should show biometric quick unlock when session exists and biometric enabled', (tester) async {
      // Arrange
      final mockUser = User(
        id: '1',
        email: 'test@example.com',
        firstName: 'John',
        lastName: 'Doe',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: true,
        
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final mockAuthResult = EnhancedAuthResult(
        success: true,
        user: mockUser,
        accessToken: 'token123',
        sessionId: 'session123',
        deviceTrusted: true,
      );

      when(() => mockAuthService.restoreSession()).thenAnswer((_) async => mockAuthResult);
      when(() => mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(() => mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => true);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      // Assert
      expect(find.text('Quick Unlock'), findsOneWidget);
      expect(find.byIcon(Icons.fingerprint), findsOneWidget);
    });

    testWidgets('should navigate to home when session restored without biometric', (tester) async {
      // Arrange
      final mockUser = User(
        id: '1',
        email: 'test@example.com',
        firstName: 'John',
        lastName: 'Doe',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: true,
        
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final mockAuthResult = EnhancedAuthResult(
        success: true,
        user: mockUser,
        accessToken: 'token123',
        sessionId: 'session123',
        deviceTrusted: true,
      );

      when(() => mockAuthService.restoreSession()).thenAnswer((_) async => mockAuthResult);
      when(() => mockAuthService.canUseBiometrics()).thenAnswer((_) async => false);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      // Assert
      verify(() => mockRouter.go('/home')).called(1);
    });

    testWidgets('should navigate to login when no session exists', (tester) async {
      // Arrange
      when(() => mockAuthService.restoreSession()).thenAnswer((_) async => null);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      // Assert
      verify(() => mockRouter.go('/landing')).called(1);
    });

    testWidgets('should handle biometric authentication success', (tester) async {
      // Arrange
      final mockUser = User(
        id: '1',
        email: 'test@example.com',
        firstName: 'John',
        lastName: 'Doe',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: true,
        
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final mockAuthResult = EnhancedAuthResult(
        success: true,
        user: mockUser,
        accessToken: 'token123',
        sessionId: 'session123',
        deviceTrusted: true,
      );

      when(() => mockAuthService.restoreSession()).thenAnswer((_) async => mockAuthResult);
      when(() => mockAuthService.canUseBiometrics()).thenAnswer((_) async => true);
      when(() => mockAuthService.isBiometricEnabled('test@example.com')).thenAnswer((_) async => true);
      when(() => mockAuthService.authenticateWithBiometrics('test@example.com')).thenAnswer((_) async => mockAuthResult);

      // Act
      await tester.pumpWidget(createWidget());
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.fingerprint));
      await tester.pumpAndSettle();

      // Assert
      verify(() => mockRouter.go('/home')).called(1);
    });
  });
}
