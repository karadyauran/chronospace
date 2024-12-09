package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	db "chronospace-be/internal/db/sqlc"
)

type IServiceRepository interface {
	CreateService(ctx context.Context, arg db.CreateServiceParams) (db.Service, error)
	DeleteService(ctx context.Context, id pgtype.UUID) error
	GetService(ctx context.Context, id pgtype.UUID) (db.Service, error)
	ListServices(ctx context.Context) ([]db.Service, error)
	UpdateService(ctx context.Context, arg db.UpdateServiceParams) (db.Service, error)
}

type ServiceService struct {
	serviceRepo IServiceRepository
}

func NewServiceService(serviceRepository IServiceRepository) *ServiceService {
	return &ServiceService{
		serviceRepo: serviceRepository,
	}
}

