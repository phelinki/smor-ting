package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	app := fiber.New(fiber.Config{
		AppName: "Smor-Ting Backend",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "smor-ting-backend",
			"version": "1.0.0",
		})
	})

	// API routes
	api := app.Group("/api/v1")
	
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/verify-otp", authHandler.VerifyOTP)
	auth.Post("/resend-otp", authHandler.ResendOTP)

	// Services routes
	services := api.Group("/services")
	services.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Services list endpoint"})
	})

	// Users routes
	users := api.Group("/users")
	users.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "User profile endpoint"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Smor-Ting backend server on port %s", port)
	log.Fatal(app.Listen(":" + port))
} 