package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestAuthClient_ValidateToken_Success(t *testing.T) {
	// Crear servidor mock
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/verify" {
			t.Errorf("Path incorrecto: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Método incorrecto: %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type incorrecto: %s", r.Header.Get("Content-Type"))
		}

		response := TokenInfo{
			Valid:  true,
			UserID: "user-123",
			Email:  "test@test.com",
			Role:   "teacher",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Crear cliente
	client := NewAuthClient(AuthClientConfig{
		BaseURL:      server.URL,
		Timeout:      5 * time.Second,
		CacheEnabled: false,
	})

	// Validar token
	info, err := client.ValidateToken(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if !info.Valid {
		t.Error("Token debería ser válido")
	}
	if info.UserID != "user-123" {
		t.Errorf("UserID incorrecto: esperado 'user-123', obtenido '%s'", info.UserID)
	}
	if info.Email != "test@test.com" {
		t.Errorf("Email incorrecto: esperado 'test@test.com', obtenido '%s'", info.Email)
	}
	if info.Role != "teacher" {
		t.Errorf("Role incorrecto: esperado 'teacher', obtenido '%s'", info.Role)
	}
}

func TestAuthClient_ValidateToken_Invalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TokenInfo{
			Valid: false,
			Error: "token expirado",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
	})

	info, err := client.ValidateToken(context.Background(), "expired-token")
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if info.Valid {
		t.Error("Token debería ser inválido")
	}
	if info.Error != "token expirado" {
		t.Errorf("Error incorrecto: esperado 'token expirado', obtenido '%s'", info.Error)
	}
}

func TestAuthClient_Cache_HitAndMiss(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		response := TokenInfo{Valid: true, UserID: "user-123"}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:      server.URL,
		CacheTTL:     5 * time.Second,
		CacheEnabled: true,
	})

	// Primera llamada - debe ir al servidor
	_, _ = client.ValidateToken(context.Background(), "cached-token")
	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("Primera llamada: esperado 1 call, obtenido %d", callCount)
	}

	// Segunda llamada con mismo token - debe usar cache
	_, _ = client.ValidateToken(context.Background(), "cached-token")
	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("Segunda llamada (cache hit): esperado 1 call, obtenido %d", callCount)
	}

	// Tercera llamada con token diferente - debe ir al servidor
	_, _ = client.ValidateToken(context.Background(), "different-token")
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("Tercera llamada (cache miss): esperado 2 calls, obtenido %d", callCount)
	}
}

func TestAuthClient_Cache_Disabled(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		response := TokenInfo{Valid: true, UserID: "user-123"}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false, // Cache deshabilitado
	})

	// Ambas llamadas deben ir al servidor
	_, _ = client.ValidateToken(context.Background(), "token")
	_, _ = client.ValidateToken(context.Background(), "token")

	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("Con cache deshabilitado: esperado 2 calls, obtenido %d", callCount)
	}
}

func TestAuthClient_CircuitBreaker_OpensOnFailures(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		// Siempre retorna error 500
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
		CircuitBreaker: CircuitBreakerConfig{
			MaxRequests: 1,
			Interval:    100 * time.Millisecond,
			Timeout:     500 * time.Millisecond,
		},
	})

	// Hacer varias llamadas que fallarán
	for i := 0; i < 10; i++ {
		_, _ = client.ValidateToken(context.Background(), "test-token")
	}

	// El circuit breaker debería estar abierto después de 3+ fallos
	// Las llamadas después de que se abre no deberían llegar al servidor
	count := atomic.LoadInt32(&callCount)
	if count >= 10 {
		t.Errorf("Circuit breaker no funcionó: se esperaban menos de 10 calls, obtenido %d", count)
	}
	t.Logf("Circuit breaker funcionó: solo %d llamadas llegaron al servidor", count)
}

func TestAuthClient_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simular latencia excesiva
		time.Sleep(500 * time.Millisecond)
		response := TokenInfo{Valid: true}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:      server.URL,
		Timeout:      100 * time.Millisecond, // Timeout corto
		CacheEnabled: false,
	})

	info, _ := client.ValidateToken(context.Background(), "test-token")

	// Debe fallar por timeout
	if info.Valid {
		t.Error("Debería haber fallado por timeout")
	}
	if info.Error == "" {
		t.Error("Debería tener un mensaje de error")
	}
}

func TestAuthClient_ServerError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
		CircuitBreaker: CircuitBreakerConfig{
			MaxRequests: 100, // Alto para que no se abra el circuito
			Interval:    1 * time.Hour,
			Timeout:     1 * time.Hour,
		},
	})

	info, _ := client.ValidateToken(context.Background(), "test-token")

	if info.Valid {
		t.Error("Token no debería ser válido con error 500")
	}
}

func TestAuthClient_FallbackValidation(t *testing.T) {
	// Servidor que siempre falla
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	server.Close() // Cerrar inmediatamente para simular servicio no disponible

	client := NewAuthClient(AuthClientConfig{
		BaseURL:         server.URL,
		CacheEnabled:    false,
		FallbackEnabled: true,
		CircuitBreaker: CircuitBreakerConfig{
			MaxRequests: 1,
			Interval:    10 * time.Millisecond,
			Timeout:     10 * time.Millisecond,
		},
	})

	// Forzar apertura del circuit breaker
	for i := 0; i < 5; i++ {
		_, _ = client.ValidateToken(context.Background(), "token")
	}

	// Ahora con circuit breaker abierto, debería usar fallback
	info, _ := client.ValidateToken(context.Background(), "test-token")

	// Fallback actual rechaza tokens por seguridad
	if info.Valid {
		t.Error("Fallback debería rechazar tokens por seguridad")
	}
	if info.Error == "" {
		t.Error("Fallback debería indicar que el servicio no está disponible")
	}
}

func TestTokenCache_GetSet(t *testing.T) {
	cache := newTokenCache(1 * time.Second)

	// Set
	info := &TokenInfo{Valid: true, UserID: "test-user"}
	cache.Set("key1", info)

	// Get existente
	result, found := cache.Get("key1")
	if !found {
		t.Error("Debería encontrar el entry en cache")
	}
	if result.UserID != "test-user" {
		t.Errorf("UserID incorrecto: %s", result.UserID)
	}

	// Get no existente
	_, found = cache.Get("nonexistent")
	if found {
		t.Error("No debería encontrar entry inexistente")
	}
}

func TestTokenCache_Expiration(t *testing.T) {
	cache := newTokenCache(50 * time.Millisecond) // TTL muy corto

	info := &TokenInfo{Valid: true, UserID: "test-user"}
	cache.Set("key1", info)

	// Debería estar en cache inmediatamente
	_, found := cache.Get("key1")
	if !found {
		t.Error("Debería encontrar el entry recién agregado")
	}

	// Esperar a que expire
	time.Sleep(100 * time.Millisecond)

	// Ya no debería estar
	_, found = cache.Get("key1")
	if found {
		t.Error("Entry debería haber expirado")
	}
}

func TestTokenCache_Stats(t *testing.T) {
	cache := newTokenCache(1 * time.Hour)

	cache.Set("key1", &TokenInfo{Valid: true})
	cache.Set("key2", &TokenInfo{Valid: true})
	cache.Set("key3", &TokenInfo{Valid: true})

	total, expired := cache.Stats()
	if total != 3 {
		t.Errorf("Total incorrecto: esperado 3, obtenido %d", total)
	}
	if expired != 0 {
		t.Errorf("Expired incorrecto: esperado 0, obtenido %d", expired)
	}
}

func TestAuthClient_HashToken(t *testing.T) {
	client := NewAuthClient(AuthClientConfig{
		BaseURL: "http://localhost",
	})

	hash1 := client.hashToken("token123")
	hash2 := client.hashToken("token123")
	hash3 := client.hashToken("different-token")

	// Mismo token = mismo hash
	if hash1 != hash2 {
		t.Error("Mismo token debería producir mismo hash")
	}

	// Diferente token = diferente hash
	if hash1 == hash3 {
		t.Error("Diferente token debería producir diferente hash")
	}

	// Hash debería ser hexadecimal de 64 caracteres (SHA256)
	if len(hash1) != 64 {
		t.Errorf("Hash debería tener 64 caracteres, tiene %d", len(hash1))
	}
}

func TestAuthClient_ConcurrentAccess(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		time.Sleep(10 * time.Millisecond) // Simular latencia
		response := TokenInfo{Valid: true, UserID: "user-123"}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: true,
		CacheTTL:     5 * time.Second,
	})

	// Lanzar múltiples goroutines concurrentes
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := client.ValidateToken(context.Background(), "concurrent-token")
			if err != nil {
				t.Errorf("Error en validación concurrente: %v", err)
			}
			done <- true
		}()
	}

	// Esperar a que todas terminen
	for i := 0; i < 100; i++ {
		<-done
	}

	// Con cache habilitado, no deberían ser 100 llamadas al servidor
	count := atomic.LoadInt32(&callCount)
	t.Logf("Llamadas al servidor con 100 requests concurrentes: %d", count)

	// Debería haber menos llamadas gracias al cache
	// (aunque no exactamente 1 debido a race conditions antes de que se cache)
	if count > 50 {
		t.Logf("Advertencia: Cache podría no estar funcionando óptimamente bajo carga")
	}
}
