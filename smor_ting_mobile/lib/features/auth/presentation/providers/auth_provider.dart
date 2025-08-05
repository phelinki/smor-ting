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
    if (state is _Error) {
      state = const AuthState.initial();
    }
  }
}

// Auth State
sealed class AuthState {
  const AuthState();

  const factory AuthState.initial() = _Initial;
  const factory AuthState.loading() = _Loading;
  const factory AuthState.authenticated(User user, String accessToken) = _Authenticated;
  const factory AuthState.requiresOTP({required String email, required User user}) = _RequiresOTP;
  const factory AuthState.error(String message) = _Error;
}

class _Initial extends AuthState {
  const _Initial();
}

class _Loading extends AuthState {
  const _Loading();
}

class _Authenticated extends AuthState {
  final User user;
  final String accessToken;
  const _Authenticated(this.user, this.accessToken);
}

class _RequiresOTP extends AuthState {
  final String email;
  final User user;
  const _RequiresOTP({required this.email, required this.user});
}

class _Error extends AuthState {
  final String message;
  const _Error(this.message);
}

 