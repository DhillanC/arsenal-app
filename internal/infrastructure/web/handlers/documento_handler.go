package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	inbound "github.com/DhillanC/arsenal-app/internal/domain/ports/inbound"
	"github.com/gin-gonic/gin"
)

// maxUploadBytes es el cap real del request body para subidas.
// ParseMultipartForm(N) sin MaxBytesReader es solo umbral de memoria — el body
// completo puede ser arbitrariamente grande y Gin lo spilea a disco.
const maxUploadBytes = 10 << 20

// DocumentoHandler maneja las peticiones HTTP para documentos
type DocumentoHandler struct {
	service inbound.DocumentoService
}

// NewDocumentoHandler crea un nuevo handler
func NewDocumentoHandler(service inbound.DocumentoService) *DocumentoHandler {
	return &DocumentoHandler{service: service}
}

// Upload maneja la subida de documentos (multipart/form-data)
func (h *DocumentoHandler) Upload(c *gin.Context) {
	replicaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de réplica inválido"})
		return
	}

	// MaxBytesReader corta el body en el límite — sin esto, ParseMultipartForm
	// solo limita memoria y permite uploads arbitrariamente grandes spilados a disco.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadBytes)
	if err := c.Request.ParseMultipartForm(maxUploadBytes); err != nil {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Archivo demasiado grande o formulario inválido (máx 10MB)"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo requerido"})
		return
	}
	defer file.Close()

	// Validar tipo MIME
	mimeType := header.Header.Get("Content-Type")
	if !isAllowedMimeType(mimeType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de archivo no permitido: " + mimeType})
		return
	}

	// Validar extensión
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedExtension(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Extensión no permitida: " + ext})
		return
	}

	// Leer archivo completo — io.ReadAll maneja el loop internamente y respeta
	// MaxBytesReader que ya acotó el body a 10MB.
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leyendo archivo"})
		return
	}

	// Crear documento
	now := time.Now()
	doc := &models.Documento{
		ReplicaID:       &replicaID,
		Tipo:            c.PostForm("tipo"),
		NombreArchivo:   header.Filename,
		MimeType:        mimeType,
		TamanoBytes:     int64(header.Size),
		FechaDocumento:  &now,
		NumeroDocumento: c.PostForm("numero_documento"),
		Notas:           c.PostForm("notas"),
	}

	if err := h.service.Create(c.Request.Context(), doc, fileBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// ListByReplica lista documentos de una réplica
func (h *DocumentoHandler) ListByReplica(c *gin.Context) {
	replicaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	docs, err := h.service.ListByReplica(c.Request.Context(), replicaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, docs)
}

// ListByReplicaAndType lista documentos filtrados por tipo
func (h *DocumentoHandler) ListByReplicaAndType(c *gin.Context) {
	replicaID, err := strconv.Atoi(c.Query("replica_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "replica_id requerido"})
		return
	}

	tipo := c.Query("tipo")
	if tipo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo requerido"})
		return
	}

	docs, err := h.service.ListByReplicaAndType(c.Request.Context(), replicaID, tipo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"replica_id": replicaID,
		"tipo":       tipo,
		"count":      len(docs),
		"results":    docs,
	})
}

// Search busca documentos por texto OCR o número de documento
func (h *DocumentoHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'q' requerido"})
		return
	}

	docs, err := h.service.SearchByOCR(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"count":   len(docs),
		"results": docs,
	})
}

// isAllowedMimeType valida tipos MIME permitidos
func isAllowedMimeType(mime string) bool {
	allowed := []string{
		"application/pdf",
		"image/jpeg",
		"image/png",
		"image/gif",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}
	for _, a := range allowed {
		if a == mime {
			return true
		}
	}
	return false
}

// isAllowedExtension valida extensiones permitidas
func isAllowedExtension(ext string) bool {
	allowed := []string{".pdf", ".jpg", ".jpeg", ".png", ".gif", ".doc", ".docx", ".xls", ".xlsx"}
	for _, a := range allowed {
		if a == ext {
			return true
		}
	}
	return false
}
