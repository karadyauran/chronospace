package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	ID          pgtype.UUID `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Type        string      `json:"type"` // "hotel" or "apartment"
	Location    string      `json:"location"`
	CreatedAt   pgtype.Time `json:"created_at"`
	UpdatedAt   pgtype.Time `json:"updated_at"`
}

type CreateServiceRequest struct {
	Name        string         `json:"name" binding:"required"`
	Description pgtype.Text    `json:"description" binding:"required"`
	Price       pgtype.Numeric `json:"price" binding:"required"`
	Type        string         `json:"type" binding:"required"`
	Location    string         `json:"location" binding:"required"`
}

type UpdateServiceRequest struct {
	Name        string         `json:"name"`
	Description pgtype.Text    `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Type        string         `json:"type"`
	Location    string         `json:"location"`
}

type ServiceResponse struct {
	ID          pgtype.UUID    `json:"id"`
	Name        string         `json:"name"`
	Description pgtype.Text    `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Type        string         `json:"type"`
	Location    string         `json:"location"`
}
