import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/services/connectivity_provider.dart';
import '../../../../core/models/kyc.dart';
import '../providers/kyc_provider.dart';

class KycPage extends ConsumerStatefulWidget {
  const KycPage({super.key});

  @override
  ConsumerState<KycPage> createState() => _KycPageState();
}

class _KycPageState extends ConsumerState<KycPage> {
  final _formKey = GlobalKey<FormState>();
  final _country = TextEditingController(text: 'LR');
  final _idType = TextEditingController(text: 'NIN');
  final _idNumber = TextEditingController();
  final _firstName = TextEditingController();
  final _lastName = TextEditingController();
  final _phone = TextEditingController();

  @override
  void dispose() {
    _country.dispose();
    _idType.dispose();
    _idNumber.dispose();
    _firstName.dispose();
    _lastName.dispose();
    _phone.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final net = ref.watch(connectivityProvider);
    final kyc = ref.watch(kycProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('KYC Verification'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              TextFormField(
                controller: _firstName,
                decoration: const InputDecoration(labelText: 'First Name'),
                validator: (v) => v == null || v.isEmpty ? 'Required' : null,
              ),
              TextFormField(
                controller: _lastName,
                decoration: const InputDecoration(labelText: 'Last Name'),
                validator: (v) => v == null || v.isEmpty ? 'Required' : null,
              ),
              TextFormField(
                controller: _phone,
                decoration: const InputDecoration(labelText: 'Phone (+231...)'),
                validator: (v) => v == null || v.isEmpty ? 'Required' : null,
                keyboardType: TextInputType.phone,
              ),
              TextFormField(
                controller: _idNumber,
                decoration: const InputDecoration(labelText: 'ID Number'),
                validator: (v) => v == null || v.isEmpty ? 'Required' : null,
              ),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: (!net.isOnline || kyc.loading)
                    ? null
                    : () async {
                        if (!_formKey.currentState!.validate()) return;
                        final req = KycRequest(
                          country: _country.text,
                          idType: _idType.text,
                          idNumber: _idNumber.text,
                          firstName: _firstName.text,
                          lastName: _lastName.text,
                          phone: _phone.text,
                        );
                        await ref.read(kycProvider.notifier).submit(req);
                        final state = ref.read(kycProvider);
                        if (state.error != null) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(content: Text(state.error!), backgroundColor: AppTheme.primaryRed),
                          );
                        } else if (state.result != null) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(content: Text('KYC submitted: ${state.result!.status}')),
                          );
                        }
                      },
                child: kyc.loading
                    ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2))
                    : const Text('Submit for Verification'),
              ),
              if (!net.isOnline)
                const Padding(
                  padding: EdgeInsets.only(top: 8),
                  child: Text('You must be online to submit KYC.', style: TextStyle(color: AppTheme.primaryRed)),
                ),
            ],
          ),
        ),
      ),
    );
  }
}


