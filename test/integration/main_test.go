//go:build integration

package integration

import (
	"fmt"
	"os"
	"testing"
)

// TestMain se ejecuta antes y después de todos los tests
// Gestiona el ciclo de vida de los contenedores compartidos
func TestMain(m *testing.M) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║     🚀 INICIANDO SUITE DE TESTS DE INTEGRACIÓN                      ║")
	fmt.Println("║     Usando contenedores compartidos para mejor performance           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Ejecutar todos los tests
	exitCode := m.Run()

	// Limpiar contenedores compartidos al final
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║     🧹 LIMPIANDO CONTENEDORES COMPARTIDOS                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")

	if err := TerminateSharedContainers(); err != nil {
		fmt.Printf("⚠️  Error al limpiar contenedores: %v\n", err)
	} else {
		fmt.Println("✅ Contenedores compartidos limpiados exitosamente")
	}

	os.Exit(exitCode)
}
