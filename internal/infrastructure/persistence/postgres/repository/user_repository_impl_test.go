//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"
	"time"

	testifySuite "github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// UserRepositoryIntegrationSuite tests de integración para UserRepository
type UserRepositoryIntegrationSuite struct {
	suite.IntegrationTestSuite
	repo repository.UserRepository
}

// SetupSuite se ejecuta UNA VEZ antes de todos los tests
func (s *UserRepositoryIntegrationSuite) SetupSuite() {
	s.IntegrationTestSuite.SetupSuite()
	// La tabla users ya debe existir por las migraciones de infrastructure
}

// SetupTest prepara cada test individual
func (s *UserRepositoryIntegrationSuite) SetupTest() {
	s.IntegrationTestSuite.SetupTest()
	s.repo = NewPostgresUserRepository(s.PostgresDB)
}

// TestUserRepositoryIntegration ejecuta la suite
func TestUserRepositoryIntegration(t *testing.T) {
	testifySuite.Run(t, new(UserRepositoryIntegrationSuite))
}

// TestFindByEmail_UserExists valida que FindByEmail retorna usuario cuando existe
func (s *UserRepositoryIntegrationSuite) TestFindByEmail_UserExists() {
	ctx := context.Background()

	// Arrange
	email := "test@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var userID string
	err := s.PostgresDB.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'John', 'Doe', 'student', true)
		RETURNING id
	`, email, string(hashedPassword)).Scan(&userID)
	s.Require().NoError(err)

	// Crear value object de email
	emailVO, errEmail := valueobject.NewEmail(email)
	s.Require().NoError(errEmail)

	// Act
	user, err := s.repo.FindByEmail(ctx, emailVO)

	// Assert
	s.NoError(err, "FindByEmail should not return error when user exists")
	s.NotNil(user)
	s.Equal(email, user.Email().String())
	s.Equal("John", user.FirstName())
	s.Equal("Doe", user.LastName())
	s.Equal(enum.SystemRoleStudent, user.Role())
	s.True(user.IsActive())
}

// TestFindByEmail_UserNotFound valida que FindByEmail retorna error cuando no existe
func (s *UserRepositoryIntegrationSuite) TestFindByEmail_UserNotFound() {
	ctx := context.Background()

	// Arrange
	emailVO, errEmail := valueobject.NewEmail("nonexistent@example.com")
	s.Require().NoError(errEmail)

	// Act
	user, err := s.repo.FindByEmail(ctx, emailVO)

	// Assert
	s.Error(err, "FindByEmail should return error when user not found")
	s.Nil(user)
	s.Contains(err.Error(), "user not found")
}

// TestFindByID_UserExists valida que FindByID retorna usuario cuando existe
func (s *UserRepositoryIntegrationSuite) TestFindByID_UserExists() {
	ctx := context.Background()

	// Arrange
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var userID string
	err := s.PostgresDB.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Jane', 'Smith', 'teacher', true)
		RETURNING id
	`, "jane@example.com", string(hashedPassword)).Scan(&userID)
	s.Require().NoError(err)

	userValueObject, err := valueobject.UserIDFromString(userID)
	s.Require().NoError(err)

	// Act
	user, err := s.repo.FindByID(ctx, userValueObject)

	// Assert
	s.NoError(err, "FindByID should not return error when user exists")
	s.NotNil(user)
	s.Equal(userID, user.ID().String())
	s.Equal("jane@example.com", user.Email().String())
	s.Equal("Jane", user.FirstName())
	s.Equal("Smith", user.LastName())
	s.Equal(enum.SystemRoleTeacher, user.Role())
}

// TestFindByID_UserNotFound valida que FindByID retorna error cuando no existe
func (s *UserRepositoryIntegrationSuite) TestFindByID_UserNotFound() {
	ctx := context.Background()

	// Arrange
	nonExistentID := valueobject.NewUserID()

	// Act
	user, err := s.repo.FindByID(ctx, nonExistentID)

	// Assert
	s.Error(err, "FindByID should return error when user not found")
	s.Nil(user)
	s.Contains(err.Error(), "user not found")
}

// TestFindByEmail_MultipleUsersWithSameEmail valida que DB enforce UNIQUE constraint
func (s *UserRepositoryIntegrationSuite) TestFindByEmail_MultipleUsersWithSameEmail() {
	ctx := context.Background()

	// Arrange
	email := "unique@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Intentar crear primer usuario
	_, err := s.PostgresDB.Exec(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'User', 'One', 'student', true)
	`, email, string(hashedPassword))
	s.Require().NoError(err)

	// Intentar crear segundo usuario con mismo email (debe fallar por UNIQUE constraint)
	_, err = s.PostgresDB.Exec(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'User', 'Two', 'student', true)
	`, email, string(hashedPassword))
	s.Error(err, "Database should enforce UNIQUE constraint on email")

	// Act - Verificar que FindByEmail devuelve el primer usuario
	emailVO, err := valueobject.NewEmail(email)
	s.Require().NoError(err)

	user, err := s.repo.FindByEmail(ctx, emailVO)

	// Assert
	s.NoError(err)
	s.NotNil(user)
	s.Equal("User", user.FirstName())
	s.Equal("One", user.LastName())
}

// TestUpdate valida que Update actualiza correctamente un usuario
func (s *UserRepositoryIntegrationSuite) TestUpdate() {
	ctx := context.Background()

	// Arrange - Crear usuario inicial
	email, _ := valueobject.NewEmail("original@example.com")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var userID string
	err := s.PostgresDB.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Original', 'Name', 'student', true)
		RETURNING id
	`, email.String(), string(hashedPassword)).Scan(&userID)
	s.Require().NoError(err)

	userIDValue, err := valueobject.UserIDFromString(userID)
	s.Require().NoError(err)

	// Reconstruir user con cambios
	updatedUser := pgentities.User{
		ID:           userIDValue.UUID().UUID,
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    "Updated",
		LastName:     "Name",
		Role:         enum.SystemRoleStudent,
		IsActive:     false, // Cambiar isActive
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Act
	err = s.repo.Update(ctx, updatedUser)

	// Assert
	s.NoError(err, "Update should not return error")

	// Verificar que se actualizó
	var firstName, lastName string
	var isActive bool
	err = s.PostgresDB.QueryRow(`
		SELECT first_name, last_name, is_active
		FROM users
		WHERE id = $1
	`, userID).Scan(&firstName, &lastName, &isActive)
	s.Require().NoError(err)

	s.Equal("Updated", firstName)
	s.Equal("Name", lastName)
	s.False(isActive)
}
