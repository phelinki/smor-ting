import 'package:dio/dio.dart';
import 'auth_service.dart';

class AuthInterceptor extends Interceptor {
  final AuthService _authService;
  static bool _isRefreshing = false;
  static final List<String> _authEndpoints = [
    '/auth/login',
    '/auth/register', 
    '/auth/refresh-token',
    '/auth/logout',
    '/auth/forgot-password',
    '/auth/reset-password',
  ];
  
  AuthInterceptor(this._authService);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    // CRITICAL: Never add auth headers to auth endpoints
    if (_isAuthEndpoint(options.path)) {
      print('Skipping auth for endpoint: ${options.path}');
      return handler.next(options);
    }
    
    // CRITICAL: Don't add auth headers if refresh is in progress
    if (_isRefreshing) {
      print('Refresh in progress, skipping auth header');
      return handler.next(options);
    }
    
    try {
      final token = await _authService.getValidToken();
      options.headers['Authorization'] = 'Bearer $token';
      print('Added auth header to: ${options.path}');
    } catch (e) {
      print('Failed to get valid token for ${options.path}: $e');
      // Continue without auth header - let the request fail naturally
    }
    
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    print('Request error: ${err.response?.statusCode} for ${err.requestOptions.path}');
    
    // Only handle 401 errors for non-auth endpoints
    if (err.response?.statusCode == 401 && 
        !_isAuthEndpoint(err.requestOptions.path) && 
        !_isRefreshing) {
      
      print('Attempting token refresh for 401 error');
      _isRefreshing = true;
      
      try {
        final newToken = await _authService.refreshToken();
        print('Token refresh successful, retrying request');
        
        // Clone and retry the original request
        final options = err.requestOptions;
        options.headers['Authorization'] = 'Bearer $newToken';
        
        // Use the same Dio instance but bypass interceptors for retry
        final response = await Dio().fetch(options);
        return handler.resolve(response);
        
      } catch (refreshError) {
        print('Token refresh failed: $refreshError');
        // Let the original error propagate
      } finally {
        _isRefreshing = false;
      }
    }
    
    handler.next(err);
  }

  bool _isAuthEndpoint(String path) {
    return _authEndpoints.any((authPath) => path.contains(authPath));
  }
}
