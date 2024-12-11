package models

import "time"

type CreateUserParams struct {
	Username              string    `json:"username"`
	FullName              string    `json:"full_name"`
	Email                 string    `json:"email"`
	Password              string    `json:"password"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type UserCreatedResponse struct {
	Message     string `json:"message"`
}
