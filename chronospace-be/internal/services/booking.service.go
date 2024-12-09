package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	db "chronospace-be/internal/db/sqlc"
)

type IBookingRepository interface {
	CreateBooking(ctx context.Context, arg db.CreateBookingParams) (db.Booking, error)
	DeleteBooking(ctx context.Context, id pgtype.UUID) error
	GetBooking(ctx context.Context, id pgtype.UUID) (db.Booking, error)
	ListBookings(ctx context.Context) ([]db.Booking, error)
	ListBookingsByUser(ctx context.Context, userID pgtype.UUID) ([]db.Booking, error)
	UpdateBooking(ctx context.Context, arg db.UpdateBookingParams) (db.Booking, error)
}

type BookingService struct {
	bookingRepo IBookingRepository
}

func NewBookingService (bookingRepository IBookingRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepository,
	}
}