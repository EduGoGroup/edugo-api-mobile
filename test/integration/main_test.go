//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
)

// TestMain se ejecuta antes y después de todos los tests
// Gestiona el ciclo de vida de los contenedores compartidos
func TestMain(m *testing.M) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║     🚀 INICIANDO SUITE DE TESTS DE INTEGRACIÓN                      ║")
	fmt.Println("║     Usando shared/testing containers para mejor performance         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Configurar containers usando shared/testing
	config := containers.NewConfig().
		WithPostgreSQL(nil).
		WithMongoDB(nil).
		WithRabbitMQ(nil).
		Build()

	manager, err := containers.GetManager(nil, config)
	if err != nil {
		fmt.Printf("❌ Error creando containers: %v\n", err)
		os.Exit(1)
	}

	// Ejecutar todos los tests
	exitCode := m.Run()

	// Limpiar contenedores compartidos al final
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║     🧹 LIMPIANDO CONTENEDORES COMPARTIDOS                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")

	ctx := context.Background()
	if err := manager.Cleanup(ctx); err != nil {
		fmt.Printf("⚠️  Error al limpiar contenedores: %v\n", err)
	} else {
		fmt.Println("✅ Contenedores compartidos limpiados exitosamente")
	}

	os.Exit(exitCode)
}
