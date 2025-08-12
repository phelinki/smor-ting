import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../../../core/theme/app_theme.dart';

/// Widget for displaying enhanced login outcomes like cooldowns, CAPTCHA, and 2FA
class EnhancedLoginOutcomes extends StatefulWidget {
  final bool showCooldown;
  final bool showCaptcha;
  final bool showTwoFactor;
  final int cooldownSeconds;
  final int remainingAttempts;
  final Function(String) onCaptchaCompleted;
  final Function(String) onTwoFactorSubmitted;

  const EnhancedLoginOutcomes({
    super.key,
    this.showCooldown = false,
    this.showCaptcha = false,
    this.showTwoFactor = false,
    this.cooldownSeconds = 0,
    this.remainingAttempts = 0,
    required this.onCaptchaCompleted,
    required this.onTwoFactorSubmitted,
  });

  @override
  State<EnhancedLoginOutcomes> createState() => _EnhancedLoginOutcomesState();
}

class _EnhancedLoginOutcomesState extends State<EnhancedLoginOutcomes> {
  Timer? _cooldownTimer;
  int _remainingCooldownSeconds = 0;
  final List<TextEditingController> _twoFactorControllers = List.generate(
    6,
    (index) => TextEditingController(),
  );
  final List<FocusNode> _twoFactorFocusNodes = List.generate(
    6,
    (index) => FocusNode(),
  );

  @override
  void initState() {
    super.initState();
    if (widget.showCooldown && widget.cooldownSeconds > 0) {
      _startCooldownTimer();
    }
  }

  @override
  void didUpdateWidget(EnhancedLoginOutcomes oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.showCooldown && widget.cooldownSeconds > 0 && !oldWidget.showCooldown) {
      _startCooldownTimer();
    } else if (!widget.showCooldown && _cooldownTimer != null) {
      _cooldownTimer?.cancel();
    }
  }

  void _startCooldownTimer() {
    _remainingCooldownSeconds = widget.cooldownSeconds;
    _cooldownTimer = Timer.periodic(const Duration(seconds: 1), (timer) {
      setState(() {
        if (_remainingCooldownSeconds > 0) {
          _remainingCooldownSeconds--;
        } else {
          timer.cancel();
        }
      });
    });
  }

  void _onTwoFactorChanged(String value, int index) {
    if (value.isNotEmpty) {
      if (index < 5) {
        _twoFactorFocusNodes[index + 1].requestFocus();
      } else {
        // Last digit entered, submit
        _submitTwoFactorCode();
      }
    }
  }

  void _submitTwoFactorCode() {
    final code = _twoFactorControllers.map((c) => c.text).join();
    if (code.length == 6) {
      widget.onTwoFactorSubmitted(code);
    }
  }

  String _formatCooldownTime(int seconds) {
    final minutes = seconds ~/ 60;
    final remainingSeconds = seconds % 60;
    if (minutes > 0) {
      return '${minutes}m ${remainingSeconds}s';
    }
    return '${remainingSeconds}s';
  }

  @override
  void dispose() {
    _cooldownTimer?.cancel();
    for (final controller in _twoFactorControllers) {
      controller.dispose();
    }
    for (final focusNode in _twoFactorFocusNodes) {
      focusNode.dispose();
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (!widget.showCooldown && !widget.showCaptcha && !widget.showTwoFactor) {
      return const SizedBox.shrink();
    }

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Card(
        elevation: 4,
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (widget.showCooldown) _buildCooldownWidget(),
              if (widget.showCaptcha) _buildCaptchaWidget(),
              if (widget.showTwoFactor) _buildTwoFactorWidget(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildCooldownWidget() {
    return Column(
      children: [
        Icon(
          Icons.lock_clock,
          size: 48,
          color: AppTheme.error,
        ),
        const SizedBox(height: 16),
        Text(
          'Account Temporarily Locked',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w600,
            color: AppTheme.error,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Too many failed login attempts. Try again in ${_formatCooldownTime(_remainingCooldownSeconds)}.',
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: 14,
            color: Colors.grey[600],
          ),
        ),
        const SizedBox(height: 16),
        LinearProgressIndicator(
          value: _remainingCooldownSeconds > 0 
              ? (_remainingCooldownSeconds / widget.cooldownSeconds) 
              : 1.0,
          backgroundColor: Colors.grey[300],
          valueColor: AlwaysStoppedAnimation<Color>(AppTheme.error),
        ),
      ],
    );
  }

  Widget _buildCaptchaWidget() {
    return Column(
      children: [
        Icon(
          Icons.security,
          size: 48,
          color: AppTheme.warning,
        ),
        const SizedBox(height: 16),
        Text(
          'Security Verification',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w600,
            color: AppTheme.warning,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Please complete the security verification below. ${widget.remainingAttempts} attempts remaining.',
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: 14,
            color: Colors.grey[600],
          ),
        ),
        const SizedBox(height: 16),
        // Placeholder for CAPTCHA widget
        Container(
          height: 120,
          decoration: BoxDecoration(
            border: Border.all(color: Colors.grey[300]!),
            borderRadius: BorderRadius.circular(8),
          ),
          child: const Center(
            child: Text(
              'CAPTCHA Widget\n(Implementation Pending)',
              textAlign: TextAlign.center,
              style: TextStyle(color: Colors.grey),
            ),
          ),
        ),
        const SizedBox(height: 16),
        ElevatedButton(
          onPressed: () {
            // Mock CAPTCHA completion
            widget.onCaptchaCompleted('mock_captcha_token');
          },
          child: const Text('Verify'),
        ),
      ],
    );
  }

  Widget _buildTwoFactorWidget() {
    return Column(
      children: [
        Icon(
          Icons.phonelink_lock,
          size: 48,
          color: AppTheme.secondaryBlue,
        ),
        const SizedBox(height: 16),
        Text(
          'Two-Factor Authentication',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w600,
            color: AppTheme.secondaryBlue,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Enter the 6-digit code from your authenticator app',
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: 14,
            color: Colors.grey[600],
          ),
        ),
        const SizedBox(height: 24),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: List.generate(6, (index) {
            return SizedBox(
              width: 40,
              height: 50,
              child: TextField(
                controller: _twoFactorControllers[index],
                focusNode: _twoFactorFocusNodes[index],
                keyboardType: TextInputType.number,
                textAlign: TextAlign.center,
                maxLength: 1,
                decoration: InputDecoration(
                  counterText: '',
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                inputFormatters: [
                  FilteringTextInputFormatter.digitsOnly,
                ],
                onChanged: (value) => _onTwoFactorChanged(value, index),
              ),
            );
          }),
        ),
        const SizedBox(height: 16),
        TextButton(
          onPressed: () {
            // TODO: Implement resend 2FA code
          },
          child: const Text('Resend Code'),
        ),
      ],
    );
  }
}
