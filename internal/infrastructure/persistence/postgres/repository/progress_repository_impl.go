package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type postgresProgressRepository struct {
	db *sql.DB
}

func NewPostgresProgressRepository(db *sql.DB) repository.ProgressRepository {
	return &postgresProgressRepository{db: db}
}

func (r *postgresProgressRepository) Save(ctx context.Context, progress *pgentities.Progress) error {
	query := `
		INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		progress.MaterialID,
		progress.UserID,
		progress.Percentage,
		progress.LastPage,
		progress.Status,
		progress.LastAccessedAt,
		progress.CreatedAt,
		progress.UpdatedAt,
	)

	return err
}

func (r *postgresProgressRepository) FindByMaterialAndUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*pgentities.Progress, error) {
	query := `
		SELECT material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at
		FROM material_progress
		WHERE material_id = $1 AND user_id = $2
	`

	var (
		matID          uuid.UUID
		uID            uuid.UUID
		percentage     int
		lastPage       int
		status         string
		lastAccessedAt time.Time
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := r.db.QueryRowContext(ctx, query, materialID.UUID(), userID.UUID()).Scan(
		&matID, &uID, &percentage, &lastPage, &status, &lastAccessedAt, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pgentities.Progress{
		MaterialID:     matID,
		UserID:         uID,
		Percentage:     percentage,
		LastPage:       lastPage,
		Status:         status,
		LastAccessedAt: lastAccessedAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

func (r *postgresProgressRepository) Update(ctx context.Context, progress *pgentities.Progress) error {
	query := `
		UPDATE material_progress
		SET percentage = $1, last_page = $2, status = $3, last_accessed_at = $4, updated_at = $5
		WHERE material_id = $6 AND user_id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		progress.Percentage,
		progress.LastPage,
		progress.Status,
		progress.LastAccessedAt,
		progress.UpdatedAt,
		progress.MaterialID,
		progress.UserID,
	)

	return err
}

// Upsert implementa operación idempotente INSERT o UPDATE usando ON CONFLICT de PostgreSQL.
// Si el registro (material_id, user_id) existe, se actualiza; si no existe, se inserta.
// Cuando percentage = 100, se establece status = 'completed' automáticamente.
// Retorna la entidad Progress actualizada.
func (r *postgresProgressRepository) Upsert(ctx context.Context, progress *pgentities.Progress) (*pgentities.Progress, error) {
	// Query UPSERT usando ON CONFLICT de PostgreSQL
	// La PRIMARY KEY (material_id, user_id) garantiza unicidad
	query := `
		INSERT INTO material_progress (
			material_id, user_id, percentage, last_page, status,
			last_accessed_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (material_id, user_id)
		DO UPDATE SET
			percentage = EXCLUDED.percentage,
			last_page = EXCLUDED.last_page,
			status = EXCLUDED.status,
			last_accessed_at = EXCLUDED.last_accessed_at,
			updated_at = EXCLUDED.updated_at
		RETURNING material_id, user_id, percentage, last_page, status,
		          last_accessed_at, created_at, updated_at
	`

	var (
		matID          uuid.UUID
		uID            uuid.UUID
		percentage     int
		lastPage       int
		status         string
		lastAccessedAt time.Time
		createdAt      time.Time
		updatedAt      time.Time
	)

	// Ejecutar query UPSERT y escanear resultado
	err := r.db.QueryRowContext(ctx, query,
		progress.MaterialID,
		progress.UserID,
		progress.Percentage,
		progress.LastPage,
		progress.Status,
		progress.LastAccessedAt,
		progress.CreatedAt,
		progress.UpdatedAt,
	).Scan(
		&matID, &uID, &percentage, &lastPage, &status,
		&lastAccessedAt, &createdAt, &updatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Retornar entidad Progress con datos retornados de la BD
	return &pgentities.Progress{
		MaterialID:     matID,
		UserID:         uID,
		Percentage:     percentage,
		LastPage:       lastPage,
		Status:         status,
		LastAccessedAt: lastAccessedAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

// CountActiveUsers cuenta usuarios únicos con actividad reciente (últimos 30 días)
// Usado para estadísticas globales del sistema
func (r *postgresProgressRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT user_id)
		FROM material_progress
		WHERE last_accessed_at >= NOW() - INTERVAL '30 days'
	`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CalculateAverageProgress calcula el promedio de progreso de todos los usuarios
// Usado para estadísticas globales del sistema
func (r *postgresProgressRepository) CalculateAverageProgress(ctx context.Context) (float64, error) {
	query := `SELECT COALESCE(AVG(percentage), 0) FROM material_progress`

	var avgProgress float64
	err := r.db.QueryRowContext(ctx, query).Scan(&avgProgress)
	if err != nil {
		return 0.0, err
	}

	return avgProgress, nil
}
