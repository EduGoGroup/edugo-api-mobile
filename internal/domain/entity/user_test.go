package entity

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconstructUser(t *testing.T) {
	t.Parallel()

	// Arrange
	userID := valueobject.NewUserID()
	email, err := valueobject.NewEmail("test@example.com")
	require.NoError(t, err)

	now := time.Now()

	// Act
	user := ReconstructUser(
		userID,
		email,
		"hashed_password",
		"John",
		"Doe",
		enum.SystemRoleStudent,
		true,
		now,
		now,
	)

	// Assert
	assert.NotNil(t, user, "ReconstructUser should return non-nil user")
	assert.Equal(t, userID, user.ID())
	assert.Equal(t, email, user.Email())
	assert.Equal(t, "hashed_password", user.PasswordHash())
	assert.Equal(t, "John", user.FirstName())
	assert.Equal(t, "Doe", user.LastName())
	assert.Equal(t, enum.SystemRoleStudent, user.Role())
	assert.True(t, user.IsActive())
	assert.Equal(t, now, user.CreatedAt())
	assert.Equal(t, now, user.UpdatedAt())
}

func TestUser_FullName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		firstName string
		lastName  string
		expected  string
	}{
		{
			name:      "nombre completo normal",
			firstName: "John",
			lastName:  "Doe",
			expected:  "John Doe",
		},
		{
			name:      "nombre con un solo carácter",
			firstName: "A",
			lastName:  "B",
			expected:  "A B",
		},
		{
			name:      "nombre con acentos",
			firstName: "José",
			lastName:  "García",
			expected:  "José García",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			userID := valueobject.NewUserID()
			email, _ := valueobject.NewEmail("test@example.com")
			user := ReconstructUser(
				userID, email, "hash", tt.firstName, tt.lastName,
				enum.SystemRoleStudent, true, time.Now(), time.Now(),
			)

			// Act
			fullName := user.FullName()

			// Assert
			assert.Equal(t, tt.expected, fullName)
		})
	}
}

func TestUser_Getters(t *testing.T) {
	t.Parallel()

	// Arrange
	userID := valueobject.NewUserID()
	email, err := valueobject.NewEmail("test@example.com")
	require.NoError(t, err)

	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	user := ReconstructUser(
		userID,
		email,
		"password_hash_123",
		"Jane",
		"Smith",
		enum.SystemRoleTeacher,
		false,
		createdAt,
		updatedAt,
	)

	// Act & Assert
	assert.Equal(t, userID, user.ID())
	assert.Equal(t, email, user.Email())
	assert.Equal(t, "password_hash_123", user.PasswordHash())
	assert.Equal(t, "Jane", user.FirstName())
	assert.Equal(t, "Smith", user.LastName())
	assert.Equal(t, "Jane Smith", user.FullName())
	assert.Equal(t, enum.SystemRoleTeacher, user.Role())
	assert.False(t, user.IsActive())
	assert.Equal(t, createdAt, user.CreatedAt())
	assert.Equal(t, updatedAt, user.UpdatedAt())
}

func TestUser_DifferentRoles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		role enum.SystemRole
	}{
		{
			name: "student role",
			role: enum.SystemRoleStudent,
		},
		{
			name: "teacher role",
			role: enum.SystemRoleTeacher,
		},
		{
			name: "admin role",
			role: enum.SystemRoleAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			userID := valueobject.NewUserID()
			email, _ := valueobject.NewEmail("test@example.com")

			// Act
			user := ReconstructUser(
				userID, email, "hash", "Test", "User",
				tt.role, true, time.Now(), time.Now(),
			)

			// Assert
			assert.Equal(t, tt.role, user.Role())
		})
	}
}
