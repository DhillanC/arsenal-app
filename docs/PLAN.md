# Arsenal App - Plan de Desarrollo

## Visión

Aplicación web/móvil (PWA) para gestión integral de réplicas airsoft. Inventario personal, trazabilidad legal (DIAN), mantenimiento técnico, registro de uso y documentación.

## Stack Técnico

| Capa | Tecnología | Justificación |
|------|-----------|---------------|
| **Frontend** | Next.js 14 (App Router) | SSR, PWA nativo, API routes integrados |
| **Estilos** | Tailwind CSS + shadcn/ui | Desarrollo rápido, componentes accesibles |
| **Backend** | Next.js API Routes + tRPC | Type-safe end-to-end, menos boilerplate |
| **Base de Datos** | SQLite | Ligero, zero-config, perfecto para uso personal/local |
| **ORM** | Drizzle ORM | Type-safe queries, migraciones SQL-first |
| **Auth** | NextAuth.js (Credentials) | Auth simple sin dependencias externas |
| **Storage** | Local filesystem + sharp | Fotos optimizadas, PDFs, videos |
| **OCR** | Tesseract.js (cliente) | Extrae texto de facturas y documentos DIAN |
| **Deploy** | Mac mini local + PM2 | Zero cost, acceso vía Tailscale/ngrok |

## Estructura del Proyecto

```
arsenal-app/
├── app/                          # Next.js App Router
│   ├── (dashboard)/              # Layout principal
│   │   ├── replicas/             # Lista de réplicas
│   │   ├── replicas/[id]/        # Ficha de réplica
│   │   ├── mantenimiento/        # Calendario de mantenimiento
│   │   ├── documentos/           # Gestión documental
│   │   └── estadisticas/         # Dashboard de uso
│   ├── api/                      # API Routes
│   │   ├── trpc/                 # tRPC router
│   │   └── auth/                 # NextAuth handlers
│   └── layout.tsx                # Root layout
├── components/                   # Componentes React
│   ├── ui/                       # shadcn/ui base
│   ├── replica-card.tsx          # Tarjeta de réplica
│   ├── timeline.tsx              # Timeline de actividades
│   └── document-uploader.tsx     # Subida con OCR
├── lib/                          # Utilidades
│   ├── db/                       # Drizzle schema + queries
│   ├── auth.ts                   # Configuración auth
│   └── ocr.ts                    # Wrapper Tesseract.js
├── public/                       # Assets estáticos
│   └── uploads/                  # Archivos subidos (gitignored)
├── types/                        # Tipos TypeScript globales
├── docs/                         # Documentación del proyecto
│   ├── adr/                      # Architecture Decision Records
│   └── wireframes/               # Mockups y diseños
├── tests/                        # Tests E2E y unitarios
└── scripts/                      # Scripts de utilidad
    └── seed.ts                   # Datos de ejemplo
```

## Modelo de Datos (SQLite)

### Tablas Principales

```sql
-- Réplicas
CREATE TABLE replicas (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  nombre TEXT NOT NULL,              -- "HK416 A5"
  marca TEXT,                        -- "VFC", "Tokyo Marui", etc.
  modelo TEXT,                       -- Modelo específico del fabricante
  tipo TEXT,                         -- "AEG", "GBB", "HPA", "Spring"
  numero_serie TEXT UNIQUE,          -- Serial del fabricante
  fecha_adquisicion DATE,
  proveedor TEXT,
  costo_adquisicion REAL,
  estado TEXT DEFAULT 'activo',      -- activo, vendido, reparacion, prestado
  fps INTEGER,                       -- Feet per second
  joules REAL,                       -- Energía
  peso_gramos INTEGER,
  longitud_mm INTEGER,
  hop_up TEXT,                       -- Tipo de hop-up
  capacidad_cargador INTEGER,
  notas TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Actividades (registro de uso, mantenimiento, etc.)
CREATE TABLE actividades (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  replica_id INTEGER NOT NULL REFERENCES replicas(id),
  fecha DATE NOT NULL,
  tipo TEXT NOT NULL,                -- compra, venta, mantenimiento, reparacion, modificacion, uso, importacion, documentacion
  descripcion TEXT NOT NULL,
  proveedor_tecnico TEXT,            -- Quién hizo el trabajo
  costo REAL,
  kilometraje_bb INTEGER,            -- BBs disparadas en esta actividad (si aplica)
  ubicacion TEXT,                    -- Campo de juego, taller, etc.
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Documentos (facturas, manuales, DIAN)
CREATE TABLE documentos (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  replica_id INTEGER REFERENCES replicas(id),
  actividad_id INTEGER REFERENCES actividades(id),
  tipo TEXT NOT NULL,                -- factura, manual, manifiesto_dian, declaracion_dian, foto, video, otro
  nombre_archivo TEXT NOT NULL,
  ruta_archivo TEXT NOT NULL,        -- Ruta relativa en public/uploads/
  mime_type TEXT,
  tamano_bytes INTEGER,
  ocr_texto TEXT,                    -- Texto extraído por OCR (si aplica)
  fecha_documento DATE,              -- Fecha que aparece en el documento
  numero_documento TEXT,             -- Número de factura, manifiesto, etc.
  notas TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Mantenimiento Programado
CREATE TABLE mantenimiento_programado (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  replica_id INTEGER NOT NULL REFERENCES replicas(id),
  tipo_tarea TEXT NOT NULL,          -- lubricacion, revision_compresion, cambio_orings, etc.
  frecuencia_dias INTEGER,           -- Cada cuántos días
  frecuencia_bb INTEGER,             -- O cada cuántas BBs
  ultima_fecha DATE,
  proxima_fecha DATE,
  completado BOOLEAN DEFAULT FALSE,
  notas TEXT
);

-- Piezas / Upgrades
CREATE TABLE piezas (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  replica_id INTEGER REFERENCES replicas(id),
  nombre TEXT NOT NULL,              -- "Piston SHS 14 teeth"
  marca TEXT,
  tipo TEXT,                         -- hop_up, piston, spring, barrel, motor, etc.
  instalada_en DATE,                 -- Cuándo se instaló
  instalada_por TEXT,                -- Quién la instaló
  costo REAL,
  notas TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Campos de Juego (log de uso)
CREATE TABLE sesiones_campo (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  fecha DATE NOT NULL,
  ubicacion TEXT,
  tipo_evento TEXT,                  -- practica, milsim, competencia
  duracion_minutos INTEGER,
  notas TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Relación réplicas-sesiones (qué réplicas se usaron)
CREATE TABLE replica_sesion (
  replica_id INTEGER REFERENCES replicas(id),
  sesion_id INTEGER REFERENCES sesiones_campo(id),
  bb_disparadas INTEGER,
  PRIMARY KEY (replica_id, sesion_id)
);
```

## Features por Fase

### Fase 1 - MVP (Semanas 1-2)
- [ ] Setup proyecto Next.js + SQLite + Drizzle
- [ ] CRUD de réplicas (lista, crear, editar, ficha)
- [ ] Subida de documentos básica (facturas, fotos)
- [ ] Registro de actividades simple (fecha, tipo, descripción)
- [ ] Timeline de actividades por réplica

### Fase 2 - Documentación DIAN (Semana 3)
- [ ] OCR de facturas y manifiestos
- [ ] Campos específicos para importaciones
- [ ] Búsqueda por número de manifiesto/serial
- [ ] Alertas de vencimiento (si aplica)

### Fase 3 - Mantenimiento (Semana 4)
- [ ] Calendario de mantenimiento
- [ ] Recordatorios (cron jobs locales)
- [ ] Checklist post-juego
- [ ] Historial de piezas cambiadas

### Fase 4 - Uso y Estadísticas (Semana 5)
- [ ] Log de sesiones de campo
- [ ] Estadísticas de uso por réplica
- [ ] Costo total de propiedad (TCO)
- [ ] Dashboard con gráficos

### Fase 5 - PWA y Mobile (Semana 6)
- [ ] Service Worker
- [ ] Instalable en iOS/Android
- [ ] Cámara nativa para fotos
- [ ] Offline-first (sync cuando hay conexión)

## Wireframes / UI Ideas

### Dashboard Principal
```
+--------------------------------------------------+
|  Arsenal App                        [+] Nueva    |
+--------------------------------------------------+
|                                                   |
|  MIS RÉPLICAS          PRÓXIMO MANTENIMIENTO      |
  +--------+--------+    +------------------------+ |
  | [foto] | [foto] |    | HK416 A5               | |
  | HK416  | M4A1   |    | Lubricación en 5 días  | |
  |   A5   |        |    +------------------------+ |
  +--------+--------+                               |
|                                                   |
|  ACTIVIDAD RECIENTE                               |
|  • 2026-05-23 - Compra HK416 A5 (Universal)      |
|  • 2026-05-24 - Desempaque + documentación       |
|                                                   |
+--------------------------------------------------+
```

### Ficha de Réplica
```
+--------------------------------------------------+
|  < HK416 A5                          [Editar]    |
+--------------------------------------------------+
|  [Foto principal]        Estado: Activo           |
|                          FPS: 380                 |
|  Serial: ABC123456       Joules: 1.2              |
|                                                   |
  [Timeline] [Docs] [Mant.] [Stats]                |
+--------------------------------------------------+
|  2026-05-23  ● Compra                             |
|              Universal de deportes SAS            |
|              [📄 Factura]                          |
|                                                   |
|  2026-05-24  ● Documentación                      |
|              Desempaque + fotos DIAN               |
|              [📷] [🎥]                            |
+--------------------------------------------------+
```

## Decisiones de Arquitectura (ADRs)

### ADR-001: SQLite sobre PostgreSQL
**Contexto:** App personal, un solo usuario, Mac mini como servidor.
**Decisión:** SQLite por simplicidad, zero-config, y portabilidad.
**Consecuencias:** No escala a multi-usuario fácilmente. Backup es copiar un archivo.

### ADR-002: Next.js full-stack sobre separación frontend/backend
**Contexto:** Desarrollo rápido, un solo desarrollador, deploy simple.
**Decisión:** Next.js con API routes + tRPC. Un solo proceso.
**Consecuencias:** Acoplamiento frontend/backend. Migrar a separación requiere refactor.

### ADR-003: Filesystem local sobre S3/Cloud
**Contexto:** Privacidad de documentos DIAN, costo cero, control total.
**Decisión:** Almacenamiento local en ~/Documents/arsenal-uploads/.
**Consecuencias:** Sin CDN. Backup manual o rsync.

### ADR-004: PWA sobre app nativa
**Contexto:** Un solo desarrollador, multiplataforma, sin App Store.
**Decisión:** PWA con Next.js. Instalable desde navegador.
**Consecuencias:** Acceso limitado a APIs nativas (cámara funciona, notificaciones push complicadas en iOS).

## Roadmap Extendido (Post-MVP)

- [ ] **Multi-usuario:** Auth real, roles (propietario, técnico, visualizador)
- [ ] **Marketplace interno:** Registro de ventas entre usuarios
- [ ] **Integración campos:** Lista de campos de juego en Colombia con reviews
- [ ] **Comunidad:** Compartir specs anónimamente para base de datos colaborativa
- [ ] **Importación desde GitHub:** Migrar datos del repo personal actual
- [ ] **Backup automático:** Sync a GitHub/GitLab como backup de documentos
- [ ] **App móvil nativa:** Expo/React Native si PWA no alcanza

## Notas de Desarrollo

- Usar `sharp` para optimización de imágenes al subir
- Videos: almacenar como están, mostrar con `<video>` tag
- OCR: procesar en cliente para no cargar servidor, guardar resultado en DB
- Tailscale para acceso remoto seguro desde cualquier lugar
- PM2 para mantener el proceso vivo en Mac mini

---

*Documento vivo - se actualiza conforme avanza el desarrollo*
