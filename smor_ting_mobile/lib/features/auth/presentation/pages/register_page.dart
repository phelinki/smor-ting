import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/constants/app_constants.dart';
import '../../../../core/models/user.dart';
import '../providers/auth_provider.dart';
import '../widgets/custom_text_field.dart';
import 'otp_verification_page.dart';

class RegisterPage extends ConsumerStatefulWidget {
  const RegisterPage({super.key});

  @override
  ConsumerState<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends ConsumerState<RegisterPage> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  
  bool _isPasswordVisible = false;
  bool _isConfirmPasswordVisible = false;
  bool _isLoading = false;
  String _selectedUserType = AppConstants.userTypeCustomer;

  @override
  void dispose() {
    _nameController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  Future<void> _handleRegister() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
    });

    try {
      // Split the full name into first and last name
      final fullName = _nameController.text.trim();
      final nameParts = fullName.split(' ');
      final firstName = nameParts.first;
      final lastName = nameParts.length > 1 ? nameParts.sublist(1).join(' ') : '';
      
      await ref.read(authNotifierProvider.notifier).register(
        firstName: firstName,
        lastName: lastName,
        email: _emailController.text.trim(),
        phone: _phoneController.text.trim(),
        password: _passwordController.text,
        role: _selectedUserType == AppConstants.userTypeAgent ? UserRole.provider : UserRole.customer,
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

    // Listen to auth state changes
    ref.listen<AuthState>(authNotifierProvider, (previous, next) {
      if (next is Authenticated) {
        context.go('/home');
      } else if (next is RequiresOTP) {
        Navigator.of(context).push(
          MaterialPageRoute(
            builder: (context) => OTPVerificationPage(
              email: next.email,
              userFullName: next.user.fullName,
            ),
          ),
        );
      }
    });

    return Scaffold(
      backgroundColor: AppTheme.white,
      appBar: AppBar(
        title: const Text('Create Account'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.go('/login'),
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(AppConstants.defaultPadding),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const SizedBox(height: 20),
                
                // Title
                Text(
                  'Join Smor-Ting',
                  style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppTheme.darkGray,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 8),
                Text(
                  'Create your account to get started',
                  style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                    color: AppTheme.gray,
                  ),
                  textAlign: TextAlign.center,
                ),
                
                const SizedBox(height: 32),
                
                // User Type Selection
                Text(
                  'I am a:',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppTheme.darkGray,
                  ),
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Expanded(
                      child: _buildUserTypeCard(
                        title: 'Customer',
                        subtitle: 'Looking for services',
                        icon: Icons.person,
                        isSelected: _selectedUserType == AppConstants.userTypeCustomer,
                        onTap: () => setState(() {
                          _selectedUserType = AppConstants.userTypeCustomer;
                        }),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: _buildUserTypeCard(
                        title: 'Agent',
                        subtitle: 'Offering services',
                        icon: Icons.build,
                        isSelected: _selectedUserType == AppConstants.userTypeAgent,
                        onTap: () => setState(() {
                          _selectedUserType = AppConstants.userTypeAgent;
                        }),
                      ),
                    ),
                  ],
                ),
                
                const SizedBox(height: 32),
                
                // Name Field
                CustomTextField(
                  controller: _nameController,
                  labelText: 'Full Name',
                  hintText: 'Enter your full name',
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
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return AppConstants.requiredFieldMessage;
                    }
                    if (value.length < 8) {
                      return AppConstants.passwordTooShortMessage;
                    }
                    return null;
                  },
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
                
                // Register Button
                ElevatedButton(
                  onPressed: _isLoading ? null : _handleRegister,
                  child: _isLoading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            valueColor: AlwaysStoppedAnimation<Color>(AppTheme.white),
                          ),
                        )
                      : const Text('Create Account'),
                ),
                
                const SizedBox(height: 24),
                
                // Login Link
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Text(
                      'Already have an account? ',
                      style: TextStyle(color: AppTheme.gray),
                    ),
                    TextButton(
                      onPressed: () {
                        context.go('/login');
                      },
                      child: const Text(
                        'Sign In',
                        style: TextStyle(
                          color: AppTheme.primaryRed,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ],
                ),
                
                const SizedBox(height: 40),
                
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

  Widget _buildUserTypeCard({
    required String title,
    required String subtitle,
    required IconData icon,
    required bool isSelected,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected ? AppTheme.primaryRed.withValues(alpha: 0.1) : AppTheme.lightGray,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isSelected ? AppTheme.primaryRed : Colors.transparent,
            width: 2,
          ),
        ),
        child: Column(
          children: [
            Icon(
              icon,
              color: isSelected ? AppTheme.primaryRed : AppTheme.gray,
              size: 32,
            ),
            const SizedBox(height: 8),
            Text(
              title,
              style: TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w600,
                color: isSelected ? AppTheme.primaryRed : AppTheme.darkGray,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 4),
            Text(
              subtitle,
              style: TextStyle(
                fontSize: 12,
                color: isSelected ? AppTheme.primaryRed : AppTheme.gray,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }
} 