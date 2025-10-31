package valueobject

import (
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/validator"
)

// Email representa un email válido
type Email struct {
	value string
}

// NewEmail crea y valida un nuevo Email
func NewEmail(value string) (Email, error) {
	normalized := validator.Normalize(value)

	if !validator.IsValidEmail(normalized) {
		return Email{}, errors.NewValidationError("invalid email format").
			WithField("email", value)
	}

	return Email{value: normalized}, nil
}

// String retorna la representación en string del email
func (e Email) String() string {
	return e.value
}

// IsZero verifica si es el valor cero
func (e Email) IsZero() bool {
	return e.value == ""
}
