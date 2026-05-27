# Arsenal App - Feature Request Backfill

Este archivo controla el backfill de feature requests historicos para features ya implementados en Arsenal App.

## Workflow

1. Crear todos los feature requests como issues abiertos en GitHub.
2. Procesar uno por corrida de cron cada 6 minutos.
3. En cada corrida:
   - elegir el primer item con estado `open`.
   - comentar el issue con commit resolutorio, archivos principales y evidencia.
   - cerrar el issue.
   - actualizar este archivo a `closed`.
   - avisar por Telegram cual feature fue documentado/cerrado.

## Estados

- `pending_issue`: feature inventariado, issue aun no creado.
- `open`: issue creado y pendiente de documentar/cerrar.
- `closed`: issue documentado y cerrado.
- `local_uncommitted`: feature implementado en cambios locales, pendiente de commit resolutorio.

## Inventario

| ID | Estado | Issue | Feature | Commit resolutorio | Archivos principales | Notas |
|---|---|---:|---|---|---|---|
| FRB-001 | closed | #1 | Foundation de proyecto Go con Docker y estructura hexagonal | ec70a1a | `Dockerfile`, `docker-compose.yml`, `Makefile`, `cmd/api/main.go`, `internal/` | Base de proyecto, runtime local y estructura inicial. |
| FRB-002 | closed | #2 | Modelo de dominio y esquema SQLite inicial | 5195c00 | `internal/domain/models/`, `internal/domain/ports/`, `internal/infrastructure/persistence/sqlite/migrations/001_initial.sql` | Entidades Replica, Actividad, Documento, Mantenimiento, piezas y sesiones. |
| FRB-003 | closed | #3 | Repositorios, servicios y API REST core | 9fd086a | `internal/domain/services/`, `internal/infrastructure/persistence/sqlite/`, `internal/infrastructure/web/handlers/`, `internal/infrastructure/web/server.go` | CRUD base para replicas, actividades y documentos. |
| FRB-004 | closed | #4 | Storage local de documentos y pruebas de integracion | 5a0b779 | `internal/infrastructure/storage/local/storage.go`, `tests/` | Persistencia de archivos y cobertura inicial de storage/repositorios. |
| FRB-005 | closed | #5 | Threat model STRIDE y DFD de seguridad | 61bf76e | `docs/SECURITY.md`, `docs/MERMAID.md` | Documentacion de amenazas y boundaries. |
| FRB-006 | closed | #6 | Runtime Docker optimizado para ARM64 y migraciones en contenedor | 568ece9 | `Dockerfile`, `PROBLEMAS_DOCKER.md` | Ajustes para correr en Mac mini/ARM64. |
| FRB-007 | closed | #7 | Configuracion por entorno y puerto configurable | 1b7de4d | `cmd/api/main.go`, `internal/infrastructure/web/server.go`, `internal/infrastructure/persistence/sqlite/replica_repository.go` | `APP_PORT`, config central y ajustes de handlers. |
| FRB-008 | closed | #8 | Migraciones SQLite embebidas en el binario | f44d3ec | `internal/infrastructure/persistence/sqlite/migrations_embed.go`, `internal/infrastructure/persistence/sqlite/db.go`, `cmd/api/main.go` | Despliegue portable sin depender de archivos externos. |
| FRB-009 | closed | #9 | Healthcheck con ping de base de datos, graceful shutdown y runtime hardening | e2512c0 | `cmd/api/main.go`, `internal/infrastructure/persistence/sqlite/db.go`, `internal/infrastructure/web/server.go` | Timeouts HTTP, WAL/busy timeout y cierre ordenado. |
| FRB-010 | closed | #10 | Defensa contra path traversal y CORS configurable | 1bb04c9 | `internal/infrastructure/storage/local/storage.go`, `tests/storage_test.go`, `tests/integration/repository_test.go` | Controles de seguridad de rutas y origenes. |
| FRB-011 | closed | #11 | Límite real de upload y patron `run()` para cierre de DB | 6c34355 | `cmd/api/main.go`, `internal/infrastructure/web/handlers/documento_handler.go` | `http.MaxBytesReader` y cleanup fiable del proceso. |
| FRB-012 | closed | #12 | Upload multipart de documentos | 75c27fc | `internal/infrastructure/web/handlers/documento_handler.go`, `internal/infrastructure/web/server.go` | Subida de archivos asociada a replicas. |
| FRB-013 | closed | #13 | OCR con Tesseract para documentos de imagen | 3b142d2 | `internal/domain/services/ocr_service.go`, `internal/infrastructure/ocr/tesseract.go`, `internal/infrastructure/web/handlers/documento_handler.go` | Integracion inicial OCR. |
| FRB-014 | closed | #14 | Busqueda full-text por contenido OCR | 8ead052 | `internal/infrastructure/web/handlers/documento_handler.go`, `internal/infrastructure/web/server.go` | Endpoint de busqueda documental. |
| FRB-015 | closed | #15 | Filtros de documentos y timeline con documentos adjuntos | 9d1eee9 | `internal/domain/services/documento_service.go`, `internal/infrastructure/persistence/sqlite/documento_repository.go`, `internal/infrastructure/web/handlers/actividad_handler.go`, `internal/infrastructure/web/handlers/documento_handler.go` | Cierra Fase 3 funcional. |
| FRB-016 | closed | #16 | Frontend HTMX + Tailwind con tema DCS | 6f1eea5 | `web/templates/`, `web/static/manifest.json`, `internal/infrastructure/web/handlers/html_handler.go`, `internal/infrastructure/web/server.go` | Dashboard, lista, detalle, formularios, PWA base y dark/light mode. |
| FRB-017 | closed | #17 | Mantenimiento programado y trazabilidad DIAN | 24d8b84 | `internal/domain/services/mantenimiento_service.go`, `internal/infrastructure/persistence/sqlite/mantenimiento_repository.go`, `internal/infrastructure/web/handlers/mantenimiento_handler.go`, `internal/infrastructure/web/server.go` | CRUD mantenimiento, proximos mantenimientos, completar tareas y busqueda por serial. |
| FRB-018 | closed | #18 | Correccion de drift: ruta Search API y validacion BB-count M6 | ed131ab | `internal/infrastructure/web/server.go`, `internal/infrastructure/web/handlers/` | Ajustes funcionales detectados despues de Fase 5. |
| FRB-019 | closed | #19 | Correccion de drift y conexion real de OCR | 1d82e09 | `internal/infrastructure/ocr/tesseract.go`, `internal/infrastructure/web/handlers/documento_handler.go`, `docs/TASKS.md` | Alinea codigo con documentacion y activa flujo OCR real. |
| FRB-020 | closed | #20 | Workflow de planeacion en GitHub | 5ccb0fb | `.github/ISSUE_TEMPLATE/`, `.github/PULL_REQUEST_TEMPLATE.md`, `docs/PROJECT_BOARD.md`, `docs/EXAMPLES.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md` | Base de templates y flujo de tablero. |
| FRB-021 | closed | #21 | Actualizacion de templates de contribucion y workflow Markdown | b198713 | `.github/ISSUE_TEMPLATE/`, `.github/workflows/markdown.yml`, `.editorconfig`, `.markdownlint.json`, `CONTRIBUTING.md` | Refinamiento del flujo de issues/PR/documentacion. |
| FRB-022 | closed | #22 | Reparacion de frontend documentos/mantenimiento/HTMX | 2c22b23 | `docs/TASKS.md`, `internal/infrastructure/web/handlers/html_handler.go`, `internal/infrastructure/web/handlers/mantenimiento_handler.go`, `internal/infrastructure/web/handlers/replica_handler.go`, `internal/infrastructure/web/server.go`, `web/templates/` | Corrige vistas HTML de documentos y mantenimiento, formularios HTMX y estructura de templates. |
| FRB-023 | closed | #23 | Opciones de colores para campos de formulario | 2c22b23 | `docs/FIELD_COLOR_OPTIONS.md` | Documenta combinaciones de color para inputs en modo claro y oscuro; no cambia estilos todavia. |

## Progreso de cron

- Job: `0c3a0e82-fe86-406e-beb5-9f8766662bdd` (`arsenal-feature-request-backfill`) completado y eliminado el 2026-05-27 08:59 PDT.
- Cadencia: cada 6 minutos.
- Regla: procesar un solo issue por corrida.
- Ultimo item procesado: `FRB-023` / issue `#23` cerrado el 2026-05-27 08:47 PDT.
