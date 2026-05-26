package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// MantenimientoRepository implementa outbound.MantenimientoRepository
type MantenimientoRepository struct {
	db *sql.DB
}

// NewMantenimientoRepository crea un nuevo repositorio de mantenimiento
func NewMantenimientoRepository(db *sql.DB) outbound.MantenimientoRepository {
	return &MantenimientoRepository{db: db}
}

// Create inserta un nuevo registro de mantenimiento
func (r *MantenimientoRepository) Create(ctx context.Context, m *models.Mantenimiento) error {
	query := `
		INSERT INTO mantenimiento (replica_id, tipo_tarea, frecuencia_dias, frecuencia_bb, ultima_fecha, proxima_fecha, completado, notas)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		m.ReplicaID, m.TipoTarea, m.FrecuenciaDias, m.FrecuenciaBB,
		m.UltimaFecha, m.ProximaFecha, m.Completado, m.Notas,
	)
	if err != nil {
		return fmt.Errorf("insertar mantenimiento: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener id: %w", err)
	}
	m.ID = int(id)
	return nil
}

// GetByID obtiene un mantenimiento por ID
func (r *MantenimientoRepository) GetByID(ctx context.Context, id int) (*models.Mantenimiento, error) {
	query := `
		SELECT id, replica_id, tipo_tarea, frecuencia_dias, frecuencia_bb, ultima_fecha, proxima_fecha, completado, notas
		FROM mantenimiento WHERE id = ?
	`
	var m models.Mantenimiento
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.ReplicaID, &m.TipoTarea, &m.FrecuenciaDias, &m.FrecuenciaBB,
		&m.UltimaFecha, &m.ProximaFecha, &m.Completado, &m.Notas,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("mantenimiento no encontrado: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("consultar mantenimiento: %w", err)
	}
	return &m, nil
}

// ListByReplica lista mantenimientos de una réplica
func (r *MantenimientoRepository) ListByReplica(ctx context.Context, replicaID int) ([]models.Mantenimiento, error) {
	query := `
		SELECT id, replica_id, tipo_tarea, frecuencia_dias, frecuencia_bb, ultima_fecha, proxima_fecha, completado, notas
		FROM mantenimiento WHERE replica_id = ? ORDER BY proxima_fecha ASC
	`
	rows, err := r.db.QueryContext(ctx, query, replicaID)
	if err != nil {
		return nil, fmt.Errorf("listar mantenimientos: %w", err)
	}
	defer rows.Close()

	var mantenimientos []models.Mantenimiento
	for rows.Next() {
		var m models.Mantenimiento
		if err := rows.Scan(&m.ID, &m.ReplicaID, &m.TipoTarea, &m.FrecuenciaDias, &m.FrecuenciaBB,
			&m.UltimaFecha, &m.ProximaFecha, &m.Completado, &m.Notas); err != nil {
			return nil, fmt.Errorf("scan mantenimiento: %w", err)
		}
		mantenimientos = append(mantenimientos, m)
	}
	return mantenimientos, rows.Err()
}

// ListProximos lista mantenimientos próximos a vencer
func (r *MantenimientoRepository) ListProximos(ctx context.Context, dias int) ([]models.Mantenimiento, error) {
	query := `
		SELECT id, replica_id, tipo_tarea, frecuencia_dias, frecuencia_bb, ultima_fecha, proxima_fecha, completado, notas
		FROM mantenimiento
		WHERE completado = 0 AND (proxima_fecha IS NULL OR proxima_fecha <= date('now', '+' || ? || ' days'))
		ORDER BY proxima_fecha ASC
	`
	rows, err := r.db.QueryContext(ctx, query, dias)
	if err != nil {
		return nil, fmt.Errorf("listar mantenimientos próximos: %w", err)
	}
	defer rows.Close()

	var mantenimientos []models.Mantenimiento
	for rows.Next() {
		var m models.Mantenimiento
		if err := rows.Scan(&m.ID, &m.ReplicaID, &m.TipoTarea, &m.FrecuenciaDias, &m.FrecuenciaBB,
			&m.UltimaFecha, &m.ProximaFecha, &m.Completado, &m.Notas); err != nil {
			return nil, fmt.Errorf("scan mantenimiento: %w", err)
		}
		mantenimientos = append(mantenimientos, m)
	}
	return mantenimientos, rows.Err()
}

// Update actualiza un mantenimiento
func (r *MantenimientoRepository) Update(ctx context.Context, m *models.Mantenimiento) error {
	query := `
		UPDATE mantenimiento SET
			replica_id = ?, tipo_tarea = ?, frecuencia_dias = ?, frecuencia_bb = ?,
			ultima_fecha = ?, proxima_fecha = ?, completado = ?, notas = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		m.ReplicaID, m.TipoTarea, m.FrecuenciaDias, m.FrecuenciaBB,
		m.UltimaFecha, m.ProximaFecha, m.Completado, m.Notas, m.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizar mantenimiento: %w", err)
	}
	return nil
}

// Delete elimina un mantenimiento
func (r *MantenimientoRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM mantenimiento WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("eliminar mantenimiento: %w", err)
	}
	return nil
}

// MarcarCompletado marca un mantenimiento como completado
func (r *MantenimientoRepository) MarcarCompletado(ctx context.Context, id int, fechaCompletado *time.Time) error {
	now := time.Now()
	if fechaCompletado != nil {
		now = *fechaCompletado
	}
	
	// Obtener el mantenimiento para calcular próxima fecha
	m, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Calcular próxima fecha basada en frecuencia
	var proximaFecha *time.Time
	if m.FrecuenciaDias > 0 {
		next := now.AddDate(0, 0, m.FrecuenciaDias)
		proximaFecha = &next
	}
	
	query := `
		UPDATE mantenimiento SET
			completado = 1, ultima_fecha = ?, proxima_fecha = ?
		WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, query, now, proximaFecha, id)
	if err != nil {
		return fmt.Errorf("marcar completado: %w", err)
	}
	return nil
}
