package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/internal/services"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// test context helper to inject a user into Locals
func withUser(app *fiber.App, u *models.User) *fiber.App {
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", u)
		return c.Next()
	})
	return app
}

func TestHomeSummary_ReturnsStatsAndWallet(t *testing.T) {
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()

	// create a user with some bookings and wallet tx
	user := &models.User{Email: "u@test.com", FirstName: "John", LastName: "Doe"}
	if err := repo.CreateUser(context.TODO(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	// add bookings: 1 upcoming, 1 active, 1 completed
	upcoming := &models.Booking{CustomerID: user.ID, Status: models.BookingConfirmed, ScheduledDate: time.Now().Add(48 * time.Hour)}
	active := &models.Booking{CustomerID: user.ID, Status: models.BookingInProgress, ScheduledDate: time.Now().Add(-1 * time.Hour)}
	completed := &models.Booking{CustomerID: user.ID, Status: models.BookingCompleted, ScheduledDate: time.Now().Add(-48 * time.Hour)}
	_ = repo.CreateBooking(context.TODO(), upcoming)
	_ = repo.CreateBooking(context.TODO(), active)
	_ = repo.CreateBooking(context.TODO(), completed)

	// add wallet tx
	_ = repo.UpdateWallet(context.TODO(), user.ID, &models.Transaction{Type: "credit", Amount: 50, Status: "completed", Reference: "t1"})

	ledger := services.NewWalletLedgerService(repo)
	h := handlers.NewDashboardHandler(repo, ledger, lg)

	app := withUser(fiber.New(), user)
	app.Get("/home-summary", h.HomeSummary)

	req := httptest.NewRequest(http.MethodGet, "/home-summary", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)

	stats := body["stats"].(map[string]interface{})
	if int(stats["upcoming_bookings"].(float64)) != 1 {
		t.Fatalf("expected 1 upcoming, got %v", stats["upcoming_bookings"])
	}
	if int(stats["active_jobs"].(float64)) != 1 {
		t.Fatalf("expected 1 active, got %v", stats["active_jobs"])
	}
	if int(stats["completed_jobs"].(float64)) != 1 {
		t.Fatalf("expected 1 completed, got %v", stats["completed_jobs"])
	}
}

func TestWalletDashboard_RecentTransactionsAndBalances(t *testing.T) {
	lg, _ := logger.New("debug", "console", "stdout")
	repo := database.NewMemoryDatabase()

	user := &models.User{Email: "u2@test.com", FirstName: "Jane"}
	if err := repo.CreateUser(context.TODO(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	// create 12 transactions to test trimming to 10
	for i := 0; i < 12; i++ {
		_ = repo.UpdateWallet(context.TODO(), user.ID, &models.Transaction{
			ID:          primitive.NewObjectID(),
			Type:        "credit",
			Amount:      float64(10 + i),
			Description: "tx",
			Reference:   "r",
			Status:      "completed",
			CreatedAt:   time.Now().Add(time.Duration(i) * time.Minute),
		})
	}

	ledger := services.NewWalletLedgerService(repo)
	h := handlers.NewDashboardHandler(repo, ledger, lg)
	app := withUser(fiber.New(), user)
	app.Get("/wallet-dashboard", h.WalletDashboard)

	req := httptest.NewRequest(http.MethodGet, "/wallet-dashboard", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	txs := body["transactions"].([]interface{})
	if len(txs) != 10 {
		t.Fatalf("expected 10 recent transactions, got %d", len(txs))
	}
}
