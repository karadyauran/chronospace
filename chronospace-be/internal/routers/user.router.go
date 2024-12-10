package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"

	"github.com/gin-gonic/gin"
)

type userRouter struct {
	userController *controllers.UserController
	config         *config.Config
}

func newUserRouter(userController *controllers.UserController, config *config.Config) *userRouter {
	return &userRouter{userController, config}
}

func (ar *userRouter) setUserRoutes(rg *gin.RouterGroup) {
	router := rg.Group("user")
	router.POST("/create")
}
