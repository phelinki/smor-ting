import 'package:json_annotation/json_annotation.dart';
import 'user.dart';
import '../services/device_fingerprint_service.dart' show DeviceFingerprint;

part 'enhanced_auth_models.g.dart';

@JsonSerializable()
class EnhancedLoginRequest {
  final String email;
  final String password;
  final bool rememberMe;
  final DeviceFingerprint deviceInfo;
  final String? captchaToken;
  final String? twoFactorCode;

  const EnhancedLoginRequest({
    required this.email,
    required this.password,
    this.rememberMe = false,
    required this.deviceInfo,
    this.captchaToken,
    this.twoFactorCode,
  });

  factory EnhancedLoginRequest.fromJson(Map<String, dynamic> json) =>
      _$EnhancedLoginRequestFromJson(json);

  Map<String, dynamic> toJson() => _$EnhancedLoginRequestToJson(this);
}

@JsonSerializable()
class EnhancedAuthResult {
  final bool success;
  final String? message;
  final User? user;
  final String? accessToken;
  final String? refreshToken;
  final String? sessionId;
  final DateTime? tokenExpiresAt;
  final DateTime? refreshExpiresAt;
  final bool requiresTwoFactor;
  final bool requiresVerification;
  final bool deviceTrusted;
  final bool isRestoredSession;

  const EnhancedAuthResult({
    required this.success,
    this.message,
    this.user,
    this.accessToken,
    this.refreshToken,
    this.sessionId,
    this.tokenExpiresAt,
    this.refreshExpiresAt,
    this.requiresTwoFactor = false,
    this.requiresVerification = false,
    this.deviceTrusted = false,
    this.isRestoredSession = false,
  });

  factory EnhancedAuthResult.fromJson(Map<String, dynamic> json) =>
      _$EnhancedAuthResultFromJson(json);

  factory EnhancedAuthResult.fromResponse(Map<String, dynamic> response) {
    return EnhancedAuthResult(
      success: response['success'] ?? false,
      message: response['message'],
      user: response['user'] != null ? User.fromJson(response['user']) : null,
      accessToken: response['access_token'],
      refreshToken: response['refresh_token'],
      sessionId: response['session_id'],
      tokenExpiresAt: response['token_expires_at'] != null 
          ? DateTime.parse(response['token_expires_at']) 
          : null,
      refreshExpiresAt: response['refresh_expires_at'] != null 
          ? DateTime.parse(response['refresh_expires_at']) 
          : null,
      requiresTwoFactor: response['requires_two_factor'] ?? false,
      requiresVerification: response['requires_verification'] ?? false,
      deviceTrusted: response['device_trusted'] ?? false,
    );
  }

  Map<String, dynamic> toJson() => _$EnhancedAuthResultToJson(this);
}
