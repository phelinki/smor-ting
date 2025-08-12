import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../../../core/models/consent.dart';
import '../../../../core/services/consent_service.dart';
import '../../../../core/theme/app_theme.dart';

/// Dialog for collecting user consent
class ConsentDialog extends ConsumerStatefulWidget {
  final String userId;
  final List<ConsentRequirement> requirements;
  final VoidCallback? onConsentsGiven;
  final VoidCallback? onCancel;

  const ConsentDialog({
    super.key,
    required this.userId,
    required this.requirements,
    this.onConsentsGiven,
    this.onCancel,
  });

  @override
  ConsumerState<ConsentDialog> createState() => _ConsentDialogState();
}

class _ConsentDialogState extends ConsumerState<ConsentDialog> {
  final Map<ConsentType, bool> _consents = {};
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    // Initialize with false for all required consents
    for (final requirement in widget.requirements) {
      _consents[requirement.type] = false;
    }
  }

  bool get _canProceed {
    // Check if all required consents are given
    for (final requirement in widget.requirements) {
      if (requirement.required && (_consents[requirement.type] != true)) {
        return false;
      }
    }
    return true;
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text(
        'Privacy & Consent',
        style: TextStyle(
          fontSize: 20,
          fontWeight: FontWeight.bold,
        ),
      ),
      content: SizedBox(
        width: double.maxFinite,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text(
              'We need your consent for the following to provide you with the best experience:',
              style: TextStyle(fontSize: 14),
            ),
            const SizedBox(height: 16),
            Flexible(
              child: ListView.builder(
                shrinkWrap: true,
                itemCount: widget.requirements.length,
                itemBuilder: (context, index) {
                  final requirement = widget.requirements[index];
                  return _buildConsentItem(requirement);
                },
              ),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: widget.onCancel,
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: _canProceed && !_isLoading ? _submitConsents : null,
          child: _isLoading
              ? const SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(strokeWidth: 2),
                )
              : const Text('Continue'),
        ),
      ],
    );
  }

  Widget _buildConsentItem(ConsentRequirement requirement) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    requirement.title,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                if (requirement.required)
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 6,
                      vertical: 2,
                    ),
                    decoration: BoxDecoration(
                      color: AppTheme.error,
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: const Text(
                      'Required',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 10,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 4),
            Text(
              requirement.description,
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey[700],
              ),
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(
                  child: CheckboxListTile(
                    contentPadding: EdgeInsets.zero,
                    title: Text(
                      'I ${requirement.required ? 'agree' : 'consent'} to ${requirement.title.toLowerCase()}',
                      style: const TextStyle(fontSize: 14),
                    ),
                    value: _consents[requirement.type] ?? false,
                    onChanged: (value) {
                      setState(() {
                        _consents[requirement.type] = value ?? false;
                      });
                    },
                    activeColor: AppTheme.secondaryBlue,
                  ),
                ),
                if (requirement.documentUrl != null)
                  IconButton(
                    icon: const Icon(Icons.open_in_new, size: 18),
                    onPressed: () => _openDocument(requirement.documentUrl!),
                    tooltip: 'Read full document',
                  ),
              ],
            ),
          ],
        ),
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

  Future<void> _submitConsents() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final consentService = ref.read(consentServiceProvider);
      await consentService.updateMultipleConsents(
        widget.userId,
        _consents,
        userAgent: 'Smor-Ting Mobile App',
        metadata: {
          'submitted_via': 'consent_dialog',
          'timestamp': DateTime.now().toIso8601String(),
        },
      );

      if (mounted) {
        widget.onConsentsGiven?.call();
        Navigator.of(context).pop(true);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to save consent: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }
}

/// Simple consent checkbox widget for inline use
class ConsentCheckbox extends StatelessWidget {
  final ConsentRequirement requirement;
  final bool value;
  final ValueChanged<bool?> onChanged;

  const ConsentCheckbox({
    super.key,
    required this.requirement,
    required this.value,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return CheckboxListTile(
      title: RichText(
        text: TextSpan(
          style: DefaultTextStyle.of(context).style,
          children: [
            TextSpan(text: 'I agree to the '),
            TextSpan(
              text: requirement.title,
              style: const TextStyle(
                color: AppTheme.secondaryBlue,
                decoration: TextDecoration.underline,
              ),
            ),
            if (requirement.required)
              const TextSpan(
                text: ' *',
                style: TextStyle(color: AppTheme.error),
              ),
          ],
        ),
      ),
      subtitle: requirement.required
          ? const Text(
              'Required to use the service',
              style: TextStyle(
                fontSize: 12,
                color: AppTheme.error,
              ),
            )
          : null,
      value: value,
      onChanged: onChanged,
      activeColor: AppTheme.secondaryBlue,
      controlAffinity: ListTileControlAffinity.leading,
    );
  }
}
