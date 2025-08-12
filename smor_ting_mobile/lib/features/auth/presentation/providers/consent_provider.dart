import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/models/consent.dart';
import '../../../../core/services/consent_service.dart';

part 'consent_provider.g.dart';

/// State for consent management
sealed class ConsentState {
  const ConsentState();
}

class ConsentLoading extends ConsentState {
  const ConsentLoading();
}

class ConsentLoaded extends ConsentState {
  final List<ConsentRequirement> requirements;
  final UserConsent? userConsent;

  const ConsentLoaded(this.requirements, this.userConsent);
}

class ConsentError extends ConsentState {
  final String message;

  const ConsentError(this.message);
}

/// Provider for managing consent state
@riverpod
class ConsentNotifier extends _$ConsentNotifier {
  @override
  ConsentState build() {
    return const ConsentLoading();
  }

  /// Load consent requirements
  Future<void> loadConsentRequirements() async {
    state = const ConsentLoading();
    
    try {
      final consentService = ref.read(consentServiceProvider);
      final requirements = await consentService.getConsentRequirements();
      state = ConsentLoaded(requirements, null);
    } catch (e) {
      state = ConsentError('Failed to load consent requirements: $e');
    }
  }

  /// Load user consent status
  Future<void> loadUserConsent(String userId) async {
    try {
      final consentService = ref.read(consentServiceProvider);
      final requirements = await consentService.getConsentRequirements();
      final userConsent = await consentService.getUserConsent(userId);
      state = ConsentLoaded(requirements, userConsent);
    } catch (e) {
      state = ConsentError('Failed to load user consent: $e');
    }
  }

  /// Update multiple consents
  Future<void> updateConsents(
    String userId,
    Map<ConsentType, bool> consents,
  ) async {
    try {
      final consentService = ref.read(consentServiceProvider);
      await consentService.updateMultipleConsents(
        userId,
        consents,
        userAgent: 'Smor-Ting Mobile App',
        metadata: {
          'submitted_via': 'auth_flow',
          'timestamp': DateTime.now().toIso8601String(),
        },
      );
      
      // Reload user consent after update
      await loadUserConsent(userId);
    } catch (e) {
      state = ConsentError('Failed to update consent: $e');
    }
  }

  /// Check if user has required consents
  Future<bool> hasRequiredConsents(String userId) async {
    try {
      final consentService = ref.read(consentServiceProvider);
      return await consentService.hasRequiredConsents(userId);
    } catch (e) {
      return false;
    }
  }

  /// Get missing required consents
  Future<List<ConsentRequirement>> getMissingConsents(String userId) async {
    try {
      final consentService = ref.read(consentServiceProvider);
      return await consentService.getMissingConsents(userId);
    } catch (e) {
      return [];
    }
  }
}
