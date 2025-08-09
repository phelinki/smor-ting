package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WalletWebhookHandler processes provider callbacks (Lonestar/MoMo)
type WalletWebhookHandler struct {
	logger *logger.Logger
	ledger *services.WalletLedgerService
}

func NewWalletWebhookHandler(logger *logger.Logger) *WalletWebhookHandler {
	return &WalletWebhookHandler{logger: logger}
}
func NewWalletWebhookHandlerWithLedger(logger *logger.Logger, ledger *services.WalletLedgerService) *WalletWebhookHandler {
	return &WalletWebhookHandler{logger: logger, ledger: ledger}
}

// Callback payload is provider-specific; accept generic map and log securely (no PII)
func (h *WalletWebhookHandler) MomoCallback(c *fiber.Ctx) error {
	var payload struct {
		Type        string  `json:"type"`
		Status      string  `json:"status"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		UserID      string  `json:"user_id"`
		ReferenceID string  `json:"referenceId"`
		ProviderRef string  `json:"provider_ref"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if h.ledger != nil && payload.UserID != "" {
		if uid, err := primitive.ObjectIDFromHex(payload.UserID); err == nil {
			entry := &models.WalletLedgerEntry{
				UserID:      uid,
				Amount:      payload.Amount,
				Currency:    payload.Currency,
				Reference:   payload.ReferenceID,
				ProviderRef: payload.ProviderRef,
				Direction:   models.LedgerCredit,
				Status:      models.LedgerPending,
			}
			switch payload.Type {
			case "topup":
				entry.Type = models.LedgerTopup
			case "payment":
				entry.Type = models.LedgerPayment
			case "escrow_hold":
				entry.Type = models.LedgerEscrowHold
				entry.IsEscrow = true
			case "escrow_release":
				entry.Type = models.LedgerEscrowRelease
			case "withdraw":
				entry.Type = models.LedgerWithdraw
				entry.Direction = models.LedgerDebit
			}
			if payload.Status == "SUCCESSFUL" {
				entry.Status = models.LedgerCompleted
				// Escrow holds should remain Pending (counted in pending/held balance)
				if entry.Type == models.LedgerEscrowHold {
					entry.Status = models.LedgerPending
				}
			} else if payload.Status == "FAILED" {
				entry.Status = models.LedgerFailed
			}
			_ = h.ledger.RecordEntry(c.Context(), entry)
		}
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
