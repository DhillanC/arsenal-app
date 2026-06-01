-- Migration 004: Add OCR status tracking for async OCR processing
-- Issue: #39 - OCR síncrono bloquea upload hasta 30s

ALTER TABLE documentos ADD COLUMN ocr_status TEXT DEFAULT 'pending' 
    CHECK(ocr_status IN ('pending', 'processing', 'completed', 'failed'));

-- Actualizar documentos existentes: si tienen ocr_texto, están completados
UPDATE documentos SET ocr_status = 'completed' WHERE ocr_texto IS NOT NULL AND ocr_texto != '';
UPDATE documentos SET ocr_status = 'pending' WHERE ocr_status IS NULL;
