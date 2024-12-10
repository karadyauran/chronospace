package controllers

import "chronospace-be/internal/services"

type Controller struct {
	AuthController     *AuthController
	BookingController  *BookingController
	ScheduleController *ScheduleController
	ServiceController  *ServiceController
}

func NewController(services services.Service) *Controller {
	return &Controller{
		AuthController: NewAuthController(*services.AuthService),
		BookingController: NewBookingController(*services.BookingService),
		ScheduleController: NewScheduleController(*services.ScheduleService),
		ServiceController: NewServiceController(*services.ServiceService),
	}
}
