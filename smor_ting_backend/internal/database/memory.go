package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/smorting/backend/internal/models"
)

type MemoryDatabase struct {
	users      map[string]*models.User
	otpRecords map[string]*models.OTPRecord
	mu         sync.RWMutex
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users:      make(map[string]*models.User),
		otpRecords: make(map[string]*models.OTPRecord),
	}
}

func (m *MemoryDatabase) Close() error {
	return nil
}

// User operations
func (m *MemoryDatabase) CreateUser(ctx context.Context, user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if user already exists
	for _, existingUser := range m.users {
		if existingUser.Email == user.Email {
			return errors.New("user with this email already exists")
		}
	}
	
	// Generate a simple ID
	user.ID = generateID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	m.users[user.ID] = user
	return nil
}

func (m *MemoryDatabase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	
	return nil, errors.New("user not found")
}

func (m *MemoryDatabase) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

func (m *MemoryDatabase) UpdateUser(ctx context.Context, user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

// OTP operations
func (m *MemoryDatabase) CreateOTP(ctx context.Context, otp *models.OTPRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Remove any existing unused OTPs for this email and purpose
	for id, existingOTP := range m.otpRecords {
		if existingOTP.Email == otp.Email && existingOTP.Purpose == otp.Purpose && !existingOTP.IsUsed {
			delete(m.otpRecords, id)
		}
	}
	
	otp.ID = generateID()
	otp.CreatedAt = time.Now()
	
	m.otpRecords[otp.ID] = otp
	return nil
}

func (m *MemoryDatabase) GetOTP(ctx context.Context, email, otpCode string) (*models.OTPRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for _, otp := range m.otpRecords {
		if otp.Email == email && otp.OTP == otpCode && !otp.IsUsed && otp.ExpiresAt.After(time.Now()) {
			return otp, nil
		}
	}
	
	return nil, errors.New("invalid or expired OTP")
}

func (m *MemoryDatabase) MarkOTPAsUsed(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	otp, exists := m.otpRecords[id]
	if !exists {
		return errors.New("OTP not found")
	}
	
	otp.IsUsed = true
	return nil
}

// Helper function to generate simple IDs
func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}