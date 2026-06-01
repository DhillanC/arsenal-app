# Cobertura de Tests — Estado Actual

**Fecha:** 2026-05-30
**Issue:** #60 [TESTS] Cobertura real muy inferior a lo declarado en TASKS.md

## Estado Actual

| Componente | Tests | Cobertura | Notas |
|-----------|-------|-----------|-------|
| **Replica Repository** | ✅ Create, GetByID, List, Update, Delete, Search | ~60% | Falta ListPaginated, Search con escape |
| **Actividad Repository** | ✅ Create, GetByID, ListByReplica | ~40% | Falta Update, Delete, Search |
| **Documento Repository** | ✅ Create, GetByID, ListByReplica | ~40% | Falta Update, Delete, ListByActividad, ListByActividades |
| **Mantenimiento Repository** | ❌ Ninguno | 0% | CRítico: Fase 5 declarada completa |
| **Replica Service** | ❌ Ninguno | 0% | Solo tests de integración HTTP |
| **Documento Service** | ❌ Ninguno | 0% | Solo tests de integración HTTP |
| **Mantenimiento Service** | ❌ Ninguno | 0% | Lógica de cálculo de próxima fecha sin test |
| **Storage** | ✅ Save, Get, Delete, Sanitización | ~70% | Buena cobertura |
| **HTTP Handlers** | ✅ Health, CORS, Replica Create/Get, Fecha inválida | ~30% | Falta upload, mantenimiento, timeline, search |
| **OCR** | ❌ Ninguno | 0% | Requiere tesseract instalado |

## Tests Agregados en Este Sprint

1. **formdata_test.go** — Verificar POST/PUT con form-data (Issue #26)
2. **documento_upload_test.go** — Upload de 1MB y rechazo de 11MB (Issue #28)

## Gaps Críticos

1. **Mantenimiento Service** — Sin tests unitarios (lógica de cálculo de próxima fecha)
2. **Documento Service** — Sin tests unitarios (solo integración HTTP)
3. **OCR** — Sin tests (requiere tesseract + pdftoppm instalados)
4. **Timeline** — Sin tests de integración (endpoint /timeline)
5. **Mantenimiento Handlers** — Sin tests de integración (create, completar, listar próximos)

## Recomendaciones

1. Priorizar tests de **Mantenimiento Service** (lógica de negocio crítica)
2. Agregar tests de **Documento Service** con mocks de storage
3. Agregar tests de **Timeline** endpoint (JSON y HTML)
4. Considerar tests de **OCR** con archivos de prueba pequeños

## TASKS.md Actualización Recomendada

- [ ] Tests unitarios dominio — **Parcial** (faltan service-level tests)
- [ ] Tests integración repositorio — **Parcial** (faltan mantenimiento, documento completo)
- [ ] Tests integración storage — ✅ **Completo**
- [ ] Tests E2E API — **Pendiente** (marcado correctamente)
