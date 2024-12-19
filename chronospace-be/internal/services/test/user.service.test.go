package test

import (
	"chronospace-be/internal/db/sqlc"
	"chronospace-be/internal/models"
	"chronospace-be/internal/services"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks the IUserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// Implement all IUserRepository methods for the mock
func (m *MockUserRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockUserRepository) CreateUserToken(ctx context.Context, arg db.CreateUserTokenParams) (db.UserToken, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.UserToken), args.Error(1)
}

func (m *MockUserRepository) GetUser(ctx context.Context, id pgtype.UUID) (db.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserToken(ctx context.Context, arg db.UpdateUserTokenParams) (db.UserToken, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.UserToken), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.User), args.Error(1)
}

// Test cases
func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo, "test-secret")

	tests := []struct {
		name    string
		input   models.CreateUserParams
		mock    func()
		wantErr bool
	}{
		{
			name: "successful registration",
			input: models.CreateUserParams{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			mock: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(db.User{}, assert.AnError)
				mockRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(db.User{}, assert.AnError)
				mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(db.User{}, nil)
			},
			wantErr: false,
		},
		{
			name: "email already exists",
			input: models.CreateUserParams{
				Username: "testuser",
				Email:    "existing@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			mock: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(db.User{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := userService.RegisterUser(context.Background(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo, "test-secret")

	tests := []struct {
		name     string
		email    string
		password string
		mock     func()
		wantErr  bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			mock: func() {
				hashedPassword := "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LdiYYWgUQqX5EnX." // hash for "password123"
				mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(db.User{
					Password: hashedPassword,
				}, nil)
				mockRepo.On("UpdateUserToken", mock.Anything, mock.Anything).Return(db.UserToken{}, nil)
			},
			wantErr: false,
		},
		{
			name:     "invalid credentials",
			email:    "wrong@example.com",
			password: "wrongpass",
			mock: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "wrong@example.com").Return(db.User{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := userService.LoginUser(context.Background(), tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Add more test functions for other service methods...

