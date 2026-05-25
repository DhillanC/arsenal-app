package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/services"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/persistence/sqlite"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/storage/local"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/web"
)

func main() {
	// Configuración
	dbPath := getEnv("DB_PATH", "./data/arsenal.db")
	appPort := getEnv("APP_PORT", "8080")

	// Inicializar base de datos
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Error inicializando DB: %v", err)
	}
	defer db.Close()

	// Ejecutar migraciones embebidas
	if err := db.RunMigrations(); err != nil {
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
	handler := web.NewHandler(config, replicaService, actividadService, documentoService)
	
	srv := &http.Server{
		Addr:    ":" + appPort,
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 Arsenal App iniciado en http://localhost:%s", appPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()

	// Esperar señal de terminación
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Apagado solicitado, drenando conexiones...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error en graceful shutdown: %v", err)
	}

	log.Println("Servidor detenido limpiamente")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
