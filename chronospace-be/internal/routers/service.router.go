package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"
	"chronospace-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

type serviceRouter struct {
	serviceController *controllers.ServiceController
	config           *config.Config
	jwtMiddleware    *middleware.JWTConfig
}

func newServiceRouter(serviceController *controllers.ServiceController, config *config.Config, jwtMiddleware *middleware.JWTConfig) *serviceRouter {
	return &serviceRouter{serviceController, config, jwtMiddleware}
}

func (sr *serviceRouter) setServiceRoutes(rg *gin.RouterGroup) {
	router := rg.Group("services")

	// Public routes
	router.GET("", sr.serviceController.ListServices)
	router.GET("/:id", sr.serviceController.GetService)

	// Protected routes
	protected := router.Group("")
	protected.Use(sr.jwtMiddleware.ValidateJWT())
	{
		protected.POST("", sr.serviceController.CreateService)
		protected.PUT("/:id", sr.serviceController.UpdateService)
		protected.DELETE("/:id", sr.serviceController.DeleteService)
	}
}