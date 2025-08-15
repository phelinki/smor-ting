import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/core/models/user.dart';
import 'package:smor_ting_mobile/features/navigation/domain/entities/navigation_result.dart';
import 'package:smor_ting_mobile/features/navigation/domain/usecases/navigation_flow_service.dart';
import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';

void main() {
  group('Navigation Flow Service', () {
    late NavigationFlowService navigationFlowService;

    setUp(() {
      navigationFlowService = NavigationFlowService();
    });

    group('getPostLoginDestination', () {
      test('should navigate customer to home dashboard after login', () {
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
        final result = navigationFlowService.getPostLoginDestination(user);

        // Assert
        expect(result.destination, '/home');
        expect(result.shouldReplace, true);
        expect(result.clearHistory, true);
      });

      test('should navigate provider to agent dashboard after login', () {
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
        final result = navigationFlowService.getPostLoginDestination(user);

        // Assert
        expect(result.destination, '/agent-dashboard');
        expect(result.shouldReplace, true);
        expect(result.clearHistory, true);
      });

      test('should navigate admin to admin dashboard after login', () {
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
        final result = navigationFlowService.getPostLoginDestination(user);

        // Assert
        expect(result.destination, '/admin-dashboard');
        expect(result.shouldReplace, true);
        expect(result.clearHistory, true);
      });

      test('should navigate to OTP verification if email not verified', () {
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

        // Act
        final result = navigationFlowService.getPostLoginDestination(user);

        // Assert
        expect(result.destination, '/verify-otp');
        expect(result.shouldReplace, true);
        expect(result.queryParameters, containsPair('email', 'unverified@test.com'));
        expect(result.queryParameters, containsPair('fullName', 'John Doe'));
      });

      test('should navigate to KYC if provider needs verification', () {
        // Arrange
        final user = User(
          id: '2',
          email: 'newprovider@test.com',
          firstName: 'New',
          lastName: 'Provider',
          phone: '+1234567890',
          role: UserRole.provider,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        // Act
        final result = navigationFlowService.getPostLoginDestination(
          user,
          requiresKyc: true,
        );

        // Assert
        expect(result.destination, '/kyc');
        expect(result.shouldReplace, true);
        expect(result.clearHistory, true);
      });
    });

    group('getPostRegistrationDestination', () {
      test('should navigate to OTP verification after customer registration', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'newcustomer@test.com',
          firstName: 'New',
          lastName: 'Customer',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: false,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        // Act
        final result = navigationFlowService.getPostRegistrationDestination(user);

        // Assert
        expect(result.destination, '/verify-otp');
        expect(result.shouldReplace, true);
        expect(result.queryParameters, containsPair('email', 'newcustomer@test.com'));
        expect(result.queryParameters, containsPair('fullName', 'New Customer'));
      });

      test('should navigate to OTP verification after provider registration', () {
        // Arrange
        final user = User(
          id: '2',
          email: 'newprovider@test.com',
          firstName: 'New',
          lastName: 'Provider',
          phone: '+1234567890',
          role: UserRole.provider,
          isEmailVerified: false,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        // Act
        final result = navigationFlowService.getPostRegistrationDestination(user);

        // Assert
        expect(result.destination, '/verify-otp');
        expect(result.shouldReplace, true);
        expect(result.queryParameters, containsPair('email', 'newprovider@test.com'));
        expect(result.queryParameters, containsPair('fullName', 'New Provider'));
      });

      test('should navigate directly to dashboard if already verified', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'verified@test.com',
          firstName: 'Already',
          lastName: 'Verified',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        // Act
        final result = navigationFlowService.getPostRegistrationDestination(user);

        // Assert
        expect(result.destination, '/home');
        expect(result.shouldReplace, true);
        expect(result.clearHistory, true);
      });
    });



    group('getPostKYCDestination', () {
      test('should navigate provider to dashboard after successful KYC', () {
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
        final result = navigationFlowService.getPostKYCDestination(user, kycSuccess: true);

        // Assert
        expect(result.destination, '/agent-dashboard');
        expect(result.shouldReplace, true);
        expect(result.clearHistory, true);
      });

      test('should stay on KYC page if verification failed', () {
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
        final result = navigationFlowService.getPostKYCDestination(user, kycSuccess: false);

        // Assert
        expect(result.destination, '/kyc');
        expect(result.shouldReplace, false);
        expect(result.message, 'KYC verification failed. Please try again.');
      });
    });

    group('getLogoutDestination', () {
      test('should navigate to landing page after logout', () {
        // Act
        final result = navigationFlowService.getLogoutDestination();

        // Assert
        expect(result.destination, '/landing');
        expect(result.shouldReplace, true);
        expect(result.clearHistory, true);
      });
    });

    group('handleDeepLink', () {
      test('should validate deep link access based on user role', () {
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
        final result = navigationFlowService.handleDeepLink(
          '/services/electronics',
          user,
        );

        // Assert
        expect(result.destination, '/services/electronics');
        expect(result.shouldReplace, false);
      });

      test('should redirect unauthorized deep link to appropriate dashboard', () {
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
        final result = navigationFlowService.handleDeepLink(
          '/agent-dashboard',
          user,
        );

        // Assert
        expect(result.destination, '/home');
        expect(result.shouldReplace, true);
        expect(result.message, 'Access denied. Redirected to your dashboard.');
      });

      test('should handle deep link for unauthenticated user', () {
        // Act
        final result = navigationFlowService.handleDeepLink(
          '/services/electronics',
          null,
        );

        // Assert
        expect(result.destination, '/landing');
        expect(result.shouldReplace, true);
        expect(result.queryParameters, containsPair('redirect', '/services/electronics'));
        expect(result.message, 'Please log in to access this page.');
      });
    });

    group('shouldShowOnboarding', () {
      test('should show onboarding for new customer users', () {
        // Arrange
        final user = User(
          id: '1',
          email: 'newcustomer@test.com',
          firstName: 'New',
          lastName: 'Customer',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        // Act
        final result = navigationFlowService.shouldShowOnboarding(user, isFirstLogin: true);

        // Assert
        expect(result, true);
      });

      test('should not show onboarding for returning users', () {
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
        final result = navigationFlowService.shouldShowOnboarding(user, isFirstLogin: false);

        // Assert
        expect(result, false);
      });

      test('should not show onboarding for admin users', () {
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
        final result = navigationFlowService.shouldShowOnboarding(user, isFirstLogin: true);

        // Assert
        expect(result, false);
      });
    });
  });
}
