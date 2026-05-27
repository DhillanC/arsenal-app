package services

import (
	"context"
	"fmt"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// ReplicaService implementa inbound.ReplicaService
type ReplicaService struct {
	repo outbound.ReplicaRepository
}

// NewReplicaService crea un nuevo servicio
func NewReplicaService(repo outbound.ReplicaRepository) inbound.ReplicaService {
	return &ReplicaService{repo: repo}
}

// Create crea una nueva réplica
func (s *ReplicaService) Create(ctx context.Context, replica *models.Replica) error {
	if replica.Nombre == "" {
		return fmt.Errorf("nombre es requerido")
	}
	if replica.Estado == "" {
		replica.Estado = "activo"
	}
	return s.repo.Create(ctx, replica)
}

// GetByID obtiene una réplica por ID
func (s *ReplicaService) GetByID(ctx context.Context, id int) (*models.Replica, error) {
	return s.repo.GetByID(ctx, id)
}

// List lista todas las réplicas
func (s *ReplicaService) List(ctx context.Context) ([]models.Replica, error) {
	return s.repo.List(ctx)
}

// Update actualiza una réplica
func (s *ReplicaService) Update(ctx context.Context, replica *models.Replica) error {
	if replica.ID == 0 {
		return fmt.Errorf("id es requerido")
	}
	return s.repo.Update(ctx, replica)
}

// Delete elimina una réplica
func (s *ReplicaService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// Search busca réplicas por nombre, número de serie, marca o modelo (trazabilidad DIAN)
func (s *ReplicaService) Search(ctx context.Context, query string) ([]models.Replica, error) {
	if query == "" {
		return s.repo.List(ctx)
	}
	return s.repo.Search(ctx, query)
}
