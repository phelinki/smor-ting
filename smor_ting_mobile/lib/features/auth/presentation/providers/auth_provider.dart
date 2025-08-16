import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/models/user.dart';
import '../../../../core/services/api_service.dart';
import '../../../../core/exceptions/auth_exceptions.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

part 'auth_provider.g.dart';

@riverpod
class AuthNotifier extends _$AuthNotifier {
  final _secureStorage = const FlutterSecureStorage();
  
  @override
  AuthState build() {
    // Start with initial state, not loading
    // The initialization will happen in the splash page
    return const AuthState.initial();
  }

  /// Initialize auth state by checking for stored tokens
  /// This should be called from the splash page
  Future<void> initializeAuthState() async {
    try {
      print('🔵 AuthProvider: Checking for stored authentication tokens...');
      
      // Check if we have stored tokens
      final accessToken = await _secureStorage.read(key: 'access_token');
      final refreshToken = await _secureStorage.read(key: 'refresh_token');
      final sessionId = await _secureStorage.read(key: 'session_id');
      
      if (accessToken == null || refreshToken == null) {
        print('🔵 AuthProvider: No stored tokens found, user needs to login');
        state = const AuthState.initial();
        return;
      }
      
      print('🔵 AuthProvider: Found stored tokens, validating...');
      
      // Try to get a valid token (this will refresh if needed)
      final apiService = ref.read(apiServiceProvider);
      final validToken = await apiService.authService.getValidToken();
      
      // Fetch user profile to restore full auth state
      final user = await apiService.getUserProfile();
      
      print('🔵 AuthProvider: Successfully restored authentication state for user: ${user.email}');
      state = AuthState.authenticated(user, validToken);
      
    } catch (e) {
      print('🔴 AuthProvider: Failed to restore auth state: $e');
      // Clear invalid tokens
      await _clearStoredTokens();
      state = const AuthState.initial();
    }
  }

  /// Clear all stored authentication tokens
  Future<void> _clearStoredTokens() async {
    await Future.wait([
      _secureStorage.delete(key: 'access_token'),
      _secureStorage.delete(key: 'refresh_token'),
      _secureStorage.delete(key: 'session_id'),
      _secureStorage.delete(key: 'token_expires_at'),
      _secureStorage.delete(key: 'refresh_expires_at'),
    ]);
  }

  Future<void> login(String email, String password) async {
    state = const AuthState.loading();
    
    try {
      final apiService = ref.read(apiServiceProvider);
      final request = LoginRequest(email: email, password: password);
      final response = await apiService.login(request);
      
      // Store tokens using auth service
      await apiService.authService.storeTokens({
        'access_token': response.accessToken!,
        'refresh_token': response.refreshToken!,
        'token_expires_at': DateTime.now().add(const Duration(hours: 1)).toIso8601String(),
        'refresh_expires_at': DateTime.now().add(const Duration(days: 7)).toIso8601String(),
        'session_id': '', // Will be set by backend if needed
      });
      
      // OTP is disabled - always go directly to authenticated state
      state = AuthState.authenticated(response.user, response.accessToken!);
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  Future<void> register({
    required String firstName,
    required String lastName,
    required String email,
    required String phone,
    required String password,
    UserRole role = UserRole.customer,
  }) async {
    state = const AuthState.loading();
    
    try {
      final apiService = ref.read(apiServiceProvider);
      final request = RegisterRequest(
        email: email,
        password: password,
        firstName: firstName,
        lastName: lastName,
        phone: phone,
        role: role,
      );
      
      final response = await apiService.register(request);
      
      // Store tokens using auth service
      await apiService.authService.storeTokens({
        'access_token': response.accessToken!,
        'refresh_token': response.refreshToken!,
        'token_expires_at': DateTime.now().add(const Duration(hours: 1)).toIso8601String(),
        'refresh_expires_at': DateTime.now().add(const Duration(days: 7)).toIso8601String(),
        'session_id': '', // Will be set by backend if needed
      });
      
      // OTP is disabled - always go directly to authenticated state
      state = AuthState.authenticated(response.user, response.accessToken!);
    } on EmailAlreadyExistsException catch (e) {
      print('🔴 AuthProvider: Caught EmailAlreadyExistsException with email: ${e.email}');
      state = AuthState.emailAlreadyExists(e.email);
    } catch (e) {
      print('🔴 AuthProvider: Caught general exception: $e');
      print('🔴 AuthProvider: Exception type: ${e.runtimeType}');
      state = AuthState.error(e.toString());
    }
  }

  void logout() async {
    try {
      // Clear stored tokens
      await _clearStoredTokens();
      
      state = const AuthState.initial();
    } catch (e) {
      print('🔴 AuthProvider: Error during logout: $e');
      // Still set to initial state even if clearing fails
      state = const AuthState.initial();
    }
  }

  void clearError() {
    if (state is Error || state is EmailAlreadyExists) {
      state = const AuthState.initial();
    }
  }

  void resetToInitial() {
    state = const AuthState.initial();
  }

  // Forgot password UX
  Future<void> requestPasswordReset(String email) async {
    try {
      final apiService = ref.read(apiServiceProvider);
      await apiService.requestPasswordReset(email);
      state = PasswordResetEmailSent(email);
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  Future<void> resetPassword(String email, String newPassword) async {
    try {
      final apiService = ref.read(apiServiceProvider);
      await apiService.resetPassword(email, newPassword);
      state = const PasswordResetSuccess();
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  void setAuthenticatedUser(User user, String accessToken) {
    state = AuthState.authenticated(user, accessToken);
  }
}

// Auth State
sealed class AuthState {
  const AuthState();

  const factory AuthState.initial() = Initial;
  const factory AuthState.loading() = Loading;
  const factory AuthState.authenticated(User user, String accessToken) = Authenticated;
  const factory AuthState.requiresOTP({required String email, required User user}) = RequiresOTP;
  const factory AuthState.error(String message) = Error;
  const factory AuthState.emailAlreadyExists(String email) = EmailAlreadyExists;
  const factory AuthState.passwordResetEmailSent(String email) = PasswordResetEmailSent;
  const factory AuthState.passwordResetSuccess() = PasswordResetSuccess;
}

class Initial extends AuthState {
  const Initial();
}

class Loading extends AuthState {
  const Loading();
}

class Authenticated extends AuthState {
  final User user;
  final String accessToken;
  const Authenticated(this.user, this.accessToken);
}

class RequiresOTP extends AuthState {
  final String email;
  final User user;
  const RequiresOTP({required this.email, required this.user});
}

class Error extends AuthState {
  final String message;
  const Error(this.message);
}

class EmailAlreadyExists extends AuthState {
  final String email;
  const EmailAlreadyExists(this.email);
}

class PasswordResetEmailSent extends AuthState {
  final String email;
  const PasswordResetEmailSent(this.email);
}

class PasswordResetSuccess extends AuthState {
  const PasswordResetSuccess();
}


