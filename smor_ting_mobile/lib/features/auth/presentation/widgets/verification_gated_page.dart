import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/models/user.dart';
import '../../../../core/services/message_service.dart';
import '../providers/auth_provider.dart';
import '../services/verification_guard_service.dart';
import 'verification_blocking_overlay.dart';
import 'verification_required_banner.dart';

/// A wrapper widget that provides verification gating for any page
class VerificationGatedPage extends ConsumerWidget {
  final Widget child;
  final String route;
  final VerificationGuardService _verificationGuard = VerificationGuardService();

  VerificationGatedPage({
    super.key,
    required this.child,
    required this.route,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authNotifierProvider);

    // If user is authenticated, check verification requirements
    if (authState is Authenticated) {
      return _buildGatedContent(context, ref, authState.user);
    }

    // If not authenticated, show child as-is
    return child;
  }

  Widget _buildGatedContent(BuildContext context, WidgetRef ref, User user) {
    final verificationRequirement = _verificationGuard.getVerificationRequirementForRoute(
      route: route,
      user: user,
    );

    if (!verificationRequirement.required) {
      // No verification required, show page as normal
      return child;
    }

    if (verificationRequirement.blocking) {
      // Show blocking overlay
      return Stack(
        children: [
          child,
          VerificationBlockingOverlay(
            showEmailVerification: verificationRequirement.emailRequired,
            showPhoneVerification: verificationRequirement.phoneRequired,
            userEmail: user.email,
            userPhone: user.phone ?? '',
            onResendEmail: () => _handleResendEmail(ref, user.email),
            onResendSms: () => _handleResendSms(ref, user.phone ?? ''),
          ),
        ],
      );
    } else {
      // Show non-blocking banner
      return Column(
        children: [
          VerificationRequiredBanner(
            emailVerified: user.isEmailVerified,
            phoneVerified: user.isEmailVerified, // TODO: Add proper phone verification field to User model
            onResendEmail: () => _handleResendEmail(ref, user.email),
            onResendSms: () => _handleResendSms(ref, user.phone),
          ),
          Expanded(child: child),
        ],
      );
    }
  }

  void _handleResendEmail(WidgetRef ref, String email) {
    // TODO: Implement email resend logic
    // This could call an API service to resend verification email
    print('Resending verification email to: $email');
    
    // Show confirmation message
    MessageService.showSuccess(
      ref.context,
      message: 'Verification email sent! Please check your inbox.',
    );
  }

  void _handleResendSms(WidgetRef ref, String phone) {
    // TODO: Implement SMS resend logic
    // This could call an API service to resend verification SMS
    print('Resending verification SMS to: $phone');
    
    // Show confirmation message
    MessageService.showSuccess(
      ref.context,
      message: 'Verification SMS sent! Please check your messages.',
    );
  }
}

/// Helper function to wrap any widget with verification gating
Widget withVerificationGating({
  required Widget child,
  required String route,
}) {
  return VerificationGatedPage(
    route: route,
    child: child,
  );
}
