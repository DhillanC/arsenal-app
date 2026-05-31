package ports

import (
	"context"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
)

// AuditLogRepository define las operaciones de persistencia para auditoría
type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	ListByEntity(ctx context.Context, entity string, entityID int) ([]models.AuditLog, error)
	ListRecent(ctx context.Context, limit int) ([]models.AuditLog, error)
}
