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

// AnswerRepositoryIntegrationSuite tests de integración para AnswerRepository
type AnswerRepositoryIntegrationSuite struct {
	suite.IntegrationTestSuite
	repo *repository.PostgresAnswerRepository
}

// SetupSuite se ejecuta UNA VEZ antes de todos los tests
func (s *AnswerRepositoryIntegrationSuite) SetupSuite() {
	s.IntegrationTestSuite.SetupSuite()
	// Crear tablas de assessment
	err := createAssessmentTables(s.PostgresDB)
	s.Require().NoError(err, "Tablas de assessment deben crearse correctamente")
}

// SetupTest prepara cada test individual
func (s *AnswerRepositoryIntegrationSuite) SetupTest() {
	s.IntegrationTestSuite.SetupTest()
	s.repo = repository.NewPostgresAnswerRepository(s.PostgresDB).(*repository.PostgresAnswerRepository)
}

// TestAnswerRepositoryIntegration ejecuta la suite
func TestAnswerRepositoryIntegration(t *testing.T) {
	testifySuite.Run(t, new(AnswerRepositoryIntegrationSuite))
}

// TestSave_BatchInsert valida que Save guarda múltiples answers en batch
func (s *AnswerRepositoryIntegrationSuite) TestSave_BatchInsert() {
	ctx := context.Background()

	// Arrange
	attemptID := uuid.New()
	now := time.Now()
	studentAnswer1 := "a"
	isCorrect1 := true
	timeSpent1 := 30
	answer1 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    0,
		StudentAnswer:    &studentAnswer1,
		IsCorrect:        &isCorrect1,
		TimeSpentSeconds: &timeSpent1,
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	studentAnswer2 := "b"
	isCorrect2 := false
	timeSpent2 := 45
	answer2 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    1,
		StudentAnswer:    &studentAnswer2,
		IsCorrect:        &isCorrect2,
		TimeSpentSeconds: &timeSpent2,
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	studentAnswer3 := "c"
	isCorrect3 := true
	timeSpent3 := 60
	answer3 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    2,
		StudentAnswer:    &studentAnswer3,
		IsCorrect:        &isCorrect3,
		TimeSpentSeconds: &timeSpent3,
		AnsweredAt:       now,
		CreatedAt:        now,
	}

	answers := []*pgentities.AssessmentAttemptAnswer{answer1, answer2, answer3}

	// Act
	err := s.repo.Save(ctx, answers)

	// Assert
	s.NoError(err, "Save debe completar batch insert sin errores")

	// Verificar que se guardaron todas
	found, err := s.repo.FindByAttemptID(ctx, attemptID)
	s.NoError(err)
	s.Equal(3, len(found), "Deben haberse guardado las 3 answers")
}

// TestSave_EmptyArray valida que Save falla con array vacío
func (s *AnswerRepositoryIntegrationSuite) TestSave_EmptyArray() {
	ctx := context.Background()

	// Act
	err := s.repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{})

	// Assert
	s.Error(err, "Save debe fallar con array vacío")
}

// TestFindByAttemptID_OrderedByCreatedAt valida que retorna answers ordenadas
func (s *AnswerRepositoryIntegrationSuite) TestFindByAttemptID_OrderedByCreatedAt() {
	ctx := context.Background()

	// Arrange - Guardar answers en orden específico
	attemptID := uuid.New()
	now := time.Now()
	studentAnswer1 := "a"
	isCorrect1 := true
	timeSpent1 := 10
	answer1 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    0,
		StudentAnswer:    &studentAnswer1,
		IsCorrect:        &isCorrect1,
		TimeSpentSeconds: &timeSpent1,
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	studentAnswer2 := "b"
	isCorrect2 := false
	timeSpent2 := 20
	answer2 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    1,
		StudentAnswer:    &studentAnswer2,
		IsCorrect:        &isCorrect2,
		TimeSpentSeconds: &timeSpent2,
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	studentAnswer3 := "c"
	isCorrect3 := true
	timeSpent3 := 30
	answer3 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionIndex:    2,
		StudentAnswer:    &studentAnswer3,
		IsCorrect:        &isCorrect3,
		TimeSpentSeconds: &timeSpent3,
		AnsweredAt:       now,
		CreatedAt:        now,
	}

	err := s.repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer1, answer2, answer3})
	s.Require().NoError(err)

	// Act
	found, err := s.repo.FindByAttemptID(ctx, attemptID)

	// Assert
	s.NoError(err)
	s.Equal(3, len(found))
	// Verificar orden por created_at ASC
	s.Equal(0, found[0].QuestionIndex)
	s.Equal(1, found[1].QuestionIndex)
	s.Equal(2, found[2].QuestionIndex)
}

// TestFindByAttemptID_Empty valida que retorna array vacío cuando no hay answers
func (s *AnswerRepositoryIntegrationSuite) TestFindByAttemptID_Empty() {
	ctx := context.Background()

	// Act
	found, err := s.repo.FindByAttemptID(ctx, uuid.New())

	// Assert
	s.NoError(err)
	s.Empty(found, "Debe retornar array vacío cuando no hay answers")
}

// TestFindByQuestionID_WithPagination valida paginación por question_index
func (s *AnswerRepositoryIntegrationSuite) TestFindByQuestionID_WithPagination() {
	ctx := context.Background()

	// Arrange - Guardar 5 respuestas para el mismo índice de pregunta (diferentes intentos)
	questionIndex := 0
	now := time.Now()

	for i := 0; i < 5; i++ {
		attemptID := uuid.New()
		studentAnswer := "a"
		isCorrect := i%2 == 0
		timeSpent := 30
		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               uuid.New(),
			AttemptID:        attemptID,
			QuestionIndex:    questionIndex,
			StudentAnswer:    &studentAnswer,
			IsCorrect:        &isCorrect,
			TimeSpentSeconds: &timeSpent,
			AnsweredAt:       now,
			CreatedAt:        now,
		}

		err := s.repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer})
		s.Require().NoError(err)
	}

	// Act - Primera página (limit=2, offset=0)
	page1, err := s.repo.FindByQuestionID(ctx, "0", 2, 0)
	s.NoError(err)
	s.Equal(2, len(page1), "Primera página debe tener 2 elementos")

	// Act - Segunda página (limit=2, offset=2)
	page2, err := s.repo.FindByQuestionID(ctx, "0", 2, 2)
	s.NoError(err)
	s.Equal(2, len(page2), "Segunda página debe tener 2 elementos")

	// Act - Tercera página (limit=2, offset=4)
	page3, err := s.repo.FindByQuestionID(ctx, "0", 2, 4)
	s.NoError(err)
	s.Equal(1, len(page3), "Tercera página debe tener 1 elemento")

	// Assert - Todas deben ser del mismo índice
	for _, answer := range page1 {
		s.Equal(questionIndex, answer.QuestionIndex)
	}
	for _, answer := range page2 {
		s.Equal(questionIndex, answer.QuestionIndex)
	}
	for _, answer := range page3 {
		s.Equal(questionIndex, answer.QuestionIndex)
	}
}

// TestFindByQuestionID_InvalidQuestionID valida que maneja question_id inválido
func (s *AnswerRepositoryIntegrationSuite) TestFindByQuestionID_InvalidQuestionID() {
	ctx := context.Background()

	// Act
	found, err := s.repo.FindByQuestionID(ctx, "", 10, 0)

	// Assert
	s.Error(err, "Debe fallar con question_id vacío")
	s.Empty(found)
}

// TestGetQuestionDifficultyStats_Success valida cálculo de estadísticas
func (s *AnswerRepositoryIntegrationSuite) TestGetQuestionDifficultyStats_Success() {
	ctx := context.Background()

	// Arrange - Guardar 10 respuestas: 7 correctas, 3 incorrectas
	questionIndex := 1
	now := time.Now()

	for i := 0; i < 10; i++ {
		attemptID := uuid.New()
		isCorrect := i < 7 // Primeras 7 son correctas
		studentAnswer := "a"
		timeSpent := 30
		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               uuid.New(),
			AttemptID:        attemptID,
			QuestionIndex:    questionIndex,
			StudentAnswer:    &studentAnswer,
			IsCorrect:        &isCorrect,
			TimeSpentSeconds: &timeSpent,
			AnsweredAt:       now,
			CreatedAt:        now,
		}

		err := s.repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer})
		s.Require().NoError(err)
	}

	// Act
	total, correct, errorRate, err := s.repo.GetQuestionDifficultyStats(ctx, "1")

	// Assert
	s.NoError(err)
	s.Equal(10, total, "Total debe ser 10")
	s.Equal(7, correct, "Correctas debe ser 7")
	s.InDelta(0.30, errorRate, 0.01, "Error rate debe ser ~30%")
}

// TestGetQuestionDifficultyStats_NoData valida stats cuando no hay datos
func (s *AnswerRepositoryIntegrationSuite) TestGetQuestionDifficultyStats_NoData() {
	ctx := context.Background()

	// Act
	total, correct, errorRate, err := s.repo.GetQuestionDifficultyStats(ctx, "nonexistent")

	// Assert
	s.NoError(err)
	s.Equal(0, total)
	s.Equal(0, correct)
	s.Equal(0.0, errorRate)
}

// TestGetQuestionDifficultyStats_AllCorrect valida stats cuando todas son correctas
func (s *AnswerRepositoryIntegrationSuite) TestGetQuestionDifficultyStats_AllCorrect() {
	ctx := context.Background()

	// Arrange - Todas correctas
	questionIndex := 2
	now := time.Now()

	for i := 0; i < 5; i++ {
		attemptID := uuid.New()
		studentAnswer := "a"
		isCorrect := true
		timeSpent := 30
		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               uuid.New(),
			AttemptID:        attemptID,
			QuestionIndex:    questionIndex,
			StudentAnswer:    &studentAnswer,
			IsCorrect:        &isCorrect,
			TimeSpentSeconds: &timeSpent,
			AnsweredAt:       now,
			CreatedAt:        now,
		}

		err := s.repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer})
		s.Require().NoError(err)
	}

	// Act
	total, correct, errorRate, err := s.repo.GetQuestionDifficultyStats(ctx, "2")

	// Assert
	s.NoError(err)
	s.Equal(5, total)
	s.Equal(5, correct)
	s.Equal(0.0, errorRate, "Error rate debe ser 0% cuando todas son correctas")
}

// TestGetQuestionDifficultyStats_AllIncorrect valida stats cuando todas son incorrectas
func (s *AnswerRepositoryIntegrationSuite) TestGetQuestionDifficultyStats_AllIncorrect() {
	ctx := context.Background()

	// Arrange - Todas incorrectas
	questionIndex := 3
	now := time.Now()

	for i := 0; i < 5; i++ {
		attemptID := uuid.New()
		studentAnswer := "wrong"
		isCorrect := false
		timeSpent := 30
		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               uuid.New(),
			AttemptID:        attemptID,
			QuestionIndex:    questionIndex,
			StudentAnswer:    &studentAnswer,
			IsCorrect:        &isCorrect,
			TimeSpentSeconds: &timeSpent,
			AnsweredAt:       now,
			CreatedAt:        now,
		}

		err := s.repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer})
		s.Require().NoError(err)
	}

	// Act
	total, correct, errorRate, err := s.repo.GetQuestionDifficultyStats(ctx, "3")

	// Assert
	s.NoError(err)
	s.Equal(5, total)
	s.Equal(0, correct)
	s.Equal(1.0, errorRate, "Error rate debe ser 100% cuando todas son incorrectas")
}
