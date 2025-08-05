package models

import (
	"time"
)

type UserRole string

const (
	CustomerRole UserRole = "customer"
	ProviderRole UserRole = "provider"
	AdminRole    UserRole = "admin"
)

type User struct {
	ID                string    `json:"id"`
	Email             string    `json:"email"`
	Password          string    `json:"-"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Phone             string    `json:"phone"`
	Role              UserRole  `json:"role"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	ProfileImage      string    `json:"profile_image"`
	Address           Address   `json:"address"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Address struct {
	Street    string  `json:"street"`
	City      string  `json:"city"`
	County    string  `json:"county"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type OTPRecord struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	OTP       string    `json:"otp"`
	Purpose   string    `json:"purpose"` // "registration", "login", "password_reset"
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	IsUsed    bool      `json:"is_used"`
}

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
	RequiresOTP  bool   `json:"requires_otp,omitempty"`
}