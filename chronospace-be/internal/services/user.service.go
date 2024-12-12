package services

import (
	db "chronospace-be/internal/db/sqlc"
	"chronospace-be/internal/models"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	GetUser(ctx context.Context, id pgtype.UUID) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	UpdateUserToken(ctx context.Context, arg db.UpdateUserTokenParams) (db.UserToken, error)
}

var (
	ErrInvalidContex         = errors.New("invalid context")
	ErrPassword8Symbols      = errors.New("password must be at least 8 characters")
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, params models.CreateUserParams) (models.UserCreatedResponse, error) {
	// Validate context
	if ctx == nil {
		return models.UserCreatedResponse{}, ErrInvalidContex
	}

	// Trim whitespace from inputs
	params.Email = strings.TrimSpace(params.Email)
	params.Username = strings.TrimSpace(params.Username)

	// Basic validation
	if len(params.Password) < 8 {
		return models.UserCreatedResponse{}, ErrPassword8Symbols
	}

	if !strings.Contains(params.Email, "@") {
		return models.UserCreatedResponse{}, ErrInvalidEmailFormat
	}

	// Check if email already exists (using case-insensitive comparison)
	if _, err := s.userRepo.GetUserByEmail(ctx, strings.ToLower(params.Email)); err == nil {
		return models.UserCreatedResponse{}, ErrEmailAlreadyExists
	}

	// Check if username already exists (using case-insensitive comparison)
	if _, err := s.userRepo.GetUserByUsername(ctx, strings.ToLower(params.Username)); err == nil {
		return models.UserCreatedResponse{}, ErrUsernameAlreadyExists
	}

	// Use password hashing cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost+2)
	if err != nil {
		return models.UserCreatedResponse{}, fmt.Errorf("error hashing password: %w", err)
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

	_, err = s.userRepo.CreateUser(ctx, db.CreateUserParams{
		Username: params.Username,
		FullName: params.FullName,
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return models.UserCreatedResponse{}, fmt.Errorf("error creating user: %w", err)
	}

	return models.UserCreatedResponse{
		Message: "user created successfully",
	}, nil
}
