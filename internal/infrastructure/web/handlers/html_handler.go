package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	"github.com/gin-gonic/gin"
)

// HTMLHandler maneja las vistas HTML del frontend
type HTMLHandler struct {
	replicaService       inbound.ReplicaService
	actividadService     inbound.ActividadService
	documentoService     inbound.DocumentoService
	mantenimientoService inbound.MantenimientoService
	uploadPath           string
}

// NewHTMLHandler crea un nuevo handler HTML
func NewHTMLHandler(
	replicaService inbound.ReplicaService,
	actividadService inbound.ActividadService,
	documentoService inbound.DocumentoService,
	mantenimientoService inbound.MantenimientoService,
	uploadPath string,
) *HTMLHandler {
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	return &HTMLHandler{
		replicaService:       replicaService,
		actividadService:     actividadService,
		documentoService:     documentoService,
		mantenimientoService: mantenimientoService,
		uploadPath:           uploadPath,
	}
}

// RegisterHTMLRoutes registra las rutas HTML
func (h *HTMLHandler) RegisterHTMLRoutes(router *gin.Engine) {
	// Static files
	router.Static("/static", "./web/static")
	router.Static("/uploads", h.uploadPath)

	// HTML routes
	router.GET("/", h.Dashboard)
	router.GET("/dashboard", h.Dashboard)
	router.GET("/replicas", h.ReplicaList)
	router.GET("/replicas/nueva", h.ReplicaCreateForm)
	router.GET("/replicas/:id", h.ReplicaDetail)
	router.GET("/replicas/:id/editar", h.ReplicaEditForm)
	router.GET("/documentos", h.DocumentList)
	router.GET("/mantenimiento", h.MantenimientoList)
}

// Dashboard muestra el dashboard principal
func (h *HTMLHandler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener réplicas
	replicas, err := h.replicaService.List(ctx)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Calcular estadísticas
	stats := calculateDashboardStats(replicas)

	// Obtener actividades recientes
	var actividadesRecientes []ActividadResumen
	var actividadesErrors []string
	for _, r := range replicas {
		acts, err := h.actividadService.ListByReplica(ctx, r.ID)
		if err != nil {
			actividadesErrors = append(actividadesErrors, fmt.Sprintf("réplica %d: %v", r.ID, err))
			continue
		}
		for _, a := range acts {
			actividadesRecientes = append(actividadesRecientes, ActividadResumen{
				Descripcion:   a.Descripcion,
				Fecha:         a.Fecha,
				Costo:         a.Costo,
				ReplicaNombre: r.Nombre,
			})
		}
	}

	// Limitar a las 10 más recientes
	if len(actividadesRecientes) > 10 {
		actividadesRecientes = actividadesRecientes[:10]
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title":                "Dashboard",
		"Stats":                stats,
		"ActividadesRecientes": actividadesRecientes,
		"ActividadesErrors":    actividadesErrors,
	})
}

// ReplicaList muestra la lista de réplicas
func (h *HTMLHandler) ReplicaList(c *gin.Context) {
	ctx := c.Request.Context()

	replicas, err := h.replicaService.List(ctx)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	stats := calculateDashboardStats(replicas)

	c.HTML(http.StatusOK, "replica_list.html", gin.H{
		"Title":    "Mis Réplicas",
		"Replicas": replicas,
		"Stats":    stats,
	})
}

// ReplicaDetail muestra la ficha de una réplica
func (h *HTMLHandler) ReplicaDetail(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "ID inválido"})
		return
	}

	replica, err := h.replicaService.GetByID(ctx, id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Réplica no encontrada"})
		return
	}

	actividades, _ := h.actividadService.ListByReplica(ctx, id)
	documentos, _ := h.documentoService.ListByReplica(ctx, id)
	mantenimientos, _ := h.mantenimientoService.ListByReplica(ctx, id)

	// Construir timeline
	timeline := buildHTMLTimeline(actividades, documentos)

	c.HTML(http.StatusOK, "replica_detail.html", gin.H{
		"Title":          replica.Nombre,
		"Replica":        replica,
		"Timeline":       timeline,
		"Documentos":     documentos,
		"Mantenimientos": mantenimientos,
	})
}

// ReplicaCreateForm muestra el formulario de creación
func (h *HTMLHandler) ReplicaCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "replica_form.html", gin.H{
		"Title":    "Nueva Réplica",
		"EditMode": false,
		"Replica":  models.Replica{Estado: "activo"},
	})
}

// ReplicaEditForm muestra el formulario de edición
func (h *HTMLHandler) ReplicaEditForm(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "ID inválido"})
		return
	}

	replica, err := h.replicaService.GetByID(ctx, id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Rplica no encontrada"})
		return
	}

	c.HTML(http.StatusOK, "replica_form.html", gin.H{
		"Title":    "Editar " + replica.Nombre,
		"EditMode": true,
		"Replica":  replica,
	})
}

// DocumentList muestra la lista de documentos
func (h *HTMLHandler) DocumentList(c *gin.Context) {
	replicas, err := h.replicaService.List(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "document_list.html", gin.H{
		"Title":    "Documentos",
		"Replicas": replicas,
	})
}

// MantenimientoList muestra los mantenimientos próximos.
func (h *HTMLHandler) MantenimientoList(c *gin.Context) {
	ctx := c.Request.Context()

	mantenimientos, err := h.mantenimientoService.ListProximos(ctx, 90)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	replicas, err := h.replicaService.List(ctx)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	replicaNames := make(map[int]string, len(replicas))
	for _, r := range replicas {
		replicaNames[r.ID] = r.Nombre
	}

	c.HTML(http.StatusOK, "mantenimiento.html", gin.H{
		"Title":          "Mantenimiento",
		"Mantenimientos": mantenimientos,
		"ReplicaNames":   replicaNames,
	})
}

// Tipos de datos para el template

type DashboardStats struct {
	Total      int
	Activas    int
	Reparacion int
	ValorTotal float64
	PorTipo    []TipoStat
	PorEstado  []EstadoStat
}

type TipoStat struct {
	Tipo       string
	Cantidad   int
	Porcentaje float64
}

type EstadoStat struct {
	Estado   string
	Cantidad int
}

type ActividadResumen struct {
	Descripcion   string
	Fecha         time.Time
	Costo         float64
	ReplicaNombre string
}

type HTMLTimelineItem struct {
	ID               int
	Fecha            time.Time
	Tipo             string
	Descripcion      string
	ProveedorTecnico string
	Costo            float64
	KilometrajeBB    int
	Ubicacion        string
	Documentos       []models.Documento
}

// calculateDashboardStats calcula las estadísticas
func calculateDashboardStats(replicas []models.Replica) DashboardStats {
	stats := DashboardStats{}

	tipoCount := make(map[string]int)
	estadoCount := make(map[string]int)

	for _, r := range replicas {
		stats.Total++
		stats.ValorTotal += r.CostoAdquisicion

		if r.Estado == "activo" {
			stats.Activas++
		} else if r.Estado == "reparacion" {
			stats.Reparacion++
		}

		tipoCount[r.Tipo]++
		estadoCount[r.Estado]++
	}

	// Calcular porcentajes para tipos
	for tipo, count := range tipoCount {
		porcentaje := 0.0
		if stats.Total > 0 {
			porcentaje = float64(count) / float64(stats.Total) * 100
		}
		stats.PorTipo = append(stats.PorTipo, TipoStat{
			Tipo:       tipo,
			Cantidad:   count,
			Porcentaje: porcentaje,
		})
	}

	// Estados
	for estado, count := range estadoCount {
		stats.PorEstado = append(stats.PorEstado, EstadoStat{
			Estado:   estado,
			Cantidad: count,
		})
	}

	return stats
}

// buildHTMLTimeline construye el timeline con actividades y documentos
func buildHTMLTimeline(actividades []models.Actividad, documentos []models.Documento) []HTMLTimelineItem {
	var timeline []HTMLTimelineItem

	for _, act := range actividades {
		item := HTMLTimelineItem{
			ID:               act.ID,
			Fecha:            act.Fecha,
			Tipo:             act.Tipo,
			Descripcion:      act.Descripcion,
			ProveedorTecnico: act.ProveedorTecnico,
			Costo:            act.Costo,
			KilometrajeBB:    act.KilometrajeBB,
			Ubicacion:        act.Ubicacion,
			Documentos:       []models.Documento{},
		}

		// Buscar documentos de esta actividad
		for _, doc := range documentos {
			if doc.ActividadID != nil && *doc.ActividadID == act.ID {
				item.Documentos = append(item.Documentos, doc)
			}
		}

		timeline = append(timeline, item)
	}

	return timeline
}
