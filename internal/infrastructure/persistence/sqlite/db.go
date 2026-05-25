package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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
