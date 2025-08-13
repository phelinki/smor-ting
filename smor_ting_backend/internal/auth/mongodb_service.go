package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"github.com/smorting/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// MongoDBClaims represents JWT claims for MongoDB
type MongoDBClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// MongoDBService represents the authentication service for MongoDB
type MongoDBService struct {
	repository database.Repository
	config     *configs.AuthConfig
	logger     *logger.Logger
}

// NewMongoDBService creates a new MongoDB authentication service
func NewMongoDBService(repository database.Repository, config *configs.AuthConfig, logger *logger.Logger) (*MongoDBService, error) {
	if repository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	if config == nil {
		return nil, fmt.Errorf("auth configuration is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &MongoDBService{
		repository: repository,
		config:     config,
		logger:     logger,
	}, nil
}

// Register registers a new user
func (s *MongoDBService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	passwordHash, err := s.hashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:           req.Email,
		Password:        passwordHash,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		Role:            req.Role,
		IsEmailVerified: false,
		ProfileImage:    "",
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		s.logger.Error("Failed to create user", err, zap.String("email", req.Email))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token", err, zap.String("user_id", user.ID.Hex()))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("User registered successfully", zap.String("email", req.Email), zap.String("user_id", user.ID.Hex()))

	return &models.AuthResponse{
		User:         *user,
		AccessToken:  token,
		RefreshToken: "", // TODO: Implement refresh token
		RequiresOTP:  false,
	}, nil
}

// Login authenticates a user
func (s *MongoDBService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("Login attempt with non-existent email", zap.String("email", req.Email))
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := s.verifyPassword(req.Password, user.Password); err != nil {
		s.logger.Warn("Login attempt with invalid password", zap.String("email", req.Email))
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token", err, zap.String("user_id", user.ID.Hex()))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("User logged in successfully", zap.String("email", req.Email), zap.String("user_id", user.ID.Hex()))

	return &models.AuthResponse{
		User:         *user,
		AccessToken:  token,
		RefreshToken: "", // TODO: Implement refresh token
		RequiresOTP:  false,
	}, nil
}

// Authenticate verifies email and password and returns the user when valid
func (s *MongoDBService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	// Get user by email
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("Authentication attempt with non-existent email", zap.String("email", email))
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := s.verifyPassword(password, user.Password); err != nil {
		s.logger.Warn("Authentication attempt with invalid password", zap.String("email", email))
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *MongoDBService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MongoDBClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*MongoDBClaims); ok && token.Valid {
		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID in token: %w", err)
		}

		user, err := s.repository.GetUserByID(context.Background(), userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// CreateOTP creates an OTP for the user
func (s *MongoDBService) CreateOTP(ctx context.Context, email, purpose string) error {
	// Generate OTP
	otp := generateOTP()

	// Create OTP record
	otpRecord := &models.OTPRecord{
		Email:     email,
		OTP:       otp,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10 minutes expiry
	}

	if err := s.repository.CreateOTP(ctx, otpRecord); err != nil {
		s.logger.Error("Failed to create OTP", err, zap.String("email", email))
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	s.logger.Info("OTP created successfully", zap.String("email", email), zap.String("purpose", purpose))
	return nil
}

// VerifyOTP verifies an OTP
func (s *MongoDBService) VerifyOTP(ctx context.Context, email, otpCode string) error {
	otp, err := s.repository.GetOTP(ctx, email, otpCode)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Mark OTP as used
	if err := s.repository.MarkOTPAsUsed(ctx, otp.ID); err != nil {
		s.logger.Error("Failed to mark OTP as used", err)
		return fmt.Errorf("failed to verify OTP: %w", err)
	}

	s.logger.Info("OTP verified successfully", zap.String("email", email))
	return nil
}

// GetLatestOTPByEmail returns the most recent unused, unexpired OTP for an email (test utility)
func (s *MongoDBService) GetLatestOTPByEmail(ctx context.Context, email string) (*models.OTPRecord, error) {
	return s.repository.GetLatestOTPByEmail(ctx, email)
}

// GetUserByEmail returns a user by email
func (s *MongoDBService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repository.GetUserByEmail(ctx, email)
}

// RequestPasswordReset creates a password reset OTP and (optionally) triggers email by caller
func (s *MongoDBService) RequestPasswordReset(ctx context.Context, email string) error {
	// Ensure user exists (avoid leaking existence via error message details)
	if _, err := s.repository.GetUserByEmail(ctx, email); err != nil {
		// Return OK to avoid user enumeration; log internally
		s.logger.Warn("Password reset requested for non-existent email", zap.String("email", email))
		return nil
	}
	// Generate OTP and persist
	otp := generateOTP()
	otpRecord := &models.OTPRecord{
		Email:     email,
		OTP:       otp,
		Purpose:   "password_reset",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := s.repository.CreateOTP(ctx, otpRecord); err != nil {
		return fmt.Errorf("failed to create reset otp: %w", err)
	}
	s.logger.Info("Password reset OTP created", zap.String("email", email))
	return nil
}

// ResetPassword verifies OTP and updates the user's password
func (s *MongoDBService) ResetPassword(ctx context.Context, email, otpCode, newPassword string) error {
	// Validate OTP
	if _, err := s.repository.GetOTP(ctx, email, otpCode); err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.config.BCryptCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	// Update user
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	user.Password = string(hash)
	user.UpdatedAt = time.Now()
	if err := s.repository.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	// Mark OTP as used
	if otp, err := s.repository.GetOTP(ctx, email, otpCode); err == nil {
		_ = s.repository.MarkOTPAsUsed(ctx, otp.ID)
	}
	s.logger.Info("Password reset successful", zap.String("email", email))
	return nil
}

// hashPassword hashes a password using bcrypt
func (s *MongoDBService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.BCryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// verifyPassword verifies a password against a hash
func (s *MongoDBService) verifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// generateToken generates a JWT token for a user
func (s *MongoDBService) generateToken(user *models.User) (string, error) {
	now := time.Now()
	claims := &MongoDBClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "smor-ting-backend",
			Subject:   user.ID.Hex(),
			ID:        generateMongoDBTokenID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// generateMongoDBTokenID generates a random token ID for MongoDB service
func generateMongoDBTokenID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// generateOTP generates a 6-digit OTP
func generateOTP() string {
	b := make([]byte, 3)
	rand.Read(b)
	// Convert to 6-digit number
	num := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", num%1000000)
}
