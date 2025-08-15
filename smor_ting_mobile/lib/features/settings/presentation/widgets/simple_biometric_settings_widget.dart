import 'package:flutter/material.dart';

/// Simplified biometric settings widget (disabled for now)
class SimpleBiometricSettingsWidget extends StatelessWidget {
  final String userEmail;

  const SimpleBiometricSettingsWidget({
    super.key,
    required this.userEmail,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Biometric Authentication',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              'Biometric authentication is currently disabled.',
              style: TextStyle(color: Colors.grey),
            ),
            const SizedBox(height: 16),
            SwitchListTile(
              title: const Text('Enable Biometric Login'),
              subtitle: const Text('Coming soon'),
              value: false,
              onChanged: null, // Disabled
            ),
          ],
        ),
      ),
    );
  }
}
