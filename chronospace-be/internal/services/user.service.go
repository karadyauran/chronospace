package services

import (
	db "chronospace-be/internal/db/sqlc"
	"chronospace-be/internal/models"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"go/token"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidContex         = errors.New("invalid context")
	ErrPassword8Symbols      = errors.New("password must be at least 8 characters")
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type IUserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	GetUser(ctx context.Context, id pgtype.UUID) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	DeleteUserRefreshToken(ctx context.Context, id pgtype.UUID) (db.User, error)
	PatchTokenAfterLogin(ctx context.Context, arg db.PatchTokenAfterLoginParams) (db.User, error)
	UpdateUserRefreshToken(ctx context.Context, arg db.UpdateUserRefreshTokenParams) (db.User, error)
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, params models.CreateUserParams) (db.User, error) {
	// Validate context
	if ctx == nil {
		return db.User{}, ErrInvalidContex
	}

	// Trim whitespace from inputs
	params.Email = strings.TrimSpace(params.Email)
	params.Username = strings.TrimSpace(params.Username)

	// Basic validation
	if len(params.Password) < 8 {
		return db.User{}, ErrPassword8Symbols
	}

	if !strings.Contains(params.Email, "@") {
		return db.User{}, ErrInvalidEmailFormat
	}

	// Check if email already exists (using case-insensitive comparison)
	if _, err := s.userRepo.GetUserByEmail(ctx, strings.ToLower(params.Email)); err == nil {
		return db.User{}, ErrEmailAlreadyExists
	}

	// Check if username already exists (using case-insensitive comparison)
	if _, err := s.userRepo.GetUserByUsername(ctx, strings.ToLower(params.Username)); err == nil {
		return db.User{}, ErrUsernameAlreadyExists
	}

	// Use password hashing cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost+2)
	if err != nil {
		return db.User{}, fmt.Errorf("error hashing password: %w", err)
	}

	// Clear the original password from memory
	passwordBytes := []byte(params.Password)
	for i := range passwordBytes {
		passwordBytes[i] = 0
	}
	params.Password = string(passwordBytes)

	// Update params with hashed password and normalized email/username
	params.Password = string(hashedPassword)
	params.Email = strings.ToLower(params.Email)
	params.Username = strings.ToLower(params.Username)

	// Create user with timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.userRepo.CreateUser(ctx, db.CreateUserParams{
		Username: params.Username,
		FullName: params.FullName,
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return db.User{}, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (s *UserService) GenerateTokens(ctx context.Context, userID pgtype.UUID) (models.Tokens, error) {
	// Set expiry times
	accessExpiry := time.Now().Add(15 * time.Minute)
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

	// Create JWT claims for access token
	accessClaims := jwt.MapClaims{
		"user_id": hex.EncodeToString(userID.Bytes[:]),
		"exp":     accessExpiry.Unix(),
		"type":    "access",
	}

	// Create JWT claims for refresh token
	refreshClaims := jwt.MapClaims{
		"user_id": hex.EncodeToString(userID.Bytes[:]),
		"exp":     refreshExpiry.Unix(),
		"type":    "refresh",
	}

	// Sign tokens with secret key
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte("your-secret-key"))
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte("your-secret-key"))
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	// Update refresh token in database
	_, err = s.userRepo.UpdateUserRefreshToken(ctx, db.UpdateUserRefreshTokenParams{
		ID:                    userID,
		RefreshToken:          pgtype.Text{String: refreshToken, Valid: true},
		RefreshTokenExpiresAt: pgtype.Timestamp{Time: refreshExpiry, Valid: true},
	})
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error updating refresh token: %w", err)
	}

	return models.Tokens{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}, nil
}
