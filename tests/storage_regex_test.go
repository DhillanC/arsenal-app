package local

import (
	"testing"

	"github.com/DhillanC/arsenal-app/internal/infrastructure/storage/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStorage_RegexRejectsDotDotFoo verifica que '..foo' sea rechazado
// (el regex por sí solo lo permitiría, pero sanitizeFilename debe rechazarlo).
// Issue #62: [TESTS] Regex de storage permite '..foo' — test pasa por suerte.
func TestStorage_RegexRejectsDotDotFoo(t *testing.T) {
	tempDir := t.TempDir()
	storage := local.NewStorage(tempDir)

	// ..foo contiene ".." y debe ser rechazado por sanitizeFilename
	_, err := storage.Save([]byte("x"), "..foo", 1)
	require.Error(t, err, "'..foo' debe ser rechazado (contiene '..')")
	assert.Contains(t, err.Error(), "path traversal")
}

// TestStorage_RegexRejectsDotDotSlash verifica que '../etc/passwd' sea rechazado.
func TestStorage_RegexRejectsDotDotSlash(t *testing.T) {
	tempDir := t.TempDir()
	storage := local.NewStorage(tempDir)

	_, err := storage.Save([]byte("x"), "../etc/passwd", 1)
	require.Error(t, err, "'../etc/passwd' debe ser rechazado")
	assert.Contains(t, err.Error(), "path traversal")
}

// TestStorage_RegexRejectsDotDotBackslash verifica que '..\windows\system32' sea rechazado.
func TestStorage_RegexRejectsDotDotBackslash(t *testing.T) {
	tempDir := t.TempDir()
	storage := local.NewStorage(tempDir)

	_, err := storage.Save([]byte("x"), "..\\windows\\system32", 1)
	require.Error(t, err, "'..\\windows\\system32' debe ser rechazado")
	assert.Contains(t, err.Error(), "path traversal")
}
