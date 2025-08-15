import 'dart:convert';
import 'dart:io';
import 'dart:async';
import 'package:crypto/crypto.dart';
import 'package:device_info_plus/device_info_plus.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:local_auth/local_auth.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:jwt_decoder/jwt_decoder.dart';
import '../models/user.dart';
import '../models/enhanced_auth_models.dart';
import 'api_service.dart';
import 'session_manager.dart';
import 'device_fingerprint_service.dart';

part 'enhanced_auth_service.g.dart';

/// Enhanced authentication service with comprehensive security features
class EnhancedAuthService {
  final ApiService _apiService;
  final SessionManager _sessionManager;
  final DeviceFingerprintService _deviceService;
  final FlutterSecureStorage _secureStorage;
  final LocalAuthentication _localAuth;

  // Debouncing mechanism for refresh tokens
  Completer<EnhancedAuthResult?>? _refreshCompleter;
  DateTime? _lastRefreshAttempt;
  static const Duration _refreshDebounceThreshold = Duration(milliseconds: 500);
  
  // Error handling and retry mechanism
  static const int _maxRetryAttempts = 3;
  static const Duration _baseRetryDelay = Duration(milliseconds: 1000);

  EnhancedAuthService(
    this._apiService,
    this._sessionManager,
    this._deviceService,
    this._secureStorage,
    this._localAuth,
  );

  /// Enhanced login with comprehensive security features
  Future<EnhancedAuthResult> enhancedLogin({
    required String email,
    required String password,
    bool rememberMe = false,
    String? captchaToken,
    String? twoFactorCode,
  }) async {
    try {
      // Get device fingerprint
      final deviceInfo = await _deviceService.generateFingerprint();
      
      // Prepare enhanced login request
      final request = EnhancedLoginRequest(
        email: email,
        password: password,
        rememberMe: rememberMe,
        deviceInfo: deviceInfo,
        captchaToken: captchaToken,
        twoFactorCode: twoFactorCode,
      );

      // Perform login
      final response = await _apiService.enhancedLogin(request);
      final authResult = EnhancedAuthResult.fromResponse(response);
      
      if (authResult.success) {
        // Store session data securely
        if (authResult.accessToken != null && authResult.refreshToken != null) {
          await _sessionManager.storeSession(SessionData(
            accessToken: authResult.accessToken!,
            refreshToken: authResult.refreshToken!,
            sessionId: authResult.sessionId!,
            user: authResult.user!,
            tokenExpiresAt: authResult.tokenExpiresAt!,
            refreshExpiresAt: authResult.refreshExpiresAt!,
            deviceTrusted: authResult.deviceTrusted,
            rememberMe: rememberMe,
          ));

          // Store biometric preference if device is trusted
          if (authResult.deviceTrusted && await _localAuth.canCheckBiometrics) {
            await _storeBiometricPreference(email, rememberMe);
          }
        }
      }

      return authResult;
    } catch (e) {
      throw AuthException('Login failed: ${e.toString()}');
    }
  }

  /// Restore session on app launch
  Future<EnhancedAuthResult?> restoreSession() async {
    try {
      final sessionData = await _sessionManager.getCurrentSession();
      if (sessionData == null) {
        return null;
      }

      // Check if session is still valid - use needsRefresh instead of direct comparison
      if (sessionData.needsRefresh) {
        // Try to refresh token
        final refreshResult = await _refreshToken();
        if (refreshResult != null) {
          return refreshResult;
        }
        
        // If refresh fails, check if we can use biometric unlock
        if (sessionData.rememberMe && sessionData.deviceTrusted) {
          final biometricResult = await _tryBiometricUnlock(sessionData);
          if (biometricResult != null) {
            return biometricResult;
          }
        }
        
        // Session is invalid and refresh failed
        await _sessionManager.clearSession();
        return null;
      }

      // Session is still valid
      _apiService.setAuthToken(sessionData.accessToken);
      return EnhancedAuthResult(
        success: true,
        user: sessionData.user,
        accessToken: sessionData.accessToken,
        refreshToken: sessionData.refreshToken,
        sessionId: sessionData.sessionId,
        deviceTrusted: sessionData.deviceTrusted,
        isRestoredSession: true,
      );
    } catch (e) {
      // Clear corrupted session data
      await _sessionManager.clearSession();
      return null;
    }
  }

  /// Refresh authentication token using the new cleaner auth service
  Future<EnhancedAuthResult?> _refreshToken() async {
    try {
      // Use the new auth service which has proper infinite loop prevention
      final newAccessToken = await _apiService.authService.refreshToken();
      
      // Get the updated session data
      final sessionData = await _sessionManager.getCurrentSession();
      if (sessionData == null) return null;

      return EnhancedAuthResult(
        success: true,
        user: sessionData.user,
        accessToken: newAccessToken,
        refreshToken: sessionData.refreshToken,
        sessionId: sessionData.sessionId,
        deviceTrusted: sessionData.deviceTrusted,
        isRestoredSession: true,
      );
    } catch (e) {
      // Clear session on failure
      await _sessionManager.clearSession();
      return null;
    }
  }

  /// Perform token refresh with retry logic and comprehensive error handling
  Future<EnhancedAuthResult?> _performRefreshWithRetry() async {
    int attemptCount = 0;
    Exception? lastException;

    while (attemptCount < _maxRetryAttempts) {
      try {
        final sessionData = await _sessionManager.getCurrentSession();
        if (sessionData == null) return null;

        // JWT VALIDATION: Validate refresh token before making API call
        if (!_isValidJwtToken(sessionData.refreshToken)) {
          await _sessionManager.clearSession();
          return null;
        }

        // JWT VALIDATION: Check if refresh token is expired
        if (DateTime.now().isAfter(sessionData.refreshExpiresAt)) {
          await _sessionManager.clearSession();
          return null;
        }

        // JWT VALIDATION: Additional JWT expiry check using decoder
        if (_isJwtExpired(sessionData.refreshToken)) {
          await _sessionManager.clearSession();
          return null;
        }

        // Attempt token refresh
        final response = await _apiService.refreshToken(
          sessionData.refreshToken,
          sessionData.sessionId,
        );

        // ERROR HANDLING: Validate response format
        if (!_isValidRefreshResponse(response)) {
          throw Exception('Invalid refresh response format');
        }

        if (response['success'] == true) {
          // JWT VALIDATION: Validate new tokens before storing
          final newAccessToken = response['access_token'] as String?;
          final newRefreshToken = response['refresh_token'] as String?;
          
          if (newAccessToken == null || newRefreshToken == null ||
              !_isValidJwtToken(newAccessToken) || !_isValidJwtToken(newRefreshToken)) {
            throw Exception('Invalid JWT tokens in refresh response');
          }

          // Update stored session with new tokens
          final updatedSession = sessionData.copyWith(
            accessToken: newAccessToken,
            refreshToken: newRefreshToken,
            tokenExpiresAt: DateTime.parse(response['token_expires_at']),
            refreshExpiresAt: DateTime.parse(response['refresh_expires_at']),
          );

          await _sessionManager.storeSession(updatedSession);
          _apiService.setAuthToken(updatedSession.accessToken);

          return EnhancedAuthResult(
            success: true,
            user: updatedSession.user,
            accessToken: updatedSession.accessToken,
            refreshToken: updatedSession.refreshToken,
            sessionId: updatedSession.sessionId,
            deviceTrusted: updatedSession.deviceTrusted,
            isRestoredSession: true,
          );
        } else {
          // ERROR HANDLING: Server returned failure
          throw Exception('Server rejected refresh request: ${response['message'] ?? 'Unknown error'}');
        }
      } catch (e) {
        lastException = e is Exception ? e : Exception(e.toString());
        attemptCount++;

        // ERROR HANDLING: Don't retry on certain errors
        if (_isNonRetryableError(e)) {
          await _sessionManager.clearSession();
          break;
        }

        // ERROR HANDLING: Exponential backoff for retries
        if (attemptCount < _maxRetryAttempts) {
          final delay = Duration(
            milliseconds: _baseRetryDelay.inMilliseconds * (1 << (attemptCount - 1))
          );
          await Future.delayed(delay);
        }
      }
    }

    // ERROR HANDLING: All retries failed
    if (lastException != null && _isSessionInvalidatingError(lastException)) {
      await _sessionManager.clearSession();
    }
    
    return null;
  }

  /// JWT VALIDATION: Check if a token has a valid JWT structure
  bool _isValidJwtToken(String token) {
    if (token.isEmpty) return false;
    
    try {
      // JWT tokens should have 3 parts separated by dots
      final parts = token.split('.');
      if (parts.length != 3) return false;
      
      // Try to decode the header and payload
      final header = utf8.decode(base64Decode(_padBase64(parts[0])));
      final payload = utf8.decode(base64Decode(_padBase64(parts[1])));
      
      // Basic validation - should be valid JSON
      jsonDecode(header);
      jsonDecode(payload);
      
      return true;
    } catch (e) {
      return false;
    }
  }

  /// JWT VALIDATION: Check if JWT token is expired using decoder
  bool _isJwtExpired(String token) {
    try {
      return JwtDecoder.isExpired(token);
    } catch (e) {
      // If we can't decode it, consider it expired
      return true;
    }
  }

  /// ERROR HANDLING: Validate refresh response structure
  bool _isValidRefreshResponse(Map<String, dynamic> response) {
    if (response['success'] == true) {
      return response.containsKey('access_token') &&
             response.containsKey('refresh_token') &&
             response.containsKey('token_expires_at') &&
             response.containsKey('refresh_expires_at') &&
             response['access_token'] is String &&
             response['refresh_token'] is String &&
             response['token_expires_at'] is String &&
             response['refresh_expires_at'] is String;
    }
    return response.containsKey('success');
  }

  /// ERROR HANDLING: Check if error should not be retried
  bool _isNonRetryableError(dynamic error) {
    final errorMessage = error.toString().toLowerCase();
    return errorMessage.contains('invalid') ||
           errorMessage.contains('expired') ||
           errorMessage.contains('unauthorized') ||
           errorMessage.contains('forbidden') ||
           errorMessage.contains('malformed');
  }

  /// ERROR HANDLING: Check if error should invalidate session
  bool _isSessionInvalidatingError(Exception error) {
    final errorMessage = error.toString().toLowerCase();
    return errorMessage.contains('invalid') ||
           errorMessage.contains('expired') ||
           errorMessage.contains('unauthorized') ||
           errorMessage.contains('revoked');
  }

  /// Utility: Pad base64 string for proper decoding
  String _padBase64(String base64) {
    final missingPadding = 4 - (base64.length % 4);
    if (missingPadding != 4) {
      return base64 + ('=' * missingPadding);
    }
    return base64;
  }

  /// Biometric authentication for quick unlock
  Future<EnhancedAuthResult?> _tryBiometricUnlock(SessionData sessionData) async {
    try {
      // Check if biometric unlock is enabled for this user
      final biometricEnabled = await _isBiometricEnabled(sessionData.user.email);
      if (!biometricEnabled) return null;

      // Check if biometrics are available
      final isAvailable = await _localAuth.canCheckBiometrics;
      if (!isAvailable) return null;

      // Get available biometric types
      final availableBiometrics = await _localAuth.getAvailableBiometrics();
      if (availableBiometrics.isEmpty) return null;

      // Perform biometric authentication
      final authenticated = await _localAuth.authenticate(
        localizedReason: 'Unlock Smor-Ting with biometric authentication',
        options: const AuthenticationOptions(
          biometricOnly: true,
          stickyAuth: true,
        ),
      );

      if (!authenticated) return null;

      // Generate new session via biometric login
      final deviceInfo = await _deviceService.generateFingerprint();
      
      // Call biometric login endpoint
      final response = await _apiService.biometricLogin(
        sessionData.user.email,
        sessionData.sessionId,
        deviceInfo,
      );

      if (response.success && response.accessToken != null) {
        final updatedSession = sessionData.copyWith(
          accessToken: response.accessToken!,
          refreshToken: response.refreshToken!,
          tokenExpiresAt: response.tokenExpiresAt!,
          refreshExpiresAt: response.refreshExpiresAt!,
        );

        await _sessionManager.storeSession(updatedSession);
        _apiService.setAuthToken(updatedSession.accessToken);

        return response;
      }

      return null;
    } catch (e) {
      return null;
    }
  }

  /// Get all active sessions for current user
  Future<List<SessionInfo>> getUserSessions() async {
    try {
      final response = await _apiService.getUserSessions();
      return (response['sessions'] as List)
          .map((session) => SessionInfo.fromJson(session))
          .toList();
    } catch (e) {
      throw AuthException('Failed to get sessions: ${e.toString()}');
    }
  }

  /// Revoke a specific session
  Future<void> revokeSession(String sessionId) async {
    try {
      await _apiService.revokeSession(sessionId);
    } catch (e) {
      throw AuthException('Failed to revoke session: ${e.toString()}');
    }
  }

  /// Sign out from all devices
  Future<void> signOutAllDevices() async {
    try {
      await _apiService.revokeAllSessions();
      await _sessionManager.clearSession();
      _apiService.clearAuthToken();
    } catch (e) {
      throw AuthException('Failed to sign out all devices: ${e.toString()}');
    }
  }

  /// Enhanced logout with session cleanup
  Future<void> enhancedLogout({bool revokeSession = true}) async {
    try {
      if (revokeSession) {
        final sessionData = await _sessionManager.getCurrentSession();
        if (sessionData != null) {
          await _apiService.revokeSession(sessionData.sessionId);
        }
      }
    } catch (e) {
      // Log error but continue with local cleanup
      debugPrint('Error revoking session: $e');
    }

    // Always clear local session data
    await _sessionManager.clearSession();
    _apiService.clearAuthToken();
  }

  /// Enable/disable biometric authentication
  Future<bool> setBiometricEnabled(String email, bool enabled) async {
    try {
      if (enabled) {
        // Check if biometrics are available
        final isAvailable = await _localAuth.canCheckBiometrics;
        if (!isAvailable) return false;

        // Test biometric authentication
        final authenticated = await _localAuth.authenticate(
          localizedReason: 'Enable biometric authentication for Smor-Ting',
          options: const AuthenticationOptions(
            biometricOnly: true,
            stickyAuth: true,
          ),
        );

        if (!authenticated) return false;
      }

      await _secureStorage.write(
        key: 'biometric_enabled_$email',
        value: enabled.toString(),
      );

      return true;
    } catch (e) {
      return false;
    }
  }

  /// Check if biometric authentication is enabled (public method)
  Future<bool> isBiometricEnabled(String email) async {
    return await _isBiometricEnabled(email);
  }

  /// Check if biometric authentication is enabled (private method)
  Future<bool> _isBiometricEnabled(String email) async {
    try {
      final value = await _secureStorage.read(key: 'biometric_enabled_$email');
      return value == 'true';
    } catch (e) {
      return false;
    }
  }

  /// Store biometric preference
  Future<void> _storeBiometricPreference(String email, bool rememberMe) async {
    if (rememberMe) {
      await _secureStorage.write(
        key: 'biometric_preference_$email',
        value: 'true',
      );
    }
  }

  /// Get available biometric types
  Future<List<BiometricType>> getAvailableBiometrics() async {
    try {
      return await _localAuth.getAvailableBiometrics();
    } catch (e) {
      return [];
    }
  }

  /// Check if device supports biometrics
  Future<bool> canUseBiometrics() async {
    try {
      return await _localAuth.canCheckBiometrics;
    } catch (e) {
      return false;
    }
  }

  /// Get current session info
  Future<SessionData?> getCurrentSession() async {
    return await _sessionManager.getCurrentSession();
  }

  /// Authenticate using biometric authentication for quick unlock
  Future<EnhancedAuthResult> authenticateWithBiometrics(String email) async {
    try {
      // Check if biometric unlock is enabled for this user
      final biometricEnabled = await _isBiometricEnabled(email);
      if (!biometricEnabled) {
        return EnhancedAuthResult(
          success: false,
          message: 'Biometric authentication is not enabled for this account',
        );
      }

      // Check if biometrics are available
      final isAvailable = await _localAuth.canCheckBiometrics;
      if (!isAvailable) {
        return EnhancedAuthResult(
          success: false,
          message: 'Biometric authentication is not available on this device',
        );
      }

      // Get available biometric types
      final availableBiometrics = await _localAuth.getAvailableBiometrics();
      if (availableBiometrics.isEmpty) {
        return EnhancedAuthResult(
          success: false,
          message: 'No biometric authentication methods are set up',
        );
      }

      // Perform biometric authentication
      final authenticated = await _localAuth.authenticate(
        localizedReason: 'Unlock Smor-Ting with biometric authentication',
        options: const AuthenticationOptions(
          biometricOnly: true,
          stickyAuth: true,
        ),
      );

      if (!authenticated) {
        return EnhancedAuthResult(
          success: false,
          message: 'Biometric authentication was cancelled or failed',
        );
      }

      // Get cached session for biometric unlock
      final cachedSession = await _sessionManager.getCurrentSession();
      if (cachedSession?.user.email != email) {
        return EnhancedAuthResult(
          success: false,
          message: 'No valid session found for biometric unlock',
        );
      }

      // Generate new session via biometric login API
      final deviceInfo = await _deviceService.generateFingerprint();
      
      // Call biometric login endpoint
      final response = await _apiService.biometricLogin(
        email,
        cachedSession!.sessionId,
        deviceInfo,
      );

      if (response.success && response.accessToken != null) {
        final updatedSession = cachedSession.copyWith(
          accessToken: response.accessToken!,
          refreshToken: response.refreshToken!,
          tokenExpiresAt: response.tokenExpiresAt!,
          refreshExpiresAt: response.refreshExpiresAt!,
        );

        await _sessionManager.storeSession(updatedSession);
        _apiService.setAuthToken(updatedSession.accessToken);

        return EnhancedAuthResult(
          success: true,
          user: response.user,
          accessToken: response.accessToken,
          refreshToken: response.refreshToken,
          sessionId: response.sessionId,
          tokenExpiresAt: response.tokenExpiresAt,
          refreshExpiresAt: response.refreshExpiresAt,
          deviceTrusted: true,
          isRestoredSession: true,
          message: 'Biometric authentication successful',
        );
      }

      return EnhancedAuthResult(
        success: false,
        message: 'Failed to authenticate with server',
      );
    } catch (e) {
      return EnhancedAuthResult(
        success: false,
        message: 'Biometric authentication failed: ${e.toString()}',
      );
    }
  }

  /// Check if user is authenticated
  Future<bool> isAuthenticated() async {
    final session = await _sessionManager.getCurrentSession();
    if (session == null) return false;
    
    // Check if token is expired
    if (DateTime.now().isAfter(session.tokenExpiresAt)) {
      // Try to refresh
      final refreshResult = await _refreshToken();
      return refreshResult != null;
    }
    
    return true;
  }

  /// Auto-refresh token on 401 errors
  Future<bool> handleTokenExpiration() async {
    final refreshResult = await _refreshToken();
    return refreshResult != null;
  }
}



/// Authentication exception
class AuthException implements Exception {
  final String message;
  AuthException(this.message);
  
  @override
  String toString() => 'AuthException: $message';
}

/// Riverpod provider for enhanced auth service
@riverpod
EnhancedAuthService enhancedAuthService(EnhancedAuthServiceRef ref) {
  return EnhancedAuthService(
    ref.read(apiServiceProvider),
    ref.read(sessionManagerProvider),
    ref.read(deviceFingerprintServiceProvider),
    const FlutterSecureStorage(
      aOptions: AndroidOptions(
        encryptedSharedPreferences: true,
      ),
      iOptions: IOSOptions(
        accessibility: KeychainAccessibility.first_unlock_this_device,
      ),
    ),
    LocalAuthentication(),
  );
}
