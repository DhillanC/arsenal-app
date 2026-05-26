package ports

import (
	"context"
	"github.com/DhillanC/arsenal-app/internal/domain/models"
)

// ReplicaService define los casos de uso para réplicas
type ReplicaService interface {
	Create(ctx context.Context, replica *models.Replica) error
	GetByID(ctx context.Context, id int) (*models.Replica, error)
	List(ctx context.Context) ([]models.Replica, error)
	Update(ctx context.Context, replica *models.Replica) error
	Delete(ctx context.Context, id int) error
}

// ActividadService define los casos de uso para actividades
type ActividadService interface {
	Create(ctx context.Context, actividad *models.Actividad) error
	GetByID(ctx context.Context, id int) (*models.Actividad, error)
	ListByReplica(ctx context.Context, replicaID int) ([]models.Actividad, error)
	Update(ctx context.Context, actividad *models.Actividad) error
	Delete(ctx context.Context, id int) error
}

// DocumentoService define los casos de uso para documentos
type DocumentoService interface {
	Create(ctx context.Context, documento *models.Documento, file []byte) error
	GetByID(ctx context.Context, id int) (*models.Documento, error)
	ListByReplica(ctx context.Context, replicaID int) ([]models.Documento, error)
	ListByReplicaAndType(ctx context.Context, replicaID int, tipo string) ([]models.Documento, error)
	ListByActividad(ctx context.Context, actividadID int) ([]models.Documento, error)
	Update(ctx context.Context, documento *models.Documento) error
	Delete(ctx context.Context, id int) error
	SearchByOCR(ctx context.Context, query string) ([]models.Documento, error)
}