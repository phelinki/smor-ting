package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConsentType represents different types of consent
type ConsentType string

const (
	ConsentTypeTermsOfService          ConsentType = "terms_of_service"
	ConsentTypePrivacyPolicy           ConsentType = "privacy_policy"
	ConsentTypeMarketingCommunications ConsentType = "marketing_communications"
	ConsentTypeDataProcessing          ConsentType = "data_processing"
	ConsentTypeBiometricData           ConsentType = "biometric_data"
	ConsentTypeLocationTracking        ConsentType = "location_tracking"
	ConsentTypeAnalytics               ConsentType = "analytics"
)

// ConsentRecord represents a single consent record for audit trail
type ConsentRecord struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Type        ConsentType            `json:"type" bson:"type"`
	Granted     bool                   `json:"granted" bson:"granted"`
	ConsentedAt time.Time              `json:"consented_at" bson:"consented_at"`
	Version     string                 `json:"version" bson:"version"`
	UserAgent   string                 `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	IPAddress   string                 `json:"ip_address,omitempty" bson:"ip_address,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
}

// UserConsent represents a user's consent status
type UserConsent struct {
	UserID      string                        `json:"user_id" bson:"user_id"`
	Consents    map[ConsentType]ConsentRecord `json:"consents" bson:"consents"`
	LastUpdated time.Time                     `json:"last_updated" bson:"last_updated"`
}

// ConsentRequirement defines what consent is required
type ConsentRequirement struct {
	Type        ConsentType `json:"type" bson:"type"`
	Title       string      `json:"title" bson:"title"`
	Description string      `json:"description" bson:"description"`
	Version     string      `json:"version" bson:"version"`
	Required    bool        `json:"required" bson:"required"`
	DocumentURL string      `json:"document_url,omitempty" bson:"document_url,omitempty"`
}

// ConsentUpdateRequest represents a request to update consent
type ConsentUpdateRequest struct {
	Type      ConsentType            `json:"type" validate:"required"`
	Granted   bool                   `json:"granted"`
	Version   string                 `json:"version" validate:"required"`
	UserAgent string                 `json:"user_agent,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ConsentBatchUpdateRequest represents a batch consent update
type ConsentBatchUpdateRequest struct {
	Updates []ConsentUpdateRequest `json:"updates" validate:"required,dive"`
}

// GetDefaultConsentRequirements returns the default consent requirements
func GetDefaultConsentRequirements() []ConsentRequirement {
	return []ConsentRequirement{
		{
			Type:        ConsentTypeTermsOfService,
			Title:       "Terms of Service",
			Description: "By using Smor-Ting, you agree to our Terms of Service",
			Version:     "1.0",
			Required:    true,
			DocumentURL: "https://smor-ting.com/terms",
		},
		{
			Type:        ConsentTypePrivacyPolicy,
			Title:       "Privacy Policy",
			Description: "We respect your privacy and handle your data according to our Privacy Policy",
			Version:     "1.0",
			Required:    true,
			DocumentURL: "https://smor-ting.com/privacy",
		},
		{
			Type:        ConsentTypeDataProcessing,
			Title:       "Data Processing",
			Description: "Allow processing of your personal data for service delivery",
			Version:     "1.0",
			Required:    true,
		},
		{
			Type:        ConsentTypeMarketingCommunications,
			Title:       "Marketing Communications",
			Description: "Receive promotional emails and notifications about new services",
			Version:     "1.0",
			Required:    false,
		},
		{
			Type:        ConsentTypeAnalytics,
			Title:       "Analytics",
			Description: "Help us improve the app by allowing anonymous usage analytics",
			Version:     "1.0",
			Required:    false,
		},
	}
}
