# Arsenal App - Resumen de Problemas Docker

## 📋 Contexto del Proyecto

**Repositorio:** https://github.com/DhillanC/arsenal-app (público)  
**Rama:** `development`  
**Stack:** Go 1.26 + SQLite + Gin + Docker  
**Estado:** Fase 1 completa (API REST funcionando), Fase 2 en progreso (seguridad)

## ✅ Lo que funciona

1. **App local:** `go run cmd/api/main.go` levanta en http://localhost:8080
2. **API REST:** CRUD de réplicas y actividades funcionando
3. **Tests:** Pasando (repository + storage)
4. **Compilación:** `go build` genera binario de 16MB

## ❌ Problema: Docker Build Timeout

### Síntoma
```bash
docker build -t arsenal-app .
# Tarda más de 10 minutos sin output, eventualmente timeout
```

### Dockerfile actual (optimizado)
```dockerfile
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-w -s" -o bin/api cmd/api/main.go

FROM alpine:3.19
RUN adduser -D -u 1000 arsenal && apk add --no-cache ca-certificates sqlite-libs
COPY --from=builder /app/bin/api /app/api
USER arsenal
EXPOSE 8080
ENTRYPOINT ["/app/api"]
```

### docker-compose.yml
```yaml
services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
      - ./uploads:/uploads
```

## 🔍 Diagnóstico realizado

| Prueba | Resultado | Tiempo |
|--------|-----------|--------|
| `docker ps` | ✅ Funciona | - |
| `docker build` | ❌ Timeout | 10+ min |
| `docker pull golang:1.26-alpine` | ❌ Timeout | 10+ min |
| `go run cmd/api/main.go` | ✅ Funciona | 2 seg |
| `go build` | ✅ Funciona | 5 seg |

## 🤔 Hipótesis

1. **Red lenta:** Descarga de imágenes base tarda demasiado
2. **Docker Desktop:** Posible problema de recursos o configuración
3. **Mac mini:** Limitaciones de CPU/RAM durante build
4. **CGO + SQLite:** Compilación C puede ser lenta en Docker

## 📁 Estructura del proyecto

```
arsenal-app/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── domain/models/             # Entidades
│   ├── domain/services/            # Lógica de negocio
│   ├── infrastructure/
│   │   ├── persistence/sqlite/   # Repositorios
│   │   ├── storage/local/        # Filesystem
│   │   └── web/                  # Handlers + Server
├── tests/integration/            # Tests
├── Dockerfile                    # Multi-stage build
├── docker-compose.yml            # Compose config
└── docs/
    ├── DOCKER.md                 # Guía Docker
    ├── SECURITY.md               # Análisis de seguridad
    └── TASKS.md                  # Tareas por fase
```

## 🎯 Lo que necesitamos

1. **Docker build funcione** (idealmente en < 2 minutos)
2. **Docker compose up** levante la app completa
3. **Health check** responda desde contenedor
4. **Volumes** persistan datos entre reinicios

## 🔧 Comandos para probar

```bash
# Clonar y probar
git clone https://github.com/DhillanC/arsenal-app.git
cd arsenal-app
git checkout development

# Probar build
docker build -t arsenal-app .

# Probar compose
docker-compose up -d

# Verificar
curl http://localhost:8080/health
```

## 📊 Especificaciones del host

- **OS:** macOS (Darwin 24.6.0, arm64)
- **Docker:** v29.4.0
- **Go:** v1.26.3 (instalado recientemente)
- **Hardware:** Mac mini (Apple Silicon)

## 📝 Notas adicionales

- La app funciona perfectamente sin Docker (modo desarrollo)
- El Dockerfile ya fue optimizado (multi-stage, cache layers, non-root user)
- Posiblemente necesitemos ajustar recursos de Docker Desktop
- Alternativa: Podman o build local + deploy manual

---

**Contacto:** Dhillan Contreras - dhillancontreras@Mac-mini-de-Dhillan.local
**Fecha:** 2026-05-25
