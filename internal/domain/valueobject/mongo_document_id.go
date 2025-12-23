package valueobject

import (
	"errors"
	"regexp"
)

// MongoDocumentID representa un ObjectId de MongoDB (24 caracteres hexadecimales)
type MongoDocumentID struct {
	value string
}

var (
	// ErrInvalidMongoDocumentID es retornado cuando el ID no es válido
	ErrInvalidMongoDocumentID = errors.New("mongo document ID must be exactly 24 hexadecimal characters")

	// hexPattern valida que el string contenga solo caracteres hexadecimales
	hexPattern = regexp.MustCompile("^[0-9a-fA-F]{24}$")
)

// NewMongoDocumentID crea un MongoDocumentID válido
func NewMongoDocumentID(value string) (MongoDocumentID, error) {
	if len(value) != 24 {
		return MongoDocumentID{}, ErrInvalidMongoDocumentID
	}

	if !hexPattern.MatchString(value) {
		return MongoDocumentID{}, ErrInvalidMongoDocumentID
	}

	return MongoDocumentID{value: value}, nil
}

// Value retorna el valor del ID
func (m MongoDocumentID) Value() string {
	return m.value
}

// String implementa Stringer
func (m MongoDocumentID) String() string {
	return m.value
}

// Equals compara dos MongoDocumentIDs
func (m MongoDocumentID) Equals(other MongoDocumentID) bool {
	return m.value == other.value
}
