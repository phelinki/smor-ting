package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/handlers"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
)

func TestWalletMethods_LinkListDelete(t *testing.T) {
	repo := database.NewMemoryDatabase()
	lg, _ := logger.New("debug", "console", "stdout")
	h := handlers.NewWalletMethodsHandler(repo, lg)

	// Seed user
	u := &models.User{Email: "msisdn@example.com"}
	_ = repo.CreateUser(nil, u)

	app := fiber.New()
	app.Post("/wallet/methods", func(c *fiber.Ctx) error {
		c.Locals("user", u)
		return h.Link(c)
	})
	app.Get("/wallet/methods", func(c *fiber.Ctx) error {
		c.Locals("user", u)
		return h.List(c)
	})
	app.Delete("/wallet/methods/:id", func(c *fiber.Ctx) error {
		c.Locals("user", u)
		return h.Delete(c)
	})

	// Link
	body := map[string]interface{}{"type": "mobile_money", "provider": "lonestar", "msisdn": "231770000000", "default": true}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/wallet/methods", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("link expected 200, got %d", resp.StatusCode)
	}

	// List
	req2 := httptest.NewRequest(http.MethodGet, "/wallet/methods", nil)
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("list expected 200, got %d", resp2.StatusCode)
	}

	var got struct {
		Data []models.LinkedPaymentMethod `json:"data"`
	}
	_ = json.NewDecoder(resp2.Body).Decode(&got)
	if len(got.Data) != 1 || got.Data[0].Provider != "lonestar" {
		t.Fatalf("unexpected list: %+v", got)
	}

	// Delete
	id := got.Data[0].ID.Hex()
	req3 := httptest.NewRequest(http.MethodDelete, "/wallet/methods/"+id, nil)
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("delete expected 200, got %d", resp3.StatusCode)
	}
}
