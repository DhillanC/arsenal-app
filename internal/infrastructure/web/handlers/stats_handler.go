package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/DhillanC/arsenal-app/internal/infrastructure/persistence/sqlite"
	"github.com/gin-gonic/gin"
)

// StatsHandler maneja peticiones de estadísticas agregadas
type StatsHandler struct {
	db *sqlite.DB
}

// NewStatsHandler crea un nuevo handler de stats
func NewStatsHandler(db *sqlite.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

// DashboardStats devuelve estadísticas agregadas para el dashboard
// GET /api/v1/stats/dashboard
func (h *StatsHandler) DashboardStats(c *gin.Context) {
	ctx := c.Request.Context()

	// Stats por tipo de réplica
	typeStats, err := h.statsByTipo(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Stats por estado
	estadoStats, err := h.statsByEstado(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Valor total del inventario
	var valorTotal float64
	err = h.db.ReadConn.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(costo_adquisicion), 0) FROM replicas WHERE estado != 'archivado'
	`).Scan(&valorTotal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"por_tipo":    typeStats,
		"por_estado":  estadoStats,
		"valor_total": valorTotal,
		"timestamp":   time.Now().UTC(),
	})
}

func (h *StatsHandler) statsByTipo(ctx context.Context) ([]gin.H, error) {
	rows, err := h.db.ReadConn.QueryContext(ctx, `
		SELECT tipo, COUNT(*) as count, COALESCE(SUM(costo_adquisicion), 0) as valor
		FROM replicas WHERE estado != 'archivado' GROUP BY tipo
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []gin.H
	for rows.Next() {
		var tipo string
		var count int
		var valor float64
		if err := rows.Scan(&tipo, &count, &valor); err != nil {
			return nil, err
		}
		result = append(result, gin.H{"tipo": tipo, "count": count, "valor": valor})
	}
	return result, rows.Err()
}

func (h *StatsHandler) statsByEstado(ctx context.Context) ([]gin.H, error) {
	rows, err := h.db.ReadConn.QueryContext(ctx, `
		SELECT estado, COUNT(*) as count FROM replicas WHERE estado != 'archivado' GROUP BY estado
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []gin.H
	for rows.Next() {
		var estado string
		var count int
		if err := rows.Scan(&estado, &count); err != nil {
			return nil, err
		}
		result = append(result, gin.H{"estado": estado, "count": count})
	}
	return result, rows.Err()
}

// Export devuelve un backup JSON de todas las réplicas
// GET /api/v1/export/json
func (h *StatsHandler) ExportJSON(c *gin.Context) {
	ctx := c.Request.Context()

	rows, err := h.db.ReadConn.QueryContext(ctx, `
		SELECT id, nombre, marca, modelo, tipo, numero_serie, fecha_adquisicion,
			proveedor, costo_adquisicion, estado, fps, joules, peso_gramos,
			longitud_mm, hop_up, capacidad_cargador, notas, created_at, updated_at
		FROM replicas WHERE estado != 'archivado'
		ORDER BY id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var replicas []gin.H
	for rows.Next() {
		var r struct {
			ID                int
			Nombre            string
			Marca             string
			Modelo            string
			Tipo              string
			NumeroSerie       string
			FechaAdquisicion  time.Time
			Proveedor         string
			CostoAdquisicion  float64
			Estado            string
			FPS               int
			Joules            float64
			PesoGramos        int
			LongitudMM        int
			HopUp             string
			CapacidadCargador int
			Notas             string
			CreatedAt         time.Time
			UpdatedAt         time.Time
		}
		if err := rows.Scan(
			&r.ID, &r.Nombre, &r.Marca, &r.Modelo, &r.Tipo, &r.NumeroSerie,
			&r.FechaAdquisicion, &r.Proveedor, &r.CostoAdquisicion, &r.Estado,
			&r.FPS, &r.Joules, &r.PesoGramos, &r.LongitudMM, &r.HopUp,
			&r.CapacidadCargador, &r.Notas, &r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		replicas = append(replicas, gin.H{
			"id":                r.ID,
			"nombre":            r.Nombre,
			"marca":             r.Marca,
			"modelo":            r.Modelo,
			"tipo":              r.Tipo,
			"numero_serie":      r.NumeroSerie,
			"fecha_adquisicion": r.FechaAdquisicion,
			"proveedor":         r.Proveedor,
			"costo_adquisicion": r.CostoAdquisicion,
			"estado":            r.Estado,
			"fps":               r.FPS,
			"joules":            r.Joules,
			"peso_gramos":       r.PesoGramos,
			"longitud_mm":       r.LongitudMM,
			"hop_up":            r.HopUp,
			"capacidad_cargador": r.CapacidadCargador,
			"notas":             r.Notas,
			"created_at":        r.CreatedAt,
			"updated_at":        r.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"export_date": time.Now().UTC(),
		"version":     "1.0",
		"replicas":    replicas,
		"count":       len(replicas),
	})
}
