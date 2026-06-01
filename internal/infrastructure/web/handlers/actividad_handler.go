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
// y el timeline combinado con documentos
type ActividadHandler struct {
	service        inbound.ActividadService
	documentoService inbound.DocumentoService
}

// NewActividadHandler crea un nuevo handler
func NewActividadHandler(service inbound.ActividadService, documentoService inbound.DocumentoService) *ActividadHandler {
	return &ActividadHandler{
		service:        service,
		documentoService: documentoService,
	}
}

// RegisterRoutes registra las rutas de actividades
func (h *ActividadHandler) RegisterRoutes(router *gin.RouterGroup) {
	actividades := router.Group("/replicas/:id/actividades")
	{
		actividades.GET("", h.ListByReplica)
		actividades.POST("", h.Create)
		actividades.GET("/timeline", h.Timeline)
	}

	// Rutas independientes para operaciones por ID de actividad
	router.GET("/actividades/:actividad_id", h.GetByID)
	router.PUT("/actividades/:actividad_id", h.Update)
	router.DELETE("/actividades/:actividad_id", h.Delete)
}

// ListByReplica lista actividades de una réplica. Soporta paginación via ?limit=20&offset=0.
func (h *ActividadHandler) ListByReplica(c *gin.Context) {
	ctx := c.Request.Context()
	replicaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de réplica inválido"})
		return
	}

	var actividades []models.Actividad
	if c.Query("limit") != "" || c.Query("offset") != "" {
		limit, offset := PaginationParams(c)
		actividades, err = h.service.ListByReplicaPaginated(ctx, replicaID, limit, offset)
	} else {
		actividades, err = h.service.ListByReplica(ctx, replicaID)
	}
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

// TimelineItem representa un item en el timeline (actividad + documentos)
type TimelineItem struct {
	ID               int                      `json:"id"`
	Fecha            time.Time                `json:"fecha"`
	Tipo             string                   `json:"tipo"`
	Descripcion      string                   `json:"descripcion"`
	ProveedorTecnico string                   `json:"proveedor_tecnico,omitempty"`
	Costo            float64                  `json:"costo,omitempty"`
	KilometrajeBB    int                      `json:"kilometraje_bb,omitempty"`
	Ubicacion        string                   `json:"ubicacion,omitempty"`
	Documentos       []models.Documento       `json:"documentos"`
}

// Timeline devuelve el timeline cronológico de una réplica
// con actividades y sus documentos asociados
func (h *ActividadHandler) Timeline(c *gin.Context) {
	ctx := c.Request.Context()
	replicaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de réplica inválido"})
		return
	}

	// Obtener actividades
	actividades, err := h.service.ListByReplica(ctx, replicaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Construir timeline con documentos
	timeline := make([]TimelineItem, 0, len(actividades))
	
	// Obtener documentos de todas las actividades en una sola query (evita N+1)
	actividadIDs := make([]int, len(actividades))
	for i, act := range actividades {
		actividadIDs[i] = act.ID
	}
	
	var docsByActividad map[int][]models.Documento
	if len(actividadIDs) > 0 {
		allDocs, err := h.documentoService.ListByActividades(ctx, actividadIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Agrupar documentos por actividad_id
		docsByActividad = make(map[int][]models.Documento)
		for _, doc := range allDocs {
			if doc.ActividadID != nil {
				docsByActividad[*doc.ActividadID] = append(docsByActividad[*doc.ActividadID], doc)
			}
		}
	}
	
	for _, act := range actividades {
		item := TimelineItem{
			ID:               act.ID,
			Fecha:            act.Fecha,
			Tipo:             act.Tipo,
			Descripcion:      act.Descripcion,
			ProveedorTecnico: act.ProveedorTecnico,
			Costo:            act.Costo,
			KilometrajeBB:    act.KilometrajeBB,
			Ubicacion:        act.Ubicacion,
			Documentos:       docsByActividad[act.ID],
		}
		timeline = append(timeline, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"replica_id": replicaID,
		"count":      len(timeline),
		"timeline":   timeline,
	})
}
