# Arsenal App - Análisis de Seguridad v1

**Fecha:** 2026-05-28  
**Versión analizada:** Post-Fase 5 (Mantenimiento + DIAN)  
**Estado:** Sin autenticación (endpoints públicos)  
**Scope:** Superficie de ataque actual antes de implementar JWT (Fase 6)

---

## 1. Resumen Ejecutivo

La aplicación Arsenal App (Go + SQLite + Gin) está en una fase **pre-autenticación**. Todos los endpoints son públicos — cualquiera con acceso de red puede crear, leer, modificar y eliminar réplicas, actividades, documentos y mantenimientos. Este análisis identifica los riesgos específicos de esta configuración y prioriza fixes antes de que la autenticación los enmascare o complique.

**Veredicto:** 🟡 **MEDIO-ALTO** — La app no está expuesta a internet (Tailscale/localhost), pero los vectores internos son reales. 5 hallazgos críticos/alto, 4 medios, 3 bajos.

---

## 2. Superficie de Ataque

### 2.1 Endpoints API (todos públicos)

| Método | Ruta | Riesgo sin auth |
|--------|------|-----------------|
| GET/POST/PUT/DELETE | `/api/v1/replicas` | CRUD completo sin restricción |
| GET/POST/PUT/DELETE | `/api/v1/replicas/:id/actividades` | CRUD completo |
| GET/POST | `/api/v1/replicas/:id/documentos` | **Upload de archivos sin auth** |
| GET | `/api/v1/documentos/filter` | Filtros de documentos |
| GET | `/api/v1/documentos/search` | Búsqueda OCR |
| GET/POST/PUT/DELETE | `/api/v1/replicas/:id/mantenimiento` | CRUD mantenimiento |
| POST | `/api/v1/mantenimiento/:id/completar` | Marcar completado |
| GET | `/api/v1/mantenimiento/proximos` | Listar próximos |
| GET | `/health` | Health check (legítimo) |

### 2.2 Endpoints HTML (frontend)

| Ruta | Riesgo |
|------|--------|
| `/` (dashboard) | Lectura pública de datos |
| `/replicas` | Lista completa |
| `/replicas/:id` | Ficha con documentos |
| `/replicas/nueva` | Formulario de creación |
| `/replicas/:id/editar` | Formulario de edición |
| `/documentos` | Lista documentos |
| `/mantenimiento` | Próximos mantenimientos |

### 2.3 Static Files

| Ruta | Riesgo |
|------|--------|
| `/static/*` | Assets públicos (legítimo) |
| `/uploads/*` | **Documentos DIAN/facturas servidos sin auth** |

---

## 3. Hallazgos por Severidad

### 🔴 CRÍTICO

#### H-001: Upload de archivos sin autenticación
**Descripción:** Cualquiera puede subir archivos a `/api/v1/replicas/:id/documentos` sin autenticarse.  
**Impacto:** Posible upload de malware, archivos de gran tamaño (aunque limitado a 10MB), o contaminación del filesystem.  
**Mitigación inmediata:**
- [ ] Implementar auth básica (HTTP Basic Auth) como medida temporal hasta JWT
- [ ] O: Restringir uploads a localhost/loopback en `main.go`
- [ ] Validar que el `replica_id` existe antes de permitir upload

**Evidencia:**
```go
// server.go — documentoRoutes sin middleware de auth
documentoRoutes.POST("", documentoHandler.Upload)  // PÚBLICO
```

#### H-002: Documentos servidos sin autenticación
**Descripción:** `/uploads/` expone todos los documentos subidos (facturas, manifiestos DIAN, fotos) sin ninguna verificación.  
**Impacto:** 🔴 **Fuga de documentos DIAN** — información legal de importación accesible por URL directa.  
**Mitigación inmediata:**
- [ ] Servir documentos a través de handler Go en vez de `router.Static`
- [ ] Agregar middleware de autenticación (incluso Basic Auth temporal)
- [ ] Agregar verificación de ownership del documento

**Evidencia:**
```go
// server.go
router.Static("/uploads", h.uploadPath)  // Cualquiera accede
```

---

### 🟠 ALTO

#### H-003: Soft-delete no implementado en réplicas
**Descripción:** `ReplicaHandler.Delete` llama `h.service.Delete(ctx, id)` que probablemente hace `DELETE FROM replicas`.  
**Impacto:** Pérdida permanente de datos sin trazabilidad. Para una app de trazabilidad legal (DIAN), esto es inaceptable.  
**Mitigación:**
- [ ] Implementar soft-delete: `UPDATE replicas SET deleted_at = NOW() WHERE id = ?`
- [ ] Filtrar `WHERE deleted_at IS NULL` en todas las queries
- [ ] Agregar campo `deleted_at` a la tabla

**Evidencia:**
```go
// replica_handler.go
c.JSON(http.StatusOK, gin.H{"message": "réplica eliminada"})  // Borrado físico
```

#### H-004: IDOR (Insecure Direct Object Reference) en todos los endpoints
**Descripción:** Cualquiera puede acceder a `/replicas/1`, `/replicas/2`, etc. sin verificar ownership.  
**Impacto:** En multi-user futuro, usuario A ve datos de usuario B. Ahora, con un solo usuario, es aceptable pero debe documentarse como deuda técnica.  
**Mitigación:**
- [ ] Documentar: "Fase 6 implementará user_id en todas las queries"
- [ ] Agregar `user_id` a JWT claims y filtrar en repositorios

#### H-005: Sin rate limiting
**Descripción:** Ningún límite de requests por IP o endpoint.  
**Impacto:** Vulnerable a brute force en búsquedas, scraping masivo, o DoS por volumen.  
**Mitigación:**
- [ ] Implementar rate limiting por IP (100 req/min) como middleware Gin
- [ ] Limitar más agresivamente en endpoints de upload (10 req/min)

---

### 🟡 MEDIO

#### H-006: Sin headers de seguridad HTTP
**Descripción:** No hay `X-Content-Type-Options`, `X-Frame-Options`, `Content-Security-Policy`, ni `Strict-Transport-Security`.  
**Impacto:** Clickjacking, MIME sniffing, XSS potencial.  
**Mitigación:**
- [ ] Agregar middleware de security headers

#### H-007: CORS permite cualquier origin en dev
**Descripción:** Si `CORS_ALLOWED_ORIGINS` está vacío, `CORSMiddleware` permite cualquier origin.  
**Impacto:** En producción sin configurar, un sitio malicioso puede hacer requests a la API.  
**Mitigación:**
- [ ] Rechazar requests con Origin no vacío cuando `allowedOrigins` está vacío en prod
- [ ] O: Requerir `CORS_ALLOWED_ORIGINS` en producción (fail closed)

#### H-008: Sin logging de auditoría
**Descripción:** No hay registro de quién hizo qué operación.  
**Impacto:** Imposible trazar modificaciones o borrados para propósitos legales.  
**Mitigación:**
- [ ] Implementar audit logging middleware (mínimo: timestamp, IP, método, ruta, body hash)

#### H-009: `documento_service.go` — Delete no elimina archivo físico
**Descripción:** `DocumentoService.Delete` solo borra el registro de DB, no el archivo del filesystem.  
**Impacto:** Archivos huérfanos acumulándose en `uploads/`. Posible fuga si el archivo contiene datos sensibles y la DB ya no lo referencia.  
**Mitigación:**
- [ ] Implementar eliminación física en `Delete`
- [ ] O: Job de cleanup periódico

---

### 🟢 BAJO

#### H-010: `html_handler.go` — errores silenciosos en carga de datos
**Descripción:** En `ReplicaDetail`, los errores de `actividadService.ListByReplica`, `documentoService.ListByReplica`, y `mantenimientoService.ListByReplica` se ignoran con `_`.  
**Impacto:** Página renderiza con datos incompletos sin notificar al usuario.  
**Mitigación:**
- [ ] Loggear errores (aunque no fallar la página)

#### H-011: Sin timeout en OCR
**Descripción:** `extractOCRText` no tiene timeout. Tesseract en una imagen grande puede colgar.  
**Impacto:** Goroutine bloqueada, posible acumulación.  
**Mitigación:**
- [ ] Agregar `context.WithTimeout` en OCR

#### H-012: `MaxMultipartMemory` no configurado en Gin
**Descripción:** `ParseMultipartForm(10<<20)` limita memoria a 10MB, pero `router.MaxMultipartMemory` (default Gin) puede ser diferente.  
**Impacto:** Inconsistencia en manejo de memoria vs disco para uploads.  
**Mitigación:**
- [ ] Explicitar `router.MaxMultipartMemory = 10 << 20` para claridad

---

## 4. Matriz de Riesgo

| Hallazgo | Severidad | Probabilidad | Riesgo | Prioridad |
|----------|-----------|--------------|--------|-----------|
| H-001 Upload sin auth | 🔴 Crítico | Alta | 🔴 Crítico | **P0** |
| H-002 Documentos públicos | 🔴 Crítico | Alta | 🔴 Crítico | **P0** |
| H-003 Sin soft-delete | 🟠 Alto | Media | 🟠 Alto | **P1** |
| H-004 IDOR | 🟠 Alto | Media | 🟠 Alto | **P1** (documentar para Fase 6) |
| H-005 Sin rate limiting | 🟠 Alto | Media | 🟠 Alto | **P1** |
| H-006 Sin security headers | 🟡 Medio | Media | 🟡 Medio | **P2** |
| H-007 CORS permissivo | 🟡 Medio | Baja (localhost) | 🟡 Medio | **P2** |
| H-008 Sin audit log | 🟡 Medio | Media | 🟡 Medio | **P2** |
| H-009 Delete sin cleanup | 🟡 Medio | Baja | 🟡 Medio | **P2** |
| H-010 Errores silenciosos | 🟢 Bajo | Media | 🟢 Bajo | **P3** |
| H-011 OCR sin timeout | 🟢 Bajo | Baja | 🟢 Bajo | **P3** |
| H-012 MaxMultipartMemory | 🟢 Bajo | Baja | 🟢 Bajo | **P3** |

---

## 5. Recomendaciones Pre-Fase 6 (Auth)

### Acciones inmediatas (antes de tocar JWT)

1. **H-001 + H-002 — Proteger uploads y documentos:**
   - Implementar middleware de autenticación básica (HTTP Basic Auth) como medida temporal
   - O: Agregar flag `REQUIRE_AUTH=true` que bloquee uploads y documentos
   - Servir `/uploads/` a través de handler con verificación

2. **H-003 — Soft-delete:**
   - Migración SQL para agregar `deleted_at` a `replicas`
   - Modificar `Delete` para soft-delete
   - Modificar `List`/`GetByID` para filtrar `deleted_at IS NULL`

3. **H-005 — Rate limiting:**
   - Middleware Gin con mapa de IPs + timestamps
   - Límite: 100 req/min general, 10 req/min uploads

### Acciones durante Fase 6

4. **H-004 — IDOR:** Implementar `user_id` en JWT claims y filtrar en repos
5. **H-008 — Audit logging:** Middleware que loguea todas las operaciones de escritura
6. **H-006 — Security headers:** Middleware con headers estándar

### Acciones post-Fase 6

7. **H-009 — Cleanup de archivos:** Implementar en `DocumentoService.Delete`
8. **H-007 — CORS hardening:** Requerir `CORS_ALLOWED_ORIGINS` en producción

---

## 6. Diagrama de Amenazas Actual (sin auth)

```mermaid
flowchart TD
    Attacker([Atacante]) -->|1. Upload malware| Upload[/api/v1/replicas/:id/documentos]
    Attacker -->|2. Leak docs| Static[/uploads/ — sin auth]
    Attacker -->|3. Delete todo| Delete[/api/v1/replicas/:id — DELETE]
    Attacker -->|4. Scrape data| List[/api/v1/replicas — GET]
    Attacker -->|5. DoS| DoS["Sin rate limiting\nrequests masivos"]

    Upload --> Storage[(Local Storage)]
    Static --> Storage
    Delete --> DB[(SQLite)]
    List --> DB
    DoS --> Server[HTTP Server]

    style Attacker fill:#ff6b6b,stroke:#c92a2a
    style Upload fill:#ffe3e3,stroke:#c92a2a
    style Static fill:#ffe3e3,stroke:#c92a2a
    style Delete fill:#ffe3e3,stroke:#c92a2a
```

---

## 7. Checklist de Verificación

- [ ] H-001: Upload requiere autenticación (o está bloqueado en red)
- [ ] H-002: `/uploads/` no es accesible públicamente
- [ ] H-003: Soft-delete implementado en réplicas
- [ ] H-004: Documentado en ADR-007 (IDOR mitigation plan)
- [ ] H-005: Rate limiting middleware activo
- [ ] H-006: Security headers presentes en todas las responses
- [ ] H-007: CORS no permite wildcard en producción
- [ ] H-008: Audit logging captura operaciones de escritura
- [ ] H-009: Delete de documento elimina archivo físico
- [ ] H-010: Errores en carga de datos secundarios son loggeados
- [ ] H-011: OCR tiene timeout configurado
- [ ] H-012: `MaxMultipartMemory` explicitado en código

---

## 8. Comparativa con Análisis STRIDE Original (Fase 2)

| Amenaza STRIDE | Estado Fase 2 | Estado v1 (ahora) | Gap |
|----------------|---------------|-------------------|-----|
| Spoofing (acceso no auth) | 🔴 Planeado (JWT Fase 6) | 🔴 **Aún sin auth** | Sin mitigación temporal |
| Tampering (modificación) | 🟠 Hash SHA-256 pendiente | 🟠 Sin soft-delete | H-003 |
| Repudiation | 🟠 Audit log pendiente | 🟠 Sin audit log | H-008 |
| Info Disclosure | 🟠 Path traversal fixed | 🔴 **Documentos públicos** | H-002 nuevo |
| DoS | 🟡 Rate limiting pendiente | 🟡 Sin rate limiting | H-005 |
| Elevation | 🟠 IDOR planeado Fase 6 | 🟠 IDOR aún presente | Aceptable pre-auth |

**Conclusión:** Los controles de Fase 2 fueron buenos (path traversal, upload cap, graceful shutdown), pero la falta de autenticación expone vectores que no estaban en el threat model original porque se asumía "single user local". La realidad es que la app corre en Docker con puertos expuestos.

---

*Análisis realizado: 2026-05-28*  
*Próxima revisión: Análisis de seguridad v2 (post-Fase 7, antes de release público)*
