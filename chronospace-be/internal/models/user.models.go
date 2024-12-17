package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserParams struct {
	Username              string    `json:"username"`
	FullName              string    `json:"full_name"`
	Email                 string    `json:"email"`
	Password              string    `json:"password"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type UserCreatedResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type UpdateUserParams struct {
	Username string
	FullName string
	Email    string
	Password string
}

type UserResponse struct {
	ID       pgtype.UUID `json:"id"`
	Username string      `json:"username"`
	FullName string      `json:"full_name"`
	Email    string      `json:"email"`
}
