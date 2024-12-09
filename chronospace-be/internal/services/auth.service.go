package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	db "chronospace-be/internal/db/sqlc"
)

type IAuthRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	GetUser(ctx context.Context, id pgtype.UUID) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) ( db.User, error)
	GetUserByUsername(ctx context.Context, username string) ( db.User, error)
	ListUsers(ctx context.Context, arg  db.ListUsersParams) ([] db.User, error)
	UpdateUser(ctx context.Context, arg  db.UpdateUserParams) ( db.User, error)
}

type AuthService struct {
	authRepo IAuthRepository
}

func NewAuthService(authRepository IAuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepository,
	}
}