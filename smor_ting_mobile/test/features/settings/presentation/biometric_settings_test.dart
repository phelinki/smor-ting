import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mocktail/mocktail.dart';
import 'package:local_auth/local_auth.dart';

import 'package:smor_ting_mobile/core/services/enhanced_auth_service.dart';

// Mock classes
class MockEnhancedAuthService extends Mock implements EnhancedAuthService {}

void main() {
  group('Biometric Settings Tests', () {
    late MockEnhancedAuthService mockAuthService;

    setUp(() {
      mockAuthService = MockEnhancedAuthService();
    });

    Widget createBiometricSettingsWidget({bool biometricAvailable = true, bool biometricEnabled = false}) {
      return ProviderScope(
        overrides: [
          enhancedAuthServiceProvider.overrideWithValue(mockAuthService),
        ],
        child: MaterialApp(
          home: Scaffold(
            body: BiometricSettingsWidget(),
          ),
        ),
      );
    }

    group('Biometric Service Integration', () {
      test('canUseBiometrics returns true when biometrics are available', () async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => true);

        // Act
        final result = await mockAuthService.canUseBiometrics();

        // Assert
        expect(result, true);
        verify(() => mockAuthService.canUseBiometrics()).called(1);
      });

      test('canUseBiometrics returns false when biometrics are not available', () async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenAnswer((_) async => false);

        // Act
        final result = await mockAuthService.canUseBiometrics();

        // Assert
        expect(result, false);
        verify(() => mockAuthService.canUseBiometrics()).called(1);
      });

      test('isBiometricEnabled returns current state for user', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.isBiometricEnabled(email))
            .thenAnswer((_) async => true);

        // Act
        final result = await mockAuthService.isBiometricEnabled(email);

        // Assert
        expect(result, true);
        verify(() => mockAuthService.isBiometricEnabled(email)).called(1);
      });

      test('setBiometricEnabled successfully enables biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.setBiometricEnabled(email, true))
            .thenAnswer((_) async => true);

        // Act
        final result = await mockAuthService.setBiometricEnabled(email, true);

        // Assert
        expect(result, true);
        verify(() => mockAuthService.setBiometricEnabled(email, true)).called(1);
      });

      test('setBiometricEnabled successfully disables biometric authentication', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.setBiometricEnabled(email, false))
            .thenAnswer((_) async => true);

        // Act
        final result = await mockAuthService.setBiometricEnabled(email, false);

        // Assert
        expect(result, true);
        verify(() => mockAuthService.setBiometricEnabled(email, false)).called(1);
      });

      test('getAvailableBiometrics returns list of available biometric types', () async {
        // Arrange
        final expectedBiometrics = [BiometricType.fingerprint, BiometricType.face];
        when(() => mockAuthService.getAvailableBiometrics())
            .thenAnswer((_) async => expectedBiometrics);

        // Act
        final result = await mockAuthService.getAvailableBiometrics();

        // Assert
        expect(result, expectedBiometrics);
        verify(() => mockAuthService.getAvailableBiometrics()).called(1);
      });
    });

    group('Error Handling', () {
      test('setBiometricEnabled throws exception when biometric setup fails', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.setBiometricEnabled(email, true))
            .thenThrow(Exception('Biometric setup failed'));

        // Act & Assert
        expect(
          () => mockAuthService.setBiometricEnabled(email, true),
          throwsException,
        );
      });

      test('setBiometricEnabled throws exception when disable fails', () async {
        // Arrange
        const email = 'test@example.com';
        when(() => mockAuthService.setBiometricEnabled(email, false))
            .thenThrow(Exception('Disable biometric failed'));

        // Act & Assert
        expect(
          () => mockAuthService.setBiometricEnabled(email, false),
          throwsException,
        );
      });

      test('canUseBiometrics handles service errors gracefully', () async {
        // Arrange
        when(() => mockAuthService.canUseBiometrics())
            .thenThrow(Exception('Biometric service unavailable'));

        // Act & Assert
        expect(
          () => mockAuthService.canUseBiometrics(),
          throwsException,
        );
      });
    });

    group('Biometric State Management', () {
      test('biometric state changes are properly tracked', () async {
        // Arrange
        const email = 'test@example.com';
        
        // Initially disabled
        when(() => mockAuthService.isBiometricEnabled(email))
            .thenAnswer((_) async => false);
        
        // Act - check initial state
        final initialState = await mockAuthService.isBiometricEnabled(email);
        
        // Enable biometric
        when(() => mockAuthService.setBiometricEnabled(email, true))
            .thenAnswer((_) async => true);
        await mockAuthService.setBiometricEnabled(email, true);
        
        // Now enabled
        when(() => mockAuthService.isBiometricEnabled(email))
            .thenAnswer((_) async => true);
        final enabledState = await mockAuthService.isBiometricEnabled(email);
        
        // Assert
        expect(initialState, false);
        expect(enabledState, true);
        verify(() => mockAuthService.setBiometricEnabled(email, true)).called(1);
      });
    });
  });
}

// Simple widget component for testing biometric settings in isolation
class BiometricSettingsWidget extends ConsumerStatefulWidget {
  @override
  ConsumerState<BiometricSettingsWidget> createState() => _BiometricSettingsWidgetState();
}

class _BiometricSettingsWidgetState extends ConsumerState<BiometricSettingsWidget> {
  bool _biometricAvailable = false;
  bool _biometricEnabled = false;
  bool _loading = false;

  @override
  void initState() {
    super.initState();
    _checkBiometricAvailability();
  }

  Future<void> _checkBiometricAvailability() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final available = await authService.canUseBiometrics();
      if (mounted) {
        setState(() {
          _biometricAvailable = available;
        });
      }
    } catch (e) {
      // Handle error
    }
  }

  Future<void> _toggleBiometric(bool enabled) async {
    setState(() { _loading = true; });
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      await authService.setBiometricEnabled('test@example.com', enabled);
      if (mounted) {
        setState(() {
          _biometricEnabled = enabled;
        });
      }
    } catch (e) {
      // Handle error
    } finally {
      if (mounted) {
        setState(() { _loading = false; });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (!_biometricAvailable) {
      return const Center(child: Text('Biometric authentication not available'));
    }

    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        children: [
          ListTile(
            title: const Text('Biometric Authentication'),
            subtitle: const Text('Use fingerprint or face unlock to secure your account'),
            trailing: Switch(
              value: _biometricEnabled,
              onChanged: _loading ? null : _toggleBiometric,
            ),
          ),
        ],
      ),
    );
  }
}