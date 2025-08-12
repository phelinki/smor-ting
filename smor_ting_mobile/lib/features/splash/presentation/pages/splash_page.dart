import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../auth/presentation/providers/enhanced_auth_provider.dart';
import '../../../auth/presentation/widgets/biometric_quick_unlock.dart';

/// Splash screen with session restoration and biometric quick unlock
class SplashPage extends ConsumerStatefulWidget {
  const SplashPage({super.key});

  @override
  ConsumerState<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends ConsumerState<SplashPage>
    with TickerProviderStateMixin {
  late AnimationController _logoAnimationController;
  late AnimationController _fadeAnimationController;
  late Animation<double> _logoAnimation;
  late Animation<double> _fadeAnimation;

  bool _showBiometricUnlock = false;
  String? _userEmail;

  @override
  void initState() {
    super.initState();
    _setupAnimations();
    _initializeApp();
  }

  void _setupAnimations() {
    _logoAnimationController = AnimationController(
      duration: const Duration(milliseconds: 1500),
      vsync: this,
    );

    _fadeAnimationController = AnimationController(
      duration: const Duration(milliseconds: 800),
      vsync: this,
    );

    _logoAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _logoAnimationController,
      curve: Curves.elasticOut,
    ));

    _fadeAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _fadeAnimationController,
      curve: Curves.easeIn,
    ));

    _logoAnimationController.forward();
  }

  Future<void> _initializeApp() async {
    // Give time for animations to start
    await Future.delayed(const Duration(milliseconds: 500));

    try {
      final authService = ref.read(enhancedAuthServiceProvider);
      final restoredSession = await authService.restoreSession();

      if (restoredSession != null && restoredSession.success && restoredSession.user != null) {
        final user = restoredSession.user!;
        _userEmail = user.email;

        // Check if biometric unlock is available and enabled
        final canUseBiometrics = await authService.canUseBiometrics();
        final isBiometricEnabled = await authService.isBiometricEnabled(user.email);

        if (canUseBiometrics && isBiometricEnabled) {
          // Show biometric unlock option
          setState(() {
            _showBiometricUnlock = true;
          });
          _fadeAnimationController.forward();
        } else {
          // Go directly to home
          _navigateToHome();
        }
      } else {
        // No session, go to landing/login
        _navigateToLanding();
      }
    } catch (e) {
      // Error restoring session, go to landing
      _navigateToLanding();
    }
  }

  void _navigateToHome() {
    if (mounted) {
      context.go('/home');
    }
  }

  void _navigateToLanding() {
    if (mounted) {
      context.go('/landing');
    }
  }

  void _navigateToLogin() {
    if (mounted) {
      context.go('/login');
    }
  }

  @override
  void dispose() {
    _logoAnimationController.dispose();
    _fadeAnimationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.secondaryBlue,
      body: SafeArea(
        child: Stack(
          children: [
            // Background gradient
            Container(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    AppTheme.secondaryBlue,
                    AppTheme.secondaryBlue.withOpacity(0.8),
                  ],
                ),
              ),
            ),

            // Main content
            Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  // Logo animation
                  AnimatedBuilder(
                    animation: _logoAnimation,
                    builder: (context, child) {
                      return Transform.scale(
                        scale: _logoAnimation.value,
                        child: Column(
                          children: [
                            Container(
                              width: 120,
                              height: 120,
                              decoration: BoxDecoration(
                                color: Colors.white,
                                borderRadius: BorderRadius.circular(24),
                                boxShadow: [
                                  BoxShadow(
                                    color: Colors.black.withOpacity(0.2),
                                    blurRadius: 20,
                                    offset: const Offset(0, 10),
                                  ),
                                ],
                              ),
                              child: Center(
                                child: Text(
                                  'ST',
                                  style: TextStyle(
                                    fontSize: 36,
                                    fontWeight: FontWeight.bold,
                                    color: AppTheme.secondaryBlue,
                                  ),
                                ),
                              ),
                            ),
                            const SizedBox(height: 24),
                            const Text(
                              'Smor-Ting',
                              style: TextStyle(
                                fontSize: 32,
                                fontWeight: FontWeight.bold,
                                color: Colors.white,
                              ),
                            ),
                            const SizedBox(height: 8),
                            Text(
                              'Service Marketplace',
                              style: TextStyle(
                                fontSize: 16,
                                color: Colors.white.withOpacity(0.8),
                              ),
                            ),
                          ],
                        ),
                      );
                    },
                  ),

                  const SizedBox(height: 60),

                  // Biometric unlock or loading indicator
                  if (_showBiometricUnlock && _userEmail != null)
                    FadeTransition(
                      opacity: _fadeAnimation,
                      child: BiometricQuickUnlock(
                        userEmail: _userEmail!,
                        onSuccess: _navigateToHome,
                        onCancel: _navigateToLogin,
                      ),
                    )
                  else
                    const CircularProgressIndicator(
                      valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                    ),
                ],
              ),
            ),

            // App version at bottom
            Positioned(
              bottom: 32,
              left: 0,
              right: 0,
              child: Center(
                child: Text(
                  'Version 1.0.0',
                  style: TextStyle(
                    color: Colors.white.withOpacity(0.6),
                    fontSize: 14,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}