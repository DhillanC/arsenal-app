package web

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	"github.com/DhillanC/arsenal-app/internal/domain/services"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/ocr"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/persistence/sqlite"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/web/handlers"
	"github.com/gin-gonic/gin"
)

// HTMLConfig configura las vistas HTML
func setupTemplates(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"documentURL": documentURL,
	})
	router.LoadHTMLGlob("web/templates/*")
}

// Config holds server configuration
type Config struct {
	Port           string
	AllowedOrigins []string
	// DB se usa para el health-check; si es nil, /health no verifica DB.
	DB *sql.DB
	// EnableTemplates carga templates HTML para frontend
	EnableTemplates bool
	// UploadPath expone los documentos cargados para el frontend HTML.
	UploadPath string
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

	// Setup HTML templates (solo si se solicita explícitamente)
	if config.EnableTemplates {
		setupTemplates(router)
	}

	// Health-check enriquecido: DB, uploads writable, OCR disponible.
	// Si la DB no responde a Ping en 2s, devuelve 503 — load balancer puede sacarnos del pool.
	router.GET("/health", func(c *gin.Context) {
		checks := gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		}
		status := http.StatusOK

		// Check DB
		if config.DB != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()
			if err := config.DB.PingContext(ctx); err != nil {
				checks["status"] = "degraded"
				checks["db"] = "unreachable"
				checks["db_error"] = err.Error()
				status = http.StatusServiceUnavailable
			} else {
				checks["db"] = "ok"
			}
		}

		// Check uploads writable
		if config.UploadPath != "" {
			testFile := filepath.Join(config.UploadPath, ".healthcheck")
			if err := os.WriteFile(testFile, []byte("ok"), 0o644); err != nil {
				checks["status"] = "degraded"
				checks["uploads"] = "not_writable"
				checks["uploads_error"] = err.Error()
				status = http.StatusServiceUnavailable
			} else {
				os.Remove(testFile)
				checks["uploads"] = "ok"
			}
		}

		// Check OCR disponible
		if ocr.IsAvailable() {
			checks["ocr"] = "ok"
		} else {
			checks["ocr"] = "not_available"
		}

		c.JSON(status, checks)
	})

	// API v1
	api := router.Group("/api/v1")
	{
		replicaHandler := handlers.NewReplicaHandler(replicaService)
		replicaHandler.RegisterRoutes(api)

		actividadHandler := handlers.NewActividadHandler(actividadService, documentoService)
		actividadHandler.RegisterRoutes(api)

		// Documentos handler
		documentoHandler := handlers.NewDocumentoHandler(documentoService)
		// Descarga directa con validación de contención (fuera del grupo /replicas)
		router.GET("/api/v1/documentos/:id/file", documentoHandler.Download)
		documentoRoutes := api.Group("/replicas/:id/documentos")
		{
			documentoRoutes.GET("", documentoHandler.ListByReplica)
			documentoRoutes.POST("", documentoHandler.Upload)
		}
		api.GET("/documentos/filter", documentoHandler.ListByReplicaAndType)
		api.GET("/documentos/search", documentoHandler.Search)
	}

	// Stats endpoint (agregados SQL)
	statsHandler := handlers.NewStatsHandler(config.DB)
	api.GET("/stats/dashboard", statsHandler.DashboardStats)
	api.GET("/export/json", statsHandler.ExportJSON)
	
	// Mantenimiento service and handler
	mantenimientoRepo := sqlite.NewMantenimientoRepository(config.DB)
	mantenimientoService := services.NewMantenimientoService(mantenimientoRepo)
	mantenimientoHandler := handlers.NewMantenimientoHandler(mantenimientoService)
	mantenimientoHandler.RegisterRoutes(api)

	// HTML Frontend routes
	if config.EnableTemplates {
		htmlHandler := handlers.NewHTMLHandler(replicaService, actividadService, documentoService, mantenimientoService, config.UploadPath)
		htmlHandler.RegisterHTMLRoutes(router)
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

func documentURL(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "#"
	}
	if strings.HasPrefix(path, "/uploads/") {
		return path
	}
	if idx := strings.Index(path, "/uploads/"); idx >= 0 {
		return path[idx:]
	}
	path = strings.TrimPrefix(path, "./")
	path = strings.TrimPrefix(path, "uploads/")
	return "/uploads/" + strings.TrimLeft(path, "/")
}
