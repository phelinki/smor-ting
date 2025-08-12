import '../../../../core/models/user.dart';

/// Configuration for user roles and their capabilities
class RoleConfig {
  final UserRole role;
  final String dashboardRoute;
  final List<String> allowedRoutes;
  final List<String> features;
  final bool requiresEmailVerification;
  final bool requiresKyc;
  final int priority; // Higher number = higher privilege

  const RoleConfig({
    required this.role,
    required this.dashboardRoute,
    required this.allowedRoutes,
    required this.features,
    this.requiresEmailVerification = true,
    this.requiresKyc = false,
    required this.priority,
  });

  /// Get role configuration for a specific role
  static RoleConfig getConfigForRole(UserRole role) {
    switch (role) {
      case UserRole.customer:
        return const CustomerRoleConfig();
      case UserRole.provider:
        return const ProviderRoleConfig();
      case UserRole.admin:
        return const AdminRoleConfig();
    }
  }

  /// Check if this role can access a specific route
  bool canAccessRoute(String route) {
    // Remove query parameters for route matching
    final cleanRoute = route.split('?').first;
    
    // Check exact matches first
    if (allowedRoutes.contains(cleanRoute)) {
      return true;
    }
    
    // Check pattern matches (e.g., /services/:id)
    for (final allowedRoute in allowedRoutes) {
      if (_matchesRoutePattern(cleanRoute, allowedRoute)) {
        return true;
      }
    }
    
    return false;
  }

  /// Check if route pattern matches (supports :param syntax)
  bool _matchesRoutePattern(String actualRoute, String patternRoute) {
    if (!patternRoute.contains(':')) {
      return actualRoute == patternRoute;
    }
    
    final actualParts = actualRoute.split('/');
    final patternParts = patternRoute.split('/');
    
    if (actualParts.length != patternParts.length) {
      return false;
    }
    
    for (int i = 0; i < patternParts.length; i++) {
      final patternPart = patternParts[i];
      final actualPart = actualParts[i];
      
      // Skip parameter parts (start with :)
      if (patternPart.startsWith(':')) {
        continue;
      }
      
      // Exact match required for non-parameter parts
      if (patternPart != actualPart) {
        return false;
      }
    }
    
    return true;
  }

  /// Check if this role has a specific feature
  bool hasFeature(String feature) {
    return features.contains(feature);
  }

  /// Get the appropriate redirect route for unauthorized access
  String getUnauthorizedRedirectRoute() {
    return dashboardRoute;
  }
}

/// Customer role configuration
class CustomerRoleConfig extends RoleConfig {
  const CustomerRoleConfig()
      : super(
          role: UserRole.customer,
          dashboardRoute: '/home',
          allowedRoutes: const [
            '/home',
            '/services',
            '/services/:categoryId',
            '/service-listings/:category',
            '/provider-profile/:providerId',
            '/profile',
            '/settings',
            '/help',
            '/about',
            '/booking-confirmation',
            '/tracking',
            '/payment-methods',
            '/bookings-history',
            '/messages',
          ],
          features: const [
            'browse_services',
            'book_services',
            'track_bookings',
            'payment_wallet',
            'view_profile',
            'messaging',
            'rate_providers',
          ],
          requiresEmailVerification: true,
          requiresKyc: false,
          priority: 1,
        );
}

/// Provider role configuration
class ProviderRoleConfig extends RoleConfig {
  const ProviderRoleConfig()
      : super(
          role: UserRole.provider,
          dashboardRoute: '/agent-dashboard',
          allowedRoutes: const [
            '/agent-dashboard',
            '/kyc',
            '/profile',
            '/settings',
            '/help',
            '/about',
            '/messages',
            '/tracking',
          ],
          features: const [
            'manage_services',
            'view_earnings',
            'kyc_verification',
            'job_management',
            'provider_analytics',
            'availability_settings',
            'view_profile',
            'messaging',
          ],
          requiresEmailVerification: true,
          requiresKyc: true,
          priority: 2,
        );
}

/// Admin role configuration
class AdminRoleConfig extends RoleConfig {
  const AdminRoleConfig()
      : super(
          role: UserRole.admin,
          dashboardRoute: '/admin-dashboard',
          allowedRoutes: const [
            '/admin-dashboard',
            '/home',
            '/services',
            '/services/:categoryId',
            '/service-listings/:category',
            '/provider-profile/:providerId',
            '/agent-dashboard',
            '/kyc',
            '/profile',
            '/settings',
            '/help',
            '/about',
            '/booking-confirmation',
            '/tracking',
            '/payment-methods',
            '/bookings-history',
            '/messages',
          ],
          features: const [
            'admin_panel',
            'user_management',
            'system_analytics',
            'manage_services',
            'browse_services',
            'book_services',
            'track_bookings',
            'payment_wallet',
            'kyc_verification',
            'job_management',
            'provider_analytics',
            'view_profile',
            'messaging',
            'content_management',
            'system_settings',
          ],
          requiresEmailVerification: true,
          requiresKyc: false,
          priority: 10,
        );
}

/// Result of user access validation
class UserAccessResult {
  final bool isAuthorized;
  final String? redirectRoute;
  final String? denialReason;
  final Map<String, String>? queryParameters;

  const UserAccessResult({
    required this.isAuthorized,
    this.redirectRoute,
    this.denialReason,
    this.queryParameters,
  });

  factory UserAccessResult.authorized() {
    return const UserAccessResult(isAuthorized: true);
  }

  factory UserAccessResult.denied({
    required String reason,
    String? redirectRoute,
    Map<String, String>? queryParameters,
  }) {
    return UserAccessResult(
      isAuthorized: false,
      denialReason: reason,
      redirectRoute: redirectRoute,
      queryParameters: queryParameters,
    );
  }
}
