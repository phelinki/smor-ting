package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
)

type KYCHandler struct {
	client *services.SmileIDClient
	logger *logger.Logger
}

func NewKYCHandler(client *services.SmileIDClient, logger *logger.Logger) *KYCHandler {
	return &KYCHandler{client: client, logger: logger}
}

func (h *KYCHandler) Submit(c *fiber.Ctx) error {
	var req services.KYCRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	res, err := h.client.SubmitKYC(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "kyc submission failed"})
	}
	return c.JSON(fiber.Map{"status": res.Status, "reference": res.Reference})
}
