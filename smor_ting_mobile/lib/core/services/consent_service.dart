import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:uuid/uuid.dart';
import '../models/consent.dart';
import 'api_service.dart';

part 'consent_service.g.dart';

/// Service for managing user consent and compliance
class ConsentService {
  final ApiService _apiService;
  final FlutterSecureStorage _secureStorage;
  final Uuid _uuid = const Uuid();

  static const String _consentKey = 'smor_ting_consent_v1';
  static const String _pendingConsentKey = 'smor_ting_pending_consent_v1';

  ConsentService(this._apiService, this._secureStorage);

  /// Get current consent requirements
  Future<List<ConsentRequirement>> getConsentRequirements() async {
    try {
      final response = await _apiService.get('/api/v1/consent/requirements');
      final requirementsList = response['requirements'] as List;
      return requirementsList
          .map((json) => ConsentRequirement.fromJson(json))
          .toList();
    } catch (e) {
      // Return default consent requirements if API fails
      return _getDefaultConsentRequirements();
    }
  }

  /// Get user's current consent status
  Future<UserConsent?> getUserConsent(String userId) async {
    try {
      // Try to get from API first
      final response = await _apiService.get('/api/v1/consent/user/$userId');
      return UserConsent.fromJson(response);
    } catch (e) {
      // Fallback to local storage
      return await _getLocalConsent(userId);
    }
  }

  /// Update user consent
  Future<void> updateConsent(
    String userId,
    ConsentType type,
    bool granted, {
    String? userAgent,
    String? ipAddress,
    Map<String, dynamic>? metadata,
  }) async {
    final requirements = await getConsentRequirements();
    final requirement = requirements.firstWhere(
      (r) => r.type == type,
      orElse: () => throw Exception('Consent requirement not found: $type'),
    );

    final request = ConsentUpdateRequest(
      type: type,
      granted: granted,
      version: requirement.version,
      userAgent: userAgent,
      ipAddress: ipAddress,
      metadata: metadata,
    );

    try {
      // Send to API
      await _apiService.post('/api/v1/consent/user/$userId', request.toJson());
      
      // Update local storage
      await _updateLocalConsent(userId, type, granted, requirement.version);
    } catch (e) {
      // Store as pending if API fails
      await _storePendingConsent(userId, request);
      rethrow;
    }
  }

  /// Batch update multiple consents
  Future<void> updateMultipleConsents(
    String userId,
    Map<ConsentType, bool> consents, {
    String? userAgent,
    String? ipAddress,
    Map<String, dynamic>? metadata,
  }) async {
    final updates = <ConsentUpdateRequest>[];
    final requirements = await getConsentRequirements();

    for (final entry in consents.entries) {
      final requirement = requirements.firstWhere(
        (r) => r.type == entry.key,
        orElse: () => throw Exception('Consent requirement not found: ${entry.key}'),
      );

      updates.add(ConsentUpdateRequest(
        type: entry.key,
        granted: entry.value,
        version: requirement.version,
        userAgent: userAgent,
        ipAddress: ipAddress,
        metadata: metadata,
      ));
    }

    try {
      // Send batch update to API
      await _apiService.post('/api/v1/consent/user/$userId/batch', {
        'updates': updates.map((u) => u.toJson()).toList(),
      });

      // Update local storage
      for (final update in updates) {
        await _updateLocalConsent(
          userId,
          update.type,
          update.granted,
          update.version,
        );
      }
    } catch (e) {
      // Store as pending if API fails
      for (final update in updates) {
        await _storePendingConsent(userId, update);
      }
      rethrow;
    }
  }

  /// Check if user has given required consents
  Future<bool> hasRequiredConsents(String userId) async {
    final requirements = await getConsentRequirements();
    final userConsent = await getUserConsent(userId);

    if (userConsent == null) return false;

    for (final requirement in requirements) {
      if (requirement.required) {
        final consent = userConsent.consents[requirement.type];
        if (consent == null || !consent.granted) {
          return false;
        }
      }
    }

    return true;
  }

  /// Get missing required consents
  Future<List<ConsentRequirement>> getMissingConsents(String userId) async {
    final requirements = await getConsentRequirements();
    final userConsent = await getUserConsent(userId);
    final missing = <ConsentRequirement>[];

    for (final requirement in requirements) {
      if (requirement.required) {
        final consent = userConsent?.consents[requirement.type];
        if (consent == null || !consent.granted) {
          missing.add(requirement);
        }
      }
    }

    return missing;
  }

  /// Sync pending consents when online
  Future<void> syncPendingConsents() async {
    try {
      final pendingJson = await _secureStorage.read(key: _pendingConsentKey);
      if (pendingJson == null) return;

      final pendingList = jsonDecode(pendingJson) as List;
      final pending = pendingList
          .map((json) => MapEntry<String, ConsentUpdateRequest>(
                json['userId'] as String,
                ConsentUpdateRequest.fromJson(json['request']),
              ))
          .toList();

      for (final entry in pending) {
        try {
          await _apiService.post(
            '/api/v1/consent/user/${entry.key}',
            entry.value.toJson(),
          );
        } catch (e) {
          // Skip failed syncs for now
        }
      }

      // Clear pending consents after sync
      await _secureStorage.delete(key: _pendingConsentKey);
    } catch (e) {
      // Ignore sync errors
    }
  }

  /// Check if consent is expired and needs refresh
  Future<bool> isConsentExpired(String userId, ConsentType type) async {
    final requirements = await getConsentRequirements();
    final userConsent = await getUserConsent(userId);

    if (userConsent == null) return true;

    final requirement = requirements.firstWhere(
      (r) => r.type == type,
      orElse: () => throw Exception('Consent requirement not found: $type'),
    );

    final consent = userConsent.consents[type];
    if (consent == null) return true;

    // Check if version has changed
    return consent.version != requirement.version;
  }

  /// Private methods

  Future<UserConsent?> _getLocalConsent(String userId) async {
    try {
      final consentJson = await _secureStorage.read(key: '${_consentKey}_$userId');
      if (consentJson == null) return null;

      return UserConsent.fromJson(jsonDecode(consentJson));
    } catch (e) {
      return null;
    }
  }

  Future<void> _updateLocalConsent(
    String userId,
    ConsentType type,
    bool granted,
    String version,
  ) async {
    final currentConsent = await _getLocalConsent(userId);
    final now = DateTime.now();

    final record = ConsentRecord(
      id: _uuid.v4(),
      type: type,
      granted: granted,
      consentedAt: now,
      version: version,
    );

    final updatedConsents = Map<ConsentType, ConsentRecord>.from(
      currentConsent?.consents ?? {},
    );
    updatedConsents[type] = record;

    final newConsent = UserConsent(
      userId: userId,
      consents: updatedConsents,
      lastUpdated: now,
    );

    await _secureStorage.write(
      key: '${_consentKey}_$userId',
      value: jsonEncode(newConsent.toJson()),
    );
  }

  Future<void> _storePendingConsent(String userId, ConsentUpdateRequest request) async {
    try {
      final pendingJson = await _secureStorage.read(key: _pendingConsentKey);
      final pendingList = pendingJson != null ? jsonDecode(pendingJson) as List : <Map<String, dynamic>>[];

      pendingList.add({
        'userId': userId,
        'request': request.toJson(),
        'timestamp': DateTime.now().toIso8601String(),
      });

      await _secureStorage.write(
        key: _pendingConsentKey,
        value: jsonEncode(pendingList),
      );
    } catch (e) {
      // Ignore storage errors for pending consents
    }
  }

  List<ConsentRequirement> _getDefaultConsentRequirements() {
    return [
      const ConsentRequirement(
        type: ConsentType.termsOfService,
        title: 'Terms of Service',
        description: 'By using Smor-Ting, you agree to our Terms of Service',
        version: '1.0',
        required: true,
        documentUrl: 'https://smor-ting.com/terms',
      ),
      const ConsentRequirement(
        type: ConsentType.privacyPolicy,
        title: 'Privacy Policy',
        description: 'We respect your privacy and handle your data according to our Privacy Policy',
        version: '1.0',
        required: true,
        documentUrl: 'https://smor-ting.com/privacy',
      ),
      const ConsentRequirement(
        type: ConsentType.dataProcessing,
        title: 'Data Processing',
        description: 'Allow processing of your personal data for service delivery',
        version: '1.0',
        required: true,
      ),
      const ConsentRequirement(
        type: ConsentType.marketingCommunications,
        title: 'Marketing Communications',
        description: 'Receive promotional emails and notifications about new services',
        version: '1.0',
        required: false,
      ),
      const ConsentRequirement(
        type: ConsentType.analytics,
        title: 'Analytics',
        description: 'Help us improve the app by allowing anonymous usage analytics',
        version: '1.0',
        required: false,
      ),
    ];
  }
}

/// Riverpod provider for consent service
@riverpod
ConsentService consentService(ConsentServiceRef ref) {
  return ConsentService(
    ref.read(apiServiceProvider),
    const FlutterSecureStorage(
      aOptions: AndroidOptions(
        encryptedSharedPreferences: true,
      ),
      iOptions: IOSOptions(
        accessibility: KeychainAccessibility.first_unlock_this_device,
      ),
    ),
  );
}
