package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// ReplicaRepository implementa outbound.ReplicaRepository
type ReplicaRepository struct {
	db *sql.DB
}

// NewReplicaRepository crea un nuevo repositorio
func NewReplicaRepository(db *sql.DB) outbound.ReplicaRepository {
	return &ReplicaRepository{db: db}
}

// Create inserta una nueva réplica
func (r *ReplicaRepository) Create(ctx context.Context, replica *models.Replica) error {
	query := `
		INSERT INTO replicas (
			nombre, marca, modelo, tipo, numero_serie, fecha_adquisicion,
			proveedor, costo_adquisicion, estado, fps, joules, peso_gramos,
			longitud_mm, hop_up, capacidad_cargador, notas
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		replica.Nombre, replica.Marca, replica.Modelo, replica.Tipo,
		replica.NumeroSerie, replica.FechaAdquisicion, replica.Proveedor,
		replica.CostoAdquisicion, replica.Estado, replica.FPS, replica.Joules,
		replica.PesoGramos, replica.LongitudMM, replica.HopUp,
		replica.CapacidadCargador, replica.Notas,
	)
	if err != nil {
		return fmt.Errorf("insertar replica: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener id: %w", err)
	}
	replica.ID = int(id)
	return nil
}

// GetByID obtiene una réplica por su ID
func (r *ReplicaRepository) GetByID(ctx context.Context, id int) (*models.Replica, error) {
	query := `
		SELECT id, nombre, marca, modelo, tipo, numero_serie, fecha_adquisicion,
			proveedor, costo_adquisicion, estado, fps, joules, peso_gramos,
			longitud_mm, hop_up, capacidad_cargador, notas, created_at, updated_at
		FROM replicas WHERE id = ?
	`

	var replica models.Replica
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&replica.ID, &replica.Nombre, &replica.Marca, &replica.Modelo,
		&replica.Tipo, &replica.NumeroSerie, &replica.FechaAdquisicion,
		&replica.Proveedor, &replica.CostoAdquisicion, &replica.Estado,
		&replica.FPS, &replica.Joules, &replica.PesoGramos, &replica.LongitudMM,
		&replica.HopUp, &replica.CapacidadCargador, &replica.Notas,
		&replica.CreatedAt, &replica.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("replica no encontrada: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("consultar replica: %w", err)
	}
	return &replica, nil
}

// List devuelve todas las réplicas activas (no archivadas)
func (r *ReplicaRepository) List(ctx context.Context) ([]models.Replica, error) {
	query := `
		SELECT id, nombre, marca, modelo, tipo, numero_serie, fecha_adquisicion,
			proveedor, costo_adquisicion, estado, fps, joules, peso_gramos,
			longitud_mm, hop_up, capacidad_cargador, notas, created_at, updated_at
		FROM replicas 
		WHERE estado != 'archivado'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("listar replicas: %w", err)
	}
	defer rows.Close()

	var replicas []models.Replica
	for rows.Next() {
		var replica models.Replica
		err := rows.Scan(
			&replica.ID, &replica.Nombre, &replica.Marca, &replica.Modelo,
			&replica.Tipo, &replica.NumeroSerie, &replica.FechaAdquisicion,
			&replica.Proveedor, &replica.CostoAdquisicion, &replica.Estado,
			&replica.FPS, &replica.Joules, &replica.PesoGramos, &replica.LongitudMM,
			&replica.HopUp, &replica.CapacidadCargador, &replica.Notas,
			&replica.CreatedAt, &replica.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan replica: %w", err)
		}
		replicas = append(replicas, replica)
	}

	return replicas, rows.Err()
}

// Update actualiza una réplica
func (r *ReplicaRepository) Update(ctx context.Context, replica *models.Replica) error {
	query := `
		UPDATE replicas SET
			nombre = ?, marca = ?, modelo = ?, tipo = ?, numero_serie = ?,
			fecha_adquisicion = ?, proveedor = ?, costo_adquisicion = ?,
			estado = ?, fps = ?, joules = ?, peso_gramos = ?, longitud_mm = ?,
			hop_up = ?, capacidad_cargador = ?, notas = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		replica.Nombre, replica.Marca, replica.Modelo, replica.Tipo,
		replica.NumeroSerie, replica.FechaAdquisicion, replica.Proveedor,
		replica.CostoAdquisicion, replica.Estado, replica.FPS, replica.Joules,
		replica.PesoGramos, replica.LongitudMM, replica.HopUp,
		replica.CapacidadCargador, replica.Notas, replica.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizar replica: %w", err)
	}
	return nil
}

// Delete elimina una réplica (soft delete cambiando estado)
func (r *ReplicaRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE replicas SET estado = 'archivado', updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("eliminar replica: %w", err)
	}
	return nil
}

// Search busca réplicas por número de serie o nombre (para trazabilidad DIAN)
func (r *ReplicaRepository) Search(ctx context.Context, query string) ([]models.Replica, error) {
	searchTerm := "%" + query + "%"
	sqlQuery := `
		SELECT id, nombre, marca, modelo, tipo, numero_serie, fecha_adquisicion,
			proveedor, costo_adquisicion, estado, fps, joules, peso_gramos,
			longitud_mm, hop_up, capacidad_cargador, notas, created_at, updated_at
		FROM replicas 
		WHERE estado != 'archivado' AND (
			nombre LIKE ? OR 
			numero_serie LIKE ? OR 
			marca LIKE ? OR
			modelo LIKE ?
		)
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("buscar replicas: %w", err)
	}
	defer rows.Close()

	var replicas []models.Replica
	for rows.Next() {
		var replica models.Replica
		err := rows.Scan(
			&replica.ID, &replica.Nombre, &replica.Marca, &replica.Modelo,
			&replica.Tipo, &replica.NumeroSerie, &replica.FechaAdquisicion,
			&replica.Proveedor, &replica.CostoAdquisicion, &replica.Estado,
			&replica.FPS, &replica.Joules, &replica.PesoGramos, &replica.LongitudMM,
			&replica.HopUp, &replica.CapacidadCargador, &replica.Notas,
			&replica.CreatedAt, &replica.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan replica: %w", err)
		}
		replicas = append(replicas, replica)
	}

	return replicas, rows.Err()
}