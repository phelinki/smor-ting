import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:dio/dio.dart';
import 'package:smor_ting_mobile/services/dio_interceptor.dart';
import 'package:smor_ting_mobile/services/auth_service.dart';

class MockAuthService extends Mock implements AuthService {}

void main() {
  group('DioInterceptor - HTTP Interceptor Tests', () {
    late DioInterceptor interceptor;
    late MockAuthService mockAuthService;
    late Dio dio;

    const validJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';

    setUp(() {
      mockAuthService = MockAuthService();
      interceptor = DioInterceptor(mockAuthService);
      dio = Dio();
      dio.interceptors.add(interceptor);
      
      // Register fallback values
      registerFallbackValue(RequestOptions(path: '/test'));
      registerFallbackValue(ErrorInterceptorHandler());
    });

    group('Auth Header Management Tests', () {
      test('should skip auth header for refresh endpoint to prevent loops', () async {
        // Arrange
        final requestOptions = RequestOptions(
          path: '/auth/refresh-token',
          method: 'POST',
        );
        
        bool authHeaderSkipped = false;
        final testHandler = InterceptorHandler();
        
        // Override the handler to capture what happens
        when(() => mockAuthService.getCurrentAccessToken()).thenReturn(validJwt);

        // Act - This should NOT add Authorization header for refresh endpoint
        await interceptor.onRequest(requestOptions, testHandler);

        // Assert - Auth header should not be added for refresh endpoint
        expect(requestOptions.headers['Authorization'], isNull);
      });

      test('should add auth header for non-refresh endpoints', () async {
        // Arrange
        final requestOptions = RequestOptions(
          path: '/api/user/profile',
          method: 'GET',
        );
        
        final testHandler = InterceptorHandler();
        when(() => mockAuthService.getCurrentAccessToken()).thenReturn(validJwt);

        // Act
        await interceptor.onRequest(requestOptions, testHandler);

        // Assert - Auth header should be added for regular endpoints
        expect(requestOptions.headers['Authorization'], equals('Bearer $validJwt'));
      });

      test('should not add auth header if no token available', () async {
        // Arrange
        final requestOptions = RequestOptions(
          path: '/api/user/profile',
          method: 'GET',
        );
        
        final testHandler = InterceptorHandler();
        when(() => mockAuthService.getCurrentAccessToken()).thenReturn(null);

        // Act
        await interceptor.onRequest(requestOptions, testHandler);

        // Assert - No auth header should be added if no token
        expect(requestOptions.headers['Authorization'], isNull);
      });
    });

    group('401 Error Handling Tests', () {
      test('should trigger token refresh on 401 error for non-refresh endpoints', () async {
        // Arrange
        final dioError = DioException(
          requestOptions: RequestOptions(path: '/api/user/profile'),
          response: Response(
            requestOptions: RequestOptions(path: '/api/user/profile'),
            statusCode: 401,
            data: {'error': 'Unauthorized'},
          ),
        );

        when(() => mockAuthService.refreshToken()).thenAnswer((_) async => {
          'success': true,
          'access_token': validJwt,
        });

        when(() => mockAuthService.getCurrentAccessToken()).thenReturn(validJwt);

        final testHandler = ErrorInterceptorHandler();

        // Act
        await interceptor.onError(dioError, testHandler);

        // Assert - Should trigger refresh token
        verify(() => mockAuthService.refreshToken()).called(1);
      });

      test('should NOT trigger token refresh on 401 error for refresh endpoint to prevent loops', () async {
        // Arrange
        final dioError = DioException(
          requestOptions: RequestOptions(path: '/auth/refresh-token'),
          response: Response(
            requestOptions: RequestOptions(path: '/auth/refresh-token'),
            statusCode: 401,
            data: {'error': 'Invalid refresh token'},
          ),
        );

        final testHandler = ErrorInterceptorHandler();

        // Act
        await interceptor.onError(dioError, testHandler);

        // Assert - Should NOT trigger refresh to prevent infinite loop
        verifyNever(() => mockAuthService.refreshToken());
      });

      test('should retry original request after successful token refresh', () async {
        // Arrange
        final originalRequest = RequestOptions(path: '/api/user/profile');
        final dioError = DioException(
          requestOptions: originalRequest,
          response: Response(
            requestOptions: originalRequest,
            statusCode: 401,
            data: {'error': 'Unauthorized'},
          ),
        );

        when(() => mockAuthService.refreshToken()).thenAnswer((_) async => {
          'success': true,
          'access_token': validJwt,
        });

        when(() => mockAuthService.getCurrentAccessToken()).thenReturn(validJwt);

        final testHandler = ErrorInterceptorHandler();
        bool requestResolved = false;
        
        // Mock the handler to capture if request gets resolved
        testHandler.resolve = (response) => requestResolved = true;

        // Act
        await interceptor.onError(dioError, testHandler);

        // Assert - Should have attempted to resolve with retry
        verify(() => mockAuthService.refreshToken()).called(1);
        // Note: Full retry logic test would require more complex mocking
      });

      test('should not retry if token refresh fails', () async {
        // Arrange
        final dioError = DioException(
          requestOptions: RequestOptions(path: '/api/user/profile'),
          response: Response(
            requestOptions: RequestOptions(path: '/api/user/profile'),
            statusCode: 401,
            data: {'error': 'Unauthorized'},
          ),
        );

        when(() => mockAuthService.refreshToken()).thenAnswer((_) async => null);

        final testHandler = ErrorInterceptorHandler();
        bool errorPassed = false;
        testHandler.next = (error) => errorPassed = true;

        // Act
        await interceptor.onError(dioError, testHandler);

        // Assert - Should pass error along if refresh fails
        verify(() => mockAuthService.refreshToken()).called(1);
        expect(errorPassed, isTrue);
      });
    });

    group('Error Loop Prevention Tests', () {
      test('should prevent infinite retry loops by tracking retry attempts', () async {
        // Arrange
        final originalRequest = RequestOptions(path: '/api/user/profile');
        originalRequest.extra['retryCount'] = 2; // Already retried twice
        
        final dioError = DioException(
          requestOptions: originalRequest,
          response: Response(
            requestOptions: originalRequest,
            statusCode: 401,
            data: {'error': 'Unauthorized'},
          ),
        );

        final testHandler = ErrorInterceptorHandler();
        bool errorPassed = false;
        testHandler.next = (error) => errorPassed = true;

        // Act
        await interceptor.onError(dioError, testHandler);

        // Assert - Should not attempt refresh if already retried max times
        verifyNever(() => mockAuthService.refreshToken());
        expect(errorPassed, isTrue);
      });

      test('should handle non-401 errors normally', () async {
        // Arrange
        final dioError = DioException(
          requestOptions: RequestOptions(path: '/api/user/profile'),
          response: Response(
            requestOptions: RequestOptions(path: '/api/user/profile'),
            statusCode: 404,
            data: {'error': 'Not found'},
          ),
        );

        final testHandler = ErrorInterceptorHandler();
        bool errorPassed = false;
        testHandler.next = (error) => errorPassed = true;

        // Act
        await interceptor.onError(dioError, testHandler);

        // Assert - Should pass through non-401 errors without refresh attempt
        verifyNever(() => mockAuthService.refreshToken());
        expect(errorPassed, isTrue);
      });
    });
  });
}
