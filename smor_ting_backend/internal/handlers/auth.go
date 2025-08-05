package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = models.CustomerRole
	}

	response, err := h.authService.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	response, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req models.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Email == "" || req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and OTP are required",
		})
	}

	if len(req.OTP) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP must be 6 digits",
		})
	}

	response, err := h.authService.VerifyOTP(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *AuthHandler) ResendOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	// For simplicity, we'll trigger a login attempt which will send OTP if email is not verified
	loginReq := models.LoginRequest{
		Email:    req.Email,
		Password: "dummy", // This will fail but trigger OTP if user exists and is not verified
	}

	// This is a simplified approach - in production you might want a dedicated resend endpoint
	_, err := h.authService.Login(c.Context(), loginReq)
	if err != nil {
		// Check if it's a password error but user exists
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Please try logging in with your correct password to receive OTP",
		})
	}

	return c.JSON(fiber.Map{
		"message": "OTP sent successfully",
	})
}