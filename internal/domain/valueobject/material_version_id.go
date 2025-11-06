package valueobject

import (
	"github.com/EduGoGroup/edugo-shared/common/types"
)

// MaterialVersionID representa el identificador único de una versión de material educativo
type MaterialVersionID struct {
	value types.UUID
}

// NewMaterialVersionID crea un nuevo MaterialVersionID
func NewMaterialVersionID() MaterialVersionID {
	return MaterialVersionID{value: types.NewUUID()}
}

// MaterialVersionIDFromString crea un MaterialVersionID desde un string
func MaterialVersionIDFromString(s string) (MaterialVersionID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return MaterialVersionID{}, err
	}
	return MaterialVersionID{value: uuid}, nil
}

// String retorna la representación en string
func (mv MaterialVersionID) String() string {
	return mv.value.String()
}

// UUID retorna el UUID subyacente
func (mv MaterialVersionID) UUID() types.UUID {
	return mv.value
}

// IsZero verifica si es el valor cero
func (mv MaterialVersionID) IsZero() bool {
	return mv.value.IsZero()
}

// Equals compara dos MaterialVersionID
func (mv MaterialVersionID) Equals(other MaterialVersionID) bool {
	return mv.value.String() == other.value.String()
}
