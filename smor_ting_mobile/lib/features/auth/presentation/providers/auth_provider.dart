import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/models/user.dart';
import '../../../../core/services/api_service.dart';

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
    } catch (e) {
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
    if (state is Error) {
      state = const AuthState.initial();
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


