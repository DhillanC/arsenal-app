# docs

## ¿Qué es?

Documentación central del proyecto: plan de desarrollo, tareas por fase, análisis de seguridad, y guías de diagramas.

## Archivos principales

| Archivo | Contenido |
|---------|-----------|
| `PLAN.md` | Plan completo de desarrollo con fases, stack técnico, ADRs |
| `TASKS.md` | Checklist detallado de tareas por fase |
| `SECURITY.md` | Análisis de amenazas STRIDE, controles implementados, DFDs |
| `SECURITY_ANALYSIS_V1.md` | **Análisis v1** — revisión post-Fase 5, sin auth (12 hallazgos con priorización) |
| `MERMAID.md` | Guía de estilo para diagramas Mermaid (compatibilidad GitHub) |
| `EXAMPLES.md` | Ejemplos de issues, PRs, y user stories |
| `MVP_BITACORA.md` | Bitácora de decisiones durante el MVP |
| `PROJECT_BOARD.md` | Guía del board de proyecto |
| `FEATURE_REQUEST_BACKFILL.md` | Template para solicitudes de features |

## Documentación de Seguridad

| Documento | Descripción |
|-----------|-------------|
| `SECURITY.md` | Análisis de amenazas STRIDE, controles implementados, DFDs |
| `SECURITY_ANALYSIS_V1.md` | **Análisis v1** — revisión post-Fase 5, sin auth (12 hallazgos con priorización) |

### Proceso de seguridad

1. **Fase 2:** Threat model STRIDE inicial → 11 fixes aplicados
2. **Fase 5 (v1):** Revisión sin auth → 12 hallazgos, priorización P0-P3
3. **Fase 6:** Implementar auth → mitigar H-001, H-002, H-004, H-008
4. **Fase 7 (v2):** Revisión final post-MVP → antes de release público

## Diagramas incluidos

- Arquitectura hexagonal (flowchart)
- DFD con trust boundaries (flowchart)
- Secuencia de autenticación (sequenceDiagram)
- Flujo de subida de documento (flowchart)
- Flujo de mantenimiento programado (flowchart)
- Health check (flowchart)
- Request completo (sequenceDiagram)
- ER de base de datos (erDiagram)
- Diagrama de amenazas actual sin auth (flowchart)

## Convención

> Cada fase completada se marca con ✅ y se actualiza el estado en `PLAN.md`, `TASKS.md`, y `README.md`.
