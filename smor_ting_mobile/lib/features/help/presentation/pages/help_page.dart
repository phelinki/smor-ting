import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class HelpPage extends ConsumerWidget {
  const HelpPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        title: const Text(
          'Help & Support',
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

            // Search Help
            Container(
              margin: const EdgeInsets.all(16),
              child: TextField(
                decoration: InputDecoration(
                  hintText: 'Search for help...',
                  prefixIcon: const Icon(Icons.search, color: Color(0xFF002868)),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: Color(0xFFE0E0E0)),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: Color(0xFFD21034), width: 2),
                  ),
                  fillColor: Colors.white,
                  filled: true,
                ),
              ),
            ),

            // Quick Actions
            _HelpSection(
              title: 'Quick Actions',
              children: [
                _HelpTile(
                  icon: Icons.support_agent,
                  title: 'Contact Support',
                  subtitle: 'Get help from our support team',
                  onTap: () => _showContactDialog(context),
                ),
                _HelpTile(
                  icon: Icons.bug_report,
                  title: 'Report a Problem',
                  subtitle: 'Let us know about any issues',
                  onTap: () => _showReportDialog(context),
                ),
                _HelpTile(
                  icon: Icons.feedback,
                  title: 'Send Feedback',
                  subtitle: 'Share your thoughts and suggestions',
                  onTap: () => _showFeedbackDialog(context),
                ),
              ],
            ),

            const SizedBox(height: 16),

            // FAQ Section
            _HelpSection(
              title: 'Frequently Asked Questions',
              children: [
                _FAQTile(
                  question: 'How do I book a service?',
                  answer: 'Browse service categories, select a service, choose your preferred agent, and complete the booking with your details.',
                ),
                _FAQTile(
                  question: 'How do I verify my email?',
                  answer: 'After registration or login, you\'ll receive a 6-digit OTP code via email. Enter this code to verify your account.',
                ),
                _FAQTile(
                  question: 'What payment methods are accepted?',
                  answer: 'We accept Flutterwave, Orange Money, MTN Mobile Money, and bank transfers in both USD and LRD.',
                ),
                _FAQTile(
                  question: 'How do I become an agent?',
                  answer: 'Register with an agent account, complete your profile, add your services, and wait for verification.',
                ),
                _FAQTile(
                  question: 'How do I cancel a booking?',
                  answer: 'Go to your bookings, select the booking you want to cancel, and tap the cancel button. Cancellation policies may apply.',
                ),
              ],
            ),

            const SizedBox(height: 16),

            // Contact Information
            _HelpSection(
              title: 'Contact Information',
              children: [
                Container(
                  padding: const EdgeInsets.all(20),
                  child: Column(
                    children: [
                      _ContactInfo(
                        icon: Icons.email,
                        title: 'Email',
                        value: 'support@smorting.com',
                      ),
                      const SizedBox(height: 16),
                      _ContactInfo(
                        icon: Icons.phone,
                        title: 'Phone',
                        value: '+231 123 456 789',
                      ),
                      const SizedBox(height: 16),
                      _ContactInfo(
                        icon: Icons.access_time,
                        title: 'Support Hours',
                        value: 'Monday - Friday: 8AM - 6PM GMT',
                      ),
                      const SizedBox(height: 16),
                      _ContactInfo(
                        icon: Icons.location_on,
                        title: 'Address',
                        value: 'Monrovia, Liberia',
                      ),
                    ],
                  ),
                ),
              ],
            ),

            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }

  void _showContactDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Contact Support'),
        content: const Text('Choose how you\'d like to contact our support team:'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              Navigator.of(context).pop();
              // TODO: Open email client
            },
            child: const Text('Email'),
          ),
          TextButton(
            onPressed: () {
              Navigator.of(context).pop();
              // TODO: Open phone dialer
            },
            child: const Text('Call'),
          ),
        ],
      ),
    );
  }

  void _showReportDialog(BuildContext context) {
    final controller = TextEditingController();
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Report a Problem'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('Please describe the problem you\'re experiencing:'),
            const SizedBox(height: 16),
            TextField(
              controller: controller,
              maxLines: 4,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                hintText: 'Describe the problem...',
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
              Navigator.of(context).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Problem report sent successfully'),
                  backgroundColor: Colors.green,
                ),
              );
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFD21034),
            ),
            child: const Text('Send Report'),
          ),
        ],
      ),
    );
  }

  void _showFeedbackDialog(BuildContext context) {
    final controller = TextEditingController();
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Send Feedback'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('We\'d love to hear your thoughts and suggestions:'),
            const SizedBox(height: 16),
            TextField(
              controller: controller,
              maxLines: 4,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                hintText: 'Share your feedback...',
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
              Navigator.of(context).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Feedback sent successfully'),
                  backgroundColor: Colors.green,
                ),
              );
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFD21034),
            ),
            child: const Text('Send Feedback'),
          ),
        ],
      ),
    );
  }
}

class _HelpSection extends StatelessWidget {
  final String title;
  final List<Widget> children;

  const _HelpSection({
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

class _HelpTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  const _HelpTile({
    required this.icon,
    required this.title,
    required this.subtitle,
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

class _FAQTile extends StatefulWidget {
  final String question;
  final String answer;

  const _FAQTile({
    required this.question,
    required this.answer,
  });

  @override
  State<_FAQTile> createState() => _FAQTileState();
}

class _FAQTileState extends State<_FAQTile> {
  bool _isExpanded = false;

  @override
  Widget build(BuildContext context) {
    return ExpansionTile(
      title: Text(
        widget.question,
        style: const TextStyle(
          fontWeight: FontWeight.w500,
          color: Color(0xFF002868),
        ),
      ),
      children: [
        Padding(
          padding: const EdgeInsets.all(16),
          child: Text(
            widget.answer,
            style: TextStyle(
              color: Colors.grey[700],
              fontSize: 14,
            ),
          ),
        ),
      ],
    );
  }
}

class _ContactInfo extends StatelessWidget {
  final IconData icon;
  final String title;
  final String value;

  const _ContactInfo({
    required this.icon,
    required this.title,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, color: const Color(0xFF002868)),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                title,
                style: const TextStyle(
                  fontWeight: FontWeight.w500,
                  color: Color(0xFF002868),
                ),
              ),
              Text(
                value,
                style: TextStyle(
                  color: Colors.grey[600],
                  fontSize: 14,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}