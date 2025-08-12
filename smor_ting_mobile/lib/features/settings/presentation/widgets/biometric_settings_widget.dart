import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:local_auth/local_auth.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/services/enhanced_auth_service.dart';

/// Widget for managing biometric authentication settings
class BiometricSettingsWidget extends ConsumerStatefulWidget {
  final String userEmail;

  const BiometricSettingsWidget({
    super.key,
    required this.userEmail,
  });

  @override
  ConsumerState<BiometricSettingsWidget> createState() => _BiometricSettingsWidgetState();
}

class _BiometricSettingsWidgetState extends ConsumerState<BiometricSettingsWidget> {
  bool _isLoading = false;
  bool _biometricAvailable = false;
  bool _biometricEnabled = false;

  @override
  void initState() {
    super.initState();
    _checkBiometricAvailability();
  }

  Future<void> _checkBiometricAvailability() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final isAvailable = await authService.canUseBiometrics();
      
      if (isAvailable) {
        final isEnabled = await authService.isBiometricEnabled(widget.userEmail);
        
        setState(() {
          _biometricAvailable = true;
          _biometricEnabled = isEnabled;
        });
      }
    } catch (e) {
      // Silently handle errors - biometric won't be shown if not available
    }
  }

  Future<void> _toggleBiometric(bool enabled) async {
    setState(() {
      _isLoading = true;
    });

    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      bool success;
      
      success = await authService.setBiometricEnabled(widget.userEmail, enabled);

      if (success) {
        setState(() {
          _biometricEnabled = enabled;
        });
        
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(
                enabled 
                  ? 'Biometric authentication enabled'
                  : 'Biometric authentication disabled'
              ),
              backgroundColor: AppTheme.successGreen,
            ),
          );
        }
      } else {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(
                enabled 
                  ? 'Failed to enable biometric authentication'
                  : 'Failed to disable biometric authentication'
              ),
              backgroundColor: AppTheme.error,
            ),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error: ${e.toString()}'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (!_biometricAvailable) {
      return const SizedBox.shrink();
    }

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  Icons.fingerprint,
                  color: AppTheme.secondaryBlue,
                  size: 24,
                ),
                const SizedBox(width: 12),
                const Expanded(
                  child: Text(
                    'Biometric Authentication',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                if (_isLoading)
                  const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                else
                  Switch(
                    value: _biometricEnabled,
                    onChanged: _toggleBiometric,
                    activeColor: AppTheme.secondaryBlue,
                  ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Use fingerprint or face recognition to quickly unlock the app',
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
