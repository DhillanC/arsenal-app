package ports

import (
	"context"
	"github.com/DhillanC/arsenal-app/internal/domain/models"
)

// ReplicaRepository define las operaciones de persistencia para réplicas
type ReplicaRepository interface {
	Create(ctx context.Context, replica *models.Replica) error
	GetByID(ctx context.Context, id int) (*models.Replica, error)
	List(ctx context.Context) ([]models.Replica, error)
	Update(ctx context.Context, replica *models.Replica) error
	Delete(ctx context.Context, id int) error
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
