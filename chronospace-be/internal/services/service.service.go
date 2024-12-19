package services

import (
	db "chronospace-be/internal/db/sqlc"
	"chronospace-be/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
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
	mapsService MapsService
}

func NewServiceService(serviceRepository IServiceRepository, maps MapsService) *ServiceService {
	return &ServiceService{
		serviceRepo: serviceRepository,
		mapsService: maps,
	}
} 

func (s *ServiceService) CreateService(ctx context.Context, req models.CreateServiceRequest) (*models.ServiceResponse, error) {
	// First validate if the location exists
	isValid, err := s.mapsService.ValidateLocation(ctx, req.Location)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, fmt.Errorf("invalid location provided: %s", req.Location)
	}

	arg := db.CreateServiceParams{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Location:    req.Location,
	}

	service, err := s.serviceRepo.CreateService(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &models.ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		Price:       service.Price,
		Location:    service.Location,
	}, nil
}

func (s *ServiceService) GetService(ctx context.Context, id pgtype.UUID) (*models.ServiceResponse, error) {
	service, err := s.serviceRepo.GetService(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		Price:       service.Price,
		Location:    service.Location,
	}, nil
}

func (s *ServiceService) UpdateService(ctx context.Context, id pgtype.UUID, req models.UpdateServiceRequest) (*models.ServiceResponse, error) {
	// If location is being updated, validate it
	if req.Location != "" {
		isValid, err := s.mapsService.ValidateLocation(ctx, req.Location)
		if err != nil {
			return nil, err
		}
		if !isValid {
			return nil, fmt.Errorf("invalid location provided: %s", req.Location)
		}
	}

	// Get existing service to merge with updates
	existingService, err := s.serviceRepo.GetService(ctx, id)
	if err != nil {
		return nil, err
	}

	// Prepare update parameters, keeping existing values if not provided in request
	arg := db.UpdateServiceParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Location:    req.Location,
	}

	// If fields are empty, keep existing values
	if req.Name == "" {
		arg.Name = existingService.Name
	}
	if !req.Description.Valid {
		arg.Description = existingService.Description
	}
	if !req.Price.Valid {
		arg.Price = existingService.Price
	}
	if req.Location == "" {
		arg.Location = existingService.Location
	}

	service, err := s.serviceRepo.UpdateService(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &models.ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		Price:       service.Price,
		Location:    service.Location,
	}, nil
}

func (s *ServiceService) DeleteService(ctx context.Context, id pgtype.UUID) error {
	// Check if service exists
	_, err := s.serviceRepo.GetService(ctx, id)
	if err != nil {
		return fmt.Errorf("service not found: %v", err)
	}

	// Delete the service
	err = s.serviceRepo.DeleteService(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete service: %v", err)
	}

	return nil
}

func (s *ServiceService) ListServices(ctx context.Context) ([]*models.ServiceResponse, error) {
	services, err := s.serviceRepo.ListServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}

	var response []*models.ServiceResponse
	for _, service := range services {
		response = append(response, &models.ServiceResponse{
			ID:          service.ID,
			Name:        service.Name,
			Description: service.Description,
			Price:       service.Price,
			Location:    service.Location,
		})
	}

	return response, nil
}
