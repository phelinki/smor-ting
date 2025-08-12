import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import '../../../../../lib/features/auth/presentation/widgets/enhanced_login_outcomes.dart';

void main() {
  group('EnhancedLoginOutcomes', () {
    Widget createWidget({
      bool showCooldown = false,
      bool showCaptcha = false,
      bool showTwoFactor = false,
      int cooldownSeconds = 0,
      int remainingAttempts = 0,
    }) {
      return MaterialApp(
        home: Scaffold(
          body: EnhancedLoginOutcomes(
            showCooldown: showCooldown,
            showCaptcha: showCaptcha,
            showTwoFactor: showTwoFactor,
            cooldownSeconds: cooldownSeconds,
            remainingAttempts: remainingAttempts,
            onCaptchaCompleted: (token) {},
            onTwoFactorSubmitted: (code) {},
          ),
        ),
      );
    }

    testWidgets('should show cooldown timer when showCooldown is true', (tester) async {
      // Act
      await tester.pumpWidget(createWidget(
        showCooldown: true,
        cooldownSeconds: 60,
      ));

      // Assert
      expect(find.text('Account Temporarily Locked'), findsOneWidget);
      expect(find.textContaining('Try again in'), findsOneWidget);
    });

    testWidgets('should show CAPTCHA when showCaptcha is true', (tester) async {
      // Act
      await tester.pumpWidget(createWidget(
        showCaptcha: true,
        remainingAttempts: 3,
      ));

      // Assert
      expect(find.text('Security Verification'), findsOneWidget);
      expect(find.textContaining('attempts remaining'), findsOneWidget);
    });

    testWidgets('should show 2FA input when showTwoFactor is true', (tester) async {
      // Act
      await tester.pumpWidget(createWidget(showTwoFactor: true));

      // Assert
      expect(find.text('Two-Factor Authentication'), findsOneWidget);
      expect(find.text('Enter the 6-digit code from your authenticator app'), findsOneWidget);
    });

    testWidgets('should not show any content when all flags are false', (tester) async {
      // Act
      await tester.pumpWidget(createWidget());

      // Assert
      expect(find.byType(EnhancedLoginOutcomes), findsOneWidget);
      expect(find.text('Account Temporarily Locked'), findsNothing);
      expect(find.text('Security Verification'), findsNothing);
      expect(find.text('Two-Factor Authentication'), findsNothing);
    });
  });
}
