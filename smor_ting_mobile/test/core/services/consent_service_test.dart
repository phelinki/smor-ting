import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:smor_ting/core/models/consent.dart';
import 'package:smor_ting/core/services/api_service.dart';
import 'package:smor_ting/core/services/consent_service.dart';

import 'consent_service_test.mocks.dart';

@GenerateMocks([ApiService, FlutterSecureStorage])
void main() {
  late ConsentService consentService;
  late MockApiService mockApiService;
  late MockFlutterSecureStorage mockSecureStorage;

  setUp(() {
    mockApiService = MockApiService();
    mockSecureStorage = MockFlutterSecureStorage();
    consentService = ConsentService(mockApiService, mockSecureStorage);
  });

  group('ConsentService', () {
    group('getConsentRequirements', () {
      test('should return requirements from API when available', () async {
        // Arrange
        final apiResponse = {
          'requirements': [
            {
              'type': 'terms_of_service',
              'title': 'Terms of Service',
              'description': 'Accept our terms',
              'version': '1.0',
              'required': true,
              'documentUrl': 'https://example.com/terms',
            },
            {
              'type': 'privacy_policy',
              'title': 'Privacy Policy',
              'description': 'Accept our privacy policy',
              'version': '1.0',
              'required': true,
              'documentUrl': 'https://example.com/privacy',
            },
          ],
        };

        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => apiResponse);

        // Act
        final requirements = await consentService.getConsentRequirements();

        // Assert
        expect(requirements, hasLength(2));
        expect(requirements[0].type, ConsentType.termsOfService);
        expect(requirements[0].title, 'Terms of Service');
        expect(requirements[0].required, true);
        expect(requirements[1].type, ConsentType.privacyPolicy);
      });

      test('should return default requirements when API fails', () async {
        // Arrange
        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenThrow(Exception('API error'));

        // Act
        final requirements = await consentService.getConsentRequirements();

        // Assert
        expect(requirements, isNotEmpty);
        expect(requirements.any((r) => r.type == ConsentType.termsOfService), true);
        expect(requirements.any((r) => r.type == ConsentType.privacyPolicy), true);
      });
    });

    group('getUserConsent', () {
      test('should return consent from API when available', () async {
        // Arrange
        const userId = 'user123';
        final apiResponse = {
          'userId': userId,
          'consents': {
            'terms_of_service': {
              'id': 'consent1',
              'type': 'terms_of_service',
              'granted': true,
              'consentedAt': '2024-01-01T00:00:00.000Z',
              'version': '1.0',
            },
          },
          'lastUpdated': '2024-01-01T00:00:00.000Z',
        };

        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => apiResponse);

        // Act
        final userConsent = await consentService.getUserConsent(userId);

        // Assert
        expect(userConsent, isNotNull);
        expect(userConsent!.userId, userId);
        expect(userConsent.consents, hasLength(1));
        expect(userConsent.consents[ConsentType.termsOfService]?.granted, true);
      });

      test('should return local consent when API fails', () async {
        // Arrange
        const userId = 'user123';
        final localConsentJson = jsonEncode({
          'userId': userId,
          'consents': {
            'privacy_policy': {
              'id': 'consent2',
              'type': 'privacy_policy',
              'granted': true,
              'consentedAt': '2024-01-01T00:00:00.000Z',
              'version': '1.0',
            },
          },
          'lastUpdated': '2024-01-01T00:00:00.000Z',
        });

        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenThrow(Exception('API error'));
        when(mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => localConsentJson);

        // Act
        final userConsent = await consentService.getUserConsent(userId);

        // Assert
        expect(userConsent, isNotNull);
        expect(userConsent!.userId, userId);
        expect(userConsent.consents[ConsentType.privacyPolicy]?.granted, true);
      });

      test('should return null when no consent found', () async {
        // Arrange
        const userId = 'user123';

        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenThrow(Exception('API error'));
        when(mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => null);

        // Act
        final userConsent = await consentService.getUserConsent(userId);

        // Assert
        expect(userConsent, isNull);
      });
    });

    group('updateConsent', () {
      test('should update consent via API and local storage', () async {
        // Arrange
        const userId = 'user123';
        const type = ConsentType.termsOfService;
        const granted = true;

        // Mock consent requirements
        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => {
              'requirements': [
                {
                  'type': 'terms_of_service',
                  'title': 'Terms of Service',
                  'description': 'Accept our terms',
                  'version': '1.0',
                  'required': true,
                },
              ],
            });

        when(mockApiService.post('/api/v1/consent/user/$userId', any))
            .thenAnswer((_) async => {});
        when(mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => null);
        when(mockSecureStorage.write(key: 'smor_ting_consent_v1_$userId', value: any))
            .thenAnswer((_) async {});

        // Act
        await consentService.updateConsent(userId, type, granted);

        // Assert
        verify(mockApiService.post('/api/v1/consent/user/$userId', any)).called(1);
        verify(mockSecureStorage.write(key: 'smor_ting_consent_v1_$userId', value: any))
            .called(1);
      });

      test('should store as pending when API fails', () async {
        // Arrange
        const userId = 'user123';
        const type = ConsentType.termsOfService;
        const granted = true;

        // Mock consent requirements
        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => {
              'requirements': [
                {
                  'type': 'terms_of_service',
                  'title': 'Terms of Service',
                  'description': 'Accept our terms',
                  'version': '1.0',
                  'required': true,
                },
              ],
            });

        when(mockApiService.post('/api/v1/consent/user/$userId', any))
            .thenThrow(Exception('API error'));
        when(mockSecureStorage.read(key: 'smor_ting_pending_consent_v1'))
            .thenAnswer((_) async => null);
        when(mockSecureStorage.write(key: 'smor_ting_pending_consent_v1', value: any))
            .thenAnswer((_) async {});

        // Act & Assert
        expect(
          () => consentService.updateConsent(userId, type, granted),
          throwsA(isA<Exception>()),
        );

        verify(mockSecureStorage.write(key: 'smor_ting_pending_consent_v1', value: any))
            .called(1);
      });
    });

    group('hasRequiredConsents', () {
      test('should return true when all required consents are given', () async {
        // Arrange
        const userId = 'user123';

        // Mock requirements with one required consent
        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => {
              'requirements': [
                {
                  'type': 'terms_of_service',
                  'title': 'Terms of Service',
                  'description': 'Accept our terms',
                  'version': '1.0',
                  'required': true,
                },
                {
                  'type': 'marketing_communications',
                  'title': 'Marketing',
                  'description': 'Marketing emails',
                  'version': '1.0',
                  'required': false,
                },
              ],
            });

        // Mock user consent with required consent granted
        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => {
              'userId': userId,
              'consents': {
                'terms_of_service': {
                  'id': 'consent1',
                  'type': 'terms_of_service',
                  'granted': true,
                  'consentedAt': '2024-01-01T00:00:00.000Z',
                  'version': '1.0',
                },
              },
              'lastUpdated': '2024-01-01T00:00:00.000Z',
            });

        // Act
        final hasRequired = await consentService.hasRequiredConsents(userId);

        // Assert
        expect(hasRequired, true);
      });

      test('should return false when required consents are missing', () async {
        // Arrange
        const userId = 'user123';

        // Mock requirements with required consent
        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => {
              'requirements': [
                {
                  'type': 'terms_of_service',
                  'title': 'Terms of Service',
                  'description': 'Accept our terms',
                  'version': '1.0',
                  'required': true,
                },
              ],
            });

        // Mock user consent without required consent
        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => {
              'userId': userId,
              'consents': {},
              'lastUpdated': '2024-01-01T00:00:00.000Z',
            });

        // Act
        final hasRequired = await consentService.hasRequiredConsents(userId);

        // Assert
        expect(hasRequired, false);
      });

      test('should return false when user consent is null', () async {
        // Arrange
        const userId = 'user123';

        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => {
              'requirements': [
                {
                  'type': 'terms_of_service',
                  'title': 'Terms of Service',
                  'description': 'Accept our terms',
                  'version': '1.0',
                  'required': true,
                },
              ],
            });

        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenThrow(Exception('API error'));
        when(mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => null);

        // Act
        final hasRequired = await consentService.hasRequiredConsents(userId);

        // Assert
        expect(hasRequired, false);
      });
    });

    group('getMissingConsents', () {
      test('should return missing required consents', () async {
        // Arrange
        const userId = 'user123';

        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => {
              'requirements': [
                {
                  'type': 'terms_of_service',
                  'title': 'Terms of Service',
                  'description': 'Accept our terms',
                  'version': '1.0',
                  'required': true,
                },
                {
                  'type': 'privacy_policy',
                  'title': 'Privacy Policy',
                  'description': 'Accept our privacy policy',
                  'version': '1.0',
                  'required': true,
                },
              ],
            });

        // Mock user consent with only one consent
        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => {
              'userId': userId,
              'consents': {
                'terms_of_service': {
                  'id': 'consent1',
                  'type': 'terms_of_service',
                  'granted': true,
                  'consentedAt': '2024-01-01T00:00:00.000Z',
                  'version': '1.0',
                },
              },
              'lastUpdated': '2024-01-01T00:00:00.000Z',
            });

        // Act
        final missing = await consentService.getMissingConsents(userId);

        // Assert
        expect(missing, hasLength(1));
        expect(missing[0].type, ConsentType.privacyPolicy);
      });
    });

    group('isConsentExpired', () {
      test('should return true when consent version is outdated', () async {
        // Arrange
        const userId = 'user123';
        const type = ConsentType.termsOfService;

        // Mock requirements with new version
        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => {
              'requirements': [
                {
                  'type': 'terms_of_service',
                  'title': 'Terms of Service',
                  'description': 'Accept our terms',
                  'version': '2.0', // New version
                  'required': true,
                },
              ],
            });

        // Mock user consent with old version
        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => {
              'userId': userId,
              'consents': {
                'terms_of_service': {
                  'id': 'consent1',
                  'type': 'terms_of_service',
                  'granted': true,
                  'consentedAt': '2024-01-01T00:00:00.000Z',
                  'version': '1.0', // Old version
                },
              },
              'lastUpdated': '2024-01-01T00:00:00.000Z',
            });

        // Act
        final isExpired = await consentService.isConsentExpired(userId, type);

        // Assert
        expect(isExpired, true);
      });

      test('should return false when consent version is current', () async {
        // Arrange
        const userId = 'user123';
        const type = ConsentType.termsOfService;

        // Mock requirements
        when(mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => {
              'requirements': [
                {
                  'type': 'terms_of_service',
                  'title': 'Terms of Service',
                  'description': 'Accept our terms',
                  'version': '1.0',
                  'required': true,
                },
              ],
            });

        // Mock user consent with current version
        when(mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => {
              'userId': userId,
              'consents': {
                'terms_of_service': {
                  'id': 'consent1',
                  'type': 'terms_of_service',
                  'granted': true,
                  'consentedAt': '2024-01-01T00:00:00.000Z',
                  'version': '1.0',
                },
              },
              'lastUpdated': '2024-01-01T00:00:00.000Z',
            });

        // Act
        final isExpired = await consentService.isConsentExpired(userId, type);

        // Assert
        expect(isExpired, false);
      });
    });
  });
}
