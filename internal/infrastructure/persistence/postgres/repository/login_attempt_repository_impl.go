package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

type postgresLoginAttemptRepository struct {
	db *sql.DB
}

// NewPostgresLoginAttemptRepository crea una nueva instancia del repositorio
func NewPostgresLoginAttemptRepository(db *sql.DB) repository.LoginAttemptRepository {
	return &postgresLoginAttemptRepository{db: db}
}

func (r *postgresLoginAttemptRepository) RecordAttempt(ctx context.Context, attempt repository.LoginAttemptData) error {
	query := `
		INSERT INTO login_attempts (identifier, attempt_type, successful, user_agent, ip_address, attempted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		attempt.Identifier,
		attempt.AttemptType,
		attempt.Successful,
		attempt.UserAgent,
		attempt.IPAddress,
		attempt.AttemptedAt,
	)

	if err != nil {
		return fmt.Errorf("error al registrar intento de login: %w", err)
	}

	return nil
}

func (r *postgresLoginAttemptRepository) CountFailedAttempts(ctx context.Context, identifier string, windowMinutes int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM login_attempts
		WHERE identifier = $1
		  AND successful = false
		  AND attempted_at > NOW() - ($2 || ' minutes')::INTERVAL
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, identifier, windowMinutes).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar intentos fallidos: %w", err)
	}

	return count, nil
}

func (r *postgresLoginAttemptRepository) IsRateLimited(ctx context.Context, identifier string, maxAttempts int, windowMinutes int) (bool, error) {
	count, err := r.CountFailedAttempts(ctx, identifier, windowMinutes)
	if err != nil {
		return false, err
	}

	return count >= maxAttempts, nil
}
