import 'package:flutter/material.dart';

import '../../../../core/theme/app_theme.dart';

/// Banner widget to show verification requirements for email and phone
class VerificationRequiredBanner extends StatelessWidget {
  final bool emailVerified;
  final bool phoneVerified;
  final VoidCallback onResendEmail;
  final VoidCallback onResendSms;

  const VerificationRequiredBanner({
    super.key,
    required this.emailVerified,
    required this.phoneVerified,
    required this.onResendEmail,
    required this.onResendSms,
  });

  @override
  Widget build(BuildContext context) {
    // Don't show banner if everything is verified
    if (emailVerified && phoneVerified) {
      return const SizedBox.shrink();
    }

    return Container(
      margin: const EdgeInsets.all(16),
      child: Column(
        children: [
          if (!emailVerified) _buildEmailVerificationBanner(context),
          if (!emailVerified && !phoneVerified) const SizedBox(height: 8),
          if (!phoneVerified) _buildPhoneVerificationBanner(context),
        ],
      ),
    );
  }

  Widget _buildEmailVerificationBanner(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.warning.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppTheme.warning.withOpacity(0.3),
          width: 1,
        ),
      ),
      child: Row(
        children: [
          Icon(
            Icons.email_outlined,
            color: AppTheme.warning,
            size: 24,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Email Verification Required',
                  style: TextStyle(
                    fontWeight: FontWeight.w600,
                    color: AppTheme.warning,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  'Please check your email and click the verification link to access all features.',
                  style: TextStyle(
                    color: Colors.grey[700],
                    fontSize: 14,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          TextButton(
            onPressed: onResendEmail,
            style: TextButton.styleFrom(
              foregroundColor: AppTheme.warning,
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            ),
            child: const Text(
              'Resend Email',
              style: TextStyle(
                fontWeight: FontWeight.w600,
                fontSize: 12,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPhoneVerificationBanner(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.secondaryBlue.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppTheme.secondaryBlue.withOpacity(0.3),
          width: 1,
        ),
      ),
      child: Row(
        children: [
          Icon(
            Icons.sms_outlined,
            color: AppTheme.secondaryBlue,
            size: 24,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Phone Verification Required',
                  style: TextStyle(
                    fontWeight: FontWeight.w600,
                    color: AppTheme.secondaryBlue,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  'Verify your phone number to enable SMS notifications and account recovery.',
                  style: TextStyle(
                    color: Colors.grey[700],
                    fontSize: 14,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          TextButton(
            onPressed: onResendSms,
            style: TextButton.styleFrom(
              foregroundColor: AppTheme.secondaryBlue,
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            ),
            child: const Text(
              'Resend SMS',
              style: TextStyle(
                fontWeight: FontWeight.w600,
                fontSize: 12,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
