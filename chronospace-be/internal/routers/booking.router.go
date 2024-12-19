package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"
	"chronospace-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

type bookingRouter struct {
	bookingController *controllers.BookingController
	config            *config.Config
	jwtMiddleware     *middleware.JWTConfig
}

func newBookingRouter(bookingController *controllers.BookingController, config *config.Config, jwtMiddleware *middleware.JWTConfig) *bookingRouter {
	return &bookingRouter{bookingController, config, jwtMiddleware}
}

func (br *bookingRouter) setBookingRoutes(rg *gin.RouterGroup) {
	router := rg.Group("bookings")

	// Public routes
	router.GET("", br.bookingController.ListBookings)
	router.GET("/:id", br.bookingController.GetBooking)

	// Protected routes
	protected := router.Group("")
	protected.Use(br.jwtMiddleware.ValidateJWT())
	{
		protected.POST("", br.bookingController.CreateBooking)
		protected.GET("/user", br.bookingController.ListUserBookings)
		protected.PUT("/:id", br.bookingController.UpdateBooking)
		protected.DELETE("/:id", br.bookingController.DeleteBooking)
	}
}
