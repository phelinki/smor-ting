import 'package:dio/dio.dart';
import '../../services/auth_service.dart';

class AuthInterceptor extends Interceptor {
  final AuthService _authService;
  
  // Track ongoing refresh to prevent multiple simultaneous refreshes
  static bool _isRefreshing = false;
  static bool _isInitializing = false; // Add this
  
  AuthInterceptor(this._authService);

  @override
  Future<void> onRequest(
    RequestOptions options, 
    RequestInterceptorHandler handler,
  ) async {
    // CRITICAL: Skip auth for these endpoints to prevent infinite loops
    if (_isAuthEndpoint(options.path) || _isInitializing) {
      return handler.next(options);
    }
    
    try {
      // Only try to get token if we're not refreshing
      if (!_isRefreshing) {
        final token = await _authService.getValidToken();
        if (token.isNotEmpty) {
          options.headers['Authorization'] = 'Bearer $token';
        }
      }
    } catch (e) {
      // No valid token available - continue without auth header
      print('No valid token available, proceeding without auth: $e');
    }
    
    handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err, 
    ErrorInterceptorHandler handler,
  ) async {
    // Only handle 401 errors and prevent refresh loops
    if (err.response?.statusCode == 401 && !_isRefreshing && !_isAuthEndpoint(err.requestOptions.path)) {
      _isRefreshing = true;
      
      try {
        // Attempt to refresh token
        final newToken = await _authService.refreshToken();
        
        // Clone the failed request and retry with new token
        final requestOptions = err.requestOptions;
        requestOptions.headers['Authorization'] = 'Bearer $newToken';
        
        final dio = Dio();
        final response = await dio.fetch(requestOptions);
        return handler.resolve(response);
        
      } catch (refreshError) {
        // Refresh failed, clear tokens and pass through original error
        print('Token refresh failed: $refreshError');
        // Let the original 401 error propagate
      } finally {
        _isRefreshing = false;
      }
    }
    
    return handler.next(err);
  }

  /// Check if the endpoint is an auth endpoint that shouldn't have auth headers
  bool _isAuthEndpoint(String path) {
    final authPaths = [
      '/auth/login',
      '/auth/register', 
      '/auth/refresh-token',
      '/auth/logout',
      '/auth/forgot-password',
      '/auth/reset-password',
    ];
    
    return authPaths.any((authPath) => path.contains(authPath));
  }

  // Add method to set initialization state
  static void setInitializing(bool value) {
    _isInitializing = value;
  }

  /// Reset refresh state for testing
  static void resetRefreshState() {
    _isRefreshing = false;
  }
}
