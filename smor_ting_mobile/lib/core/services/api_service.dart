import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/user.dart';

class ApiService {
  late final Dio _dio;
  
  ApiService({String? baseUrl}) {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl ?? 'http://localhost:8080/api/v1',
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {
        'Content-Type': 'application/json',
      },
    ));

    // Add interceptors for logging and error handling
    _dio.interceptors.add(LogInterceptor(
      requestBody: true,
      responseBody: true,
      logPrint: (obj) => print(obj),
    ));

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
  }

  // Clear authorization token
  void clearAuthToken() {
    _dio.options.headers.remove('Authorization');
  }

  // Auth endpoints
  Future<AuthResponse> register(RegisterRequest request) async {
    try {
      final response = await _dio.post('/auth/register', data: request.toJson());
      return AuthResponse.fromJson(response.data);
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

  Future<AuthResponse> verifyOTP(VerifyOTPRequest request) async {
    try {
      final response = await _dio.post('/auth/verify-otp', data: request.toJson());
      return AuthResponse.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> resendOTP(String email) async {
    try {
      await _dio.post('/auth/resend-otp', data: {'email': email});
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

  String _handleError(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return 'Connection timeout. Please check your internet connection.';
      case DioExceptionType.badResponse:
        final statusCode = error.response?.statusCode;
        final message = error.response?.data?['error'] ?? 'An error occurred';
        if (statusCode == 401) {
          return 'Unauthorized. Please login again.';
        } else if (statusCode == 404) {
          return 'Resource not found.';
        } else if (statusCode == 500) {
          return 'Server error. Please try again later.';
        }
        return message;
      case DioExceptionType.cancel:
        return 'Request cancelled.';
      case DioExceptionType.connectionError:
        return 'No internet connection. Please check your network.';
      default:
        return 'An unexpected error occurred.';
    }
  }
}

// Provider for ApiService
final apiServiceProvider = Provider<ApiService>((ref) {
  return ApiService();
});