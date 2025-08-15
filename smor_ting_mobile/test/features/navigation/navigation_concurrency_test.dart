import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../lib/core/models/user.dart';
import '../../../lib/features/auth/presentation/providers/enhanced_auth_provider.dart';
import '../../../lib/features/navigation/presentation/enhanced_app_router.dart';

void main() {
  group('Navigation Concurrency Tests', () {
    test('should not have concurrent navigation conflicts', () async {
      // Create a container with initial auth state
      final container = ProviderContainer(
        overrides: [
          enhancedAuthNotifierProvider.overrideWith(
            () => TestEnhancedAuthNotifier(const EnhancedAuthState.initial()),
          ),
        ],
      );

      addTearDown(container.dispose);

      // Get the router
      final router = container.read(enhancedAppRouterProvider);

      // Initial state should work
      expect(router.routerDelegate.currentConfiguration.uri.toString(), '/');

      // Simulate auth state change
      final notifier = container.read(enhancedAuthNotifierProvider.notifier) as TestEnhancedAuthNotifier;
      
      final user = User(
        id: '123',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '1234567890',
        role: UserRole.customer,
        isEmailVerified: true,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      // Change auth state to authenticated
      notifier.updateState(EnhancedAuthState.authenticated(
        user: user,
        accessToken: 'token123',
        sessionId: 'session123',
        deviceTrusted: true,
      ));

      // Give the router a moment to process the state change
      await Future.delayed(Duration.zero);

      // Router should handle the navigation automatically without conflicts
      expect(router.routerDelegate.currentConfiguration.uri.toString(), isNot(equals('/')));
    });

    test('should handle rapid auth state changes gracefully', () async {
      final container = ProviderContainer(
        overrides: [
          enhancedAuthNotifierProvider.overrideWith(
            () => TestEnhancedAuthNotifier(const EnhancedAuthState.unauthenticated()),
          ),
        ],
      );

      addTearDown(container.dispose);

      final router = container.read(enhancedAppRouterProvider);
      final notifier = container.read(enhancedAuthNotifierProvider.notifier) as TestEnhancedAuthNotifier;

      final user = User(
        id: '123',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '1234567890',
        role: UserRole.customer,
        isEmailVerified: true,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      // Rapid state changes that previously caused concurrency issues
      notifier.updateState(const EnhancedAuthState.loading());
      await Future.delayed(Duration.zero);
      
      notifier.updateState(EnhancedAuthState.authenticated(
        user: user,
        accessToken: 'token123',
        sessionId: 'session123',
        deviceTrusted: true,
      ));
      await Future.delayed(Duration.zero);

      notifier.updateState(const EnhancedAuthState.unauthenticated());
      await Future.delayed(Duration.zero);

      // Should not throw navigation assertion errors
      expect(() => router.routerDelegate.currentConfiguration, returnsNormally);
    });
  });
}

// Test implementation of EnhancedAuthNotifier for testing
class TestEnhancedAuthNotifier extends EnhancedAuthNotifier {
  EnhancedAuthState _state;

  TestEnhancedAuthNotifier(this._state);

  @override
  EnhancedAuthState build() => _state;

  void updateState(EnhancedAuthState newState) {
    _state = newState;
    // Trigger a rebuild by calling the internal state setter
    ref.notifyListeners();
  }
}
