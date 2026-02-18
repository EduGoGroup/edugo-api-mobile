package repository

import (
	"context"
)

// MenuResource represents a resource visible in navigation menus
type MenuResource struct {
	ID            string
	Key           string
	DisplayName   string
	Icon          *string
	ParentID      *string
	SortOrder     int
	Scope         string
	IsMenuVisible bool
	IsActive      bool
}

// ResourceScreenMapping maps a resource to its default screen
type ResourceScreenMapping struct {
	ResourceKey string
	ScreenKey   string
	ScreenType  string
	IsDefault   bool
	SortOrder   int
}

// ResourceReader reads resources for navigation building
type ResourceReader interface {
	// GetMenuResources returns active, menu-visible resources ordered by sort_order
	GetMenuResources(ctx context.Context) ([]*MenuResource, error)

	// GetResourceScreenMappings returns default screen mappings for given resource keys
	GetResourceScreenMappings(ctx context.Context, resourceKeys []string) ([]*ResourceScreenMapping, error)
}
