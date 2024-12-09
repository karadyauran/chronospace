package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	db "chronospace-be/internal/db/sqlc"
)

type IScheduleRepository interface {
	CreateSchedule(ctx context.Context, arg db.CreateScheduleParams) (db.Schedule, error)
	DeleteSchedule(ctx context.Context, id pgtype.UUID) error
	GetScheduleByID(ctx context.Context, id pgtype.UUID) (db.Schedule, error)
	ListSchedules(ctx context.Context) ([]db.Schedule, error)
	UpdateSchedule(ctx context.Context, arg db.UpdateScheduleParams) (db.Schedule, error)
}

type ScheduleService struct {
	scheduleRepo IScheduleRepository
}

func NewScheduleService(scheduleRepository IScheduleRepository) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepository,
	}
}
