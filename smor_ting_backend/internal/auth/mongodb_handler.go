package auth

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"go.uber.org/zap"
)

// MongoDBHandler represents the authentication HTTP handler for MongoDB
type MongoDBHandler struct {
	service *MongoDBService
	logger  *logger.Logger
}

// NewMongoDBHandler creates a new MongoDB authentication handler
func NewMongoDBHandler(service *MongoDBService, logger *logger.Logger) (*MongoDBHandler, error) {
	if service == nil {
		return nil, fmt.Errorf("auth service is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &MongoDBHandler{
		service: service,
		logger:  logger,
	}, nil
}

// Register handles user registration
func (h *MongoDBHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse register request body", err)
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: "Failed to parse request body",
		})
	}

	// Validate request
	if err := h.validateRegisterRequest(&req); err != nil {
		h.logger.Warn("Invalid register request", zap.String("error", err.Error()))
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	// Register user
	response, err := h.service.Register(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to register user", err, zap.String("email", req.Email))

		// Handle specific errors
		if err.Error() == fmt.Sprintf("user with email %s already exists", req.Email) {
			return c.Status(http.StatusConflict).JSON(ErrorResponse{
				Error:   "User already exists",
				Message: "A user with this email already exists",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Registration failed",
			Message: "Failed to register user",
		})
	}

	h.logger.Info("User registered successfully", zap.String("email", req.Email))
	return c.Status(http.StatusCreated).JSON(response)
}

// Login handles user login
func (h *MongoDBHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse login request body", err)
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: "Failed to parse request body",
		})
	}

	// Validate request
	if err := h.validateLoginRequest(&req); err != nil {
		h.logger.Warn("Invalid login request", zap.String("error", err.Error()))
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	// Login user
	response, err := h.service.Login(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to login user", err, zap.String("email", req.Email))

		// Handle specific errors
		if err.Error() == "invalid credentials" {
			return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
				Error:   "Invalid credentials",
				Message: "Invalid email or password",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Login failed",
			Message: "Failed to authenticate user",
		})
	}

	h.logger.Info("User logged in successfully", zap.String("email", req.Email))
	return c.Status(http.StatusOK).JSON(response)
}

// ValidateToken validates a JWT token
func (h *MongoDBHandler) ValidateToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse token validation request", err)
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: "Failed to parse request body",
		})
	}

	// Validate token
	user, err := h.service.ValidateToken(req.Token)
	if err != nil {
		h.logger.Error("Failed to validate token", err)
		return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Invalid token",
			Message: "The provided token is invalid or expired",
		})
	}

	h.logger.Info("Token validated successfully", zap.String("user_id", user.ID.Hex()))
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"valid": true,
		"user":  user,
	})
}

// validateRegisterRequest validates a registration request
func (h *MongoDBHandler) validateRegisterRequest(req *models.RegisterRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if req.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if req.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if req.Phone == "" {
		return fmt.Errorf("phone is required")
	}
	if req.Role == "" {
		return fmt.Errorf("role is required")
	}
	// Validate role is one of the allowed values
	if req.Role != models.CustomerRole && req.Role != models.ProviderRole && req.Role != models.AdminRole {
		return fmt.Errorf("role must be 'customer', 'provider', or 'admin'")
	}
	return nil
}

// validateLoginRequest validates a login request
func (h *MongoDBHandler) validateLoginRequest(req *models.LoginRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
