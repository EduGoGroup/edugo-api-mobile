package repository

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/types/enum"
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
