package models

import (
	"time"
)

// Replica representa una réplica airsoft en el inventario
type Replica struct {
	ID                int       `json:"id" db:"id"`
	Nombre            string    `json:"nombre" db:"nombre"`
	Marca             string    `json:"marca" db:"marca"`
	Modelo            string    `json:"modelo" db:"modelo"`
	Tipo              string    `json:"tipo" db:"tipo"` // AEG, GBB, HPA, Spring
	NumeroSerie       string    `json:"numero_serie" db:"numero_serie"`
	FechaAdquisicion  time.Time `json:"fecha_adquisicion" db:"fecha_adquisicion"`
	Proveedor         string    `json:"proveedor" db:"proveedor"`
	CostoAdquisicion  float64   `json:"costo_adquisicion" db:"costo_adquisicion"`
	Estado            string    `json:"estado" db:"estado"` // activo, vendido, reparacion, prestado
	FPS               int       `json:"fps" db:"fps"`
	Joules            float64   `json:"joules" db:"joules"`
	PesoGramos        int       `json:"peso_gramos" db:"peso_gramos"`
	LongitudMM        int       `json:"longitud_mm" db:"longitud_mm"`
	HopUp             string    `json:"hop_up" db:"hop_up"`
	CapacidadCargador int       `json:"capacidad_cargador" db:"capacidad_cargador"`
	Notas             string    `json:"notas" db:"notas"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Actividad representa un evento en la vida de una réplica
type Actividad struct {
	ID              int       `json:"id" db:"id"`
	ReplicaID       int       `json:"replica_id" db:"replica_id"`
	Fecha           time.Time `json:"fecha" db:"fecha"`
	Tipo            string    `json:"tipo" db:"tipo"` // compra, venta, mantenimiento, reparacion, modificacion, uso, importacion, documentacion
	Descripcion     string    `json:"descripcion" db:"descripcion"`
	ProveedorTecnico string   `json:"proveedor_tecnico" db:"proveedor_tecnico"`
	Costo           float64   `json:"costo" db:"costo"`
	KilometrajeBB   int       `json:"kilometraje_bb" db:"kilometraje_bb"`
	Ubicacion       string    `json:"ubicacion" db:"ubicacion"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Documento representa un archivo asociado a una réplica o actividad
type Documento struct {
	ID             int       `json:"id" db:"id"`
	ReplicaID      *int      `json:"replica_id,omitempty" db:"replica_id"`
	ActividadID    *int      `json:"actividad_id,omitempty" db:"actividad_id"`
	Tipo           string    `json:"tipo" db:"tipo"` // factura, manual, manifiesto_dian, declaracion_dian, foto, video, otro
	NombreArchivo  string    `json:"nombre_archivo" db:"nombre_archivo"`
	RutaArchivo    string    `json:"ruta_archivo" db:"ruta_archivo"`
	MimeType       string    `json:"mime_type" db:"mime_type"`
	TamanoBytes    int64     `json:"tamano_bytes" db:"tamano_bytes"`
	OCRTexto       string    `json:"ocr_texto,omitempty" db:"ocr_texto"`
	FechaDocumento *time.Time `json:"fecha_documento,omitempty" db:"fecha_documento"`
	NumeroDocumento string   `json:"numero_documento" db:"numero_documento"`
	Notas          string    `json:"notas" db:"notas"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Mantenimiento representa una tarea de mantenimiento programada
type Mantenimiento struct {
	ID             int        `json:"id" db:"id"`
	ReplicaID      int        `json:"replica_id" db:"replica_id"`
	TipoTarea      string     `json:"tipo_tarea" db:"tipo_tarea"`
	FrecuenciaDias int        `json:"frecuencia_dias" db:"frecuencia_dias"`
	UltimaFecha    *time.Time `json:"ultima_fecha,omitempty" db:"ultima_fecha"`
	ProximaFecha   *time.Time `json:"proxima_fecha,omitempty" db:"proxima_fecha"`
	Completado     bool       `json:"completado" db:"completado"`
	Notas          string     `json:"notas" db:"notas"`
}

// Tipos de actividad válidos
var TiposActividad = []string{
	"compra",
	"venta",
	"importacion",
	"mantenimiento",
	"reparacion",
	"modificacion",
	"uso",
	"documentacion",
}

// Tipos de documento válidos
var TiposDocumento = []string{
	"factura",
	"manual",
	"manifiesto_dian",
	"declaracion_dian",
	"foto",
	"video",
	"otro",
}