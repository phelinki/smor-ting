import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

enum BookingStatus {
  confirmed,
  enRoute,
  arrived,
  inProgress,
  completed,
}

class RealTimeTrackingPage extends ConsumerStatefulWidget {
  final String bookingId;
  final String providerName;
  final String serviceName;
  final String address;

  const RealTimeTrackingPage({
    super.key,
    required this.bookingId,
    required this.providerName,
    required this.serviceName,
    required this.address,
  });

  @override
  ConsumerState<RealTimeTrackingPage> createState() => _RealTimeTrackingPageState();
}

class _RealTimeTrackingPageState extends ConsumerState<RealTimeTrackingPage> {
  BookingStatus currentStatus = BookingStatus.confirmed;
  bool isProviderVisible = true;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: Column(
        children: [
          // Map View (Placeholder)
          Expanded(
            flex: 2,
            child: Container(
              width: double.infinity,
              decoration: BoxDecoration(
                color: const Color(0xFFF2F2F7),
                borderRadius: const BorderRadius.only(
                  bottomLeft: Radius.circular(20),
                  bottomRight: Radius.circular(20),
                ),
              ),
              child: Stack(
                children: [
                  // Map placeholder
                  Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Container(
                          width: 80,
                          height: 80,
                          decoration: BoxDecoration(
                            color: const Color(0xFF007AFF),
                            borderRadius: BorderRadius.circular(40),
                          ),
                          child: const Icon(
                            Icons.location_on,
                            color: Colors.white,
                            size: 40,
                          ),
                        ),
                        const SizedBox(height: 16),
                        const Text(
                          'Map View',
                          style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.w600,
                            color: Colors.black,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          'Provider location will be shown here',
                          style: TextStyle(
                            fontSize: 14,
                            color: Colors.grey[600],
                          ),
                        ),
                      ],
                    ),
                  ),
                  // Provider location indicator
                  if (isProviderVisible)
                    Positioned(
                      top: 100,
                      right: 50,
                      child: Container(
                        padding: const EdgeInsets.all(8),
                        decoration: BoxDecoration(
                          color: const Color(0xFFFF3B30),
                          borderRadius: BorderRadius.circular(20),
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withOpacity(0.2),
                              blurRadius: 8,
                              offset: const Offset(0, 2),
                            ),
                          ],
                        ),
                        child: const Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(
                              Icons.person,
                              color: Colors.white,
                              size: 16,
                            ),
                            SizedBox(width: 4),
                            Text(
                              'Provider',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 12,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                ],
              ),
            ),
          ),
          
          // Status and Details
          Expanded(
            flex: 1,
            child: Container(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Status Indicator
                  _buildStatusIndicator(),
                  const SizedBox(height: 16),
                  
                  // Booking Details
                  _buildBookingDetails(),
                  const SizedBox(height: 16),
                  
                  // Action Buttons
                  _buildActionButtons(),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusIndicator() {
    return Container(
      padding: const EdgeInsets.all(16),
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
      child: Column(
        children: [
          Row(
            children: [
              Icon(
                _getStatusIcon(),
                color: _getStatusColor(),
                size: 24,
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Text(
                  _getStatusText(),
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: _getStatusColor(),
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          LinearProgressIndicator(
            value: _getProgressValue(),
            backgroundColor: const Color(0xFFE5E5EA),
            valueColor: AlwaysStoppedAnimation<Color>(_getStatusColor()),
          ),
          const SizedBox(height: 8),
          Text(
            _getStatusDescription(),
            style: const TextStyle(
              fontSize: 12,
              color: Color(0xFF8E8E93),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBookingDetails() {
    return Container(
      padding: const EdgeInsets.all(16),
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
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Booking Details',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: Colors.black,
            ),
          ),
          const SizedBox(height: 12),
          _buildDetailRow('Service', widget.serviceName),
          _buildDetailRow('Provider', widget.providerName),
          _buildDetailRow('Address', widget.address),
          _buildDetailRow('Booking ID', widget.bookingId),
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text(
              label,
              style: const TextStyle(
                fontSize: 14,
                color: Color(0xFF8E8E93),
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                fontSize: 14,
                color: Colors.black,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActionButtons() {
    return Row(
      children: [
        Expanded(
          child: ElevatedButton.icon(
            onPressed: () {
              // Handle message action
              _showMessageDialog();
            },
            icon: const Icon(Icons.message, size: 20),
            label: const Text('Message'),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF007AFF),
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              padding: const EdgeInsets.symmetric(vertical: 12),
            ),
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: ElevatedButton.icon(
            onPressed: () {
              // Handle call action
              _showCallDialog();
            },
            icon: const Icon(Icons.call, size: 20),
            label: const Text('Call'),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF34C759),
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              padding: const EdgeInsets.symmetric(vertical: 12),
            ),
          ),
        ),
      ],
    );
  }

  IconData _getStatusIcon() {
    switch (currentStatus) {
      case BookingStatus.confirmed:
        return Icons.schedule;
      case BookingStatus.enRoute:
        return Icons.directions_car;
      case BookingStatus.arrived:
        return Icons.location_on;
      case BookingStatus.inProgress:
        return Icons.build;
      case BookingStatus.completed:
        return Icons.check_circle;
    }
  }

  Color _getStatusColor() {
    switch (currentStatus) {
      case BookingStatus.confirmed:
        return const Color(0xFF007AFF);
      case BookingStatus.enRoute:
        return const Color(0xFF007AFF);
      case BookingStatus.arrived:
        return const Color(0xFF34C759);
      case BookingStatus.inProgress:
        return const Color(0xFF007AFF);
      case BookingStatus.completed:
        return const Color(0xFF34C759);
    }
  }

  String _getStatusText() {
    switch (currentStatus) {
      case BookingStatus.confirmed:
        return 'Booking Confirmed';
      case BookingStatus.enRoute:
        return 'Provider En Route';
      case BookingStatus.arrived:
        return 'Provider Arrived';
      case BookingStatus.inProgress:
        return 'Service In Progress';
      case BookingStatus.completed:
        return 'Service Completed';
    }
  }

  double _getProgressValue() {
    switch (currentStatus) {
      case BookingStatus.confirmed:
        return 0.2;
      case BookingStatus.enRoute:
        return 0.4;
      case BookingStatus.arrived:
        return 0.6;
      case BookingStatus.inProgress:
        return 0.8;
      case BookingStatus.completed:
        return 1.0;
    }
  }

  String _getStatusDescription() {
    switch (currentStatus) {
      case BookingStatus.confirmed:
        return 'Your booking has been confirmed and provider will be assigned soon';
      case BookingStatus.enRoute:
        return 'Provider is on the way to your location';
      case BookingStatus.arrived:
        return 'Provider has arrived at your location';
      case BookingStatus.inProgress:
        return 'Provider is currently working on your service';
      case BookingStatus.completed:
        return 'Service has been completed successfully';
    }
  }

  void _showMessageDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Message Provider'),
          content: const Text('Chat functionality will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('OK'),
            ),
          ],
        );
      },
    );
  }

  void _showCallDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Call Provider'),
          content: const Text('Call functionality will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('OK'),
            ),
          ],
        );
      },
    );
  }
} 