package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"
	"github.com/gin-gonic/gin"
)

type serviceRouter struct {
	serviceController *controllers.ServiceController
	config           *config.Config
}

func newServiceRouter(serviceController *controllers.ServiceController, config *config.Config) *serviceRouter {
	return &serviceRouter{serviceController, config}
}

func (sr *serviceRouter) setServiceRoutes(rg *gin.RouterGroup) {
	router := rg.Group("service")
	router.POST("/create")
}