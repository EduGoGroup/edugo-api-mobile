// Package client proporciona clientes HTTP para comunicación con otros servicios.
// AuthClient permite validar tokens JWT contra api-admin como autoridad central.
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

	"github.com/sony/gobreaker"
)

// TokenInfo contiene la información de un token validado por api-admin
type TokenInfo struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// AuthClientConfig configuración del cliente de autenticación
type AuthClientConfig struct {
	BaseURL           string               // URL base de api-admin (ej: http://localhost:8081)
	Timeout           time.Duration        // Timeout para requests HTTP (default: 5s)
	CacheTTL          time.Duration        // TTL del cache de validaciones (default: 60s)
	CacheEnabled      bool                 // Habilitar cache de validaciones
	CircuitBreaker    CircuitBreakerConfig // Configuración del circuit breaker
	FallbackEnabled   bool                 // Habilitar fallback si api-admin no responde
	FallbackJWTSecret string               // Secret para validación local (fallback)
}

// CircuitBreakerConfig configuración del circuit breaker
type CircuitBreakerConfig struct {
	MaxRequests uint32        // Máximo de requests en estado half-open
	Interval    time.Duration // Intervalo para resetear contadores
	Timeout     time.Duration // Tiempo que permanece abierto antes de half-open
}

// AuthClient cliente para validar tokens con api-admin
type AuthClient struct {
	baseURL        string
	httpClient     *http.Client
	cache          *tokenCache
	circuitBreaker *gobreaker.CircuitBreaker
	config         AuthClientConfig
}

// NewAuthClient crea una nueva instancia del cliente de autenticación
func NewAuthClient(config AuthClientConfig) *AuthClient {
	// Valores por defecto
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 60 * time.Second
	}

	// Configurar circuit breaker
	cbSettings := gobreaker.Settings{
		Name:        "auth-service",
		MaxRequests: config.CircuitBreaker.MaxRequests,
		Interval:    config.CircuitBreaker.Interval,
		Timeout:     config.CircuitBreaker.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Abrir circuito si hay 60% de fallos con al menos 3 requests
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			fmt.Printf("[AuthClient] Circuit breaker '%s': %s -> %s\n", name, from, to)
		},
	}

	// Valores por defecto del circuit breaker
	if cbSettings.MaxRequests == 0 {
		cbSettings.MaxRequests = 3
	}
	if cbSettings.Interval == 0 {
		cbSettings.Interval = 10 * time.Second
	}
	if cbSettings.Timeout == 0 {
		cbSettings.Timeout = 30 * time.Second
	}

	return &AuthClient{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		cache:          newTokenCache(config.CacheTTL),
		circuitBreaker: gobreaker.NewCircuitBreaker(cbSettings),
		config:         config,
	}
}

// ValidateToken valida un token JWT con api-admin
// Utiliza cache y circuit breaker para resiliencia
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	// 1. Verificar cache primero
	cacheKey := c.hashToken(token)
	if c.config.CacheEnabled {
		if cached, found := c.cache.Get(cacheKey); found {
			return cached, nil
		}
	}

	// 2. Llamar a api-admin con circuit breaker
	result, err := c.circuitBreaker.Execute(func() (interface{}, error) {
		return c.doValidateToken(ctx, token)
	})

	if err != nil {
		// 3. Fallback si está habilitado y api-admin no responde
		if c.config.FallbackEnabled {
			return c.fallbackValidation(token)
		}
		return &TokenInfo{Valid: false, Error: fmt.Sprintf("auth service error: %v", err)}, nil
	}

	info := result.(*TokenInfo)

	// 4. Guardar en cache si el token es válido
	if c.config.CacheEnabled && info.Valid {
		c.cache.Set(cacheKey, info)
	}

	return info, nil
}

// doValidateToken realiza la llamada HTTP a api-admin /v1/auth/verify
func (c *AuthClient) doValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	url := c.baseURL + "/v1/auth/verify"

	// Preparar request body
	reqBody := map[string]string{"token": token}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Crear request con contexto
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Ejecutar request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling auth service: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Leer response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Verificar status code
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("auth service error: status %d", resp.StatusCode)
	}

	// Parsear response
	var info TokenInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &info, nil
}

// fallbackValidation validación local básica si api-admin no responde
// NOTA: Esta es una validación de emergencia con funcionalidad limitada
func (c *AuthClient) fallbackValidation(token string) (*TokenInfo, error) {
	// En modo fallback, rechazamos el token por seguridad
	// En producción, se podría implementar validación JWT local si se comparte el secret
	return &TokenInfo{
		Valid: false,
		Error: "auth service unavailable, fallback validation denied",
	}, nil
}

// hashToken genera un hash SHA256 del token para usar como cache key
// Esto evita almacenar el token completo en memoria
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

	// Iniciar limpieza periódica de entries expirados
	go cache.cleanupLoop()

	return cache
}

// Get obtiene un entry del cache si existe y no ha expirado
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

// Set guarda un entry en el cache con el TTL configurado
func (c *tokenCache) Set(key string, info *TokenInfo) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = &cacheEntry{
		info:      info,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// cleanupLoop limpia entries expirados cada minuto
func (c *tokenCache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup elimina todos los entries expirados
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

// Stats retorna estadísticas del cache (para métricas)
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
