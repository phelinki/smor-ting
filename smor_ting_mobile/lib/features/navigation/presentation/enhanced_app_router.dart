import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/models/user.dart';
import '../../auth/presentation/providers/auth_provider.dart';
import '../domain/usecases/role_detection_service.dart';
import '../domain/usecases/navigation_flow_service.dart';

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

/// Enhanced app router with role-based access control and proper navigation flows
final enhancedAppRouterProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(enhancedAuthNotifierProvider);
  final roleDetectionService = RoleDetectionService();
  final navigationFlowService = NavigationFlowService();
  
  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final currentPath = state.matchedLocation;
      
      // Handle authentication-based redirects
      return _handleAuthRedirect(
        authState,
        currentPath,
        roleDetectionService,
        navigationFlowService,
      );
    },
    routes: [
      // Splash and onboarding routes
      GoRoute(
        path: '/',
        builder: (context, state) => const SplashPage(),
      ),
      GoRoute(
        path: '/onboarding',
        builder: (context, state) => const OnboardingPage(),
      ),
      
      // Public authentication routes
      GoRoute(
        path: '/landing',
        builder: (context, state) => const LandingPage(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const SimpleLoginPage(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const NewRegisterPage(),
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

      
      // Customer routes (requires customer role or higher)
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
          // Create mock category for now
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
      
      // Common protected routes (all authenticated users)
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
      
      // Provider routes (requires provider role or higher)
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
      GoRoute(
        path: '/kyc',
        builder: (context, state) => const KycPage(),
      ),
      
      // Admin routes (requires admin role)
      GoRoute(
        path: '/admin-dashboard',
        builder: (context, state) => const AdminDashboardPage(),
      ),
    ],
  );
});

/// Handle authentication-based redirects with role validation
String? _handleAuthRedirect(
  EnhancedAuthState authState,
  String currentPath,
  RoleDetectionService roleDetectionService,
  NavigationFlowService navigationFlowService,
) {
  // Skip redirect during loading to avoid bouncing
  final isLoading = authState.maybeWhen(
    loading: () => true,
    orElse: () => false,
  ) ?? false;
  if (isLoading) {
    return null;
  }
  
  // Define public routes that don't require authentication
  const publicRoutes = [
    '/',
    '/landing',
    '/login',
    '/register',
    '/forgot-password',
    '/reset-password',
    '/onboarding',
    '/agent-login',
  ];
  
  final isPublicRoute = publicRoutes.contains(currentPath) ||
      currentPath.startsWith('/reset-password');
      // EMAIL OTP REMOVED: /verify-otp is no longer a public route
  
  // Handle different enhanced authentication states
  return authState.when(
    initial: () {
      // Initial state - redirect to landing if not on public route
      if (!isPublicRoute) {
        return '/landing';
      }
      return null;
    },
    loading: () {
      // Authentication in progress - keep current route to avoid bouncing
      return null;
    },
    unauthenticated: () {
      // Unauthenticated user - redirect to landing if not on public route
      if (!isPublicRoute) {
        return '/landing';
      }
      return null;
    },
    authenticated: (user, accessToken, sessionId, deviceTrusted, isRestoredSession, requiresVerification) {
      // Check if user is on a public route after authentication
      if (isPublicRoute) {
        // Redirect to appropriate dashboard based on role
        return roleDetectionService.getDashboardRouteForRole(user.role);
      }
      
      // Check role-based access for protected routes
      final accessResult = roleDetectionService.validateUserAccess(user, currentPath);
      
      if (!accessResult.isAuthorized) {
        // Redirect to appropriate dashboard or handling page
        return accessResult.redirectRoute ?? 
               roleDetectionService.getDashboardRouteForRole(user.role);
      }
      
      return null;
    },
    requiresTwoFactor: (email, tempUser, deviceTrusted) {
      // 2FA DISABLED: Skip 2FA and go directly to dashboard
      return roleDetectionService.getDashboardRouteForRole(tempUser.role);
    },
    requiresCaptcha: (email, remainingAttempts, lockoutInfo) {
      // CAPTCHA DISABLED: Stay on current page (should not happen in dev)
      return null;
    },
    lockedOut: (lockoutInfo, message) {
      // LOCKOUT DISABLED: Stay on current page (should not happen in dev)
      return null;
    },
    error: (message, canRetry) {
      // On error, redirect to landing if not on public route
      if (!isPublicRoute) {
        return '/landing';
      }
      return null;
    },
                requiresVerification: (user, email) {
              // EMAIL VERIFICATION DISABLED: Skip verification and go to dashboard
              return roleDetectionService.getDashboardRouteForRole(user.role);
            },
  );
}

/// Custom dashboard page that redirects based on role (if needed)
class AdminDashboardPage extends StatefulWidget {
  const AdminDashboardPage({super.key});

  @override
  State<AdminDashboardPage> createState() => _AdminDashboardPageState();
}

class _AdminDashboardPageState extends State<AdminDashboardPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Admin Dashboard'),
        backgroundColor: Colors.red[700],
        foregroundColor: Colors.white,
      ),
      body: const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.admin_panel_settings,
              size: 100,
              color: Colors.red,
            ),
            SizedBox(height: 20),
            Text(
              'Admin Dashboard',
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
            SizedBox(height: 10),
            Text(
              'Welcome to the admin panel',
              style: TextStyle(
                fontSize: 16,
                color: Colors.grey,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
