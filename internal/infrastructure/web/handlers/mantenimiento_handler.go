package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	"github.com/gin-gonic/gin"
)

// MantenimientoHandler maneja las peticiones HTTP para mantenimiento
type MantenimientoHandler struct {
	service inbound.MantenimientoService
}

// NewMantenimientoHandler crea un nuevo handler
func NewMantenimientoHandler(service inbound.MantenimientoService) *MantenimientoHandler {
	return &MantenimientoHandler{service: service}
}

// RegisterRoutes registra las rutas de mantenimiento
func (h *MantenimientoHandler) RegisterRoutes(router *gin.RouterGroup) {
	mant := router.Group("/replicas/:id/mantenimiento")
	{
		mant.GET("", h.ListByReplica)
		mant.POST("", h.Create)
	}

	// Rutas independientes
	router.GET("/mantenimiento/proximos", h.ListProximos)
	router.GET("/mantenimiento/:mantenimiento_id", h.GetByID)
	router.PUT("/mantenimiento/:mantenimiento_id", h.Update)
	router.DELETE("/mantenimiento/:mantenimiento_id", h.Delete)
	router.POST("/mantenimiento/:mantenimiento_id/completar", h.MarcarCompletado)
}

// ListByReplica lista mantenimientos de una réplica
func (h *MantenimientoHandler) ListByReplica(c *gin.Context) {
	replicaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de réplica inválido"})
		return
	}

	mantenimientos, err := h.service.ListByReplica(c.Request.Context(), replicaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mantenimientos)
}

// Create crea un nuevo mantenimiento
func (h *MantenimientoHandler) Create(c *gin.Context) {
	replicaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de réplica inválido"})
		return
	}

	var req struct {
		TipoTarea      string `json:"tipo_tarea" form:"tipo_tarea" binding:"required"`
		FrecuenciaDias int    `json:"frecuencia_dias" form:"frecuencia_dias"`
		FrecuenciaBB   int    `json:"frecuencia_bb" form:"frecuencia_bb"`
		UltimaFecha    string `json:"ultima_fecha" form:"ultima_fecha"`
		ProximaFecha   string `json:"proxima_fecha" form:"proxima_fecha"`
		Notas          string `json:"notas" form:"notas"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := &models.Mantenimiento{
		ReplicaID:      replicaID,
		TipoTarea:      req.TipoTarea,
		FrecuenciaDias: req.FrecuenciaDias,
		FrecuenciaBB:   req.FrecuenciaBB,
		Notas:          req.Notas,
	}

	if req.UltimaFecha != "" {
		fecha, err := time.Parse("2006-01-02", req.UltimaFecha)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ultima_fecha inválida"})
			return
		}
		m.UltimaFecha = &fecha
	}

	if req.ProximaFecha != "" {
		fecha, err := time.Parse("2006-01-02", req.ProximaFecha)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "proxima_fecha inválida"})
			return
		}
		m.ProximaFecha = &fecha
	}

	if err := h.service.Create(c.Request.Context(), m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, m)
}

// GetByID obtiene un mantenimiento por ID
func (h *MantenimientoHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("mantenimiento_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	m, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

// ListProximos lista mantenimientos próximos
func (h *MantenimientoHandler) ListProximos(c *gin.Context) {
	dias := 30 // default
	if d := c.Query("dias"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			dias = parsed
		}
	}

	mantenimientos, err := h.service.ListProximos(c.Request.Context(), dias)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dias":           dias,
		"count":          len(mantenimientos),
		"mantenimientos": mantenimientos,
	})
}

// Update actualiza un mantenimiento
func (h *MantenimientoHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("mantenimiento_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		TipoTarea      string `json:"tipo_tarea" form:"tipo_tarea"`
		FrecuenciaDias int    `json:"frecuencia_dias" form:"frecuencia_dias"`
		FrecuenciaBB   int    `json:"frecuencia_bb" form:"frecuencia_bb"`
		UltimaFecha    string `json:"ultima_fecha" form:"ultima_fecha"`
		ProximaFecha   string `json:"proxima_fecha" form:"proxima_fecha"`
		Completado     bool   `json:"completado" form:"completado"`
		Notas          string `json:"notas" form:"notas"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener mantenimiento existente para preservar replica_id
	existing, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mantenimiento no encontrado"})
		return
	}

	m := &models.Mantenimiento{
		ID:             id,
		ReplicaID:      existing.ReplicaID, // Preservar replica_id original
		TipoTarea:      req.TipoTarea,
		FrecuenciaDias: req.FrecuenciaDias,
		FrecuenciaBB:   req.FrecuenciaBB,
		Completado:     req.Completado,
		Notas:          req.Notas,
	}

	if req.UltimaFecha != "" {
		fecha, err := time.Parse("2006-01-02", req.UltimaFecha)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ultima_fecha inválida"})
			return
		}
		m.UltimaFecha = &fecha
	}

	if req.ProximaFecha != "" {
		fecha, err := time.Parse("2006-01-02", req.ProximaFecha)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "proxima_fecha inválida"})
			return
		}
		m.ProximaFecha = &fecha
	}

	if err := h.service.Update(c.Request.Context(), m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

// Delete elimina un mantenimiento
func (h *MantenimientoHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("mantenimiento_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "mantenimiento eliminado"})
}

// MarcarCompletado marca un mantenimiento como completado
func (h *MantenimientoHandler) MarcarCompletado(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("mantenimiento_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		FechaCompletado string `json:"fecha_completado" form:"fecha_completado"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fecha *time.Time
	if req.FechaCompletado != "" {
		f, err := time.Parse("2006-01-02", req.FechaCompletado)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fecha_completado inválida"})
			return
		}
		fecha = &f
	}

	if err := h.service.MarcarCompletado(c.Request.Context(), id, fecha); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "mantenimiento completado"})
}
