import 'package:freezed_annotation/freezed_annotation.dart';

part 'consent.freezed.dart';
part 'consent.g.dart';

/// Type of consent required from user
enum ConsentType {
  @JsonValue('terms_of_service')
  termsOfService,
  @JsonValue('privacy_policy') 
  privacyPolicy,
  @JsonValue('marketing_communications')
  marketingCommunications,
  @JsonValue('data_processing')
  dataProcessing,
  @JsonValue('biometric_data')
  biometricData,
  @JsonValue('location_tracking')
  locationTracking,
  @JsonValue('analytics')
  analytics,
}

/// Consent record for audit trail
@freezed
class ConsentRecord with _$ConsentRecord {
  const factory ConsentRecord({
    required String id,
    required ConsentType type,
    required bool granted,
    required DateTime consentedAt,
    required String version,
    String? userAgent,
    String? ipAddress,
    Map<String, dynamic>? metadata,
  }) = _ConsentRecord;

  factory ConsentRecord.fromJson(Map<String, dynamic> json) =>
      _$ConsentRecordFromJson(json);
}

/// User consent status
@freezed
class UserConsent with _$UserConsent {
  const factory UserConsent({
    required String userId,
    required Map<ConsentType, ConsentRecord> consents,
    required DateTime lastUpdated,
  }) = _UserConsent;

  factory UserConsent.fromJson(Map<String, dynamic> json) =>
      _$UserConsentFromJson(json);
}

/// Consent requirement definition
@freezed
class ConsentRequirement with _$ConsentRequirement {
  const factory ConsentRequirement({
    required ConsentType type,
    required String title,
    required String description,
    required String version,
    required bool required,
    String? documentUrl,
  }) = _ConsentRequirement;

  factory ConsentRequirement.fromJson(Map<String, dynamic> json) =>
      _$ConsentRequirementFromJson(json);
}

/// Request to update user consent
@freezed
class ConsentUpdateRequest with _$ConsentUpdateRequest {
  const factory ConsentUpdateRequest({
    required ConsentType type,
    required bool granted,
    required String version,
    String? userAgent,
    String? ipAddress,
    Map<String, dynamic>? metadata,
  }) = _ConsentUpdateRequest;

  factory ConsentUpdateRequest.fromJson(Map<String, dynamic> json) =>
      _$ConsentUpdateRequestFromJson(json);
}
