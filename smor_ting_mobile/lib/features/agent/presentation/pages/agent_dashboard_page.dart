import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../auth/presentation/providers/auth_provider.dart';

class Job {
  final String id;
  final String customerName;
  final String serviceName;
  final DateTime dateTime;
  final String address;
  final double amount;
  final String status;

  Job({
    required this.id,
    required this.customerName,
    required this.serviceName,
    required this.dateTime,
    required this.address,
    required this.amount,
    required this.status,
  });
}

class AgentDashboardPage extends ConsumerStatefulWidget {
  const AgentDashboardPage({super.key});

  @override
  ConsumerState<AgentDashboardPage> createState() => _AgentDashboardPageState();
}

class _AgentDashboardPageState extends ConsumerState<AgentDashboardPage> {
  // Mock data - in real app this would come from API
  final double totalEarnings = 2850.75;
  final double pendingPayout = 425.50;
  final double rating = 4.8;
  final int reviewCount = 127;
  
  final List<Job> todaysJobs = [
    Job(
      id: '1',
      customerName: 'Sarah Johnson',
      serviceName: 'Plumbing Repair',
      dateTime: DateTime.now().add(const Duration(hours: 2)),
      address: '123 Cantonments Rd, Accra',
      amount: 150.00,
      status: 'pending',
    ),
    Job(
      id: '2',
      customerName: 'Michael Asante',
      serviceName: 'Electrical Installation',
      dateTime: DateTime.now().add(const Duration(hours: 4)),
      address: '456 East Legon, Accra',
      amount: 275.00,
      status: 'pending',
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authNotifierProvider);
    final agentName = authState is Authenticated ? authState.user.firstName : 'Agent';

    return Scaffold(
      backgroundColor: AppTheme.lightGray,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header
              _buildHeader(agentName),
              const SizedBox(height: 24),
              
              // Stats Cards
              _buildStatsCards(),
              const SizedBox(height: 24),
              
              // Quick Actions
              _buildQuickActions(),
              const SizedBox(height: 24),
              
              // Today's Jobs
              _buildTodaysJobs(),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader(String agentName) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Good ${_getGreeting()},',
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: AppTheme.textSecondary,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              agentName,
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
                color: AppTheme.textPrimary,
              ),
            ),
            const SizedBox(height: 8),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(
                color: AppTheme.secondaryBlue.withOpacity(0.1),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Text(
                'Today',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: AppTheme.secondaryBlue,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
          ],
        ),
        Row(
          children: [
            IconButton(
              onPressed: () => _showNotifications(),
              icon: Stack(
                children: [
                  const Icon(
                    Icons.notifications_outlined,
                    color: AppTheme.textSecondary,
                    size: 24,
                  ),
                  Positioned(
                    right: 0,
                    top: 0,
                    child: Container(
                      width: 8,
                      height: 8,
                      decoration: const BoxDecoration(
                        color: AppTheme.primaryRed,
                        shape: BoxShape.circle,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            IconButton(
              onPressed: () => _showProfileMenu(),
              icon: Container(
                width: 32,
                height: 32,
                decoration: BoxDecoration(
                  color: AppTheme.secondaryBlue,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: const Icon(
                  Icons.person,
                  color: AppTheme.white,
                  size: 20,
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildStatsCards() {
    return Column(
      children: [
        // Earnings Overview
        Container(
          width: double.infinity,
          padding: const EdgeInsets.all(20),
          decoration: BoxDecoration(
            color: AppTheme.white,
            borderRadius: BorderRadius.circular(16),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.05),
                blurRadius: 10,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Earnings Overview',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                  color: AppTheme.textPrimary,
                ),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Total Earnings',
                          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: AppTheme.textSecondary,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          'GH₵${totalEarnings.toStringAsFixed(2)}',
                          style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                            fontWeight: FontWeight.bold,
                            color: AppTheme.successGreen,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Container(
                    width: 1,
                    height: 40,
                    color: AppTheme.lightGray,
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Pending Payout',
                          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: AppTheme.textSecondary,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          'GH₵${pendingPayout.toStringAsFixed(2)}',
                          style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                            fontWeight: FontWeight.bold,
                            color: AppTheme.textSecondary,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              SizedBox(
                width: double.infinity,
                child: OutlinedButton(
                  onPressed: () => _viewFullReport(),
                  style: OutlinedButton.styleFrom(
                    foregroundColor: AppTheme.secondaryBlue,
                    side: const BorderSide(color: AppTheme.secondaryBlue),
                    padding: const EdgeInsets.symmetric(vertical: 12),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text('View Full Report'),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 16),
        
        // Rating Card
        Container(
          width: double.infinity,
          padding: const EdgeInsets.all(20),
          decoration: BoxDecoration(
            color: AppTheme.white,
            borderRadius: BorderRadius.circular(16),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.05),
                blurRadius: 10,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Your Rating',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                  color: AppTheme.textPrimary,
                ),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Row(
                    children: List.generate(5, (index) {
                      return Icon(
                        index < rating.floor() ? Icons.star : Icons.star_border,
                        color: AppTheme.secondaryBlue,
                        size: 24,
                      );
                    }),
                  ),
                  const SizedBox(width: 12),
                  Text(
                    '$rating',
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                      color: AppTheme.secondaryBlue,
                    ),
                  ),
                  const SizedBox(width: 8),
                  Text(
                    'based on $reviewCount reviews',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: AppTheme.textSecondary,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildQuickActions() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Quick Actions',
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.w600,
            color: AppTheme.textPrimary,
          ),
        ),
        const SizedBox(height: 16),
        Row(
          children: [
            Expanded(
              child: ElevatedButton.icon(
                onPressed: () => _setAvailability(),
                icon: const Icon(Icons.schedule, size: 20),
                label: const Text('Set Availability'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppTheme.secondaryBlue,
                  foregroundColor: AppTheme.white,
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: ElevatedButton.icon(
                onPressed: () => _withdrawEarnings(),
                icon: const Icon(Icons.account_balance_wallet, size: 20),
                label: const Text('Withdraw Earnings'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppTheme.primaryRed,
                  foregroundColor: AppTheme.white,
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildTodaysJobs() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Today\'s Jobs',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
                color: AppTheme.textPrimary,
              ),
            ),
            TextButton(
              onPressed: () => _viewAllJobs(),
              child: const Text('View All Jobs'),
            ),
          ],
        ),
        const SizedBox(height: 16),
        ...todaysJobs.take(2).map((job) => _buildJobCard(job)),
      ],
    );
  }

  Widget _buildJobCard(Job job) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                _formatTime(job.dateTime),
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                  color: AppTheme.secondaryBlue,
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: AppTheme.warning.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  'GH₵${job.amount.toStringAsFixed(2)}',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: AppTheme.warning,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            job.serviceName,
            style: Theme.of(context).textTheme.titleSmall?.copyWith(
              fontWeight: FontWeight.w600,
              color: AppTheme.textPrimary,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            'Customer: ${job.customerName}',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: AppTheme.textSecondary,
            ),
          ),
          const SizedBox(height: 4),
          Row(
            children: [
              const Icon(
                Icons.location_on,
                size: 16,
                color: AppTheme.textSecondary,
              ),
              const SizedBox(width: 4),
              Expanded(
                child: Text(
                  job.address,
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: AppTheme.textSecondary,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: () => _declineJob(job.id),
                  style: OutlinedButton.styleFrom(
                    foregroundColor: AppTheme.primaryRed,
                    side: const BorderSide(color: AppTheme.primaryRed),
                    padding: const EdgeInsets.symmetric(vertical: 12),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text('Decline'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: ElevatedButton(
                  onPressed: () => _acceptJob(job.id),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppTheme.successGreen,
                    foregroundColor: AppTheme.white,
                    padding: const EdgeInsets.symmetric(vertical: 12),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text('Accept'),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _getGreeting() {
    final hour = DateTime.now().hour;
    if (hour < 12) return 'Morning';
    if (hour < 17) return 'Afternoon';
    return 'Evening';
  }

  String _formatTime(DateTime dateTime) {
    final hour = dateTime.hour;
    final minute = dateTime.minute.toString().padLeft(2, '0');
    final period = hour >= 12 ? 'PM' : 'AM';
    final displayHour = hour > 12 ? hour - 12 : (hour == 0 ? 12 : hour);
    return '$displayHour:$minute $period';
  }

  // Action methods
  void _showNotifications() {
    // TODO: Implement notifications
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Notifications feature coming soon!')),
    );
  }

  void _showProfileMenu() {
    showModalBottomSheet(
      context: context,
      builder: (context) => Container(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const Icon(Icons.person),
              title: const Text('Profile'),
              onTap: () {
                Navigator.pop(context);
                // Navigate to profile
              },
            ),
            ListTile(
              leading: const Icon(Icons.settings),
              title: const Text('Settings'),
              onTap: () {
                Navigator.pop(context);
                // Navigate to settings
              },
            ),
            ListTile(
              leading: const Icon(Icons.logout, color: AppTheme.primaryRed),
              title: const Text('Logout', style: TextStyle(color: AppTheme.primaryRed)),
              onTap: () {
                Navigator.pop(context);
                ref.read(authNotifierProvider.notifier).logout();
              },
            ),
          ],
        ),
      ),
    );
  }

  void _viewFullReport() {
    // TODO: Navigate to earnings report
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Earnings report feature coming soon!')),
    );
  }

  void _setAvailability() {
    // TODO: Show availability settings
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Set availability feature coming soon!')),
    );
  }

  void _withdrawEarnings() {
    // TODO: Show withdrawal options
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Withdraw earnings feature coming soon!')),
    );
  }

  void _viewAllJobs() {
    // TODO: Navigate to all jobs
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('View all jobs feature coming soon!')),
    );
  }

  void _acceptJob(String jobId) {
    // TODO: Accept job API call
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('Job $jobId accepted!'),
        backgroundColor: AppTheme.successGreen,
      ),
    );
  }

  void _declineJob(String jobId) {
    // TODO: Decline job API call
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('Job $jobId declined!'),
        backgroundColor: AppTheme.primaryRed,
      ),
    );
  }
}
