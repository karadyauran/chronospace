package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Schedule struct {
	ID        pgtype.UUID `json:"id"`
	StartTime pgtype.Time `json:"start_time"`
	EndTime   pgtype.Time `json:"end_time"`
	CreatedAt pgtype.Time `json:"created_at"`
	UpdatedAt pgtype.Time `json:"updated_at"`
}

type CreateScheduleRequest struct {
	StartTime pgtype.Time `json:"start_time" binding:"required"`
	EndTime   pgtype.Time `json:"end_time" binding:"required"`
}

type UpdateScheduleRequest struct {
	StartTime pgtype.Time `json:"start_time"`
	EndTime   pgtype.Time `json:"end_time"`
}

type ScheduleResponse struct {
	ID        pgtype.UUID `json:"id"`
	TimeStart pgtype.Time `json:"time_start"`
	TimeEnd   pgtype.Time `json:"time_end"`
	Status    string      `json:"status"`
}
