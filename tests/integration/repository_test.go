package sqlite_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/persistence/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlite.DB {
	db, err := sqlite.NewDB(":memory:")
	require.NoError(t, err)
	
	// Usar ruta absoluta desde el proyecto
	migrationsDir := "../../internal/infrastructure/persistence/sqlite/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Fallback para cuando se ejecuta desde diferentes directorios
		migrationsDir = "internal/infrastructure/persistence/sqlite/migrations"
	}
	
	err = db.RunMigrations(migrationsDir)
	require.NoError(t, err)
	
	return db
}

func TestReplicaRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := sqlite.NewReplicaRepository(db.Conn)
	ctx := context.Background()
	
	t.Run("Create", func(t *testing.T) {
		replica := &models.Replica{
			Nombre:           "HK416 A5",
			Marca:            "VFC",
			Modelo:           "HK416 A5",
			Tipo:             "AEG",
			NumeroSerie:      "TEST001",
			FechaAdquisicion: time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC),
			Proveedor:        "Universal de deportes SAS",
			CostoAdquisicion: 350000,
			Estado:           "activo",
			FPS:              380,
			Joules:           1.2,
			PesoGramos:       3200,
			LongitudMM:       900,
			HopUp:            "Ajustable",
			CapacidadCargador: 120,
			Notas:            "Test replica",
		}
		
		err := repo.Create(ctx, replica)
		require.NoError(t, err)
		assert.NotZero(t, replica.ID)
	})
	
	t.Run("GetByID", func(t *testing.T) {
		// Crear primero
		replica := &models.Replica{
			Nombre:           "M4A1",
			Marca:            "Tokyo Marui",
			Tipo:             "AEG",
			Estado:           "activo",
			FechaAdquisicion: time.Now(),
		}
		err := repo.Create(ctx, replica)
		require.NoError(t, err)
		
		// Obtener
		found, err := repo.GetByID(ctx, replica.ID)
		require.NoError(t, err)
		assert.Equal(t, replica.Nombre, found.Nombre)
		assert.Equal(t, replica.Marca, found.Marca)
	})
	
	t.Run("List", func(t *testing.T) {
		// Crear varias
		for i := 0; i < 3; i++ {
			replica := &models.Replica{
				Nombre:           fmt.Sprintf("Replica %d", i),
				Tipo:             "AEG",
				Estado:           "activo",
				NumeroSerie:      fmt.Sprintf("SERIE%d", i),
				FechaAdquisicion: time.Now(),
			}
			err := repo.Create(ctx, replica)
			require.NoError(t, err)
		}
		
		replicas, err := repo.List(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(replicas), 3)
	})
	
	t.Run("Update", func(t *testing.T) {
		replica := &models.Replica{
			Nombre:           "Original",
			Tipo:             "AEG",
			Estado:           "activo",
			NumeroSerie:      "UPDATE001",
			FechaAdquisicion: time.Now(),
		}
		err := repo.Create(ctx, replica)
		require.NoError(t, err)
		
		replica.Nombre = "Updated"
		err = repo.Update(ctx, replica)
		require.NoError(t, err)
		
		found, err := repo.GetByID(ctx, replica.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated", found.Nombre)
	})
	
	t.Run("Delete", func(t *testing.T) {
		replica := &models.Replica{
			Nombre:           "ToDelete",
			Tipo:             "AEG",
			Estado:           "activo",
			NumeroSerie:      "DELETE001",
			FechaAdquisicion: time.Now(),
		}
		err := repo.Create(ctx, replica)
		require.NoError(t, err)
		
		err = repo.Delete(ctx, replica.ID)
		require.NoError(t, err)
		
		found, err := repo.GetByID(ctx, replica.ID)
		require.NoError(t, err)
		assert.Equal(t, "archivado", found.Estado)
	})
}

func TestActividadRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := sqlite.NewActividadRepository(db.Conn)
	ctx := context.Background()
	
	t.Run("Create and ListByReplica", func(t *testing.T) {
		// Primero crear una réplica
		replicaRepo := sqlite.NewReplicaRepository(db.Conn)
		replica := &models.Replica{
			Nombre:           "HK416 A5",
			Tipo:             "AEG",
			Estado:           "activo",
			FechaAdquisicion: time.Now(),
		}
		err := replicaRepo.Create(ctx, replica)
		require.NoError(t, err)
		
		actividad := &models.Actividad{
			ReplicaID:        replica.ID,
			Fecha:            time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC),
			Tipo:             "compra",
			Descripcion:      "Adquisición HK416 A5",
			ProveedorTecnico: "Universal de deportes SAS",
			Costo:            350000,
		}
		
		err = repo.Create(ctx, actividad)
		require.NoError(t, err)
		assert.NotZero(t, actividad.ID)
		
		actividades, err := repo.ListByReplica(ctx, replica.ID)
		require.NoError(t, err)
		assert.Len(t, actividades, 1)
		assert.Equal(t, "compra", actividades[0].Tipo)
	})
}

func TestDocumentoRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := sqlite.NewDocumentoRepository(db.Conn)
	ctx := context.Background()
	
	t.Run("Create and SearchByOCR", func(t *testing.T) {
		// Primero crear una réplica
		replicaRepo := sqlite.NewReplicaRepository(db.Conn)
		replica := &models.Replica{
			Nombre:           "HK416 A5",
			Tipo:             "AEG",
			Estado:           "activo",
			FechaAdquisicion: time.Now(),
		}
		err := replicaRepo.Create(ctx, replica)
		require.NoError(t, err)
		
		doc := &models.Documento{
			ReplicaID:       intPtr(replica.ID),
			Tipo:            "factura",
			NombreArchivo:   "factura.pdf",
			RutaArchivo:     "/uploads/1/2026-05/factura.pdf",
			MimeType:        "application/pdf",
			TamanoBytes:     1024,
			OCRTexto:        "Factura de compra HK416 A5",
			NumeroDocumento: "3BAQ13218",
		}
		
		err = repo.Create(ctx, doc)
		require.NoError(t, err)
		assert.NotZero(t, doc.ID)
		
		// Buscar por OCR
		results, err := repo.SearchByOCR(ctx, "HK416")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "factura", results[0].Tipo)
	})
}

func intPtr(i int) *int {
	return &i
}
