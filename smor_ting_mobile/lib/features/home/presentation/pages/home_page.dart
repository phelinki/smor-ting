import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/models/service.dart';

import '../../../auth/presentation/providers/auth_provider.dart';
import '../../../wallet/presentation/pages/wallet_page.dart';

class HomePage extends ConsumerStatefulWidget {
  const HomePage({super.key});

  @override
  ConsumerState<HomePage> createState() => _HomePageState();
}

class _HomePageState extends ConsumerState<HomePage> {
  int _currentIndex = 0;

  final List<Widget> _pages = [
    const _ServicesPage(),
    const _BookingsPage(),
    const _MessagesPage(),
    const WalletPage(),
    const _ProfilePage(),
  ];

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authNotifierProvider);
    
    if (authState is Authenticated) {
      return Scaffold(
        backgroundColor: AppTheme.white,
        body: SafeArea(
          child: _pages[_currentIndex],
        ),
        bottomNavigationBar: BottomNavigationBar(
          currentIndex: _currentIndex,
          onTap: (index) {
            setState(() {
              _currentIndex = index;
            });
          },
          type: BottomNavigationBarType.fixed,
          selectedItemColor: AppTheme.primaryRed,
          unselectedItemColor: AppTheme.textSecondary,
          items: const [
            BottomNavigationBarItem(
              icon: Icon(Icons.home_outlined),
              activeIcon: Icon(Icons.home),
              label: 'Home',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.calendar_today_outlined),
              activeIcon: Icon(Icons.calendar_today),
              label: 'Bookings',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.message_outlined),
              activeIcon: Icon(Icons.message),
              label: 'Messages',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.account_balance_wallet_outlined),
              activeIcon: Icon(Icons.account_balance_wallet),
              label: 'Wallet',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.person_outline),
              activeIcon: Icon(Icons.person),
              label: 'Profile',
            ),
          ],
        ),
      );
    }
    
    // If not authenticated, show loading or redirect
    return const Scaffold(
      body: Center(
        child: CircularProgressIndicator(),
      ),
    );
  }
}

// Main Services/Home Page
class _ServicesPage extends ConsumerWidget {
  const _ServicesPage();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authNotifierProvider);
    final userName = authState is Authenticated ? authState.user.firstName : 'User';

    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
              // Top Navigation Bar
              Padding(
                padding: const EdgeInsets.fromLTRB(16.0, 16.0, 16.0, 16.0),
                child: Row(
              children: [
                // Profile Icon
                GestureDetector(
                  onTap: () {
                    // TODO: Navigate to profile or show menu
                  },
                  child: Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: AppTheme.secondaryBlue,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: const Icon(
                      Icons.person,
                      color: AppTheme.white,
                      size: 24,
                    ),
                  ),
                ),
                
                const SizedBox(width: 12),
                
                // Greeting
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Good morning,',
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: AppTheme.textSecondary,
                        ),
                      ),
                      Text(
                        userName,
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
                
                // Notification Bell
                IconButton(
                  onPressed: () {
                    // TODO: Navigate to notifications
                  },
                  icon: const Icon(
                    Icons.notifications_outlined,
                    color: AppTheme.secondaryBlue,
                    size: 24,
                  ),
                ),
              ],
            ),
          ),
          
          // Search Bar
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            child: Container(
              decoration: BoxDecoration(
                color: AppTheme.lightGray,
                borderRadius: BorderRadius.circular(12),
              ),
              child: TextField(
                decoration: InputDecoration(
                  hintText: 'Search for services...',
                  hintStyle: TextStyle(
                    color: AppTheme.textSecondary,
                    fontSize: 16,
                  ),
                  prefixIcon: const Icon(
                    Icons.search,
                    color: AppTheme.textSecondary,
                  ),
                  border: InputBorder.none,
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 16,
                  ),
                ),
                onTap: () {
                  // TODO: Navigate to search page
                },
              ),
            ),
          ),
          
          const SizedBox(height: 24),
          
          // Hero Banner
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            child: Container(
              decoration: BoxDecoration(
                color: AppTheme.secondaryBlue,
                borderRadius: BorderRadius.circular(16),
                boxShadow: [
                  BoxShadow(
                    color: AppTheme.secondaryBlue.withOpacity(0.3),
                    blurRadius: 10,
                    offset: const Offset(0, 4),
                  ),
                ],
              ),
              child: Padding(
                padding: const EdgeInsets.all(24.0),
                child: Row(
                  children: [
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(
                            'Need a Service?',
                            style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                              color: AppTheme.white,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Book trusted professionals for your home and business needs',
                            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                              color: AppTheme.white.withOpacity(0.9),
                            ),
                          ),
                          const SizedBox(height: 16),
                          ElevatedButton(
                            onPressed: () {
                              // TODO: Navigate to service categories
                            },
                            style: ElevatedButton.styleFrom(
                              backgroundColor: AppTheme.white,
                              foregroundColor: AppTheme.secondaryBlue,
                              padding: const EdgeInsets.symmetric(
                                horizontal: 20,
                                vertical: 12,
                              ),
                            ),
                            child: const Text('Book Now'),
                          ),
                        ],
                      ),
                    ),
                    const Icon(
                      Icons.handyman,
                      size: 60,
                      color: AppTheme.white,
                    ),
                  ],
                ),
              ),
            ),
          ),
          
          const SizedBox(height: 32),
          
          // Service Categories Section
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            child: Text(
              'Service Categories',
              style: Theme.of(context).textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          
          const SizedBox(height: 16),
          
          // Service Categories Grid
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            child: GridView.builder(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 3,
                crossAxisSpacing: 12,
                mainAxisSpacing: 12,
                childAspectRatio: 1.0,
              ),
              itemCount: _serviceCategories.length,
              itemBuilder: (context, index) {
                final category = _serviceCategories[index];
                return _buildServiceCategoryCard(context, category);
              },
            ),
          ),
          
          const SizedBox(height: 32),
          
          // Add bottom padding to prevent overflow with bottom navigation bar
          const SizedBox(height: 160),
        ],
      ),
    );
  }

  Widget _buildServiceCategoryCard(BuildContext context, ServiceCategory category) {
    return GestureDetector(
      onTap: () {
        context.push('/service-listings/${category.name.toLowerCase()}');
      },
      child: Container(
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
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              _getIconData(category.icon),
              size: 32,
              color: AppTheme.secondaryBlue,
            ),
            const SizedBox(height: 8),
            Text(
              category.name,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                fontWeight: FontWeight.w600,
                color: AppTheme.textPrimary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  IconData _getIconData(String iconName) {
    switch (iconName) {
      case 'plumbing':
        return Icons.plumbing;
      case 'electrical_services':
        return Icons.electrical_services;
      case 'cleaning_services':
        return Icons.cleaning_services;
      case 'handyman':
        return Icons.handyman;
      case 'format_paint':
        return Icons.format_paint;
      case 'ac_unit':
        return Icons.ac_unit;
      case 'grass':
        return Icons.grass;
      case 'local_shipping':
        return Icons.local_shipping;
      case 'security':
        return Icons.security;
      default:
        return Icons.build;
    }
  }
}

// Other page classes
class _BookingsPage extends StatelessWidget {
  const _BookingsPage();

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(
        child: Text('Bookings Page'),
      ),
    );
  }
}

class _MessagesPage extends StatelessWidget {
  const _MessagesPage();

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(
        child: Text('Messages Page'),
      ),
    );
  }
}

class _ProfilePage extends ConsumerWidget {
  const _ProfilePage();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authNotifierProvider);
    
    if (authState is Authenticated) {
      return Scaffold(
        backgroundColor: AppTheme.white,
        appBar: AppBar(
          title: const Text('Profile'),
          backgroundColor: AppTheme.white,
          foregroundColor: AppTheme.textPrimary,
          elevation: 0,
        ),
        body: SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(24.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Profile Header
                Center(
                  child: Column(
                    children: [
                      Container(
                        width: 100,
                        height: 100,
                        decoration: BoxDecoration(
                          color: AppTheme.secondaryBlue,
                          borderRadius: BorderRadius.circular(50),
                        ),
                        child: const Icon(
                          Icons.person,
                          color: AppTheme.white,
                          size: 50,
                        ),
                      ),
                      const SizedBox(height: 16),
                      Text(
                        authState.user.fullName,
                        style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        authState.user.email,
                        style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          color: AppTheme.textSecondary,
                        ),
                      ),
                      const SizedBox(height: 32),
                    ],
                  ),
                ),
                
                // Profile Options
                _buildProfileOption(
                  icon: Icons.person_outline,
                  title: 'Edit Profile',
                  onTap: () {
                    // TODO: Navigate to edit profile
                  },
                ),
                _buildProfileOption(
                  icon: Icons.location_on_outlined,
                  title: 'Address Book',
                  onTap: () {
                    // TODO: Navigate to address book
                  },
                ),
                _buildProfileOption(
                  icon: Icons.settings_outlined,
                  title: 'Settings',
                  onTap: () {
                    // TODO: Navigate to settings
                  },
                ),
                _buildProfileOption(
                  icon: Icons.help_outline,
                  title: 'Help & Support',
                  onTap: () {
                    // TODO: Navigate to help
                  },
                ),
                const SizedBox(height: 32),
                
                // Logout Button
                SizedBox(
                  width: double.infinity,
                  child: TextButton(
                    onPressed: () {
                      ref.read(authNotifierProvider.notifier).logout();
                    },
                    child: const Text(
                      'Logout',
                      style: TextStyle(
                        color: AppTheme.primaryRed,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
                
                // Delete Account Button
                SizedBox(
                  width: double.infinity,
                  child: TextButton(
                    onPressed: () {
                      // TODO: Show delete account confirmation
                    },
                    child: const Text(
                      'Delete Account',
                      style: TextStyle(
                        color: AppTheme.primaryRed,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      );
    }
    
    return const Scaffold(
      body: Center(
        child: CircularProgressIndicator(),
      ),
    );
  }

  Widget _buildProfileOption({
    required IconData icon,
    required String title,
    required VoidCallback onTap,
  }) {
    return ListTile(
      leading: Icon(icon, color: AppTheme.secondaryBlue),
      title: Text(title),
      trailing: const Icon(Icons.arrow_forward_ios, size: 16, color: AppTheme.textSecondary),
      onTap: onTap,
      contentPadding: EdgeInsets.zero,
    );
  }
}

// Sample data
final List<ServiceCategory> _serviceCategories = [
  ServiceCategory(
    id: 'plumbing',
    name: 'Plumbing',
    description: 'Plumbing services',
    icon: 'plumbing',
    color: '#2196F3',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
  ServiceCategory(
    id: 'electrical',
    name: 'Electrical',
    description: 'Electrical services',
    icon: 'electrical_services',
    color: '#FF9800',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
  ServiceCategory(
    id: 'cleaning',
    name: 'Cleaning',
    description: 'Cleaning services',
    icon: 'cleaning_services',
    color: '#4CAF50',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
  ServiceCategory(
    id: 'carpentry',
    name: 'Carpentry',
    description: 'Carpentry services',
    icon: 'handyman',
    color: '#795548',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
  ServiceCategory(
    id: 'painting',
    name: 'Painting',
    description: 'Painting services',
    icon: 'format_paint',
    color: '#9C27B0',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
  ServiceCategory(
    id: 'hvac',
    name: 'HVAC',
    description: 'HVAC services',
    icon: 'ac_unit',
    color: '#00BCD4',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
  ServiceCategory(
    id: 'gardening',
    name: 'Gardening',
    description: 'Gardening services',
    icon: 'grass',
    color: '#8BC34A',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
  ServiceCategory(
    id: 'moving',
    name: 'Moving',
    description: 'Moving services',
    icon: 'local_shipping',
    color: '#FF5722',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
  ServiceCategory(
    id: 'security',
    name: 'Security',
    description: 'Security services',
    icon: 'security',
    color: '#607D8B',
    isActive: true,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  ),
]; 