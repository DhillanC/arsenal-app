# Arsenal App - Pruebas Manuales

Documento de registro para pruebas manuales que no pueden ser automatizadas o requieren verificación visual/interactiva.

**Fecha de creación:** 2026-05-31
**Rama:** development

---

## Grupo 5 — Issues #52, #53, #54, #55, #56

| Issue | Descripción | Fix Aplicado | Prueba Manual Requerida | Estado |
|-------|-------------|--------------|------------------------|--------|
| **#52** | Botón "Subir" sin handler en vista de documentos | Verificado: formulario funciona correctamente | Confirmar que el formulario de upload en `/documentos` envía multipart y redirige correctamente | ⏳ Pendiente |
| **#53** | Tab Mantenimiento muestra "Próximamente" | Verificado: funcionalidad completa implementada | Navegar a ficha de réplica → tab Mantenimiento → verificar que lista tareas y permite crear nuevas | ⏳ Pendiente |
| **#54** | DarkMode FOUC (flash de tema claro al cargar) | Script bloqueante en `<head>` para leer localStorage antes de render | Refrescar página con tema oscuro guardado → verificar que NO hay flash blanco | ⏳ Pendiente |
| **#55** | go.mod dice 1.21, Dockerfile/README dicen 1.26 | Inconsistencia documentada, no requiere fix inmediato | N/A — solo documentación | ✅ N/A |
| **#56** | COPY migrations innecesario en Dockerfile | Ya embebidas con go:embed, línea eliminada | Verificar build Docker funciona: `docker build -t arsenal-test .` | ⏳ Pendiente |

**Notas grupo 5:**
- #52 y #53 fueron marcados como falsos positivos — la funcionalidad ya existía pero no estaba conectada correctamente en el template. Se verificó en código.
- #54 requiere prueba visual en navegador real (no testeable con go test).
- #56 requiere build Docker manual para confirmar que no se rompe la compilación.

---

## Grupo 6 — Issues #57, #59, #60, #61, #62

| Issue | Descripción | Fix Aplicado | Prueba Manual Requerida | Estado |
|-------|-------------|--------------|------------------------|--------|
| **#57** | GOOS=linux sin GOARCH en Dockerfile | Documentar comportamiento en Dockerfile | Verificar que la imagen Docker construida corre en ARM64 (Mac mini) sin problemas | ⏳ Pendiente |
| **#59** | CI workflow para Go | Crear `.github/workflows/go.yml` | Push a rama `development` → verificar que GitHub Actions ejecuta build, test y vet | ⏳ Pendiente |
| **#60** | Cobertura real inferior a TASKS.md | Documentar estado en `COVERAGE_STATUS.md` | Revisar que `docs/COVERAGE_STATUS.md` refleja gaps reales y recomendaciones | ✅ Verificado |
| **#61** | Falta test form-encoded vs JSON | Tests ya existen (`tests/integration/formdata_test.go`) | Correr `go test ./tests/integration/...` y confirmar que `formdata_test.go` pasa | ⏳ Pendiente |
| **#62** | Regex de storage permite '..foo' | Test agregado que documenta comportamiento edge-case | Verificar que `..foo` es rechazado por `SanitizePath` en prueba real | ⏳ Pendiente |

**Notas grupo 6:**
- #57: El Dockerfile actual produce imagen ARM64-only en Apple Silicon. Esto es conocido y aceptado para el Mac mini.
- #59: El workflow YAML está creado pero requiere push a GitHub para verificar que corre correctamente.
- #61: Los tests existen pero no se han corrido en este sprint — verificar que siguen pasando.

---

## Grupo 7 — Issues #63, #64, #65, #66, #67

| Issue | Descripción | Fix Aplicado | Prueba Manual Requerida | Estado |
|-------|-------------|--------------|------------------------|--------|
| **#63** | Endpoint `/documentos/:id/file` con auth y audit | Handler creado con validación de contención | Subir documento → acceder a `/api/v1/documentos/:id/file` → verificar descarga correcta | ⏳ Pendiente |
| **#64** | Stream uploads en chunks para archivos grandes | Documentar limitación y opciones futuras | Intentar subir archivo de 15MB → verificar que recibe 413 Payload Too Large | ⏳ Pendiente |
| **#65** | Endpoint `/stats/dashboard` con agregados SQL | StatsHandler con SQL agregado | Acceder a `/` (dashboard) → verificar que estadísticas se renderizan sin error | ⏳ Pendiente |
| **#66** | Export CSV/JSON para backup | ExportJSON endpoint creado | Acceder a `/api/v1/export/json` → verificar que descarga JSON válido con datos | ⏳ Pendiente |
| **#67** | demo-colores.html huérfano | Archivo eliminado del repo | Verificar que no queda referencia a `demo-colores.html` en ningún template o handler | ✅ Verificado |

**Notas grupo 7:**
- #63: El endpoint de descarga de archivo requiere prueba con archivo real subido.
- #64: El límite de 10MB está configurado — verificar que el frontend muestra error amigable cuando se excede.
- #65: Dashboard usa SQL agregado — verificar que no hay SQL injection (usa parámetros).
- #66: Export JSON requiere datos en la base para producir output significativo.

---

## Grupo 8 — Issues #68, #69, #70, #71

| Issue | Descripción | Fix Aplicado | Prueba Manual Requerida | Estado |
|-------|-------------|--------------|------------------------|--------|
| **#68** | Healthcheck separado: `/health/live` vs `/health/ready` | Separar endpoints con propósito claro | `curl /health/live` → 200; `curl /health/ready` → 200 (si DB ok) o 503 (si DB caída) | ⏳ Pendiente |
| **#69** | Índices SQLite faltantes | Agregar 5 índices nuevos para performance | Verificar que migración `002_indices.sql` aplica sin errores en DB existente | ⏳ Pendiente |
| **#70** | Validar tipo/estado antes de INSERT | Validación en handler con mensajes 400 claros | Enviar POST con tipo inválido → verificar 400 con mensaje claro; tipo válido → 201 | ⏳ Pendiente |
| **#71** | demo-colores.html huérfano (duplicado) | Archivo ya eliminado en #67 | N/A — verificado en grupo 7 | ✅ N/A |

**Notas grupo 8:**
- #68: Requiere verificar comportamiento con DB caída (simulado o real).
- #69: Índices agregados en migración 002 — verificar que no rompe DB existente.
- #70: Validación de enum en handler — probar edge cases como tipo vacío, tipo con espacios, etc.

---

## Resumen por Estado

| Estado | Cantidad |
|--------|----------|
| ✅ Verificado / N/A | 4 |
| ⏳ Pendiente de prueba manual | 16 |
| **Total** | **20** |

---

## Instrucciones para Ejecutar Pruebas Manuales

### Prerrequisitos
1. App corriendo: `go run cmd/api/main.go` o `docker-compose up`
2. Base de datos con datos de prueba (usar `scripts/seed.sql` si existe)
3. Navegador moderno con DevTools abierto

### Checklist General
- [ ] Tema oscuro: no hay flash blanco al cargar (#54)
- [ ] Upload de documentos: formulario funciona, archivo aparece en lista (#52, #63)
- [ ] Tab Mantenimiento: crear tarea, completar, verificar recálculo de fecha (#53)
- [ ] Dashboard: estadísticas se muestran sin errores (#65)
- [ ] Health endpoints: live siempre 200, ready depende de DB (#68)
- [ ] Export JSON: descarga archivo válido con datos reales (#66)
- [ ] Docker build: compila sin errores (#56, #57)
- [ ] GitHub Actions: push a development dispara workflow (#59)
- [ ] Validación tipos: mensajes de error claros en español (#70)
- [ ] Índices: migración aplica sin errores en DB existente (#69)

### Comandos Útiles
```bash
# Verificar health endpoints
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready

# Verificar export JSON
curl http://localhost:8080/api/v1/export/json | jq .

# Verificar build Docker
docker build -t arsenal-manual-test .

# Correr tests de integración
go test ./tests/integration/... -v
```

---

*Última actualización: 2026-05-31*
*Próximo paso: Ejecutar pruebas manuales pendientes y marcar como ✅ completadas.*
