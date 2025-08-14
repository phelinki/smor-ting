import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/core/models/user.dart';
import 'package:smor_ting_mobile/features/auth/presentation/providers/auth_provider.dart';

void main() {
  group('OTP Disabled Tests', () {
    test('AuthResponse should never require OTP', () {
      // Test that AuthResponse model works without OTP requirement
      final authResponse = AuthResponse(
        user: User(
          id: 'test123',
          email: 'test@example.com',
          firstName: 'Test',
          lastName: 'User',
          phone: '+1234567890',
          role: UserRole.customer,
          isEmailVerified: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        accessToken: 'token123',
        refreshToken: 'refresh123',
        requiresOTP: false, // This should always be false
      );

      // Verify OTP is disabled
      expect(authResponse.requiresOTP, false, reason: 'OTP should be disabled');
      expect(authResponse.accessToken, isNotNull, reason: 'Access token should be provided immediately');
      expect(authResponse.refreshToken, isNotNull, reason: 'Refresh token should be provided immediately');
    });

    test('AuthState should go directly to authenticated, not RequiresOTP', () {
      // Test that we don't use the RequiresOTP auth state
      final user = User(
        id: 'test123',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: false, // Even if not verified, should not require OTP
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      // Should go directly to authenticated state
      final authenticatedState = AuthState.authenticated(user, 'token123');

      // Verify it's an authenticated state, not RequiresOTP
      expect(authenticatedState, isA<Authenticated>(), reason: 'Should be authenticated state');
      expect(authenticatedState, isNot(isA<RequiresOTP>()), reason: 'Should not be RequiresOTP state');
      
      if (authenticatedState is Authenticated) {
        expect(authenticatedState.user.email, user.email, reason: 'Should have same user');
        expect(authenticatedState.accessToken, 'token123', reason: 'Should have access token');
      }
    });

    test('User model isEmailVerified field should not affect login', () {
      // Test that user can login regardless of email verification status
      final unverifiedUser = User(
        id: 'test123',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: false, // Not verified
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final verifiedUser = User(
        id: 'test123',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '+1234567890',
        role: UserRole.customer,
        isEmailVerified: true, // Verified
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      // Both should be valid for authentication
      expect(unverifiedUser.email, isNotEmpty, reason: 'Unverified user should have email');
      expect(verifiedUser.email, isNotEmpty, reason: 'Verified user should have email');
      
      // Verification status should not matter for login flow
      expect(unverifiedUser.isEmailVerified, false, reason: 'Test user is intentionally unverified');
      expect(verifiedUser.isEmailVerified, true, reason: 'Test user is intentionally verified');
    });
  });
}
