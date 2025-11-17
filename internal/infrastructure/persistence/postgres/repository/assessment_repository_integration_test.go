//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/postgres/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
)

// AssessmentRepositoryIntegrationSuite tests de integración para AssessmentRepository
type AssessmentRepositoryIntegrationSuite struct {
	suite.IntegrationTestSuite
	repo *repository.PostgresAssessmentRepository
}

// SetupSuite se ejecuta UNA VEZ antes de todos los tests
func (s *AssessmentRepositoryIntegrationSuite) SetupSuite() {
	s.IntegrationTestSuite.SetupSuite()
	// Crear tablas de assessment
	err := createAssessmentTables(s.PostgresDB)
	s.Require().NoError(err, "Tablas de assessment deben crearse correctamente")
}

// SetupTest prepara cada test individual
func (s *AssessmentRepositoryIntegrationSuite) SetupTest() {
	s.IntegrationTestSuite.SetupTest()
	s.repo = repository.NewPostgresAssessmentRepository(s.PostgresDB).(*repository.PostgresAssessmentRepository)
}

// TestAssessmentRepositoryIntegration ejecuta la suite
func TestAssessmentRepositoryIntegration(t *testing.T) {
	testifySuite.Run(t, new(AssessmentRepositoryIntegrationSuite))
}

// TestSave_Insert valida que Save inserta un nuevo assessment
func (s *AssessmentRepositoryIntegrationSuite) TestSave_Insert() {
	ctx := context.Background()

	// Arrange
	maxAttempts := 3
	timeLimit := 60
	assessment := &entities.Assessment{
		ID:               uuid.New(),
		MaterialID:       uuid.New(),
		MongoDocumentID:  "507f1f77bcf86cd799439011",
		Title:            "Test Assessment",
		TotalQuestions:   5,
		PassThreshold:    70,
		MaxAttempts:      &maxAttempts,
		TimeLimitMinutes: &timeLimit,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	// Act
	err := s.repo.Save(ctx, assessment)

	// Assert
	s.NoError(err, "Save debe insertar sin errores")

	// Verificar que se insertó
	found, err := s.repo.FindByID(ctx, assessment.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(assessment.ID, found.ID)
	s.Equal(assessment.Title, found.Title)
	s.Equal(assessment.TotalQuestions, found.TotalQuestions)
	s.Equal(*assessment.MaxAttempts, *found.MaxAttempts)
	s.Equal(*assessment.TimeLimitMinutes, *found.TimeLimitMinutes)
}

// TestSave_Update valida que Save actualiza un assessment existente (UPSERT)
func (s *AssessmentRepositoryIntegrationSuite) TestSave_Update() {
	ctx := context.Background()

	// Arrange - Insertar assessment inicial
	maxAttempts := 3
	assessment := &entities.Assessment{
		ID:               uuid.New(),
		MaterialID:       uuid.New(),
		MongoDocumentID:  "507f1f77bcf86cd799439011",
		Title:            "Original Title",
		TotalQuestions:   5,
		PassThreshold:    70,
		MaxAttempts:      &maxAttempts,
		TimeLimitMinutes: nil,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	err := s.repo.Save(ctx, assessment)
	s.Require().NoError(err)

	// Act - Actualizar el mismo assessment (UPSERT)
	assessment.Title = "Updated Title"
	assessment.TotalQuestions = 10
	newTimeLimit := 90
	assessment.TimeLimitMinutes = &newTimeLimit
	assessment.UpdatedAt = time.Now().UTC()

	err = s.repo.Save(ctx, assessment)

	// Assert
	s.NoError(err, "Save debe actualizar sin errores")

	// Verificar que se actualizó
	found, err := s.repo.FindByID(ctx, assessment.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal("Updated Title", found.Title)
	s.Equal(10, found.TotalQuestions)
	s.Equal(90, *found.TimeLimitMinutes)
}

// TestSave_WithNullValues valida que Save maneja NULL values correctamente
func (s *AssessmentRepositoryIntegrationSuite) TestSave_WithNullValues() {
	ctx := context.Background()

	// Arrange - Assessment sin MaxAttempts ni TimeLimitMinutes
	assessment := &entities.Assessment{
		ID:               uuid.New(),
		MaterialID:       uuid.New(),
		MongoDocumentID:  "507f1f77bcf86cd799439011",
		Title:            "Assessment sin límites",
		TotalQuestions:   5,
		PassThreshold:    70,
		MaxAttempts:      nil,
		TimeLimitMinutes: nil,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	// Act
	err := s.repo.Save(ctx, assessment)

	// Assert
	s.NoError(err)

	// Verificar que se guardó con NULL values
	found, err := s.repo.FindByID(ctx, assessment.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Nil(found.MaxAttempts, "MaxAttempts debe ser nil")
	s.Nil(found.TimeLimitMinutes, "TimeLimitMinutes debe ser nil")
}

// TestFindByID_NotFound valida que FindByID retorna nil cuando no encuentra
func (s *AssessmentRepositoryIntegrationSuite) TestFindByID_NotFound() {
	ctx := context.Background()

	// Act
	found, err := s.repo.FindByID(ctx, uuid.New())

	// Assert
	s.NoError(err)
	s.Nil(found, "Debe retornar nil cuando no encuentra")
}

// TestFindByMaterialID_Success valida que FindByMaterialID encuentra por material
func (s *AssessmentRepositoryIntegrationSuite) TestFindByMaterialID_Success() {
	ctx := context.Background()

	// Arrange - Insertar assessment
	materialID := uuid.New()
	assessment := &entities.Assessment{
		ID:               uuid.New(),
		MaterialID:       materialID,
		MongoDocumentID:  "507f1f77bcf86cd799439011",
		Title:            "Test Assessment",
		TotalQuestions:   5,
		PassThreshold:    70,
		MaxAttempts:      nil,
		TimeLimitMinutes: nil,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	err := s.repo.Save(ctx, assessment)
	s.Require().NoError(err)

	// Act
	found, err := s.repo.FindByMaterialID(ctx, materialID)

	// Assert
	s.NoError(err)
	s.NotNil(found)
	s.Equal(assessment.ID, found.ID)
	s.Equal(materialID, found.MaterialID)
}

// TestFindByMaterialID_NotFound valida que retorna nil cuando no encuentra por material
func (s *AssessmentRepositoryIntegrationSuite) TestFindByMaterialID_NotFound() {
	ctx := context.Background()

	// Act
	found, err := s.repo.FindByMaterialID(ctx, uuid.New())

	// Assert
	s.NoError(err)
	s.Nil(found)
}

// TestDelete_Success valida que Delete elimina correctamente
func (s *AssessmentRepositoryIntegrationSuite) TestDelete_Success() {
	ctx := context.Background()

	// Arrange - Insertar assessment
	assessment := &entities.Assessment{
		ID:               uuid.New(),
		MaterialID:       uuid.New(),
		MongoDocumentID:  "507f1f77bcf86cd799439011",
		Title:            "To Delete",
		TotalQuestions:   5,
		PassThreshold:    70,
		MaxAttempts:      nil,
		TimeLimitMinutes: nil,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	err := s.repo.Save(ctx, assessment)
	s.Require().NoError(err)

	// Act
	err = s.repo.Delete(ctx, assessment.ID)

	// Assert
	s.NoError(err)

	// Verificar que se eliminó
	found, err := s.repo.FindByID(ctx, assessment.ID)
	s.NoError(err)
	s.Nil(found, "Assessment debe estar eliminado")
}

// TestDelete_NotFound valida que Delete falla cuando no encuentra el registro
func (s *AssessmentRepositoryIntegrationSuite) TestDelete_NotFound() {
	ctx := context.Background()

	// Act
	err := s.repo.Delete(ctx, uuid.New())

	// Assert
	s.Error(err, "Delete debe fallar cuando no encuentra el registro")
	assert.Contains(s.T(), err.Error(), "not found")
}
