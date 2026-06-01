# cmd/api

## ¿Qué es?

Punto de entrada de la aplicación. Contiene `main.go` que inicializa el servidor HTTP con Gin, configura la base de datos SQLite, y arranca todos los componentes.

## Responsabilidades

- Leer variables de entorno (`APP_PORT`, `DB_PATH`, `UPLOAD_PATH`, `CORS_ALLOWED_ORIGINS`)
- Abrir conexión SQLite con WAL mode y `busy_timeout`
- Ejecutar migraciones embebidas (`go:embed`)
- Inicializar repositorios, servicios, y handlers
- Configurar router Gin con middleware (logging, CORS, recovery)
- Graceful shutdown con `signal.NotifyContext`

## Archivos clave

- `main.go` — Entry point, setup del container

## Relación con otras carpetas

- Usa `internal/domain/services` para lógica de negocio
- Usa `internal/infrastructure` para adaptadores (DB, storage, OCR)
- Sirve `web/` para templates estáticos y HTMX
