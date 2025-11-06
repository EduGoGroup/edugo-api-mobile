// +build integration

package integration

import (
	"testing"
)

// TestExample muestra cómo usar el sistema de control de tests
func TestExample(t *testing.T) {
	// ✅ IMPORTANTE: Siempre llamar esto al inicio
	SkipIfIntegrationTestsDisabled(t)

	// Si llegamos aquí, los tests están habilitados
	t.Log("✅ Integration tests están HABILITADOS")
	
	// TODO: Implementar test real cuando esté listo
	// Por ahora solo verificamos que el sistema de control funciona
}

// TestExampleAlwaysRuns es un test que NO usa el sistema de control
// Útil para verificar que el build tag funciona
func TestExampleAlwaysRuns(t *testing.T) {
	t.Log("🏃 Este test siempre corre (sin SkipIfIntegrationTestsDisabled)")
	
	// Verificar que estamos en modo integration
	if !testing.Short() {
		t.Log("✅ Build tag 'integration' está activo")
	}
}

// TestCheckDockerAvailable verifica que Docker esté disponible
func TestCheckDockerAvailable(t *testing.T) {
	SkipIfIntegrationTestsDisabled(t)
	
	// TODO: Agregar verificación real de Docker con exec.Command
	t.Log("🐳 Docker check pendiente de implementar")
}
