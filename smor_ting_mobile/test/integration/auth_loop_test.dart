import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:smor_ting_mobile/services/auth_service.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class MockApiService extends Mock implements ApiService {}
class MockFlutterSecureStorage extends Mock implements FlutterSecureStorage {}

void main() {
  group('Auth Service Integration - Infinite Loop Prevention', () {
    late AuthService authService;
    late MockApiService mockApiService;
    late MockFlutterSecureStorage mockSecureStorage;

    // Valid JWT token for testing (expires in year 2030)
    const validJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE5MDkzMzkyMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';

    setUp(() {
      mockApiService = MockApiService();
      mockSecureStorage = MockFlutterSecureStorage();
      authService = AuthService(
        apiService: mockApiService,
        secureStorage: mockSecureStorage,
      );
      
      // Register fallback values
      registerFallbackValue(<String, dynamic>{});
    });

    test('CRITICAL: Should prevent infinite loops during concurrent refresh requests', () async {
      // Arrange - Simulate expired token that needs refresh
      when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => 'expired_token');
      when(() => mockSecureStorage.read(key: 'refresh_token')).thenAnswer((_) async => validJwt);
      when(() => mockSecureStorage.read(key: 'session_id')).thenAnswer((_) async => 'test-session');
      
      // Mock slow API response to allow concurrent calls
      when(() => mockApiService.refreshToken(any(), any())).thenAnswer((_) async {
        await Future.delayed(const Duration(milliseconds: 200));
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

      print('🧪 Testing concurrent refresh requests...');
      
      // Act - Make 10 concurrent calls to getValidToken (simulating multiple API calls)
      final stopwatch = Stopwatch()..start();
      final futures = List.generate(10, (i) async {
        print('🔄 Starting request $i');
        try {
          final result = await authService.getValidToken();
          print('✅ Request $i completed: ${result.substring(0, 20)}...');
          return result;
        } catch (e) {
          print('❌ Request $i failed: $e');
          rethrow;
        }
      });
      
      final results = await Future.wait(futures);
      stopwatch.stop();
      
      print('⏱️  Total time: ${stopwatch.elapsedMilliseconds}ms');
      
      // Assert - Only ONE API call should have been made despite 10 concurrent requests
      verify(() => mockApiService.refreshToken(any(), any())).called(1);
      print('✅ PASS: Only 1 API call made for 10 concurrent requests');
      
      // All results should be the same valid token
      for (int i = 0; i < results.length; i++) {
        expect(results[i], equals(validJwt), reason: 'Request $i should return the same token');
      }
      print('✅ PASS: All requests returned the same token');
      
      // Should complete reasonably quickly (under 1 second for 10 concurrent calls)
      expect(stopwatch.elapsedMilliseconds, lessThan(1000), 
        reason: 'Should complete quickly without retry loops');
      print('✅ PASS: Completed in ${stopwatch.elapsedMilliseconds}ms (under 1 second)');
      
      print('🎉 INFINITE LOOP PREVENTION TEST PASSED!');
    });

    test('Should handle token validation correctly with buffer', () async {
      // Arrange - Token that expires in 3 minutes (within 5-minute buffer)
      const expiringSoonJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE3MzE4MjQ2MDB9.invalid_signature_but_we_only_check_exp';
      
      when(() => mockSecureStorage.read(key: 'access_token')).thenAnswer((_) async => expiringSoonJwt);
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

      print('🧪 Testing token buffer validation...');
      
      // Act
      final result = await authService.getValidToken();
      
      // Assert - Should trigger refresh due to 5-minute buffer
      verify(() => mockApiService.refreshToken(any(), any())).called(1);
      expect(result, equals(validJwt));
      
      print('✅ PASS: Token buffer validation working correctly');
    });
  });
}
