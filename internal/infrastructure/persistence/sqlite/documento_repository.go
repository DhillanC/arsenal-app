package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// DocumentoRepository implementa outbound.DocumentoRepository
type DocumentoRepository struct {
	db *sql.DB
}

// NewDocumentoRepository crea un nuevo repositorio
func NewDocumentoRepository(db *sql.DB) outbound.DocumentoRepository {
	return &DocumentoRepository{db: db}
}

// Create inserta un nuevo documento
func (r *DocumentoRepository) Create(ctx context.Context, documento *models.Documento) error {
	query := `
		INSERT INTO documentos (
			replica_id, actividad_id, tipo, nombre_archivo, ruta_archivo,
			mime_type, tamano_bytes, ocr_texto, fecha_documento, numero_documento, notas
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		documento.ReplicaID, documento.ActividadID, documento.Tipo,
		documento.NombreArchivo, documento.RutaArchivo, documento.MimeType,
		documento.TamanoBytes, documento.OCRTexto, documento.FechaDocumento,
		documento.NumeroDocumento, documento.Notas,
	)
	if err != nil {
		return fmt.Errorf("insertar documento: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener id: %w", err)
	}
	documento.ID = int(id)
	return nil
}

// GetByID obtiene un documento por su ID
func (r *DocumentoRepository) GetByID(ctx context.Context, id int) (*models.Documento, error) {
	query := `
		SELECT id, replica_id, actividad_id, tipo, nombre_archivo, ruta_archivo,
			mime_type, tamano_bytes, ocr_texto, fecha_documento, numero_documento,
			notas, created_at
		FROM documentos WHERE id = ?
	`

	var documento models.Documento
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&documento.ID, &documento.ReplicaID, &documento.ActividadID,
		&documento.Tipo, &documento.NombreArchivo, &documento.RutaArchivo,
		&documento.MimeType, &documento.TamanoBytes, &documento.OCRTexto,
		&documento.FechaDocumento, &documento.NumeroDocumento,
		&documento.Notas, &documento.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("documento no encontrado: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("consultar documento: %w", err)
	}
	return &documento, nil
}

// ListByReplica devuelve documentos de una réplica
func (r *DocumentoRepository) ListByReplica(ctx context.Context, replicaID int) ([]models.Documento, error) {
	query := `
		SELECT id, replica_id, actividad_id, tipo, nombre_archivo, ruta_archivo,
			mime_type, tamano_bytes, ocr_texto, fecha_documento, numero_documento,
			notas, created_at
		FROM documentos WHERE replica_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, replicaID)
	if err != nil {
		return nil, fmt.Errorf("listar documentos: %w", err)
	}
	defer rows.Close()

	var documentos []models.Documento
	for rows.Next() {
		var documento models.Documento
		err := rows.Scan(
			&documento.ID, &documento.ReplicaID, &documento.ActividadID,
			&documento.Tipo, &documento.NombreArchivo, &documento.RutaArchivo,
			&documento.MimeType, &documento.TamanoBytes, &documento.OCRTexto,
			&documento.FechaDocumento, &documento.NumeroDocumento,
			&documento.Notas, &documento.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan documento: %w", err)
		}
		documentos = append(documentos, documento)
	}

	return documentos, rows.Err()
}

// Update actualiza un documento
func (r *DocumentoRepository) Update(ctx context.Context, documento *models.Documento) error {
	query := `
		UPDATE documentos SET
			replica_id = ?, actividad_id = ?, tipo = ?, nombre_archivo = ?,
			ruta_archivo = ?, mime_type = ?, tamano_bytes = ?, ocr_texto = ?,
			fecha_documento = ?, numero_documento = ?, notas = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		documento.ReplicaID, documento.ActividadID, documento.Tipo,
		documento.NombreArchivo, documento.RutaArchivo, documento.MimeType,
		documento.TamanoBytes, documento.OCRTexto, documento.FechaDocumento,
		documento.NumeroDocumento, documento.Notas, documento.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizar documento: %w", err)
	}
	return nil
}

// Delete elimina un documento
func (r *DocumentoRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM documentos WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("eliminar documento: %w", err)
	}
	return nil
}

// ListByReplicaAndType devuelve documentos de una réplica filtrados por tipo
func (r *DocumentoRepository) ListByReplicaAndType(ctx context.Context, replicaID int, tipo string) ([]models.Documento, error) {
	query := `
		SELECT id, replica_id, actividad_id, tipo, nombre_archivo, ruta_archivo,
			mime_type, tamano_bytes, ocr_texto, fecha_documento, numero_documento,
			notas, created_at
		FROM documentos WHERE replica_id = ? AND tipo = ? ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, replicaID, tipo)
	if err != nil {
		return nil, fmt.Errorf("listar documentos por tipo: %w", err)
	}
	defer rows.Close()

	var documentos []models.Documento
	for rows.Next() {
		var documento models.Documento
		err := rows.Scan(
			&documento.ID, &documento.ReplicaID, &documento.ActividadID,
			&documento.Tipo, &documento.NombreArchivo, &documento.RutaArchivo,
			&documento.MimeType, &documento.TamanoBytes, &documento.OCRTexto,
			&documento.FechaDocumento, &documento.NumeroDocumento,
			&documento.Notas, &documento.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan documento: %w", err)
		}
		documentos = append(documentos, documento)
	}

	return documentos, rows.Err()
}

// ListByActividad devuelve documentos asociados a una actividad
func (r *DocumentoRepository) ListByActividad(ctx context.Context, actividadID int) ([]models.Documento, error) {
	query := `
		SELECT id, replica_id, actividad_id, tipo, nombre_archivo, ruta_archivo,
			mime_type, tamano_bytes, ocr_texto, fecha_documento, numero_documento,
			notas, created_at
		FROM documentos WHERE actividad_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, actividadID)
	if err != nil {
		return nil, fmt.Errorf("listar documentos por actividad: %w", err)
	}
	defer rows.Close()

	var documentos []models.Documento
	for rows.Next() {
		var documento models.Documento
		err := rows.Scan(
			&documento.ID, &documento.ReplicaID, &documento.ActividadID,
			&documento.Tipo, &documento.NombreArchivo, &documento.RutaArchivo,
			&documento.MimeType, &documento.TamanoBytes, &documento.OCRTexto,
			&documento.FechaDocumento, &documento.NumeroDocumento,
			&documento.Notas, &documento.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan documento: %w", err)
		}
		documentos = append(documentos, documento)
	}

	return documentos, rows.Err()
}

// SearchByOCR busca documentos por texto OCR (búsqueda simple con LIKE)
func (r *DocumentoRepository) SearchByOCR(ctx context.Context, query string) ([]models.Documento, error) {
	// Escapar caracteres LIKE para prevenir wildcard injection
	searchTerm := "%" + strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`).Replace(query) + "%"
	sqlQuery := `
		SELECT id, replica_id, actividad_id, tipo, nombre_archivo, ruta_archivo,
			mime_type, tamano_bytes, ocr_texto, fecha_documento, numero_documento,
			notas, created_at
		FROM documentos
		WHERE ocr_texto LIKE ? ESCAPE '\' OR numero_documento LIKE ? ESCAPE '\'
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, sqlQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("buscar documentos OCR: %w", err)
	}
	defer rows.Close()

	var documentos []models.Documento
	for rows.Next() {
		var documento models.Documento
		err := rows.Scan(
			&documento.ID, &documento.ReplicaID, &documento.ActividadID,
			&documento.Tipo, &documento.NombreArchivo, &documento.RutaArchivo,
			&documento.MimeType, &documento.TamanoBytes, &documento.OCRTexto,
			&documento.FechaDocumento, &documento.NumeroDocumento,
			&documento.Notas, &documento.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan documento: %w", err)
		}
		documentos = append(documentos, documento)
	}

	return documentos, rows.Err()
}
