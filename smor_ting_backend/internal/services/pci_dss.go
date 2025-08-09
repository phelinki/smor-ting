package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// PCIDSSService provides PCI-DSS compliant payment processing
type PCIDSSService struct {
	encryptionKey []byte
	logger        *zap.Logger
	store         PaymentTokenStore
	tokenTTL      time.Duration
}

// PaymentToken represents a tokenized payment method
type PaymentToken struct {
	TokenID     string    `json:"token_id"`
	TokenType   string    `json:"token_type"` // "card", "bank_account", "mobile_money"
	LastFour    string    `json:"last_four"`
	ExpiryMonth string    `json:"expiry_month,omitempty"`
	ExpiryYear  string    `json:"expiry_year,omitempty"`
	Brand       string    `json:"brand,omitempty"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsed    time.Time `json:"last_used"`
}

// PaymentRequest represents a PCI-DSS compliant payment request
type PaymentRequest struct {
	Amount      float64                `json:"amount" validate:"required,gt=0"`
	Currency    string                 `json:"currency" validate:"required"`
	TokenID     string                 `json:"token_id" validate:"required"`
	Description string                 `json:"description"`
	Reference   string                 `json:"reference"`
	MerchantID  string                 `json:"merchant_id"`
	CustomerID  string                 `json:"customer_id"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentResponse represents a PCI-DSS compliant payment response
type PaymentResponse struct {
	TransactionID string                 `json:"transaction_id"`
	Status        string                 `json:"status"` // "pending", "completed", "failed"
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Reference     string                 `json:"reference"`
	GatewayRef    string                 `json:"gateway_ref"`
	CreatedAt     time.Time              `json:"created_at"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// SensitivePaymentData represents sensitive payment information (encrypted)
type SensitivePaymentData struct {
	CardNumber    string `json:"card_number,omitempty"`
	CVV           string `json:"cvv,omitempty"`
	ExpiryMonth   string `json:"expiry_month,omitempty"`
	ExpiryYear    string `json:"expiry_year,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	RoutingNumber string `json:"routing_number,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty"`
}

// NewPCIDSSService creates a new PCI-DSS compliant payment service
func NewPCIDSSService(encryptionKey []byte, logger *zap.Logger) (*PCIDSSService, error) {
	if len(encryptionKey) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes (256 bits)")
	}

	return &PCIDSSService{
		encryptionKey: encryptionKey,
		logger:        logger,
		store:         NewMemoryPaymentTokenStore(),
		tokenTTL:      30 * 24 * time.Hour, // default 30 days
	}, nil
}

// SetTokenStore allows injection of a persistent store (Mongo, etc.)
func (p *PCIDSSService) SetTokenStore(store PaymentTokenStore) {
	if store != nil {
		p.store = store
	}
}

// SetTokenTTL configures token expiry duration
func (p *PCIDSSService) SetTokenTTL(ttl time.Duration) {
	if ttl > 0 {
		p.tokenTTL = ttl
	}
}

// TokenizePaymentMethod tokenizes sensitive payment data
func (p *PCIDSSService) TokenizePaymentMethod(sensitiveData *SensitivePaymentData, userID string) (*PaymentToken, error) {
	// Generate unique token ID
	tokenID, err := p.generateTokenID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token ID: %w", err)
	}

	// Encrypt sensitive data
	encryptedData, err := p.encryptSensitiveData(sensitiveData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt sensitive data: %w", err)
	}

	// Create payment token
	token := &PaymentToken{
		TokenID:   tokenID,
		TokenType: p.determineTokenType(sensitiveData),
		LastFour:  p.extractLastFour(sensitiveData),
		Brand:     p.determineBrand(sensitiveData),
		IsDefault: false,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	// Store encrypted data securely with expiry and minimal metadata
	err = p.storeEncryptedData(tokenID, encryptedData, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to store encrypted data: %w", err)
	}
	// Try to persist metadata for retrieval without secret
	if p.store != nil {
		// best-effort: update record with metadata fields
		if rec, err := p.store.Get(tokenID); err == nil && rec != nil {
			rec.TokenType = token.TokenType
			rec.LastFour = token.LastFour
			rec.Brand = token.Brand
		}
	}

	p.logger.Info("Payment method tokenized",
		zap.String("token_id", tokenID),
		zap.String("user_id", userID),
		zap.String("token_type", token.TokenType),
	)

	return token, nil
}

// ProcessPayment processes a payment using a tokenized payment method
func (p *PCIDSSService) ProcessPayment(req *PaymentRequest) (*PaymentResponse, error) {
	// Validate payment request
	if err := p.validatePaymentRequest(req); err != nil {
		return nil, fmt.Errorf("invalid payment request: %w", err)
	}

	// Retrieve encrypted payment data
	encryptedData, err := p.retrieveEncryptedData(req.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment data: %w", err)
	}

	// Decrypt sensitive data for processing
	sensitiveData, err := p.decryptSensitiveData(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt payment data: %w", err)
	}

	// Process payment through secure gateway
	gatewayResponse, err := p.processWithGateway(sensitiveData, req)
	if err != nil {
		return nil, fmt.Errorf("payment processing failed: %w", err)
	}

	// Create payment response
	response := &PaymentResponse{
		TransactionID: p.generateTransactionID(),
		Status:        gatewayResponse.Status,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Reference:     req.Reference,
		GatewayRef:    gatewayResponse.GatewayRef,
		CreatedAt:     time.Now(),
		Metadata:      req.Metadata,
	}

	// Log payment attempt (without sensitive data)
	p.logger.Info("Payment processed",
		zap.String("transaction_id", response.TransactionID),
		zap.String("token_id", req.TokenID),
		zap.String("status", response.Status),
		zap.Float64("amount", req.Amount),
		zap.String("currency", req.Currency),
	)

	return response, nil
}

// ValidatePaymentToken validates a payment token
func (p *PCIDSSService) ValidatePaymentToken(tokenID string) (*PaymentToken, error) {
	// Retrieve token information
	token, err := p.retrieveToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment token: %w", err)
	}

	// Check if token is expired or revoked
	if p.isTokenExpired(token) {
		return nil, fmt.Errorf("payment token expired")
	}

	return token, nil
}

// DeletePaymentToken securely deletes a payment token
func (p *PCIDSSService) DeletePaymentToken(tokenID string) error {
	// Securely delete encrypted data
	err := p.deleteEncryptedData(tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete encrypted data: %w", err)
	}

	// Delete token record
	err = p.deleteToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	p.logger.Info("Payment token deleted", zap.String("token_id", tokenID))
	return nil
}

// encryptSensitiveData encrypts sensitive payment data
func (p *PCIDSSService) encryptSensitiveData(data *SensitivePaymentData) (string, error) {
	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal sensitive data: %w", err)
	}

	// Encrypt using AES-256-GCM
	block, err := aes.NewCipher(p.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptSensitiveData decrypts sensitive payment data
func (p *PCIDSSService) decryptSensitiveData(encryptedData string) (*SensitivePaymentData, error) {
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Decrypt using AES-256-GCM
	block, err := aes.NewCipher(p.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	// Unmarshal JSON
	var data SensitivePaymentData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sensitive data: %w", err)
	}

	return &data, nil
}

// generateTokenID generates a unique token ID
func (p *PCIDSSService) generateTokenID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateTransactionID generates a unique transaction ID
func (p *PCIDSSService) generateTransactionID() string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("txn_%x", hash[:8])
}

// determineTokenType determines the type of payment token
func (p *PCIDSSService) determineTokenType(data *SensitivePaymentData) string {
	if data.CardNumber != "" {
		return "card"
	} else if data.AccountNumber != "" {
		return "bank_account"
	} else if data.PhoneNumber != "" {
		return "mobile_money"
	}
	return "unknown"
}

// extractLastFour extracts the last four digits
func (p *PCIDSSService) extractLastFour(data *SensitivePaymentData) string {
	if data.CardNumber != "" && len(data.CardNumber) >= 4 {
		return data.CardNumber[len(data.CardNumber)-4:]
	}
	if data.AccountNumber != "" && len(data.AccountNumber) >= 4 {
		return data.AccountNumber[len(data.AccountNumber)-4:]
	}
	return ""
}

// determineBrand determines the card brand
func (p *PCIDSSService) determineBrand(data *SensitivePaymentData) string {
	if data.CardNumber == "" {
		return ""
	}

	// Simple brand detection (implement more sophisticated logic)
	firstDigit := data.CardNumber[0]
	switch firstDigit {
	case '4':
		return "visa"
	case '5':
		return "mastercard"
	case '3':
		return "amex"
	default:
		return "unknown"
	}
}

// validatePaymentRequest validates a payment request
func (p *PCIDSSService) validatePaymentRequest(req *PaymentRequest) error {
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if req.TokenID == "" {
		return fmt.Errorf("token ID is required")
	}
	return nil
}

// processWithGateway processes payment through secure gateway
func (p *PCIDSSService) processWithGateway(data *SensitivePaymentData, req *PaymentRequest) (*GatewayResponse, error) {
	// TODO: Implement actual payment gateway integration
	// This is a placeholder for demonstration purposes

	return &GatewayResponse{
		Status:     "completed",
		GatewayRef: fmt.Sprintf("gw_%d", time.Now().Unix()),
	}, nil
}

// GatewayResponse represents a payment gateway response
type GatewayResponse struct {
	Status     string `json:"status"`
	GatewayRef string `json:"gateway_ref"`
}

// Database operations via store abstraction
func (p *PCIDSSService) storeEncryptedData(tokenID, encryptedData, userID string) error {
	if p.store == nil {
		return fmt.Errorf("payment token store not configured")
	}
	expiresAt := time.Now().Add(p.tokenTTL)
	return p.store.Save(tokenID, userID, encryptedData, expiresAt)
}

func (p *PCIDSSService) retrieveEncryptedData(tokenID string) (string, error) {
	if p.store == nil {
		return "", fmt.Errorf("payment token store not configured")
	}
	rec, err := p.store.Get(tokenID)
	if err != nil {
		return "", err
	}
	return rec.EncryptedData, nil
}

func (p *PCIDSSService) deleteEncryptedData(tokenID string) error {
	if p.store == nil {
		return fmt.Errorf("payment token store not configured")
	}
	return p.store.Delete(tokenID)
}

func (p *PCIDSSService) retrieveToken(tokenID string) (*PaymentToken, error) {
	if p.store == nil {
		return nil, fmt.Errorf("payment token store not configured")
	}
	rec, err := p.store.Get(tokenID)
	if err != nil {
		return nil, err
	}
	return &PaymentToken{
		TokenID:   tokenID,
		TokenType: rec.TokenType,
		LastFour:  rec.LastFour,
		Brand:     rec.Brand,
		CreatedAt: rec.CreatedAt,
		LastUsed:  rec.LastUsed,
		IsDefault: false,
	}, nil
}

func (p *PCIDSSService) deleteToken(tokenID string) error {
	if p.store == nil {
		return fmt.Errorf("payment token store not configured")
	}
	return p.store.Delete(tokenID)
}

func (p *PCIDSSService) isTokenExpired(token *PaymentToken) bool {
	if p.store == nil {
		return true
	}
	rec, err := p.store.Get(token.TokenID)
	if err != nil {
		return true
	}
	return time.Now().After(rec.ExpiresAt)
}
