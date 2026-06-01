package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	"github.com/DhillanC/arsenal-app/internal/domain/services"
	"github.com/gin-gonic/gin"
)

// AuditConfig configura qué rutas y métodos auditar
type AuditConfig struct {
	// SkipPaths es un mapa de rutas a saltar (e.g., healthchecks)
	SkipPaths map[string]bool
	// SkipMethods es un mapa de métodos HTTP a saltar (e.g., GET para no auditar views)
	SkipMethods map[string]bool
}

// DefaultAuditConfig devuelve una configuración sensata por defecto
func DefaultAuditConfig() AuditConfig {
	return AuditConfig{
		SkipPaths: map[string]bool{
			"/health":       true,
			"/health/live":  true,
			"/health/ready": true,
			"/static":       true,
			"/uploads":      true,
		},
		SkipMethods: map[string]bool{
			"GET":    true, // No auditar lecturas por defecto (puede cambiarse)
			"HEAD":   true,
			"OPTIONS": true,
		},
	}
}

// AuditMiddleware crea un middleware Gin que registra operaciones CRUD
func AuditMiddleware(auditService *services.AuditLogService, config AuditConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si la ruta debe saltarse
		path := c.Request.URL.Path
		for skipPath := range config.SkipPaths {
			if len(path) >= len(skipPath) && path[:len(skipPath)] == skipPath {
				c.Next()
				return
			}
		}

		// Verificar si el método debe saltarse
		if config.SkipMethods[c.Request.Method] {
			c.Next()
			return
		}

		// Capturar body para extraer detalles (solo para POST/PUT/PATCH)
		var bodyBytes []byte
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
				// Restaurar body para que los handlers puedan leerlo
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Ejecutar el handler
		c.Next()

		// Solo auditar si la respuesta fue exitosa (2xx)
		if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
			return
		}

		// Determinar entidad y acción
		entity, action := detectEntityAndAction(path, c.Request.Method)
		if entity == "" {
			return // No es una entidad que auditemos
		}

		// Extraer entity_id de la URL si existe
		entityID := 0
		if id := c.Param("id"); id != "" {
			// Intentar convertir a int, si falla dejar en 0
			if parsed, err := strconv.Atoi(id); err == nil {
				entityID = parsed
			}
		}
		if mantenimientoID := c.Param("mantenimiento_id"); mantenimientoID != "" && entityID == 0 {
			if parsed, err := strconv.Atoi(mantenimientoID); err == nil {
				entityID = parsed
			}
		}

		// Construir detalles JSON
		details := make(map[string]any)
		if len(bodyBytes) > 0 {
			// Intentar parsear como JSON
			var bodyMap map[string]any
			if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
				// Filtrar campos sensibles
				for k, v := range bodyMap {
					if k != "password" && k != "token" && k != "secret" {
						details[k] = v
					}
				}
			}
		}

		detailsJSON := ""
		if len(details) > 0 {
			if b, err := json.Marshal(details); err == nil {
				detailsJSON = string(b)
			}
		}

		// Crear registro de auditoría
		log := &models.AuditLog{
			Ts:          time.Now(),
			Action:      action,
			Entity:      entity,
			EntityID:    entityID,
			DetailsJSON: detailsJSON,
			IPAddress:   c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
		}

		// Guardar de forma asíncrona (no bloquear la respuesta)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = auditService.Create(ctx, log)
		}()
	}
}

// detectEntityAndAction determina la entidad y acción basado en la ruta y método
func detectEntityAndAction(path, method string) (string, string) {
	var entity, action string

	// Detectar entidad
	switch {
	case contains(path, "/replicas"):
		entity = "replica"
	case contains(path, "/documentos"):
		entity = "documento"
	case contains(path, "/mantenimiento"):
		entity = "mantenimiento"
	case contains(path, "/actividades"):
		entity = "actividad"
	default:
		return "", ""
	}

	// Detectar acción
	switch method {
	case "POST":
		action = "CREATE"
	case "PUT", "PATCH":
		action = "UPDATE"
	case "DELETE":
		action = "DELETE"
	default:
		action = "VIEW"
	}

	return entity, action
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && s[:len(substr)] == substr || findSubstr(s, substr))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
