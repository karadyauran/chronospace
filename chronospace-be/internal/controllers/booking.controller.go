package controllers

import "chronospace-be/internal/services"

type BookingController struct {
	bookingService services.BookingService
}

func NewBookingController(bookingService services.BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}