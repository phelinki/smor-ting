package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap/zaptest"
)

func TestConsentHandler_GetRequirements(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	handler := NewConsentHandler(logger)

	app := fiber.New()
	app.Get("/api/v1/consent/requirements", handler.GetRequirements)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/v1/consent/requirements", nil)
	resp, err := app.Test(req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	requirements, ok := response["requirements"].([]interface{})
	if !ok {
		t.Fatal("Expected requirements array in response")
	}

	if len(requirements) == 0 {
		t.Error("Expected at least one consent requirement")
	}

	// Check that default requirements include required ones
	hasTerms := false
	hasPrivacy := false
	for _, req := range requirements {
		reqMap := req.(map[string]interface{})
		reqType := reqMap["type"].(string)
		if reqType == string(models.ConsentTypeTermsOfService) {
			hasTerms = true
			if !reqMap["required"].(bool) {
				t.Error("Terms of Service should be required")
			}
		}
		if reqType == string(models.ConsentTypePrivacyPolicy) {
			hasPrivacy = true
			if !reqMap["required"].(bool) {
				t.Error("Privacy Policy should be required")
			}
		}
	}

	if !hasTerms {
		t.Error("Expected Terms of Service requirement")
	}
	if !hasPrivacy {
		t.Error("Expected Privacy Policy requirement")
	}
}

func TestConsentHandler_UpdateUserConsent(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	handler := NewConsentHandler(logger)

	app := fiber.New()
	app.Post("/api/v1/consent/user/:userId", handler.UpdateUserConsent)

	userID := primitive.NewObjectID().Hex()
	updateRequest := models.ConsentUpdateRequest{
		Type:      models.ConsentTypeTermsOfService,
		Granted:   true,
		Version:   "1.0",
		UserAgent: "Test Agent",
		IPAddress: "127.0.0.1",
		Metadata: map[string]interface{}{
			"source": "test",
		},
	}

	requestBody, _ := json.Marshal(updateRequest)

	// Act
	req := httptest.NewRequest(http.MethodPost, "/api/v1/consent/user/"+userID, bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["message"] != "Consent updated successfully" {
		t.Errorf("Expected success message, got: %v", response["message"])
	}
}

func TestConsentHandler_UpdateUserConsent_InvalidRequest(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	handler := NewConsentHandler(logger)

	app := fiber.New()
	app.Post("/api/v1/consent/user/:userId", handler.UpdateUserConsent)

	userID := primitive.NewObjectID().Hex()
	invalidRequest := map[string]interface{}{
		"type":    "", // Invalid empty type
		"granted": true,
	}

	requestBody, _ := json.Marshal(invalidRequest)

	// Act
	req := httptest.NewRequest(http.MethodPost, "/api/v1/consent/user/"+userID, bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got: %d", resp.StatusCode)
	}
}

func TestConsentHandler_GetUserConsent(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	handler := NewConsentHandler(logger)

	app := fiber.New()
	app.Get("/api/v1/consent/user/:userId", handler.GetUserConsent)

	userID := primitive.NewObjectID().Hex()

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/v1/consent/user/"+userID, nil)
	resp, err := app.Test(req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// For a new user, we expect empty consent
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}

	var response models.UserConsent
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.UserID != userID {
		t.Errorf("Expected user ID %s, got: %s", userID, response.UserID)
	}

	if len(response.Consents) != 0 {
		t.Errorf("Expected empty consents for new user, got: %d", len(response.Consents))
	}
}

func TestConsentHandler_BatchUpdateUserConsent(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	handler := NewConsentHandler(logger)

	app := fiber.New()
	app.Post("/api/v1/consent/user/:userId/batch", handler.BatchUpdateUserConsent)

	userID := primitive.NewObjectID().Hex()
	batchRequest := models.ConsentBatchUpdateRequest{
		Updates: []models.ConsentUpdateRequest{
			{
				Type:    models.ConsentTypeTermsOfService,
				Granted: true,
				Version: "1.0",
			},
			{
				Type:    models.ConsentTypePrivacyPolicy,
				Granted: true,
				Version: "1.0",
			},
			{
				Type:    models.ConsentTypeMarketingCommunications,
				Granted: false,
				Version: "1.0",
			},
		},
	}

	requestBody, _ := json.Marshal(batchRequest)

	// Act
	req := httptest.NewRequest(http.MethodPost, "/api/v1/consent/user/"+userID+"/batch", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["message"] != "Batch consent updated successfully" {
		t.Errorf("Expected success message, got: %v", response["message"])
	}

	if response["updated_count"] != float64(3) {
		t.Errorf("Expected 3 updates, got: %v", response["updated_count"])
	}
}

func TestConsentHandler_BatchUpdateUserConsent_InvalidRequest(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	handler := NewConsentHandler(logger)

	app := fiber.New()
	app.Post("/api/v1/consent/user/:userId/batch", handler.BatchUpdateUserConsent)

	userID := primitive.NewObjectID().Hex()
	invalidRequest := map[string]interface{}{
		"updates": []map[string]interface{}{
			{
				"type":    "", // Invalid empty type
				"granted": true,
			},
		},
	}

	requestBody, _ := json.Marshal(invalidRequest)

	// Act
	req := httptest.NewRequest(http.MethodPost, "/api/v1/consent/user/"+userID+"/batch", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got: %d", resp.StatusCode)
	}
}
