// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Booking struct {
	ID        pgtype.UUID `json:"id"`
	UserID    pgtype.UUID `json:"user_id"`
	ServiceID pgtype.UUID `json:"service_id"`
	Date      pgtype.Date `json:"date"`
	Time      pgtype.Time `json:"time"`
	Status    string      `json:"status"`
}

type Schedule struct {
	ID        pgtype.UUID `json:"id"`
	ServiceID pgtype.UUID `json:"service_id"`
	Date      pgtype.Date `json:"date"`
	TimeStart pgtype.Time `json:"time_start"`
	TimeEnd   pgtype.Time `json:"time_end"`
	Status    string      `json:"status"`
}

type Service struct {
	ID          pgtype.UUID    `json:"id"`
	Name        string         `json:"name"`
	Description pgtype.Text    `json:"description"`
	Location    string         `json:"location"`
	Price       pgtype.Numeric `json:"price"`
}

type User struct {
	ID        pgtype.UUID      `json:"id"`
	Username  string           `json:"username"`
	FullName  string           `json:"full_name"`
	Email     string           `json:"email"`
	Password  string           `json:"password"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type UserToken struct {
	ID                    pgtype.UUID      `json:"id"`
	UserID                pgtype.UUID      `json:"user_id"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt pgtype.Timestamp `json:"refresh_token_expires_at"`
	CreatedAt             pgtype.Timestamp `json:"created_at"`
}
