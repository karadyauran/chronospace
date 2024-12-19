package routers

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"
	"chronospace-be/internal/middleware"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "chronospace-be/docs"
)

type Router struct {
	Gin    *gin.Engine
	config *config.Config

	authRouter    *userRouter
	bookingRouter *bookingRouter
}

func NewRouter(config *config.Config, controller *controllers.Controller, jwtMiddleware *middleware.JWTConfig) *Router {
	ginRouter := gin.Default()

	ginRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.WebappBaseUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	ginRouter.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	return &Router{
		Gin:           ginRouter,
		config:        config,
		authRouter:    newUserRouter(controller.UserController, config, jwtMiddleware),
		bookingRouter: newBookingRouter(controller.BookingController, config, jwtMiddleware),
	}
}

func (r *Router) SetRoutes() {
	api := r.Gin.Group("/v1/api")
	r.authRouter.setUserRoutes(api)

	if r.config.EnvType != "prod" {
		r.Gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
