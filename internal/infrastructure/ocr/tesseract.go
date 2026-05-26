package ocr

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
)

// OCRClient encapsula el cliente de Tesseract
type OCRClient struct {
	available bool
}

// NewOCRClient crea un nuevo cliente OCR
func NewOCRClient() (*OCRClient, error) {
	if !IsAvailable() {
		return &OCRClient{available: false}, nil
	}
	return &OCRClient{available: true}, nil
}

// ExtractText extrae texto de una imagen o PDF
func (o *OCRClient) ExtractText(filePath string) (string, error) {
	if !o.available {
		return "", fmt.Errorf("tesseract no está instalado")
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return "", fmt.Errorf("PDF OCR requiere conversión a imagen primero")
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff":
		return o.extractFromImage(filePath)
	default:
		return "", fmt.Errorf("formato no soportado para OCR: %s", ext)
	}
}

// extractFromImage extrae texto de una imagen (placeholder)
func (o *OCRClient) extractFromImage(filePath string) (string, error) {
	// Verificar que el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo no encontrado: %s", filePath)
	}

	// Abrir imagen para verificar formato
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("abrir imagen: %w", err)
	}
	defer file.Close()

	// Decodificar para verificar que es una imagen válida
	_, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("decodificar imagen: %w", err)
	}

	// Placeholder: en producción usaría gosseract
	return fmt.Sprintf("[OCR Placeholder - imagen %s válida]", format), nil
}

// Close libera recursos del cliente OCR
func (o *OCRClient) Close() error {
	return nil
}

// IsAvailable verifica si Tesseract está instalado
func IsAvailable() bool {
	_, err := os.Stat("/usr/bin/tesseract")
	if err != nil {
		_, err = os.Stat("/opt/homebrew/bin/tesseract")
	}
	return err == nil
}
