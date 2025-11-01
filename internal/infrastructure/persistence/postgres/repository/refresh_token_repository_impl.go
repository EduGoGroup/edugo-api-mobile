package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/google/uuid"
)

type postgresRefreshTokenRepository struct {
	db *sql.DB
}

// NewPostgresRefreshTokenRepository crea una nueva instancia del repositorio
func NewPostgresRefreshTokenRepository(db *sql.DB) repository.RefreshTokenRepository {
	return &postgresRefreshTokenRepository{db: db}
}

func (r *postgresRefreshTokenRepository) Store(ctx context.Context, token repository.RefreshTokenData) error {
	// Serializar client_info a JSON
	clientInfoJSON, err := json.Marshal(token.ClientInfo)
	if err != nil {
		return fmt.Errorf("error al serializar client_info: %w", err)
	}

	query := `
		INSERT INTO refresh_tokens (id, token_hash, user_id, client_info, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.db.ExecContext(ctx, query,
		token.ID,
		token.TokenHash,
		token.UserID,
		clientInfoJSON,
		token.ExpiresAt,
		token.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al guardar refresh token: %w", err)
	}

	return nil
}

func (r *postgresRefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*repository.RefreshTokenData, error) {
	query := `
		SELECT id, token_hash, user_id, client_info, expires_at, created_at, revoked_at, replaced_by
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var token repository.RefreshTokenData
	var clientInfoJSON []byte

	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.TokenHash,
		&token.UserID,
		&clientInfoJSON,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
		&token.ReplacedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Token no encontrado
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar refresh token: %w", err)
	}

	// Deserializar client_info
	if len(clientInfoJSON) > 0 {
		if err := json.Unmarshal(clientInfoJSON, &token.ClientInfo); err != nil {
			return nil, fmt.Errorf("error al deserializar client_info: %w", err)
		}
	}

	return &token, nil
}

func (r *postgresRefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1 AND revoked_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("error al revocar token: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("token no encontrado o ya está revocado")
	}

	return nil
}

func (r *postgresRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error al revocar tokens del usuario: %w", err)
	}

	return nil
}

func (r *postgresRefreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	// Eliminar tokens expirados hace más de 30 días
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW()
		  AND (revoked_at IS NOT NULL OR expires_at < NOW() - INTERVAL '30 days')
	`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("error al eliminar tokens expirados: %w", err)
	}

	count, _ := result.RowsAffected()
	return count, nil
}
