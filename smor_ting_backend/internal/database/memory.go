package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smorting/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemoryDatabase struct {
	users                map[string]*models.User
	otpRecords           map[string]*models.OTPRecord
	services             map[string]*models.Service
	bookings             map[string]*models.Booking
	deviceSessions       map[string]*models.DeviceSession
	securityEvents       map[string]*models.SecurityEvent
	syncCheckpoints      map[string]*models.SyncCheckpoint
	syncMetrics          map[string]*models.SyncMetrics
	syncStatuses         map[string]*models.SyncStatus
	syncQueueItems       map[string]*models.SyncQueueItem
	backgroundSyncStatus map[string]*models.BackgroundSyncStatus
	mu                   sync.RWMutex
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users:                make(map[string]*models.User),
		otpRecords:           make(map[string]*models.OTPRecord),
		services:             make(map[string]*models.Service),
		bookings:             make(map[string]*models.Booking),
		deviceSessions:       make(map[string]*models.DeviceSession),
		securityEvents:       make(map[string]*models.SecurityEvent),
		syncCheckpoints:      make(map[string]*models.SyncCheckpoint),
		syncMetrics:          make(map[string]*models.SyncMetrics),
		syncStatuses:         make(map[string]*models.SyncStatus),
		syncQueueItems:       make(map[string]*models.SyncQueueItem),
		backgroundSyncStatus: make(map[string]*models.BackgroundSyncStatus),
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

	// Generate a simple ID if not already provided by caller (tests may set a specific ID)
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastSyncAt = time.Now()
	user.Version = 1
	user.IsOffline = false

	// Initialize wallet with default currency if not already set
	if user.Wallet.Currency == "" {
		user.Wallet = models.Wallet{
			Balance:     0,
			Currency:    "USD",
			LastUpdated: time.Now(),
		}
	} else {
		// Preserve existing wallet settings but ensure LastUpdated is set
		if user.Wallet.LastUpdated.IsZero() {
			user.Wallet.LastUpdated = time.Now()
		}
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

	// Set expiry if not already set
	if otp.ExpiresAt.IsZero() {
		otp.ExpiresAt = time.Now().Add(10 * time.Minute)
	}

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

	var services []models.Service
	for _, service := range m.services {
		if service.ProviderID == userID && service.LastSyncAt.After(lastSyncAt) {
			services = append(services, *service)
		}
	}

	return map[string]interface{}{
		"bookings": bookings,
		"services": services,
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
	start := time.Now()

	// Get unsynced data using the basic method
	data, err := m.GetUnsyncedData(ctx, req.UserID, req.LastSyncAt)
	if err != nil {
		return nil, err
	}

	// Count total records
	recordsCount := 0
	for _, value := range data {
		if slice, ok := value.([]models.Booking); ok {
			recordsCount += len(slice)
		} else if slice, ok := value.([]models.Service); ok {
			recordsCount += len(slice)
		} else {
			recordsCount += 1 // user object
		}
	}

	// Generate a checkpoint (simplified)
	checkpoint := fmt.Sprintf("checkpoint_%d", time.Now().Unix())

	// Calculate data size (rough estimate)
	dataSize := int64(recordsCount * 1024) // Rough estimate of 1KB per record

	return &models.SyncResponse{
		Data:         data,
		Checkpoint:   checkpoint,
		LastSyncAt:   time.Now(),
		HasMore:      false, // Simplified - always false for memory database
		Compressed:   req.Compression,
		DataSize:     dataSize,
		RecordsCount: recordsCount,
		SyncDuration: time.Since(start),
	}, nil
}

func (m *MemoryDatabase) GetChunkedUnsyncedData(ctx context.Context, req *models.ChunkedSyncRequest) (*models.ChunkedSyncResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var allData []interface{}

	// Collect all user data
	if user, exists := m.users[req.UserID.Hex()]; exists {
		allData = append(allData, user)
	}

	// Collect bookings
	for _, booking := range m.bookings {
		if booking.CustomerID == req.UserID {
			allData = append(allData, booking)
		}
	}

	// Collect services
	for _, service := range m.services {
		if service.ProviderID == req.UserID {
			allData = append(allData, service)
		}
	}

	// Apply chunking
	startIndex := req.ChunkIndex * req.ChunkSize
	endIndex := startIndex + req.ChunkSize
	if endIndex > len(allData) {
		endIndex = len(allData)
	}

	var chunkData []interface{}
	if startIndex < len(allData) {
		chunkData = allData[startIndex:endIndex]
	}

	totalChunks := (len(allData) + req.ChunkSize - 1) / req.ChunkSize
	hasMore := req.ChunkIndex < totalChunks-1
	nextChunk := req.ChunkIndex
	if hasMore {
		nextChunk++
	}

	checkpoint := fmt.Sprintf("chunk_%d_%d", req.ChunkIndex, time.Now().Unix())

	return &models.ChunkedSyncResponse{
		Data:         chunkData,
		HasMore:      hasMore,
		NextChunk:    nextChunk,
		ResumeToken:  fmt.Sprintf("resume_%d", nextChunk),
		TotalChunks:  totalChunks,
		Checkpoint:   checkpoint,
		Compressed:   false,
		DataSize:     int64(len(chunkData) * 1024), // Rough estimate
		RecordsCount: len(chunkData),
	}, nil
}

func (m *MemoryDatabase) GetSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.SyncStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if sync status already exists
	if status, exists := m.syncStatuses[userID.Hex()]; exists {
		return status, nil
	}

	// Get user to create initial sync status
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

	// Create default sync status
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

	// Store the initial status
	m.syncStatuses[userID.Hex()] = status

	return status, nil
}

// Device session operations
func (m *MemoryDatabase) CreateDeviceSession(ctx context.Context, session *models.DeviceSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session.ID.IsZero() {
		session.ID = primitive.NewObjectID()
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	session.LastActivity = time.Now()

	m.deviceSessions[session.ID.Hex()] = session
	return nil
}

func (m *MemoryDatabase) GetDeviceSession(ctx context.Context, sessionID string) (*models.DeviceSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.deviceSessions[sessionID]
	if !exists {
		return nil, errors.New("device session not found")
	}

	return session, nil
}

func (m *MemoryDatabase) GetDeviceSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.DeviceSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, session := range m.deviceSessions {
		if session.RefreshToken == refreshToken && session.IsActive {
			return session, nil
		}
	}

	return nil, errors.New("device session not found for refresh token")
}

func (m *MemoryDatabase) GetDeviceSessionByDeviceID(ctx context.Context, deviceID string) (*models.DeviceSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, session := range m.deviceSessions {
		if session.DeviceID == deviceID && session.IsActive {
			return session, nil
		}
	}

	return nil, errors.New("device session not found for device ID")
}

func (m *MemoryDatabase) GetUserDeviceSessions(ctx context.Context, userID string) ([]models.DeviceSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var sessions []models.DeviceSession
	for _, session := range m.deviceSessions {
		if session.UserID == userObjectID {
			sessions = append(sessions, *session)
		}
	}

	return sessions, nil
}

func (m *MemoryDatabase) UpdateDeviceSessionActivity(ctx context.Context, sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.deviceSessions[sessionID]
	if !exists {
		return errors.New("device session not found")
	}

	session.UpdateActivity()
	return nil
}

func (m *MemoryDatabase) RevokeDeviceSession(ctx context.Context, sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.deviceSessions[sessionID]
	if !exists {
		return errors.New("device session not found")
	}

	session.RevokeSession()
	return nil
}

func (m *MemoryDatabase) RevokeAllUserTokens(ctx context.Context, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	for _, session := range m.deviceSessions {
		if session.UserID == userObjectID && session.IsActive {
			session.RevokeSession()
		}
	}

	return nil
}

func (m *MemoryDatabase) RotateRefreshToken(ctx context.Context, sessionID string, newRefreshToken string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.deviceSessions[sessionID]
	if !exists {
		return errors.New("device session not found")
	}

	session.RefreshToken = newRefreshToken
	session.UpdateActivity()
	return nil
}

func (m *MemoryDatabase) CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, session := range m.deviceSessions {
		if session.IsExpired(maxAge) {
			session.RevokeSession()
			// In a real implementation, you might want to delete expired sessions
			// For testing, we'll just revoke them
			m.deviceSessions[id] = session
		}
	}

	return nil
}

// Security event operations
func (m *MemoryDatabase) LogSecurityEvent(ctx context.Context, event *models.SecurityEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if event.ID.IsZero() {
		event.ID = primitive.NewObjectID()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	m.securityEvents[event.ID.Hex()] = event
	return nil
}

func (m *MemoryDatabase) GetUserSecurityEvents(ctx context.Context, userID string, limit int) ([]models.SecurityEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var events []models.SecurityEvent
	count := 0
	for _, event := range m.securityEvents {
		if event.UserID == userObjectID {
			events = append(events, *event)
			count++
			if count >= limit {
				break
			}
		}
	}

	return events, nil
}

func (m *MemoryDatabase) GetSecurityEventsByType(ctx context.Context, userID string, eventType models.SecurityEventType, limit int) ([]models.SecurityEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var events []models.SecurityEvent
	count := 0
	for _, event := range m.securityEvents {
		if event.UserID == userObjectID && event.EventType == eventType {
			events = append(events, *event)
			count++
			if count >= limit {
				break
			}
		}
	}

	return events, nil
}

// Sync status operations
func (m *MemoryDatabase) UpdateSyncStatus(ctx context.Context, status *models.SyncStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	status.UpdatedAt = time.Now()
	m.syncStatuses[status.UserID.Hex()] = status
	return nil
}

// Sync checkpoint operations
func (m *MemoryDatabase) CreateSyncCheckpoint(ctx context.Context, checkpoint *models.SyncCheckpoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if checkpoint.ID.IsZero() {
		checkpoint.ID = primitive.NewObjectID()
	}
	if checkpoint.CreatedAt.IsZero() {
		checkpoint.CreatedAt = time.Now()
	}
	checkpoint.UpdatedAt = time.Now()

	m.syncCheckpoints[checkpoint.UserID.Hex()] = checkpoint
	return nil
}

func (m *MemoryDatabase) GetSyncCheckpoint(ctx context.Context, userID primitive.ObjectID) (*models.SyncCheckpoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	checkpoint, exists := m.syncCheckpoints[userID.Hex()]
	if !exists {
		return nil, errors.New("sync checkpoint not found")
	}

	return checkpoint, nil
}

func (m *MemoryDatabase) UpdateSyncCheckpoint(ctx context.Context, checkpoint *models.SyncCheckpoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	checkpoint.UpdatedAt = time.Now()
	m.syncCheckpoints[checkpoint.UserID.Hex()] = checkpoint
	return nil
}

// Sync metrics operations
func (m *MemoryDatabase) CreateSyncMetrics(ctx context.Context, metrics *models.SyncMetrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metrics.ID.IsZero() {
		metrics.ID = primitive.NewObjectID()
	}
	if metrics.CreatedAt.IsZero() {
		metrics.CreatedAt = time.Now()
	}

	m.syncMetrics[metrics.ID.Hex()] = metrics
	return nil
}

func (m *MemoryDatabase) GetRecentSyncMetrics(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var metrics []models.SyncMetrics
	count := 0
	for _, metric := range m.syncMetrics {
		if metric.UserID == userID {
			metrics = append(metrics, *metric)
			count++
			if count >= limit {
				break
			}
		}
	}

	// Sort by creation time (most recent first)
	for i := 0; i < len(metrics)-1; i++ {
		for j := i + 1; j < len(metrics); j++ {
			if metrics[i].CreatedAt.Before(metrics[j].CreatedAt) {
				metrics[i], metrics[j] = metrics[j], metrics[i]
			}
		}
	}

	return metrics, nil
}

// Background sync queue operations
func (m *MemoryDatabase) CreateSyncQueueItem(ctx context.Context, item *models.SyncQueueItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = time.Now()

	m.syncQueueItems[item.ID.Hex()] = item
	return nil
}

func (m *MemoryDatabase) GetSyncQueueItem(ctx context.Context, itemID primitive.ObjectID) (*models.SyncQueueItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.syncQueueItems[itemID.Hex()]
	if !exists {
		return nil, errors.New("sync queue item not found")
	}

	return item, nil
}

func (m *MemoryDatabase) UpdateSyncQueueItem(ctx context.Context, item *models.SyncQueueItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item.UpdatedAt = time.Now()
	m.syncQueueItems[item.ID.Hex()] = item
	return nil
}

func (m *MemoryDatabase) GetPendingSyncQueueItems(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncQueueItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var items []models.SyncQueueItem
	count := 0

	// Convert to slice for sorting
	var allItems []*models.SyncQueueItem
	for _, item := range m.syncQueueItems {
		if item.UserID == userID && item.Status == models.SyncQueuePending {
			allItems = append(allItems, item)
		}
	}

	// Sort by priority (highest first)
	for i := 0; i < len(allItems)-1; i++ {
		for j := i + 1; j < len(allItems); j++ {
			if allItems[i].Priority < allItems[j].Priority {
				allItems[i], allItems[j] = allItems[j], allItems[i]
			}
		}
	}

	// Apply limit and convert to value slice
	for _, item := range allItems {
		if count >= limit {
			break
		}
		items = append(items, *item)
		count++
	}

	return items, nil
}

func (m *MemoryDatabase) GetConflictQueueItems(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.SyncQueueItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var items []models.SyncQueueItem
	count := 0

	for _, item := range m.syncQueueItems {
		if item.UserID == userID && item.Type == models.SyncTypeConflict {
			items = append(items, *item)
			count++
			if count >= limit {
				break
			}
		}
	}

	return items, nil
}

func (m *MemoryDatabase) CleanupCompletedQueueItems(ctx context.Context, olderThan time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoffTime := time.Now().Add(-olderThan)
	var deletedCount int64

	for id, item := range m.syncQueueItems {
		if item.Status == models.SyncQueueCompleted &&
			item.CompletedAt != nil &&
			item.CompletedAt.Before(cutoffTime) {
			delete(m.syncQueueItems, id)
			deletedCount++
		}
	}

	return deletedCount, nil
}

// Background sync status operations
func (m *MemoryDatabase) GetBackgroundSyncStatus(ctx context.Context, userID primitive.ObjectID) (*models.BackgroundSyncStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if status already exists
	if status, exists := m.backgroundSyncStatus[userID.Hex()]; exists {
		return status, nil
	}

	// Create default status
	status := &models.BackgroundSyncStatus{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		IsEnabled:        true,
		LastSyncAt:       time.Now(),
		PendingItems:     0,
		FailedItems:      0,
		ConflictItems:    0,
		AutoRetryEnabled: true,
		NextScheduledRun: time.Now().Add(5 * time.Minute),
		UpdatedAt:        time.Now(),
	}

	m.backgroundSyncStatus[userID.Hex()] = status
	return status, nil
}

func (m *MemoryDatabase) UpdateBackgroundSyncStatus(ctx context.Context, status *models.BackgroundSyncStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	status.UpdatedAt = time.Now()
	m.backgroundSyncStatus[status.UserID.Hex()] = status
	return nil
}

// Setup indexes (no-op for memory database)
func (m *MemoryDatabase) SetupIndexes(ctx context.Context) error {
	return nil
}
