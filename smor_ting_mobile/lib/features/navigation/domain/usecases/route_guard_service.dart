import '../../../../core/models/user.dart';
import '../../../auth/presentation/providers/auth_provider.dart';
import '../entities/route_guard_result.dart';
import 'role_detection_service.dart';

/// Service for guarding routes based on authentication and authorization
class RouteGuardService {
  final RoleDetectionService _roleDetectionService;

  RouteGuardService({RoleDetectionService? roleDetectionService})
      : _roleDetectionService = roleDetectionService ?? RoleDetectionService();

  /// List of public routes that don't require authentication
  static const List<String> _publicRoutes = [
    '/',
    '/landing',
    '/login',
    '/register',
    '/forgot-password',
    '/reset-password',
    '/onboarding',
    '/agent-login',
  ];

  /// Check if a route requires authentication
  bool isPublicRoute(String route) {
    // Remove query parameters for route matching
    final cleanRoute = route.split('?').first;
    return _publicRoutes.contains(cleanRoute);
  }

  /// Check route access based on authentication state
  RouteGuardResult checkRouteAccess({
    required String route,
    required AuthState authState,
  }) {
    final cleanRoute = route.split('?').first;

    // Allow access to public routes
    if (isPublicRoute(cleanRoute)) {
      return RouteGuardResult.allowed();
    }

    // Handle different authentication states
    switch (authState) {
      case Initial():
        return _handleUnauthenticatedAccess(cleanRoute);
      case Loading():
        return RouteGuardResult.pending(reason: 'Authentication in progress');
      case Authenticated():
        return _handleAuthenticatedAccess(cleanRoute, authState.user);
      case RequiresOTP():
        return _handleOTPRequiredAccess(cleanRoute, authState.email);
      case Error():
        return _handleErrorAccess(cleanRoute, authState.message);
      case EmailAlreadyExists():
        return _handleUnauthenticatedAccess(cleanRoute);
      case PasswordResetEmailSent():
        return _handleUnauthenticatedAccess(cleanRoute);
      case PasswordResetSuccess():
        return _handleUnauthenticatedAccess(cleanRoute);
    }
  }

  /// Handle access for unauthenticated users
  RouteGuardResult _handleUnauthenticatedAccess(String route) {
    return RouteGuardResult.denied(
      redirectRoute: '/landing',
      reason: 'Authentication required',
    );
  }

  /// Handle access for authenticated users
  RouteGuardResult _handleAuthenticatedAccess(String route, User user) {
    // Check email verification
    if (!user.isEmailVerified) {
      return RouteGuardResult.denied(
        redirectRoute: '/verify-otp',
        reason: 'Email verification required',
        queryParameters: {
          'email': user.email,
          'fullName': user.fullName,
        },
      );
    }

    // Check role-based access
    final accessResult = _roleDetectionService.validateUserAccess(user, route);
    
    if (!accessResult.isAuthorized) {
      return RouteGuardResult.denied(
        redirectRoute: accessResult.redirectRoute!,
        reason: accessResult.denialReason!,
        queryParameters: accessResult.queryParameters,
      );
    }

    return RouteGuardResult.allowed();
  }

  /// Handle access when OTP verification is required
  RouteGuardResult _handleOTPRequiredAccess(String route, String email) {
    return RouteGuardResult.denied(
      redirectRoute: '/verify-otp',
      reason: 'OTP verification required',
      queryParameters: {'email': email},
    );
  }

  /// Handle access when there's an authentication error
  RouteGuardResult _handleErrorAccess(String route, String message) {
    return RouteGuardResult.denied(
      redirectRoute: '/landing',
      reason: 'Authentication error: $message',
    );
  }

  /// Get redirect route based on authentication state
  String? getRedirectRouteForAuthState(AuthState authState, String targetRoute) {
    switch (authState) {
      case Initial():
        return '/landing';
      case Loading():
        return null; // Don't redirect during loading
      case Authenticated():
        return _getRedirectForAuthenticatedUser(authState.user, targetRoute);
      case RequiresOTP():
        return '/verify-otp';
      case Error():
        return '/landing';
      case EmailAlreadyExists():
        return '/landing';
      case PasswordResetEmailSent():
        return '/landing';
      case PasswordResetSuccess():
        return '/login';
    }
  }

  /// Get redirect route for authenticated user trying to access unauthorized route
  String? _getRedirectForAuthenticatedUser(User user, String targetRoute) {
    // If user is trying to access a route they're not authorized for
    final accessResult = _roleDetectionService.validateUserAccess(user, targetRoute);
    
    if (!accessResult.isAuthorized) {
      return accessResult.redirectRoute;
    }
    
    return null;
  }

  /// Check if user should be redirected from current route after login
  bool shouldRedirectAfterLogin(String currentRoute, User user) {
    // If user is on a public route after login, redirect to dashboard
    if (isPublicRoute(currentRoute)) {
      return true;
    }

    // If user is on a route they don't have access to, redirect
    final accessResult = _roleDetectionService.validateUserAccess(user, currentRoute);
    return !accessResult.isAuthorized;
  }

  /// Get appropriate redirect route after login
  String getPostLoginRedirectRoute(User user, {String? intendedRoute}) {
    // If there's an intended route and user can access it, go there
    if (intendedRoute != null && !isPublicRoute(intendedRoute)) {
      final accessResult = _roleDetectionService.validateUserAccess(user, intendedRoute);
      if (accessResult.isAuthorized) {
        return intendedRoute;
      }
    }

    // Otherwise, redirect to appropriate dashboard
    return _roleDetectionService.getDashboardRouteForRole(user.role);
  }

  /// Check if route requires specific role
  bool routeRequiresRole(String route, UserRole requiredRole) {
    final isAuthorized = _roleDetectionService.isRoleAuthorizedForRoute(requiredRole, route);
    return isAuthorized;
  }

  /// Get minimum role required for route
  UserRole? getMinimumRoleForRoute(String route) {
    // Check from lowest to highest privilege
    for (final role in [UserRole.customer, UserRole.provider, UserRole.admin]) {
      if (_roleDetectionService.isRoleAuthorizedForRoute(role, route)) {
        return role;
      }
    }
    return null; // No role can access this route
  }

  /// Check if route is role-specific (only accessible by specific roles)
  bool isRoleSpecificRoute(String route) {
    final customerAccess = _roleDetectionService.isRoleAuthorizedForRoute(UserRole.customer, route);
    final providerAccess = _roleDetectionService.isRoleAuthorizedForRoute(UserRole.provider, route);
    final adminAccess = _roleDetectionService.isRoleAuthorizedForRoute(UserRole.admin, route);

    // If not all roles can access it, it's role-specific
    return !(customerAccess && providerAccess && adminAccess);
  }

  /// Get accessible routes for a user
  List<String> getAccessibleRoutesForUser(User user) {
    final role = user.role;
    final config = _roleDetectionService.getRoleConfig(user);
    return config.allowedRoutes;
  }

  /// Check if user can access any dashboard
  bool canAccessAnyDashboard(User user) {
    final dashboards = ['/home', '/agent-dashboard', '/admin-dashboard'];
    for (final dashboard in dashboards) {
      final accessResult = _roleDetectionService.validateUserAccess(user, dashboard);
      if (accessResult.isAuthorized) {
        return true;
      }
    }
    return false;
  }

  /// Get denied access message for user and route
  String getDeniedAccessMessage(User user, String route) {
    final userRole = _roleDetectionService.getRoleDisplayName(user.role);
    return 'Access denied: $userRole users cannot access this page.';
  }
}
