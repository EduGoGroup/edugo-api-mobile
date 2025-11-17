//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/postgres/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
)

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
	answer1, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)
	answer2, _ := entities.NewAnswer(attemptID, "q2", "b", false, 45)
	answer3, _ := entities.NewAnswer(attemptID, "q3", "c", true, 60)

	attempt := &entities.Attempt{
		ID:               attemptID,
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            66, // 2 correctas de 3 = (2*100)/3 = 66 (división entera)
		MaxScore:         100,
		TimeSpentSeconds: 135,
		StartedAt:        time.Now().UTC().Add(-3 * time.Minute),
		CompletedAt:      time.Now().UTC(),
		CreatedAt:        time.Now().UTC(),
		Answers:          []*entities.Answer{answer1, answer2, answer3},
		IdempotencyKey:   nil,
	}

	// Act
	err := s.repo.Save(ctx, attempt)

	// Assert
	s.NoError(err, "Save debe completar transacción sin errores")

	// Verificar que se guardaron attempt Y answers (transacción atómica)
	found, err := s.repo.FindByID(ctx, attemptID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(attemptID, found.ID)
	s.Equal(66, found.Score, "Score debe coincidir con el calculado")
	s.Equal(3, len(found.Answers), "Todas las answers deben haberse guardado")
	s.Equal("q1", found.Answers[0].QuestionID)
	s.Equal(true, found.Answers[0].IsCorrect)
	s.Equal("q2", found.Answers[1].QuestionID)
	s.Equal(false, found.Answers[1].IsCorrect)
}

// TestSave_WithIdempotencyKey valida que Save maneja idempotency_key correctamente
func (s *AttemptRepositoryIntegrationSuite) TestSave_WithIdempotencyKey() {
	ctx := context.Background()

	// Arrange
	attemptID := uuid.New()
	answer1, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)
	idempotencyKey := "attempt-" + attemptID.String()

	attempt := &entities.Attempt{
		ID:               attemptID,
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            100,
		MaxScore:         100,
		TimeSpentSeconds: 30,
		StartedAt:        time.Now().UTC().Add(-1 * time.Minute),
		CompletedAt:      time.Now().UTC(),
		CreatedAt:        time.Now().UTC(),
		Answers:          []*entities.Answer{answer1},
		IdempotencyKey:   &idempotencyKey,
	}

	// Act
	err := s.repo.Save(ctx, attempt)

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
	answer1, _ := entities.NewAnswer(attemptID, "q1", "a", true, 20)
	answer2, _ := entities.NewAnswer(attemptID, "q2", "b", false, 30)
	answer3, _ := entities.NewAnswer(attemptID, "q3", "a", true, 40)
	answer4, _ := entities.NewAnswer(attemptID, "q4", "d", true, 25)

	attempt := &entities.Attempt{
		ID:               attemptID,
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            75,
		MaxScore:         100,
		TimeSpentSeconds: 115,
		StartedAt:        time.Now().UTC().Add(-2 * time.Minute),
		CompletedAt:      time.Now().UTC(),
		CreatedAt:        time.Now().UTC(),
		Answers:          []*entities.Answer{answer1, answer2, answer3, answer4},
		IdempotencyKey:   nil,
	}

	err := s.repo.Save(ctx, attempt)
	s.Require().NoError(err)

	// Act
	found, err := s.repo.FindByID(ctx, attemptID)

	// Assert
	s.NoError(err)
	s.NotNil(found)
	s.Equal(4, len(found.Answers), "Debe cargar todas las answers con JOIN")

	// Verificar orden por created_at
	s.Equal("q1", found.Answers[0].QuestionID)
	s.Equal("q2", found.Answers[1].QuestionID)
	s.Equal("q3", found.Answers[2].QuestionID)
	s.Equal("q4", found.Answers[3].QuestionID)
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

	for i := 0; i < 3; i++ {
		attemptID := uuid.New()
		answer, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)

		attempt := &entities.Attempt{
			ID:               attemptID,
			AssessmentID:     assessmentID,
			StudentID:        studentID,
			Score:            100, // 1 correcta de 1 = (1*100)/1 = 100
			MaxScore:         100,
			TimeSpentSeconds: 30,
			StartedAt:        time.Now().UTC().Add(-time.Duration(i+1) * time.Minute),
			CompletedAt:      time.Now().UTC().Add(-time.Duration(i) * time.Minute),
			CreatedAt:        time.Now().UTC().Add(-time.Duration(i) * time.Minute),
			Answers:          []*entities.Answer{answer},
			IdempotencyKey:   nil,
		}

		err := s.repo.Save(ctx, attempt)
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
		s.True(found[0].CompletedAt.After(found[1].CompletedAt) || found[0].CompletedAt.Equal(found[1].CompletedAt),
			"Debe estar ordenado por completed_at DESC")
	}

	// Verificar que cada attempt tiene sus answers
	for _, attempt := range found {
		s.NotEmpty(attempt.Answers, "Cada attempt debe tener sus answers cargadas")
		s.Equal(100, attempt.Score, "Todos deben tener score 100")
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

	for i := 0; i < 2; i++ {
		attemptID := uuid.New()
		answer, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)

		attempt := &entities.Attempt{
			ID:               attemptID,
			AssessmentID:     assessmentID,
			StudentID:        studentID,
			Score:            100, // 1 correcta de 1 = (1*100)/1 = 100
			MaxScore:         100,
			TimeSpentSeconds: 30,
			StartedAt:        time.Now().UTC().Add(-1 * time.Minute),
			CompletedAt:      time.Now().UTC(),
			CreatedAt:        time.Now().UTC(),
			Answers:          []*entities.Answer{answer},
			IdempotencyKey:   nil,
		}

		err := s.repo.Save(ctx, attempt)
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

	for i := 0; i < 5; i++ {
		attemptID := uuid.New()
		answer, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)

		attempt := &entities.Attempt{
			ID:               attemptID,
			AssessmentID:     uuid.New(),
			StudentID:        studentID,
			Score:            100, // 1 correcta de 1 = (1*100)/1 = 100
			MaxScore:         100,
			TimeSpentSeconds: 30,
			StartedAt:        time.Now().UTC().Add(-time.Duration(i+1) * time.Minute),
			CompletedAt:      time.Now().UTC().Add(-time.Duration(i) * time.Minute),
			CreatedAt:        time.Now().UTC().Add(-time.Duration(i) * time.Minute),
			Answers:          []*entities.Answer{answer},
			IdempotencyKey:   nil,
		}

		err := s.repo.Save(ctx, attempt)
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
