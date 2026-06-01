-- Migration 005: Permitir 'skipped' en ocr_status para documentos donde OCR no aplica
-- (mime types no soportados o OCR_ENABLED=false). Antes se marcaba 'completed' con
-- ocr_texto vacío, lo cual mentía sobre el estado real y rompía dashboards que
-- filtran por completados.
--
-- SQLite no permite ALTER de un CHECK existente, así que reconstruimos la tabla.

PRAGMA foreign_keys = OFF;

CREATE TABLE documentos_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    replica_id INTEGER,
    actividad_id INTEGER,
    tipo TEXT NOT NULL CHECK(tipo IN ('factura', 'manual', 'manifiesto_dian', 'declaracion_dian', 'foto', 'video', 'otro')),
    nombre_archivo TEXT NOT NULL,
    ruta_archivo TEXT NOT NULL,
    mime_type TEXT,
    tamano_bytes INTEGER,
    ocr_texto TEXT,
    ocr_status TEXT DEFAULT 'pending' CHECK(ocr_status IN ('pending', 'processing', 'completed', 'failed', 'skipped')),
    fecha_documento DATE,
    numero_documento TEXT,
    notas TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (replica_id) REFERENCES replicas(id) ON DELETE SET NULL,
    FOREIGN KEY (actividad_id) REFERENCES actividades(id) ON DELETE SET NULL
);

INSERT INTO documentos_new (
    id, replica_id, actividad_id, tipo, nombre_archivo, ruta_archivo,
    mime_type, tamano_bytes, ocr_texto, ocr_status, fecha_documento,
    numero_documento, notas, created_at
)
SELECT
    id, replica_id, actividad_id, tipo, nombre_archivo, ruta_archivo,
    mime_type, tamano_bytes, ocr_texto, ocr_status, fecha_documento,
    numero_documento, notas, created_at
FROM documentos;

DROP TABLE documentos;
ALTER TABLE documentos_new RENAME TO documentos;

CREATE INDEX IF NOT EXISTS idx_documentos_replica ON documentos(replica_id);
CREATE INDEX IF NOT EXISTS idx_documentos_actividad ON documentos(actividad_id);
CREATE INDEX IF NOT EXISTS idx_documentos_ocr ON documentos(ocr_texto);
CREATE INDEX IF NOT EXISTS idx_documentos_tipo ON documentos(tipo);

PRAGMA foreign_keys = ON;
