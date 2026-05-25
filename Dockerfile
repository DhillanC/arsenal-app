# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Instalar dependencias de build (CGO para SQLite)
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

# Copiar solo go.mod y go.sum primero (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Build con optimizaciones
# -ldflags="-w -s" reduce tamaño del binario
# -trimpath remueve paths absolutos
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-w -s" \
    -o bin/api cmd/api/main.go

# Runtime stage
FROM alpine:3.19

# Crear usuario no-root para seguridad
RUN adduser -D -u 1000 arsenal && \
    apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

# Crear directorios con permisos correctos
RUN mkdir -p /data /uploads && \
    chown -R arsenal:arsenal /app /data /uploads

# Copiar binario desde builder
COPY --from=builder /app/bin/api /app/api

# Cambiar a usuario no-root
USER arsenal

# Puerto expuesto
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Volumes para persistencia
VOLUME ["/data", "/uploads"]

ENTRYPOINT ["/app/api"]
