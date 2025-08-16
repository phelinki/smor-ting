import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/theme/app_theme.dart';
import 'agent_dashboard_page.dart';
import 'agent_jobs_page.dart';
import 'agent_messages_page.dart';
import 'agent_wallet_page.dart';
import 'agent_profile_page.dart';

class AgentMainPage extends ConsumerStatefulWidget {
  const AgentMainPage({super.key});

  @override
  ConsumerState<AgentMainPage> createState() => _AgentMainPageState();
}

class _AgentMainPageState extends ConsumerState<AgentMainPage> {
  int _currentIndex = 0;

  final List<Widget> _pages = [
    const AgentDashboardPage(),
    const AgentJobsPage(),
    const AgentMessagesPage(),
    const AgentWalletPage(),
    const AgentProfilePage(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[_currentIndex],
      bottomNavigationBar: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        currentIndex: _currentIndex,
        onTap: (index) {
          setState(() {
            _currentIndex = index;
          });
        },
        selectedItemColor: AppTheme.primaryRed,
        unselectedItemColor: AppTheme.textSecondary,
        backgroundColor: AppTheme.white,
        elevation: 8,
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.dashboard),
            label: 'Dashboard',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.calendar_today),
            label: 'Jobs',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.message),
            label: 'Messages',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.account_balance_wallet),
            label: 'Wallet',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.person),
            label: 'Profile',
          ),
        ],
      ),
    );
  }
}
