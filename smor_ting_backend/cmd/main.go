package main

// Force Railway deployment - removed railway.toml configuration
import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/migrations"
	pkgDatabase "github.com/smorting/backend/pkg/database"
	"github.com/smorting/backend/pkg/logger"
	"github.com/smorting/backend/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// App represents the application
type App struct {
	config          *configs.Config
	logger          *logger.Logger
	mongoDB         *pkgDatabase.MongoDB
	repository      database.Repository
	migrator        *migrations.Migrator
	changeStreamSvc *services.ChangeStreamService
	authSvc         *auth.MongoDBService
	authHdl         *auth.MongoDBHandler
	// New security services
	jwtService          *services.JWTRefreshService
	encryptionService   *services.EncryptionService
	pciService          *services.PCIDSSService
	auditService        *services.AuditService
	authHandler         *handlers.AuthHandler
	enhancedAuthHandler *handlers.EnhancedAuthHandler
	server              *fiber.App
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	loggerInstance, err := logger.New(config.Logging.Level, config.Logging.Format, config.Logging.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize default logger
	if err := logger.InitDefault(config.Logging.Level, config.Logging.Format, config.Logging.Output); err != nil {
		return nil, fmt.Errorf("failed to initialize default logger: %w", err)
	}

	app := &App{
		config: config,
		logger: loggerInstance,
	}

	// Initialize database
	if err := app.initializeDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repository
	if err := app.initializeRepository(); err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Run migrations
	if err := app.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize change stream service
	if err := app.initializeChangeStreamService(); err != nil {
		return nil, fmt.Errorf("failed to initialize change stream service: %w", err)
	}

	// Initialize auth service
	if err := app.initializeAuthService(); err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %w", err)
	}

	// Initialize auth handler
	if err := app.initializeAuthHandler(); err != nil {
		return nil, fmt.Errorf("failed to initialize auth handler: %w", err)
	}

	// Initialize security services
	if err := app.initializeSecurityServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize security services: %w", err)
	}

	// Initialize server
	if err := app.initializeServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return app, nil
}

// initializeDatabase initializes the database connection
func (a *App) initializeDatabase() error {
	// Ensure MongoDB driver is set; do not force in-memory in development
	if a.config.IsDevelopment() {
		a.config.Database.Driver = "mongodb"
		// Respect DB_IN_MEMORY env default (false) to keep dev close to prod
	}

	mongoDB, err := pkgDatabase.NewMongoDB(&a.config.Database, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB connection: %w", err)
	}

	a.mongoDB = mongoDB
	a.logger.Info("MongoDB initialized successfully", zap.Bool("in_memory", a.mongoDB.IsInMemory()))
	return nil
}

// initializeRepository initializes the data repository
func (a *App) initializeRepository() error {
	repo, err := database.NewRepository(&a.config.Database, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	a.repository = repo
	a.logger.Info("Repository initialized successfully")
	return nil
}

// runMigrations runs database migrations
func (a *App) runMigrations() error {
	if a.config.Database.InMemory {
		a.logger.Info("Skipping migrations for in-memory database")
		return nil
	}

	migrator := migrations.NewMigrator(a.mongoDB.GetDB(), a.logger)
	if err := migrator.RunMigrations(context.Background()); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.migrator = migrator
	a.logger.Info("Migrations completed successfully")
	return nil
}

// initializeChangeStreamService initializes the change stream service
func (a *App) initializeChangeStreamService() error {
	if a.config.Database.InMemory || a.config.IsDevelopment() || os.Getenv("DISABLE_CHANGE_STREAMS") == "true" {
		a.logger.Info("Skipping change stream service for this environment",
			zap.Bool("in_memory", a.config.Database.InMemory),
			zap.Bool("development", a.config.IsDevelopment()),
			zap.String("disable_env", os.Getenv("DISABLE_CHANGE_STREAMS")),
		)
		return nil
	}

	changeStreamSvc := services.NewChangeStreamService(a.mongoDB.GetDB(), a.logger)
	if err := changeStreamSvc.StartChangeStream(); err != nil {
		// Log and continue in case environment is not a replica set
		a.logger.Warn("Failed to start change stream; continuing without it",
			zap.Error(err),
		)
		return nil
	}

	a.changeStreamSvc = changeStreamSvc
	a.logger.Info("Change stream service initialized successfully")
	return nil
}

// initializeAuthService initializes the authentication service
func (a *App) initializeAuthService() error {
	authSvc, err := auth.NewMongoDBService(a.repository, &a.config.Auth, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create auth service: %w", err)
	}

	a.authSvc = authSvc
	a.logger.Info("Auth service initialized successfully")
	return nil
}

// initializeAuthHandler initializes the authentication handler
func (a *App) initializeAuthHandler() error {
	authHdl, err := auth.NewMongoDBHandler(a.authSvc, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create auth handler: %w", err)
	}

	a.authHdl = authHdl
	a.logger.Info("Auth handler initialized successfully")
	return nil
}

// initializeSecurityServices initializes all security services
func (a *App) initializeSecurityServices() error {
	// Decode encryption keys from base64; in development fallback to raw bytes for DX
	accessSecret, err := base64.StdEncoding.DecodeString(a.config.Auth.JWTAccessSecret)
	if err != nil {
		if a.config.IsDevelopment() {
			accessSecret = []byte(a.config.Auth.JWTAccessSecret)
		} else {
			return fmt.Errorf("failed to decode JWT access secret: %w", err)
		}
	}

	refreshSecret, err := base64.StdEncoding.DecodeString(a.config.Auth.JWTRefreshSecret)
	if err != nil {
		if a.config.IsDevelopment() {
			refreshSecret = []byte(a.config.Auth.JWTRefreshSecret)
		} else {
			return fmt.Errorf("failed to decode JWT refresh secret: %w", err)
		}
	}

	encryptionKey, err := base64.StdEncoding.DecodeString(a.config.Security.EncryptionKey)
	if err != nil {
		if a.config.IsDevelopment() {
			encryptionKey = []byte(a.config.Security.EncryptionKey)
		} else {
			return fmt.Errorf("failed to decode encryption key: %w", err)
		}
	}

	paymentEncryptionKey, err := base64.StdEncoding.DecodeString(a.config.Security.PaymentEncryptionKey)
	if err != nil {
		if a.config.IsDevelopment() {
			paymentEncryptionKey = []byte(a.config.Security.PaymentEncryptionKey)
		} else {
			return fmt.Errorf("failed to decode payment encryption key: %w", err)
		}
	}

	// Initialize JWT refresh service
	jwtService := services.NewJWTRefreshService(accessSecret, refreshSecret, a.logger.Logger)
	// Hook persistent revocation store
	if a.mongoDB != nil && a.mongoDB.GetDB() != nil {
		if store, err := services.NewMongoRevocationStore(a.mongoDB.GetDB(), a.logger); err == nil {
			jwtService.SetRevocationStore(store)
		} else {
			a.logger.Error("Falling back to in-memory revocation store", err)
		}
	}
	a.jwtService = jwtService

	// Initialize encryption service
	encryptionService, err := services.NewEncryptionService(encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to create encryption service: %w", err)
	}
	a.encryptionService = encryptionService

	// Initialize PCI-DSS service
	pciService, err := services.NewPCIDSSService(paymentEncryptionKey, a.logger.Logger)
	if err != nil {
		return fmt.Errorf("failed to create PCI-DSS service: %w", err)
	}
	// Prefer Mongo-backed token store
	if a.mongoDB != nil && a.mongoDB.GetDB() != nil {
		if store, err := services.NewMongoPaymentTokenStore(a.mongoDB.GetDB(), a.logger); err == nil {
			pciService.SetTokenStore(store)
		} else {
			a.logger.Warn("Falling back to in-memory payment token store", zap.Error(err))
		}
	}
	a.pciService = pciService

	// Initialize enhanced auth handler with Mongo-backed service
	authHandler := handlers.NewAuthHandler(jwtService, encryptionService, a.logger, a.authSvc)
	a.authHandler = authHandler

	// Initialize audit service for compliance and security monitoring
	auditService := services.NewAuditService(a.mongoDB.GetDB(), a.logger.Logger)
	a.auditService = auditService

	// Initialize enhanced auth components backed by MongoDB for sessions/devices
	sessionStore := services.NewMongoSessionStore(a.mongoDB.GetDB(), a.logger.Logger)
	deviceStore := services.NewMongoDeviceStore(a.mongoDB.GetDB(), a.logger.Logger)
	enhancedCore := services.NewEnhancedAuthService(a.jwtService, sessionStore, deviceStore, services.NewStubOTPService(a.logger.Logger), a.logger.Logger)

	// Adapt core service to handler interfaces
	adapter := services.NewEnhancedAuthServiceAdapter(enhancedCore, a.repository, a.authSvc, a.logger.Logger)
	userService := services.NewUserServiceAdapter(a.authSvc)
	captchaService := services.NewStubCaptchaService(a.logger.Logger)

	// EMAIL OTP REMOVED: Use stub OTP service (disabled)
	otpSvc := &OtpAdapter{core: services.NewStubOTPService(a.logger.Logger)}
	enhancedAuthHandler := handlers.NewEnhancedAuthHandler(
		adapter,
		userService,
		otpSvc,
		captchaService,
		auditService,
		a.logger.Logger,
	)
	a.enhancedAuthHandler = enhancedAuthHandler

	a.logger.Info("Security services initialized successfully",
		zap.String("jwt_service", "enabled"),
		zap.String("encryption_service", "enabled"),
		zap.String("pci_service", "enabled"),
		zap.String("audit_service", "enabled"),
		zap.String("enhanced_auth_handler", "enabled"),
	)
	return nil
}

// initializeServer initializes the Fiber server
func (a *App) initializeServer() error {
	// Create Fiber app with configuration
	app := fiber.New(fiber.Config{
		AppName:      "Smor-Ting Backend",
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Log the error
			a.logger.Error("HTTP error occurred", err,
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("ip", c.IP()),
			)

			// Return structured error response
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":   "Request failed",
				"message": err.Error(),
				"path":    c.Path(),
			})
		},
	})

	// Middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// CORS middleware
	corsMiddleware, err := middleware.NewCORSMiddleware(&a.config.CORS, a.logger)
	if err != nil {
		a.logger.Error("Failed to create CORS middleware, using default", err)
		app.Use(middleware.DevelopmentCORS(a.logger))
	} else {
		app.Use(corsMiddleware.Configure())
	}

	// Global rate limiter for API and GraphQL routes (exclude health/docs by scoping below)
	apiLimiter := limiter.New(limiter.Config{
		Max:        a.config.Security.RateLimitRequests,
		Expiration: a.config.Security.RateLimitWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Respect Cloudflare's connecting IP header when present
			if cfip := c.Get("CF-Connecting-IP"); cfip != "" {
				return cfip
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	})

	// Auth middleware - JWT-based using Mongo-backed repository
	jwtMiddleware, err := middleware.NewJWTAuthMiddleware(a.jwtService, a.repository, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create jwt auth middleware: %w", err)
	}

	// Audit middleware for compliance and security monitoring
	auditMiddleware := middleware.NewAuditMiddleware(a.auditService, a.logger.Logger)

	// Setup routes
	a.setupRoutes(app, jwtMiddleware, auditMiddleware, apiLimiter)

	a.server = app
	a.logger.Info("Server initialized successfully")
	return nil
}

// setupRoutes sets up all application routes
func (a *App) setupRoutes(app *fiber.App, authMiddleware *middleware.JWTAuthMiddleware, auditMiddleware *middleware.AuditMiddleware, apiLimiter fiber.Handler) {
	// Health check endpoint
	app.Get("/health", a.healthCheck)

	// API documentation endpoints removed with GraphQL

	// Basic WebSocket echo endpoint for real-time connections
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()
		_ = conn.WriteMessage(websocket.TextMessage, []byte("connected"))
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// Echo message back to client
			if err := conn.WriteMessage(msgType, msg); err != nil {
				break
			}
		}
	}))

	// API routes (protected by rate limiter)
	api := app.Group("/api/v1", apiLimiter)
	// KYC routes (auth required in real app; leaving public for test scaffolding)
	kycClient := services.NewSmileIDClient(a.config.KYC.BaseURL, a.config.KYC.PartnerID, a.config.KYC.APIKey)
	kycHandler := handlers.NewKYCHandler(kycClient, a.logger)
	api.Post("/kyc/submit", kycHandler.Submit)
	// Webhooks (no auth)
	webhooks := api.Group("/webhooks")
	ledgerSvc := services.NewWalletLedgerService(a.repository)
	// Attach secure store for encrypted MongoDB Atlas system-of-record
	if a.mongoDB != nil && a.mongoDB.GetDB() != nil && a.encryptionService != nil {
		secure := services.NewWalletLedgerSecureStore(a.mongoDB.GetDB(), a.encryptionService, a.logger)
		// expose via unexported field by assignment
		// go doesn't allow direct assignment to unexported; we are in same package main, but struct in services package
		// provide a small helper to set it
		ledgerSvc = services.AttachSecureStore(ledgerSvc, secure)
	}
	walletWebhook := handlers.NewWalletWebhookHandlerWithLedger(a.logger, ledgerSvc)
	webhooks.Post("/momo", walletWebhook.MomoCallback)

	// Auth routes (no authentication required)
	auth := api.Group("/auth")
	auth.Post("/register", a.authHandler.Register)                  // Using enhanced auth handler for consistency
	auth.Post("/login", a.enhancedAuthHandler.EnhancedLogin)        // Enhanced login to match mobile client expectations
	auth.Post("/validate", a.authHandler.ValidateToken)             // Using enhanced auth handler
	auth.Post("/refresh", a.authHandler.RefreshToken)               // Legacy refresh endpoint
	auth.Post("/refresh-token", a.enhancedAuthHandler.RefreshToken) // Mobile-expected enhanced refresh endpoint
	auth.Post("/revoke", a.authHandler.RevokeToken)                 // New revoke endpoint
	auth.Get("/token-info", a.authHandler.GetTokenInfo)             // New token info endpoint
	// EMAIL OTP ENDPOINTS DISABLED - Removed email OTP functionality
	// auth.Post("/verify-otp", a.authHandler.VerifyOTP)
	// auth.Post("/resend-otp", a.authHandler.ResendOTP)
	auth.Post("/request-password-reset", a.authHandler.RequestPasswordReset)
	auth.Post("/reset-password", a.authHandler.ResetPassword)
	// Biometric authentication endpoint
	auth.Post("/biometric-login", a.enhancedAuthHandler.BiometricLogin)

	// Enhanced auth session management endpoints (require authentication)
	auth.Get("/sessions", authMiddleware.Authenticate(), a.enhancedAuthHandler.GetSessions)
	auth.Delete("/sessions/:id", authMiddleware.Authenticate(), a.enhancedAuthHandler.RevokeSession)
	auth.Delete("/sessions/all", authMiddleware.Authenticate(), a.enhancedAuthHandler.RevokeAllSessions)

	// EMAIL OTP TEST ENDPOINT DISABLED - Removed email OTP functionality 
	// auth.Get("/test/get-latest-otp", a.authHandler.TestGetLatestOTP)

	// Protected routes (authentication required)
	// Apply authentication middleware to each route group individually

	// PROTECTED ROUTES - Apply middleware directly to each route for maximum security

	// Users routes - PROTECTED
	api.Get("/users/profile", authMiddleware.Authenticate(), a.getUserProfile)

	// Services routes - PROTECTED with RBAC and audit logging
	api.Get("/services", authMiddleware.Authenticate(), a.getServices)
	api.Post("/services",
		authMiddleware.RequireRoles(models.ProviderRole, models.AdminRole),
		auditMiddleware.Audit(middleware.AuditConfig{
			Action:   services.ActionServiceCreate,
			Resource: "services",
		}),
		a.createService)
	api.Put("/services/:id",
		authMiddleware.RequireRoles(models.ProviderRole, models.AdminRole),
		auditMiddleware.AuditWithResourceID(services.ActionServiceUpdate, "services", "id"),
		a.updateService)
	api.Delete("/services/:id",
		authMiddleware.RequireRoles(models.AdminRole),
		auditMiddleware.AdminActionAudit(services.ActionServiceDelete, "services"),
		a.deleteService)
	api.Get("/services/:id", authMiddleware.Authenticate(), a.getService)

	// Payment routes - PROTECTED (PCI-DSS compliant) with audit logging
	api.Post("/payments/tokenize",
		authMiddleware.RequireRoles(models.CustomerRole),
		auditMiddleware.AuditSensitiveOperation(services.ActionPaymentProcess, "payments"),
		a.tokenizePaymentMethod)
	api.Post("/payments/process",
		authMiddleware.RequireRoles(models.CustomerRole),
		auditMiddleware.AuditSensitiveOperation(services.ActionPaymentProcess, "payments"),
		a.processPayment)
	api.Get("/payments/validate",
		authMiddleware.RequireRoles(models.CustomerRole, models.AdminRole),
		auditMiddleware.Audit(middleware.AuditConfig{Action: services.ActionPaymentProcess, Resource: "payments"}),
		a.validatePaymentToken)
	api.Delete("/payments/token",
		authMiddleware.RequireRoles(models.CustomerRole, models.AdminRole),
		auditMiddleware.AuditSensitiveOperation(services.ActionPaymentProcess, "payments"),
		a.deletePaymentToken)

	// Sync routes - PROTECTED for offline-first functionality
	api.Post("/sync/data", authMiddleware.Authenticate(), a.syncData)
	api.Get("/sync/unsynced", authMiddleware.Authenticate(), a.getUnsyncedData)
	api.Post("/sync/data/checkpoint", authMiddleware.Authenticate(), a.getUnsyncedDataWithCheckpoint)
	api.Post("/sync/data/chunked", authMiddleware.Authenticate(), a.getChunkedUnsyncedData)
	api.Get("/sync/status/:user_id", authMiddleware.Authenticate(), a.getSyncStatus)
	api.Post("/sync/decompress", authMiddleware.Authenticate(), a.decompressData)

	// Wallet routes - PROTECTED with RBAC and audit logging (sensitive financial operations)
	momoClient := services.NewMomoClient(a.config.Momo.BaseURL, a.config.Momo.TargetEnvironment, a.config.Momo.APIUser, a.config.Momo.APIKey, a.config.Momo.SubscriptionKeyCollection, a.config.Momo.SubscriptionKeyDisbursement)
	walletHandler := handlers.NewWalletHandlerWithLedger(momoClient, a.logger, ledgerSvc)
	api.Post("/wallet/topup",
		authMiddleware.RequireRoles(models.CustomerRole, models.ProviderRole),
		auditMiddleware.AuditSensitiveOperation(services.ActionWalletCreate, "wallet"),
		walletHandler.Topup)
	api.Post("/wallet/pay",
		authMiddleware.RequireRoles(models.CustomerRole),
		auditMiddleware.AuditSensitiveOperation(services.ActionPaymentProcess, "wallet"),
		walletHandler.PayEscrow)
	api.Post("/wallet/withdraw",
		authMiddleware.RequireRoles(models.CustomerRole, models.ProviderRole),
		auditMiddleware.AuditSensitiveOperation(services.ActionPaymentRefund, "wallet"),
		walletHandler.Withdraw)
	api.Get("/wallet/balances",
		authMiddleware.Authenticate(),
		auditMiddleware.Audit(middleware.AuditConfig{Action: services.ActionWalletCreate, Resource: "wallet"}),
		walletHandler.Balances)

	// Composite dashboards - PROTECTED
	dashboardHandler := handlers.NewDashboardHandler(a.repository, ledgerSvc, a.logger)
	api.Get("/home/summary", authMiddleware.Authenticate(), dashboardHandler.HomeSummary)
	api.Get("/wallet/dashboard", authMiddleware.Authenticate(), dashboardHandler.WalletDashboard)

	a.logger.Info("Routes configured successfully")
}

// healthCheck handles health check requests
func (a *App) healthCheck(c *fiber.Ctx) error {
	// Check database health
	dbHealth := "healthy"
	if err := a.mongoDB.HealthCheck(); err != nil {
		dbHealth = "unhealthy"
		a.logger.Error("Database health check failed", err)
	}

	environment := a.environmentLabel()

	return c.JSON(fiber.Map{
		"status":      "healthy",
		"service":     "smor-ting-backend",
		"version":     "1.0.0",
		"timestamp":   time.Now().UTC(),
		"database":    dbHealth,
		"environment": environment,
		"security": fiber.Map{
			"aes_256_encryption": "enabled",
			"jwt_refresh":        "enabled",
			"pci_dss_compliance": "enabled",
		},
	})
}

// environmentLabel returns a human-readable environment string for responses
func (a *App) environmentLabel() string {
	if a.config == nil {
		return "unknown"
	}
	if a.config.IsStaging() {
		return "staging"
	}
	if a.config.IsDevelopment() {
		return "development"
	}
	if a.config.IsProduction() {
		return "production"
	}
	return "unknown"
}

// apiDocs serves API documentation
// API docs and Swagger/Redoc removed with GraphQL cleanup

// Payment handlers
func (a *App) tokenizePaymentMethod(c *fiber.Ctx) error {
	var req services.SensitivePaymentData
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, _ := middleware.GetUserFromContextModels(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID := user.ID.Hex()
	token, err := a.pciService.TokenizePaymentMethod(&req, userID)
	if err != nil {
		a.logger.Error("Failed to tokenize payment method", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to tokenize payment method",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment method tokenized successfully",
		"data":    token,
	})
}

func (a *App) processPayment(c *fiber.Ctx) error {
	var req services.PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := a.pciService.ProcessPayment(&req)
	if err != nil {
		a.logger.Error("Failed to process payment", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process payment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment processed successfully",
		"data":    response,
	})
}

func (a *App) validatePaymentToken(c *fiber.Ctx) error {
	tokenID := c.Query("token_id")
	if tokenID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token ID is required",
		})
	}

	// ensure auth
	if u, _ := middleware.GetUserFromContextModels(c); u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	token, err := a.pciService.ValidatePaymentToken(tokenID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid payment token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment token is valid",
		"data":    token,
	})
}

func (a *App) deletePaymentToken(c *fiber.Ctx) error {
	tokenID := c.Query("token_id")
	if tokenID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token ID is required",
		})
	}

	err := a.pciService.DeletePaymentToken(tokenID)
	if err != nil {
		a.logger.Error("Failed to delete payment token", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete payment token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment token deleted successfully",
	})
}

// Placeholder handlers for protected routes
func (a *App) getUserProfile(c *fiber.Ctx) error {
	user, ok := middleware.GetUserFromContextModels(c)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "User not found in context",
		})
	}
	return c.JSON(fiber.Map{
		"message": "User profile endpoint",
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"role":    user.Role,
	})
}

func (a *App) getServices(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Services list endpoint",
		"data":    []fiber.Map{},
	})
}

func (a *App) createService(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Create service endpoint",
	})
}

func (a *App) getService(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Get service endpoint",
		"id":      id,
	})
}

func (a *App) updateService(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Update service endpoint",
		"id":      id,
	})
}

func (a *App) deleteService(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Delete service endpoint",
		"id":      id,
	})
}

// Enhanced sync handlers for offline-first functionality
func (a *App) syncData(c *fiber.Ctx) error {
	var req models.SyncRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate user ID
	if req.UserID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Use sync service for enhanced functionality
	syncService := services.NewSyncService(a.repository, a.auditService, a.logger.Logger)
	response, err := syncService.GetUnsyncedDataWithCheckpoint(c.Context(), &req)
	if err != nil {
		a.logger.Error("Failed to sync data", err,
			zap.String("user_id", req.UserID.Hex()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sync data",
		})
	}

	return c.JSON(response)
}

func (a *App) getUnsyncedData(c *fiber.Ctx) error {
	// Get user ID from query params or body
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Get last sync time from query params
	lastSyncStr := c.Query("last_sync_at")
	var lastSyncAt time.Time
	if lastSyncStr != "" {
		lastSyncAt, err = time.Parse(time.RFC3339, lastSyncStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid last_sync_at format",
			})
		}
	} else {
		lastSyncAt = time.Now().Add(-24 * time.Hour) // Default to 24 hours ago
	}

	// Get unsynced data
	data, err := a.repository.GetUnsyncedData(c.Context(), userID, lastSyncAt)
	if err != nil {
		a.logger.Error("Failed to get unsynced data", err,
			zap.String("user_id", userID.Hex()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unsynced data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Unsynced data retrieved successfully",
		"data":    data,
		"user_id": userID.Hex(),
	})
}

// New enhanced sync endpoints
func (a *App) getUnsyncedDataWithCheckpoint(c *fiber.Ctx) error {
	var req models.SyncRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate user ID
	if req.UserID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Use sync service for enhanced functionality
	syncService := services.NewSyncService(a.repository, a.auditService, a.logger.Logger)
	response, err := syncService.GetUnsyncedDataWithCheckpoint(c.Context(), &req)
	if err != nil {
		a.logger.Error("Failed to get unsynced data with checkpoint", err,
			zap.String("user_id", req.UserID.Hex()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unsynced data",
		})
	}

	return c.JSON(response)
}

func (a *App) getChunkedUnsyncedData(c *fiber.Ctx) error {
	var req models.ChunkedSyncRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate user ID
	if req.UserID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Use sync service for enhanced functionality
	syncService := services.NewSyncService(a.repository, a.auditService, a.logger.Logger)
	response, err := syncService.GetChunkedUnsyncedData(c.Context(), &req)
	if err != nil {
		a.logger.Error("Failed to get chunked unsynced data", err,
			zap.String("user_id", req.UserID.Hex()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get chunked data",
		})
	}

	return c.JSON(response)
}

func (a *App) getSyncStatus(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Use sync service for enhanced functionality
	syncService := services.NewSyncService(a.repository, a.auditService, a.logger.Logger)
	status, err := syncService.GetSyncStatus(c.Context(), userID)
	if err != nil {
		a.logger.Error("Failed to get sync status", err,
			zap.String("user_id", userID.Hex()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get sync status",
		})
	}

	return c.JSON(status)
}

func (a *App) decompressData(c *fiber.Ctx) error {
	// Get compressed data from request body
	var req struct {
		CompressedData string `json:"compressed_data"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.CompressedData == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Compressed data is required",
		})
	}

	// Decode base64 compressed data
	compressedBytes, err := base64.StdEncoding.DecodeString(req.CompressedData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid compressed data format",
		})
	}

	// Decompress data
	syncService := services.NewSyncService(a.repository, a.auditService, a.logger.Logger)
	decompressedData, err := syncService.DecompressData(compressedBytes)
	if err != nil {
		a.logger.Error("Failed to decompress data", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decompress data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data decompressed successfully",
		"data":    decompressedData,
	})
}

// Start starts the application
func (a *App) Start() error {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go a.handleGracefulShutdown(ctx, cancel)

	// Start server
	addr := fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port)
	environment := "production"
	if a.config.IsDevelopment() {
		environment = "development"
	}

	a.logger.Info("Starting server",
		zap.String("address", addr),
		zap.String("environment", environment),
		zap.Bool("in_memory_db", a.mongoDB.IsInMemory()),
		zap.String("security", "enabled"),
	)

	if err := a.server.Listen(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// handleGracefulShutdown handles graceful shutdown
func (a *App) handleGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	a.logger.Info("Received shutdown signal, starting graceful shutdown")

	// Cancel context
	cancel()

	// Shutdown server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := a.server.ShutdownWithContext(shutdownCtx); err != nil {
		a.logger.Error("Failed to shutdown server gracefully", err)
	}

	// Stop change stream service
	if a.changeStreamSvc != nil {
		if err := a.changeStreamSvc.StopChangeStream(); err != nil {
			a.logger.Error("Failed to stop change stream service", err)
		}
	}

	// Close repository
	if err := a.repository.Close(); err != nil {
		a.logger.Error("Failed to close repository", err)
	}

	// Close database connection
	if err := a.mongoDB.Close(); err != nil {
		a.logger.Error("Failed to close database connection", err)
	}

	// Sync logger
	if err := a.logger.Sync(); err != nil {
		log.Printf("Failed to sync logger: %v", err)
	}

	a.logger.Info("Graceful shutdown completed")
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Create application
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Start application
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
