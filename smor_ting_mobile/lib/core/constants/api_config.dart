import 'package:flutter/foundation.dart';

class ApiConfig {
  // Environment-based API URLs
  static const String _devBaseUrl = 'http://127.0.0.1:8080/api/v1';
  static const String _stagingBaseUrl = 'https://api.smor-ting.com/api/v1';
  static const String _productionBaseUrl = 'https://api.smor-ting.com/api/v1';
  
  // Environment selection
  // - In release/TestFlight/App Store builds (kReleaseMode), use production
  // - In debug/profile (simulator/local dev), use development
  static Environment get _currentEnvironment =>
      kReleaseMode ? Environment.production : Environment.development;
  
  static String get baseUrl {
    switch (_currentEnvironment) {
      case Environment.development:
        return _devBaseUrl;
      case Environment.staging:
        return _stagingBaseUrl;
      case Environment.production:
        return _productionBaseUrl;
    }
  }
  
  static String get environmentName {
    switch (_currentEnvironment) {
      case Environment.development:
        return 'Development';
      case Environment.staging:
        return 'Staging';
      case Environment.production:
        return 'Production';
    }
  }
  
  // API Endpoints
  static const String authRegister = '/auth/register';
  static const String authLogin = '/auth/login';
  static const String authVerifyOTP = '/auth/verify-otp';
  static const String authResendOTP = '/auth/resend-otp';
  
  static const String userProfile = '/users/profile';
  static const String userUpdateProfile = '/users/profile';
  
  static const String servicesCategories = '/services/categories';
  static const String servicesList = '/services';
  static const String serviceDetails = '/services';
  
  static const String bookingsCreate = '/bookings';
  static const String bookingsList = '/bookings';
  static const String bookingsDetails = '/bookings';
  
  static const String walletBalance = '/wallet/balance';
  static const String walletTransactions = '/wallet/transactions';
  
  static const String syncData = '/sync/data';
  static const String syncUnsynced = '/sync/unsynced';
  
  // Timeouts
  static const int connectTimeoutSeconds = 30;
  static const int receiveTimeoutSeconds = 30;
  
  // Retry configuration
  static const int maxRetries = 3;
  static const int retryDelaySeconds = 2;
}

enum Environment {
  development,
  staging,
  production,
} 