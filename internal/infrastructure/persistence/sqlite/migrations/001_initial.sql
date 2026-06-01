CREATE TABLE IF NOT EXISTS replicas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    marca TEXT,
    modelo TEXT,
    tipo TEXT CHECK(tipo IN ('AEG', 'GBB', 'HPA', 'Spring', 'Otro')),
    numero_serie TEXT UNIQUE,
    fecha_adquisicion DATE,
    proveedor TEXT,
    costo_adquisicion REAL,
    estado TEXT DEFAULT 'activo' CHECK(estado IN ('activo', 'vendido', 'reparacion', 'prestado', 'archivado')),
    fps INTEGER,
    joules REAL,
    peso_gramos INTEGER,
    longitud_mm INTEGER,
    hop_up TEXT,
    capacidad_cargador INTEGER,
    notas TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS actividades (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    replica_id INTEGER NOT NULL,
    fecha DATE NOT NULL,
    tipo TEXT NOT NULL CHECK(tipo IN ('compra', 'venta', 'importacion', 'mantenimiento', 'reparacion', 'modificacion', 'uso', 'documentacion')),
    descripcion TEXT NOT NULL,
    proveedor_tecnico TEXT,
    costo REAL,
    kilometraje_bb INTEGER,
    ubicacion TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (replica_id) REFERENCES replicas(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS documentos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    replica_id INTEGER,
    actividad_id INTEGER,
    tipo TEXT NOT NULL CHECK(tipo IN ('factura', 'manual', 'manifiesto_dian', 'declaracion_dian', 'foto', 'video', 'otro')),
    nombre_archivo TEXT NOT NULL,
    ruta_archivo TEXT NOT NULL,
    mime_type TEXT,
    tamano_bytes INTEGER,
    ocr_texto TEXT,
    fecha_documento DATE,
    numero_documento TEXT,
    notas TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (replica_id) REFERENCES replicas(id) ON DELETE SET NULL,
    FOREIGN KEY (actividad_id) REFERENCES actividades(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS mantenimiento (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    replica_id INTEGER NOT NULL,
    tipo_tarea TEXT NOT NULL,
    frecuencia_dias INTEGER,
    frecuencia_bb INTEGER,
    ultima_fecha DATE,
    proxima_fecha DATE,
    completado BOOLEAN DEFAULT FALSE,
    notas TEXT,
    FOREIGN KEY (replica_id) REFERENCES replicas(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS piezas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    replica_id INTEGER,
    nombre TEXT NOT NULL,
    marca TEXT,
    tipo TEXT,
    instalada_en DATE,
    instalada_por TEXT,
    costo REAL,
    notas TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (replica_id) REFERENCES replicas(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sesiones (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fecha DATE NOT NULL,
    ubicacion TEXT,
    tipo_evento TEXT CHECK(tipo_evento IN ('practica', 'milsim', 'competencia', 'otro')),
    duracion_minutos INTEGER,
    notas TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS replica_sesion (
    replica_id INTEGER NOT NULL,
    sesion_id INTEGER NOT NULL,
    bb_disparadas INTEGER,
    PRIMARY KEY (replica_id, sesion_id),
    FOREIGN KEY (replica_id) REFERENCES replicas(id) ON DELETE CASCADE,
    FOREIGN KEY (sesion_id) REFERENCES sesiones(id) ON DELETE CASCADE
);

-- Índices para búsquedas comunes
CREATE INDEX IF NOT EXISTS idx_replicas_estado ON replicas(estado);
CREATE INDEX IF NOT EXISTS idx_replicas_marca ON replicas(marca);
CREATE INDEX IF NOT EXISTS idx_replicas_tipo ON replicas(tipo); -- Para stats/dashboard GROUP BY tipo
CREATE INDEX IF NOT EXISTS idx_replicas_estado_tipo ON replicas(estado, tipo); -- Para dashboard stats combinados
CREATE INDEX IF NOT EXISTS idx_actividades_replica ON actividades(replica_id);
CREATE INDEX IF NOT EXISTS idx_actividades_fecha ON actividades(fecha);
CREATE INDEX IF NOT EXISTS idx_actividades_replica_fecha ON actividades(replica_id, fecha DESC); -- Para timeline ordenado
CREATE INDEX IF NOT EXISTS idx_documentos_replica ON documentos(replica_id);
CREATE INDEX IF NOT EXISTS idx_documentos_actividad ON documentos(actividad_id); -- Para N+1 fix (ListByActividades)
CREATE INDEX IF NOT EXISTS idx_documentos_ocr ON documentos(ocr_texto); -- FTS en futuro
CREATE INDEX IF NOT EXISTS idx_documentos_tipo ON documentos(tipo); -- Para filtros por tipo
CREATE INDEX IF NOT EXISTS idx_mantenimiento_replica ON mantenimiento(replica_id);
CREATE INDEX IF NOT EXISTS idx_mantenimiento_proxima ON mantenimiento(proxima_fecha);