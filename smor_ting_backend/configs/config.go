package configs

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	CORS     CORSConfig
	Logging  LoggingConfig
	Security SecurityConfig
	Momo     MomoConfig
	KYC      KYCConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
	// For in-memory database (testing/development)
	InMemory bool
	// MongoDB Atlas specific
	ConnectionString string
	AtlasCluster     bool
}

// AuthConfig holds authentication-related configuration
type AuthConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
	BCryptCost    int
	// New JWT configuration
	JWTAccessSecret  string
	JWTRefreshSecret string
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	EncryptionKey        string
	PaymentEncryptionKey string
	RateLimitRequests    int
	RateLimitWindow      time.Duration
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowHeaders     []string
	AllowMethods     []string
	AllowCredentials bool
	MaxAge           int
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level      string
	Format     string
	Output     string
	TimeFormat string
}

// MomoConfig holds MTN MoMo API configuration
type MomoConfig struct {
	BaseURL                     string
	TargetEnvironment           string // e.g., "sandbox" or "production"
	APIUser                     string // UUID
	APIKey                      string // Secret (treat as secret)
	SubscriptionKeyCollection   string // Ocp-Apim-Subscription-Key for Collection
	SubscriptionKeyDisbursement string // for Disbursement
	CallbackHost                string // public URL for callbacks/webhooks
}

// KYCConfig holds SmileID configuration
type KYCConfig struct {
	BaseURL     string
	PartnerID   string
	APIKey      string
	CallbackURL string
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			Driver:           getEnv("DB_DRIVER", "mongodb"),
			Host:             getEnv("DB_HOST", "localhost"),
			Port:             getEnv("DB_PORT", "27017"),
			Username:         getEnv("DB_USERNAME", ""),
			Password:         getEnv("DB_PASSWORD", ""),
			Database:         getEnv("DB_NAME", "smor_ting"),
			SSLMode:          getEnv("DB_SSL_MODE", "disable"),
			InMemory:         getBoolEnv("DB_IN_MEMORY", false),
			ConnectionString: getEnv("MONGODB_URI", ""),
			AtlasCluster:     getBoolEnv("MONGODB_ATLAS", false),
		},
		Auth: AuthConfig{
			JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			JWTExpiration:    getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
			BCryptCost:       getIntEnv("BCRYPT_COST", 12),
			JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "your-32-byte-access-secret-key-change-in-production"),
			JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-32-byte-refresh-secret-key-change-in-production"),
		},
		Security: SecurityConfig{
			EncryptionKey:        getEnv("ENCRYPTION_KEY", "your-32-byte-encryption-key-change-in-production"),
			PaymentEncryptionKey: getEnv("PAYMENT_ENCRYPTION_KEY", "your-32-byte-payment-encryption-key-change-in-production"),
			RateLimitRequests:    getIntEnv("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:      getDurationEnv("RATE_LIMIT_WINDOW", 1*time.Minute),
		},
		CORS: CORSConfig{
			AllowOrigins:     getStringSliceEnv("CORS_ALLOW_ORIGINS", getDefaultCORSOrigins()),
			AllowHeaders:     getStringSliceEnv("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "CF-Connecting-IP"}),
			AllowMethods:     getStringSliceEnv("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowCredentials: getBoolEnv("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getIntEnv("CORS_MAX_AGE", 86400),
		},
		Logging: LoggingConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			Output:     getEnv("LOG_OUTPUT", "stdout"),
			TimeFormat: getEnv("LOG_TIME_FORMAT", "2006-01-02T15:04:05Z07:00"),
		},
		Momo: MomoConfig{
			BaseURL:                     getEnv("MOMO_BASE_URL", ""),
			TargetEnvironment:           getEnv("MOMO_TARGET_ENV", "sandbox"),
			APIUser:                     getEnv("MOMO_API_USER", ""),
			APIKey:                      getEnv("MOMO_API_KEY", ""),
			SubscriptionKeyCollection:   getEnv("MOMO_SUB_KEY_COLLECTION", ""),
			SubscriptionKeyDisbursement: getEnv("MOMO_SUB_KEY_DISBURSEMENT", ""),
			CallbackHost:                getEnv("MOMO_CALLBACK_HOST", ""),
		},
		KYC: KYCConfig{
			BaseURL:     getEnv("SMILEID_BASE_URL", ""),
			PartnerID:   getEnv("SMILEID_PARTNER_ID", ""),
			APIKey:      getEnv("SMILEID_API_KEY", ""),
			CallbackURL: getEnv("SMILEID_CALLBACK_URL", ""),
		},
	}

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if c.Auth.JWTAccessSecret == "" {
		return fmt.Errorf("JWT access secret is required")
	}

	if c.Auth.JWTRefreshSecret == "" {
		return fmt.Errorf("JWT refresh secret is required")
	}

	if c.Security.EncryptionKey == "" {
		return fmt.Errorf("encryption key is required")
	}

	if c.Security.PaymentEncryptionKey == "" {
		return fmt.Errorf("payment encryption key is required")
	}

	// In production and staging, fail closed if any critical secret is missing, default, or not valid base64
	if c.IsProduction() || c.IsStaging() {
		// helper closure to check base64 length
		mustBeBase64 := func(name, value string) error {
			if value == "" {
				return fmt.Errorf("%s is required in production", name)
			}
			decoded, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				return fmt.Errorf("%s must be base64-encoded: %w", name, err)
			}
			if len(decoded) < 32 {
				return fmt.Errorf("%s must decode to at least 32 bytes", name)
			}
			return nil
		}

		// Must not be default placeholders
		if c.Auth.JWTAccessSecret == "your-32-byte-access-secret-key-change-in-production" {
			return fmt.Errorf("JWT_ACCESS_SECRET default value is not allowed in production")
		}
		if c.Auth.JWTRefreshSecret == "your-32-byte-refresh-secret-key-change-in-production" {
			return fmt.Errorf("JWT_REFRESH_SECRET default value is not allowed in production")
		}
		if c.Security.EncryptionKey == "your-32-byte-encryption-key-change-in-production" {
			return fmt.Errorf("ENCRYPTION_KEY default value is not allowed in production")
		}
		if c.Security.PaymentEncryptionKey == "your-32-byte-payment-encryption-key-change-in-production" {
			return fmt.Errorf("PAYMENT_ENCRYPTION_KEY default value is not allowed in production")
		}

		if err := mustBeBase64("JWT_ACCESS_SECRET", c.Auth.JWTAccessSecret); err != nil {
			return err
		}
		if err := mustBeBase64("JWT_REFRESH_SECRET", c.Auth.JWTRefreshSecret); err != nil {
			return err
		}
		if err := mustBeBase64("ENCRYPTION_KEY", c.Security.EncryptionKey); err != nil {
			return err
		}
		if err := mustBeBase64("PAYMENT_ENCRYPTION_KEY", c.Security.PaymentEncryptionKey); err != nil {
			return err
		}
	}

	// Check for default values and warn
	if c.Auth.JWTSecret == "your-secret-key-change-in-production" {
		fmt.Println("WARNING: Using default JWT secret. Change JWT_SECRET in production!")
	}

	if c.Auth.JWTAccessSecret == "your-32-byte-access-secret-key-change-in-production" {
		fmt.Println("WARNING: Using default JWT access secret. Change JWT_ACCESS_SECRET in production!")
	}

	if c.Auth.JWTRefreshSecret == "your-32-byte-refresh-secret-key-change-in-production" {
		fmt.Println("WARNING: Using default JWT refresh secret. Change JWT_REFRESH_SECRET in production!")
	}

	if c.Security.EncryptionKey == "your-32-byte-encryption-key-change-in-production" {
		fmt.Println("WARNING: Using default encryption key. Change ENCRYPTION_KEY in production!")
	}

	if c.Security.PaymentEncryptionKey == "your-32-byte-payment-encryption-key-change-in-production" {
		fmt.Println("WARNING: Using default payment encryption key. Change PAYMENT_ENCRYPTION_KEY in production!")
	}

	// MTN MoMo and KYC configuration checks
	if c.IsProduction() || c.IsStaging() {
		if c.Momo.BaseURL == "" || c.Momo.APIUser == "" || c.Momo.APIKey == "" {
			return fmt.Errorf("MoMo configuration is required in production")
		}
		if c.Momo.SubscriptionKeyCollection == "" && c.Momo.SubscriptionKeyDisbursement == "" {
			return fmt.Errorf("at least one MoMo subscription key is required in production")
		}
		if c.KYC.BaseURL == "" || c.KYC.PartnerID == "" || c.KYC.APIKey == "" {
			return fmt.Errorf("SmileID KYC configuration is required in production")
		}
	}

	return nil
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return os.Getenv("ENV") == "development" || os.Getenv("ENV") == ""
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return os.Getenv("ENV") == "production"
}

// IsStaging returns true if the application is running in staging mode
func (c *Config) IsStaging() bool {
	return os.Getenv("ENV") == "staging"
}

// Helper functions for environment variable parsing

// getDefaultCORSOrigins returns appropriate CORS origins based on environment
func getDefaultCORSOrigins() []string {
	env := os.Getenv("ENV")
	switch env {
	case "production":
		return []string{
			"https://smor-ting.com",
			"https://www.smor-ting.com",
			"https://api.smor-ting.com",
		}
	case "staging":
		return []string{
			"https://staging.smor-ting.com",
			"https://api-staging.smor-ting.com",
			"http://localhost:3000",
		}
	default: // development
		return []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:8080",
			"http://127.0.0.1:8080",
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated values
		return []string{value} // For now, just return as single value
	}
	return defaultValue
}
