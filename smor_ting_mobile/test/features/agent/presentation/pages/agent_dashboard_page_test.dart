import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/features/agent/presentation/pages/agent_dashboard_page.dart';
import 'package:smor_ting_mobile/core/theme/app_theme.dart';

void main() {
  group('AgentDashboardPage', () {
    testWidgets('should display greeting header', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      expect(find.textContaining('Good'), findsOneWidget);
      expect(find.textContaining('John!'), findsOneWidget);
      expect(find.text('Here\'s what\'s happening with your business today.'), findsOneWidget);
    });

    testWidgets('should display earnings overview section', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      expect(find.text('Earnings Overview'), findsOneWidget);
      expect(find.text('Total Earnings'), findsOneWidget);
      expect(find.text('Pending Payout'), findsOneWidget);
      expect(find.text('\$1250.00'), findsOneWidget);
      expect(find.text('\$350.00'), findsOneWidget);
      expect(find.text('View Full Report'), findsOneWidget);
    });

    testWidgets('should display rating card', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      expect(find.text('Your Rating'), findsOneWidget);
      expect(find.text('4.8 (156 reviews)'), findsOneWidget);
      expect(find.byIcon(Icons.star), findsWidgets);
    });

    testWidgets('should display job statistics', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      expect(find.text('Job Statistics'), findsOneWidget);
      expect(find.text('Total Jobs'), findsOneWidget);
      expect(find.text('Completed'), findsOneWidget);
      expect(find.text('Cancelled'), findsOneWidget);
      expect(find.text('45'), findsOneWidget);
      expect(find.text('38'), findsOneWidget);
      expect(find.text('2'), findsOneWidget);
    });

    testWidgets('should display recent job requests', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      expect(find.text('Recent Job Requests'), findsOneWidget);
      expect(find.text('View All'), findsOneWidget);
      expect(find.text('Plumbing Repair'), findsOneWidget);
      expect(find.text('Emergency plumbing repair needed'), findsOneWidget);
      expect(find.text('Customer: Sarah Johnson'), findsOneWidget);
      expect(find.text('123 Main St, Monrovia, Liberia'), findsOneWidget);
      expect(find.text('\$75.00'), findsOneWidget);
    });

    testWidgets('should display job action buttons', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      expect(find.text('Decline'), findsOneWidget);
      expect(find.text('Accept'), findsOneWidget);
    });

    testWidgets('should have refresh functionality', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      expect(find.byType(RefreshIndicator), findsOneWidget);
    });

    testWidgets('should display correct icons', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      expect(find.byIcon(Icons.work), findsOneWidget);
      expect(find.byIcon(Icons.check_circle), findsOneWidget);
      expect(find.byIcon(Icons.cancel), findsOneWidget);
      expect(find.byIcon(Icons.location_on), findsOneWidget);
    });

    testWidgets('should have proper styling and layout', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      // Verify the page has proper structure
      expect(find.byType(Scaffold), findsOneWidget);
      expect(find.byType(SingleChildScrollView), findsOneWidget);
      expect(find.byType(Column), findsWidgets);
    });

    testWidgets('should display job amount with correct styling', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      final amountText = tester.widget<Text>(find.text('\$75.00'));
      expect(amountText.style?.color, equals(AppTheme.warning));
      expect(amountText.style?.fontWeight, equals(FontWeight.w600));
    });

    testWidgets('should display earnings with correct styling', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentDashboardPage(),
          ),
        ),
      );

      final totalEarningsText = tester.widget<Text>(find.text('\$1250.00'));
      expect(totalEarningsText.style?.color, equals(AppTheme.successGreen));
      expect(totalEarningsText.style?.fontWeight, equals(FontWeight.bold));
    });
  });
}
