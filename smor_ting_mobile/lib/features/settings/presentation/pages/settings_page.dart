import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:local_auth/local_auth.dart';

import '../../../../core/services/enhanced_auth_service.dart';
import '../../../../core/models/user.dart';
import '../../../auth/presentation/providers/enhanced_auth_provider.dart';
import '../widgets/biometric_settings_widget.dart';

class SettingsPage extends ConsumerStatefulWidget {
  const SettingsPage({super.key});

  @override
  ConsumerState<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends ConsumerState<SettingsPage> {
  bool _pushNotifications = true;
  bool _emailNotifications = true;
  bool _smsNotifications = false;
  bool _locationServices = true;
  bool _darkMode = false;

  String _language = 'English';
  String _currency = 'USD';



  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        title: const Text(
          'Settings',
          style: TextStyle(
            fontWeight: FontWeight.w600,
            color: Color(0xFF002868),
          ),
        ),
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF002868)),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            const SizedBox(height: 16),

            // Notifications Section
            _SettingsSection(
              title: 'Notifications',
              children: [
                _SettingsSwitchTile(
                  title: 'Push Notifications',
                  subtitle: 'Receive notifications about bookings and updates',
                  value: _pushNotifications,
                  onChanged: (value) {
                    setState(() {
                      _pushNotifications = value;
                    });
                  },
                ),
                _SettingsSwitchTile(
                  title: 'Email Notifications',
                  subtitle: 'Receive email updates about your account',
                  value: _emailNotifications,
                  onChanged: (value) {
                    setState(() {
                      _emailNotifications = value;
                    });
                  },
                ),
                _SettingsSwitchTile(
                  title: 'SMS Notifications',
                  subtitle: 'Receive text messages for important updates',
                  value: _smsNotifications,
                  onChanged: (value) {
                    setState(() {
                      _smsNotifications = value;
                    });
                  },
                ),
              ],
            ),

            const SizedBox(height: 16),

            // Biometric Authentication Section
            Consumer(
              builder: (context, ref, child) {
                final authState = ref.watch(enhancedAuthNotifierProvider);
                return authState.maybeWhen(
                  authenticated: (user, _, __, ___, ____, _____) {
                    return BiometricSettingsWidget(userEmail: user.email);
                  },
                  orElse: () => const SizedBox.shrink(),
                ) ?? const SizedBox.shrink();
              },
            ),

            const SizedBox(height: 16),

            // Privacy & Security Section
            _SettingsSection(
              title: 'Privacy & Security',
              children: [
                _SettingsSwitchTile(
                  title: 'Location Services',
                  subtitle: 'Allow app to access your location for better service',
                  value: _locationServices,
                  onChanged: (value) {
                    setState(() {
                      _locationServices = value;
                    });
                  },
                ),
                _SettingsTile(
                  title: 'Change Password',
                  subtitle: 'Update your account password',
                  icon: Icons.lock_outline,
                  onTap: () {
                    _showChangePasswordDialog();
                  },
                ),
                _SettingsTile(
                  title: 'KYC Verification',
                  subtitle: 'Verify your identity with SmileID',
                  icon: Icons.verified_user,
                  onTap: () {
                    context.push('/kyc');
                  },
                ),
                _SettingsTile(
                  title: 'Two-Factor Authentication',
                  subtitle: 'Add an extra layer of security',
                  icon: Icons.security,
                  onTap: () {
                    // TODO: Navigate to 2FA setup
                  },
                ),
                _SettingsTile(
                  title: 'Manage Sessions',
                  subtitle: 'View active sessions and sign out from devices',
                  icon: Icons.devices,
                  onTap: () {
                    context.push('/sessions');
                  },
                ),
              ],
            ),

            const SizedBox(height: 16),

            // App Preferences Section
            _SettingsSection(
              title: 'App Preferences',
              children: [
                _SettingsSwitchTile(
                  title: 'Dark Mode',
                  subtitle: 'Use dark theme for better viewing in low light',
                  value: _darkMode,
                  onChanged: (value) {
                    setState(() {
                      _darkMode = value;
                    });
                    // TODO: Implement theme switching
                  },
                ),
                _SettingsDropdownTile(
                  title: 'Language',
                  subtitle: 'Choose your preferred language',
                  value: _language,
                  options: const ['English', 'French', 'Kpelle', 'Bassa'],
                  onChanged: (value) {
                    setState(() {
                      _language = value!;
                    });
                  },
                ),
                _SettingsDropdownTile(
                  title: 'Currency',
                  subtitle: 'Set your preferred currency',
                  value: _currency,
                  options: const ['USD', 'LRD'],
                  onChanged: (value) {
                    setState(() {
                      _currency = value!;
                    });
                  },
                ),
              ],
            ),

            const SizedBox(height: 16),

            // Support Section
            _SettingsSection(
              title: 'Support',
              children: [
                _SettingsTile(
                  title: 'Help Center',
                  subtitle: 'Get help and find answers to common questions',
                  icon: Icons.help_outline,
                  onTap: () {
                    // TODO: Navigate to help center
                  },
                ),
                _SettingsTile(
                  title: 'Contact Support',
                  subtitle: 'Get in touch with our support team',
                  icon: Icons.support_agent,
                  onTap: () {
                    _showContactSupportDialog();
                  },
                ),
                _SettingsTile(
                  title: 'Report a Problem',
                  subtitle: 'Let us know about any issues you\'re experiencing',
                  icon: Icons.bug_report,
                  onTap: () {
                    // TODO: Navigate to problem reporting
                  },
                ),
              ],
            ),

            const SizedBox(height: 16),

            // Legal Section
            _SettingsSection(
              title: 'Legal',
              children: [
                _SettingsTile(
                  title: 'Terms of Service',
                  subtitle: 'Read our terms and conditions',
                  icon: Icons.description,
                  onTap: () {
                    // TODO: Show terms of service
                  },
                ),
                _SettingsTile(
                  title: 'Privacy Policy',
                  subtitle: 'Learn how we protect your privacy',
                  icon: Icons.privacy_tip,
                  onTap: () {
                    // TODO: Show privacy policy
                  },
                ),
                _SettingsTile(
                  title: 'Licenses',
                  subtitle: 'View open source licenses',
                  icon: Icons.code,
                  onTap: () {
                    showLicensePage(
                      context: context,
                      applicationName: 'Smor-Ting',
                      applicationVersion: '1.0.0',
                    );
                  },
                ),
              ],
            ),

            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }



  void _showChangePasswordDialog() {
    final currentPasswordController = TextEditingController();
    final newPasswordController = TextEditingController();
    final confirmPasswordController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Change Password'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: currentPasswordController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: 'Current Password',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: newPasswordController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: 'New Password',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: confirmPasswordController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: 'Confirm New Password',
                border: OutlineInputBorder(),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              // TODO: Implement password change
              Navigator.of(context).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Password changed successfully'),
                  backgroundColor: Colors.green,
                ),
              );
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFD21034),
            ),
            child: const Text('Change Password'),
          ),
        ],
      ),
    );
  }

  void _showContactSupportDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Contact Support'),
        content: const Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Get in touch with our support team:'),
            SizedBox(height: 16),
            Row(
              children: [
                Icon(Icons.email, color: Color(0xFF002868)),
                SizedBox(width: 8),
                Text('support@smorting.com'),
              ],
            ),
            SizedBox(height: 8),
            Row(
              children: [
                Icon(Icons.phone, color: Color(0xFF002868)),
                SizedBox(width: 8),
                Text('+231 123 456 789'),
              ],
            ),
            SizedBox(height: 8),
            Row(
              children: [
                Icon(Icons.access_time, color: Color(0xFF002868)),
                SizedBox(width: 8),
                Text('Mon-Fri: 8AM-6PM GMT'),
              ],
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Close'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.of(context).pop();
              // TODO: Open email client or phone dialer
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFD21034),
            ),
            child: const Text('Contact Now'),
          ),
        ],
      ),
    );
  }
}

class _SettingsSection extends StatelessWidget {
  final String title;
  final List<Widget> children;

  const _SettingsSection({
    required this.title,
    required this.children,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(20),
            child: Text(
              title,
              style: const TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w600,
                color: Color(0xFF002868),
              ),
            ),
          ),
          ...children,
        ],
      ),
    );
  }
}

class _SettingsTile extends StatelessWidget {
  final String title;
  final String subtitle;
  final IconData icon;
  final VoidCallback onTap;

  const _SettingsTile({
    required this.title,
    required this.subtitle,
    required this.icon,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Icon(icon, color: const Color(0xFF002868)),
      title: Text(
        title,
        style: const TextStyle(
          fontWeight: FontWeight.w500,
          color: Color(0xFF002868),
        ),
      ),
      subtitle: Text(
        subtitle,
        style: TextStyle(
          color: Colors.grey[600],
          fontSize: 12,
        ),
      ),
      trailing: const Icon(Icons.arrow_forward_ios, size: 16, color: Colors.grey),
      onTap: onTap,
    );
  }
}

class _SettingsSwitchTile extends StatelessWidget {
  final String title;
  final String subtitle;
  final bool value;
  final ValueChanged<bool> onChanged;

  const _SettingsSwitchTile({
    required this.title,
    required this.subtitle,
    required this.value,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return SwitchListTile(
      title: Text(
        title,
        style: const TextStyle(
          fontWeight: FontWeight.w500,
          color: Color(0xFF002868),
        ),
      ),
      subtitle: Text(
        subtitle,
        style: TextStyle(
          color: Colors.grey[600],
          fontSize: 12,
        ),
      ),
      value: value,
      onChanged: onChanged,
      activeColor: const Color(0xFFD21034),
    );
  }
}

class _SettingsDropdownTile extends StatelessWidget {
  final String title;
  final String subtitle;
  final String value;
  final List<String> options;
  final ValueChanged<String?> onChanged;

  const _SettingsDropdownTile({
    required this.title,
    required this.subtitle,
    required this.value,
    required this.options,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      title: Text(
        title,
        style: const TextStyle(
          fontWeight: FontWeight.w500,
          color: Color(0xFF002868),
        ),
      ),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            subtitle,
            style: TextStyle(
              color: Colors.grey[600],
              fontSize: 12,
            ),
          ),
          const SizedBox(height: 8),
          DropdownButtonFormField<String>(
            value: value,
            decoration: InputDecoration(
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
                borderSide: const BorderSide(color: Color(0xFFE0E0E0)),
              ),
              contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            ),
            items: options.map((String option) {
              return DropdownMenuItem<String>(
                value: option,
                child: Text(option),
              );
            }).toList(),
            onChanged: onChanged,
          ),
        ],
      ),
    );
  }
}