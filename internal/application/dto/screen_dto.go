package dto

import (
	"encoding/json"
	"time"
)

// CombinedScreenDTO es la respuesta combinada de pantalla para el frontend
// Replica la estructura de screenconfig.CombinedScreenDTO para evitar dependencia de modulo no publicado
type CombinedScreenDTO struct {
	ScreenID        string          `json:"screenId"`
	ScreenKey       string          `json:"screenKey"`
	ScreenName      string          `json:"screenName"`
	Pattern         string          `json:"pattern"`
	Version         int             `json:"version"`
	Template        json.RawMessage `json:"template"`
	DataEndpoint    string          `json:"dataEndpoint,omitempty"`
	DataConfig      json.RawMessage `json:"dataConfig,omitempty"`
	Actions         json.RawMessage `json:"actions"`
	UserPreferences json.RawMessage `json:"userPreferences,omitempty"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

// ResourceScreenDTO representa una pantalla vinculada a un recurso
type ResourceScreenDTO struct {
	ResourceID  string `json:"resource_id"`
	ResourceKey string `json:"resource_key"`
	ScreenKey   string `json:"screen_key"`
	ScreenType  string `json:"screen_type"`
	IsDefault   bool   `json:"is_default"`
}
