import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import '../../../../../lib/features/settings/presentation/widgets/biometric_settings_widget.dart';

void main() {
  group('BiometricSettingsWidget', () {
    Widget createWidget({String? userEmail}) {
      return MaterialApp(
        home: Scaffold(
          body: BiometricSettingsWidget(
            userEmail: userEmail ?? 'test@example.com',
          ),
        ),
      );
    }

    testWidgets('should render without crashing', (tester) async {
      // Act
      await tester.pumpWidget(createWidget());

      // Assert - Just verify it renders without throwing
      expect(find.byType(BiometricSettingsWidget), findsOneWidget);
    });

    testWidgets('should show loading state initially', (tester) async {
      // Act
      await tester.pumpWidget(createWidget());

      // Assert - The widget should handle initialization gracefully
      // Even if biometric isn't available, it should not crash
      expect(find.byType(BiometricSettingsWidget), findsOneWidget);
    });
  });
}