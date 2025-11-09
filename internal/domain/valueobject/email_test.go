package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmail_ValidEmails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "email simple válido",
			input:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "email con subdomain",
			input:    "user@mail.example.com",
			expected: "user@mail.example.com",
		},
		{
			name:     "email con números",
			input:    "user123@example.com",
			expected: "user123@example.com",
		},
		{
			name:     "email con guiones",
			input:    "test-user@example.com",
			expected: "test-user@example.com",
		},
		{
			name:     "email con puntos",
			input:    "first.last@example.com",
			expected: "first.last@example.com",
		},
		{
			name:     "email con espacios (normalizado)",
			input:    "  test@example.com  ",
			expected: "test@example.com",
		},
		{
			name:     "email uppercase (normalizado a lowercase)",
			input:    "TEST@EXAMPLE.COM",
			expected: "test@example.com",
		},
		{
			name:     "email mixto (normalizado)",
			input:    "  TeSt@ExAmPle.CoM  ",
			expected: "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			email, err := NewEmail(tt.input)

			// Assert
			require.NoError(t, err, "NewEmail should not return error for valid email")
			assert.Equal(t, tt.expected, email.String(), "Email should be normalized correctly")
			assert.False(t, email.IsZero(), "Valid email should not be zero")
		})
	}
}

func TestNewEmail_InvalidEmails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "email vacío",
			input: "",
		},
		{
			name:  "email solo espacios",
			input: "   ",
		},
		{
			name:  "sin @",
			input: "testexample.com",
		},
		{
			name:  "sin dominio",
			input: "test@",
		},
		{
			name:  "sin usuario",
			input: "@example.com",
		},
		{
			name:  "sin TLD",
			input: "test@example",
		},
		{
			name:  "múltiples @",
			input: "test@@example.com",
		},
		{
			name:  "caracteres inválidos con espacios",
			input: "test user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			email, err := NewEmail(tt.input)

			// Assert
			require.Error(t, err, "NewEmail should return error for invalid email")
			assert.Contains(t, err.Error(), "invalid email format", "Error should mention invalid email format")
			assert.True(t, email.IsZero(), "Invalid email should be zero value")
		})
	}
}

func TestEmail_String(t *testing.T) {
	t.Parallel()

	// Arrange
	email, err := NewEmail("test@example.com")
	require.NoError(t, err)

	// Act
	result := email.String()

	// Assert
	assert.Equal(t, "test@example.com", result, "String() should return email value")
}

func TestEmail_IsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		email    Email
		expected bool
	}{
		{
			name:     "email vacío es zero",
			email:    Email{value: ""},
			expected: true,
		},
		{
			name:     "email válido no es zero",
			email:    Email{value: "test@example.com"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			result := tt.email.IsZero()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmail_Normalization(t *testing.T) {
	t.Parallel()

	// Arrange & Act
	email1, err1 := NewEmail("Test@Example.Com")
	email2, err2 := NewEmail("test@example.com")
	email3, err3 := NewEmail("  TEST@EXAMPLE.COM  ")

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)

	assert.Equal(t, email1.String(), email2.String(), "Emails should normalize to same value")
	assert.Equal(t, email2.String(), email3.String(), "Emails should normalize to same value")
	assert.Equal(t, "test@example.com", email1.String(), "All should normalize to lowercase trimmed")
}
