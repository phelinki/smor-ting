/// Result of route guard check
class RouteGuardResult {
  final bool isAllowed;
  final String? redirectRoute;
  final String? denialReason;
  final Map<String, String>? queryParameters;
  final bool clearHistory;

  const RouteGuardResult({
    required this.isAllowed,
    this.redirectRoute,
    this.denialReason,
    this.queryParameters,
    this.clearHistory = false,
  });

  /// Factory for allowed access
  factory RouteGuardResult.allowed() {
    return const RouteGuardResult(isAllowed: true);
  }

  /// Factory for denied access with redirect
  factory RouteGuardResult.denied({
    required String redirectRoute,
    required String reason,
    Map<String, String>? queryParameters,
    bool clearHistory = false,
  }) {
    return RouteGuardResult(
      isAllowed: false,
      redirectRoute: redirectRoute,
      denialReason: reason,
      queryParameters: queryParameters,
      clearHistory: clearHistory,
    );
  }

  /// Factory for temporary denial (during loading, etc.)
  factory RouteGuardResult.pending({String? reason}) {
    return RouteGuardResult(
      isAllowed: false,
      denialReason: reason ?? 'Access check in progress',
    );
  }

  @override
  String toString() {
    return 'RouteGuardResult(isAllowed: $isAllowed, redirectRoute: $redirectRoute, denialReason: $denialReason)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is RouteGuardResult &&
        other.isAllowed == isAllowed &&
        other.redirectRoute == redirectRoute &&
        other.denialReason == denialReason &&
        other.clearHistory == clearHistory;
  }

  @override
  int get hashCode {
    return isAllowed.hashCode ^
        redirectRoute.hashCode ^
        denialReason.hashCode ^
        clearHistory.hashCode;
  }
}
