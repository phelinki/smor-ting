package configs

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
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
			JWTSecret:        getEnv("JWT_SECRET", ""),
			JWTExpiration:    getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
			BCryptCost:       getIntEnv("BCRYPT_COST", 12),
			JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "abcdefghijklmnopqrstuvwxyz123456"),  // Exactly 32 bytes for development
			JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "zyxwvutsrqponmlkjihgfedcba654321"), // Exactly 32 bytes for development
		},
		Security: SecurityConfig{
			EncryptionKey:        getEnv("ENCRYPTION_KEY", "12345678901234567890123456789012"),         // Exactly 32 bytes for development
			PaymentEncryptionKey: getEnv("PAYMENT_ENCRYPTION_KEY", "12345678901234567890123456789012"), // Exactly 32 bytes for development
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

	// Security validation for JWT secrets
	// Allow legacy JWT_SECRET to be optional when access/refresh secrets are provided
	if c.Auth.JWTAccessSecret == "" || c.Auth.JWTRefreshSecret == "" {
		if c.Auth.JWTSecret == "" {
			return fmt.Errorf("JWT secret is required")
		}
		if err := c.validateSecretSecurity("JWT_SECRET", c.Auth.JWTSecret); err != nil {
			return err
		}
	}

	if c.Auth.JWTAccessSecret == "" {
		return fmt.Errorf("JWT access secret is required")
	}
	if err := c.validateSecretSecurity("JWT_ACCESS_SECRET", c.Auth.JWTAccessSecret); err != nil {
		return err
	}
	if c.IsProduction() || c.IsStaging() {
		if err := c.validateBase64Key("JWT_ACCESS_SECRET", c.Auth.JWTAccessSecret); err != nil {
			return err
		}
	}

	if c.Auth.JWTRefreshSecret == "" {
		return fmt.Errorf("JWT refresh secret is required")
	}
	if err := c.validateSecretSecurity("JWT_REFRESH_SECRET", c.Auth.JWTRefreshSecret); err != nil {
		return err
	}
	if c.IsProduction() || c.IsStaging() {
		if err := c.validateBase64Key("JWT_REFRESH_SECRET", c.Auth.JWTRefreshSecret); err != nil {
			return err
		}
	}

	if c.Security.EncryptionKey == "" {
		return fmt.Errorf("encryption key is required")
	}
	if err := c.validateSecretSecurity("ENCRYPTION_KEY", c.Security.EncryptionKey); err != nil {
		return err
	}
	if c.IsProduction() || c.IsStaging() {
		if err := c.validateBase64Key("ENCRYPTION_KEY", c.Security.EncryptionKey); err != nil {
			return err
		}
	}

	if c.Security.PaymentEncryptionKey == "" {
		return fmt.Errorf("payment encryption key is required")
	}
	if err := c.validateSecretSecurity("PAYMENT_ENCRYPTION_KEY", c.Security.PaymentEncryptionKey); err != nil {
		return err
	}
	if c.IsProduction() || c.IsStaging() {
		if err := c.validateBase64Key("PAYMENT_ENCRYPTION_KEY", c.Security.PaymentEncryptionKey); err != nil {
			return err
		}
	}

	// Database security validation
	if c.Database.ConnectionString != "" {
		if err := c.validateDatabaseConnection(); err != nil {
			return err
		}
	}

	// Production-specific validations
	if c.IsProduction() || c.IsStaging() {
		// Database configuration is required in production
		if c.Database.ConnectionString == "" && !c.Database.InMemory {
			return fmt.Errorf("MONGODB_URI is required in production")
		}

		// MoMo API validation in production
		if c.Momo.BaseURL == "" || c.Momo.APIUser == "" || c.Momo.APIKey == "" {
			return fmt.Errorf("MTN MoMo API configuration is required in production")
		}

		// KYC validation in production
		if c.KYC.BaseURL == "" || c.KYC.PartnerID == "" || c.KYC.APIKey == "" {
			return fmt.Errorf("SmileID KYC configuration is required in production")
		}
	}

	return nil
}

// validateSecretSecurity checks if a secret meets security requirements
func (c *Config) validateSecretSecurity(name, value string) error {
	// List of insecure default values that should never be used in production
	insecureDefaults := []string{
		"your-secret-key",
		"your-secret-key-change-in-production",
		"your-32-byte-access-secret-key-change-in-production",
		"your-32-byte-refresh-secret-key-change-in-production",
		"your-32-byte-encryption-key-change-in-production",
		"your-32-byte-payment-encryption-key-change-in-production",
		"YOUR_JWT_SECRET_MIN_32_CHARS",
		"YOUR_JWT_SECRET_MIN_32_CHARS_GENERATE_WITH_OPENSSL",
		"YOUR_ACCESS_SECRET_BASE64_ENCODED",
		"YOUR_REFRESH_SECRET_BASE64_ENCODED",
		"YOUR_ENCRYPTION_KEY_32_BYTES",
		"YOUR_PAYMENT_ENCRYPTION_KEY_32_BYTES",
		"YOUR_PRODUCTION_JWT_SECRET_MIN_32_CHARS",
		"change-this-in-production",
		"default",
		"secret",
		"password",
		"123456",
	}

	// Allow development-only defaults (but only in development mode)
	developmentDefaults := []string{
		"abcdefghijklmnopqrstuvwxyz123456", // JWT access secret
		"zyxwvutsrqponmlkjihgfedcba654321", // JWT refresh secret
		"12345678901234567890123456789012", // Encryption keys
	}

	// Check if the value is one of the insecure defaults
	for _, insecure := range insecureDefaults {
		if value == insecure {
			return fmt.Errorf("%s contains an insecure default value. Please set a secure value", name)
		}
	}

	// In production/staging, also check development defaults are not used
	if c.IsProduction() || c.IsStaging() {
		for _, devDefault := range developmentDefaults {
			if value == devDefault {
				return fmt.Errorf("%s is using a development default value. Please set a secure production value", name)
			}
		}
	}

	// Minimum length requirement
	if len(value) < 32 {
		return fmt.Errorf("%s must be at least 32 characters long for security", name)
	}

	// In production/staging, enforce stronger requirements
	if c.IsProduction() || c.IsStaging() {
		// Check for common weak patterns
		if isWeakSecret(value) {
			return fmt.Errorf("%s appears to be weak. Use a cryptographically secure random value", name)
		}
	}

	return nil
}

// validateBase64Key ensures the provided value is base64-encoded 32 bytes
func (c *Config) validateBase64Key(name, value string) error {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return fmt.Errorf("%s must be base64-encoded: %w", name, err)
	}
	if len(decoded) != 32 {
		return fmt.Errorf("%s must decode to 32 bytes", name)
	}
	return nil
}

// validateDatabaseConnection checks database connection string security
func (c *Config) validateDatabaseConnection() error {
	uri := c.Database.ConnectionString

	// Check for insecure placeholder values
	insecurePlaceholders := []string{
		"your_password",
		"your_username",
		"YOUR_PASSWORD",
		"YOUR_USERNAME",
		"password",
		"username",
		"smorting_user",
		"cluster0.xxxxx.mongodb.net",
		"YOUR_CLUSTER.mongodb.net",
	}

	for _, placeholder := range insecurePlaceholders {
		if strings.Contains(uri, placeholder) {
			return fmt.Errorf("MONGODB_URI contains placeholder values. Please use actual credentials")
		}
	}

	// In production, ensure we're using mongodb+srv (Atlas) or secure connection
	if c.IsProduction() {
		if !strings.HasPrefix(uri, "mongodb+srv://") && !strings.Contains(uri, "ssl=true") {
			return fmt.Errorf("production database connection must use SSL/TLS")
		}
	}

	return nil
}

// isWeakSecret checks if a secret follows weak patterns
func isWeakSecret(secret string) bool {
	// Convert to lowercase for checking
	lower := strings.ToLower(secret)

	// Check for dictionary words or common patterns
	weakPatterns := []string{
		"password", "secret", "admin", "user", "test", "demo",
		"123456", "qwerty", "abc", "letmein", "welcome",
		"smor", "ting", "smorting", "api", "key",
	}

	for _, pattern := range weakPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	// Check if it's all the same character repeated
	if len(secret) > 0 {
		firstChar := secret[0]
		allSame := true
		for _, char := range secret {
			if byte(char) != firstChar {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}

	// Check for simple patterns like "abcd", "1234", etc.
	if isSequentialPattern(secret) {
		return true
	}

	return false
}

// isSequentialPattern checks for sequential character patterns
func isSequentialPattern(s string) bool {
	if len(s) < 4 {
		return false
	}

	// Check for increasing sequences
	increasing := true
	decreasing := true

	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1]+1 {
			increasing = false
		}
		if s[i] != s[i-1]-1 {
			decreasing = false
		}
	}

	return increasing || decreasing
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
