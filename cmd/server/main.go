package main

import (
	"context"
	"errors"
	"fleet-pulse-users-service/docs"
	"fleet-pulse-users-service/internal/api"
	"fleet-pulse-users-service/internal/config"
	"fleet-pulse-users-service/internal/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Fleet Pulse Users Service API
// @version 1.0
// @description API for users, authentication, and refresh tokens
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8000
// @BasePath /v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cfg := config.Get()

	databaseConnection := db.DatabaseConnection()

	router := gin.Default()
	v1Group := router.Group("/v1")
	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	api.AddHealthRoutes(router, databaseConnection)

	api.AddUserRoutes(v1Group, databaseConnection)
	api.AddAuthRoutes(v1Group, databaseConnection)

	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
