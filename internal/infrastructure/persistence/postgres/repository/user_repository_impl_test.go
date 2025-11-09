//go:build integration

package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"golang.org/x/crypto/bcrypt"
)

// setupTestDB crea un testcontainer de PostgreSQL para tests de repository
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	ctx := context.Background()

	// Levantar PostgreSQL testcontainer
	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("edugo"),
		postgres.WithUsername("edugo_user"),
		postgres.WithPassword("edugo_pass"),
	)
	require.NoError(t, err, "Failed to start PostgreSQL testcontainer")

	// Obtener connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Conectar a la base de datos
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	// Crear schema de users
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			role VARCHAR(50) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	require.NoError(t, err, "Failed to create users table")

	cleanup := func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return db, cleanup
}

func TestUserRepository_FindByEmail_UserExists(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresUserRepository(db)

	// Crear usuario de prueba
	email := "test@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var userID string
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'John', 'Doe', 'student', true)
		RETURNING id
	`, email, string(hashedPassword)).Scan(&userID)
	require.NoError(t, err)

	// Crear value object de email
	emailVO, errEmail := valueobject.NewEmail(email)
	require.NoError(t, errEmail)

	// Act
	user, err := repo.FindByEmail(context.Background(), emailVO)

	// Assert
	require.NoError(t, err, "FindByEmail should not return error when user exists")
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email().String())
	assert.Equal(t, "John", user.FirstName())
	assert.Equal(t, "Doe", user.LastName())
	assert.Equal(t, enum.SystemRoleStudent, user.Role())
	assert.True(t, user.IsActive())
}

func TestUserRepository_FindByEmail_UserNotFound(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresUserRepository(db)

	// Crear value object de email
	emailVO, errEmail := valueobject.NewEmail("nonexistent@example.com")
	require.NoError(t, errEmail)

	// Act
	user, err := repo.FindByEmail(context.Background(), emailVO)

	// Assert
	require.Error(t, err, "FindByEmail should return error when user not found")
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_FindByID_UserExists(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresUserRepository(db)

	// Crear usuario de prueba
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var userID string
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Jane', 'Smith', 'teacher', true)
		RETURNING id
	`, "jane@example.com", string(hashedPassword)).Scan(&userID)
	require.NoError(t, err)

	userValueObject, err := valueobject.UserIDFromString(userID)
	require.NoError(t, err)

	// Act
	user, err := repo.FindByID(context.Background(), userValueObject)

	// Assert
	require.NoError(t, err, "FindByID should not return error when user exists")
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID().String())
	assert.Equal(t, "jane@example.com", user.Email().String())
	assert.Equal(t, "Jane", user.FirstName())
	assert.Equal(t, "Smith", user.LastName())
	assert.Equal(t, enum.SystemRoleTeacher, user.Role())
}

func TestUserRepository_FindByID_UserNotFound(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresUserRepository(db)

	// Crear un ID que no existe
	nonExistentID := valueobject.NewUserID()

	// Act
	user, err := repo.FindByID(context.Background(), nonExistentID)

	// Assert
	require.Error(t, err, "FindByID should return error when user not found")
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_FindByEmail_MultipleUsersWithSameEmail(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresUserRepository(db)

	email := "unique@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Intentar crear primer usuario
	_, err := db.Exec(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'User', 'One', 'student', true)
	`, email, string(hashedPassword))
	require.NoError(t, err)

	// Intentar crear segundo usuario con mismo email (debe fallar por UNIQUE constraint)
	_, err = db.Exec(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'User', 'Two', 'student', true)
	`, email, string(hashedPassword))
	require.Error(t, err, "Database should enforce UNIQUE constraint on email")

	// Act - Verificar que FindByEmail devuelve el primer usuario
	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)

	user, err := repo.FindByEmail(context.Background(), emailVO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "User", user.FirstName())
	assert.Equal(t, "One", user.LastName())
}

func TestUserRepository_Update(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresUserRepository(db)

	// Crear usuario inicial
	email, _ := valueobject.NewEmail("original@example.com")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var userID string
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Original', 'Name', 'student', true)
		RETURNING id
	`, email.String(), string(hashedPassword)).Scan(&userID)
	require.NoError(t, err)

	userIDValue, err := valueobject.UserIDFromString(userID)
	require.NoError(t, err)

	// Reconstruir user con cambios
	updatedUser := entity.ReconstructUser(
		userIDValue,
		email,
		string(hashedPassword),
		"Updated",
		"Name",
		enum.SystemRoleStudent,
		false, // Cambiar isActive
		time.Now(),
		time.Now(),
	)

	// Act
	err = repo.Update(context.Background(), updatedUser)

	// Assert
	require.NoError(t, err, "Update should not return error")

	// Verificar que se actualizó
	var firstName, lastName string
	var isActive bool
	err = db.QueryRow(`
		SELECT first_name, last_name, is_active
		FROM users
		WHERE id = $1
	`, userID).Scan(&firstName, &lastName, &isActive)
	require.NoError(t, err)

	assert.Equal(t, "Updated", firstName)
	assert.Equal(t, "Name", lastName)
	assert.False(t, isActive)
}
