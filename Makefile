.PHONY: dev build test migrate seed docker-build docker-up docker-down clean

# Desarrollo local
dev:
	go run cmd/api/main.go

# Build
build:
	go build -o bin/api cmd/api/main.go

# Tests
test:
	go test ./... -v

# Docker
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

# Limpieza
clean:
	rm -rf bin/ tmp/
	go clean -cache

# Instalar dependencias
deps:
	go mod download
	go mod tidy

# Lint
lint:
	golangci-lint run

# Formatear
fmt:
	go fmt ./...

# Todos los checks
check: fmt lint test
	@echo "All checks passed!"