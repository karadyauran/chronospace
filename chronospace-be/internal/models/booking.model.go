package models

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

type CreateBookingParams struct {
	UserID    pgtype.UUID `json:"user_id"`
	ServiceID pgtype.UUID `json:"service_id"`
	Date      pgtype.Date `json:"date"`
	Time      pgtype.Time `json:"time"`
	Status    string      `json:"status"`
}

type UpdateBookingParams struct {
	ID        pgtype.UUID `json:"id"`
	Date      pgtype.Date `json:"date"`
	Time      pgtype.Time `json:"time"`
	Status    string      `json:"status"`
}
