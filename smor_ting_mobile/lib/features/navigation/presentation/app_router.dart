import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../auth/presentation/providers/auth_provider.dart';
import '../../auth/presentation/pages/login_page.dart';
import '../../auth/presentation/pages/register_page.dart';
import '../../auth/presentation/pages/phone_verification_page.dart';
import '../../home/presentation/pages/home_page.dart';
import '../../splash/presentation/pages/splash_page.dart';
import '../../splash/presentation/pages/onboarding_page.dart';
import '../../services/presentation/pages/service_listings_page.dart';
import '../../services/presentation/pages/provider_profile_page.dart';

final appRouterProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authNotifierProvider);
  
  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final isAuthenticated = authState is _Authenticated;
      final isAuthRoute = state.matchedLocation == '/login' || 
                         state.matchedLocation == '/register' ||
                         state.matchedLocation == '/phone-verification' ||
                         state.matchedLocation == '/onboarding';
      
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
        path: '/onboarding',
        builder: (context, state) => const OnboardingPage(),
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
        path: '/phone-verification',
        builder: (context, state) => const PhoneVerificationPage(),
      ),
      GoRoute(
        path: '/home',
        builder: (context, state) => const HomePage(),
      ),
      GoRoute(
        path: '/service-listings/:category',
        builder: (context, state) {
          final category = state.pathParameters['category']!;
          return ServiceListingsPage(category: category);
        },
      ),
      GoRoute(
        path: '/provider-profile/:providerId',
        builder: (context, state) {
          final providerId = state.pathParameters['providerId']!;
          return ProviderProfilePage(providerId: providerId);
        },
      ),
    ],
  );
}); 