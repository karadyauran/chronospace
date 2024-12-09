// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateBooking(ctx context.Context, arg CreateBookingParams) (Booking, error)
	CreateSchedule(ctx context.Context, arg CreateScheduleParams) (Schedule, error)
	CreateService(ctx context.Context, arg CreateServiceParams) (Service, error)
	DeleteBooking(ctx context.Context, id pgtype.UUID) error
	DeleteSchedule(ctx context.Context, id pgtype.UUID) error
	DeleteService(ctx context.Context, id pgtype.UUID) error
	GetBooking(ctx context.Context, id pgtype.UUID) (Booking, error)
	GetScheduleByID(ctx context.Context, id pgtype.UUID) (Schedule, error)
	GetService(ctx context.Context, id pgtype.UUID) (Service, error)
	ListBookings(ctx context.Context) ([]Booking, error)
	ListBookingsByUser(ctx context.Context, userID pgtype.UUID) ([]Booking, error)
	ListSchedules(ctx context.Context) ([]Schedule, error)
	ListSchedulesByService(ctx context.Context, serviceID pgtype.UUID) ([]Schedule, error)
	ListServices(ctx context.Context) ([]Service, error)
	UpdateBooking(ctx context.Context, arg UpdateBookingParams) (Booking, error)
	UpdateSchedule(ctx context.Context, arg UpdateScheduleParams) (Schedule, error)
	UpdateService(ctx context.Context, arg UpdateServiceParams) (Service, error)
}

var _ Querier = (*Queries)(nil)