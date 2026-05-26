package services

import (
	"context"
	"fmt"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// MantenimientoService implementa inbound.MantenimientoService
type MantenimientoService struct {
	repo outbound.MantenimientoRepository
}

// NewMantenimientoService crea un nuevo servicio
func NewMantenimientoService(repo outbound.MantenimientoRepository) inbound.MantenimientoService {
	return &MantenimientoService{repo: repo}
}

// Create crea un nuevo registro de mantenimiento
func (s *MantenimientoService) Create(ctx context.Context, m *models.Mantenimiento) error {
	if m.TipoTarea == "" {
		return fmt.Errorf("tipo_tarea es requerido")
	}
	if m.ReplicaID == 0 {
		return fmt.Errorf("replica_id es requerido")
	}
	
	// Calcular próxima fecha si hay frecuencia y última fecha
	if m.FrecuenciaDias > 0 && m.UltimaFecha != nil {
		next := m.UltimaFecha.AddDate(0, 0, m.FrecuenciaDias)
		m.ProximaFecha = &next
	}
	
	return s.repo.Create(ctx, m)
}

// GetByID obtiene un mantenimiento por ID
func (s *MantenimientoService) GetByID(ctx context.Context, id int) (*models.Mantenimiento, error) {
	return s.repo.GetByID(ctx, id)
}

// ListByReplica lista mantenimientos de una réplica
func (s *MantenimientoService) ListByReplica(ctx context.Context, replicaID int) ([]models.Mantenimiento, error) {
	return s.repo.ListByReplica(ctx, replicaID)
}

// ListProximos lista mantenimientos próximos a vencer
func (s *MantenimientoService) ListProximos(ctx context.Context, dias int) ([]models.Mantenimiento, error) {
	if dias <= 0 {
		dias = 30 // default 30 días
	}
	return s.repo.ListProximos(ctx, dias)
}

// Update actualiza un mantenimiento
func (s *MantenimientoService) Update(ctx context.Context, m *models.Mantenimiento) error {
	if m.ID == 0 {
		return fmt.Errorf("id es requerido")
	}
	
	// Recalcular próxima fecha si cambió la frecuencia o última fecha
	if m.FrecuenciaDias > 0 && m.UltimaFecha != nil {
		next := m.UltimaFecha.AddDate(0, 0, m.FrecuenciaDias)
		m.ProximaFecha = &next
	}
	
	return s.repo.Update(ctx, m)
}

// Delete elimina un mantenimiento
func (s *MantenimientoService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// MarcarCompletado marca un mantenimiento como completado
func (s *MantenimientoService) MarcarCompletado(ctx context.Context, id int, fechaCompletado *time.Time) error {
	return s.repo.MarcarCompletado(ctx, id, fechaCompletado)
}
