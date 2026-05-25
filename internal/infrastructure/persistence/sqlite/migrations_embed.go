package sqlite

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations ejecuta las migraciones SQL embebidas
func (db *DB) RunMigrations() error {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("leer migraciones embebidas: %w", err)
	}

	// Ordenar por nombre
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		content, err := migrationsFS.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("leer migración %s: %w", entry.Name(), err)
		}

		if _, err := db.Conn.Exec(string(content)); err != nil {
			return fmt.Errorf("ejecutar migración %s: %w", entry.Name(), err)
		}
	}

	return nil
}
