import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../../../core/models/consent.dart';
import '../../../../core/theme/app_theme.dart';
import '../providers/consent_provider.dart';

/// Inline consent widget for auth flows
class InlineConsentWidget extends ConsumerStatefulWidget {
  final Map<ConsentType, bool> consents;
  final ValueChanged<Map<ConsentType, bool>> onConsentsChanged;
  final bool showOptionalConsents;

  const InlineConsentWidget({
    super.key,
    required this.consents,
    required this.onConsentsChanged,
    this.showOptionalConsents = false,
  });

  @override
  ConsumerState<InlineConsentWidget> createState() => _InlineConsentWidgetState();
}

class _InlineConsentWidgetState extends ConsumerState<InlineConsentWidget> {
  @override
  void initState() {
    super.initState();
    // Load consent requirements on widget initialization
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(consentNotifierProvider.notifier).loadConsentRequirements();
    });
  }

  @override
  Widget build(BuildContext context) {
    final consentState = ref.watch(consentNotifierProvider);

    return switch (consentState) {
      ConsentLoading() => const Center(
          child: Padding(
            padding: EdgeInsets.all(16.0),
            child: CircularProgressIndicator(),
          ),
        ),
      ConsentError(:final message) => Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: AppTheme.error.withValues(alpha: 0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Text(
            'Error loading consent requirements: $message',
            style: const TextStyle(color: AppTheme.error),
          ),
        ),
      ConsentLoaded(:final requirements, userConsent: _) => _buildConsentList(requirements),
    };
  }

  Widget _buildConsentList(List<ConsentRequirement> requirements) {
    final filteredRequirements = widget.showOptionalConsents
        ? requirements
        : requirements.where((r) => r.required).toList();

    if (filteredRequirements.isEmpty) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Privacy & Consent',
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.w600,
            color: AppTheme.textPrimary,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'To create your account, we need your consent for:',
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            color: AppTheme.textSecondary,
          ),
        ),
        const SizedBox(height: 16),
        ...filteredRequirements.map(_buildConsentItem),
      ],
    );
  }

  Widget _buildConsentItem(ConsentRequirement requirement) {
    final isConsented = widget.consents[requirement.type] ?? false;

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Checkbox(
            value: isConsented,
            onChanged: (value) {
              final updatedConsents = Map<ConsentType, bool>.from(widget.consents);
              updatedConsents[requirement.type] = value ?? false;
              widget.onConsentsChanged(updatedConsents);
            },
            activeColor: AppTheme.secondaryBlue,
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: RichText(
                        text: TextSpan(
                          style: DefaultTextStyle.of(context).style,
                          children: [
                            const TextSpan(text: 'I agree to the '),
                            TextSpan(
                              text: requirement.title.toLowerCase(),
                              style: const TextStyle(
                                color: AppTheme.secondaryBlue,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                            if (requirement.required)
                              const TextSpan(
                                text: ' *',
                                style: TextStyle(
                                  color: AppTheme.error,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                          ],
                        ),
                      ),
                    ),
                    if (requirement.documentUrl != null)
                      IconButton(
                        icon: const Icon(
                          Icons.open_in_new,
                          size: 16,
                          color: AppTheme.secondaryBlue,
                        ),
                        onPressed: () => _openDocument(requirement.documentUrl!),
                        constraints: const BoxConstraints(
                          minWidth: 24,
                          minHeight: 24,
                        ),
                        padding: const EdgeInsets.all(4),
                        tooltip: 'Read full document',
                      ),
                  ],
                ),
                if (requirement.description.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(top: 4),
                    child: Text(
                      requirement.description,
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: AppTheme.textSecondary,
                      ),
                    ),
                  ),
                if (requirement.required)
                  Padding(
                    padding: const EdgeInsets.only(top: 4),
                    child: Text(
                      'Required to use the service',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: AppTheme.error,
                        fontWeight: FontWeight.w500,
                        fontSize: 11,
                      ),
                    ),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _openDocument(String url) async {
    try {
      final uri = Uri.parse(url);
      if (await canLaunchUrl(uri)) {
        await launchUrl(uri, mode: LaunchMode.externalApplication);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Could not open document: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }
}

/// Simple terms and privacy footer widget
class TermsPrivacyFooter extends StatelessWidget {
  const TermsPrivacyFooter({super.key});

  @override
  Widget build(BuildContext context) {
    return RichText(
      text: TextSpan(
        style: Theme.of(context).textTheme.bodySmall?.copyWith(
          color: AppTheme.textSecondary,
        ),
        children: const [
          TextSpan(text: 'By continuing, you agree to our '),
          TextSpan(
            text: 'Terms of Service',
            style: TextStyle(
              color: AppTheme.secondaryBlue,
              decoration: TextDecoration.underline,
            ),
          ),
          TextSpan(text: ' and '),
          TextSpan(
            text: 'Privacy Policy',
            style: TextStyle(
              color: AppTheme.secondaryBlue,
              decoration: TextDecoration.underline,
            ),
          ),
        ],
      ),
      textAlign: TextAlign.center,
    );
  }
}
