import 'package:json_annotation/json_annotation.dart';

part 'user.g.dart';

enum UserRole {
  @JsonValue('customer')
  customer,
  @JsonValue('provider')
  provider,
  @JsonValue('admin')
  admin,
}

@JsonSerializable()
class User {
  @JsonKey(fromJson: _idFromJson, toJson: _idToJson)
  final String id;
  final String email;
  @JsonKey(name: 'first_name')
  final String firstName;
  @JsonKey(name: 'last_name')
  final String lastName;
  @JsonKey(defaultValue: '')
  final String phone;
  @JsonKey(defaultValue: UserRole.customer)
  final UserRole role;
  @JsonKey(name: 'is_email_verified', defaultValue: false)
  final bool isEmailVerified;
  @JsonKey(name: 'profile_image')
  final String? profileImage;
  final Address? address;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const User({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.phone,
    required this.role,
    required this.isEmailVerified,
    this.profileImage,
    this.address,
    required this.createdAt,
    required this.updatedAt,
  });

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
  Map<String, dynamic> toJson() => _$UserToJson(this);

  String get fullName => '$firstName $lastName';
}

// Custom converters for ID field
String _idFromJson(dynamic value) {
  if (value is int) {
    return value.toString();
  }
  return value as String;
}

String _idToJson(String value) => value;

@JsonSerializable()
class Address {
  final String street;
  final String city;
  final String county;
  final String country;
  final double latitude;
  final double longitude;

  const Address({
    required this.street,
    required this.city,
    required this.county,
    required this.country,
    required this.latitude,
    required this.longitude,
  });

  factory Address.fromJson(Map<String, dynamic> json) => _$AddressFromJson(json);
  Map<String, dynamic> toJson() => _$AddressToJson(this);
}

@JsonSerializable()
class LoginRequest {
  final String email;
  final String password;

  const LoginRequest({
    required this.email,
    required this.password,
  });

  factory LoginRequest.fromJson(Map<String, dynamic> json) => _$LoginRequestFromJson(json);
  Map<String, dynamic> toJson() => _$LoginRequestToJson(this);
}

@JsonSerializable()
class RegisterRequest {
  final String email;
  final String password;
  @JsonKey(name: 'first_name')
  final String firstName;
  @JsonKey(name: 'last_name')
  final String lastName;
  final String phone;
  final UserRole role;

  const RegisterRequest({
    required this.email,
    required this.password,
    required this.firstName,
    required this.lastName,
    required this.phone,
    required this.role,
  });

  factory RegisterRequest.fromJson(Map<String, dynamic> json) => _$RegisterRequestFromJson(json);
  Map<String, dynamic> toJson() => _$RegisterRequestToJson(this);
}

@JsonSerializable()
class VerifyOTPRequest {
  final String email;
  final String otp;

  const VerifyOTPRequest({
    required this.email,
    required this.otp,
  });

  factory VerifyOTPRequest.fromJson(Map<String, dynamic> json) => _$VerifyOTPRequestFromJson(json);
  Map<String, dynamic> toJson() => _$VerifyOTPRequestToJson(this);
}

@JsonSerializable()
class AuthResponse {
  final User user;
  @JsonKey(name: 'access_token')
  final String? accessToken;
  @JsonKey(name: 'refresh_token')
  final String? refreshToken;
  @JsonKey(name: 'session_id')
  final String? sessionId;  // Add this line
  @JsonKey(name: 'requires_otp', defaultValue: false)
  final bool requiresOTP;

  const AuthResponse({
    required this.user,
    this.accessToken,
    this.refreshToken,
    this.sessionId,  // Add this line
    this.requiresOTP = false,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) => _$AuthResponseFromJson(json);
  Map<String, dynamic> toJson() => _$AuthResponseToJson(this);
}