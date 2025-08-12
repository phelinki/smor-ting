import 'package:flutter/material.dart';

import '../theme/app_theme.dart';
import 'error_message_service.dart';

/// Service for showing consistent messages throughout the app
class MessageService {
  static const Duration _defaultDuration = Duration(seconds: 4);
  static const Duration _longDuration = Duration(seconds: 6);
  static const Duration _shortDuration = Duration(seconds: 2);

  /// Show error message with appropriate styling
  static void showError(
    BuildContext context, {
    String? message,
    String? errorCode,
    bool canRetry = true,
    VoidCallback? onRetry,
    Duration? duration,
  }) {
    final errorMessage = ErrorMessageService.getErrorMessage(errorCode, fallbackMessage: message);
    final actionMessage = ErrorMessageService.getActionMessage(errorCode);
    final isRetryable = ErrorMessageService.isRetryable(errorCode) && canRetry;
    
    _showSnackBar(
      context,
      message: errorMessage,
      backgroundColor: AppTheme.error,
      icon: Icons.error_outline,
      actionLabel: isRetryable ? 'Retry' : null,
      onAction: isRetryable ? onRetry : null,
      duration: duration ?? _defaultDuration,
      actionMessage: actionMessage,
    );
  }

  /// Show warning message
  static void showWarning(
    BuildContext context, {
    required String message,
    String? actionLabel,
    VoidCallback? onAction,
    Duration? duration,
  }) {
    _showSnackBar(
      context,
      message: message,
      backgroundColor: AppTheme.warning,
      icon: Icons.warning_outlined,
      actionLabel: actionLabel,
      onAction: onAction,
      duration: duration ?? _defaultDuration,
    );
  }

  /// Show success message
  static void showSuccess(
    BuildContext context, {
    required String message,
    String? actionLabel,
    VoidCallback? onAction,
    Duration? duration,
  }) {
    _showSnackBar(
      context,
      message: message,
      backgroundColor: AppTheme.successGreen,
      icon: Icons.check_circle_outline,
      actionLabel: actionLabel,
      onAction: onAction,
      duration: duration ?? _shortDuration,
    );
  }

  /// Show info message
  static void showInfo(
    BuildContext context, {
    required String message,
    String? actionLabel,
    VoidCallback? onAction,
    Duration? duration,
  }) {
    _showSnackBar(
      context,
      message: message,
      backgroundColor: AppTheme.secondaryBlue,
      icon: Icons.info_outline,
      actionLabel: actionLabel,
      onAction: onAction,
      duration: duration ?? _defaultDuration,
    );
  }

  /// Show message based on error severity
  static void showErrorBySeverity(
    BuildContext context, {
    required String? errorCode,
    String? fallbackMessage,
    bool canRetry = true,
    VoidCallback? onRetry,
  }) {
    final severity = ErrorMessageService.getSeverity(errorCode);
    final message = ErrorMessageService.getErrorMessage(errorCode, fallbackMessage: fallbackMessage);
    
    switch (severity) {
      case ErrorSeverity.info:
        showInfo(context, message: message);
        break;
      case ErrorSeverity.warning:
        showWarning(context, message: message);
        break;
      case ErrorSeverity.error:
        showError(
          context,
          message: message,
          errorCode: errorCode,
          canRetry: canRetry,
          onRetry: onRetry,
        );
        break;
    }
  }

  /// Show verification required message
  static void showVerificationRequired(
    BuildContext context, {
    required bool emailRequired,
    required bool phoneRequired,
    VoidCallback? onVerify,
  }) {
    String message;
    if (emailRequired && phoneRequired) {
      message = 'Email and phone verification required to access this feature.';
    } else if (emailRequired) {
      message = 'Email verification required to access this feature.';
    } else {
      message = 'Phone verification required to access this feature.';
    }

    showWarning(
      context,
      message: message,
      actionLabel: 'Verify',
      onAction: onVerify,
      duration: _longDuration,
    );
  }

  /// Show loading message (usually as inline widget, but can be snackbar)
  static void showLoading(
    BuildContext context, {
    String message = 'Loading...',
  }) {
    _showSnackBar(
      context,
      message: message,
      backgroundColor: AppTheme.gray,
      icon: null, // Will show loading indicator
      showLoading: true,
      duration: const Duration(seconds: 30), // Long duration for loading
    );
  }

  /// Show network error with appropriate retry
  static void showNetworkError(
    BuildContext context, {
    VoidCallback? onRetry,
  }) {
    showError(
      context,
      errorCode: 'NET_001',
      canRetry: true,
      onRetry: onRetry,
    );
  }

  /// Show authentication error
  static void showAuthError(
    BuildContext context, {
    String? errorCode,
    VoidCallback? onLogin,
  }) {
    showError(
      context,
      errorCode: errorCode ?? 'AUTH_005',
      canRetry: false,
      onRetry: onLogin,
    );
  }

  /// Internal method to show snackbar with consistent styling
  static void _showSnackBar(
    BuildContext context, {
    required String message,
    required Color backgroundColor,
    IconData? icon,
    String? actionLabel,
    VoidCallback? onAction,
    Duration? duration,
    String? actionMessage,
    bool showLoading = false,
  }) {
    // Dismiss any existing snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    final snackBar = SnackBar(
      content: Row(
        children: [
          if (showLoading)
            SizedBox(
              width: 20,
              height: 20,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
              ),
            )
          else if (icon != null)
            Icon(
              icon,
              color: Colors.white,
              size: 20,
            ),
          if (icon != null || showLoading) const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  message,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                if (actionMessage != null) ...[
                  const SizedBox(height: 4),
                  Text(
                    actionMessage,
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.8),
                      fontSize: 12,
                    ),
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
      backgroundColor: backgroundColor,
      duration: duration ?? _defaultDuration,
      behavior: SnackBarBehavior.floating,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      margin: const EdgeInsets.all(16),
      action: (actionLabel != null && onAction != null)
          ? SnackBarAction(
              label: actionLabel,
              textColor: Colors.white,
              onPressed: onAction,
            )
          : null,
    );

    ScaffoldMessenger.of(context).showSnackBar(snackBar);
  }

  /// Dismiss current snackbar
  static void dismiss(BuildContext context) {
    ScaffoldMessenger.of(context).hideCurrentSnackBar();
  }
}

/// Extension to make it easier to show messages from any widget
extension MessageServiceExtension on BuildContext {
  void showError({
    String? message,
    String? errorCode,
    bool canRetry = true,
    VoidCallback? onRetry,
  }) {
    MessageService.showError(
      this,
      message: message,
      errorCode: errorCode,
      canRetry: canRetry,
      onRetry: onRetry,
    );
  }

  void showSuccess({required String message}) {
    MessageService.showSuccess(this, message: message);
  }

  void showWarning({required String message}) {
    MessageService.showWarning(this, message: message);
  }

  void showInfo({required String message}) {
    MessageService.showInfo(this, message: message);
  }

  void showNetworkError({VoidCallback? onRetry}) {
    MessageService.showNetworkError(this, onRetry: onRetry);
  }
}
