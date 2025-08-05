import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/models/service.dart';

class ServicesListPage extends ConsumerStatefulWidget {
  final ServiceCategory category;

  const ServicesListPage({
    super.key,
    required this.category,
  });

  @override
  ConsumerState<ServicesListPage> createState() => _ServicesListPageState();
}

class _ServicesListPageState extends ConsumerState<ServicesListPage> {
  final TextEditingController _searchController = TextEditingController();
  String _searchQuery = '';
  String _sortBy = 'rating'; // rating, price_low, price_high, newest

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  // Mock data for demonstration
  final List<Service> _mockServices = [
    Service(
      id: '1',
      name: 'Emergency Plumbing Repair',
      description: 'Quick fixes for leaks, clogs, and pipe repairs. Available 24/7 for emergency situations.',
      categoryId: '1',
      providerId: 'p1',
      price: 75.0,
      currency: 'USD',
      duration: 120,
      images: ['https://example.com/plumbing1.jpg'],
      isActive: true,
      rating: 4.8,
      reviewCount: 156,
      createdAt: DateTime.now().subtract(const Duration(days: 30)),
      updatedAt: DateTime.now(),
    ),
    Service(
      id: '2',
      name: 'Bathroom Installation',
      description: 'Complete bathroom setup including toilet, sink, shower installation and plumbing connections.',
      categoryId: '1',
      providerId: 'p2',
      price: 450.0,
      currency: 'USD',
      duration: 480,
      images: ['https://example.com/plumbing2.jpg'],
      isActive: true,
      rating: 4.9,
      reviewCount: 89,
      createdAt: DateTime.now().subtract(const Duration(days: 15)),
      updatedAt: DateTime.now(),
    ),
    Service(
      id: '3',
      name: 'Drain Cleaning Service',
      description: 'Professional drain cleaning for kitchen, bathroom, and outdoor drains using modern equipment.',
      categoryId: '1',
      providerId: 'p3',
      price: 95.0,
      currency: 'USD',
      duration: 90,
      images: ['https://example.com/plumbing3.jpg'],
      isActive: true,
      rating: 4.6,
      reviewCount: 203,
      createdAt: DateTime.now().subtract(const Duration(days: 45)),
      updatedAt: DateTime.now(),
    ),
    Service(
      id: '4',
      name: 'Water Heater Repair',
      description: 'Repair and maintenance of water heaters, including electric and gas units.',
      categoryId: '1',
      providerId: 'p1',
      price: 125.0,
      currency: 'USD',
      duration: 180,
      images: ['https://example.com/plumbing4.jpg'],
      isActive: true,
      rating: 4.7,
      reviewCount: 78,
      createdAt: DateTime.now().subtract(const Duration(days: 60)),
      updatedAt: DateTime.now(),
    ),
  ];

  List<Service> get _filteredAndSortedServices {
    var services = _mockServices.where((service) => service.categoryId == widget.category.id).toList();
    
    // Filter by search query
    if (_searchQuery.isNotEmpty) {
      services = services.where((service) =>
        service.name.toLowerCase().contains(_searchQuery.toLowerCase()) ||
        service.description.toLowerCase().contains(_searchQuery.toLowerCase())
      ).toList();
    }

    // Sort services
    switch (_sortBy) {
      case 'rating':
        services.sort((a, b) => b.rating.compareTo(a.rating));
        break;
      case 'price_low':
        services.sort((a, b) => a.price.compareTo(b.price));
        break;
      case 'price_high':
        services.sort((a, b) => b.price.compareTo(a.price));
        break;
      case 'newest':
        services.sort((a, b) => b.createdAt.compareTo(a.createdAt));
        break;
    }

    return services;
  }

  void _showSortOptions() {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return Container(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Sort by',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.w600,
                  color: const Color(0xFF002868),
                ),
              ),
              const SizedBox(height: 20),
              _SortOption(
                title: 'Highest Rated',
                value: 'rating',
                groupValue: _sortBy,
                onChanged: (value) {
                  setState(() {
                    _sortBy = value!;
                  });
                  Navigator.pop(context);
                },
              ),
              _SortOption(
                title: 'Price: Low to High',
                value: 'price_low',
                groupValue: _sortBy,
                onChanged: (value) {
                  setState(() {
                    _sortBy = value!;
                  });
                  Navigator.pop(context);
                },
              ),
              _SortOption(
                title: 'Price: High to Low',
                value: 'price_high',
                groupValue: _sortBy,
                onChanged: (value) {
                  setState(() {
                    _sortBy = value!;
                  });
                  Navigator.pop(context);
                },
              ),
              _SortOption(
                title: 'Newest First',
                value: 'newest',
                groupValue: _sortBy,
                onChanged: (value) {
                  setState(() {
                    _sortBy = value!;
                  });
                  Navigator.pop(context);
                },
              ),
              const SizedBox(height: 20),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final services = _filteredAndSortedServices;

    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        title: Text(
          widget.category.name,
          style: const TextStyle(
            fontWeight: FontWeight.w600,
            color: Color(0xFF002868),
          ),
        ),
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF002868)),
          onPressed: () => Navigator.of(context).pop(),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.sort, color: Color(0xFF002868)),
            onPressed: _showSortOptions,
          ),
        ],
      ),
      body: Column(
        children: [
          // Search Bar
          Container(
            color: Colors.white,
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Search ${widget.category.name.toLowerCase()} services...',
                prefixIcon: const Icon(Icons.search, color: Color(0xFF002868)),
                suffixIcon: _searchQuery.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          setState(() {
                            _searchQuery = '';
                          });
                        },
                      )
                    : null,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: const BorderSide(color: Color(0xFFE0E0E0)),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: const BorderSide(color: Color(0xFFD21034), width: 2),
                ),
                fillColor: Colors.grey[50],
                filled: true,
              ),
              onChanged: (value) {
                setState(() {
                  _searchQuery = value;
                });
              },
            ),
          ),

          // Services List
          Expanded(
            child: services.isEmpty
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.search_off,
                          size: 64,
                          color: Colors.grey[400],
                        ),
                        const SizedBox(height: 16),
                        Text(
                          'No services found',
                          style: theme.textTheme.titleLarge?.copyWith(
                            color: Colors.grey[600],
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          'Try adjusting your search or filters',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: Colors.grey[500],
                          ),
                        ),
                      ],
                    ),
                  )
                : ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: services.length,
                    itemBuilder: (context, index) {
                      final service = services[index];
                      
                      return GestureDetector(
                        onTap: () {
                          context.push('/service/${service.id}', extra: service);
                        },
                        child: Container(
                          margin: const EdgeInsets.only(bottom: 16),
                          decoration: BoxDecoration(
                            color: Colors.white,
                            borderRadius: BorderRadius.circular(16),
                            boxShadow: [
                              BoxShadow(
                                color: Colors.black.withOpacity(0.05),
                                blurRadius: 10,
                                offset: const Offset(0, 2),
                              ),
                            ],
                          ),
                          child: Padding(
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Row(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    // Service Image Placeholder
                                    Container(
                                      width: 80,
                                      height: 80,
                                      decoration: BoxDecoration(
                                        color: const Color(0xFF002868).withOpacity(0.1),
                                        borderRadius: BorderRadius.circular(12),
                                      ),
                                      child: const Icon(
                                        Icons.plumbing,
                                        color: Color(0xFF002868),
                                        size: 40,
                                      ),
                                    ),
                                    const SizedBox(width: 16),
                                    
                                    // Service Details
                                    Expanded(
                                      child: Column(
                                        crossAxisAlignment: CrossAxisAlignment.start,
                                        children: [
                                          Text(
                                            service.name,
                                            style: theme.textTheme.titleMedium?.copyWith(
                                              fontWeight: FontWeight.w600,
                                              color: const Color(0xFF002868),
                                            ),
                                          ),
                                          const SizedBox(height: 4),
                                          Text(
                                            service.description,
                                            style: theme.textTheme.bodySmall?.copyWith(
                                              color: Colors.grey[600],
                                            ),
                                            maxLines: 2,
                                            overflow: TextOverflow.ellipsis,
                                          ),
                                          const SizedBox(height: 8),
                                          
                                          // Rating and Reviews
                                          Row(
                                            children: [
                                              const Icon(
                                                Icons.star,
                                                color: Colors.amber,
                                                size: 16,
                                              ),
                                              const SizedBox(width: 4),
                                              Text(
                                                service.rating.toString(),
                                                style: theme.textTheme.bodySmall?.copyWith(
                                                  fontWeight: FontWeight.w600,
                                                ),
                                              ),
                                              const SizedBox(width: 4),
                                              Text(
                                                '(${service.reviewCount} reviews)',
                                                style: theme.textTheme.bodySmall?.copyWith(
                                                  color: Colors.grey[600],
                                                ),
                                              ),
                                            ],
                                          ),
                                        ],
                                      ),
                                    ),
                                  ],
                                ),
                                
                                const SizedBox(height: 16),
                                
                                // Price and Duration
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                  children: [
                                    Column(
                                      crossAxisAlignment: CrossAxisAlignment.start,
                                      children: [
                                        Text(
                                          'Starting from',
                                          style: theme.textTheme.bodySmall?.copyWith(
                                            color: Colors.grey[600],
                                          ),
                                        ),
                                        Text(
                                          service.formattedPrice,
                                          style: theme.textTheme.titleMedium?.copyWith(
                                            fontWeight: FontWeight.bold,
                                            color: const Color(0xFFD21034),
                                          ),
                                        ),
                                      ],
                                    ),
                                    Column(
                                      crossAxisAlignment: CrossAxisAlignment.end,
                                      children: [
                                        Text(
                                          'Duration',
                                          style: theme.textTheme.bodySmall?.copyWith(
                                            color: Colors.grey[600],
                                          ),
                                        ),
                                        Text(
                                          service.formattedDuration,
                                          style: theme.textTheme.bodyMedium?.copyWith(
                                            fontWeight: FontWeight.w600,
                                            color: const Color(0xFF002868),
                                          ),
                                        ),
                                      ],
                                    ),
                                  ],
                                ),
                              ],
                            ),
                          ),
                        ),
                      );
                    },
                  ),
          ),
        ],
      ),
    );
  }
}

class _SortOption extends StatelessWidget {
  final String title;
  final String value;
  final String groupValue;
  final ValueChanged<String?> onChanged;

  const _SortOption({
    required this.title,
    required this.value,
    required this.groupValue,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return RadioListTile<String>(
      title: Text(title),
      value: value,
      groupValue: groupValue,
      onChanged: onChanged,
      activeColor: const Color(0xFFD21034),
      contentPadding: EdgeInsets.zero,
    );
  }
}