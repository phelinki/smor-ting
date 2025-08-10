import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/constants/app_constants.dart';
import '../../../../core/models/user.dart';
import '../providers/auth_provider.dart';
import '../widgets/custom_text_field.dart';

class NewRegisterPage extends ConsumerStatefulWidget {
  const NewRegisterPage({super.key});

  @override
  ConsumerState<NewRegisterPage> createState() => _NewRegisterPageState();
}

class _NewRegisterPageState extends ConsumerState<NewRegisterPage> {
  final _formKey = GlobalKey<FormState>();
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();
  final _addressController = TextEditingController();
  final _phoneController = TextEditingController();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  
  bool _isPasswordVisible = false;
  bool _isConfirmPasswordVisible = false;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    // Clear any existing auth errors when the page loads
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(authNotifierProvider.notifier).clearError();
    });
  }

  @override
  void dispose() {
    _firstNameController.dispose();
    _lastNameController.dispose();
    _addressController.dispose();
    _phoneController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  Future<void> _handleRegister(UserRole role) async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
    });

    try {
      await ref.read(authNotifierProvider.notifier).register(
        firstName: _firstNameController.text.trim(),
        lastName: _lastNameController.text.trim(),
        email: _emailController.text.trim(),
        phone: _phoneController.text.trim(),
        password: _passwordController.text,
        role: role,
      );
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  String? _validatePhone(String? value) {
    if (value == null || value.isEmpty) {
      return AppConstants.requiredFieldMessage;
    }
    
    // Remove any formatting
    String cleanPhone = value.replaceAll(RegExp(r'[^\d]'), '');
    
    // Check if it's a valid Liberia phone number
    if (!RegExp(AppConstants.liberiaPhonePattern).hasMatch(cleanPhone)) {
      return AppConstants.invalidPhoneMessage;
    }
    
    return null;
  }

  String? _validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return AppConstants.requiredFieldMessage;
    }
    if (value.length < 8) {
      return AppConstants.passwordTooShortMessage;
    }
    // Check for at least one uppercase, one lowercase, one number, and one symbol
    if (!RegExp(r'^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]').hasMatch(value)) {
      return 'Password must contain:\n• One uppercase letter\n• One lowercase letter\n• One number\n• One symbol (@\$!%*?&)';
    }
    return null;
  }

  String? _validateConfirmPassword(String? value) {
    if (value == null || value.isEmpty) {
      return AppConstants.requiredFieldMessage;
    }
    if (value != _passwordController.text) {
      return AppConstants.passwordMismatchMessage;
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authNotifierProvider);

    // Listen to auth state changes for navigation
    ref.listen<AuthState>(authNotifierProvider, (previous, next) {
      if (next is Authenticated) {
        final role = next.user.role;
        if (role == UserRole.provider || role == UserRole.admin) {
          context.go('/agent-dashboard');
        } else {
          context.go('/home');
        }
      } else if (next is RequiresOTP) {
        context.go('/verify-otp?email=${next.email}&fullName=${next.user.fullName}');
      } else if (next is Error) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(next.message),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    });

    return Scaffold(
      backgroundColor: AppTheme.white,
      appBar: AppBar(
        title: const Text('Create Account'),
        backgroundColor: AppTheme.white,
        foregroundColor: AppTheme.textPrimary,
        elevation: 0,
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Create your account',
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                
                const SizedBox(height: 8),
                
                Text(
                  'Fill in your details to get started',
                  style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                    color: AppTheme.gray,
                  ),
                ),
                
                const SizedBox(height: 32),
                
                // First Name Field
                CustomTextField(
                  controller: _firstNameController,
                  labelText: 'First Name',
                  hintText: 'Enter your first name',
                  keyboardType: TextInputType.name,
                  prefixIcon: Icons.person_outlined,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return AppConstants.requiredFieldMessage;
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 20),
                
                // Last Name Field
                CustomTextField(
                  controller: _lastNameController,
                  labelText: 'Last Name',
                  hintText: 'Enter your last name',
                  keyboardType: TextInputType.name,
                  prefixIcon: Icons.person_outlined,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return AppConstants.requiredFieldMessage;
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 20),
                
                // Address Field
                CustomTextField(
                  controller: _addressController,
                  labelText: 'Address',
                  hintText: 'Enter your address',
                  keyboardType: TextInputType.streetAddress,
                  prefixIcon: Icons.location_on_outlined,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return AppConstants.requiredFieldMessage;
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 20),
                
                // Phone Field
                CustomTextField(
                  controller: _phoneController,
                  labelText: 'Phone Number',
                  hintText: 'Enter your phone number',
                  keyboardType: TextInputType.phone,
                  prefixIcon: Icons.phone_outlined,
                  inputFormatters: [
                    FilteringTextInputFormatter.digitsOnly,
                  ],
                  validator: _validatePhone,
                ),
                
                const SizedBox(height: 20),
                
                // Email Field
                CustomTextField(
                  controller: _emailController,
                  labelText: 'Email',
                  hintText: 'Enter your email',
                  keyboardType: TextInputType.emailAddress,
                  prefixIcon: Icons.email_outlined,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return AppConstants.requiredFieldMessage;
                    }
                    if (!RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(value)) {
                      return AppConstants.invalidEmailMessage;
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 20),
                
                // Password Field
                CustomTextField(
                  controller: _passwordController,
                  labelText: 'Password',
                  hintText: 'Create a password',
                  obscureText: !_isPasswordVisible,
                  prefixIcon: Icons.lock_outlined,
                  suffixIcon: IconButton(
                    icon: Icon(
                      _isPasswordVisible ? Icons.visibility : Icons.visibility_off,
                      color: AppTheme.gray,
                    ),
                    onPressed: () {
                      setState(() {
                        _isPasswordVisible = !_isPasswordVisible;
                      });
                    },
                  ),
                  validator: _validatePassword,
                ),
                
                // Password requirements helper text
                Padding(
                  padding: const EdgeInsets.only(top: 8.0, left: 16.0),
                  child: Text(
                    'Password must contain:\n• One uppercase letter\n• One lowercase letter\n• One number\n• One symbol (@\$!%*?&)',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: AppTheme.textSecondary,
                      fontSize: 12,
                    ),
                  ),
                ),
                
                const SizedBox(height: 20),
                
                // Confirm Password Field
                CustomTextField(
                  controller: _confirmPasswordController,
                  labelText: 'Confirm Password',
                  hintText: 'Confirm your password',
                  obscureText: !_isConfirmPasswordVisible,
                  prefixIcon: Icons.lock_outlined,
                  suffixIcon: IconButton(
                    icon: Icon(
                      _isConfirmPasswordVisible ? Icons.visibility : Icons.visibility_off,
                      color: AppTheme.gray,
                    ),
                    onPressed: () {
                      setState(() {
                        _isConfirmPasswordVisible = !_isConfirmPasswordVisible;
                      });
                    },
                  ),
                  validator: _validateConfirmPassword,
                ),
                
                const SizedBox(height: 32),
                
                // Register as Customer Button
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: _isLoading ? null : () => _handleRegister(UserRole.customer),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppTheme.primaryRed,
                      foregroundColor: AppTheme.white,
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: _isLoading
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor: AlwaysStoppedAnimation<Color>(AppTheme.white),
                            ),
                          )
                        : const Text(
                            'Register as Customer',
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                  ),
                ),
                
                const SizedBox(height: 12),
                
                // Register as Agent Button
                SizedBox(
                  width: double.infinity,
                  child: OutlinedButton(
                    onPressed: _isLoading ? null : () => _handleRegister(UserRole.provider),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppTheme.primaryRed,
                      side: const BorderSide(color: AppTheme.primaryRed, width: 2),
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: const Text(
                      'Register as Agent',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
                
                const SizedBox(height: 24),
                
                // Cancel Button
                SizedBox(
                  width: double.infinity,
                  child: TextButton(
                    onPressed: () {
                      context.go('/landing');
                    },
                    child: const Text(
                      'Cancel',
                      style: TextStyle(
                        color: AppTheme.gray,
                        fontSize: 16,
                      ),
                    ),
                  ),
                ),
                
                const SizedBox(height: 24),
                
                // Error Message
                if (authState is Error)
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: AppTheme.error.withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(color: AppTheme.error.withValues(alpha: 0.3)),
                    ),
                    child: Text(
                      (authState as Error).message,
                      style: const TextStyle(
                        color: AppTheme.error,
                        fontSize: 14,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }
} 