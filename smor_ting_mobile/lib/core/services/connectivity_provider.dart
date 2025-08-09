import 'dart:async';

import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class ConnectivityState {
  final bool isOnline;
  const ConnectivityState(this.isOnline);
}

class ConnectivityNotifier extends StateNotifier<ConnectivityState> {
  final Connectivity _connectivity;
  StreamSubscription<List<ConnectivityResult>>? _sub;

  ConnectivityNotifier(this._connectivity) : super(const ConnectivityState(true)) {
    _init();
  }

  Future<void> _init() async {
    final result = await _connectivity.checkConnectivity();
    state = ConnectivityState(result.any((r) => r != ConnectivityResult.none));
    _sub = _connectivity.onConnectivityChanged.listen((results) {
      final online = results.any((r) => r != ConnectivityResult.none);
      state = ConnectivityState(online);
    });
  }

  @override
  void dispose() {
    _sub?.cancel();
    super.dispose();
  }
}

final connectivityProvider = StateNotifierProvider<ConnectivityNotifier, ConnectivityState>((ref) {
  return ConnectivityNotifier(Connectivity());
});


