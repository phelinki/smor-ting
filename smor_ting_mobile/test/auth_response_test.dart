import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:smor_ting_mobile/core/models/user.dart';

void main() {
  group('AuthResponse JSON Parsing', () {
    test('should parse AuthResponse with requires_otp false', () {
      // Arrange - This is the exact JSON format returned by the backend
      const jsonString = '''
      {
        "user": {
          "id": "6898ac7c3dc75dac76cd788c",
          "email": "libworker@smorting.com",
          "first_name": "Agent",
          "last_name": "Test",
          "phone": "231999999999",
          "role": "provider",
          "is_email_verified": false,
          "profile_image": "",
          "address": {
            "street": "",
            "city": "",
            "county": "",
            "country": "",
            "latitude": 0,
            "longitude": 0
          },
          "created_at": "2025-08-10T14:28:12.834255791Z",
          "updated_at": "2025-08-10T14:28:12.834255912Z"
        },
        "access_token": "test_access_token",
        "refresh_token": "test_refresh_token",
        "requires_otp": false
      }
      ''';

      final json = Map<String, dynamic>.from(
        jsonDecode(jsonString) as Map<String, dynamic>,
      );

      // Act
      final authResponse = AuthResponse.fromJson(json);

      // Assert
      expect(authResponse.user.email, equals('libworker@smorting.com'));
      expect(authResponse.user.firstName, equals('Agent'));
      expect(authResponse.user.lastName, equals('Test'));
      expect(authResponse.user.role, equals(UserRole.provider));
      expect(authResponse.accessToken, equals('test_access_token'));
      expect(authResponse.refreshToken, equals('test_refresh_token'));
      expect(authResponse.requiresOTP, equals(false));
    });

    test('should parse AuthResponse with requires_otp true', () {
      // Arrange
      final json = {
        'user': {
          'id': '6898ac7c3dc75dac76cd788c',
          'email': 'test@example.com',
          'first_name': 'Test',
          'last_name': 'User',
          'phone': '231999999999',
          'role': 'customer',
          'is_email_verified': false,
          'profile_image': '',
          'address': {
            'street': '',
            'city': '',
            'county': '',
            'country': '',
            'latitude': 0,
            'longitude': 0,
          },
          'created_at': '2025-08-10T14:28:12.834255791Z',
          'updated_at': '2025-08-10T14:28:12.834255912Z',
        },
        'access_token': null,
        'refresh_token': null,
        'requires_otp': true,
      };

      // Act
      final authResponse = AuthResponse.fromJson(json);

      // Assert
      expect(authResponse.user.email, equals('test@example.com'));
      expect(authResponse.user.role, equals(UserRole.customer));
      expect(authResponse.accessToken, isNull);
      expect(authResponse.refreshToken, isNull);
      expect(authResponse.requiresOTP, equals(true));
    });

    test('should handle missing requires_otp field gracefully', () {
      // Arrange - JSON without requires_otp field (for backward compatibility)
      final json = {
        'user': {
          'id': '6898ac7c3dc75dac76cd788c',
          'email': 'test@example.com',
          'first_name': 'Test',
          'last_name': 'User',
          'phone': '231999999999',
          'role': 'customer',
          'is_email_verified': false,
          'profile_image': '',
          'address': {
            'street': '',
            'city': '',
            'county': '',
            'country': '',
            'latitude': 0,
            'longitude': 0,
          },
          'created_at': '2025-08-10T14:28:12.834255791Z',
          'updated_at': '2025-08-10T14:28:12.834255912Z',
        },
        'access_token': 'test_access_token',
        'refresh_token': 'test_refresh_token',
        // Note: requires_otp field is missing
      };

      // Act
      final authResponse = AuthResponse.fromJson(json);

      // Assert - Should default to false when field is missing
      expect(authResponse.user.email, equals('test@example.com'));
      expect(authResponse.accessToken, equals('test_access_token'));
      expect(authResponse.requiresOTP, equals(false)); // Should default to false
    });

    test('should serialize AuthResponse to JSON correctly', () {
      // Arrange
      final user = User(
        id: '6898ac7c3dc75dac76cd788c',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        phone: '231999999999',
        role: UserRole.provider,
        isEmailVerified: false,
        profileImage: '',
        address: const Address(
          street: '',
          city: '',
          county: '',
          country: '',
          latitude: 0,
          longitude: 0,
        ),
        createdAt: DateTime.parse('2025-08-10T14:28:12.834255791Z'),
        updatedAt: DateTime.parse('2025-08-10T14:28:12.834255912Z'),
      );

      final authResponse = AuthResponse(
        user: user,
        accessToken: 'test_access_token',
        refreshToken: 'test_refresh_token',
        requiresOTP: false,
      );

      // Act
      final json = authResponse.toJson();

      // Assert
      expect(json['requires_otp'], equals(false));
      expect(json['access_token'], equals('test_access_token'));
      expect(json['refresh_token'], equals('test_refresh_token'));
      expect(json['user'], isA<User>());
    });
  });
}
