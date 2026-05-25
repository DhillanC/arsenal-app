# Arsenal App - Arquitectura Hexagonal

App de gestión de réplicas airsoft - Inventario, mantenimiento, documentación DIAN y registro de uso.

## Desarrollo

- Rama principal: `development`
- Plan completo: [docs/PLAN.md](docs/PLAN.md)

## Stack

- Go 1.21+ (backend)
- SQLite (base de datos)
- HTMX + Alpine.js (frontend)
- Tailwind CSS (estilos)

## Estado

🚧 En fase de planificación y setup inicial.

---

## Diagrama de Arquitectura Hexagonal

```mermaid
flowchart TD
    subgraph Infraestructura["🟦 Infraestructura (Adaptadores)"]
        Web[Web HTTP<br/>Gin/Echo]
        DB[(SQLite<br/>Database)]
        Storage[Local Storage<br/>Filesystem]
        OCR[Tesseract<br/>OCR Engine]
    end

    subgraph Aplicacion["🟨 Aplicación (Casos de Uso)"]
        Commands[Commands<br/>Escritura]
        Queries[Queries<br/>Lectura]
        Handlers[Handlers<br/>Orquestación]
    end

    subgraph Dominio["🟥 Dominio (Núcleo)"]
        Replica[Replica<br/>Entity]
        Actividad[Actividad<br/>Entity]
        Documento[Documento<br/>Entity]
        Mantenimiento[Mantenimiento<br/>Entity]
        
        subgraph Puertos["Puertos (Interfaces)"]
            Inbound[Inbound Ports<br/>Service Interfaces]
            Outbound[Outbound Ports<br/>Repository Interfaces]
        end
    end

    Usuario([Usuario]) --> Web
    Web --> Commands
    Web --> Queries
    Commands --> Inbound
    Queries --> Inbound
    Inbound --> Replica
    Inbound --> Actividad
    Inbound --> Documento
    Inbound --> Mantenimiento
    
    Replica --> Outbound
    Actividad --> Outbound
    Documento --> Outbound
    Mantenimiento --> Outbound
    
    Outbound --> DB
    Outbound --> Storage
    Outbound --> OCR
```

## Estructura de Carpetas

```mermaid
flowchart LR
    subgraph Root["arsenal-app/"]
        CMD["cmd/api/<br/>main.go"]
        Internal["internal/"]
        Pkg["pkg/"]
        Scripts["scripts/"]
        Docs["docs/"]
        Tests["tests/"]
    end

    subgraph Domain["domain/"]
        Models["models/<br/>entities"]
        Ports["ports/<br/>interfaces"]
        Services["services/<br/>business logic"]
    end

    subgraph App["application/"]
        Commands["commands/"]
        Queries["queries/"]
        AppHandlers["handlers/"]
    end

    subgraph Infra["infrastructure/"]
        Persistence["persistence/<br/>sqlite"]
        StorageInfra["storage/<br/>local"]
        OCRInfra["ocr/<br/>tesseract"]
        WebInfra["web/<br/>server"]
    end

    Internal --> Domain
    Internal --> App
    Internal --> Infra
```

## Flujo de Datos: Subida de Documento

```mermaid
sequenceDiagram
    actor Usuario
    participant Web as Web Handler
    participant Cmd as Command
    participant Dom as Dominio
    participant Storage as Local Storage
    participant OCR as Tesseract OCR
    participant Repo as SQLite Repo

    Usuario->>Web: POST /documentos<br/>multipart/form-data
    Web->>Cmd: SubirDocumento(file)
    
    par Guardar archivo
        Cmd->>Storage: Guardar en filesystem
        Storage-->>Cmd: ruta_archivo
    and Procesar OCR
        Cmd->>OCR: Extraer texto
        OCR-->>Cmd: ocr_texto
    end
    
    Cmd->>Dom: Crear Documento
    Dom-->>Cmd: documento validado
    
    Cmd->>Repo: Insertar(documento)
    Repo-->>Cmd: id generado
    
    Cmd-->>Web: Documento creado
    Web-->>Usuario: 201 + metadatos
```

## Modelo de Datos

```mermaid
erDiagram
    REPLICA ||--o{ ACTIVIDAD : tiene
    REPLICA ||--o{ DOCUMENTO : posee
    REPLICA ||--o{ MANTENIMIENTO : requiere
    REPLICA ||--o{ PIEZA : contiene
    SESION ||--o{ REPLICA_USO : registra
    
    REPLICA {
        int id PK
        string nombre
        string marca
        string modelo
        string tipo
        string numero_serie
        date fecha_adquisicion
        string proveedor
        float costo_adquisicion
        string estado
        int fps
        float joules
        int peso_gramos
        int longitud_mm
        string hop_up
        int capacidad_cargador
        text notas
    }
    
    ACTIVIDAD {
        int id PK
        int replica_id FK
        date fecha
        string tipo
        text descripcion
        string proveedor_tecnico
        float costo
        int kilometraje_bb
        string ubicacion
    }
    
    DOCUMENTO {
        int id PK
        int replica_id FK
        int actividad_id FK
        string tipo
        string nombre_archivo
        string ruta_archivo
        string mime_type
        int tamano_bytes
        text ocr_texto
        date fecha_documento
        string numero_documento
    }
    
    MANTENIMIENTO {
        int id PK
        int replica_id FK
        string tipo_tarea
        int frecuencia_dias
        int frecuencia_bb
        date ultima_fecha
        date proxima_fecha
        boolean completado
    }
    
    PIEZA {
        int id PK
        int replica_id FK
        string nombre
        string marca
        string tipo
        date instalada_en
        string instalada_por
        float costo
    }
    
    SESION {
        int id PK
        date fecha
        string ubicacion
        string tipo_evento
        int duracion_minutos
    }
    
    REPLICA_USO {
        int replica_id FK
        int sesion_id FK
        int bb_disparadas
    }
```

---

*Repositorio privado - Digital Consultancy Solutions*