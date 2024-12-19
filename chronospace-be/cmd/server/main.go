// main.go
//
// Package main is the entry point for the Chronospace backend server.
//
// @title Chronospace API
// @version 1.0
// @description This is the backend server for the Chronospace application.
//
// @contact.name API Support
// @contact.url http://www.chronospace.com/support
// @contact.email support@chronospace.com
//
// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
//
// The main function initializes the configuration, database connection, services, controllers, middleware, and routers.
// It also starts the HTTP server and handles graceful shutdown on interrupt signals.
package main

import (
	"chronospace-be/internal/config"
	"chronospace-be/internal/controllers"
	"chronospace-be/internal/middleware"
	"chronospace-be/internal/routers"
	"chronospace-be/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	database "chronospace-be/internal/db"
)

func main() {
	newConfig, err := config.LoadConfig("./")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load config: %v\n", err)
		os.Exit(1)
	}

	newPool, err := database.NewPostgresDB(&newConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer newPool.Close()

	newService := services.NewService(newPool, newConfig.SecretKey, newConfig.GoogleAPI)
	newController := controllers.NewController(*newService)

	jwtMiddleware := middleware.NewJWTMiddleware(newConfig.SecretKey)

	newRouter := routers.NewRouter(&newConfig, newController, jwtMiddleware)
	newRouter.SetRoutes()

	newServer := &http.Server{
		Addr:    ":" + newConfig.ServerPort,
		Handler: newRouter.Gin,
	}

	// Start the server in a separate goroutine
	go func() {
		log.Printf("Server is running on port %s\n", newConfig.ServerPort)
		if err := newServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server start failed: %s\n", err)
		}
	}()

	// Set up signal catching
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown initiated...")

	// Context for graceful shutdown with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := newServer.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Error: %v", err)
	}

	// Waiting for the shutdown context to be done or timeout
	<-ctx.Done()
	log.Println("Server shutdown completed or timed out")

	log.Println("Server exiting")
	os.Exit(0)
}
