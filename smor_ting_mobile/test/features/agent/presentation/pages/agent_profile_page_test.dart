import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/features/agent/presentation/pages/agent_profile_page.dart';

void main() {
  group('AgentProfilePage', () {
    testWidgets('should display profile page title', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentProfilePage(),
          ),
        ),
      );

      expect(find.text('Profile'), findsOneWidget);
    });

    testWidgets('should display profile information sections', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentProfilePage(),
          ),
        ),
      );

      expect(find.text('Personal Information'), findsOneWidget);
      expect(find.text('Account Settings'), findsOneWidget);
      expect(find.text('Support'), findsOneWidget);
    });

    testWidgets('should display profile action items', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentProfilePage(),
          ),
        ),
      );

      expect(find.text('Edit Profile'), findsOneWidget);
      expect(find.text('Change Password'), findsOneWidget);
      expect(find.text('Notifications'), findsOneWidget);
      expect(find.text('Privacy'), findsOneWidget);
      expect(find.text('Help & Support'), findsOneWidget);
      expect(find.text('About'), findsOneWidget);
      expect(find.text('Logout'), findsOneWidget);
    });

    testWidgets('should display agent statistics', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentProfilePage(),
          ),
        ),
      );

      expect(find.text('Total Jobs'), findsOneWidget);
      expect(find.text('Rating'), findsOneWidget);
      expect(find.text('Member Since'), findsOneWidget);
    });

    testWidgets('should have logout functionality', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentProfilePage(),
          ),
        ),
      );

      expect(find.text('Logout'), findsOneWidget);
      expect(find.byIcon(Icons.logout), findsOneWidget);
    });

    testWidgets('should display profile avatar', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentProfilePage(),
          ),
        ),
      );

      expect(find.byType(CircleAvatar), findsOneWidget);
    });
  });
}
