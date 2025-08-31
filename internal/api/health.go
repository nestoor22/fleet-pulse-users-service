package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Checks    map[string]Health `json:"checks"`
}

type Health struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

// HealthCheckHandler godoc
// @Summary Health check endpoint
// @Description Check the health status of the service and its dependencies
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func HealthCheckHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		checks := make(map[string]Health)
		overallStatus := StatusHealthy

		// Check database connectivity
		dbHealth := checkDatabase(db)
		checks["database"] = dbHealth
		if dbHealth.Status != StatusHealthy {
			overallStatus = StatusUnhealthy
		}

		// You can add more health checks here
		// checks["redis"] = checkRedis()
		// checks["external_api"] = checkExternalAPI()

		response := HealthResponse{
			Status:    overallStatus,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Service:   "fleet-pulse-users-service",
			Version:   "1.0.0", // You might want to make this configurable
			Checks:    checks,
		}

		statusCode := http.StatusOK
		if overallStatus == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, response)
	}
}

// ReadinessCheckHandler godoc
// @Summary Readiness check endpoint
// @Description Check if the service is ready to accept requests
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /ready [get]
func ReadinessCheckHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		checks := make(map[string]Health)
		overallStatus := StatusHealthy

		dbHealth := checkDatabase(db)
		checks["database"] = dbHealth
		if dbHealth.Status != StatusHealthy {
			overallStatus = StatusUnhealthy
		}

		response := HealthResponse{
			Status:    overallStatus,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Service:   "fleet-pulse-users-service",
			Version:   "1.0.0",
			Checks:    checks,
		}

		statusCode := http.StatusOK
		if overallStatus == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, response)
	}
}

// LivenessCheckHandler godoc
// @Summary Liveness check endpoint
// @Description Check if the service is alive and running
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /live [get]
func LivenessCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Status:    StatusHealthy,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Service:   "fleet-pulse-users-service",
			Version:   "1.0.0",
			Checks: map[string]Health{
				"service": {
					Status:  StatusHealthy,
					Message: "Service is running",
				},
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func checkDatabase(db *gorm.DB) Health {
	if db == nil {
		return Health{
			Status:  StatusUnhealthy,
			Message: "Database connection is nil",
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return Health{
			Status:  StatusUnhealthy,
			Message: "Failed to get database instance: " + err.Error(),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return Health{
			Status:  StatusUnhealthy,
			Message: "Database ping failed: " + err.Error(),
		}
	}

	return Health{
		Status:  StatusHealthy,
		Message: "Database connection is healthy",
	}
}

func AddHealthRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/health", HealthCheckHandler(db))
	router.GET("/ready", ReadinessCheckHandler(db))
	router.GET("/live", LivenessCheckHandler())
}
