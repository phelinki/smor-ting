import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/models/user.dart';
import '../../../navigation/domain/usecases/navigation_flow_service.dart';
import '../../../navigation/domain/usecases/role_detection_service.dart';
import 'auth_provider.dart';

/// Enhanced auth notifier that provides role-based navigation
final enhancedAuthNotifierProvider = Provider<EnhancedAuthNotifier>((ref) {
  final authNotifier = ref.read(authNotifierProvider.notifier);
  final navigationFlowService = NavigationFlowService();
  final roleDetectionService = RoleDetectionService();
  
  return EnhancedAuthNotifier(
    authNotifier: authNotifier,
    navigationFlowService: navigationFlowService,
    roleDetectionService: roleDetectionService,
  );
});

/// Enhanced authentication notifier with role-based navigation logic
class EnhancedAuthNotifier {
  final AuthNotifier _authNotifier;
  final NavigationFlowService _navigationFlowService;
  final RoleDetectionService _roleDetectionService;

  EnhancedAuthNotifier({
    required AuthNotifier authNotifier,
    required NavigationFlowService navigationFlowService,
    required RoleDetectionService roleDetectionService,
  })  : _authNotifier = authNotifier,
        _navigationFlowService = navigationFlowService,
        _roleDetectionService = roleDetectionService;

  /// Enhanced login with role-based navigation
  Future<void> loginWithNavigation({
    required String email,
    required String password,
    required GoRouter router,
  }) async {
    try {
      // Perform login
      await _authNotifier.login(email, password);
      
      // Handle post-login navigation based on auth state
      await _handlePostLoginNavigation(router);
    } catch (e) {
      // Error is already handled by the auth notifier
      rethrow;
    }
  }

  /// Enhanced registration with role-based navigation
  Future<void> registerWithNavigation({
    required String firstName,
    required String lastName,
    required String email,
    required String phone,
    required String password,
    UserRole role = UserRole.customer,
    required GoRouter router,
  }) async {
    try {
      // Perform registration
      await _authNotifier.register(
        firstName: firstName,
        lastName: lastName,
        email: email,
        phone: phone,
        password: password,
        role: role,
      );
      
      // Handle post-registration navigation
      await _handlePostRegistrationNavigation(router);
    } catch (e) {
      // Error is already handled by the auth notifier
      rethrow;
    }
  }

  /// Enhanced OTP verification with role-based navigation
  Future<void> verifyOTPWithNavigation({
    required String email,
    required String otp,
    required GoRouter router,
    bool requiresKyc = false,
  }) async {
    try {
      // Perform OTP verification
      await _authNotifier.verifyOTP(email, otp);
      
      // Handle post-OTP navigation
      await _handlePostOTPNavigation(router, requiresKyc: requiresKyc);
    } catch (e) {
      // Error is already handled by the auth notifier
      rethrow;
    }
  }

  /// Handle navigation after successful login
  Future<void> _handlePostLoginNavigation(GoRouter router) async {
    // Get the current auth state to access user info
    final authState = _authNotifier.state;
    
    if (authState is Authenticated) {
      final user = authState.user;
      
      // Determine if user needs KYC (for providers)
      final needsKyc = _roleDetectionService.requiresKyc(user) && 
                      !_isKycCompleted(user); // You'd implement this check
      
      // Get navigation destination
      final navigationResult = _navigationFlowService.getPostLoginDestination(
        user,
        requiresKyc: needsKyc,
      );
      
      // Navigate based on result
      if (navigationResult.clearHistory) {
        router.go(navigationResult.fullUri);
      } else if (navigationResult.shouldReplace) {
        router.pushReplacement(navigationResult.fullUri);
      } else {
        router.push(navigationResult.fullUri);
      }
    }
  }

  /// Handle navigation after successful registration
  Future<void> _handlePostRegistrationNavigation(GoRouter router) async {
    final authState = _authNotifier.state;
    
    if (authState is Authenticated) {
      final user = authState.user;
      
      final navigationResult = _navigationFlowService.getPostRegistrationDestination(user);
      
      if (navigationResult.clearHistory) {
        router.go(navigationResult.fullUri);
      } else if (navigationResult.shouldReplace) {
        router.pushReplacement(navigationResult.fullUri);
      } else {
        router.push(navigationResult.fullUri);
      }
    }
    // EMAIL OTP REMOVED: No longer handle RequiresOTP state
    // All users go directly to dashboard after registration
  }

  /// Handle navigation after successful OTP verification
  Future<void> _handlePostOTPNavigation(GoRouter router, {bool requiresKyc = false}) async {
    final authState = _authNotifier.state;
    
    if (authState is Authenticated) {
      final user = authState.user;
      
      final navigationResult = _navigationFlowService.getPostOTPVerificationDestination(
        user,
        requiresKyc: requiresKyc,
      );
      
      if (navigationResult.clearHistory) {
        router.go(navigationResult.fullUri);
      } else if (navigationResult.shouldReplace) {
        router.pushReplacement(navigationResult.fullUri);
      } else {
        router.push(navigationResult.fullUri);
      }
    }
  }

  /// Handle logout with navigation
  Future<void> logoutWithNavigation(GoRouter router) async {
    // Perform logout
    _authNotifier.logout();
    
    // Get logout navigation destination
    final navigationResult = _navigationFlowService.getLogoutDestination();
    
    // Navigate to landing page and clear history
    router.go(navigationResult.destination);
  }

  /// Check if user has access to a specific route
  bool canAccessRoute(User user, String route) {
    final accessResult = _roleDetectionService.validateUserAccess(user, route);
    return accessResult.isAuthorized;
  }

  /// Get appropriate dashboard for user
  String getDashboardForUser(User user) {
    return _roleDetectionService.getDashboardRouteForRole(user.role);
  }

  /// Handle deep link access
  Future<void> handleDeepLink({
    required String deepLink,
    required GoRouter router,
    User? user,
  }) async {
    final navigationResult = _navigationFlowService.handleDeepLink(deepLink, user);
    
    if (navigationResult.clearHistory) {
      router.go(navigationResult.fullUri);
    } else if (navigationResult.shouldReplace) {
      router.pushReplacement(navigationResult.fullUri);
    } else {
      router.push(navigationResult.fullUri);
    }
  }

  /// Check if user should see onboarding
  bool shouldShowOnboarding(User user, {required bool isFirstLogin}) {
    return _navigationFlowService.shouldShowOnboarding(user, isFirstLogin: isFirstLogin);
  }

  /// Handle role change navigation
  Future<void> handleRoleChange({
    required User user,
    required UserRole previousRole,
    required GoRouter router,
  }) async {
    final navigationResult = _navigationFlowService.handleRoleChangeNavigation(user, previousRole);
    
    if (navigationResult.clearHistory) {
      router.go(navigationResult.fullUri);
    } else if (navigationResult.shouldReplace) {
      router.pushReplacement(navigationResult.fullUri);
    } else {
      router.push(navigationResult.fullUri);
    }
  }

  /// Check if user profile is complete
  bool isProfileComplete(User user) {
    return !_navigationFlowService.shouldCompleteProfile(user);
  }

  /// Handle incomplete profile navigation
  Future<void> handleIncompleteProfile({
    required User user,
    required GoRouter router,
  }) async {
    if (_navigationFlowService.shouldCompleteProfile(user)) {
      final navigationResult = _navigationFlowService.getCompleteProfileDestination(user);
      
      if (navigationResult.shouldReplace) {
        router.pushReplacement(navigationResult.fullUri);
      } else {
        router.push(navigationResult.fullUri);
      }
    }
  }

  /// Get available features for user role
  List<String> getAvailableFeatures(User user) {
    return _roleDetectionService.getAvailableFeatures(user.role);
  }

  /// Check if user has specific feature
  bool hasFeature(User user, String feature) {
    return _roleDetectionService.roleHasFeature(user.role, feature);
  }

  /// Get role display information
  String getRoleDisplayName(User user) {
    return _roleDetectionService.getRoleDisplayName(user.role);
  }

  /// Get role description
  String getRoleDescription(User user) {
    return _roleDetectionService.getRoleDescription(user.role);
  }

  /// Placeholder for KYC completion check
  /// In a real app, this would check against the backend
  bool _isKycCompleted(User user) {
    // This is a placeholder - in reality you'd check with your backend
    // or store KYC status in the user model
    return false;
  }

  /// Handle session restoration with navigation
  Future<void> handleSessionRestoration({
    required User user,
    required GoRouter router,
    String? lastRoute,
  }) async {
    final navigationResult = _navigationFlowService.handleSessionRestoration(user, lastRoute);
    
    if (navigationResult.clearHistory) {
      router.go(navigationResult.fullUri);
    } else if (navigationResult.shouldReplace) {
      router.pushReplacement(navigationResult.fullUri);
    } else {
      router.push(navigationResult.fullUri);
    }
  }

  /// Handle unauthorized access attempt
  Future<void> handleUnauthorizedAccess({
    required User user,
    required String attemptedRoute,
    required GoRouter router,
  }) async {
    final navigationResult = _navigationFlowService.handleUnauthorizedAccess(user, attemptedRoute);
    
    if (navigationResult.shouldReplace) {
      router.pushReplacement(navigationResult.fullUri);
    } else {
      router.push(navigationResult.fullUri);
    }
  }
}
