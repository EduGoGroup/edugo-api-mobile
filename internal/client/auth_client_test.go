package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

const (
	testJWTSecret = "test-secret-key-at-least-32-chars-long-for-security"
	testJWTIssuer = "edugo-central"
)

// generateTestToken genera un token JWT válido para pruebas con contexto RBAC
func generateTestToken(t *testing.T, userID, email string, role enum.SystemRole, expiresIn time.Duration) string {
	t.Helper()
	manager := auth.NewJWTManager(testJWTSecret, testJWTIssuer)

	// Crear contexto RBAC básico para tests
	activeContext := &auth.UserContext{
		RoleID:      "role-" + string(role),
		RoleName:    string(role),
		SchoolID:    "test-school-123",
		SchoolName:  "Test School",
		Permissions: []string{"read", "write"},
	}

	token, _, err := manager.GenerateTokenWithContext(userID, email, activeContext, expiresIn)
	if err != nil {
		t.Fatalf("Error generando token de prueba: %v", err)
	}
	return token
}

func TestAuthClient_ValidateToken_Local_Success(t *testing.T) {
	// Crear cliente con validación local
	client := NewAuthClient(AuthClientConfig{
		JWTSecret:    testJWTSecret,
		JWTIssuer:    testJWTIssuer,
		CacheEnabled: false,
	})

	// Generar token válido
	token := generateTestToken(t, "user-123", "test@test.com", enum.SystemRoleTeacher, 15*time.Minute)

	// Validar token
	info, err := client.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if !info.Valid {
		t.Errorf("Token debería ser válido, error: %s", info.Error)
	}
	if info.UserID != "user-123" {
		t.Errorf("UserID incorrecto: esperado 'user-123', obtenido '%s'", info.UserID)
	}
	if info.Email != "test@test.com" {
		t.Errorf("Email incorrecto: esperado 'test@test.com', obtenido '%s'", info.Email)
	}
	if info.ActiveContext == nil || info.ActiveContext.RoleName != "teacher" {
		t.Errorf("Role incorrecto: esperado 'teacher', obtenido '%v'", info.ActiveContext)
	}
}

func TestAuthClient_ValidateToken_Local_Expired(t *testing.T) {
	client := NewAuthClient(AuthClientConfig{
		JWTSecret:    testJWTSecret,
		JWTIssuer:    testJWTIssuer,
		CacheEnabled: false,
	})

	// Generar token que ya expiró (duración negativa)
	manager := auth.NewJWTManager(testJWTSecret, testJWTIssuer)
	activeContext := &auth.UserContext{
		RoleID:      "role-student",
		RoleName:    "student",
		SchoolID:    "test-school-123",
		Permissions: []string{"read"},
	}
	token, _, _ := manager.GenerateTokenWithContext("user-123", "test@test.com", activeContext, -1*time.Hour)

	info, err := client.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if info.Valid {
		t.Error("Token expirado debería ser inválido")
	}
	if info.Error == "" {
		t.Error("Debería tener mensaje de error")
	}
}

func TestAuthClient_ValidateToken_Local_InvalidSecret(t *testing.T) {
	// Cliente con un secret diferente
	client := NewAuthClient(AuthClientConfig{
		JWTSecret:    "different-secret-at-least-32-characters-long",
		JWTIssuer:    testJWTIssuer,
		CacheEnabled: false,
	})

	// Token generado con el secret original
	token := generateTestToken(t, "user-123", "test@test.com", enum.SystemRoleStudent, 15*time.Minute)

	info, err := client.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if info.Valid {
		t.Error("Token con secret diferente debería ser inválido")
	}
}

func TestAuthClient_ValidateToken_NoJWTSecret(t *testing.T) {
	// Cliente sin JWT secret configurado
	client := NewAuthClient(AuthClientConfig{
		CacheEnabled: false,
		// Sin JWTSecret ni RemoteEnabled
	})

	info, err := client.ValidateToken(context.Background(), "any-token")
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if info.Valid {
		t.Error("Sin método de validación, token debería ser inválido")
	}
	if info.Error == "" {
		t.Error("Debería indicar que no hay método de validación disponible")
	}
}

func TestAuthClient_ValidateToken_Remote_Success(t *testing.T) {
	// Crear servidor mock para validación remota
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/verify" {
			t.Errorf("Path incorrecto: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Método incorrecto: %s", r.Method)
		}

		response := TokenInfo{
			Valid:  true,
			UserID: "user-456",
			Email:  "remote@test.com",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Cliente con validación remota habilitada (sin JWT secret)
	client := NewAuthClient(AuthClientConfig{
		BaseURL:       server.URL,
		RemoteEnabled: true,
		Timeout:       5 * time.Second,
		CacheEnabled:  false,
	})

	info, err := client.ValidateToken(context.Background(), "remote-token")
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if !info.Valid {
		t.Errorf("Token debería ser válido (remoto), error: %s", info.Error)
	}
	if info.UserID != "user-456" {
		t.Errorf("UserID incorrecto: esperado 'user-456', obtenido '%s'", info.UserID)
	}
}

func TestAuthClient_ValidateToken_Fallback_ToRemote(t *testing.T) {
	var remoteCallCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&remoteCallCount, 1)
		response := TokenInfo{
			Valid:  true,
			UserID: "fallback-user",
			Email:  "fallback@test.com",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Cliente con validación local (secret incorrecto) y fallback remoto
	client := NewAuthClient(AuthClientConfig{
		JWTSecret:       "wrong-secret-at-least-32-characters-long",
		JWTIssuer:       testJWTIssuer,
		BaseURL:         server.URL,
		RemoteEnabled:   true,
		FallbackEnabled: true,
		Timeout:         5 * time.Second,
		CacheEnabled:    false,
	})

	// Token generado con el secret correcto (diferente al del cliente)
	token := generateTestToken(t, "user-123", "test@test.com", enum.SystemRoleTeacher, 15*time.Minute)

	info, err := client.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	// Debería usar fallback remoto y ser válido
	if !info.Valid {
		t.Errorf("Token debería ser válido via fallback, error: %s", info.Error)
	}
	if atomic.LoadInt32(&remoteCallCount) != 1 {
		t.Errorf("Debería haber llamado al servidor remoto como fallback")
	}
	if info.UserID != "fallback-user" {
		t.Errorf("Debería tener datos del servidor remoto")
	}
}

func TestAuthClient_Cache_Local_HitAndMiss(t *testing.T) {
	client := NewAuthClient(AuthClientConfig{
		JWTSecret:    testJWTSecret,
		JWTIssuer:    testJWTIssuer,
		CacheTTL:     5 * time.Second,
		CacheEnabled: true,
	})

	token := generateTestToken(t, "cache-user", "cache@test.com", enum.SystemRoleStudent, 15*time.Minute)

	// Primera llamada
	info1, _ := client.ValidateToken(context.Background(), token)
	if !info1.Valid {
		t.Fatalf("Primera validación debería ser exitosa")
	}

	// Segunda llamada con mismo token - debería usar caché
	info2, _ := client.ValidateToken(context.Background(), token)
	if !info2.Valid {
		t.Error("Segunda validación (cache hit) debería ser exitosa")
	}

	// Token diferente - cache miss
	token2 := generateTestToken(t, "other-user", "other@test.com", enum.SystemRoleTeacher, 15*time.Minute)
	info3, _ := client.ValidateToken(context.Background(), token2)
	if !info3.Valid {
		t.Error("Tercera validación (cache miss) debería ser exitosa")
	}
}

func TestAuthClient_Cache_Disabled(t *testing.T) {
	client := NewAuthClient(AuthClientConfig{
		JWTSecret:    testJWTSecret,
		JWTIssuer:    testJWTIssuer,
		CacheEnabled: false,
	})

	token := generateTestToken(t, "no-cache-user", "nocache@test.com", enum.SystemRoleStudent, 15*time.Minute)

	// Ambas llamadas deberían funcionar sin caché
	info1, _ := client.ValidateToken(context.Background(), token)
	info2, _ := client.ValidateToken(context.Background(), token)

	if !info1.Valid || !info2.Valid {
		t.Error("Ambas validaciones deberían ser exitosas")
	}
}

func TestAuthClient_CircuitBreaker_Remote(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:       server.URL,
		RemoteEnabled: true,
		CacheEnabled:  false,
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

	count := atomic.LoadInt32(&callCount)
	if count >= 10 {
		t.Errorf("Circuit breaker no funcionó: se esperaban menos de 10 calls, obtenido %d", count)
	}
	t.Logf("Circuit breaker funcionó: solo %d llamadas llegaron al servidor", count)
}

func TestAuthClient_Timeout_Remote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		response := TokenInfo{Valid: true}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:       server.URL,
		RemoteEnabled: true,
		Timeout:       100 * time.Millisecond,
		CacheEnabled:  false,
	})

	info, _ := client.ValidateToken(context.Background(), "test-token")

	if info.Valid {
		t.Error("Debería haber fallado por timeout")
	}
	if info.Error == "" {
		t.Error("Debería tener un mensaje de error")
	}
}

func TestTokenCache_GetSet(t *testing.T) {
	cache := newTokenCache(1 * time.Second)

	info := &TokenInfo{Valid: true, UserID: "test-user"}
	cache.Set("key1", info)

	result, found := cache.Get("key1")
	if !found {
		t.Error("Debería encontrar el entry en cache")
	}
	if result.UserID != "test-user" {
		t.Errorf("UserID incorrecto: %s", result.UserID)
	}

	_, found = cache.Get("nonexistent")
	if found {
		t.Error("No debería encontrar entry inexistente")
	}
}

func TestTokenCache_Expiration(t *testing.T) {
	cache := newTokenCache(50 * time.Millisecond)

	info := &TokenInfo{Valid: true, UserID: "test-user"}
	cache.Set("key1", info)

	_, found := cache.Get("key1")
	if !found {
		t.Error("Debería encontrar el entry recién agregado")
	}

	time.Sleep(100 * time.Millisecond)

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
		JWTSecret: testJWTSecret,
	})

	hash1 := client.hashToken("token123")
	hash2 := client.hashToken("token123")
	hash3 := client.hashToken("different-token")

	if hash1 != hash2 {
		t.Error("Mismo token debería producir mismo hash")
	}

	if hash1 == hash3 {
		t.Error("Diferente token debería producir diferente hash")
	}

	if len(hash1) != 64 {
		t.Errorf("Hash debería tener 64 caracteres, tiene %d", len(hash1))
	}
}

func TestAuthClient_ConcurrentAccess_Local(t *testing.T) {
	client := NewAuthClient(AuthClientConfig{
		JWTSecret:    testJWTSecret,
		JWTIssuer:    testJWTIssuer,
		CacheEnabled: true,
		CacheTTL:     5 * time.Second,
	})

	token := generateTestToken(t, "concurrent-user", "concurrent@test.com", enum.SystemRoleTeacher, 15*time.Minute)

	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			info, err := client.ValidateToken(context.Background(), token)
			if err != nil {
				t.Errorf("Error en validación concurrente: %v", err)
			}
			if !info.Valid {
				t.Errorf("Token debería ser válido en validación concurrente")
			}
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}
