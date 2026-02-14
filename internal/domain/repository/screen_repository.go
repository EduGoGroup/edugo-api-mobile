package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CombinedScreen representa el resultado de la consulta JOIN entre template, instancia y preferencias
type CombinedScreen struct {
	ID              string
	ScreenKey       string
	Name            string
	Pattern         string
	Version         int
	Definition      json.RawMessage
	SlotData        json.RawMessage
	Actions         json.RawMessage
	DataEndpoint    string
	DataConfig      json.RawMessage
	UserPreferences json.RawMessage
	LastUpdated     time.Time
}

// ResourceScreenInfo representa la info de una pantalla vinculada a un recurso
type ResourceScreenInfo struct {
	ResourceID  string
	ResourceKey string
	ScreenKey   string
	ScreenType  string
	IsDefault   bool
}

// ScreenRepository define operaciones de lectura para pantallas
type ScreenRepository interface {
	// GetCombinedScreen carga template + instancia + preferencias de usuario en una sola consulta
	GetCombinedScreen(ctx context.Context, screenKey string, userID uuid.UUID) (*CombinedScreen, error)

	// GetScreensForResource retorna todas las configuraciones de pantalla vinculadas a un recurso
	GetScreensForResource(ctx context.Context, resourceKey string) ([]*ResourceScreenInfo, error)

	// GetUserPreferences retorna las preferencias especificas del usuario para una pantalla
	GetUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID) (json.RawMessage, error)

	// SaveUserPreferences guarda las preferencias especificas del usuario
	SaveUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID, prefs json.RawMessage) error
}
