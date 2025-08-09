package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
)

// WalletHandler exposes wallet endpoints backed by MoMo
type WalletHandler struct {
	momo   services.MomoAPI
	logger *logger.Logger
	ledger *services.WalletLedgerService
}

func NewWalletHandler(momo services.MomoAPI, logger *logger.Logger) *WalletHandler {
	return &WalletHandler{momo: momo, logger: logger}
}
func NewWalletHandlerWithLedger(momo services.MomoAPI, logger *logger.Logger, ledger *services.WalletLedgerService) *WalletHandler {
	return &WalletHandler{momo: momo, logger: logger, ledger: ledger}
}

// Online-only guard wrapper
func (h *WalletHandler) ensureOnline(c *fiber.Ctx) error {
	if err := h.momo.EnsureOnline(c.Context()); err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error":   "Wallet requires internet",
			"message": "Please connect to the internet to use wallet",
		})
	}
	return nil
}

type TopupRequest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Msisdn   string `json:"msisdn"`
}

func (h *WalletHandler) Topup(c *fiber.Ctx) error {
	if err := h.ensureOnline(c); err != nil {
		return err
	}
	var req TopupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	body := services.RequestToPay{
		Amount:       req.Amount,
		Currency:     req.Currency,
		ExternalId:   "wallet_topup",
		Payer:        services.Party{PartyIdType: "MSISDN", PartyId: req.Msisdn},
		PayerMessage: "Wallet topup",
		PayeeNote:    "Smor-Ting Wallet",
	}
	ref, err := h.momo.RequestToPay(c.Context(), body)
	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "topup failed"})
	}
	return c.JSON(fiber.Map{"reference_id": ref, "status": "PENDING"})
}

type PayRequest struct {
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Msisdn     string `json:"msisdn"`
	BookingRef string `json:"booking_ref"`
}

func (h *WalletHandler) PayEscrow(c *fiber.Ctx) error {
	if err := h.ensureOnline(c); err != nil {
		return err
	}
	var req PayRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	// Initiate R2P into escrow (handled by our ledger on webhook confirmation)
	body := services.RequestToPay{Amount: req.Amount, Currency: req.Currency, ExternalId: req.BookingRef, Payer: services.Party{PartyIdType: "MSISDN", PartyId: req.Msisdn}, PayerMessage: "Task escrow", PayeeNote: req.BookingRef}
	ref, err := h.momo.RequestToPay(c.Context(), body)
	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "payment failed"})
	}
	return c.JSON(fiber.Map{"reference_id": ref, "status": "PENDING"})
}

type WithdrawRequest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Msisdn   string `json:"msisdn"`
}

func (h *WalletHandler) Withdraw(c *fiber.Ctx) error {
	if err := h.ensureOnline(c); err != nil {
		return err
	}
	var req WithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	ref, err := h.momo.Transfer(c.Context(), services.TransferRequest{Amount: req.Amount, Currency: req.Currency, ExternalId: "wallet_withdraw", Payee: services.Party{PartyIdType: "MSISDN", PartyId: req.Msisdn}})
	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "withdrawal failed"})
	}
	return c.JSON(fiber.Map{"reference_id": ref, "status": "PENDING"})
}

// Placeholder: balances and transactions will pull from repository
func (h *WalletHandler) Balances(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	if h.ledger != nil {
		if bal, err := h.ledger.ComputeBalances(c.Context(), user.ID); err == nil {
			return c.JSON(bal)
		}
	}
	return c.JSON(fiber.Map{"available": user.Wallet.Balance, "pending": 0, "total": user.Wallet.Balance, "currency": user.Wallet.Currency})
}
