# tests

## ¿Qué es?

Suite de tests automatizados. Incluye tests unitarios, de integración de repositorio, de integración HTTP, y de storage.

## Estructura

| Carpeta | Tipo de test |
|---------|-------------|
| `integration/` | Tests de integración HTTP (health, CORS, upload 413, path traversal) |
| `unit/` | Tests unitarios de dominio (entidades, validaciones) |

## Cómo correr

```bash
# Todos los tests
go test ./...

# Solo integración
go test ./tests/integration/...

# Con verbose
go test -v ./...
```

## Cobertura actual

- ✅ 13 funciones de test, todas PASS
- ✅ Health check (200/503)
- ✅ CORS (allow/block/preflight)
- ✅ Upload size limit (413)
- ✅ Path traversal defense (rechazo de `../`)
- ✅ CRUD de réplicas, actividades, documentos

## Gaps conocidos

- [ ] Tests de render HTML para rutas principales
- [ ] Tests E2E con navegador
- [ ] Tests de service layer (actualmente solo handlers + repos)
