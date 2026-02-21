package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/EduGoGroup/edugo-shared/screenconfig"
)

// ScreenService define las operaciones de negocio para pantallas
type ScreenService interface {
	GetScreen(ctx context.Context, screenKey string, userID uuid.UUID, platform string) (*dto.CombinedScreenDTO, error)
	GetNavigationConfig(ctx context.Context, userID uuid.UUID, platform string, permissions []string) (*NavigationConfigDTO, error)
	SaveUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID, prefs json.RawMessage) error
	GetScreensForResource(ctx context.Context, resourceKey string) ([]*dto.ResourceScreenDTO, error)
}

// NavigationConfigDTO representa la configuracion de navegacion
type NavigationConfigDTO struct {
	BottomNav   []NavItemDTO `json:"bottomNav"`
	DrawerItems []NavItemDTO `json:"drawerItems"`
	Version     int          `json:"version"`
}

// NavItemDTO representa un item de navegacion
type NavItemDTO struct {
	Key       string       `json:"key"`
	Label     string       `json:"label"`
	Icon      string       `json:"icon,omitempty"`
	ScreenKey string       `json:"screenKey,omitempty"`
	SortOrder int          `json:"sortOrder"`
	Children  []NavItemDTO `json:"children,omitempty"`
}

// screenCache es una entrada de cache con TTL
type screenCache struct {
	data      *dto.CombinedScreenDTO
	expiresAt time.Time
}

type screenService struct {
	repo           repository.ScreenRepository
	resourceReader repository.ResourceReader
	logger         logger.Logger

	mu    sync.RWMutex
	cache map[string]*screenCache
	ttl   time.Duration
}

// NewScreenService crea una nueva instancia del servicio de pantallas
func NewScreenService(
	repo repository.ScreenRepository,
	resourceReader repository.ResourceReader,
	logger logger.Logger,
) ScreenService {
	return &screenService{
		repo:           repo,
		resourceReader: resourceReader,
		logger:         logger,
		cache:          make(map[string]*screenCache),
		ttl:            1 * time.Hour,
	}
}

// GetScreen retorna la definicion de pantalla combinada para el renderizado del frontend
func (s *screenService) GetScreen(ctx context.Context, screenKey string, userID uuid.UUID, platform string) (*dto.CombinedScreenDTO, error) {
	// 1. Verificar cache (incluyendo platform para evitar servir overrides incorrectos)
	cacheKey := fmt.Sprintf("screen:%s:platform:%s", screenKey, platform)
	if cached := s.getCached(cacheKey); cached != nil {
		s.logger.Info("screen cache hit", "screen_key", screenKey)
		// Fusionar preferencias de usuario sobre la copia cacheada
		return s.withUserPreferences(ctx, cached, screenKey, userID)
	}

	// 2. Cargar desde BD
	combined, err := s.repo.GetCombinedScreen(ctx, screenKey, userID)
	if err != nil {
		s.logger.Error("failed to get combined screen", "screen_key", screenKey, "error", err)
		return nil, errors.NewDatabaseError("get screen", err)
	}
	if combined == nil {
		return nil, errors.NewNotFoundError("screen")
	}

	// 3. Resolver referencias de slot
	resolvedDefinition := screenconfig.ResolveSlots(combined.Definition, combined.SlotData)

	// 4. Aplicar platformOverrides si se proporciona platform
	if platform != "" {
		resolvedDefinition = screenconfig.ApplyPlatformOverrides(resolvedDefinition, platform)
	}

	// 5. Construir DTO
	result := &dto.CombinedScreenDTO{
		ScreenID:        combined.ID,
		ScreenKey:       combined.ScreenKey,
		ScreenName:      combined.Name,
		Pattern:         screenconfig.Pattern(combined.Pattern),
		Version:         combined.Version,
		Template:        resolvedDefinition,
		DataEndpoint:    combined.DataEndpoint,
		DataConfig:      combined.DataConfig,
		Actions:         combined.Actions,
		HandlerKey:      combined.HandlerKey,
		UserPreferences: combined.UserPreferences,
		UpdatedAt:       combined.LastUpdated,
	}

	// 6. Cachear resultado (sin preferencias de usuario para compartir entre usuarios)
	s.setCache(cacheKey, result)

	s.logger.Info("screen loaded from database",
		"screen_key", screenKey,
		"pattern", combined.Pattern,
		"version", combined.Version,
	)

	return result, nil
}

// GetNavigationConfig retorna la estructura de navegacion dinamica basada en permisos del usuario
func (s *screenService) GetNavigationConfig(ctx context.Context, userID uuid.UUID, platform string, permissions []string) (*NavigationConfigDTO, error) {
	// 1. Obtener todos los recursos visibles en menu
	resources, err := s.resourceReader.GetMenuResources(ctx)
	if err != nil {
		s.logger.Error("failed to get menu resources", "error", err)
		return nil, errors.NewDatabaseError("get menu resources", err)
	}

	if len(resources) == 0 {
		return &NavigationConfigDTO{
			BottomNav:   []NavItemDTO{},
			DrawerItems: []NavItemDTO{},
			Version:     1,
		}, nil
	}

	// 2. Filtrar recursos por permisos del usuario
	// 2a. Determinar qué resource keys tienen algún permiso del usuario
	userResourceKeys := make(map[string]bool)
	for _, res := range resources {
		permPrefix := res.Key + ":"
		for _, perm := range permissions {
			if strings.HasPrefix(perm, permPrefix) {
				userResourceKeys[res.Key] = true
				break
			}
		}
		// System scope resources are always visible
		if res.Scope == "system" {
			userResourceKeys[res.Key] = true
		}
	}

	// 2b. Build lookup maps for parent traversal
	resourceByKey := make(map[string]*repository.MenuResource)
	resourceByID := make(map[string]*repository.MenuResource)
	for _, res := range resources {
		resourceByKey[res.Key] = res
		if res.ID != "" {
			resourceByID[res.ID] = res
		}
	}

	// 2c. Include parent resources (walk up the hierarchy)
	allowedKeys := make(map[string]bool)
	for key := range userResourceKeys {
		allowedKeys[key] = true
		res := resourceByKey[key]
		for res != nil && res.ParentID != nil {
			parent := resourceByID[*res.ParentID]
			if parent != nil {
				allowedKeys[parent.Key] = true
				res = parent
			} else {
				break
			}
		}
	}

	// 2d. Filter resources to only allowed ones
	var allowedResources []*repository.MenuResource
	var resourceKeys []string
	for _, res := range resources {
		if allowedKeys[res.Key] {
			allowedResources = append(allowedResources, res)
			resourceKeys = append(resourceKeys, res.Key)
		}
	}

	if len(allowedResources) == 0 {
		return &NavigationConfigDTO{
			BottomNav:   []NavItemDTO{},
			DrawerItems: []NavItemDTO{},
			Version:     1,
		}, nil
	}

	// 3. Obtener mappings de pantalla para los recursos permitidos
	mappings, err := s.resourceReader.GetResourceScreenMappings(ctx, resourceKeys)
	if err != nil {
		s.logger.Error("failed to get resource screen mappings", "error", err)
		return nil, errors.NewDatabaseError("get resource screen mappings", err)
	}

	screenMap := make(map[string]string) // resourceKey -> screenKey
	for _, m := range mappings {
		if m.IsDefault {
			screenMap[m.ResourceKey] = m.ScreenKey
		}
	}

	// 4. Construir arbol de navegacion
	bottomNav, drawerItems := buildNavigationTree(allowedResources, screenMap, platform)

	return &NavigationConfigDTO{
		BottomNav:   bottomNav,
		DrawerItems: drawerItems,
		Version:     1,
	}, nil
}

// toMenuNodes convierte MenuResource del repositorio a screenconfig.MenuNode del shared.
// Los campos opcionales Icon y ParentID se asignan solo si no son nil.
func toMenuNodes(resources []*repository.MenuResource) []screenconfig.MenuNode {
	nodes := make([]screenconfig.MenuNode, len(resources))
	for i, res := range resources {
		node := screenconfig.MenuNode{
			ID:          res.ID,
			Key:         res.Key,
			DisplayName: res.DisplayName,
			SortOrder:   res.SortOrder,
			Scope:       res.Scope,
		}
		if res.Icon != nil {
			node.Icon = *res.Icon
		}
		if res.ParentID != nil {
			node.ParentID = *res.ParentID
		}
		nodes[i] = node
	}
	return nodes
}

// toNavItems convierte screenconfig.MenuTreeItem a NavItemDTO especifico de mobile.
// El campo DisplayName se mapea a Label para mantener el contrato JSON del endpoint.
func toNavItems(items []screenconfig.MenuTreeItem) []NavItemDTO {
	result := make([]NavItemDTO, len(items))
	for i, item := range items {
		navItem := NavItemDTO{
			Key:       item.Key,
			Label:     item.DisplayName,
			Icon:      item.Icon,
			ScreenKey: item.ScreenKey,
			SortOrder: item.SortOrder,
		}
		if len(item.Children) > 0 {
			navItem.Children = toNavItems(item.Children)
		}
		result[i] = navItem
	}
	return result
}

// buildNavigationTree construye los items de navegacion separados en bottomNav y drawer.
// Usa screenconfig.BuildMenuTree para la construccion del arbol, luego aplica logica de
// separacion especifica de mobile (max 5 bottomNav para mobile, todo en drawer para desktop/web).
func buildNavigationTree(resources []*repository.MenuResource, screenMap map[string]string, platform string) ([]NavItemDTO, []NavItemDTO) {
	nodes := toMenuNodes(resources)
	treeItems := screenconfig.BuildMenuTree(nodes, nil, screenMap)
	allItems := toNavItems(treeItems)

	// Separar: mobile obtiene max 5 en bottomNav, el resto en drawer.
	// desktop/web obtiene todo en drawer (sin bottom nav).
	maxBottomNav := 5
	if platform == "desktop" || platform == "web" {
		maxBottomNav = 0
	}

	bottomNav := make([]NavItemDTO, 0)
	drawerItems := make([]NavItemDTO, 0)
	for i, item := range allItems {
		if i < maxBottomNav {
			bottomNav = append(bottomNav, item)
		} else {
			drawerItems = append(drawerItems, item)
		}
	}

	return bottomNav, drawerItems
}

// SaveUserPreferences almacena las preferencias de pantalla especificas del usuario
func (s *screenService) SaveUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID, prefs json.RawMessage) error {
	if err := s.repo.SaveUserPreferences(ctx, screenKey, userID, prefs); err != nil {
		s.logger.Error("failed to save user preferences",
			"screen_key", screenKey,
			"user_id", userID.String(),
			"error", err,
		)
		return errors.NewDatabaseError("save user preferences", err)
	}

	s.logger.Info("user preferences saved",
		"screen_key", screenKey,
		"user_id", userID.String(),
	)

	return nil
}

// GetScreensForResource retorna las pantallas vinculadas a un recurso
func (s *screenService) GetScreensForResource(ctx context.Context, resourceKey string) ([]*dto.ResourceScreenDTO, error) {
	infos, err := s.repo.GetScreensForResource(ctx, resourceKey)
	if err != nil {
		s.logger.Error("failed to get screens for resource", "resource_key", resourceKey, "error", err)
		return nil, errors.NewDatabaseError("get screens for resource", err)
	}

	dtos := make([]*dto.ResourceScreenDTO, len(infos))
	for i, info := range infos {
		dtos[i] = &dto.ResourceScreenDTO{
			ResourceID:  info.ResourceID,
			ResourceKey: info.ResourceKey,
			ScreenKey:   info.ScreenKey,
			ScreenType:  info.ScreenType,
			IsDefault:   info.IsDefault,
		}
	}

	return dtos, nil
}

// getCached obtiene una entrada del cache si existe y no ha expirado
func (s *screenService) getCached(key string) *dto.CombinedScreenDTO {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.cache[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil
	}

	return entry.data
}

// setCache almacena una entrada en el cache con TTL
func (s *screenService) setCache(key string, data *dto.CombinedScreenDTO) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[key] = &screenCache{
		data:      data,
		expiresAt: time.Now().Add(s.ttl),
	}
}

// withUserPreferences fusiona las preferencias del usuario sobre una pantalla cacheada
func (s *screenService) withUserPreferences(ctx context.Context, cached *dto.CombinedScreenDTO, screenKey string, userID uuid.UUID) (*dto.CombinedScreenDTO, error) {
	prefs, err := s.repo.GetUserPreferences(ctx, screenKey, userID)
	if err != nil {
		s.logger.Warn("failed to get user preferences for cached screen",
			"screen_key", screenKey,
			"error", err,
		)
		return cached, nil
	}

	// Crear copia con las preferencias del usuario
	result := *cached
	result.UserPreferences = prefs
	return &result, nil
}
