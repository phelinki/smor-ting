import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:smor_ting_mobile/core/models/consent.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:smor_ting_mobile/core/services/consent_service.dart';
import 'package:dio/dio.dart';

class MockApiService extends Mock implements ApiService {}
class MockFlutterSecureStorage extends Mock implements FlutterSecureStorage {}

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
              'type': 'termsOfService',
              'title': 'Terms of Service',
              'description': 'Accept our terms',
              'version': '1.0',
              'required': true,
              'document_url': 'https://example.com/terms'
            }
          ]
        };
        
        when(() => mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => Response(
                  requestOptions: RequestOptions(path: ''),
                  data: apiResponse,
                  statusCode: 200,
                ));

        // Act
        final result = await consentService.getConsentRequirements();

        // Assert
        expect(result, isNotEmpty);
        expect(result.first.type, ConsentType.termsOfService);
        expect(result.first.title, 'Terms of Service');
        expect(result.first.required, true);
      });

      test('should return empty list when API returns no requirements', () async {
        // Arrange
        when(() => mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => Response(
                  requestOptions: RequestOptions(path: ''),
                  data: {'requirements': []},
                  statusCode: 200,
                ));

        // Act
        final result = await consentService.getConsentRequirements();

        // Assert
        expect(result, isEmpty);
      });

      test('should throw exception when API call fails', () async {
        // Arrange
        when(() => mockApiService.get('/api/v1/consent/requirements'))
            .thenThrow(DioException(requestOptions: RequestOptions(path: '')));

        // Act & Assert
        expect(
          () => consentService.getConsentRequirements(),
          throwsException,
        );
      });
    });

    group('getUserConsent', () {
      const userId = 'test-user';

      test('should return user consent from API', () async {
        // Arrange
        final apiResponse = {
          'user_id': userId,
          'consents': {
            'termsOfService': {
              'id': 'consent-1',
              'type': 'termsOfService',
              'granted': true,
              'consented_at': '2023-01-01T00:00:00Z',
              'version': '1.0'
            }
          },
          'last_updated': '2023-01-01T00:00:00Z'
        };

        when(() => mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => Response(
                  requestOptions: RequestOptions(path: ''),
                  data: apiResponse,
                  statusCode: 200,
                ));
        
        when(() => mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => null);

        // Act
        final result = await consentService.getUserConsent(userId);

        // Assert
        expect(result?.userId, userId);
        expect(result?.consents.containsKey(ConsentType.termsOfService), true);
      });

      test('should return cached consent when available', () async {
        // Arrange
        final cachedConsent = {
          'user_id': userId,
          'consents': {
            'termsOfService': {
              'id': 'consent-1',
              'type': 'termsOfService',
              'granted': true,
              'consented_at': '2023-01-01T00:00:00Z',
              'version': '1.0'
            }
          },
          'last_updated': '2023-01-01T00:00:00Z'
        };

        when(() => mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => json.encode(cachedConsent));

        // Act
        final result = await consentService.getUserConsent(userId);

        // Assert
        expect(result?.userId, userId);
        verifyNever(() => mockApiService.get(any()));
      });
    });

    group('updateConsent', () {
      const userId = 'test-user';
      const type = ConsentType.termsOfService;

      test('should update consent successfully', () async {
        // Arrange
        when(() => mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => Response(
                  requestOptions: RequestOptions(path: ''),
                  data: {
                    'requirements': [
                      {
                        'type': 'termsOfService',
                        'title': 'Terms of Service',
                        'description': 'Accept our terms',
                        'version': '1.0',
                        'required': true,
                      }
                    ]
                  },
                  statusCode: 200,
                ));

        when(() => mockApiService.post(
              '/api/v1/consent/user/$userId',
              data: any(named: 'data'),
            )).thenAnswer((_) async => Response(
              requestOptions: RequestOptions(path: ''),
              statusCode: 200,
            ));

        when(() => mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => null);
        
        when(() => mockSecureStorage.write(
              key: 'smor_ting_consent_v1_$userId',
              value: any(named: 'value'),
            )).thenAnswer((_) async {});

        // Act
        await consentService.updateConsent(userId, type, true);

        // Assert
        verify(() => mockApiService.post(
          '/api/v1/consent/user/$userId',
          data: any(named: 'data'),
        )).called(1);
        verify(() => mockSecureStorage.write(
          key: 'smor_ting_consent_v1_$userId',
          value: any(named: 'value'),
        )).called(1);
      });
    });

    group('getMissingConsents', () {
      const userId = 'test-user';

      test('should return missing consents correctly', () async {
        // Arrange
        final requirements = [
          const ConsentRequirement(
            type: ConsentType.termsOfService,
            title: 'Terms of Service',
            description: 'Accept our terms',
            version: '1.0',
            required: true,
          ),
          const ConsentRequirement(
            type: ConsentType.privacyPolicy,
            title: 'Privacy Policy',
            description: 'Accept our privacy policy',
            version: '1.0',
            required: true,
          ),
        ];

        final userConsent = UserConsent(
          userId: userId,
          consents: {
            ConsentType.termsOfService: ConsentRecord(
              id: 'consent-1',
              type: ConsentType.termsOfService,
              granted: true,
              consentedAt: DateTime.parse('2023-01-01T00:00:00Z'),
              version: '1.0',
            ),
          },
          lastUpdated: DateTime.parse('2023-01-01T00:00:00Z'),
        );

        when(() => mockApiService.get('/api/v1/consent/requirements'))
            .thenAnswer((_) async => Response(
                  requestOptions: RequestOptions(path: ''),
                  data: {
                    'requirements': requirements
                        .map((r) => r.toJson())
                        .toList(),
                  },
                  statusCode: 200,
                ));

        when(() => mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => Response(
                  requestOptions: RequestOptions(path: ''),
                  data: userConsent.toJson(),
                  statusCode: 200,
                ));
        
        when(() => mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => null);

        // Act
        final missing = await consentService.getMissingConsents(userId);

        // Assert
        expect(missing, hasLength(1));
        expect(missing[0].type, ConsentType.privacyPolicy);
      });
    });

    group('isConsentExpired', () {
      const userId = 'test-user';
      const type = ConsentType.termsOfService;

      test('should return false for non-expired consent', () async {
        // Arrange
        final userConsent = UserConsent(
          userId: userId,
          consents: {
            type: ConsentRecord(
              id: 'consent-1',
              type: type,
              granted: true,
              consentedAt: DateTime.now().subtract(const Duration(days: 30)),
              version: '1.0',
            ),
          },
          lastUpdated: DateTime.now(),
        );

        when(() => mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => Response(
                  requestOptions: RequestOptions(path: ''),
                  data: userConsent.toJson(),
                  statusCode: 200,
                ));
        
        when(() => mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => null);

        // Act
        final isExpired = await consentService.isConsentExpired(userId, type);

        // Assert
        expect(isExpired, false);
      });

      test('should return true for expired consent', () async {
        // Arrange
        final userConsent = UserConsent(
          userId: userId,
          consents: {
            type: ConsentRecord(
              id: 'consent-1',
              type: type,
              granted: true,
              consentedAt: DateTime.now().subtract(const Duration(days: 400)),
              version: '1.0',
            ),
          },
          lastUpdated: DateTime.now(),
        );

        when(() => mockApiService.get('/api/v1/consent/user/$userId'))
            .thenAnswer((_) async => Response(
                  requestOptions: RequestOptions(path: ''),
                  data: userConsent.toJson(),
                  statusCode: 200,
                ));
        
        when(() => mockSecureStorage.read(key: 'smor_ting_consent_v1_$userId'))
            .thenAnswer((_) async => null);

        // Act
        final isExpired = await consentService.isConsentExpired(userId, type);

        // Assert
        expect(isExpired, true);
      });
    });
  });
}