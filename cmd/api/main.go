package main

import (
	"context"
	"fmt"
	"log/slog"
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
	initLogger()
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

// run contiene toda la lógica de arranque. Devolver error en vez de log.Fatalf
// garantiza que los `defer` (notablemente db.Close) corran en todos los paths.
func run() error {
	dbPath := getEnv("DB_PATH", "./data/arsenal.db")
	appPort := getEnv("APP_PORT", "8080")
	uploadPath := getEnv("UPLOAD_PATH", "./uploads")
	allowedOrigins := parseOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))

	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		return fmt.Errorf("inicializando DB: %w", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		return fmt.Errorf("migraciones: %w", err)
	}

	replicaRepo := sqlite.NewReplicaRepository(db.Conn)
	actividadRepo := sqlite.NewActividadRepository(db.Conn)
	documentoRepo := sqlite.NewDocumentoRepository(db.Conn)

	replicaService := services.NewReplicaService(replicaRepo)
	actividadService := services.NewActividadService(actividadRepo)

	storage := local.NewStorage(uploadPath)
	documentoService := services.NewDocumentoService(documentoRepo, storage)

	config := web.Config{
		Port:            appPort,
		AllowedOrigins:  allowedOrigins,
		DB:              db.Conn,
		EnableTemplates: true,
		UploadPath:      uploadPath,
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

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Arsenal App iniciado", "port", appPort, "cors", allowedOrigins)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		slog.Info("Apagado solicitado, drenando conexiones...")
	case err := <-serverErr:
		return fmt.Errorf("servidor: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	slog.Info("Servidor detenido limpiamente")
	return nil
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
