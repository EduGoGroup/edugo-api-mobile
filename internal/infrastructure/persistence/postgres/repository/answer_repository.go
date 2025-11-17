package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"
	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
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
func (r *PostgresAnswerRepository) FindByAttemptID(ctx context.Context, attemptID uuid.UUID) ([]*entities.Answer, error) {
	query := `
		SELECT id, attempt_id, question_id, selected_answer_id,
		       is_correct, time_spent_seconds, created_at
		FROM assessment_attempt_answer
		WHERE attempt_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, attemptID.String())
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding answers: %w", err)
	}
	defer rows.Close()

	var answers []*entities.Answer
	for rows.Next() {
		var (
			answerIDStr      string
			attemptIDStr     string
			questionID       string
			selectedAnswerID string
			isCorrect        bool
			timeSpentSecs    int
			createdAt        time.Time
		)

		err := rows.Scan(
			&answerIDStr, &attemptIDStr, &questionID, &selectedAnswerID,
			&isCorrect, &timeSpentSecs, &createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
		}

		answerID, _ := uuid.Parse(answerIDStr)
		attemptIDParsed, _ := uuid.Parse(attemptIDStr)

		answer := &entities.Answer{
			ID:               answerID,
			AttemptID:        attemptIDParsed,
			QuestionID:       questionID,
			SelectedAnswerID: selectedAnswerID,
			IsCorrect:        isCorrect,
			TimeSpentSeconds: timeSpentSecs,
			CreatedAt:        createdAt,
		}

		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating answers: %w", err)
	}

	return answers, nil
}

// Save guarda una o múltiples respuestas (batch insert)
func (r *PostgresAnswerRepository) Save(ctx context.Context, answers []*entities.Answer) error {
	if len(answers) == 0 {
		return fmt.Errorf("postgres: no answers to save")
	}

	// Validar cada answer
	for _, answer := range answers {
		if err := answer.Validate(); err != nil {
			return fmt.Errorf("postgres: invalid answer: %w", err)
		}
	}

	query := `
		INSERT INTO assessment_attempt_answer (
			id, attempt_id, question_id, selected_answer_id,
			is_correct, time_spent_seconds, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// Usar transacción para batch insert
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres: error starting transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("postgres: error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, answer := range answers {
		_, err := stmt.ExecContext(ctx,
			answer.ID,
			answer.AttemptID,
			answer.QuestionID,
			answer.SelectedAnswerID,
			answer.IsCorrect,
			answer.TimeSpentSeconds,
			answer.CreatedAt,
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

// FindByQuestionID busca todas las respuestas para una pregunta específica
// Útil para analytics: identificar preguntas difíciles
func (r *PostgresAnswerRepository) FindByQuestionID(ctx context.Context, questionID string, limit, offset int) ([]*entities.Answer, error) {
	if questionID == "" {
		return []*entities.Answer{}, domainErrors.ErrInvalidQuestionID
	}

	query := `
		SELECT id, attempt_id, question_id, selected_answer_id,
		       is_correct, time_spent_seconds, created_at
		FROM assessment_attempt_answer
		WHERE question_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, questionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding answers: %w", err)
	}
	defer rows.Close()

	var answers []*entities.Answer
	for rows.Next() {
		var (
			answerIDStr      string
			attemptIDStr     string
			qID              string
			selectedAnswerID string
			isCorrect        bool
			timeSpentSecs    int
			createdAt        time.Time
		)

		err := rows.Scan(
			&answerIDStr, &attemptIDStr, &qID, &selectedAnswerID,
			&isCorrect, &timeSpentSecs, &createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
		}

		answerID, _ := uuid.Parse(answerIDStr)
		attemptID, _ := uuid.Parse(attemptIDStr)

		answer := &entities.Answer{
			ID:               answerID,
			AttemptID:        attemptID,
			QuestionID:       qID,
			SelectedAnswerID: selectedAnswerID,
			IsCorrect:        isCorrect,
			TimeSpentSeconds: timeSpentSecs,
			CreatedAt:        createdAt,
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
