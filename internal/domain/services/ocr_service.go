package services

import (
	"fmt"
	"strings"

	"github.com/DhillanC/arsenal-app/internal/domain/models"
	"github.com/DhillanC/arsenal-app/internal/infrastructure/ocr"
)

// OCRService maneja la extracción de texto de documentos
type OCRService struct {
	client *ocr.OCRClient
}

// NewOCRService crea un nuevo servicio OCR
func NewOCRService() (*OCRService, error) {
	if !ocr.IsAvailable() {
		return nil, fmt.Errorf("tesseract no está instalado")
	}

	client, err := ocr.NewOCRClient()
	if err != nil {
		return nil, fmt.Errorf("inicializar OCR: %w", err)
	}

	return &OCRService{client: client}, nil
}

// ProcessDocument extrae OCR de un documento y lo guarda
func (s *OCRService) ProcessDocument(doc *models.Documento, filePath string) error {
	if doc == nil {
		return fmt.Errorf("documento es nil")
	}

	// Extraer texto
	text, err := s.client.ExtractText(filePath)
	if err != nil {
		// Si falla OCR, guardar error pero no bloquear
		doc.OCRTexto = fmt.Sprintf("[OCR Error: %v]", err)
		return nil
	}

	// Limpiar y truncar si es muy largo
	text = cleanOCRText(text)
	if len(text) > 10000 {
		text = text[:10000] + "... [truncado]"
	}

	doc.OCRTexto = text
	return nil
}

// SearchText busca texto en el OCR de documentos
func (s *OCRService) SearchText(docs []models.Documento, query string) []models.Documento {
	query = strings.ToLower(query)
	var results []models.Documento

	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc.OCRTexto), query) {
			results = append(results, doc)
		}
	}

	return results
}

// Close libera recursos
func (s *OCRService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// cleanOCRText limpia el texto extraído
func cleanOCRText(text string) string {
	// Eliminar espacios múltiples
	text = strings.Join(strings.Fields(text), " ")
	// Eliminar caracteres de control
	text = strings.ReplaceAll(text, "\x00", "")
	return strings.TrimSpace(text)
}