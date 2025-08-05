package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
<<<<<<< HEAD
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/auth"
	"github.com/smorting/backend/pkg/database"
	"github.com/smorting/backend/pkg/logger"
	"github.com/smorting/backend/pkg/middleware"
	"go.uber.org/zap"
)

// App represents the application
type App struct {
	config  *configs.Config
	logger  *logger.Logger
	db      *database.Database
	authSvc *auth.Service
	authHdl *auth.Handler
	server  *fiber.App
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

	// Initialize auth service
	if err := app.initializeAuthService(); err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %w", err)
	}

	// Initialize auth handler
	if err := app.initializeAuthHandler(); err != nil {
		return nil, fmt.Errorf("failed to initialize auth handler: %w", err)
	}

	// Initialize server
	if err := app.initializeServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return app, nil
}

// initializeDatabase initializes the database connection
func (a *App) initializeDatabase() error {
	db, err := database.New(&a.config.Database, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	a.db = db
	a.logger.Info("Database initialized successfully", zap.Bool("in_memory", a.db.IsInMemory()))
	return nil
}

// initializeAuthService initializes the authentication service
func (a *App) initializeAuthService() error {
	authSvc, err := auth.NewService(a.db.GetDB(), &a.config.Auth, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create auth service: %w", err)
	}

	a.authSvc = authSvc
	a.logger.Info("Auth service initialized successfully")
	return nil
}

// initializeAuthHandler initializes the authentication handler
func (a *App) initializeAuthHandler() error {
	authHdl, err := auth.NewHandler(a.authSvc, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create auth handler: %w", err)
	}

	a.authHdl = authHdl
	a.logger.Info("Auth handler initialized successfully")
	return nil
}

// initializeServer initializes the Fiber server
func (a *App) initializeServer() error {
	// Create Fiber app with configuration
=======
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/services"
)

func main() {
	// Initialize in-memory database
	db := database.NewMemoryDatabase()
	defer db.Close()

	// Initialize services
	emailService := services.NewEmailService()
	authService := services.NewAuthService(db, emailService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Create new Fiber app
>>>>>>> origin/cursor/check-previous-ask-status-7bc9
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

	// Auth middleware
	authMiddleware, err := middleware.NewAuthMiddleware(a.authSvc, a.logger)
	if err != nil {
		return fmt.Errorf("failed to create auth middleware: %w", err)
	}

	// Setup routes
	a.setupRoutes(app, authMiddleware)

	a.server = app
	a.logger.Info("Server initialized successfully")
	return nil
}

// setupRoutes sets up all application routes
func (a *App) setupRoutes(app *fiber.App, authMiddleware *middleware.AuthMiddleware) {
	// Health check endpoint
	app.Get("/health", a.healthCheck)

	// API documentation
	app.Get("/docs", a.apiDocs)
	app.Get("/swagger", a.swaggerDocs)

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (no authentication required)
	auth := api.Group("/auth")
	auth.Post("/register", a.authHdl.Register)
	auth.Post("/login", a.authHdl.Login)
	auth.Post("/validate", a.authHdl.ValidateToken)
	auth.Post("/verify-otp", a.authHdl.VerifyOTP)
	auth.Post("/resend-otp", a.authHdl.ResendOTP)

	// Protected routes (authentication required)
	protected := api.Group("/")
	protected.Use(authMiddleware.Authenticate())

	// Users routes
	users := protected.Group("/users")
	users.Get("/profile", a.getUserProfile)

	// Services routes
	services := protected.Group("/services")
	services.Get("/", a.getServices)
	services.Post("/", a.createService)
	services.Get("/:id", a.getService)
	services.Put("/:id", a.updateService)
	services.Delete("/:id", a.deleteService)

	a.logger.Info("Routes configured successfully")
}

// healthCheck handles health check requests
func (a *App) healthCheck(c *fiber.Ctx) error {
	// Check database health
	dbHealth := "healthy"
	if err := a.db.HealthCheck(); err != nil {
		dbHealth = "unhealthy"
		a.logger.Error("Database health check failed", err)
	}

	environment := "production"
	if a.config.IsDevelopment() {
		environment = "development"
	}

	return c.JSON(fiber.Map{
		"status":      "healthy",
		"service":     "smor-ting-backend",
		"version":     "1.0.0",
		"timestamp":   time.Now().UTC(),
		"database":    dbHealth,
		"environment": environment,
	})
}

// apiDocs serves API documentation
func (a *App) apiDocs(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "API Documentation",
		"version": "1.0.0",
		"endpoints": fiber.Map{
			"health": "/health",
			"auth": fiber.Map{
				"register": "POST /api/v1/auth/register",
				"login":    "POST /api/v1/auth/login",
				"validate": "POST /api/v1/auth/validate",
				"verify-otp": "POST /api/v1/auth/verify-otp",
				"resend-otp": "POST /api/v1/auth/resend-otp",
			},
			"users": fiber.Map{
				"profile": "GET /api/v1/users/profile",
			},
			"services": fiber.Map{
				"list":   "GET /api/v1/services",
				"create": "POST /api/v1/services",
				"get":    "GET /api/v1/services/:id",
				"update": "PUT /api/v1/services/:id",
				"delete": "DELETE /api/v1/services/:id",
			},
		},
	})
}

// swaggerDocs serves Swagger documentation
func (a *App) swaggerDocs(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Swagger Documentation",
		"swagger": "2.0",
		"info": fiber.Map{
			"title":       "Smor-Ting Backend API",
			"version":     "1.0.0",
			"description": "A robust, production-ready backend API for the Smor-Ting platform",
		},
		"host":     fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port),
		"basePath": "/api/v1",
		"schemes":  []string{"http", "https"},
	})
}

// Placeholder handlers for protected routes
func (a *App) getUserProfile(c *fiber.Ctx) error {
	user, _ := middleware.GetUserFromContext(c)
	return c.JSON(fiber.Map{
		"message": "User profile endpoint",
		"user_id": user.ID,
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
		zap.Bool("in_memory_db", a.db.IsInMemory()),
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

	// Close database connection
	if err := a.db.Close(); err != nil {
		a.logger.Error("Failed to close database connection", err)
	}

	// Sync logger
	if err := a.logger.Sync(); err != nil {
		log.Printf("Failed to sync logger: %v", err)
	}

	a.logger.Info("Graceful shutdown completed")
}

func main() {
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
