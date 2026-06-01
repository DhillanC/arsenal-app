# .github

## ¿Qué es?

Configuración de GitHub: templates de PR/issue, workflows de CI/CD, y automatizaciones.

## Subcarpetas

| Carpeta | Contenido |
|---------|-----------|
| `workflows/` | GitHub Actions (CI/CD) — *pendiente implementar* |
| `ISSUE_TEMPLATE/` | Templates para bug reports y feature requests |

## Templates disponibles

- `PULL_REQUEST_TEMPLATE.md` — Checklist de PR
- `ISSUE_TEMPLATE/` — Formularios de bug report y feature request

## CI/CD Pendiente (Fase 7)

- [ ] Workflow de test automático (`go test ./...`)
- [ ] Workflow de build Docker
- [ ] Workflow de lint (`go vet`, `golangci-lint`)
- [ ] Release automático con tags
