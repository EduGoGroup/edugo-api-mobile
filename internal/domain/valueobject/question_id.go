package valueobject

import "errors"

// QuestionID representa el ID de una pregunta en MongoDB
type QuestionID struct {
	value string
}

var (
	// ErrInvalidQuestionID es retornado cuando el ID no es válido
	ErrInvalidQuestionID = errors.New("question ID cannot be empty")
)

// NewQuestionID crea un QuestionID válido
func NewQuestionID(value string) (QuestionID, error) {
	if value == "" {
		return QuestionID{}, ErrInvalidQuestionID
	}

	return QuestionID{value: value}, nil
}

// Value retorna el valor del ID
func (q QuestionID) Value() string {
	return q.value
}

// String implementa Stringer
func (q QuestionID) String() string {
	return q.value
}

// Equals compara dos QuestionIDs
func (q QuestionID) Equals(other QuestionID) bool {
	return q.value == other.value
}
