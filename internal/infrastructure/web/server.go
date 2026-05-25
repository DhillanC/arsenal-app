package web

import (
	"net/http"
	"time"

	inbound "github.com/digital-consultory-solutions/arsenal-app/internal/domain/ports/inbound"
	"github.com/digital-consultory-solutions/arsenal-app/internal/infrastructure/web/handlers"
	"github.com/gin-gonic/gin"
)

// Server encapsula el servidor HTTP
type Server struct {
	router *gin.Engine
	port   string
}

// NewServer crea un nuevo servidor con las dependencias inyectadas
func NewServer(
	replicaService inbound.ReplicaService,
	actividadService inbound.ActividadService,
	documentoService inbound.DocumentoService,
) *Server {
	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// API v1
	api := router.Group("/api/v1")
	{
		// Réplicas
		replicaHandler := handlers.NewReplicaHandler(replicaService)
		replicaHandler.RegisterRoutes(api)

		// Actividades
		actividadHandler := handlers.NewActividadHandler(actividadService)
		actividadHandler.RegisterRoutes(api)
	}

	return &Server{
		router: router,
		port:   "8080",
	}
}

// Run inicia el servidor
func (s *Server) Run() error {
	return s.router.Run(":" + s.port)
}

// CORSMiddleware configura CORS para el frontend
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
