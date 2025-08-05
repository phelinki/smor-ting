import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/constants/app_constants.dart';
import '../providers/auth_provider.dart';
import '../widgets/custom_text_field.dart';
import 'otp_verification_page.dart';

class LoginPage extends ConsumerStatefulWidget {
  const LoginPage({super.key});

  @override
  ConsumerState<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends ConsumerState<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _phoneController = TextEditingController();
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();
  final _addressController = TextEditingController();
  bool _isLoading = false;
  bool _showEmailLogin = false;

  @override
  void dispose() {
    _phoneController.dispose();
    _firstNameController.dispose();
    _lastNameController.dispose();
    _addressController.dispose();
    super.dispose();
  }

  Future<void> _handlePhoneLogin() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
    });

    try {
      // TODO: Implement phone authentication
      await Future.delayed(const Duration(seconds: 2)); // Simulate API call
      if (mounted) {
        context.go('/phone-verification');
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _handleEmailLogin() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
    });

    try {
      await ref.read(authNotifierProvider.notifier).login(
        _phoneController.text.trim(),
        'password', // TODO: Add password field for email login
      );
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
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
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const SizedBox(height: 40),
                
                // Header
                Center(
                  child: Column(
                    children: [
                      Text(
                        'Welcome to Smor-Ting',
                        style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: AppTheme.secondaryBlue,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Connect with trusted service providers',
                        style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          color: AppTheme.textSecondary,
                        ),
                        textAlign: TextAlign.center,
                      ),
                    ],
                  ),
                ),
                
                const SizedBox(height: 48),
                
                if (!_showEmailLogin) ...[
                  // Phone Number Field
                  CustomTextField(
                    controller: _phoneController,
                    labelText: 'Phone Number',
                    hintText: 'Enter your phone number',
                    keyboardType: TextInputType.phone,
                    prefixIcon: Icons.phone_outlined,
                    inputFormatters: [
                      FilteringTextInputFormatter.digitsOnly,
                    ],
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return AppConstants.requiredFieldMessage;
                      }
                      if (value.length < 10) {
                        return 'Please enter a valid phone number';
                      }
                      return null;
                    },
                  ),
                  
                  const SizedBox(height: 20),
                  
                  // First Name Field
                  Container(
                    decoration: BoxDecoration(
                      color: AppTheme.primaryRed,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: TextFormField(
                      controller: _firstNameController,
                      style: const TextStyle(color: AppTheme.white),
                      decoration: const InputDecoration(
                        labelText: 'First Name',
                        labelStyle: TextStyle(color: AppTheme.white),
                        hintText: 'Enter your first name',
                        hintStyle: TextStyle(color: AppTheme.white),
                        prefixIcon: Icon(Icons.person_outline, color: AppTheme.white),
                        border: InputBorder.none,
                        contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'First name is required';
                        }
                        return null;
                      },
                    ),
                  ),
                  
                  const SizedBox(height: 20),
                  
                  // Last Name Field
                  Container(
                    decoration: BoxDecoration(
                      color: AppTheme.primaryRed,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: TextFormField(
                      controller: _lastNameController,
                      style: const TextStyle(color: AppTheme.white),
                      decoration: const InputDecoration(
                        labelText: 'Last Name',
                        labelStyle: TextStyle(color: AppTheme.white),
                        hintText: 'Enter your last name',
                        hintStyle: TextStyle(color: AppTheme.white),
                        prefixIcon: Icon(Icons.person_outline, color: AppTheme.white),
                        border: InputBorder.none,
                        contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Last name is required';
                        }
                        return null;
                      },
                    ),
                  ),
                  
                  const SizedBox(height: 20),
                  
                  // Address Field
                  Container(
                    decoration: BoxDecoration(
                      color: AppTheme.secondaryBlue,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: TextFormField(
                      controller: _addressController,
                      style: const TextStyle(color: AppTheme.white),
                      decoration: const InputDecoration(
                        labelText: 'Address',
                        labelStyle: TextStyle(color: AppTheme.white),
                        hintText: 'Enter your address',
                        hintStyle: TextStyle(color: AppTheme.white),
                        prefixIcon: Icon(Icons.location_on_outlined, color: AppTheme.white),
                        border: InputBorder.none,
                        contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Address is required';
                        }
                        return null;
                      },
                    ),
                  ),
                ] else ...[
                  // Email Field (for email login option)
                  CustomTextField(
                    controller: _phoneController,
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
                ],
                
                const SizedBox(height: 32),
                
                // Primary Login Button
                ElevatedButton(
                  onPressed: _isLoading ? null : (_showEmailLogin ? _handleEmailLogin : _handlePhoneLogin),
                  child: _isLoading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            valueColor: AlwaysStoppedAnimation<Color>(AppTheme.white),
                          ),
                        )
                      : Text(_showEmailLogin ? 'Continue with Email' : 'Continue with Phone'),
                ),
                
                const SizedBox(height: 16),
                
                // Secondary Login Button
                OutlinedButton(
                  onPressed: () {
                    setState(() {
                      _showEmailLogin = !_showEmailLogin;
                      _phoneController.clear();
                    });
                  },
                  child: Text(_showEmailLogin ? 'Continue with Phone' : 'Continue with Email'),
                ),
                
                const SizedBox(height: 32),
                
                // User Type Selection
                Row(
                  children: [
                    Expanded(
                      child: TextButton(
                        onPressed: () {
                          // TODO: Set user type to customer and continue
                        },
                        child: const Text(
                          'I\'m a Customer',
                          style: TextStyle(color: AppTheme.secondaryBlue),
                        ),
                      ),
                    ),
                    Container(
                      width: 1,
                      height: 20,
                      color: AppTheme.textSecondary,
                    ),
                    Expanded(
                      child: TextButton(
                        onPressed: () {
                          // TODO: Navigate to agent registration
                        },
                        child: const Text(
                          'I\'m a Service Provider',
                          style: TextStyle(color: AppTheme.secondaryBlue),
                        ),
                      ),
                    ),
                  ],
                ),
                
                const SizedBox(height: 24),
                
                // Terms and Privacy
                Text(
                  'By continuing, you agree to our Terms of Service and Privacy Policy',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: AppTheme.textSecondary,
                  ),
                  textAlign: TextAlign.center,
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
} 