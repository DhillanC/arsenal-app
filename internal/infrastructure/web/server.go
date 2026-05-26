package web

import (
	"context"
	"database/sql"
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
	// DB se usa para el health-check; si es nil, /health no verifica DB.
	DB *sql.DB
}

// NewHandler creates a Gin engine with all routes configured
func NewHandler(
	config Config,
	replicaService inbound.ReplicaService,
	actividadService inbound.ActividadService,
	documentoService inbound.DocumentoService,
) http.Handler {
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware(config.AllowedOrigins))

	// Health-check con verificación de DB.
	// Si la DB no responde a Ping en 2s, devuelve 503 — load balancer puede sacarnos del pool.
	router.GET("/health", func(c *gin.Context) {
		if config.DB != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()
			if err := config.DB.PingContext(ctx); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "degraded",
					"db":     "unreachable",
					"error":  err.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"db":        "ok",
			"timestamp": time.Now().UTC(),
		})
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

// CORSMiddleware configures CORS for the frontend.
//
// Si allowedOrigins está vacío: permite cualquier Origin (modo dev). Esto debe
// configurarse explícitamente en prod vía CORS_ALLOWED_ORIGINS.
// Si la lista contiene orígenes: solo refleja el Origin del request si está en la lista.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

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
