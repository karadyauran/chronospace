package services

import (
	db "chronospace-be/internal/db/sqlc"
	"chronospace-be/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
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

func (s *ScheduleService) CreateSchedule(ctx context.Context, req models.CreateScheduleRequest) (models.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.CreateSchedule(ctx, db.CreateScheduleParams{
		TimeStart: req.StartTime,
		TimeEnd:   req.EndTime,
	})

	if err != nil {
		return models.ScheduleResponse{}, err
	}

	return models.ScheduleResponse{
		ID:        schedule.ID,
		TimeStart: schedule.TimeStart,
		TimeEnd:   schedule.TimeEnd,
		Status:    schedule.Status,
	}, err
}

func (s *ScheduleService) GetSchedule(ctx context.Context, id pgtype.UUID) (db.Schedule, error) {
	return s.scheduleRepo.GetScheduleByID(ctx, id)
}

func (s *ScheduleService) ListSchedules(ctx context.Context) ([]db.Schedule, error) {
	return s.scheduleRepo.ListSchedules(ctx)
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, id pgtype.UUID, req models.UpdateScheduleRequest) (models.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.UpdateSchedule(ctx, db.UpdateScheduleParams{
		ID:        id,
		TimeStart: req.StartTime,
		TimeEnd:   req.EndTime,
	})

	if err != nil {
		return models.ScheduleResponse{}, err
	}

	return models.ScheduleResponse{
		ID:        schedule.ID,
		TimeStart: schedule.TimeStart,
		TimeEnd:   schedule.TimeEnd,
		Status:    schedule.Status,
	}, err
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, id pgtype.UUID) error {
	return s.scheduleRepo.DeleteSchedule(ctx, id)
}
