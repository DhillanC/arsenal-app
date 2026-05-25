package web

import (
	"net/http"
	"os"
	"time"

	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/web/handlers"
	"github.com/gin-gonic/gin"
)

// Config holds server configuration
type Config struct {
	Port           string
	AllowedOrigins []string
}

// NewHandler creates a Gin engine with all routes configured
func NewHandler(
	config Config,
	replicaService inbound.ReplicaService,
	actividadService inbound.ActividadService,
	documentoService inbound.DocumentoService,
) http.Handler {
	// Set Gin mode based on env
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware(config.AllowedOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// API v1
	api := router.Group("/api/v1")
	{
		replicaHandler := handlers.NewReplicaHandler(replicaService)
		replicaHandler.RegisterRoutes(api)

		actividadHandler := handlers.NewActividadHandler(actividadService)
		actividadHandler.RegisterRoutes(api)
	}

	return router
}

// CORSMiddleware configures CORS for the frontend
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow wildcard in dev (empty list), otherwise check against allowed origins
		if len(allowedOrigins) == 0 || contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Vary", "Origin")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
