package valueobjects

import "fmt"

// Score representa un puntaje de 0 a 100
type Score struct {
	value int
}

// NewScore crea un Score válido
func NewScore(value int) (Score, error) {
	if value < 0 || value > 100 {
		return Score{}, fmt.Errorf("score must be between 0 and 100, got %d", value)
	}
	return Score{value: value}, nil
}

// Value retorna el valor del score
func (s Score) Value() int {
	return s.value
}

// IsPassing verifica si el score aprueba con un threshold
func (s Score) IsPassing(threshold int) bool {
	return s.value >= threshold
}

// IsFailing verifica si el score reprueba
func (s Score) IsFailing(threshold int) bool {
	return !s.IsPassing(threshold)
}

// String implementa Stringer
func (s Score) String() string {
	return fmt.Sprintf("%d%%", s.value)
}

// Equals compara dos scores
func (s Score) Equals(other Score) bool {
	return s.value == other.value
}
