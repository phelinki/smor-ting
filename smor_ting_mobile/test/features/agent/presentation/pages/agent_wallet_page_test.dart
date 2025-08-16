import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/features/agent/presentation/pages/agent_wallet_page.dart';

void main() {
  group('AgentWalletPage', () {
    testWidgets('should display wallet page title', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentWalletPage(),
          ),
        ),
      );

      expect(find.text('Wallet'), findsOneWidget);
    });

    testWidgets('should display balance information', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentWalletPage(),
          ),
        ),
      );

      expect(find.text('Available Balance'), findsOneWidget);
      expect(find.text('Pending Payout'), findsOneWidget);
    });

    testWidgets('should display action buttons', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentWalletPage(),
          ),
        ),
      );

      expect(find.text('Withdraw'), findsOneWidget);
      expect(find.text('Transaction History'), findsOneWidget);
    });

    testWidgets('should display recent transactions', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentWalletPage(),
          ),
        ),
      );

      expect(find.text('Recent Transactions'), findsOneWidget);
    });

    testWidgets('should display empty state when no transactions', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentWalletPage(),
          ),
        ),
      );

      expect(find.text('No transactions yet'), findsOneWidget);
      expect(find.text('Your transaction history will appear here.'), findsOneWidget);
    });

    testWidgets('should have refresh functionality', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentWalletPage(),
          ),
        ),
      );

      expect(find.byType(RefreshIndicator), findsOneWidget);
    });
  });
}
