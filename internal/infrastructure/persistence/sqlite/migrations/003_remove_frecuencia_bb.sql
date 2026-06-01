-- Migration 003: Remove frecuencia_bb from mantenimiento table
-- Issue: #67 - frecuencia_bb was dead data (column existed but code rejected it)

ALTER TABLE mantenimiento DROP COLUMN frecuencia_bb;
