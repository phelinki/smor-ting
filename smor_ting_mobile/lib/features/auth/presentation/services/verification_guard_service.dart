import '../../../../core/models/user.dart';

/// Service to determine which verification requirements should be shown/enforced
class VerificationGuardService {
  /// Routes that require email verification
  static const List<String> _emailVerificationRequiredRoutes = [
    '/wallet',
    '/payment',
    '/booking',
    '/profile/edit',
    '/settings/security',
  ];

  /// Routes that require phone verification  
  static const List<String> _phoneVerificationRequiredRoutes = [
    '/booking',
    '/emergency-services',
    '/wallet/withdraw',
    '/two-factor-setup',
  ];

  /// Routes that require both email and phone verification
  static const List<String> _fullVerificationRequiredRoutes = [
    '/kyc',
    '/agent-application',
    '/become-provider',
    '/admin',
  ];

  /// Routes that are accessible even without verification (verification optional)
  static const List<String> _verificationOptionalRoutes = [
    '/home',
    '/services',
    '/search',
    '/help',
    '/about',
    '/settings',
    '/profile',
  ];

  /// Check what verification is required for a given route
  VerificationRequirement checkVerificationRequirement({
    required String route,
    required User user,
  }) {
    final cleanRoute = route.split('?').first;

    // Check if route requires full verification
    if (_fullVerificationRequiredRoutes.any((r) => cleanRoute.startsWith(r))) {
      if (!user.isEmailVerified || !user.isPhoneVerified) {
        return VerificationRequirement(
          required: true,
          blocking: true,
          emailRequired: !user.isEmailVerified,
          phoneRequired: !user.isPhoneVerified,
          reason: 'This feature requires both email and phone verification for security.',
        );
      }
    }

    // Check if route requires email verification
    if (_emailVerificationRequiredRoutes.any((r) => cleanRoute.startsWith(r))) {
      if (!user.isEmailVerified) {
        return VerificationRequirement(
          required: true,
          blocking: true,
          emailRequired: true,
          phoneRequired: false,
          reason: 'Email verification is required to access this feature.',
        );
      }
    }

    // Check if route requires phone verification
    if (_phoneVerificationRequiredRoutes.any((r) => cleanRoute.startsWith(r))) {
      if (!user.isPhoneVerified) {
        return VerificationRequirement(
          required: true,
          blocking: true,
          emailRequired: false,
          phoneRequired: true,
          reason: 'Phone verification is required to access this feature.',
        );
      }
    }

    // For optional verification routes, show banner but don't block
    if (_verificationOptionalRoutes.any((r) => cleanRoute.startsWith(r))) {
      if (!user.isEmailVerified || !user.isPhoneVerified) {
        return VerificationRequirement(
          required: true,
          blocking: false,
          emailRequired: !user.isEmailVerified,
          phoneRequired: !user.isPhoneVerified,
          reason: 'Complete verification to access all features and enhance security.',
        );
      }
    }

    // No verification required
    return VerificationRequirement(
      required: false,
      blocking: false,
      emailRequired: false,
      phoneRequired: false,
      reason: null,
    );
  }

  /// Check if a route should show verification prompts (non-blocking)
  bool shouldShowVerificationBanner({
    required String route,
    required User user,
  }) {
    final requirement = checkVerificationRequirement(route: route, user: user);
    return requirement.required && !requirement.blocking;
  }

  /// Check if a route should be blocked due to verification requirements
  bool shouldBlockRoute({
    required String route,
    required User user,
  }) {
    final requirement = checkVerificationRequirement(route: route, user: user);
    return requirement.required && requirement.blocking;
  }

  /// Get the verification requirement for overlay display
  VerificationRequirement getVerificationRequirementForRoute({
    required String route,
    required User user,
  }) {
    return checkVerificationRequirement(route: route, user: user);
  }
}

/// Represents verification requirements for a route
class VerificationRequirement {
  final bool required;
  final bool blocking;
  final bool emailRequired;
  final bool phoneRequired;
  final String? reason;

  VerificationRequirement({
    required this.required,
    required this.blocking,
    required this.emailRequired,
    required this.phoneRequired,
    required this.reason,
  });

  bool get hasAnyRequirement => emailRequired || phoneRequired;
  bool get requiresBoth => emailRequired && phoneRequired;
}
