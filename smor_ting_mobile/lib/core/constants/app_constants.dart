class AppConstants {
  // App Information
  static const String appName = 'Smor-Ting';
  static const String appVersion = '1.0.0';
  static const String appDescription = 'Handyman and Service Marketplace for Liberia';
  
  // API Configuration
  // Note: API configuration has been moved to ApiConfig class in api_config.dart
  // Use ApiConfig.baseUrl and ApiConfig.apiBaseUrl instead
  
  // Liberia-specific settings
  static const String defaultCountry = 'Liberia';
  static const String defaultCurrency = 'USD';
  static const String defaultLanguage = 'en';
  static const String defaultTimezone = 'Africa/Monrovia';
  
  // Payment Methods
  static const List<String> supportedPaymentMethods = [
    'flutterwave',
    'orange_money',
    'mtn_mobile_money',
    'lonestar_cell_mtn',
    'bank_transfer',
  ];
  
  // Service Categories
  static const List<String> serviceCategories = [
    'Electrical',
    'Plumbing',
    'Cleaning',
    'Yardwork',
    'Carpentry',
    'Painting',
    'HVAC',
    'Security',
    'Moving',
    'Other',
  ];
  
  // App Features
  static const bool enableOfflineMode = true;
  static const bool enableBiometricAuth = true;
  static const bool enablePushNotifications = true;
  static const bool enableLocationServices = true;
  
  // Cache Configuration
  static const int cacheExpirationHours = 24;
  static const int maxCacheSize = 100; // MB
  
  // Security
  static const int sessionTimeoutMinutes = 30;
  static const int maxLoginAttempts = 5;
  static const int lockoutDurationMinutes = 15;
  
  // File Upload
  static const int maxImageSize = 5; // MB
  static const int maxFileSize = 10; // MB
  static const List<String> allowedImageFormats = ['jpg', 'jpeg', 'png', 'webp'];
  static const List<String> allowedFileFormats = ['pdf', 'doc', 'docx'];
  
  // Pagination
  static const int defaultPageSize = 20;
  static const int maxPageSize = 100;
  
  // Timeouts
  static const int connectionTimeoutSeconds = 30;
  static const int receiveTimeoutSeconds = 60;
  
  // Error Messages
  static const String networkErrorMessage = 'Please check your internet connection and try again.';
  static const String serverErrorMessage = 'Something went wrong. Please try again later.';
  static const String offlineErrorMessage = 'You are currently offline. Some features may be limited.';
  static const String authenticationErrorMessage = 'Please log in to continue.';
  
  // Success Messages
  static const String bookingSuccessMessage = 'Your booking has been confirmed!';
  static const String paymentSuccessMessage = 'Payment completed successfully!';
  static const String profileUpdateMessage = 'Profile updated successfully!';
  
  // Validation Messages
  static const String requiredFieldMessage = 'This field is required.';
  static const String invalidEmailMessage = 'Please enter a valid email address.';
  static const String invalidPhoneMessage = 'Please enter a valid Liberian phone number.';
  static const String phoneFormatHint = 'Format: +231xxxxxxxxx or 77xxxxxxx (9 digits)';
  static const String passwordTooShortMessage = 'Password must be at least 8 characters.';
  static const String passwordMismatchMessage = 'Passwords do not match.';
  
  // Liberia Phone Number Patterns
  static const String liberiaPhonePattern = r'^(\+231|231)?[0-9]{9}$';
  static const String liberiaPhonePrefix = '+231';
  
  // Service Status
  static const List<String> serviceStatuses = [
    'pending',
    'accepted',
    'in_progress',
    'completed',
    'cancelled',
  ];
  
  // User Types
  static const String userTypeCustomer = 'customer';
  static const String userTypeAgent = 'agent';
  static const String userTypeAdmin = 'admin';
  
  // Verification Status
  static const String verificationStatusPending = 'pending';
  static const String verificationStatusApproved = 'approved';
  static const String verificationStatusRejected = 'rejected';
  
  // Notification Types
  static const String notificationTypeBooking = 'booking';
  static const String notificationTypePayment = 'payment';
  static const String notificationTypeStatus = 'status';
  static const String notificationTypeMessage = 'message';
  
  // Storage Keys
  static const String authTokenKey = 'auth_token';
  static const String userDataKey = 'user_data';
  static const String appSettingsKey = 'app_settings';
  static const String offlineDataKey = 'offline_data';
  
  // Animation Durations
  static const Duration shortAnimationDuration = Duration(milliseconds: 200);
  static const Duration mediumAnimationDuration = Duration(milliseconds: 300);
  static const Duration longAnimationDuration = Duration(milliseconds: 500);
  
  // UI Constants
  static const double defaultPadding = 16.0;
  static const double smallPadding = 8.0;
  static const double largePadding = 24.0;
  static const double defaultRadius = 12.0;
  static const double smallRadius = 8.0;
  static const double largeRadius = 16.0;
  
  // Map Configuration
  static const double defaultZoom = 15.0;
  static const double minZoom = 10.0;
  static const double maxZoom = 20.0;
  static const double searchRadiusKm = 50.0;
  
  // Liberia Coordinates (Monrovia)
  static const double defaultLatitude = 6.3004;
  static const double defaultLongitude = -10.7969;
} 