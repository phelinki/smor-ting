import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import '../../../../../lib/features/auth/presentation/widgets/verification_required_banner.dart';

void main() {
  group('VerificationRequiredBanner', () {
    Widget createWidget({
      bool emailVerified = false,
      bool phoneVerified = false,
      VoidCallback? onResendEmail,
      VoidCallback? onResendSms,
    }) {
      return MaterialApp(
        home: Scaffold(
          body: VerificationRequiredBanner(
            emailVerified: emailVerified,
            phoneVerified: phoneVerified,
            onResendEmail: onResendEmail ?? () {},
            onResendSms: onResendSms ?? () {},
          ),
        ),
      );
    }

    testWidgets('should show email verification banner when email not verified', (tester) async {
      // Act
      await tester.pumpWidget(createWidget(
        emailVerified: false,
        phoneVerified: true,
      ));

      // Assert
      expect(find.text('Email Verification Required'), findsOneWidget);
      expect(find.text('Resend Email'), findsOneWidget);
      expect(find.byIcon(Icons.email_outlined), findsOneWidget);
    });

    testWidgets('should show phone verification banner when phone not verified', (tester) async {
      // Act
      await tester.pumpWidget(createWidget(
        emailVerified: true,
        phoneVerified: false,
      ));

      // Assert
      expect(find.text('Phone Verification Required'), findsOneWidget);
      expect(find.text('Resend SMS'), findsOneWidget);
      expect(find.byIcon(Icons.sms_outlined), findsOneWidget);
    });

    testWidgets('should show both banners when neither is verified', (tester) async {
      // Act
      await tester.pumpWidget(createWidget(
        emailVerified: false,
        phoneVerified: false,
      ));

      // Assert
      expect(find.text('Email Verification Required'), findsOneWidget);
      expect(find.text('Phone Verification Required'), findsOneWidget);
    });

    testWidgets('should not show banner when both are verified', (tester) async {
      // Act
      await tester.pumpWidget(createWidget(
        emailVerified: true,
        phoneVerified: true,
      ));

      // Assert
      expect(find.byType(VerificationRequiredBanner), findsOneWidget);
      expect(find.text('Email Verification Required'), findsNothing);
      expect(find.text('Phone Verification Required'), findsNothing);
    });

    testWidgets('should call onResendEmail when email resend button is tapped', (tester) async {
      // Arrange
      bool emailResendCalled = false;
      
      // Act
      await tester.pumpWidget(createWidget(
        emailVerified: false,
        phoneVerified: true,
        onResendEmail: () => emailResendCalled = true,
      ));

      await tester.tap(find.text('Resend Email'));
      await tester.pumpAndSettle();

      // Assert
      expect(emailResendCalled, isTrue);
    });

    testWidgets('should call onResendSms when SMS resend button is tapped', (tester) async {
      // Arrange
      bool smsResendCalled = false;
      
      // Act
      await tester.pumpWidget(createWidget(
        emailVerified: true,
        phoneVerified: false,
        onResendSms: () => smsResendCalled = true,
      ));

      await tester.tap(find.text('Resend SMS'));
      await tester.pumpAndSettle();

      // Assert
      expect(smsResendCalled, isTrue);
    });
  });
}
