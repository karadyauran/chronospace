package services

import (
	db "chronospace-be/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	UserService         *UserService
	BookingService      *BookingService
	ServiceService      *ServiceService
	ScheduleService     *ScheduleService
	NotificationService *NotificationService
}

func NewService(pool *pgxpool.Pool) *Service {
	queries := db.New(pool)
	return &Service{
		UserService:         NewUserService(queries),
		BookingService:      NewBookingService(queries),
		ServiceService:      NewServiceService(queries),
		ScheduleService:     NewScheduleService(queries),
		NotificationService: NewNotificationService(),
	}
}
