package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	CustomerRole UserRole = "customer"
	ProviderRole UserRole = "provider"
	AdminRole    UserRole = "admin"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"-" bson:"password"`
	FirstName       string             `json:"first_name" bson:"first_name"`
	LastName        string             `json:"last_name" bson:"last_name"`
	Phone           string             `json:"phone" bson:"phone"`
	Role            UserRole           `json:"role" bson:"role"`
	IsEmailVerified bool               `json:"is_email_verified" bson:"is_email_verified"`
	ProfileImage    string             `json:"profile_image" bson:"profile_image"`
	Address         Address            `json:"address" bson:"address"`
	// Embedded documents for better performance
	Bookings []Booking `json:"bookings,omitempty" bson:"bookings,omitempty"`
	Services []Service `json:"services,omitempty" bson:"services,omitempty"`
	Wallet   Wallet    `json:"wallet,omitempty" bson:"wallet,omitempty"`
	// Offline-first fields
	LastSyncAt time.Time `json:"last_sync_at" bson:"last_sync_at"`
	IsOffline  bool      `json:"is_offline" bson:"is_offline"`
	Version    int       `json:"version" bson:"version"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

type Address struct {
	Street    string  `json:"street" bson:"street"`
	City      string  `json:"city" bson:"city"`
	County    string  `json:"county" bson:"county"`
	Country   string  `json:"country" bson:"country"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type Wallet struct {
	Balance      float64       `json:"balance" bson:"balance"`
	Currency     string        `json:"currency" bson:"currency"`
	Transactions []Transaction `json:"transactions,omitempty" bson:"transactions,omitempty"`
	LastUpdated  time.Time     `json:"last_updated" bson:"last_updated"`
}

type Transaction struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type        string             `json:"type" bson:"type"` // "credit", "debit"
	Amount      float64            `json:"amount" bson:"amount"`
	Description string             `json:"description" bson:"description"`
	Reference   string             `json:"reference" bson:"reference"`
	Status      string             `json:"status" bson:"status"` // "pending", "completed", "failed"
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type OTPRecord struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	OTP       string             `json:"otp" bson:"otp"`
	Purpose   string             `json:"purpose" bson:"purpose"` // "registration", "login", "password_reset"
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	IsUsed    bool               `json:"is_used" bson:"is_used"`
}

// Request/Response models
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=6"`
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	Phone     string   `json:"phone" validate:"required"`
	Role      UserRole `json:"role" validate:"required"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	RequiresOTP  bool   `json:"requires_otp"`
}
