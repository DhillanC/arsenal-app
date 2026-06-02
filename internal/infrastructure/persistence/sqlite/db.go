package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DB encapsula la conexión a SQLite con pools separados para lectura y escritura.
// WriteConn: MaxOpenConns=1 para serializar escrituras (WAL mode).
// ReadConn:  MaxOpenConns=4 para queries de lectura concurrentes.
type DB struct {
	WriteConn *sql.DB // Escrituras + transacciones
	ReadConn  *sql.DB // SELECTs exclusivamente
	Path      string
}

// NewDB crea una nueva conexión SQLite con pools separados.
// Para tests en memoria, usar ":memory:" como dbPath (solo WriteConn, ReadConn=WriteConn).
func NewDB(dbPath string) (*DB, error) {
	// Saltar MkdirAll para DBs en memoria.
	if !strings.HasPrefix(dbPath, ":") && !strings.Contains(dbPath, ":memory:") {
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("crear directorio db: %w", err)
		}
	}

	// Pool de escritura: serializado para evitar "database is locked" en WAL.
	writeDSN := dbPath + "?_journal_mode=WAL&_foreign_keys=ON&_busy_timeout=5000"
	writeConn, err := sql.Open("sqlite3", writeDSN)
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite escritura: %w", err)
	}
	writeConn.SetMaxOpenConns(1)
	if err := writeConn.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite escritura: %w", err)
	}

	// Para DB en memoria, reutilizar la misma conexión (no se puede abrir segunda conexión).
	var readConn *sql.DB
	if strings.HasPrefix(dbPath, ":") || strings.Contains(dbPath, ":memory:") {
		readConn = writeConn
	} else {
		// Pool de lectura: modo read-only con query_only para evitar escrituras accidentales.
		readDSN := dbPath + "?mode=ro&_query_only=on&_journal_mode=WAL&_busy_timeout=5000"
		readConn, err = sql.Open("sqlite3", readDSN)
		if err != nil {
			_ = writeConn.Close()
			return nil, fmt.Errorf("abrir sqlite lectura: %w", err)
		}
		readConn.SetMaxOpenConns(4)
		readConn.SetMaxIdleConns(2)
		if err := readConn.Ping(); err != nil {
			_ = writeConn.Close()
			_ = readConn.Close()
			return nil, fmt.Errorf("ping sqlite lectura: %w", err)
		}
	}

	return &DB{
		WriteConn: writeConn,
		ReadConn:  readConn,
		Path:      dbPath,
	}, nil
}

// Conn devuelve la conexión de escritura (compatibilidad hacia atrás).
// Deprecated: usar WriteConn o ReadConn explícitamente.
func (db *DB) Conn() *sql.DB {
	return db.WriteConn
}

// Close cierra ambas conexiones
func (db *DB) Close() error {
	var errs []error
	if err := db.WriteConn.Close(); err != nil {
		errs = append(errs, err)
	}
	// No cerrar dos veces si es la misma conexión (memoria).
	if db.ReadConn != db.WriteConn {
		if err := db.ReadConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("cerrar db: %v", errs)
	}
	return nil
}
