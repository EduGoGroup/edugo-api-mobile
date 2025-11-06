-- =====================================================
-- Script: 01_create_schema.sql
-- Descripción: Schema completo de la base de datos EduGo Mobile
-- Autor: EduGo Team
-- Fecha: 2025-11-05
-- Versión: 1.0.0
-- =====================================================

-- Habilitar extensión para UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- Tabla: users
-- Descripción: Usuarios del sistema móvil (students, teachers, guardians)
-- =====================================================

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL,                      -- 'student', 'teacher', 'guardian'
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT check_role CHECK (role IN ('student', 'teacher', 'guardian', 'admin'))
);

-- Índices para users
CREATE INDEX IF NOT EXISTS idx_users_email
    ON users(email);

CREATE INDEX IF NOT EXISTS idx_users_role
    ON users(role)
    WHERE is_active = true;

-- Comentarios de documentación para users
COMMENT ON TABLE users IS
    'Usuarios del sistema móvil. Incluye estudiantes, profesores y tutores.';

COMMENT ON COLUMN users.email IS
    'Email del usuario. Se usa como identificador único para login.';

COMMENT ON COLUMN users.password_hash IS
    'Hash bcrypt de la contraseña. NUNCA se guarda en texto plano.';

COMMENT ON COLUMN users.role IS
    'Rol del usuario en el sistema: student, teacher, guardian, admin';

COMMENT ON COLUMN users.is_active IS
    'Indica si el usuario está activo. false = cuenta desactivada/suspendida';

-- =====================================================
-- Tabla: materials
-- Descripción: Materiales educativos (PDFs) subidos por profesores
-- =====================================================

CREATE TABLE IF NOT EXISTS materials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject_id VARCHAR(100),                        -- ID de materia (viene de API Admin)
    s3_key VARCHAR(500),                            -- Ruta del archivo en S3
    s3_url VARCHAR(1000),                           -- URL pública del archivo
    status VARCHAR(50) NOT NULL DEFAULT 'draft',    -- 'draft', 'published', 'archived'
    processing_status VARCHAR(50) NOT NULL DEFAULT 'pending',  -- 'pending', 'processing', 'completed', 'failed'
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT check_status CHECK (status IN ('draft', 'published', 'archived')),
    CONSTRAINT check_processing_status CHECK (processing_status IN ('pending', 'processing', 'completed', 'failed'))
);

-- Índices para materials
CREATE INDEX IF NOT EXISTS idx_materials_author_id
    ON materials(author_id);

CREATE INDEX IF NOT EXISTS idx_materials_subject_id
    ON materials(subject_id)
    WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_materials_status
    ON materials(status)
    WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_materials_created_at
    ON materials(created_at DESC);

-- Comentarios de documentación para materials
COMMENT ON TABLE materials IS
    'Materiales educativos (PDFs) subidos por profesores. Almacenados en AWS S3.';

COMMENT ON COLUMN materials.author_id IS
    'ID del profesor que creó el material. CASCADE delete si se elimina el profesor.';

COMMENT ON COLUMN materials.subject_id IS
    'ID de la materia asociada. Obtenido de la API Admin central.';

COMMENT ON COLUMN materials.s3_key IS
    'Ruta interna del archivo en el bucket S3: materials/{author_id}/{material_id}.pdf';

COMMENT ON COLUMN materials.s3_url IS
    'URL pública firmada para acceder al PDF. Se regenera periódicamente.';

COMMENT ON COLUMN materials.status IS
    'Estado de publicación: draft = borrador, published = visible para estudiantes, archived = oculto';

COMMENT ON COLUMN materials.processing_status IS
    'Estado del procesamiento del PDF: pending = recién subido, processing = extrayendo texto, completed = listo, failed = error';

COMMENT ON COLUMN materials.is_deleted IS
    'Soft delete. true = eliminado lógicamente (no se muestra pero se conserva por auditoría)';

-- =====================================================
-- Tabla: material_progress
-- Descripción: Progreso de lectura de materiales por estudiante
-- =====================================================

CREATE TABLE IF NOT EXISTS material_progress (
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    percentage INTEGER NOT NULL DEFAULT 0,          -- Porcentaje de lectura (0-100)
    last_page INTEGER NOT NULL DEFAULT 0,           -- Última página leída
    status VARCHAR(50) NOT NULL DEFAULT 'not_started',  -- 'not_started', 'in_progress', 'completed'
    last_accessed_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (material_id, user_id),
    CONSTRAINT check_percentage CHECK (percentage >= 0 AND percentage <= 100),
    CONSTRAINT check_last_page CHECK (last_page >= 0),
    CONSTRAINT check_progress_status CHECK (status IN ('not_started', 'in_progress', 'completed'))
);

-- Índices para material_progress
CREATE INDEX IF NOT EXISTS idx_material_progress_user_id
    ON material_progress(user_id);

CREATE INDEX IF NOT EXISTS idx_material_progress_status
    ON material_progress(user_id, status);

CREATE INDEX IF NOT EXISTS idx_material_progress_last_accessed
    ON material_progress(user_id, last_accessed_at DESC);

-- Comentarios de documentación para material_progress
COMMENT ON TABLE material_progress IS
    'Rastrea el progreso de lectura de cada estudiante en cada material.';

COMMENT ON COLUMN material_progress.percentage IS
    'Porcentaje de avance en la lectura del PDF (0 a 100). Calculado por la app móvil.';

COMMENT ON COLUMN material_progress.last_page IS
    'Número de la última página leída. Permite retomar lectura desde donde se dejó.';

COMMENT ON COLUMN material_progress.status IS
    'Estado del progreso: not_started = no iniciado, in_progress = leyendo, completed = terminado (100%)';

COMMENT ON COLUMN material_progress.last_accessed_at IS
    'Timestamp del último acceso al material. Útil para "continuar leyendo" en la app.';

-- =====================================================
-- Triggers automáticos para updated_at
-- =====================================================

-- Función genérica para actualizar updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para users
CREATE TRIGGER trigger_update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para materials
CREATE TRIGGER trigger_update_materials_updated_at
    BEFORE UPDATE ON materials
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para material_progress
CREATE TRIGGER trigger_update_material_progress_updated_at
    BEFORE UPDATE ON material_progress
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- Vistas útiles para consultas comunes
-- =====================================================

-- Vista: Materiales publicados con información del autor
CREATE OR REPLACE VIEW published_materials_with_author AS
SELECT
    m.id,
    m.title,
    m.description,
    m.subject_id,
    m.s3_url,
    m.created_at,
    m.updated_at,
    u.id as author_id,
    u.first_name as author_first_name,
    u.last_name as author_last_name,
    u.email as author_email
FROM materials m
JOIN users u ON m.author_id = u.id
WHERE m.status = 'published'
  AND m.is_deleted = false
  AND u.is_active = true;

COMMENT ON VIEW published_materials_with_author IS
    'Materiales publicados con información del autor. Usada por la app móvil para listar materiales disponibles.';

-- Vista: Progreso reciente por estudiante
CREATE OR REPLACE VIEW recent_student_progress AS
SELECT
    mp.user_id,
    mp.material_id,
    m.title as material_title,
    mp.percentage,
    mp.status,
    mp.last_accessed_at,
    u.first_name,
    u.last_name
FROM material_progress mp
JOIN materials m ON mp.material_id = m.id
JOIN users u ON mp.user_id = u.id
WHERE mp.last_accessed_at > NOW() - INTERVAL '30 days'
  AND m.is_deleted = false
ORDER BY mp.last_accessed_at DESC;

COMMENT ON VIEW recent_student_progress IS
    'Progreso reciente de estudiantes (últimos 30 días). Útil para dashboard de "continuar leyendo".';

-- =====================================================
-- Funciones auxiliares
-- =====================================================

-- Función: Obtener estadísticas de un material
CREATE OR REPLACE FUNCTION get_material_stats(p_material_id UUID)
RETURNS TABLE (
    total_views INTEGER,
    completed_count INTEGER,
    in_progress_count INTEGER,
    avg_percentage NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::INTEGER as total_views,
        COUNT(*) FILTER (WHERE status = 'completed')::INTEGER as completed_count,
        COUNT(*) FILTER (WHERE status = 'in_progress')::INTEGER as in_progress_count,
        ROUND(AVG(percentage), 2) as avg_percentage
    FROM material_progress
    WHERE material_id = p_material_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_material_stats(UUID) IS
    'Obtiene estadísticas de uso de un material: vistas, completados, en progreso, promedio de avance.';

-- =====================================================
-- Grants de permisos (ajustar según roles de BD)
-- =====================================================

-- Si usas un usuario específico para la aplicación, descomentar:
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO edugo_api_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO edugo_api_user;
-- GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO edugo_api_user;

-- =====================================================
-- Verificación de integridad
-- =====================================================

-- Verificar que las tablas se crearon correctamente
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        RAISE EXCEPTION 'Tabla users no fue creada correctamente';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'materials') THEN
        RAISE EXCEPTION 'Tabla materials no fue creada correctamente';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'material_progress') THEN
        RAISE EXCEPTION 'Tabla material_progress no fue creada correctamente';
    END IF;

    RAISE NOTICE 'Schema creado exitosamente: users, materials, material_progress';
END $$;
