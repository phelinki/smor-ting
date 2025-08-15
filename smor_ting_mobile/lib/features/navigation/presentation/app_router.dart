import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../auth/presentation/providers/auth_provider.dart';
import '../../../core/models/user.dart';
import '../../auth/presentation/pages/landing_page.dart';
import '../../auth/presentation/pages/simple_login_page.dart';
import '../../auth/presentation/pages/new_register_page.dart';

import '../../auth/presentation/pages/forgot_password_page.dart';
import '../../auth/presentation/pages/reset_password_page.dart';
import '../../home/presentation/pages/home_page.dart';
import '../../splash/presentation/pages/splash_page.dart';
import '../../splash/presentation/pages/onboarding_page.dart';
import '../../services/presentation/pages/service_categories_page.dart';
import '../../services/presentation/pages/services_list_page.dart';
import '../../services/presentation/pages/service_listings_page.dart';
import '../../services/presentation/pages/provider_profile_page.dart';
import '../../profile/presentation/pages/profile_page.dart';
import '../../settings/presentation/pages/settings_page.dart';
import '../../help/presentation/pages/help_page.dart';
import '../../about/presentation/pages/about_page.dart';
import '../../booking/presentation/pages/booking_confirmation_page.dart';
import '../../tracking/presentation/pages/real_time_tracking_page.dart';
import '../../payment/presentation/pages/payment_methods_page.dart';
import '../../bookings/presentation/pages/bookings_history_page.dart';
import '../../messages/presentation/pages/messages_page.dart';
import '../../agent/presentation/pages/agent_login_page.dart';
import '../../agent/presentation/pages/agent_verification_page.dart';
import '../../agent/presentation/pages/agent_dashboard_page.dart';
import '../../kyc/presentation/pages/kyc_page.dart';
import '../../../core/models/service.dart';

final appRouterProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authNotifierProvider);
  
  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final isAuthenticated = authState is Authenticated;
      // While auth is in progress, do not redirect to avoid bouncing
      if (authState is Loading) {
        return null;
      }
      final isAuthRoute = state.matchedLocation == '/landing' || 
                         state.matchedLocation == '/login' || 
                         state.matchedLocation == '/register' ||
                         state.matchedLocation == '/onboarding' ||
                         state.matchedLocation == '/agent-login';
      
      // If not authenticated and not on auth route, redirect to landing
      if (!isAuthenticated && !isAuthRoute && state.matchedLocation != '/') {
        return '/landing';
      }
      
      // If authenticated and on auth route, redirect based on user role
      if (isAuthenticated && isAuthRoute) {
        final userRole = authState.user.role;
        if (userRole == UserRole.provider || userRole == UserRole.admin) {
          return '/agent-dashboard';
        }
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
        path: '/landing',
        builder: (context, state) => const LandingPage(),
      ),
      GoRoute(
        path: '/onboarding',
        builder: (context, state) => const OnboardingPage(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const SimpleLoginPage(),
      ),
      GoRoute(
        path: '/forgot-password',
        builder: (context, state) => const ForgotPasswordPage(),
      ),
      GoRoute(
        path: '/reset-password',
        builder: (context, state) {
          final email = state.uri.queryParameters['email'] ?? '';
          return ResetPasswordPage(email: email);
        },
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const NewRegisterPage(),
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
      GoRoute(
        path: '/booking-confirmation',
        builder: (context, state) {
          final serviceId = state.uri.queryParameters['serviceId'] ?? '';
          final providerId = state.uri.queryParameters['providerId'] ?? '';
          final serviceName = state.uri.queryParameters['serviceName'] ?? '';
          final providerName = state.uri.queryParameters['providerName'] ?? '';
          final price = double.tryParse(state.uri.queryParameters['price'] ?? '0') ?? 0.0;
          return BookingConfirmationPage(
            serviceId: serviceId,
            providerId: providerId,
            serviceName: serviceName,
            providerName: providerName,
            price: price,
          );
        },
      ),
      GoRoute(
        path: '/tracking',
        builder: (context, state) {
          final bookingId = state.uri.queryParameters['bookingId'] ?? '';
          final providerName = state.uri.queryParameters['providerName'] ?? '';
          final serviceName = state.uri.queryParameters['serviceName'] ?? '';
          final address = state.uri.queryParameters['address'] ?? '';
          return RealTimeTrackingPage(
            bookingId: bookingId,
            providerName: providerName,
            serviceName: serviceName,
            address: address,
          );
        },
      ),
      GoRoute(
        path: '/payment-methods',
        builder: (context, state) => const PaymentMethodsPage(),
      ),
      GoRoute(
        path: '/bookings-history',
        builder: (context, state) => const BookingsHistoryPage(),
      ),
      GoRoute(
        path: '/messages',
        builder: (context, state) => const MessagesPage(),
      ),
      // Agent Routes
      GoRoute(
        path: '/kyc',
        builder: (context, state) => const KycPage(),
      ),
      GoRoute(
        path: '/agent-login',
        builder: (context, state) => const AgentLoginPage(),
      ),
      GoRoute(
        path: '/agent-verification',
        builder: (context, state) => const AgentVerificationPage(),
      ),
      GoRoute(
        path: '/agent-dashboard',
        builder: (context, state) => const AgentDashboardPage(),
      ),
    ],
  );
}); 