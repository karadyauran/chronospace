package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"

	"github.com/gin-gonic/gin"
)

type bookingRouter struct {
	bookingController *controllers.BookingController
	config           *config.Config
}

func newBookingRouter(bookingController *controllers.BookingController, config *config.Config) *bookingRouter {
	return &bookingRouter{bookingController, config}
}

func (br *bookingRouter) setBookingRoutes(rg *gin.RouterGroup) {
	router := rg.Group("booking")
	router.POST("/create")
	router.GET("/list")
	router.GET("/:id")
	router.PUT("/:id")
	router.DELETE("/:id")
}