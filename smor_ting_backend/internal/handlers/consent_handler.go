package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// ConsentHandler handles consent-related operations
type ConsentHandler struct {
	logger *zap.Logger
	// In-memory storage for demo purposes
	// In production, this would use a database service
	userConsents map[string]*models.UserConsent
}

// NewConsentHandler creates a new consent handler
func NewConsentHandler(logger *zap.Logger) *ConsentHandler {
	return &ConsentHandler{
		logger:       logger,
		userConsents: make(map[string]*models.UserConsent),
	}
}

// GetRequirements returns the consent requirements
func (h *ConsentHandler) GetRequirements(c *fiber.Ctx) error {
	requirements := models.GetDefaultConsentRequirements()

	return c.JSON(fiber.Map{
		"requirements": requirements,
	})
}

// GetUserConsent returns a user's consent status
func (h *ConsentHandler) GetUserConsent(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Get or create user consent record
	userConsent, exists := h.userConsents[userID]
	if !exists {
		userConsent = &models.UserConsent{
			UserID:      userID,
			Consents:    make(map[models.ConsentType]models.ConsentRecord),
			LastUpdated: time.Now(),
		}
		h.userConsents[userID] = userConsent
	}

	return c.JSON(userConsent)
}

// UpdateUserConsent updates a single consent for a user
func (h *ConsentHandler) UpdateUserConsent(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	var req models.ConsentUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Type == "" || req.Version == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Type and version are required",
		})
	}

	// Get client IP if not provided
	if req.IPAddress == "" {
		req.IPAddress = c.IP()
	}

	// Get or create user consent record
	userConsent, exists := h.userConsents[userID]
	if !exists {
		userConsent = &models.UserConsent{
			UserID:      userID,
			Consents:    make(map[models.ConsentType]models.ConsentRecord),
			LastUpdated: time.Now(),
		}
	}

	// Create consent record
	consentRecord := models.ConsentRecord{
		ID:          primitive.NewObjectID(),
		Type:        req.Type,
		Granted:     req.Granted,
		ConsentedAt: time.Now(),
		Version:     req.Version,
		UserAgent:   req.UserAgent,
		IPAddress:   req.IPAddress,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
	}

	// Update consent
	userConsent.Consents[req.Type] = consentRecord
	userConsent.LastUpdated = time.Now()
	h.userConsents[userID] = userConsent

	h.logger.Info("User consent updated",
		zap.String("user_id", userID),
		zap.String("consent_type", string(req.Type)),
		zap.Bool("granted", req.Granted),
		zap.String("version", req.Version),
	)

	return c.JSON(fiber.Map{
		"message":    "Consent updated successfully",
		"consent_id": consentRecord.ID.Hex(),
	})
}

// BatchUpdateUserConsent updates multiple consents for a user
func (h *ConsentHandler) BatchUpdateUserConsent(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	var req models.ConsentBatchUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if len(req.Updates) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one update is required",
		})
	}

	for _, update := range req.Updates {
		if update.Type == "" || update.Version == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Type and version are required for all updates",
			})
		}
	}

	// Get or create user consent record
	userConsent, exists := h.userConsents[userID]
	if !exists {
		userConsent = &models.UserConsent{
			UserID:      userID,
			Consents:    make(map[models.ConsentType]models.ConsentRecord),
			LastUpdated: time.Now(),
		}
	}

	clientIP := c.IP()
	now := time.Now()
	updatedCount := 0

	// Process each consent update
	for _, update := range req.Updates {
		// Use client IP if not provided
		if update.IPAddress == "" {
			update.IPAddress = clientIP
		}

		// Create consent record
		consentRecord := models.ConsentRecord{
			ID:          primitive.NewObjectID(),
			Type:        update.Type,
			Granted:     update.Granted,
			ConsentedAt: now,
			Version:     update.Version,
			UserAgent:   update.UserAgent,
			IPAddress:   update.IPAddress,
			Metadata:    update.Metadata,
			CreatedAt:   now,
		}

		// Update consent
		userConsent.Consents[update.Type] = consentRecord
		updatedCount++

		h.logger.Info("Batch consent updated",
			zap.String("user_id", userID),
			zap.String("consent_type", string(update.Type)),
			zap.Bool("granted", update.Granted),
			zap.String("version", update.Version),
		)
	}

	userConsent.LastUpdated = now
	h.userConsents[userID] = userConsent

	return c.JSON(fiber.Map{
		"message":       "Batch consent updated successfully",
		"updated_count": updatedCount,
	})
}

// CheckRequiredConsents checks if a user has all required consents
func (h *ConsentHandler) CheckRequiredConsents(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	requirements := models.GetDefaultConsentRequirements()
	userConsent, exists := h.userConsents[userID]

	missingConsents := []models.ConsentRequirement{}
	hasAllRequired := true

	for _, requirement := range requirements {
		if requirement.Required {
			if !exists {
				hasAllRequired = false
				missingConsents = append(missingConsents, requirement)
				continue
			}

			consent, hasConsent := userConsent.Consents[requirement.Type]
			if !hasConsent || !consent.Granted || consent.Version != requirement.Version {
				hasAllRequired = false
				missingConsents = append(missingConsents, requirement)
			}
		}
	}

	return c.JSON(fiber.Map{
		"has_all_required": hasAllRequired,
		"missing_consents": missingConsents,
	})
}
