package ports

import (
	"context"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
)

// ReplicaRepository define las operaciones de persistencia para réplicas
type ReplicaRepository interface {
	Create(ctx context.Context, replica *models.Replica) error
	GetByID(ctx context.Context, id int) (*models.Replica, error)
	List(ctx context.Context) ([]models.Replica, error)
	Update(ctx context.Context, replica *models.Replica) error
	Delete(ctx context.Context, id int) error
	// Search busca réplicas por número de serie (DIAN)
	Search(ctx context.Context, query string) ([]models.Replica, error)
}

// ActividadRepository define las operaciones de persistencia para actividades
type ActividadRepository interface {
	Create(ctx context.Context, actividad *models.Actividad) error
	GetByID(ctx context.Context, id int) (*models.Actividad, error)
	ListByReplica(ctx context.Context, replicaID int) ([]models.Actividad, error)
	Update(ctx context.Context, actividad *models.Actividad) error
	Delete(ctx context.Context, id int) error
}

// DocumentoRepository define las operaciones de persistencia para documentos
type DocumentoRepository interface {
	Create(ctx context.Context, documento *models.Documento) error
	GetByID(ctx context.Context, id int) (*models.Documento, error)
	ListByReplica(ctx context.Context, replicaID int) ([]models.Documento, error)
	ListByReplicaAndType(ctx context.Context, replicaID int, tipo string) ([]models.Documento, error)
	ListByActividad(ctx context.Context, actividadID int) ([]models.Documento, error)
	ListByActividades(ctx context.Context, actividadIDs []int) ([]models.Documento, error)
	Update(ctx context.Context, documento *models.Documento) error
	Delete(ctx context.Context, id int) error
	SearchByOCR(ctx context.Context, query string) ([]models.Documento, error)
}

// Storage define las operaciones de almacenamiento de archivos
type Storage interface {
	Save(file []byte, filename string, replicaID int) (string, error)
	Get(path string) ([]byte, error)
	Delete(path string) error
}

// OCR define las operaciones de reconocimiento óptico de caracteres
type OCR interface {
	ExtractText(imagePath string) (string, error)
	IsAvailable() bool
}

// MantenimientoRepository define las operaciones de persistencia para mantenimiento
type MantenimientoRepository interface {
	Create(ctx context.Context, mantenimiento *models.Mantenimiento) error
	GetByID(ctx context.Context, id int) (*models.Mantenimiento, error)
	ListByReplica(ctx context.Context, replicaID int) ([]models.Mantenimiento, error)
	ListProximos(ctx context.Context, dias int) ([]models.Mantenimiento, error)
	Update(ctx context.Context, mantenimiento *models.Mantenimiento) error
	Delete(ctx context.Context, id int) error
	MarcarCompletado(ctx context.Context, id int, fechaCompletado *time.Time) error
}
