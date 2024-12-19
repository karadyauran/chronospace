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

func NewService(pool *pgxpool.Pool, secretKey string, gmapsKey string) *Service {
	queries := db.New(pool)

	return &Service{
		UserService:         NewUserService(queries, secretKey),
		BookingService:      NewBookingService(queries),
		ServiceService:      NewServiceService(queries, *NewMapsService(gmapsKey)),
		ScheduleService:     NewScheduleService(queries),
		NotificationService: NewNotificationService(),
	}
}
