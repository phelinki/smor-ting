package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
)

// DashboardHandler serves composite, summary-style endpoints for the home and wallet screens
type DashboardHandler struct {
	repo   database.Repository
	ledger *services.WalletLedgerService
	logger *logger.Logger
}

func NewDashboardHandler(repo database.Repository, ledger *services.WalletLedgerService, logger *logger.Logger) *DashboardHandler {
	return &DashboardHandler{repo: repo, ledger: ledger, logger: logger}
}

// HomeSummary aggregates a minimal set of information for the app home screen
// Response shape is stable and optimized for mobile consumption
func (h *DashboardHandler) HomeSummary(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	bookings, _ := h.repo.GetUserBookings(c.Context(), user.ID)

	now := time.Now()
	var upcoming, active, completed int
	for _, b := range bookings {
		switch b.Status {
		case models.BookingInProgress:
			active++
		case models.BookingCompleted:
			completed++
		case models.BookingPending, models.BookingConfirmed:
			if b.ScheduledDate.After(now) {
				upcoming++
			}
		}
	}

	// Wallet balances (fallback to embedded wallet if ledger unavailable)
	wallet := fiber.Map{
		"available": user.Wallet.Balance,
		"pending":   0.0,
		"total":     user.Wallet.Balance,
		"currency":  user.Wallet.Currency,
	}
	if h.ledger != nil {
		if bal, err := h.ledger.ComputeBalances(c.Context(), user.ID); err == nil {
			wallet = fiber.Map{
				"available": bal.Available,
				"pending":   bal.PendingHeld,
				"total":     bal.Total,
				"currency":  bal.Currency,
			}
		}
	}

	return c.JSON(fiber.Map{
		"greeting": "Welcome back, " + user.FirstName,
		"stats": fiber.Map{
			"upcoming_bookings": upcoming,
			"active_jobs":       active,
			"completed_jobs":    completed,
		},
		"wallet": wallet,
	})
}

// WalletDashboard returns balances and a compact, recent transaction list for the wallet screen
func (h *DashboardHandler) WalletDashboard(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Refresh balances from ledger when available
	wallet := fiber.Map{
		"available": user.Wallet.Balance,
		"pending":   0.0,
		"total":     user.Wallet.Balance,
		"currency":  user.Wallet.Currency,
	}
	if h.ledger != nil {
		if bal, err := h.ledger.ComputeBalances(c.Context(), user.ID); err == nil {
			wallet = fiber.Map{
				"available": bal.Available,
				"pending":   bal.PendingHeld,
				"total":     bal.Total,
				"currency":  bal.Currency,
			}
		}
	}

	// Pull latest user to ensure we have recent transactions state
	latestUser, err := h.repo.GetUserByID(c.Context(), user.ID)
	if err == nil {
		user = latestUser
	}

	txs := user.Wallet.Transactions
	sort.Slice(txs, func(i, j int) bool { return txs[i].CreatedAt.After(txs[j].CreatedAt) })
	if len(txs) > 10 {
		txs = txs[:10]
	}

	// Project a compact transaction representation
	recent := make([]fiber.Map, 0, len(txs))
	for _, t := range txs {
		recent = append(recent, fiber.Map{
			"id":          t.ID.Hex(),
			"type":        t.Type,
			"amount":      t.Amount,
			"status":      t.Status,
			"reference":   t.Reference,
			"description": t.Description,
			"created_at":  t.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"wallet":       wallet,
		"transactions": recent,
	})
}
