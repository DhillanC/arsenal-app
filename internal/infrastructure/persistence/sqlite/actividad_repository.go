package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/digital-consultory-solutions/arsenal-app/internal/domain/models"
	"github.com/digital-consultory-solutions/arsenal-app/internal/domain/ports/outbound"
)

// ActividadRepository implementa outbound.ActividadRepository
type ActividadRepository struct {
	db *sql.DB
}

// NewActividadRepository crea un nuevo repositorio
func NewActividadRepository(db *sql.DB) outbound.ActividadRepository {
	return &ActividadRepository{db: db}
}

// Create inserta una nueva actividad
func (r *ActividadRepository) Create(ctx context.Context, actividad *models.Actividad) error {
	query := `
		INSERT INTO actividades (
			replica_id, fecha, tipo, descripcion, proveedor_tecnico,
			costo, kilometraje_bb, ubicacion
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		actividad.ReplicaID, actividad.Fecha, actividad.Tipo,
		actividad.Descripcion, actividad.ProveedorTecnico,
		actividad.Costo, actividad.KilometrajeBB, actividad.Ubicacion,
	)
	if err != nil {
		return fmt.Errorf("insertar actividad: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener id: %w", err)
	}
	actividad.ID = int(id)
	return nil
}

// GetByID obtiene una actividad por su ID
func (r *ActividadRepository) GetByID(ctx context.Context, id int) (*models.Actividad, error) {
	query := `
		SELECT id, replica_id, fecha, tipo, descripcion, proveedor_tecnico,
			costo, kilometraje_bb, ubicacion, created_at
		FROM actividades WHERE id = ?
	`

	var actividad models.Actividad
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&actividad.ID, &actividad.ReplicaID, &actividad.Fecha, &actividad.Tipo,
		&actividad.Descripcion, &actividad.ProveedorTecnico,
		&actividad.Costo, &actividad.KilometrajeBB, &actividad.Ubicacion,
		&actividad.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("actividad no encontrada: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("consultar actividad: %w", err)
	}
	return &actividad, nil
}

// ListByReplica devuelve actividades de una réplica
func (r *ActividadRepository) ListByReplica(ctx context.Context, replicaID int) ([]models.Actividad, error) {
	query := `
		SELECT id, replica_id, fecha, tipo, descripcion, proveedor_tecnico,
			costo, kilometraje_bb, ubicacion, created_at
		FROM actividades WHERE replica_id = ? ORDER BY fecha DESC
	`

	rows, err := r.db.QueryContext(ctx, query, replicaID)
	if err != nil {
		return nil, fmt.Errorf("listar actividades: %w", err)
	}
	defer rows.Close()

	var actividades []models.Actividad
	for rows.Next() {
		var actividad models.Actividad
		err := rows.Scan(
			&actividad.ID, &actividad.ReplicaID, &actividad.Fecha, &actividad.Tipo,
			&actividad.Descripcion, &actividad.ProveedorTecnico,
			&actividad.Costo, &actividad.KilometrajeBB, &actividad.Ubicacion,
			&actividad.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan actividad: %w", err)
		}
		actividades = append(actividades, actividad)
	}

	return actividades, rows.Err()
}

// Update actualiza una actividad
func (r *ActividadRepository) Update(ctx context.Context, actividad *models.Actividad) error {
	query := `
		UPDATE actividades SET
			replica_id = ?, fecha = ?, tipo = ?, descripcion = ?,
			proveedor_tecnico = ?, costo = ?, kilometraje_bb = ?, ubicacion = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		actividad.ReplicaID, actividad.Fecha, actividad.Tipo,
		actividad.Descripcion, actividad.ProveedorTecnico,
		actividad.Costo, actividad.KilometrajeBB, actividad.Ubicacion,
		actividad.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizar actividad: %w", err)
	}
	return nil
}

// Delete elimina una actividad
func (r *ActividadRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM actividades WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("eliminar actividad: %w", err)
	}
	return nil
}
