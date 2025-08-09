package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletMethodsHandler struct {
	repo    database.Repository
	logger  *logger.Logger
	methods map[string][]models.LinkedPaymentMethod
}

func NewWalletMethodsHandler(repo database.Repository, logger *logger.Logger) *WalletMethodsHandler {
	return &WalletMethodsHandler{repo: repo, logger: logger, methods: make(map[string][]models.LinkedPaymentMethod)}
}

type linkReq struct {
	Type     string `json:"type"`
	Provider string `json:"provider"`
	Msisdn   string `json:"msisdn"`
	Default  bool   `json:"default"`
}

func (h *WalletMethodsHandler) Link(c *fiber.Ctx) error {
	var req linkReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.Provider == "" || req.Msisdn == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "provider and msisdn required"})
	}
	user, _ := c.Locals("user").(*models.User)
	method := models.LinkedPaymentMethod{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		Type:      models.PaymentMethodMobileMoney,
		Provider:  req.Provider,
		Msisdn:    maskMsisdn(req.Msisdn),
		IsDefault: req.Default,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Store in-memory by user
	uid := user.ID.Hex()
	h.methods[uid] = append(h.methods[uid], method)
	return c.JSON(fiber.Map{"data": []models.LinkedPaymentMethod{method}})
}

func (h *WalletMethodsHandler) List(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	return c.JSON(fiber.Map{"data": h.methods[user.ID.Hex()]})
}

func (h *WalletMethodsHandler) Delete(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	id := c.Params("id")
	list := h.methods[user.ID.Hex()]
	out := make([]models.LinkedPaymentMethod, 0, len(list))
	for _, m := range list {
		if m.ID.Hex() != id {
			out = append(out, m)
		}
	}
	h.methods[user.ID.Hex()] = out
	return c.JSON(fiber.Map{"status": "deleted"})
}

func maskMsisdn(msisdn string) string {
	if len(msisdn) < 4 {
		return msisdn
	}
	return "***" + msisdn[len(msisdn)-4:]
}
