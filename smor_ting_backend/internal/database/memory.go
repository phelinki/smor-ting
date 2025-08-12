package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemoryDatabase struct {
	users      map[string]*models.User
	otpRecords map[string]*models.OTPRecord
	services   map[string]*models.Service
	bookings   map[string]*models.Booking
	mu         sync.RWMutex
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users:      make(map[string]*models.User),
		otpRecords: make(map[string]*models.OTPRecord),
		services:   make(map[string]*models.Service),
		bookings:   make(map[string]*models.Booking),
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
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastSyncAt = time.Now()
	user.Version = 1
	user.IsOffline = false

	// Initialize wallet
	user.Wallet = models.Wallet{
		Balance:     0,
		Currency:    "LRD",
		LastUpdated: time.Now(),
	}

	m.users[user.ID.Hex()] = user
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

func (m *MemoryDatabase) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[id.Hex()]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (m *MemoryDatabase) UpdateUser(ctx context.Context, user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.ID.Hex()]; !exists {
		return errors.New("user not found")
	}

	user.UpdatedAt = time.Now()
	user.LastSyncAt = time.Now()
	user.Version++
	m.users[user.ID.Hex()] = user
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

	otp.ID = primitive.NewObjectID()
	otp.CreatedAt = time.Now()
	otp.IsUsed = false

	m.otpRecords[otp.ID.Hex()] = otp
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

func (m *MemoryDatabase) MarkOTPAsUsed(ctx context.Context, id primitive.ObjectID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	otp, exists := m.otpRecords[id.Hex()]
	if !exists {
		return errors.New("OTP not found")
	}

	otp.IsUsed = true
	return nil
}

// GetLatestOTPByEmail returns the most recent unused, unexpired OTP for an email
func (m *MemoryDatabase) GetLatestOTPByEmail(ctx context.Context, email string) (*models.OTPRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var latest *models.OTPRecord
	for _, otp := range m.otpRecords {
		if otp.Email == email && !otp.IsUsed && otp.ExpiresAt.After(time.Now()) {
			if latest == nil || otp.CreatedAt.After(latest.CreatedAt) {
				latest = otp
			}
		}
	}
	if latest == nil {
		return nil, errors.New("no otp found")
	}
	return latest, nil
}

// Service operations
func (m *MemoryDatabase) CreateService(ctx context.Context, service *models.Service) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	service.ID = primitive.NewObjectID()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.LastSyncAt = time.Now()
	service.Version = 1

	m.services[service.ID.Hex()] = service
	return nil
}

func (m *MemoryDatabase) GetServices(ctx context.Context, categoryID *primitive.ObjectID, location *models.Address, radius float64) ([]models.Service, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var services []models.Service
	for _, service := range m.services {
		if !service.IsActive {
			continue
		}

		if categoryID != nil && service.CategoryID != *categoryID {
			continue
		}

		services = append(services, *service)
	}

	return services, nil
}

// Booking operations
func (m *MemoryDatabase) CreateBooking(ctx context.Context, booking *models.Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	booking.ID = primitive.NewObjectID()
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	booking.LastSyncAt = time.Now()
	booking.Version = 1

	m.bookings[booking.ID.Hex()] = booking

	// Update user's bookings array
	if user, exists := m.users[booking.CustomerID.Hex()]; exists {
		user.Bookings = append(user.Bookings, *booking)
	}

	return nil
}

func (m *MemoryDatabase) GetUserBookings(ctx context.Context, userID primitive.ObjectID) ([]models.Booking, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var bookings []models.Booking
	for _, booking := range m.bookings {
		if booking.CustomerID == userID {
			bookings = append(bookings, *booking)
		}
	}

	return bookings, nil
}

func (m *MemoryDatabase) UpdateBookingStatus(ctx context.Context, bookingID primitive.ObjectID, status models.BookingStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	booking, exists := m.bookings[bookingID.Hex()]
	if !exists {
		return errors.New("booking not found")
	}

	booking.Status = status
	booking.UpdatedAt = time.Now()
	booking.LastSyncAt = time.Now()
	booking.Version++

	return nil
}

// Wallet operations
func (m *MemoryDatabase) UpdateWallet(ctx context.Context, userID primitive.ObjectID, transaction *models.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, exists := m.users[userID.Hex()]
	if !exists {
		return errors.New("user not found")
	}

	transaction.ID = primitive.NewObjectID()
	transaction.CreatedAt = time.Now()

	user.Wallet.Transactions = append(user.Wallet.Transactions, *transaction)
	user.Wallet.Balance += transaction.Amount
	user.Wallet.LastUpdated = time.Now()

	return nil
}

// Offline-first sync operations
func (m *MemoryDatabase) GetUnsyncedData(ctx context.Context, userID primitive.ObjectID, lastSyncAt time.Time) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[userID.Hex()]
	if !exists {
		return nil, errors.New("user not found")
	}

	var bookings []models.Booking
	for _, booking := range m.bookings {
		if booking.CustomerID == userID && booking.LastSyncAt.After(lastSyncAt) {
			bookings = append(bookings, *booking)
		}
	}

	return map[string]interface{}{
		"bookings": bookings,
		"user":     user,
	}, nil
}

func (m *MemoryDatabase) SyncData(ctx context.Context, userID primitive.ObjectID, data map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, exists := m.users[userID.Hex()]
	if !exists {
		return errors.New("user not found")
	}

	user.LastSyncAt = time.Now()
	user.IsOffline = false

	return nil
}

// Enhanced sync operations with checkpoint and compression
func (m *MemoryDatabase) GetUnsyncedDataWithCheckpoint(ctx context.Context, req *models.SyncRequest) (*models.SyncResponse, error) {
	// This method delegates to the sync service
	// In a real implementation, you would implement this directly in the repository
	// For now, we'll return a placeholder response
	return &models.SyncResponse{
		Data:         make(map[string]interface{}),
		Checkpoint:   "",
		LastSyncAt:   time.Now(),
		HasMore:      false,
		Compressed:   false,
		DataSize:     0,
		RecordsCount: 0,
		SyncDuration: 0,
	}, nil
}

func (m *MemoryDatabase) GetChunkedUnsyncedData(ctx context.Context, req *models.ChunkedSyncRequest) (*models.ChunkedSyncResponse, error) {
	// This method delegates to the sync service
	// In a real implementation, you would implement this directly in the repository
	// For now, we'll return a placeholder response
	return &models.ChunkedSyncResponse{
		Data:         make([]interface{}, 0),
		HasMore:      false,
		NextChunk:    0,
		ResumeToken:  "",
		TotalChunks:  0,
		Checkpoint:   "",
		Compressed:   false,
		DataSize:     0,
		RecordsCount: 0,
	}, nil
}

func (m *MemoryDatabase) GetSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.SyncStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[userID.Hex()]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Count pending changes (simplified for memory database)
	pendingCount := 0
	for _, booking := range m.bookings {
		if booking.CustomerID == userID && booking.LastSyncAt.Before(user.LastSyncAt) {
			pendingCount++
		}
	}

	status := &models.SyncStatus{
		UserID:          userID,
		IsOnline:        !user.IsOffline,
		LastSyncAt:      user.LastSyncAt,
		PendingChanges:  pendingCount,
		SyncInProgress:  false,     // Will be set by client
		ConnectionType:  "unknown", // Will be set by client
		ConnectionSpeed: "unknown", // Will be set by client
		UpdatedAt:       time.Now(),
	}

	return status, nil
}

// Setup indexes (no-op for memory database)
func (m *MemoryDatabase) SetupIndexes(ctx context.Context) error {
	return nil
}
