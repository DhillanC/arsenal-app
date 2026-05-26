package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/services"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/persistence/sqlite"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/storage/local"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/web"
)

func main() {
	dbPath := getEnv("DB_PATH", "./data/arsenal.db")
	appPort := getEnv("APP_PORT", "8080")
	uploadPath := getEnv("UPLOAD_PATH", "./uploads")
	allowedOrigins := parseOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))

	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Error inicializando DB: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Error en migraciones: %v", err)
	}

	replicaRepo := sqlite.NewReplicaRepository(db.Conn)
	actividadRepo := sqlite.NewActividadRepository(db.Conn)
	documentoRepo := sqlite.NewDocumentoRepository(db.Conn)

	replicaService := services.NewReplicaService(replicaRepo)
	actividadService := services.NewActividadService(actividadRepo)

	storage := local.NewStorage(uploadPath)
	documentoService := services.NewDocumentoService(documentoRepo, storage)

	config := web.Config{
		Port:           appPort,
		AllowedOrigins: allowedOrigins,
		DB:             db.Conn,
	}
	handler := web.NewHandler(config, replicaService, actividadService, documentoService)

	// Timeouts protegen contra Slowloris / clientes lentos.
	// IdleTimeout cierra keep-alives ociosos para liberar FDs bajo carga.
	srv := &http.Server{
		Addr:              ":" + appPort,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("Arsenal App iniciado en http://localhost:%s (CORS=%v)", appPort, allowedOrigins)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()

	// Esperar señal de terminación con context cancelado al recibirla.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Apagado solicitado, drenando conexiones...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		// No usar Fatalf — saltea los defers (db.Close). Loguear y salir limpio.
		log.Printf("Error en graceful shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Servidor detenido limpiamente")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseOrigins convierte "https://a.com,https://b.com" en []string{"https://a.com","https://b.com"}.
// Strings vacíos se descartan para evitar agregar un "" como Origin permitido por accidente.
func parseOrigins(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
