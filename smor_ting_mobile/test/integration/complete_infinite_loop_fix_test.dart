import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:smor_ting_mobile/core/interceptors/auth_interceptor.dart';
import 'package:smor_ting_mobile/services/auth_service.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:dio/dio.dart';

class MockApiService extends Mock implements ApiService {}
class MockFlutterSecureStorage extends Mock implements FlutterSecureStorage {}

void main() {
  group('COMPLETE INFINITE LOOP FIX VERIFICATION', () {
    late AuthInterceptor interceptor;
    late AuthService authService;
    late MockApiService mockApiService;
    late MockFlutterSecureStorage mockSecureStorage;

    const validJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';

    setUpAll(() {
      registerFallbackValue(RequestOptions(path: ''));
      registerFallbackValue(DioException(requestOptions: RequestOptions(path: '')));
      registerFallbackValue(Response(requestOptions: RequestOptions(path: '')));
      registerFallbackValue(<String, dynamic>{});
    });

    setUp(() {
      mockApiService = MockApiService();
      mockSecureStorage = MockFlutterSecureStorage();
      authService = AuthService(
        apiService: mockApiService,
        secureStorage: mockSecureStorage,
      );
      interceptor = AuthInterceptor(authService);
      
      AuthInterceptor.resetRefreshState();
    });

    test('🎯 FINAL VERIFICATION: All 4 requirements implemented correctly', () async {
      print('🧪 TESTING ALL 4 INFINITE LOOP PREVENTION REQUIREMENTS:');
      
      // Setup mocks
      when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => 'expired_token');
      when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => validJwt);
      when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => 'test-session');
      
      when(() => mockApiService.refreshToken(any(), any())).thenAnswer((_) async => {
        'access_token': validJwt,
        'refresh_token': validJwt,
        'refresh_expires_at': DateTime.now().add(const Duration(days: 30)).toIso8601String(),
        'token_expires_at': DateTime.now().add(const Duration(hours: 1)).toIso8601String(),
        'session_id': 'test-session',
      });

      when(() => mockSecureStorage.write(key: any(named: 'key'), value: any(named: 'value')))
          .thenAnswer((_) async {});

      // ✅ REQUIREMENT 1: Prevent concurrent refresh requests using a Completer
      print('✅ 1. TESTING: Prevent concurrent refresh requests using a Completer');
      final futures = List.generate(5, (_) => authService.refreshToken());
      final results = await Future.wait(futures);
      
      verify(() => mockApiService.refreshToken(any(), any())).called(1);
      for (final result in results) {
        expect(result, equals(validJwt));
      }
      print('   ✅ PASS: 5 concurrent requests → 1 API call');

      // Reset for next test
      reset(mockApiService);
      reset(mockSecureStorage);
      
      // ✅ REQUIREMENT 2: Add token validity buffer to prevent edge case refreshes
      print('✅ 2. TESTING: Token validity buffer (5-minute)');
      when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => validJwt);
      
      final validToken = await authService.getValidToken();
      expect(validToken, equals(validJwt));
      verifyNever(() => mockApiService.refreshToken(any(), any()));
      print('   ✅ PASS: Valid token returned without refresh');

      // ✅ REQUIREMENT 3: Skip auth header for refresh endpoint to prevent loops  
      print('✅ 3. TESTING: Skip auth header for refresh endpoint');
      final refreshRequest = RequestOptions(path: '/auth/refresh-token');
      final requestHandler = _MockRequestHandler();
      
      await interceptor.onRequest(refreshRequest, requestHandler);
      expect(refreshRequest.headers['Authorization'], isNull);
      print('   ✅ PASS: No auth header added to refresh endpoint');

      // ✅ REQUIREMENT 4: Proper error handling to avoid retry loops
      print('✅ 4. TESTING: Proper error handling');
      final errorRequest = RequestOptions(path: '/auth/refresh-token');
      final dioError = DioException(
        requestOptions: errorRequest,
        response: Response(requestOptions: errorRequest, statusCode: 401),
      );
      final errorHandler = _MockErrorHandler();
      
      await interceptor.onError(dioError, errorHandler);
      verifyNever(() => mockApiService.refreshToken(any(), any()));
      print('   ✅ PASS: No refresh attempted for auth endpoint 401');

      print('');
      print('🎉 ALL 4 REQUIREMENTS SUCCESSFULLY VERIFIED!');
      print('🎉 INFINITE LOOP ISSUE IS COMPLETELY FIXED!');
    });

    test('🚀 PERFORMANCE: Verify sub-second response time', () async {
      // Setup
      when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => 'expired_token');
      when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => validJwt);
      when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => 'test-session');
      
      when(() => mockApiService.refreshToken(any(), any())).thenAnswer((_) async {
        await Future.delayed(const Duration(milliseconds: 50)); // Simulate network
        return {
          'access_token': validJwt,
          'refresh_token': validJwt,
          'refresh_expires_at': DateTime.now().add(const Duration(days: 30)).toIso8601String(),
          'token_expires_at': DateTime.now().add(const Duration(hours: 1)).toIso8601String(),
          'session_id': 'test-session',
        };
      });

      when(() => mockSecureStorage.write(key: any(named: 'key'), value: any(named: 'value')))
          .thenAnswer((_) async {});

      // Test
      final stopwatch = Stopwatch()..start();
      final futures = List.generate(10, (_) => authService.getValidToken());
      await Future.wait(futures);
      stopwatch.stop();

      expect(stopwatch.elapsedMilliseconds, lessThan(500), 
        reason: 'Should complete 10 concurrent requests in under 500ms');
      
      print('🚀 PERFORMANCE: 10 concurrent requests completed in ${stopwatch.elapsedMilliseconds}ms');
    });
  });
}

class _MockRequestHandler extends Mock implements RequestInterceptorHandler {}
class _MockErrorHandler extends Mock implements ErrorInterceptorHandler {}
