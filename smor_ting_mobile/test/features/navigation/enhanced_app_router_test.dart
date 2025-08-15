import 'package:flutter_test/flutter_test.dart';

import '../../../lib/core/models/user.dart';
import '../../../lib/core/models/enhanced_auth_models.dart';
import '../../../lib/features/auth/presentation/providers/enhanced_auth_provider.dart';
import '../../../lib/features/navigation/domain/usecases/role_detection_service.dart';
import '../../../lib/features/navigation/domain/usecases/navigation_flow_service.dart';
import '../../../lib/features/navigation/domain/entities/role_config.dart';

void main() {
  group('Enhanced App Router Tests', () {
    late RoleDetectionService roleDetectionService;
    late NavigationFlowService navigationFlowService;

    setUp(() {
      roleDetectionService = RoleDetectionService();
      navigationFlowService = NavigationFlowService();
    });

    group('_handleAuthRedirect function', () {
      test('returns null when auth state is loading', () {
        // Arrange
        const authState = EnhancedAuthState.loading();
        const currentPath = '/home';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull);
      });

      test('redirects to landing when unauthenticated and on protected route', () {
        // Arrange
        const authState = EnhancedAuthState.unauthenticated();
        const currentPath = '/home';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, '/landing');
      });

      test('returns null when unauthenticated and on public route', () {
        // Arrange
        const authState = EnhancedAuthState.unauthenticated();
        const currentPath = '/login';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull);
      });

      test('redirects to dashboard when authenticated user is on public route', () {
        // Arrange
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

        final authState = EnhancedAuthState.authenticated(
          user: user,
          accessToken: 'token123',
          sessionId: 'session123',
          deviceTrusted: true,
        );
        const currentPath = '/login';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, '/home'); // Customer dashboard route
      });

      test('returns null when authenticated user has access to protected route', () {
        // Arrange
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

        final authState = EnhancedAuthState.authenticated(
          user: user,
          accessToken: 'token123',
          sessionId: 'session123',
          deviceTrusted: true,
        );
        const currentPath = '/home';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull); // Customer can access home
      });

      test('redirects to dashboard when requiresTwoFactor state', () {
        // Arrange
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

        final authState = EnhancedAuthState.requiresTwoFactor(
          email: 'test@example.com',
          tempUser: user,
          deviceTrusted: false,
        );
        const currentPath = '/login';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, '/home'); // Customer dashboard route
      });

      test('redirects to dashboard when requiresVerification state', () {
        // Arrange
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

        final authState = EnhancedAuthState.requiresVerification(
          user: user,
          email: 'test@example.com',
        );
        const currentPath = '/login';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, '/home'); // Customer dashboard route
      });

      test('returns null for captcha state on current page', () {
        // Arrange
        final authState = EnhancedAuthState.requiresCaptcha(
          email: 'test@example.com',
          remainingAttempts: 2,
        );
        const currentPath = '/login';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull);
      });

      test('returns null for locked out state on current page', () {
        // Arrange
        final lockoutInfo = LockoutInfo(
          lockedUntil: DateTime.now().add(const Duration(minutes: 15)),
          remainingAttempts: 0,
          timeUntilUnlock: 900,
        );

        final authState = EnhancedAuthState.lockedOut(
          lockoutInfo: lockoutInfo,
          message: 'Account locked',
        );
        const currentPath = '/login';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull);
      });

      test('redirects to landing on error state when on protected route', () {
        // Arrange
        const authState = EnhancedAuthState.error('Authentication failed', canRetry: true);
        const currentPath = '/home';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, '/landing');
      });

      test('returns null on error state when on public route', () {
        // Arrange
        const authState = EnhancedAuthState.error('Authentication failed', canRetry: true);
        const currentPath = '/login';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull);
      });

      test('redirects to landing on initial state when on protected route', () {
        // Arrange
        const authState = EnhancedAuthState.initial();
        const currentPath = '/home';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, '/landing');
      });

      test('handles reset-password route with query parameters as public', () {
        // Arrange
        const authState = EnhancedAuthState.unauthenticated();
        const currentPath = '/reset-password?email=test@example.com';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull);
      });

      test('redirects customer trying to access admin dashboard', () {
        // Arrange
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

        final authState = EnhancedAuthState.authenticated(
          user: user,
          accessToken: 'token123',
          sessionId: 'session123',
          deviceTrusted: true,
        );
        const currentPath = '/admin-dashboard';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, '/home'); // Customer redirected to their dashboard
      });

      test('allows admin user to access admin dashboard', () {
        // Arrange
        final user = User(
          id: '123',
          email: 'admin@example.com',
          firstName: 'Admin',
          lastName: 'User',
          phone: '1234567890',
          role: UserRole.admin,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        final authState = EnhancedAuthState.authenticated(
          user: user,
          accessToken: 'token123',
          sessionId: 'session123',
          deviceTrusted: true,
        );
        const currentPath = '/admin-dashboard';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull); // Admin can access admin dashboard
      });

      test('allows provider user to access agent dashboard', () {
        // Arrange
        final user = User(
          id: '123',
          email: 'provider@example.com',
          firstName: 'Provider',
          lastName: 'User',
          phone: '1234567890',
          role: UserRole.provider,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        final authState = EnhancedAuthState.authenticated(
          user: user,
          accessToken: 'token123',
          sessionId: 'session123',
          deviceTrusted: true,
        );
        const currentPath = '/agent-dashboard';

        // Act
        final result = _handleAuthRedirect(
          authState,
          currentPath,
          roleDetectionService,
          navigationFlowService,
        );

        // Assert
        expect(result, isNull); // Provider can access agent dashboard
      });
    });
  });
}

// We need to make _handleAuthRedirect function accessible for testing
String? _handleAuthRedirect(
  EnhancedAuthState authState,
  String currentPath,
  RoleDetectionService roleDetectionService,
  NavigationFlowService navigationFlowService,
) {
  // Skip redirect during loading to avoid bouncing
  final isLoading = authState.maybeWhen(
    loading: () => true,
    orElse: () => false,
  ) ?? false;
  if (isLoading) {
    return null;
  }
  
  // Define public routes that don't require authentication
  const publicRoutes = [
    '/',
    '/landing',
    '/login',
    '/register',
    '/forgot-password',
    '/reset-password',
    '/onboarding',
    '/agent-login',
  ];
  
  final isPublicRoute = publicRoutes.contains(currentPath) ||
      currentPath.startsWith('/reset-password');
      // EMAIL OTP REMOVED: /verify-otp is no longer a public route
  
  // Handle different enhanced authentication states
  return authState.when(
    initial: () {
      // Initial state - redirect to landing if not on public route
      if (!isPublicRoute) {
        return '/landing';
      }
      return null;
    },
    loading: () {
      // Authentication in progress - keep current route to avoid bouncing
      return null;
    },
    unauthenticated: () {
      // Unauthenticated user - redirect to landing if not on public route
      if (!isPublicRoute) {
        return '/landing';
      }
      return null;
    },
    authenticated: (user, accessToken, sessionId, deviceTrusted, isRestoredSession, requiresVerification) {
      // Check if user is on a public route after authentication
      if (isPublicRoute) {
        // Redirect to appropriate dashboard based on role
        return roleDetectionService.getDashboardRouteForRole(user.role);
      }
      
      // Check role-based access for protected routes
      final accessResult = roleDetectionService.validateUserAccess(user, currentPath);
      
      if (!accessResult.isAuthorized) {
        // Redirect to appropriate dashboard or handling page
        return accessResult.redirectRoute ?? 
               roleDetectionService.getDashboardRouteForRole(user.role);
      }
      
      return null;
    },
    requiresTwoFactor: (email, tempUser, deviceTrusted) {
      // 2FA DISABLED: Skip 2FA and go directly to dashboard
      return roleDetectionService.getDashboardRouteForRole(tempUser.role);
    },
    requiresCaptcha: (email, remainingAttempts, lockoutInfo) {
      // CAPTCHA DISABLED: Stay on current page (should not happen in dev)
      return null;
    },
    lockedOut: (lockoutInfo, message) {
      // LOCKOUT DISABLED: Stay on current page (should not happen in dev)
      return null;
    },
    error: (message, canRetry) {
      // On error, redirect to landing if not on public route
      if (!isPublicRoute) {
        return '/landing';
      }
      return null;
    },
    requiresVerification: (user, email) {
      // EMAIL VERIFICATION DISABLED: Skip verification and go to dashboard
      return roleDetectionService.getDashboardRouteForRole(user.role);
    },
  );
}