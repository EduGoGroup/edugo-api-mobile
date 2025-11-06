-- ============================================================
-- Migration: 05_indexes_materials.sql
-- Description: Agregar índice descendente en materials.updated_at
--              para optimizar queries de listado cronológico
-- Author: Claude Code / EduGo Team
-- Date: 2025-11-05
-- ============================================================

-- Objetivo:
-- Mejorar performance de queries que ordenan materiales por fecha
-- de actualización más reciente (patrón común en la aplicación).
--
-- Queries beneficiadas:
-- 1. SELECT * FROM materials ORDER BY updated_at DESC LIMIT N;
-- 2. SELECT * FROM materials WHERE course_id = X ORDER BY updated_at DESC;
-- 3. SELECT * FROM materials WHERE type = 'Y' ORDER BY updated_at DESC;
--
-- Mejora esperada: 5-10x más rápido (de 50-200ms a 5-20ms)

-- Crear índice descendente de forma idempotente
CREATE INDEX IF NOT EXISTS idx_materials_updated_at
ON materials(updated_at DESC);

-- Verificación:
-- Después de ejecutar este script, verificar con:
-- SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';
--
-- Validar uso del índice con:
-- EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20;
-- Debe mostrar: "Index Scan using idx_materials_updated_at"

-- Rollback (si es necesario):
-- DROP INDEX IF EXISTS idx_materials_updated_at;
