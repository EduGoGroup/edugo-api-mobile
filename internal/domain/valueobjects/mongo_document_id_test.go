package valueobjects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobjects"
)

func TestNewMongoDocumentID_Success(t *testing.T) {
	validIDs := []string{
		"507f1f77bcf86cd799439011",
		"123456789012345678901234",
		"ABCDEF1234567890abcdef12",
		"000000000000000000000000",
		"ffffffffffffffffffffffff",
		"FFFFFFFFFFFFFFFFFFFFFFFF",
	}

	for _, id := range validIDs {
		t.Run(id, func(t *testing.T) {
			mongoID, err := valueobjects.NewMongoDocumentID(id)
			require.NoError(t, err)
			assert.Equal(t, id, mongoID.Value())
			assert.Equal(t, id, mongoID.String())
		})
	}
}

func TestNewMongoDocumentID_InvalidLength(t *testing.T) {
	invalidIDs := []string{
		"",                              // empty
		"123",                           // too short
		"12345678901234567890123",       // 23 characters
		"1234567890123456789012345",     // 25 characters
		"507f1f77bcf86cd799439011EXTRA", // too long
	}

	for _, id := range invalidIDs {
		t.Run(id, func(t *testing.T) {
			_, err := valueobjects.NewMongoDocumentID(id)
			assert.ErrorIs(t, err, valueobjects.ErrInvalidMongoDocumentID)
		})
	}
}

func TestNewMongoDocumentID_InvalidCharacters(t *testing.T) {
	invalidIDs := []string{
		"507f1f77bcf86cd799439G11", // G is not hex
		"507f1f77bcf86cd79943901!", // special character
		"507f1f77bcf86cd79943901 ", // space
		"507f1f77bcf86cd79943901-", // dash
		"507f1f77bcf86cd79943901_", // underscore
	}

	for _, id := range invalidIDs {
		t.Run(id, func(t *testing.T) {
			_, err := valueobjects.NewMongoDocumentID(id)
			assert.ErrorIs(t, err, valueobjects.ErrInvalidMongoDocumentID)
		})
	}
}

func TestMongoDocumentID_Equals(t *testing.T) {
	id1, _ := valueobjects.NewMongoDocumentID("507f1f77bcf86cd799439011")
	id2, _ := valueobjects.NewMongoDocumentID("507f1f77bcf86cd799439011")
	id3, _ := valueobjects.NewMongoDocumentID("123456789012345678901234")

	assert.True(t, id1.Equals(id2))
	assert.False(t, id1.Equals(id3))
}

func TestMongoDocumentID_CaseInsensitive(t *testing.T) {
	// Ambas deberían ser válidas
	id1, err1 := valueobjects.NewMongoDocumentID("abcdef1234567890abcdef12")
	id2, err2 := valueobjects.NewMongoDocumentID("ABCDEF1234567890ABCDEF12")

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Pero no son iguales porque el value object preserva el caso original
	assert.False(t, id1.Equals(id2))
}
