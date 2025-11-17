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

// PostgresAttemptRepository implementa repositories.AttemptRepository para PostgreSQL
type PostgresAttemptRepository struct {
	db *sql.DB
}

// NewPostgresAttemptRepository crea una nueva instancia del repositorio
func NewPostgresAttemptRepository(db *sql.DB) repositories.AttemptRepository {
	return &PostgresAttemptRepository{db: db}
}

// FindByID busca un intento por ID
func (r *PostgresAttemptRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Attempt, error) {
	// 1. Query para el attempt
	attemptQuery := `
		SELECT id, assessment_id, student_id, score, max_score,
		       time_spent_seconds, started_at, completed_at, created_at,
		       idempotency_key
		FROM assessment_attempt
		WHERE id = $1
	`

	var (
		idStr           string
		assessmentIDStr string
		studentIDStr    string
		score           int
		maxScore        int
		timeSpent       int
		startedAt       time.Time
		completedAt     time.Time
		createdAt       time.Time
		idempotencyKey  sql.NullString
	)

	err := r.db.QueryRowContext(ctx, attemptQuery, id.String()).Scan(
		&idStr, &assessmentIDStr, &studentIDStr, &score, &maxScore,
		&timeSpent, &startedAt, &completedAt, &createdAt, &idempotencyKey,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding attempt: %w", err)
	}

	assessmentID, _ := uuid.Parse(assessmentIDStr)
	studentID, _ := uuid.Parse(studentIDStr)

	// 2. Query para las respuestas (answers)
	answersQuery := `
		SELECT id, attempt_id, question_id, selected_answer_id,
		       is_correct, time_spent_seconds, created_at
		FROM assessment_attempt_answer
		WHERE attempt_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, answersQuery, id.String())
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
			answerCreatedAt  time.Time
		)

		err := rows.Scan(
			&answerIDStr, &attemptIDStr, &questionID, &selectedAnswerID,
			&isCorrect, &timeSpentSecs, &answerCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
		}

		answerID, _ := uuid.Parse(answerIDStr)
		attemptID, _ := uuid.Parse(attemptIDStr)

		answer := &entities.Answer{
			ID:               answerID,
			AttemptID:        attemptID,
			QuestionID:       questionID,
			SelectedAnswerID: selectedAnswerID,
			IsCorrect:        isCorrect,
			TimeSpentSeconds: timeSpentSecs,
			CreatedAt:        answerCreatedAt,
		}

		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating answers: %w", err)
	}

	// Construir attempt con answers
	var idempotencyKeyPtr *string
	if idempotencyKey.Valid {
		idempotencyKeyPtr = &idempotencyKey.String
	}

	attempt := &entities.Attempt{
		ID:               id,
		AssessmentID:     assessmentID,
		StudentID:        studentID,
		Score:            score,
		MaxScore:         maxScore,
		TimeSpentSeconds: timeSpent,
		StartedAt:        startedAt,
		CompletedAt:      completedAt,
		CreatedAt:        createdAt,
		Answers:          answers,
		IdempotencyKey:   idempotencyKeyPtr,
	}

	return attempt, nil
}

// FindByStudentAndAssessment busca intentos de un estudiante en una evaluación
func (r *PostgresAttemptRepository) FindByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) ([]*entities.Attempt, error) {
	query := `
		SELECT id, assessment_id, student_id, score, max_score,
		       time_spent_seconds, started_at, completed_at, created_at,
		       idempotency_key
		FROM assessment_attempt
		WHERE student_id = $1 AND assessment_id = $2
		ORDER BY completed_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, studentID.String(), assessmentID.String())
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding attempts: %w", err)
	}
	defer rows.Close()

	var attempts []*entities.Attempt
	for rows.Next() {
		var (
			idStr           string
			assessmentIDStr string
			studentIDStr    string
			score           int
			maxScore        int
			timeSpent       int
			startedAt       time.Time
			completedAt     time.Time
			createdAt       time.Time
			idempotencyKey  sql.NullString
		)

		err := rows.Scan(
			&idStr, &assessmentIDStr, &studentIDStr, &score, &maxScore,
			&timeSpent, &startedAt, &completedAt, &createdAt, &idempotencyKey,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres: error scanning attempt: %w", err)
		}

		attemptID, _ := uuid.Parse(idStr)
		parsedAssessmentID, _ := uuid.Parse(assessmentIDStr)
		parsedStudentID, _ := uuid.Parse(studentIDStr)

		var idempotencyKeyPtr *string
		if idempotencyKey.Valid {
			idempotencyKeyPtr = &idempotencyKey.String
		}

		// Cargar las respuestas para este intento
		answersQuery := `
			SELECT id, attempt_id, question_id, selected_answer_id,
			       is_correct, time_spent_seconds, created_at
			FROM assessment_attempt_answer
			WHERE attempt_id = $1
			ORDER BY created_at ASC
		`

		answerRows, err := r.db.QueryContext(ctx, answersQuery, attemptID.String())
		if err != nil {
			return nil, fmt.Errorf("postgres: error finding answers: %w", err)
		}

		var answers []*entities.Answer
		for answerRows.Next() {
			var (
				answerIDStr      string
				attemptIDStr     string
				questionID       string
				selectedAnswerID string
				isCorrect        bool
				timeSpentSecs    int
				answerCreatedAt  time.Time
			)

			err := answerRows.Scan(
				&answerIDStr, &attemptIDStr, &questionID, &selectedAnswerID,
				&isCorrect, &timeSpentSecs, &answerCreatedAt,
			)
			if err != nil {
				answerRows.Close()
				return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
			}

			answerID, _ := uuid.Parse(answerIDStr)
			answerAttemptID, _ := uuid.Parse(attemptIDStr)

			answer := &entities.Answer{
				ID:               answerID,
				AttemptID:        answerAttemptID,
				QuestionID:       questionID,
				SelectedAnswerID: selectedAnswerID,
				IsCorrect:        isCorrect,
				TimeSpentSeconds: timeSpentSecs,
				CreatedAt:        answerCreatedAt,
			}

			answers = append(answers, answer)
		}
		answerRows.Close()

		attempt := &entities.Attempt{
			ID:               attemptID,
			AssessmentID:     parsedAssessmentID,
			StudentID:        parsedStudentID,
			Score:            score,
			MaxScore:         maxScore,
			TimeSpentSeconds: timeSpent,
			StartedAt:        startedAt,
			CompletedAt:      completedAt,
			CreatedAt:        createdAt,
			Answers:          answers,
			IdempotencyKey:   idempotencyKeyPtr,
		}

		attempts = append(attempts, attempt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating attempts: %w", err)
	}

	return attempts, nil
}

// Save guarda un intento (solo INSERT, no UPDATE - inmutable)
// IMPORTANTE: Debe guardar el attempt Y sus answers en una transacción atómica
func (r *PostgresAttemptRepository) Save(ctx context.Context, attempt *entities.Attempt) error {
	if attempt == nil {
		return fmt.Errorf("postgres: attempt cannot be nil")
	}

	if err := attempt.Validate(); err != nil {
		return fmt.Errorf("postgres: invalid attempt: %w", err)
	}

	if len(attempt.Answers) == 0 {
		return domainErrors.ErrNoAnswersProvided
	}

	// 1. Iniciar transacción
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres: error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Ignorar error si ya se hizo Commit

	// 2. INSERT del attempt
	attemptQuery := `
		INSERT INTO assessment_attempt (
			id, assessment_id, student_id, score, max_score,
			time_spent_seconds, started_at, completed_at, created_at,
			idempotency_key
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var idempotencyKey interface{}
	if attempt.IdempotencyKey != nil {
		idempotencyKey = *attempt.IdempotencyKey
	}

	_, err = tx.ExecContext(ctx, attemptQuery,
		attempt.ID,
		attempt.AssessmentID,
		attempt.StudentID,
		attempt.Score,
		attempt.MaxScore,
		attempt.TimeSpentSeconds,
		attempt.StartedAt,
		attempt.CompletedAt,
		attempt.CreatedAt,
		idempotencyKey,
	)

	if err != nil {
		return fmt.Errorf("postgres: error inserting attempt: %w", err)
	}

	// 3. INSERT de todas las answers (batch)
	answerQuery := `
		INSERT INTO assessment_attempt_answer (
			id, attempt_id, question_id, selected_answer_id,
			is_correct, time_spent_seconds, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, answer := range attempt.Answers {
		_, err = tx.ExecContext(ctx, answerQuery,
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

	// 4. Commit de la transacción
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("postgres: error committing transaction: %w", err)
	}

	return nil
}

// CountByStudentAndAssessment cuenta intentos de un estudiante
func (r *PostgresAttemptRepository) CountByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM assessment_attempt
		WHERE student_id = $1 AND assessment_id = $2
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, studentID.String(), assessmentID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("postgres: error counting attempts: %w", err)
	}

	return count, nil
}

// FindByStudent busca todos los intentos de un estudiante (historial)
func (r *PostgresAttemptRepository) FindByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*entities.Attempt, error) {
	query := `
		SELECT id, assessment_id, student_id, score, max_score,
		       time_spent_seconds, started_at, completed_at, created_at,
		       idempotency_key
		FROM assessment_attempt
		WHERE student_id = $1
		ORDER BY completed_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, studentID.String(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("postgres: error finding attempts: %w", err)
	}
	defer rows.Close()

	var attempts []*entities.Attempt
	for rows.Next() {
		var (
			idStr           string
			assessmentIDStr string
			studentIDStr    string
			score           int
			maxScore        int
			timeSpent       int
			startedAt       time.Time
			completedAt     time.Time
			createdAt       time.Time
			idempotencyKey  sql.NullString
		)

		err := rows.Scan(
			&idStr, &assessmentIDStr, &studentIDStr, &score, &maxScore,
			&timeSpent, &startedAt, &completedAt, &createdAt, &idempotencyKey,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres: error scanning attempt: %w", err)
		}

		attemptID, _ := uuid.Parse(idStr)
		assessmentID, _ := uuid.Parse(assessmentIDStr)
		parsedStudentID, _ := uuid.Parse(studentIDStr)

		var idempotencyKeyPtr *string
		if idempotencyKey.Valid {
			idempotencyKeyPtr = &idempotencyKey.String
		}

		// Cargar las respuestas para este intento
		answersQuery := `
			SELECT id, attempt_id, question_id, selected_answer_id,
			       is_correct, time_spent_seconds, created_at
			FROM assessment_attempt_answer
			WHERE attempt_id = $1
			ORDER BY created_at ASC
		`

		answerRows, err := r.db.QueryContext(ctx, answersQuery, attemptID.String())
		if err != nil {
			return nil, fmt.Errorf("postgres: error finding answers: %w", err)
		}

		var answers []*entities.Answer
		for answerRows.Next() {
			var (
				answerIDStr      string
				attemptIDStr     string
				questionID       string
				selectedAnswerID string
				isCorrect        bool
				timeSpentSecs    int
				answerCreatedAt  time.Time
			)

			err := answerRows.Scan(
				&answerIDStr, &attemptIDStr, &questionID, &selectedAnswerID,
				&isCorrect, &timeSpentSecs, &answerCreatedAt,
			)
			if err != nil {
				answerRows.Close()
				return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
			}

			answerID, _ := uuid.Parse(answerIDStr)
			answerAttemptID, _ := uuid.Parse(attemptIDStr)

			answer := &entities.Answer{
				ID:               answerID,
				AttemptID:        answerAttemptID,
				QuestionID:       questionID,
				SelectedAnswerID: selectedAnswerID,
				IsCorrect:        isCorrect,
				TimeSpentSeconds: timeSpentSecs,
				CreatedAt:        answerCreatedAt,
			}

			answers = append(answers, answer)
		}
		answerRows.Close()

		attempt := &entities.Attempt{
			ID:               attemptID,
			AssessmentID:     assessmentID,
			StudentID:        parsedStudentID,
			Score:            score,
			MaxScore:         maxScore,
			TimeSpentSeconds: timeSpent,
			StartedAt:        startedAt,
			CompletedAt:      completedAt,
			CreatedAt:        createdAt,
			Answers:          answers,
			IdempotencyKey:   idempotencyKeyPtr,
		}

		attempts = append(attempts, attempt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating attempts: %w", err)
	}

	return attempts, nil
}
