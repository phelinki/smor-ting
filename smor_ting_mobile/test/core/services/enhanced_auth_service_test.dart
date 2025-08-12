import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:local_auth/local_auth.dart';
import 'package:smor_ting_mobile/core/services/enhanced_auth_service.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:smor_ting_mobile/core/services/session_manager.dart';
import 'package:smor_ting_mobile/core/services/device_fingerprint_service.dart';
import 'package:smor_ting_mobile/core/models/user.dart';

import 'enhanced_auth_service_test.mocks.dart';

@GenerateMocks([
  ApiService,
  SessionManager,
  DeviceFingerprintService,
  FlutterSecureStorage,
  LocalAuthentication,
])
void main() {
  group('EnhancedAuthService Tests', () {
    late EnhancedAuthService authService;
    late MockApiService mockApiService;
    late MockSessionManager mockSessionManager;
    late MockDeviceFingerprintService mockDeviceService;
    late MockFlutterSecureStorage mockSecureStorage;
    late MockLocalAuthentication mockLocalAuth;

    setUp(() {
      mockApiService = MockApiService();
      mockSessionManager = MockSessionManager();
      mockDeviceService = MockDeviceFingerprintService();
      mockSecureStorage = MockFlutterSecureStorage();
      mockLocalAuth = MockLocalAuthentication();

      authService = EnhancedAuthService(
        mockApiService,
        mockSessionManager,
        mockDeviceService,
        mockSecureStorage,
        mockLocalAuth,
      );
    });

    group('Enhanced Login', () {
      test('should perform successful login with remember me', () async {
        // Arrange
        final deviceFingerprint = DeviceFingerprint(
          deviceId: 'test_device_123',
          platform: 'Android',
          osVersion: 'Android 11',
          appVersion: '1.0.0',
          isJailbroken: false,
          attestationData: 'valid_attestation',
        );

        final expectedUser = User(
          id: 'user123',
          email: 'test@example.com',
          firstName: 'Test',
          lastName: 'User',
          role: UserRole.customer,
          isVerified: true,
        );

        final mockResponse = {
          'success': true,
          'user': expectedUser.toJson(),
          'access_token': 'access_token_123',
          'refresh_token': 'refresh_token_123',
          'session_id': 'session_123',
          'token_expires_at': DateTime.now().add(Duration(minutes: 30)).toIso8601String(),
          'refresh_expires_at': DateTime.now().add(Duration(days: 7)).toIso8601String(),
          'device_trusted': true,
          'requires_verification': false,
        };

        when(mockDeviceService.generateFingerprint())
            .thenAnswer((_) async => deviceFingerprint);

        when(mockApiService.enhancedLogin(any))
            .thenAnswer((_) async => mockResponse);

        when(mockSessionManager.storeSession(any))
            .thenAnswer((_) async {});

        when(mockLocalAuth.canCheckBiometrics)
            .thenAnswer((_) async => true);

        when(mockSecureStorage.write(key: anyNamed('key'), value: anyNamed('value')))
            .thenAnswer((_) async {});

        // Act
        final result = await authService.enhancedLogin(
          email: 'test@example.com',
          password: 'password123',
          rememberMe: true,
        );

        // Assert
        expect(result.success, true);
        expect(result.user?.email, 'test@example.com');
        expect(result.accessToken, 'access_token_123');
        expect(result.deviceTrusted, true);
        expect(result.isRestoredSession, false);

        // Verify session was stored
        verify(mockSessionManager.storeSession(any)).called(1);
        
        // Verify biometric preference was stored (since device is trusted and remember me is true)
        verify(mockSecureStorage.write(
          key: argThat(contains('biometric_preference_')),
          value: 'true',
        )).called(1);
      });

      test('should handle two-factor authentication requirement', () async {
        // Arrange
        final deviceFingerprint = DeviceFingerprint(
          deviceId: 'untrusted_device',
          platform: 'Android',
          osVersion: 'Android 11',
          appVersion: '1.0.0',
          isJailbroken: true, // Jailbroken device requires 2FA
          attestationData: 'compromised',
        );

        final mockResponse = {
          'success': true,
          'requires_two_factor': true,
          'device_trusted': false,
          'user': {
            'id': 'user123',
            'email': 'test@example.com',
            'first_name': 'Test',
            'last_name': 'User',
            'role': 'customer',
            'is_verified': true,
          },
        };

        when(mockDeviceService.generateFingerprint())
            .thenAnswer((_) async => deviceFingerprint);

        when(mockApiService.enhancedLogin(any))
            .thenAnswer((_) async => mockResponse);

        // Act
        final result = await authService.enhancedLogin(
          email: 'test@example.com',
          password: 'password123',
        );

        // Assert
        expect(result.success, true);
        expect(result.requiresTwoFactor, true);
        expect(result.deviceTrusted, false);
        expect(result.accessToken, isNull);

        // Verify session was not stored yet
        verifyNever(mockSessionManager.storeSession(any));
      });

      test('should handle CAPTCHA requirement', () async {
        // Arrange
        final deviceFingerprint = DeviceFingerprint(
          deviceId: 'test_device',
          platform: 'Android',
          osVersion: 'Android 11',
          appVersion: '1.0.0',
          isJailbroken: false,
          attestationData: 'valid',
        );

        final mockResponse = {
          'success': false,
          'requires_captcha': true,
          'remaining_attempts': 2,
          'message': 'CAPTCHA verification required',
        };

        when(mockDeviceService.generateFingerprint())
            .thenAnswer((_) async => deviceFingerprint);

        when(mockApiService.enhancedLogin(any))
            .thenAnswer((_) async => mockResponse);

        // Act
        final result = await authService.enhancedLogin(
          email: 'test@example.com',
          password: 'wrong_password',
        );

        // Assert
        expect(result.success, false);
        expect(result.requiresCaptcha, true);
        expect(result.remainingAttempts, 2);
      });

      test('should handle account lockout', () async {
        // Arrange
        final deviceFingerprint = DeviceFingerprint(
          deviceId: 'test_device',
          platform: 'Android',
          osVersion: 'Android 11',
          appVersion: '1.0.0',
          isJailbroken: false,
          attestationData: 'valid',
        );

        final mockResponse = {
          'success': false,
          'lockout_info': {
            'email_locked': true,
            'email_lockout_remaining': 900000, // 15 minutes
            'email_attempts': 5,
            'ip_locked': false,
            'ip_lockout_remaining': 0,
            'ip_attempts': 3,
          },
          'message': 'Account temporarily locked',
        };

        when(mockDeviceService.generateFingerprint())
            .thenAnswer((_) async => deviceFingerprint);

        when(mockApiService.enhancedLogin(any))
            .thenAnswer((_) async => mockResponse);

        // Act
        final result = await authService.enhancedLogin(
          email: 'test@example.com',
          password: 'wrong_password',
        );

        // Assert
        expect(result.success, false);
        expect(result.lockoutInfo?.emailLocked, true);
        expect(result.lockoutInfo?.emailAttempts, 5);
      });
    });

    group('Session Restoration', () {
      test('should restore valid session', () async {
        // Arrange
        final sessionData = SessionData(
          accessToken: 'valid_access_token',
          refreshToken: 'valid_refresh_token',
          sessionId: 'session_123',
          user: User(
            id: 'user123',
            email: 'test@example.com',
            firstName: 'Test',
            lastName: 'User',
            role: UserRole.customer,
            isVerified: true,
          ),
          tokenExpiresAt: DateTime.now().add(Duration(minutes: 15)), // Still valid
          refreshExpiresAt: DateTime.now().add(Duration(days: 6)),
          deviceTrusted: true,
          rememberMe: true,
        );

        when(mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        // Act
        final result = await authService.restoreSession();

        // Assert
        expect(result?.success, true);
        expect(result?.user?.email, 'test@example.com');
        expect(result?.isRestoredSession, true);
        expect(result?.deviceTrusted, true);

        // Verify API service token was set
        verify(mockApiService.setAuthToken('valid_access_token')).called(1);
      });

      test('should refresh token when access token is expired', () async {
        // Arrange
        final sessionData = SessionData(
          accessToken: 'expired_access_token',
          refreshToken: 'valid_refresh_token',
          sessionId: 'session_123',
          user: User(
            id: 'user123',
            email: 'test@example.com',
            firstName: 'Test',
            lastName: 'User',
            role: UserRole.customer,
            isVerified: true,
          ),
          tokenExpiresAt: DateTime.now().subtract(Duration(minutes: 5)), // Expired
          refreshExpiresAt: DateTime.now().add(Duration(days: 6)), // Still valid
          deviceTrusted: true,
          rememberMe: true,
        );

        final refreshResponse = {
          'success': true,
          'access_token': 'new_access_token',
          'refresh_token': 'new_refresh_token',
          'token_expires_at': DateTime.now().add(Duration(minutes: 30)).toIso8601String(),
          'refresh_expires_at': DateTime.now().add(Duration(days: 7)).toIso8601String(),
        };

        when(mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(mockApiService.refreshToken(any, any))
            .thenAnswer((_) async => refreshResponse);

        when(mockSessionManager.storeSession(any))
            .thenAnswer((_) async {});

        // Act
        final result = await authService.restoreSession();

        // Assert
        expect(result?.success, true);
        expect(result?.accessToken, 'new_access_token');
        expect(result?.isRestoredSession, true);

        // Verify refresh was called
        verify(mockApiService.refreshToken('valid_refresh_token', 'session_123')).called(1);
        verify(mockSessionManager.storeSession(any)).called(1);
      });

      test('should try biometric unlock when remember me is enabled', () async {
        // Arrange
        final sessionData = SessionData(
          accessToken: 'expired_access_token',
          refreshToken: 'expired_refresh_token',
          sessionId: 'session_123',
          user: User(
            id: 'user123',
            email: 'test@example.com',
            firstName: 'Test',
            lastName: 'User',
            role: UserRole.customer,
            isVerified: true,
          ),
          tokenExpiresAt: DateTime.now().subtract(Duration(minutes: 5)), // Expired
          refreshExpiresAt: DateTime.now().subtract(Duration(hours: 1)), // Also expired
          deviceTrusted: true,
          rememberMe: true, // Remember me enabled
        );

        final biometricResponse = {
          'success': true,
          'user': sessionData.user.toJson(),
          'access_token': 'biometric_access_token',
          'refresh_token': 'biometric_refresh_token',
          'session_id': 'new_session_123',
          'token_expires_at': DateTime.now().add(Duration(minutes: 30)).toIso8601String(),
          'refresh_expires_at': DateTime.now().add(Duration(days: 7)).toIso8601String(),
          'device_trusted': true,
        };

        when(mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(mockApiService.refreshToken(any, any))
            .thenThrow(Exception('Refresh token expired'));

        when(mockSecureStorage.read(key: 'biometric_enabled_test@example.com'))
            .thenAnswer((_) async => 'true');

        when(mockLocalAuth.canCheckBiometrics)
            .thenAnswer((_) async => true);

        when(mockLocalAuth.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint]);

        when(mockLocalAuth.authenticate(
          localizedReason: anyNamed('localizedReason'),
          options: anyNamed('options'),
        )).thenAnswer((_) async => true);

        when(mockDeviceService.generateFingerprint())
            .thenAnswer((_) async => DeviceFingerprint(
              deviceId: 'device123',
              platform: 'Android',
              osVersion: 'Android 11',
              appVersion: '1.0.0',
              isJailbroken: false,
              attestationData: 'valid',
            ));

        when(mockApiService.biometricLogin(any, any, any))
            .thenAnswer((_) async => EnhancedAuthResult.fromResponse(biometricResponse));

        when(mockSessionManager.storeSession(any))
            .thenAnswer((_) async {});

        // Act
        final result = await authService.restoreSession();

        // Assert
        expect(result?.success, true);
        expect(result?.accessToken, 'biometric_access_token');
        expect(result?.isRestoredSession, true);

        // Verify biometric authentication was attempted
        verify(mockLocalAuth.authenticate(
          localizedReason: anyNamed('localizedReason'),
          options: anyNamed('options'),
        )).called(1);
      });

      test('should clear session when all restoration methods fail', () async {
        // Arrange
        final sessionData = SessionData(
          accessToken: 'expired_access_token',
          refreshToken: 'expired_refresh_token',
          sessionId: 'session_123',
          user: User(
            id: 'user123',
            email: 'test@example.com',
            firstName: 'Test',
            lastName: 'User',
            role: UserRole.customer,
            isVerified: true,
          ),
          tokenExpiresAt: DateTime.now().subtract(Duration(minutes: 5)), // Expired
          refreshExpiresAt: DateTime.now().subtract(Duration(hours: 1)), // Also expired
          deviceTrusted: false, // Not trusted
          rememberMe: false, // Remember me not enabled
        );

        when(mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(mockApiService.refreshToken(any, any))
            .thenThrow(Exception('Refresh token expired'));

        when(mockSessionManager.clearSession())
            .thenAnswer((_) async {});

        // Act
        final result = await authService.restoreSession();

        // Assert
        expect(result, isNull);

        // Verify session was cleared
        verify(mockSessionManager.clearSession()).called(1);
      });
    });

    group('Session Management', () {
      test('should get user sessions', () async {
        // Arrange
        final mockSessions = [
          {
            'session_id': 'session_1',
            'device_info': {
              'device_id': 'device_1',
              'platform': 'Android',
              'os_version': 'Android 11',
              'app_version': '1.0.0',
              'is_jailbroken': false,
              'attestation_data': 'valid',
            },
            'ip_address': '192.168.1.1',
            'user_agent': 'App/1.0.0',
            'is_remembered': true,
            'last_activity': DateTime.now().toIso8601String(),
            'created_at': DateTime.now().subtract(Duration(days: 1)).toIso8601String(),
            'expires_at': DateTime.now().add(Duration(days: 29)).toIso8601String(),
          },
        ];

        when(mockApiService.getUserSessions())
            .thenAnswer((_) async => {'sessions': mockSessions});

        // Act
        final sessions = await authService.getUserSessions();

        // Assert
        expect(sessions.length, 1);
        expect(sessions[0].sessionId, 'session_1');
        expect(sessions[0].deviceInfo.platform, 'Android');
        expect(sessions[0].isRemembered, true);
      });

      test('should revoke specific session', () async {
        // Arrange
        const sessionId = 'session_to_revoke';

        when(mockApiService.revokeSession(sessionId))
            .thenAnswer((_) async {});

        // Act
        await authService.revokeSession(sessionId);

        // Assert
        verify(mockApiService.revokeSession(sessionId)).called(1);
      });

      test('should sign out all devices', () async {
        // Arrange
        when(mockApiService.revokeAllSessions())
            .thenAnswer((_) async {});

        when(mockSessionManager.clearSession())
            .thenAnswer((_) async {});

        // Act
        await authService.signOutAllDevices();

        // Assert
        verify(mockApiService.revokeAllSessions()).called(1);
        verify(mockSessionManager.clearSession()).called(1);
        verify(mockApiService.clearAuthToken()).called(1);
      });
    });

    group('Biometric Authentication', () {
      test('should enable biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';

        when(mockLocalAuth.canCheckBiometrics)
            .thenAnswer((_) async => true);

        when(mockLocalAuth.authenticate(
          localizedReason: anyNamed('localizedReason'),
          options: anyNamed('options'),
        )).thenAnswer((_) async => true);

        when(mockSecureStorage.write(
          key: 'biometric_enabled_$email',
          value: 'true',
        )).thenAnswer((_) async {});

        // Act
        final result = await authService.setBiometricEnabled(email, true);

        // Assert
        expect(result, true);
        verify(mockLocalAuth.authenticate(
          localizedReason: anyNamed('localizedReason'),
          options: anyNamed('options'),
        )).called(1);
        verify(mockSecureStorage.write(
          key: 'biometric_enabled_$email',
          value: 'true',
        )).called(1);
      });

      test('should disable biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';

        when(mockSecureStorage.write(
          key: 'biometric_enabled_$email',
          value: 'false',
        )).thenAnswer((_) async {});

        // Act
        final result = await authService.setBiometricEnabled(email, false);

        // Assert
        expect(result, true);
        verify(mockSecureStorage.write(
          key: 'biometric_enabled_$email',
          value: 'false',
        )).called(1);
        
        // Should not authenticate when disabling
        verifyNever(mockLocalAuth.authenticate(
          localizedReason: anyNamed('localizedReason'),
          options: anyNamed('options'),
        ));
      });

      test('should check biometric availability', () async {
        // Arrange
        when(mockLocalAuth.canCheckBiometrics)
            .thenAnswer((_) async => true);

        // Act
        final canUse = await authService.canUseBiometrics();

        // Assert
        expect(canUse, true);
      });

      test('should get available biometric types', () async {
        // Arrange
        when(mockLocalAuth.getAvailableBiometrics())
            .thenAnswer((_) async => [BiometricType.fingerprint, BiometricType.face]);

        // Act
        final biometrics = await authService.getAvailableBiometrics();

        // Assert
        expect(biometrics.length, 2);
        expect(biometrics, contains(BiometricType.fingerprint));
        expect(biometrics, contains(BiometricType.face));
      });
    });

    group('Token Refresh', () {
      test('should handle token expiration automatically', () async {
        // Arrange
        final sessionData = SessionData(
          accessToken: 'expired_token',
          refreshToken: 'valid_refresh_token',
          sessionId: 'session_123',
          user: User(
            id: 'user123',
            email: 'test@example.com',
            firstName: 'Test',
            lastName: 'User',
            role: UserRole.customer,
            isVerified: true,
          ),
          tokenExpiresAt: DateTime.now().subtract(Duration(minutes: 1)), // Expired
          refreshExpiresAt: DateTime.now().add(Duration(days: 6)), // Valid
          deviceTrusted: true,
          rememberMe: false,
        );

        final refreshResponse = {
          'success': true,
          'access_token': 'new_access_token',
          'refresh_token': 'new_refresh_token',
          'token_expires_at': DateTime.now().add(Duration(minutes: 30)).toIso8601String(),
          'refresh_expires_at': DateTime.now().add(Duration(days: 7)).toIso8601String(),
        };

        when(mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => sessionData);

        when(mockApiService.refreshToken(any, any))
            .thenAnswer((_) async => refreshResponse);

        when(mockSessionManager.storeSession(any))
            .thenAnswer((_) async {});

        // Act
        final result = await authService.handleTokenExpiration();

        // Assert
        expect(result, true);
        verify(mockApiService.refreshToken('valid_refresh_token', 'session_123')).called(1);
      });

      test('should return false when refresh token is expired', () async {
        // Arrange
        when(mockSessionManager.getCurrentSession())
            .thenAnswer((_) async => null);

        // Act
        final result = await authService.handleTokenExpiration();

        // Assert
        expect(result, false);
      });
    });
  });
}
