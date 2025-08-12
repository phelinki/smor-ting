import '../../../../core/models/user.dart';
import '../entities/role_config.dart';

/// Service for detecting and managing user roles
class RoleDetectionService {
  /// Get the role of a user
  UserRole getUserRole(User user) {
    return user.role;
  }

  /// Get the dashboard route for a specific role
  String getDashboardRouteForRole(UserRole role) {
    final config = RoleConfig.getConfigForRole(role);
    return config.dashboardRoute;
  }

  /// Check if a role is authorized to access a specific route
  bool isRoleAuthorizedForRoute(UserRole role, String route) {
    final config = RoleConfig.getConfigForRole(role);
    return config.canAccessRoute(route);
  }

  /// Get available features for a role
  List<String> getAvailableFeatures(UserRole role) {
    final config = RoleConfig.getConfigForRole(role);
    return config.features;
  }

  /// Validate user access to a route
  UserAccessResult validateUserAccess(User user, String route) {
    // Check email verification requirement
    if (!user.isEmailVerified) {
      return UserAccessResult.denied(
        reason: 'Email verification required',
        redirectRoute: '/verify-otp',
        queryParameters: {
          'email': user.email,
          'fullName': user.fullName,
        },
      );
    }

    // Check role-based access
    final role = getUserRole(user);
    final isAuthorized = isRoleAuthorizedForRoute(role, route);

    if (!isAuthorized) {
      final dashboardRoute = getDashboardRouteForRole(role);
      return UserAccessResult.denied(
        reason: 'Insufficient permissions for this route',
        redirectRoute: dashboardRoute,
      );
    }

    return UserAccessResult.authorized();
  }

  /// Check if a role has a specific feature
  bool roleHasFeature(UserRole role, String feature) {
    final config = RoleConfig.getConfigForRole(role);
    return config.hasFeature(feature);
  }

  /// Get role priority (higher = more privileged)
  int getRolePriority(UserRole role) {
    final config = RoleConfig.getConfigForRole(role);
    return config.priority;
  }

  /// Check if role A has higher privilege than role B
  bool isRoleHigherThan(UserRole roleA, UserRole roleB) {
    return getRolePriority(roleA) > getRolePriority(roleB);
  }

  /// Get the role configuration for a user
  RoleConfig getRoleConfig(User user) {
    return RoleConfig.getConfigForRole(user.role);
  }

  /// Check if user needs KYC verification
  bool requiresKyc(User user) {
    final config = getRoleConfig(user);
    return config.requiresKyc;
  }

  /// Check if user needs email verification
  bool requiresEmailVerification(User user) {
    final config = getRoleConfig(user);
    return config.requiresEmailVerification;
  }

  /// Get unauthorized redirect route for a user
  String getUnauthorizedRedirectRoute(User user) {
    final config = getRoleConfig(user);
    return config.getUnauthorizedRedirectRoute();
  }

  /// Check if route is a dashboard route
  bool isDashboardRoute(String route) {
    const dashboardRoutes = [
      '/home',
      '/agent-dashboard',
      '/admin-dashboard',
    ];
    return dashboardRoutes.contains(route);
  }

  /// Get appropriate dashboard for user after verification
  String getPostVerificationDashboard(User user, {bool kycCompleted = false}) {
    final role = getUserRole(user);
    
    // For providers, check if they need KYC
    if (role == UserRole.provider && requiresKyc(user) && !kycCompleted) {
      return '/kyc';
    }
    
    return getDashboardRouteForRole(role);
  }

  /// Check if user can access admin features
  bool canAccessAdminFeatures(User user) {
    return user.role == UserRole.admin;
  }

  /// Check if user can access provider features
  bool canAccessProviderFeatures(User user) {
    return user.role == UserRole.provider || user.role == UserRole.admin;
  }

  /// Check if user can access customer features
  bool canAccessCustomerFeatures(User user) {
    return true; // All roles can access customer features
  }

  /// Get role display name
  String getRoleDisplayName(UserRole role) {
    switch (role) {
      case UserRole.customer:
        return 'Customer';
      case UserRole.provider:
        return 'Service Provider';
      case UserRole.admin:
        return 'Administrator';
    }
  }

  /// Get role description
  String getRoleDescription(UserRole role) {
    switch (role) {
      case UserRole.customer:
        return 'Browse and book services from providers';
      case UserRole.provider:
        return 'Offer services and manage bookings';
      case UserRole.admin:
        return 'Manage platform and all users';
    }
  }
}
