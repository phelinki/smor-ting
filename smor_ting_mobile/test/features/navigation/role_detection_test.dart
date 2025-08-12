import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/core/models/user.dart';
import 'package:smor_ting_mobile/features/navigation/domain/entities/role_config.dart';
import 'package:smor_ting_mobile/features/navigation/domain/usecases/role_detection_service.dart';

void main() {
  group('Role Detection Service', () {
    late RoleDetectionService roleDetectionService;

    setUp(() {
      roleDetectionService = RoleDetectionService();
    });

    group('getUserRole', () {
      test('should return customer role for customer user', () {
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

        // Act
        final result = roleDetectionService.getUserRole(user);

        // Assert
        expect(result, UserRole.customer);
      });

      test('should return provider role for provider user', () {
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

        // Act
        final result = roleDetectionService.getUserRole(user);

        // Assert
        expect(result, UserRole.provider);
      });

      test('should return admin role for admin user', () {
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

        // Act
        final result = roleDetectionService.getUserRole(user);

        // Assert
        expect(result, UserRole.admin);
      });
    });

    group('getDashboardRouteForRole', () {
      test('should return /home for customer role', () {
        // Act
        final result = roleDetectionService.getDashboardRouteForRole(UserRole.customer);

        // Assert
        expect(result, '/home');
      });

      test('should return /agent-dashboard for provider role', () {
        // Act
        final result = roleDetectionService.getDashboardRouteForRole(UserRole.provider);

        // Assert
        expect(result, '/agent-dashboard');
      });

      test('should return /admin-dashboard for admin role', () {
        // Act
        final result = roleDetectionService.getDashboardRouteForRole(UserRole.admin);

        // Assert
        expect(result, '/admin-dashboard');
      });
    });

    group('isRoleAuthorizedForRoute', () {
      test('should allow customer access to customer routes', () {
        // Arrange
        const allowedRoutes = ['/home', '/services', '/profile', '/bookings-history'];

        for (final route in allowedRoutes) {
          // Act
          final result = roleDetectionService.isRoleAuthorizedForRoute(
            UserRole.customer,
            route,
          );

          // Assert
          expect(result, true, reason: 'Customer should access $route');
        }
      });

      test('should deny customer access to provider/admin routes', () {
        // Arrange
        const deniedRoutes = ['/agent-dashboard', '/admin-dashboard', '/kyc'];

        for (final route in deniedRoutes) {
          // Act
          final result = roleDetectionService.isRoleAuthorizedForRoute(
            UserRole.customer,
            route,
          );

          // Assert
          expect(result, false, reason: 'Customer should not access $route');
        }
      });

      test('should allow provider access to provider routes', () {
        // Arrange
        const allowedRoutes = ['/agent-dashboard', '/kyc', '/profile'];

        for (final route in allowedRoutes) {
          // Act
          final result = roleDetectionService.isRoleAuthorizedForRoute(
            UserRole.provider,
            route,
          );

          // Assert
          expect(result, true, reason: 'Provider should access $route');
        }
      });

      test('should deny provider access to admin routes', () {
        // Arrange
        const deniedRoutes = ['/admin-dashboard'];

        for (final route in deniedRoutes) {
          // Act
          final result = roleDetectionService.isRoleAuthorizedForRoute(
            UserRole.provider,
            route,
          );

          // Assert
          expect(result, false, reason: 'Provider should not access $route');
        }
      });

      test('should allow admin access to all routes', () {
        // Arrange
        const allRoutes = [
          '/home',
          '/services',
          '/agent-dashboard',
          '/admin-dashboard',
          '/kyc',
          '/profile',
        ];

        for (final route in allRoutes) {
          // Act
          final result = roleDetectionService.isRoleAuthorizedForRoute(
            UserRole.admin,
            route,
          );

          // Assert
          expect(result, true, reason: 'Admin should access $route');
        }
      });
    });

    group('getAvailableFeatures', () {
      test('should return customer features for customer role', () {
        // Act
        final result = roleDetectionService.getAvailableFeatures(UserRole.customer);

        // Assert
        expect(result, contains('browse_services'));
        expect(result, contains('book_services'));
        expect(result, contains('track_bookings'));
        expect(result, contains('payment_wallet'));
        expect(result, isNot(contains('manage_services')));
        expect(result, isNot(contains('admin_panel')));
      });

      test('should return provider features for provider role', () {
        // Act
        final result = roleDetectionService.getAvailableFeatures(UserRole.provider);

        // Assert
        expect(result, contains('manage_services'));
        expect(result, contains('view_earnings'));
        expect(result, contains('kyc_verification'));
        expect(result, contains('job_management'));
        expect(result, isNot(contains('admin_panel')));
      });

      test('should return admin features for admin role', () {
        // Act
        final result = roleDetectionService.getAvailableFeatures(UserRole.admin);

        // Assert
        expect(result, contains('admin_panel'));
        expect(result, contains('user_management'));
        expect(result, contains('system_analytics'));
        expect(result, contains('manage_services'));
        expect(result, contains('browse_services'));
      });
    });

    group('validateUserAccess', () {
      test('should validate successful access for authorized user and route', () {
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

        // Act
        final result = roleDetectionService.validateUserAccess(user, '/home');

        // Assert
        expect(result.isAuthorized, true);
        expect(result.redirectRoute, isNull);
        expect(result.denialReason, isNull);
      });

      test('should return denial with redirect for unauthorized user and route', () {
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

        // Act
        final result = roleDetectionService.validateUserAccess(user, '/admin-dashboard');

        // Assert
        expect(result.isAuthorized, false);
        expect(result.redirectRoute, '/home'); // Should redirect to customer dashboard
        expect(result.denialReason, 'Insufficient permissions for this route');
      });

      test('should handle unverified email appropriately', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'customer@test.com',
          firstName: 'John',
          lastName: 'Doe',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: false,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        // Act
        final result = roleDetectionService.validateUserAccess(user, '/services');

        // Assert
        expect(result.isAuthorized, false);
        expect(result.redirectRoute, '/verify-otp');
        expect(result.denialReason, 'Email verification required');
      });
    });
  });
}
