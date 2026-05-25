package web

import (
	"net/http"
	"time"

	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/web/handlers"
	"github.com/gin-gonic/gin"
)

// Config holds server configuration
type Config struct {
	Port            string
	AllowedOrigins  []string
}

// Server encapsula el servidor HTTP
type Server struct {
	router *gin.Engine
	config Config
}

// NewServer crea un nuevo servidor con las dependencias inyectadas
func NewServer(
	config Config,
	replicaService inbound.ReplicaService,
	actividadService inbound.ActividadService,
	documentoService inbound.DocumentoService,
) *Server {
	// Set Gin mode based on env
	if gin.Mode() == gin.DebugMode {
		// already set, keep default
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware(config.AllowedOrigins))

	// Health check with DB ping
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

	return &Server{
		router: router,
		config: config,
	}
}

// Run inicia el servidor
func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}

// CORSMiddleware configura CORS para el frontend
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
