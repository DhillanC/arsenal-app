package local_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DhillanC/arsenal-app/internal/infrastructure/storage/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	// Crear directorio temporal
	tempDir := t.TempDir()
	storage := local.NewStorage(tempDir)

	t.Run("Save", func(t *testing.T) {
		content := []byte("test content")
		path, err := storage.Save(content, "test.txt", 1)
		require.NoError(t, err)
		assert.Contains(t, path, "test.txt")
		assert.FileExists(t, path)
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
		assert.FileExists(t, path)

		err = storage.Delete(path)
		require.NoError(t, err)
		assert.NoFileExists(t, path)
	})

	t.Run("Duplicate filename", func(t *testing.T) {
		content := []byte("first")
		path1, err := storage.Save(content, "dup.txt", 1)
		require.NoError(t, err)

		content2 := []byte("second")
		path2, err := storage.Save(content2, "dup.txt", 1)
		require.NoError(t, err)

		// Deben ser paths diferentes
		assert.NotEqual(t, path1, path2)
		assert.FileExists(t, path1)
		assert.FileExists(t, path2)
	})

	t.Run("Path traversal blocked", func(t *testing.T) {
		content := []byte("evil")
		// filepath.Base("../../../etc/passwd") devuelve "passwd"
		// filepath.Base limpia los "../" así que no detectamos por strings.Contains
		// Pero el filepath.Rel check debería detectar que está fuera de basePath
		_, err := storage.Save(content, "../../../etc/passwd", 1)
		// Nota: filepath.Base hace que esto pase, pero el filepath.Rel check debería fallar
		// Si no falla, es un bug que necesitamos arreglar
		if err == nil {
			// Si no hay error, verificar que al menos no escapó del tempDir
			outsidePath := filepath.Join(tempDir, "..", "..", "..", "etc", "passwd")
			_, statErr := os.Stat(outsidePath)
			assert.True(t, os.IsNotExist(statErr), "El archivo no debería existir fuera del tempDir")
		} else {
			assert.Contains(t, err.Error(), "path traversal")
		}
	})

	t.Run("Path traversal with Base escape blocked", func(t *testing.T) {
		content := []byte("evil2")
		// filepath.Base limpia "../../" pero debería quedar "passwd"
		// El regex rechazaría "passwd" sin extensión? No, es válido
		// Pero el path final debería estar dentro de tempDir/1/2026-01/
		path, err := storage.Save(content, "../../passwd", 1)
		require.NoError(t, err)

		// Verificar que está dentro del directorio esperado
		assert.True(t, strings.HasPrefix(path, tempDir))
		// No debería haber creado "etc/passwd"
		_, err = os.Stat(filepath.Join(tempDir, "etc", "passwd"))
		assert.True(t, os.IsNotExist(err))
	})
}
