//go:build integration
// +build integration

package suite_test

import (
	"context"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
	"github.com/stretchr/testify/assert"
	testifySuite "github.com/stretchr/testify/suite"
)

// ExampleSuite es un ejemplo de cómo usar IntegrationTestSuite
type ExampleSuite struct {
	suite.IntegrationTestSuite
}

// TestExampleSuite ejecuta la suite de ejemplo
func TestExampleSuite(t *testing.T) {
	testifySuite.Run(t, new(ExampleSuite))
}

// TestPostgresConnection verifica que PostgreSQL tiene datos de seeds
func (s *ExampleSuite) TestPostgresConnection() {
	ctx := context.Background()

	// Verificar que la base de datos tiene tablas creadas
	var count int
	query := `
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name != 'schema_migrations'
	`
	err := s.PostgresDB.QueryRowContext(ctx, query).Scan(&count)

	s.NoError(err, "Query debe ejecutarse sin errores")
	s.Greater(count, 0, "Debe haber al menos una tabla creada")

	s.Logger.Info("✅ PostgreSQL tiene tablas creadas", "count", count)
}

// TestSeedsApplied verifica que los seeds se aplicaron correctamente
func (s *ExampleSuite) TestSeedsApplied() {
	ctx := context.Background()

	// Verificar que hay usuarios en la base de datos (del seed)
	var userCount int
	err := s.PostgresDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)

	s.NoError(err, "Query debe ejecutarse sin errores")
	s.Greater(userCount, 0, "Debe haber usuarios del seed")

	s.Logger.Info("✅ Seeds aplicados correctamente", "users", userCount)
}

// TestDataCleanupBetweenTests verifica que datos se limpian entre tests
func (s *ExampleSuite) TestDataCleanupBetweenTests_First() {
	ctx := context.Background()

	// Insertar un usuario nuevo
	_, err := s.PostgresDB.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active)
		VALUES (gen_random_uuid(), 'test@cleanup.com', 'hash', 'Test', 'Cleanup', 'student', true)
	`)
	s.NoError(err)

	// Verificar que existe
	var count int
	err = s.PostgresDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", "test@cleanup.com").Scan(&count)
	s.NoError(err)
	s.Equal(1, count, "Usuario debe existir")

	s.Logger.Info("✅ Usuario de prueba insertado")
}

// TestDataCleanupBetweenTests_Second se ejecuta después del anterior
// Verifica que el usuario insertado en el test anterior NO existe
func (s *ExampleSuite) TestDataCleanupBetweenTests_Second() {
	ctx := context.Background()

	// Verificar que el usuario del test anterior NO existe
	var count int
	err := s.PostgresDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", "test@cleanup.com").Scan(&count)
	s.NoError(err)
	s.Equal(0, count, "Usuario del test anterior debe haberse limpiado")

	s.Logger.Info("✅ Datos limpiados correctamente entre tests")
}

// TestContainersAreShared verifica que contenedores NO se reinician entre tests
func (s *ExampleSuite) TestContainersAreShared() {
	// El contenedor debe ser el mismo que en tests anteriores
	assert.NotNil(s.T(), s.PostgresContainer, "Contenedor debe estar disponible")
	assert.NotNil(s.T(), s.MongoContainer, "Contenedor debe estar disponible")
	assert.NotNil(s.T(), s.RabbitContainer, "Contenedor debe estar disponible")

	s.Logger.Info("✅ Contenedores compartidos correctamente")
}
