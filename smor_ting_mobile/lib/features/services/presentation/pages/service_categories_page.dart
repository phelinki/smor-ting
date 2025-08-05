import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/models/service.dart';

class ServiceCategoriesPage extends ConsumerStatefulWidget {
  const ServiceCategoriesPage({super.key});

  @override
  ConsumerState<ServiceCategoriesPage> createState() => _ServiceCategoriesPageState();
}

class _ServiceCategoriesPageState extends ConsumerState<ServiceCategoriesPage> {
  final TextEditingController _searchController = TextEditingController();
  String _searchQuery = '';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  // Mock data for demonstration - this would come from a provider in real implementation
  final List<ServiceCategory> _mockCategories = [
    ServiceCategory(
      id: '1',
      name: 'Plumbing',
      description: 'Water pipes, drainage, and bathroom repairs',
      icon: 'plumbing',
      color: '#2196F3',
      isActive: true,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    ),
    ServiceCategory(
      id: '2',
      name: 'Electrical',
      description: 'Wiring, lighting, and electrical installations',
      icon: 'electrical',
      color: '#FF9800',
      isActive: true,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    ),
    ServiceCategory(
      id: '3',
      name: 'Carpentry',
      description: 'Furniture, doors, windows, and woodwork',
      icon: 'carpentry',
      color: '#8BC34A',
      isActive: true,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    ),
    ServiceCategory(
      id: '4',
      name: 'Painting',
      description: 'Interior and exterior painting services',
      icon: 'painting',
      color: '#E91E63',
      isActive: true,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    ),
    ServiceCategory(
      id: '5',
      name: 'Cleaning',
      description: 'House cleaning and maintenance services',
      icon: 'cleaning',
      color: '#9C27B0',
      isActive: true,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    ),
    ServiceCategory(
      id: '6',
      name: 'Gardening',
      description: 'Landscaping, lawn care, and garden maintenance',
      icon: 'gardening',
      color: '#4CAF50',
      isActive: true,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    ),
    ServiceCategory(
      id: '7',
      name: 'Roofing',
      description: 'Roof repairs, installation, and maintenance',
      icon: 'roofing',
      color: '#795548',
      isActive: true,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    ),
    ServiceCategory(
      id: '8',
      name: 'Appliance Repair',
      description: 'Repair and maintenance of home appliances',
      icon: 'appliance',
      color: '#607D8B',
      isActive: true,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    ),
  ];

  List<ServiceCategory> get _filteredCategories {
    if (_searchQuery.isEmpty) {
      return _mockCategories;
    }
    return _mockCategories.where((category) =>
      category.name.toLowerCase().contains(_searchQuery.toLowerCase()) ||
      category.description.toLowerCase().contains(_searchQuery.toLowerCase())
    ).toList();
  }

  IconData _getIconForCategory(String iconName) {
    switch (iconName) {
      case 'plumbing':
        return Icons.plumbing;
      case 'electrical':
        return Icons.electrical_services;
      case 'carpentry':
        return Icons.carpenter;
      case 'painting':
        return Icons.format_paint;
      case 'cleaning':
        return Icons.cleaning_services;
      case 'gardening':
        return Icons.grass;
      case 'roofing':
        return Icons.roofing;
      case 'appliance':
        return Icons.kitchen;
      default:
        return Icons.build;
    }
  }

  Color _getColorFromHex(String hexColor) {
    hexColor = hexColor.replaceAll('#', '');
    return Color(int.parse('FF$hexColor', radix: 16));
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final filteredCategories = _filteredCategories;

    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        title: const Text(
          'Service Categories',
          style: TextStyle(
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
                hintText: 'Search services...',
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

          // Categories Grid
          Expanded(
            child: filteredCategories.isEmpty
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
                          'Try searching with different keywords',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: Colors.grey[500],
                          ),
                        ),
                      ],
                    ),
                  )
                : GridView.builder(
                    padding: const EdgeInsets.all(16),
                    gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: 2,
                      childAspectRatio: 1.1,
                      crossAxisSpacing: 16,
                      mainAxisSpacing: 16,
                    ),
                    itemCount: filteredCategories.length,
                    itemBuilder: (context, index) {
                      final category = filteredCategories[index];
                      final categoryColor = _getColorFromHex(category.color);
                      
                      return GestureDetector(
                        onTap: () {
                          context.push('/services/${category.id}', extra: category);
                        },
                        child: Container(
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
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Container(
                                width: 60,
                                height: 60,
                                decoration: BoxDecoration(
                                  color: categoryColor.withOpacity(0.1),
                                  borderRadius: BorderRadius.circular(16),
                                ),
                                child: Icon(
                                  _getIconForCategory(category.icon),
                                  size: 30,
                                  color: categoryColor,
                                ),
                              ),
                              const SizedBox(height: 12),
                              Text(
                                category.name,
                                style: theme.textTheme.titleMedium?.copyWith(
                                  fontWeight: FontWeight.w600,
                                  color: const Color(0xFF002868),
                                ),
                                textAlign: TextAlign.center,
                              ),
                              const SizedBox(height: 4),
                              Padding(
                                padding: const EdgeInsets.symmetric(horizontal: 8),
                                child: Text(
                                  category.description,
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: Colors.grey[600],
                                  ),
                                  textAlign: TextAlign.center,
                                  maxLines: 2,
                                  overflow: TextOverflow.ellipsis,
                                ),
                              ),
                            ],
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