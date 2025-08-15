import 'dart:async';
import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../core/services/api_service.dart';

class AuthService {
  final ApiService _apiService;
  final FlutterSecureStorage _secureStorage;
  
  // Prevent multiple simultaneous refresh attempts
  Completer<String>? _refreshCompleter;
  
  AuthService({
    required ApiService apiService,
    required FlutterSecureStorage secureStorage,
  }) : _apiService = apiService, _secureStorage = secureStorage;

  Future<String> getValidToken() async {
    final accessToken = await _secureStorage.read(key: 'access_token');
    
    if (accessToken == null) {
      throw Exception('No access token found');
    }
    
    // Check if token is still valid (with buffer time)
    if (await _isTokenValid(accessToken)) {
      return accessToken;
    }
    
    // If not valid, refresh it
    return await refreshToken();
  }

  Future<String> refreshToken() async {
    // If a refresh is already in progress, wait for it
    if (_refreshCompleter != null && !_refreshCompleter!.isCompleted) {
      return await _refreshCompleter!.future;
    }
    
    // Start new refresh process
    _refreshCompleter = Completer<String>();
    
    try {
      final refreshToken = await _secureStorage.read(key: 'refresh_token');
      final sessionId = await _secureStorage.read(key: 'session_id');
      
      if (refreshToken == null) {
        throw Exception('No refresh token found');
      }
      
      final response = await _apiService.refreshToken(refreshToken, sessionId ?? '');
      
      // Store new tokens
      await _storeTokens(response);
      
      // Complete the future with new access token
      final newAccessToken = response['access_token'] as String;
      _refreshCompleter!.complete(newAccessToken);
      return newAccessToken;
      
    } catch (error) {
      _refreshCompleter!.completeError(error);
      rethrow;
    } finally {
      _refreshCompleter = null;
    }
  }

  Future<bool> _isTokenValid(String token) async {
    try {
      // Decode JWT and check expiration
      final parts = token.split('.');
      if (parts.length != 3) return false;
      
      final payload = jsonDecode(
        utf8.decode(base64Url.decode(base64Url.normalize(parts[1])))
      );
      
      final exp = payload['exp'] as int;
      final expiryDate = DateTime.fromMillisecondsSinceEpoch(exp * 1000);
      
      // Add 5 minute buffer to prevent edge cases
      final bufferTime = DateTime.now().add(const Duration(minutes: 5));
      
      return expiryDate.isAfter(bufferTime);
    } catch (e) {
      return false;
    }
  }

  Future<void> _storeTokens(Map<String, dynamic> response) async {
    await Future.wait([
      _secureStorage.write(key: 'access_token', value: response['access_token']),
      _secureStorage.write(key: 'refresh_token', value: response['refresh_token']),
      _secureStorage.write(
        key: 'refresh_expires_at', 
        value: response['refresh_expires_at']
      ),
      _secureStorage.write(
        key: 'token_expires_at', 
        value: response['token_expires_at']
      ),
      _secureStorage.write(
        key: 'session_id', 
        value: response['session_id'] ?? ''
      ),
    ]);
  }

  // Legacy methods for compatibility
  String? getCurrentAccessToken() => null; // Will be handled by getValidToken()
  void setCachedAccessToken(String? token) {} // Not needed with new approach
  Future<Map<String, dynamic>?> refreshTokenIfNeeded() async {
    try {
      await getValidToken();
      return {'success': true};
    } catch (e) {
      return null;
    }
  }
}