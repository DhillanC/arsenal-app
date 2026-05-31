package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
)

// AuditLogRepository implementa outbound.AuditLogRepository con SQLite
type AuditLogRepository struct {
	db *sql.DB
}

// NewAuditLogRepository crea un nuevo repositorio de auditoría
func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create inserta un nuevo registro de auditoría
func (r *AuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_log (ts, action, entity, entity_id, user_id, details_json, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		log.Ts,
		log.Action,
		log.Entity,
		log.EntityID,
		log.UserID,
		log.DetailsJSON,
		log.IPAddress,
		log.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("insertar audit_log: %w", err)
	}
	return nil
}

// ListByEntity lista registros de auditoría por entidad
func (r *AuditLogRepository) ListByEntity(ctx context.Context, entity string, entityID int) ([]models.AuditLog, error) {
	query := `
		SELECT id, ts, action, entity, entity_id, user_id, details_json, ip_address, user_agent
		FROM audit_log
		WHERE entity = ? AND entity_id = ?
		ORDER BY ts DESC
	`
	rows, err := r.db.QueryContext(ctx, query, entity, entityID)
	if err != nil {
		return nil, fmt.Errorf("listar audit_log por entidad: %w", err)
	}
	defer rows.Close()

	return scanAuditLogs(rows)
}

// ListRecent lista los N registros más recientes
func (r *AuditLogRepository) ListRecent(ctx context.Context, limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	query := `
		SELECT id, ts, action, entity, entity_id, user_id, details_json, ip_address, user_agent
		FROM audit_log
		ORDER BY ts DESC
		LIMIT ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("listar audit_log recientes: %w", err)
	}
	defer rows.Close()

	return scanAuditLogs(rows)
}

func scanAuditLogs(rows *sql.Rows) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var userID sql.NullInt64
		var detailsJSON, ipAddress, userAgent sql.NullString
		var ts string

		err := rows.Scan(
			&log.ID,
			&ts,
			&log.Action,
			&log.Entity,
			&log.EntityID,
			&userID,
			&detailsJSON,
			&ipAddress,
			&userAgent,
		)
		if err != nil {
			return nil, fmt.Errorf("scan audit_log: %w", err)
		}

		// Parse timestamp
		log.Ts, _ = time.Parse("2006-01-02 15:04:05", ts)

		if userID.Valid {
			uid := int(userID.Int64)
			log.UserID = &uid
		}
		if detailsJSON.Valid {
			log.DetailsJSON = detailsJSON.String
		}
		if ipAddress.Valid {
			log.IPAddress = ipAddress.String
		}
		if userAgent.Valid {
			log.UserAgent = userAgent.String
		}

		logs = append(logs, log)
	}
	return logs, rows.Err()
}
