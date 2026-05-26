package services

import (
	"context"
	"fmt"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// DocumentoService implementa inbound.DocumentoService
type DocumentoService struct {
	repo    outbound.DocumentoRepository
	storage outbound.Storage
}

// NewDocumentoService crea un nuevo servicio
func NewDocumentoService(repo outbound.DocumentoRepository, storage outbound.Storage) inbound.DocumentoService {
	return &DocumentoService{repo: repo, storage: storage}
}

// Create crea un nuevo documento con archivo
func (s *DocumentoService) Create(ctx context.Context, documento *models.Documento, file []byte) error {
	if documento.Tipo == "" {
		return fmt.Errorf("tipo es requerido")
	}
	if documento.NombreArchivo == "" {
		return fmt.Errorf("nombre_archivo es requerido")
	}

	// Guardar archivo en storage
	if len(file) > 0 {
		ruta, err := s.storage.Save(file, documento.NombreArchivo, *documento.ReplicaID)
		if err != nil {
			return fmt.Errorf("guardar archivo: %w", err)
		}
		documento.RutaArchivo = ruta
		documento.TamanoBytes = int64(len(file))
	}

	return s.repo.Create(ctx, documento)
}

// GetByID obtiene un documento por ID
func (s *DocumentoService) GetByID(ctx context.Context, id int) (*models.Documento, error) {
	return s.repo.GetByID(ctx, id)
}

// ListByReplica lista documentos de una réplica
func (s *DocumentoService) ListByReplica(ctx context.Context, replicaID int) ([]models.Documento, error) {
	return s.repo.ListByReplica(ctx, replicaID)
}

// ListByReplicaAndType lista documentos de una réplica filtrados por tipo
func (s *DocumentoService) ListByReplicaAndType(ctx context.Context, replicaID int, tipo string) ([]models.Documento, error) {
	return s.repo.ListByReplicaAndType(ctx, replicaID, tipo)
}

// ListByActividad lista documentos asociados a una actividad
func (s *DocumentoService) ListByActividad(ctx context.Context, actividadID int) ([]models.Documento, error) {
	return s.repo.ListByActividad(ctx, actividadID)
}

// Update actualiza un documento
func (s *DocumentoService) Update(ctx context.Context, documento *models.Documento) error {
	if documento.ID == 0 {
		return fmt.Errorf("id es requerido")
	}
	return s.repo.Update(ctx, documento)
}

// Delete elimina un documento
func (s *DocumentoService) Delete(ctx context.Context, id int) error {
	// TODO: Eliminar archivo del storage también
	return s.repo.Delete(ctx, id)
}

// SearchByOCR busca documentos por texto OCR
func (s *DocumentoService) SearchByOCR(ctx context.Context, query string) ([]models.Documento, error) {
	return s.repo.SearchByOCR(ctx, query)
}
