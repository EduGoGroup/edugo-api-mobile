package valueobject

import (
	"github.com/edugo/shared/pkg/types"
)

// MaterialID representa el identificador único de un material educativo
type MaterialID struct {
	value types.UUID
}

// NewMaterialID crea un nuevo MaterialID
func NewMaterialID() MaterialID {
	return MaterialID{value: types.NewUUID()}
}

// MaterialIDFromString crea un MaterialID desde un string
func MaterialIDFromString(s string) (MaterialID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return MaterialID{}, err
	}
	return MaterialID{value: uuid}, nil
}

// String retorna la representación en string
func (m MaterialID) String() string {
	return m.value.String()
}

// UUID retorna el UUID subyacente
func (m MaterialID) UUID() types.UUID {
	return m.value
}

// IsZero verifica si es el valor cero
func (m MaterialID) IsZero() bool {
	return m.value.IsZero()
}

// Equals compara dos MaterialID
func (m MaterialID) Equals(other MaterialID) bool {
	return m.value.String() == other.value.String()
}
