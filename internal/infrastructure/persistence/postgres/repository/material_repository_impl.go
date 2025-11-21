package repository

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

type postgresMaterialRepository struct {
	db *sql.DB
}

func NewPostgresMaterialRepository(db *sql.DB) repository.MaterialRepository {
	return &postgresMaterialRepository{db: db}
}

func (r *postgresMaterialRepository) Create(ctx context.Context, material *entity.Material) error {
	query := `
		INSERT INTO materials (
			id, title, description, author_id, subject_id,
			status, processing_status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		material.ID().String(),
		material.Title(),
		material.Description(),
		material.AuthorID().String(),
		material.SubjectID(),
		material.Status().String(),
		material.ProcessingStatus().String(),
		material.CreatedAt(),
		material.UpdatedAt(),
	)

	return err
}

func (r *postgresMaterialRepository) FindByID(ctx context.Context, id valueobject.MaterialID) (*entity.Material, error) {
	query := `
		SELECT id, title, description, author_id, subject_id, s3_key, s3_url,
		       status, processing_status, created_at, updated_at
		FROM materials
		WHERE id = $1 AND is_deleted = false
	`

	var (
		idStr            string
		title            string
		description      string
		authorIDStr      string
		subjectID        sql.NullString
		s3Key            sql.NullString
		s3URL            sql.NullString
		statusStr        string
		processingStatus string
		createdAt        sql.NullTime
		updatedAt        sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &title, &description, &authorIDStr, &subjectID, &s3Key, &s3URL,
		&statusStr, &processingStatus, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	materialID, _ := valueobject.MaterialIDFromString(idStr)
	authorID, _ := valueobject.UserIDFromString(authorIDStr)

	return entity.ReconstructMaterial(
		materialID,
		title,
		description,
		authorID,
		subjectID.String,
		s3Key.String,
		s3URL.String,
		enum.MaterialStatus(statusStr),
		enum.ProcessingStatus(processingStatus),
		createdAt.Time,
		updatedAt.Time,
	), nil
}

func (r *postgresMaterialRepository) Update(ctx context.Context, material *entity.Material) error {
	query := `
		UPDATE materials
		SET title = $1, description = $2, s3_key = $3, s3_url = $4,
		    status = $5, processing_status = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(ctx, query,
		material.Title(),
		material.Description(),
		material.S3Key(),
		material.S3URL(),
		material.Status().String(),
		material.ProcessingStatus().String(),
		material.UpdatedAt(),
		material.ID().String(),
	)

	return err
}

func (r *postgresMaterialRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entity.Material, error) {
	query := `
		SELECT id, title, description, author_id, subject_id, s3_key, s3_url,
		       status, processing_status, created_at, updated_at
		FROM materials
		WHERE is_deleted = false
	`

	args := []interface{}{}

	if filters.Status != nil {
		query += ` AND status = $1`
		args = append(args, filters.Status.String())
	}

	query += ` ORDER BY created_at DESC LIMIT 50`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var materials []*entity.Material
	for rows.Next() {
		var (
			idStr            string
			title            string
			description      string
			authorIDStr      string
			subjectID        sql.NullString
			s3Key            sql.NullString
			s3URL            sql.NullString
			statusStr        string
			processingStatus string
			createdAt        sql.NullTime
			updatedAt        sql.NullTime
		)

		err := rows.Scan(
			&idStr, &title, &description, &authorIDStr, &subjectID, &s3Key, &s3URL,
			&statusStr, &processingStatus, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		materialID, _ := valueobject.MaterialIDFromString(idStr)
		authorID, _ := valueobject.UserIDFromString(authorIDStr)

		material := entity.ReconstructMaterial(
			materialID,
			title,
			description,
			authorID,
			subjectID.String,
			s3Key.String,
			s3URL.String,
			enum.MaterialStatus(statusStr),
			enum.ProcessingStatus(processingStatus),
			createdAt.Time,
			updatedAt.Time,
		)

		materials = append(materials, material)
	}

	return materials, rows.Err()
}

func (r *postgresMaterialRepository) FindByAuthor(ctx context.Context, authorID valueobject.UserID) ([]*entity.Material, error) {
	// Similar a List pero filtrando por author_id
	// Implementación completa se puede agregar después
	return nil, nil
}

func (r *postgresMaterialRepository) UpdateStatus(ctx context.Context, id valueobject.MaterialID, status enum.MaterialStatus) error {
	query := `UPDATE materials SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status.String(), id.String())
	return err
}

func (r *postgresMaterialRepository) UpdateProcessingStatus(ctx context.Context, id valueobject.MaterialID, status enum.ProcessingStatus) error {
	query := `UPDATE materials SET processing_status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status.String(), id.String())
	return err
}

// FindByIDWithVersions busca un material por ID incluyendo su historial completo de versiones
// Ejecuta un LEFT JOIN con material_versions para obtener todas las versiones en una sola query
func (r *postgresMaterialRepository) FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*entity.Material, []*entity.MaterialVersion, error) {
	// Query SQL con LEFT JOIN para obtener material + versiones
	// Ordena versiones por version_number DESC para mostrar más recientes primero
	query := `
		SELECT
			m.id, m.title, m.description, m.author_id, m.subject_id, m.s3_key, m.s3_url,
			m.status, m.processing_status, m.created_at, m.updated_at,
			v.id as version_id, v.version_number, v.title as version_title,
			v.content_url as version_url, v.changed_by, v.created_at as version_created_at
		FROM materials m
		LEFT JOIN material_versions v ON m.id = v.material_id
		WHERE m.id = $1 AND m.is_deleted = false
		ORDER BY v.version_number DESC
	`

	rows, err := r.db.QueryContext(ctx, query, id.String())
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var (
		material *entity.Material
		versions []*entity.MaterialVersion
	)

	// Iterar sobre resultados
	for rows.Next() {
		// Variables para material (se repite en cada fila, pero solo lo usamos una vez)
		var (
			materialIDStr    string
			title            string
			description      string
			authorIDStr      string
			subjectID        sql.NullString
			s3Key            sql.NullString
			s3URL            sql.NullString
			statusStr        string
			processingStatus string
			createdAt        sql.NullTime
			updatedAt        sql.NullTime
		)

		// Variables para versión (puede ser NULL si material no tiene versiones)
		var (
			versionIDStr     sql.NullString
			versionNumber    sql.NullInt32
			versionTitle     sql.NullString
			versionURL       sql.NullString
			versionChangedBy sql.NullString
			versionCreatedAt sql.NullTime
		)

		// Scan de la fila
		err := rows.Scan(
			&materialIDStr, &title, &description, &authorIDStr, &subjectID, &s3Key, &s3URL,
			&statusStr, &processingStatus, &createdAt, &updatedAt,
			&versionIDStr, &versionNumber, &versionTitle, &versionURL,
			&versionChangedBy, &versionCreatedAt,
		)
		if err != nil {
			return nil, nil, err
		}

		// Construir entidad Material solo una vez (en la primera iteración)
		if material == nil {
			matID, _ := valueobject.MaterialIDFromString(materialIDStr)
			authID, _ := valueobject.UserIDFromString(authorIDStr)

			material = entity.ReconstructMaterial(
				matID,
				title,
				description,
				authID,
				subjectID.String,
				s3Key.String,
				s3URL.String,
				enum.MaterialStatus(statusStr),
				enum.ProcessingStatus(processingStatus),
				createdAt.Time,
				updatedAt.Time,
			)
		}

		// Si existe versión en esta fila, agregarla al array
		// (Puede no existir si material no tiene versiones, LEFT JOIN retorna NULL)
		if versionIDStr.Valid {
			verID, _ := valueobject.MaterialVersionIDFromString(versionIDStr.String)
			matID, _ := valueobject.MaterialIDFromString(materialIDStr)
			changedByID, _ := valueobject.UserIDFromString(versionChangedBy.String)

			version := entity.ReconstructMaterialVersion(
				verID,
				matID,
				int(versionNumber.Int32),
				versionTitle.String,
				versionURL.String,
				changedByID,
				versionCreatedAt.Time,
			)

			versions = append(versions, version)
		}
	}

	// Verificar errores durante la iteración
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	// Si no se encontró el material, retornar nil (sin error)
	if material == nil {
		return nil, nil, nil
	}

	// Retornar material con su array de versiones (puede ser vacío)
	return material, versions, nil
}

// CountPublishedMaterials cuenta el total de materiales publicados
// Usado para estadísticas globales del sistema
func (r *postgresMaterialRepository) CountPublishedMaterials(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM materials WHERE status = $1 AND is_deleted = false`

	var count int64
	err := r.db.QueryRowContext(ctx, query, enum.MaterialStatusPublished.String()).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
