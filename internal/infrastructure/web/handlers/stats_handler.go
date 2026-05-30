package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StatsHandler maneja peticiones de estadísticas agregadas
type StatsHandler struct {
	db *sql.DB
}

// NewStatsHandler crea un nuevo handler de stats
func NewStatsHandler(db *sql.DB) *StatsHandler {
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
	err = h.db.QueryRowContext(ctx, `
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
	rows, err := h.db.QueryContext(ctx, `
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
	rows, err := h.db.QueryContext(ctx, `
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
