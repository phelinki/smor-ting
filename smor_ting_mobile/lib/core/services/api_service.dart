import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/user.dart';
import '../models/enhanced_auth_models.dart';
import '../constants/api_config.dart';
import '../models/kyc.dart';
import '../exceptions/auth_exceptions.dart';
import 'device_fingerprint_service.dart';
import '../../services/auth_service.dart';
import '../../services/dio_interceptor.dart';
import 'session_manager.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiService {
  late final Dio _dio;
  final bool _loggingEnabled;
  late final AuthService _authService;

  bool get loggingEnabled => _loggingEnabled;

  ApiService({String? baseUrl, bool? enableLogging, SessionManager? sessionManager}) : _loggingEnabled = enableLogging ?? !kReleaseMode {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl ?? ApiConfig.apiBaseUrl,
      connectTimeout: Duration(seconds: ApiConfig.connectTimeoutSeconds),
      receiveTimeout: Duration(seconds: ApiConfig.receiveTimeoutSeconds),
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'SmorTing-Mobile/${ApiConfig.environmentName}',
      },
    ));

    // Initialize AuthService for token refresh
    final secureStorage = const FlutterSecureStorage();
    _authService = AuthService(
      apiService: this,
      secureStorage: secureStorage,
    );

    // Add interceptors for logging and error handling
    if (_loggingEnabled) {
      _dio.interceptors.add(LogInterceptor(
        requestBody: true,
        responseBody: true,
        logPrint: (obj) => print(obj),
      ));
    }

    // Add custom auth interceptor to handle 401 errors and prevent infinite loops
    _dio.interceptors.add(AuthInterceptor(_authService));

    // Add general error interceptor
    _dio.interceptors.add(InterceptorsWrapper(
      onError: (error, handler) {
        print('API Error: ${error.message}');
        handler.next(error);
      },
    ));
  }

  // Set authorization token
  void setAuthToken(String token) {
    _dio.options.headers['Authorization'] = 'Bearer $token';
    // Also update cached token
    _authService.setCachedAccessToken(token);
  }

  // Clear authorization token
  void clearAuthToken() {
    _dio.options.headers.remove('Authorization');
    // Also clear cached token
    _authService.setCachedAccessToken(null);
  }

  // Get auth service for token management
  AuthService get authService => _authService;

  // Auth endpoints
  Future<AuthResponse> register(RegisterRequest request) async {
    try {
      final response = await _dio.post('/auth/register', data: request.toJson());
      return AuthResponse.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Enhanced authentication endpoints
  Future<Map<String, dynamic>> enhancedLogin(EnhancedLoginRequest request) async {
    try {
      final response = await _dio.post('/auth/login', data: request.toJson());
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Map<String, dynamic>> refreshToken(String refreshToken, String sessionId) async {
    try {
      final response = await _dio.post('/auth/refresh-token', data: {
        'refresh_token': refreshToken,
        'session_id': sessionId,
      });
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Map<String, dynamic>> getUserSessions() async {
    try {
      final response = await _dio.get('/auth/sessions');
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> revokeSession(String sessionId) async {
    try {
      await _dio.delete('/auth/sessions/$sessionId');
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> revokeAllSessions() async {
    try {
      await _dio.delete('/auth/sessions/all');
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<EnhancedAuthResult> biometricLogin(String email, String sessionId, DeviceFingerprint deviceInfo) async {
    try {
      final response = await _dio.post('/auth/biometric-login', data: {
        'email': email,
        'session_id': sessionId,
        'device_info': deviceInfo.toJson(),
      });
      return EnhancedAuthResult.fromResponse(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<AuthResponse> login(LoginRequest request) async {
    try {
      final response = await _dio.post('/auth/login', data: request.toJson());
      return AuthResponse.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }



  // Forgot password
  Future<void> requestPasswordReset(String email) async {
    try {
      await _dio.post('/auth/request-password-reset', data: {'email': email});
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> resetPassword(String email, String newPassword) async {
    try {
      await _dio.post('/auth/reset-password', data: {
        'email': email,
        'new_password': newPassword,
      });
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }
  // User endpoints
  Future<User> getUserProfile() async {
    try {
      final response = await _dio.get('/users/profile');
      return User.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Wallet endpoints
  Future<Map<String, dynamic>> getWalletBalances() async {
    try {
      final response = await _dio.get('/wallet/balances');
      return Map<String, dynamic>.from(response.data as Map);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // KYC endpoints (SmileID-backed via backend)
  Future<KycResponse> submitKyc(KycRequest request) async {
    try {
      final response = await _dio.post('/kyc/submit', data: request.toJson());
      return KycResponse.fromJson(response.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<User> updateUserProfile(Map<String, dynamic> data) async {
    try {
      final response = await _dio.put('/users/profile', data: data);
      return User.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Service endpoints
  Future<List<dynamic>> getServiceCategories() async {
    try {
      final response = await _dio.get('/services/categories');
      return response.data as List<dynamic>;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<List<dynamic>> getServices({String? categoryId}) async {
    try {
      final queryParams = categoryId != null ? {'category_id': categoryId} : <String, dynamic>{};
      final response = await _dio.get('/services', queryParameters: queryParams);
      return response.data as List<dynamic>;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<dynamic> getService(String serviceId) async {
    try {
      final response = await _dio.get('/services/$serviceId');
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Booking endpoints
  Future<dynamic> createBooking(Map<String, dynamic> bookingData) async {
    try {
      final response = await _dio.post('/bookings', data: bookingData);
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<List<dynamic>> getUserBookings() async {
    try {
      final response = await _dio.get('/bookings');
      return response.data as List<dynamic>;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<dynamic> getBooking(String bookingId) async {
    try {
      final response = await _dio.get('/bookings/$bookingId');
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<dynamic> updateBookingStatus(String bookingId, String status) async {
    try {
      final response = await _dio.patch('/bookings/$bookingId/status', 
        data: {'status': status});
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Provider endpoints
  Future<List<dynamic>> getServiceProviders({String? categoryId}) async {
    try {
      final queryParams = categoryId != null ? {'category_id': categoryId} : <String, dynamic>{};
      final response = await _dio.get('/providers', queryParameters: queryParams);
      return response.data as List<dynamic>;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<dynamic> getServiceProvider(String providerId) async {
    try {
      final response = await _dio.get('/providers/$providerId');
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Review endpoints
  Future<List<dynamic>> getServiceReviews(String serviceId) async {
    try {
      final response = await _dio.get('/services/$serviceId/reviews');
      return response.data as List<dynamic>;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<dynamic> createReview(Map<String, dynamic> reviewData) async {
    try {
      final response = await _dio.post('/reviews', data: reviewData);
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Never _handleError(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        throw const AuthException('Connection timeout. Please check your internet connection.');
      case DioExceptionType.badResponse:
        final statusCode = error.response?.statusCode;
        final responseData = error.response?.data;
        final message = responseData?['error'] ?? 'An error occurred';
        
        if (statusCode == 409) {
          print('🔴 API Service: Got 409 response with message: $message');
          print('🔴 API Service: Response data: $responseData');
          
          // Handle specific 409 conflicts
          if (message.toLowerCase().contains('already exists') || 
              message.toLowerCase().contains('user already exists')) {
            // Extract email from request if available
            final email = _extractEmailFromError(error) ?? 'unknown';
            print('🔴 API Service: Throwing EmailAlreadyExistsException with email: $email');
            throw EmailAlreadyExistsException(email);
          }
        } else if (statusCode == 401) {
          if (message.toLowerCase().contains('invalid credentials') ||
              message.toLowerCase().contains('invalid email or password')) {
            throw const InvalidCredentialsException();
          }
          throw const AuthException('Unauthorized. Please login again.');
        } else if (statusCode == 404) {
          throw const AuthException('Resource not found.');
        } else if (statusCode == 500) {
          throw const AuthException('Server error. Please try again later.');
        }
        
        throw AuthException(message);
      case DioExceptionType.cancel:
        throw const AuthException('Request cancelled.');
      case DioExceptionType.connectionError:
        throw const AuthException('No internet connection. Please check your network.');
      default:
        throw const AuthException('An unexpected error occurred.');
    }
  }

  String? _extractEmailFromError(DioException error) {
    try {
      // Try to extract email from the request data
      final requestData = error.requestOptions.data;
      if (requestData is Map<String, dynamic>) {
        return requestData['email'] as String?;
      }
    } catch (e) {
      // Ignore extraction errors
    }
    return null;
  }

  // Generic HTTP methods
  Future<Response<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
    ProgressCallback? onReceiveProgress,
  }) async {
    try {
      return await _dio.get<T>(
        path,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
        onReceiveProgress: onReceiveProgress,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Response<T>> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
    ProgressCallback? onSendProgress,
    ProgressCallback? onReceiveProgress,
  }) async {
    try {
      return await _dio.post<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
        onSendProgress: onSendProgress,
        onReceiveProgress: onReceiveProgress,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Response<T>> put<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
    ProgressCallback? onSendProgress,
    ProgressCallback? onReceiveProgress,
  }) async {
    try {
      return await _dio.put<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
        onSendProgress: onSendProgress,
        onReceiveProgress: onReceiveProgress,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Response<T>> delete<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.delete<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }
}

// Provider for ApiService
final apiServiceProvider = Provider<ApiService>((ref) {
  return ApiService();
});