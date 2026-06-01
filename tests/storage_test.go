package local_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DhillanC/arsenal-app/internal/infrastructure/storage/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	tempDir := t.TempDir()
	storage := local.NewStorage(tempDir)

	t.Run("Save", func(t *testing.T) {
		content := []byte("test content")
		path, err := storage.Save(content, "test.txt", 1)
		require.NoError(t, err)
		assert.Contains(t, path, "test.txt")
		// Ahora devuelve path relativo, verificar que existe bajo tempDir
		absPath := filepath.Join(tempDir, path)
		assert.FileExists(t, absPath)
	})

	t.Run("Get", func(t *testing.T) {
		content := []byte("hello world")
		path, err := storage.Save(content, "hello.txt", 1)
		require.NoError(t, err)

		data, err := storage.Get(path)
		require.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("Delete", func(t *testing.T) {
		content := []byte("to delete")
		path, err := storage.Save(content, "delete.txt", 1)
		require.NoError(t, err)
		absPath := filepath.Join(tempDir, path)
		assert.FileExists(t, absPath)

		err = storage.Delete(path)
		require.NoError(t, err)
		assert.NoFileExists(t, absPath)
	})

	t.Run("Duplicate filename", func(t *testing.T) {
		content := []byte("first")
		path1, err := storage.Save(content, "dup.txt", 1)
		require.NoError(t, err)

		content2 := []byte("second")
		path2, err := storage.Save(content2, "dup.txt", 1)
		require.NoError(t, err)

		assert.NotEqual(t, path1, path2)
		// Verificar que ambos existen bajo tempDir
		assert.FileExists(t, filepath.Join(tempDir, path1))
		assert.FileExists(t, filepath.Join(tempDir, path2))
	})
}

// TestStorage_SanitizeRejection verifica que el storage RECHAZA entradas
// peligrosas en lugar de neutralizarlas silenciosamente. Esto es importante:
// la versión anterior aceptaba "../../../etc/passwd" como "passwd" y lo
// guardaba bajo basePath. Funcional, pero opaca para auditoría.
func TestStorage_SanitizeRejection(t *testing.T) {
	tempDir := t.TempDir()
	storage := local.NewStorage(tempDir)

	cases := []struct {
		name     string
		filename string
		wantErr  string
	}{
		{"empty filename", "", "vacío"},
		{"forward slash traversal", "../../../etc/passwd", "separadores"},
		{"backslash traversal", "..\\..\\windows\\system32", "separadores"},
		{"null byte injection", "evil\x00.txt", "separadores"},
		{"only dots", "..", "relativos"},
		{"hidden traversal", "foo..bar", "relativos"},
		{"single dot", ".", "relativos"},
		{"non-ascii", "archivo_ñ.txt", "inválido"},
		{"space in name", "my file.txt", "inválido"},
		{"shell metachar", "a;rm -rf b.txt", "inválido"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := storage.Save([]byte("x"), tc.filename, 1)
			require.Error(t, err, "esperaba rechazo de %q", tc.filename)
			assert.Contains(t, err.Error(), tc.wantErr,
				"mensaje de error debe mencionar la causa, got: %v", err)
		})
	}

	// Independientemente del path lógico que el atacante intentó, /etc/passwd
	// no debe existir bajo tempDir y nada debe haberse creado fuera.
	outsidePath := filepath.Join(tempDir, "..", "..", "..", "etc", "passwd")
	_, statErr := os.Stat(outsidePath)
	assert.True(t, os.IsNotExist(statErr), "ningún archivo debe escapar de tempDir")
}

// TestStorage_AcceptsValidNames verifica que nombres legítimos pasan el filtro.
func TestStorage_AcceptsValidNames(t *testing.T) {
	tempDir := t.TempDir()
	storage := local.NewStorage(tempDir)

	valid := []string{
		"factura.pdf",
		"foto_2026-05.jpg",
		"A1B2C3.txt",
		"archivo-con-guiones.dat",
		"sin_extension",
		".dotfile",
	}

	for _, name := range valid {
		t.Run(name, func(t *testing.T) {
			path, err := storage.Save([]byte("x"), name, 1)
			require.NoError(t, err, "nombre válido rechazado: %q", name)
			// Ahora devuelve path relativo, verificar que existe bajo tempDir
			assert.FileExists(t, filepath.Join(tempDir, path),
				"archivo debe quedar bajo tempDir, got %q", path)
		})
	}
}
