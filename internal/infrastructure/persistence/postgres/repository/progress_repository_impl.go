package repository

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

type postgresProgressRepository struct {
	db *sql.DB
}

func NewPostgresProgressRepository(db *sql.DB) repository.ProgressRepository {
	return &postgresProgressRepository{db: db}
}

func (r *postgresProgressRepository) Save(ctx context.Context, progress *entity.Progress) error {
	query := `
		INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		progress.MaterialID().String(),
		progress.UserID().String(),
		progress.Percentage(),
		progress.LastPage(),
		progress.Status().String(),
		progress.LastAccessedAt(),
		progress.CreatedAt(),
		progress.UpdatedAt(),
	)

	return err
}

func (r *postgresProgressRepository) FindByMaterialAndUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*entity.Progress, error) {
	query := `
		SELECT material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at
		FROM material_progress
		WHERE material_id = $1 AND user_id = $2
	`

	var (
		matIDStr       string
		userIDStr      string
		percentage     int
		lastPage       int
		statusStr      string
		lastAccessedAt sql.NullTime
		createdAt      sql.NullTime
		updatedAt      sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, materialID.String(), userID.String()).Scan(
		&matIDStr, &userIDStr, &percentage, &lastPage, &statusStr, &lastAccessedAt, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	matID, _ := valueobject.MaterialIDFromString(matIDStr)
	uID, _ := valueobject.UserIDFromString(userIDStr)

	return entity.ReconstructProgress(
		matID,
		uID,
		percentage,
		lastPage,
		enum.ProgressStatus(statusStr),
		lastAccessedAt.Time,
		createdAt.Time,
		updatedAt.Time,
	), nil
}

func (r *postgresProgressRepository) Update(ctx context.Context, progress *entity.Progress) error {
	query := `
		UPDATE material_progress
		SET percentage = $1, last_page = $2, status = $3, last_accessed_at = $4, updated_at = $5
		WHERE material_id = $6 AND user_id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		progress.Percentage(),
		progress.LastPage(),
		progress.Status().String(),
		progress.LastAccessedAt(),
		progress.UpdatedAt(),
		progress.MaterialID().String(),
		progress.UserID().String(),
	)

	return err
}

// Upsert implementa operación idempotente INSERT o UPDATE usando ON CONFLICT de PostgreSQL.
// Si el registro (material_id, user_id) existe, se actualiza; si no existe, se inserta.
// Cuando percentage = 100, se establece completed_at automáticamente (gracias al trigger).
// Retorna la entidad Progress actualizada.
func (r *postgresProgressRepository) Upsert(ctx context.Context, progress *entity.Progress) (*entity.Progress, error) {
	// Query UPSERT usando ON CONFLICT de PostgreSQL
	// La PRIMARY KEY (material_id, user_id) garantiza unicidad
	// completed_at se maneja automáticamente por trigger update_material_progress_completed_at
	query := `
		INSERT INTO material_progress (
			material_id, user_id, percentage, last_page, status,
			last_accessed_at, created_at, updated_at, completed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (material_id, user_id)
		DO UPDATE SET
			percentage = EXCLUDED.percentage,
			last_page = EXCLUDED.last_page,
			status = EXCLUDED.status,
			last_accessed_at = EXCLUDED.last_accessed_at,
			updated_at = EXCLUDED.updated_at,
			completed_at = CASE
				WHEN EXCLUDED.percentage = 100 AND EXCLUDED.status = 'completed' THEN
					COALESCE(material_progress.completed_at, EXCLUDED.updated_at)
				WHEN EXCLUDED.percentage < 100 THEN
					NULL
				ELSE
					material_progress.completed_at
			END
		RETURNING material_id, user_id, percentage, last_page, status,
		          last_accessed_at, created_at, updated_at, completed_at
	`

	// Calcular completed_at basado en percentage
	var completedAt sql.NullTime
	if progress.Percentage() == 100 {
		completedAt = sql.NullTime{Time: progress.UpdatedAt(), Valid: true}
	}

	var (
		matIDStr            string
		userIDStr           string
		percentage          int
		lastPage            int
		statusStr           string
		lastAccessedAt      sql.NullTime
		createdAt           sql.NullTime
		updatedAt           sql.NullTime
		returnedCompletedAt sql.NullTime
	)

	// Ejecutar query UPSERT y escanear resultado
	err := r.db.QueryRowContext(ctx, query,
		progress.MaterialID().String(),
		progress.UserID().String(),
		progress.Percentage(),
		progress.LastPage(),
		progress.Status().String(),
		progress.LastAccessedAt(),
		progress.CreatedAt(),
		progress.UpdatedAt(),
		completedAt,
	).Scan(
		&matIDStr, &userIDStr, &percentage, &lastPage, &statusStr,
		&lastAccessedAt, &createdAt, &updatedAt, &returnedCompletedAt,
	)

	if err != nil {
		return nil, err
	}

	// Reconstruir entidad Progress con datos retornados de la BD
	matID, _ := valueobject.MaterialIDFromString(matIDStr)
	uID, _ := valueobject.UserIDFromString(userIDStr)

	return entity.ReconstructProgress(
		matID,
		uID,
		percentage,
		lastPage,
		enum.ProgressStatus(statusStr),
		lastAccessedAt.Time,
		createdAt.Time,
		updatedAt.Time,
	), nil
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
