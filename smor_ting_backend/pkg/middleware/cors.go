package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/pkg/logger"
	"go.uber.org/zap"
)

// CORSMiddleware represents CORS middleware with configuration
type CORSMiddleware struct {
	config *configs.CORSConfig
	logger *logger.Logger
}

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware(config *configs.CORSConfig, logger *logger.Logger) (*CORSMiddleware, error) {
	if config == nil {
		return nil, fmt.Errorf("CORS configuration is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &CORSMiddleware{
		config: config,
		logger: logger,
	}, nil
}

// Configure returns a configured CORS middleware
func (cm *CORSMiddleware) Configure() fiber.Handler {
	// Convert string slices to individual strings for Fiber CORS
	allowOrigins := cm.config.AllowOrigins
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	allowHeaders := cm.config.AllowHeaders
	if len(allowHeaders) == 0 {
		allowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}

	allowMethods := cm.config.AllowMethods
	if len(allowMethods) == 0 {
		allowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}

	cm.logger.Info("Configuring CORS middleware",
		zap.Strings("allow_origins", allowOrigins),
		zap.Strings("allow_headers", allowHeaders),
		zap.Strings("allow_methods", allowMethods),
		zap.Bool("allow_credentials", cm.config.AllowCredentials),
		zap.Int("max_age", cm.config.MaxAge),
	)

	// Convert slices to comma-separated strings for Fiber CORS
	originsStr := ""
	if len(allowOrigins) > 0 {
		originsStr = allowOrigins[0]
		for i := 1; i < len(allowOrigins); i++ {
			originsStr += "," + allowOrigins[i]
		}
	}

	headersStr := ""
	if len(allowHeaders) > 0 {
		headersStr = allowHeaders[0]
		for i := 1; i < len(allowHeaders); i++ {
			headersStr += "," + allowHeaders[i]
		}
	}

	methodsStr := ""
	if len(allowMethods) > 0 {
		methodsStr = allowMethods[0]
		for i := 1; i < len(allowMethods); i++ {
			methodsStr += "," + allowMethods[i]
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     originsStr,
		AllowHeaders:     headersStr,
		AllowMethods:     methodsStr,
		AllowCredentials: cm.config.AllowCredentials,
		MaxAge:           cm.config.MaxAge,
		ExposeHeaders:    "Content-Length,Content-Type",
	})
}

// DevelopmentCORS returns CORS configuration suitable for development
func DevelopmentCORS(logger *logger.Logger) fiber.Handler {
	config := &configs.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	middleware, err := NewCORSMiddleware(config, logger)
	if err != nil {
		logger.Error("Failed to create development CORS middleware", err)
		// Fallback to basic CORS
		return cors.New()
	}

	return middleware.Configure()
}

// ProductionCORS returns CORS configuration suitable for production
func ProductionCORS(allowedOrigins []string, logger *logger.Logger) fiber.Handler {
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"https://smor-ting.com", "https://www.smor-ting.com"}
	}

	config := &configs.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	middleware, err := NewCORSMiddleware(config, logger)
	if err != nil {
		logger.Error("Failed to create production CORS middleware", err)
		// Fallback to basic CORS
		return cors.New()
	}

	return middleware.Configure()
}
