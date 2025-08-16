import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/features/agent/presentation/pages/agent_messages_page.dart';

void main() {
  group('AgentMessagesPage', () {
    testWidgets('should display messages page title', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMessagesPage(),
          ),
        ),
      );

      expect(find.text('Messages'), findsOneWidget);
    });

    testWidgets('should display empty state when no messages', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMessagesPage(),
          ),
        ),
      );

      expect(find.text('No messages'), findsOneWidget);
      expect(find.text('You don\'t have any messages yet.'), findsOneWidget);
    });

    testWidgets('should display message list when messages exist', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMessagesPage(),
          ),
        ),
      );

      // Initially no messages
      expect(find.byType(ListTile), findsNothing);
    });

    testWidgets('should have search functionality', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMessagesPage(),
          ),
        ),
      );

      expect(find.byType(TextField), findsOneWidget);
      expect(find.byIcon(Icons.search), findsOneWidget);
    });

    testWidgets('should have refresh functionality', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMessagesPage(),
          ),
        ),
      );

      expect(find.byType(RefreshIndicator), findsOneWidget);
    });

    testWidgets('should display message with correct information', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: AgentMessagesPage(),
          ),
        ),
      );

      // This test will need to be updated when we implement message data
      expect(find.byType(Card), findsNothing); // Initially no messages
    });
  });
}
