package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// DB encapsula la conexión a SQLite
type DB struct {
	Conn *sql.DB
	Path string
}

// NewDB crea una nueva conexión SQLite
func NewDB(dbPath string) (*DB, error) {
	// Crear directorio si no existe
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("crear directorio db: %w", err)
	}

	// Abrir conexión
	conn, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=ON")
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite: %w", err)
	}

	// Verificar conexión
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	return &DB{
		Conn: conn,
		Path: dbPath,
	}, nil
}

// Close cierra la conexión
func (db *DB) Close() error {
	return db.Conn.Close()
}

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
