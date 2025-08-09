package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptionService provides AES-256 encryption for sensitive data
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(key []byte) (*EncryptionService, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes (256 bits)")
	}

	return &EncryptionService{
		key: key,
	}, nil
}

// Encrypt encrypts data using AES-256-GCM
func (e *EncryptionService) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-256-GCM
func (e *EncryptionService) Decrypt(encryptedData string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(e.key)
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

	return plaintext, nil
}

// EncryptString encrypts a string and returns base64 encoded result
func (e *EncryptionService) EncryptString(plaintext string) (string, error) {
	return e.Encrypt([]byte(plaintext))
}

// DecryptString decrypts a base64 encoded string
func (e *EncryptionService) DecryptString(encryptedData string) (string, error) {
	plaintext, err := e.Decrypt(encryptedData)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// GenerateEncryptionKey generates a random 32-byte key for AES-256
func GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}
	return key, nil
}

// EncryptWalletData encrypts sensitive wallet information
func (e *EncryptionService) EncryptWalletData(walletData map[string]interface{}) (map[string]interface{}, error) {
	encryptedData := make(map[string]interface{})

	// Encrypt sensitive fields
	sensitiveFields := []string{"balance", "transactions", "payment_methods"}

	for key, value := range walletData {
		if contains(sensitiveFields, key) {
			// Convert to string for encryption
			valueStr := fmt.Sprintf("%v", value)
			encrypted, err := e.EncryptString(valueStr)
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt %s: %w", key, err)
			}
			encryptedData[key] = encrypted
		} else {
			encryptedData[key] = value
		}
	}

	return encryptedData, nil
}

// DecryptWalletData decrypts sensitive wallet information
func (e *EncryptionService) DecryptWalletData(encryptedData map[string]interface{}) (map[string]interface{}, error) {
	decryptedData := make(map[string]interface{})

	// Decrypt sensitive fields
	sensitiveFields := []string{"balance", "transactions", "payment_methods"}

	for key, value := range encryptedData {
		if contains(sensitiveFields, key) {
			if encryptedStr, ok := value.(string); ok {
				decrypted, err := e.DecryptString(encryptedStr)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt %s: %w", key, err)
				}
				decryptedData[key] = decrypted
			} else {
				decryptedData[key] = value
			}
		} else {
			decryptedData[key] = value
		}
	}

	return decryptedData, nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
