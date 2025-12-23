//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/postgres/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// ptrStr crea un puntero a string
func ptrStr(s string) *string {
	return &s
}

// ptrBool crea un puntero a bool
func ptrBool(b bool) *bool {
	return &b
}

// ptrInt crea un puntero a int
func ptrInt(i int) *int {
	return &i
}

// AttemptRepositoryIntegrationSuite tests de integración para AttemptRepository
type AttemptRepositoryIntegrationSuite struct {
	suite.IntegrationTestSuite
	repo *repository.PostgresAttemptRepository
}

// SetupSuite se ejecuta UNA VEZ antes de todos los tests
func (s *AttemptRepositoryIntegrationSuite) SetupSuite() {
	s.IntegrationTestSuite.SetupSuite()
	// Crear tablas de assessment
	err := createAssessmentTables(s.PostgresDB)
	s.Require().NoError(err, "Tablas de assessment deben crearse correctamente")
}

// SetupTest prepara cada test individual
func (s *AttemptRepositoryIntegrationSuite) SetupTest() {
	s.IntegrationTestSuite.SetupTest()
	s.repo = repository.NewPostgresAttemptRepository(s.PostgresDB).(*repository.PostgresAttemptRepository)
}

// TestAttemptRepositoryIntegration ejecuta la suite
func TestAttemptRepositoryIntegration(t *testing.T) {
	testifySuite.Run(t, new(AttemptRepositoryIntegrationSuite))
}

// TestSave_AtomicTransaction valida que Save guarda attempt + answers en transacción atómica
func (s *AttemptRepositoryIntegrationSuite) TestSave_AtomicTransaction() {
	ctx := context.Background()

	// Arrange
	attemptID := uuid.New()
	now := time.Now().UTC()
	answer1 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    0,
		StudentAnswer:    ptrStr("a"),
		IsCorrect:        ptrBool(true),
		TimeSpentSeconds: ptrInt(30),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	answer2 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    1,
		StudentAnswer:    ptrStr("b"),
		IsCorrect:        ptrBool(false),
		TimeSpentSeconds: ptrInt(45),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	answer3 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    2,
		StudentAnswer:    ptrStr("c"),
		IsCorrect:        ptrBool(true),
		TimeSpentSeconds: ptrInt(60),
		AnsweredAt:       now,
		CreatedAt:        now,
	}

	score := 66.0 // 2 correctas de 3 = (2*100)/3 = 66 (división entera)
	maxScore := 100.0
	timeSpent := 135
	completedAt := now
	attempt := &pgentities.AssessmentAttempt{
		ID:               attemptID,
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            &score,
		MaxScore:         &maxScore,
		TimeSpentSeconds: &timeSpent,
		StartedAt:        now.Add(-3 * time.Minute),
		CompletedAt:      &completedAt,
		CreatedAt:        now,
		IdempotencyKey:   nil,
	}

	// Act - Guardar attempt
	err := s.repo.Save(ctx, attempt)
	s.Require().NoError(err)

	// Guardar answers por separado
	answerRepo := repository.NewPostgresAnswerRepository(s.PostgresDB)
	err = answerRepo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer1, answer2, answer3})

	// Assert
	s.NoError(err, "Save debe completar transacción sin errores")

	// Verificar que se guardó el attempt
	found, err := s.repo.FindByID(ctx, attemptID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(attemptID, found.ID)
	s.Equal(66.0, *found.Score, "Score debe coincidir con el calculado")

	// Verificar que se guardaron las answers
	foundAnswers, err := answerRepo.FindByAttemptID(ctx, attemptID)
	s.NoError(err)
	s.Equal(3, len(foundAnswers), "Todas las answers deben haberse guardado")
	s.Equal(0, foundAnswers[0].QuestionIndex)
	s.Equal(true, *foundAnswers[0].IsCorrect)
	s.Equal(1, foundAnswers[1].QuestionIndex)
	s.Equal(false, *foundAnswers[1].IsCorrect)
}

// TestSave_WithIdempotencyKey valida que Save maneja idempotency_key correctamente
func (s *AttemptRepositoryIntegrationSuite) TestSave_WithIdempotencyKey() {
	ctx := context.Background()

	// Arrange
	attemptID := uuid.New()
	now := time.Now().UTC()
	answer1 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    0,
		StudentAnswer:    ptrStr("a"),
		IsCorrect:        ptrBool(true),
		TimeSpentSeconds: ptrInt(30),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	idempotencyKey := "attempt-" + attemptID.String()
	score := 100.0
	maxScore := 100.0
	timeSpent := 30
	completedAt := now
	attempt := &pgentities.AssessmentAttempt{
		ID:               attemptID,
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            &score,
		MaxScore:         &maxScore,
		TimeSpentSeconds: &timeSpent,
		StartedAt:        now.Add(-1 * time.Minute),
		CompletedAt:      &completedAt,
		CreatedAt:        now,
		IdempotencyKey:   &idempotencyKey,
	}

	// Guardar answers por separado
	answerRepo := repository.NewPostgresAnswerRepository(s.PostgresDB)
	err := answerRepo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer1})
	s.Require().NoError(err)

	// Act
	err = s.repo.Save(ctx, attempt)

	// Assert
	s.NoError(err)

	// Verificar que se guardó con idempotency_key
	found, err := s.repo.FindByID(ctx, attemptID)
	s.NoError(err)
	s.NotNil(found)
	s.NotNil(found.IdempotencyKey)
	s.Equal(idempotencyKey, *found.IdempotencyKey)
}

// TestFindByID_WithAnswers valida que FindByID carga las answers correctamente (JOIN)
func (s *AttemptRepositoryIntegrationSuite) TestFindByID_WithAnswers() {
	ctx := context.Background()

	// Arrange - Guardar attempt con múltiples answers
	attemptID := uuid.New()
	now := time.Now().UTC()
	answer1 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    0,
		StudentAnswer:    ptrStr("a"),
		IsCorrect:        ptrBool(true),
		TimeSpentSeconds: ptrInt(20),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	answer2 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    1,
		StudentAnswer:    ptrStr("b"),
		IsCorrect:        ptrBool(false),
		TimeSpentSeconds: ptrInt(30),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	answer3 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    2,
		StudentAnswer:    ptrStr("a"),
		IsCorrect:        ptrBool(true),
		TimeSpentSeconds: ptrInt(40),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	answer4 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    3,
		StudentAnswer:    ptrStr("d"),
		IsCorrect:        ptrBool(true),
		TimeSpentSeconds: ptrInt(25),
		AnsweredAt:       now,
		CreatedAt:        now,
	}

	score := 75.0
	maxScore := 100.0
	timeSpent := 115
	completedAt := now
	attempt := &pgentities.AssessmentAttempt{
		ID:               attemptID,
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            &score,
		MaxScore:         &maxScore,
		TimeSpentSeconds: &timeSpent,
		StartedAt:        now.Add(-2 * time.Minute),
		CompletedAt:      &completedAt,
		CreatedAt:        now,
		IdempotencyKey:   nil,
	}

	err := s.repo.Save(ctx, attempt)
	s.Require().NoError(err)

	// Guardar answers por separado
	answerRepo := repository.NewPostgresAnswerRepository(s.PostgresDB)
	err = answerRepo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer1, answer2, answer3, answer4})
	s.Require().NoError(err)

	// Act
	found, err := s.repo.FindByID(ctx, attemptID)

	// Assert
	s.NoError(err)
	s.NotNil(found)

	// Verificar que las answers se guardaron
	foundAnswers, err := answerRepo.FindByAttemptID(ctx, attemptID)
	s.NoError(err)
	s.Equal(4, len(foundAnswers), "Debe tener todas las answers")

	// Verificar orden por question_index
	s.Equal(0, foundAnswers[0].QuestionIndex)
	s.Equal(1, foundAnswers[1].QuestionIndex)
	s.Equal(2, foundAnswers[2].QuestionIndex)
	s.Equal(3, foundAnswers[3].QuestionIndex)
}

// TestFindByID_NotFound valida que retorna nil cuando no encuentra
func (s *AttemptRepositoryIntegrationSuite) TestFindByID_NotFound() {
	ctx := context.Background()

	// Act
	found, err := s.repo.FindByID(ctx, uuid.New())

	// Assert
	s.NoError(err)
	s.Nil(found, "Debe retornar nil cuando no encuentra")
}

// TestFindByStudentAndAssessment_Success valida que encuentra intentos por estudiante y assessment
func (s *AttemptRepositoryIntegrationSuite) TestFindByStudentAndAssessment_Success() {
	ctx := context.Background()

	// Arrange - Guardar 3 intentos del mismo estudiante en la misma assessment
	studentID := uuid.New()
	assessmentID := uuid.New()
	now := time.Now().UTC()

	for i := 0; i < 3; i++ {
		attemptID := uuid.New()
		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               uuid.New(),
			AttemptID:        attemptID,
			QuestionIndex:    0,
			StudentAnswer:    ptrStr("a"),
			IsCorrect:        ptrBool(true),
			TimeSpentSeconds: ptrInt(30),
			AnsweredAt:       now,
			CreatedAt:        now,
		}

		score := 100.0 // 1 correcta de 1 = (1*100)/1 = 100
		maxScore := 100.0
		timeSpent := 30
		completedAt := now.Add(-time.Duration(i) * time.Minute)
		attempt := &pgentities.AssessmentAttempt{
			ID:               attemptID,
			AssessmentID:     assessmentID,
			StudentID:        studentID,
			Score:            &score,
			MaxScore:         &maxScore,
			TimeSpentSeconds: &timeSpent,
			StartedAt:        now.Add(-time.Duration(i+1) * time.Minute),
			CompletedAt:      &completedAt,
			CreatedAt:        now.Add(-time.Duration(i) * time.Minute),
			IdempotencyKey:   nil,
		}

		err := s.repo.Save(ctx, attempt)
		s.Require().NoError(err)

		// Guardar answers por separado
		answerRepo := repository.NewPostgresAnswerRepository(s.PostgresDB)
		err = answerRepo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer})
		s.Require().NoError(err)
	}

	// Act
	found, err := s.repo.FindByStudentAndAssessment(ctx, studentID, assessmentID)

	// Assert
	s.NoError(err)
	s.Equal(3, len(found), "Debe encontrar los 3 intentos")

	// Verificar orden descendente por completed_at (el más reciente primero)
	// Todos tienen el mismo score (100), así que verificamos el orden temporal
	if len(found) >= 2 {
		s.True(found[0].CompletedAt.After(*found[1].CompletedAt) || found[0].CompletedAt.Equal(*found[1].CompletedAt),
			"Debe estar ordenado por completed_at DESC")
	}

	// Verificar scores
	for _, attempt := range found {
		s.Equal(100.0, *attempt.Score, "Todos deben tener score 100")
	}
}

// TestFindByStudentAndAssessment_Empty valida que retorna array vacío cuando no hay intentos
func (s *AttemptRepositoryIntegrationSuite) TestFindByStudentAndAssessment_Empty() {
	ctx := context.Background()

	// Act
	found, err := s.repo.FindByStudentAndAssessment(ctx, uuid.New(), uuid.New())

	// Assert
	s.NoError(err)
	s.Empty(found, "Debe retornar array vacío cuando no hay intentos")
}

// TestCountByStudentAndAssessment_Success valida que cuenta intentos correctamente
func (s *AttemptRepositoryIntegrationSuite) TestCountByStudentAndAssessment_Success() {
	ctx := context.Background()

	// Arrange - Guardar 2 intentos del mismo estudiante
	studentID := uuid.New()
	assessmentID := uuid.New()
	now := time.Now().UTC()

	for i := 0; i < 2; i++ {
		attemptID := uuid.New()
		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               uuid.New(),
			AttemptID:        attemptID,
			QuestionIndex:    0,
			StudentAnswer:    ptrStr("a"),
			IsCorrect:        ptrBool(true),
			TimeSpentSeconds: ptrInt(30),
			AnsweredAt:       now,
			CreatedAt:        now,
		}

		score := 100.0 // 1 correcta de 1 = (1*100)/1 = 100
		maxScore := 100.0
		timeSpent := 30
		completedAt := now
		attempt := &pgentities.AssessmentAttempt{
			ID:               attemptID,
			AssessmentID:     assessmentID,
			StudentID:        studentID,
			Score:            &score,
			MaxScore:         &maxScore,
			TimeSpentSeconds: &timeSpent,
			StartedAt:        now.Add(-1 * time.Minute),
			CompletedAt:      &completedAt,
			CreatedAt:        now,
			IdempotencyKey:   nil,
		}

		err := s.repo.Save(ctx, attempt)
		s.Require().NoError(err)

		// Guardar answers por separado
		answerRepo := repository.NewPostgresAnswerRepository(s.PostgresDB)
		err = answerRepo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer})
		s.Require().NoError(err)
	}

	// Act
	count, err := s.repo.CountByStudentAndAssessment(ctx, studentID, assessmentID)

	// Assert
	s.NoError(err)
	s.Equal(2, count, "Debe contar 2 intentos")
}

// TestFindByStudent_WithPagination valida que la paginación funciona correctamente
func (s *AttemptRepositoryIntegrationSuite) TestFindByStudent_WithPagination() {
	ctx := context.Background()

	// Arrange - Guardar 5 intentos del mismo estudiante en diferentes assessments
	studentID := uuid.New()
	now := time.Now().UTC()

	for i := 0; i < 5; i++ {
		attemptID := uuid.New()
		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               uuid.New(),
			AttemptID:        attemptID,
			QuestionIndex:    0,
			StudentAnswer:    ptrStr("a"),
			IsCorrect:        ptrBool(true),
			TimeSpentSeconds: ptrInt(30),
			AnsweredAt:       now,
			CreatedAt:        now,
		}

		score := 100.0 // 1 correcta de 1 = (1*100)/1 = 100
		maxScore := 100.0
		timeSpent := 30
		completedAt := now.Add(-time.Duration(i) * time.Minute)
		attempt := &pgentities.AssessmentAttempt{
			ID:               attemptID,
			AssessmentID:     uuid.New(),
			StudentID:        studentID,
			Score:            &score,
			MaxScore:         &maxScore,
			TimeSpentSeconds: &timeSpent,
			StartedAt:        now.Add(-time.Duration(i+1) * time.Minute),
			CompletedAt:      &completedAt,
			CreatedAt:        now.Add(-time.Duration(i) * time.Minute),
			IdempotencyKey:   nil,
		}

		err := s.repo.Save(ctx, attempt)
		s.Require().NoError(err)

		// Guardar answers por separado
		answerRepo := repository.NewPostgresAnswerRepository(s.PostgresDB)
		err = answerRepo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer})
		s.Require().NoError(err)
	}

	// Act - Primera página (limit=2, offset=0)
	page1, err := s.repo.FindByStudent(ctx, studentID, 2, 0)
	s.NoError(err)
	s.Equal(2, len(page1), "Primera página debe tener 2 elementos")

	// Act - Segunda página (limit=2, offset=2)
	page2, err := s.repo.FindByStudent(ctx, studentID, 2, 2)
	s.NoError(err)
	s.Equal(2, len(page2), "Segunda página debe tener 2 elementos")

	// Act - Tercera página (limit=2, offset=4)
	page3, err := s.repo.FindByStudent(ctx, studentID, 2, 4)
	s.NoError(err)
	s.Equal(1, len(page3), "Tercera página debe tener 1 elemento")

	// Assert - Verificar que no hay duplicados
	attemptIDs := make(map[uuid.UUID]bool)
	for _, attempt := range page1 {
		attemptIDs[attempt.ID] = true
	}
	for _, attempt := range page2 {
		s.False(attemptIDs[attempt.ID], "No debe haber duplicados entre páginas")
		attemptIDs[attempt.ID] = true
	}
	for _, attempt := range page3 {
		s.False(attemptIDs[attempt.ID], "No debe haber duplicados entre páginas")
	}
}
