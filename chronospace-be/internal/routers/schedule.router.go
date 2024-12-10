package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"

	"github.com/gin-gonic/gin"
)

type scheduleRouter struct {
	scheduleController *controllers.ScheduleController
	config            *config.Config
}

func newScheduleRouter(scheduleController *controllers.ScheduleController, config *config.Config) *scheduleRouter {
	return &scheduleRouter{scheduleController, config}
}

func (sr *scheduleRouter) setScheduleRoutes(rg *gin.RouterGroup) {
	router := rg.Group("schedule")
	router.POST("/create")
	router.GET("/list")
	router.PUT("/update/:id")
	router.DELETE("/delete/:id")
}