package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// PostgresAssessmentRepository implementa repositories.AssessmentRepository para PostgreSQL
type PostgresAssessmentRepository struct {
	db *sql.DB
}

// NewPostgresAssessmentRepository crea una nueva instancia del repositorio
func NewPostgresAssessmentRepository(db *sql.DB) repositories.AssessmentRepository {
	return &PostgresAssessmentRepository{db: db}
}

// FindByID busca una evaluación por ID
func (r *PostgresAssessmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*pgentities.Assessment, error) {
	query := `
		SELECT id, material_id, mongo_document_id, questions_count, total_questions,
		       title, pass_threshold, max_attempts, time_limit_minutes, status,
		       created_at, updated_at
		FROM assessment
		WHERE id = $1
	`

	var (
		idStr          string
		materialIDStr  string
		mongoDocID     string
		questionsCount int
		totalQuestions sql.NullInt32
		title          sql.NullString
		passThreshold  sql.NullInt32
		maxAttempts    sql.NullInt32
		timeLimitMins  sql.NullInt32
		status         string
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &materialIDStr, &mongoDocID, &questionsCount, &totalQuestions,
		&title, &passThreshold, &maxAttempts, &timeLimitMins, &status,
		&createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No encontrado
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding assessment: %w", err)
	}

	materialID, _ := uuid.Parse(materialIDStr)

	assessment := &pgentities.Assessment{
		ID:              id,
		MaterialID:      materialID,
		MongoDocumentID: mongoDocID,
		QuestionsCount:  questionsCount,
		Status:          status,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}

	if totalQuestions.Valid {
		val := int(totalQuestions.Int32)
		assessment.TotalQuestions = &val
	}

	if title.Valid {
		assessment.Title = &title.String
	}

	if passThreshold.Valid {
		val := int(passThreshold.Int32)
		assessment.PassThreshold = &val
	}

	if maxAttempts.Valid {
		val := int(maxAttempts.Int32)
		assessment.MaxAttempts = &val
	}

	if timeLimitMins.Valid {
		val := int(timeLimitMins.Int32)
		assessment.TimeLimitMinutes = &val
	}

	return assessment, nil
}

// FindByMaterialID busca una evaluación por material ID
func (r *PostgresAssessmentRepository) FindByMaterialID(ctx context.Context, materialID uuid.UUID) (*pgentities.Assessment, error) {
	query := `
		SELECT id, material_id, mongo_document_id, questions_count, total_questions,
		       title, pass_threshold, max_attempts, time_limit_minutes, status,
		       created_at, updated_at
		FROM assessment
		WHERE material_id = $1
	`

	var (
		idStr          string
		materialIDStr  string
		mongoDocID     string
		questionsCount int
		totalQuestions sql.NullInt32
		title          sql.NullString
		passThreshold  sql.NullInt32
		maxAttempts    sql.NullInt32
		timeLimitMins  sql.NullInt32
		status         string
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := r.db.QueryRowContext(ctx, query, materialID.String()).Scan(
		&idStr, &materialIDStr, &mongoDocID, &questionsCount, &totalQuestions,
		&title, &passThreshold, &maxAttempts, &timeLimitMins, &status,
		&createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No encontrado
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding assessment by material: %w", err)
	}

	id, _ := uuid.Parse(idStr)
	parsedMaterialID, _ := uuid.Parse(materialIDStr)

	assessment := &pgentities.Assessment{
		ID:              id,
		MaterialID:      parsedMaterialID,
		MongoDocumentID: mongoDocID,
		QuestionsCount:  questionsCount,
		Status:          status,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}

	if totalQuestions.Valid {
		val := int(totalQuestions.Int32)
		assessment.TotalQuestions = &val
	}

	if title.Valid {
		assessment.Title = &title.String
	}

	if passThreshold.Valid {
		val := int(passThreshold.Int32)
		assessment.PassThreshold = &val
	}

	if maxAttempts.Valid {
		val := int(maxAttempts.Int32)
		assessment.MaxAttempts = &val
	}

	if timeLimitMins.Valid {
		val := int(timeLimitMins.Int32)
		assessment.TimeLimitMinutes = &val
	}

	return assessment, nil
}

// Save guarda una evaluación (INSERT o UPDATE)
func (r *PostgresAssessmentRepository) Save(ctx context.Context, assessment *pgentities.Assessment) error {
	if assessment == nil {
		return fmt.Errorf("postgres: assessment cannot be nil")
	}

	query := `
		INSERT INTO assessment (
			id, material_id, mongo_document_id, questions_count, total_questions,
			title, pass_threshold, max_attempts, time_limit_minutes, status,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			questions_count = EXCLUDED.questions_count,
			total_questions = EXCLUDED.total_questions,
			title = EXCLUDED.title,
			pass_threshold = EXCLUDED.pass_threshold,
			max_attempts = EXCLUDED.max_attempts,
			time_limit_minutes = EXCLUDED.time_limit_minutes,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`

	var totalQuestions, title, passThreshold, maxAttempts, timeLimitMins interface{}

	if assessment.TotalQuestions != nil {
		totalQuestions = *assessment.TotalQuestions
	}
	if assessment.Title != nil {
		title = *assessment.Title
	}
	if assessment.PassThreshold != nil {
		passThreshold = *assessment.PassThreshold
	}
	if assessment.MaxAttempts != nil {
		maxAttempts = *assessment.MaxAttempts
	}
	if assessment.TimeLimitMinutes != nil {
		timeLimitMins = *assessment.TimeLimitMinutes
	}

	_, err := r.db.ExecContext(ctx, query,
		assessment.ID,
		assessment.MaterialID,
		assessment.MongoDocumentID,
		assessment.QuestionsCount,
		totalQuestions,
		title,
		passThreshold,
		maxAttempts,
		timeLimitMins,
		assessment.Status,
		assessment.CreatedAt,
		assessment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("postgres: error saving assessment: %w", err)
	}

	return nil
}

// Delete elimina una evaluación
func (r *PostgresAssessmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM assessment WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres: error deleting assessment: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("postgres: assessment not found")
	}

	return nil
}
