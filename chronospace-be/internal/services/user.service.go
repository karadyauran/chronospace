package services

import (
	db "chronospace-be/internal/db/sqlc"
	"chronospace-be/internal/models"
	"chronospace-be/internal/utils"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	err2 "chronospace-be/internal/models/enums"
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
		return models.UserCreatedResponse{}, err2.ErrInvalidContex
	}

	// Trim whitespace from inputs
	params.Email = strings.TrimSpace(params.Email)
	params.Username = strings.TrimSpace(params.Username)

	// Basic validation
	if len(params.Password) < 8 {
		return models.UserCreatedResponse{}, err2.ErrPassword8Symbols
	}

	if !strings.Contains(params.Email, "@") {
		return models.UserCreatedResponse{}, err2.ErrInvalidEmailFormat
	}

	// Check if email already exists (using case-insensitive comparison)
	if _, err := s.userRepo.GetUserByEmail(ctx, strings.ToLower(params.Email)); err == nil {
		return models.UserCreatedResponse{}, err2.ErrEmailAlreadyExists
	}

	// Check if username already exists (using case-insensitive comparison)
	if _, err := s.userRepo.GetUserByUsername(ctx, strings.ToLower(params.Username)); err == nil {
		return models.UserCreatedResponse{}, err2.ErrUsernameAlreadyExists
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

func (s *UserService) LoginUser(ctx context.Context, email, password string) (models.LoginResponse, error) {
	if ctx == nil {
		return models.LoginResponse{}, err2.ErrInvalidContex
	}

	// Trim and normalize email
	email = strings.ToLower(strings.TrimSpace(email))

	// Basic email validation
	if !strings.Contains(email, "@") {
		return models.LoginResponse{}, err2.ErrInvalidEmailFormat
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.LoginResponse{}, err2.ErrInvalidCredentials
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.LoginResponse{}, err2.ErrInvalidCredentials
	}

	// Generate JWT token
	tokens, err := utils.GenerateTokens(ctx, user.ID)
	if err != nil {
		return models.LoginResponse{}, err2.ErrGeneratingToken
	}

	// Store token in database
	_, err = s.userRepo.UpdateUserToken(ctx, db.UpdateUserTokenParams{
		ID:                    user.ID,
		RefreshToken:          tokens.RefreshToken,
		RefreshTokenExpiresAt: pgtype.Timestamp{Time: tokens.RefreshExpiry, Valid: true},
	})
	if err != nil {
		return models.LoginResponse{}, err2.ErrStoringToken
	}

	return models.LoginResponse{
		AccessToken: tokens.AccessToken,
	}, nil
}

func (s *UserService) LogoutUser(ctx context.Context, userID pgtype.UUID) error {
	if ctx == nil {
		return err2.ErrInvalidContex
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Clear the user's refresh token
	_, err := s.userRepo.UpdateUserToken(ctx, db.UpdateUserTokenParams{
		ID:                    userID,
		RefreshToken:          "",
		RefreshTokenExpiresAt: pgtype.Timestamp{Valid: false},
	})
	if err != nil {
		return err2.ErrCleaningToken
	}

	return nil
}

func (s *UserService) GetUser(ctx context.Context, userID pgtype.UUID) (models.UserResponse, error) {
	if ctx == nil {
		return models.UserResponse{}, err2.ErrInvalidContex
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return models.UserResponse{}, err2.ErrUserNotFound
	}

	return models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID pgtype.UUID, params models.UpdateUserParams) (models.UserResponse, error) {
	if ctx == nil {
		return models.UserResponse{}, err2.ErrInvalidContex
	}

	// Trim whitespace from inputs
	params.Email = strings.TrimSpace(params.Email)
	params.Username = strings.TrimSpace(params.Username)

	// Basic validation
	if len(params.Password) > 0 && len(params.Password) < 8 {
		return models.UserResponse{}, err2.ErrPassword8Symbols
	}

	if !strings.Contains(params.Email, "@") {
		return models.UserResponse{}, err2.ErrInvalidEmailFormat
	}

	// Check if email already exists (using case-insensitive comparison)
	if _, err := s.userRepo.GetUserByEmail(ctx, strings.ToLower(params.Email)); err == nil {
		return models.UserResponse{}, err2.ErrEmailAlreadyExists
	}

	// Check if username already exists (using case-insensitive comparison)
	if _, err := s.userRepo.GetUserByUsername(ctx, strings.ToLower(params.Username)); err == nil {
		return models.UserResponse{}, err2.ErrUsernameAlreadyExists
	}

	// Use password hashing cost if password is provided
	if len(params.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost+2)
		if err != nil {
			return models.UserResponse{}, fmt.Errorf("error hashing password: %w", err)
		}
		params.Password = string(hashedPassword)
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	updatedUser, err := s.userRepo.UpdateUser(ctx, db.UpdateUserParams{
		ID:       userID,
		Username: params.Username,
		FullName: params.FullName,
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("error updating user: %w", err)
	}

	return models.UserResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		FullName: updatedUser.FullName,
		Email:    updatedUser.Email,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID pgtype.UUID) error {
	if ctx == nil {
		return err2.ErrInvalidContex
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return err2.ErrDeletingUser
	}

	return nil
}

func (s *UserService) ListUsers(ctx context.Context, params db.ListUsersParams) ([]models.UserResponse, error) {
	if ctx == nil {
		return nil, err2.ErrInvalidContex
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	users, err := s.userRepo.ListUsers(ctx, params)
	if err != nil {
		return nil, err2.ErrListingUsers
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		})
	}

	return userResponses, nil
}
