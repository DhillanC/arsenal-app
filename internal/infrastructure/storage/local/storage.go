package local

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	outbound "github.com/DhillanC/arsenal-app/internal/domain/ports/outbound"
)

// validFilenameRegex permite solo caracteres seguros
var validFilenameRegex = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// Storage implementa outbound.Storage usando filesystem local
type Storage struct {
	basePath string
}

// NewStorage crea un nuevo storage local
func NewStorage(basePath string) outbound.Storage {
	return &Storage{basePath: basePath}
}

// Save guarda un archivo en el filesystem con protección contra path traversal
func (s *Storage) Save(file []byte, filename string, replicaID int) (string, error) {
	// Sanitizar filename: solo nombre base, no paths
	filename = filepath.Base(filename)

	// filepath.Base("../../../etc/passwd") devuelve "passwd", que pasaría el regex
	// Pero strings.Contains con "/" detecta el intento de path traversal
	if strings.Contains(filename, "/") || strings.Contains(filename, string(os.PathSeparator)) || !validFilenameRegex.MatchString(filename) {
		return "", fmt.Errorf("nombre de archivo inválido: solo letras, números, puntos, guiones y guiones bajos")
	}

	// Crear estructura de directorios: uploads/replica_id/YYYY-MM/
	yearMonth := time.Now().Format("2006-01")
	dir := filepath.Join(s.basePath, strconv.Itoa(replicaID), yearMonth)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("crear directorio: %w", err)
	}

	// Defensa en profundidad: verificar que el path final está dentro de basePath
	path := filepath.Join(dir, filename)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolver path absoluto: %w", err)
	}
	absBase, err := filepath.Abs(s.basePath)
	if err != nil {
		return "", fmt.Errorf("resolver base path: %w", err)
	}
	rel, err := filepath.Rel(absBase, absPath)
	if err != nil || rel == ".." || rel[:3] == ".."+string(filepath.Separator) {
		return "", fmt.Errorf("path traversal detectado")
	}

	// Generar nombre único si ya existe (sufijo aleatorio, no timestamp)
	if _, err := os.Stat(path); err == nil {
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		suffix := make([]byte, 4)
		rand.Read(suffix)
		filename = fmt.Sprintf("%s_%s%s", name, hex.EncodeToString(suffix), ext)
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
