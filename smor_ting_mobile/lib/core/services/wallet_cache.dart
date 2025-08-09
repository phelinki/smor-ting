import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:hive_flutter/hive_flutter.dart';

class WalletCache {
  static const _boxName = 'wallet_cache';
  static const _keyBalances = 'balances';
  static const _encKeyName = 'wallet_cache_key_v1';

  final FlutterSecureStorage _secureStorage;
  Box<String>? _box;

  WalletCache(this._secureStorage);

  Future<void> init() async {
    await Hive.initFlutter();
    final keyHex = await _getOrCreateKey();
    final keyBytes = base64Decode(keyHex);
    _box = await Hive.openBox<String>(_boxName, encryptionCipher: HiveAesCipher(keyBytes));
  }

  Future<String> _getOrCreateKey() async {
    final existing = await _secureStorage.read(key: _encKeyName);
    if (existing != null) return existing;
    final key = Hive.generateSecureKey();
    final b64 = base64Encode(key);
    await _secureStorage.write(key: _encKeyName, value: b64);
    return b64;
  }

  Future<void> saveBalances(Map<String, dynamic> balances) async {
    await _box?.put(_keyBalances, jsonEncode(balances));
  }

  Map<String, dynamic>? getBalances() {
    final raw = _box?.get(_keyBalances);
    if (raw == null) return null;
    return jsonDecode(raw) as Map<String, dynamic>;
  }
}


