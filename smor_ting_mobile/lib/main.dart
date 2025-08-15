import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';

import 'core/theme/app_theme.dart';
import 'core/constants/app_constants.dart';
import 'features/navigation/presentation/app_router.dart';
import 'core/services/enhanced_auth_service.dart';
import 'core/models/enhanced_auth_models.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize Hive for local storage
  await Hive.initFlutter();
  
  // Add this flag to prevent multiple initialization attempts
  bool isInitializing = false;
  
  runApp(ProviderScope(
    child: Consumer(
      builder: (context, ref, child) {
        return FutureBuilder<EnhancedAuthResult?>(
          future: isInitializing ? null : (() {
            isInitializing = true;
            return ref.read(enhancedAuthServiceProvider).restoreSession();
          })(),
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return MaterialApp(
                home: Scaffold(
                  body: Center(child: CircularProgressIndicator()),
                ),
              );
            }
            
            return SmorTingApp();
          },
        );
      },
    ),
  ));
}

class SmorTingApp extends ConsumerWidget {
  const SmorTingApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);
    
    return MaterialApp.router(
      title: AppConstants.appName,
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      routerConfig: router,
      builder: (context, child) {
        return MediaQuery(
          data: MediaQuery.of(context).copyWith(textScaler: const TextScaler.linear(1.0)),
          child: child!,
        );
      },
    );
  }
}
