package models

import (
	db "chronospace-be/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserParams struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ListUsersParams struct {
	Limit  int32 `form:"limit,default=10"`
	Offset int32 `form:"offset,default=0"`
}

func (p ListUsersParams) ToDBParams() db.ListUsersParams {
	return db.ListUsersParams{
		Limit:  p.Limit,
		Offset: p.Offset,
	}
}
