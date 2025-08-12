import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:smor_ting_mobile/core/models/user.dart';
import 'package:smor_ting_mobile/features/navigation/domain/entities/route_guard_result.dart';
import 'package:smor_ting_mobile/features/navigation/domain/usecases/route_guard_service.dart';
import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';

void main() {
  group('Route Guard Service', () {
    late RouteGuardService routeGuardService;

    setUp(() {
      routeGuardService = RouteGuardService();
    });

    group('checkRouteAccess', () {
      test('should allow access to public routes for unauthenticated users', () {
        // Arrange
        const publicRoutes = [
          '/landing',
          '/login',
          '/register',
          '/forgot-password',
          '/onboarding',
        ];

        for (final route in publicRoutes) {
          // Act
          final result = routeGuardService.checkRouteAccess(
            route: route,
            authState: const Initial(),
          );

          // Assert
          expect(result.isAllowed, true, reason: 'Public route $route should be accessible');
          expect(result.redirectRoute, isNull);
        }
      });

      test('should deny access to protected routes for unauthenticated users', () {
        // Arrange
        const protectedRoutes = [
          '/home',
          '/services',
          '/agent-dashboard',
          '/admin-dashboard',
          '/profile',
          '/settings',
        ];

        for (final route in protectedRoutes) {
          // Act
          final result = routeGuardService.checkRouteAccess(
            route: route,
            authState: const Initial(),
          );

          // Assert
          expect(result.isAllowed, false, reason: 'Protected route $route should not be accessible');
          expect(result.redirectRoute, '/landing');
          expect(result.denialReason, 'Authentication required');
        }
      });

      test('should allow access to customer routes for authenticated customer', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'customer@test.com',
          firstName: 'John',
          lastName: 'Doe',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = Authenticated(user, 'token');

        const customerRoutes = [
          '/home',
          '/services',
          '/profile',
          '/settings',
          '/bookings-history',
          '/payment-methods',
        ];

        for (final route in customerRoutes) {
          // Act
          final result = routeGuardService.checkRouteAccess(
            route: route,
            authState: authState,
          );

          // Assert
          expect(result.isAllowed, true, reason: 'Customer should access $route');
          expect(result.redirectRoute, isNull);
        }
      });

      test('should deny access to provider routes for customer user', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'customer@test.com',
          firstName: 'John',
          lastName: 'Doe',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = Authenticated(user, 'token');

        const providerRoutes = [
          '/agent-dashboard',
          '/kyc',
        ];

        for (final route in providerRoutes) {
          // Act
          final result = routeGuardService.checkRouteAccess(
            route: route,
            authState: authState,
          );

          // Assert
          expect(result.isAllowed, false, reason: 'Customer should not access $route');
          expect(result.redirectRoute, '/home');
          expect(result.denialReason, 'Insufficient role permissions');
        }
      });

      test('should allow access to provider routes for provider user', () {
        // Arrange
        final user = User(
          id: '2',
          email: 'provider@test.com',
          firstName: 'Jane',
          lastName: 'Smith',
          phone: '+1234567890',
          role: UserRole.provider,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = Authenticated(user, 'token');

        const providerRoutes = [
          '/agent-dashboard',
          '/kyc',
          '/profile',
          '/settings',
        ];

        for (final route in providerRoutes) {
          // Act
          final result = routeGuardService.checkRouteAccess(
            route: route,
            authState: authState,
          );

          // Assert
          expect(result.isAllowed, true, reason: 'Provider should access $route');
          expect(result.redirectRoute, isNull);
        }
      });

      test('should deny access to admin routes for provider user', () {
        // Arrange
        final user = User(
          id: '2',
          email: 'provider@test.com',
          firstName: 'Jane',
          lastName: 'Smith',
          phone: '+1234567890',
          role: UserRole.provider,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = Authenticated(user, 'token');

        const adminRoutes = [
          '/admin-dashboard',
        ];

        for (final route in adminRoutes) {
          // Act
          final result = routeGuardService.checkRouteAccess(
            route: route,
            authState: authState,
          );

          // Assert
          expect(result.isAllowed, false, reason: 'Provider should not access $route');
          expect(result.redirectRoute, '/agent-dashboard');
          expect(result.denialReason, 'Insufficient role permissions');
        }
      });

      test('should allow access to all routes for admin user', () {
        // Arrange
        final user = User(
          id: '3',
          email: 'admin@test.com',
          firstName: 'Admin',
          lastName: 'User',
          phone: '+1234567890',
          role: UserRole.admin,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = Authenticated(user, 'token');

        const allRoutes = [
          '/home',
          '/services',
          '/agent-dashboard',
          '/admin-dashboard',
          '/kyc',
          '/profile',
          '/settings',
        ];

        for (final route in allRoutes) {
          // Act
          final result = routeGuardService.checkRouteAccess(
            route: route,
            authState: authState,
          );

          // Assert
          expect(result.isAllowed, true, reason: 'Admin should access $route');
          expect(result.redirectRoute, isNull);
        }
      });

      test('should handle email verification requirement', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'unverified@test.com',
          firstName: 'John',
          lastName: 'Doe',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: false,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = Authenticated(user, 'token');

        // Act
        final result = routeGuardService.checkRouteAccess(
          route: '/services',
          authState: authState,
        );

        // Assert
        expect(result.isAllowed, false);
        expect(result.redirectRoute, '/verify-otp');
        expect(result.denialReason, 'Email verification required');
      });

      test('should handle RequiresOTP auth state', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'user@test.com',
          firstName: 'John',
          lastName: 'Doe',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: false,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = RequiresOTP(email: 'user@test.com', user: user);

        // Act
        final result = routeGuardService.checkRouteAccess(
          route: '/home',
          authState: authState,
        );

        // Assert
        expect(result.isAllowed, false);
        expect(result.redirectRoute, '/verify-otp');
        expect(result.denialReason, 'OTP verification required');
      });

      test('should handle Loading auth state', () {
        // Arrange
        const authState = Loading();

        // Act
        final result = routeGuardService.checkRouteAccess(
          route: '/home',
          authState: authState,
        );

        // Assert
        expect(result.isAllowed, false);
        expect(result.redirectRoute, isNull); // No redirect during loading
        expect(result.denialReason, 'Authentication in progress');
      });
    });

    group('isPublicRoute', () {
      test('should identify public routes correctly', () {
        // Arrange
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

        for (final route in publicRoutes) {
          // Act
          final result = routeGuardService.isPublicRoute(route);

          // Assert
          expect(result, true, reason: '$route should be public');
        }
      });

      test('should identify protected routes correctly', () {
        // Arrange
        const protectedRoutes = [
          '/home',
          '/services',
          '/agent-dashboard',
          '/admin-dashboard',
          '/profile',
          '/settings',
          '/bookings-history',
        ];

        for (final route in protectedRoutes) {
          // Act
          final result = routeGuardService.isPublicRoute(route);

          // Assert
          expect(result, false, reason: '$route should be protected');
        }
      });
    });

    group('getRedirectRouteForAuthState', () {
      test('should redirect unauthenticated users to landing page', () {
        // Act
        final result = routeGuardService.getRedirectRouteForAuthState(
          const Initial(),
          '/home',
        );

        // Assert
        expect(result, '/landing');
      });

      test('should redirect to OTP verification when required', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'user@test.com',
          firstName: 'John',
          lastName: 'Doe',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: false,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = RequiresOTP(email: 'user@test.com', user: user);

        // Act
        final result = routeGuardService.getRedirectRouteForAuthState(
          authState,
          '/home',
        );

        // Assert
        expect(result, '/verify-otp');
      });

      test('should redirect based on user role when accessing wrong dashboard', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'customer@test.com',
          firstName: 'John',
          lastName: 'Doe',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        final authState = Authenticated(user, 'token');

        // Act
        final result = routeGuardService.getRedirectRouteForAuthState(
          authState,
          '/agent-dashboard',
        );

        // Assert
        expect(result, '/home');
      });
    });
  });
}
