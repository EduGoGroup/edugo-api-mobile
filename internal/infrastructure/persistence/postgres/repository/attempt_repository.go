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
// TODO: STUB - Conectar con base de datos real + JOIN con answers (claude-local)
func (r *PostgresAttemptRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Attempt, error) {
	// STUB: Retornar mock data
	if id == uuid.Nil {
		return nil, domainErrors.ErrInvalidAttemptID
	}

	// Mock attempt con respuestas
	answer1, _ := entities.NewAnswer(id, "q1", "a", true, 60)
	answer2, _ := entities.NewAnswer(id, "q2", "b", false, 45)
	answer3, _ := entities.NewAnswer(id, "q3", "c", true, 30)

	mockAttempt := &entities.Attempt{
		ID:               id,
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            67, // 2/3 correctas
		MaxScore:         100,
		TimeSpentSeconds: 135,
		StartedAt:        time.Now().UTC().Add(-5 * time.Minute),
		CompletedAt:      time.Now().UTC(),
		CreatedAt:        time.Now().UTC(),
		Answers:          []*entities.Answer{answer1, answer2, answer3},
		IdempotencyKey:   nil,
	}

	return mockAttempt, nil

	// Query SQL real (con JOIN):
	/*
		// 1. Query para el attempt
		attemptQuery := `
			SELECT id, assessment_id, student_id, score, max_score,
			       time_spent_seconds, started_at, completed_at, created_at,
			       idempotency_key
			FROM assessment_attempt
			WHERE id = $1
		`

		var (
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
			&assessmentIDStr, &studentIDStr, &score, &maxScore,
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
	*/
}

// FindByStudentAndAssessment busca intentos de un estudiante en una evaluación
// TODO: STUB - Conectar con base de datos real (claude-local)
func (r *PostgresAttemptRepository) FindByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) ([]*entities.Attempt, error) {
	// STUB: Retornar array vacío
	if studentID == uuid.Nil || assessmentID == uuid.Nil {
		return []*entities.Attempt{}, nil
	}

	// Mock: Retornar un intento
	attemptID := uuid.New()
	answer1, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)

	mockAttempt := &entities.Attempt{
		ID:               attemptID,
		AssessmentID:     assessmentID,
		StudentID:        studentID,
		Score:            100,
		MaxScore:         100,
		TimeSpentSeconds: 30,
		StartedAt:        time.Now().UTC().Add(-1 * time.Minute),
		CompletedAt:      time.Now().UTC(),
		CreatedAt:        time.Now().UTC(),
		Answers:          []*entities.Answer{answer1},
		IdempotencyKey:   nil,
	}

	return []*entities.Attempt{mockAttempt}, nil

	// Query SQL real:
	/*
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
			// Scan y cargar cada attempt con sus answers...
			// Similar a FindByID
		}

		return attempts, nil
	*/
}

// Save guarda un intento (solo INSERT, no UPDATE - inmutable)
// IMPORTANTE: Debe guardar el attempt Y sus answers en una transacción atómica
// TODO: STUB - Conectar con base de datos real + TRANSACCIÓN (claude-local)
func (r *PostgresAttemptRepository) Save(ctx context.Context, attempt *entities.Attempt) error {
	// STUB: Validar y simular guardado exitoso
	if attempt == nil {
		return fmt.Errorf("postgres: attempt cannot be nil")
	}

	if err := attempt.Validate(); err != nil {
		return fmt.Errorf("postgres: invalid attempt: %w", err)
	}

	if len(attempt.Answers) == 0 {
		return domainErrors.ErrNoAnswersProvided
	}

	// Mock: Simulamos guardado exitoso
	return nil

	// Query SQL real (con TRANSACCIÓN):
	/*
		// 1. Iniciar transacción
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("postgres: error starting transaction: %w", err)
		}
		defer tx.Rollback() // Rollback si no se hace Commit

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
	*/
}

// CountByStudentAndAssessment cuenta intentos de un estudiante
// TODO: STUB - Conectar con base de datos real (claude-local)
func (r *PostgresAttemptRepository) CountByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (int, error) {
	// STUB: Retornar mock count
	if studentID == uuid.Nil || assessmentID == uuid.Nil {
		return 0, nil
	}

	// Mock: Simulamos que tiene 1 intento
	return 1, nil

	// Query SQL real:
	/*
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
	*/
}

// FindByStudent busca todos los intentos de un estudiante (historial)
// TODO: STUB - Conectar con base de datos real + PAGINACIÓN (claude-local)
func (r *PostgresAttemptRepository) FindByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*entities.Attempt, error) {
	// STUB: Retornar array vacío
	if studentID == uuid.Nil {
		return []*entities.Attempt{}, nil
	}

	// Mock: Retornar un intento
	attemptID := uuid.New()
	answer1, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)

	mockAttempt := &entities.Attempt{
		ID:               attemptID,
		AssessmentID:     uuid.New(),
		StudentID:        studentID,
		Score:            100,
		MaxScore:         100,
		TimeSpentSeconds: 30,
		StartedAt:        time.Now().UTC().Add(-1 * time.Hour),
		CompletedAt:      time.Now().UTC().Add(-59 * time.Minute),
		CreatedAt:        time.Now().UTC().Add(-59 * time.Minute),
		Answers:          []*entities.Answer{answer1},
		IdempotencyKey:   nil,
	}

	return []*entities.Attempt{mockAttempt}, nil

	// Query SQL real (con paginación):
	/*
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
			// Scan y cargar cada attempt...
			// Para cada attempt, cargar sus answers (sub-query o batch)
		}

		return attempts, nil
	*/
}
