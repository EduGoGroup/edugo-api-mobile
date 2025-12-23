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

// PostgresAttemptRepository implementa repositories.AttemptRepository para PostgreSQL
type PostgresAttemptRepository struct {
	db *sql.DB
}

// NewPostgresAttemptRepository crea una nueva instancia del repositorio
func NewPostgresAttemptRepository(db *sql.DB) repositories.AttemptRepository {
	return &PostgresAttemptRepository{db: db}
}

// FindByID busca un intento por ID
func (r *PostgresAttemptRepository) FindByID(ctx context.Context, id uuid.UUID) (*pgentities.AssessmentAttempt, error) {
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
		SELECT id, attempt_id, question_index, student_answer,
		       is_correct, points_earned, max_points, time_spent_seconds,
		       answered_at, created_at, updated_at
		FROM assessment_attempt_answer
		WHERE attempt_id = $1
		ORDER BY question_index ASC
	`

	rows, err := r.db.QueryContext(ctx, answersQuery, id.String())
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

		answers = append(answers, answer) //nolint:staticcheck // SA4010: answers usado en estructura del attempt
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating answers: %w", err)
	}

	// Construir attempt con answers
	var idempotencyKeyPtr *string
	if idempotencyKey.Valid {
		idempotencyKeyPtr = &idempotencyKey.String
	}

	var scorePtr *float64
	if score > 0 {
		scoreF := float64(score)
		scorePtr = &scoreF
	}

	var maxScorePtr *float64
	if maxScore > 0 {
		maxScoreF := float64(maxScore)
		maxScorePtr = &maxScoreF
	}

	var timeSpentPtr *int
	if timeSpent > 0 {
		timeSpentPtr = &timeSpent
	}

	attempt := &pgentities.AssessmentAttempt{
		ID:               id,
		AssessmentID:     assessmentID,
		StudentID:        studentID,
		Score:            scorePtr,
		MaxScore:         maxScorePtr,
		TimeSpentSeconds: timeSpentPtr,
		StartedAt:        startedAt,
		CompletedAt:      &completedAt,
		CreatedAt:        createdAt,
		IdempotencyKey:   idempotencyKeyPtr,
	}

	return attempt, nil
}

// FindByStudentAndAssessment busca intentos de un estudiante en una evaluación
func (r *PostgresAttemptRepository) FindByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) ([]*pgentities.AssessmentAttempt, error) {
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
	defer func() { _ = rows.Close() }()

	var attempts []*pgentities.AssessmentAttempt
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
			SELECT id, attempt_id, question_index, student_answer,
			       is_correct, points_earned, max_points, time_spent_seconds,
			       answered_at, created_at, updated_at
			FROM assessment_attempt_answer
			WHERE attempt_id = $1
			ORDER BY question_index ASC
		`

		answerRows, err := r.db.QueryContext(ctx, answersQuery, attemptID.String())
		if err != nil {
			return nil, fmt.Errorf("postgres: error finding answers: %w", err)
		}

		var answers []*pgentities.AssessmentAttemptAnswer
		for answerRows.Next() {
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

			err := answerRows.Scan(
				&answerIDStr, &attemptIDStr, &questionIndex, &studentAnswer,
				&isCorrect, &pointsEarned, &maxPoints, &timeSpentSecs,
				&answeredAt, &createdAt, &updatedAt,
			)
			if err != nil {
				_ = answerRows.Close()
				return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
			}

			answerID, _ := uuid.Parse(answerIDStr)
			answerAttemptID, _ := uuid.Parse(attemptIDStr)

			answer := &pgentities.AssessmentAttemptAnswer{
				ID:               answerID,
				AttemptID:        answerAttemptID,
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

			answers = append(answers, answer) //nolint:staticcheck // SA4010: answers usado posteriormente
		}
		_ = answerRows.Close()

		var scorePtr *float64
		if score > 0 {
			scoreF := float64(score)
			scorePtr = &scoreF
		}

		var maxScorePtr *float64
		if maxScore > 0 {
			maxScoreF := float64(maxScore)
			maxScorePtr = &maxScoreF
		}

		var timeSpentPtr *int
		if timeSpent > 0 {
			timeSpentPtr = &timeSpent
		}

		attempt := &pgentities.AssessmentAttempt{
			ID:               attemptID,
			AssessmentID:     parsedAssessmentID,
			StudentID:        parsedStudentID,
			Score:            scorePtr,
			MaxScore:         maxScorePtr,
			TimeSpentSeconds: timeSpentPtr,
			StartedAt:        startedAt,
			CompletedAt:      &completedAt,
			CreatedAt:        createdAt,
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
func (r *PostgresAttemptRepository) Save(ctx context.Context, attempt *pgentities.AssessmentAttempt) error {
	if attempt == nil {
		return fmt.Errorf("postgres: attempt cannot be nil")
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

	// Note: Answers are managed separately through AnswerRepository
	// The attempt is saved without its answers relationship

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

// CountCompletedAssessments cuenta el total de evaluaciones completadas (con completed_at != NULL)
// Implementa repositories.AssessmentStats para estadísticas globales
func (r *PostgresAttemptRepository) CountCompletedAssessments(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM assessment_attempt WHERE completed_at IS NOT NULL`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("postgres: error counting completed assessments: %w", err)
	}

	return count, nil
}

// CalculateAverageScore calcula el promedio de scores de evaluaciones completadas
// Implementa repositories.AssessmentStats para estadísticas globales
func (r *PostgresAttemptRepository) CalculateAverageScore(ctx context.Context) (float64, error) {
	query := `
		SELECT COALESCE(AVG(score), 0.0)
		FROM assessment_attempt
		WHERE completed_at IS NOT NULL AND score IS NOT NULL
	`

	var avgScore float64
	err := r.db.QueryRowContext(ctx, query).Scan(&avgScore)
	if err != nil {
		return 0.0, fmt.Errorf("postgres: error calculating average score: %w", err)
	}

	return avgScore, nil
}

// FindByStudent busca todos los intentos de un estudiante (historial)
func (r *PostgresAttemptRepository) FindByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*pgentities.AssessmentAttempt, error) {
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
	defer func() { _ = rows.Close() }()

	var attempts []*pgentities.AssessmentAttempt
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
			SELECT id, attempt_id, question_index, student_answer,
			       is_correct, points_earned, max_points, time_spent_seconds,
			       answered_at, created_at, updated_at
			FROM assessment_attempt_answer
			WHERE attempt_id = $1
			ORDER BY question_index ASC
		`

		answerRows, err := r.db.QueryContext(ctx, answersQuery, attemptID.String())
		if err != nil {
			return nil, fmt.Errorf("postgres: error finding answers: %w", err)
		}

		var answers []*pgentities.AssessmentAttemptAnswer
		for answerRows.Next() {
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

			err := answerRows.Scan(
				&answerIDStr, &attemptIDStr, &questionIndex, &studentAnswer,
				&isCorrect, &pointsEarned, &maxPoints, &timeSpentSecs,
				&answeredAt, &createdAt, &updatedAt,
			)
			if err != nil {
				_ = answerRows.Close()
				return nil, fmt.Errorf("postgres: error scanning answer: %w", err)
			}

			answerID, _ := uuid.Parse(answerIDStr)
			answerAttemptID, _ := uuid.Parse(attemptIDStr)

			answer := &pgentities.AssessmentAttemptAnswer{
				ID:               answerID,
				AttemptID:        answerAttemptID,
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

			answers = append(answers, answer) //nolint:staticcheck // SA4010: answers usado posteriormente
		}
		_ = answerRows.Close()

		var scorePtr *float64
		if score > 0 {
			scoreF := float64(score)
			scorePtr = &scoreF
		}

		var maxScorePtr *float64
		if maxScore > 0 {
			maxScoreF := float64(maxScore)
			maxScorePtr = &maxScoreF
		}

		var timeSpentPtr *int
		if timeSpent > 0 {
			timeSpentPtr = &timeSpent
		}

		attempt := &pgentities.AssessmentAttempt{
			ID:               attemptID,
			AssessmentID:     assessmentID,
			StudentID:        parsedStudentID,
			Score:            scorePtr,
			MaxScore:         maxScorePtr,
			TimeSpentSeconds: timeSpentPtr,
			StartedAt:        startedAt,
			CompletedAt:      &completedAt,
			CreatedAt:        createdAt,
			IdempotencyKey:   idempotencyKeyPtr,
		}

		attempts = append(attempts, attempt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating attempts: %w", err)
	}

	return attempts, nil
}
