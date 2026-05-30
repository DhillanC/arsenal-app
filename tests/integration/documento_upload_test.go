package sqlite_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDocumento_UploadLargeFile verifica que archivos grandes (hasta 10MB)
// se suban completos sin truncamiento (Issue #28).
func TestDocumento_UploadLargeFile(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear réplica primero
	payload := map[string]any{
		"nombre":            "Test Upload",
		"tipo":              "AEG",
		"estado":            "activo",
		"fecha_adquisicion": "2026-05-30",
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

	// Crear archivo grande (1MB para prueba rápida) con extensión .pdf
	largeContent := make([]byte, 1<<20) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	_ = w.WriteField("tipo", "manifiesto_dian")
	_ = w.WriteField("numero_documento", "TEST-001")
	
	// Crear parte con MIME type explícito para PDF
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "large_test.pdf"))
	h.Set("Content-Type", "application/pdf")
	fw, err := w.CreatePart(h)
	require.NoError(t, err)
	_, err = fw.Write(largeContent)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/replicas/"+strconv.Itoa(replicaID)+"/documentos", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equalf(t, http.StatusCreated, rec.Code, "body: %s", rec.Body.String())
	
	var doc map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &doc))
	assert.Equal(t, float64(len(largeContent)), doc["tamano_bytes"])
}

// TestDocumento_UploadExceedsMaxSize verifica que archivos > 10MB son rechazados.
func TestDocumento_UploadExceedsMaxSize(t *testing.T) {
	handler, db := setupHTTPTest(t, nil)
	defer db.Close()

	// Crear réplica primero
	payload := map[string]any{
		"nombre":            "Test Upload",
		"tipo":              "AEG",
		"estado":            "activo",
		"fecha_adquisicion": "2026-05-30",
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

	// Crear archivo de 11MB (debe ser rechazado) con extensión .pdf
	largeContent := make([]byte, 11<<20) // 11MB

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	_ = w.WriteField("tipo", "manifiesto_dian")
	
	// Crear parte con MIME type explícito para PDF
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "too_large.pdf"))
	h.Set("Content-Type", "application/pdf")
	fw, err := w.CreatePart(h)
	require.NoError(t, err)
	_, err = fw.Write(largeContent)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/replicas/"+strconv.Itoa(replicaID)+"/documentos", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	assert.Contains(t, rec.Body.String(), "demasiado grande")
}
