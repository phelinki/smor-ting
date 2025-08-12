# iOS Appium Setup and Execution Guide

## Prerequisites

- Xcode 16+ with Command Line Tools (detected)
- iOS Simulator runtime installed (e.g., iOS 18.5)
- Flutter SDK (detected)
- Node.js and Appium 2.x (detected)
- Python 3 with required packages installed (`requirements.txt`)

## One-time setup

1) Install Appium drivers

```bash
cd smor_ting_mobile/appium
npm run appium:install
```

2) Verify iOS simulators

```bash
xcrun simctl list devices available | grep iPhone
```

3) Build Flutter app for iOS simulator

```bash
cd ../
flutter clean && flutter pub get
flutter build ios --simulator --debug
```

This produces: `smor_ting_mobile/build/ios/iphonesimulator/Runner.app`.

If you prefer to use a different app path, export it:

```bash
export APP_PATH="/absolute/path/to/Runner.app"
```

## Running iOS tests

From `smor_ting_mobile/appium`:

```bash
./scripts/run_ios_tests.sh --device "iPhone 16 Pro Max" --ios-version 18.5
```

Options:

- `--suite <all|auth|registration|login>`: limit to a suite
- `--markers <marker>`: run pytest markers, e.g. `--markers "smoke"`
- `--parallel`: run tests in parallel

Reports will be generated under `smor_ting_mobile/appium/reports/`:

- Latest HTML: `reports/ios-report.html`
- JUnit XML: `reports/ios-junit.xml`
- Appium log: `reports/appium.log`
- Screenshots: `reports/screenshots/`

## Environment variables

- `APP_PATH`: absolute path to `.app` for simulator build
- `IOS_DEVICE_NAME`: simulator name (defaults to iPhone 13)
- `IOS_VERSION`: simulator iOS version (defaults to 16.4)
- `APPIUM_HOST`/`APPIUM_PORT`: Appium server host/port
- `ENVIRONMENT`: `local|ci|staging|production` (affects timeouts and options)

## Known build blockers (must be fixed in app code before tests can run)

While setting up iOS automation, the Flutter iOS simulator build failed with the following compile-time issues:

```text
Error (Xcode): lib/features/auth/presentation/providers/auth_provider.dart:146:8: Error: 'logout' is already declared in this scope.
```

```text
Error (Xcode): lib/features/navigation/presentation/enhanced_app_router.dart:264:11: Error: The type 'AuthState' is not exhaustively matched by the switch cases since it doesn't match 'Loading()'.
```

Please resolve these compile errors in the Flutter app. After a successful `flutter build ios --simulator`, re-run:

```bash
cd smor_ting_mobile/appium
./scripts/run_ios_tests.sh --device "iPhone 16 Pro Max" --ios-version 18.5
```

## Notes

- iOS simulator testing requires a simulator `.app` build (not an `.ipa`).
- You can pass a prebuilt app via `APP_PATH` to skip local build.
- Tests adhere to security, performance, and usability priorities [[memory:5639049]].


