# Análisis de Review - PR #106: Dynamic UI Fase 1

**Branch:** `feature/dynamic-ui-phase1` → `dev`  
**PR:** https://github.com/EduGoGroup/edugo-api-mobile/pull/106  
**Total Comentarios Copilot:** 7 comentarios  
**Estado:** ✅ PR listo para merge (no hay correcciones críticas)

---

## ⚠️ IMPORTANTE: Comentario #1 de Copilot es FALSO POSITIVO

### Comentario Incorrecto de Copilot

**Ubicación:** `internal/infrastructure/http/router/router.go:191`  
**Comentario:** Cambiar permiso de `PermissionScreensRead` a `PermissionScreensUpdate`

### ❌ Por Qué Está Equivocado Copilot

**1. El permiso sugerido NO EXISTE:**

Según `edugo-shared/common/types/enum/permission.go`, los permisos disponibles son:

```go
// Screen instances
PermissionScreenInstancesRead   Permission = "screen_instances:read"
PermissionScreenInstancesCreate Permission = "screen_instances:create"
PermissionScreenInstancesUpdate Permission = "screen_instances:update"  ← EXISTE
PermissionScreenInstancesDelete Permission = "screen_instances:delete"

// Screens (solo lectura combinada)
PermissionScreensRead Permission = "screens:read"

// ❌ NO EXISTE: PermissionScreensUpdate
// ❌ NO EXISTE: PermissionScreensWrite
```

**2. El código actual es CORRECTO:**

```go
// Código actual en router.go (CORRECTO ✅)
screens.PUT("/:screenKey/preferences",
    middleware.RequirePermission(enum.PermissionScreenInstancesUpdate),  // ✅ CORRECTO
    c.Handlers.ScreenHandler.SavePreferences,
)
```

**3. Justificación:**

- `SavePreferences` guarda preferencias en tabla `screen_user_preferences`
- Esta tabla está relacionada con `screen_instances` (no con `screen_templates`)
- El permiso `PermissionScreenInstancesUpdate` es el apropiado para modificar instancias de pantalla
- `PermissionScreensRead` es solo lectura (usado correctamente en endpoints GET)

### ✅ Conclusión

**NO HAY CORRECCIÓN NECESARIA.** El código actual usa el permiso correcto. Copilot sugirió un permiso que no existe en el sistema.

---

## 📋 Resumen de Comentarios de Copilot

| # | Tipo | Descripción | Acción |
|---|------|-------------|--------|
| 1 | ❌ Falso Positivo | Cambiar permiso a PermissionScreensUpdate (no existe) | DESCARTAR |
| 2 | 🔵 Documentación | Agregar nota en Swagger sobre navegación estática | Opcional (Fase futura) |
| 3 | 🟡 Deuda Técnica | Cache sin eviction automática | DT-001 (Fase 2) |
| 4 | 🟡 Deuda Técnica | TTL hardcodeado | DT-002 (Fase 2) |
| 5 | 🟡 Deuda Técnica | Validación de preferencias | DT-003 (Fase 2) |
| 6 | 🟡 Deuda Técnica | MD5 → SHA-256 para ETags | DT-004 (Fase 3) |
| 7 | ⚪ Preferencia | Role como string vacío vs nil | Descartar (funciona bien) |

---

## 🎯 Deuda Técnica Identificada

Este documento contiene los tickets de deuda técnica válidos sugeridos a partir de los comentarios de review.

**Total de Items:** 4 tickets  
**Esfuerzo Total Estimado:** 6.5 - 9.5 horas  
**Sprint Sugerido:** Fase 2-3 de Dynamic UI

---

## DT-001: Implementar Cache con Eviction Automática

### Metadata
- **ID:** DT-001
- **Título:** Implementar cache con eviction automática para ScreenService
- **Prioridad:** Media
- **Esfuerzo:** 3-5 horas
- **Sprint Sugerido:** Fase 2 o Fase 3
- **Labels:** `tech-debt`, `performance`, `enhancement`

### Problema

**Archivo:** `internal/application/service/screen_service.go:67`

La cache en memoria (`map[string]*screenCache`) no tiene mecanismo de limpieza de entradas expiradas. Con el tiempo, esto puede causar memory leaks ya que las entradas expiradas nunca se eliminan del mapa.

**Código Actual:**
```go
type screenService struct {
    repo   repository.ScreenRepository
    logger logger.Logger

    mu    sync.RWMutex
    cache map[string]*screenCache  // ← Crece indefinidamente
    ttl   time.Duration
}
```

**Problema:**
- Entradas expiradas nunca se eliminan
- El mapa crece indefinidamente en memoria
- Potencial memory leak en entornos de alta carga

### Solución Propuesta

**Opción A (Recomendada):** Usar biblioteca de cache con TTL automático

```bash
go get github.com/patrickmn/go-cache
```

```go
import "github.com/patrickmn/go-cache"

type screenService struct {
    repo   repository.ScreenRepository
    logger logger.Logger
    cache  *cache.Cache  // ✅ Cache con eviction automática
}

func NewScreenService(repo repository.ScreenRepository, logger logger.Logger) ScreenService {
    return &screenService{
        repo:   repo,
        logger: logger,
        cache:  cache.New(1*time.Hour, 10*time.Minute),  // TTL=1h, cleanup=10min
    }
}

func (s *screenService) getCached(key string) *dto.CombinedScreenDTO {
    if cached, found := s.cache.Get(key); found {
        return cached.(*dto.CombinedScreenDTO)
    }
    return nil
}

func (s *screenService) setCache(key string, data *dto.CombinedScreenDTO) {
    s.cache.Set(key, data, cache.DefaultExpiration)
}
```

**Opción B:** Implementar goroutine de limpieza manual

```go
func NewScreenService(repo repository.ScreenRepository, logger logger.Logger) ScreenService {
    svc := &screenService{
        repo:   repo,
        logger: logger,
        cache:  make(map[string]*screenCache),
        ttl:    1 * time.Hour,
    }
    
    // Goroutine de limpieza cada 10 minutos
    go svc.cleanupExpiredEntries()
    
    return svc
}

func (s *screenService) cleanupExpiredEntries() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        s.mu.Lock()
        now := time.Now()
        for key, entry := range s.cache {
            if now.After(entry.expiresAt) {
                delete(s.cache, key)
                s.logger.Debug("cache entry evicted", "key", key)
            }
        }
        s.mu.Unlock()
    }
}
```

### Criterios de Aceptación

- [ ] Cache no crece indefinidamente
- [ ] Entradas expiradas se eliminan automáticamente
- [ ] No hay memory leaks después de 24h de operación
- [ ] Tests de integración verifican limpieza de cache
- [ ] Métricas de uso de memoria son estables

### Testing

```go
func TestScreenService_CacheEviction(t *testing.T) {
    // Mock con TTL de 100ms para test rápido
    svc := NewScreenService(mockRepo, mockLogger)
    
    // Cachear 100 pantallas
    for i := 0; i < 100; i++ {
        svc.setCache(fmt.Sprintf("screen-%d", i), testScreen)
    }
    
    // Esperar 200ms (suficiente para expirar)
    time.Sleep(200 * time.Millisecond)
    
    // Verificar que cache está vacía después de eviction
    assert.Equal(t, 0, svc.cache.Len())
}
```

---

## DT-002: Mover TTL de Cache a Configuración

### Metadata
- **ID:** DT-002
- **Título:** Mover TTL de cache a configuración externa
- **Prioridad:** Baja
- **Esfuerzo:** 1 hora
- **Sprint Sugerido:** Fase 2-3 (puede agruparse con DT-001)
- **Labels:** `tech-debt`, `config`, `enhancement`

### Problema

**Archivo:** `internal/application/service/screen_service.go:67`

El TTL de cache está hardcodeado a 1 hora en el código. Esto dificulta ajustar el rendimiento sin recompilar.

**Código Actual:**
```go
func NewScreenService(repo repository.ScreenRepository, logger logger.Logger) ScreenService {
    return &screenService{
        // ...
        ttl:    1 * time.Hour,  // ← Hardcodeado
    }
}
```

### Solución Propuesta

**1. Agregar configuración en config.yaml:**

```yaml
# config/config.yaml
cache:
  screen:
    ttl: 1h           # Producción
    cleanup: 10m      # Intervalo de limpieza

# config/config-dev.yaml
cache:
  screen:
    ttl: 5m           # Desarrollo (cache corto para testing)
    cleanup: 1m
```

**2. Actualizar struct de config:**

```go
// internal/config/config.go
type Config struct {
    // ... campos existentes
    Cache CacheConfig `yaml:"cache"`
}

type CacheConfig struct {
    Screen ScreenCacheConfig `yaml:"screen"`
}

type ScreenCacheConfig struct {
    TTL     time.Duration `yaml:"ttl"`
    Cleanup time.Duration `yaml:"cleanup"`
}
```

**3. Pasar configuración al servicio:**

```go
// internal/container/services.go
func (c *Container) ScreenService() service.ScreenService {
    return service.NewScreenService(
        c.ScreenRepository(),
        c.Logger(),
        c.Config.Cache.Screen.TTL,      // ← Desde config
        c.Config.Cache.Screen.Cleanup,  // ← Desde config
    )
}
```

**4. Actualizar constructor:**

```go
func NewScreenService(
    repo repository.ScreenRepository,
    logger logger.Logger,
    cacheTTL time.Duration,
    cleanupInterval time.Duration,
) ScreenService {
    return &screenService{
        repo:   repo,
        logger: logger,
        cache:  cache.New(cacheTTL, cleanupInterval),
    }
}
```

### Criterios de Aceptación

- [ ] TTL de cache se lee desde config.yaml
- [ ] Intervalo de limpieza se lee desde config
- [ ] Diferentes valores por ambiente (dev, qa, prod)
- [ ] Tests usan valores de TTL cortos (100ms) para rapidez
- [ ] Documentación actualizada en SETUP.md

---

## DT-003: Validar Schema de Preferencias de Usuario

### Metadata
- **ID:** DT-003
- **Título:** Agregar validación de schema para preferencias de pantalla
- **Prioridad:** Media
- **Esfuerzo:** 2-3 horas
- **Sprint Sugerido:** Fase 2
- **Labels:** `tech-debt`, `validation`, `api`

### Problema

**Archivo:** `internal/infrastructure/http/handler/screen_handler.go:179`

No se valida que el JSON en `prefs` tenga una estructura válida antes de guardarlo. Aunque `ShouldBindJSON` valida sintaxis JSON, no valida contenido o estructura.

**Código Actual:**
```go
var prefs json.RawMessage
if err := c.ShouldBindJSON(&prefs); err != nil {
    // Solo valida que sea JSON válido, no la estructura
    c.JSON(http.StatusBadRequest, ErrorResponse{...})
    return
}

// Se guarda sin validar contenido ❌
if err := h.service.SaveUserPreferences(ctx, screenKey, userID, prefs); err != nil {
    // ...
}
```

### Solución Propuesta

**Opción A (Recomendada):** Validar que sea objeto JSON no vacío

```go
var prefs json.RawMessage
if err := c.ShouldBindJSON(&prefs); err != nil {
    h.logger.Warn("invalid json syntax", "error", err)
    c.JSON(http.StatusBadRequest, ErrorResponse{
        Error: "JSON inválido",
        Code:  "INVALID_JSON",
    })
    return
}

// Validar que sea objeto JSON no vacío
var prefsObj map[string]interface{}
if err := json.Unmarshal(prefs, &prefsObj); err != nil {
    h.logger.Warn("invalid preferences structure", "error", err)
    c.JSON(http.StatusBadRequest, ErrorResponse{
        Error: "las preferencias deben ser un objeto JSON",
        Code:  "INVALID_PREFERENCES_STRUCTURE",
    })
    return
}

if len(prefsObj) == 0 {
    h.logger.Warn("empty preferences payload", "screen_key", screenKey, "user_id", userIDStr)
    c.JSON(http.StatusBadRequest, ErrorResponse{
        Error: "las preferencias no pueden estar vacías",
        Code:  "EMPTY_PREFERENCES",
    })
    return
}
```

**Opción B (Más estricta):** JSON Schema Validation

```bash
go get github.com/xeipuuv/gojsonschema
```

```go
// Define schema esperado
const preferencesSchema = `{
  "type": "object",
  "properties": {
    "theme": {"type": "string", "enum": ["light", "dark"]},
    "layout": {"type": "string"},
    "filters": {"type": "object"}
  },
  "required": ["theme"],
  "additionalProperties": true
}`

func validatePreferences(prefs json.RawMessage) error {
    schema := gojsonschema.NewStringLoader(preferencesSchema)
    doc := gojsonschema.NewBytesLoader(prefs)
    
    result, err := gojsonschema.Validate(schema, doc)
    if err != nil {
        return err
    }
    
    if !result.Valid() {
        return fmt.Errorf("invalid preferences: %v", result.Errors())
    }
    
    return nil
}
```

### Decisiones de Negocio Requeridas

Antes de implementar, definir:

1. **¿Las preferencias pueden estar vacías?**
   - Opción A: No, deben tener al menos 1 campo
   - Opción B: Sí, `{}` es válido (borrar preferencias)

2. **¿Qué campos son válidos en preferencias?**
   - Opción A: Cualquier JSON válido (flexible)
   - Opción B: Solo campos predefinidos (theme, layout, filters, etc.)

3. **¿Se valida contra schema estricto?**
   - Opción A: Sí, JSON Schema formal
   - Opción B: Solo validación básica (objeto no vacío)

### Criterios de Aceptación

- [ ] Preferencias vacías son rechazadas (si negocio lo requiere)
- [ ] JSON inválido retorna 400 con mensaje claro
- [ ] Arrays y primitivos son rechazados (solo objetos)
- [ ] Tests de validación para casos edge
- [ ] Documentación Swagger actualizada con ejemplos

---

## DT-004: Migrar ETags de MD5 a SHA-256

### Metadata
- **ID:** DT-004
- **Título:** Migrar generación de ETags de MD5 a SHA-256
- **Prioridad:** Baja
- **Esfuerzo:** 30 minutos
- **Sprint Sugerido:** Fase 3 o Backlog
- **Labels:** `tech-debt`, `security`, `best-practices`

### Problema

**Archivo:** `internal/infrastructure/http/handler/screen_handler.go:4`

Se usa MD5 para generar ETags. Aunque MD5 es aceptable para ETags (no es uso criptográfico), genera warnings de seguridad y no sigue mejores prácticas modernas.

**Código Actual:**
```go
import (
    "crypto/md5"  // ← Deprecado
    // ...
)

func generateETag(data []byte) string {
    hash := md5.Sum(data)
    return fmt.Sprintf("\"%x\"", hash)
}
```

### Solución Propuesta

**Cambiar a SHA-256:**

```go
import (
    "crypto/sha256"  // ✅ Algoritmo moderno
    // ...
)

func generateETag(data []byte) string {
    hash := sha256.Sum256(data)
    // Truncar a 16 bytes para ETags más cortos (opcional)
    return fmt.Sprintf("\"%x\"", hash[:16])
}
```

### Consideraciones de Compatibilidad

**IMPORTANTE:** Cambiar el algoritmo de hash invalidará todos los ETags existentes en clientes.

**Impacto:**
- Clientes con ETags cacheados recibirán respuestas 200 con datos completos
- Primer request después del cambio será full response (no 304 Not Modified)
- Requests subsecuentes funcionarán normalmente con nuevos ETags

**Mitigación:**
- Desplegar cambio en horario de baja carga
- Notificar a equipo frontend sobre invalidación de cache
- Considerar estrategia de migración gradual (opcional):

```go
func generateETag(data []byte, useLegacy bool) string {
    if useLegacy {
        hash := md5.Sum(data)
        return fmt.Sprintf("\"md5-%x\"", hash)
    }
    hash := sha256.Sum256(data)
    return fmt.Sprintf("\"sha256-%x\"", hash[:16])
}

// Aceptar ambos formatos durante período de transición
func matchesETag(requestETag string, currentData []byte) bool {
    if strings.HasPrefix(requestETag, "\"md5-") {
        return requestETag == generateETag(currentData, true)
    }
    return requestETag == generateETag(currentData, false)
}
```

### Criterios de Aceptación

- [ ] ETags generados con SHA-256
- [ ] Tests actualizados con nuevos valores de hash
- [ ] Documentación de API actualizada
- [ ] Equipo frontend notificado sobre invalidación de cache
- [ ] No hay warnings de MD5 deprecado en builds

### Testing

```go
func TestGenerateETag_SHA256(t *testing.T) {
    data := []byte(`{"screenKey": "test"}`)
    etag := generateETag(data)
    
    // Verificar formato
    assert.True(t, strings.HasPrefix(etag, "\""))
    assert.True(t, strings.HasSuffix(etag, "\""))
    
    // Verificar que usa SHA-256 (hash más largo que MD5)
    assert.Greater(t, len(etag), 34) // SHA-256 truncado a 16 bytes = 32 chars + 2 quotes
    
    // Verificar consistencia
    etag2 := generateETag(data)
    assert.Equal(t, etag, etag2)
    
    // Verificar que datos diferentes generan ETags diferentes
    etag3 := generateETag([]byte(`{"screenKey": "other"}`))
    assert.NotEqual(t, etag, etag3)
}
```

---

## Resumen y Priorización

### Por Prioridad

| Prioridad | Tickets | Esfuerzo Total |
|-----------|---------|----------------|
| **Media** | DT-001, DT-003 | 5-8 horas |
| **Baja** | DT-002, DT-004 | 1.5 horas |

### Por Sprint Sugerido

| Sprint | Tickets | Descripción |
|--------|---------|-------------|
| **Fase 2** | DT-001, DT-002, DT-003 | Mejoras de arquitectura y validación |
| **Fase 3** | DT-004 | Mejoras de seguridad (best practices) |

### Recomendación de Implementación

**Sprint Fase 2:**
1. Implementar DT-001 + DT-002 juntos (cache con config)
2. Implementar DT-003 después de definir requisitos de negocio

**Sprint Fase 3:**
3. Implementar DT-004 (baja prioridad, puede posponerse)

---

## Plantilla de Issue para GitHub

```markdown
## [DT-XXX] Título del Ticket

**Origen:** PR #106 - Review comment de GitHub Copilot  
**Prioridad:** Media/Baja  
**Esfuerzo:** X horas  
**Sprint Sugerido:** Fase X  

### Problema
[Descripción del problema actual]

### Solución Propuesta
[Descripción de la solución]

### Criterios de Aceptación
- [ ] Item 1
- [ ] Item 2

### Referencias
- PR: #106
- Copilot Comment: #XXXXXXX
- Archivo: `path/to/file.go:line`
```

---

**Generado por:** Claude Code (Análisis de PR)  
**Fecha:** 2026-02-14  
**Versión:** 1.0
