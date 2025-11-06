-- =====================================================
-- Script: 04_material_versions.sql
-- Descripción: Tabla de versionado histórico de materiales educativos
-- Autor: EduGo Team
-- Fecha: 2025-11-05
-- Versión: 1.0.0
-- =====================================================

-- =====================================================
-- Tabla: material_versions
-- Descripción: Almacena historial completo de versiones de materiales
-- =====================================================

CREATE TABLE IF NOT EXISTS material_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content_url TEXT NOT NULL,
    changed_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraint para prevenir versiones duplicadas del mismo material
    CONSTRAINT unique_material_version UNIQUE(material_id, version_number),

    -- Constraint para garantizar que version_number sea positivo
    CONSTRAINT check_version_positive CHECK (version_number > 0)
);

-- =====================================================
-- Índices de Performance
-- =====================================================

-- Índice para JOINs frecuentes con tabla materials
-- Usado en query: SELECT m.*, v.* FROM materials m LEFT JOIN material_versions v ON m.id = v.material_id
CREATE INDEX IF NOT EXISTS idx_material_versions_material_id
    ON material_versions(material_id);

-- Índice para ordenar versiones por fecha de creación (DESC)
-- Usado para obtener últimas versiones primero
CREATE INDEX IF NOT EXISTS idx_material_versions_created_at
    ON material_versions(created_at DESC);

-- Índice compuesto para queries que filtran por material y ordenan por versión
-- Optimiza: WHERE material_id = $1 ORDER BY version_number DESC
CREATE INDEX IF NOT EXISTS idx_material_versions_material_version
    ON material_versions(material_id, version_number DESC);

-- =====================================================
-- Comentarios de Documentación
-- =====================================================

COMMENT ON TABLE material_versions IS
    'Historial de versiones de materiales educativos. Cada cambio significativo genera una nueva versión.';

COMMENT ON COLUMN material_versions.material_id IS
    'ID del material al que pertenece esta versión. CASCADE delete si se elimina el material.';

COMMENT ON COLUMN material_versions.version_number IS
    'Número secuencial de versión. Comienza en 1 y se incrementa con cada cambio.';

COMMENT ON COLUMN material_versions.title IS
    'Título del material en esta versión específica. Puede cambiar entre versiones.';

COMMENT ON COLUMN material_versions.content_url IS
    'URL del contenido en S3 para esta versión. Cada versión apunta a un archivo diferente.';

COMMENT ON COLUMN material_versions.changed_by IS
    'ID del usuario (profesor) que creó esta versión. Se usa para auditoría.';

COMMENT ON COLUMN material_versions.created_at IS
    'Fecha y hora de creación de esta versión. Inmutable una vez creada.';
