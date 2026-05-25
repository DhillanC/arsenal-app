package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	"github.com/gin-gonic/gin"
)

// ActividadHandler maneja las peticiones HTTP para actividades
type ActividadHandler struct {
	service inbound.ActividadService
}

// NewActividadHandler crea un nuevo handler
func NewActividadHandler(service inbound.ActividadService) *ActividadHandler {
	return &ActividadHandler{service: service}
}

// RegisterRoutes registra las rutas de actividades
func (h *ActividadHandler) RegisterRoutes(router *gin.RouterGroup) {
	actividades := router.Group("/replicas/:id/actividades")
	{
		actividades.GET("", h.ListByReplica)
		actividades.POST("", h.Create)
	}

	// Rutas independientes para operaciones por ID de actividad
	router.GET("/actividades/:actividad_id", h.GetByID)
	router.PUT("/actividades/:actividad_id", h.Update)
	router.DELETE("/actividades/:actividad_id", h.Delete)
}

// ListByReplica lista actividades de una réplica
func (h *ActividadHandler) ListByReplica(c *gin.Context) {
	ctx := c.Request.Context()
	replicaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de réplica inválido"})
		return
	}

	actividades, err := h.service.ListByReplica(ctx, replicaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actividades)
}

// Create crea una nueva actividad
func (h *ActividadHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	replicaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de réplica inválido"})
		return
	}

	var req struct {
		Fecha            string  `json:"fecha" binding:"required"`
		Tipo             string  `json:"tipo" binding:"required"`
		Descripcion      string  `json:"descripcion" binding:"required"`
		ProveedorTecnico string  `json:"proveedor_tecnico"`
		Costo            float64 `json:"costo"`
		KilometrajeBB    int     `json:"kilometraje_bb"`
		Ubicacion        string  `json:"ubicacion"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fechaAdq, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inválida, formato esperado: YYYY-MM-DD"})
		return
	}

	actividad := &models.Actividad{
		ReplicaID:        replicaID,
		Fecha:            fechaAdq,
		Tipo:             req.Tipo,
		Descripcion:      req.Descripcion,
		ProveedorTecnico: req.ProveedorTecnico,
		Costo:            req.Costo,
		KilometrajeBB:    req.KilometrajeBB,
		Ubicacion:        req.Ubicacion,
	}

	if err := h.service.Create(ctx, actividad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, actividad)
}

// GetByID obtiene una actividad por ID
func (h *ActividadHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("actividad_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	actividad, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actividad)
}

// Update actualiza una actividad
func (h *ActividadHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("actividad_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req struct {
		Fecha            string  `json:"fecha"`
		Tipo             string  `json:"tipo"`
		Descripcion      string  `json:"descripcion"`
		ProveedorTecnico string  `json:"proveedor_tecnico"`
		Costo            float64 `json:"costo"`
		KilometrajeBB    int     `json:"kilometraje_bb"`
		Ubicacion        string  `json:"ubicacion"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inválida, formato esperado: YYYY-MM-DD"})
		return
	}

	actividad := &models.Actividad{
		ID:               id,
		Fecha:            fecha,
		Tipo:             req.Tipo,
		Descripcion:      req.Descripcion,
		ProveedorTecnico: req.ProveedorTecnico,
		Costo:            req.Costo,
		KilometrajeBB:    req.KilometrajeBB,
		Ubicacion:        req.Ubicacion,
	}

	if err := h.service.Update(ctx, actividad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actividad)
}

// Delete elimina una actividad
func (h *ActividadHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("actividad_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "actividad eliminada"})
}
