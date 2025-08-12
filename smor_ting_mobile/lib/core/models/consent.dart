// TODO: Temporarily simplified consent models for compilation
// import 'package:freezed_annotation/freezed_annotation.dart';

// part 'consent.freezed.dart';
// part 'consent.g.dart';

/// Type of consent required from user
enum ConsentType {
  termsOfService,
  privacyPolicy,
  marketingCommunications,
  dataProcessing,
  biometricData,
  locationTracking,
  analytics,
}

/// Consent record for audit trail
class ConsentRecord {
  final String id;
  final ConsentType type;
  final bool granted;
  final DateTime consentedAt;
  final String version;
  final String? userAgent;
  final String? ipAddress;
  final Map<String, dynamic>? metadata;

  const ConsentRecord({
    required this.id,
    required this.type,
    required this.granted,
    required this.consentedAt,
    required this.version,
    this.userAgent,
    this.ipAddress,
    this.metadata,
  });

  factory ConsentRecord.fromJson(Map<String, dynamic> json) => ConsentRecord(
    id: json['id'] as String,
    type: ConsentType.values.firstWhere((e) => e.name == json['type']),
    granted: json['granted'] as bool,
    consentedAt: DateTime.parse(json['consented_at'] as String),
    version: json['version'] as String,
    userAgent: json['user_agent'] as String?,
    ipAddress: json['ip_address'] as String?,
    metadata: json['metadata'] as Map<String, dynamic>?,
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'type': type.name,
    'granted': granted,
    'consented_at': consentedAt.toIso8601String(),
    'version': version,
    'user_agent': userAgent,
    'ip_address': ipAddress,
    'metadata': metadata,
  };
}

/// User consent status
class UserConsent {
  final String userId;
  final Map<ConsentType, ConsentRecord> consents;
  final DateTime lastUpdated;

  const UserConsent({
    required this.userId,
    required this.consents,
    required this.lastUpdated,
  });

  factory UserConsent.fromJson(Map<String, dynamic> json) => UserConsent(
    userId: json['user_id'] as String,
    consents: (json['consents'] as Map<String, dynamic>).map(
      (k, v) => MapEntry(
        ConsentType.values.firstWhere((e) => e.name == k),
        ConsentRecord.fromJson(v as Map<String, dynamic>),
      ),
    ),
    lastUpdated: DateTime.parse(json['last_updated'] as String),
  );

  Map<String, dynamic> toJson() => {
    'user_id': userId,
    'consents': consents.map((k, v) => MapEntry(k.name, v.toJson())),
    'last_updated': lastUpdated.toIso8601String(),
  };
}

/// Consent requirement definition
class ConsentRequirement {
  final ConsentType type;
  final String title;
  final String description;
  final String version;
  final bool required;
  final String? documentUrl;

  const ConsentRequirement({
    required this.type,
    required this.title,
    required this.description,
    required this.version,
    required this.required,
    this.documentUrl,
  });

  factory ConsentRequirement.fromJson(Map<String, dynamic> json) => ConsentRequirement(
    type: ConsentType.values.firstWhere((e) => e.name == json['type']),
    title: json['title'] as String,
    description: json['description'] as String,
    version: json['version'] as String,
    required: json['required'] as bool,
    documentUrl: json['document_url'] as String?,
  );

  Map<String, dynamic> toJson() => {
    'type': type.name,
    'title': title,
    'description': description,
    'version': version,
    'required': required,
    'document_url': documentUrl,
  };
}

/// Request to update user consent
class ConsentUpdateRequest {
  final ConsentType type;
  final bool granted;
  final String version;
  final String? userAgent;
  final String? ipAddress;
  final Map<String, dynamic>? metadata;

  const ConsentUpdateRequest({
    required this.type,
    required this.granted,
    required this.version,
    this.userAgent,
    this.ipAddress,
    this.metadata,
  });

  factory ConsentUpdateRequest.fromJson(Map<String, dynamic> json) => ConsentUpdateRequest(
    type: ConsentType.values.firstWhere((e) => e.name == json['type']),
    granted: json['granted'] as bool,
    version: json['version'] as String,
    userAgent: json['user_agent'] as String?,
    ipAddress: json['ip_address'] as String?,
    metadata: json['metadata'] as Map<String, dynamic>?,
  );

  Map<String, dynamic> toJson() => {
    'type': type.name,
    'granted': granted,
    'version': version,
    'user_agent': userAgent,
    'ip_address': ipAddress,
    'metadata': metadata,
  };
}
