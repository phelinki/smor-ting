import 'package:flutter_test/flutter_test.dart';
import 'package:riverpod/riverpod.dart';

import 'package:smor_ting_mobile/core/models/kyc.dart';
import 'package:smor_ting_mobile/core/services/api_service.dart';
import 'package:smor_ting_mobile/features/kyc/presentation/providers/kyc_provider.dart';

class FakeApi extends ApiService {
  FakeApi(): super(baseUrl: 'http://localhost')
  ;
  @override
  Future<KycResponse> submitKyc(KycRequest request) async {
    return KycResponse(status: 'PENDING', reference: 'ref-123');
  }
}

void main() {
  test('kyc provider submits and returns result', () async {
    final container = ProviderContainer(overrides: [
      apiServiceProvider.overrideWithValue(FakeApi()),
    ]);
    addTearDown(container.dispose);

    final notifier = container.read(kycProvider.notifier);
    await notifier.submit(KycRequest(
      country: 'LR', idType: 'NIN', idNumber: '123', firstName: 'A', lastName: 'B', phone: '231770000000'));

    final state = container.read(kycProvider);
    expect(state.loading, false);
    expect(state.error, isNull);
    expect(state.result, isNotNull);
    expect(state.result!.reference, 'ref-123');
  });
}


