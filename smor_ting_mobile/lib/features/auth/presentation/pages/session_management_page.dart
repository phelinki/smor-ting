import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/services/enhanced_auth_service.dart';
import '../../../../core/theme/app_theme.dart';
import '../providers/enhanced_auth_provider.dart';

/// Comprehensive session management page
class SessionManagementPage extends ConsumerStatefulWidget {
  const SessionManagementPage({super.key});

  @override
  ConsumerState<SessionManagementPage> createState() => _SessionManagementPageState();
}

class _SessionManagementPageState extends ConsumerState<SessionManagementPage> {
  List<SessionInfo>? _sessions;
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadSessions();
  }

  Future<void> _loadSessions() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final authNotifier = ref.read(enhancedAuthNotifierProvider.notifier);
      final sessions = await authNotifier.getUserSessions();
      
      setState(() {
        _sessions = sessions;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  Future<void> _revokeSession(String sessionId, String deviceName) async {
    final confirmed = await _showConfirmationDialog(
      'Sign Out Device',
      'Are you sure you want to sign out from $deviceName?',
    );

    if (!confirmed) return;

    try {
      final authNotifier = ref.read(enhancedAuthNotifierProvider.notifier);
      await authNotifier.revokeSession(sessionId);
      
      // Refresh sessions list
      await _loadSessions();
      
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Signed out from $deviceName'),
            backgroundColor: AppTheme.successGreen,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to sign out: ${e.toString()}'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  Future<void> _signOutAllDevices() async {
    final confirmed = await _showConfirmationDialog(
      'Sign Out All Devices',
      'This will sign you out from all devices including this one. You will need to sign in again.',
    );

    if (!confirmed) return;

    try {
      final authNotifier = ref.read(enhancedAuthNotifierProvider.notifier);
      await authNotifier.signOutAllDevices();
      
      if (mounted) {
        // Navigate to login page
        context.go('/landing');
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to sign out all devices: ${e.toString()}'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  Future<bool> _showConfirmationDialog(String title, String message) async {
    final result = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(title),
        content: Text(message),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppTheme.error,
              foregroundColor: AppTheme.white,
            ),
            child: const Text('Confirm'),
          ),
        ],
      ),
    );
    
    return result ?? false;
  }

  String _getDeviceDisplayName(SessionInfo session) {
    final platform = session.deviceInfo.platformDisplayName;
    final deviceName = session.deviceInfo.platform;
    return '$platform ($deviceName)';
  }

  String _getLocationInfo(SessionInfo session) {
    return session.ipAddress;
  }

  String _getLastActiveTime(SessionInfo session) {
    final now = DateTime.now();
    final difference = now.difference(session.lastActivity);
    
    if (difference.inMinutes < 1) {
      return 'Just now';
    } else if (difference.inMinutes < 60) {
      return '${difference.inMinutes}m ago';
    } else if (difference.inHours < 24) {
      return '${difference.inHours}h ago';
    } else {
      return '${difference.inDays}d ago';
    }
  }

  IconData _getDeviceIcon(SessionInfo session) {
    final platform = session.deviceInfo.platform.toLowerCase();
    
    if (platform.contains('ios') || platform.contains('iphone')) {
      return Icons.phone_iphone;
    } else if (platform.contains('android')) {
      return Icons.phone_android;
    } else if (platform.contains('web')) {
      return Icons.computer;
    } else {
      return Icons.devices;
    }
  }

  Color _getTrustIndicatorColor(SessionInfo session) {
    if (session.deviceInfo.isCompromised) {
      return AppTheme.error;
    } else if (session.deviceInfo.calculateTrustScore() > 0.8) {
      return AppTheme.success;
    } else {
      return AppTheme.warning;
    }
  }

  Widget _buildTrustIndicator(SessionInfo session) {
    final trustScore = session.deviceInfo.calculateTrustScore();
    final isCompromised = session.deviceInfo.isCompromised;
    
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(
          isCompromised ? Icons.security : Icons.verified_user,
          size: 16,
          color: _getTrustIndicatorColor(session),
        ),
        const SizedBox(width: 4),
        Text(
          isCompromised ? 'Untrusted' : 
          trustScore > 0.8 ? 'Trusted' : 'Limited Trust',
          style: TextStyle(
            fontSize: 12,
            color: _getTrustIndicatorColor(session),
            fontWeight: FontWeight.w500,
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.backgroundLight,
      appBar: AppBar(
        title: const Text(
          'Active Sessions',
          style: TextStyle(
            color: AppTheme.textPrimary,
            fontWeight: FontWeight.w600,
          ),
        ),
        backgroundColor: AppTheme.white,
        elevation: 0,
        iconTheme: const IconThemeData(color: AppTheme.textPrimary),
        actions: [
          IconButton(
            onPressed: _loadSessions,
            icon: const Icon(Icons.refresh),
            tooltip: 'Refresh',
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? _buildErrorState()
              : _buildSessionsList(),
      bottomNavigationBar: _sessions != null && _sessions!.isNotEmpty
          ? _buildBottomActions()
          : null,
    );
  }

  Widget _buildErrorState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.error_outline,
            size: 64,
            color: AppTheme.error.withOpacity(0.5),
          ),
          const SizedBox(height: 16),
          Text(
            'Failed to load sessions',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w500,
              color: AppTheme.textSecondary,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            _error!,
            textAlign: TextAlign.center,
            style: TextStyle(
              fontSize: 14,
              color: AppTheme.textSecondary,
            ),
          ),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: _loadSessions,
            icon: const Icon(Icons.refresh),
            label: const Text('Retry'),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppTheme.primary,
              foregroundColor: AppTheme.white,
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSessionsList() {
    if (_sessions == null || _sessions!.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.devices_other,
              size: 64,
              color: AppTheme.textSecondary.withOpacity(0.5),
            ),
            const SizedBox(height: 16),
            Text(
              'No active sessions',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w500,
                color: AppTheme.textSecondary,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'You are not signed in on any other devices',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 14,
                color: AppTheme.textSecondary,
              ),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _loadSessions,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _sessions!.length,
        itemBuilder: (context, index) {
          final session = _sessions![index];
          final isCurrentSession = index == 0; // Assume first is current
          
          return Card(
            margin: const EdgeInsets.only(bottom: 12),
            elevation: 2,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
              side: isCurrentSession
                  ? BorderSide(color: AppTheme.primary, width: 2)
                  : BorderSide.none,
            ),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(
                        _getDeviceIcon(session),
                        size: 24,
                        color: AppTheme.primary,
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Expanded(
                                  child: Text(
                                    _getDeviceDisplayName(session),
                                    style: const TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.w600,
                                      color: AppTheme.textPrimary,
                                    ),
                                  ),
                                ),
                                if (isCurrentSession) ...[
                                  Container(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 8, vertical: 4),
                                    decoration: BoxDecoration(
                                      color: AppTheme.primary,
                                      borderRadius: BorderRadius.circular(12),
                                    ),
                                    child: const Text(
                                      'Current',
                                      style: TextStyle(
                                        fontSize: 12,
                                        color: AppTheme.white,
                                        fontWeight: FontWeight.w500,
                                      ),
                                    ),
                                  ),
                                ],
                              ],
                            ),
                            const SizedBox(height: 4),
                            _buildTrustIndicator(session),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  _buildSessionDetails(session),
                  if (!isCurrentSession) ...[
                    const SizedBox(height: 12),
                    SizedBox(
                      width: double.infinity,
                      child: OutlinedButton.icon(
                        onPressed: () => _revokeSession(
                          session.sessionId,
                          _getDeviceDisplayName(session),
                        ),
                        icon: const Icon(Icons.logout, size: 16),
                        label: const Text('Sign Out'),
                        style: OutlinedButton.styleFrom(
                          foregroundColor: AppTheme.error,
                          side: BorderSide(color: AppTheme.error),
                          padding: const EdgeInsets.symmetric(vertical: 8),
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildSessionDetails(SessionInfo session) {
    return Column(
      children: [
        _buildDetailRow(
          Icons.location_on_outlined,
          'Location',
          _getLocationInfo(session),
        ),
        const SizedBox(height: 8),
        _buildDetailRow(
          Icons.access_time,
          'Last Active',
          _getLastActiveTime(session),
        ),
        const SizedBox(height: 8),
        _buildDetailRow(
          Icons.smartphone_outlined,
          'OS Version',
          session.deviceInfo.osVersion,
        ),
        if (session.isRemembered) ...[
          const SizedBox(height: 8),
          _buildDetailRow(
            Icons.bookmark_outline,
            'Session Type',
            'Remember Me Enabled',
          ),
        ],
      ],
    );
  }

  Widget _buildDetailRow(IconData icon, String label, String value) {
    return Row(
      children: [
        Icon(
          icon,
          size: 16,
          color: AppTheme.textSecondary,
        ),
        const SizedBox(width: 8),
        Text(
          '$label:',
          style: TextStyle(
            fontSize: 14,
            color: AppTheme.textSecondary,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            value,
            style: const TextStyle(
              fontSize: 14,
              color: AppTheme.textPrimary,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildBottomActions() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.white,
        border: Border(
          top: BorderSide(
            color: AppTheme.borderLight,
            width: 1,
          ),
        ),
      ),
      child: SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            SizedBox(
              width: double.infinity,
              child: ElevatedButton.icon(
                onPressed: _signOutAllDevices,
                icon: const Icon(Icons.logout),
                label: const Text('Sign Out All Devices'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppTheme.error,
                  foregroundColor: AppTheme.white,
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'This will sign you out from all devices including this one',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 12,
                color: AppTheme.textSecondary,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
