//go:build legacy_sql_auth

package auth

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/pkg/logger"
	"go.uber.org/zap"
)

// Handler represents the authentication HTTP handler
type Handler struct {
	service *Service
	logger  *logger.Logger
}

// NewHandler creates a new authentication handler
func NewHandler(service *Service, logger *logger.Logger) (*Handler, error) {
	if service == nil {
		return nil, fmt.Errorf("auth service is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &Handler{
		service: service,
		logger:  logger,
	}, nil
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email, password, first name, and last name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

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
	response, err := h.service.Register(&req)
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
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest

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
	response, err := h.service.Login(&req)
	if err != nil {
		h.logger.Error("Failed to login user", err, zap.String("email", req.Email))

		// Handle specific errors
		if err.Error() == "invalid credentials" {
			return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
				Error:   "Invalid credentials",
				Message: "Email or password is incorrect",
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
// @Summary Validate JWT token
// @Description Validate a JWT token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} User
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/validate [post]
func (h *Handler) ValidateToken(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		h.logger.Warn("Missing Authorization header")
		return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Missing token",
			Message: "Authorization header is required",
		})
	}

	// Extract token from "Bearer <token>" format
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		h.logger.Warn("Invalid Authorization header format")
		return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Invalid token format",
			Message: "Token must be in 'Bearer <token>' format",
		})
	}

	// Validate token
	user, err := h.service.ValidateToken(token)
	if err != nil {
		h.logger.Error("Failed to validate token", err)
		return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Invalid token",
			Message: "Token is invalid or expired",
		})
	}

	h.logger.Info("Token validated successfully", zap.Int("user_id", user.ID))
	return c.Status(http.StatusOK).JSON(user)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// validateRegisterRequest validates a registration request
func (h *Handler) validateRegisterRequest(req *RegisterRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if req.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if req.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	return nil
}

// validateLoginRequest validates a login request
func (h *Handler) validateLoginRequest(req *LoginRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
