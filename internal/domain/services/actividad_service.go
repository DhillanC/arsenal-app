package services

import (
	"context"
	"fmt"

	"github.com/digital-consultory-solutions/arsenal-app/internal/domain/models"
	inbound "github.com/digital-consultory-solutions/arsenal-app/internal/domain/ports/inbound"
	outbound "github.com/digital-consultory-solutions/arsenal-app/internal/domain/ports/outbound"
)

// ActividadService implementa inbound.ActividadService
type ActividadService struct {
	repo outbound.ActividadRepository
}

// NewActividadService crea un nuevo servicio
func NewActividadService(repo outbound.ActividadRepository) inbound.ActividadService {
	return &ActividadService{repo: repo}
}

// Create crea una nueva actividad
func (s *ActividadService) Create(ctx context.Context, actividad *models.Actividad) error {
	if actividad.ReplicaID == 0 {
		return fmt.Errorf("replica_id es requerido")
	}
	if actividad.Tipo == "" {
		return fmt.Errorf("tipo es requerido")
	}
	if actividad.Descripcion == "" {
		return fmt.Errorf("descripcion es requerida")
	}
	return s.repo.Create(ctx, actividad)
}

// GetByID obtiene una actividad por ID
func (s *ActividadService) GetByID(ctx context.Context, id int) (*models.Actividad, error) {
	return s.repo.GetByID(ctx, id)
}

// ListByReplica lista actividades de una réplica
func (s *ActividadService) ListByReplica(ctx context.Context, replicaID int) ([]models.Actividad, error) {
	return s.repo.ListByReplica(ctx, replicaID)
}

// Update actualiza una actividad
func (s *ActividadService) Update(ctx context.Context, actividad *models.Actividad) error {
	if actividad.ID == 0 {
		return fmt.Errorf("id es requerido")
	}
	return s.repo.Update(ctx, actividad)
}

// Delete elimina una actividad
func (s *ActividadService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
