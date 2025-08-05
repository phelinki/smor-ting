import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';

class ServiceListingsPage extends ConsumerStatefulWidget {
  final String category;
  
  const ServiceListingsPage({
    super.key,
    required this.category,
  });

  @override
  ConsumerState<ServiceListingsPage> createState() => _ServiceListingsPageState();
}

class _ServiceListingsPageState extends ConsumerState<ServiceListingsPage> {
  String _sortBy = 'rating';
  bool _showFilters = false;

  @override
  Widget build(BuildContext context) {
    final categoryName = widget.category.replaceFirst(widget.category[0], widget.category[0].toUpperCase());
    
    return Scaffold(
      backgroundColor: AppTheme.white,
      appBar: AppBar(
        backgroundColor: AppTheme.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios, color: AppTheme.secondaryBlue),
          onPressed: () => context.pop(),
        ),
        title: Text(
          '$categoryName Services',
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.bold,
            color: AppTheme.textPrimary,
          ),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.tune, color: AppTheme.secondaryBlue),
            onPressed: () {
              setState(() {
                _showFilters = !_showFilters;
              });
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Search and Filter Bar
          if (_showFilters) _buildFilterSection(),
          
          // Sort Options
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: const BoxDecoration(
              color: AppTheme.lightGray,
              border: Border(
                bottom: BorderSide(color: AppTheme.textSecondary, width: 0.5),
              ),
            ),
            child: Row(
              children: [
                const Text(
                  'Sort by:',
                  style: TextStyle(
                    color: AppTheme.textSecondary,
                    fontSize: 14,
                  ),
                ),
                const SizedBox(width: 8),
                DropdownButton<String>(
                  value: _sortBy,
                  underline: const SizedBox(),
                  style: const TextStyle(
                    color: AppTheme.secondaryBlue,
                    fontWeight: FontWeight.w600,
                  ),
                  items: const [
                    DropdownMenuItem(value: 'rating', child: Text('Rating')),
                    DropdownMenuItem(value: 'price', child: Text('Price')),
                    DropdownMenuItem(value: 'distance', child: Text('Distance')),
                    DropdownMenuItem(value: 'reviews', child: Text('Reviews')),
                  ],
                  onChanged: (value) {
                    if (value != null) {
                      setState(() {
                        _sortBy = value;
                      });
                    }
                  },
                ),
                const Spacer(),
                Text(
                  '${_sampleProviders.length} providers found',
                  style: const TextStyle(
                    color: AppTheme.textSecondary,
                    fontSize: 14,
                  ),
                ),
              ],
            ),
          ),
          
          // Provider List
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _sampleProviders.length,
              itemBuilder: (context, index) {
                final provider = _sampleProviders[index];
                return _buildProviderCard(provider);
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFilterSection() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: const BoxDecoration(
        color: AppTheme.lightGray,
        border: Border(
          bottom: BorderSide(color: AppTheme.textSecondary, width: 0.5),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Filters',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: AppTheme.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          
          // Price Range
          Row(
            children: [
              const Text('Price Range: '),
              Expanded(
                child: RangeSlider(
                  values: const RangeValues(25, 150),
                  min: 0,
                  max: 300,
                  divisions: 12,
                  activeColor: AppTheme.secondaryBlue,
                  labels: const RangeLabels('\$25', '\$150'),
                  onChanged: (values) {
                    // TODO: Update price filter
                  },
                ),
              ),
            ],
          ),
          
          // Availability
          Row(
            children: [
              const Text('Available today'),
              const Spacer(),
              Switch(
                value: true,
                activeColor: AppTheme.successGreen,
                onChanged: (value) {
                  // TODO: Update availability filter
                },
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildProviderCard(ServiceProvider provider) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: AppTheme.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.08),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            // Provider Image
            Container(
              width: 60,
              height: 60,
              decoration: BoxDecoration(
                color: AppTheme.lightGray,
                borderRadius: BorderRadius.circular(30),
              ),
              child: Icon(
                Icons.person,
                color: AppTheme.secondaryBlue,
                size: 30,
              ),
            ),
            
            const SizedBox(width: 16),
            
            // Provider Details
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Name and Verification
                  Row(
                    children: [
                      Text(
                        provider.name,
                        style: Theme.of(context).textTheme.titleSmall?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: AppTheme.textPrimary,
                        ),
                      ),
                      if (provider.isVerified) ...[
                        const SizedBox(width: 4),
                        const Icon(
                          Icons.verified,
                          size: 16,
                          color: AppTheme.secondaryBlue,
                        ),
                      ],
                    ],
                  ),
                  
                  const SizedBox(height: 4),
                  
                  // Rating and Reviews
                  Row(
                    children: [
                      const Icon(
                        Icons.star,
                        size: 16,
                        color: AppTheme.warning,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        provider.rating.toString(),
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: AppTheme.secondaryBlue,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '(${provider.reviewCount} reviews)',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: AppTheme.textSecondary,
                        ),
                      ),
                    ],
                  ),
                  
                  const SizedBox(height: 4),
                  
                  // Services Offered
                  Text(
                    provider.services.join(', '),
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: AppTheme.textSecondary,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  
                  const SizedBox(height: 8),
                  
                  // Price and Availability
                  Row(
                    children: [
                      Text(
                        'Starts at \$${provider.startingPrice}',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: AppTheme.secondaryBlue,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const Spacer(),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: provider.isAvailable 
                              ? AppTheme.successGreen.withValues(alpha: 0.1)
                              : AppTheme.textSecondary.withValues(alpha: 0.1),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          provider.isAvailable ? 'Available' : 'Busy',
                          style: TextStyle(
                            color: provider.isAvailable 
                                ? AppTheme.successGreen 
                                : AppTheme.textSecondary,
                            fontSize: 12,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            
            const SizedBox(width: 16),
            
            // Book Button
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              decoration: BoxDecoration(
                color: AppTheme.primaryRed,
                borderRadius: BorderRadius.circular(12),
              ),
              child: GestureDetector(
                onTap: () {
                  context.push('/provider-profile/${provider.id}');
                },
                child: Text(
                  'Book',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: AppTheme.white,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// Data class for service providers
class ServiceProvider {
  final String id;
  final String name;
  final double rating;
  final int reviewCount;
  final List<String> services;
  final int startingPrice;
  final bool isVerified;
  final bool isAvailable;

  ServiceProvider({
    required this.id,
    required this.name,
    required this.rating,
    required this.reviewCount,
    required this.services,
    required this.startingPrice,
    required this.isVerified,
    required this.isAvailable,
  });
}

// Sample data
final List<ServiceProvider> _sampleProviders = [
  ServiceProvider(
    id: '1',
    name: 'John Martinez',
    rating: 4.8,
    reviewCount: 127,
    services: ['Kitchen Plumbing', 'Bathroom Repair', 'Pipe Installation'],
    startingPrice: 75,
    isVerified: true,
    isAvailable: true,
  ),
  ServiceProvider(
    id: '2',
    name: 'Sarah Johnson',
    rating: 4.9,
    reviewCount: 89,
    services: ['Emergency Plumbing', 'Water Heater Repair'],
    startingPrice: 85,
    isVerified: true,
    isAvailable: false,
  ),
  ServiceProvider(
    id: '3',
    name: 'Mike Chen',
    rating: 4.7,
    reviewCount: 156,
    services: ['General Plumbing', 'Drain Cleaning', 'Fixture Installation'],
    startingPrice: 65,
    isVerified: true,
    isAvailable: true,
  ),
  ServiceProvider(
    id: '4',
    name: 'Lisa Thompson',
    rating: 4.6,
    reviewCount: 93,
    services: ['Residential Plumbing', 'Leak Detection'],
    startingPrice: 70,
    isVerified: false,
    isAvailable: true,
  ),
];