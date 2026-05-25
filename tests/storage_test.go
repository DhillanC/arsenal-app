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
}
