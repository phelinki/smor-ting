import 'package:flutter/material.dart';

class AppTheme {
  // Apple HIG Colors
  static const Color primaryRed = Color(0xFFFF3B30);  // Apple's red for primary actions
  static const Color secondaryBlue = Color(0xFF007AFF);  // Apple's blue for secondary elements
  static const Color successGreen = Color(0xFF34C759);  // Apple's green for success states
  static const Color white = Color(0xFFFFFFFF);
  static const Color lightGray = Color(0xFFF2F2F7);  // Apple's light gray for backgrounds
  
  // Text colors
  static const Color textPrimary = Color(0xFF000000);  // Black for primary text
  static const Color textSecondary = Color(0xFF8E8E93);  // Gray for secondary text
  
  // Additional colors
  static const Color darkGray = Color(0xFF424242);
  static const Color warning = Color(0xFFFF9500);  // Apple's orange
  static const Color error = Color(0xFFFF3B30);  // Using red for errors

  // Light Theme
  static ThemeData lightTheme = ThemeData(
    useMaterial3: true,
    fontFamily: 'SF Pro Display',  // Apple's font family
    colorScheme: ColorScheme.fromSeed(
      seedColor: primaryRed,
      brightness: Brightness.light,
      primary: primaryRed,
      secondary: secondaryBlue,
      surface: white,
      error: error,
      onPrimary: white,
      onSecondary: white,
      onSurface: textPrimary,
    ),
    scaffoldBackgroundColor: white,
    appBarTheme: const AppBarTheme(
      backgroundColor: white,
      foregroundColor: textPrimary,
      elevation: 0,
      centerTitle: true,
      titleTextStyle: TextStyle(
        color: textPrimary,
        fontSize: 18,
        fontWeight: FontWeight.w600,
        fontFamily: 'SF Pro Display',
      ),
      iconTheme: IconThemeData(
        color: secondaryBlue,
        size: 24,
      ),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: primaryRed,
        foregroundColor: white,
        elevation: 0,
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          fontFamily: 'SF Pro Display',
        ),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: secondaryBlue,
        backgroundColor: white,
        side: const BorderSide(color: secondaryBlue, width: 1.5),
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          fontFamily: 'SF Pro Display',
        ),
      ),
    ),
    textButtonTheme: TextButtonThemeData(
      style: TextButton.styleFrom(
        foregroundColor: secondaryBlue,
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          fontFamily: 'SF Pro Display',
        ),
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: lightGray,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: BorderSide.none,
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: BorderSide.none,
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: secondaryBlue, width: 2),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: error, width: 2),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
      hintStyle: const TextStyle(
        color: textSecondary,
        fontSize: 16,
        fontFamily: 'SF Pro Display',
      ),
      labelStyle: const TextStyle(
        color: textSecondary,
        fontSize: 16,
        fontFamily: 'SF Pro Display',
      ),
    ),
    cardTheme: const CardTheme(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.all(Radius.circular(16)),
      ),
      color: white,
      shadowColor: Color(0x1F000000),
    ),
    bottomNavigationBarTheme: const BottomNavigationBarThemeData(
      backgroundColor: white,
      selectedItemColor: primaryRed,
      unselectedItemColor: textSecondary,
      type: BottomNavigationBarType.fixed,
      elevation: 8,
      selectedLabelStyle: TextStyle(
        fontSize: 12,
        fontWeight: FontWeight.w600,
        fontFamily: 'SF Pro Display',
      ),
      unselectedLabelStyle: TextStyle(
        fontSize: 12,
        fontWeight: FontWeight.w400,
        fontFamily: 'SF Pro Display',
      ),
    ),
    textTheme: const TextTheme(
      headlineLarge: TextStyle(
        fontSize: 32,
        fontWeight: FontWeight.bold,
        color: textPrimary,
        fontFamily: 'SF Pro Display',
      ),
      headlineMedium: TextStyle(
        fontSize: 28,
        fontWeight: FontWeight.bold,
        color: textPrimary,
        fontFamily: 'SF Pro Display',
      ),
      headlineSmall: TextStyle(
        fontSize: 24,
        fontWeight: FontWeight.w600,
        color: textPrimary,
        fontFamily: 'SF Pro Display',
      ),
      titleLarge: TextStyle(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: textPrimary,
        fontFamily: 'SF Pro Display',
      ),
      titleMedium: TextStyle(
        fontSize: 18,
        fontWeight: FontWeight.w600,
        color: textPrimary,
        fontFamily: 'SF Pro Display',
      ),
      titleSmall: TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        color: textPrimary,
        fontFamily: 'SF Pro Display',
      ),
      bodyLarge: TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w400,
        color: textPrimary,
        fontFamily: 'SF Pro Display',
      ),
      bodyMedium: TextStyle(
        fontSize: 14,
        fontWeight: FontWeight.w400,
        color: textPrimary,
        fontFamily: 'SF Pro Display',
      ),
      bodySmall: TextStyle(
        fontSize: 12,
        fontWeight: FontWeight.w400,
        color: textSecondary,
        fontFamily: 'SF Pro Display',
      ),
    ),
  );

  // Dark Theme
  static ThemeData darkTheme = ThemeData(
    useMaterial3: true,
    fontFamily: 'SF Pro Display',
    colorScheme: ColorScheme.fromSeed(
      seedColor: primaryRed,
      brightness: Brightness.dark,
      primary: primaryRed,
      secondary: secondaryBlue,
      surface: const Color(0xFF1C1C1E),
      error: error,
      onPrimary: white,
      onSecondary: white,
      onSurface: white,
    ),
    scaffoldBackgroundColor: const Color(0xFF000000),
    appBarTheme: const AppBarTheme(
      backgroundColor: Color(0xFF1C1C1E),
      foregroundColor: white,
      elevation: 0,
      centerTitle: true,
      titleTextStyle: TextStyle(
        color: white,
        fontSize: 18,
        fontWeight: FontWeight.w600,
        fontFamily: 'SF Pro Display',
      ),
      iconTheme: IconThemeData(
        color: secondaryBlue,
        size: 24,
      ),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: primaryRed,
        foregroundColor: white,
        elevation: 0,
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          fontFamily: 'SF Pro Display',
        ),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: secondaryBlue,
        backgroundColor: const Color(0xFF1C1C1E),
        side: const BorderSide(color: secondaryBlue, width: 1.5),
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          fontFamily: 'SF Pro Display',
        ),
      ),
    ),
    textButtonTheme: TextButtonThemeData(
      style: TextButton.styleFrom(
        foregroundColor: secondaryBlue,
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          fontFamily: 'SF Pro Display',
        ),
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: const Color(0xFF2C2C2E),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: BorderSide.none,
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: BorderSide.none,
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: secondaryBlue, width: 2),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: error, width: 2),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
      hintStyle: const TextStyle(
        color: textSecondary,
        fontSize: 16,
        fontFamily: 'SF Pro Display',
      ),
      labelStyle: const TextStyle(
        color: textSecondary,
        fontSize: 16,
        fontFamily: 'SF Pro Display',
      ),
    ),
    cardTheme: const CardTheme(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.all(Radius.circular(16)),
      ),
      color: Color(0xFF1C1C1E),
      shadowColor: Color(0x3F000000),
    ),
    bottomNavigationBarTheme: const BottomNavigationBarThemeData(
      backgroundColor: Color(0xFF1C1C1E),
      selectedItemColor: primaryRed,
      unselectedItemColor: textSecondary,
      type: BottomNavigationBarType.fixed,
      elevation: 8,
      selectedLabelStyle: TextStyle(
        fontSize: 12,
        fontWeight: FontWeight.w600,
        fontFamily: 'SF Pro Display',
      ),
      unselectedLabelStyle: TextStyle(
        fontSize: 12,
        fontWeight: FontWeight.w400,
        fontFamily: 'SF Pro Display',
      ),
    ),
    textTheme: const TextTheme(
      headlineLarge: TextStyle(
        fontSize: 32,
        fontWeight: FontWeight.bold,
        color: white,
        fontFamily: 'SF Pro Display',
      ),
      headlineMedium: TextStyle(
        fontSize: 28,
        fontWeight: FontWeight.bold,
        color: white,
        fontFamily: 'SF Pro Display',
      ),
      headlineSmall: TextStyle(
        fontSize: 24,
        fontWeight: FontWeight.w600,
        color: white,
        fontFamily: 'SF Pro Display',
      ),
      titleLarge: TextStyle(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: white,
        fontFamily: 'SF Pro Display',
      ),
      titleMedium: TextStyle(
        fontSize: 18,
        fontWeight: FontWeight.w600,
        color: white,
        fontFamily: 'SF Pro Display',
      ),
      titleSmall: TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        color: white,
        fontFamily: 'SF Pro Display',
      ),
      bodyLarge: TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w400,
        color: white,
        fontFamily: 'SF Pro Display',
      ),
      bodyMedium: TextStyle(
        fontSize: 14,
        fontWeight: FontWeight.w400,
        color: white,
        fontFamily: 'SF Pro Display',
      ),
      bodySmall: TextStyle(
        fontSize: 12,
        fontWeight: FontWeight.w400,
        color: textSecondary,
        fontFamily: 'SF Pro Display',
      ),
    ),
  );
} 