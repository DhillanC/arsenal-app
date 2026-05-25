package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/digital-consultory-solutions/arsenal-app/internal/domain/models"
	inbound "github.com/digital-consultory-solutions/arsenal-app/internal/domain/ports/inbound"
	"github.com/gin-gonic/gin"
)

// ReplicaHandler maneja las peticiones HTTP para réplicas
type ReplicaHandler struct {
	service inbound.ReplicaService
}

// NewReplicaHandler crea un nuevo handler
func NewReplicaHandler(service inbound.ReplicaService) *ReplicaHandler {
	return &ReplicaHandler{service: service}
}

// RegisterRoutes registra las rutas de réplicas
func (h *ReplicaHandler) RegisterRoutes(router *gin.RouterGroup) {
	replicas := router.Group("/replicas")
	{
		replicas.GET("", h.List)
		replicas.POST("", h.Create)
		replicas.GET("/:id", h.GetByID)
		replicas.PUT("/:id", h.Update)
		replicas.DELETE("/:id", h.Delete)
	}
}

// List devuelve todas las réplicas
func (h *ReplicaHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	replicas, err := h.service.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, replicas)
}

// Create crea una nueva réplica
func (h *ReplicaHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	
	var req struct {
		Nombre            string    `json:"nombre" binding:"required"`
		Marca             string    `json:"marca"`
		Modelo            string    `json:"modelo"`
		Tipo              string    `json:"tipo"`
		NumeroSerie       string    `json:"numero_serie"`
		FechaAdquisicion  string    `json:"fecha_adquisicion"`
		Proveedor         string    `json:"proveedor"`
		CostoAdquisicion  float64   `json:"costo_adquisicion"`
		Estado            string    `json:"estado"`
		FPS               int       `json:"fps"`
		Joules            float64   `json:"joules"`
		PesoGramos        int       `json:"peso_gramos"`
		LongitudMM        int       `json:"longitud_mm"`
		HopUp             string    `json:"hop_up"`
		CapacidadCargador int       `json:"capacidad_cargador"`
		Notas             string    `json:"notas"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fechaAdq, _ := time.Parse("2006-01-02", req.FechaAdquisicion)

	replica := &models.Replica{
		Nombre:            req.Nombre,
		Marca:             req.Marca,
		Modelo:            req.Modelo,
		Tipo:              req.Tipo,
		NumeroSerie:       req.NumeroSerie,
		FechaAdquisicion:  fechaAdq,
		Proveedor:         req.Proveedor,
		CostoAdquisicion:  req.CostoAdquisicion,
		Estado:            req.Estado,
		FPS:               req.FPS,
		Joules:            req.Joules,
		PesoGramos:        req.PesoGramos,
		LongitudMM:        req.LongitudMM,
		HopUp:             req.HopUp,
		CapacidadCargador: req.CapacidadCargador,
		Notas:             req.Notas,
	}

	if err := h.service.Create(ctx, replica); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, replica)
}

// GetByID obtiene una réplica por ID
func (h *ReplicaHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	replica, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, replica)
}

// Update actualiza una réplica
func (h *ReplicaHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req struct {
		Nombre            string    `json:"nombre"`
		Marca             string    `json:"marca"`
		Modelo            string    `json:"modelo"`
		Tipo              string    `json:"tipo"`
		NumeroSerie       string    `json:"numero_serie"`
		FechaAdquisicion  string    `json:"fecha_adquisicion"`
		Proveedor         string    `json:"proveedor"`
		CostoAdquisicion  float64   `json:"costo_adquisicion"`
		Estado            string    `json:"estado"`
		FPS               int       `json:"fps"`
		Joules            float64   `json:"joules"`
		PesoGramos        int       `json:"peso_gramos"`
		LongitudMM        int       `json:"longitud_mm"`
		HopUp             string    `json:"hop_up"`
		CapacidadCargador int       `json:"capacidad_cargador"`
		Notas             string    `json:"notas"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fechaAdq, _ := time.Parse("2006-01-02", req.FechaAdquisicion)

	replica := &models.Replica{
		ID:                id,
		Nombre:            req.Nombre,
		Marca:             req.Marca,
		Modelo:            req.Modelo,
		Tipo:              req.Tipo,
		NumeroSerie:       req.NumeroSerie,
		FechaAdquisicion:  fechaAdq,
		Proveedor:         req.Proveedor,
		CostoAdquisicion:  req.CostoAdquisicion,
		Estado:            req.Estado,
		FPS:               req.FPS,
		Joules:            req.Joules,
		PesoGramos:        req.PesoGramos,
		LongitudMM:        req.LongitudMM,
		HopUp:             req.HopUp,
		CapacidadCargador: req.CapacidadCargador,
		Notas:             req.Notas,
	}

	if err := h.service.Update(ctx, replica); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, replica)
}

// Delete elimina una réplica
func (h *ReplicaHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "réplica eliminada"})
}