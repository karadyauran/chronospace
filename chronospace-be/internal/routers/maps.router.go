package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"
	"chronospace-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

type mapsRouter struct {
	mapsController *controllers.MapsController
	config         *config.Config
	jwtMiddleware  *middleware.JWTConfig
}

func newMapsRouter(mapsController *controllers.MapsController, config *config.Config, jwtMiddleware *middleware.JWTConfig) *mapsRouter {
	return &mapsRouter{mapsController, config, jwtMiddleware}
}

func (mr *mapsRouter) setMapsRoutes(rg *gin.RouterGroup) {
	router := rg.Group("maps")
	router.GET("/search", mr.mapsController.SearchPlaces)
}
