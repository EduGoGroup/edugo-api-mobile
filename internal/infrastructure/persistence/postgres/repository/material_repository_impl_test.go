//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
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

func (s *MaterialRepositoryIntegrationSuite) getSeedSchoolAndAuthor() (uuid.UUID, valueobject.UserID) {
	var schoolIDStr string
	err := s.PostgresDB.QueryRow(`SELECT id::text FROM schools LIMIT 1`).Scan(&schoolIDStr)
	s.Require().NoError(err)

	var authorIDStr string
	err = s.PostgresDB.QueryRow(`SELECT id::text FROM users LIMIT 1`).Scan(&authorIDStr)
	s.Require().NoError(err)

	schoolID, err := uuid.Parse(schoolIDStr)
	s.Require().NoError(err)

	authorID, err := valueobject.UserIDFromString(authorIDStr)
	s.Require().NoError(err)

	return schoolID, authorID
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
	schoolID, authorID := s.getSeedSchoolAndAuthor()
	materialID := valueobject.NewMaterialID()
	now := time.Now()
	material := &pgentities.Material{
		ID:                    materialID.UUID().UUID,
		SchoolID:              schoolID,
		UploadedByTeacherID:   authorID.UUID().UUID,
		AcademicUnitID:        nil,
		Title:                 "Test Material",
		Description:           ptr("Description"),
		Subject:               ptr("Mathematics"),
		Grade:                 ptr("10th"),
		FileURL:               "",
		FileType:              "",
		FileSizeBytes:         0,
		Status:                "uploaded",
		ProcessingStartedAt:   nil,
		ProcessingCompletedAt: nil,
		IsPublic:              false,
		CreatedAt:             now,
		UpdatedAt:             now,
		DeletedAt:             nil,
	}

	// Act
	err := s.repo.Create(ctx, material)

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
	schoolID, authorID := s.getSeedSchoolAndAuthor()

	_, err := s.PostgresDB.Exec(`
		INSERT INTO materials (id, school_id, uploaded_by_teacher_id, title, description, subject, grade, file_url, file_type, file_size_bytes, status, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`, materialID.UUID().UUID, schoolID, authorID.UUID().UUID, "Test Material", "Test Description", "Math", "10th", "https://example.com/file.pdf", "pdf", 1024, "ready", true, time.Now(), time.Now())
	s.Require().NoError(err)

	// Act
	material, err := s.repo.FindByID(ctx, materialID)

	// Assert
	s.NoError(err, "FindByID should not return error when material exists")
	s.NotNil(material)
	s.Equal("Test Material", material.Title)
	s.NotNil(material.Description)
	s.Equal("Test Description", *material.Description)
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
	schoolID, authorID := s.getSeedSchoolAndAuthor()

	// Crear 2 materiales del mismo autor
	for i := 1; i <= 2; i++ {
		materialID := valueobject.NewMaterialID()
		_, err := s.PostgresDB.Exec(`
			INSERT INTO materials (id, school_id, uploaded_by_teacher_id, title, description, subject, grade, file_url, file_type, file_size_bytes, status, is_public, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`, materialID.UUID().UUID, schoolID, authorID.UUID().UUID, fmt.Sprintf("Material %d", i), "Description", "Math", "10th", "https://example.com/file.pdf", "pdf", 1024, "ready", true, time.Now(), time.Now())
		s.Require().NoError(err)
	}

	// Act
	materials, err := s.repo.FindByAuthor(ctx, authorID)

	// Assert
	s.NoError(err)
	s.Len(materials, 2, "Should find 2 materials")
}
