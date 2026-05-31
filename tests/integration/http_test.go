package sqlite_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DhillanC/arsenal-app/internal/domain/services"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/persistence/sqlite"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/storage/local"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupHTTPTest arma el stack completo (DB en memoria + servicios + handler gin)
// para probar el contrato HTTP sin depender de archivos en disco.
func setupHTTPTest(t *testing.T, origins []string) (http.Handler, *sqlite.DB) {
	t.Helper()
	db, err := sqlite.NewDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, db.RunMigrations())

	replicaService := services.NewReplicaService(sqlite.NewReplicaRepository(db.Conn))
	actividadService := services.NewActividadService(sqlite.NewActividadRepository(db.Conn))
	storage := local.NewStorage(t.TempDir())
	documentoService := services.NewDocumentoService(sqlite.NewDocumentoRepository(db.Conn), storage)

	handler := web.NewHandler(web.Config{
		Port:           "8080",
		AllowedOrigins: origins,
		DB:             db.Conn,
	}, replicaService, actividadService, documentoService)
	return handler, db
}

func TestHealth_HealthyDB(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ready", body["status"])
	assert.Equal(t, "ok", body["db"])
}

func TestHealth_DegradedWhenDBClosed(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	require.NoError(t, db.Close()) // forzar fallo del Ping

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "not_ready", body["status"])
}

func TestReplica_CreateAndGet_Roundtrip(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	payload := map[string]any{
		"nombre":            "HK416 A5",
		"tipo":              "AEG",
		"estado":            "activo",
		"fecha_adquisicion": "2026-05-25",
	}
	raw, _ := json.Marshal(payload)

	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", bytes.NewReader(raw))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()
	handler.ServeHTTP(postW, postReq)
	require.Equalf(t, http.StatusCreated, postW.Code, "body: %s", postW.Body.String())

	var created map[string]any
	require.NoError(t, json.Unmarshal(postW.Body.Bytes(), &created))
	id, ok := created["id"].(float64)
	require.True(t, ok, "respuesta debe tener id numérico, got %v", created)

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/replicas", nil)
	getW := httptest.NewRecorder()
	handler.ServeHTTP(getW, getReq)
	require.Equal(t, http.StatusOK, getW.Code)

	var list []map[string]any
	require.NoError(t, json.Unmarshal(getW.Body.Bytes(), &list))
	assert.Len(t, list, 1)
	assert.Equal(t, "HK416 A5", list[0]["nombre"])
	assert.Equal(t, id, list[0]["id"])
}

func TestReplica_Create_InvalidFechaAdquisicion(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	payload := map[string]any{
		"nombre":            "HK416 A5",
		"tipo":              "AEG",
		"estado":            "activo",
		"fecha_adquisicion": "2026-13-45", // inválida
	}
	raw, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "fecha_adquisicion")
}

func TestCORS_AllowsListedOrigin(t *testing.T) {
	handler, db := setupHTTPTest(t, []string{"https://app.example.com"})
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://app.example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_BlocksUnlistedOrigin(t *testing.T) {
	handler, db := setupHTTPTest(t, []string{"https://app.example.com"})
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_PreflightShortCircuits(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/replicas", nil)
	req.Header.Set("Origin", "https://localhost:5173")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "https://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
}
