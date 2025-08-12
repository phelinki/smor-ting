import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:local_auth/local_auth.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../../core/services/enhanced_auth_service.dart';
import '../../../../core/models/user.dart';

part 'biometric_auth_provider.g.dart';

// Biometric authentication states
abstract class BiometricAuthState {
  const BiometricAuthState();
}

class BiometricAuthInitial extends BiometricAuthState {
  const BiometricAuthInitial();
}

class BiometricAuthLoading extends BiometricAuthState {
  const BiometricAuthLoading();
}

class BiometricAuthSuccess extends BiometricAuthState {
  final User user;
  final String accessToken;
  final String sessionId;
  final bool deviceTrusted;

  const BiometricAuthSuccess({
    required this.user,
    required this.accessToken,
    required this.sessionId,
    required this.deviceTrusted,
  });
}

class BiometricAuthError extends BiometricAuthState {
  final String message;
  final bool canRetry;

  const BiometricAuthError(this.message, {this.canRetry = true});
}

@riverpod
class BiometricAuth extends _$BiometricAuth {
  @override
  BiometricAuthState build() {
    return const BiometricAuthInitial();
  }

  /// Check if the device supports biometric authentication
  Future<bool> canUseBiometrics() async {
    final authService = ref.read(enhancedAuthServiceProvider);
    return await authService.canUseBiometrics();
  }

  /// Get available biometric types on the device
  Future<List<BiometricType>> getAvailableBiometrics() async {
    final authService = ref.read(enhancedAuthServiceProvider);
    return await authService.getAvailableBiometrics();
  }

  /// Check if biometric authentication is enabled for the user
  Future<bool> isBiometricEnabled(String email) async {
    final authService = ref.read(enhancedAuthServiceProvider);
    return await authService.isBiometricEnabled(email);
  }

  /// Enable biometric authentication for the user
  Future<bool> enableBiometric(String email) async {
    final authService = ref.read(enhancedAuthServiceProvider);
    return await authService.setBiometricEnabled(email, true);
  }

  /// Disable biometric authentication for the user
  Future<bool> disableBiometric(String email) async {
    final authService = ref.read(enhancedAuthServiceProvider);
    return await authService.setBiometricEnabled(email, false);
  }

  /// Authenticate using biometric authentication
  Future<void> authenticateWithBiometrics(String email) async {
    state = const BiometricAuthLoading();
    
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final result = await authService.authenticateWithBiometrics(email);
      
      if (result.success && result.user != null) {
        state = BiometricAuthSuccess(
          user: result.user!,
          accessToken: result.accessToken!,
          sessionId: result.sessionId!,
          deviceTrusted: result.deviceTrusted,
        );
      } else {
        state = BiometricAuthError(
          result.message ?? 'Biometric authentication failed',
          canRetry: true,
        );
      }
    } catch (e) {
      state = BiometricAuthError(
        'Biometric authentication failed: ${e.toString()}',
        canRetry: true,
      );
    }
  }

  /// Reset the authentication state
  void reset() {
    state = const BiometricAuthInitial();
  }
}

// Provider for quick access to biometric capabilities
@riverpod
Future<bool> biometricAvailable(BiometricAvailableRef ref) async {
  final authService = ref.watch(enhancedAuthServiceProvider);
  return await authService.canUseBiometrics();
}

@riverpod
Future<List<BiometricType>> availableBiometrics(AvailableBiometricsRef ref) async {
  final authService = ref.watch(enhancedAuthServiceProvider);
  return await authService.getAvailableBiometrics();
}
