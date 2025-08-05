import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../auth/presentation/providers/auth_provider.dart';
import '../../auth/presentation/pages/login_page.dart';
import '../../auth/presentation/pages/register_page.dart';
import '../../auth/presentation/pages/otp_verification_page.dart';
import '../../home/presentation/pages/home_page.dart';
import '../../splash/presentation/pages/splash_page.dart';
import '../../services/presentation/pages/service_categories_page.dart';
import '../../services/presentation/pages/services_list_page.dart';
import '../../profile/presentation/pages/profile_page.dart';
import '../../settings/presentation/pages/settings_page.dart';
import '../../help/presentation/pages/help_page.dart';
import '../../about/presentation/pages/about_page.dart';
import '../../../core/models/service.dart';

final appRouterProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authNotifierProvider);
  
  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final isAuthenticated = authState is _Authenticated;
      final isAuthRoute = state.matchedLocation == '/login' || 
                         state.matchedLocation == '/register' ||
                         state.matchedLocation == '/verify-otp';
      
      // If not authenticated and not on auth route, redirect to login
      if (!isAuthenticated && !isAuthRoute && state.matchedLocation != '/') {
        return '/login';
      }
      
      // If authenticated and on auth route, redirect to home
      if (isAuthenticated && isAuthRoute) {
        return '/home';
      }
      
      return null;
    },
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const SplashPage(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const RegisterPage(),
      ),
      GoRoute(
        path: '/verify-otp',
        builder: (context, state) {
          final email = state.uri.queryParameters['email'] ?? '';
          final fullName = state.uri.queryParameters['fullName'] ?? '';
          return OTPVerificationPage(
            email: email,
            userFullName: fullName,
          );
        },
      ),
      GoRoute(
        path: '/home',
        builder: (context, state) => const HomePage(),
      ),
      GoRoute(
        path: '/services',
        builder: (context, state) => const ServiceCategoriesPage(),
      ),
      GoRoute(
        path: '/services/:categoryId',
        builder: (context, state) {
          final categoryId = state.pathParameters['categoryId']!;
          // You'll need to pass the actual category object here
          // For now, creating a mock category
          final category = ServiceCategory(
            id: categoryId,
            name: 'Services',
            description: 'Service category',
            icon: 'build',
            color: '#2196F3',
            isActive: true,
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          );
          return ServicesListPage(category: category);
        },
      ),
      GoRoute(
        path: '/profile',
        builder: (context, state) => const ProfilePage(),
      ),
      GoRoute(
        path: '/settings',
        builder: (context, state) => const SettingsPage(),
      ),
      GoRoute(
        path: '/help',
        builder: (context, state) => const HelpPage(),
      ),
      GoRoute(
        path: '/about',
        builder: (context, state) => const AboutPage(),
      ),
    ],
  );
}); 