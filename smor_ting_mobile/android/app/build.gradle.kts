import java.io.FileInputStream
import java.util.Properties

plugins {
    id("com.android.application")
    id("kotlin-android")
    // The Flutter Gradle Plugin must be applied after the Android and Kotlin Gradle plugins.
    id("dev.flutter.flutter-gradle-plugin")
}

// Load keystore properties
val keystorePropertiesFile = rootProject.file("key.properties")
val keystoreProperties = Properties()
if (keystorePropertiesFile.exists()) {
    keystoreProperties.load(FileInputStream(keystorePropertiesFile))
}

android {
    namespace = "com.smorting.app.smor_ting_mobile"
    compileSdk = flutter.compileSdkVersion
    // Pin to NDK version required by plugins (per Flutter build output recommendation)
    ndkVersion = "27.0.12077973"

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }

    kotlinOptions {
        jvmTarget = JavaVersion.VERSION_11.toString()
    }

    defaultConfig {
        // TODO: Specify your own unique Application ID (https://developer.android.com/studio/build/application-id.html).
        applicationId = "com.smorting.app.smor_ting_mobile"
        // You can update the following values to match your application needs.
        // For more information, see: https://flutter.dev/to/review-gradle-config.
        minSdk = flutter.minSdkVersion
        targetSdk = flutter.targetSdkVersion
        versionCode = flutter.versionCode
        versionName = flutter.versionName
    }

    signingConfigs {
        create("release") {
            if (keystorePropertiesFile.exists() && keystoreProperties["storeFile"] != null) {
                val keystoreFile = file(keystoreProperties["storeFile"] as String)
                if (keystoreFile.exists()) {
                    keyAlias = keystoreProperties["keyAlias"] as String?
                    keyPassword = keystoreProperties["keyPassword"] as String?
                    storeFile = keystoreFile
                    storePassword = keystoreProperties["storePassword"] as String?
                }
            }
        }
    }

    buildTypes {
        release {
            // Use release signing config if available, otherwise use debug
            signingConfig = if (keystorePropertiesFile.exists() && 
                               keystoreProperties["storeFile"] != null && 
                               file(keystoreProperties["storeFile"] as String).exists()) {
                signingConfigs.getByName("release")
            } else {
                signingConfigs.getByName("debug")
            }
            // Explicitly disable minification and resource shrinking for testing
            isMinifyEnabled = false
            isShrinkResources = false
            // proguardFiles(getDefaultProguardFile("proguard-android-optimize.txt"), "proguard-rules.pro")
        }
        debug {
            isMinifyEnabled = false
            isShrinkResources = false
        }
    }
    
    // Suppress known lint issue as per Flutter fix guidance (Groovy lintOptions -> Kotlin DSL)
    lint {
        checkReleaseBuilds = false
    }
}

dependencies {
    // Updated Play Core library for Android 14 compatibility
    implementation("com.google.android.play:core-ktx:1.8.1")
}

flutter {
    source = "../.."
}
