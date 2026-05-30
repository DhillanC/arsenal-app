package main

import (
	"log/slog"
	"os"
)

// initLogger inicializa slog con formato JSON para producción
// y formato texto legible para desarrollo.
func initLogger() {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
}
