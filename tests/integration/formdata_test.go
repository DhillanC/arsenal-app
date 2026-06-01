package sqlite_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReplica_Create_WithFormData verifica que el endpoint acepte
// application/x-www-form-urlencoded (lo que envía HTMX por defecto).
// Issue #26: HTMX form-data vs JSON.
func TestReplica_Create_WithFormData(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	form := url.Values{}
	form.Set("nombre", "HK416 FormData")
	form.Set("tipo", "AEG")
	form.Set("estado", "activo")
	form.Set("fecha_adquisicion", "2026-05-25")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equalf(t, http.StatusCreated, w.Code, "body: %s", w.Body.String())
	assert.Contains(t, w.Body.String(), "HK416 FormData")
}

// TestReplica_Create_WithMultipartForm verifica que el endpoint acepte
// multipart/form-data (otro formato común de formularios HTML).
func TestReplica_Create_WithMultipartForm(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	var b strings.Builder
	w := multipart.NewWriter(&b)
	_ = w.WriteField("nombre", "HK416 Multipart")
	_ = w.WriteField("tipo", "AEG")
	_ = w.WriteField("estado", "activo")
	_ = w.WriteField("fecha_adquisicion", "2026-05-25")
	require.NoError(t, w.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", strings.NewReader(b.String()))
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equalf(t, http.StatusCreated, rec.Code, "body: %s", rec.Body.String())
	assert.Contains(t, rec.Body.String(), "HK416 Multipart")
}

// TestReplica_Update_WithFormData verifica que PUT también acepte form-data.
func TestReplica_Update_WithFormData(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear primero con JSON
	payload := map[string]any{
		"nombre":            "Original",
		"tipo":              "AEG",
		"estado":            "activo",
		"fecha_adquisicion": "2026-05-25",
	}
	raw, _ := json.Marshal(payload)
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", bytes.NewReader(raw))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()
	handler.ServeHTTP(postW, postReq)
	require.Equal(t, http.StatusCreated, postW.Code)

	var created map[string]any
	require.NoError(t, json.Unmarshal(postW.Body.Bytes(), &created))
	id := int(created["id"].(float64))

	// Actualizar con form-data
	form := url.Values{}
	form.Set("nombre", "Actualizado FormData")
	form.Set("tipo", "GBB")
	form.Set("estado", "activo")

	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/replicas/"+strconv.Itoa(id), strings.NewReader(form.Encode()))
	putReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	putW := httptest.NewRecorder()
	handler.ServeHTTP(putW, putReq)

	require.Equalf(t, http.StatusOK, putW.Code, "body: %s", putW.Body.String())
	assert.Contains(t, putW.Body.String(), "Actualizado FormData")
}
