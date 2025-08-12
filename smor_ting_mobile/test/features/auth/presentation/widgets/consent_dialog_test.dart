import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:smor_ting/core/models/consent.dart';
import 'package:smor_ting/core/services/consent_service.dart';
import 'package:smor_ting/features/auth/presentation/widgets/consent_dialog.dart';

import 'consent_dialog_test.mocks.dart';

@GenerateMocks([ConsentService])
void main() {
  late MockConsentService mockConsentService;

  setUp(() {
    mockConsentService = MockConsentService();
  });

  Widget createTestWidget({
    required List<ConsentRequirement> requirements,
    VoidCallback? onConsentsGiven,
    VoidCallback? onCancel,
  }) {
    return ProviderScope(
      overrides: [
        consentServiceProvider.overrideWithValue(mockConsentService),
      ],
      child: MaterialApp(
        home: Scaffold(
          body: ConsentDialog(
            userId: 'test_user',
            requirements: requirements,
            onConsentsGiven: onConsentsGiven,
            onCancel: onCancel,
          ),
        ),
      ),
    );
  }

  group('ConsentDialog', () {
    testWidgets('should display all consent requirements', (tester) async {
      // Arrange
      final requirements = [
        const ConsentRequirement(
          type: ConsentType.termsOfService,
          title: 'Terms of Service',
          description: 'Accept our terms of service',
          version: '1.0',
          required: true,
          documentUrl: 'https://example.com/terms',
        ),
        const ConsentRequirement(
          type: ConsentType.privacyPolicy,
          title: 'Privacy Policy',
          description: 'Accept our privacy policy',
          version: '1.0',
          required: true,
          documentUrl: 'https://example.com/privacy',
        ),
        const ConsentRequirement(
          type: ConsentType.marketingCommunications,
          title: 'Marketing Communications',
          description: 'Receive marketing emails',
          version: '1.0',
          required: false,
        ),
      ];

      // Act
      await tester.pumpWidget(createTestWidget(requirements: requirements));

      // Assert
      expect(find.text('Privacy & Consent'), findsOneWidget);
      expect(find.text('Terms of Service'), findsOneWidget);
      expect(find.text('Privacy Policy'), findsOneWidget);
      expect(find.text('Marketing Communications'), findsOneWidget);
      expect(find.text('Required'), findsNWidgets(2)); // Two required items
      expect(find.byType(Checkbox), findsNWidgets(3));
    });

    testWidgets('should show external link icons for requirements with document URLs', (tester) async {
      // Arrange
      final requirements = [
        const ConsentRequirement(
          type: ConsentType.termsOfService,
          title: 'Terms of Service',
          description: 'Accept our terms',
          version: '1.0',
          required: true,
          documentUrl: 'https://example.com/terms',
        ),
        const ConsentRequirement(
          type: ConsentType.marketingCommunications,
          title: 'Marketing',
          description: 'Marketing emails',
          version: '1.0',
          required: false,
          // No document URL
        ),
      ];

      // Act
      await tester.pumpWidget(createTestWidget(requirements: requirements));

      // Assert
      expect(find.byIcon(Icons.open_in_new), findsOneWidget); // Only one has URL
    });

    testWidgets('should disable continue button when required consents not given', (tester) async {
      // Arrange
      final requirements = [
        const ConsentRequirement(
          type: ConsentType.termsOfService,
          title: 'Terms of Service',
          description: 'Accept our terms',
          version: '1.0',
          required: true,
        ),
      ];

      // Act
      await tester.pumpWidget(createTestWidget(requirements: requirements));

      // Assert
      final continueButton = find.text('Continue');
      expect(continueButton, findsOneWidget);
      
      final buttonWidget = tester.widget<ElevatedButton>(
        find.byType(ElevatedButton).first,
      );
      expect(buttonWidget.onPressed, isNull); // Button should be disabled
    });

    testWidgets('should enable continue button when all required consents are given', (tester) async {
      // Arrange
      final requirements = [
        const ConsentRequirement(
          type: ConsentType.termsOfService,
          title: 'Terms of Service',
          description: 'Accept our terms',
          version: '1.0',
          required: true,
        ),
        const ConsentRequirement(
          type: ConsentType.marketingCommunications,
          title: 'Marketing',
          description: 'Marketing emails',
          version: '1.0',
          required: false,
        ),
      ];

      // Act
      await tester.pumpWidget(createTestWidget(requirements: requirements));

      // Check the required consent checkbox
      await tester.tap(find.byType(Checkbox).first);
      await tester.pump();

      // Assert
      final buttonWidget = tester.widget<ElevatedButton>(
        find.byType(ElevatedButton).first,
      );
      expect(buttonWidget.onPressed, isNotNull); // Button should be enabled
    });

    testWidgets('should call onCancel when cancel button is pressed', (tester) async {
      // Arrange
      bool cancelCalled = false;
      final requirements = [
        const ConsentRequirement(
          type: ConsentType.termsOfService,
          title: 'Terms of Service',
          description: 'Accept our terms',
          version: '1.0',
          required: true,
        ),
      ];

      // Act
      await tester.pumpWidget(createTestWidget(
        requirements: requirements,
        onCancel: () => cancelCalled = true,
      ));

      await tester.tap(find.text('Cancel'));
      await tester.pump();

      // Assert
      expect(cancelCalled, true);
    });

    testWidgets('should submit consents when continue button is pressed', (tester) async {
      // Arrange
      bool consentsCalled = false;
      final requirements = [
        const ConsentRequirement(
          type: ConsentType.termsOfService,
          title: 'Terms of Service',
          description: 'Accept our terms',
          version: '1.0',
          required: true,
        ),
      ];

      when(mockConsentService.updateMultipleConsents(
        any,
        any,
        userAgent: anyNamed('userAgent'),
        metadata: anyNamed('metadata'),
      )).thenAnswer((_) async {});

      // Act
      await tester.pumpWidget(createTestWidget(
        requirements: requirements,
        onConsentsGiven: () => consentsCalled = true,
      ));

      // Check the required consent
      await tester.tap(find.byType(Checkbox).first);
      await tester.pump();

      // Tap continue button
      await tester.tap(find.text('Continue'));
      await tester.pump();

      // Assert
      verify(mockConsentService.updateMultipleConsents(
        'test_user',
        {ConsentType.termsOfService: true},
        userAgent: 'Smor-Ting Mobile App',
        metadata: anyNamed('metadata'),
      )).called(1);
    });

    testWidgets('should show loading indicator when submitting consents', (tester) async {
      // Arrange
      final requirements = [
        const ConsentRequirement(
          type: ConsentType.termsOfService,
          title: 'Terms of Service',
          description: 'Accept our terms',
          version: '1.0',
          required: true,
        ),
      ];

      // Make the API call hang to test loading state
      when(mockConsentService.updateMultipleConsents(
        any,
        any,
        userAgent: anyNamed('userAgent'),
        metadata: anyNamed('metadata'),
      )).thenAnswer((_) => Future.delayed(const Duration(seconds: 10)));

      // Act
      await tester.pumpWidget(createTestWidget(requirements: requirements));

      // Check the consent and tap continue
      await tester.tap(find.byType(Checkbox).first);
      await tester.pump();
      await tester.tap(find.text('Continue'));
      await tester.pump();

      // Assert
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      
      // Button should be disabled while loading
      final buttonWidget = tester.widget<ElevatedButton>(
        find.byType(ElevatedButton).first,
      );
      expect(buttonWidget.onPressed, isNull);
    });

    testWidgets('should show error snackbar when consent submission fails', (tester) async {
      // Arrange
      final requirements = [
        const ConsentRequirement(
          type: ConsentType.termsOfService,
          title: 'Terms of Service',
          description: 'Accept our terms',
          version: '1.0',
          required: true,
        ),
      ];

      when(mockConsentService.updateMultipleConsents(
        any,
        any,
        userAgent: anyNamed('userAgent'),
        metadata: anyNamed('metadata'),
      )).thenThrow(Exception('Network error'));

      // Act
      await tester.pumpWidget(createTestWidget(requirements: requirements));

      // Check consent and submit
      await tester.tap(find.byType(Checkbox).first);
      await tester.pump();
      await tester.tap(find.text('Continue'));
      await tester.pump();

      // Assert
      expect(find.text('Failed to save consent: Exception: Network error'), findsOneWidget);
    });
  });

  group('ConsentCheckbox', () {
    testWidgets('should display consent title and handle value changes', (tester) async {
      // Arrange
      const requirement = ConsentRequirement(
        type: ConsentType.termsOfService,
        title: 'Terms of Service',
        description: 'Accept our terms',
        version: '1.0',
        required: true,
      );

      bool value = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: StatefulBuilder(
              builder: (context, setState) {
                return ConsentCheckbox(
                  requirement: requirement,
                  value: value,
                  onChanged: (newValue) {
                    setState(() {
                      value = newValue ?? false;
                    });
                  },
                );
              },
            ),
          ),
        ),
      );

      // Assert initial state
      expect(find.text('Terms of Service'), findsOneWidget);
      expect(find.text('*'), findsOneWidget); // Required indicator
      expect(find.text('Required to use the service'), findsOneWidget);

      // Act - tap checkbox
      await tester.tap(find.byType(Checkbox));
      await tester.pump();

      // Assert checkbox is now checked
      final checkboxWidget = tester.widget<Checkbox>(find.byType(Checkbox));
      expect(checkboxWidget.value, true);
    });

    testWidgets('should not show required indicators for optional consent', (tester) async {
      // Arrange
      const requirement = ConsentRequirement(
        type: ConsentType.marketingCommunications,
        title: 'Marketing Communications',
        description: 'Receive emails',
        version: '1.0',
        required: false,
      );

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ConsentCheckbox(
              requirement: requirement,
              value: false,
              onChanged: (_) {},
            ),
          ),
        ),
      );

      // Assert
      expect(find.text('Marketing Communications'), findsOneWidget);
      expect(find.text('*'), findsNothing); // No required indicator
      expect(find.text('Required to use the service'), findsNothing);
    });
  });
}
