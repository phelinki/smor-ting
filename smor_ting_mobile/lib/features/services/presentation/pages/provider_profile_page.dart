import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';

class ProviderProfilePage extends ConsumerStatefulWidget {
  final String providerId;
  
  const ProviderProfilePage({
    super.key,
    required this.providerId,
  });

  @override
  ConsumerState<ProviderProfilePage> createState() => _ProviderProfilePageState();
}

class _ProviderProfilePageState extends ConsumerState<ProviderProfilePage> {
  DateTime? _selectedDate;
  TimeOfDay? _selectedTime;
  bool _isDateTimeSelected = false;

  // Agent's working hours (9 AM to 6 PM)
  final List<TimeOfDay> _availableTimes = [
    const TimeOfDay(hour: 9, minute: 0),
    const TimeOfDay(hour: 10, minute: 0),
    const TimeOfDay(hour: 11, minute: 0),
    const TimeOfDay(hour: 12, minute: 0),
    const TimeOfDay(hour: 13, minute: 0),
    const TimeOfDay(hour: 14, minute: 0),
    const TimeOfDay(hour: 15, minute: 0),
    const TimeOfDay(hour: 16, minute: 0),
    const TimeOfDay(hour: 17, minute: 0),
    const TimeOfDay(hour: 18, minute: 0),
  ];

  void _selectDate() async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: DateTime.now().add(const Duration(days: 1)),
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 30)),
      builder: (context, child) {
        return Theme(
          data: Theme.of(context).copyWith(
            colorScheme: const ColorScheme.light(
              primary: AppTheme.primaryRed,
              onPrimary: AppTheme.white,
            ),
          ),
          child: child!,
        );
      },
    );
    if (picked != null) {
      setState(() {
        _selectedDate = picked;
        _updateBookingState();
      });
    }
  }

  void _selectTime() async {
    final TimeOfDay? picked = await showTimePicker(
      context: context,
      initialTime: const TimeOfDay(hour: 9, minute: 0),
      builder: (context, child) {
        return Theme(
          data: Theme.of(context).copyWith(
            colorScheme: const ColorScheme.light(
              primary: AppTheme.primaryRed,
              onPrimary: AppTheme.white,
            ),
          ),
          child: child!,
        );
      },
    );
    if (picked != null) {
      setState(() {
        _selectedTime = picked;
        _updateBookingState();
      });
    }
  }

  void _updateBookingState() {
    setState(() {
      _isDateTimeSelected = _selectedDate != null && _selectedTime != null;
    });
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year}';
  }

  String _formatTime(TimeOfDay time) {
    return '${time.hour.toString().padLeft(2, '0')}:${time.minute.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    // TODO: Fetch provider data based on providerId
    final provider = _sampleProvider; // Using sample data for now
    
    return Scaffold(
      backgroundColor: AppTheme.white,
      body: CustomScrollView(
        slivers: [
          // App Bar with Cover Photo
          SliverAppBar(
            expandedHeight: 200,
            pinned: true,
            backgroundColor: AppTheme.secondaryBlue,
            leading: IconButton(
              icon: const Icon(Icons.arrow_back_ios, color: AppTheme.white),
              onPressed: () => context.pop(),
            ),
            actions: [
              IconButton(
                icon: const Icon(Icons.favorite_border, color: AppTheme.white),
                onPressed: () {
                  // TODO: Add to favorites
                },
              ),
              IconButton(
                icon: const Icon(Icons.share, color: AppTheme.white),
                onPressed: () {
                  // TODO: Share provider
                },
              ),
            ],
            flexibleSpace: FlexibleSpaceBar(
              background: Container(
                decoration: const BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topCenter,
                    end: Alignment.bottomCenter,
                    colors: [
                      AppTheme.secondaryBlue,
                      Color(0xFF0056CC),
                    ],
                  ),
                ),
                child: const Center(
                  child: Icon(
                    Icons.handyman,
                    size: 80,
                    color: AppTheme.white,
                  ),
                ),
              ),
            ),
          ),
          
          // Profile Content
          SliverToBoxAdapter(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Provider Header
                Container(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Profile Picture and Basic Info
                      Row(
                        children: [
                          Container(
                            width: 80,
                            height: 80,
                            decoration: BoxDecoration(
                              color: AppTheme.lightGray,
                              borderRadius: BorderRadius.circular(40),
                              border: Border.all(
                                color: AppTheme.white,
                                width: 4,
                              ),
                            ),
                            child: const Icon(
                              Icons.person,
                              size: 40,
                              color: AppTheme.textSecondary,
                            ),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  provider.name,
                                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Row(
                                  children: [
                                    Icon(
                                      Icons.star,
                                      size: 16,
                                      color: AppTheme.warning,
                                    ),
                                    const SizedBox(width: 4),
                                    Text(
                                      '${provider.rating} (${provider.reviewCount} reviews)',
                                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                        color: AppTheme.textSecondary,
                                      ),
                                    ),
                                  ],
                                ),
                                const SizedBox(height: 4),
                                Row(
                                  children: [
                                    Icon(
                                      Icons.verified,
                                      size: 16,
                                      color: AppTheme.successGreen,
                                    ),
                                    const SizedBox(width: 4),
                                    Text(
                                      'Verified Agent',
                                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                        color: AppTheme.successGreen,
                                        fontWeight: FontWeight.w600,
                                      ),
                                    ),
                                  ],
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                      
                      const SizedBox(height: 24),
                      
                      // Date and Time Selection Section
                      _buildDateTimeSelection(),
                      
                      const SizedBox(height: 24),
                      
                      // Action Buttons
                      Row(
                        children: [
                          Expanded(
                            child: ElevatedButton(
                              onPressed: _isDateTimeSelected ? () {
                                context.push('/booking-confirmation/${provider.id}');
                              } : null,
                              style: ElevatedButton.styleFrom(
                                backgroundColor: _isDateTimeSelected 
                                    ? AppTheme.primaryRed 
                                    : AppTheme.lightGray,
                                foregroundColor: _isDateTimeSelected 
                                    ? AppTheme.white 
                                    : AppTheme.textSecondary,
                                padding: const EdgeInsets.symmetric(vertical: 16),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                              ),
                              child: Text(
                                _isDateTimeSelected ? 'Book Now' : 'Select Date & Time',
                                style: const TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: OutlinedButton(
                              onPressed: () {
                                // TODO: Navigate to messages
                              },
                              style: OutlinedButton.styleFrom(
                                foregroundColor: AppTheme.primaryRed,
                                side: const BorderSide(color: AppTheme.primaryRed),
                                padding: const EdgeInsets.symmetric(vertical: 16),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                              ),
                              child: const Text(
                                'Message',
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                
                const Divider(height: 1),
                
                // Services Offered Section
                _buildSection(
                  title: 'Services Offered',
                  child: Column(
                    children: provider.services.map((service) {
                      return _buildServiceItem(service);
                    }).toList(),
                  ),
                ),
                
                const Divider(height: 1),
                
                // Pricing Section
                _buildSection(
                  title: 'Pricing',
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Starting at \$${provider.startingPrice}/hour',
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          color: AppTheme.secondaryBlue,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Final pricing may vary based on job complexity and materials needed.',
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: AppTheme.textSecondary,
                        ),
                      ),
                    ],
                  ),
                ),
                
                const Divider(height: 1),
                
                // Availability Section
                _buildSection(
                  title: 'Availability',
                  child: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: provider.isAvailable 
                          ? AppTheme.successGreen.withOpacity(0.1)
                          : AppTheme.primaryRed.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Row(
                      children: [
                        Icon(
                          provider.isAvailable ? Icons.check_circle : Icons.schedule,
                          color: provider.isAvailable 
                              ? AppTheme.successGreen 
                              : AppTheme.primaryRed,
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                provider.isAvailable ? 'Available Now' : 'Currently Busy',
                                style: Theme.of(context).textTheme.titleSmall?.copyWith(
                                  fontWeight: FontWeight.bold,
                                  color: provider.isAvailable 
                                      ? AppTheme.successGreen 
                                      : AppTheme.primaryRed,
                                ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                'Working hours: 9:00 AM - 6:00 PM',
                                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                  color: AppTheme.textSecondary,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                
                const Divider(height: 1),
                
                // About Section
                _buildSection(
                  title: 'About',
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        provider.description,
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: AppTheme.textPrimary,
                          height: 1.5,
                        ),
                      ),
                      const SizedBox(height: 16),
                      Row(
                        children: [
                          Icon(
                            Icons.work,
                            size: 16,
                            color: AppTheme.textSecondary,
                          ),
                          const SizedBox(width: 8),
                          Text(
                            '${provider.yearsOfExperience} years of experience',
                            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                              color: AppTheme.textSecondary,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                
                const Divider(height: 1),
                
                // Reviews Section
                _buildSection(
                  title: 'Reviews',
                  child: Column(
                    children: _sampleReviews.map((review) {
                      return _buildReviewItem(review);
                    }).toList(),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDateTimeSelection() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.lightGray.withOpacity(0.3),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppTheme.lightGray,
          width: 1,
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Select Appointment Date & Time',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
              color: AppTheme.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          
          // Date Selection
          Row(
            children: [
              Expanded(
                child: InkWell(
                  onTap: _selectDate,
                  child: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: AppTheme.white,
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(
                        color: _selectedDate != null 
                            ? AppTheme.primaryRed 
                            : AppTheme.lightGray,
                        width: 2,
                      ),
                    ),
                    child: Row(
                      children: [
                        Icon(
                          Icons.calendar_today,
                          color: _selectedDate != null 
                              ? AppTheme.primaryRed 
                              : AppTheme.textSecondary,
                          size: 20,
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Text(
                            _selectedDate != null 
                                ? _formatDate(_selectedDate!)
                                : 'Select Date',
                            style: TextStyle(
                              color: _selectedDate != null 
                                  ? AppTheme.textPrimary 
                                  : AppTheme.textSecondary,
                              fontWeight: _selectedDate != null 
                                  ? FontWeight.w600 
                                  : FontWeight.normal,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: InkWell(
                  onTap: _selectTime,
                  child: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: AppTheme.white,
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(
                        color: _selectedTime != null 
                            ? AppTheme.primaryRed 
                            : AppTheme.lightGray,
                        width: 2,
                      ),
                    ),
                    child: Row(
                      children: [
                        Icon(
                          Icons.access_time,
                          color: _selectedTime != null 
                              ? AppTheme.primaryRed 
                              : AppTheme.textSecondary,
                          size: 20,
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Text(
                            _selectedTime != null 
                                ? _formatTime(_selectedTime!)
                                : 'Select Time',
                            style: TextStyle(
                              color: _selectedTime != null 
                                  ? AppTheme.textPrimary 
                                  : AppTheme.textSecondary,
                              fontWeight: _selectedTime != null 
                                  ? FontWeight.w600 
                                  : FontWeight.normal,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ],
          ),
          
          if (_selectedDate != null && _selectedTime != null) ...[
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: AppTheme.successGreen.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
                border: Border.all(
                  color: AppTheme.successGreen.withOpacity(0.3),
                ),
              ),
              child: Row(
                children: [
                  Icon(
                    Icons.check_circle,
                    color: AppTheme.successGreen,
                    size: 16,
                  ),
                  const SizedBox(width: 8),
                  Text(
                    'Appointment scheduled for ${_formatDate(_selectedDate!)} at ${_formatTime(_selectedTime!)}',
                    style: TextStyle(
                      color: AppTheme.successGreen,
                      fontWeight: FontWeight.w600,
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildSection({
    required String title,
    required Widget child,
  }) {
    return Container(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          child,
        ],
      ),
    );
  }

  Widget _buildServiceItem(ProviderService service) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppTheme.lightGray,
          width: 1,
        ),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  service.name,
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  service.description,
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: AppTheme.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          Text(
            '\$${service.price}',
            style: Theme.of(context).textTheme.titleSmall?.copyWith(
              color: AppTheme.secondaryBlue,
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildReviewItem(Review review) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppTheme.lightGray,
          width: 1,
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(
                child: Text(
                  review.customerName,
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              Row(
                children: List.generate(5, (index) {
                  return Icon(
                    index < review.rating ? Icons.star : Icons.star_border,
                    size: 16,
                    color: AppTheme.warning,
                  );
                }),
              ),
            ],
          ),
          const SizedBox(height: 4),
          Text(
            review.date,
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
              color: AppTheme.textSecondary,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            review.comment,
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: AppTheme.textPrimary,
              height: 1.4,
            ),
          ),
        ],
      ),
    );
  }
}

// Data classes
class ProviderService {
  final String name;
  final String description;
  final int price;

  ProviderService({
    required this.name,
    required this.description,
    required this.price,
  });
}

class Review {
  final String customerName;
  final String date;
  final int rating;
  final String comment;

  Review({
    required this.customerName,
    required this.date,
    required this.rating,
    required this.comment,
  });
}

class DetailedProvider {
  final String id;
  final String name;
  final double rating;
  final int reviewCount;
  final List<ProviderService> services;
  final int startingPrice;
  final bool isVerified;
  final bool isAvailable;
  final int yearsOfExperience;
  final String description;

  DetailedProvider({
    required this.id,
    required this.name,
    required this.rating,
    required this.reviewCount,
    required this.services,
    required this.startingPrice,
    required this.isVerified,
    required this.isAvailable,
    required this.yearsOfExperience,
    required this.description,
  });
}

// Sample data
final DetailedProvider _sampleProvider = DetailedProvider(
  id: '1',
  name: 'John Martinez',
  rating: 4.8,
  reviewCount: 127,
  services: [
    ProviderService(
      name: 'Kitchen Plumbing Repair',
      description: 'Fix leaks, unclog drains, repair faucets',
      price: 75,
    ),
    ProviderService(
      name: 'Bathroom Plumbing',
      description: 'Toilet repair, shower installation, pipe work',
      price: 85,
    ),
    ProviderService(
      name: 'Emergency Plumbing',
      description: '24/7 emergency plumbing services',
      price: 120,
    ),
  ],
  startingPrice: 75,
  isVerified: true,
  isAvailable: true,
  yearsOfExperience: 8,
  description: 'Licensed and insured plumber with over 8 years of experience. Specializing in residential and commercial plumbing services. I take pride in delivering quality work and excellent customer service. Available for emergency calls and scheduled appointments.',
);

final List<Review> _sampleReviews = [
  Review(
    customerName: 'Sarah Wilson',
    date: '2 days ago',
    rating: 5,
    comment: 'John did an excellent job fixing our kitchen sink. He was professional, punctual, and cleaned up after himself. Highly recommend!',
  ),
  Review(
    customerName: 'Mike Chen',
    date: '1 week ago',
    rating: 5,
    comment: 'Great service! Fixed our bathroom leak quickly and explained everything clearly. Fair pricing too.',
  ),
  Review(
    customerName: 'Lisa Davis',
    date: '2 weeks ago',
    rating: 4,
    comment: 'Good work overall. Arrived on time and completed the job efficiently. Would use again.',
  ),
];