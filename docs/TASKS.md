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
- [ ] Healthcheck endpoint

### Base de Datos
- [ ] Schema SQL inicial (migración 001)
- [ ] Conexión SQLite con WAL mode
- [ ] Runner de migraciones

### Dominio (Núcleo Hexagonal)
- [ ] Entidad Replica
- [ ] Entidad Actividad
- [ ] Entidad Documento
- [ ] Entidad Mantenimiento
- [ ] Puertos (interfaces repository)
- [ ] Puertos (interfaces service)

### Repositorio SQLite
- [ ] Implementar ReplicaRepository
- [ ] Implementar ActividadRepository
- [ ] Implementar DocumentoRepository

### API REST
- [ ] Setup Gin router
- [ ] Middleware (logging, CORS, recovery)
- [ ] Handlers réplicas (CRUD)
- [ ] Handlers actividades
- [ ] Handlers documentos

### Tests
- [ ] Tests unitarios dominio
- [ ] Tests integración repositorio
- [ ] Tests E2E API

---

## Fase 2: Core Features ⭐

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
**Progreso:** 5/25 tareas completadas (20%)

*Última actualización: 2026-05-25*