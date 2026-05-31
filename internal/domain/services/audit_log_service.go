package services

import (
	"context"
	"fmt"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// AuditLogService maneja la lógica de negocio para auditoría
type AuditLogService struct {
	repo outbound.AuditLogRepository
}

// NewAuditLogService crea un nuevo servicio de auditoría
func NewAuditLogService(repo outbound.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

// Create registra una nueva entrada de auditoría
func (s *AuditLogService) Create(ctx context.Context, log *models.AuditLog) error {
	if log.Action == "" {
		return fmt.Errorf("action es requerido")
	}
	if log.Entity == "" {
		return fmt.Errorf("entity es requerido")
	}
	return s.repo.Create(ctx, log)
}

// ListByEntity lista registros de auditoría por entidad
func (s *AuditLogService) ListByEntity(ctx context.Context, entity string, entityID int) ([]models.AuditLog, error) {
	return s.repo.ListByEntity(ctx, entity, entityID)
}

// ListRecent lista los N registros más recientes
func (s *AuditLogService) ListRecent(ctx context.Context, limit int) ([]models.AuditLog, error) {
	return s.repo.ListRecent(ctx, limit)
}
