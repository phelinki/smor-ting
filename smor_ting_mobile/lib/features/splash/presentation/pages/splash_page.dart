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

class _SplashPageState extends ConsumerState<SplashPage> {
  bool _hasInitialized = false;



  @override
  void initState() {
    super.initState();
    _initializeApp();
  }

  Future<void> _initializeApp() async {
    if (_hasInitialized) return;
    _hasInitialized = true;
    
    print('🔵 SplashPage: Starting app initialization...');
    
    // Initialize auth state
    await ref.read(authNotifierProvider.notifier).initializeAuthState();
    
    print('🔵 SplashPage: Auth initialization complete');
    
    // Wait a moment to ensure state updates
    await Future.delayed(const Duration(milliseconds: 100));
    
    // The GoRouter redirect will handle navigation based on auth state
    // No manual navigation needed here
  }

  // Navigation methods removed - GoRouter handles all navigation based on auth state

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircularProgressIndicator(),
            SizedBox(height: 16),
            Text('Initializing app...'),
          ],
        ),
      ),
    );
  }
}