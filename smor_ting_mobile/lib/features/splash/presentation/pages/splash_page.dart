import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/models/user.dart';
import '../../../auth/presentation/providers/auth_provider.dart';



/// Splash screen with app initialization
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
    // Provide a small idle window for first frame to stabilize in tests/automation
    // and expose a stable semantics node for landing detection.
    Future.microtask(() {});
  }

  Future<void> _initializeApp() async {
    // Wait for animations to start
    await Future.delayed(const Duration(milliseconds: 500));
    
    try {
      print('🔵 SplashPage: Starting app initialization...');
      
      // Initialize auth state - this will check for stored tokens and restore session
      final authNotifier = ref.read(authNotifierProvider.notifier);
      await authNotifier.initializeAuthState();
      
      print('🔵 SplashPage: Auth initialization complete');
      
      // Give the app time to stabilize
      await Future.delayed(const Duration(milliseconds: 1000));
      
      // Force navigation based on auth state
      final authState = ref.read(authNotifierProvider);
      if (authState is Authenticated) {
        final userRole = authState.user.role;
        if (userRole == UserRole.provider || userRole == UserRole.admin) {
          context.go('/agent-dashboard');
        } else {
          context.go('/home');
        }
      } else {
        context.go('/landing');
      }
      
      print('🔵 SplashPage: Navigation triggered');
    } catch (e) {
      print('🔴 SplashPage: Error during initialization: $e');
      // On error, go to landing page
      context.go('/landing');
    }
  }

  // Navigation methods removed - GoRouter handles all navigation based on auth state

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

                  // Loading indicator
                  Semantics(
                    label: 'splash_loading',
                    child: CircularProgressIndicator(
                      valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                    ),
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