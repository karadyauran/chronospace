package controllers

import "chronospace-be/internal/services"

type Controller struct {
	UserController     *UserController
	BookingController  *BookingController
	ScheduleController *ScheduleController
	ServiceController  *ServiceController
	MapsController     *MapsController
}

func NewController(services services.Service) *Controller {
	return &Controller{
		UserController:     NewUserController(*services.UserService),
		BookingController:  NewBookingController(*services.BookingService),
		ScheduleController: NewScheduleController(*services.ScheduleService),
		ServiceController:  NewServiceController(*services.ServiceService),
		MapsController:     NewMapsController(*&services.MapsService),
	}
}
