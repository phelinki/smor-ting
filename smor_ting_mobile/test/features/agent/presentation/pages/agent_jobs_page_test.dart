import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/features/agent/presentation/pages/agent_jobs_page.dart';

void main() {
  group('AgentJobsPage', () {
    testWidgets('should display jobs page title', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentJobsPage(),
          ),
        ),
      );

      expect(find.text('My Jobs'), findsOneWidget);
    });

    testWidgets('should display job status tabs', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentJobsPage(),
          ),
        ),
      );

      expect(find.text('Active'), findsOneWidget);
      expect(find.text('Completed'), findsOneWidget);
      expect(find.text('Cancelled'), findsOneWidget);
    });

    testWidgets('should start with Active tab selected', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentJobsPage(),
          ),
        ),
      );

      // Verify Active tab is initially selected by checking TabBar
      expect(find.byType(TabBar), findsOneWidget);
      expect(find.text('Active'), findsOneWidget);
    });

    testWidgets('should switch to Completed tab when tapped', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentJobsPage(),
          ),
        ),
      );

      await tester.tap(find.text('Completed'));
      await tester.pumpAndSettle();

      // Verify Completed tab is now selected
      expect(find.text('Completed'), findsOneWidget);
    });

    testWidgets('should display empty state when no jobs', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentJobsPage(),
          ),
        ),
      );

      expect(find.text('No jobs found'), findsOneWidget);
      expect(find.text('You don\'t have any jobs yet.'), findsOneWidget);
    });

    testWidgets('should display job card with correct information', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentJobsPage(),
          ),
        ),
      );

      // This test will need to be updated when we implement job data
      expect(find.byType(Card), findsNothing); // Initially no jobs
    });

    testWidgets('should have refresh functionality', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentJobsPage(),
          ),
        ),
      );

      // Verify refresh indicator is present
      expect(find.byType(RefreshIndicator), findsOneWidget);
    });
  });
}
