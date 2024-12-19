package services

import (
	db "chronospace-be/internal/db/sqlc"
	"chronospace-be/internal/models"
	"context"

	err2 "chronospace-be/internal/models/enums"

	"github.com/jackc/pgx/v5/pgtype"
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

func NewBookingService(bookingRepository IBookingRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepository,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, params models.CreateBookingParams) (models.Booking, error) {
	if !params.UserID.Valid || !params.ServiceID.Valid {
		return models.Booking{}, err2.ErrBookingInvalidInput
	}

	booking, err := s.bookingRepo.CreateBooking(ctx, db.CreateBookingParams{
		UserID:    params.UserID,
		ServiceID: params.ServiceID,
		Date:      params.Date,
		Time:      params.Time,
		Status:    params.Status,
	})
	if err != nil {
		return models.Booking{}, err
	}

	return models.Booking{
		ID:        booking.ID,
		UserID:    booking.UserID,
		ServiceID: booking.ServiceID,
		Date:      booking.Date,
		Time:      booking.Time,
		Status:    booking.Status,
	}, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id pgtype.UUID) (models.Booking, error) {
	if !id.Valid {
		return models.Booking{}, err2.ErrBookingInvalidInput
	}

	booking, err := s.bookingRepo.GetBooking(ctx, id)
	if err != nil {
		return models.Booking{}, err
	}

	return models.Booking{
		ID:        booking.ID,
		UserID:    booking.UserID,
		ServiceID: booking.ServiceID,
		Date:      booking.Date,
		Time:      booking.Time,
		Status:    booking.Status,
	}, nil
}

func (s *BookingService) ListBookings(ctx context.Context) ([]models.Booking, error) {
	bookings, err := s.bookingRepo.ListBookings(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.Booking, len(bookings))
	for i, booking := range bookings {
		result[i] = models.Booking{
			ID:        booking.ID,
			UserID:    booking.UserID,
			ServiceID: booking.ServiceID,
			Date:      booking.Date,
			Time:      booking.Time,
			Status:    booking.Status,
		}
	}
	return result, nil
}

func (s *BookingService) ListBookingsByUser(ctx context.Context, userID pgtype.UUID) ([]models.Booking, error) {
	if !userID.Valid {
		return nil, err2.ErrBookingInvalidInput
	}

	bookings, err := s.bookingRepo.ListBookingsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]models.Booking, len(bookings))
	for i, booking := range bookings {
		result[i] = models.Booking{
			ID:        booking.ID,
			UserID:    booking.UserID,
			ServiceID: booking.ServiceID,
			Date:      booking.Date,
			Time:      booking.Time,
			Status:    booking.Status,
		}
	}
	return result, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, params models.UpdateBookingParams) (models.Booking, error) {
	if !params.ID.Valid {
		return models.Booking{}, err2.ErrBookingInvalidInput
	}

	booking, err := s.bookingRepo.UpdateBooking(ctx, db.UpdateBookingParams{
		ID:     params.ID,
		Date:   params.Date,
		Time:   params.Time,
		Status: params.Status,
	})
	if err != nil {
		return models.Booking{}, err
	}

	return models.Booking{
		ID:        booking.ID,
		UserID:    booking.UserID,
		ServiceID: booking.ServiceID,
		Date:      booking.Date,
		Time:      booking.Time,
		Status:    booking.Status,
	}, nil
}

func (s *BookingService) DeleteBooking(ctx context.Context, id pgtype.UUID) error {
	if !id.Valid {
		return err2.ErrBookingInvalidInput
	}

	return s.bookingRepo.DeleteBooking(ctx, id)
}
