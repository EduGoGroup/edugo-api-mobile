//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	testifySuite "github.com/stretchr/testify/suite"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// ptr crea un puntero a un valor
func ptr[T any](v T) *T {
	return &v
}

// MaterialRepositoryIntegrationSuite tests de integración para MaterialRepository
type MaterialRepositoryIntegrationSuite struct {
	suite.IntegrationTestSuite
	repo repository.MaterialRepository
}

// SetupSuite se ejecuta UNA VEZ antes de todos los tests
func (s *MaterialRepositoryIntegrationSuite) SetupSuite() {
	s.IntegrationTestSuite.SetupSuite()
	// La tabla materials ya debe existir por las migraciones de infrastructure
}

// SetupTest prepara cada test individual
func (s *MaterialRepositoryIntegrationSuite) SetupTest() {
	s.IntegrationTestSuite.SetupTest()
	s.repo = NewPostgresMaterialRepository(s.PostgresDB)
}

// TestMaterialRepositoryIntegration ejecuta la suite
func TestMaterialRepositoryIntegration(t *testing.T) {
	testifySuite.Run(t, new(MaterialRepositoryIntegrationSuite))
}

// TestCreate valida que Create inserta un material correctamente
func (s *MaterialRepositoryIntegrationSuite) TestCreate() {
	ctx := context.Background()

	// Arrange
	authorID := valueobject.NewUserID()
	now := time.Now()
	material := &pgentities.Material{
		ID:                  valueobject.NewMaterialID(),
		Title:               "Test Material",
		Description:         ptr("Description"),
		SchoolID:            nil,
		UploadedByTeacherID: authorID,
		FileURL:             "",
		FileType:            "",
		FileSizeBytes:       0,
		Status:              "draft",
		IsPublic:            false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Act
	err := s.repo.Create(ctx, *material)

	// Assert
	s.NoError(err, "Create should not return error")

	// Verificar que se creó en DB
	var count int
	s.PostgresDB.QueryRow("SELECT COUNT(*) FROM materials WHERE id = $1", material.ID.String()).Scan(&count)
	s.Equal(1, count, "Material should be in database")
}

// TestFindByID_MaterialExists valida que FindByID retorna material cuando existe
func (s *MaterialRepositoryIntegrationSuite) TestFindByID_MaterialExists() {
	ctx := context.Background()

	// Arrange - Crear material directamente en DB
	materialID := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()

	_, err := s.PostgresDB.Exec(`
		INSERT INTO materials (id, title, description, author_id, subject_id, status, processing_status)
		VALUES ($1, $2, $3, $4, $5, 'published', 'completed')
	`, materialID.String(), "Test Material", "Test Description", authorID.String(), "subject-1")
	s.Require().NoError(err)

	// Act
	material, err := s.repo.FindByID(ctx, materialID)

	// Assert
	s.NoError(err, "FindByID should not return error when material exists")
	s.NotNil(material)
	s.Equal("Test Material", material.Title())
	s.Equal("Test Description", material.Description())
}

// TestFindByID_MaterialNotFound valida que FindByID retorna nil cuando no existe
func (s *MaterialRepositoryIntegrationSuite) TestFindByID_MaterialNotFound() {
	ctx := context.Background()

	// Arrange
	nonExistentID := valueobject.NewMaterialID()

	// Act
	material, err := s.repo.FindByID(ctx, nonExistentID)

	// Assert
	s.NoError(err, "FindByID should not error on not found")
	s.Nil(material, "Material should be nil when not found")
}

// TestFindByAuthor valida que FindByAuthor retorna materiales del autor
func (s *MaterialRepositoryIntegrationSuite) TestFindByAuthor() {
	ctx := context.Background()

	// Arrange
	authorID := valueobject.NewUserID()

	// Crear 2 materiales del mismo autor
	for i := 1; i <= 2; i++ {
		materialID := valueobject.NewMaterialID()
		_, err := s.PostgresDB.Exec(`
			INSERT INTO materials (id, title, description, author_id, status, processing_status)
			VALUES ($1, $2, $3, $4, 'published', 'completed')
		`, materialID.String(), fmt.Sprintf("Material %d", i), "Description", authorID.String())
		s.Require().NoError(err)
	}

	// Act
	materials, err := s.repo.FindByAuthor(ctx, authorID)

	// Assert
	s.NoError(err)
	s.Len(materials, 2, "Should find 2 materials")
}
