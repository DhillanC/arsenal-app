package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DB encapsula la conexión a SQLite
type DB struct {
	Conn *sql.DB
	Path string
}

// NewDB crea una nueva conexión SQLite.
// Para tests en memoria, usar ":memory:" como dbPath.
func NewDB(dbPath string) (*DB, error) {
	// Saltar MkdirAll para DBs en memoria (":memory:", "file::memory:?...").
	if !strings.HasPrefix(dbPath, ":") && !strings.Contains(dbPath, ":memory:") {
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("crear directorio db: %w", err)
		}
	}

	// busy_timeout=5000ms evita "database is locked" bajo concurrencia ligera.
	conn, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=ON&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite: %w", err)
	}

	// SQLite solo permite un escritor a la vez; serializar evita errores
	// de "database is locked" en cargas de escritura concurrentes.
	conn.SetMaxOpenConns(1)

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
