import 'package:flutter/foundation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/models/user.dart';
import '../../../../core/services/enhanced_auth_service.dart';
import '../../../../core/services/session_manager.dart';

part 'enhanced_auth_provider.g.dart';

/// Enhanced authentication state with comprehensive security features
@riverpod
class EnhancedAuthNotifier extends _$EnhancedAuthNotifier {
  @override
  EnhancedAuthState build() {
    // Try to restore session on startup
    _restoreSession();
    return const EnhancedAuthState.initial();
  }

  /// Restore session from secure storage
  Future<void> _restoreSession() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final result = await authService.restoreSession();
      
      if (result != null && result.success) {
        state = EnhancedAuthState.authenticated(
          user: result.user!,
          accessToken: result.accessToken!,
          sessionId: result.sessionId!,
          deviceTrusted: result.deviceTrusted,
          isRestoredSession: result.isRestoredSession,
          requiresVerification: result.requiresVerification,
        );
      } else {
        state = const EnhancedAuthState.unauthenticated();
      }
    } catch (e) {
      debugPrint('Session restoration failed: $e');
      state = const EnhancedAuthState.unauthenticated();
    }
  }

  /// Enhanced login with comprehensive security
  Future<void> enhancedLogin({
    required String email,
    required String password,
    bool rememberMe = false,
    String? captchaToken,
    String? twoFactorCode,
  }) async {
    state = const EnhancedAuthState.loading();
    
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final result = await authService.enhancedLogin(
        email: email,
        password: password,
        rememberMe: rememberMe,
        captchaToken: captchaToken,
        twoFactorCode: twoFactorCode,
      );

      if (result.success) {
        // Handle different authentication states
        if (result.requiresTwoFactor) {
          state = EnhancedAuthState.requiresTwoFactor(
            email: email,
            tempUser: result.user!,
            deviceTrusted: result.deviceTrusted,
          );
        } else if (result.requiresCaptcha) {
          state = EnhancedAuthState.requiresCaptcha(
            email: email,
            remainingAttempts: result.remainingAttempts ?? 0,
            lockoutInfo: result.lockoutInfo,
          );
        } else if (result.user != null && result.accessToken != null) {
          state = EnhancedAuthState.authenticated(
            user: result.user!,
            accessToken: result.accessToken!,
            sessionId: result.sessionId!,
            deviceTrusted: result.deviceTrusted,
            requiresVerification: result.requiresVerification,
          );
        }
      } else {
        // Handle authentication failure
        if (result.lockoutInfo != null) {
          state = EnhancedAuthState.lockedOut(
            lockoutInfo: result.lockoutInfo!,
            message: result.message ?? 'Account temporarily locked',
          );
        } else if (result.requiresCaptcha) {
          state = EnhancedAuthState.requiresCaptcha(
            email: email,
            remainingAttempts: result.remainingAttempts ?? 0,
          );
        } else {
          state = EnhancedAuthState.error(
            result.message ?? 'Authentication failed',
            canRetry: true,
          );
        }
      }
    } catch (e) {
      state = EnhancedAuthState.error(
        'Login failed: ${e.toString()}',
        canRetry: true,
      );
    }
  }

  /// Verify two-factor authentication code
  Future<void> verifyTwoFactorCode(String email, String code) async {
    state = const EnhancedAuthState.loading();
    
    try {
      await enhancedLogin(
        email: email,
        password: '', // Password already verified
        twoFactorCode: code,
      );
    } catch (e) {
      state = EnhancedAuthState.error(
        'Two-factor verification failed: ${e.toString()}',
        canRetry: true,
      );
    }
  }

  /// Biometric authentication
  Future<void> authenticateWithBiometrics(String email) async {
    state = const EnhancedAuthState.loading();
    
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final sessionManager = ref.read(sessionManagerProvider);
      
      // Try to get cached session for biometric unlock
      final cachedSession = await sessionManager.getCurrentSession();
      if (cachedSession?.user.email != email) {
        throw Exception('No valid session for biometric unlock');
      }
      
      // This would trigger the biometric unlock flow in the auth service
      final result = await authService.restoreSession();
      
      if (result != null && result.success) {
        state = EnhancedAuthState.authenticated(
          user: result.user!,
          accessToken: result.accessToken!,
          sessionId: result.sessionId!,
          deviceTrusted: result.deviceTrusted,
          isRestoredSession: true,
        );
      } else {
        state = EnhancedAuthState.error(
          'Biometric authentication failed',
          canRetry: true,
        );
      }
    } catch (e) {
      state = EnhancedAuthState.error(
        'Biometric authentication failed: ${e.toString()}',
        canRetry: true,
      );
    }
  }

  /// Get user sessions for session management UI
  Future<List<SessionInfo>> getUserSessions() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      return await authService.getUserSessions();
    } catch (e) {
      throw Exception('Failed to get sessions: ${e.toString()}');
    }
  }

  /// Revoke a specific session
  Future<void> revokeSession(String sessionId) async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      await authService.revokeSession(sessionId);
    } catch (e) {
      throw Exception('Failed to revoke session: ${e.toString()}');
    }
  }

  /// Sign out from all devices
  Future<void> signOutAllDevices() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      await authService.signOutAllDevices();
      state = const EnhancedAuthState.unauthenticated();
    } catch (e) {
      throw Exception('Failed to sign out all devices: ${e.toString()}');
    }
  }

  /// Enhanced logout
  Future<void> enhancedLogout({bool revokeSession = true}) async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      await authService.enhancedLogout(revokeSession: revokeSession);
      state = const EnhancedAuthState.unauthenticated();
    } catch (e) {
      // Even if logout fails on server, clear local state
      state = const EnhancedAuthState.unauthenticated();
    }
  }

  /// Check if current session needs verification
  bool get requiresVerification {
    return state.maybeWhen(
      authenticated: (_, __, ___, ____, _____, requiresVerification) => 
          requiresVerification ?? false,
      orElse: () => false,
    );
  }

  /// Check if device is trusted
  bool get isDeviceTrusted {
    return state.maybeWhen(
      authenticated: (_, __, ___, deviceTrusted, ____, _____) => deviceTrusted,
      requiresTwoFactor: (_, __, deviceTrusted) => deviceTrusted,
      orElse: () => false,
    );
  }

  /// Get current user
  User? get currentUser {
    return state.maybeWhen(
      authenticated: (user, _, __, ___, ____, _____) => user,
      orElse: () => null,
    );
  }

  /// Check if user is authenticated
  bool get isAuthenticated {
    return state.maybeWhen(
      authenticated: (_, __, ___, ____, _____, ______) => true,
      orElse: () => false,
    );
  }

  /// Clear error state
  void clearError() {
    state.maybeWhen(
      error: (_, __) => state = const EnhancedAuthState.unauthenticated(),
      orElse: () {},
    );
  }

  /// Handle token expiration with automatic refresh
  Future<bool> handleTokenExpiration() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final refreshed = await authService.handleTokenExpiration();
      
      if (refreshed) {
        // Session was refreshed, restore authentication state
        await _restoreSession();
        return true;
      } else {
        // Refresh failed, logout user
        state = const EnhancedAuthState.unauthenticated();
        return false;
      }
    } catch (e) {
      state = const EnhancedAuthState.unauthenticated();
      return false;
    }
  }

  /// Set biometric authentication enabled/disabled
  Future<bool> setBiometricEnabled(bool enabled) async {
    try {
      final user = currentUser;
      if (user == null) return false;
      
      final authService = ref.read(enhancedAuthServiceProvider);
      return await authService.setBiometricEnabled(user.email, enabled);
    } catch (e) {
      return false;
    }
  }

  /// Check if biometric authentication is available
  Future<bool> canUseBiometrics() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      return await authService.canUseBiometrics();
    } catch (e) {
      return false;
    }
  }

  /// Get available biometric types
  Future<List<String>> getAvailableBiometrics() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final biometrics = await authService.getAvailableBiometrics();
      return biometrics.map((b) => b.name).toList();
    } catch (e) {
      return [];
    }
  }
}

/// Enhanced authentication state
sealed class EnhancedAuthState {
  const EnhancedAuthState();

  // Initial state (checking for existing session)
  const factory EnhancedAuthState.initial() = _Initial;
  
  // Loading state (authentication in progress)
  const factory EnhancedAuthState.loading() = _Loading;
  
  // Unauthenticated state
  const factory EnhancedAuthState.unauthenticated() = _Unauthenticated;
  
  // Authenticated state with comprehensive information
  const factory EnhancedAuthState.authenticated({
    required User user,
    required String accessToken,
    required String sessionId,
    required bool deviceTrusted,
    bool isRestoredSession,
    bool? requiresVerification,
  }) = _Authenticated;
  
  // Requires two-factor authentication
  const factory EnhancedAuthState.requiresTwoFactor({
    required String email,
    required User tempUser,
    required bool deviceTrusted,
  }) = _RequiresTwoFactor;
  
  // Requires CAPTCHA verification
  const factory EnhancedAuthState.requiresCaptcha({
    required String email,
    required int remainingAttempts,
    LockoutInfo? lockoutInfo,
  }) = _RequiresCaptcha;
  
  // Account locked out due to brute force protection
  const factory EnhancedAuthState.lockedOut({
    required LockoutInfo lockoutInfo,
    required String message,
  }) = _LockedOut;
  
  // Error state with retry capability
  const factory EnhancedAuthState.error(
    String message, {
    bool canRetry,
  }) = _Error;
  
  // Email verification required
  const factory EnhancedAuthState.requiresVerification({
    required User user,
    required String email,
  }) = _RequiresVerification;
}

// State implementations
class _Initial extends EnhancedAuthState {
  const _Initial();
}

class _Loading extends EnhancedAuthState {
  const _Loading();
}

class _Unauthenticated extends EnhancedAuthState {
  const _Unauthenticated();
}

class _Authenticated extends EnhancedAuthState {
  final User user;
  final String accessToken;
  final String sessionId;
  final bool deviceTrusted;
  final bool isRestoredSession;
  final bool? requiresVerification;

  const _Authenticated({
    required this.user,
    required this.accessToken,
    required this.sessionId,
    required this.deviceTrusted,
    this.isRestoredSession = false,
    this.requiresVerification,
  });
}

class _RequiresTwoFactor extends EnhancedAuthState {
  final String email;
  final User tempUser;
  final bool deviceTrusted;

  const _RequiresTwoFactor({
    required this.email,
    required this.tempUser,
    required this.deviceTrusted,
  });
}

class _RequiresCaptcha extends EnhancedAuthState {
  final String email;
  final int remainingAttempts;
  final LockoutInfo? lockoutInfo;

  const _RequiresCaptcha({
    required this.email,
    required this.remainingAttempts,
    this.lockoutInfo,
  });
}

class _LockedOut extends EnhancedAuthState {
  final LockoutInfo lockoutInfo;
  final String message;

  const _LockedOut({
    required this.lockoutInfo,
    required this.message,
  });
}

class _Error extends EnhancedAuthState {
  final String message;
  final bool canRetry;

  const _Error(this.message, {this.canRetry = false});
}

class _RequiresVerification extends EnhancedAuthState {
  final User user;
  final String email;

  const _RequiresVerification({
    required this.user,
    required this.email,
  });
}

/// Extension for pattern matching on auth state
extension EnhancedAuthStateExtensions on EnhancedAuthState {
  T when<T>({
    required T Function() initial,
    required T Function() loading,
    required T Function() unauthenticated,
    required T Function(User user, String accessToken, String sessionId, 
        bool deviceTrusted, bool isRestoredSession, bool? requiresVerification) authenticated,
    required T Function(String email, User tempUser, bool deviceTrusted) requiresTwoFactor,
    required T Function(String email, int remainingAttempts, LockoutInfo? lockoutInfo) requiresCaptcha,
    required T Function(LockoutInfo lockoutInfo, String message) lockedOut,
    required T Function(String message, bool canRetry) error,
    required T Function(User user, String email) requiresVerification,
  }) {
    switch (this) {
      case _Initial():
        return initial();
      case _Loading():
        return loading();
      case _Unauthenticated():
        return unauthenticated();
      case _Authenticated(:final user, :final accessToken, :final sessionId, 
            :final deviceTrusted, :final isRestoredSession, :final requiresVerification):
        return authenticated(user, accessToken, sessionId, deviceTrusted, 
            isRestoredSession, requiresVerification);
      case _RequiresTwoFactor(:final email, :final tempUser, :final deviceTrusted):
        return requiresTwoFactor(email, tempUser, deviceTrusted);
      case _RequiresCaptcha(:final email, :final remainingAttempts, :final lockoutInfo):
        return requiresCaptcha(email, remainingAttempts, lockoutInfo);
      case _LockedOut(:final lockoutInfo, :final message):
        return lockedOut(lockoutInfo, message);
      case _Error(:final message, :final canRetry):
        return error(message, canRetry);
      case _RequiresVerification(:final user, :final email):
        return requiresVerification(user, email);
    }
  }

  T? maybeWhen<T>({
    T Function()? initial,
    T Function()? loading,
    T Function()? unauthenticated,
    T Function(User user, String accessToken, String sessionId, 
        bool deviceTrusted, bool isRestoredSession, bool? requiresVerification)? authenticated,
    T Function(String email, User tempUser, bool deviceTrusted)? requiresTwoFactor,
    T Function(String email, int remainingAttempts, LockoutInfo? lockoutInfo)? requiresCaptcha,
    T Function(LockoutInfo lockoutInfo, String message)? lockedOut,
    T Function(String message, bool canRetry)? error,
    T Function(User user, String email)? requiresVerification,
    required T Function() orElse,
  }) {
    switch (this) {
      case _Initial():
        return initial?.call() ?? orElse();
      case _Loading():
        return loading?.call() ?? orElse();
      case _Unauthenticated():
        return unauthenticated?.call() ?? orElse();
      case _Authenticated(:final user, :final accessToken, :final sessionId, 
            :final deviceTrusted, :final isRestoredSession, :final requiresVerification):
        return authenticated?.call(user, accessToken, sessionId, deviceTrusted, 
            isRestoredSession, requiresVerification) ?? orElse();
      case _RequiresTwoFactor(:final email, :final tempUser, :final deviceTrusted):
        return requiresTwoFactor?.call(email, tempUser, deviceTrusted) ?? orElse();
      case _RequiresCaptcha(:final email, :final remainingAttempts, :final lockoutInfo):
        return requiresCaptcha?.call(email, remainingAttempts, lockoutInfo) ?? orElse();
      case _LockedOut(:final lockoutInfo, :final message):
        return lockedOut?.call(lockoutInfo, message) ?? orElse();
      case _Error(:final message, :final canRetry):
        return error?.call(message, canRetry) ?? orElse();
      case _RequiresVerification(:final user, :final email):
        return requiresVerification?.call(user, email) ?? orElse();
    }
  }
}
