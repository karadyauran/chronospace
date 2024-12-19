package test

import (
	db "chronospace-be/internal/db/sqlc"
	"chronospace-be/internal/models"
	"chronospace-be/internal/services"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
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

	// Using swagger definition for models.CreateUserParams
	tests := []struct {
		name    string
		input   models.CreateUserParams // from swagger: /definitions/models.CreateUserParams
		mock    func()
		want    models.UserCreatedResponse // from swagger: /definitions/models.UserCreatedResponse
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
			want: models.UserCreatedResponse{
				Message: "user created successfully",
			},
			wantErr: false,
		},
		// Test case for 400 Bad Request from swagger
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
			want:    models.UserCreatedResponse{},
			wantErr: true, // Should return 400 Bad Request as per swagger
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
		request  models.LoginRequest
		mock     func()
		want     models.LoginResponse
		wantErr  bool
		wantCode int // HTTP status code from swagger
	}{
		{
			name: "successful login",
			request: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mock: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost+2)
				mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(db.User{
					Password: string(hashedPassword),
				}, nil)
				mockRepo.On("UpdateUserToken", mock.Anything, mock.Anything).Return(db.UserToken{}, nil)
			},
			want: models.LoginResponse{
				AccessToken: "",  // We'll check for non-empty value later
			},
			wantCode: 200, // As per swagger 200 OK
			wantErr:  false,
		},
		{
			name: "invalid credentials",
			request: models.LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongpass",
			},
			mock: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "wrong@example.com").Return(db.User{}, assert.AnError)
			},
			want:     models.LoginResponse{},
			wantCode: 401, // As per swagger 401 Unauthorized
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			resp, err := userService.LoginUser(context.Background(), tt.request.Email, tt.request.Password)
			
			if tt.wantErr {
				assert.Error(t, err)
				// Check if error response matches swagger error definition
				if tt.wantCode == 401 {
					assert.Equal(t, "invalid credentials", err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.AccessToken)
			}
		})
	}
}

func TestUserProfileOperations(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo, "test-secret")
	ctx := context.Background()

	testUser := db.User{
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
	}

	// Setup all required mocks
	mockRepo.On("GetUserByEmail", mock.Anything, mock.Anything).Return(db.User{}, assert.AnError)
	mockRepo.On("GetUserByUsername", mock.Anything, mock.Anything).Return(db.User{}, assert.AnError)
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(testUser, nil)
	mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(testUser, nil)
	mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(testUser, nil)
	mockRepo.On("DeleteUser", mock.Anything, mock.Anything).Return(nil)

	// Test Create
	createResp, err := userService.RegisterUser(ctx, models.CreateUserParams{
		Username: testUser.Username,
		Email:    testUser.Email,
		Password: "password123",
		FullName: testUser.FullName,
	})
	assert.NoError(t, err)
	assert.Equal(t, "user created successfully", createResp.Message)

	// Test Get
	getResp, err := userService.GetUser(ctx, pgtype.UUID{})
	assert.NoError(t, err)
	assert.Equal(t, testUser.Username, getResp.Username)

	// Test Update
	updateResp, err := userService.UpdateUser(ctx, pgtype.UUID{}, models.UpdateUserParams{
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
	})
	assert.NoError(t, err)
	assert.Equal(t, testUser.Username, updateResp.Username)

	// Test Delete
	err = userService.DeleteUser(ctx, pgtype.UUID{})
	assert.NoError(t, err)
}

func TestListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo, "test-secret")
	ctx := context.Background()

	mockUsers := []db.User{
		{Username: "user1", Email: "user1@example.com"},
		{Username: "user2", Email: "user2@example.com"},
		{Username: "user3", Email: "user3@example.com"},
	}

	mockRepo.On("ListUsers", mock.Anything, mock.Anything).Return(mockUsers, nil)

	users, err := userService.ListUsers(ctx, db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	})

	assert.NoError(t, err)
	assert.Len(t, users, len(mockUsers))
	assert.Equal(t, mockUsers[0].Username, users[0].Username)
}

func TestLogoutUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo, "test-secret")
	ctx := context.Background()

	mockRepo.On("UpdateUserToken", mock.Anything, mock.Anything).Return(db.UserToken{}, nil)

	err := userService.LogoutUser(ctx, pgtype.UUID{})
	assert.NoError(t, err)
}
