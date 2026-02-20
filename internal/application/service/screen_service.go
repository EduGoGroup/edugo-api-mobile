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
	resolvedDefinition := resolveSlots(combined.Definition, combined.SlotData)

	// 4. Aplicar platformOverrides si se proporciona platform
	if platform != "" {
		resolvedDefinition = applyPlatformOverrides(resolvedDefinition, platform)
	}

	// 5. Construir DTO
	result := &dto.CombinedScreenDTO{
		ScreenID:        combined.ID,
		ScreenKey:       combined.ScreenKey,
		ScreenName:      combined.Name,
		Pattern:         combined.Pattern,
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
	var allowedResources []*repository.MenuResource
	var resourceKeys []string
	for _, res := range resources {
		permKey := res.Key + ":read"
		if hasPermission(permissions, permKey) || res.Scope == "system" {
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

// hasPermission verifica si una lista de permisos contiene un permiso requerido
func hasPermission(perms []string, required string) bool {
	for _, p := range perms {
		if p == required {
			return true
		}
	}
	return false
}

// buildNavigationTree construye los items de navegacion separados en bottomNav y drawer
func buildNavigationTree(resources []*repository.MenuResource, screenMap map[string]string, platform string) ([]NavItemDTO, []NavItemDTO) {
	// Separar top-level (sin parent) de hijos
	topLevel := make([]*repository.MenuResource, 0)
	children := make(map[string][]*repository.MenuResource) // parentID -> children

	for _, res := range resources {
		if res.ParentID == nil {
			topLevel = append(topLevel, res)
		} else {
			children[*res.ParentID] = append(children[*res.ParentID], res)
		}
	}

	// Convertir a NavItemDTO
	var allItems []NavItemDTO
	for _, res := range topLevel {
		item := NavItemDTO{
			Key:       res.Key,
			Label:     res.DisplayName,
			SortOrder: res.SortOrder,
		}
		if res.Icon != nil {
			item.Icon = *res.Icon
		}
		if sk, ok := screenMap[res.Key]; ok {
			item.ScreenKey = sk
		}
		// Agregar hijos
		if kids, ok := children[res.ID]; ok {
			for _, kid := range kids {
				childItem := NavItemDTO{
					Key:       kid.Key,
					Label:     kid.DisplayName,
					SortOrder: kid.SortOrder,
				}
				if kid.Icon != nil {
					childItem.Icon = *kid.Icon
				}
				if sk, ok := screenMap[kid.Key]; ok {
					childItem.ScreenKey = sk
				}
				item.Children = append(item.Children, childItem)
			}
		}
		allItems = append(allItems, item)
	}

	// Separar: mobile obtiene max 5 en bottomNav, el resto en drawer
	// desktop/web obtiene todo en drawer (sin bottom nav)
	maxBottomNav := 5
	if platform == "desktop" || platform == "web" {
		maxBottomNav = 0
	}

	var bottomNav, drawerItems []NavItemDTO
	for i, item := range allItems {
		if i < maxBottomNav {
			bottomNav = append(bottomNav, item)
		} else {
			drawerItems = append(drawerItems, item)
		}
	}

	if bottomNav == nil {
		bottomNav = []NavItemDTO{}
	}
	if drawerItems == nil {
		drawerItems = []NavItemDTO{}
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

// resolveSlots reemplaza las referencias "slot:xxx" en el template con valores de slotData
func resolveSlots(definition json.RawMessage, slotData json.RawMessage) json.RawMessage {
	if len(slotData) == 0 || string(slotData) == "null" || string(slotData) == "{}" {
		return definition
	}

	var slots map[string]interface{}
	if err := json.Unmarshal(slotData, &slots); err != nil {
		return definition
	}

	if len(slots) == 0 {
		return definition
	}

	var defMap interface{}
	if err := json.Unmarshal(definition, &defMap); err != nil {
		return definition
	}

	resolved := resolveValue(defMap, slots)

	result, err := json.Marshal(resolved)
	if err != nil {
		return definition
	}

	return result
}

// resolveValue resuelve recursivamente referencias slot:xxx en un valor JSON
func resolveValue(value interface{}, slots map[string]interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "slot:") {
			slotKey := strings.TrimPrefix(v, "slot:")
			if slotValue, ok := slots[slotKey]; ok {
				return slotValue
			}
		}
		return v
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			result[key] = resolveValue(val, slots)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = resolveValue(val, slots)
		}
		return result
	default:
		return v
	}
}

// applyPlatformOverrides aplica overrides especificos de plataforma al template
func applyPlatformOverrides(definition json.RawMessage, platform string) json.RawMessage {
	var defMap map[string]interface{}
	if err := json.Unmarshal(definition, &defMap); err != nil {
		return definition
	}

	overrides, ok := defMap["platformOverrides"]
	if !ok {
		return definition
	}

	overridesMap, ok := overrides.(map[string]interface{})
	if !ok {
		return definition
	}

	// Resolver override con cadena de fallback (ej: ios -> mobile -> sin override)
	resolvedKey, found := screenconfig.ResolvePlatformOverrideKey(
		screenconfig.Platform(platform), overridesMap,
	)
	if !found {
		return definition
	}

	platformOverride, ok := overridesMap[resolvedKey]
	if !ok {
		return definition
	}

	platformMap, ok := platformOverride.(map[string]interface{})
	if !ok {
		return definition
	}

	// Aplicar overrides de zonas
	if zoneOverrides, ok := platformMap["zones"]; ok {
		if zonesMap, ok := zoneOverrides.(map[string]interface{}); ok {
			if zones, ok := defMap["zones"]; ok {
				if zonesArr, ok := zones.([]interface{}); ok {
					for i, zone := range zonesArr {
						if zoneMap, ok := zone.(map[string]interface{}); ok {
							zoneID, _ := zoneMap["id"].(string)
							if override, ok := zonesMap[zoneID]; ok {
								if overrideMap, ok := override.(map[string]interface{}); ok {
									for k, v := range overrideMap {
										zoneMap[k] = v
									}
									zonesArr[i] = zoneMap
								}
							}
						}
					}
				}
			}
		}
	}

	// Remover platformOverrides del resultado final
	delete(defMap, "platformOverrides")

	result, err := json.Marshal(defMap)
	if err != nil {
		return definition
	}

	return result
}
