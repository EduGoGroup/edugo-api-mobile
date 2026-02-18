package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// PostgresScreenRepository implementa repository.ScreenRepository para PostgreSQL
type PostgresScreenRepository struct {
	db *sql.DB
}

// NewPostgresScreenRepository crea una nueva instancia del repositorio
func NewPostgresScreenRepository(db *sql.DB) repository.ScreenRepository {
	return &PostgresScreenRepository{db: db}
}

// GetCombinedScreen carga template + instancia + preferencias de usuario en una sola consulta JOIN
func (r *PostgresScreenRepository) GetCombinedScreen(ctx context.Context, screenKey string, userID uuid.UUID) (*repository.CombinedScreen, error) {
	query := `
		SELECT si.id, si.screen_key, si.name, st.pattern, st.version, st.definition,
		       si.slot_data, si.actions, si.data_endpoint, si.data_config,
		       si.handler_key,
		       COALESCE(sup.preferences, '{}'::jsonb) as user_preferences,
		       GREATEST(si.updated_at, st.updated_at) as last_updated
		FROM ui_config.screen_instances si
		JOIN ui_config.screen_templates st ON si.template_id = st.id
		LEFT JOIN ui_config.screen_user_preferences sup
		    ON sup.screen_instance_id = si.id AND sup.user_id = $2
		WHERE si.screen_key = $1 AND si.is_active = true AND st.is_active = true
	`

	var (
		id              string
		sKey            string
		name            string
		pattern         string
		version         int
		definition      json.RawMessage
		slotData        json.RawMessage
		actions         json.RawMessage
		dataEndpoint    sql.NullString
		dataConfig      json.RawMessage
		handlerKey      sql.NullString
		userPreferences json.RawMessage
		lastUpdated     time.Time
	)

	err := r.db.QueryRowContext(ctx, query, screenKey, userID).Scan(
		&id, &sKey, &name, &pattern, &version, &definition,
		&slotData, &actions, &dataEndpoint, &dataConfig,
		&handlerKey, &userPreferences, &lastUpdated,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: error getting combined screen: %w", err)
	}

	result := &repository.CombinedScreen{
		ID:              id,
		ScreenKey:       sKey,
		Name:            name,
		Pattern:         pattern,
		Version:         version,
		Definition:      definition,
		SlotData:        slotData,
		Actions:         actions,
		DataConfig:      dataConfig,
		UserPreferences: userPreferences,
		LastUpdated:     lastUpdated,
	}

	if dataEndpoint.Valid {
		result.DataEndpoint = dataEndpoint.String
	}
	if handlerKey.Valid {
		result.HandlerKey = &handlerKey.String
	}

	return result, nil
}

// GetScreensForResource retorna todas las configuraciones de pantalla vinculadas a un recurso
func (r *PostgresScreenRepository) GetScreensForResource(ctx context.Context, resourceKey string) ([]*repository.ResourceScreenInfo, error) {
	query := `
		SELECT rs.id, rs.resource_key, rs.screen_key, rs.screen_type, rs.is_default
		FROM ui_config.resource_screens rs
		JOIN ui_config.screen_instances si ON si.screen_key = rs.screen_key
		WHERE rs.resource_key = $1 AND rs.is_active = true AND si.is_active = true
		ORDER BY rs.sort_order
	`

	rows, err := r.db.QueryContext(ctx, query, resourceKey)
	if err != nil {
		return nil, fmt.Errorf("postgres: error getting screens for resource: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []*repository.ResourceScreenInfo
	for rows.Next() {
		info := &repository.ResourceScreenInfo{}
		if err := rows.Scan(&info.ResourceID, &info.ResourceKey, &info.ScreenKey, &info.ScreenType, &info.IsDefault); err != nil {
			return nil, fmt.Errorf("postgres: error scanning resource screen info: %w", err)
		}
		results = append(results, info)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating resource screens: %w", err)
	}

	return results, nil
}

// GetUserPreferences retorna las preferencias especificas del usuario para una pantalla
func (r *PostgresScreenRepository) GetUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID) (json.RawMessage, error) {
	query := `
		SELECT sup.preferences
		FROM ui_config.screen_user_preferences sup
		JOIN ui_config.screen_instances si ON si.id = sup.screen_instance_id
		WHERE si.screen_key = $1 AND sup.user_id = $2
	`

	var prefs json.RawMessage
	err := r.db.QueryRowContext(ctx, query, screenKey, userID).Scan(&prefs)

	if errors.Is(err, sql.ErrNoRows) {
		return json.RawMessage("{}"), nil
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: error getting user preferences: %w", err)
	}

	return prefs, nil
}

// SaveUserPreferences guarda las preferencias especificas del usuario
func (r *PostgresScreenRepository) SaveUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID, prefs json.RawMessage) error {
	query := `
		INSERT INTO ui_config.screen_user_preferences (id, screen_instance_id, user_id, preferences, created_at, updated_at)
		SELECT gen_random_uuid(), si.id, $2, $3, NOW(), NOW()
		FROM ui_config.screen_instances si
		WHERE si.screen_key = $1 AND si.is_active = true
		ON CONFLICT (screen_instance_id, user_id) DO UPDATE SET
			preferences = EXCLUDED.preferences,
			updated_at = NOW()
	`

	result, err := r.db.ExecContext(ctx, query, screenKey, userID, prefs)
	if err != nil {
		return fmt.Errorf("postgres: error saving user preferences: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("postgres: screen instance not found for key %q", screenKey)
	}

	return nil
}
