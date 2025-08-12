import '../constants/app_constants.dart';

/// Service for mapping backend error codes to user-friendly messages
class ErrorMessageService {
  static const Map<String, String> _errorCodeMap = {
    // Authentication Errors
    'AUTH_001': 'Invalid email or password. Please check your credentials and try again.',
    'AUTH_002': 'Your account has been temporarily locked due to multiple failed login attempts.',
    'AUTH_003': 'This email address is already registered. Please use a different email or try logging in.',
    'AUTH_004': 'Your account is not verified. Please check your email for verification instructions.',
    'AUTH_005': 'Your session has expired. Please log in again.',
    'AUTH_006': 'Two-factor authentication code is required.',
    'AUTH_007': 'Invalid two-factor authentication code. Please try again.',
    'AUTH_008': 'CAPTCHA verification is required to continue.',
    'AUTH_009': 'Invalid CAPTCHA. Please try again.',
    'AUTH_010': 'Password reset token is invalid or has expired.',
    'AUTH_011': 'Email verification token is invalid or has expired.',
    'AUTH_012': 'Phone verification code is invalid or has expired.',
    'AUTH_013': 'Password does not meet security requirements.',
    'AUTH_014': 'Current password is incorrect.',
    
    // User Errors
    'USER_001': 'User not found. Please check your credentials.',
    'USER_002': 'User profile is incomplete. Please complete your profile.',
    'USER_003': 'Email address is required.',
    'USER_004': 'Phone number is required.',
    'USER_005': 'Invalid email format.',
    'USER_006': 'Invalid phone number format.',
    'USER_007': 'User account is suspended.',
    'USER_008': 'User account is deactivated.',
    
    // Permission Errors
    'PERM_001': 'You do not have permission to access this resource.',
    'PERM_002': 'Admin privileges required.',
    'PERM_003': 'Provider verification required.',
    'PERM_004': 'Email verification required to access this feature.',
    'PERM_005': 'Phone verification required to access this feature.',
    'PERM_006': 'KYC verification required.',
    'PERM_007': 'Account verification pending.',
    
    // Service Errors
    'SERV_001': 'Service not found.',
    'SERV_002': 'Service is not available in your area.',
    'SERV_003': 'Service provider is not available.',
    'SERV_004': 'Service booking limit exceeded.',
    'SERV_005': 'Service is temporarily unavailable.',
    
    // Booking Errors
    'BOOK_001': 'Booking not found.',
    'BOOK_002': 'Cannot cancel booking. Cancellation period has expired.',
    'BOOK_003': 'Cannot modify booking. Service is already in progress.',
    'BOOK_004': 'Booking time slot is no longer available.',
    'BOOK_005': 'Insufficient funds for booking.',
    'BOOK_006': 'Provider has declined the booking.',
    'BOOK_007': 'Booking is already completed.',
    'BOOK_008': 'Cannot book same service within 24 hours.',
    
    // Payment Errors
    'PAY_001': 'Payment failed. Please try again.',
    'PAY_002': 'Insufficient funds in your wallet.',
    'PAY_003': 'Payment method not supported.',
    'PAY_004': 'Payment processing error. Please contact support.',
    'PAY_005': 'Refund processing failed.',
    'PAY_006': 'Transaction not found.',
    'PAY_007': 'Payment already processed.',
    'PAY_008': 'Invalid payment amount.',
    'PAY_009': 'Payment gateway error.',
    'PAY_010': 'Card verification failed.',
    
    // Wallet Errors
    'WALLET_001': 'Wallet not found.',
    'WALLET_002': 'Insufficient wallet balance.',
    'WALLET_003': 'Wallet transaction limit exceeded.',
    'WALLET_004': 'Invalid withdrawal amount.',
    'WALLET_005': 'Withdrawal not allowed to this account.',
    'WALLET_006': 'Daily transaction limit exceeded.',
    
    // File Upload Errors
    'FILE_001': 'File size exceeds maximum limit.',
    'FILE_002': 'File format not supported.',
    'FILE_003': 'File upload failed. Please try again.',
    'FILE_004': 'Too many files uploaded.',
    'FILE_005': 'File contains malicious content.',
    
    // Validation Errors
    'VAL_001': 'Required field is missing.',
    'VAL_002': 'Invalid input format.',
    'VAL_003': 'Input exceeds maximum length.',
    'VAL_004': 'Input is too short.',
    'VAL_005': 'Invalid date format.',
    'VAL_006': 'Date is in the past.',
    'VAL_007': 'Date is too far in the future.',
    
    // Network Errors
    'NET_001': 'Network connection failed. Please check your internet connection.',
    'NET_002': 'Request timeout. Please try again.',
    'NET_003': 'Server is temporarily unavailable.',
    'NET_004': 'Too many requests. Please wait and try again.',
    
    // KYC Errors
    'KYC_001': 'KYC verification failed.',
    'KYC_002': 'Invalid identification document.',
    'KYC_003': 'Document quality is too poor.',
    'KYC_004': 'Document has expired.',
    'KYC_005': 'KYC verification is pending.',
    'KYC_006': 'KYC verification was rejected.',
    
    // Location Errors
    'LOC_001': 'Location service is disabled.',
    'LOC_002': 'Unable to determine your location.',
    'LOC_003': 'Service not available in your location.',
    'LOC_004': 'Invalid address format.',
    
    // Notification Errors
    'NOTIF_001': 'Failed to send notification.',
    'NOTIF_002': 'Notification preferences not set.',
    'NOTIF_003': 'Push notification permission denied.',
    
    // General System Errors
    'SYS_001': 'An unexpected error occurred. Please try again.',
    'SYS_002': 'System is under maintenance. Please try again later.',
    'SYS_003': 'Feature is temporarily disabled.',
    'SYS_004': 'Database connection error.',
    'SYS_005': 'External service unavailable.',
  };

  /// Get user-friendly error message for a given error code
  static String getErrorMessage(String? errorCode, {String? fallbackMessage}) {
    if (errorCode == null || errorCode.isEmpty) {
      return fallbackMessage ?? AppConstants.serverErrorMessage;
    }

    return _errorCodeMap[errorCode] ?? fallbackMessage ?? _getGenericErrorMessage(errorCode);
  }

  /// Get generic error message based on error code prefix
  static String _getGenericErrorMessage(String errorCode) {
    final prefix = errorCode.split('_').first;
    
    switch (prefix) {
      case 'AUTH':
        return 'Authentication error. Please try logging in again.';
      case 'USER':
        return 'User account error. Please contact support if the issue persists.';
      case 'PERM':
        return 'You do not have permission to perform this action.';
      case 'SERV':
        return 'Service error. Please try again later.';
      case 'BOOK':
        return 'Booking error. Please try again or contact support.';
      case 'PAY':
        return 'Payment error. Please check your payment method and try again.';
      case 'WALLET':
        return 'Wallet error. Please check your balance and try again.';
      case 'FILE':
        return 'File upload error. Please check your file and try again.';
      case 'VAL':
        return 'Invalid input. Please check your information and try again.';
      case 'NET':
        return AppConstants.networkErrorMessage;
      case 'KYC':
        return 'Identity verification error. Please contact support.';
      case 'LOC':
        return 'Location error. Please check your location settings.';
      case 'NOTIF':
        return 'Notification error. Please check your notification settings.';
      case 'SYS':
        return AppConstants.serverErrorMessage;
      default:
        return AppConstants.serverErrorMessage;
    }
  }

  /// Get contextual action message for error codes
  static String? getActionMessage(String? errorCode) {
    if (errorCode == null) return null;

    switch (errorCode) {
      case 'AUTH_001':
        return 'Try resetting your password if you\'ve forgotten it.';
      case 'AUTH_002':
        return 'Please wait 15 minutes before trying again.';
      case 'AUTH_003':
        return 'Use the "Sign In" option instead.';
      case 'AUTH_004':
        return 'Check your email and click the verification link.';
      case 'AUTH_005':
        return 'Please sign in again to continue.';
      case 'NET_001':
        return 'Check your internet connection and try again.';
      case 'PAY_002':
        return 'Please add funds to your wallet.';
      case 'BOOK_004':
        return 'Please select a different time slot.';
      case 'FILE_001':
        return 'Please choose a smaller file.';
      case 'LOC_001':
        return 'Enable location services in your device settings.';
      default:
        return null;
    }
  }

  /// Check if error code indicates a retry-able error
  static bool isRetryable(String? errorCode) {
    if (errorCode == null) return true;

    final nonRetryableCodes = [
      'AUTH_003', // Email already exists
      'AUTH_004', // Account not verified
      'AUTH_007', // Invalid 2FA code
      'AUTH_013', // Password requirements
      'USER_007', // Account suspended
      'USER_008', // Account deactivated
      'PERM_001', // No permission
      'PERM_002', // Admin required
      'FILE_002', // Unsupported format
      'FILE_005', // Malicious content
      'KYC_006', // KYC rejected
    ];

    return !nonRetryableCodes.contains(errorCode);
  }

  /// Get severity level for error code
  static ErrorSeverity getSeverity(String? errorCode) {
    if (errorCode == null) return ErrorSeverity.error;

    final warningSeverityCodes = [
      'AUTH_004', // Account not verified
      'AUTH_006', // 2FA required
      'AUTH_008', // CAPTCHA required
      'PERM_004', // Email verification required
      'PERM_005', // Phone verification required
      'PERM_006', // KYC required
      'WALLET_002', // Insufficient balance
      'FILE_001', // File too large
      'LOC_001', // Location disabled
    ];

    final infoCodes = [
      'AUTH_005', // Session expired
      'BOOK_007', // Already completed
      'PAY_007', // Already processed
      'KYC_005', // KYC pending
    ];

    if (infoCodes.contains(errorCode)) {
      return ErrorSeverity.info;
    } else if (warningSeverityCodes.contains(errorCode)) {
      return ErrorSeverity.warning;
    } else {
      return ErrorSeverity.error;
    }
  }
}

/// Error severity levels
enum ErrorSeverity {
  info,
  warning,
  error,
}
