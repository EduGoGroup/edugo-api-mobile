-- =====================================================
-- Tabla: login_attempts
-- Descripción: Registra intentos de login para rate limiting y auditoría
-- Autor: EduGo Team
-- Fecha: 2025-10-31
-- =====================================================

-- Crear tabla login_attempts
CREATE TABLE IF NOT EXISTS login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier VARCHAR(255) NOT NULL,  -- Email o IP del intento
    attempt_type VARCHAR(20) NOT NULL,  -- 'email' o 'ip'
    attempted_at TIMESTAMP DEFAULT NOW(),
    successful BOOLEAN NOT NULL,
    user_agent TEXT,
    ip_address INET,
    CONSTRAINT check_attempt_type CHECK (attempt_type IN ('email', 'ip'))
);

-- =====================================================
-- Índices para optimización de rate limiting
-- =====================================================

-- Índice compuesto para búsqueda de intentos recientes por identifier
-- Usa DESC en attempted_at para queries de "últimos N intentos"
CREATE INDEX IF NOT EXISTS idx_login_attempts_identifier_time
    ON login_attempts(identifier, attempted_at DESC);

-- Índice parcial para intentos fallidos (los que interesan para rate limit)
-- Solo indexa intentos fallidos recientes para ahorrar espacio
CREATE INDEX IF NOT EXISTS idx_login_attempts_failed_recent
    ON login_attempts(identifier, attempted_at DESC)
    WHERE successful = false AND attempted_at > NOW() - INTERVAL '1 hour';

-- Índice para limpieza de datos antiguos
CREATE INDEX IF NOT EXISTS idx_login_attempts_cleanup
    ON login_attempts(attempted_at)
    WHERE successful = false;

-- =====================================================
-- Comentarios de documentación
-- =====================================================

COMMENT ON TABLE login_attempts IS
    'Registra todos los intentos de login (exitosos y fallidos) para rate limiting y auditoría de seguridad';

COMMENT ON COLUMN login_attempts.identifier IS
    'Email del usuario o dirección IP del intento. Permite rate limiting por ambos criterios.';

COMMENT ON COLUMN login_attempts.attempt_type IS
    'Tipo de identificador: "email" para rate limit por usuario, "ip" para rate limit por dirección IP';

COMMENT ON COLUMN login_attempts.successful IS
    'true = login exitoso, false = credenciales incorrectas o error';

COMMENT ON COLUMN login_attempts.user_agent IS
    'User-Agent del navegador/app que hizo el intento. Útil para detectar bots.';

-- =====================================================
-- Función para contar intentos fallidos recientes
-- =====================================================

CREATE OR REPLACE FUNCTION count_failed_attempts(
    p_identifier VARCHAR(255),
    p_minutes INTEGER DEFAULT 15
) RETURNS INTEGER AS $$
DECLARE
    failed_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO failed_count
    FROM login_attempts
    WHERE identifier = p_identifier
      AND successful = false
      AND attempted_at > NOW() - (p_minutes || ' minutes')::INTERVAL;

    RETURN failed_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION count_failed_attempts(VARCHAR, INTEGER) IS
    'Cuenta intentos fallidos de login en los últimos N minutos para un identifier (email o IP)';

-- =====================================================
-- Función de limpieza automática (housekeeping)
-- =====================================================

CREATE OR REPLACE FUNCTION cleanup_old_login_attempts()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Eliminar intentos fallidos de hace más de 7 días
    -- Mantener intentos exitosos por 30 días para auditoría
    DELETE FROM login_attempts
    WHERE (successful = false AND attempted_at < NOW() - INTERVAL '7 days')
       OR (successful = true AND attempted_at < NOW() - INTERVAL '30 days');

    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_old_login_attempts() IS
    'Limpia intentos de login antiguos: fallidos >7 días, exitosos >30 días. Retorna cantidad eliminada.';

-- =====================================================
-- Vista para monitoreo de intentos sospechosos
-- =====================================================

CREATE OR REPLACE VIEW suspicious_login_attempts AS
SELECT
    identifier,
    attempt_type,
    COUNT(*) as attempt_count,
    MAX(attempted_at) as last_attempt,
    COUNT(*) FILTER (WHERE successful = false) as failed_count,
    COUNT(*) FILTER (WHERE successful = true) as success_count,
    ARRAY_AGG(DISTINCT user_agent) as user_agents,
    ARRAY_AGG(DISTINCT ip_address::text) as ip_addresses
FROM login_attempts
WHERE attempted_at > NOW() - INTERVAL '1 hour'
GROUP BY identifier, attempt_type
HAVING COUNT(*) FILTER (WHERE successful = false) >= 3  -- Al menos 3 fallos
ORDER BY failed_count DESC, last_attempt DESC;

COMMENT ON VIEW suspicious_login_attempts IS
    'Vista de actividad sospechosa: identifiers con múltiples intentos fallidos en la última hora';

-- =====================================================
-- Función para verificar si un identifier está bloqueado
-- =====================================================

CREATE OR REPLACE FUNCTION is_rate_limited(
    p_identifier VARCHAR(255),
    p_max_attempts INTEGER DEFAULT 5,
    p_window_minutes INTEGER DEFAULT 15
) RETURNS BOOLEAN AS $$
DECLARE
    failed_count INTEGER;
BEGIN
    SELECT count_failed_attempts(p_identifier, p_window_minutes)
    INTO failed_count;

    RETURN failed_count >= p_max_attempts;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION is_rate_limited(VARCHAR, INTEGER, INTEGER) IS
    'Verifica si un identifier (email o IP) está bloqueado por rate limiting.
     Default: max 5 intentos en 15 minutos.';

-- =====================================================
-- Trigger para logging automático (opcional)
-- =====================================================

CREATE OR REPLACE FUNCTION log_suspicious_attempt()
RETURNS TRIGGER AS $$
BEGIN
    -- Si es un intento fallido y ya hay varios, loguear
    IF NEW.successful = false THEN
        DECLARE
            recent_failures INTEGER;
        BEGIN
            SELECT count_failed_attempts(NEW.identifier, 15) INTO recent_failures;

            IF recent_failures >= 3 THEN
                RAISE NOTICE 'Suspicious activity detected: % has % failed attempts',
                    NEW.identifier, recent_failures;
            END IF;
        END;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_suspicious_attempts
    AFTER INSERT ON login_attempts
    FOR EACH ROW
    EXECUTE FUNCTION log_suspicious_attempt();

COMMENT ON TRIGGER trigger_log_suspicious_attempts ON login_attempts IS
    'Trigger que loguea intentos sospechosos en tiempo real (3+ fallos en 15 minutos)';
