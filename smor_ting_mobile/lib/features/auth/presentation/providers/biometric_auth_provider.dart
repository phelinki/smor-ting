import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:local_auth/local_auth.dart';
import '../../../../core/services/enhanced_auth_service.dart';
import '../../../../core/models/enhanced_auth_models.dart' as models;

// States
abstract class BiometricAuthState {}
class BiometricAuthInitial extends BiometricAuthState {}
class BiometricAuthSuccess extends BiometricAuthState {
  final models.EnhancedAuthResult result;
  BiometricAuthSuccess(this.result);
  String get accessToken => result.accessToken ?? '';
  // expose user but avoid type import errors by using dynamic map through result
  get user => result.user!;
}
class BiometricAuthError extends BiometricAuthState {
  final String message;
  BiometricAuthError(this.message);
}

class BiometricAuthNotifier extends StateNotifier<BiometricAuthState> {
  final EnhancedAuthService _authService;
  final LocalAuthentication _localAuth;
  BiometricAuthNotifier(this._authService, this._localAuth) : super(BiometricAuthInitial());

  Future<bool> canUseBiometrics() => _authService.canUseBiometrics();
  Future<List<BiometricType>> getAvailableBiometrics() => _authService.getAvailableBiometrics();
  Future<bool> isBiometricEnabled(String email) => _authService.isBiometricEnabled(email);
  Future<bool> enableBiometric(String email) => _authService.setBiometricEnabled(email, true);
  Future<bool> disableBiometric(String email) => _authService.setBiometricEnabled(email, false);

  Future<void> authenticateWithBiometrics(String email) async {
    try {
      final models.EnhancedAuthResult result = await _authService.authenticateWithBiometrics(email);
      if (result.success) {
        state = BiometricAuthSuccess(result);
      } else {
        state = BiometricAuthError(result.message ?? 'Biometric authentication failed');
      }
    } catch (e) {
      state = BiometricAuthError(e.toString());
    }
  }
}

final biometricAuthProvider = StateNotifierProvider<BiometricAuthNotifier, BiometricAuthState>((ref) {
  final auth = ref.read(enhancedAuthServiceProvider);
  final local = LocalAuthentication();
  return BiometricAuthNotifier(auth, local);
});


