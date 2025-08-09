package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/smorting/backend/internal/database"
	"github.com/smorting/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db           *database.MemoryDatabase
	emailService *EmailService
	jwtSecret    string
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(db *database.MemoryDatabase, emailService *EmailService) *AuthService {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	return &AuthService{
		db:           db,
		emailService: emailService,
		jwtSecret:    jwtSecret,
	}
}

func (a *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:           req.Email,
		Password:        string(hashedPassword),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		Role:            req.Role,
		IsEmailVerified: false,
	}

	// Insert user
	err = a.db.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate and send OTP
	otp, err := a.generateOTP()
	if err != nil {
		return nil, err
	}

	err = a.saveOTP(ctx, req.Email, otp, "registration")
	if err != nil {
		return nil, err
	}

	// Send OTP email
	err = a.emailService.SendOTP(req.Email, otp, "registration")
	if err != nil {
		// Log error but don't fail the registration
		fmt.Printf("Failed to send OTP email: %v\n", err)
	}

	return &models.AuthResponse{
		User:        *user,
		RequiresOTP: true,
	}, nil
}

func (a *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	// Find user
	user, err := a.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// If email is not verified, send OTP
	if !user.IsEmailVerified {
		otp, err := a.generateOTP()
		if err != nil {
			return nil, err
		}

		err = a.saveOTP(ctx, user.Email, otp, "login")
		if err != nil {
			return nil, err
		}

		// Send OTP email
		err = a.emailService.SendOTP(user.Email, otp, "login")
		if err != nil {
			fmt.Printf("Failed to send OTP email: %v\n", err)
		}

		return &models.AuthResponse{
			User:        *user,
			RequiresOTP: true,
		}, nil
	}

	// Generate tokens
	accessToken, refreshToken, err := a.generateTokens(*user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RequiresOTP:  false,
	}, nil
}

func (a *AuthService) VerifyOTP(ctx context.Context, req models.VerifyOTPRequest) (*models.AuthResponse, error) {
	// Find OTP record
	otpRecord, err := a.db.GetOTP(ctx, req.Email, req.OTP)
	if err != nil {
		return nil, errors.New("invalid or expired OTP")
	}

	// Mark OTP as used
	err = a.db.MarkOTPAsUsed(ctx, otpRecord.ID)
	if err != nil {
		return nil, err
	}

	// Find and update user
	user, err := a.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// Mark email as verified
	user.IsEmailVerified = true
	user.UpdatedAt = time.Now()
	err = a.db.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, refreshToken, err := a.generateTokens(*user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RequiresOTP:  false,
	}, nil
}

func (a *AuthService) generateOTP() (string, error) {
	// Generate 6-digit OTP
	max := big.NewInt(999999)
	min := big.NewInt(100000)
	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Add(n, min).Int64()), nil
}

func (a *AuthService) saveOTP(ctx context.Context, email, otp, purpose string) error {
	// Create new OTP record
	otpRecord := &models.OTPRecord{
		Email:     email,
		OTP:       otp,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10 minutes expiry
		IsUsed:    false,
	}

	return a.db.CreateOTP(ctx, otpRecord)
}

func (a *AuthService) generateTokens(user models.User) (string, string, error) {
	// Access token (expires in 24 hours)
	accessClaims := Claims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh token (expires in 7 days)
	refreshClaims := Claims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
