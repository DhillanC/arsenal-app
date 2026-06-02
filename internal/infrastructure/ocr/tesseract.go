package ocr

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const defaultOCRTimeout = 30 * time.Second

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
		return o.extractFromPDF(filePath)
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff":
		return o.extractFromImage(filePath)
	default:
		return "", fmt.Errorf("formato no soportado para OCR: %s", ext)
	}
}

// extractFromImage extrae texto de una imagen usando el binario local de Tesseract.
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

	bin, err := resolveTesseract()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultOCRTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, filePath, "stdout")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("tesseract timeout después de %s", defaultOCRTimeout)
	}
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("tesseract OCR falló: %s", msg)
	}

	text := strings.TrimSpace(string(out))
	if text == "" {
		return "", fmt.Errorf("tesseract no extrajo texto de imagen %s", format)
	}
	return text, nil
}

// extractFromPDF convierte PDF a imagen y luego extrae texto con Tesseract.
func (o *OCRClient) extractFromPDF(filePath string) (string, error) {
	// Verificar que pdftoppm está disponible
	pdftoppm, err := exec.LookPath("pdftoppm")
	if err != nil {
		return "", fmt.Errorf("pdftoppm no está instalado (poppler-utils): %w", err)
	}

	// Crear directorio temporal para imágenes
	tmpDir, err := os.MkdirTemp("", "arsenal-ocr-*")
	if err != nil {
		return "", fmt.Errorf("crear temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Convertir PDF a PNG (primera página)
	ctx, cancel := context.WithTimeout(context.Background(), defaultOCRTimeout)
	defer cancel()

	outPrefix := filepath.Join(tmpDir, "page")
	cmd := exec.CommandContext(ctx, pdftoppm, "-png", "-f", "1", "-l", "1", filePath, outPrefix)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftoppm falló: %s", stderr.String())
	}

	// Buscar archivo generado
	pngFile := outPrefix + "-1.png"
	if _, err := os.Stat(pngFile); os.IsNotExist(err) {
		return "", fmt.Errorf("pdftoppm no generó imagen")
	}

	// Extraer texto de la imagen
	return o.extractFromImage(pngFile)
}

// Close libera recursos del cliente OCR
func (o *OCRClient) Close() error {
	return nil
}

// IsAvailable verifica si Tesseract está instalado
func IsAvailable() bool {
	_, err := resolveTesseract()
	return err == nil
}

func resolveTesseract() (string, error) {
	if path, err := exec.LookPath("tesseract"); err == nil {
		return path, nil
	}

	for _, path := range []string{"/opt/homebrew/bin/tesseract", "/usr/local/bin/tesseract", "/usr/bin/tesseract"} {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("tesseract no está instalado")
}
