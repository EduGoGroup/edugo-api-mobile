package valueobject

import "fmt"

// TimeSpent representa el tiempo gastado en segundos (0-7200)
type TimeSpent struct {
	seconds int
}

// NewTimeSpent crea un TimeSpent válido
func NewTimeSpent(seconds int) (TimeSpent, error) {
	if seconds < 0 {
		return TimeSpent{}, fmt.Errorf("time spent cannot be negative, got %d", seconds)
	}

	if seconds > 7200 { // Máximo 2 horas
		return TimeSpent{}, fmt.Errorf("time spent cannot exceed 7200 seconds (2 hours), got %d", seconds)
	}

	return TimeSpent{seconds: seconds}, nil
}

// Seconds retorna el valor en segundos
func (t TimeSpent) Seconds() int {
	return t.seconds
}

// Minutes retorna el valor en minutos
func (t TimeSpent) Minutes() int {
	return t.seconds / 60
}

// String implementa Stringer
func (t TimeSpent) String() string {
	minutes := t.seconds / 60
	seconds := t.seconds % 60
	return fmt.Sprintf("%dm%ds", minutes, seconds)
}

// Equals compara dos TimeSpents
func (t TimeSpent) Equals(other TimeSpent) bool {
	return t.seconds == other.seconds
}
