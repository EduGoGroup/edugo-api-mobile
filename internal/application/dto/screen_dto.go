package dto

import "github.com/EduGoGroup/edugo-shared/screenconfig"

// CombinedScreenDTO es un alias al tipo del shared library screenconfig.CombinedScreenDTO.
// Evita romper codigo existente que referencia dto.CombinedScreenDTO.
// Nota: el campo Pattern es screenconfig.Pattern (tipo string), no string plano.
type CombinedScreenDTO = screenconfig.CombinedScreenDTO

// ResourceScreenDTO es un alias al tipo del shared library screenconfig.ResourceScreenDTO.
type ResourceScreenDTO = screenconfig.ResourceScreenDTO
