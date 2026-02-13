// Package client proporciona clientes HTTP para comunicación con otros servicios.
// AuthClient valida tokens JWT de forma local usando el mismo secret que api-admin,
// con fallback opcional a validación remota.
package client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/sony/gobreaker"
)

// TokenInfo contiene la información de un token validado
type TokenInfo struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"` // Deprecated: usar ActiveContext
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Error     string    `json:"error,omitempty"`

	// RBAC: Contexto activo del usuario (tokens nuevos)
	ActiveContext *auth.UserContext `json:"active_context,omitempty"`
}

// AuthClientConfig configuración del cliente de autenticación
type AuthClientConfig struct {
	// Configuración para validación LOCAL (preferida)
	JWTSecret string // Secret para validación JWT local (MISMO que api-admin)
	JWTIssuer string // Issuer esperado (debe ser "edugo-central" para compatibilidad con api-admin)

	// Configuración para validación REMOTA (fallback opcional)
	BaseURL         string        // URL base de api-admin (ej: http://localhost:8082)
	Timeout         time.Duration // Timeout para requests HTTP (default: 5s)
	RemoteEnabled   bool          // Habilitar validación remota como fallback
	CircuitBreaker  CircuitBreakerConfig
	FallbackEnabled bool // Si falla validación local, intentar remota

	// Cache
	CacheTTL     time.Duration // TTL del cache de validaciones (default: 60s)
	CacheEnabled bool          // Habilitar cache de validaciones
}

// CircuitBreakerConfig configuración del circuit breaker para llamadas remotas
type CircuitBreakerConfig struct {
	MaxRequests uint32        // Máximo de requests en estado half-open
	Interval    time.Duration // Intervalo para resetear contadores
	Timeout     time.Duration // Tiempo que permanece abierto antes de half-open
}

// AuthClient cliente para validar tokens JWT
// Prioriza validación LOCAL usando el mismo JWT secret que api-admin
// Opcionalmente puede usar validación REMOTA como fallback
type AuthClient struct {
	jwtManager     *auth.JWTManager // Para validación local
	baseURL        string
	httpClient     *http.Client
	cache          *tokenCache
	circuitBreaker *gobreaker.CircuitBreaker
	config         AuthClientConfig
}

// NewAuthClient crea una nueva instancia del cliente de autenticación
// Requiere JWTSecret y JWTIssuer para validación local
func NewAuthClient(config AuthClientConfig) *AuthClient {
	// Valores por defecto
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 60 * time.Second
	}
	if config.JWTIssuer == "" {
		config.JWTIssuer = "edugo-central" // Issuer por defecto compatible con api-admin
	}

	// Crear JWTManager para validación local
	var jwtManager *auth.JWTManager
	if config.JWTSecret != "" {
		jwtManager = auth.NewJWTManager(config.JWTSecret, config.JWTIssuer)
	}

	// Configurar circuit breaker para llamadas remotas (solo si RemoteEnabled)
	var cb *gobreaker.CircuitBreaker
	if config.RemoteEnabled && config.BaseURL != "" {
		cbSettings := gobreaker.Settings{
			Name:        "auth-service",
			MaxRequests: config.CircuitBreaker.MaxRequests,
			Interval:    config.CircuitBreaker.Interval,
			Timeout:     config.CircuitBreaker.Timeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= 3 && failureRatio >= 0.6
			},
			OnStateChange: func(name string, from, to gobreaker.State) {
				fmt.Printf("[AuthClient] Circuit breaker '%s': %s -> %s\n", name, from, to)
			},
		}

		if cbSettings.MaxRequests == 0 {
			cbSettings.MaxRequests = 3
		}
		if cbSettings.Interval == 0 {
			cbSettings.Interval = 10 * time.Second
		}
		if cbSettings.Timeout == 0 {
			cbSettings.Timeout = 30 * time.Second
		}

		cb = gobreaker.NewCircuitBreaker(cbSettings)
	}

	return &AuthClient{
		jwtManager: jwtManager,
		baseURL:    config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		cache:          newTokenCache(config.CacheTTL),
		circuitBreaker: cb,
		config:         config,
	}
}

// ValidateToken valida un token JWT
// Estrategia:
// 1. Si hay cache habilitado, verificar cache primero
// 2. Validar localmente usando JWTManager (preferido)
// 3. Si falla y FallbackEnabled, intentar validación remota
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	// 1. Verificar cache primero
	cacheKey := c.hashToken(token)
	if c.config.CacheEnabled {
		if cached, found := c.cache.Get(cacheKey); found {
			return cached, nil
		}
	}

	// 2. Validar localmente (método preferido)
	if c.jwtManager != nil {
		info, err := c.validateTokenLocally(token)
		if err == nil && info.Valid {
			// Token válido localmente
			if c.config.CacheEnabled {
				c.cache.Set(cacheKey, info)
			}
			return info, nil
		}

		// Si la validación local falló y no hay fallback, retornar el error
		if !c.config.FallbackEnabled || !c.config.RemoteEnabled {
			if info != nil {
				return info, nil
			}
			return &TokenInfo{
				Valid: false,
				Error: fmt.Sprintf("validación local falló: %v", err),
			}, nil
		}
	}

	// 3. Fallback a validación remota (si está habilitada)
	if c.config.RemoteEnabled && c.baseURL != "" && c.circuitBreaker != nil {
		result, err := c.circuitBreaker.Execute(func() (interface{}, error) {
			return c.doValidateTokenRemote(ctx, token)
		})

		if err != nil {
			return &TokenInfo{
				Valid: false,
				Error: fmt.Sprintf("validación remota falló: %v", err),
			}, nil
		}

		info := result.(*TokenInfo)
		if c.config.CacheEnabled && info.Valid {
			c.cache.Set(cacheKey, info)
		}
		return info, nil
	}

	// No hay forma de validar el token
	return &TokenInfo{
		Valid: false,
		Error: "no hay método de validación disponible (JWTSecret no configurado y RemoteEnabled=false)",
	}, nil
}

// validateTokenLocally valida el token usando el JWTManager local
func (c *AuthClient) validateTokenLocally(token string) (*TokenInfo, error) {
	if c.jwtManager == nil {
		return nil, fmt.Errorf("JWTManager no configurado")
	}

	claims, err := c.jwtManager.ValidateToken(token)
	if err != nil {
		return &TokenInfo{
			Valid: false,
			Error: err.Error(),
		}, err
	}

	// Extraer ExpiresAt de los claims
	var expiresAt time.Time
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}

	info := &TokenInfo{
		Valid:         true,
		UserID:        claims.UserID,
		Email:         claims.Email,
		Role:          string(claims.Role),
		ExpiresAt:     expiresAt,
		ActiveContext: claims.ActiveContext,
	}

	// Si hay ActiveContext, extraer role del contexto RBAC
	if claims.ActiveContext != nil {
		info.Role = claims.ActiveContext.RoleName
	}

	return info, nil
}

// doValidateTokenRemote realiza la llamada HTTP a api-admin /v1/auth/verify
func (c *AuthClient) doValidateTokenRemote(ctx context.Context, token string) (*TokenInfo, error) {
	url := c.baseURL + "/v1/auth/verify"

	reqBody := map[string]string{"token": token}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling auth service: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("auth service error: status %d", resp.StatusCode)
	}

	var info TokenInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &info, nil
}

// hashToken genera un hash SHA256 del token para usar como cache key
func (c *AuthClient) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ============================================
// Token Cache Implementation
// ============================================

type tokenCache struct {
	entries map[string]*cacheEntry
	ttl     time.Duration
	mutex   sync.RWMutex
}

type cacheEntry struct {
	info      *TokenInfo
	expiresAt time.Time
}

func newTokenCache(ttl time.Duration) *tokenCache {
	cache := &tokenCache{
		entries: make(map[string]*cacheEntry),
		ttl:     ttl,
	}

	go cache.cleanupLoop()

	return cache
}

func (c *tokenCache) Get(key string) (*TokenInfo, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.info, true
}

func (c *tokenCache) Set(key string, info *TokenInfo) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = &cacheEntry{
		info:      info,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *tokenCache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

func (c *tokenCache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.expiresAt) {
			delete(c.entries, key)
		}
	}
}

// Stats retorna estadísticas del cache
func (c *tokenCache) Stats() (total int, expired int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	now := time.Now()
	total = len(c.entries)
	for _, entry := range c.entries {
		if now.After(entry.expiresAt) {
			expired++
		}
	}
	return total, expired
}
