package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/lib/pq"
)

// PostgresResourceRepository implementa repository.ResourceReader para PostgreSQL
type PostgresResourceRepository struct {
	db *sql.DB
}

// NewPostgresResourceRepository crea una nueva instancia del repositorio de recursos
func NewPostgresResourceRepository(db *sql.DB) repository.ResourceReader {
	return &PostgresResourceRepository{db: db}
}

// GetMenuResources retorna los recursos activos y visibles en menu, ordenados por sort_order
func (r *PostgresResourceRepository) GetMenuResources(ctx context.Context) ([]*repository.MenuResource, error) {
	query := `
		SELECT id::text, key, display_name, icon, parent_id::text, sort_order, scope
		FROM auth.resources
		WHERE is_active = true AND is_menu_visible = true
		ORDER BY sort_order ASC, key ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres: error getting menu resources: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var resources []*repository.MenuResource
	for rows.Next() {
		res := &repository.MenuResource{}
		var parentID sql.NullString
		var icon sql.NullString

		if err := rows.Scan(&res.ID, &res.Key, &res.DisplayName, &icon, &parentID, &res.SortOrder, &res.Scope); err != nil {
			return nil, fmt.Errorf("postgres: error scanning menu resource: %w", err)
		}

		if icon.Valid {
			res.Icon = &icon.String
		}
		if parentID.Valid {
			res.ParentID = &parentID.String
		}

		resources = append(resources, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating menu resources: %w", err)
	}

	return resources, nil
}

// GetResourceScreenMappings retorna los mappings de pantalla default para las claves de recurso dadas
func (r *PostgresResourceRepository) GetResourceScreenMappings(ctx context.Context, resourceKeys []string) ([]*repository.ResourceScreenMapping, error) {
	if len(resourceKeys) == 0 {
		return []*repository.ResourceScreenMapping{}, nil
	}

	query := `
		SELECT resource_key, screen_key, screen_type, is_default, sort_order
		FROM ui_config.resource_screens
		WHERE resource_key = ANY($1) AND is_active = true AND is_default = true
		ORDER BY sort_order ASC
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(resourceKeys))
	if err != nil {
		return nil, fmt.Errorf("postgres: error getting resource screen mappings: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var mappings []*repository.ResourceScreenMapping
	for rows.Next() {
		m := &repository.ResourceScreenMapping{}
		if err := rows.Scan(&m.ResourceKey, &m.ScreenKey, &m.ScreenType, &m.IsDefault, &m.SortOrder); err != nil {
			return nil, fmt.Errorf("postgres: error scanning resource screen mapping: %w", err)
		}
		mappings = append(mappings, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: error iterating resource screen mappings: %w", err)
	}

	return mappings, nil
}
