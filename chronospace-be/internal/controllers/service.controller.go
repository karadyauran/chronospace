package controllers

import "chronospace-be/internal/services"

type ServiceController struct {
	ServiceService services.ServiceService
}

func NewServiceController(serviceService services.ServiceService) *ServiceController {
	return &ServiceController{
		ServiceService: serviceService,
	}
}