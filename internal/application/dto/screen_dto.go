package dto

import "github.com/EduGoGroup/edugo-shared/screenconfig"

// CombinedScreenDTO is a type alias for the shared library's CombinedScreenDTO.
// This avoids breaking existing code that references dto.CombinedScreenDTO.
// Note: Pattern field is screenconfig.Pattern (a string type), not plain string.
type CombinedScreenDTO = screenconfig.CombinedScreenDTO

// ResourceScreenDTO is a type alias for the shared library's ResourceScreenDTO.
type ResourceScreenDTO = screenconfig.ResourceScreenDTO
