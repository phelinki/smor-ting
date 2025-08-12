import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/services/enhanced_auth_service.dart';

/// Widget for biometric quick unlock on app launch
class BiometricQuickUnlock extends ConsumerStatefulWidget {
  final String userEmail;
  final VoidCallback onSuccess;
  final VoidCallback onCancel;

  const BiometricQuickUnlock({
    super.key,
    required this.userEmail,
    required this.onSuccess,
    required this.onCancel,
  });

  @override
  ConsumerState<BiometricQuickUnlock> createState() => _BiometricQuickUnlockState();
}

class _BiometricQuickUnlockState extends ConsumerState<BiometricQuickUnlock> {
  bool _isLoading = false;
  bool _showBiometric = false;

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
        
        if (mounted) {
          setState(() {
            _showBiometric = isEnabled;
          });
        }
      }
    } catch (e) {
      // Silently handle errors - will show password login instead
    }
  }

  Future<void> _authenticateWithBiometrics() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final result = await authService.authenticateWithBiometrics(widget.userEmail);
      
      if (result.success) {
        widget.onSuccess();
      } else {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(result.message ?? 'Authentication failed'),
              backgroundColor: AppTheme.error,
            ),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Authentication failed: ${e.toString()}'),
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
    if (!_showBiometric) {
      return const SizedBox.shrink();
    }

    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            'Quick Unlock',
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.w600,
              color: AppTheme.secondaryBlue,
            ),
          ),
          const SizedBox(height: 16),
          Text(
            'Use biometric authentication to unlock',
            style: TextStyle(
              fontSize: 14,
              color: Colors.grey[600],
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          GestureDetector(
            onTap: _isLoading ? null : _authenticateWithBiometrics,
            child: Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                              color: AppTheme.secondaryBlue.withOpacity(0.1),
              shape: BoxShape.circle,
              border: Border.all(
                color: AppTheme.secondaryBlue,
                width: 2,
              ),
              ),
              child: _isLoading
                  ? const CircularProgressIndicator()
                  : Icon(
                      Icons.fingerprint,
                      size: 40,
                      color: AppTheme.secondaryBlue,
                    ),
            ),
          ),
          const SizedBox(height: 24),
          TextButton(
            onPressed: widget.onCancel,
            child: const Text('Use Password Instead'),
          ),
        ],
      ),
    );
  }
}
