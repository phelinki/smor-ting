import 'package:dio/dio.dart';
import 'auth_service.dart';

class AuthInterceptor extends Interceptor {
  final AuthService _authService;
  
  AuthInterceptor(this._authService);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    // Skip auth for refresh token endpoint to prevent loop
    if (options.path.contains('/auth/refresh-token')) {
      return handler.next(options);
    }
    
    try {
      final token = await _authService.getValidToken();
      options.headers['Authorization'] = 'Bearer $token';
    } catch (e) {
      // Handle auth error (redirect to login, etc.)
      // For now, continue without token
    }
    
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      // Skip refresh attempt for auth endpoints to prevent loops
      if (err.requestOptions.path.contains('/auth/refresh-token') ||
          err.requestOptions.path.contains('/auth/login') ||
          err.requestOptions.path.contains('/auth/register')) {
        return handler.next(err);
      }
      
      try {
        // Try to refresh token
        final newToken = await _authService.refreshToken();
        
        // Retry original request with new token
        final options = err.requestOptions;
        options.headers['Authorization'] = 'Bearer $newToken';
        
        final dio = Dio();
        final response = await dio.fetch(options);
        return handler.resolve(response);
      } catch (refreshError) {
        // Refresh failed, redirect to login
        // For now, pass the error along
        return handler.next(err);
      }
    }
    
    handler.next(err);
  }
}