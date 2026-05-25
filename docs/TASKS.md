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
- [ ] Handlers documentos
- [x] Entry point main.go

### Storage
- [x] Storage local (filesystem)

### Tests
- [x] Tests unitarios dominio
- [x] Tests integración repositorio
- [ ] Tests E2E API

---

## Fase 2: Core Features ⭐

### Seguridad (Nuevo)
- [ ] Análisis de amenazas (Threat Modeling)
- [ ] Diagramas de flujo de datos con trust boundaries
- [ ] Identificación de activos críticos (PHI, documentos DIAN)
- [ ] Controles de seguridad (auth, autorización, sanitización)
- [ ] Rate limiting y protección contra abuso
- [ ] Logging de auditoría (quién hizo qué y cuándo)
- [ ] Backup y recuperación de datos
- [ ] Encriptación de datos sensibles en reposo

### Documentos
- [ ] Subida de archivos (multipart)
- [ ] Storage local con organización por réplica
- [ ] OCR con Tesseract (gosseract)
- [ ] Metadatos extraídos en DB

### Actividades
- [ ] Timeline cronológico
- [ ] Filtros por tipo
- [ ] Búsqueda básica

### Validación
- [ ] Validación de campos (go-playground)
- [ ] Sanitización de inputs

---

## Fase 3: Mantenimiento & DIAN 🔧

### Mantenimiento Programado
- [ ] CRUD tareas de mantenimiento
- [ ] Cálculo de próximas fechas
- [ ] Alertas/recordatorios

### DIAN
- [ ] Campos específicos importación
- [ ] Búsqueda por número manifiesto
- [ ] Trazabilidad completa

---

## Fase 4: UI/UX 🎨

### Frontend
- [ ] Setup HTMX + Tailwind
- [ ] Templates HTML base
- [ ] Página lista réplicas
- [ ] Página ficha réplica
- [ ] Formularios

### PWA
- [ ] Manifest.json
- [ ] Service Worker
- [ ] Offline indicators

---

## Fase 5: Polish ✨

### Dashboard
- [ ] Estadísticas generales
- [ ] Gráficos de uso
- [ ] Costo total de propiedad

### Exportar
- [ ] Backup JSON
- [ ] Export CSV

### Deploy
- [ ] Documentación deploy Mac mini
- [ ] PM2 config
- [ ] Tailscale access

---

## Estado General

**Fase actual:** 1 - Foundation
**Progreso:** 15/25 tareas completadas (60%)

*Última actualización: 2026-05-25*