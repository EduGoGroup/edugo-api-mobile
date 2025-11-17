//go:build integration
// +build integration

package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	testifySuite "github.com/stretchr/testify/suite"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// ProgressRepositoryIntegrationSuite tests de integración para ProgressRepository
type ProgressRepositoryIntegrationSuite struct {
	suite.IntegrationTestSuite
	repo repository.ProgressRepository
}

// SetupSuite se ejecuta UNA VEZ antes de todos los tests
func (s *ProgressRepositoryIntegrationSuite) SetupSuite() {
	s.IntegrationTestSuite.SetupSuite()
	// Las tablas users, materials, material_progress ya deben existir por las migraciones
}

// SetupTest prepara cada test individual
func (s *ProgressRepositoryIntegrationSuite) SetupTest() {
	s.IntegrationTestSuite.SetupTest()
	s.repo = NewPostgresProgressRepository(s.PostgresDB)
}

// TestProgressRepositoryIntegration ejecuta la suite
func TestProgressRepositoryIntegration(t *testing.T) {
	testifySuite.Run(t, new(ProgressRepositoryIntegrationSuite))
}

// seedUserAndMaterial crea un usuario y material de prueba usando los datos ya seeded
func (s *ProgressRepositoryIntegrationSuite) seedUserAndMaterial() (valueobject.UserID, valueobject.MaterialID) {
	ctx := context.Background()

	// Usar el usuario y material ya creados por la suite (via infrastructure seeds)
	// O crear nuevos para mayor aislamiento
	var userIDStr string
	err := s.PostgresDB.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Test', 'User', 'student', true)
		RETURNING id
	`, "test@example.com", "hashedpassword").Scan(&userIDStr)
	s.Require().NoError(err, "Failed to create test user")

	userID, err := valueobject.UserIDFromString(userIDStr)
	s.Require().NoError(err, "Failed to parse user ID")

	var materialIDStr string
	err = s.PostgresDB.QueryRowContext(ctx, `
		INSERT INTO materials (title, description, author_id, status, processing_status)
		VALUES ($1, $2, $3, 'published', 'completed')
		RETURNING id
	`, "Test Material", "Test Description", userIDStr).Scan(&materialIDStr)
	s.Require().NoError(err, "Failed to create test material")

	materialID, err := valueobject.MaterialIDFromString(materialIDStr)
	s.Require().NoError(err, "Failed to parse material ID")

	return userID, materialID
}

// TestUpsert_CreateNewProgress valida que Upsert crea nuevo progreso
func (s *ProgressRepositoryIntegrationSuite) TestUpsert_CreateNewProgress() {
	ctx := context.Background()

	// Arrange
	userID, materialID := s.seedUserAndMaterial()

	progress := entity.NewProgress(materialID, userID)
	err := progress.UpdateProgress(25, 5)
	s.Require().NoError(err)

	// Act
	result, err := s.repo.Upsert(ctx, progress)

	// Assert
	s.NoError(err, "Upsert should not return error when creating new progress")
	s.NotNil(result)
	s.Equal(materialID.String(), result.MaterialID().String())
	s.Equal(userID.String(), result.UserID().String())
	s.Equal(25, result.Percentage())
	s.Equal(5, result.LastPage())
	s.Equal(enum.ProgressStatusInProgress, result.Status())

	// Verificar que se creó en DB
	var count int
	err = s.PostgresDB.QueryRow(`
		SELECT COUNT(*) FROM material_progress
		WHERE material_id = $1 AND user_id = $2
	`, materialID.String(), userID.String()).Scan(&count)
	s.Require().NoError(err)
	s.Equal(1, count, "Progress should be in database")
}

// TestUpsert_UpdateExistingProgress valida que Upsert actualiza progreso existente
func (s *ProgressRepositoryIntegrationSuite) TestUpsert_UpdateExistingProgress() {
	ctx := context.Background()

	// Arrange
	userID, materialID := s.seedUserAndMaterial()

	// Crear progreso inicial
	initialProgress := entity.NewProgress(materialID, userID)
	err := initialProgress.UpdateProgress(25, 5)
	s.Require().NoError(err)

	_, err = s.repo.Upsert(ctx, initialProgress)
	s.Require().NoError(err)

	// Esperar un momento para asegurar que updated_at sea diferente
	time.Sleep(10 * time.Millisecond)

	// Actualizar progreso
	updatedProgress := entity.ReconstructProgress(
		materialID,
		userID,
		50,
		10,
		enum.ProgressStatusInProgress,
		time.Now(),
		initialProgress.CreatedAt(),
		time.Now(),
	)

	// Act
	result, err := s.repo.Upsert(ctx, updatedProgress)

	// Assert
	s.NoError(err, "Upsert should not return error when updating existing progress")
	s.NotNil(result)
	s.Equal(50, result.Percentage())
	s.Equal(10, result.LastPage())
	s.Equal(enum.ProgressStatusInProgress, result.Status())

	// Verificar que solo hay un registro en DB
	var count int
	err = s.PostgresDB.QueryRow(`
		SELECT COUNT(*) FROM material_progress
		WHERE material_id = $1 AND user_id = $2
	`, materialID.String(), userID.String()).Scan(&count)
	s.Require().NoError(err)
	s.Equal(1, count, "Should still have only one progress record")

	// Verificar que el porcentaje se actualizó
	var percentage int
	err = s.PostgresDB.QueryRow(`
		SELECT percentage FROM material_progress
		WHERE material_id = $1 AND user_id = $2
	`, materialID.String(), userID.String()).Scan(&percentage)
	s.Require().NoError(err)
	s.Equal(50, percentage)
}

// TestUpsert_CompleteProgress valida que Upsert maneja progreso completado
func (s *ProgressRepositoryIntegrationSuite) TestUpsert_CompleteProgress() {
	ctx := context.Background()

	// Arrange
	userID, materialID := s.seedUserAndMaterial()

	progress := entity.NewProgress(materialID, userID)
	err := progress.UpdateProgress(100, 20)
	s.Require().NoError(err)

	// Act
	result, err := s.repo.Upsert(ctx, progress)

	// Assert
	s.NoError(err, "Upsert should not return error when completing progress")
	s.NotNil(result)
	s.Equal(100, result.Percentage())
	s.Equal(enum.ProgressStatusCompleted, result.Status())

	// Verificar que completed_at se estableció
	var completedAt sql.NullTime
	err = s.PostgresDB.QueryRow(`
		SELECT completed_at FROM material_progress
		WHERE material_id = $1 AND user_id = $2
	`, materialID.String(), userID.String()).Scan(&completedAt)
	s.Require().NoError(err)
	s.True(completedAt.Valid, "completed_at should be set when percentage is 100")
}

// TestFindByMaterialAndUser_ProgressExists valida que FindByMaterialAndUser retorna progreso
func (s *ProgressRepositoryIntegrationSuite) TestFindByMaterialAndUser_ProgressExists() {
	ctx := context.Background()

	// Arrange
	userID, materialID := s.seedUserAndMaterial()

	now := time.Now()
	_, err := s.PostgresDB.Exec(`
		INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, materialID.String(), userID.String(), 75, 15, "in_progress", now, now, now)
	s.Require().NoError(err)

	// Act
	progress, err := s.repo.FindByMaterialAndUser(ctx, materialID, userID)

	// Assert
	s.NoError(err, "FindByMaterialAndUser should not return error when progress exists")
	s.NotNil(progress)
	s.Equal(materialID.String(), progress.MaterialID().String())
	s.Equal(userID.String(), progress.UserID().String())
	s.Equal(75, progress.Percentage())
	s.Equal(15, progress.LastPage())
	s.Equal(enum.ProgressStatusInProgress, progress.Status())
}

// TestFindByMaterialAndUser_ProgressNotFound valida que FindByMaterialAndUser retorna nil cuando no existe
func (s *ProgressRepositoryIntegrationSuite) TestFindByMaterialAndUser_ProgressNotFound() {
	ctx := context.Background()

	// Arrange
	userID, materialID := s.seedUserAndMaterial()
	// No crear ningún progreso

	// Act
	progress, err := s.repo.FindByMaterialAndUser(ctx, materialID, userID)

	// Assert
	s.NoError(err, "FindByMaterialAndUser should not return error when progress not found")
	s.Nil(progress, "Progress should be nil when not found")
}

// TestFindByMaterialAndUser_DifferentUser valida que progreso es específico por usuario
func (s *ProgressRepositoryIntegrationSuite) TestFindByMaterialAndUser_DifferentUser() {
	ctx := context.Background()

	// Arrange
	userID1, materialID := s.seedUserAndMaterial()

	// Crear segundo usuario
	var userID2Str string
	err := s.PostgresDB.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Test2', 'User2', 'student', true)
		RETURNING id
	`, "test2@example.com", "hashedpassword").Scan(&userID2Str)
	s.Require().NoError(err)

	userID2, err := valueobject.UserIDFromString(userID2Str)
	s.Require().NoError(err)

	// Crear progreso para userID1
	now := time.Now()
	_, err = s.PostgresDB.Exec(`
		INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, materialID.String(), userID1.String(), 50, 10, "in_progress", now, now, now)
	s.Require().NoError(err)

	// Act - Buscar progreso para userID2 (no existe)
	progress, err := s.repo.FindByMaterialAndUser(ctx, materialID, userID2)

	// Assert
	s.NoError(err)
	s.Nil(progress, "Should not find progress for different user")
}
