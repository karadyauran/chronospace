package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"
	"chronospace-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

type scheduleRouter struct {
	scheduleController *controllers.ScheduleController
	config            *config.Config
	jwtMiddleware     *middleware.JWTConfig
}

func newScheduleRouter(scheduleController *controllers.ScheduleController, config *config.Config, jwtMiddleware *middleware.JWTConfig) *scheduleRouter {
	return &scheduleRouter{scheduleController, config, jwtMiddleware}
}

func (sr *scheduleRouter) setScheduleRoutes(rg *gin.RouterGroup) {
	router := rg.Group("schedules")

	// Public routes
	router.GET("", sr.scheduleController.ListSchedules)
	router.GET("/:id", sr.scheduleController.GetSchedule)

	// Protected routes
	protected := router.Group("")
	protected.Use(sr.jwtMiddleware.ValidateJWT())
	{
		protected.POST("", sr.scheduleController.CreateSchedule)
		protected.PUT("/:id", sr.scheduleController.UpdateSchedule)
		protected.DELETE("/:id", sr.scheduleController.DeleteSchedule)
	}
}