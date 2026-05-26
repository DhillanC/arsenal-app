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

// validFilenameRegex permite solo caracteres seguros para nombres ya saneados.
var validFilenameRegex = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// Storage implementa outbound.Storage usando filesystem local
type Storage struct {
	basePath string
}

// NewStorage crea un nuevo storage local
func NewStorage(basePath string) outbound.Storage {
	return &Storage{basePath: basePath}
}

// sanitizeFilename rechaza intentos de path traversal en lugar de neutralizarlos.
//
// Política: si el input contiene cualquier separador de path, segmentos "..", o
// el byte nulo, devolvemos error. La detección debe ocurrir ANTES de filepath.Base
// — de lo contrario Base ya recortó la evidencia y el log de seguridad no ve el
// intento.
//
// Tras la validación, exigimos que el nombre coincida con el whitelist regex
// (solo letras, dígitos, ".", "_", "-").
func sanitizeFilename(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("nombre de archivo vacío")
	}
	if strings.ContainsAny(name, "/\\\x00") {
		return "", fmt.Errorf("path traversal detectado: separadores de path no permitidos")
	}
	if name == "." || name == ".." || strings.Contains(name, "..") {
		return "", fmt.Errorf("path traversal detectado: segmentos relativos no permitidos")
	}
	// filepath.Base es defensa en profundidad — a esta altura el nombre ya
	// no debería contener nada extraño, pero confirmamos.
	cleaned := filepath.Base(name)
	if cleaned != name {
		return "", fmt.Errorf("nombre de archivo inválido")
	}
	if !validFilenameRegex.MatchString(cleaned) {
		return "", fmt.Errorf("nombre de archivo inválido: solo letras, números, '.', '_', '-'")
	}
	return cleaned, nil
}

// Save guarda un archivo en el filesystem.
// El filename se valida estrictamente: cualquier intento de path traversal
// devuelve error en lugar de saneamiento silencioso.
func (s *Storage) Save(file []byte, filename string, replicaID int) (string, error) {
	safe, err := sanitizeFilename(filename)
	if err != nil {
		return "", err
	}

	yearMonth := time.Now().Format("2006-01")
	dir := filepath.Join(s.basePath, strconv.Itoa(replicaID), yearMonth)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("crear directorio: %w", err)
	}

	path := filepath.Join(dir, safe)

	// Defensa en profundidad: confirmar contención del path resuelto bajo basePath.
	absBase, err := filepath.Abs(s.basePath)
	if err != nil {
		return "", fmt.Errorf("resolver base path: %w", err)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolver path: %w", err)
	}
	rel, err := filepath.Rel(absBase, absPath)
	// strings.HasPrefix es seguro contra rel cortos (no panic como rel[:3]).
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal detectado: contención fallida")
	}

	// Si el archivo existe, usar sufijo único aleatorio. No usar time.Now() en
	// segundos — colisiona en uploads simultáneos.
	if _, err := os.Stat(path); err == nil {
		ext := filepath.Ext(safe)
		base := strings.TrimSuffix(safe, ext)
		suffix, err := randomHex(6)
		if err != nil {
			return "", fmt.Errorf("generar sufijo único: %w", err)
		}
		safe = fmt.Sprintf("%s_%s%s", base, suffix, ext)
		path = filepath.Join(dir, safe)
	}

	if err := os.WriteFile(path, file, 0o644); err != nil {
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

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("crypto/rand: %w", err)
	}
	return hex.EncodeToString(b), nil
}
