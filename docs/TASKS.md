# Arsenal App - Plan de Tareas

## Fase 1: Foundation 🏗️

### Docker & Project Setup
- [x] Crear estructura de carpetas hexagonal
- [x] Dockerfile multi-stage
- [x] docker-compose.yml
- [x] Makefile con comandos
- [x] .env.example
- [x] go.mod con dependencias

### Configuración
- [ ] Configurar Viper (.env, flags)
- [ ] Logger estructurado (zap/logrus)
- [x] Healthcheck endpoint

### Base de Datos
- [x] Schema SQL inicial (migración 001)
- [x] Conexión SQLite con WAL mode
- [x] Runner de migraciones

### Dominio (Núcleo Hexagonal)
- [x] Entidad Replica
- [x] Entidad Actividad
- [x] Entidad Documento
- [x] Entidad Mantenimiento
- [x] Puertos (interfaces repository)
- [x] Puertos (interfaces service)

### Repositorio SQLite
- [x] Implementar ReplicaRepository
- [x] Implementar ActividadRepository
- [x] Implementar DocumentoRepository

### Servicios (Aplicación)
- [x] ReplicaService
- [x] ActividadService
- [x] DocumentoService

### API REST
- [x] Setup Gin router
- [x] Middleware (logging, CORS, recovery)
- [x] Handlers réplicas (CRUD)
- [x] Handlers actividades
- [x] Handlers documentos (upload + list + search)
- [x] Entry point main.go

### Storage
- [x] Storage local (filesystem)

### Tests
- [x] Tests unitarios dominio
- [x] Tests integración repositorio
- [x] Tests integración storage
- [ ] Tests E2E API

---

## Fase 2: Core Ops + Seguridad ✅

### Seguridad (Completado)
- [x] Análisis de amenazas (Threat Modeling) - STRIDE completado
- [x] Diagramas de flujo de datos con trust boundaries
- [x] 11 fixes de seguridad aplicados:
  - [x] APP_PORT env var honored
  - [x] Soft-delete leak fixed
  - [x] time.Parse error handling
  - [x] Migrations embebidas (go:embed)
  - [x] Graceful shutdown
  - [x] SQLite busy_timeout + MaxOpenConns(1)
  - [x] Health check con DB ping
  - [x] Gin ReleaseMode en production
  - [x] Path traversal defense
  - [x] Docker target: builder eliminado
  - [x] CORS configurable
- [ ] Controles pendientes (post-auth):
  - [ ] File hash verification (SHA-256)
  - [ ] Max upload size limit (10MB)
- [ ] Backup y recuperación de datos
- [ ] Encriptación de datos sensibles en reposo

### Actividades
- [x] Timeline cronológico (API lista por réplica)
- [ ] Filtros por tipo
- [ ] Búsqueda básica

### Validación
- [x] Validación de campos (go-playground validator via Gin binding)
- [x] Sanitización de inputs (path traversal defense)

---

## Fase 3: Gestión de Documentos 📄

### Subida de Archivos
- [x] Handler multipart para documentos
- [x] Validación MIME type
- [x] Límite de tamaño (10MB)
- [x] Organización por réplica en filesystem

### OCR
- [x] Integración Tesseract (gosseract)
- [x] Extracción de texto en upload
- [x] Almacenar OCR en DB

### Búsqueda
- [x] Búsqueda full-text por contenido OCR
- [ ] Filtros por tipo de documento
- [ ] Timeline con documentos adjuntos

---

## Fase 4: Frontend Web 🎨

### HTMX + Tailwind
- [ ] Setup HTMX + Tailwind
- [ ] Templates HTML base
- [ ] Página lista réplicas
- [ ] Página ficha réplica
- [ ] Formularios

### Dashboard
- [ ] Estadísticas generales
- [ ] Gráficos de uso
- [ ] Costo total de propiedad

### PWA
- [ ] Manifest.json
- [ ] Service Worker
- [ ] Offline indicators

---

## Fase 5: Autenticación y Seguridad API 🔐

### JWT Authentication
- [ ] Login/registro de usuarios
- [ ] Middleware de auth
- [ ] Password hashing (bcrypt)

### Rate Limiting
- [ ] Límite por IP (100 req/min)
- [ ] Límite por usuario autenticado

### Audit Logging
- [ ] Quién hizo qué y cuándo
- [ ] Immutable audit store
- [ ] Queries de auditoría

---

## Fase 6: Mantenimiento & DIAN 🔧

### Mantenimiento Programado
- [ ] CRUD tareas de mantenimiento
- [ ] Cálculo de próximas fechas
- [ ] Alertas/recordatorios

### DIAN
- [ ] Campos específicos importación
- [ ] Búsqueda por número manifiesto
- [ ] Trazabilidad completa

---

## Fase 7: Polish ✨

### Exportar
- [ ] Backup JSON
- [ ] Export CSV

### Deploy
- [ ] Documentación deploy Mac mini
- [ ] PM2 config
- [ ] Tailscale access
- [ ] GitHub Actions CI/CD
- [ ] Release v1.0.0

---

## Estado General

**Fase actual:** 3 - Gestión de Documentos
**Progreso:** 18/35 tareas completadas (51%)

*Última actualización: 2026-05-25*