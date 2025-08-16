import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/theme/app_theme.dart';

class AgentJobsPage extends ConsumerStatefulWidget {
  const AgentJobsPage({super.key});

  @override
  ConsumerState<AgentJobsPage> createState() => _AgentJobsPageState();
}

class _AgentJobsPageState extends ConsumerState<AgentJobsPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  String _selectedStatus = 'Active';

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.lightGray,
      appBar: AppBar(
        title: const Text(
          'My Jobs',
          style: TextStyle(
            fontWeight: FontWeight.w600,
            color: AppTheme.textPrimary,
          ),
        ),
        backgroundColor: AppTheme.white,
        elevation: 0,
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppTheme.primaryRed,
          unselectedLabelColor: AppTheme.textSecondary,
          indicatorColor: AppTheme.primaryRed,
          onTap: (index) {
            setState(() {
              _selectedStatus = ['Active', 'Completed', 'Cancelled'][index];
            });
          },
          tabs: const [
            Tab(text: 'Active'),
            Tab(text: 'Completed'),
            Tab(text: 'Cancelled'),
          ],
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          // TODO: Implement refresh logic
          await Future.delayed(const Duration(seconds: 1));
        },
        child: TabBarView(
          controller: _tabController,
          children: [
            _buildJobList('Active'),
            _buildJobList('Completed'),
            _buildJobList('Cancelled'),
          ],
        ),
      ),
    );
  }

  Widget _buildJobList(String status) {
    // TODO: Implement actual job data fetching
    final jobs = <Map<String, dynamic>>[]; // Empty for now

    if (jobs.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.work_outline,
              size: 64,
              color: AppTheme.textSecondary,
            ),
            const SizedBox(height: 16),
            Text(
              'No jobs found',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                color: AppTheme.textPrimary,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'You don\'t have any jobs yet.',
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: AppTheme.textSecondary,
              ),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: jobs.length,
      itemBuilder: (context, index) {
        final job = jobs[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: ListTile(
            title: Text(job['title'] ?? ''),
            subtitle: Text(job['description'] ?? ''),
            trailing: Text(
              '\$${job['amount']?.toStringAsFixed(2) ?? '0.00'}',
              style: Theme.of(context).textTheme.titleSmall?.copyWith(
                fontWeight: FontWeight.bold,
                color: AppTheme.successGreen,
              ),
            ),
          ),
        );
      },
    );
  }
}
