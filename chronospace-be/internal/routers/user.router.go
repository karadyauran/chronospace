package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"
	"chronospace-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

type userRouter struct {
	userController *controllers.UserController
	config         *config.Config
	jwtMiddleware  *middleware.JWTConfig
}

func newUserRouter(userController *controllers.UserController, config *config.Config, jwtMiddleware *middleware.JWTConfig) *userRouter {
	return &userRouter{userController, config, jwtMiddleware}
}

func (ar *userRouter) setUserRoutes(rg *gin.RouterGroup) {
	router := rg.Group("users")

	// Public routes
	router.POST("/register", ar.userController.Register)
	router.POST("/login", ar.userController.Login)

	// Protected routes
	protected := router.Group("")
	protected.Use(ar.jwtMiddleware.ValidateJWT())
	{
		protected.POST("/logout", ar.userController.Logout)
		protected.GET("/:id", ar.userController.GetUser)
		protected.PUT("/:id", ar.userController.UpdateUser)
		protected.DELETE("/:id", ar.userController.DeleteUser)
		protected.GET("", ar.userController.ListUsers)
	}
}
