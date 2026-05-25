FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependencias de build
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copiar go.mod y go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

# Imagen final
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /root/

# Crear directorios para datos
RUN mkdir -p /data /uploads

# Copiar binario
COPY --from=builder /app/bin/api .

# Puerto
EXPOSE 8080

# Volumes para persistencia
VOLUME ["/data", "/uploads"]

CMD ["./api"]
