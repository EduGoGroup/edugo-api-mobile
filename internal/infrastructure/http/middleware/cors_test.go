package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter crea un router de prueba con Gin
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TestCORS_HeadersSet verifica que todos los headers CORS se configuren correctamente
func TestCORS_HeadersSet(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Verificar headers CORS
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
}

// TestCORS_AllowOrigin verifica el header Access-Control-Allow-Origin
func TestCORS_AllowOrigin(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"),
		"Debería permitir todos los orígenes")
}

// TestCORS_AllowCredentials verifica el header Access-Control-Allow-Credentials
func TestCORS_AllowCredentials(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"),
		"Debería permitir credenciales")
}

// TestCORS_AllowHeaders verifica el header Access-Control-Allow-Headers
func TestCORS_AllowHeaders(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	allowedHeaders := w.Header().Get("Access-Control-Allow-Headers")

	// Verificar headers importantes
	assert.Contains(t, allowedHeaders, "Content-Type", "Debería permitir Content-Type")
	assert.Contains(t, allowedHeaders, "Authorization", "Debería permitir Authorization")
	assert.Contains(t, allowedHeaders, "Accept-Encoding", "Debería permitir Accept-Encoding")
	assert.Contains(t, allowedHeaders, "X-CSRF-Token", "Debería permitir X-CSRF-Token")
	assert.Contains(t, allowedHeaders, "Cache-Control", "Debería permitir Cache-Control")
	assert.Contains(t, allowedHeaders, "X-Requested-With", "Debería permitir X-Requested-With")
}

// TestCORS_AllowMethods verifica el header Access-Control-Allow-Methods
func TestCORS_AllowMethods(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	allowedMethods := w.Header().Get("Access-Control-Allow-Methods")

	// Verificar métodos HTTP permitidos
	assert.Contains(t, allowedMethods, "POST", "Debería permitir POST")
	assert.Contains(t, allowedMethods, "GET", "Debería permitir GET")
	assert.Contains(t, allowedMethods, "PUT", "Debería permitir PUT")
	assert.Contains(t, allowedMethods, "DELETE", "Debería permitir DELETE")
	assert.Contains(t, allowedMethods, "PATCH", "Debería permitir PATCH")
	assert.Contains(t, allowedMethods, "OPTIONS", "Debería permitir OPTIONS")
}

// TestCORS_PreflightRequest verifica el manejo de preflight requests (OPTIONS)
func TestCORS_PreflightRequest(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())
	router.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	// Verificar que la respuesta sea 204 No Content
	assert.Equal(t, http.StatusNoContent, w.Code,
		"Las preflight requests deberían retornar 204")

	// Verificar que los headers CORS estén configurados
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

// TestCORS_PreflightWithDifferentPaths verifica preflight en diferentes rutas
func TestCORS_PreflightWithDifferentPaths(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected int
	}{
		{
			name:     "Preflight en /api/users",
			path:     "/api/users",
			expected: http.StatusNoContent,
		},
		{
			name:     "Preflight en /api/materials",
			path:     "/api/materials",
			expected: http.StatusNoContent,
		},
		{
			name:     "Preflight en /health",
			path:     "/health",
			expected: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(CORS())
			router.GET(tc.path, func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("OPTIONS", tc.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

// TestCORS_ActualRequest verifica que las peticiones normales (no preflight) funcionen correctamente
func TestCORS_ActualRequest(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	// Verificar que la petición normal pase correctamente
	assert.Equal(t, http.StatusOK, w.Code, "Las peticiones normales deberían pasar")

	// Verificar que los headers CORS estén configurados
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORS_DifferentHTTPMethods verifica CORS en diferentes métodos HTTP
func TestCORS_DifferentHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(CORS())

			// Registrar handler para el método
			router.Handle(method, "/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, "/test", nil)
			router.ServeHTTP(w, req)

			// Verificar headers CORS
			assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestCORS_WithCustomOriginHeader verifica que CORS funcione con header Origin personalizado
func TestCORS_WithCustomOriginHeader(t *testing.T) {
	testCases := []struct {
		name         string
		originHeader string
	}{
		{"Origin localhost", "http://localhost:3000"},
		{"Origin production", "https://edugo.com"},
		{"Origin dev", "https://dev.edugo.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(CORS())
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", tc.originHeader)
			router.ServeHTTP(w, req)

			// Verificar que CORS responda correctamente
			assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestCORS_WithAuthorizationHeader verifica que CORS permita el header Authorization
func TestCORS_WithAuthorizationHeader(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

// TestCORS_MultipleMiddleware verifica que CORS funcione con otros middleware
func TestCORS_MultipleMiddleware(t *testing.T) {
	router := setupTestRouter()

	// CORS debe ser el primer middleware
	router.Use(CORS())
	router.Use(gin.Recovery())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORS_ResponseNotModified verifica que CORS no modifique el body de la respuesta
func TestCORS_ResponseNotModified(t *testing.T) {
	router := setupTestRouter()
	router.Use(CORS())

	expectedBody := `{"message":"test","value":123}`
	router.GET("/test", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(expectedBody))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedBody, w.Body.String(), "CORS no debería modificar el body")
}

// BenchmarkCORS_Middleware benchmark del middleware CORS
func BenchmarkCORS_Middleware(b *testing.B) {
	router := setupTestRouter()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// BenchmarkCORS_PreflightRequest benchmark de preflight requests
func BenchmarkCORS_PreflightRequest(b *testing.B) {
	router := setupTestRouter()
	router.Use(CORS())
	router.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	req, _ := http.NewRequest("OPTIONS", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
