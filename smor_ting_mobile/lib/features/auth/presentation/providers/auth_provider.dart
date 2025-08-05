import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'auth_provider.g.dart';

@riverpod
class AuthNotifier extends _$AuthNotifier {
  @override
  AuthState build() {
    return const AuthState.initial();
  }

  Future<void> login(String email, String password) async {
    state = const AuthState.loading();
    
    try {
      // TODO: Implement actual login logic
      await Future.delayed(const Duration(seconds: 2)); // Simulate API call
      
      // For now, just simulate successful login
      state = AuthState.authenticated(
        User(
          id: '1',
          email: email,
          name: 'Test User',
          phone: '+231123456789',
          userType: 'customer',
        ),
      );
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  Future<void> register(String name, String email, String phone, String password) async {
    state = const AuthState.loading();
    
    try {
      // TODO: Implement actual registration logic
      await Future.delayed(const Duration(seconds: 2)); // Simulate API call
      
      // For now, just simulate successful registration
      state = AuthState.authenticated(
        User(
          id: '1',
          email: email,
          name: name,
          phone: phone,
          userType: 'customer',
        ),
      );
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  void logout() {
    state = const AuthState.initial();
  }
}

// Auth State
sealed class AuthState {
  const AuthState();

  const factory AuthState.initial() = Initial;
  const factory AuthState.loading() = Loading;
  factory AuthState.authenticated(User user) = Authenticated;
  const factory AuthState.error(String message) = Error;
}

class Initial extends AuthState {
  const Initial();
}

class Loading extends AuthState {
  const Loading();
}

class Authenticated extends AuthState {
  final User user;
  Authenticated(this.user);
}

class Error extends AuthState {
  final String message;
  const Error(this.message);
}

// User Model
class User {
  final String id;
  final String email;
  final String name;
  final String phone;
  final String userType;
  final String? profileImage;
  final bool isVerified;
  final DateTime createdAt;

  User({
    required this.id,
    required this.email,
    required this.name,
    required this.phone,
    required this.userType,
    this.profileImage,
    this.isVerified = false,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime(2024, 1, 1);

  User copyWith({
    String? id,
    String? email,
    String? name,
    String? phone,
    String? userType,
    String? profileImage,
    bool? isVerified,
    DateTime? createdAt,
  }) {
    return User(
      id: id ?? this.id,
      email: email ?? this.email,
      name: name ?? this.name,
      phone: phone ?? this.phone,
      userType: userType ?? this.userType,
      profileImage: profileImage ?? this.profileImage,
      isVerified: isVerified ?? this.isVerified,
      createdAt: createdAt ?? this.createdAt,
    );
  }
} 