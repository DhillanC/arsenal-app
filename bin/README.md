# bin

## ¿Qué es?

Directorio para binarios compilados y scripts auxiliares del proyecto.

## Uso típico

```bash
# Compilar el binario principal
go build -o bin/arsenal ./cmd/api

# Ejecutar
./bin/arsenal
```

## Nota

Este directorio está en `.gitignore` para no versionar binarios compilados. Se crea automáticamente al compilar.
