import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/user.dart';

part 'session_manager.g.dart';

/// Manages secure session storage and retrieval
class SessionManager {
  final FlutterSecureStorage _secureStorage;
  
  static const String _sessionKey = 'smor_ting_session_v2';
  static const String _biometricTokenKey = 'smor_ting_biometric_token';
  static const String _rememberMeKey = 'smor_ting_remember_me';

  SessionManager(this._secureStorage);

  /// Store session data securely
  Future<void> storeSession(SessionData sessionData) async {
    try {
      final sessionJson = jsonEncode(sessionData.toJson());
      await _secureStorage.write(key: _sessionKey, value: sessionJson);
      
      // Store remember me preference separately
      await _secureStorage.write(
        key: _rememberMeKey, 
        value: sessionData.rememberMe.toString(),
      );
    } catch (e) {
      throw SessionException('Failed to store session: ${e.toString()}');
    }
  }

  /// Get current session data
  Future<SessionData?> getCurrentSession() async {
    try {
      final sessionJson = await _secureStorage.read(key: _sessionKey);
      if (sessionJson == null) return null;
      
      final sessionMap = jsonDecode(sessionJson) as Map<String, dynamic>;
      return SessionData.fromJson(sessionMap);
    } catch (e) {
      // If session data is corrupted, clear it
      await clearSession();
      return null;
    }
  }

  /// Check if remember me is enabled
  Future<bool> isRememberMeEnabled() async {
    try {
      final rememberMe = await _secureStorage.read(key: _rememberMeKey);
      return rememberMe == 'true';
    } catch (e) {
      return false;
    }
  }

  /// Store biometric token for quick unlock
  Future<void> storeBiometricToken(String token, String userId) async {
    try {
      final biometricData = {
        'token': token,
        'user_id': userId,
        'created_at': DateTime.now().toIso8601String(),
      };
      
      await _secureStorage.write(
        key: _biometricTokenKey,
        value: jsonEncode(biometricData),
      );
    } catch (e) {
      throw SessionException('Failed to store biometric token: ${e.toString()}');
    }
  }

  /// Get biometric token for quick unlock
  Future<Map<String, dynamic>?> getBiometricToken() async {
    try {
      final tokenJson = await _secureStorage.read(key: _biometricTokenKey);
      if (tokenJson == null) return null;
      
      final tokenData = jsonDecode(tokenJson) as Map<String, dynamic>;
      
      // Check if biometric token is expired (7 days)
      final createdAt = DateTime.parse(tokenData['created_at']);
      if (DateTime.now().difference(createdAt).inDays > 7) {
        await _secureStorage.delete(key: _biometricTokenKey);
        return null;
      }
      
      return tokenData;
    } catch (e) {
      await _secureStorage.delete(key: _biometricTokenKey);
      return null;
    }
  }

  /// Clear all session data
  Future<void> clearSession() async {
    try {
      await Future.wait([
        _secureStorage.delete(key: _sessionKey),
        _secureStorage.delete(key: _biometricTokenKey),
        _secureStorage.delete(key: _rememberMeKey),
      ]);
    } catch (e) {
      // Best effort cleanup
    }
  }

  /// Update session with new tokens
  Future<void> updateTokens({
    required String accessToken,
    required String refreshToken,
    required DateTime tokenExpiresAt,
    required DateTime refreshExpiresAt,
  }) async {
    try {
      final currentSession = await getCurrentSession();
      if (currentSession == null) {
        throw SessionException('No current session to update');
      }
      
      final updatedSession = currentSession.copyWith(
        accessToken: accessToken,
        refreshToken: refreshToken,
        tokenExpiresAt: tokenExpiresAt,
        refreshExpiresAt: refreshExpiresAt,
      );
      
      await storeSession(updatedSession);
    } catch (e) {
      throw SessionException('Failed to update tokens: ${e.toString()}');
    }
  }

  /// Check if session exists and is valid
  Future<bool> hasValidSession() async {
    try {
      final session = await getCurrentSession();
      if (session == null) return false;
      
      // Check if access token is expired
      final now = DateTime.now();
      if (now.isAfter(session.tokenExpiresAt)) {
        // Check if refresh token is still valid
        return now.isBefore(session.refreshExpiresAt);
      }
      
      return true;
    } catch (e) {
      return false;
    }
  }

  /// Get session expiry information
  Future<SessionExpiryInfo?> getSessionExpiryInfo() async {
    try {
      final session = await getCurrentSession();
      if (session == null) return null;
      
      final now = DateTime.now();
      
      return SessionExpiryInfo(
        accessTokenExpired: now.isAfter(session.tokenExpiresAt),
        refreshTokenExpired: now.isAfter(session.refreshExpiresAt),
        accessTokenExpiresAt: session.tokenExpiresAt,
        refreshTokenExpiresAt: session.refreshExpiresAt,
        timeUntilAccessTokenExpiry: session.tokenExpiresAt.difference(now),
        timeUntilRefreshTokenExpiry: session.refreshExpiresAt.difference(now),
      );
    } catch (e) {
      return null;
    }
  }

  /// Clear expired sessions
  Future<void> cleanupExpiredSessions() async {
    try {
      final session = await getCurrentSession();
      if (session == null) return;
      
      final now = DateTime.now();
      
      // If refresh token is expired, clear the session
      if (now.isAfter(session.refreshExpiresAt)) {
        await clearSession();
      }
    } catch (e) {
      // Best effort cleanup
      await clearSession();
    }
  }

  /// Get all stored user emails (for quick login)
  Future<List<String>> getStoredUserEmails() async {
    try {
      final keys = await _secureStorage.readAll();
      final userEmails = <String>[];
      
      for (final key in keys.keys) {
        if (key.startsWith('biometric_enabled_')) {
          final email = key.replaceFirst('biometric_enabled_', '');
          if (keys[key] == 'true') {
            userEmails.add(email);
          }
        }
      }
      
      return userEmails;
    } catch (e) {
      return [];
    }
  }
}

/// Session data model
class SessionData {
  final String accessToken;
  final String refreshToken;
  final String sessionId;
  final User user;
  final DateTime tokenExpiresAt;
  final DateTime refreshExpiresAt;
  final bool deviceTrusted;
  final bool rememberMe;
  final DateTime createdAt;

  SessionData({
    required this.accessToken,
    required this.refreshToken,
    required this.sessionId,
    required this.user,
    required this.tokenExpiresAt,
    required this.refreshExpiresAt,
    required this.deviceTrusted,
    required this.rememberMe,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  Map<String, dynamic> toJson() {
    return {
      'access_token': accessToken,
      'refresh_token': refreshToken,
      'session_id': sessionId,
      'user': user.toJson(),
      'token_expires_at': tokenExpiresAt.toIso8601String(),
      'refresh_expires_at': refreshExpiresAt.toIso8601String(),
      'device_trusted': deviceTrusted,
      'remember_me': rememberMe,
      'created_at': createdAt.toIso8601String(),
    };
  }

  factory SessionData.fromJson(Map<String, dynamic> json) {
    return SessionData(
      accessToken: json['access_token'],
      refreshToken: json['refresh_token'],
      sessionId: json['session_id'],
      user: User.fromJson(json['user']),
      tokenExpiresAt: DateTime.parse(json['token_expires_at']),
      refreshExpiresAt: DateTime.parse(json['refresh_expires_at']),
      deviceTrusted: json['device_trusted'] ?? false,
      rememberMe: json['remember_me'] ?? false,
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at'])
          : DateTime.now(),
    );
  }

  SessionData copyWith({
    String? accessToken,
    String? refreshToken,
    String? sessionId,
    User? user,
    DateTime? tokenExpiresAt,
    DateTime? refreshExpiresAt,
    bool? deviceTrusted,
    bool? rememberMe,
    DateTime? createdAt,
  }) {
    return SessionData(
      accessToken: accessToken ?? this.accessToken,
      refreshToken: refreshToken ?? this.refreshToken,
      sessionId: sessionId ?? this.sessionId,
      user: user ?? this.user,
      tokenExpiresAt: tokenExpiresAt ?? this.tokenExpiresAt,
      refreshExpiresAt: refreshExpiresAt ?? this.refreshExpiresAt,
      deviceTrusted: deviceTrusted ?? this.deviceTrusted,
      rememberMe: rememberMe ?? this.rememberMe,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  /// Check if access token needs refresh (within 5 minutes of expiry)
  bool get needsRefresh {
    final now = DateTime.now();
    final refreshThreshold = tokenExpiresAt.subtract(const Duration(minutes: 5));
    return now.isAfter(refreshThreshold);
  }

  /// Check if session is expired
  bool get isExpired {
    return DateTime.now().isAfter(refreshExpiresAt);
  }

  /// Get time until session expires
  Duration get timeUntilExpiry {
    return refreshExpiresAt.difference(DateTime.now());
  }
}

/// Session expiry information
class SessionExpiryInfo {
  final bool accessTokenExpired;
  final bool refreshTokenExpired;
  final DateTime accessTokenExpiresAt;
  final DateTime refreshTokenExpiresAt;
  final Duration timeUntilAccessTokenExpiry;
  final Duration timeUntilRefreshTokenExpiry;

  SessionExpiryInfo({
    required this.accessTokenExpired,
    required this.refreshTokenExpired,
    required this.accessTokenExpiresAt,
    required this.refreshTokenExpiresAt,
    required this.timeUntilAccessTokenExpiry,
    required this.timeUntilRefreshTokenExpiry,
  });

  bool get needsRefresh => accessTokenExpired && !refreshTokenExpired;
  bool get sessionExpired => refreshTokenExpired;
}

/// Session exception
class SessionException implements Exception {
  final String message;
  SessionException(this.message);
  
  @override
  String toString() => 'SessionException: $message';
}

/// Riverpod provider for session manager
@riverpod
SessionManager sessionManager(SessionManagerRef ref) {
  return SessionManager(
    const FlutterSecureStorage(
      aOptions: AndroidOptions(
        encryptedSharedPreferences: true,
      ),
      iOptions: IOSOptions(
        accessibility: KeychainAccessibility.first_unlock_this_device,
      ),
    ),
  );
}
