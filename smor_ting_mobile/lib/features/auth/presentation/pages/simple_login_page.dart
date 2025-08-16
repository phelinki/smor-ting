import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/auth_provider.dart';

class SimpleLoginPage extends ConsumerStatefulWidget {
  const SimpleLoginPage({super.key});

  @override
  ConsumerState<SimpleLoginPage> createState() => _SimpleLoginPageState();
}

class _SimpleLoginPageState extends ConsumerState<SimpleLoginPage> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();



  Future<void> _handleLogin() async {
    try {
      final authNotifier = ref.read(authNotifierProvider.notifier);
      await authNotifier.login(
        _emailController.text.trim(),
        _passwordController.text,
      );
      
      // User is authenticated, GoRouter will handle navigation
      // No manual navigation needed
      
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Login failed: ${e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    ref.listen(authNotifierProvider, (previous, next) {
      next.when(
        initial: () {},
        loading: () {},
        authenticated: (user, token) {
          // User is authenticated, GoRouter will handle navigation
          // No manual navigation needed
        },
        error: (message) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(message)),
          );
        },
        emailAlreadyExists: (email) {},
        passwordResetEmailSent: (email) {},
        passwordResetSuccess: () {},
      );
    });

    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            TextField(
              controller: _emailController,
              decoration: InputDecoration(labelText: 'Email'),
            ),
            SizedBox(height: 16),
            TextField(
              controller: _passwordController,
              decoration: InputDecoration(labelText: 'Password'),
              obscureText: true,
            ),
            SizedBox(height: 24),
            ElevatedButton(
              onPressed: () {
                ref.read(authNotifierProvider.notifier).login(
                  _emailController.text,
                  _passwordController.text,
                );
              },
              child: Text('Login'),
            ),
            SizedBox(height: 16),
            TextButton(
              onPressed: () {
                // Use GoRouter for navigation
                context.go('/forgot-password');
              },
              child: Text('Forgot Password?'),
            ),
            TextButton(
              onPressed: () {
                // Use GoRouter for navigation
                context.go('/register');
              },
              child: Text('Register Here'),
            ),
          ],
        ),
      ),
    );
  }
}
