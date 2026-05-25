package main

import (
	"log"
	"os"

	"github.com/DhillanC/arsenal-app/internal/domain/services"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/persistence/sqlite"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/storage/local"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/web"
)

func main() {
	// Configuración
	dbPath := getEnv("DB_PATH", "./data/arsenal.db")
	appPort := getEnv("APP_PORT", "8080")
	migrationsDir := "./internal/infrastructure/persistence/sqlite/migrations"

	// Inicializar base de datos
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Error inicializando DB: %v", err)
	}
	defer db.Close()

	// Ejecutar migraciones
	if err := db.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Error en migraciones: %v", err)
	}

	// Repositorios
	replicaRepo := sqlite.NewReplicaRepository(db.Conn)
	actividadRepo := sqlite.NewActividadRepository(db.Conn)
	documentoRepo := sqlite.NewDocumentoRepository(db.Conn)

	// Servicios (capa de aplicación)
	replicaService := services.NewReplicaService(replicaRepo)
	actividadService := services.NewActividadService(actividadRepo)
	
	// Storage
	uploadPath := getEnv("UPLOAD_PATH", "./uploads")
	storage := local.NewStorage(uploadPath)
	documentoService := services.NewDocumentoService(documentoRepo, storage)

	// Servidor HTTP
	config := web.Config{
		Port: appPort,
	}
	server := web.NewServer(config, replicaService, actividadService, documentoService)
	
	log.Printf("🚀 Arsenal App iniciado en http://localhost:%s", appPort)
	if err := server.Run(); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
