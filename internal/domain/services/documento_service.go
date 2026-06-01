package services

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/ocr"
)

// DocumentoService implementa inbound.DocumentoService
type DocumentoService struct {
	repo    outbound.DocumentoRepository
	storage outbound.Storage
	ocrWg   sync.WaitGroup // Para esperar OCR en tests
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

		// Marcar para OCR async si aplica. "skipped" indica que OCR no se
		// intentó (mime no soportado o OCR_ENABLED=false). Distinto de
		// "completed" — que implica que OCR sí corrió y dejó texto (o vacío).
		if shouldRunOCR(documento) {
			documento.OCRStatus = "pending"
		} else {
			documento.OCRStatus = "skipped"
		}
	}

	// Guardar documento inmediatamente (sin esperar OCR)
	if err := s.repo.Create(ctx, documento); err != nil {
		return err
	}

	// Lanzar OCR en background si aplica
	if documento.OCRStatus == "pending" {
		s.ocrWg.Add(1)
		go func() {
			defer s.ocrWg.Done()
			s.processOCRAsync(documento.ID, documento.RutaArchivo)
		}()
	}

	return nil
}

// GetByID obtiene un documento por ID
func (s *DocumentoService) GetByID(ctx context.Context, id int) (*models.Documento, error) {
	return s.repo.GetByID(ctx, id)
}

// ListByReplica lista documentos de una réplica
func (s *DocumentoService) ListByReplica(ctx context.Context, replicaID int) ([]models.Documento, error) {
	return s.repo.ListByReplica(ctx, replicaID)
}

// ListByReplicaPaginated lista documentos de una réplica con paginación
func (s *DocumentoService) ListByReplicaPaginated(ctx context.Context, replicaID int, limit, offset int) ([]models.Documento, error) {
	return s.repo.ListByReplicaPaginated(ctx, replicaID, limit, offset)
}

// ListByReplicaAndType lista documentos de una réplica filtrados por tipo
func (s *DocumentoService) ListByReplicaAndType(ctx context.Context, replicaID int, tipo string) ([]models.Documento, error) {
	return s.repo.ListByReplicaAndType(ctx, replicaID, tipo)
}

// ListByActividad lista documentos asociados a una actividad
func (s *DocumentoService) ListByActividad(ctx context.Context, actividadID int) ([]models.Documento, error) {
	return s.repo.ListByActividad(ctx, actividadID)
}

// ListByActividades lista documentos asociados a múltiples actividades (batch).
func (s *DocumentoService) ListByActividades(ctx context.Context, actividadIDs []int) ([]models.Documento, error) {
	return s.repo.ListByActividades(ctx, actividadIDs)
}

// Update actualiza un documento
func (s *DocumentoService) Update(ctx context.Context, documento *models.Documento) error {
	if documento.ID == 0 {
		return fmt.Errorf("id es requerido")
	}
	return s.repo.Update(ctx, documento)
}

// Delete elimina un documento y su archivo asociado
func (s *DocumentoService) Delete(ctx context.Context, id int) error {
	// Obtener documento para saber qué archivo borrar
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("obtener documento: %w", err)
	}

	// Borrar archivo del storage si existe
	if doc.RutaArchivo != "" {
		if err := s.storage.Delete(doc.RutaArchivo); err != nil {
			// Loguear pero no fallar si el archivo ya no existe
			fmt.Printf("advertencia: no se pudo borrar archivo %s: %v\n", doc.RutaArchivo, err)
		}
	}

	return s.repo.Delete(ctx, id)
}

// SearchByOCR busca documentos por texto OCR
func (s *DocumentoService) SearchByOCR(ctx context.Context, query string) ([]models.Documento, error) {
	return s.repo.SearchByOCR(ctx, query)
}

// WaitForOCR espera a que terminen todos los jobs OCR pendientes.
// Útil en tests para evitar que la DB se cierre antes de que termine el OCR async.
func (s *DocumentoService) WaitForOCR() {
	s.ocrWg.Wait()
}

// processOCRAsync procesa OCR en background y actualiza el documento
func (s *DocumentoService) processOCRAsync(id int, filePath string) {
	ctx := context.Background()

	// Verificar Tesseract primero para evitar actualizar DB si no está disponible
	client, err := ocr.NewOCRClient()
	if err != nil {
		fmt.Printf("[OCR] tesseract no disponible para doc %d: %v\n", id, err)
		return
	}
	defer client.Close()

	// Actualizar estado a "processing"
	if err := s.repo.UpdateOCRStatus(ctx, id, "processing", ""); err != nil {
		fmt.Printf("[OCR] error marcando processing para doc %d: %v\n", id, err)
		return
	}

	text, err := client.ExtractText(filePath)
	if err != nil {
		s.repo.UpdateOCRStatus(ctx, id, "failed", "")
		fmt.Printf("[OCR] error extrayendo texto para doc %d: %v\n", id, err)
		return
	}

	// Limpiar y truncar
	text = CleanOCRText(text)
	if len(text) > 10000 {
		text = text[:10000] + "... [truncado]"
	}

	// Guardar resultado
	if err := s.repo.UpdateOCRStatus(ctx, id, "completed", text); err != nil {
		fmt.Printf("[OCR] error guardando resultado para doc %d: %v\n", id, err)
		return
	}

	fmt.Printf("[OCR] completado para doc %d (%d chars)\n", id, len(text))
}

func shouldRunOCR(documento *models.Documento) bool {
	if os.Getenv("OCR_ENABLED") == "false" {
		return false
	}
	switch documento.MimeType {
	case "image/jpeg", "image/png", "image/gif", "application/pdf":
		return true
	default:
		return false
	}
}
