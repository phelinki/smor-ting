import '../../../../core/models/user.dart';
import '../entities/navigation_result.dart';
import 'role_detection_service.dart';

/// Service for managing navigation flows throughout the app
class NavigationFlowService {
  final RoleDetectionService _roleDetectionService;

  NavigationFlowService({RoleDetectionService? roleDetectionService})
      : _roleDetectionService = roleDetectionService ?? RoleDetectionService();

  /// Get destination after successful login
  NavigationResult getPostLoginDestination(User user, {bool requiresKyc = false}) {
    // EMAIL OTP REMOVED: Skip email verification check
    // All users go directly to dashboard regardless of email verification status
    
    // Check KYC requirement for providers
    if (user.role == UserRole.provider && requiresKyc) {
      return NavigationResult.replace('/kyc', clearHistory: true);
    }

    // Navigate to appropriate dashboard
    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.replace(dashboardRoute, clearHistory: true);
  }

  /// Get destination after successful registration
  NavigationResult getPostRegistrationDestination(User user) {
    // EMAIL OTP REMOVED: Skip email verification completely
    // All users go directly to dashboard after registration
    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.replace(dashboardRoute, clearHistory: true);
  }

  /// Get destination after OTP verification
  NavigationResult getPostOTPVerificationDestination(User user, {bool requiresKyc = false}) {
    // For providers who need KYC
    if (user.role == UserRole.provider && requiresKyc) {
      return NavigationResult.replace('/kyc', clearHistory: true);
    }

    // Navigate to appropriate dashboard
    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.replace(dashboardRoute, clearHistory: true);
  }

  /// Get destination after KYC completion
  NavigationResult getPostKYCDestination(User user, {required bool kycSuccess}) {
    if (!kycSuccess) {
      return NavigationResult.withMessage(
        destination: '/kyc',
        message: 'KYC verification failed. Please try again.',
        shouldReplace: false,
      );
    }

    // KYC successful - go to provider dashboard
    return NavigationResult.replace('/agent-dashboard', clearHistory: true);
  }

  /// Get destination after logout
  NavigationResult getLogoutDestination() {
    return NavigationResult.replace('/landing', clearHistory: true);
  }

  /// Handle deep link navigation
  NavigationResult handleDeepLink(String deepLink, User? user) {
    // If user is not authenticated, redirect to login with deep link saved
    if (user == null) {
      return NavigationResult.withParameters(
        destination: '/landing',
        queryParameters: {'redirect': deepLink},
        shouldReplace: true,
        message: 'Please log in to access this page.',
      );
    }

    // Check if user can access the deep link
    final accessResult = _roleDetectionService.validateUserAccess(user, deepLink);
    
    if (accessResult.isAuthorized) {
      return NavigationResult.go(deepLink);
    }

    // User can't access the deep link, redirect to their dashboard
    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.withMessage(
      destination: dashboardRoute,
      message: 'Access denied. Redirected to your dashboard.',
      shouldReplace: true,
    );
  }

  /// Check if onboarding should be shown
  bool shouldShowOnboarding(User user, {required bool isFirstLogin}) {
    // Don't show onboarding for admin users
    if (user.role == UserRole.admin) {
      return false;
    }

    // Show onboarding for first-time users
    return isFirstLogin;
  }

  /// Get onboarding completion destination
  NavigationResult getPostOnboardingDestination(User user) {
    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.replace(dashboardRoute, clearHistory: true);
  }

  /// Get navigation result for role-based redirection
  NavigationResult getRoleBasedRedirection(User user, String currentRoute) {
    // If user is on a public route, redirect to dashboard
    const publicRoutes = ['/', '/landing', '/login', '/register'];
    if (publicRoutes.contains(currentRoute)) {
      final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
      return NavigationResult.replace(dashboardRoute);
    }

    // Check if user can access current route
    final accessResult = _roleDetectionService.validateUserAccess(user, currentRoute);
    
    if (!accessResult.isAuthorized) {
      return NavigationResult.replace(accessResult.redirectRoute!);
    }

    // User can stay on current route
    return NavigationResult.go(currentRoute);
  }

  /// Handle session restoration navigation
  NavigationResult handleSessionRestoration(User user, String? lastRoute) {
    // If there's a last route and user can access it, go there
    if (lastRoute != null && lastRoute != '/') {
      final accessResult = _roleDetectionService.validateUserAccess(user, lastRoute);
      if (accessResult.isAuthorized) {
        return NavigationResult.go(lastRoute);
      }
    }

    // Otherwise, go to dashboard
    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.replace(dashboardRoute);
  }

  /// Get navigation for password reset completion
  NavigationResult getPostPasswordResetDestination() {
    return NavigationResult.withMessage(
      destination: '/login',
      message: 'Password reset successful. Please log in with your new password.',
      shouldReplace: true,
      clearHistory: true,
    );
  }

  /// Get navigation for account verification completion
  NavigationResult getPostAccountVerificationDestination(User user) {
    // Check if provider needs KYC
    if (user.role == UserRole.provider && _roleDetectionService.requiresKyc(user)) {
      return NavigationResult.replace('/kyc', clearHistory: true);
    }

    // Go to appropriate dashboard
    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.replace(dashboardRoute, clearHistory: true);
  }

  /// Handle navigation when user role changes
  NavigationResult handleRoleChangeNavigation(User user, UserRole previousRole) {
    // If role changed, redirect to new role's dashboard
    if (user.role != previousRole) {
      final newDashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
      return NavigationResult.withMessage(
        destination: newDashboardRoute,
        message: 'Your role has been updated. Welcome to your new dashboard!',
        shouldReplace: true,
        clearHistory: true,
      );
    }

    // No change needed
    final currentDashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.go(currentDashboardRoute);
  }

  /// Get navigation for unauthorized access attempts
  NavigationResult handleUnauthorizedAccess(User user, String attemptedRoute) {
    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    final userRoleName = _roleDetectionService.getRoleDisplayName(user.role);
    
    return NavigationResult.withMessage(
      destination: dashboardRoute,
      message: 'Access denied: $userRoleName users cannot access the requested page.',
      shouldReplace: true,
    );
  }

  /// Check if user should be forced to complete profile
  bool shouldCompleteProfile(User user) {
    // Check if essential profile fields are missing
    return user.firstName.isEmpty || 
           user.lastName.isEmpty || 
           user.phone.isEmpty;
  }

  /// Get navigation for incomplete profile
  NavigationResult getCompleteProfileDestination(User user) {
    return NavigationResult.withMessage(
      destination: '/profile',
      message: 'Please complete your profile to continue.',
      shouldReplace: true,
    );
  }

  /// Get initial route based on app state
  String getInitialRoute() {
    return '/';
  }

  /// Check if route is a modal/overlay route
  bool isModalRoute(String route) {
    const modalRoutes = [
      '/verify-otp',
      '/forgot-password',
      '/reset-password',
    ];
    return modalRoutes.contains(route);
  }

  /// Get navigation for error recovery
  NavigationResult getErrorRecoveryNavigation(User? user, String errorType) {
    if (user == null) {
      return NavigationResult.withMessage(
        destination: '/landing',
        message: 'An error occurred. Please log in again.',
        shouldReplace: true,
        clearHistory: true,
      );
    }

    final dashboardRoute = _roleDetectionService.getDashboardRouteForRole(user.role);
    return NavigationResult.withMessage(
      destination: dashboardRoute,
      message: 'An error occurred. Redirected to your dashboard.',
      shouldReplace: true,
    );
  }
}
