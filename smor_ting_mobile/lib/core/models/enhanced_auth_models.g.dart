// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'enhanced_auth_models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

EnhancedLoginRequest _$EnhancedLoginRequestFromJson(
        Map<String, dynamic> json) =>
    EnhancedLoginRequest(
      email: json['email'] as String,
      password: json['password'] as String,
      rememberMe: json['rememberMe'] as bool? ?? false,
      deviceInfo: DeviceFingerprint.fromJson(
          json['deviceInfo'] as Map<String, dynamic>),
      captchaToken: json['captchaToken'] as String?,
      twoFactorCode: json['twoFactorCode'] as String?,
    );

Map<String, dynamic> _$EnhancedLoginRequestToJson(
        EnhancedLoginRequest instance) =>
    <String, dynamic>{
      'email': instance.email,
      'password': instance.password,
      'rememberMe': instance.rememberMe,
      'deviceInfo': instance.deviceInfo,
      'captchaToken': instance.captchaToken,
      'twoFactorCode': instance.twoFactorCode,
    };

EnhancedAuthResult _$EnhancedAuthResultFromJson(Map<String, dynamic> json) =>
    EnhancedAuthResult(
      success: json['success'] as bool,
      message: json['message'] as String?,
      user: json['user'] == null
          ? null
          : User.fromJson(json['user'] as Map<String, dynamic>),
      accessToken: json['accessToken'] as String?,
      refreshToken: json['refreshToken'] as String?,
      sessionId: json['sessionId'] as String?,
      tokenExpiresAt: json['tokenExpiresAt'] == null
          ? null
          : DateTime.parse(json['tokenExpiresAt'] as String),
      refreshExpiresAt: json['refreshExpiresAt'] == null
          ? null
          : DateTime.parse(json['refreshExpiresAt'] as String),
      requiresTwoFactor: json['requiresTwoFactor'] as bool? ?? false,
      requiresVerification: json['requiresVerification'] as bool? ?? false,
      deviceTrusted: json['deviceTrusted'] as bool? ?? false,
      isRestoredSession: json['isRestoredSession'] as bool? ?? false,
      requiresCaptcha: json['requiresCaptcha'] as bool? ?? false,
      remainingAttempts: (json['remainingAttempts'] as num?)?.toInt(),
      lockoutInfo: json['lockoutInfo'] == null
          ? null
          : LockoutInfo.fromJson(json['lockoutInfo'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$EnhancedAuthResultToJson(EnhancedAuthResult instance) =>
    <String, dynamic>{
      'success': instance.success,
      'message': instance.message,
      'user': instance.user,
      'accessToken': instance.accessToken,
      'refreshToken': instance.refreshToken,
      'sessionId': instance.sessionId,
      'tokenExpiresAt': instance.tokenExpiresAt?.toIso8601String(),
      'refreshExpiresAt': instance.refreshExpiresAt?.toIso8601String(),
      'requiresTwoFactor': instance.requiresTwoFactor,
      'requiresVerification': instance.requiresVerification,
      'deviceTrusted': instance.deviceTrusted,
      'isRestoredSession': instance.isRestoredSession,
      'requiresCaptcha': instance.requiresCaptcha,
      'remainingAttempts': instance.remainingAttempts,
      'lockoutInfo': instance.lockoutInfo,
    };

SessionInfo _$SessionInfoFromJson(Map<String, dynamic> json) => SessionInfo(
      sessionId: json['sessionId'] as String,
      deviceName: json['deviceName'] as String,
      deviceType: json['deviceType'] as String,
      ipAddress: json['ipAddress'] as String,
      lastActivity: DateTime.parse(json['lastActivity'] as String),
      createdAt: DateTime.parse(json['createdAt'] as String),
      isCurrent: json['isCurrent'] as bool? ?? false,
    );

Map<String, dynamic> _$SessionInfoToJson(SessionInfo instance) =>
    <String, dynamic>{
      'sessionId': instance.sessionId,
      'deviceName': instance.deviceName,
      'deviceType': instance.deviceType,
      'ipAddress': instance.ipAddress,
      'lastActivity': instance.lastActivity.toIso8601String(),
      'createdAt': instance.createdAt.toIso8601String(),
      'isCurrent': instance.isCurrent,
    };

LockoutInfo _$LockoutInfoFromJson(Map<String, dynamic> json) => LockoutInfo(
      lockedUntil: DateTime.parse(json['lockedUntil'] as String),
      remainingAttempts: (json['remainingAttempts'] as num).toInt(),
      timeUntilUnlock: (json['timeUntilUnlock'] as num?)?.toInt(),
    );

Map<String, dynamic> _$LockoutInfoToJson(LockoutInfo instance) =>
    <String, dynamic>{
      'lockedUntil': instance.lockedUntil.toIso8601String(),
      'remainingAttempts': instance.remainingAttempts,
      'timeUntilUnlock': instance.timeUntilUnlock,
    };
