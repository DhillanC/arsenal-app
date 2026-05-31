-- Migration 002: Audit log table for CRUD operations tracking
-- Pre-requisito legal para manejo de manifiestos DIAN

CREATE TABLE IF NOT EXISTS audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ts DATETIME DEFAULT CURRENT_TIMESTAMP,
    action TEXT NOT NULL,      -- CREATE, UPDATE, DELETE, VIEW
    entity TEXT NOT NULL,      -- replica, documento, mantenimiento, actividad
    entity_id INTEGER,
    user_id INTEGER,           -- NULL hasta que exista auth (Fase 6)
    details_json TEXT,         -- JSON con cambios relevantes
    ip_address TEXT,
    user_agent TEXT
);

-- Índices para consultas comunes de auditoría
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_log(entity, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_ts ON audit_log(ts);
CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_log(action);
