package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service represents the authentication service
type Service struct {
	db     *sql.DB
	config *configs.AuthConfig
	logger *logger.Logger
}

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	PasswordHash string    `json:"-"` // Never expose password hash
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// Claims represents JWT claims
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewService creates a new authentication service
func NewService(db *sql.DB, config *configs.AuthConfig, logger *logger.Logger) (*Service, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is required")
	}
	if config == nil {
		return nil, fmt.Errorf("auth configuration is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &Service{
		db:     db,
		config: config,
		logger: logger,
	}, nil
}

// Register registers a new user
func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	exists, err := s.userExists(req.Email)
	if err != nil {
		s.logger.Error("Failed to check if user exists", err, zap.String("email", req.Email))
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	passwordHash, err := s.hashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.createUser(req, passwordHash)
	if err != nil {
		s.logger.Error("Failed to create user", err, zap.String("email", req.Email))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token", err, zap.Int("user_id", user.ID))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("User registered successfully", zap.String("email", req.Email), zap.Int("user_id", user.ID))

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// Login authenticates a user
func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.getUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("Login attempt with non-existent email", zap.String("email", req.Email))
			return nil, fmt.Errorf("invalid credentials")
		}
		s.logger.Error("Failed to get user by email", err, zap.String("email", req.Email))
		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}

	// Verify password
	if err := s.verifyPassword(req.Password, user.PasswordHash); err != nil {
		s.logger.Warn("Login attempt with invalid password", zap.String("email", req.Email))
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token", err, zap.Int("user_id", user.ID))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("User logged in successfully", zap.String("email", req.Email), zap.Int("user_id", user.ID))

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *Service) ValidateToken(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		user, err := s.getUserByID(claims.UserID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// userExists checks if a user with the given email exists
func (s *Service) userExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
	err := s.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// createUser creates a new user in the database
func (s *Service) createUser(req *RegisterRequest, passwordHash string) (*User, error) {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name)
		VALUES (?, ?, ?, ?)
		RETURNING id, email, first_name, last_name, created_at, updated_at
	`

	user := &User{}
	err := s.db.QueryRow(query, req.Email, passwordHash, req.FirstName, req.LastName).
		Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// getUserByEmail retrieves a user by email
func (s *Service) getUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, created_at, updated_at FROM users WHERE email = ?`

	user := &User{}
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// getUserByID retrieves a user by ID
func (s *Service) getUserByID(id int) (*User, error) {
	query := `SELECT id, email, first_name, last_name, created_at, updated_at FROM users WHERE id = ?`

	user := &User{}
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// hashPassword hashes a password using bcrypt
func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.BCryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// verifyPassword verifies a password against a hash
func (s *Service) verifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// generateToken generates a JWT token for a user
func (s *Service) generateToken(user *User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "smor-ting-backend",
			Subject:   fmt.Sprintf("%d", user.ID),
			ID:        generateTokenID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// generateTokenID generates a random token ID
func generateTokenID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
