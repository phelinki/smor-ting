// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

User _$UserFromJson(Map<String, dynamic> json) => User(
      id: json['id'] as String,
      email: json['email'] as String,
      firstName: json['first_name'] as String,
      lastName: json['last_name'] as String,
      phone: json['phone'] as String,
      role: $enumDecode(_$UserRoleEnumMap, json['role']),
      isEmailVerified: json['is_email_verified'] as bool,
      profileImage: json['profile_image'] as String?,
      address: json['address'] == null
          ? null
          : Address.fromJson(json['address'] as Map<String, dynamic>),
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );

Map<String, dynamic> _$UserToJson(User instance) => <String, dynamic>{
      'id': instance.id,
      'email': instance.email,
      'first_name': instance.firstName,
      'last_name': instance.lastName,
      'phone': instance.phone,
      'role': _$UserRoleEnumMap[instance.role]!,
      'is_email_verified': instance.isEmailVerified,
      'profile_image': instance.profileImage,
      'address': instance.address,
      'created_at': instance.createdAt.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
    };

const _$UserRoleEnumMap = {
  UserRole.customer: 'customer',
  UserRole.provider: 'provider',
  UserRole.admin: 'admin',
};

Address _$AddressFromJson(Map<String, dynamic> json) => Address(
      street: json['street'] as String,
      city: json['city'] as String,
      county: json['county'] as String,
      country: json['country'] as String,
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
    );

Map<String, dynamic> _$AddressToJson(Address instance) => <String, dynamic>{
      'street': instance.street,
      'city': instance.city,
      'county': instance.county,
      'country': instance.country,
      'latitude': instance.latitude,
      'longitude': instance.longitude,
    };

LoginRequest _$LoginRequestFromJson(Map<String, dynamic> json) => LoginRequest(
      email: json['email'] as String,
      password: json['password'] as String,
    );

Map<String, dynamic> _$LoginRequestToJson(LoginRequest instance) =>
    <String, dynamic>{
      'email': instance.email,
      'password': instance.password,
    };

RegisterRequest _$RegisterRequestFromJson(Map<String, dynamic> json) =>
    RegisterRequest(
      email: json['email'] as String,
      password: json['password'] as String,
      firstName: json['first_name'] as String,
      lastName: json['last_name'] as String,
      phone: json['phone'] as String,
      role: $enumDecode(_$UserRoleEnumMap, json['role']),
    );

Map<String, dynamic> _$RegisterRequestToJson(RegisterRequest instance) =>
    <String, dynamic>{
      'email': instance.email,
      'password': instance.password,
      'first_name': instance.firstName,
      'last_name': instance.lastName,
      'phone': instance.phone,
      'role': _$UserRoleEnumMap[instance.role]!,
    };

VerifyOTPRequest _$VerifyOTPRequestFromJson(Map<String, dynamic> json) =>
    VerifyOTPRequest(
      email: json['email'] as String,
      otp: json['otp'] as String,
    );

Map<String, dynamic> _$VerifyOTPRequestToJson(VerifyOTPRequest instance) =>
    <String, dynamic>{
      'email': instance.email,
      'otp': instance.otp,
    };

AuthResponse _$AuthResponseFromJson(Map<String, dynamic> json) => AuthResponse(
      user: User.fromJson(json['user'] as Map<String, dynamic>),
      accessToken: json['access_token'] as String?,
      refreshToken: json['refresh_token'] as String?,
      requiresOTP: json['requires_otp'] as bool,
    );

Map<String, dynamic> _$AuthResponseToJson(AuthResponse instance) =>
    <String, dynamic>{
      'user': instance.user,
      'access_token': instance.accessToken,
      'refresh_token': instance.refreshToken,
      'requires_otp': instance.requiresOTP,
    };
