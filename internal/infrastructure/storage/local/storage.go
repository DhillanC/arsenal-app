package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// Storage implementa outbound.Storage usando filesystem local
type Storage struct {
	basePath string
}

// NewStorage crea un nuevo storage local
func NewStorage(basePath string) outbound.Storage {
	return &Storage{basePath: basePath}
}

// Save guarda un archivo en el filesystem
func (s *Storage) Save(file []byte, filename string, replicaID int) (string, error) {
	// Crear estructura de directorios: uploads/replica_id/YYYY-MM/filename
	yearMonth := time.Now().Format("2006-01")
	dir := filepath.Join(s.basePath, strconv.Itoa(replicaID), yearMonth)
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("crear directorio: %w", err)
	}

	// Generar nombre único si ya existe
	path := filepath.Join(dir, filename)
	if _, err := os.Stat(path); err == nil {
		// Archivo existe, agregar timestamp
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		filename = fmt.Sprintf("%s_%d%s", name, time.Now().Unix(), ext)
		path = filepath.Join(dir, filename)
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return "", fmt.Errorf("escribir archivo: %w", err)
	}

	return path, nil
}

// Get lee un archivo del filesystem
func (s *Storage) Get(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Delete elimina un archivo del filesystem
func (s *Storage) Delete(path string) error {
	return os.Remove(path)
}
