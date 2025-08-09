import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/models/user.dart';
import '../../../../core/constants/app_constants.dart';
import '../providers/auth_provider.dart';
import '../widgets/custom_text_field.dart';

class NewLoginPage extends ConsumerStatefulWidget {
  const NewLoginPage({super.key});

  @override
  ConsumerState<NewLoginPage> createState() => _NewLoginPageState();
}

class _NewLoginPageState extends ConsumerState<NewLoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  
  bool _isPasswordVisible = false;
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
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
    });

    try {
      await ref.read(authNotifierProvider.notifier).login(
        _usernameController.text.trim(),
        _passwordController.text,
      );
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  String? _validateUsername(String? value) {
    if (value == null || value.isEmpty) {
      return AppConstants.requiredFieldMessage;
    }
    if (value.length < 6) {
      return 'Username must be at least 6 characters long.';
    }
    return null;
  }

  String? _validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return AppConstants.requiredFieldMessage;
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authNotifierProvider);

    // Listen to auth state changes
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
        title: const Text('Sign In'),
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
                  'Welcome back!',
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                
                const SizedBox(height: 8),
                
                Text(
                  'Sign in to your account to continue',
                  style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                    color: AppTheme.gray,
                  ),
                ),
                
                const SizedBox(height: 32),
                
                // Username Field
                CustomTextField(
                  controller: _usernameController,
                  labelText: 'Username or Email',
                  hintText: 'Enter your username or email',
                  keyboardType: TextInputType.text,
                  prefixIcon: Icons.person_outlined,
                  suffixIcon: Tooltip(
                    message: 'Username must be at least 6 characters. You can use your email address as your username.',
                    child: const Icon(
                      Icons.help_outline,
                      color: AppTheme.gray,
                      size: 20,
                    ),
                  ),
                  validator: _validateUsername,
                ),
                
                const SizedBox(height: 20),
                
                // Password Field
                CustomTextField(
                  controller: _passwordController,
                  labelText: 'Password',
                  hintText: 'Enter your password',
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
                
                const SizedBox(height: 16),
                
                // Forgot Password Link
                Align(
                  alignment: Alignment.centerRight,
                  child: TextButton(
                    onPressed: () {
                      // TODO: Navigate to forgot password page
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Forgot password feature coming soon!'),
                        ),
                      );
                    },
                    child: const Text(
                      'Forgot Password?',
                      style: TextStyle(
                        color: AppTheme.primaryRed,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
                
                const SizedBox(height: 32),
                
                // Sign In Button
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: _isLoading ? null : _handleLogin,
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
                            'Sign In',
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