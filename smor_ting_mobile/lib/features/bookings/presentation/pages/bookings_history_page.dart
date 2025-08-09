import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

enum BookingStatus {
  upcoming,
  completed,
  cancelled,
}

class Booking {
  final String id;
  final String serviceName;
  final String providerName;
  final DateTime dateTime;
  final double amount;
  BookingStatus status;
  final String address;
  final double? rating;

  Booking({
    required this.id,
    required this.serviceName,
    required this.providerName,
    required this.dateTime,
    required this.amount,
    required this.status,
    required this.address,
    this.rating,
  });
}

class BookingsHistoryPage extends ConsumerStatefulWidget {
  const BookingsHistoryPage({super.key});

  @override
  ConsumerState<BookingsHistoryPage> createState() => _BookingsHistoryPageState();
}

class _BookingsHistoryPageState extends ConsumerState<BookingsHistoryPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  BookingStatus selectedStatus = BookingStatus.upcoming;

  List<Booking> bookings = [
    Booking(
      id: '1',
      serviceName: 'Plumbing Service',
      providerName: 'John Smith',
      dateTime: DateTime.now().add(const Duration(days: 2)),
      amount: 75.00,
      status: BookingStatus.upcoming,
      address: '123 Main St, Accra',
    ),
    Booking(
      id: '2',
      serviceName: 'Electrical Repair',
      providerName: 'Mike Johnson',
      dateTime: DateTime.now().subtract(const Duration(days: 1)),
      amount: 120.00,
      status: BookingStatus.completed,
      address: '456 Oak Ave, Accra',
      rating: 4.5,
    ),
    Booking(
      id: '3',
      serviceName: 'Cleaning Service',
      providerName: 'Sarah Wilson',
      dateTime: DateTime.now().subtract(const Duration(days: 3)),
      amount: 60.00,
      status: BookingStatus.cancelled,
      address: '789 Pine St, Accra',
    ),
    Booking(
      id: '4',
      serviceName: 'Carpentry Work',
      providerName: 'David Brown',
      dateTime: DateTime.now().subtract(const Duration(days: 5)),
      amount: 150.00,
      status: BookingStatus.completed,
      address: '321 Elm St, Accra',
      rating: 5.0,
    ),
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _tabController.addListener(() {
      setState(() {
        selectedStatus = BookingStatus.values[_tabController.index];
      });
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF2F2F7),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF007AFF)),
          onPressed: () => context.pop(),
        ),
        title: const Text(
          'My Bookings',
          style: TextStyle(
            color: Colors.black,
            fontSize: 18,
            fontWeight: FontWeight.w600,
          ),
        ),
        centerTitle: true,
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: const Color(0xFFFF3B30),
          labelColor: const Color(0xFFFF3B30),
          unselectedLabelColor: const Color(0xFF8E8E93),
          labelStyle: const TextStyle(
            fontWeight: FontWeight.w600,
            fontSize: 14,
          ),
          tabs: const [
            Tab(text: 'Upcoming'),
            Tab(text: 'Completed'),
            Tab(text: 'Cancelled'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildBookingsList(BookingStatus.upcoming),
          _buildBookingsList(BookingStatus.completed),
          _buildBookingsList(BookingStatus.cancelled),
        ],
      ),
    );
  }

  Widget _buildBookingsList(BookingStatus status) {
    final filteredBookings = bookings.where((booking) => booking.status == status).toList();
    
    if (filteredBookings.isEmpty) {
      return _buildEmptyState(status);
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: filteredBookings.length,
      itemBuilder: (context, index) {
        final booking = filteredBookings[index];
        return _buildBookingCard(booking);
      },
    );
  }

  Widget _buildEmptyState(BookingStatus status) {
    String message;
    IconData icon;
    
    switch (status) {
      case BookingStatus.upcoming:
        message = 'No upcoming bookings';
        icon = Icons.schedule;
        break;
      case BookingStatus.completed:
        message = 'No completed bookings';
        icon = Icons.check_circle;
        break;
      case BookingStatus.cancelled:
        message = 'No cancelled bookings';
        icon = Icons.cancel;
        break;
    }

    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            icon,
            size: 64,
            color: const Color(0xFF8E8E93),
          ),
          const SizedBox(height: 16),
          Text(
            message,
            style: const TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w600,
              color: Color(0xFF8E8E93),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Your ${status.name} bookings will appear here',
            style: const TextStyle(
              fontSize: 14,
              color: Color(0xFF8E8E93),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBookingCard(Booking booking) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: const Color(0xFFE5E5EA)),
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
            // Header with service name and status
            Row(
              children: [
                Expanded(
                  child: Text(
                    booking.serviceName,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: Colors.black,
                    ),
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: _getStatusColor(booking.status),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    _getStatusText(booking.status),
                    style: const TextStyle(
                      fontSize: 10,
                      color: Colors.white,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            
            // Provider info
            Row(
              children: [
                const Icon(
                  Icons.person,
                  size: 16,
                  color: Color(0xFF8E8E93),
                ),
                const SizedBox(width: 8),
                Text(
                  booking.providerName,
                  style: const TextStyle(
                    fontSize: 14,
                    color: Color(0xFF8E8E93),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            
            // Date and time
            Row(
              children: [
                const Icon(
                  Icons.schedule,
                  size: 16,
                  color: Color(0xFF8E8E93),
                ),
                const SizedBox(width: 8),
                Text(
                  _formatDateTime(booking.dateTime),
                  style: const TextStyle(
                    fontSize: 14,
                    color: Color(0xFF8E8E93),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            
            // Address
            Row(
              children: [
                const Icon(
                  Icons.location_on,
                  size: 16,
                  color: Color(0xFF8E8E93),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    booking.address,
                    style: const TextStyle(
                      fontSize: 14,
                      color: Color(0xFF8E8E93),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            
            // Amount and actions
            Row(
              children: [
                Text(
                  '\$${booking.amount.toStringAsFixed(2)}',
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w600,
                    color: Color(0xFFFF3B30),
                  ),
                ),
                const Spacer(),
                if (booking.rating != null) ...[
                  Row(
                    children: [
                      const Icon(
                        Icons.star,
                        size: 16,
                        color: Color(0xFF007AFF),
                      ),
                      const SizedBox(width: 4),
                      Text(
                        booking.rating!.toString(),
                        style: const TextStyle(
                          fontSize: 14,
                          color: Color(0xFF007AFF),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(width: 12),
                ],
                if (booking.status == BookingStatus.upcoming)
                  ElevatedButton(
                    onPressed: () {
                      _showBookingActions(booking);
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFFFF3B30),
                      foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    ),
                    child: const Text(
                      'Actions',
                      style: TextStyle(fontSize: 12),
                    ),
                  ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Color _getStatusColor(BookingStatus status) {
    switch (status) {
      case BookingStatus.upcoming:
        return const Color(0xFF007AFF);
      case BookingStatus.completed:
        return const Color(0xFF34C759);
      case BookingStatus.cancelled:
        return const Color(0xFFFF3B30);
    }
  }

  String _getStatusText(BookingStatus status) {
    switch (status) {
      case BookingStatus.upcoming:
        return 'Upcoming';
      case BookingStatus.completed:
        return 'Completed';
      case BookingStatus.cancelled:
        return 'Cancelled';
    }
  }

  String _formatDateTime(DateTime dateTime) {
    final now = DateTime.now();
    final difference = dateTime.difference(now).inDays;
    
    if (difference == 0) {
      return 'Today at ${_formatTime(dateTime)}';
    } else if (difference == 1) {
      return 'Tomorrow at ${_formatTime(dateTime)}';
    } else if (difference > 0) {
      return '${dateTime.day}/${dateTime.month}/${dateTime.year} at ${_formatTime(dateTime)}';
    } else {
      return '${dateTime.day}/${dateTime.month}/${dateTime.year} at ${_formatTime(dateTime)}';
    }
  }

  String _formatTime(DateTime dateTime) {
    final hour = dateTime.hour;
    final minute = dateTime.minute;
    final period = hour >= 12 ? 'PM' : 'AM';
    final displayHour = hour > 12 ? hour - 12 : (hour == 0 ? 12 : hour);
    return '${displayHour.toString().padLeft(2, '0')}:${minute.toString().padLeft(2, '0')} $period';
  }

  void _showBookingActions(Booking booking) {
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
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: const Color(0xFFE5E5EA),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 20),
              const Text(
                'Booking Actions',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w600,
                  color: Colors.black,
                ),
              ),
              const SizedBox(height: 20),
              _buildActionButton(
                'Reschedule',
                Icons.schedule,
                const Color(0xFF007AFF),
                () {
                  Navigator.of(context).pop();
                  _showRescheduleDialog(booking);
                },
              ),
              const SizedBox(height: 12),
              _buildActionButton(
                'Cancel Booking',
                Icons.cancel,
                const Color(0xFFFF3B30),
                () {
                  Navigator.of(context).pop();
                  _showCancelDialog(booking);
                },
              ),
              const SizedBox(height: 12),
              _buildActionButton(
                'Contact Provider',
                Icons.message,
                const Color(0xFF34C759),
                () {
                  Navigator.of(context).pop();
                  _showContactDialog(booking);
                },
              ),
              const SizedBox(height: 20),
            ],
          ),
        );
      },
    );
  }

  Widget _buildActionButton(String title, IconData icon, Color color, VoidCallback onTap) {
    return InkWell(
      onTap: onTap,
      child: Container(
        width: double.infinity,
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          border: Border.all(color: color),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          children: [
            Icon(icon, color: color, size: 20),
            const SizedBox(width: 12),
            Text(
              title,
              style: TextStyle(
                fontSize: 16,
                color: color,
                fontWeight: FontWeight.w500,
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showRescheduleDialog(Booking booking) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Reschedule Booking'),
          content: const Text('Reschedule functionality will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                _showSnackBar('Booking rescheduled successfully');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Reschedule'),
            ),
          ],
        );
      },
    );
  }

  void _showCancelDialog(Booking booking) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Cancel Booking'),
          content: const Text('Are you sure you want to cancel this booking?'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('No'),
            ),
            ElevatedButton(
              onPressed: () {
                setState(() {
                  booking.status = BookingStatus.cancelled;
                });
                Navigator.of(context).pop();
                _showSnackBar('Booking cancelled successfully');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Yes'),
            ),
          ],
        );
      },
    );
  }

  void _showContactDialog(Booking booking) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Contact Provider'),
          content: const Text('Contact functionality will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                _showSnackBar('Contacting provider...');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Contact'),
            ),
          ],
        );
      },
    );
  }

  void _showSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: const Color(0xFF34C759),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
      ),
    );
  }
} 