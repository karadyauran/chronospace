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
	UserID      pgtype.UUID `json:"user_id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	FullName    string      `json:"full_name"`
	AccessToken string      `json:"access_token"`
}
