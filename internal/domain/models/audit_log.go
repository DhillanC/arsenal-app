package models

import "time"

// AuditLog representa un registro de auditoría para operaciones CRUD
type AuditLog struct {
	ID         int       `json:"id"`
	Ts         time.Time `json:"ts"`
	Action     string    `json:"action"`     // CREATE, UPDATE, DELETE, VIEW
	Entity     string    `json:"entity"`     // replica, documento, mantenimiento, actividad
	EntityID   int       `json:"entity_id"`
	UserID     *int      `json:"user_id,omitempty"`
	DetailsJSON string   `json:"details_json,omitempty"`
	IPAddress  string    `json:"ip_address,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
}
