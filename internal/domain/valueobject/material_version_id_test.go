package valueobject

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMaterialVersionID(t *testing.T) {
	t.Parallel()

	// Act
	id := NewMaterialVersionID()

	// Assert
	assert.NotEmpty(t, id.String(), "NewMaterialVersionID should generate a non-empty ID")
	assert.False(t, id.IsZero(), "NewMaterialVersionID should not be zero")

	// Verificar que es un UUID válido
	_, err := uuid.Parse(id.String())
	assert.NoError(t, err, "Generated ID should be a valid UUID")
}

func TestNewMaterialVersionID_Uniqueness(t *testing.T) {
	t.Parallel()

	// Arrange & Act
	id1 := NewMaterialVersionID()
	id2 := NewMaterialVersionID()

	// Assert
	assert.NotEqual(t, id1.String(), id2.String(), "Each NewMaterialVersionID should generate unique ID")
	assert.False(t, id1.Equals(id2), "Different IDs should not be equal")
}

func TestMaterialVersionIDFromString_ValidUUID(t *testing.T) {
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
			id, err := MaterialVersionIDFromString(tt.input)

			// Assert
			require.NoError(t, err, "MaterialVersionIDFromString should not return error for valid UUID")
			assert.NotEmpty(t, id.String(), "ID should have a value")
			assert.False(t, id.IsZero(), "Valid ID should not be zero")
		})
	}
}

func TestMaterialVersionIDFromString_InvalidUUID(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			id, err := MaterialVersionIDFromString(tt.input)

			// Assert
			require.Error(t, err, "MaterialVersionIDFromString should return error for invalid UUID")
			assert.True(t, id.IsZero(), "Invalid ID should be zero value")
		})
	}
}

func TestMaterialVersionID_String(t *testing.T) {
	t.Parallel()

	// Arrange
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	id, err := MaterialVersionIDFromString(validUUID)
	require.NoError(t, err)

	// Act
	result := id.String()

	// Assert
	assert.Equal(t, validUUID, result, "String() should return the UUID string")
}

func TestMaterialVersionID_UUID(t *testing.T) {
	t.Parallel()

	// Arrange
	id := NewMaterialVersionID()

	// Act
	uuid := id.UUID()

	// Assert
	assert.NotNil(t, uuid, "UUID() should return non-nil UUID")
	assert.Equal(t, id.String(), uuid.String(), "UUID() should match String()")
}

func TestMaterialVersionID_IsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       MaterialVersionID
		expected bool
	}{
		{
			name:     "ID vacío es zero",
			id:       MaterialVersionID{},
			expected: true,
		},
		{
			name: "ID válido no es zero",
			id: func() MaterialVersionID {
				id, _ := MaterialVersionIDFromString("123e4567-e89b-12d3-a456-426614174000")
				return id
			}(),
			expected: false,
		},
		{
			name:     "ID generado no es zero",
			id:       NewMaterialVersionID(),
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

func TestMaterialVersionID_Equals(t *testing.T) {
	t.Parallel()

	// Arrange
	uuid1 := "123e4567-e89b-12d3-a456-426614174000"
	uuid2 := "987e6543-e21b-34d5-a678-123456789abc"

	id1a, err := MaterialVersionIDFromString(uuid1)
	require.NoError(t, err)

	id1b, err := MaterialVersionIDFromString(uuid1)
	require.NoError(t, err)

	id2, err := MaterialVersionIDFromString(uuid2)
	require.NoError(t, err)

	tests := []struct {
		name     string
		id1      MaterialVersionID
		id2      MaterialVersionID
		expected bool
	}{
		{
			name:     "mismo UUID son iguales",
			id1:      id1a,
			id2:      id1b,
			expected: true,
		},
		{
			name:     "diferentes UUID no son iguales",
			id1:      id1a,
			id2:      id2,
			expected: false,
		},
		{
			name:     "ID con sí mismo es igual",
			id1:      id1a,
			id2:      id1a,
			expected: true,
		},
		{
			name:     "zero con zero son iguales",
			id1:      MaterialVersionID{},
			id2:      MaterialVersionID{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			result := tt.id1.Equals(tt.id2)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaterialVersionID_RoundTrip(t *testing.T) {
	t.Parallel()

	// Arrange
	originalUUID := "123e4567-e89b-12d3-a456-426614174000"

	// Act
	id, err := MaterialVersionIDFromString(originalUUID)
	require.NoError(t, err)

	backToString := id.String()

	// Assert
	assert.Equal(t, originalUUID, backToString, "Round-trip conversion should preserve UUID")
}
