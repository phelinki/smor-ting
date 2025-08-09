import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/models/kyc.dart';
import '../../../../core/services/api_service.dart';

class KycState {
  final bool loading;
  final String? error;
  final KycResponse? result;
  const KycState({this.loading = false, this.error, this.result});

  KycState copyWith({bool? loading, String? error, KycResponse? result}) =>
      KycState(loading: loading ?? this.loading, error: error, result: result ?? this.result);
}

class KycNotifier extends StateNotifier<KycState> {
  final ApiService _api;
  KycNotifier(this._api) : super(const KycState());

  Future<void> submit(KycRequest req) async {
    state = const KycState(loading: true);
    try {
      final res = await _api.submitKyc(req);
      state = KycState(loading: false, result: res);
    } catch (e) {
      state = KycState(loading: false, error: e.toString());
    }
  }
}

final kycProvider = StateNotifierProvider<KycNotifier, KycState>((ref) {
  final api = ref.read(apiServiceProvider);
  return KycNotifier(api);
});


