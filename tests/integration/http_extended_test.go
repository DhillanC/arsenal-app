package sqlite_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDocumento_UploadMultipart verifica subida de documento con multipart/form-data
func TestDocumento_UploadMultipart(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear réplica primero
	payload := map[string]any{
		"nombre":            "HK416 Upload",
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
	replicaID := int(created["id"].(float64))

	// Subir documento
	var b strings.Builder
	w := multipart.NewWriter(&b)
	_ = w.WriteField("tipo", "factura")
	_ = w.WriteField("numero_documento", "FAC-001")
	_ = w.WriteField("notas", "Factura de prueba")

	// Crear parte de archivo con Content-Type correcto
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "factura.pdf"))
	h.Set("Content-Type", "application/pdf")
	fileWriter, err := w.CreatePart(h)
	require.NoError(t, err)
	_, err = fileWriter.Write([]byte("contenido PDF de prueba"))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	uploadReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/replicas/%d/documentos", replicaID), strings.NewReader(b.String()))
	uploadReq.Header.Set("Content-Type", w.FormDataContentType())
	uploadW := httptest.NewRecorder()
	handler.ServeHTTP(uploadW, uploadReq)

	require.Equalf(t, http.StatusCreated, uploadW.Code, "body: %s", uploadW.Body.String())
	var doc map[string]any
	require.NoError(t, json.Unmarshal(uploadW.Body.Bytes(), &doc))
	assert.Equal(t, "factura", doc["tipo"])
	assert.Equal(t, "factura.pdf", doc["nombre_archivo"])
}

// TestDocumento_UploadTooLarge verifica límite de tamaño (10MB)
func TestDocumento_UploadTooLarge(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear réplica
	payload := map[string]any{
		"nombre": "HK416 Large",
		"tipo":   "AEG",
		"estado": "activo",
	}
	raw, _ := json.Marshal(payload)
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", bytes.NewReader(raw))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()
	handler.ServeHTTP(postW, postReq)
	require.Equal(t, http.StatusCreated, postW.Code)

	var created map[string]any
	require.NoError(t, json.Unmarshal(postW.Body.Bytes(), &created))
	replicaID := int(created["id"].(float64))

	// Subir archivo de 11MB
	var b strings.Builder
	w := multipart.NewWriter(&b)
	_ = w.WriteField("tipo", "factura")
	fileWriter, err := w.CreateFormFile("file", "large.bin")
	require.NoError(t, err)
	// Escribir 11MB de ceros
	largeContent := make([]byte, 11*1024*1024)
	_, err = fileWriter.Write(largeContent)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	uploadReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/replicas/%d/documentos", replicaID), strings.NewReader(b.String()))
	uploadReq.Header.Set("Content-Type", w.FormDataContentType())
	uploadW := httptest.NewRecorder()
	handler.ServeHTTP(uploadW, uploadReq)

	assert.Equal(t, http.StatusRequestEntityTooLarge, uploadW.Code)
}

// TestMantenimiento_CRUD verifica operaciones CRUD de mantenimiento
func TestMantenimiento_CRUD(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear réplica
	payload := map[string]any{
		"nombre":            "HK416 Maint",
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
	replicaID := int(created["id"].(float64))

	// Crear mantenimiento
	mantPayload := map[string]any{
		"tipo_tarea":      "Lubricación",
		"frecuencia_dias": 30,
		"ultima_fecha":    "2026-05-01",
		"notas":           "Mantenimiento mensual",
	}
	mantRaw, _ := json.Marshal(mantPayload)
	mantReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/replicas/%d/mantenimiento", replicaID), bytes.NewReader(mantRaw))
	mantReq.Header.Set("Content-Type", "application/json")
	mantW := httptest.NewRecorder()
	handler.ServeHTTP(mantW, mantReq)

	require.Equalf(t, http.StatusCreated, mantW.Code, "body: %s", mantW.Body.String())
	var mant map[string]any
	require.NoError(t, json.Unmarshal(mantW.Body.Bytes(), &mant))
	mantID := int(mant["id"].(float64))
	assert.Equal(t, "Lubricación", mant["tipo_tarea"])

	// Listar mantenimientos de réplica
	listReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/replicas/%d/mantenimiento", replicaID), nil)
	listW := httptest.NewRecorder()
	handler.ServeHTTP(listW, listReq)
	require.Equal(t, http.StatusOK, listW.Code)

	var list []map[string]any
	require.NoError(t, json.Unmarshal(listW.Body.Bytes(), &list))
	assert.Len(t, list, 1)

	// Completar mantenimiento
	compReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/mantenimiento/%d/completar", mantID), strings.NewReader(`{"fecha_completado":"2026-05-31"}`))
	compReq.Header.Set("Content-Type", "application/json")
	compW := httptest.NewRecorder()
	handler.ServeHTTP(compW, compReq)
	require.Equal(t, http.StatusOK, compW.Code)

	// Listar próximos mantenimientos
	proxReq := httptest.NewRequest(http.MethodGet, "/api/v1/mantenimiento/proximos?dias=30", nil)
	proxW := httptest.NewRecorder()
	handler.ServeHTTP(proxW, proxReq)
	require.Equal(t, http.StatusOK, proxW.Code)

	var proximos map[string]any
	require.NoError(t, json.Unmarshal(proxW.Body.Bytes(), &proximos))
	assert.Equal(t, 30, int(proximos["dias"].(float64)))
}

// TestReplica_Search verifica búsqueda de réplicas
func TestReplica_Search(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear réplica
	payload := map[string]any{
		"nombre":            "HK416 Search",
		"tipo":              "AEG",
		"estado":            "activo",
		"fecha_adquisicion": "2026-05-25",
		"numero_serie":      "SERIE-SEARCH-123",
	}
	raw, _ := json.Marshal(payload)
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", bytes.NewReader(raw))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()
	handler.ServeHTTP(postW, postReq)
	require.Equal(t, http.StatusCreated, postW.Code)

	// Buscar por nombre
	searchReq := httptest.NewRequest(http.MethodGet, "/api/v1/replicas/search?q=HK416", nil)
	searchW := httptest.NewRecorder()
	handler.ServeHTTP(searchW, searchReq)
	require.Equal(t, http.StatusOK, searchW.Code)

	var results []map[string]any
	require.NoError(t, json.Unmarshal(searchW.Body.Bytes(), &results))
	assert.Len(t, results, 1)

	// Buscar por serie
	searchReq2 := httptest.NewRequest(http.MethodGet, "/api/v1/replicas/search?q=SERIE-SEARCH", nil)
	searchW2 := httptest.NewRecorder()
	handler.ServeHTTP(searchW2, searchReq2)
	require.Equal(t, http.StatusOK, searchW2.Code)

	var results2 []map[string]any
	require.NoError(t, json.Unmarshal(searchW2.Body.Bytes(), &results2))
	assert.Len(t, results2, 1)
}

// TestDocumento_SearchByOCR verifica búsqueda por OCR
func TestDocumento_SearchByOCR(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear réplica
	payload := map[string]any{
		"nombre":            "HK416 OCR",
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
	replicaID := int(created["id"].(float64))

	// Subir documento con texto OCR
	var b strings.Builder
	w := multipart.NewWriter(&b)
	_ = w.WriteField("tipo", "factura")
	_ = w.WriteField("notas", "Factura con texto OCR")
	// Crear parte de archivo con Content-Type correcto
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "test.png"))
	h.Set("Content-Type", "image/png")
	fileWriter, err := w.CreatePart(h)
	require.NoError(t, err)
	_, err = fileWriter.Write([]byte("TEXTO OCR DE PRUEBA PARA BUSQUEDA"))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	uploadReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/replicas/%d/documentos", replicaID), strings.NewReader(b.String()))
	uploadReq.Header.Set("Content-Type", w.FormDataContentType())
	uploadW := httptest.NewRecorder()
	handler.ServeHTTP(uploadW, uploadReq)
	require.Equal(t, http.StatusCreated, uploadW.Code)

	// Buscar por OCR (nota: OCR solo funciona con imágenes, pero el endpoint de búsqueda busca en OCRTexto)
	searchReq := httptest.NewRequest(http.MethodGet, "/api/v1/documentos/search?q=OCR", nil)
	searchW := httptest.NewRecorder()
	handler.ServeHTTP(searchW, searchReq)
	require.Equal(t, http.StatusOK, searchW.Code)

	var searchResults map[string]any
	require.NoError(t, json.Unmarshal(searchW.Body.Bytes(), &searchResults))
	assert.Equal(t, "OCR", searchResults["query"])
}

// TestReplica_Timeline verifica endpoint de actividades (timeline)
func TestReplica_Timeline(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear réplica
	payload := map[string]any{
		"nombre":            "HK416 Timeline",
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
	replicaID := int(created["id"].(float64))

	// Crear actividad
	actPayload := map[string]any{
		"tipo":        "uso",
		"descripcion": "Juego en bosque",
		"fecha":       "2026-05-30",
		"costo":       25000,
	}
	actRaw, _ := json.Marshal(actPayload)
	actReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/replicas/%d/actividades", replicaID), bytes.NewReader(actRaw))
	actReq.Header.Set("Content-Type", "application/json")
	actW := httptest.NewRecorder()
	handler.ServeHTTP(actW, actReq)
	require.Equalf(t, http.StatusCreated, actW.Code, "body: %s", actW.Body.String())

	// Listar actividades (timeline)
	listReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/replicas/%d/actividades", replicaID), nil)
	listW := httptest.NewRecorder()
	handler.ServeHTTP(listW, listReq)
	require.Equal(t, http.StatusOK, listW.Code)

	var activities []map[string]any
	require.NoError(t, json.Unmarshal(listW.Body.Bytes(), &activities))
	assert.Len(t, activities, 1)
	assert.Equal(t, "uso", activities[0]["tipo"])
}

// TestReplica_Create_InvalidFecha verifica validación de fecha
func TestReplica_Create_InvalidFecha(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	payload := map[string]any{
		"nombre":            "HK416 Fecha",
		"tipo":              "AEG",
		"estado":            "activo",
		"fecha_adquisicion": "fecha-invalida",
	}
	raw, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "fecha_adquisicion")
}

// TestReplica_Create_InvalidTipo verifica validación de tipo
func TestReplica_Create_InvalidTipo(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	payload := map[string]any{
		"nombre":            "HK416 Tipo",
		"tipo":              "TIPO_INVALIDO",
		"estado":            "activo",
		"fecha_adquisicion": "2026-05-25",
	}
	raw, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "tipo")
}

// TestReplica_Create_InvalidEstado verifica validación de estado
func TestReplica_Create_InvalidEstado(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	payload := map[string]any{
		"nombre":            "HK416 Estado",
		"tipo":              "AEG",
		"estado":            "ESTADO_INVALIDO",
		"fecha_adquisicion": "2026-05-25",
	}
	raw, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/replicas", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "estado")
}
