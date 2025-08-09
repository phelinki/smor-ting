package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceCategory struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Icon        string             `json:"icon" bson:"icon"`
	Color       string             `json:"color" bson:"color"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Service struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CategoryID  primitive.ObjectID `json:"category_id" bson:"category_id"`
	ProviderID  primitive.ObjectID `json:"provider_id" bson:"provider_id"`
	Price       float64            `json:"price" bson:"price"`
	Currency    string             `json:"currency" bson:"currency"`
	Duration    int                `json:"duration" bson:"duration"` // in minutes
	Images      []string           `json:"images" bson:"images"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	Rating      float64            `json:"rating" bson:"rating"`
	ReviewCount int                `json:"review_count" bson:"review_count"`
	// Embedded reviews for better performance
	Reviews []Review `json:"reviews,omitempty" bson:"reviews,omitempty"`
	// Location for geospatial queries
	Location Address `json:"location" bson:"location"`
	// Offline-first fields
	LastSyncAt time.Time `json:"last_sync_at" bson:"last_sync_at"`
	Version    int       `json:"version" bson:"version"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

type ServiceProvider struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	BusinessName   string             `json:"business_name" bson:"business_name"`
	Description    string             `json:"description" bson:"description"`
	Experience     int                `json:"experience" bson:"experience"` // years
	Certifications []string           `json:"certifications" bson:"certifications"`
	ServiceAreas   []string           `json:"service_areas" bson:"service_areas"`
	IsVerified     bool               `json:"is_verified" bson:"is_verified"`
	Rating         float64            `json:"rating" bson:"rating"`
	ReviewCount    int                `json:"review_count" bson:"review_count"`
	CompletedJobs  int                `json:"completed_jobs" bson:"completed_jobs"`
	// Embedded services for better performance
	Services []Service `json:"services,omitempty" bson:"services,omitempty"`
	// Offline-first fields
	LastSyncAt time.Time `json:"last_sync_at" bson:"last_sync_at"`
	Version    int       `json:"version" bson:"version"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

type BookingStatus string

const (
	BookingPending    BookingStatus = "pending"
	BookingConfirmed  BookingStatus = "confirmed"
	BookingInProgress BookingStatus = "in_progress"
	BookingCompleted  BookingStatus = "completed"
	BookingCancelled  BookingStatus = "cancelled"
)

type Booking struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID    primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	ProviderID    primitive.ObjectID `json:"provider_id" bson:"provider_id"`
	ServiceID     primitive.ObjectID `json:"service_id" bson:"service_id"`
	Status        BookingStatus      `json:"status" bson:"status"`
	ScheduledDate time.Time          `json:"scheduled_date" bson:"scheduled_date"`
	CompletedDate *time.Time         `json:"completed_date,omitempty" bson:"completed_date,omitempty"`
	Address       Address            `json:"address" bson:"address"`
	Notes         string             `json:"notes" bson:"notes"`
	TotalAmount   float64            `json:"total_amount" bson:"total_amount"`
	Currency      string             `json:"currency" bson:"currency"`
	PaymentStatus string             `json:"payment_status" bson:"payment_status"`
	// Embedded service details for offline access
	Service Service `json:"service" bson:"service"`
	// Payment information
	Payment Payment `json:"payment" bson:"payment"`
	// Tracking information
	Tracking Tracking `json:"tracking,omitempty" bson:"tracking,omitempty"`
	// Offline-first fields
	LastSyncAt time.Time `json:"last_sync_at" bson:"last_sync_at"`
	Version    int       `json:"version" bson:"version"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

type Payment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Method    string             `json:"method" bson:"method"` // "mobile_money", "card", "wallet"
	Amount    float64            `json:"amount" bson:"amount"`
	Currency  string             `json:"currency" bson:"currency"`
	Status    string             `json:"status" bson:"status"` // "pending", "completed", "failed"
	Reference string             `json:"reference" bson:"reference"`
	Provider  string             `json:"provider" bson:"provider"` // "mpesa", "airtel_money", etc.
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Tracking struct {
	Status           string    `json:"status" bson:"status"`
	Location         Address   `json:"location" bson:"location"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	EstimatedArrival time.Time `json:"estimated_arrival" bson:"estimated_arrival"`
}

type Review struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BookingID  primitive.ObjectID `json:"booking_id" bson:"booking_id"`
	CustomerID primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	ProviderID primitive.ObjectID `json:"provider_id" bson:"provider_id"`
	ServiceID  primitive.ObjectID `json:"service_id" bson:"service_id"`
	Rating     int                `json:"rating" bson:"rating"` // 1-5
	Comment    string             `json:"comment" bson:"comment"`
	// Offline-first fields
	LastSyncAt time.Time `json:"last_sync_at" bson:"last_sync_at"`
	Version    int       `json:"version" bson:"version"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}
