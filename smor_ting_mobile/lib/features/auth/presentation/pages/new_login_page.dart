import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:local_auth/local_auth.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/models/user.dart';
import '../../../../core/constants/app_constants.dart';
import '../../../../core/services/enhanced_auth_service.dart';
import '../../../../core/services/message_service.dart';
import '../providers/auth_provider.dart';
import '../providers/enhanced_auth_provider.dart';
import '../widgets/custom_text_field.dart';
import '../widgets/enhanced_login_outcomes.dart';

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
  bool _biometricAvailable = false;
  bool _biometricEnabled = false;
  
  // Enhanced login outcome states
  bool _showCooldown = false;
  bool _showCaptcha = false;
  bool _showTwoFactor = false;
  int _cooldownSeconds = 0;
  int _remainingAttempts = 0;

  @override
  void initState() {
    super.initState();
    // Clear any existing auth errors when the page loads
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(authNotifierProvider.notifier).clearError();
      _checkBiometricAvailability();
    });
  }

  Future<void> _checkBiometricAvailability() async {
    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final isAvailable = await authService.canUseBiometrics();
      final availableBiometrics = await authService.getAvailableBiometrics();
      
      if (mounted) {
        setState(() {
          _biometricAvailable = isAvailable && availableBiometrics.isNotEmpty;
        });
        
        // Check if any user has biometric enabled (simplified check)
        if (_biometricAvailable) {
          // In a real app, check if current user has biometric enabled
          const userEmail = 'user@example.com'; // TODO: Get from session or input
          final isEnabled = await authService.isBiometricEnabled(userEmail);
          setState(() {
            _biometricEnabled = isEnabled;
          });
        }
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _biometricAvailable = false;
          _biometricEnabled = false;
        });
      }
    }
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
      _resetEnhancedStates();
    });

    try {
      final enhancedAuthNotifier = ref.read(enhancedAuthNotifierProvider.notifier);
      await enhancedAuthNotifier.enhancedLogin(
        email: _usernameController.text.trim(),
        password: _passwordController.text,
        rememberMe: true,
      );
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  void _resetEnhancedStates() {
    setState(() {
      _showCooldown = false;
      _showCaptcha = false;
      _showTwoFactor = false;
      _cooldownSeconds = 0;
      _remainingAttempts = 0;
    });
  }

  void _handleCaptchaCompleted(String token) async {
    setState(() {
      _isLoading = true;
    });

    try {
      final enhancedAuthNotifier = ref.read(enhancedAuthNotifierProvider.notifier);
      await enhancedAuthNotifier.enhancedLogin(
        email: _usernameController.text.trim(),
        password: _passwordController.text,
        rememberMe: true,
        captchaToken: token,
      );
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  void _handleTwoFactorSubmitted(String code) async {
    setState(() {
      _isLoading = true;
    });

    try {
      final enhancedAuthNotifier = ref.read(enhancedAuthNotifierProvider.notifier);
      await enhancedAuthNotifier.enhancedLogin(
        email: _usernameController.text.trim(),
        password: _passwordController.text,
        rememberMe: true,
        twoFactorCode: code,
      );
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _handleBiometricLogin() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      const userEmail = 'user@example.com'; // TODO: Get from session or input
      
      final result = await authService.authenticateWithBiometrics(userEmail);
      
      if (result.success && result.user != null) {
        // Update auth state
        ref.read(authNotifierProvider.notifier).setAuthenticatedUser(
          result.user!,
          result.accessToken!,
        );
        
        // Navigate to appropriate page based on user role
        if (mounted) {
          switch (result.user!.role) {
            case UserRole.customer:
              context.go('/home');
              break;
            case UserRole.provider:
              context.go('/agent-home');
              break;
            case UserRole.admin:
              context.go('/admin-home');
              break;
          }
        }
      } else {
        if (mounted) {
          MessageService.showError(
            context,
            message: result.message ?? 'Biometric authentication failed',
            canRetry: true,
            onRetry: _handleBiometricLogin,
          );
        }
      }
    } catch (e) {
      if (mounted) {
        MessageService.showError(
          context,
          message: 'Biometric authentication failed',
          canRetry: true,
          onRetry: _handleBiometricLogin,
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
    final enhancedAuthState = ref.watch(enhancedAuthNotifierProvider);

    // Listen to enhanced auth state changes
    ref.listen<EnhancedAuthState>(enhancedAuthNotifierProvider, (previous, next) {
      next.when(
        initial: () {},
        loading: () {},
        unauthenticated: () {},
        authenticated: (user, accessToken, sessionId, deviceTrusted, isRestoredSession, requiresVerification) {
          if (requiresVerification == true) {
            context.go('/verify-otp?email=${user.email}&fullName=${user.fullName}');
          } else {
            final role = user.role;
            if (role == UserRole.provider || role == UserRole.admin) {
              context.go('/agent-dashboard');
            } else {
              context.go('/home');
            }
          }
        },
        requiresTwoFactor: (email, tempUser, deviceTrusted) {
          setState(() {
            _showTwoFactor = true;
            _showCaptcha = false;
            _showCooldown = false;
          });
        },
        requiresCaptcha: (email, remainingAttempts, lockoutInfo) {
          setState(() {
            _showCaptcha = true;
            _showTwoFactor = false;
            _showCooldown = false;
            _remainingAttempts = remainingAttempts;
          });
        },
        lockedOut: (lockoutInfo, message) {
          setState(() {
            _showCooldown = true;
            _showCaptcha = false;
            _showTwoFactor = false;
            _cooldownSeconds = lockoutInfo.timeUntilUnlock ?? 300; // Default 5 minutes
          });
        },
        error: (message, canRetry) {
          MessageService.showError(
            context,
            message: message,
            canRetry: canRetry,
            onRetry: canRetry ? () {
              if (!_isLoading) {
                _handleLogin();
              }
            } : null,
          );
        },
        requiresVerification: (user, email) {
          context.go('/verify-otp?email=$email&fullName=${user.fullName}');
        },
      );
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
                Semantics(
                  label: 'login_email',
                  textField: true,
                  child: CustomTextField(
                    key: const ValueKey('login_email'),
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
                ),
                
                const SizedBox(height: 20),
                
                // Password Field
                Semantics(
                  label: 'login_password',
                  textField: true,
                  child: CustomTextField(
                    key: const ValueKey('login_password'),
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
                ),
                
                const SizedBox(height: 16),
                
                // Forgot Password Link
                Align(
                  alignment: Alignment.centerRight,
                  child: Semantics(
                    label: 'login_forgot_password',
                    button: true,
                    child: TextButton(
                      onPressed: () {
                        context.go('/forgot-password');
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
                ),
                
                const SizedBox(height: 32),

                // Enhanced Login Outcomes
                EnhancedLoginOutcomes(
                  showCooldown: _showCooldown,
                  showCaptcha: _showCaptcha,
                  showTwoFactor: _showTwoFactor,
                  cooldownSeconds: _cooldownSeconds,
                  remainingAttempts: _remainingAttempts,
                  onCaptchaCompleted: _handleCaptchaCompleted,
                  onTwoFactorSubmitted: _handleTwoFactorSubmitted,
                ),
                
                // Sign In Button
                SizedBox(
                  width: double.infinity,
                  child: Semantics(
                    label: 'login_submit',
                    button: true,
                    child: ElevatedButton(
                      key: const ValueKey('login_submit'),
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
                ),
                
                // Biometric Unlock Button
                if (_biometricAvailable && _biometricEnabled) ...[
                  const SizedBox(height: 16),
                  
                  Row(
                    children: [
                      Expanded(child: Divider(color: Colors.grey[300])),
                      const Padding(
                        padding: EdgeInsets.symmetric(horizontal: 16),
                        child: Text(
                          'or',
                          style: TextStyle(
                            color: AppTheme.gray,
                            fontSize: 14,
                          ),
                        ),
                      ),
                      Expanded(child: Divider(color: Colors.grey[300])),
                    ],
                  ),
                  
                  const SizedBox(height: 16),
                  
                  SizedBox(
                    width: double.infinity,
                    child: OutlinedButton.icon(
                      onPressed: _isLoading ? null : _handleBiometricLogin,
                      icon: const Icon(
                        Icons.fingerprint,
                        color: AppTheme.primaryRed,
                        size: 24,
                      ),
                      label: const Text(
                        'Unlock with Biometrics',
                        style: TextStyle(
                          color: AppTheme.primaryRed,
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        side: const BorderSide(color: AppTheme.primaryRed),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                    ),
                  ),
                ],
                
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
                
                // Register Link
                Align(
                  alignment: Alignment.center,
                  child: Semantics(
                    label: 'login_register_link',
                    button: true,
                    child: TextButton(
                      onPressed: () {
                        context.go('/register');
                      },
                      child: const Text(
                        "Don't have an account? Create one",
                        style: TextStyle(
                          color: AppTheme.primaryRed,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ),
                ),

                const SizedBox(height: 24),
              ],
            ),
          ),
        ),
      ),
    );
  }
} 