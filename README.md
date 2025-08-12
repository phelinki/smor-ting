# Smor-Ting - Handyman and Service Marketplace for Liberia

A comprehensive offline-first marketplace platform connecting customers with verified service providers in Liberia.

## 🎯 Project Overview

Smor-Ting is a mobile-first application built with Flutter (mobile) and Go (backend) that enables customers to discover, book, and pay for various services while providing service providers with tools to manage their business, track earnings, and grow their clientele.

## 🏗️ Architecture

### Mobile App (Flutter)
- **State Management**: Riverpod
- **Navigation**: GoRouter
- **Local Storage**: Hive + SQLite
- **Networking**: Dio + Retrofit
- **UI**: Material Design 3 with custom Liberia theme

### Backend (Go)
- **Framework**: Fiber
- **Database**: MongoDB (with offline-first capabilities)
- **Authentication**: JWT + Biometric
- **Payment**: Flutterwave, Orange Money, MTN Mobile Money

## 🎨 Design System

### Color Palette (Liberia Flag Colors)
- **Primary Red**: `#D21034`
- **Primary Blue**: `#002868`
- **White**: `#FFFFFF`
- **Additional Colors**: Light/Dark variants and semantic colors

### Typography
- **Font Family**: Poppins
- **Weights**: Regular, Medium, SemiBold, Bold

## 🚀 Getting Started

### Prerequisites
- Flutter SDK (3.2.3 or higher)
- Go 1.21 or higher
- Android Studio / Xcode for mobile development
- MongoDB (for backend)

### Mobile App Setup

1. **Navigate to mobile directory**:
   ```bash
   cd smor_ting_mobile
   ```

2. **Install dependencies**:
   ```bash
   flutter pub get
   ```

3. **Run code generation**:
   ```bash
   flutter packages pub run build_runner build
   ```

4. **Run the app**:
   ```bash
   flutter run
   ```

### Backend Setup

1. **Navigate to backend directory**:
   ```bash
   cd smor_ting_backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the server**:
   ```bash
   go run cmd/main.go
   ```

The backend will start on `http://localhost:8080`

## 📱 Features

### Customer Features
- ✅ Service discovery and booking
- ✅ Provider selection and reviews
- ✅ Multiple payment methods
- ✅ Real-time service tracking
- ✅ Offline-first functionality

### Service Provider Features
- ✅ Professional dashboard
- ✅ Service management
- ✅ Booking management
- ✅ Earnings tracking
- ✅ Verification system

### Admin Features
- ✅ User management
- ✅ Service oversight
- ✅ Financial management
- ✅ Analytics dashboard

## 🔒 Security Features

- Multi-factor authentication (2FA)
- Biometric authentication
- End-to-end encryption
- PCI DSS compliant payments
- Real-time fraud detection

## 💳 Payment Integration

- **International**: Flutterwave
- **Mobile Money**: Orange Money, MTN Mobile Money, Lonestar Cell MTN
- **Bank Transfers**: Local bank integration
- **Digital Wallet**: In-app wallet with LRD/USD support

## 🌍 Liberia-Specific Features

- Liberia phone number validation
- Local currency (LRD) support
- Liberia timezone (Africa/Monrovia)
- Local service categories
- Mobile money optimization

## 📊 Offline-First Architecture

- Local SQLite database for offline data
- Automatic sync when online
- Conflict resolution for offline edits
- Cached frequently accessed data
- Optimized for slow network conditions

## 🛠️ Development

### Project Structure
```
smor-ting/
├── smor_ting_mobile/          # Flutter mobile app
│   ├── lib/
│   │   ├── core/             # Core utilities, theme, constants
│   │   ├── features/         # Feature modules
│   │   └── main.dart         # App entry point
│   └── assets/               # Images, fonts, animations
├── smor_ting_backend/         # Go backend
│   ├── cmd/                  # Application entry points
│   ├── internal/             # Private application code
│   ├── pkg/                  # Public libraries
│   └── api/                  # API definitions
└── shared/                   # Shared assets and documentation
```

### Code Generation
The project uses several code generators:
- **Riverpod**: State management code generation
- **Retrofit**: API client generation
- **JSON Serializable**: JSON serialization
- **Hive**: Local storage code generation

Run code generation after dependency changes:
```bash
flutter packages pub run build_runner build
```

## 🧪 Testing

### Mobile App Testing
```bash
# Unit tests
flutter test

# Integration tests
flutter test integration_test/
```

### Backend Testing
```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 📦 Deployment

### Mobile App
- **Android**: Google Play Store
- **iOS**: App Store
- **Build Commands**:
  ```bash
  # Android
  flutter build apk --release
  
  # iOS
  flutter build ios --release
  ```

### Backend
- **Docker**: Containerized deployment
- **Kubernetes**: Orchestration
- **Environment**: Production-ready with health checks

## 🔧 Configuration

### Environment Variables
```bash
# Backend
PORT=8080
MONGODB_URI=mongodb://localhost:27017/smorting
JWT_SECRET=YOUR_JWT_SECRET_MIN_32_CHARS

# Mobile App
API_BASE_URL=http://localhost:8080/api/v1
```

## 📈 Monitoring & Analytics

- Application Performance Monitoring (APM)
- Error tracking and alerting
- User behavior analytics
- Business metrics dashboard
- Security incident detection

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is proprietary software. All rights reserved.

## 📞 Support

For support and questions:
- Email: support@smorting.com
- Documentation: [Coming Soon]
- Issues: [GitHub Issues]

## 🗺️ Roadmap

### Phase 1 (MVP) - Q1 2024
- ✅ Basic authentication
- ✅ Service discovery
- ✅ Booking system
- ✅ Payment integration
- ✅ Offline functionality

### Phase 2 - Q2 2024
- AI-powered recommendations
- Video consultations
- Advanced analytics
- White-label platform

### Phase 3 - Q3 2024
- IoT device integration
- AR service visualization
- Blockchain reputation system
- International expansion

---

**Built with ❤️ for Liberia** # smor-ting
