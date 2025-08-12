// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'enhanced_auth_models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

EnhancedLoginRequest _$EnhancedLoginRequestFromJson(Map<String, dynamic> json) =>
    EnhancedLoginRequest(
      email: json['email'] as String,
      password: json['password'] as String,
      rememberMe: json['remember_me'] as bool? ?? false,
      deviceInfo: DeviceFingerprint.fromJson(json['device_info'] as Map<String, dynamic>),
      captchaToken: json['captcha_token'] as String?,
      twoFactorCode: json['two_factor_code'] as String?,
    );

Map<String, dynamic> _$EnhancedLoginRequestToJson(EnhancedLoginRequest instance) => <String, dynamic>{
      'email': instance.email,
      'password': instance.password,
      'remember_me': instance.rememberMe,
      'device_info': instance.deviceInfo.toJson(),
      'captcha_token': instance.captchaToken,
      'two_factor_code': instance.twoFactorCode,
    };

EnhancedAuthResult _$EnhancedAuthResultFromJson(Map<String, dynamic> json) =>
    EnhancedAuthResult(
      success: json['success'] as bool,
      message: json['message'] as String?,
      user: json['user'] == null ? null : User.fromJson(json['user'] as Map<String, dynamic>),
      accessToken: json['access_token'] as String?,
      refreshToken: json['refresh_token'] as String?,
      sessionId: json['session_id'] as String?,
      tokenExpiresAt: json['token_expires_at'] == null
          ? null
          : DateTime.parse(json['token_expires_at'] as String),
      refreshExpiresAt: json['refresh_expires_at'] == null
          ? null
          : DateTime.parse(json['refresh_expires_at'] as String),
      requiresTwoFactor: json['requires_two_factor'] as bool? ?? false,
      requiresVerification: json['requires_verification'] as bool? ?? false,
      deviceTrusted: json['device_trusted'] as bool? ?? false,
      isRestoredSession: json['is_restored_session'] as bool? ?? false,
    );

Map<String, dynamic> _$EnhancedAuthResultToJson(EnhancedAuthResult instance) => <String, dynamic>{
      'success': instance.success,
      'message': instance.message,
      'user': instance.user?.toJson(),
      'access_token': instance.accessToken,
      'refresh_token': instance.refreshToken,
      'session_id': instance.sessionId,
      'token_expires_at': instance.tokenExpiresAt?.toIso8601String(),
      'refresh_expires_at': instance.refreshExpiresAt?.toIso8601String(),
      'requires_two_factor': instance.requiresTwoFactor,
      'requires_verification': instance.requiresVerification,
      'device_trusted': instance.deviceTrusted,
      'is_restored_session': instance.isRestoredSession,
    };

SessionInfo _$SessionInfoFromJson(Map<String, dynamic> json) => SessionInfo(
      sessionId: json['session_id'] as String,
      deviceName: json['device_name'] as String,
      deviceType: json['device_type'] as String,
      ipAddress: json['ip_address'] as String,
      lastActivity: DateTime.parse(json['last_activity'] as String),
      createdAt: DateTime.parse(json['created_at'] as String),
      isCurrent: json['is_current'] as bool? ?? false,
    );

Map<String, dynamic> _$SessionInfoToJson(SessionInfo instance) => <String, dynamic>{
      'session_id': instance.sessionId,
      'device_name': instance.deviceName,
      'device_type': instance.deviceType,
      'ip_address': instance.ipAddress,
      'last_activity': instance.lastActivity.toIso8601String(),
      'created_at': instance.createdAt.toIso8601String(),
      'is_current': instance.isCurrent,
    };

LockoutInfo _$LockoutInfoFromJson(Map<String, dynamic> json) => LockoutInfo(
      lockedUntil: DateTime.parse(json['locked_until'] as String),
      remainingAttempts: json['remaining_attempts'] as int,
      timeUntilUnlock: json['time_until_unlock'] as int?,
    );

Map<String, dynamic> _$LockoutInfoToJson(LockoutInfo instance) => <String, dynamic>{
      'locked_until': instance.lockedUntil.toIso8601String(),
      'remaining_attempts': instance.remainingAttempts,
      'time_until_unlock': instance.timeUntilUnlock,
    };
