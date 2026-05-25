# Arsenal App - Docker Build Guide

## Optimizaciones realizadas en Dockerfile

### 1. Multi-stage build
- **Builder stage:** Compila con todas las dependencias de build
- **Runtime stage:** Solo incluye el binario y librerías necesarias
- **Resultado:** Imagen final ~20MB vs ~300MB si fuera single-stage

### 2. Cache layers
- `go.mod` y `go.sum` se copian primero
- `go mod download` se ejecuta antes de copiar el código fuente
- **Resultado:** Si solo cambias código (no dependencias), el build usa cache

### 3. Optimizaciones de build
```dockerfile
# Reduce tamaño del binario
-ldflags="-w -s"

# Remueve paths absolutos (reproducible builds)
-trimpath

# Compila solo para Linux AMD64 (más rápido)
GOOS=linux GOARCH=amd64
```

### 4. Seguridad
- **Non-root user:** UID 1000 (`arsenal`)
- **No new privileges:** Previene escalación de privilegios
- **Read-only filesystem:** Posible con tmpfs para /tmp
- **Health checks:** Verifica que la app responde

### 5. docker-compose.yml mejorado
- **Resource limits:** CPU 1 core, RAM 512MB
- **Security opt:** `no-new-privileges:true`
- **Health check:** Más rápido (5s timeout vs 10s)
- **Named volumes:** Mejor gestión de persistencia

## Comandos

```bash
# Build rápido (aprovecha cache)
docker-compose build

# Build sin cache (desde cero)
docker-compose build --no-cache

# Levantar
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Escalar (si fuera multi-instance)
docker-compose up -d --scale api=2
```

## Troubleshooting

### Build lento primera vez
```bash
# Descarga base images primero
docker pull golang:1.26-alpine
docker pull alpine:3.19
```

### Imagen muy grande
```bash
# Ver capas
docker history arsenal-app

# Tamaño
docker images arsenal-app
```

### Permisos en volumes
```bash
# Fix ownership
sudo chown -R 1000:1000 ./data ./uploads
```

## Siguientes optimizaciones

- [ ] Distroless image (sin shell, más seguro)
- [ ] BuildKit para paralelización
- [ ] Cache mounts para `go build`
- [ ] Multi-arch (ARM64 para Mac)

---

*Documentación Docker - Arsenal App*
