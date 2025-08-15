import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/models/user.dart';
import '../../../../core/services/session_manager.dart';
import '../../../../services/auth_service.dart';
import '../../../../core/services/api_service.dart';

part 'simple_auth_provider.g.dart';

/// Simple authentication state
@freezed
class SimpleAuthState with _$SimpleAuthState {
  const factory SimpleAuthState.initial() = _Initial;
  const factory SimpleAuthState.loading() = _Loading;
  const factory SimpleAuthState.authenticated({
    required User user,
    required String accessToken,
    required String sessionId,
  }) = _Authenticated;
  const factory SimpleAuthState.unauthenticated() = _Unauthenticated;
  const factory SimpleAuthState.error(String message) = _Error;
}

/// Simple authentication notifier
@riverpod
class SimpleAuthNotifier extends _$SimpleAuthNotifier {
  @override
  SimpleAuthState build() {
    // Check for existing session on startup
    _checkSession();
    return const SimpleAuthState.initial();
  }

  /// Check if user has valid session
  Future<void> _checkSession() async {
    try {
      final sessionManager = ref.read(sessionManagerProvider);
      final hasValid = await sessionManager.hasValidSession();
      
      if (hasValid) {
        final session = await sessionManager.getCurrentSession();
        if (session != null) {
          state = SimpleAuthState.authenticated(
            user: session.user,
            accessToken: session.accessToken,
            sessionId: session.sessionId,
          );
          
          // Set token in API service
          ref.read(apiServiceProvider).setAuthToken(session.accessToken);
          return;
        }
      }
      
      state = const SimpleAuthState.unauthenticated();
    } catch (e) {
      state = const SimpleAuthState.unauthenticated();
    }
  }

  /// Simple login
  Future<void> login(String email, String password) async {
    state = const SimpleAuthState.loading();
    
    try {
      final apiService = ref.read(apiServiceProvider);
      
      // Basic login request
      final response = await apiService.login(email, password);
      
      if (response['success'] == true) {
        final userData = response['user'];
        final user = User.fromJson(userData);
        
        // Store session
        final sessionManager = ref.read(sessionManagerProvider);
        final sessionData = SessionData(
          accessToken: response['access_token'],
          refreshToken: response['refresh_token'],
          sessionId: response['session_id'] ?? '',
          user: user,
          tokenExpiresAt: DateTime.parse(response['token_expires_at']),
          refreshExpiresAt: DateTime.parse(response['refresh_expires_at']),
          deviceTrusted: false,
          rememberMe: false,
        );
        
        await sessionManager.storeSession(sessionData);
        apiService.setAuthToken(sessionData.accessToken);
        
        state = SimpleAuthState.authenticated(
          user: user,
          accessToken: sessionData.accessToken,
          sessionId: sessionData.sessionId,
        );
      } else {
        state = SimpleAuthState.error(response['message'] ?? 'Login failed');
      }
    } catch (e) {
      state = SimpleAuthState.error('Login failed: ${e.toString()}');
    }
  }

  /// Logout
  Future<void> logout() async {
    try {
      final sessionManager = ref.read(sessionManagerProvider);
      await sessionManager.clearSession();
      
      final apiService = ref.read(apiServiceProvider);
      apiService.clearAuthToken();
      
      state = const SimpleAuthState.unauthenticated();
    } catch (e) {
      // Still set to unauthenticated even if logout fails
      state = const SimpleAuthState.unauthenticated();
    }
  }
}

/// Provider for simple auth state
@riverpod
SimpleAuthState simpleAuthState(SimpleAuthStateRef ref) {
  return ref.watch(simpleAuthNotifierProvider);
}
