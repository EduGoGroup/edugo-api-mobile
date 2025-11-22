package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

type postgresMaterialRepository struct {
	db *sql.DB
}

func NewPostgresMaterialRepository(db *sql.DB) repository.MaterialRepository {
	return &postgresMaterialRepository{db: db}
}

func (r *postgresMaterialRepository) Create(ctx context.Context, material *pgentities.Material) error {
	query := `
		INSERT INTO materials (
			id, school_id, uploaded_by_teacher_id, academic_unit_id,
			title, description, subject, grade, file_url, file_type,
			file_size_bytes, status, is_public, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.ExecContext(ctx, query,
		material.ID,
		material.SchoolID,
		material.UploadedByTeacherID,
		material.AcademicUnitID,
		material.Title,
		material.Description,
		material.Subject,
		material.Grade,
		material.FileURL,
		material.FileType,
		material.FileSizeBytes,
		material.Status,
		material.IsPublic,
		material.CreatedAt,
		material.UpdatedAt,
	)

	return err
}

func (r *postgresMaterialRepository) FindByID(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, error) {
	query := `
		SELECT id, school_id, uploaded_by_teacher_id, academic_unit_id,
		       title, description, subject, grade, file_url, file_type,
		       file_size_bytes, status, processing_started_at, processing_completed_at,
		       is_public, created_at, updated_at, deleted_at
		FROM materials
		WHERE id = $1 AND deleted_at IS NULL
	`

	var (
		materialID            uuid.UUID
		schoolID              uuid.UUID
		uploadedByTeacherID   uuid.UUID
		academicUnitID        sql.NullString
		title                 string
		description           sql.NullString
		subject               sql.NullString
		grade                 sql.NullString
		fileURL               string
		fileType              string
		fileSizeBytes         int64
		status                string
		processingStartedAt   sql.NullTime
		processingCompletedAt sql.NullTime
		isPublic              bool
		createdAt             time.Time
		updatedAt             time.Time
		deletedAt             sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id.UUID()).Scan(
		&materialID, &schoolID, &uploadedByTeacherID, &academicUnitID,
		&title, &description, &subject, &grade, &fileURL, &fileType,
		&fileSizeBytes, &status, &processingStartedAt, &processingCompletedAt,
		&isPublic, &createdAt, &updatedAt, &deletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	material := &pgentities.Material{
		ID:                  materialID,
		SchoolID:            schoolID,
		UploadedByTeacherID: uploadedByTeacherID,
		Title:               title,
		FileURL:             fileURL,
		FileType:            fileType,
		FileSizeBytes:       fileSizeBytes,
		Status:              status,
		IsPublic:            isPublic,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}

	if academicUnitID.Valid {
		unitID, _ := uuid.Parse(academicUnitID.String)
		material.AcademicUnitID = &unitID
	}

	if description.Valid {
		material.Description = &description.String
	}

	if subject.Valid {
		material.Subject = &subject.String
	}

	if grade.Valid {
		material.Grade = &grade.String
	}

	if processingStartedAt.Valid {
		material.ProcessingStartedAt = &processingStartedAt.Time
	}

	if processingCompletedAt.Valid {
		material.ProcessingCompletedAt = &processingCompletedAt.Time
	}

	if deletedAt.Valid {
		material.DeletedAt = &deletedAt.Time
	}

	return material, nil
}

func (r *postgresMaterialRepository) Update(ctx context.Context, material *pgentities.Material) error {
	query := `
		UPDATE materials
		SET title = $1, description = $2, subject = $3, grade = $4,
		    file_url = $5, status = $6, is_public = $7, updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.ExecContext(ctx, query,
		material.Title,
		material.Description,
		material.Subject,
		material.Grade,
		material.FileURL,
		material.Status,
		material.IsPublic,
		material.UpdatedAt,
		material.ID,
	)

	return err
}

func (r *postgresMaterialRepository) List(ctx context.Context, filters repository.ListFilters) ([]*pgentities.Material, error) {
	query := `
		SELECT id, school_id, uploaded_by_teacher_id, academic_unit_id,
		       title, description, subject, grade, file_url, file_type,
		       file_size_bytes, status, processing_started_at, processing_completed_at,
		       is_public, created_at, updated_at
		FROM materials
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argPos := 1

	if filters.Status != nil {
		query += ` AND status = $` + string(rune('0'+argPos))
		args = append(args, *filters.Status)
		// argPos no se incrementa porque no hay más filtros después
	}

	query += ` ORDER BY created_at DESC LIMIT 50`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var materials []*pgentities.Material
	for rows.Next() {
		var (
			materialID            uuid.UUID
			schoolID              uuid.UUID
			uploadedByTeacherID   uuid.UUID
			academicUnitID        sql.NullString
			title                 string
			description           sql.NullString
			subject               sql.NullString
			grade                 sql.NullString
			fileURL               string
			fileType              string
			fileSizeBytes         int64
			status                string
			processingStartedAt   sql.NullTime
			processingCompletedAt sql.NullTime
			isPublic              bool
			createdAt             time.Time
			updatedAt             time.Time
		)

		err := rows.Scan(
			&materialID, &schoolID, &uploadedByTeacherID, &academicUnitID,
			&title, &description, &subject, &grade, &fileURL, &fileType,
			&fileSizeBytes, &status, &processingStartedAt, &processingCompletedAt,
			&isPublic, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		material := &pgentities.Material{
			ID:                  materialID,
			SchoolID:            schoolID,
			UploadedByTeacherID: uploadedByTeacherID,
			Title:               title,
			FileURL:             fileURL,
			FileType:            fileType,
			FileSizeBytes:       fileSizeBytes,
			Status:              status,
			IsPublic:            isPublic,
			CreatedAt:           createdAt,
			UpdatedAt:           updatedAt,
		}

		if academicUnitID.Valid {
			unitID, _ := uuid.Parse(academicUnitID.String)
			material.AcademicUnitID = &unitID
		}

		if description.Valid {
			material.Description = &description.String
		}

		if subject.Valid {
			material.Subject = &subject.String
		}

		if grade.Valid {
			material.Grade = &grade.String
		}

		if processingStartedAt.Valid {
			material.ProcessingStartedAt = &processingStartedAt.Time
		}

		if processingCompletedAt.Valid {
			material.ProcessingCompletedAt = &processingCompletedAt.Time
		}

		materials = append(materials, material)
	}

	return materials, rows.Err()
}

func (r *postgresMaterialRepository) FindByAuthor(ctx context.Context, authorID valueobject.UserID) ([]*pgentities.Material, error) {
	query := `
		SELECT id, school_id, uploaded_by_teacher_id, academic_unit_id,
		       title, description, subject, grade, file_url, file_type,
		       file_size_bytes, status, is_public, created_at, updated_at
		FROM materials
		WHERE uploaded_by_teacher_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, authorID.UUID())
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var materials []*pgentities.Material
	for rows.Next() {
		var (
			materialID          uuid.UUID
			schoolID            uuid.UUID
			uploadedByTeacherID uuid.UUID
			academicUnitID      sql.NullString
			title               string
			description         sql.NullString
			subject             sql.NullString
			grade               sql.NullString
			fileURL             string
			fileType            string
			fileSizeBytes       int64
			status              string
			isPublic            bool
			createdAt           time.Time
			updatedAt           time.Time
		)

		err := rows.Scan(
			&materialID, &schoolID, &uploadedByTeacherID, &academicUnitID,
			&title, &description, &subject, &grade, &fileURL, &fileType,
			&fileSizeBytes, &status, &isPublic, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		material := &pgentities.Material{
			ID:                  materialID,
			SchoolID:            schoolID,
			UploadedByTeacherID: uploadedByTeacherID,
			Title:               title,
			FileURL:             fileURL,
			FileType:            fileType,
			FileSizeBytes:       fileSizeBytes,
			Status:              status,
			IsPublic:            isPublic,
			CreatedAt:           createdAt,
			UpdatedAt:           updatedAt,
		}

		if academicUnitID.Valid {
			unitID, _ := uuid.Parse(academicUnitID.String)
			material.AcademicUnitID = &unitID
		}

		if description.Valid {
			material.Description = &description.String
		}

		if subject.Valid {
			material.Subject = &subject.String
		}

		if grade.Valid {
			material.Grade = &grade.String
		}

		materials = append(materials, material)
	}

	return materials, rows.Err()
}

func (r *postgresMaterialRepository) UpdateStatus(ctx context.Context, id valueobject.MaterialID, status enum.MaterialStatus) error {
	query := `UPDATE materials SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status.String(), id.UUID())
	return err
}

func (r *postgresMaterialRepository) UpdateProcessingStatus(ctx context.Context, id valueobject.MaterialID, status enum.ProcessingStatus) error {
	query := `UPDATE materials SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status.String(), id.UUID())
	return err
}

func (r *postgresMaterialRepository) FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, []*pgentities.MaterialVersion, error) {
	// Por ahora solo retorna el material sin versiones
	// TODO: Implementar join con material_versions cuando se necesite
	material, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return material, nil, nil
}

func (r *postgresMaterialRepository) CountPublishedMaterials(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM materials WHERE status = 'published' AND deleted_at IS NULL`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
