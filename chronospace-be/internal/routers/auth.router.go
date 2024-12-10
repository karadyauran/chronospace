package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	authController *controllers.AuthController
	config         *config.Config
}

func newAuthRouter(authController *controllers.AuthController, config *config.Config) *authRouter {
	return &authRouter{authController, config}
}

func (ar *authRouter) setAuthRoutes(rg *gin.RouterGroup) {
	router := rg.Group("auth")
	router.POST("/register")
}