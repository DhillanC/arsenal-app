package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	"github.com/DhillanC/arsenal-app/internal/domain/services"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/persistence/sqlite"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/storage/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupServiceTest(t *testing.T) (inbound.ReplicaService, inbound.ActividadService, inbound.DocumentoService, inbound.MantenimientoService, *sqlite.DB) {
	t.Helper()
	db, err := sqlite.NewDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, db.RunMigrations())

	replicaRepo := sqlite.NewReplicaRepository(db)
	actividadRepo := sqlite.NewActividadRepository(db)
	documentoRepo := sqlite.NewDocumentoRepository(db)
	mantenimientoRepo := sqlite.NewMantenimientoRepository(db)
	storage := local.NewStorage(t.TempDir())

	replicaService := services.NewReplicaService(replicaRepo)
	actividadService := services.NewActividadService(actividadRepo)
	documentoService := services.NewDocumentoService(documentoRepo, storage)
	mantenimientoService := services.NewMantenimientoService(mantenimientoRepo)

	return replicaService, actividadService, documentoService, mantenimientoService, db
}

// createTestReplica crea una réplica de prueba y devuelve su ID
func createTestReplica(t *testing.T, svc inbound.ReplicaService) int {
	t.Helper()
	ctx := context.Background()
	replica := &models.Replica{
		Nombre:           fmt.Sprintf("Test-%d", time.Now().UnixNano()),
		Marca:            "VFC",
		Modelo:           "HK416",
		Tipo:             "AEG",
		Estado:           "activo",
		FechaAdquisicion: time.Now(),
		NumeroSerie:      fmt.Sprintf("TEST-%d", time.Now().UnixNano()),
	}
	require.NoError(t, svc.Create(ctx, replica))
	return replica.ID
}

// ========== ReplicaService Tests ==========

func TestReplicaService_Create(t *testing.T) {
	replicaService, _, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	replica := &models.Replica{
		Nombre:           "HK416 Test",
		Marca:            "VFC",
		Modelo:           "HK416 A5",
		Tipo:             "AEG",
		Estado:           "activo",
		FechaAdquisicion: time.Now(),
	}

	err := replicaService.Create(ctx, replica)
	require.NoError(t, err)
	assert.NotZero(t, replica.ID)
}

func TestReplicaService_Create_Validation(t *testing.T) {
	replicaService, _, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	// Sin nombre
	replica := &models.Replica{Tipo: "AEG", Estado: "activo"}
	err := replicaService.Create(ctx, replica)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nombre")
}

func TestReplicaService_GetByID(t *testing.T) {
	replicaService, _, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	// Crear
	replica := &models.Replica{
		Nombre: "M4A1 Test",
		Tipo:   "AEG",
		Estado: "activo",
	}
	require.NoError(t, replicaService.Create(ctx, replica))

	// Obtener
	found, err := replicaService.GetByID(ctx, replica.ID)
	require.NoError(t, err)
	assert.Equal(t, replica.Nombre, found.Nombre)
}

func TestReplicaService_GetByID_NotFound(t *testing.T) {
	replicaService, _, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	_, err := replicaService.GetByID(ctx, 99999)
	assert.Error(t, err)
}

func TestReplicaService_Update(t *testing.T) {
	replicaService, _, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	// Crear
	replica := &models.Replica{
		Nombre: "Original",
		Tipo:   "AEG",
		Estado: "activo",
	}
	require.NoError(t, replicaService.Create(ctx, replica))

	// Actualizar
	replica.Nombre = "Actualizado"
	replica.Estado = "reparacion"
	err := replicaService.Update(ctx, replica)
	require.NoError(t, err)

	// Verificar
	found, err := replicaService.GetByID(ctx, replica.ID)
	require.NoError(t, err)
	assert.Equal(t, "Actualizado", found.Nombre)
	assert.Equal(t, "reparacion", found.Estado)
}

func TestReplicaService_Delete(t *testing.T) {
	replicaService, _, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	// Crear
	replica := &models.Replica{
		Nombre:      "Para Borrar",
		Tipo:        "AEG",
		Estado:      "activo",
		NumeroSerie: fmt.Sprintf("DEL-%d", time.Now().UnixNano()),
	}
	require.NoError(t, replicaService.Create(ctx, replica))

	// Borrar
	err := replicaService.Delete(ctx, replica.ID)
	require.NoError(t, err)

	// Verificar que no existe (o está marcada como eliminada)
	found, err := replicaService.GetByID(ctx, replica.ID)
	// Si hay soft delete, puede que no dé error pero el estado cambie
	// o puede que dé error si es hard delete
	if err == nil {
		// Soft delete - verificar que el estado cambió
		assert.Equal(t, "archivado", found.Estado)
	} else {
		// Hard delete - error esperado
		assert.Error(t, err)
	}
}

func TestReplicaService_List(t *testing.T) {
	replicaService, _, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	// Crear varias
	for i := 0; i < 3; i++ {
		replica := &models.Replica{
			Nombre:      fmt.Sprintf("Replica %d", i),
			Tipo:        "AEG",
			Estado:      "activo",
			NumeroSerie: fmt.Sprintf("LIST-%d-%d", i, time.Now().UnixNano()),
		}
		require.NoError(t, replicaService.Create(ctx, replica))
	}

	// Listar
	list, err := replicaService.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3)
}

func TestReplicaService_Search(t *testing.T) {
	replicaService, _, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	// Crear
	replica := &models.Replica{
		Nombre:      "HK416 DMR",
		Marca:       "VFC",
		NumeroSerie: "SERIE123",
		Tipo:        "AEG",
		Estado:      "activo",
	}
	require.NoError(t, replicaService.Create(ctx, replica))

	// Buscar por nombre
	results, err := replicaService.Search(ctx, "HK416")
	require.NoError(t, err)
	assert.Len(t, results, 1)

	// Buscar por serie
	results, err = replicaService.Search(ctx, "SERIE123")
	require.NoError(t, err)
	assert.Len(t, results, 1)

	// Buscar sin resultados
	results, err = replicaService.Search(ctx, "NOEXISTE")
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

// ========== ActividadService Tests ==========

func TestActividadService_Create(t *testing.T) {
	replicaService, actividadService, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	act := &models.Actividad{
		ReplicaID:   replicaID,
		Tipo:        "mantenimiento",
		Descripcion: "Limpieza general",
		Fecha:       time.Now(),
		Costo:       50000,
	}

	err := actividadService.Create(ctx, act)
	require.NoError(t, err)
	assert.NotZero(t, act.ID)
}

func TestActividadService_Create_Validation(t *testing.T) {
	replicaService, actividadService, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	// Sin descripción
	act := &models.Actividad{ReplicaID: replicaID, Tipo: "mantenimiento"}
	err := actividadService.Create(ctx, act)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "descripci")
}

func TestActividadService_ListByReplica(t *testing.T) {
	replicaService, actividadService, _, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	// Crear actividades
	for i := 0; i < 3; i++ {
		act := &models.Actividad{
			ReplicaID:   replicaID,
			Tipo:        "mantenimiento",
			Descripcion: fmt.Sprintf("Actividad %d", i),
			Fecha:       time.Now(),
		}
		require.NoError(t, actividadService.Create(ctx, act))
	}

	// Listar
	list, err := actividadService.ListByReplica(ctx, replicaID)
	require.NoError(t, err)
	assert.Len(t, list, 3)
}

// ========== MantenimientoService Tests ==========

func TestMantenimientoService_Create(t *testing.T) {
	replicaService, _, _, mantenimientoService, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	m := &models.Mantenimiento{
		ReplicaID:      replicaID,
		TipoTarea:      "Lubricación",
		FrecuenciaDias: 30,
	}

	err := mantenimientoService.Create(ctx, m)
	require.NoError(t, err)
	assert.NotZero(t, m.ID)
}

func TestMantenimientoService_Create_CalculaProximaFecha(t *testing.T) {
	replicaService, _, _, mantenimientoService, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	ultimaFecha := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	m := &models.Mantenimiento{
		ReplicaID:      replicaID,
		TipoTarea:      "Lubricación",
		FrecuenciaDias: 30,
		UltimaFecha:    &ultimaFecha,
	}

	err := mantenimientoService.Create(ctx, m)
	require.NoError(t, err)
	require.NotNil(t, m.ProximaFecha)
	assert.Equal(t, 31, m.ProximaFecha.Day()) // 1 + 30 días = 31
}

func TestMantenimientoService_Create_Validation(t *testing.T) {
	_, _, _, mantenimientoService, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	// Sin tipo_tarea
	m := &models.Mantenimiento{ReplicaID: 1}
	err := mantenimientoService.Create(ctx, m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tipo_tarea")
}

func TestMantenimientoService_MarcarCompletado(t *testing.T) {
	replicaService, _, _, mantenimientoService, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	// Crear
	m := &models.Mantenimiento{
		ReplicaID:      replicaID,
		TipoTarea:      "Limpieza",
		FrecuenciaDias: 30,
	}
	require.NoError(t, mantenimientoService.Create(ctx, m))

	// Completar
	fechaCompletado := time.Now()
	err := mantenimientoService.MarcarCompletado(ctx, m.ID, &fechaCompletado)
	require.NoError(t, err)

	// Verificar
	found, err := mantenimientoService.GetByID(ctx, m.ID)
	require.NoError(t, err)
	assert.True(t, found.Completado)
	require.NotNil(t, found.ProximaFecha)
}

func TestMantenimientoService_ListProximos(t *testing.T) {
	replicaService, _, _, mantenimientoService, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	// Crear mantenimiento con próxima fecha cercana
	proximaFecha := time.Now().AddDate(0, 0, 5)
	m := &models.Mantenimiento{
		ReplicaID:      replicaID,
		TipoTarea:      "Lubricación",
		FrecuenciaDias: 30,
		ProximaFecha:   &proximaFecha,
	}
	require.NoError(t, mantenimientoService.Create(ctx, m))

	// Listar próximos
	list, err := mantenimientoService.ListProximos(ctx, 30)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

// ========== DocumentoService Tests ==========

func TestDocumentoService_Create(t *testing.T) {
	replicaService, _, documentoService, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	doc := &models.Documento{
		ReplicaID:     &replicaID,
		Tipo:          "factura",
		NombreArchivo: "test.pdf",
		MimeType:      "application/pdf",
	}

	err := documentoService.Create(ctx, doc, []byte("contenido de prueba"))
	require.NoError(t, err)
	assert.NotZero(t, doc.ID)
	assert.NotEmpty(t, doc.RutaArchivo)

	// Esperar OCR async en tests para evitar que la DB se cierre antes
	if ds, ok := documentoService.(*services.DocumentoService); ok {
		ds.WaitForOCR()
	}
}

func TestDocumentoService_Create_Validation(t *testing.T) {
	_, _, documentoService, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()

	// Sin tipo
	doc := &models.Documento{NombreArchivo: "test.pdf"}
	err := documentoService.Create(ctx, doc, []byte("contenido"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tipo")

	// Esperar OCR async en tests para evitar que la DB se cierre antes
	if ds, ok := documentoService.(*services.DocumentoService); ok {
		ds.WaitForOCR()
	}
}

func TestDocumentoService_ListByReplica(t *testing.T) {
	replicaService, _, documentoService, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	// Crear documentos
	for i := 0; i < 3; i++ {
		doc := &models.Documento{
			ReplicaID:     &replicaID,
			Tipo:          "factura",
			NombreArchivo: fmt.Sprintf("factura%d.pdf", i),
			MimeType:      "application/pdf",
		}
		require.NoError(t, documentoService.Create(ctx, doc, []byte("contenido")))
	}

	// Listar
	list, err := documentoService.ListByReplica(ctx, replicaID)
	require.NoError(t, err)
	assert.Len(t, list, 3)

	// Esperar OCR async en tests para evitar que la DB se cierre antes
	if ds, ok := documentoService.(*services.DocumentoService); ok {
		ds.WaitForOCR()
	}
}

func TestDocumentoService_Delete(t *testing.T) {
	replicaService, _, documentoService, _, db := setupServiceTest(t)
	defer db.Close()
	ctx := context.Background()
	replicaID := createTestReplica(t, replicaService)

	doc := &models.Documento{
		ReplicaID:     &replicaID,
		Tipo:          "factura",
		NombreArchivo: "borrar.pdf",
		MimeType:      "application/pdf",
	}
	require.NoError(t, documentoService.Create(ctx, doc, []byte("contenido")))

	// Borrar
	err := documentoService.Delete(ctx, doc.ID)
	require.NoError(t, err)

	// Verificar que no existe
	_, err = documentoService.GetByID(ctx, doc.ID)
	assert.Error(t, err)

	// Esperar OCR async en tests para evitar que la DB se cierre antes
	if ds, ok := documentoService.(*services.DocumentoService); ok {
		ds.WaitForOCR()
	}
}
