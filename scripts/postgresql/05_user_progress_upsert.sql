-- =====================================================
-- Script: 05_user_progress_upsert.sql
-- Descripción: Mejoras en tabla material_progress para UPSERT idempotente
-- Autor: EduGo Team
-- Fecha: 2025-11-05
-- Versión: 1.0.0
-- =====================================================

-- NOTA: La tabla material_progress YA tiene PRIMARY KEY (material_id, user_id)
-- lo cual garantiza unicidad. Este script agrega campos y índices adicionales
-- para optimizar operaciones UPSERT.

-- =====================================================
-- Agregar campo completed_at para tracking de completitud
-- =====================================================

-- Agregar columna completed_at si no existe
-- Se usa para rastrear la fecha exacta cuando un usuario completó un material (percentage = 100)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'material_progress'
        AND column_name = 'completed_at'
    ) THEN
        ALTER TABLE material_progress
        ADD COLUMN completed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

        COMMENT ON COLUMN material_progress.completed_at IS
            'Fecha y hora cuando el usuario completó el material (percentage = 100). NULL si no completado.';
    END IF;
END $$;

-- =====================================================
-- Índices de Performance para UPSERT
-- =====================================================

-- Índice compuesto para operación UPSERT
-- Usado en query: INSERT ... ON CONFLICT (user_id, material_id) DO UPDATE ...
-- La PRIMARY KEY ya cubre (material_id, user_id), creamos índice inverso para queries que filtran por user_id primero
CREATE INDEX IF NOT EXISTS idx_material_progress_user_material
    ON material_progress(user_id, material_id);

-- Índice para filtrar por fecha de actualización
-- Usado en queries de "usuarios activos" (last_accessed_at >= NOW() - INTERVAL '30 days')
CREATE INDEX IF NOT EXISTS idx_material_progress_last_accessed_recent
    ON material_progress(last_accessed_at DESC)
    WHERE last_accessed_at >= NOW() - INTERVAL '30 days';

-- Índice para queries de materiales completados recientemente
CREATE INDEX IF NOT EXISTS idx_material_progress_completed
    ON material_progress(user_id, completed_at DESC)
    WHERE completed_at IS NOT NULL;

-- =====================================================
-- Actualizar valores de completed_at para registros existentes
-- =====================================================

-- Inicializar completed_at para registros que ya están con percentage = 100
-- pero no tienen completed_at asignado (por ser creados antes de esta migración)
UPDATE material_progress
SET completed_at = updated_at
WHERE percentage = 100
  AND status = 'completed'
  AND completed_at IS NULL;

-- =====================================================
-- Trigger para mantener completed_at sincronizado
-- =====================================================

-- Función para actualizar completed_at automáticamente
CREATE OR REPLACE FUNCTION update_material_progress_completed_at()
RETURNS TRIGGER AS $$
BEGIN
    -- Si el progreso llega a 100%, establecer completed_at
    IF NEW.percentage = 100 AND NEW.status = 'completed' AND OLD.completed_at IS NULL THEN
        NEW.completed_at = NOW();
    END IF;

    -- Si el progreso baja de 100% (re-lectura), limpiar completed_at
    IF NEW.percentage < 100 AND OLD.completed_at IS NOT NULL THEN
        NEW.completed_at = NULL;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger si no existe
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'trigger_sync_material_progress_completed_at'
    ) THEN
        CREATE TRIGGER trigger_sync_material_progress_completed_at
            BEFORE UPDATE ON material_progress
            FOR EACH ROW
            EXECUTE FUNCTION update_material_progress_completed_at();
    END IF;
END $$;

-- =====================================================
-- Comentarios adicionales de documentación
-- =====================================================

COMMENT ON TRIGGER trigger_sync_material_progress_completed_at ON material_progress IS
    'Sincroniza automáticamente el campo completed_at cuando percentage llega a 100% o baja de 100%.';

COMMENT ON FUNCTION update_material_progress_completed_at() IS
    'Función de trigger para mantener completed_at sincronizado con percentage y status.';
