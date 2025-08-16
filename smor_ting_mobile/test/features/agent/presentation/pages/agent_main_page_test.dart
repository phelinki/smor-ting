import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/features/agent/presentation/pages/agent_main_page.dart';
import 'package:smor_ting_mobile/core/theme/app_theme.dart';

void main() {
  group('AgentMainPage', () {
    testWidgets('should display bottom navigation with 5 items', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      // Verify bottom navigation bar exists
      expect(find.byType(BottomNavigationBar), findsOneWidget);
      
      // Verify all 5 navigation items are present
      expect(find.text('Dashboard'), findsOneWidget);
      expect(find.text('Jobs'), findsOneWidget);
      expect(find.text('Messages'), findsOneWidget);
      expect(find.text('Wallet'), findsOneWidget);
      expect(find.text('Profile'), findsOneWidget);
    });

    testWidgets('should start with Dashboard as selected tab', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      // Verify Dashboard is initially selected
      final bottomNav = tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar));
      expect(bottomNav.currentIndex, equals(0));
    });

    testWidgets('should switch to Jobs tab when tapped', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      // Tap on Jobs tab
      await tester.tap(find.text('Jobs'));
      await tester.pumpAndSettle();

      // Verify Jobs tab is now selected
      final bottomNav = tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar));
      expect(bottomNav.currentIndex, equals(1));
    });

    testWidgets('should switch to Messages tab when tapped', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      // Tap on Messages tab
      await tester.tap(find.text('Messages'));
      await tester.pumpAndSettle();

      // Verify Messages tab is now selected
      final bottomNav = tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar));
      expect(bottomNav.currentIndex, equals(2));
    });

    testWidgets('should switch to Wallet tab when tapped', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      // Tap on Wallet tab
      await tester.tap(find.text('Wallet'));
      await tester.pumpAndSettle();

      // Verify Wallet tab is now selected
      final bottomNav = tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar));
      expect(bottomNav.currentIndex, equals(3));
    });

    testWidgets('should switch to Profile tab when tapped', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      // Tap on Profile tab
      await tester.tap(find.text('Profile'));
      await tester.pumpAndSettle();

      // Verify Profile tab is now selected
      final bottomNav = tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar));
      expect(bottomNav.currentIndex, equals(4));
    });

    testWidgets('should display correct icons for each tab', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      // Verify icons are present (using findAll to handle multiple instances)
      expect(find.byIcon(Icons.dashboard), findsWidgets);
      expect(find.byIcon(Icons.calendar_today), findsWidgets);
      expect(find.byIcon(Icons.message), findsWidgets);
      expect(find.byIcon(Icons.account_balance_wallet), findsWidgets);
      expect(find.byIcon(Icons.person), findsWidgets);
    });

    testWidgets('should have correct navigation bar styling', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      final bottomNav = tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar));
      
      // Verify styling properties
      expect(bottomNav.type, equals(BottomNavigationBarType.fixed));
      expect(bottomNav.selectedItemColor, equals(AppTheme.primaryRed));
      expect(bottomNav.unselectedItemColor, equals(AppTheme.textSecondary));
      expect(bottomNav.backgroundColor, equals(AppTheme.white));
      expect(bottomNav.elevation, equals(8));
    });

    testWidgets('should maintain state when switching between tabs', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMainPage(),
          ),
        ),
      );

      // Switch to Jobs tab
      await tester.tap(find.text('Jobs'));
      await tester.pumpAndSettle();
      expect(tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar)).currentIndex, equals(1));

      // Switch back to Dashboard
      await tester.tap(find.text('Dashboard'));
      await tester.pumpAndSettle();
      expect(tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar)).currentIndex, equals(0));

      // Switch to Messages
      await tester.tap(find.text('Messages'));
      await tester.pumpAndSettle();
      expect(tester.widget<BottomNavigationBar>(find.byType(BottomNavigationBar)).currentIndex, equals(2));
    });
  });
}
