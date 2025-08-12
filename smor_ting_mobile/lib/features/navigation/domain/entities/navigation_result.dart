/// Result of navigation operation
class NavigationResult {
  final String destination;
  final bool shouldReplace;
  final bool clearHistory;
  final Map<String, String> queryParameters;
  final String? message;
  final Map<String, dynamic>? extra;

  const NavigationResult({
    required this.destination,
    this.shouldReplace = false,
    this.clearHistory = false,
    this.queryParameters = const {},
    this.message,
    this.extra,
  });

  /// Factory for simple navigation
  factory NavigationResult.go(String destination) {
    return NavigationResult(destination: destination);
  }

  /// Factory for replace navigation
  factory NavigationResult.replace(String destination, {bool clearHistory = false}) {
    return NavigationResult(
      destination: destination,
      shouldReplace: true,
      clearHistory: clearHistory,
    );
  }

  /// Factory for navigation with parameters
  factory NavigationResult.withParameters({
    required String destination,
    Map<String, String> queryParameters = const {},
    bool shouldReplace = false,
    bool clearHistory = false,
    String? message,
  }) {
    return NavigationResult(
      destination: destination,
      queryParameters: queryParameters,
      shouldReplace: shouldReplace,
      clearHistory: clearHistory,
      message: message,
    );
  }

  /// Factory for navigation with message
  factory NavigationResult.withMessage({
    required String destination,
    required String message,
    bool shouldReplace = false,
    bool clearHistory = false,
  }) {
    return NavigationResult(
      destination: destination,
      message: message,
      shouldReplace: shouldReplace,
      clearHistory: clearHistory,
    );
  }

  /// Get the full URI with query parameters
  String get fullUri {
    if (queryParameters.isEmpty) {
      return destination;
    }
    
    final uri = Uri.parse(destination);
    final newUri = uri.replace(
      queryParameters: {
        ...uri.queryParameters,
        ...queryParameters,
      },
    );
    
    return newUri.toString();
  }

  /// Copy with new values
  NavigationResult copyWith({
    String? destination,
    bool? shouldReplace,
    bool? clearHistory,
    Map<String, String>? queryParameters,
    String? message,
    Map<String, dynamic>? extra,
  }) {
    return NavigationResult(
      destination: destination ?? this.destination,
      shouldReplace: shouldReplace ?? this.shouldReplace,
      clearHistory: clearHistory ?? this.clearHistory,
      queryParameters: queryParameters ?? this.queryParameters,
      message: message ?? this.message,
      extra: extra ?? this.extra,
    );
  }

  @override
  String toString() {
    return 'NavigationResult(destination: $destination, shouldReplace: $shouldReplace, clearHistory: $clearHistory, message: $message)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is NavigationResult &&
        other.destination == destination &&
        other.shouldReplace == shouldReplace &&
        other.clearHistory == clearHistory &&
        _mapEquals(other.queryParameters, queryParameters) &&
        other.message == message;
  }

  @override
  int get hashCode {
    return destination.hashCode ^
        shouldReplace.hashCode ^
        clearHistory.hashCode ^
        queryParameters.hashCode ^
        message.hashCode;
  }

  bool _mapEquals(Map<String, String> a, Map<String, String> b) {
    if (a.length != b.length) return false;
    for (final key in a.keys) {
      if (a[key] != b[key]) return false;
    }
    return true;
  }
}
