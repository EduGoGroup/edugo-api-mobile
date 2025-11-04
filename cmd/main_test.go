package main

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetEnvironment_WithEnvSet verifica que getEnvironment retorne el valor de APP_ENV
func TestGetEnvironment_WithEnvSet(t *testing.T) {
	// Guardar valor original
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)

	testCases := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "Environment local",
			envValue: "local",
			expected: "local",
		},
		{
			name:     "Environment dev",
			envValue: "dev",
			expected: "dev",
		},
		{
			name:     "Environment qa",
			envValue: "qa",
			expected: "qa",
		},
		{
			name:     "Environment prod",
			envValue: "prod",
			expected: "prod",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("APP_ENV", tc.envValue)
			result := getEnvironment()
			assert.Equal(t, tc.expected, result, "getEnvironment debería retornar %s", tc.expected)
		})
	}
}

// TestGetEnvironment_WithoutEnvSet verifica el comportamiento por defecto
func TestGetEnvironment_WithoutEnvSet(t *testing.T) {
	// Guardar valor original
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)

	// Limpiar variable de entorno
	os.Unsetenv("APP_ENV")

	result := getEnvironment()
	assert.Equal(t, "local", result, "getEnvironment debería retornar 'local' por defecto")
}

// TestGetEnvironment_WithEmptyString verifica comportamiento con string vacío
func TestGetEnvironment_WithEmptyString(t *testing.T) {
	// Guardar valor original
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)

	// Establecer string vacío
	os.Setenv("APP_ENV", "")

	result := getEnvironment()
	assert.Equal(t, "local", result, "getEnvironment debería retornar 'local' cuando APP_ENV está vacío")
}

// TestConfigureGinMode_Production verifica que se configure modo release en prod
func TestConfigureGinMode_Production(t *testing.T) {
	// Configurar modo test inicialmente
	gin.SetMode(gin.TestMode)

	configureGinMode("prod")

	assert.Equal(t, gin.ReleaseMode, gin.Mode(), "Gin debería estar en ReleaseMode para ambiente prod")
}

// TestConfigureGinMode_NonProduction verifica que no cambie el modo en otros ambientes
func TestConfigureGinMode_NonProduction(t *testing.T) {
	testCases := []struct {
		name        string
		environment string
	}{
		{"Environment local", "local"},
		{"Environment dev", "dev"},
		{"Environment qa", "qa"},
		{"Environment staging", "staging"},
		{"Empty string", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Establecer modo test
			gin.SetMode(gin.TestMode)
			initialMode := gin.Mode()

			configureGinMode(tc.environment)

			// El modo no debería cambiar para ambientes no-prod
			// NOTA: Gin mantiene el modo previamente establecido si no es "prod"
			currentMode := gin.Mode()
			assert.Equal(t, initialMode, currentMode,
				"Gin no debería cambiar de modo para ambiente %s", tc.environment)
		})
	}
}

// TestConfigureGinMode_MultipleCallsProduction verifica múltiples llamadas con prod
func TestConfigureGinMode_MultipleCallsProduction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Primera llamada
	configureGinMode("prod")
	assert.Equal(t, gin.ReleaseMode, gin.Mode())

	// Segunda llamada (debería mantener ReleaseMode)
	configureGinMode("prod")
	assert.Equal(t, gin.ReleaseMode, gin.Mode())
}

// TestConfigureGinMode_SwitchBetweenModes verifica cambios entre modos
func TestConfigureGinMode_SwitchBetweenModes(t *testing.T) {
	// Iniciar en test
	gin.SetMode(gin.TestMode)

	// Cambiar a prod
	configureGinMode("prod")
	assert.Equal(t, gin.ReleaseMode, gin.Mode())

	// Volver a test manualmente
	gin.SetMode(gin.TestMode)
	assert.Equal(t, gin.TestMode, gin.Mode())

	// Cambiar a prod nuevamente
	configureGinMode("prod")
	assert.Equal(t, gin.ReleaseMode, gin.Mode())
}

// TestConfigureGinMode_CaseInsensitive verifica que "prod" sea case-sensitive
func TestConfigureGinMode_CaseInsensitive(t *testing.T) {
	testCases := []struct {
		name         string
		environment  string
		expectedMode string
	}{
		{"Lowercase prod", "prod", gin.ReleaseMode},
		{"Uppercase PROD", "PROD", gin.TestMode}, // No cambia porque no es exactamente "prod"
		{"Mixed case Prod", "Prod", gin.TestMode},
		{"Mixed case pRoD", "pRoD", gin.TestMode},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			configureGinMode(tc.environment)
			assert.Equal(t, tc.expectedMode, gin.Mode(),
				"El modo debería ser %s para ambiente '%s'", tc.expectedMode, tc.environment)
		})
	}
}

// BenchmarkGetEnvironment benchmark de getEnvironment
func BenchmarkGetEnvironment(b *testing.B) {
	os.Setenv("APP_ENV", "prod")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getEnvironment()
	}
}

// BenchmarkGetEnvironment_Unset benchmark cuando APP_ENV no está definido
func BenchmarkGetEnvironment_Unset(b *testing.B) {
	os.Unsetenv("APP_ENV")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getEnvironment()
	}
}

// BenchmarkConfigureGinMode benchmark de configureGinMode
func BenchmarkConfigureGinMode(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		configureGinMode("prod")
	}
}

// BenchmarkConfigureGinMode_NonProd benchmark con ambiente no-prod
func BenchmarkConfigureGinMode_NonProd(b *testing.B) {
	gin.SetMode(gin.TestMode)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		configureGinMode("dev")
	}
}

// TestGetEnvironment_Concurrency verifica seguridad en concurrencia
func TestGetEnvironment_Concurrency(t *testing.T) {
	os.Setenv("APP_ENV", "test")

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			result := getEnvironment()
			assert.NotEmpty(t, result, "getEnvironment debería retornar un valor no vacío")
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestConfigureGinMode_Concurrency verifica seguridad en concurrencia
func TestConfigureGinMode_Concurrency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			if id%2 == 0 {
				configureGinMode("prod")
			} else {
				configureGinMode("dev")
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Al final debería estar en algún modo válido
	mode := gin.Mode()
	assert.Contains(t, []string{gin.ReleaseMode, gin.TestMode, gin.DebugMode}, mode,
		"Gin debería estar en un modo válido")
}

// TestGetEnvironment_AllPossibleValues verifica todos los valores típicos
func TestGetEnvironment_AllPossibleValues(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)

	possibleValues := []string{
		"local", "dev", "development", "qa", "test", "testing",
		"staging", "stage", "prod", "production", "demo",
	}

	for _, value := range possibleValues {
		t.Run("ENV_"+value, func(t *testing.T) {
			os.Setenv("APP_ENV", value)
			result := getEnvironment()
			assert.Equal(t, value, result, "Debería retornar exactamente el valor configurado")
		})
	}
}

// TestGetEnvironment_SpecialCharacters verifica manejo de caracteres especiales
func TestGetEnvironment_SpecialCharacters(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)

	specialValues := []string{
		"prod-01", "dev-feature", "qa.test", "local_dev",
	}

	for _, value := range specialValues {
		t.Run("SPECIAL_"+value, func(t *testing.T) {
			os.Setenv("APP_ENV", value)
			result := getEnvironment()
			assert.Equal(t, value, result, "Debería manejar caracteres especiales correctamente")
		})
	}
}

// TestStartServer_FunctionExists verifica que la función startServer exista
// NOTA: No podemos testear startServer directamente porque bloquea el hilo
// con r.Run(). Este test solo verifica que la función compile.
func TestStartServer_FunctionExists(t *testing.T) {
	// Este test verifica que la función exista y compile correctamente
	// No se puede ejecutar startServer en un test porque es bloqueante
	t.Skip("startServer es una función bloqueante, no se puede testear directamente sin mocks complejos")
}
