# Arsenal App - Bitacora MVP

Documento centralizado para entender que compone el MVP de Arsenal App, que features ya fueron implementados, con que commits quedaron resueltos, y que tareas reales faltan por cerrar.

Fuentes consolidadas:

- `docs/TASKS.md`
- `docs/FEATURE_REQUEST_BACKFILL.md`
- `docs/FIELD_COLOR_OPTIONS.md`
- Issues historicos `#1` a `#23` generados por el backfill de feature requests

Ultima actualizacion: 2026-05-27

## Estado ejecutivo

El MVP base de Arsenal App ya tiene foundation, backend core, documentos, OCR de imagenes, frontend HTMX, mantenimiento programado, trazabilidad DIAN y workflow de planeacion GitHub.

Estado actual:

- Rama activa: `development`
- Backfill historico de features: completado
- Issues de backfill: `#1` a `#23`, todos cerrados
- Cron de backfill: completado y eliminado
- Fase funcional actual: Fase 6 - Autenticacion y Seguridad API
- Ultima fase completada: Fase 5 - Mantenimiento y DIAN

## Que compone el MVP actual

### 1. Foundation tecnica

Incluye estructura de proyecto Go, Docker, Makefile, configuracion por entorno, SQLite, migraciones, arquitectura hexagonal, servicios, repositorios y API REST.

Features cerrados:

| Feature | Issue | Commit | Estado |
|---|---:|---|---|
| Foundation de proyecto Go con Docker y estructura hexagonal | #1 | `ec70a1a` | Cerrado |
| Modelo de dominio y esquema SQLite inicial | #2 | `5195c00` | Cerrado |
| Repositorios, servicios y API REST core | #3 | `9fd086a` | Cerrado |
| Storage local de documentos y pruebas de integracion | #4 | `5a0b779` | Cerrado |
| Runtime Docker optimizado para ARM64 y migraciones en contenedor | #6 | `568ece9` | Cerrado |
| Configuracion por entorno y puerto configurable | #7 | `1b7de4d` | Cerrado |
| Migraciones SQLite embebidas en el binario | #8 | `f44d3ec` | Cerrado |

Resultado:

- App portable para Mac mini y Docker.
- Dominio separado de infraestructura.
- SQLite local como base single-user.
- Configuracion por variables de entorno.
- Migraciones embebidas para despliegue mas confiable.

### 2. Seguridad y operacion base

Incluye threat model, DFD, hardening runtime, path traversal defense, CORS configurable, upload cap real y healthcheck con DB ping.

Features cerrados:

| Feature | Issue | Commit | Estado |
|---|---:|---|---|
| Threat model STRIDE y DFD de seguridad | #5 | `61bf76e` | Cerrado |
| Healthcheck con ping de base de datos, graceful shutdown y runtime hardening | #9 | `e2512c0` | Cerrado |
| Defensa contra path traversal y CORS configurable | #10 | `1bb04c9` | Cerrado |
| Limite real de upload y patron `run()` para cierre de DB | #11 | `6c34355` | Cerrado |

Resultado:

- La app ya tiene baseline de seguridad razonable para uso single-user/local.
- El runtime no depende de `log.Fatalf` para cerrar recursos.
- Los uploads tienen limite real con `http.MaxBytesReader`.
- Storage defiende contra path traversal.

### 3. Gestion de documentos

Incluye upload multipart, validacion MIME, storage por replica, OCR con Tesseract para imagenes, busqueda por OCR, filtros y timeline con documentos.

Features cerrados:

| Feature | Issue | Commit | Estado |
|---|---:|---|---|
| Upload multipart de documentos | #12 | `75c27fc` | Cerrado |
| OCR con Tesseract para documentos de imagen | #13 | `3b142d2` | Cerrado |
| Busqueda full-text por contenido OCR | #14 | `8ead052` | Cerrado |
| Filtros de documentos y timeline con documentos adjuntos | #15 | `9d1eee9` | Cerrado |
| Correccion de drift y conexion real de OCR | #19 | `1d82e09` | Cerrado |

Resultado:

- Los documentos pueden subirse y asociarse a replicas.
- Imagenes pueden pasar por OCR.
- El contenido OCR queda disponible para busqueda.
- Actividades pueden mostrar documentos adjuntos via API.

### 4. Frontend web y PWA

Incluye HTMX, Tailwind, tema DCS, dashboard, lista de replicas, ficha de replica, formularios, PWA manifest, service worker placeholder, rutas HTML de documentos y mantenimiento.

Features cerrados:

| Feature | Issue | Commit | Estado |
|---|---:|---|---|
| Frontend HTMX + Tailwind con tema DCS | #16 | `6f1eea5` | Cerrado |
| Correccion de drift: ruta Search API y validacion BB-count M6 | #18 | `ed131ab` | Cerrado |
| Reparacion de frontend documentos/mantenimiento/HTMX | #22 | `2c22b23` | Cerrado |
| Opciones de colores para campos de formulario | #23 | `2c22b23` | Cerrado |

Resultado:

- Dashboard inicial funcional.
- Rutas HTML principales conectadas.
- `/documentos` existe.
- `/mantenimiento` existe.
- Ficha de replica tiene documentos y mantenimiento conectados.
- Formularios HTMX fueron adaptados para form data.
- Colores de inputs estan documentados, pero aun no aplicados al CSS.

### 5. Mantenimiento y trazabilidad DIAN

Incluye CRUD de mantenimiento programado, proximos mantenimientos, completar tareas con recalculo, busqueda por numero de serie y tipos documentales DIAN.

Features cerrados:

| Feature | Issue | Commit | Estado |
|---|---:|---|---|
| Mantenimiento programado y trazabilidad DIAN | #17 | `24d8b84` | Cerrado |

Resultado:

- Mantenimientos por replica.
- Vista HTML para proximos mantenimientos.
- Creacion/listado desde ficha de replica.
- Completar mantenimiento desde UI.
- Busqueda por numero de serie para trazabilidad.

### 6. Workflow de planeacion GitHub

Incluye templates de issues, PR template, epic, task, user story, research spike, workflow Markdown, docs de tablero y reglas de contribucion.

Features cerrados:

| Feature | Issue | Commit | Estado |
|---|---:|---|---|
| Workflow de planeacion en GitHub | #20 | `5ccb0fb` | Cerrado |
| Actualizacion de templates de contribucion y workflow Markdown | #21 | `b198713` | Cerrado |

Resultado:

- El repo ya tiene estructura para planear trabajo futuro en GitHub.
- Hay templates para bugs, features, epics, tasks, user stories y research spikes.
- Markdown workflow existe para validar documentacion.

## Backfill historico

El backfill se ejecuto para reconstruir trazabilidad de features despues de haber implementado buena parte del MVP.

Resumen:

- Issues creados: `#1` a `#23`
- Issues cerrados: `#1` a `#23`
- Cron usado: `arsenal-feature-request-backfill`
- Job ID: `0c3a0e82-fe86-406e-beb5-9f8766662bdd`
- Cadencia usada: 6 minutos
- Ultimo item cerrado: `FRB-023` / issue `#23`
- Cron eliminado: 2026-05-27 08:59 PDT

Commits del cierre automatizado:

| Rango | Proposito |
|---|---|
| `528d4ff` a `ef89463` | Cierre progresivo de `FRB-004` a `FRB-023` |
| `d18d83a` | Marca el cron como completado y eliminado |

Nota: `FRB-001` a `FRB-003` tambien quedaron cerrados en GitHub y reflejados en `docs/FEATURE_REQUEST_BACKFILL.md`.

## Pendientes reales convertidos a feature candidates

Estos son los pendientes vivos que deberian convertirse en nuevos issues de GitHub cuando decidamos avanzar el backlog real.

### P0 - Usabilidad inmediata

| ID | Feature candidate | Estado | Fuente |
|---|---|---|---|
| MVP-P0-001 | Aplicar colores de inputs para modo claro y oscuro | Pendiente | `docs/FIELD_COLOR_OPTIONS.md` |
| MVP-P0-002 | Mejorar render de resultados de busqueda de documentos | Pendiente | `docs/TASKS.md` |
| MVP-P0-003 | Agregar tests de render HTML para rutas principales | Pendiente | `docs/TASKS.md` |

Detalle:

- El problema de inputs ya esta analizado.
- Recomendacion actual: opcion 1, neutro profesional.
- Rutas prioritarias para tests HTML: `/`, `/replicas`, `/replicas/nueva`, `/documentos`, `/mantenimiento`.

### P1 - Integridad de documentos y datos

| ID | Feature candidate | Estado | Fuente |
|---|---|---|---|
| MVP-P1-001 | Calcular y persistir hash SHA-256 de archivos subidos | Pendiente | `docs/TASKS.md` |
| MVP-P1-002 | Borrar archivo fisico al eliminar documento | Pendiente | `docs/TASKS.md` |
| MVP-P1-003 | Implementar backup y recuperacion de datos | Pendiente | `docs/TASKS.md` |
| MVP-P1-004 | Exportar backup JSON | Pendiente | `docs/TASKS.md` |
| MVP-P1-005 | Exportar CSV | Pendiente | `docs/TASKS.md` |
| MVP-P1-006 | Implementar OCR de PDF mediante conversion previa a imagen | Pendiente | `docs/TASKS.md` |

Detalle:

- `DocumentoService.Delete` todavia necesita eliminar tambien el archivo del storage.
- OCR PDF debe evitar introducir una dependencia pesada si `soffice`, `pdftoppm` o tooling local ya resuelve el paso de conversion.
- Backup/export deberia priorizar SQLite + uploads, no solo tablas.

### P2 - Observabilidad y hardening

| ID | Feature candidate | Estado | Fuente |
|---|---|---|---|
| MVP-P2-001 | Migrar logging a `log/slog` | Pendiente | `docs/TASKS.md` |
| MVP-P2-002 | Decidir alcance de cifrado en reposo | Pendiente | `docs/TASKS.md` |
| MVP-P2-003 | Agregar tests E2E API | Pendiente | `docs/TASKS.md` |

Detalle:

- `log/slog` debe entrar antes de audit logging para no duplicar mecanismos.
- Cifrado en reposo requiere decision: permisos locales + backup seguro, o cifrado aplicativo.
- E2E API deberia cubrir rutas criticas antes del release `v1.0.0`.

### P3 - Actividades, uso y mantenimiento

| ID | Feature candidate | Estado | Fuente |
|---|---|---|---|
| MVP-P3-001 | Filtros por tipo para actividades | Pendiente | `docs/TASKS.md` |
| MVP-P3-002 | Busqueda basica de actividades | Pendiente | `docs/TASKS.md` |
| MVP-P3-003 | Recordatorios locales de mantenimiento | Pendiente | `docs/TASKS.md` |
| MVP-P3-004 | UI y servicios para piezas y upgrades | Pendiente | `docs/TASKS.md` |
| MVP-P3-005 | UI y servicios para sesiones de campo | Pendiente | `docs/TASKS.md` |

Detalle:

- Las tablas de piezas, sesiones de campo y relacion replica-sesion ya existen en la migracion inicial.
- Falta exponer casos de uso, repositorios/servicios completos, handlers y UI.

### P4 - Autenticacion y seguridad API

| ID | Feature candidate | Estado | Fuente |
|---|---|---|---|
| MVP-P4-001 | Login y registro de usuarios | Pendiente | `docs/TASKS.md` |
| MVP-P4-002 | Middleware JWT | Pendiente | `docs/TASKS.md` |
| MVP-P4-003 | Password hashing con bcrypt | Pendiente | `docs/TASKS.md` |
| MVP-P4-004 | Rate limiting por IP | Pendiente | `docs/TASKS.md` |
| MVP-P4-005 | Rate limiting por usuario autenticado | Pendiente | `docs/TASKS.md` |
| MVP-P4-006 | Audit logging | Pendiente | `docs/TASKS.md` |

Detalle:

- Esta fase es mas relevante si Arsenal deja de ser single-user local.
- Para uso personal en Mac mini, puede quedar despues del polish operativo si no hay exposicion publica.

### P5 - Deploy y release

| ID | Feature candidate | Estado | Fuente |
|---|---|---|---|
| MVP-P5-001 | Documentacion deploy Mac mini | Pendiente | `docs/TASKS.md` |
| MVP-P5-002 | Configuracion de servicio para Mac mini | Pendiente | `docs/TASKS.md` |
| MVP-P5-003 | Tailscale access | Pendiente | `docs/TASKS.md` |
| MVP-P5-004 | GitHub Actions CI/CD | Pendiente | `docs/TASKS.md` |
| MVP-P5-005 | Release `v1.0.0` | Pendiente | `docs/TASKS.md` |

Detalle:

- Para el Mac mini, decidir entre LaunchAgent local y Docker Compose con restart policy.
- Tailscale debe ser la ruta preferida de acceso remoto antes de abrir puertos publicos.
- CI/CD minimo: tests Go, markdown lint y build Docker.

## Orden recomendado

1. `MVP-P0-001`: aplicar colores de inputs.
2. `MVP-P0-003`: tests HTML para rutas principales.
3. `MVP-P0-002`: render de busqueda de documentos.
4. `MVP-P1-002`: borrar archivo fisico al eliminar documento.
5. `MVP-P1-001`: hash SHA-256 de uploads.
6. `MVP-P1-003`: backup y recuperacion.
7. `MVP-P5-001` y `MVP-P5-002`: deploy Mac mini.
8. `MVP-P5-004`: GitHub Actions CI/CD.
9. `MVP-P5-005`: release `v1.0.0`.
10. Fase 6 auth/rate limit/audit, si el uso deja de ser estrictamente local/single-user.

## Criterio de MVP v1.0.0

Para marcar `v1.0.0`, el minimo razonable es:

- UI usable en claro y oscuro.
- Rutas HTML principales con tests.
- Upload/documentos sin leaks obvios de storage.
- Backup/export minimo.
- Deploy Mac mini documentado.
- CI de tests Go y markdown.
- Acceso remoto via Tailscale documentado.

Autenticacion puede quedar fuera de `v1.0.0` si la app corre solo en red privada/Tailscale para uso personal. Si se expone a terceros o a internet, autenticacion y rate limiting pasan a ser bloqueantes.
