package valueobject

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserID(t *testing.T) {
	t.Parallel()

	// Act
	id := NewUserID()

	// Assert
	assert.NotEmpty(t, id.String(), "NewUserID should generate a non-empty ID")
	assert.False(t, id.IsZero(), "NewUserID should not be zero")

	// Verificar que es un UUID válido
	_, err := uuid.Parse(id.String())
	assert.NoError(t, err, "Generated ID should be a valid UUID")
}

func TestNewUserID_Uniqueness(t *testing.T) {
	t.Parallel()

	// Arrange & Act
	id1 := NewUserID()
	id2 := NewUserID()

	// Assert
	assert.NotEqual(t, id1.String(), id2.String(), "Each NewUserID should generate unique ID")
}

func TestUserIDFromString_ValidUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "UUID v4 válido",
			input: "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:  "UUID generado",
			input: uuid.New().String(),
		},
		{
			name:  "UUID lowercase",
			input: "a1b2c3d4-e5f6-4789-a012-3456789abcde",
		},
		{
			name:  "UUID uppercase",
			input: "A1B2C3D4-E5F6-4789-A012-3456789ABCDE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			id, err := UserIDFromString(tt.input)

			// Assert
			require.NoError(t, err, "UserIDFromString should not return error for valid UUID")
			assert.NotEmpty(t, id.String(), "ID should have a value")
			assert.False(t, id.IsZero(), "Valid ID should not be zero")
		})
	}
}

func TestUserIDFromString_InvalidUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "string vacío",
			input: "",
		},
		{
			name:  "no es UUID",
			input: "not-a-uuid",
		},
		{
			name:  "UUID incompleto",
			input: "123e4567-e89b-12d3",
		},
		{
			name:  "UUID con caracteres inválidos",
			input: "123e4567-e89b-12d3-a456-42661417400g",
		},
		{
			name:  "texto random",
			input: "this-is-not-a-valid-uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			id, err := UserIDFromString(tt.input)

			// Assert
			require.Error(t, err, "UserIDFromString should return error for invalid UUID")
			assert.True(t, id.IsZero(), "Invalid ID should be zero value")
		})
	}
}

func TestUserID_String(t *testing.T) {
	t.Parallel()

	// Arrange
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	id, err := UserIDFromString(validUUID)
	require.NoError(t, err)

	// Act
	result := id.String()

	// Assert
	assert.Equal(t, validUUID, result, "String() should return the UUID string")
}

func TestUserID_UUID(t *testing.T) {
	t.Parallel()

	// Arrange
	id := NewUserID()

	// Act
	uuid := id.UUID()

	// Assert
	assert.NotNil(t, uuid, "UUID() should return non-nil UUID")
	assert.Equal(t, id.String(), uuid.String(), "UUID() should match String()")
}

func TestUserID_IsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       UserID
		expected bool
	}{
		{
			name:     "ID vacío es zero",
			id:       UserID{},
			expected: true,
		},
		{
			name: "ID válido no es zero",
			id: func() UserID {
				id, _ := UserIDFromString("123e4567-e89b-12d3-a456-426614174000")
				return id
			}(),
			expected: false,
		},
		{
			name:     "ID generado no es zero",
			id:       NewUserID(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			result := tt.id.IsZero()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserID_RoundTrip(t *testing.T) {
	t.Parallel()

	// Arrange
	originalUUID := "123e4567-e89b-12d3-a456-426614174000"

	// Act
	id, err := UserIDFromString(originalUUID)
	require.NoError(t, err)

	backToString := id.String()

	// Assert
	assert.Equal(t, originalUUID, backToString, "Round-trip conversion should preserve UUID")
}
