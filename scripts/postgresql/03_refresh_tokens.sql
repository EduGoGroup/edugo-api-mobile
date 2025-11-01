-- =====================================================
-- Tabla: refresh_tokens
-- Descripción: Almacena refresh tokens para renovación de access tokens
-- Autor: EduGo Team
-- Fecha: 2024-10-31
-- =====================================================

-- Crear tabla refresh_tokens
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA256 del token (no se guarda token original)
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    client_info JSONB,                       -- Info del cliente: {"ip": "192.168.1.1", "user_agent": "..."}
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    revoked_at TIMESTAMP,
    replaced_by UUID REFERENCES refresh_tokens(id),  -- Para token rotation
    CONSTRAINT check_expires_future CHECK (expires_at > created_at)
);

-- =====================================================
-- Índices para optimización de queries
-- =====================================================

-- Índice por user_id (para buscar todos los tokens de un usuario)
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id
    ON refresh_tokens(user_id);

-- Índice por token_hash (para búsqueda rápida de tokens)
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash
    ON refresh_tokens(token_hash);

-- Índice por expires_at (para limpieza de tokens expirados)
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at
    ON refresh_tokens(expires_at);

-- Índice compuesto para búsqueda de tokens válidos
-- Usa WHERE para índice parcial (solo tokens no revocados)
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_valid
    ON refresh_tokens(token_hash, user_id)
    WHERE revoked_at IS NULL;

-- =====================================================
-- Comentarios de documentación
-- =====================================================

COMMENT ON TABLE refresh_tokens IS
    'Almacena refresh tokens para renovación de access tokens JWT. Los tokens se guardan hasheados (SHA256) por seguridad.';

COMMENT ON COLUMN refresh_tokens.token_hash IS
    'SHA256 del token en texto plano. El token original NO se guarda por seguridad.';

COMMENT ON COLUMN refresh_tokens.client_info IS
    'Información del cliente en formato JSON: {"ip": "192.168.1.1", "user_agent": "Mozilla/5.0...", "device": "mobile"}';

COMMENT ON COLUMN refresh_tokens.revoked_at IS
    'Timestamp cuando el token fue revocado manualmente (logout o revoke-all). NULL = token válido.';

COMMENT ON COLUMN refresh_tokens.replaced_by IS
    'ID del nuevo token que reemplazó a este (cuando se usa token rotation). NULL si no fue reemplazado.';

-- =====================================================
-- Función de limpieza automática (housekeeping)
-- =====================================================

-- Función para eliminar tokens expirados
CREATE OR REPLACE FUNCTION cleanup_expired_refresh_tokens()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM refresh_tokens
    WHERE expires_at < NOW()
      AND (revoked_at IS NOT NULL OR expires_at < NOW() - INTERVAL '30 days');

    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_expired_refresh_tokens() IS
    'Elimina tokens expirados hace más de 30 días. Retorna el número de tokens eliminados.';

-- =====================================================
-- Trigger para prevenir modificación de token_hash
-- =====================================================

CREATE OR REPLACE FUNCTION prevent_token_hash_modification()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.token_hash IS DISTINCT FROM NEW.token_hash THEN
        RAISE EXCEPTION 'token_hash cannot be modified after creation';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_prevent_token_hash_modification
    BEFORE UPDATE ON refresh_tokens
    FOR EACH ROW
    EXECUTE FUNCTION prevent_token_hash_modification();

-- =====================================================
-- Vista para tokens activos (útil para debugging)
-- =====================================================

CREATE OR REPLACE VIEW active_refresh_tokens AS
SELECT
    rt.id,
    rt.user_id,
    u.email,
    rt.client_info,
    rt.created_at,
    rt.expires_at,
    (rt.expires_at - NOW()) AS time_until_expiry,
    rt.replaced_by IS NOT NULL AS was_rotated
FROM refresh_tokens rt
JOIN "user" u ON u.id = rt.user_id
WHERE rt.revoked_at IS NULL
  AND rt.expires_at > NOW()
ORDER BY rt.created_at DESC;

COMMENT ON VIEW active_refresh_tokens IS
    'Vista de refresh tokens activos (no revocados y no expirados) con información del usuario.';

-- =====================================================
-- Grants (si usas roles específicos)
-- =====================================================

-- GRANT SELECT, INSERT, UPDATE, DELETE ON refresh_tokens TO edugo_api_user;
-- GRANT SELECT ON active_refresh_tokens TO edugo_api_user;
-- GRANT EXECUTE ON FUNCTION cleanup_expired_refresh_tokens() TO edugo_api_user;
