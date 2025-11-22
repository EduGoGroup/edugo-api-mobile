package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// PostgresAnswerRepository implementa repositories.AnswerRepository para PostgreSQL
type PostgresAnswerRepository struct {
	db *sql.DB
}

// NewPostgresAnswerRepository crea una nueva instancia del repositorio
func NewPostgresAnswerRepository(db *sql.DB) repositories.AnswerRepository {
	return &PostgresAnswerRepository{db: db}
}

// FindByAttemptID busca todas las respuestas de un intento
func (r *PostgresAnswerRepository) FindByAttemptID(ctx context.Context, attemptID uuid.UUID) ([]*pgentities.AssessmentAttemptAnswer, error) {
	query := `
		SELECT id, attempt_id, question_index, student_answer,
		       is_correct, points_earned, max_points, time_spent_seconds,
		       answered_at, created_at, updated_at
		FROM assessment_attempt_answer
		WHERE attempt_id = $1
		ORDER BY question_index ASC
	`

	rows, err := r.db.QueryContext(ctx, query, attemptID.String())
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding answers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var answers []*pgentities.AssessmentAttemptAnswer
	for rows.Next() {
		var (
			answerIDStr   string
			attemptIDStr  string
			questionIndex int
			studentAnswer *string
			isCorrect     *bool
			pointsEarned  *float64
			maxPoints     *float64
			timeSpentSecs *int
			answeredAt    time.Time
			createdAt     time.Time
			updatedAt     time.Time
		)

		err := rows.Scan(
			&answerIDStr, &attemptIDStr, &questionIndex, &studentAnswer,
			&isCorrect, &pointsEarned, &maxPoints, &timeSpentSecs,
			&answeredAt, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
		}

		answerID, _ := uuid.Parse(answerIDStr)
		attemptIDParsed, _ := uuid.Parse(attemptIDStr)

		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               answerID,
			AttemptID:        attemptIDParsed,
			QuestionIndex:    questionIndex,
			StudentAnswer:    studentAnswer,
			IsCorrect:        isCorrect,
			PointsEarned:     pointsEarned,
			MaxPoints:        maxPoints,
			TimeSpentSeconds: timeSpentSecs,
			AnsweredAt:       answeredAt,
			CreatedAt:        createdAt,
			UpdatedAt:        updatedAt,
		}

		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating answers: %w", err)
	}

	return answers, nil
}

// Save guarda una o múltiples respuestas (batch insert)
func (r *PostgresAnswerRepository) Save(ctx context.Context, answers []*pgentities.AssessmentAttemptAnswer) error {
	if len(answers) == 0 {
		return fmt.Errorf("postgres: no answers to save")
	}

	query := `
		INSERT INTO assessment_attempt_answer (
			id, attempt_id, question_index, student_answer,
			is_correct, points_earned, max_points, time_spent_seconds,
			answered_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	// Usar transacción para batch insert
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres: error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Ignorar error si ya se hizo Commit

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("postgres: error preparing statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, answer := range answers {
		_, err := stmt.ExecContext(ctx,
			answer.ID,
			answer.AttemptID,
			answer.QuestionIndex,
			answer.StudentAnswer,
			answer.IsCorrect,
			answer.PointsEarned,
			answer.MaxPoints,
			answer.TimeSpentSeconds,
			answer.AnsweredAt,
			answer.CreatedAt,
			answer.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("postgres: error inserting answer: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("postgres: error committing transaction: %w", err)
	}

	return nil
}

// FindByQuestionID busca todas las respuestas para un índice de pregunta específico
// Útil para analytics: identificar preguntas difíciles
func (r *PostgresAnswerRepository) FindByQuestionID(ctx context.Context, questionID string, limit, offset int) ([]*pgentities.AssessmentAttemptAnswer, error) {
	if questionID == "" {
		return []*pgentities.AssessmentAttemptAnswer{}, domainErrors.ErrInvalidQuestionID
	}

	// Nota: Esta función debería recibir questionIndex (int) en lugar de questionID (string)
	// Por ahora mantenemos la firma de la interfaz pero esto es un code smell
	query := `
		SELECT id, attempt_id, question_index, student_answer,
		       is_correct, points_earned, max_points, time_spent_seconds,
		       answered_at, created_at, updated_at
		FROM assessment_attempt_answer
		WHERE question_index::text = $1
		ORDER BY answered_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, questionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding answers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var answers []*pgentities.AssessmentAttemptAnswer
	for rows.Next() {
		var (
			answerIDStr   string
			attemptIDStr  string
			questionIndex int
			studentAnswer *string
			isCorrect     *bool
			pointsEarned  *float64
			maxPoints     *float64
			timeSpentSecs *int
			answeredAt    time.Time
			createdAt     time.Time
			updatedAt     time.Time
		)

		err := rows.Scan(
			&answerIDStr, &attemptIDStr, &questionIndex, &studentAnswer,
			&isCorrect, &pointsEarned, &maxPoints, &timeSpentSecs,
			&answeredAt, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
		}

		answerID, _ := uuid.Parse(answerIDStr)
		attemptID, _ := uuid.Parse(attemptIDStr)

		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               answerID,
			AttemptID:        attemptID,
			QuestionIndex:    questionIndex,
			StudentAnswer:    studentAnswer,
			IsCorrect:        isCorrect,
			PointsEarned:     pointsEarned,
			MaxPoints:        maxPoints,
			TimeSpentSeconds: timeSpentSecs,
			AnsweredAt:       answeredAt,
			CreatedAt:        createdAt,
			UpdatedAt:        updatedAt,
		}

		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating answers: %w", err)
	}

	return answers, nil
}

// Analytics: GetQuestionDifficultyStats obtiene estadísticas de dificultad de una pregunta
// Esta es una función helper para analytics (no en la interfaz del repositorio)
func (r *PostgresAnswerRepository) GetQuestionDifficultyStats(ctx context.Context, questionID string) (totalAnswers int, correctAnswers int, errorRate float64, err error) {
	if questionID == "" {
		return 0, 0, 0.0, domainErrors.ErrInvalidQuestionID
	}

	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE is_correct = true) as correct
		FROM assessment_attempt_answer
		WHERE question_id = $1
	`

	var total, correct int
	err = r.db.QueryRowContext(ctx, query, questionID).Scan(&total, &correct)
	if err != nil {
		return 0, 0, 0.0, fmt.Errorf("postgres: error calculating stats: %w", err)
	}

	var errorRateCalc float64
	if total > 0 {
		errorRateCalc = float64(total-correct) / float64(total)
	}

	return total, correct, errorRateCalc, nil
}
