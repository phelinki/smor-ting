import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/models/user.dart';
import '../../../../core/services/api_service.dart';
import '../../../../core/exceptions/auth_exceptions.dart';

part 'auth_provider.g.dart';

@riverpod
class AuthNotifier extends _$AuthNotifier {
  @override
  AuthState build() {
    return const AuthState.initial();
  }

  Future<void> login(String email, String password) async {
    state = const AuthState.loading();
    
    try {
      final apiService = ref.read(apiServiceProvider);
      final request = LoginRequest(email: email, password: password);
      final response = await apiService.login(request);
      
      if (response.requiresOTP) {
        state = AuthState.requiresOTP(email: email, user: response.user);
      } else {
        apiService.setAuthToken(response.accessToken!);
        state = AuthState.authenticated(response.user, response.accessToken!);
      }
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
      
      if (response.requiresOTP) {
        state = AuthState.requiresOTP(email: email, user: response.user);
      } else {
        apiService.setAuthToken(response.accessToken!);
        state = AuthState.authenticated(response.user, response.accessToken!);
      }
    } on EmailAlreadyExistsException catch (e) {
      print('🔴 AuthProvider: Caught EmailAlreadyExistsException with email: ${e.email}');
      state = AuthState.emailAlreadyExists(e.email);
    } catch (e) {
      print('🔴 AuthProvider: Caught general exception: $e');
      print('🔴 AuthProvider: Exception type: ${e.runtimeType}');
      state = AuthState.error(e.toString());
    }
  }

  Future<void> verifyOTP(String email, String otp) async {
    state = const AuthState.loading();
    
    try {
      final apiService = ref.read(apiServiceProvider);
      final request = VerifyOTPRequest(email: email, otp: otp);
      final response = await apiService.verifyOTP(request);
      
      apiService.setAuthToken(response.accessToken!);
      state = AuthState.authenticated(response.user, response.accessToken!);
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  Future<void> resendOTP(String email) async {
    try {
      final apiService = ref.read(apiServiceProvider);
      await apiService.resendOTP(email);
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  void logout() {
    final apiService = ref.read(apiServiceProvider);
    apiService.clearAuthToken();
    state = const AuthState.initial();
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

  Future<void> resetPassword(String email, String otp, String newPassword) async {
    try {
      final apiService = ref.read(apiServiceProvider);
      await apiService.resetPassword(email, otp, newPassword);
      state = const PasswordResetSuccess();
    } catch (e) {
      state = AuthState.error(e.toString());
    }
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


