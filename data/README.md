# data

## ¿Qué es?

Directorio para datos persistentes de la aplicación: base de datos SQLite y archivos subidos.

## Estructura esperada

```
data/
├── arsenal.db          # Base de datos SQLite (WAL mode)
├── arsenal.db-shm      # Shared memory (WAL)
├── arsenal.db-wal      # Write-ahead log
└── uploads/            # Documentos subidos (organizados por réplica)
    └── {replica_id}/
        ├── factura_001.pdf
        ├── foto_001.jpg
        └── ...
```

## Seguridad

- **Permisos recomendados:**
  ```bash
  chmod 700 data/
  chmod 600 data/arsenal.db
  chmod 700 data/uploads/
  ```
- **Backup:** Copiar `arsenal.db` + `uploads/` para backup completo
- **No versionar:** Este directorio está en `.gitignore`

## Variables de entorno

- `DB_PATH` — Ruta al archivo SQLite (default: `./data/arsenal.db`)
- `UPLOAD_PATH` — Ruta a uploads (default: `./data/uploads`)
